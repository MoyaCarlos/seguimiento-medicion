package features

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// pasosAsignacion guarda el estado de los escenarios de HU-06. Reutiliza la
// base y el proyecto del Background, y los services de Sprint de HU-05.
type pasosAsignacion struct {
	base             *scenarioContext
	sprints          *repository.SQLiteSprintRepository
	backlog          *repository.SQLiteBacklogRepository
	crear            *service.CrearSprint
	iniciar          *service.IniciarSprint
	cerrar           *service.CerrarSprint
	asignar          *service.AsignarHistoriaASprint
	consultar        *service.ConsultarProductBacklog
	consultarBacklog *service.ConsultarSprintBacklog

	estadoPlan string // estado pedido para el Sprint a planificar
	planID     int64  // 0 hasta que sprintAPlanificar lo crea
	otroSprint int64  // Sprint del otro proyecto
	historias  map[string]int64
	antes      domain.BacklogItem
	err        error

	sprintsPorEstado map[string]int64
	sprintBacklog    []domain.BacklogItem
	errConsulta      error
}

func (p *pasosAsignacion) preparar() {
	proyectos := repository.NewSQLiteProjectRepository(p.base.db)
	p.sprints = repository.NewSQLiteSprintRepository(p.base.db)
	p.backlog = repository.NewSQLiteBacklogRepository(p.base.db)
	p.crear = service.NewCrearSprint(proyectos, p.sprints)
	p.iniciar = service.NewIniciarSprint(p.sprints)
	p.cerrar = service.NewCerrarSprint(p.sprints, p.backlog)
	p.asignar = service.NewAsignarHistoriaASprint(p.backlog, p.sprints)
	p.consultar = service.NewConsultarProductBacklog(proyectos, p.backlog)
	p.consultarBacklog = service.NewConsultarSprintBacklog(p.sprints, p.backlog)
	p.estadoPlan, p.planID, p.otroSprint = "", 0, 0
	p.historias = map[string]int64{}
	p.antes, p.err = domain.BacklogItem{}, nil
	p.sprintsPorEstado = map[string]int64{}
	p.sprintBacklog, p.errConsulta = nil, nil
}

// iniciarSprint lleva un Sprint Pendiente a Activo con el Goal y las fechas de HU-05.
func (p *pasosAsignacion) iniciarSprint(id int64) error {
	_, err := p.iniciar.Ejecutar(context.Background(), service.IniciarSprintInput{
		SprintID: id, SprintGoal: "Goal del Sprint", FechaInicio: inicioSprint, FechaFin: finSprint,
	})
	return err
}

// sprintAPlanificar crea el Sprint a planificar la primera vez que un paso lo
// necesita, llevándolo al estado pedido con los services de HU-05.
func (p *pasosAsignacion) sprintAPlanificar() (int64, error) {
	if p.planID != 0 {
		return p.planID, nil
	}
	id, err := p.crearSprintEnEstado(p.estadoPlan)
	if err != nil {
		return 0, err
	}
	p.planID = id
	return id, nil
}

// crearSprintEnEstado crea un Sprint y lo lleva al estado pedido con los
// services de HU-05: nunca se siembra el estado del Sprint a mano.
func (p *pasosAsignacion) crearSprintEnEstado(estado string) (int64, error) {
	s, err := p.crear.Ejecutar(context.Background(), p.base.proyectoID)
	if err != nil {
		return 0, err
	}
	switch estado {
	case string(domain.SprintPendiente):
		return s.ID, nil
	case string(domain.SprintActivo):
		return s.ID, p.iniciarSprint(s.ID)
	case string(domain.SprintFinalizado):
		if err := p.iniciarSprint(s.ID); err != nil {
			return 0, err
		}
		_, err := p.cerrar.Ejecutar(context.Background(), s.ID)
		return s.ID, err
	default:
		return 0, fmt.Errorf("estado de Sprint no soportado: %q", estado)
	}
}

func (p *pasosAsignacion) tieneSprintAPlanificar(estado string) error {
	p.estadoPlan = estado
	return nil
}

func (p *pasosAsignacion) nuevaHistoria(titulo string, puntos *int) (int64, error) {
	return p.nuevaHistoriaCon(titulo, domain.PrioridadMust, puntos)
}

func (p *pasosAsignacion) nuevaHistoriaCon(titulo string, prioridad domain.Prioridad, puntos *int) (int64, error) {
	item, err := domain.NewBacklogItem(p.base.proyectoID, titulo, "Descripción", prioridad, nil)
	if err != nil {
		return 0, err
	}
	item.EstimacionSP = puntos
	g, err := p.backlog.Guardar(context.Background(), item)
	p.historias[titulo] = g.ID
	return g.ID, err
}

