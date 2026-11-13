package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// pasosMetricas guarda el estado de los escenarios de HU-03. Siembra el historial
// con SQL, porque un Sprint cerrado no se puede armar en otro orden con las reglas
// de la API, y ejercita la consulta y las ediciones con los services reales.
type pasosMetricas struct {
	db          *sql.DB
	calcular    *service.CalcularMetricas
	editar      *service.EditarProyecto
	proyectos   *repository.SQLiteProjectRepository
	historiasBD *repository.SQLiteBacklogRepository
	editarHist  *service.EditarHistoria
	registrar   *service.RegistrarEstimacion

	proyectoID   int64
	scrumMaster  int64
	auxiliar     int64
	sprints      map[string]int64
	historias    map[string]int64
	metricas     domain.MetricasProyecto
	errConsulta  error
	errEdicion   error
	errCargaHist error
}

func (p *pasosMetricas) preparar() error {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return err
	}
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		return err
	}
	proyectos := repository.NewSQLiteProjectRepository(db)
	historias := repository.NewSQLiteBacklogRepository(db)
	sprints := repository.NewSQLiteSprintRepository(db)

	*p = pasosMetricas{
		db:          db,
		proyectos:   proyectos,
		historiasBD: historias,
		calcular:    service.NewCalcularMetricas(proyectos, repository.NewSQLiteMetricasRepository(db)),
		editar:      service.NewEditarProyecto(proyectos),
		editarHist:  service.NewEditarHistoria(historias, sprints, repository.NewSQLiteProjectRepository(db)),
		registrar:   service.NewRegistrarEstimacion(historias, proyectos, repository.NewSQLiteVotacionRepository(db), sprints),
		sprints:     map[string]int64{},
		historias:   map[string]int64{},
	}
	return nil
}

func (p *pasosMetricas) insertar(consulta string, args ...any) (int64, error) {
	resultado, err := p.db.Exec(consulta, args...)
	if err != nil {
		return 0, fmt.Errorf("sembrar %q: %w", consulta, err)
	}
	return resultado.LastInsertId()
}

func (p *pasosMetricas) crearProyecto(horas any) error {
	id, err := p.insertar("INSERT INTO proyectos (nombre, descripcion, horas_por_story_point, creado_en) VALUES ('Proyecto', '', ?, '2026-10-01T00:00:00Z')", horas)
	if err != nil {
		return err
	}
	p.proyectoID = id
	p.scrumMaster, err = p.insertar("INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES ('Carlos', 'carlos', '2026-10-01T00:00:00Z')")
	if err != nil {
		return err
	}
	_, err = p.insertar("INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, 'scrum_master', '2026-10-01T00:00:00Z')", p.proyectoID, p.scrumMaster)
	return err
}

func (p *pasosMetricas) proyectoConHoras(horas float64) error { return p.crearProyecto(horas) }
func (p *pasosMetricas) proyectoSinHoras() error              { return p.crearProyecto(nil) }

// sprint crea el Sprint la primera vez que se lo nombra. Nace Finalizado: dos
// Sprints no pueden estar Activos a la vez mientras se siembran sus datos.
func (p *pasosMetricas) sprint(nombre string) (int64, error) {
	if id, existe := p.sprints[nombre]; existe {
		return id, nil
	}
	id, err := p.insertar(
		"INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (?, 'Objetivo', '2026-10-05', '2026-10-19', 'Finalizado')",
		p.proyectoID)
	if err != nil {
		return 0, err
	}
	p.sprints[nombre] = id
	return id, nil
}

func (p *pasosMetricas) sprintFinalizado(nombre string) error {
	id, err := p.sprint(nombre)
	if err != nil {
		return err
	}
	var estado string
	if err := p.db.QueryRow("SELECT estado FROM sprints WHERE id = ?", id).Scan(&estado); err != nil {
		return err
	}
	if estado != string(domain.SprintFinalizado) {
		return fmt.Errorf("el Sprint %s debía estar finalizado, está %s", nombre, estado)
	}
	return nil
}

func (p *pasosMetricas) sprintActivo(nombre string) error {
	id, err := p.sprint(nombre)
	if err != nil {
		return err
	}
	_, err = p.db.Exec("UPDATE sprints SET estado = 'Activo' WHERE id = ?", id)
	return err
}

