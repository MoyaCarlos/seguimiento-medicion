package features

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/report"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// generadorPDF abstrae el generador de reportes para poder inyectar uno que falla
// en el escenario de error real (HU-10).
type generadorPDF interface {
	GenerarSprint(domain.ReporteSprint) ([]byte, error)
	GenerarProyecto(domain.ReporteProyecto) ([]byte, error)
}

type generadorQueFalla struct{}

func (generadorQueFalla) GenerarSprint(domain.ReporteSprint) ([]byte, error) {
	return nil, errors.New("Maroto no pudo armar el documento")
}

func (generadorQueFalla) GenerarProyecto(domain.ReporteProyecto) ([]byte, error) {
	return nil, errors.New("Maroto no pudo armar el documento")
}

// pasosReporte ejercita los reportes de HU-10 sobre una base SQLite en memoria,
// con las migraciones y los services reales. Siembra el historial con SQL para
// armar Sprints finalizados con sus historias, horas y defectos.
type pasosReporte struct {
	db       *sql.DB
	generar  *service.GenerarReporte
	ahora    time.Time
	usuario  int64
	proyecto int64
	sprints  map[string]int64
	auxiliar int64

	generador       generadorPDF
	reporteSprint   domain.ReporteSprint
	reporteProyecto domain.ReporteProyecto
	pdf             []byte
	err             error
}

func (p *pasosReporte) preparar() error {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return err
	}
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		return err
	}
	proyectos := repository.NewSQLiteProjectRepository(db)
	sprints := repository.NewSQLiteSprintRepository(db)
	historias := repository.NewSQLiteBacklogRepository(db)
	metricas := repository.NewSQLiteMetricasRepository(db)
	defectos := repository.NewSQLiteDefectoRepository(db)

	*p = pasosReporte{
		db:        db,
		generar:   service.NewGenerarReporte(proyectos, sprints, historias, historias, historias, metricas, defectos),
		ahora:     time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC),
		sprints:   map[string]int64{},
		generador: report.Generador{},
	}
	return nil
}

func (p *pasosReporte) insertar(consulta string, args ...any) (int64, error) {
	resultado, err := p.db.Exec(consulta, args...)
	if err != nil {
		return 0, fmt.Errorf("sembrar %q: %w", consulta, err)
	}
	return resultado.LastInsertId()
}

func (p *pasosReporte) existeProyecto(id int) error {
	proyecto, err := p.insertar(
		"INSERT INTO proyectos (id, nombre, descripcion, creado_en) VALUES (?, 'Proyecto', '', '2026-10-01T00:00:00Z')",
		id)
	if err != nil {
		return err
	}
	p.proyecto = proyecto
	return nil
}

func (p *pasosReporte) proyectoConHoras(horas float64) error {
	_, err := p.db.Exec("UPDATE proyectos SET horas_por_story_point = ? WHERE id = ?", horas, p.proyecto)
	return err
}

func (p *pasosReporte) proyectoConIntegrante(nombre, rol string) error {
	usuario, err := p.insertar(
		"INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES (?, ?, '2026-10-01T00:00:00Z')",
		nombre, domain.NormalizarNombre(nombre))
	if err != nil {
		return err
	}
	_, err = p.insertar(
		"INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, ?, '2026-10-01T00:00:00Z')",
		p.proyecto, usuario, rol)
	if err != nil {
		return err
	}
	p.usuario = usuario
	return nil
}

func (p *pasosReporte) crearSprint(nombre, estado, objetivo string) error {
	var (
		id  int64
		err error
	)
	if estado == string(domain.SprintPendiente) {
		id, err = p.insertar(
			"INSERT INTO sprints (proyecto_id, estado) VALUES (?, 'Pendiente')",
			p.proyecto)
	} else {
		id, err = p.insertar(
			`INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado)
			 VALUES (?, ?, '2026-10-01', '2026-10-15', ?)`,
			p.proyecto, objetivo, estado)
	}
	if err != nil {
		return err
	}
	p.sprints[nombre] = id
	return nil
}

func (p *pasosReporte) sprintFinalizado(nombre, objetivo string) error {
	return p.crearSprint(nombre, string(domain.SprintFinalizado), objetivo)
}

func (p *pasosReporte) sprintEnEstado(nombre, estado string) error {
	return p.crearSprint(nombre, estado, "Objetivo")
}

