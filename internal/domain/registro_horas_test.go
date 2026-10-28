package domain

import (
	"testing"
	"time"
)

func dia(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// Sprint Activo del 2026-10-05 al 2026-10-19 con una historia asignada; "hoy"
// es el 2026-10-10 (hora local).
func escenarioDeCarga() (BacklogItem, Sprint, time.Time) {
	sprintID := int64(7)
	historia := BacklogItem{ID: 3, ProyectoID: 1, Estado: EstadoNueva, SprintID: &sprintID}
	sprint := Sprint{ID: sprintID, ProyectoID: 1, SprintGoal: "g", FechaInicio: dia(2026, 10, 5), FechaFin: dia(2026, 10, 19), Estado: SprintActivo}
	return historia, sprint, civilLocal(2026, 10, 10)
}

func TestMinutosDesdeHoras_Validas(t *testing.T) {
	casos := []struct {
		horas    string
		esperado int
	}{
		{"4", 240},
		{"1.5", 90},
		{"0.25", 15},
		{"0.5", 30},
		{"24", 1440},
	}
	for _, c := range casos {
		t.Run(c.horas, func(t *testing.T) {
			minutos, err := MinutosDesdeHoras(c.horas)
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if minutos != c.esperado {
				t.Errorf("se esperaban %d minutos, se obtuvo %d", c.esperado, minutos)
			}
		})
	}
}

func TestNewRegistroHoras_Valido(t *testing.T) {
	historia, sprint, ahora := escenarioDeCarga()

	r, err := NewRegistroHoras(historia, sprint, 2, "4", "  Implementar  ", "2026-10-09", ahora)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if r.ProyectoID != 1 || r.HistoriaID != 3 || r.SprintID != 7 || r.UsuarioID != 2 {
		t.Errorf("ids: se obtuvo proyecto %d, historia %d, sprint %d, usuario %d", r.ProyectoID, r.HistoriaID, r.SprintID, r.UsuarioID)
	}
	if r.Minutos != 240 {
		t.Errorf("se esperaban 240 minutos, se obtuvo %d", r.Minutos)
	}
	if r.Actividad != "Implementar" {
		t.Errorf("se esperaba la actividad sin espacios externos, se obtuvo %q", r.Actividad)
	}
	if !r.Fecha.Equal(dia(2026, 10, 9)) || r.Fecha.Location() != time.UTC {
		t.Errorf("se esperaba la fecha 2026-10-09 a las 00:00 UTC, se obtuvo %v", r.Fecha)
	}
}

func TestNewRegistroHoras_AceptaLosBordesDelRango(t *testing.T) {
	historia, sprint, ahora := escenarioDeCarga()
	for _, fecha := range []string{"2026-10-10", "2026-10-05"} {
		t.Run(fecha, func(t *testing.T) {
			r, err := NewRegistroHoras(historia, sprint, 2, "1", "Implementar", fecha, ahora)
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if r.Fecha.Format(time.DateOnly) != fecha {
				t.Errorf("se esperaba la fecha %s, se obtuvo %v", fecha, r.Fecha)
			}
		})
	}
}
