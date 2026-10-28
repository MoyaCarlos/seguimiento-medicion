package domain

import (
	"errors"
	"strings"
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

func sprintEntre(inicio, fin time.Time) Sprint {
	return Sprint{ID: 7, ProyectoID: 1, SprintGoal: "g", FechaInicio: inicio, FechaFin: fin, Estado: SprintActivo}
}

func TestNewRegistroHoras_Rechazos(t *testing.T) {
	historia, enCurso, ahora := escenarioDeCarga()
	vencido := sprintEntre(dia(2026, 9, 21), dia(2026, 10, 5))
	adelantado := sprintEntre(dia(2026, 10, 11), dia(2026, 10, 25))
	actividadLarga := strings.Repeat("á", 501)

	casos := []struct {
		nombre                  string
		sprint                  Sprint
		horas, actividad, fecha string
		campo, mensaje          string
	}{
		{"horas vacías", enCurso, "", "x", "2026-10-10", "horas", "las horas son obligatorias"},
		{"horas no numéricas", enCurso, "abc", "x", "2026-10-10", "horas", "las horas deben ser un número"},
		{"horas como fracción", enCurso, "1/2", "x", "2026-10-10", "horas", "las horas deben ser un número"},
		{"horas en notación exponencial", enCurso, "1e1", "x", "2026-10-10", "horas", "las horas deben ser un número"},
		{"horas negativas", enCurso, "-2", "x", "2026-10-10", "horas", "las horas deben ser mayores que cero"},
		{"horas cero", enCurso, "0", "x", "2026-10-10", "horas", "las horas deben ser mayores que cero"},
		{"horas apenas sobre el tope", enCurso, "24.01", "x", "2026-10-10", "horas", "las horas no pueden superar 24 por carga"},
		{"horas sobre el tope", enCurso, "25", "x", "2026-10-10", "horas", "las horas no pueden superar 24 por carga"},
		{"horas que no dan minutos enteros", enCurso, "0.01", "x", "2026-10-10", "horas", "las horas deben equivaler a minutos enteros"},
		{"actividad vacía", enCurso, "1", "", "2026-10-10", "actividad", "la actividad es obligatoria"},
		{"actividad solo con espacios", enCurso, "1", "   ", "2026-10-10", "actividad", "la actividad es obligatoria"},
		{"actividad de 501 caracteres", enCurso, "1", actividadLarga, "2026-10-10", "actividad", "la actividad no puede superar los 500 caracteres"},
		{"fecha vacía", enCurso, "1", "x", "", "fecha", "formato esperado AAAA-MM-DD"},
		{"fecha sin cero", enCurso, "1", "x", "2026-10-9", "fecha", "formato esperado AAAA-MM-DD"},
		{"fecha basura", enCurso, "1", "x", "basura", "fecha", "formato esperado AAAA-MM-DD"},
		{"fecha inexistente", enCurso, "1", "x", "2026-02-30", "fecha", "formato esperado AAAA-MM-DD"},
		{"fecha de mañana", enCurso, "1", "x", "2026-10-11", "fecha", "la fecha no puede ser posterior a hoy"},
		{"fecha anterior al inicio del Sprint", enCurso, "1", "x", "2026-10-04", "fecha", "la fecha no puede ser anterior al inicio del Sprint"},
		{"Sprint vencido, fecha posterior al fin", vencido, "1", "x", "2026-10-06", "fecha", "la fecha no puede ser posterior al fin del Sprint"},
		{"Sprint vencido, fecha posterior a hoy y al fin", vencido, "1", "x", "2026-10-11", "fecha", "la fecha no puede ser posterior a hoy"},
		{"Sprint adelantado, fecha de hoy", adelantado, "1", "x", "2026-10-10", "fecha", "la fecha no puede ser anterior al inicio del Sprint"},
		{"Sprint adelantado, fecha de inicio", adelantado, "1", "x", "2026-10-11", "fecha", "la fecha no puede ser posterior a hoy"},
		{"horas, actividad y fecha inválidas", enCurso, "-1", "", "basura", "horas", "las horas deben ser mayores que cero"},
		{"actividad y fecha inválidas", enCurso, "1", " ", "basura", "actividad", "la actividad es obligatoria"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := NewRegistroHoras(historia, c.sprint, 2, c.horas, c.actividad, c.fecha, ahora)

			esperado := ValidationError{Campo: c.campo, Mensaje: c.mensaje}
			var verr ValidationError
			if !errors.As(err, &verr) || verr != esperado {
				t.Errorf("se esperaba %v, se obtuvo %v", esperado, err)
			}
		})
	}
}

func TestNewRegistroHoras_AceptaLosBordes(t *testing.T) {
	historia, enCurso, ahora := escenarioDeCarga()
	vencido := sprintEntre(dia(2026, 9, 21), dia(2026, 10, 5))
	casos := []struct {
		nombre                  string
		sprint                  Sprint
		horas, actividad, fecha string
	}{
		{"actividad de 500 caracteres", enCurso, "1", strings.Repeat("á", 500), "2026-10-10"},
		{"Sprint vencido, fecha igual al fin", vencido, "1", "x", "2026-10-05"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if _, err := NewRegistroHoras(historia, c.sprint, 2, c.horas, c.actividad, c.fecha, ahora); err != nil {
				t.Errorf("no se esperaba error: %v", err)
			}
		})
	}
}

func TestValidarSprintParaImputar(t *testing.T) {
	sprintID := int64(7)
	otroSprint := int64(8)
	conSprint := func(id *int64, estado Estado) BacklogItem {
		return BacklogItem{ID: 3, ProyectoID: 1, Estado: estado, SprintID: id}
	}
	sprint := func(estado EstadoSprint) Sprint { return Sprint{ID: sprintID, ProyectoID: 1, Estado: estado} }

	casos := []struct {
		nombre   string
		historia BacklogItem
		sprint   Sprint
		esperado error
	}{
		{"historia sin Sprint", conSprint(nil, EstadoNueva), Sprint{}, ErrHistoriaFueraDelSprintActivo},
		{"Sprint Pendiente", conSprint(&sprintID, EstadoNueva), sprint(SprintPendiente), ErrHistoriaFueraDelSprintActivo},
		{"Sprint Finalizado", conSprint(&sprintID, EstadoNueva), sprint(SprintFinalizado), ErrHistoriaFueraDelSprintActivo},
		{"historia asignada a otro Sprint", conSprint(&otroSprint, EstadoNueva), sprint(SprintActivo), ErrHistoriaFueraDelSprintActivo},
		{"Sprint Activo de la historia", conSprint(&sprintID, EstadoNueva), sprint(SprintActivo), nil},
		{"historia Completada en el Sprint Activo", conSprint(&sprintID, EstadoCompletada), sprint(SprintActivo), nil},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if err := ValidarSprintParaImputar(c.historia, c.sprint); !errors.Is(err, c.esperado) {
				t.Errorf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
		})
	}
}

func TestTotalMinutos(t *testing.T) {
	carga := func(minutos int) RegistroHorasDetalle {
		return RegistroHorasDetalle{RegistroHoras: RegistroHoras{Minutos: minutos}}
	}
	casos := []struct {
		nombre    string
		registros []RegistroHorasDetalle
		esperado  int
	}{
		{"0,5 h + 1,5 h + 2 h = 4 h exactas", []RegistroHorasDetalle{carga(30), carga(90), carga(120)}, 240},
		{"sin cargas", []RegistroHorasDetalle{}, 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := TotalMinutos(c.registros); got != c.esperado {
				t.Errorf("se esperaban %d minutos, se obtuvo %d", c.esperado, got)
			}
		})
	}
}