func (p *pasosReporte) sprintEnEstadoConObjetivo(nombre, estado, objetivo string) error {
	return p.crearSprint(nombre, estado, objetivo)
}

func (p *pasosReporte) historiaEnSprint(titulo, sprint, puntos, estado string) error {
	sprintID, existe := p.sprints[sprint]
	if !existe {
		return fmt.Errorf("no existe el Sprint %q", sprint)
	}
	var finalizacion any
	if estado == string(domain.EstadoCompletada) {
		finalizacion = "2026-10-10"
	}
	_, err := p.insertar(
		`INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp, sprint_id, fecha_finalizacion)
		 VALUES (?, ?, 'Detalle', 'M', ?, ?, ?, ?)`,
		p.proyecto, titulo, estado, puntos, sprintID, finalizacion)
	return err
}

func (p *pasosReporte) historiaArrastrada(titulo, sprint, puntos string) error {
	sprintID, existe := p.sprints[sprint]
	if !existe {
		return fmt.Errorf("no existe el Sprint %q", sprint)
	}
	historiasPuntos := p.numero(puntos)
	historia, err := p.insertar(
		`INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp)
		 VALUES (?, ?, 'Detalle', 'M', 'Nueva', ?)`,
		p.proyecto, titulo, historiasPuntos)
	if err != nil {
		return err
	}
	_, err = p.insertar(
		"INSERT INTO historias_arrastradas (sprint_id, historia_id, proyecto_id, estimacion_sp) VALUES (?, ?, ?, ?)",
		sprintID, historia, p.proyecto, historiasPuntos)
	return err
}

func (p *pasosReporte) historiaEnBacklog(titulo, prioridad string, puntos int) error {
	_, err := p.insertar(
		`INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp)
		 VALUES (?, ?, 'Detalle', ?, 'Nueva', ?)`,
		p.proyecto, titulo, prioridad, puntos)
	return err
}

// auxiliar sostiene las horas y los defectos sin sumar historias comprometidas.
func (p *pasosReporte) historiaAuxiliar() (int64, error) {
	if p.auxiliar != 0 {
		return p.auxiliar, nil
	}
	id, err := p.insertar(
		`INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado)
		 VALUES (?, 'Auxiliar', 'Detalle', 'W', 'Nueva')`,
		p.proyecto)
	p.auxiliar = id
	return id, err
}

func (p *pasosReporte) registrarHoras(sprint string, horas int) error {
	sprintID, existe := p.sprints[sprint]
	if !existe {
		return fmt.Errorf("no existe el Sprint %q", sprint)
	}
	historia, err := p.historiaAuxiliar()
	if err != nil {
		return err
	}
	_, err = p.insertar(
		`INSERT INTO registros_horas (proyecto_id, historia_id, usuario_id, sprint_id, fecha, actividad, minutos, creado_en)
		 VALUES (?, ?, ?, ?, '2026-10-10', 'Trabajo', ?, '2026-10-10T00:00:00Z')`,
		p.proyecto, historia, p.usuario, sprintID, horas*60)
	return err
}

func (p *pasosReporte) defectoDetectado(sprint, descripcion, severidad string) error {
	sprintID, existe := p.sprints[sprint]
	if !existe {
		return fmt.Errorf("no existe el Sprint %q", sprint)
	}
	historia, err := p.historiaAuxiliar()
	if err != nil {
		return err
	}
	_, err = p.insertar(
		`INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, creado_en)
		 VALUES (?, ?, ?, ?, 'Pendiente', ?, '2026-10-10T00:00:00Z')`,
		p.proyecto, historia, descripcion, severidad, sprintID)
	return err
}

func (p *pasosReporte) defectoResuelto(sprint, descripcion, severidad string) error {
	sprintID, existe := p.sprints[sprint]
	if !existe {
		return fmt.Errorf("no existe el Sprint %q", sprint)
	}
	resultado, err := p.db.Exec(
		`UPDATE defectos SET estado = 'Resuelto', sprint_resolucion_id = ?
		 WHERE sprint_deteccion_id = ? AND descripcion = ? AND estado = 'Pendiente'`,
		sprintID, sprintID, descripcion)
	if err != nil {
		return err
	}
	afectadas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if afectadas > 0 {
		return nil
	}
	historia, err := p.historiaAuxiliar()
	if err != nil {
		return err
	}
	_, err = p.insertar(
		`INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id, creado_en)
		 VALUES (?, ?, ?, ?, 'Resuelto', ?, ?, '2026-10-10T00:00:00Z')`,
		p.proyecto, historia, descripcion, severidad, sprintID, sprintID)
	return err
}

