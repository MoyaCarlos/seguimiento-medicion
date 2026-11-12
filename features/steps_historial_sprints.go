package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// pasosHistorial guarda el estado de los escenarios de HU-12. Igual que HU-03,
// siembra los Sprints cerrados con SQL y ejercita las consultas, la edición y
// el cálculo de métricas con los services reales.
type pasosHistorial struct {
	db         *sql.DB
	listar     *service.ConsultarHistorialSprints
	detallar   *service.ConsultarDetalleSprint
	calcular   *service.CalcularMetricas
	editarHist *service.EditarHistoria

	proyectoID  int64
	usuarioID   int64
	auxiliar    int64
	sprintAjeno int64
	sprints     map[string]int64
	historias   map[string]int64
	fotoAntes   string
	historial   []domain.SprintHistorial
	detalle     domain.DetalleSprint
	errConsulta error
}

func (p *pasosHistorial) preparar() error {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return err
	}
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		return err
	}
	proyectos := repository.NewSQLiteProjectRepository(db)
	sprints := repository.NewSQLiteSprintRepository(db)
	calcular := service.NewCalcularMetricas(proyectos, repository.NewSQLiteMetricasRepository(db))

	*p = pasosHistorial{
		db:         db,
		calcular:   calcular,
		listar:     service.NewConsultarHistorialSprints(calcular, sprints),
		detallar:   service.NewConsultarDetalleSprint(proyectos, sprints, calcular, repository.NewSQLiteHistorialSprintRepository(db)),
		editarHist: service.NewEditarHistoria(repository.NewSQLiteBacklogRepository(db), sprints, repository.NewSQLiteProjectRepository(db)),
		sprints:    map[string]int64{},
		historias:  map[string]int64{},
	}
	return nil
}

func (p *pasosHistorial) insertar(consulta string, args ...any) (int64, error) {
	resultado, err := p.db.Exec(consulta, args...)
	if err != nil {
		return 0, fmt.Errorf("sembrar %q: %w", consulta, err)
	}
	return resultado.LastInsertId()
}

func (p *pasosHistorial) crearProyecto(horas any) (int64, error) {
	return p.insertar("INSERT INTO proyectos (nombre, descripcion, horas_por_story_point, creado_en) VALUES ('Proyecto', '', ?, '2026-10-01T00:00:00Z')", horas)
}

func (p *pasosHistorial) proyectoConHoras(horas float64) error {
	id, err := p.crearProyecto(horas)
	if err != nil {
		return err
	}
	p.proyectoID = id
	p.usuarioID, err = p.insertar("INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES ('Carlos', 'carlos', '2026-10-01T00:00:00Z')")
	if err != nil {
		return err
	}
	_, err = p.insertar("INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, 'scrum_master', '2026-10-01T00:00:00Z')", p.proyectoID, p.usuarioID)
	return err
}

func (p *pasosHistorial) proyectoSinHoras() error {
	_, err := p.db.Exec("UPDATE proyectos SET horas_por_story_point = NULL WHERE id = ?", p.proyectoID)
	return err
}