func (p *pasosAsignacion) existeHistoria(titulo string) error {
	_, err := p.nuevaHistoria(titulo, nil)
	return err
}

func (p *pasosAsignacion) existeHistoriaConPuntos(titulo string, puntos int) error {
	_, err := p.nuevaHistoria(titulo, &puntos)
	return err
}

// El otro Sprint se crea antes que el Sprint a planificar: los índices únicos
// de 001_esquema_base.sql impiden dos Sprints Pendiente (o dos Activo) en un
// proyecto. La historia se siembra con SQL porque HU-11 todavía no existe.
func (p *pasosAsignacion) existeCompletadaEnOtroSprint(titulo string) error {
	if p.planID != 0 {
		return errors.New("el otro Sprint debe crearse antes que el Sprint a planificar")
	}
	id, err := p.nuevaHistoria(titulo, nil)
	if err != nil {
		return err
	}
	otro, err := p.crear.Ejecutar(context.Background(), p.base.proyectoID)
	if err != nil {
		return err
	}
	if err := p.iniciarSprint(otro.ID); err != nil {
		return err
	}
	_, err = p.base.db.Exec(
		"UPDATE historias SET sprint_id = ?, estado = 'Completada', fecha_finalizacion = '2026-10-20' WHERE id = ?", otro.ID, id)
	return err
}

func (p *pasosAsignacion) yaAsignadaAlSprint(titulo string) error {
	sprintID, err := p.sprintAPlanificar()
	if err != nil {
		return err
	}
	_, err = p.base.db.Exec("UPDATE historias SET sprint_id = ? WHERE id = ?", sprintID, p.historias[titulo])
	return err
}

func (p *pasosAsignacion) existeOtroProyectoConSprint() error {
	proyecto, err := domain.NewProject("Otro proyecto", "", nil, nil)
	if err != nil {
		return err
	}
	proyecto, err = repository.NewSQLiteProjectRepository(p.base.db).Guardar(context.Background(), proyecto)
	if err != nil {
		return err
	}
	s, err := p.crear.Ejecutar(context.Background(), proyecto.ID)
	p.otroSprint = s.ID
	return err
}

func (p *pasosAsignacion) intentarAsignar(historiaID, sprintID int64) {
	p.antes, _ = p.backlog.ObtenerPorID(context.Background(), historiaID)
	_, p.err = p.asignar.Ejecutar(context.Background(), historiaID, sprintID)
}

func (p *pasosAsignacion) asignoAlSprintAPlanificar(titulo string) error {
	sprintID, err := p.sprintAPlanificar()
	if err != nil {
		return err
	}
	p.intentarAsignar(p.historias[titulo], sprintID)
	return nil
}

func (p *pasosAsignacion) asignoAlSprintDelOtroProyecto(titulo string) error {
	p.intentarAsignar(p.historias[titulo], p.otroSprint)
	return nil
}

func (p *pasosAsignacion) asignoASprintInexistente(titulo string, sprintID int) error {
	p.intentarAsignar(p.historias[titulo], int64(sprintID))
	return nil
}

func (p *pasosAsignacion) asignoHistoriaYSprintInexistentes(historiaID, sprintID int) error {
	p.intentarAsignar(int64(historiaID), int64(sprintID))
	return nil
}

func (p *pasosAsignacion) sinErrorDeAsignacion() error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error al asignar: %w", p.err)
	}
	return nil
}

func (p *pasosAsignacion) yaNoEstaEnProductBacklog(titulo string) error {
	if err := p.sinErrorDeAsignacion(); err != nil {
		return err
	}
	historias, err := p.consultar.Ejecutar(context.Background(), p.base.proyectoID)
	if err != nil {
		return err
	}
	for _, h := range historias {
		if h.ID == p.historias[titulo] {
			return fmt.Errorf("la historia %q sigue en el Product Backlog", titulo)
		}
	}
	return nil
}

func (p *pasosAsignacion) formaParteDelSprintBacklog(titulo string) error {
	historias, err := p.backlog.ListarPorSprint(context.Background(), p.planID)
	if err != nil {
		return err
	}
	for _, h := range historias {
		if h.ID == p.historias[titulo] {
			return nil
		}
	}
	return fmt.Errorf("la historia %q no está en el Sprint Backlog", titulo)
}

func (p *pasosAsignacion) conservaEstadoYPuntos(titulo string, puntos int) error {
	if err := p.sinErrorDeAsignacion(); err != nil {
		return err
	}
	h, err := p.backlog.ObtenerPorID(context.Background(), p.historias[titulo])
	if err != nil {
		return err
	}
	if h.Estado != domain.EstadoNueva || h.EstimacionSP == nil || *h.EstimacionSP != puntos {
		return fmt.Errorf("se esperaba estado Nueva y %d Story Points, quedó %q y %v", puntos, h.Estado, h.EstimacionSP)
	}
	return nil
}

