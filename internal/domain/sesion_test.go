package domain

import (
	"testing"
	"time"
)

func TestDuracionSesion_Es7Dias(t *testing.T) {
	if DuracionSesion != 7*24*time.Hour {
		t.Errorf("se esperaban 7 días, se obtuvo %v", time.Duration(DuracionSesion))
	}
}

func TestSesionVigente(t *testing.T) {
	vence := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	cerrada := vence.Add(-24 * time.Hour)

	casos := []struct {
		nombre    string
		cerradaEn *time.Time
		ahora     time.Time
		esperado  bool
	}{
		{"antes de vencer", nil, vence.Add(-time.Nanosecond), true},
		{"en el instante exacto de vencer", nil, vence, false},
		{"después de vencer", nil, vence.Add(time.Second), false},
		{"cerrada sin haber vencido", &cerrada, vence.Add(-time.Hour), false},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			s := Sesion{VenceEn: vence, CerradaEn: c.cerradaEn}
			if got := s.Vigente(c.ahora); got != c.esperado {
				t.Errorf("Vigente(%v) = %v, se esperaba %v", c.ahora, got, c.esperado)
			}
		})
	}
}
