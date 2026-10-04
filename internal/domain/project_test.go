package domain

import (
	"strings"
	"testing"
	"time"
)

func TestNewProject_Valido(t *testing.T) {
	inicio := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)

	conFechas, err := NewProject("  Software Metrics  ", "  Descripción  ", &inicio, &fin)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if conFechas.Nombre != "Software Metrics" || conFechas.Descripcion != "Descripción" {
		t.Errorf("no se recortaron los campos: %+v", conFechas)
	}
	if conFechas.FechaInicio == nil || !conFechas.FechaInicio.Equal(inicio) {
		t.Errorf("fecha de inicio inesperada: %v", conFechas.FechaInicio)
	}
	if conFechas.FechaFin == nil || !conFechas.FechaFin.Equal(fin) {
		t.Errorf("fecha de fin inesperada: %v", conFechas.FechaFin)
	}

	sinFechas, err := NewProject("Proyecto corto", "", nil, nil)
	if err != nil {
		t.Fatalf("no se esperaba error sin fechas: %v", err)
	}
	if sinFechas.FechaInicio != nil || sinFechas.FechaFin != nil {
		t.Errorf("se esperaban fechas nulas: %+v", sinFechas)
	}
}

func TestNewProject_NombreObligatorio(t *testing.T) {
	for _, nombre := range []string{"", "   "} {
		_, err := NewProject(nombre, "", nil, nil)
		if err == nil {
			t.Fatalf("se esperaba error para el nombre %q", nombre)
		}
		verr, ok := err.(ValidationError)
		if !ok || verr.Campo != "nombre" {
			t.Fatalf("se esperaba ValidationError de nombre, se obtuvo %v", err)
		}
	}
}

func TestNewProject_NombreDemasiadoLargo(t *testing.T) {
	_, err := NewProject(strings.Repeat("a", projectNameMaxLen+1), "", nil, nil)
	verr, ok := err.(ValidationError)
	if !ok || verr.Campo != "nombre" {
		t.Fatalf("se esperaba ValidationError de nombre, se obtuvo %v", err)
	}
}

func TestNewProject_DescripcionDemasiadoLarga(t *testing.T) {
	_, err := NewProject("Proyecto", strings.Repeat("d", projectDescriptionMaxLen+1), nil, nil)
	verr, ok := err.(ValidationError)
	if !ok || verr.Campo != "descripcion" {
		t.Fatalf("se esperaba ValidationError de descripcion, se obtuvo %v", err)
	}
}

func TestNewProject_FechaFinAnteriorAlInicio(t *testing.T) {
	inicio := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	_, err := NewProject("Proyecto", "", &inicio, &fin)
	verr, ok := err.(ValidationError)
	if !ok || verr.Campo != "fecha_fin" {
		t.Fatalf("se esperaba ValidationError de fecha_fin, se obtuvo %v", err)
	}
}

func TestNewProject_FechaFinIgualAlInicioEsValida(t *testing.T) {
	fecha := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	if _, err := NewProject("Proyecto", "", &fecha, &fecha); err != nil {
		t.Fatalf("la fecha de fin igual al inicio debería ser válida: %v", err)
	}
}

func TestNewProject_NombresRepetidosPermitidos(t *testing.T) {
	for i := 0; i < 2; i++ {
		if _, err := NewProject("Mismo nombre", "", nil, nil); err != nil {
			t.Fatalf("se permiten nombres repetidos: %v", err)
		}
	}
}

func TestConDatosEditados_ConservaIdentidad(t *testing.T) {
	creado := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	base := Project{ID: 1, Nombre: "Original", Descripcion: "desc", CreadoEn: creado}
	fin := time.Date(2026, 12, 15, 0, 0, 0, 0, time.UTC)

	editado, err := base.ConDatosEditados("  Nuevo  ", "", nil, &fin)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if editado.Nombre != "Nuevo" || editado.Descripcion != "" {
		t.Errorf("campos no actualizados: %+v", editado)
	}
	if editado.FechaInicio != nil || editado.FechaFin == nil || !editado.FechaFin.Equal(fin) {
		t.Errorf("fechas inesperadas: inicio=%v fin=%v", editado.FechaInicio, editado.FechaFin)
	}
	if editado.ID != 1 || !editado.CreadoEn.Equal(creado) {
		t.Errorf("no se conservó la identidad: %+v", editado)
	}
}

func TestConDatosEditados_ReutilizaValidacionesDelAlta(t *testing.T) {
	inicio := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	base, _ := NewProject("Original", "", nil, nil)

	casos := []struct {
		nombre      string
		name        string
		description string
		start       *time.Time
		end         *time.Time
		campo       string
	}{
		{"nombre vacío", "   ", "", nil, nil, "nombre"},
		{"nombre largo", strings.Repeat("a", projectNameMaxLen+1), "", nil, nil, "nombre"},
		{"descripción larga", "Proyecto", strings.Repeat("d", projectDescriptionMaxLen+1), nil, nil, "descripcion"},
		{"fin anterior al inicio", "Proyecto", "", &inicio, &fin, "fecha_fin"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := base.ConDatosEditados(c.name, c.description, c.start, c.end)
			verr, ok := err.(ValidationError)
			if !ok || verr.Campo != c.campo {
				t.Fatalf("se esperaba ValidationError de %q, se obtuvo %v", c.campo, err)
			}
		})
	}
}
