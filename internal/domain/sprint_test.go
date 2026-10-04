package domain

import (
	"errors"
	"testing"
	"time"
)

var (
	hoy          = time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	enDosSemanas = hoy.AddDate(0, 0, 14)
)

func esperarErrorDeValidacion(t *testing.T, err error, campo string) {
	t.Helper()
	var verr ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("se esperaba ValidationError en %q, se obtuvo %v", campo, err)
	}
	if verr.Campo != campo {
		t.Errorf("se esperaba campo %q, se obtuvo %q", campo, verr.Campo)
	}
}

func sprintActivo(t *testing.T) Sprint {
	t.Helper()
	s, err := NewSprint(1)
	if err != nil {
		t.Fatalf("no se esperaba error al crear: %v", err)
	}
	if err := s.Iniciar("Entregar el MVP", hoy, enDosSemanas); err != nil {
		t.Fatalf("no se esperaba error al iniciar: %v", err)
	}
	return s
}

func TestNewSprint_QuedaPendiente(t *testing.T) {
	s, err := NewSprint(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if s.Estado != SprintPendiente {
		t.Errorf("se esperaba estado %q, se obtuvo %q", SprintPendiente, s.Estado)
	}
	if s.ProyectoID != 1 {
		t.Errorf("se esperaba proyecto 1, se obtuvo %d", s.ProyectoID)
	}
	if s.SprintGoal != "" || !s.FechaInicio.IsZero() || !s.FechaFin.IsZero() {
		t.Errorf("un Sprint recién creado no debe tener Goal ni fechas: %+v", s)
	}
}

func TestNewSprint_ProyectoObligatorio(t *testing.T) {
	for _, id := range []int64{0, -1} {
		_, err := NewSprint(id)
		esperarErrorDeValidacion(t, err, "proyecto_id")
	}
}

func TestSprint_Iniciar_PasaAActivo(t *testing.T) {
	s := sprintActivo(t)
	if s.Estado != SprintActivo {
		t.Errorf("se esperaba estado %q, se obtuvo %q", SprintActivo, s.Estado)
	}
	if s.SprintGoal != "Entregar el MVP" || !s.FechaInicio.Equal(hoy) || !s.FechaFin.Equal(enDosSemanas) {
		t.Errorf("no se guardaron el Goal y las fechas: %+v", s)
	}
}

func TestSprint_Iniciar_GoalObligatorio(t *testing.T) {
	s, _ := NewSprint(1)
	err := s.Iniciar("   ", hoy, enDosSemanas)
	esperarErrorDeValidacion(t, err, "sprint_goal")
	if s.Estado != SprintPendiente {
		t.Errorf("el Sprint no debía cambiar de estado, quedó %q", s.Estado)
	}
}

func TestSprint_Iniciar_FechaInicioObligatoria(t *testing.T) {
	s, _ := NewSprint(1)
	err := s.Iniciar("Goal", time.Time{}, enDosSemanas)
	esperarErrorDeValidacion(t, err, "fecha_inicio")
}

func TestSprint_Iniciar_FechaFinDebeSerPosterior(t *testing.T) {
	for nombre, fin := range map[string]time.Time{
		"igual":    hoy,
		"anterior": hoy.AddDate(0, 0, -1),
	} {
		t.Run(nombre, func(t *testing.T) {
			s, _ := NewSprint(1)
			err := s.Iniciar("Goal", hoy, fin)
			esperarErrorDeValidacion(t, err, "fecha_fin")
			if s.Estado != SprintPendiente {
				t.Errorf("el Sprint no debía cambiar de estado, quedó %q", s.Estado)
			}
		})
	}
}

func TestSprint_Iniciar_SoloDesdePendiente(t *testing.T) {
	s := sprintActivo(t)
	err := s.Iniciar("Otro goal", hoy, enDosSemanas)
	esperarErrorDeValidacion(t, err, "estado")
}

func TestSprint_Cerrar_PasaAFinalizado(t *testing.T) {
	s := sprintActivo(t)
	if err := s.Cerrar(); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if s.Estado != SprintFinalizado {
		t.Errorf("se esperaba estado %q, se obtuvo %q", SprintFinalizado, s.Estado)
	}
}

func TestSprint_Cerrar_SoloDesdeActivo(t *testing.T) {
	pendiente, _ := NewSprint(1)
	esperarErrorDeValidacion(t, pendiente.Cerrar(), "estado")

	finalizado := sprintActivo(t)
	_ = finalizado.Cerrar()
	esperarErrorDeValidacion(t, finalizado.Cerrar(), "estado")
}