// sprint crea el Sprint la primera vez que se lo nombra. Nace Finalizado, como
// en HU-03: dos Sprints no pueden estar Activos a la vez mientras se siembran.
func (p *pasosHistorial) sprint(nombre string) (int64, error) {
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

func (p *pasosHistorial) sprintFinalizadoConGoal(nombre, goal string) error {
	id, err := p.sprint(nombre)
	if err != nil {
		return err
	}
	_, err = p.db.Exec("UPDATE sprints SET sprint_goal = ?, fecha_inicio = '2026-10-05', fecha_fin = '2026-10-19', estado = 'Finalizado' WHERE id = ?", goal, id)
	return err
}

func (p *pasosHistorial) sprintEnEstado(nombre, estado string) error {
	id, err := p.sprint(nombre)
	if err != nil {
		return err
	}
	if estado == "pendiente" {
		_, err = p.db.Exec("UPDATE sprints SET sprint_goal = NULL, fecha_inicio = NULL, fecha_fin = NULL, estado = 'Pendiente' WHERE id = ?", id)
		return err
	}
	_, err = p.db.Exec("UPDATE sprints SET estado = 'Activo' WHERE id = ?", id)
	return err
}

func (p *pasosHistorial) historiaCompletada(nombreSprint, titulo string, puntos any) (int64, error) {
	sprintID, err := p.sprint(nombreSprint)
	if err != nil {
		return 0, err
	}
	return p.insertar(
		`INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp, sprint_id, fecha_finalizacion)
		 VALUES (?, ?, 'Detalle', 'M', 'Completada', ?, ?, '2026-10-10')`,
		p.proyectoID, titulo, puntos, sprintID)
}

// historiaArrastrada deja la historia como queda al cerrar el Sprint: de vuelta
// en el Product Backlog y con la foto de sus Story Points (D2c).
func (p *pasosHistorial) historiaArrastrada(nombreSprint, titulo string, puntos any) (int64, error) {
	sprintID, err := p.sprint(nombreSprint)
	if err != nil {
		return 0, err
	}
	historia, err := p.insertar(
		"INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp) VALUES (?, ?, 'Detalle', 'M', 'Nueva', ?)",
		p.proyectoID, titulo, puntos)
	if err != nil {
		return 0, err
	}
	_, err = p.insertar("INSERT INTO historias_arrastradas (sprint_id, historia_id, proyecto_id, estimacion_sp) VALUES (?, ?, ?, ?)",
		sprintID, historia, p.proyectoID, puntos)
	return historia, err
}

func (p *pasosHistorial) sprintConCompletada(nombre string, puntos int) error {
	_, err := p.historiaCompletada(nombre, "Historia", puntos)
	return err
}

func (p *pasosHistorial) sprintConArrastrada(nombre string, puntos int) error {
	_, err := p.historiaArrastrada(nombre, "Historia", puntos)
	return err
}

func (p *pasosHistorial) sprintConCompletadaNombrada(nombre, titulo string, puntos int) error {
	id, err := p.historiaCompletada(nombre, titulo, puntos)
	p.historias[titulo] = id
	return err
}

func (p *pasosHistorial) sprintConArrastradaNombrada(nombre, titulo string, puntos int) error {
	id, err := p.historiaArrastrada(nombre, titulo, puntos)
	p.historias[titulo] = id
	return err
}

func (p *pasosHistorial) imputar(historiaID, sprintID int64, minutos int) error {
	_, err := p.insertar(
		`INSERT INTO registros_horas (proyecto_id, historia_id, usuario_id, sprint_id, fecha, actividad, minutos, creado_en)
		 VALUES (?, ?, ?, ?, '2026-10-10', 'Trabajo', ?, '2026-10-10T00:00:00Z')`,
		p.proyectoID, historiaID, p.usuarioID, sprintID, minutos)
	return err
}

// historiaAuxiliar sostiene horas y defectos sin sumar historias comprometidas.
func (p *pasosHistorial) historiaAuxiliar() (int64, error) {
	if p.auxiliar != 0 {
		return p.auxiliar, nil
	}
	id, err := p.insertar("INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado) VALUES (?, 'Auxiliar', '', 'M', 'Nueva')", p.proyectoID)
	p.auxiliar = id
	return id, err
}

func (p *pasosHistorial) sprintConHoras(nombre string, horas int) error {
	sprintID, err := p.sprint(nombre)
	if err != nil {
		return err
	}
	historia, err := p.historiaAuxiliar()
	if err != nil {
		return err
	}
	return p.imputar(historia, sprintID, horas*60)
}

func (p *pasosHistorial) historiaConMinutos(titulo string, minutos int, nombreSprint string) error {
	sprintID, err := p.sprint(nombreSprint)
	if err != nil {
		return err
	}
	return p.imputar(p.historias[titulo], sprintID, minutos)
}

func (p *pasosHistorial) insertarDefecto(descripcion string, historiaID, deteccion int64, resolucion any) error {
	estado := "Pendiente"
	if resolucion != nil {
		estado = "Resuelto"
	}
	_, err := p.insertar(
		`INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id, creado_en)
		 VALUES (?, ?, ?, 'Alta', ?, ?, ?, '2026-10-10T00:00:00Z')`,
		p.proyectoID, historiaID, descripcion, estado, deteccion, resolucion)
	return err
}

func (p *pasosHistorial) sprintConDefectos(nombre string, detectados, resueltos int) error {
	sprintID, err := p.sprint(nombre)
	if err != nil {
		return err
	}
	historia, err := p.historiaAuxiliar()
	if err != nil {
		return err
	}
	for i := range detectados {
		var resolucion any
		if i < resueltos {
			resolucion = sprintID
		}
		if err := p.insertarDefecto("Falla", historia, sprintID, resolucion); err != nil {
			return err
		}
	}
	return nil
}

func (p *pasosHistorial) defectoEntreSprints(descripcion, deteccion, resolucion string) error {
	detectado, err := p.sprint(deteccion)
	if err != nil {
		return err
	}
	resuelto, err := p.sprint(resolucion)
	if err != nil {
		return err
	}
	historia, err := p.historiaAuxiliar()
	if err != nil {
		return err
	}
	return p.insertarDefecto(descripcion, historia, detectado, resuelto)
}

// historiaPasaATener edita los Story Points con el caso de uso real: la
// historia arrastrada volvió al Product Backlog, así que D2b lo permite.
func (p *pasosHistorial) historiaPasaATener(titulo string, puntos int) error {
	actor, err := integranteDeLaHistoria(p.db, p.historias[titulo])
	if err != nil {
		return err
	}
	_, err = p.editarHist.Ejecutar(context.Background(), service.EditarHistoriaInput{
		ID: p.historias[titulo], Titulo: titulo, Descripcion: "Detalle", EstimacionSP: &puntos, UsuarioID: actor,
	})
	return err
}

func (p *pasosHistorial) historiaSeCompleta(titulo, nombreSprint string) error {
	sprintID, err := p.sprint(nombreSprint)
	if err != nil {
		return err
	}
	_, err = p.db.Exec("UPDATE historias SET sprint_id = ?, estado = 'Completada', fecha_finalizacion = '2026-10-20' WHERE id = ?", sprintID, p.historias[titulo])
	return err
}

func (p *pasosHistorial) sprintActivoConDatos(nombre string) error {
	if err := p.sprintEnEstado(nombre, "activo"); err != nil {
		return err
	}
	sprintID := p.sprints[nombre]
	historia, err := p.insertar(
		"INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp, sprint_id) VALUES (?, 'En curso', '', 'M', 'En progreso', 3, ?)",
		p.proyectoID, sprintID)
	if err != nil {
		return err
	}
	if err := p.imputar(historia, sprintID, 45); err != nil {
		return err
	}
	if err := p.insertarDefecto("Pendiente", historia, sprintID, nil); err != nil {
		return err
	}
	p.fotoAntes, err = p.fotoDelSprint(sprintID)
	return err
}

// fotoDelSprint vuelca todas las filas del Sprint y de sus historias, horas y
// defectos, para comparar antes y después de consultar.
func (p *pasosHistorial) fotoDelSprint(sprintID int64) (string, error) {
	consultas := []string{
		"SELECT * FROM sprints WHERE id = ?",
		"SELECT * FROM historias WHERE sprint_id = ? ORDER BY id",
		"SELECT * FROM registros_horas WHERE sprint_id = ? ORDER BY id",
		"SELECT * FROM defectos WHERE sprint_deteccion_id = ? ORDER BY id",
	}
	var foto strings.Builder
	for _, consulta := range consultas {
		filas, err := p.db.Query(consulta, sprintID)
		if err != nil {
			return "", err
		}
		columnas, err := filas.Columns()
		if err != nil {
			filas.Close()
			return "", err
		}
		for filas.Next() {
			valores := make([]any, len(columnas))
			punteros := make([]any, len(columnas))
			for i := range valores {
				punteros[i] = &valores[i]
			}
			if err := filas.Scan(punteros...); err != nil {
				filas.Close()
				return "", err
			}
			fmt.Fprintf(&foto, "%v\n", valores)
		}
		filas.Close()
	}
	return foto.String(), nil
}

func (p *pasosHistorial) otroProyectoConSprint() error {
	otro, err := p.crearProyecto(nil)
	if err != nil {
		return err
	}
	p.sprintAjeno, err = p.insertar(
		"INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (?, 'Ajeno', '2026-10-05', '2026-10-19', 'Finalizado')", otro)
	return err
}

func (p *pasosHistorial) consultarHistorialDe(proyectoID int64) error {
	p.historial, p.errConsulta = p.listar.Ejecutar(context.Background(), proyectoID)
	return nil
}

func (p *pasosHistorial) consultarHistorial() error { return p.consultarHistorialDe(p.proyectoID) }

func (p *pasosHistorial) consultarHistorialInexistente() error {
	return p.consultarHistorialDe(p.proyectoID + 1000)
}

func (p *pasosHistorial) consultarDetalleDe(sprintID int64) error {
	p.detalle, p.errConsulta = p.detallar.Ejecutar(context.Background(), p.proyectoID, sprintID)
	return nil
}

func (p *pasosHistorial) consultarDetalle(nombre string) error {
	id, existe := p.sprints[nombre]
	if !existe {
		return fmt.Errorf("el Sprint %s no se sembró", nombre)
	}
	return p.consultarDetalleDe(id)
}

func (p *pasosHistorial) consultarDetalleInexistente() error { return p.consultarDetalleDe(99999) }

func (p *pasosHistorial) consultarDetalleAjeno() error { return p.consultarDetalleDe(p.sprintAjeno) }

func (p *pasosHistorial) sinError() error {
	if p.errConsulta != nil {
		return fmt.Errorf("no se esperaba error en la consulta: %w", p.errConsulta)
	}
	return nil
}

func (p *pasosHistorial) historialLista(nombres string) error {
	if err := p.sinError(); err != nil {
		return err
	}
	var esperados []int64
	for _, nombre := range strings.Split(nombres, ",") {
		esperados = append(esperados, p.sprints[strings.TrimSpace(nombre)])
	}
	var obtenidos []int64
	for _, h := range p.historial {
		obtenidos = append(obtenidos, h.Sprint.ID)
	}
	if fmt.Sprint(obtenidos) != fmt.Sprint(esperados) {
		return fmt.Errorf("se esperaban los Sprints %v (%s), se obtuvieron %v", esperados, nombres, obtenidos)
	}
	return nil
}

func (p *pasosHistorial) historialVacio() error {
	if err := p.sinError(); err != nil {
		return err
	}
	if p.historial == nil || len(p.historial) != 0 {
		return fmt.Errorf("se esperaba un historial vacío, se obtuvo %#v", p.historial)
	}
	return nil
}

func datosDelSprintCompletos(s domain.Sprint) error {
	if s.Numero <= 0 || s.SprintGoal == "" || s.FechaInicio.IsZero() || s.FechaFin.IsZero() {
		return fmt.Errorf("el Sprint %d no informa número, Sprint Goal y fechas: %+v", s.ID, s)
	}
	return nil
}

func (p *pasosHistorial) cadaSprintInformaSusDatos() error {
	if err := p.sinError(); err != nil {
		return err
	}
	for _, h := range p.historial {
		if err := datosDelSprintCompletos(h.Sprint); err != nil {
			return err
		}
	}
	return nil
}

func (p *pasosHistorial) metricasDelHistorial(nombre string) (domain.Metricas, error) {
	if err := p.sinError(); err != nil {
		return domain.Metricas{}, err
	}
	for _, h := range p.historial {
		if h.Sprint.ID == p.sprints[nombre] {
			return h.Metricas, nil
		}
	}
	return domain.Metricas{}, fmt.Errorf("el historial no incluye el Sprint %s", nombre)
}

func (p *pasosHistorial) sprintInformaMetricas(nombre string, tabla *godog.Table) error {
	m, err := p.metricasDelHistorial(nombre)
	if err != nil {
		return err
	}
	return compararMetricas(m, tabla)
}

func (p *pasosHistorial) sprintSinDatosDeHoras(nombre string) error {
	m, err := p.metricasDelHistorial(nombre)
	if err != nil {
		return err
	}
	if m.HorasEstimadas != nil || m.DesviacionHoras != nil || m.DesviacionPorcentaje != nil {
		return fmt.Errorf("se esperaban horas estimadas y desviación sin datos, se obtuvo %+v", m)
	}
	return nil
}

func (p *pasosHistorial) metricasDeHU03() (map[int64]domain.Metricas, error) {
	calculadas, err := p.calcular.Ejecutar(context.Background(), p.proyectoID)
	if err != nil {
		return nil, err
	}
	porSprint := map[int64]domain.Metricas{}
	for _, s := range calculadas.Sprints {
		porSprint[s.SprintID] = s.Metricas
	}
	return porSprint, nil
}

func (p *pasosHistorial) metricasCoincidenConHU03() error {
	if err := p.sinError(); err != nil {
		return err
	}
	esperadas, err := p.metricasDeHU03()
	if err != nil {
		return err
	}
	if len(p.historial) != len(esperadas) {
		return fmt.Errorf("el historial tiene %d Sprints y el cálculo de métricas %d", len(p.historial), len(esperadas))
	}
	for _, h := range p.historial {
		if !reflect.DeepEqual(h.Metricas, esperadas[h.Sprint.ID]) {
			return fmt.Errorf("Sprint %d: el historial informa %+v y HU-03 %+v", h.Sprint.ID, h.Metricas, esperadas[h.Sprint.ID])
		}
	}
	return nil
}

func (p *pasosHistorial) detalleInformaMetricasDeHU03() error {
	if err := p.sinError(); err != nil {
		return err
	}
	esperadas, err := p.metricasDeHU03()
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(p.detalle.Metricas, esperadas[p.detalle.Sprint.ID]) {
		return fmt.Errorf("el detalle informa %+v y HU-03 %+v", p.detalle.Metricas, esperadas[p.detalle.Sprint.ID])
	}
	return nil
}

func (p *pasosHistorial) detalleInformaDatos(goal string) error {
	if err := p.detalleInformaMetricasDeHU03(); err != nil {
		return err
	}
	if p.detalle.Sprint.SprintGoal != goal {
		return fmt.Errorf("se esperaba el Sprint Goal %q, se obtuvo %q", goal, p.detalle.Sprint.SprintGoal)
	}
	return datosDelSprintCompletos(p.detalle.Sprint)
}

func siONo(valor bool) string {
	if valor {
		return "sí"
	}
	return "no"
}

func (p *pasosHistorial) detalleInformaHistorias(tabla *godog.Table) error {
	if err := p.sinError(); err != nil {
		return err
	}
	var obtenidas []string
	for _, h := range p.detalle.Historias {
		puntos := "sin estimar"
		if h.EstimacionSP != nil {
			puntos = strconv.Itoa(*h.EstimacionSP)
		}
		obtenidas = append(obtenidas, strings.Join([]string{h.Titulo, puntos, siONo(h.Completada)}, " | "))
	}
	return compararFilas(tabla, obtenidas)
}

func (p *pasosHistorial) detalleInformaEsfuerzo(tabla *godog.Table) error {
	if err := p.sinError(); err != nil {
		return err
	}
	var obtenidas []string
	for _, h := range p.detalle.Historias {
		horas := strconv.FormatFloat(h.Horas, 'f', -1, 64)
		obtenidas = append(obtenidas, strings.Join([]string{h.Titulo, strconv.Itoa(h.Minutos), horas}, " | "))
	}
	return compararFilas(tabla, obtenidas)
}

func (p *pasosHistorial) minutosCoincidenConHorasReales() error {
	total := 0
	for _, h := range p.detalle.Historias {
		total += h.Minutos
	}
	horas := math.Round(float64(total)/60*100) / 100
	if horas != p.detalle.Metricas.HorasReales {
		return fmt.Errorf("la suma de minutos (%d = %v h) no coincide con las horas reales %v", total, horas, p.detalle.Metricas.HorasReales)
	}
	return nil
}

func (p *pasosHistorial) detalleInformaDefectos(tabla *godog.Table) error {
	if err := p.sinError(); err != nil {
		return err
	}
	var obtenidas []string
	for _, d := range p.detalle.Defectos {
		obtenidas = append(obtenidas, strings.Join([]string{d.Descripcion, siONo(d.DetectadoEnSprint), siONo(d.ResueltoEnSprint)}, " | "))
	}
	return compararFilas(tabla, obtenidas)
}

// compararFilas compara las filas de la tabla, sin el encabezado, con las
// obtenidas, en el mismo orden.
func compararFilas(tabla *godog.Table, obtenidas []string) error {
	var esperadas []string
	for _, fila := range tabla.Rows[1:] {
		var celdas []string
		for _, celda := range fila.Cells {
			celdas = append(celdas, celda.Value)
		}
		esperadas = append(esperadas, strings.Join(celdas, " | "))
	}
	if strings.Join(obtenidas, "\n") != strings.Join(esperadas, "\n") {
		return fmt.Errorf("se esperaba:\n%s\nse obtuvo:\n%s", strings.Join(esperadas, "\n"), strings.Join(obtenidas, "\n"))
	}
	return nil
}

func (p *pasosHistorial) detalleVacio() error {
	if err := p.sinError(); err != nil {
		return err
	}
	if p.detalle.Historias == nil || len(p.detalle.Historias) != 0 || p.detalle.Defectos == nil || len(p.detalle.Defectos) != 0 {
		return fmt.Errorf("se esperaban historias y defectos vacíos, se obtuvo %+v y %+v", p.detalle.Historias, p.detalle.Defectos)
	}
	return nil
}

func (p *pasosHistorial) sprintNoCambio(nombre string) error {
	despues, err := p.fotoDelSprint(p.sprints[nombre])
	if err != nil {
		return err
	}
	if despues != p.fotoAntes {
		return fmt.Errorf("el Sprint %s cambió:\nantes:\n%s\ndespués:\n%s", nombre, p.fotoAntes, despues)
	}
	return nil
}

func (p *pasosHistorial) consultaFalla(mensaje string) error {
	if p.errConsulta == nil {
		return errors.New("se esperaba un error en la consulta del historial")
	}
	if p.errConsulta.Error() != mensaje {
		return fmt.Errorf("se esperaba el error %q, se obtuvo %q", mensaje, p.errConsulta.Error())
	}
	return nil
}

func (p *pasosHistorial) consultaFallaEnCampo(campo, mensaje string) error {
	var verr domain.ValidationError
	if !errors.As(p.errConsulta, &verr) {
		return fmt.Errorf("se esperaba un error de validación, se obtuvo %v", p.errConsulta)
	}
	if verr.Campo != campo || verr.Mensaje != mensaje {
		return fmt.Errorf("se esperaba el error %q en %q, se obtuvo %q en %q", mensaje, campo, verr.Mensaje, verr.Campo)
	}
	return nil
}

// InitializeScenarioHistorial registra los pasos de HU-12 (@HU-12).
func InitializeScenarioHistorial(ctx *godog.ScenarioContext) {
	p := &pasosHistorial{}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return goctx, p.preparar()
	})
	ctx.After(func(goctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if p.db != nil {
			return goctx, p.db.Close()
		}
		return goctx, nil
	})

	ctx.Step(`^un proyecto del historial con ([\d.]+) horas por Story Point$`, p.proyectoConHoras)
	ctx.Step(`^el proyecto del historial no tiene horas por Story Point$`, p.proyectoSinHoras)
	ctx.Step(`^el Sprint del historial "([^"]*)" con Sprint Goal "([^"]*)" está finalizado$`, p.sprintFinalizadoConGoal)
	ctx.Step(`^el Sprint del historial "([^"]*)" está (activo|pendiente)$`, p.sprintEnEstado)
	ctx.Step(`^el Sprint del historial "([^"]*)" está activo con historias, horas y defectos$`, p.sprintActivoConDatos)
	ctx.Step(`^el Sprint del historial "([^"]*)" tiene una historia completada de (\d+) Story Points$`, p.sprintConCompletada)
	ctx.Step(`^el Sprint del historial "([^"]*)" tiene una historia arrastrada de (\d+) Story Points$`, p.sprintConArrastrada)
	ctx.Step(`^el Sprint del historial "([^"]*)" tiene la historia completada "([^"]*)" de (\d+) Story Points$`, p.sprintConCompletadaNombrada)
	ctx.Step(`^el Sprint del historial "([^"]*)" tiene la historia arrastrada "([^"]*)" de (\d+) Story Points$`, p.sprintConArrastradaNombrada)
	ctx.Step(`^el Sprint del historial "([^"]*)" tiene (\d+) horas imputadas$`, p.sprintConHoras)
	ctx.Step(`^el Sprint del historial "([^"]*)" tiene (\d+) defectos detectados y (\d+) resueltos?$`, p.sprintConDefectos)
	ctx.Step(`^la historia "([^"]*)" del historial tiene (\d+) minutos imputados en el Sprint "([^"]*)"$`, p.historiaConMinutos)
	ctx.Step(`^la historia "([^"]*)" del historial pasa a tener (\d+) Story Points$`, p.historiaPasaATener)
	ctx.Step(`^la historia "([^"]*)" del historial se completa en el Sprint "([^"]*)"$`, p.historiaSeCompleta)
	ctx.Step(`^un defecto del historial "([^"]*)" detectado en el Sprint "([^"]*)" y resuelto en el Sprint "([^"]*)"$`, p.defectoEntreSprints)
	ctx.Step(`^otro proyecto del historial con un Sprint finalizado$`, p.otroProyectoConSprint)

	ctx.Step(`^consulto el historial de Sprints del proyecto$`, p.consultarHistorial)
	ctx.Step(`^consulto el historial de Sprints de un proyecto inexistente$`, p.consultarHistorialInexistente)
	ctx.Step(`^consulto el detalle del Sprint "([^"]*)" en el historial$`, p.consultarDetalle)
	ctx.Step(`^consulto el detalle de un Sprint inexistente en el historial$`, p.consultarDetalleInexistente)
	ctx.Step(`^consulto en el historial del proyecto el detalle del Sprint del otro proyecto$`, p.consultarDetalleAjeno)

	ctx.Step(`^el historial lista los Sprints "([^"]*)" en ese orden$`, p.historialLista)
	ctx.Step(`^el historial no lista ningún Sprint$`, p.historialVacio)
	ctx.Step(`^cada Sprint del historial informa su número, Sprint Goal, fecha de inicio y fecha de fin$`, p.cadaSprintInformaSusDatos)
	ctx.Step(`^el Sprint "([^"]*)" del historial informa estas métricas:$`, p.sprintInformaMetricas)
	ctx.Step(`^el Sprint "([^"]*)" del historial informa sin datos las horas estimadas y la desviación$`, p.sprintSinDatosDeHoras)
	ctx.Step(`^las métricas de cada Sprint del historial coinciden con las del cálculo de métricas del proyecto$`, p.metricasCoincidenConHU03)
	ctx.Step(`^el detalle informa el número, el Sprint Goal "([^"]*)", las fechas y las métricas del Sprint$`, p.detalleInformaDatos)
	ctx.Step(`^el detalle informa estas historias comprometidas:$`, p.detalleInformaHistorias)
	ctx.Step(`^el detalle informa este esfuerzo por historia:$`, p.detalleInformaEsfuerzo)
	ctx.Step(`^la suma de los minutos del detalle coincide con las horas reales del Sprint$`, p.minutosCoincidenConHorasReales)
	ctx.Step(`^el detalle informa estos defectos:$`, p.detalleInformaDefectos)
	ctx.Step(`^el detalle no informa historias ni defectos$`, p.detalleVacio)
	ctx.Step(`^el detalle informa las métricas del cálculo de métricas para ese Sprint$`, p.detalleInformaMetricasDeHU03)
	ctx.Step(`^el Sprint "([^"]*)" del historial, sus historias, sus horas y sus defectos no cambiaron$`, p.sprintNoCambio)
	ctx.Step(`^la consulta del historial falla con "([^"]*)"$`, p.consultaFalla)
	ctx.Step(`^la consulta del historial falla en el campo "([^"]*)" con "([^"]*)"$`, p.consultaFallaEnCampo)
}
