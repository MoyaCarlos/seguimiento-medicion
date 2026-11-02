package domain

import "testing"

func TestEstado_EsValida(t *testing.T) {
	casos := []struct {
		nombre   string
		estado   Estado
		esperado bool
	}{
		{"Nueva es válida", EstadoNueva, true},
		{"En progreso es válida", EstadoEnProgreso, true},
		{"Completada es válida", EstadoCompletada, true},
		{"un estado desconocido no es válido", Estado("Pendiente"), false},
		{"el estado vacío no es válido", Estado(""), false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := c.estado.EsValida(); got != c.esperado {
				t.Errorf("Estado(%q).EsValida() = %v, se esperaba %v", c.estado, got, c.esperado)
			}
		})
	}
}

func TestEstadoEnProgreso_Valor(t *testing.T) {
	if EstadoEnProgreso != "En progreso" {
		t.Errorf("EstadoEnProgreso = %q, se esperaba %q", EstadoEnProgreso, "En progreso")
	}
}
