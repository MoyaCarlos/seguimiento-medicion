package features

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// relojEstadoHistoria fija la fecha de finalización de los escenarios de HU-11.
var relojEstadoHistoria = time.Date(2026, 10, 29, 0, 0, 0, 0, time.UTC)

// pasosEstadoHistoria guarda el estado de los escenarios de HU-11. Reutiliza la
// base y el proyecto del Background, y los services de Sprint de HU-05 y HU-06.
type pasosEstadoHistoria struct {
	base    *scenarioContext
	sprints *repository.SQLiteSprintRepository
	backlog *repository.SQLiteBacklogRepository
	crear   *service.CrearSprint
	iniciar *service.IniciarSprint
	cerrar  *service.CerrarSprint
	asignar *service.AsignarHistoriaASprint
	cambiar *service.CambiarEstadoHistoria

	sprintID  int64
	historias map[string]int64
	antes     domain.BacklogItem
	err       error
}

func (p *pasosEstadoHistoria) preparar() {
	proyectos := repository.NewSQLiteProjectRepository(p.base.db)
	p.sprints = repository.NewSQLiteSprintRepository(p.base.db)
	p.backlog = repository.NewSQLiteBacklogRepository(p.base.db)
	p.crear = service.NewCrearSprint(proyectos, p.sprints)
	p.iniciar = service.NewIniciarSprint(p.sprints)
	p.cerrar = service.NewCerrarSprint(p.sprints, p.backlog)
	p.asignar = service.NewAsignarHistoriaASprint(p.backlog, p.sprints)
	p.cambiar = service.NewCambiarEstadoHistoria(p.backlog, p.sprints, func() time.Time { return relojEstadoHistoria })
	p.sprintID = 0
	p.historias = map[string]int64{}
	p.antes, p.err = domain.BacklogItem{}, nil
}

func (p *pasosEstadoHistoria) historia(titulo string) (domain.BacklogItem, error) {
	id, ok := p.historias[titulo]
	if !ok {
		return domain.BacklogItem{}, fmt.Errorf("la historia %q no fue creada en el escenario", titulo)
	}
	return p.backlog.ObtenerPorID(context.Background(), id)
}

func (p *pasosEstadoHistoria) cambiarEstado(titulo string, estado domain.Estado) error {
	_, err := p.cambiar.Ejecutar(context.Background(), p.historias[titulo], estado)
	return err
}