func (p *pasosReporte) defectoEntreSprints(descripcion, deteccion, resolucion string) error {
	detectado, existe := p.sprints[deteccion]
	if !existe {
		return fmt.Errorf("no existe el Sprint %q", deteccion)
	}
	resuelto, existe := p.sprints[resolucion]
	if !existe {
		return fmt.Errorf("no existe el Sprint %q", resolucion)
	}
	historia, err := p.historiaAuxiliar()
	if err != nil {
		return err
	}
	_, err = p.insertar(
		`INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id, creado_en)
		 VALUES (?, ?, ?, 'Alta', 'Resuelto', ?, ?, '2026-10-10T00:00:00Z')`,
		p.proyecto, historia, descripcion, detectado, resuelto)
	return err
}

func (p *pasosReporte) numero(valor string) any {
	if valor == "" {
		return nil
	}
	numero, err := strconv.Atoi(valor)
	if err != nil {
		return nil
	}
	return numero
}

func (p *pasosReporte) reiniciarSolicitud() {
	p.reporteSprint = domain.ReporteSprint{}
	p.reporteProyecto = domain.ReporteProyecto{}
	p.pdf = nil
	p.err = nil
}

func (p *pasosReporte) solicitarSprint(nombre string) error {
	p.reiniciarSolicitud()
	id, existe := p.sprints[nombre]
	if !existe {
		return fmt.Errorf("no existe el Sprint %q", nombre)
	}
	return p.ejecutarSprint(id)
}

func (p *pasosReporte) solicitarSprintInexistente(id int) error {
	p.reiniciarSolicitud()
	return p.ejecutarSprint(int64(id))
}

func (p *pasosReporte) ejecutarSprint(id int64) error {
	reporte, err := p.generar.EjecutarSprint(context.Background(), id, p.ahora)
	p.reporteSprint = reporte
	if err != nil {
		p.err = err
		return nil
	}
	contenido, err := p.generador.GenerarSprint(reporte)
	if err != nil {
		p.err = domain.ErrReporteNoGenerado
		return nil
	}
	p.pdf = contenido
	return nil
}

func (p *pasosReporte) solicitarProyecto() error {
	p.reiniciarSolicitud()
	return p.ejecutarProyecto(p.proyecto)
}

func (p *pasosReporte) solicitarProyectoInexistente(id int) error {
	p.reiniciarSolicitud()
	return p.ejecutarProyecto(int64(id))
}

func (p *pasosReporte) solicitarProyectoConGeneradorQueFalla() error {
	p.generador = generadorQueFalla{}
	return p.solicitarProyecto()
}

func (p *pasosReporte) ejecutarProyecto(id int64) error {
	reporte, err := p.generar.EjecutarProyecto(context.Background(), id, p.ahora)
	p.reporteProyecto = reporte
	if err != nil {
		p.err = err
		return nil
	}
	contenido, err := p.generador.GenerarProyecto(reporte)
	if err != nil {
		p.err = domain.ErrReporteNoGenerado
		return nil
	}
	p.pdf = contenido
	return nil
}

func (p *pasosReporte) pdfValido() error {
	if p.err != nil {
		return fmt.Errorf("se esperaba un PDF, se obtuvo el error %v", p.err)
	}
	if !bytes.HasPrefix(p.pdf, []byte("%PDF-")) {
		return fmt.Errorf("el contenido no empieza con %%PDF-")
	}
	if !bytes.Contains(p.pdf, []byte("%%EOF")) {
		return fmt.Errorf("el PDF no contiene %%EOF")
	}
	return nil
}

func (p *pasosReporte) pdfDeSprint(nombre string) error {
	if err := p.pdfValido(); err != nil {
		return err
	}
	if p.reporteSprint.Sprint.Sprint.ID != p.sprints[nombre] {
		return fmt.Errorf("el reporte no es del Sprint %q", nombre)
	}
	return nil
}

