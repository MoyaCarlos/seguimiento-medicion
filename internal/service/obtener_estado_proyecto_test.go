package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

func TestObtenerEstadoProyecto_Ejecutar(t *testing.T) {
	ctx := context.Background()
	inicio := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	antes := inicio.AddDate(0, 0, -1)
	enRango := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	despues := fin.AddDate(0, 0, 1)

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
		{"sin sprints iniciados, ahora igual al inicio (borde) => en curso", true, nil, inicio, domain.ProyectoEnCurso},
		{"sin sprints iniciados, ahora dentro del rango => en curso", true, nil, enRango, domain.ProyectoEnCurso},
		{"sin sprints iniciados, ahora igual al fin (borde) => en curso", true, nil, fin, domain.ProyectoEnCurso},
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