func (p *pasosMetricas) historiaCompletada(sprintID int64, puntos any) (int64, error) {
	return p.insertar(
		`INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp, sprint_id, fecha_finalizacion)
		 VALUES (?, 'Historia', 'Detalle', 'M', 'Completada', ?, ?, '2026-10-10')`,
		p.proyectoID, puntos, sprintID)
}

func (p *pasosMetricas) historiaEnBacklog(puntos any) (int64, error) {
	return p.insertar(
		"INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp) VALUES (?, 'Historia', 'Detalle', 'M', 'Nueva', ?)",
		p.proyectoID, puntos)
}

func (p *pasosMetricas) sprintConCompletada(nombre string, puntos int) error {
	id, err := p.sprint(nombre)
	if err != nil {
		return err
	}
	_, err = p.historiaCompletada(id, puntos)
	return err
}

func (p *pasosMetricas) arrastrar(nombreSprint string, puntos any) (int64, error) {
	sprintID, err := p.sprint(nombreSprint)
	if err != nil {
		return 0, err
	}
	historia, err := p.historiaEnBacklog(puntos)
	if err != nil {
		return 0, err
	}
	_, err = p.insertar("INSERT INTO historias_arrastradas (sprint_id, historia_id, proyecto_id, estimacion_sp) VALUES (?, ?, ?, ?)",
		sprintID, historia, p.proyectoID, puntos)
	return historia, err
}

func (p *pasosMetricas) sprintConArrastrada(nombre string, puntos int) error {
	_, err := p.arrastrar(nombre, puntos)
	return err
}

// historiaAuxiliar sostiene las horas y los defectos sin sumar historias
// comprometidas a ningún Sprint.
func (p *pasosMetricas) historiaAuxiliar() (int64, error) {
	if p.auxiliar != 0 {
		return p.auxiliar, nil
	}
	id, err := p.historiaEnBacklog(nil)
	p.auxiliar = id
	return id, err
}

func (p *pasosMetricas) sprintConHoras(nombre string, horas int) error {
	sprintID, err := p.sprint(nombre)
	if err != nil {
		return err
	}
	historia, err := p.historiaAuxiliar()
	if err != nil {
		return err
	}
	_, err = p.insertar(
		`INSERT INTO registros_horas (proyecto_id, historia_id, usuario_id, sprint_id, fecha, actividad, minutos, creado_en)
		 VALUES (?, ?, ?, ?, '2026-10-10', 'Trabajo', ?, '2026-10-10T00:00:00Z')`,
		p.proyectoID, historia, p.scrumMaster, sprintID, horas*60)
	return err
}

func (p *pasosMetricas) insertarDefecto(deteccion int64, resolucion any) error {
	historia, err := p.historiaAuxiliar()
	if err != nil {
		return err
	}
	estado := "Pendiente"
	if resolucion != nil {
		estado = "Resuelto"
	}
	_, err = p.insertar(
		`INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id, creado_en)
		 VALUES (?, ?, 'Falla', 'Alta', ?, ?, ?, '2026-10-10T00:00:00Z')`,
		p.proyectoID, historia, estado, deteccion, resolucion)
	return err
}

func (p *pasosMetricas) sprintConDetectados(nombre string, cantidad int) error {
	sprintID, err := p.sprint(nombre)
	if err != nil {
		return err
	}
	for range cantidad {
		if err := p.insertarDefecto(sprintID, nil); err != nil {
			return err
		}
	}
	return nil
}

// sprintConResueltos resuelve en el Sprint defectos ya detectados en él.
func (p *pasosMetricas) sprintConResueltos(nombre string, cantidad int) error {
	sprintID, err := p.sprint(nombre)
	if err != nil {
		return err
	}
	resultado, err := p.db.Exec(
		`UPDATE defectos SET estado = 'Resuelto', sprint_resolucion_id = ?
		 WHERE id IN (SELECT id FROM defectos WHERE sprint_deteccion_id = ? AND estado = 'Pendiente' LIMIT ?)`,
		sprintID, sprintID, cantidad)
	if err != nil {
		return err
	}
	if afectadas, _ := resultado.RowsAffected(); int(afectadas) != cantidad {
		return fmt.Errorf("se resolvieron %d defectos, se esperaban %d", afectadas, cantidad)
	}
	return nil
}