func (p *pasosReporte) pdfDeProyecto() error {
	if err := p.pdfValido(); err != nil {
		return err
	}
	if p.reporteProyecto.Proyecto.ID != p.proyecto {
		return fmt.Errorf("el reporte no es del proyecto %d", p.proyecto)
	}
	return nil
}

func (p *pasosReporte) sprintIncluyeCompletada(titulo string) error {
	for _, historia := range p.reporteSprint.Sprint.Planificadas {
		if historia.Titulo == titulo && historia.Completada {
			return nil
		}
	}
	return fmt.Errorf("el reporte no incluye la historia completada %q: %+v", titulo, p.reporteSprint.Sprint.Planificadas)
}

func (p *pasosReporte) sprintIncluyePlanificada(titulo string) error {
	for _, historia := range p.reporteSprint.Sprint.Planificadas {
		if historia.Titulo == titulo && !historia.Completada {
			return nil
		}
	}
	return fmt.Errorf("el reporte no incluye la historia planificada %q como no completada: %+v", titulo, p.reporteSprint.Sprint.Planificadas)
}

func (p *pasosReporte) sprintInformaHoras(horas int) error {
	if p.reporteSprint.Sprint.Metricas.HorasReales != float64(horas) {
		return fmt.Errorf("se esperaban %d horas reales, se obtuvieron %v", horas, p.reporteSprint.Sprint.Metricas.HorasReales)
	}
	return nil
}

func (p *pasosReporte) sprintInformaDefectos() error {
	if len(p.reporteSprint.Sprint.DefectosDetectados) == 0 {
		return fmt.Errorf("se esperaba al menos un defecto detectado")
	}
	if len(p.reporteSprint.Sprint.DefectosResueltos) == 0 {
		return fmt.Errorf("se esperaba al menos un defecto resuelto")
	}
	return nil
}

func (p *pasosReporte) sprintSinDefectos() error {
	if len(p.reporteSprint.Sprint.DefectosDetectados) != 0 || len(p.reporteSprint.Sprint.DefectosResueltos) != 0 {
		return fmt.Errorf("no debía informar defectos: %d detectados, %d resueltos",
			len(p.reporteSprint.Sprint.DefectosDetectados), len(p.reporteSprint.Sprint.DefectosResueltos))
	}
	return nil
}

func (p *pasosReporte) solicitudRechazada(mensaje string) error {
	if p.err == nil {
		return fmt.Errorf("se esperaba el error %q, no hubo error", mensaje)
	}
	if obtenido := mensajeDeError(p.err); obtenido != mensaje {
		return fmt.Errorf("se esperaba el error %q, se obtuvo %q", mensaje, obtenido)
	}
	return nil
}

func mensajeDeError(err error) string {
	var verr domain.ValidationError
	if errors.As(err, &verr) {
		return verr.Mensaje
	}
	return err.Error()
}

func (p *pasosReporte) sinArchivo() error {
	if len(p.pdf) != 0 {
		return fmt.Errorf("no se debía descargar un archivo, se obtuvieron %d bytes", len(p.pdf))
	}
	return nil
}

func (p *pasosReporte) sprintListaDefecto(titulo, estado string) error {
	defectos := p.reporteSprint.Sprint.DefectosResueltos
	if estado == "detectado" {
		defectos = p.reporteSprint.Sprint.DefectosDetectados
	}
	for _, defecto := range defectos {
		if defecto.Descripcion == titulo {
			return nil
		}
	}
	return fmt.Errorf("el reporte no lista el defecto %q como %s", titulo, estado)
}

func (p *pasosReporte) proyectoIncluyeDatos() error {
	if p.reporteProyecto.Proyecto.ID != p.proyecto {
		return fmt.Errorf("datos del proyecto inesperados: %+v", p.reporteProyecto.Proyecto)
	}
	if len(p.reporteProyecto.Integrantes) != 1 || p.reporteProyecto.Integrantes[0].Nombre != "Ana" {
		return fmt.Errorf("integrantes inesperados: %+v", p.reporteProyecto.Integrantes)
	}
	if len(p.reporteProyecto.ProductBacklog) == 0 {
		return fmt.Errorf("se esperaba el Product Backlog con historias")
	}
	return nil
}

