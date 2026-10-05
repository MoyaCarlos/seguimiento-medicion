package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// civilLocal devuelve la medianoche local del día civil dado, para que los casos
// de borde de fecha sean deterministas e independientes de la zona horaria del host.
func civilLocal(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func TestObtenerEstadoProyecto_Ejecutar(t *testing.T) {
	ctx := context.Background()
	// Fechas del proyecto: día UTC (00:00), igual que HU-04.
	inicio := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	// Instantes de consulta: días civiles en hora local.
	antes := civilLocal(2026, 2, 28)
	inicioHoy := civilLocal(2026, 3, 1)
	enRango := civilLocal(2026, 6, 15)
	finHoy := civilLocal(2026, 11, 30)
	despues := civilLocal(2026, 12, 1)

	casos := []struct {
		nombre    string
		conFechas bool
		sprints   []domain.EstadoSprint
		ahora     time.Time
		esperado  domain.EstadoProyecto
	}{
		{"sprint activo con fecha fin vencida => en curso", true, []domain.EstadoSprint{domain.SprintActivo}, despues, domain.ProyectoEnCurso},
		{"sin activos y todos finalizados => finalizado", true, []domain.EstadoSprint{domain.SprintFinalizado, domain.SprintFinalizado}, enRango, domain.ProyectoFinalizado},
		{"sin sprints iniciados, ahora antes del inicio => planificado", true, nil, antes, domain.ProyectoPlanificado},
		{"sin sprints iniciados, ahora igual al inicio (borde) => en curso", true, nil, inicioHoy, domain.ProyectoEnCurso},
		{"sin sprints iniciados, ahora dentro del rango => en curso", true, nil, enRango, domain.ProyectoEnCurso},
		{"sin sprints iniciados, ahora igual al fin (borde) => en curso", true, nil, finHoy, domain.ProyectoEnCurso},
		{"sin sprints iniciados, ahora después del fin => finalizado", true, nil, despues, domain.ProyectoFinalizado},
		{"solo sprints pendientes con inicio futuro => planificado", true, []domain.EstadoSprint{domain.SprintPendiente}, antes, domain.ProyectoPlanificado},
		{"sin sprints y sin fechas => planificado", false, nil, enRango, domain.ProyectoPlanificado},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			proyectos := repository.NewProjectRepositoryEnMemoria()
			sprints := repository.NewSprintRepositoryEnMemoria()

			proyecto := domain.Project{Nombre: "Proyecto"}
			if c.conFechas {
				proyecto.FechaInicio = &inicio
				proyecto.FechaFin = &fin
			}
			guardado, err := proyectos.Guardar(ctx, proyecto)
			if err != nil {
				t.Fatalf("no se esperaba error al guardar el proyecto: %v", err)
			}
			for _, estado := range c.sprints {
				if _, err := sprints.Guardar(ctx, domain.Sprint{ProyectoID: guardado.ID, Estado: estado}); err != nil {
					t.Fatalf("no se esperaba error al guardar el Sprint: %v", err)
				}
			}

			servicio := NewObtenerEstadoProyecto(proyectos, sprints)
			got, err := servicio.Ejecutar(ctx, guardado.ID, c.ahora)
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if got != c.esperado {
				t.Errorf("Ejecutar() = %q, se esperaba %q", got, c.esperado)
			}
		})
	}
}

func TestObtenerEstadoProyecto_Ejecutar_ProyectoInexistente(t *testing.T) {
	servicio := NewObtenerEstadoProyecto(
		repository.NewProjectRepositoryEnMemoria(),
		repository.NewSprintRepositoryEnMemoria(),
	)

	_, err := servicio.Ejecutar(context.Background(), 999, time.Now())
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
}

func TestObtenerEstadoProyecto_RecalculaOnDemand(t *testing.T) {
	ctx := context.Background()
	inicio := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	proyectos := repository.NewProjectRepositoryEnMemoria()
	sprints := repository.NewSprintRepositoryEnMemoria()
	guardado, err := proyectos.Guardar(ctx, domain.Project{Nombre: "Proyecto", FechaInicio: &inicio, FechaFin: &fin})
	if err != nil {
		t.Fatalf("no se esperaba error al guardar el proyecto: %v", err)
	}
	servicio := NewObtenerEstadoProyecto(proyectos, sprints)

	antes := inicio.AddDate(0, 0, -1)
	estadoAntes, err := servicio.Ejecutar(ctx, guardado.ID, antes)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if estadoAntes != domain.ProyectoPlanificado {
		t.Errorf("con ahora anterior al inicio se esperaba %q, se obtuvo %q", domain.ProyectoPlanificado, estadoAntes)
	}

	dentro := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	estadoDentro, err := servicio.Ejecutar(ctx, guardado.ID, dentro)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if estadoDentro != domain.ProyectoEnCurso {
		t.Errorf("con ahora dentro del rango se esperaba %q, se obtuvo %q", domain.ProyectoEnCurso, estadoDentro)
	}
	if estadoAntes == estadoDentro {
		t.Errorf("el estado debía cambiar entre consultas, se obtuvo %q en ambas", estadoDentro)
	}
}