func (p *pasosAsignacion) rechazadaConMensaje(mensaje string) error {
	var verr domain.ValidationError
	if !errors.As(p.err, &verr) || verr.Mensaje != mensaje {
		return fmt.Errorf("se esperaba el rechazo %q, se obtuvo %v", mensaje, p.err)
	}
	return nil
}

func (p *pasosAsignacion) rechazadaConError(mensaje string) error {
	if p.err == nil || !strings.Contains(p.err.Error(), mensaje) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo %v", mensaje, p.err)
	}
	return nil
}

func (p *pasosAsignacion) quedaSinCambios(titulo string) error {
	h, err := p.backlog.ObtenerPorID(context.Background(), p.historias[titulo])
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(h, p.antes) {
		return fmt.Errorf("la historia %q cambió: antes %+v, ahora %+v", titulo, p.antes, h)
	}
	return nil
}

func (p *pasosAsignacion) sprintConHistoria(estado, titulo string) error {
	sprintID, err := p.crearSprintEnEstado(estado)
	if err != nil {
		return err
	}
	p.sprintsPorEstado[estado] = sprintID
	historiaID, err := p.nuevaHistoriaCon(titulo, domain.PrioridadMust, nil)
	if err != nil {
		return err
	}
	if _, err := p.asignar.Ejecutar(context.Background(), historiaID, sprintID); err != nil {
		return fmt.Errorf("no se pudo asignar %q al Sprint %q: %w", titulo, estado, err)
	}
	return nil
}

func (p *pasosAsignacion) sprintConHistorias(estado string, tabla *godog.Table) error {
	p.estadoPlan = estado
	sprintID, err := p.sprintAPlanificar()
	if err != nil {
		return err
	}
	p.sprintsPorEstado[estado] = sprintID
	for _, fila := range tabla.Rows[1:] {
		titulo, prioridad := fila.Cells[0].Value, domain.Prioridad(fila.Cells[1].Value)
		historiaID, err := p.nuevaHistoriaCon(titulo, prioridad, nil)
		if err != nil {
			return err
		}
		if _, err := p.asignar.Ejecutar(context.Background(), historiaID, sprintID); err != nil {
			return fmt.Errorf("no se pudo asignar %q al Sprint %q: %w", titulo, estado, err)
		}
	}
	return nil
}

func (p *pasosAsignacion) sprintSinHistorias(estado string) error {
	p.estadoPlan = estado
	sprintID, err := p.sprintAPlanificar()
	if err != nil {
		return err
	}
	p.sprintsPorEstado[estado] = sprintID
	return nil
}

func (p *pasosAsignacion) consultarBacklogDelSprint(estado string) error {
	sprintID, ok := p.sprintsPorEstado[estado]
	if !ok {
		return fmt.Errorf("no existe un Sprint en estado %q", estado)
	}
	p.sprintBacklog, p.errConsulta = p.consultarBacklog.Ejecutar(context.Background(), sprintID)
	return nil
}

func (p *pasosAsignacion) consultarBacklogDeSprintInexistente(id int) error {
	p.sprintBacklog, p.errConsulta = p.consultarBacklog.Ejecutar(context.Background(), int64(id))
	return nil
}

func (p *pasosAsignacion) veoSoloLasHistoriasDelSprint(estado string, tabla *godog.Table) error {
	if p.errConsulta != nil {
		return fmt.Errorf("no se esperaba error al consultar el Sprint Backlog: %w", p.errConsulta)
	}
	esperadas := tabla.Rows[1:]
	if len(p.sprintBacklog) != len(esperadas) {
		return fmt.Errorf("se esperaban %d historias, se obtuvieron %d: %+v", len(esperadas), len(p.sprintBacklog), p.sprintBacklog)
	}
	for i, fila := range esperadas {
		titulo, prioridad := fila.Cells[0].Value, domain.Prioridad(fila.Cells[1].Value)
		h := p.sprintBacklog[i]
		if h.Titulo != titulo || h.Prioridad != prioridad {
			return fmt.Errorf("en la posición %d se esperaba %q (%s), se obtuvo %q (%s)",
				i+1, titulo, prioridad, h.Titulo, h.Prioridad)
		}
		if h.SprintID == nil || *h.SprintID != p.sprintsPorEstado[estado] {
			return fmt.Errorf("la historia %q no pertenece al Sprint %q", titulo, estado)
		}
	}
	return nil
}

