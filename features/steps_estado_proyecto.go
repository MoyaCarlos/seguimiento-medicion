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

// pasosEstado guarda el estado de los escenarios de HU-13. Reutiliza la base y
// el proyecto reales del Background de HU-01 (no redefine ese paso).
type pasosEstado struct {
	base     *scenarioContext
	sprints  *repository.SQLiteSprintRepository
	backlog  *repository.SQLiteBacklogRepository
	crear    *service.CrearSprint
	iniciar  *service.IniciarSprint
	cerrar   *service.CerrarSprint
	consulta *service.ObtenerEstadoProyecto
	hoy      time.Time
	sprint   domain.Sprint
	estado   domain.EstadoProyecto
	err      error
}

func (p *pasosEstado) preparar() {
	p.hoy = hoyCivil()
	p.sprints = repository.NewSQLiteSprintRepository(p.base.db)
	p.backlog = repository.NewSQLiteBacklogRepository(p.base.db)
	proyectos := repository.NewSQLiteProjectRepository(p.base.db)
	p.crear = service.NewCrearSprint(proyectos, p.sprints)
	p.iniciar = service.NewIniciarSprint(p.sprints)
	p.cerrar = service.NewCerrarSprint(p.sprints, p.backlog)
	p.consulta = service.NewObtenerEstadoProyecto(proyectos, p.sprints)
	p.sprint, p.estado, p.err = domain.Sprint{}, "", nil
}

// hoyCivil devuelve la medianoche UTC del día civil actual en la zona horaria
// local del servidor. Se usa solo para derivar las fechas del proyecto de los
// escenarios; la consulta en sí pasa time.Now() real (con hora) para que el
// borde inclusivo [inicio, fin] se pruebe contra la hora del día.
func hoyCivil() time.Time {
	ahora := time.Now()
	return time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 0, 0, 0, 0, time.UTC)
}

func (p *pasosEstado) fijarFechas(inicio, fin *time.Time) error {
	_, err := p.base.db.Exec(
		"UPDATE projects SET start_date = ?, end_date = ? WHERE id = ?",
		fechaOpcionalRFC3339(inicio), fechaOpcionalRFC3339(fin), p.base.proyectoID,
	)
	return err
}

func fechaOpcionalRFC3339(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.RFC3339)
}

func (p *pasosEstado) fechasCargadas() error {
	inicio := p.hoy.AddDate(0, 0, -7)
	fin := p.hoy.AddDate(0, 0, 7)
	return p.fijarFechas(&inicio, &fin)
}

func (p *pasosEstado) inicioFuturo() error {
	inicio := p.hoy.AddDate(0, 0, 7)
	return p.fijarFechas(&inicio, nil)
}

func (p *pasosEstado) fechaActualEnRango() error {
	inicio := p.hoy.AddDate(0, 0, -7)
	fin := p.hoy.AddDate(0, 0, 7)
	return p.fijarFechas(&inicio, &fin)
}

// Borde inferior inclusivo: hoy == fecha_inicio ⇒ "En curso".
func (p *pasosEstado) fechaActualComoInicio() error {
	inicio := p.hoy
	fin := p.hoy.AddDate(0, 0, 14)
	return p.fijarFechas(&inicio, &fin)
}

// Borde superior inclusivo: hoy == fecha_fin ⇒ sigue "En curso".
func (p *pasosEstado) fechaActualComoFin() error {
	inicio := p.hoy.AddDate(0, 0, -14)
	fin := p.hoy
	return p.fijarFechas(&inicio, &fin)
}

func (p *pasosEstado) crearActivo() error {
	pendiente, err := p.crear.Ejecutar(context.Background(), p.base.proyectoID)
	if err != nil {
		return err
	}
	activo, err := p.iniciar.Ejecutar(context.Background(), service.IniciarSprintInput{
		SprintID:    pendiente.ID,
		SprintGoal:  "Goal del Sprint",
		FechaInicio: p.hoy,
		FechaFin:    p.hoy.AddDate(0, 0, 14),
	})
	if err != nil {
		return err
	}
	p.sprint = activo
	return nil
}

func (p *pasosEstado) cuentaConActivo() error {
	return p.crearActivo()
}

func (p *pasosEstado) todosFinalizados() error {
	if err := p.crearActivo(); err != nil {
		return err
	}
	_, err := p.cerrar.Ejecutar(context.Background(), p.sprint.ID)
	return err
}

