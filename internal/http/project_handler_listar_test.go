package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProjectHandler_Listar_ConProyectos(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	primero := crearProyectoViaHTTP(t, handler, "Primero")
	segundo := crearProyectoViaHTTP(t, handler, "Segundo")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	handler.Listar(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var respuesta []proyectoResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &respuesta); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if len(respuesta) != 2 {
		t.Fatalf("se esperaban 2 proyectos, se obtuvieron %d", len(respuesta))
	}
	if respuesta[0].ID != primero || respuesta[1].ID != segundo {
		t.Errorf("se esperaban los ids %d y %d, se obtuvo %+v", primero, segundo, respuesta)
	}
}

func TestProjectHandler_Listar_Vacio(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	handler.Listar(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if body := strings.TrimSpace(rec.Body.String()); body != "[]" {
		t.Errorf("se esperaba un arreglo vacío, se obtuvo %q", body)
	}
}
