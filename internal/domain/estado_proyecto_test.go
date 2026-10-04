package domain

import (
	"testing"
	"time"
)

func tiempoPtr(t time.Time) *time.Time { return &t }

func TestCalcularEstadoProyecto_TablaDeDecision(t *testing.T) {
	inicio := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	antes := inicio.AddDate(0, 0, -1)
	enRango := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	despues := fin.AddDate(0, 0, 1)

	conFechas := Project{FechaInicio: tiempoPtr(inicio), FechaFin: tiempoPtr(fin)}
	sprint := func(e EstadoSprint) Sprint { return Sprint{Estado: e} }

	casos := []struct {
		nombre   string
		proyecto Project
		sprints  []Sprint
		ahora    time.Time
		esperado EstadoProyecto
	}{
		{"sprint activo => en curso aunque la fecha de fin esté vencida", conFechas, []Sprint{sprint(SprintActivo)}, despues, ProyectoEnCurso},
		{"sprint activo dentro del rango => en curso", conFechas, []Sprint{sprint(SprintActivo)}, enRango, ProyectoEnCurso},
		{"sin activos y todos finalizados => finalizado", conFechas, []Sprint{sprint(SprintFinalizado), sprint(SprintFinalizado)}, enRango, ProyectoFinalizado},
		{"todos finalizados con fecha de fin vencida => finalizado", conFechas, []Sprint{sprint(SprintFinalizado)}, despues, ProyectoFinalizado},
		{"sin sprints iniciados, ahora antes del inicio => planificado", conFechas, nil, antes, ProyectoPlanificado},
		{"sin sprints iniciados, ahora igual al inicio (borde inclusivo) => en curso", conFechas, nil, inicio, ProyectoEnCurso},
		{"sin sprints iniciados, ahora dentro del rango => en curso", conFechas, nil, enRango, ProyectoEnCurso},
		{"sin sprints iniciados, ahora igual al fin (borde inclusivo) => en curso", conFechas, nil, fin, ProyectoEnCurso},
		{"sin sprints iniciados, ahora después del fin => finalizado", conFechas, nil, despues, ProyectoFinalizado},
		{"solo sprints pendientes con inicio futuro => planificado", conFechas, []Sprint{sprint(SprintPendiente)}, antes, ProyectoPlanificado},
		{"sin sprints y sin fecha de inicio => planificado", Project{}, nil, enRango, ProyectoPlanificado},
		{"fecha de fin sin fecha de inicio => planificado (no iniciado por fecha)", Project{FechaFin: tiempoPtr(fin)}, nil, enRango, ProyectoPlanificado},
		{"solo fecha de inicio, ahora dentro del rango abierto => en curso", Project{FechaInicio: tiempoPtr(inicio)}, nil, enRango, ProyectoEnCurso},
		{"solo fecha de inicio, ahora posterior al inicio => en curso", Project{FechaInicio: tiempoPtr(inicio)}, nil, despues, ProyectoEnCurso},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got := CalcularEstadoProyecto(c.proyecto, c.sprints, c.ahora)
			if got != c.esperado {
				t.Errorf("CalcularEstadoProyecto() = %q, se esperaba %q", got, c.esperado)
			}
		})
	}
}