func (p *pasosReporte) proyectoIncluyeSprintsFinalizados(cantidad int) error {
	if len(p.reporteProyecto.Sprints) != cantidad {
		return fmt.Errorf("se esperaban %d Sprints finalizados, se obtuvieron %d", cantidad, len(p.reporteProyecto.Sprints))
	}
	for i := 1; i < len(p.reporteProyecto.Sprints); i++ {
		if p.reporteProyecto.Sprints[i-1].Sprint.Numero >= p.reporteProyecto.Sprints[i].Sprint.Numero {
			return fmt.Errorf("los Sprints no están en orden: %+v", p.reporteProyecto.Sprints)
		}
	}
	return nil
}

func (p *pasosReporte) proyectoIncluyeSprintActivo(nombre string) error {
	if p.reporteProyecto.SprintActivo == nil {
		return fmt.Errorf("se esperaba el Sprint Activo %q", nombre)
	}
	if p.reporteProyecto.SprintActivo.Sprint.ID != p.sprints[nombre] {
		return fmt.Errorf("el Sprint Activo no es %q: %+v", nombre, p.reporteProyecto.SprintActivo.Sprint)
	}
	if len(p.reporteProyecto.SprintActivo.Historias) == 0 {
		return fmt.Errorf("se esperaban las historias del Sprint Activo %q", nombre)
	}
	return nil
}

func (p *pasosReporte) proyectoNoSumaSprintActivo() error {
	var horas float64
	var completados, detectados int
	for _, sprint := range p.reporteProyecto.Sprints {
		horas += sprint.Metricas.HorasReales
		completados += sprint.Metricas.SPCompletados
		detectados += sprint.Metricas.DefectosDetectados
	}
	metricas := p.reporteProyecto.Metricas
	if metricas.HorasReales != horas || metricas.SPCompletados != completados || metricas.DefectosDetectados != detectados {
		return fmt.Errorf("las métricas del proyecto suman el Sprint Activo: proyecto=%+v, suma de finalizados=(horas %v, SP %d, defectos %d)",
			metricas, horas, completados, detectados)
	}
	return nil
}

func (p *pasosReporte) proyectoSinSprintActivo() error {
	if p.reporteProyecto.SprintActivo != nil {
		return fmt.Errorf("no se esperaba Sprint Activo: %+v", p.reporteProyecto.SprintActivo)
	}
	return nil
}

func (p *pasosReporte) proyectoSinHistoriasEnBacklog() error {
	if len(p.reporteProyecto.ProductBacklog) != 0 {
		return fmt.Errorf("no se esperaban historias en el Product Backlog: %+v", p.reporteProyecto.ProductBacklog)
	}
	return nil
}

func (p *pasosReporte) sprintDelProyectoSinHorasNiDefectos(nombre string) error {
	for _, sprint := range p.reporteProyecto.Sprints {
		if sprint.Sprint.ID != p.sprints[nombre] {
			continue
		}
		if sprint.Metricas.HorasReales != 0 {
			return fmt.Errorf("se esperaban 0 horas reales, se obtuvieron %v", sprint.Metricas.HorasReales)
		}
		if len(sprint.DefectosDetectados) != 0 || len(sprint.DefectosResueltos) != 0 {
			return fmt.Errorf("no se esperaban defectos en el Sprint %q", nombre)
		}
		return nil
	}
	return fmt.Errorf("el reporte no incluye el Sprint %q", nombre)
}

func (p *pasosReporte) proyectoListaDefecto(descripcion, deteccion, resolucion string) error {
	detectado := p.buscarSprintDelProyecto(deteccion)
	if detectado == nil {
		return fmt.Errorf("el reporte no incluye el Sprint %q", deteccion)
	}
	resuelto := p.buscarSprintDelProyecto(resolucion)
	if resuelto == nil {
		return fmt.Errorf("el reporte no incluye el Sprint %q", resolucion)
	}
	for _, defecto := range detectado.DefectosDetectados {
		if defecto.Descripcion == descripcion && defecto.SprintDeteccion == detectado.Sprint.Numero {
			for _, resueltoDefecto := range resuelto.DefectosResueltos {
				if resueltoDefecto.Descripcion == descripcion && resueltoDefecto.SprintResolucion != nil && *resueltoDefecto.SprintResolucion == resuelto.Sprint.Numero {
					return nil
				}
			}
			return fmt.Errorf("el defecto %q no figura como resuelto en %q", descripcion, resolucion)
		}
	}
	return fmt.Errorf("el defecto %q no figura como detectado en %q", descripcion, deteccion)
}

