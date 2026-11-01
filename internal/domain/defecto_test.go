package domain

import (
	"strings"
	"testing"
)

var sprintDeDeteccion = Sprint{ID: 7, ProyectoID: 1, Estado: SprintActivo}

func TestNewDefecto_Valido(t *testing.T) {
	d, err := NewDefecto(1, 3, sprintDeDeteccion, "  El total no se actualiza  ", "Alta")

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if d.ProyectoID != 1 || d.HistoriaID != 3 || d.Descripcion != "El total no se actualiza" ||
		d.Severidad != SeveridadAlta || d.Estado != DefectoPendiente ||
		d.SprintDeteccionID != sprintDeDeteccion.ID || d.SprintResolucionID != nil {
		t.Errorf("defecto inesperado: %+v", d)
	}
}

func TestNewDefecto_SeveridadesValidas(t *testing.T) {
	for _, severidad := range []string{"Alta", "Media", "Baja"} {
		t.Run(severidad, func(t *testing.T) {
			d, err := NewDefecto(1, 3, sprintDeDeteccion, "Descripción", severidad)
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if d.Severidad != Severidad(severidad) {
				t.Errorf("se esperaba la severidad %s, se obtuvo %s", severidad, d.Severidad)
			}
		})
	}
}

func TestNewDefecto_DescripcionDe2000Caracteres(t *testing.T) {
	if _, err := NewDefecto(1, 3, sprintDeDeteccion, strings.Repeat("á", 2000), "Baja"); err != nil {
		t.Fatalf("no se esperaba error con 2000 caracteres: %v", err)
	}
}

func TestNewDefecto_Rechazos(t *testing.T) {
	casos := []struct {
		nombre      string
		descripcion string
		severidad   string
		campo       string
		mensaje     string
	}{
		{"descripción vacía", "", "Alta", "descripcion", "la descripción es obligatoria"},
		{"descripción solo con espacios", "   ", "Alta", "descripcion", "la descripción es obligatoria"},
		{"descripción de 2001 caracteres", strings.Repeat("á", 2001), "Alta", "descripcion", "la descripción no puede superar los 2000 caracteres"},
		{"severidad vacía", "Descripción", "", "severidad", "la severidad debe ser Alta, Media o Baja"},
		{"severidad en minúsculas", "Descripción", "alta", "severidad", "la severidad debe ser Alta, Media o Baja"},
		{"severidad con espacios", "Descripción", " Alta", "severidad", "la severidad debe ser Alta, Media o Baja"},
		{"severidad fuera de la escala", "Descripción", "Crítica", "severidad", "la severidad debe ser Alta, Media o Baja"},
		{"descripción antes que severidad", "", "Crítica", "descripcion", "la descripción es obligatoria"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := NewDefecto(1, 3, sprintDeDeteccion, c.descripcion, c.severidad)
			esperarErrorDeValidacion(t, err, c.campo)
			if verr, ok := err.(ValidationError); ok && verr.Mensaje != c.mensaje {
				t.Errorf("se esperaba el mensaje %q, se obtuvo %q", c.mensaje, verr.Mensaje)
			}
		})
	}
}
