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

// pasosPlazoSprint guarda el estado de los escenarios de #168. Reutiliza la
// base y el proyecto (sin fechas) del Background de HU-01.
type pasosPlazoSprint struct {
	base      *scenarioContext
	proyectos *repository.SQLiteProjectRepository
	sprints   *repository.SQLiteSprintRepository
	iniciar   *service.IniciarSprint
	sprint    domain.Sprint
	err       error
}

func (p *pasosPlazoSprint) preparar() {
	p.proyectos = repository.NewSQLiteProjectRepository(p.base.db)
	p.sprints = repository.NewSQLiteSprintRepository(p.base.db)
	p.iniciar = service.NewIniciarSprint(p.sprints, p.proyectos)
	p.sprint, p.err = domain.Sprint{}, nil
}

func (p *pasosPlazoSprint) sprintPendiente() error {
	s, err := domain.NewSprint(p.base.proyectoID)
	if err != nil {
		return err
	}
	p.sprint, err = p.sprints.Guardar(context.Background(), s)
	return err
}

func (p *pasosPlazoSprint) fijarPlazo(inicio, fin *time.Time) error {
	ctx := context.Background()
	proyecto, err := p.proyectos.ObtenerPorID(ctx, p.base.proyectoID)
	if err != nil {
		return err
	}
	editado, err := proyecto.ConDatosEditados(proyecto.Nombre, proyecto.Descripcion, inicio, fin)
	if err != nil {
		return err
	}
	return p.proyectos.Actualizar(ctx, editado)
}

func (p *pasosPlazoSprint) plazoEmpieza(inicio string) error {
	dia, err := time.Parse(time.DateOnly, inicio)
	if err != nil {
		return err
	}
	return p.fijarPlazo(&dia, nil)
}

func (p *pasosPlazoSprint) plazoTermina(fin string) error {
	dia, err := time.Parse(time.DateOnly, fin)
	if err != nil {
		return err
	}
	return p.fijarPlazo(nil, &dia)
}

func (p *pasosPlazoSprint) sinFechas() error {
	return p.fijarPlazo(nil, nil)
}

func (p *pasosPlazoSprint) iniciaDelAl(inicio, fin string) error {
	desde, err := time.Parse(time.DateOnly, inicio)
	if err != nil {
		return err
	}
	hasta, err := time.Parse(time.DateOnly, fin)
	if err != nil {
		return err
	}
	_, p.err = p.iniciar.Ejecutar(context.Background(), service.IniciarSprintInput{
		SprintID: p.sprint.ID, SprintGoal: "Goal del Sprint", FechaInicio: desde, FechaFin: hasta,
		UsuarioID: scrumMasterDelSprint(p.base.db, p.sprint.ID),
	})
	return nil
}

func (p *pasosPlazoSprint) rechazoEnCampo(campo string) error {
	var verr domain.ValidationError
	if !errors.As(p.err, &verr) {
		return fmt.Errorf("se esperaba un rechazo en %q, se obtuvo %v", campo, p.err)
	}
	if verr.Campo != campo {
		return fmt.Errorf("se esperaba el campo %q, se obtuvo %q (%s)", campo, verr.Campo, verr.Mensaje)
	}
	return nil
}

func (p *pasosPlazoSprint) releer() (domain.Sprint, error) {
	return p.sprints.ObtenerPorID(context.Background(), p.sprint.ID)
}

func (p *pasosPlazoSprint) siguePendienteSinDatos() error {
	s, err := p.releer()
	if err != nil {
		return err
	}
	if s.Estado != domain.SprintPendiente || s.SprintGoal != "" || !s.FechaInicio.IsZero() || !s.FechaFin.IsZero() {
		return fmt.Errorf("el Sprint debía seguir Pendiente, sin Goal ni fechas: %+v", s)
	}
	return nil
}

func (p *pasosPlazoSprint) quedaActivo() error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error al iniciar: %v", p.err)
	}
	s, err := p.releer()
	if err != nil {
		return err
	}
	if s.Estado != domain.SprintActivo {
		return fmt.Errorf("se esperaba el Sprint Activo, está %q", s.Estado)
	}
	return nil
}

// InitializeScenarioSprintDentroDelProyecto registra los pasos de #168.
func InitializeScenarioSprintDentroDelProyecto(ctx *godog.ScenarioContext, base *scenarioContext) {
	p := &pasosPlazoSprint{base: base}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		p.preparar()
		return goctx, nil
	})

	ctx.Step(`^el proyecto tiene un Sprint pendiente de iniciar$`, p.sprintPendiente)
	ctx.Step(`^que el plazo del proyecto empieza el "(\d{4}-\d{2}-\d{2})"$`, p.plazoEmpieza)
	ctx.Step(`^que el plazo del proyecto termina el "(\d{4}-\d{2}-\d{2})"$`, p.plazoTermina)
	ctx.Step(`^que el proyecto no tiene fechas de plazo$`, p.sinFechas)
	ctx.Step(`^el Scrum Master inicia el Sprint del "(\d{4}-\d{2}-\d{2})" al "(\d{4}-\d{2}-\d{2})"$`, p.iniciaDelAl)
	ctx.Step(`^el inicio se rechaza en el campo "([^"]*)"$`, p.rechazoEnCampo)
	ctx.Step(`^el Sprint sigue pendiente, sin Sprint Goal ni fechas$`, p.siguePendienteSinDatos)
	ctx.Step(`^el Sprint queda activo$`, p.quedaActivo)
}