func (p *pasosMetricas) defectoEntreSprints(deteccion, resolucion string) error {
	detectado, err := p.sprint(deteccion)
	if err != nil {
		return err
	}
	resuelto, err := p.sprint(resolucion)
	if err != nil {
		return err
	}
	return p.insertarDefecto(detectado, resuelto)
}

func (p *pasosMetricas) historiaArrastradaConPuntos(nombre string, puntos int, sprint string) error {
	id, err := p.arrastrar(sprint, puntos)
	p.historias[nombre] = id
	return err
}

func (p *pasosMetricas) historiaArrastradaSinEstimar(nombre, sprint string) error {
	id, err := p.arrastrar(sprint, nil)
	p.historias[nombre] = id
	return err
}

func (p *pasosMetricas) historiaSinEstimarCompletada(nombre, sprint string) error {
	sprintID, err := p.sprint(sprint)
	if err != nil {
		return err
	}
	p.historias[nombre], err = p.historiaCompletada(sprintID, nil)
	return err
}

func (p *pasosMetricas) historiaSinEstimarEnSprintActivo(nombre, sprint string) error {
	if err := p.sprintActivo(sprint); err != nil {
		return err
	}
	id, err := p.historiaEnBacklog(nil)
	if err != nil {
		return err
	}
	p.historias[nombre] = id
	_, err = p.db.Exec("UPDATE historias SET sprint_id = ? WHERE id = ?", p.sprints[sprint], id)
	return err
}

func (p *pasosMetricas) historiaSeCompleta(nombre, sprint string) error {
	sprintID, err := p.sprint(sprint)
	if err != nil {
		return err
	}
	_, err = p.db.Exec("UPDATE historias SET sprint_id = ?, estado = 'Completada', fecha_finalizacion = '2026-10-20' WHERE id = ?", sprintID, p.historias[nombre])
	return err
}

func (p *pasosMetricas) rondaRevelada(nombre string) error {
	_, err := p.insertar("INSERT INTO rondas_votacion (historia_id, numero, revelada_en) VALUES (?, 1, '2026-10-27T15:04:05Z')", p.historias[nombre])
	return err
}

func (p *pasosMetricas) editarHistoria(nombre, titulo string, puntos *int) error {
	actual, err := p.historiasBD.ObtenerPorID(context.Background(), p.historias[nombre])
	if err != nil {
		return err
	}
	actor, err := integranteDeLaHistoria(p.db, p.historias[nombre])
	if err != nil {
		return err
	}
	_, err = p.editarHist.Ejecutar(context.Background(), service.EditarHistoriaInput{
		ID: p.historias[nombre], Titulo: titulo, Descripcion: "Detalle", EstimacionSP: puntos,
		EstimacionLeida: actual.EstimacionSP, UsuarioID: actor,
	})
	return err
}

func (p *pasosMetricas) historiaSeEdita(nombre string, puntos int) error {
	return p.editarHistoria(nombre, "Historia", &puntos)
}

func (p *pasosMetricas) cargarPuntos(nombre string, puntos int, via string) error {
	if via == "votación" {
		_, err := p.registrar.Ejecutar(context.Background(), p.historias[nombre], p.scrumMaster, puntos)
		return err
	}
	return p.editarHistoria(nombre, "Historia", &puntos)
}

func (p *pasosMetricas) intentarCargarPorEdicion(puntos int, nombre string) error {
	p.errCargaHist = p.cargarPuntos(nombre, puntos, "edición")
	return nil
}

func (p *pasosMetricas) intentarRegistrar(puntos int, nombre string) error {
	p.errCargaHist = p.cargarPuntos(nombre, puntos, "votación")
	return nil
}

func (p *pasosMetricas) cambiarTitulo(nombre, titulo string) error {
	return p.editarHistoria(nombre, titulo, nil)
}

func (p *pasosMetricas) consultar() error {
	p.metricas, p.errConsulta = p.calcular.Ejecutar(context.Background(), p.proyectoID)
	return p.errConsulta
}

func (p *pasosMetricas) consultarInexistente(id int) error {
	_, p.errConsulta = p.calcular.Ejecutar(context.Background(), int64(id))
	return nil
}