// finalizadoYPendiente deja el proyecto con un Sprint Finalizado y otro
// Pendiente (sin ningún Activo), para el caso [Finalizado, Pendiente] de FR-016.
func (p *pasosEstado) finalizadoYPendiente() error {
	if err := p.crearActivo(); err != nil {
		return err
	}
	if _, err := p.cerrar.Ejecutar(context.Background(), p.sprint.ID); err != nil {
		return err
	}
	_, err := p.crear.Ejecutar(context.Background(), p.base.proyectoID)
	return err
}

func (p *pasosEstado) sinSprintsIniciados() error {
	delProyecto, err := p.sprints.ListarPorProyecto(context.Background(), p.base.proyectoID)
	if err != nil {
		return err
	}
	for _, s := range delProyecto {
		if s.Estado == domain.SprintActivo || s.Estado == domain.SprintFinalizado {
			return fmt.Errorf("se esperaba que el proyecto no tenga Sprints iniciados, hay uno en %q", s.Estado)
		}
	}
	return nil
}

func (p *pasosEstado) consultoEstado() error {
	p.estado, p.err = p.consulta.Ejecutar(context.Background(), p.base.proyectoID, time.Now())
	return nil
}

func (p *pasosEstado) consultoEstadoDeInexistente() error {
	p.estado, p.err = p.consulta.Ejecutar(context.Background(), 999, time.Now())
	return nil
}

func (p *pasosEstado) consulteYEstadoEra(esperado string) error {
	if err := p.consultoEstado(); err != nil {
		return err
	}
	return p.estadoEsperado(esperado)
}

func (p *pasosEstado) inicioSprintYConsulto() error {
	if err := p.crearActivo(); err != nil {
		return err
	}
	return p.consultoEstado()
}

func (p *pasosEstado) cierroUltimoYConsulto() error {
	if _, err := p.cerrar.Ejecutar(context.Background(), p.sprint.ID); err != nil {
		return err
	}
	return p.consultoEstado()
}

func (p *pasosEstado) estadoEsperado(esperado string) error {
	if p.err != nil {
		return fmt.Errorf("no se esperaba error al consultar el estado: %w", p.err)
	}
	if string(p.estado) != esperado {
		return fmt.Errorf("se esperaba estado %q, se obtuvo %q", esperado, p.estado)
	}
	return nil
}

func (p *pasosEstado) proyectoNoEncontrado() error {
	if !errors.Is(p.err, domain.ErrProyectoNoEncontrado) {
		return fmt.Errorf("se esperaba proyecto no encontrado, se obtuvo %v", p.err)
	}
	return nil
}

// InitializeScenarioEstado registra los pasos de HU-13 sobre el contexto de
// HU-01, para compartir el proyecto creado por el Antecedentes.
func InitializeScenarioEstado(ctx *godog.ScenarioContext, base *scenarioContext) {
	p := &pasosEstado{base: base}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		p.preparar()
		return goctx, nil
	})

	ctx.Step(`^que el proyecto tiene fecha de inicio y fecha de fin cargadas$`, p.fechasCargadas)
	ctx.Step(`^que el proyecto tiene una fecha de inicio futura$`, p.inicioFuturo)
	ctx.Step(`^que el proyecto tiene la fecha actual dentro de su rango de fechas$`, p.fechaActualEnRango)
	ctx.Step(`^que el proyecto tiene la fecha actual como fecha de inicio$`, p.fechaActualComoInicio)
	ctx.Step(`^que el proyecto tiene la fecha actual como fecha de fin$`, p.fechaActualComoFin)
	ctx.Step(`^que el proyecto cuenta con un Sprint en estado "Activo"$`, p.cuentaConActivo)
	ctx.Step(`^que el proyecto tiene todos sus Sprints en estado "Finalizado"$`, p.todosFinalizados)
	ctx.Step(`^que el proyecto tiene un Sprint en estado "Finalizado" y uno "Pendiente"$`, p.finalizadoYPendiente)
	ctx.Step(`^que el proyecto no tiene Sprints iniciados$`, p.sinSprintsIniciados)
	ctx.Step(`^que consulté el estado del proyecto y era "([^"]*)"$`, p.consulteYEstadoEra)
	ctx.Step(`^consulto el estado del proyecto$`, p.consultoEstado)
	ctx.Step(`^consulto el estado de un proyecto que no existe$`, p.consultoEstadoDeInexistente)
	ctx.Step(`^inicio un Sprint del proyecto y consulto su estado nuevamente$`, p.inicioSprintYConsulto)
	ctx.Step(`^cierro su último Sprint y consulto su estado nuevamente$`, p.cierroUltimoYConsulto)
	ctx.Step(`^el estado del proyecto es "([^"]*)"$`, p.estadoEsperado)
	ctx.Step(`^el sistema responde "proyecto no encontrado"$`, p.proyectoNoEncontrado)
}