func TestObtenerEstadoProyecto_ReflejaCambioDeSprint(t *testing.T) {
	ctx := context.Background()
	proyectos := repository.NewProjectRepositoryEnMemoria()
	sprints := repository.NewSprintRepositoryEnMemoria()
	guardado, err := proyectos.Guardar(ctx, domain.Project{Nombre: "Proyecto"})
	if err != nil {
		t.Fatalf("no se esperaba error al guardar el proyecto: %v", err)
	}
	sprint, err := sprints.Guardar(ctx, domain.Sprint{ProyectoID: guardado.ID, Estado: domain.SprintActivo})
	if err != nil {
		t.Fatalf("no se esperaba error al guardar el Sprint: %v", err)
	}
	servicio := NewObtenerEstadoProyecto(proyectos, sprints)

	estado, err := servicio.Ejecutar(ctx, guardado.ID, time.Now())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if estado != domain.ProyectoEnCurso {
		t.Fatalf("con un Sprint activo se esperaba %q, se obtuvo %q", domain.ProyectoEnCurso, estado)
	}

	sprint.Estado = domain.SprintFinalizado
	if err := sprints.Actualizar(ctx, sprint); err != nil {
		t.Fatalf("no se esperaba error al finalizar el Sprint: %v", err)
	}
	estado, err = servicio.Ejecutar(ctx, guardado.ID, time.Now())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if estado != domain.ProyectoFinalizado {
		t.Errorf("tras finalizar el último Sprint se esperaba %q, se obtuvo %q", domain.ProyectoFinalizado, estado)
	}
}

type spyProjectRepository struct {
	repository.ProjectRepository
	guardadas    int
	actualizadas int
}

func (s *spyProjectRepository) Guardar(ctx context.Context, p domain.Project) (domain.Project, error) {
	s.guardadas++
	return s.ProjectRepository.Guardar(ctx, p)
}

func (s *spyProjectRepository) Actualizar(ctx context.Context, p domain.Project) error {
	s.actualizadas++
	return s.ProjectRepository.Actualizar(ctx, p)
}

type spySprintRepository struct {
	repository.SprintRepository
	guardadas    int
	actualizadas int
}

func (s *spySprintRepository) Guardar(ctx context.Context, sp domain.Sprint) (domain.Sprint, error) {
	s.guardadas++
	return s.SprintRepository.Guardar(ctx, sp)
}

func (s *spySprintRepository) Actualizar(ctx context.Context, sp domain.Sprint) error {
	s.actualizadas++
	return s.SprintRepository.Actualizar(ctx, sp)
}

func TestObtenerEstadoProyecto_EsSoloLectura(t *testing.T) {
	ctx := context.Background()
	proyectosBase := repository.NewProjectRepositoryEnMemoria()
	sprintsBase := repository.NewSprintRepositoryEnMemoria()
	guardado, err := proyectosBase.Guardar(ctx, domain.Project{Nombre: "Proyecto"})
	if err != nil {
		t.Fatalf("no se esperaba error al guardar el proyecto: %v", err)
	}
	if _, err := sprintsBase.Guardar(ctx, domain.Sprint{ProyectoID: guardado.ID, Estado: domain.SprintActivo}); err != nil {
		t.Fatalf("no se esperaba error al guardar el Sprint: %v", err)
	}

	proyectos := &spyProjectRepository{ProjectRepository: proyectosBase}
	sprints := &spySprintRepository{SprintRepository: sprintsBase}
	servicio := NewObtenerEstadoProyecto(proyectos, sprints)

	if _, err := servicio.Ejecutar(ctx, guardado.ID, time.Now()); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if _, err := servicio.Ejecutar(ctx, 999, time.Now()); !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrProyectoNoEncontrado, se obtuvo %v", err)
	}

	if proyectos.guardadas != 0 || proyectos.actualizadas != 0 {
		t.Errorf("la consulta no debe escribir proyectos: guardadas=%d actualizadas=%d", proyectos.guardadas, proyectos.actualizadas)
	}
	if sprints.guardadas != 0 || sprints.actualizadas != 0 {
		t.Errorf("la consulta no debe escribir Sprints: guardadas=%d actualizadas=%d", sprints.guardadas, sprints.actualizadas)
	}
}