func (p *pasosReporte) buscarSprintDelProyecto(nombre string) *domain.SprintReporte {
	id, existe := p.sprints[nombre]
	if !existe {
		return nil
	}
	for i := range p.reporteProyecto.Sprints {
		if p.reporteProyecto.Sprints[i].Sprint.ID == id {
			return &p.reporteProyecto.Sprints[i]
		}
	}
	return nil
}

func (p *pasosReporte) proyectoIncluyeHistoriaEnBacklog(titulo string) error {
	for _, historia := range p.reporteProyecto.ProductBacklog {
		if historia.Titulo == titulo {
			return nil
		}
	}
	return fmt.Errorf("el reporte no incluye la historia %q en el Product Backlog", titulo)
}

func (p *pasosReporte) proyectoIncluyeHistoriaSoloEnActivo(titulo string) error {
	if p.reporteProyecto.SprintActivo == nil {
		return fmt.Errorf("se esperaba el Sprint Activo para contener %q", titulo)
	}
	enActivo := false
	for _, historia := range p.reporteProyecto.SprintActivo.Historias {
		if historia.Titulo == titulo {
			enActivo = true
			break
		}
	}
	if !enActivo {
		return fmt.Errorf("la historia %q no está en el Sprint Activo", titulo)
	}
	for _, historia := range p.reporteProyecto.ProductBacklog {
		if historia.Titulo == titulo {
			return fmt.Errorf("la historia %q también está en el Product Backlog", titulo)
		}
	}
	for _, sprint := range p.reporteProyecto.Sprints {
		for _, historia := range sprint.Planificadas {
			if historia.Titulo == titulo {
				return fmt.Errorf("la historia %q también está en el Sprint %d", titulo, sprint.Sprint.Numero)
			}
		}
	}
	return nil
}

func (p *pasosReporte) informaSinDatos() error {
	if !p.reporteProyecto.SinSprintsFinalizados {
		return fmt.Errorf("se esperaba SinSprintsFinalizados=true")
	}
	if len(p.reporteProyecto.Sprints) != 0 {
		return fmt.Errorf("no se esperaban Sprints finalizados")
	}
	metricas := p.reporteProyecto.Metricas
	if metricas.HorasReales != 0 || metricas.SPPlanificados != 0 || metricas.SPCompletados != 0 ||
		metricas.DefectosDetectados != 0 || metricas.DefectosResueltos != 0 {
		return fmt.Errorf("se esperaban métricas sin datos, se obtuvo %+v", metricas)
	}
	return nil
}

