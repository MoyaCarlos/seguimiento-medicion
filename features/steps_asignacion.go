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
	base      *scenarioContext
	sprints   *repository.SQLiteSprintRepository
	backlog   *repository.SQLiteBacklogRepository
	crear     *service.CrearSprint
	iniciar   *service.IniciarSprint
	cerrar    *service.CerrarSprint
	asignar   *service.AsignarHistoriaASprint
	consultar *service.ConsultarProductBacklog

	estadoPlan string // estado pedido para el Sprint a planificar
	planID     int64  // 0 hasta que sprintAPlanificar lo crea
	otroSprint int64  // Sprint del otro proyecto
	historias  map[string]int64
	antes      domain.BacklogItem
	err        error
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
	p.estadoPlan, p.planID, p.otroSprint = "", 0, 0
	p.historias = map[string]int64{}
	p.antes, p.err = domain.BacklogItem{}, nil
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
	s, err := p.crear.Ejecutar(context.Background(), p.base.proyectoID)
	if err != nil {
		return 0, err
	}
	if p.estadoPlan != string(domain.SprintPendiente) {
		if err := p.iniciarSprint(s.ID); err != nil {
			return 0, err
		}
	}
	if p.estadoPlan == string(domain.SprintFinalizado) {
		if _, err := p.cerrar.Ejecutar(context.Background(), s.ID); err != nil {
			return 0, err
		}
	}
	p.planID = s.ID
	return s.ID, nil
}

func (p *pasosAsignacion) tieneSprintAPlanificar(estado string) error {
	p.estadoPlan = estado
	return nil
}

func (p *pasosAsignacion) nuevaHistoria(titulo string, puntos *int) (int64, error) {
	item, err := domain.NewBacklogItem(p.base.proyectoID, titulo, "Descripción", domain.PrioridadMust, nil)
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
}
