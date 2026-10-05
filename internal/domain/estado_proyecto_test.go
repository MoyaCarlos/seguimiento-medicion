package domain

import (
	"testing"
	"time"
)

func tiempoPtr(t time.Time) *time.Time { return &t }

// civilLocal devuelve el instante correspondiente a una fecha civil dada en la
// zona horaria local, para que los casos de borde sean deterministas.
func civilLocal(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func TestCalcularEstadoProyecto_TablaDeDecision(t *testing.T) {
	// Fechas del proyecto: se persisten como día UTC (00:00), igual que HU-04.
	inicio := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)

	// Instantes de consulta, expresados como días civiles en hora local.
	antes := civilLocal(2026, 2, 28)
	inicioHoy := civilLocal(2026, 3, 1)
	enRango := civilLocal(2026, 6, 15)
	finHoy := civilLocal(2026, 11, 30)
	despues := civilLocal(2026, 12, 1)

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
		{"sin sprints iniciados, ahora igual al inicio (borde inclusivo) => en curso", conFechas, nil, inicioHoy, ProyectoEnCurso},
		{"sin sprints iniciados, ahora dentro del rango => en curso", conFechas, nil, enRango, ProyectoEnCurso},
		{"sin sprints iniciados, ahora igual al fin (borde inclusivo) => en curso", conFechas, nil, finHoy, ProyectoEnCurso},
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

func TestCalcularEstadoProyecto_BordesConHoraDelDia(t *testing.T) {
	inicio := time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	conFechas := Project{FechaInicio: tiempoPtr(inicio), FechaFin: tiempoPtr(fin)}

	casos := []struct {
		nombre   string
		ahora    time.Time
		esperado EstadoProyecto
	}{
		{"fin 30/11 a las 10:00 => en curso", time.Date(2026, 11, 30, 10, 0, 0, 0, time.Local), ProyectoEnCurso},
		{"fin 30/11 a las 23:59:59 => en curso", time.Date(2026, 11, 30, 23, 59, 59, 0, time.Local), ProyectoEnCurso},
		{"1/12 a las 00:00:01 => finalizado", time.Date(2026, 12, 1, 0, 0, 1, 0, time.Local), ProyectoFinalizado},
		{"inicio 10/10 a las 00:00:01 => en curso", time.Date(2026, 10, 10, 0, 0, 1, 0, time.Local), ProyectoEnCurso},
		{"9/10 a las 23:59:59 => planificado", time.Date(2026, 10, 9, 23, 59, 59, 0, time.Local), ProyectoPlanificado},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got := CalcularEstadoProyecto(conFechas, nil, c.ahora)
			if got != c.esperado {
				t.Errorf("CalcularEstadoProyecto() = %q, se esperaba %q", got, c.esperado)
			}
		})
	}
}