// tableroConHistoria prepara un Sprint en el estado pedido con la historia ya
// asignada, usando solo los services: la historia se asigna con el Sprint
// Pendiente y se avanza una vez activo. Un Sprint Finalizado solo conserva
// historias Completadas, porque el cierre devuelve las demás al Product Backlog.
func (p *pasosEstadoHistoria) tableroConHistoria(estadoSprint, titulo string, puntos *int, estado string) error {
	ctx := context.Background()
	item, err := domain.NewBacklogItem(p.base.proyectoID, titulo, "Descripción", domain.PrioridadMust, nil)
	if err != nil {
		return err
	}
	item.EstimacionSP = puntos
	guardada, err := p.backlog.Guardar(ctx, item)
	if err != nil {
		return err
	}
	p.historias[titulo] = guardada.ID

	sprint, err := p.crear.Ejecutar(ctx, p.base.proyectoID)
	if err != nil {
		return err
	}
	p.sprintID = sprint.ID
	if _, err := p.asignar.Ejecutar(ctx, guardada.ID, sprint.ID); err != nil {
		return err
	}
	if estadoSprint == "Pendiente" {
		if estado != string(domain.EstadoNueva) {
			return fmt.Errorf("una historia de un Sprint Pendiente solo puede estar Nueva, no %q", estado)
		}
		return nil
	}
	if _, err := p.iniciar.Ejecutar(ctx, service.IniciarSprintInput{
		SprintID: sprint.ID, SprintGoal: "Goal del Sprint", FechaInicio: inicioSprint, FechaFin: finSprint,
	}); err != nil {
		return err
	}
	switch domain.Estado(estado) {
	case domain.EstadoNueva:
	case domain.EstadoEnProgreso:
		if err := p.cambiarEstado(titulo, domain.EstadoEnProgreso); err != nil {
			return err
		}
	case domain.EstadoCompletada:
		for _, e := range []domain.Estado{domain.EstadoEnProgreso, domain.EstadoCompletada} {
			if err := p.cambiarEstado(titulo, e); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("estado de partida desconocido %q", estado)
	}
	if estadoSprint != "Finalizado" {
		return nil
	}
	if estado != string(domain.EstadoCompletada) {
		return fmt.Errorf("una historia %q no permanece en un Sprint Finalizado: el cierre la devuelve al Product Backlog", estado)
	}
	_, err = p.cerrar.Ejecutar(ctx, sprint.ID)
	return err
}

func (p *pasosEstadoHistoria) tableroConHistoriaSinPuntos(estadoSprint, titulo, estado string) error {
	return p.tableroConHistoria(estadoSprint, titulo, nil, estado)
}

func (p *pasosEstadoHistoria) tableroConHistoriaConPuntos(estadoSprint, titulo string, puntos int, estado string) error {
	return p.tableroConHistoria(estadoSprint, titulo, &puntos, estado)
}

func (p *pasosEstadoHistoria) pongoEnEstado(titulo, estado string) error {
	antes, err := p.historia(titulo)
	if err != nil {
		return err
	}
	p.antes = antes
	p.err = p.cambiarEstado(titulo, domain.Estado(estado))
	return nil
}

func (p *pasosEstadoHistoria) consultoDeNuevo(titulo string) error {
	_, err := p.historia(titulo)
	return err
}

func (p *pasosEstadoHistoria) cierroElSprint() error {
	_, err := p.cerrar.Ejecutar(context.Background(), p.sprintID)
	return err
}

func (p *pasosEstadoHistoria) figuraEnEstado(titulo, estado string) error {
	h, err := p.historia(titulo)
	if err != nil {
		return err
	}
	if h.Estado != domain.Estado(estado) {
		return fmt.Errorf("se esperaba la historia %q en estado %q, está en %q", titulo, estado, h.Estado)
	}
	return nil
}

// sigueEnElSprintConSusDatos compara con la historia tal como estaba antes del
// cambio: solo pueden variar el estado y la fecha de finalización.
func (p *pasosEstadoHistoria) sigueEnElSprintConSusDatos(titulo string) error {
	h, err := p.historia(titulo)
	if err != nil {
		return err
	}
	a := p.antes
	if h.SprintID == nil || a.SprintID == nil || *h.SprintID != *a.SprintID || *h.SprintID != p.sprintID {
		return fmt.Errorf("la historia %q ya no está en el Sprint activo", titulo)
	}
	if h.Titulo != a.Titulo || h.Descripcion != a.Descripcion || h.Prioridad != a.Prioridad || h.ProyectoID != a.ProyectoID {
		return fmt.Errorf("los datos de la historia %q cambiaron: antes %+v, ahora %+v", titulo, a, h)
	}
	return nil
}

func (p *pasosEstadoHistoria) tieneFechaDeFinalizacion(titulo string) error {
	h, err := p.historia(titulo)
	if err != nil {
		return err
	}
	if h.FechaFinalizacion == nil || h.FechaFinalizacion.Format("2006-01-02") != "2026-10-29" {
		return fmt.Errorf("se esperaba la fecha de finalización 2026-10-29, es %v", h.FechaFinalizacion)
	}
	return nil
}

func (p *pasosEstadoHistoria) cuentaPuntosCompletados(esperados int) error {
	historias, err := p.backlog.ListarPorSprint(context.Background(), p.sprintID)
	if err != nil {
		return err
	}
	total := 0
	for _, h := range historias {
		if h.Estado == domain.EstadoCompletada && h.EstimacionSP != nil {
			total += *h.EstimacionSP
		}
	}
	if total != esperados {
		return fmt.Errorf("se esperaban %d Story Points completados, hay %d", esperados, total)
	}
	return nil
}

func (p *pasosEstadoHistoria) permaneceEnElSprintCerrado(titulo string) error {
	h, err := p.historia(titulo)
	if err != nil {
		return err
	}
	if h.SprintID == nil || *h.SprintID != p.sprintID {
		return fmt.Errorf("la historia %q no permanece en el Sprint cerrado", titulo)
	}
	s, err := p.sprints.ObtenerPorID(context.Background(), p.sprintID)
	if err != nil {
		return err
	}
	if s.Estado != domain.SprintFinalizado {
		return fmt.Errorf("el Sprint está %q, se esperaba Finalizado", s.Estado)
	}
	return nil
}

func (p *pasosEstadoHistoria) noApareceEnProductBacklog(titulo string) error {
	historias, err := p.backlog.ListarSinSprint(context.Background(), p.base.proyectoID)
	if err != nil {
		return err
	}
	for _, h := range historias {
		if h.ID == p.historias[titulo] {
			return fmt.Errorf("la historia %q aparece en el Product Backlog", titulo)
		}
	}
	return nil
}

func (p *pasosEstadoHistoria) rechazadoConMensaje(mensaje string) error {
	var verr domain.ValidationError
	if !errors.As(p.err, &verr) || verr.Campo != "estado" || verr.Mensaje != mensaje {
		return fmt.Errorf("se esperaba el rechazo del campo estado con %q, se obtuvo %v", mensaje, p.err)
	}
	return nil
}

// InitializeScenarioEstadoHistoria registra los pasos de HU-11 sobre el
// contexto de HU-01 (base de datos y proyecto del Background).
func InitializeScenarioEstadoHistoria(ctx *godog.ScenarioContext, base *scenarioContext) {
	p := &pasosEstadoHistoria{base: base}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		p.preparar()
		return goctx, nil
	})

	ctx.Step(`^que el tablero del proyecto tiene un Sprint (activo|Pendiente|Finalizado) con la historia "([^"]*)" en estado "([^"]*)"$`, p.tableroConHistoriaSinPuntos)
	ctx.Step(`^que el tablero del proyecto tiene un Sprint (activo|Pendiente|Finalizado) con la historia "([^"]*)" de (\d+) Story Points en estado "([^"]*)"$`, p.tableroConHistoriaConPuntos)
	ctx.Step(`^(?:que )?pongo la historia "([^"]*)" del tablero en estado "([^"]*)"$`, p.pongoEnEstado)
	ctx.Step(`^consulto de nuevo la historia "([^"]*)" del tablero$`, p.consultoDeNuevo)
	ctx.Step(`^se cierra el Sprint activo del tablero$`, p.cierroElSprint)
	ctx.Step(`^la historia "([^"]*)" del tablero figura en estado "([^"]*)"$`, p.figuraEnEstado)
	ctx.Step(`^la historia "([^"]*)" del tablero sigue en el Sprint activo con todos sus datos$`, p.sigueEnElSprintConSusDatos)
	ctx.Step(`^la historia "([^"]*)" del tablero tiene registrada su fecha de finalización$`, p.tieneFechaDeFinalizacion)
	ctx.Step(`^el Sprint activo del tablero cuenta (\d+) Story Points completados$`, p.cuentaPuntosCompletados)
	ctx.Step(`^la historia "([^"]*)" del tablero permanece en el Sprint cerrado$`, p.permaneceEnElSprintCerrado)
	ctx.Step(`^la historia "([^"]*)" del tablero no aparece en el Product Backlog del proyecto$`, p.noApareceEnProductBacklog)
	ctx.Step(`^el cambio de estado del tablero se rechaza con el mensaje "([^"]*)"$`, p.rechazadoConMensaje)
}