func (p *pasosMetricas) editarProyecto(enviadas bool, valor *string) error {
	_, err := p.editar.Ejecutar(context.Background(), service.EditarProyectoInput{
		ID: p.proyectoID, Nombre: "Proyecto", HorasPorSPEnviadas: enviadas, HorasPorSP: valor,
		UsuarioID: scrumMasterDelProyecto(p.db, p.proyectoID),
	})
	return err
}

func (p *pasosMetricas) definirHoras(valor string) error { return p.editarProyecto(true, &valor) }
func (p *pasosMetricas) dejarHorasVacias() error         { return p.editarProyecto(true, nil) }
func (p *pasosMetricas) editarSinMencionarHoras() error  { return p.editarProyecto(false, nil) }

func (p *pasosMetricas) intentarEditarHoras(valor string) error {
	p.errEdicion = p.editarProyecto(true, &valor)
	return nil
}

func campoDelError(err error) string {
	var verr domain.ValidationError
	if errors.As(err, &verr) {
		return verr.Campo
	}
	return ""
}

func (p *pasosMetricas) edicionRechazada(campo string) error {
	if obtenido := campoDelError(p.errEdicion); obtenido != campo {
		return fmt.Errorf("se esperaba un error sobre %q, se obtuvo %v", campo, p.errEdicion)
	}
	return nil
}

func (p *pasosMetricas) cargaRechazada(campo string) error {
	if obtenido := campoDelError(p.errCargaHist); obtenido != campo {
		return fmt.Errorf("se esperaba un error sobre %q, se obtuvo %v", campo, p.errCargaHist)
	}
	return nil
}

func (p *pasosMetricas) consultaRechazada(mensaje string) error {
	if p.errConsulta == nil || p.errConsulta.Error() != mensaje {
		return fmt.Errorf("se esperaba el error %q, se obtuvo %v", mensaje, p.errConsulta)
	}
	return nil
}

func (p *pasosMetricas) tieneHoras(horas float64) error {
	proyecto, err := p.proyectos.ObtenerPorID(context.Background(), p.proyectoID)
	if err != nil {
		return err
	}
	if proyecto.HorasPorStoryPoint == nil || *proyecto.HorasPorStoryPoint != horas {
		return fmt.Errorf("se esperaban %v horas por Story Point, se obtuvo %v", horas, proyecto.HorasPorStoryPoint)
	}
	return nil
}

func (p *pasosMetricas) noTieneHoras() error {
	proyecto, err := p.proyectos.ObtenerPorID(context.Background(), p.proyectoID)
	if err != nil {
		return err
	}
	if proyecto.HorasPorStoryPoint != nil {
		return fmt.Errorf("se esperaba sin definir, se obtuvo %v", *proyecto.HorasPorStoryPoint)
	}
	return nil
}

func (p *pasosMetricas) historiaGuardada(nombre string) (domain.BacklogItem, error) {
	return p.historiasBD.ObtenerPorID(context.Background(), p.historias[nombre])
}

func (p *pasosMetricas) historiaSigueSinEstimar(nombre string) error {
	historia, err := p.historiaGuardada(nombre)
	if err != nil {
		return err
	}
	if historia.EstimacionSP != nil {
		return fmt.Errorf("se esperaba sin estimar, tiene %d Story Points", *historia.EstimacionSP)
	}
	return nil
}

func (p *pasosMetricas) historiaTienePuntos(nombre string, puntos int) error {
	historia, err := p.historiaGuardada(nombre)
	if err != nil {
		return err
	}
	if historia.EstimacionSP == nil || *historia.EstimacionSP != puntos {
		return fmt.Errorf("se esperaban %d Story Points, se obtuvo %v", puntos, historia.EstimacionSP)
	}
	return nil
}

func (p *pasosMetricas) historiaSeLlama(nombre, titulo string) error {
	historia, err := p.historiaGuardada(nombre)
	if err != nil {
		return err
	}
	if historia.Titulo != titulo {
		return fmt.Errorf("se esperaba el título %q, se obtuvo %q", titulo, historia.Titulo)
	}
	return nil
}

func entero(valor int) *float64 {
	convertido := float64(valor)
	return &convertido
}

