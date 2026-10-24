package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestProjectHandler_Listar_ConProyectos(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	primero := crearProyectoViaHTTP(t, handler, "Primero")
	segundo := crearProyectoViaHTTPConFechas(t, handler, "Segundo", "2026-03-01", "2026-11-30")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	handler.Listar(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var respuesta []json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &respuesta); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if len(respuesta) != 2 {
		t.Fatalf("se esperaban 2 proyectos, se obtuvieron %d", len(respuesta))
	}

	var ids []int64
	for _, elemento := range respuesta {
		var proyecto proyectoResponse
		if err := json.Unmarshal(elemento, &proyecto); err != nil {
			t.Fatalf("elemento JSON inválido: %v", err)
		}
		ids = append(ids, proyecto.ID)
	}
	if ids[0] != primero || ids[1] != segundo {
		t.Errorf("se esperaban los ids %d y %d, se obtuvo %v", primero, segundo, ids)
	}

	for _, elemento := range respuesta {
		var proyecto proyectoResponse
		if err := json.Unmarshal(elemento, &proyecto); err != nil {
			t.Fatalf("elemento JSON inválido: %v", err)
		}
		individual := obtenerProyectoViaHTTP(t, handler, proyecto.ID)
		if !mismoJSON(elemento, individual) {
			t.Errorf(
				"el proyecto %d del listado difiere de GET /projects/%d:\nlistado: %s\nindividual: %s",
				proyecto.ID, proyecto.ID, elemento, individual,
			)
		}
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

func TestProjectHandler_Listar_ErrorDePersistencia(t *testing.T) {
	handler, repo := nuevoProjectHandlerDePrueba()
	repo.err = errPersistencia

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	handler.Listar(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("se esperaba 500, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	verificarErrorSinCampo(t, rec)
}

func crearProyectoViaHTTPConFechas(t *testing.T, handler *ProjectHandler, nombre, inicio, fin string) int64 {
	t.Helper()
	rec := httptest.NewRecorder()
	cuerpo := `{"nombre":"` + nombre + `","creador":"Ana","fecha_inicio":"` + inicio + `","fecha_fin":"` + fin + `"}`
	req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(cuerpo))
	handler.Crear(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("no se pudo preparar el proyecto: %d (%s)", rec.Code, rec.Body.String())
	}
	var resp proyectoResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta inválida al preparar el proyecto: %v", err)
	}
	return resp.ID
}

func obtenerProyectoViaHTTP(t *testing.T, handler *ProjectHandler, id int64) json.RawMessage {
	t.Helper()
	idStr := strconv.FormatInt(id, 10)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects/"+idStr, nil)
	req.SetPathValue("id", idStr)
	handler.Obtener(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("no se pudo consultar el proyecto %d: %d (%s)", id, rec.Code, rec.Body.String())
	}
	var resp json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta inválida al consultar el proyecto %d: %v", id, err)
	}
	return resp
}

func mismoJSON(a, b json.RawMessage) bool {
	var izquierda, derecha any
	if err := json.Unmarshal(a, &izquierda); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &derecha); err != nil {
		return false
	}
	return reflect.DeepEqual(izquierda, derecha)
}