func (p *pasosAsignacion) veoSprintBacklogVacio() error {
	if p.errConsulta != nil {
		return fmt.Errorf("no se esperaba error al consultar el Sprint Backlog: %w", p.errConsulta)
	}
	if len(p.sprintBacklog) != 0 {
		return fmt.Errorf("se esperaba un Sprint Backlog vacío, se obtuvo %+v", p.sprintBacklog)
	}
	return nil
}

func (p *pasosAsignacion) consultaDelSprintBacklogRechazadaConError(mensaje string) error {
	if p.errConsulta == nil || !strings.Contains(p.errConsulta.Error(), mensaje) {
		return fmt.Errorf("se esperaba el error %q, se obtuvo %v", mensaje, p.errConsulta)
	}
	return nil
}

// InitializeScenarioAsignacion registra los pasos de HU-06 sobre el contexto de
// HU-01 (base de datos y proyecto del Background).
func InitializeScenarioAsignacion(ctx *godog.ScenarioContext, base *scenarioContext) {
	p := &pasosAsignacion{base: base}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		p.preparar()
		return goctx, nil
	})

	ctx.Step(`^que el proyecto tiene un Sprint a planificar en estado "([^"]*)"$`, p.tieneSprintAPlanificar)
	ctx.Step(`^que existe la historia "([^"]*)" en el Product Backlog del proyecto$`, p.existeHistoria)
	ctx.Step(`^que existe la historia "([^"]*)" con (\d+) Story Points en el Product Backlog del proyecto$`, p.existeHistoriaConPuntos)
	ctx.Step(`^que existe la historia "([^"]*)" en estado "Completada" en otro Sprint del proyecto$`, p.existeCompletadaEnOtroSprint)
	ctx.Step(`^que la historia "([^"]*)" ya está asignada al Sprint a planificar$`, p.yaAsignadaAlSprint)
	ctx.Step(`^que existe otro proyecto con un Sprint en estado "Pendiente"$`, p.existeOtroProyectoConSprint)
	ctx.Step(`^asigno la historia "([^"]*)" al Sprint a planificar$`, p.asignoAlSprintAPlanificar)
	ctx.Step(`^asigno la historia "([^"]*)" al Sprint del otro proyecto$`, p.asignoAlSprintDelOtroProyecto)
	ctx.Step(`^asigno la historia "([^"]*)" al Sprint con identificador (\d+), que no existe$`, p.asignoASprintInexistente)
	ctx.Step(`^asigno la historia con identificador (\d+), que no existe, al Sprint con identificador (\d+), que tampoco existe$`, p.asignoHistoriaYSprintInexistentes)
	ctx.Step(`^la historia "([^"]*)" ya no aparece en el Product Backlog del proyecto$`, p.yaNoEstaEnProductBacklog)
	ctx.Step(`^la historia "([^"]*)" forma parte del Sprint Backlog del Sprint a planificar$`, p.formaParteDelSprintBacklog)
	ctx.Step(`^la historia "([^"]*)" conserva el estado "Nueva" y sus (\d+) Story Points$`, p.conservaEstadoYPuntos)
	ctx.Step(`^la asignación se rechaza con el mensaje "([^"]*)"$`, p.rechazadaConMensaje)
	ctx.Step(`^la asignación se rechaza con error "([^"]*)"$`, p.rechazadaConError)
	ctx.Step(`^la historia "([^"]*)" queda sin cambios$`, p.quedaSinCambios)

	ctx.Step(`^que el proyecto tiene un Sprint "([^"]*)" con la historia "([^"]*)" en su Sprint Backlog$`, p.sprintConHistoria)
	ctx.Step(`^que el proyecto tiene un Sprint "([^"]*)" con estas historias, en este orden de creación:$`, p.sprintConHistorias)
	ctx.Step(`^que el proyecto tiene un Sprint "([^"]*)" sin historias en su Sprint Backlog$`, p.sprintSinHistorias)
	ctx.Step(`^consulto el Sprint Backlog del Sprint "([^"]*)"$`, p.consultarBacklogDelSprint)
	ctx.Step(`^consulto el Sprint Backlog del Sprint con identificador (\d+), que no existe$`, p.consultarBacklogDeSprintInexistente)
	ctx.Step(`^veo solo las historias del Sprint "([^"]*)", en este orden:$`, p.veoSoloLasHistoriasDelSprint)
	ctx.Step(`^veo un Sprint Backlog vacío, sin error$`, p.veoSprintBacklogVacio)
	ctx.Step(`^la consulta del Sprint Backlog se rechaza con error "([^"]*)"$`, p.consultaDelSprintBacklogRechazadaConError)
}