func valoresDe(m domain.Metricas) map[string]*float64 {
	return map[string]*float64{
		"Story Points planificados": entero(m.SPPlanificados),
		"Story Points completados":  entero(m.SPCompletados),
		"velocidad":                 &m.Velocidad,
		"historias completadas (%)": &m.PorcentajeCompletadas,
		"horas estimadas":           m.HorasEstimadas,
		"horas reales":              &m.HorasReales,
		"desviación (h)":            m.DesviacionHoras,
		"desviación (%)":            m.DesviacionPorcentaje,
		"defectos detectados":       entero(m.DefectosDetectados),
		"defectos resueltos":        entero(m.DefectosResueltos),
	}
}

func compararMetricas(m domain.Metricas, tabla *godog.Table) error {
	valores := valoresDe(m)
	for _, fila := range tabla.Rows {
		nombre, esperado := fila.Cells[0].Value, fila.Cells[1].Value
		obtenido, existe := valores[nombre]
		if !existe {
			return fmt.Errorf("métrica desconocida: %q", nombre)
		}
		if esperado == "sin datos" {
			if obtenido != nil {
				return fmt.Errorf("%s: se esperaba sin datos, se obtuvo %v", nombre, *obtenido)
			}
			continue
		}
		numero, err := strconv.ParseFloat(esperado, 64)
		if err != nil {
			return fmt.Errorf("%s: valor esperado inválido %q", nombre, esperado)
		}
		if obtenido == nil || *obtenido != numero {
			return fmt.Errorf("%s: se esperaba %v, se obtuvo %v", nombre, numero, describirValor(obtenido))
		}
	}
	return nil
}

func describirValor(valor *float64) string {
	if valor == nil {
		return "sin datos"
	}
	return strconv.FormatFloat(*valor, 'f', -1, 64)
}

func (p *pasosMetricas) proyectoInforma(tabla *godog.Table) error {
	return compararMetricas(p.metricas.Proyecto, tabla)
}

func (p *pasosMetricas) sprintInforma(nombre string, tabla *godog.Table) error {
	for _, s := range p.metricas.Sprints {
		if s.SprintID == p.sprints[nombre] {
			return compararMetricas(s.Metricas, tabla)
		}
	}
	return fmt.Errorf("las métricas no incluyen el Sprint %s: %+v", nombre, p.metricas.Sprints)
}

func (p *pasosMetricas) soloIncluye(nombre string) error {
	if len(p.metricas.Sprints) != 1 || p.metricas.Sprints[0].SprintID != p.sprints[nombre] {
		return fmt.Errorf("se esperaba solo el Sprint %s, se obtuvo %+v", nombre, p.metricas.Sprints)
	}
	return nil
}

func (p *pasosMetricas) noIncluyeSprints() error {
	if p.metricas.Sprints == nil || len(p.metricas.Sprints) != 0 {
		return fmt.Errorf("se esperaba una lista vacía de Sprints, se obtuvo %#v", p.metricas.Sprints)
	}
	return nil
}