// InitializeScenarioReporte registra los pasos de HU-10 (@HU-10).
func InitializeScenarioReporte(ctx *godog.ScenarioContext) {
	p := &pasosReporte{}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return goctx, p.preparar()
	})
	ctx.After(func(goctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if p.db != nil {
			return goctx, p.db.Close()
		}
		return goctx, nil
	})

	ctx.Step(`^un proyecto de reporte con identificador (\d+)$`, p.existeProyecto)
	ctx.Step(`^el proyecto de reporte tiene ([\d.]+) horas por Story Point$`, p.proyectoConHoras)
	ctx.Step(`^el proyecto de reporte tiene al integrante "([^"]*)" con rol "([^"]*)"$`, p.proyectoConIntegrante)
	ctx.Step(`^el Sprint "([^"]*)" finalizado con el objetivo "([^"]*)"$`, p.sprintFinalizado)
	ctx.Step(`^el Sprint "([^"]*)" en estado "([^"]*)"$`, p.sprintEnEstado)
	ctx.Step(`^el Sprint "([^"]*)" en estado "([^"]*)" con el objetivo "([^"]*)"$`, p.sprintEnEstadoConObjetivo)
	ctx.Step(`^la historia "([^"]*)" está en el Sprint "([^"]*)" con (\d+) Story Points y estado "([^"]*)"$`, p.historiaEnSprint)
	ctx.Step(`^la historia "([^"]*)" se arrastró al cerrar el Sprint "([^"]*)" con (\d+) Story Points$`, p.historiaArrastrada)
	ctx.Step(`^la historia "([^"]*)" está en el Product Backlog con prioridad "([^"]*)" y (\d+) Story Points$`, p.historiaEnBacklog)
	ctx.Step(`^el Sprint "([^"]*)" tiene (\d+) horas registradas$`, p.registrarHoras)
	ctx.Step(`^el Sprint "([^"]*)" tiene el defecto "([^"]*)" con severidad "([^"]*)" detectado$`, p.defectoDetectado)
	ctx.Step(`^el Sprint "([^"]*)" tiene el defecto "([^"]*)" con severidad "([^"]*)" resuelto$`, p.defectoResuelto)
	ctx.Step(`^el defecto "([^"]*)" fue detectado en el Sprint "([^"]*)" y resuelto en el Sprint "([^"]*)"$`, p.defectoEntreSprints)

	ctx.Step(`^solicito el reporte del Sprint "([^"]*)"$`, p.solicitarSprint)
	ctx.Step(`^solicito el reporte del Sprint con identificador (\d+), que no existe$`, p.solicitarSprintInexistente)
	ctx.Step(`^solicito el reporte del proyecto$`, p.solicitarProyecto)
	ctx.Step(`^solicito el reporte del proyecto con identificador (\d+), que no existe$`, p.solicitarProyectoInexistente)
	ctx.Step(`^solicito el reporte del proyecto con un generador que falla$`, p.solicitarProyectoConGeneradorQueFalla)

	ctx.Step(`^obtengo un PDF válido con el reporte del Sprint "([^"]*)"$`, p.pdfDeSprint)
	ctx.Step(`^obtengo un PDF válido con el reporte del proyecto$`, p.pdfDeProyecto)
	ctx.Step(`^el reporte del Sprint incluye la historia completada "([^"]*)"$`, p.sprintIncluyeCompletada)
	ctx.Step(`^el reporte del Sprint incluye la historia planificada "([^"]*)" como no completada$`, p.sprintIncluyePlanificada)
	ctx.Step(`^el reporte del Sprint informa (\d+) horas reales$`, p.sprintInformaHoras)
	ctx.Step(`^el reporte del Sprint informa los defectos detectados y resueltos$`, p.sprintInformaDefectos)
	ctx.Step(`^el reporte del Sprint no informa defectos$`, p.sprintSinDefectos)
	ctx.Step(`^el reporte del Sprint lista el defecto "([^"]*)" como (detectado|resuelto)$`, p.sprintListaDefecto)
	ctx.Step(`^el reporte del Sprint "([^"]*)" informa 0 horas reales y sin defectos$`, p.sprintDelProyectoSinHorasNiDefectos)

	ctx.Step(`^la solicitud se rechaza con el mensaje "([^"]*)"$`, p.solicitudRechazada)
	ctx.Step(`^no se descarga ningún archivo$`, p.sinArchivo)

	ctx.Step(`^el reporte incluye los datos del proyecto, sus integrantes y el Product Backlog$`, p.proyectoIncluyeDatos)
	ctx.Step(`^el reporte incluye (\d+) Sprints finalizados en orden$`, p.proyectoIncluyeSprintsFinalizados)
	ctx.Step(`^el reporte incluye el Sprint Activo "([^"]*)" con sus historias$`, p.proyectoIncluyeSprintActivo)
	ctx.Step(`^el reporte no suma el Sprint Activo a las métricas$`, p.proyectoNoSumaSprintActivo)
	ctx.Step(`^el reporte no incluye la sección del Sprint Activo$`, p.proyectoSinSprintActivo)
	ctx.Step(`^el reporte no incluye historias en el Product Backlog$`, p.proyectoSinHistoriasEnBacklog)
	ctx.Step(`^el reporte del proyecto lista el defecto "([^"]*)" como detectado en el Sprint "([^"]*)" y resuelto en el Sprint "([^"]*)"$`, p.proyectoListaDefecto)
	ctx.Step(`^el reporte incluye la historia "([^"]*)" en el Product Backlog$`, p.proyectoIncluyeHistoriaEnBacklog)
	ctx.Step(`^el reporte incluye la historia "([^"]*)" solo en el Sprint Activo$`, p.proyectoIncluyeHistoriaSoloEnActivo)
	ctx.Step(`^el reporte informa las métricas, el esfuerzo y los defectos como "sin datos"$`, p.informaSinDatos)
}
