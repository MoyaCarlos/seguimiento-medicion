package domain

import (
	"errors"
	"strings"
	"testing"
)

func intPtr(v int) *int { return &v }

func TestNewBacklogItem_CreacionValida(t *testing.T) {
	item, err := NewBacklogItem(1, "  Como usuario quiero iniciar sesión  ", "  Autenticación  ", PrioridadMust, intPtr(13))
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo %v", err)
	}
	if item.Titulo != "Como usuario quiero iniciar sesión" {
		t.Errorf("se esperaba el título normalizado, se obtuvo %q", item.Titulo)
	}
	if item.Descripcion != "Autenticación" {
		t.Errorf("se esperaba la descripción normalizada, se obtuvo %q", item.Descripcion)
	}
	if item.Estado != EstadoNueva {
		t.Errorf("se esperaba estado %q, se obtuvo %q", EstadoNueva, item.Estado)
	}
	if item.EstimacionSP != nil {
		t.Errorf("se esperaba estimación sin asignar (nil), se obtuvo %v", *item.EstimacionSP)
	}
	if item.ValorNegocio == nil || *item.ValorNegocio != 13 {
		t.Errorf("se esperaba valor de negocio 13")
	}
}

func TestNewBacklogItem_Validaciones(t *testing.T) {
	casos := []struct {
		nombre        string
		proyectoID    int64
		titulo        string
		descripcion   string
		prioridad     Prioridad
		valorNegocio  *int
		campoEsperado string
	}{
		{"proyecto cero", 0, "Título", "Descripción", PrioridadMust, nil, "proyecto_id"},
		{"proyecto negativo", -5, "Título", "Descripción", PrioridadMust, nil, "proyecto_id"},
		{"titulo vacio", 1, "", "Descripción", PrioridadMust, nil, "titulo"},
		{"titulo solo espacios", 1, "   ", "Descripción", PrioridadMust, nil, "titulo"},
		{"titulo demasiado largo", 1, strings.Repeat("a", 201), "Descripción", PrioridadMust, nil, "titulo"},
		{"descripcion vacia", 1, "Título", "", PrioridadMust, nil, "descripcion"},
		{"descripcion solo espacios", 1, "Título", "   ", PrioridadMust, nil, "descripcion"},
		{"descripcion demasiado larga", 1, "Título", strings.Repeat("a", 2001), PrioridadMust, nil, "descripcion"},
		{"prioridad invalida", 1, "Título", "Descripción", Prioridad("Urgente"), nil, "prioridad"},
		{"valor de negocio invalido", 1, "Título", "Descripción", PrioridadMust, intPtr(4), "valor_negocio"},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			_, err := NewBacklogItem(tc.proyectoID, tc.titulo, tc.descripcion, tc.prioridad, tc.valorNegocio)
			if err == nil {
				t.Fatalf("se esperaba un error de validación")
			}
			var verr ValidationError
			if !errors.As(err, &verr) {
				t.Fatalf("se esperaba ValidationError, se obtuvo %T", err)
			}
			if verr.Campo != tc.campoEsperado {
				t.Errorf("se esperaba campo %q, se obtuvo %q", tc.campoEsperado, verr.Campo)
			}
			if verr.Mensaje == "" {
				t.Errorf("el mensaje de validación no debe estar vacío")
			}
		})
	}
}

func TestNewBacklogItem_ValorNegocioFibonacci(t *testing.T) {
	for _, v := range []int{1, 2, 3, 5, 8, 13, 21} {
		if _, err := NewBacklogItem(1, "Título", "Descripción", PrioridadMust, intPtr(v)); err != nil {
			t.Errorf("valor de negocio %d debería ser válido, se obtuvo %v", v, err)
		}
	}

	item, err := NewBacklogItem(1, "Título", "Descripción", PrioridadMust, nil)
	if err != nil {
		t.Fatalf("valor de negocio nil debería ser válido, se obtuvo %v", err)
	}
	if item.ValorNegocio != nil {
		t.Errorf("se esperaba valor de negocio nil")
	}
}

func TestPrioridad_EsValida(t *testing.T) {
	validas := []Prioridad{PrioridadMust, PrioridadShould, PrioridadCould, PrioridadWont}
	for _, p := range validas {
		if !p.EsValida() {
			t.Errorf("la prioridad %q debería ser válida", p)
		}
	}
	if Prioridad("X").EsValida() {
		t.Errorf("la prioridad %q no debería ser válida", "X")
	}
}