// InitializeScenarioMetricas registra los pasos de HU-03 (@HU-03).
func InitializeScenarioMetricas(ctx *godog.ScenarioContext) {
	p := &pasosMetricas{}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return goctx, p.preparar()
	})
	ctx.After(func(goctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if p.db != nil {
			return goctx, p.db.Close()
		}
		return goctx, nil
	})

	ctx.Step(`^un proyecto de métricas con ([\d.]+) horas por Story Point$`, p.proyectoConHoras)
	ctx.Step(`^un proyecto de métricas sin horas por Story Point$`, p.proyectoSinHoras)
	ctx.Step(`^el Sprint de métricas "([^"]*)" tiene una historia completada de (\d+) Story Points$`, p.sprintConCompletada)
	ctx.Step(`^el Sprint de métricas "([^"]*)" tiene una historia arrastrada de (\d+) Story Points$`, p.sprintConArrastrada)
	ctx.Step(`^el Sprint de métricas "([^"]*)" tiene (\d+) horas imputadas$`, p.sprintConHoras)
	ctx.Step(`^el Sprint de métricas "([^"]*)" tiene (\d+) defectos? detectados?$`, p.sprintConDetectados)
	ctx.Step(`^el Sprint de métricas "([^"]*)" tiene (\d+) defectos? resueltos?$`, p.sprintConResueltos)
	ctx.Step(`^el Sprint de métricas "([^"]*)" está finalizado$`, p.sprintFinalizado)
	ctx.Step(`^el Sprint de métricas "([^"]*)" está activo$`, p.sprintActivo)
	ctx.Step(`^un defecto de métricas detectado en el Sprint "([^"]*)" y resuelto en el Sprint "([^"]*)"$`, p.defectoEntreSprints)
	ctx.Step(`^la historia de métricas "([^"]*)" de (\d+) Story Points fue arrastrada en el Sprint "([^"]*)"$`, p.historiaArrastradaConPuntos)
	ctx.Step(`^la historia de métricas "([^"]*)" sin estimar fue arrastrada en el Sprint "([^"]*)"$`, p.historiaArrastradaSinEstimar)
	ctx.Step(`^la historia de métricas "([^"]*)" sin estimar se completó en el Sprint finalizado "([^"]*)"$`, p.historiaSinEstimarCompletada)
	ctx.Step(`^la historia de métricas "([^"]*)" sin estimar está en el Sprint activo "([^"]*)"$`, p.historiaSinEstimarEnSprintActivo)
	ctx.Step(`^la historia de métricas "([^"]*)" se edita a (\d+) Story Points en el Product Backlog$`, p.historiaSeEdita)
	ctx.Step(`^la historia de métricas "([^"]*)" se completa en el Sprint "([^"]*)"$`, p.historiaSeCompleta)
	ctx.Step(`^la historia de métricas "([^"]*)" tiene su ronda de votación revelada$`, p.rondaRevelada)

	ctx.Step(`^consulto las métricas del proyecto$`, p.consultar)
	ctx.Step(`^consulto las métricas del proyecto con identificador (\d+), que no existe$`, p.consultarInexistente)
	ctx.Step(`^edito el proyecto de métricas y defino ([\d.]+) horas por Story Point$`, p.definirHoras)
	ctx.Step(`^edito el proyecto de métricas y dejo vacías las horas por Story Point$`, p.dejarHorasVacias)
	ctx.Step(`^edito el proyecto de métricas sin mencionar las horas por Story Point$`, p.editarSinMencionarHoras)
	ctx.Step(`^intento editar el proyecto de métricas con las horas por Story Point (.+)$`, p.intentarEditarHoras)
	ctx.Step(`^intento cargar (\d+) Story Points a la historia de métricas "([^"]*)" por edición$`, p.intentarCargarPorEdicion)
	ctx.Step(`^el Scrum Master intenta registrar (\d+) como estimación acordada de la historia de métricas "([^"]*)"$`, p.intentarRegistrar)
	ctx.Step(`^cambio el título de la historia de métricas "([^"]*)" a "([^"]*)" sin cargar Story Points$`, p.cambiarTitulo)
	ctx.Step(`^cargo (\d+) Story Points a la historia de métricas "([^"]*)" por (edición|votación)$`, func(puntos int, nombre, via string) error {
		return p.cargarPuntos(nombre, puntos, via)
	})

	ctx.Step(`^el Sprint "([^"]*)" informa estas métricas:$`, p.sprintInforma)
	ctx.Step(`^el proyecto informa estas métricas:$`, p.proyectoInforma)
	ctx.Step(`^las métricas incluyen solo el Sprint "([^"]*)"$`, p.soloIncluye)
	ctx.Step(`^las métricas no incluyen Sprints$`, p.noIncluyeSprints)
	ctx.Step(`^la consulta de métricas se rechaza con error "([^"]*)"$`, p.consultaRechazada)
	ctx.Step(`^la edición del proyecto de métricas se rechaza con un error sobre el campo "([^"]*)"$`, p.edicionRechazada)
	ctx.Step(`^la carga se rechaza con una advertencia sobre el campo "([^"]*)"$`, p.cargaRechazada)
	ctx.Step(`^el proyecto de métricas tiene ([\d.]+) horas por Story Point$`, p.tieneHoras)
	ctx.Step(`^el proyecto de métricas no tiene horas por Story Point$`, p.noTieneHoras)
	ctx.Step(`^la historia de métricas "([^"]*)" sigue sin estimar$`, p.historiaSigueSinEstimar)
	ctx.Step(`^la historia de métricas "([^"]*)" tiene (\d+) Story Points$`, p.historiaTienePuntos)
	ctx.Step(`^la historia de métricas "([^"]*)" se llama "([^"]*)"$`, p.historiaSeLlama)
}
