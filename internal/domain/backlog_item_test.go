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

func TestEsEstimacionValida(t *testing.T) {
	for _, v := range []int{1, 2, 3, 5, 8, 13, 21} {
		if !EsEstimacionValida(v) {
			t.Errorf("%d Story Points deberían ser válidos", v)
		}
	}
	for _, v := range []int{0, -1, 4, 22, 34} {
		if EsEstimacionValida(v) {
			t.Errorf("%d Story Points no deberían ser válidos", v)
		}
	}
}

func strPtr(v string) *string { return &v }

func historiaParaEditar() BacklogItem {
	sprintID := int64(7)
	return BacklogItem{
		ID:           3,
		ProyectoID:   1,
		Titulo:       "Título original",
		Descripcion:  "Descripción original",
		Prioridad:    PrioridadShould,
		Estado:       Estado("En progreso"),
		ValorNegocio: intPtr(8),
		EstimacionSP: intPtr(3),
		SprintID:     &sprintID,
	}
}

func TestBacklogItem_ConDatosEditados_Valida(t *testing.T) {
	original := historiaParaEditar()

	editada, err := original.ConDatosEditados("  Título nuevo  ", "  Descripción nueva  ", strPtr("  Dado, cuando, entonces  "), intPtr(5))
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo %v", err)
	}

	if editada.Titulo != "Título nuevo" || editada.Descripcion != "Descripción nueva" {
		t.Errorf("se esperaban título y descripción normalizados, se obtuvo %q %q", editada.Titulo, editada.Descripcion)
	}
	if editada.CriteriosAceptacion == nil || *editada.CriteriosAceptacion != "Dado, cuando, entonces" {
		t.Errorf("se esperaban los criterios normalizados, se obtuvo %v", editada.CriteriosAceptacion)
	}
	if editada.EstimacionSP == nil || *editada.EstimacionSP != 5 {
		t.Errorf("se esperaban 5 Story Points, se obtuvo %v", editada.EstimacionSP)
	}
	if editada.ID != original.ID || editada.ProyectoID != original.ProyectoID || editada.Prioridad != original.Prioridad ||
		editada.Estado != original.Estado || editada.ValorNegocio != original.ValorNegocio || editada.SprintID != original.SprintID {
		t.Errorf("la edición modificó campos que se conservan: %+v", editada)
	}
}

func TestBacklogItem_ConDatosEditados_CamposOpcionalesVacios(t *testing.T) {
	for _, criterios := range []*string{nil, strPtr(""), strPtr("   ")} {
		editada, err := historiaParaEditar().ConDatosEditados("Título", "Descripción", criterios, nil)
		if err != nil {
			t.Fatalf("no se esperaba error, se obtuvo %v", err)
		}
		if editada.CriteriosAceptacion != nil {
			t.Errorf("se esperaban criterios vacíos (nil), se obtuvo %q", *editada.CriteriosAceptacion)
		}
		if editada.EstimacionSP != nil {
			t.Errorf("se esperaba la historia sin estimación, se obtuvo %d", *editada.EstimacionSP)
		}
	}
}

func TestBacklogItem_ConDatosEditados_LimitesAceptados(t *testing.T) {
	_, err := historiaParaEditar().ConDatosEditados(strings.Repeat("a", 200), strings.Repeat("b", 2000), strPtr(strings.Repeat("c", 2000)), intPtr(21))
	if err != nil {
		t.Fatalf("los valores en el límite deberían aceptarse, se obtuvo %v", err)
	}
}

func TestBacklogItem_ConDatosEditados_Validaciones(t *testing.T) {
	casos := []struct {
		nombre          string
		titulo          string
		descripcion     string
		criterios       *string
		estimacion      *int
		campoEsperado   string
		mensajeEsperado string
	}{
		{"titulo vacio", "", "Descripción", nil, nil, "titulo", "el título es obligatorio"},
		{"titulo solo espacios", "   ", "Descripción", nil, nil, "titulo", "el título es obligatorio"},
		{"titulo demasiado largo", strings.Repeat("a", 201), "Descripción", nil, nil, "titulo", "el título no puede superar los 200 caracteres"},
		{"descripcion vacia", "Título", "  ", nil, nil, "descripcion", "la descripción es obligatoria"},
		{"descripcion demasiado larga", "Título", strings.Repeat("a", 2001), nil, nil, "descripcion", "la descripción no puede superar los 2000 caracteres"},
		{"criterios demasiado largos", "Título", "Descripción", strPtr(" " + strings.Repeat("a", 2001) + " "), nil, "criterios_aceptacion", "los criterios de aceptación no pueden superar los 2000 caracteres"},
		{"estimacion fuera de la escala", "Título", "Descripción", nil, intPtr(4), "estimacion_sp", ErrEstimacionInvalida.Mensaje},
		{"estimacion cero", "Título", "Descripción", nil, intPtr(0), "estimacion_sp", ErrEstimacionInvalida.Mensaje},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			_, err := historiaParaEditar().ConDatosEditados(tc.titulo, tc.descripcion, tc.criterios, tc.estimacion)
			var verr ValidationError
			if !errors.As(err, &verr) {
				t.Fatalf("se esperaba ValidationError, se obtuvo %v", err)
			}
			if verr.Campo != tc.campoEsperado || verr.Mensaje != tc.mensajeEsperado {
				t.Errorf("se esperaba %q/%q, se obtuvo %q/%q", tc.campoEsperado, tc.mensajeEsperado, verr.Campo, verr.Mensaje)
			}
		})
	}
}

func TestBacklogItem_ConDatosEditados_InvalidaNoModificaLaOriginal(t *testing.T) {
	original := historiaParaEditar()

	_, err := original.ConDatosEditados("", strings.Repeat("a", 2001), nil, intPtr(4))
	if err == nil {
		t.Fatalf("se esperaba un error de validación")
	}

	if original.Titulo != "Título original" || original.Descripcion != "Descripción original" ||
		original.CriteriosAceptacion != nil || original.EstimacionSP == nil || *original.EstimacionSP != 3 {
		t.Errorf("la historia original fue modificada: %+v", original)
	}
}
