package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type stubBacklogRepository struct {
	guardados []domain.BacklogItem
	backlog   []domain.BacklogItem
	errListar error
}

func (s *stubBacklogRepository) Guardar(_ context.Context, item domain.BacklogItem) (domain.BacklogItem, error) {
	item.ID = int64(len(s.guardados) + 1)
	s.guardados = append(s.guardados, item)
	return item, nil
}

func (s *stubBacklogRepository) ListarSinSprint(_ context.Context, _ int64) ([]domain.BacklogItem, error) {
	return s.backlog, s.errListar
}

func nuevoHandlerDePrueba() (*BacklogHandler, *stubBacklogRepository) {
	repo := &stubBacklogRepository{}
	proyectos := &stubProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1}}}
	return NewBacklogHandler(
		service.NewCrearHistoriaBacklog(repo, proyectos),
		service.NewConsultarProductBacklog(proyectos, repo),
	), repo
}

func construirPeticion(cuerpo string) *nethttp.Request {
	req := httptest.NewRequest(nethttp.MethodPost, "/backlog", bytes.NewBufferString(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestBacklogHandler_Crear_Exitosa(t *testing.T) {
	handler, repo := nuevoHandlerDePrueba()
	rec := httptest.NewRecorder()

	handler.Crear(rec, construirPeticion(`{"proyecto_id":1,"titulo":"Iniciar sesión","descripcion":"Autenticación","prioridad":"M","valor_negocio":13}`))

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp historiaResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.ID != 1 || resp.Estado != "Nueva" {
		t.Errorf("respuesta inesperada: %+v", resp)
	}
	if resp.EstimacionSP != nil {
		t.Errorf("se esperaba estimacion_sp nulo")
	}
	if resp.ValorNegocio == nil || *resp.ValorNegocio != 13 {
		t.Errorf("se esperaba valor_negocio 13")
	}
	var crudo map[string]json.RawMessage
	_ = json.Unmarshal(rec.Body.Bytes(), &crudo)
	if v, ok := crudo["criterios_aceptacion"]; !ok || string(v) != "null" {
		t.Errorf("se esperaba criterios_aceptacion null, se obtuvo %q (presente: %v)", v, ok)
	}
	if len(repo.guardados) != 1 {
		t.Errorf("se esperaba 1 historia persistida")
	}
}

func TestBacklogHandler_Crear_TituloVacio(t *testing.T) {
	handler, repo := nuevoHandlerDePrueba()
	rec := httptest.NewRecorder()

	handler.Crear(rec, construirPeticion(`{"proyecto_id":1,"titulo":"","descripcion":"Descripción","prioridad":"M"}`))

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.Campo != "titulo" {
		t.Errorf("se esperaba campo titulo, se obtuvo %q", resp.Campo)
	}
	if len(repo.guardados) != 0 {
		t.Errorf("no se debía persistir ninguna historia")
	}
}

func TestBacklogHandler_Crear_ProyectoAusente(t *testing.T) {
	handler, _ := nuevoHandlerDePrueba()
	rec := httptest.NewRecorder()

	handler.Crear(rec, construirPeticion(`{"titulo":"Sin proyecto","descripcion":"Descripción","prioridad":"M"}`))

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.Campo != "proyecto_id" {
		t.Errorf("se esperaba campo proyecto_id, se obtuvo %q", resp.Campo)
	}
}

func TestBacklogHandler_Crear_JSONInvalido(t *testing.T) {
	handler, _ := nuevoHandlerDePrueba()
	rec := httptest.NewRecorder()

	handler.Crear(rec, construirPeticion(`{no-es-json`))

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func consultarBacklog(handler *BacklogHandler, id string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodGet, "/projects/"+id+"/backlog", nil)
	req.SetPathValue("id", id)
	handler.ConsultarBacklog(rec, req)
	return rec
}

func TestBacklogHandler_ConsultarBacklog_DevuelveLasHistorias(t *testing.T) {
	handler, repo := nuevoHandlerDePrueba()
	criterios, sp, valor := "Dado, cuando, entonces", 5, 13
	repo.backlog = []domain.BacklogItem{
		{ID: 2, ProyectoID: 1, Titulo: "Iniciar sesión", Descripcion: "Con email", CriteriosAceptacion: &criterios,
			Prioridad: domain.PrioridadMust, Estado: domain.EstadoNueva, ValorNegocio: &valor, EstimacionSP: &sp},
		{ID: 1, ProyectoID: 1, Titulo: "Cerrar sesión", Descripcion: "Salir",
			Prioridad: domain.PrioridadShould, Estado: domain.EstadoNueva},
	}

	rec := consultarBacklog(handler, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp []map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if len(resp) != 2 || string(resp[0]["id"]) != "2" || string(resp[1]["id"]) != "1" {
		t.Fatalf("se esperaban las historias 2 y 1 en ese orden: %s", rec.Body.String())
	}
	for _, campo := range []string{"id", "proyecto_id", "titulo", "descripcion", "criterios_aceptacion", "prioridad", "estado", "valor_negocio", "estimacion_sp"} {
		if _, ok := resp[0][campo]; !ok {
			t.Errorf("falta el campo %q en la respuesta", campo)
		}
	}
	if string(resp[0]["criterios_aceptacion"]) != `"Dado, cuando, entonces"` || string(resp[0]["estimacion_sp"]) != "5" {
		t.Errorf("criterios o estimación inesperados: %s", rec.Body.String())
	}
	if string(resp[1]["criterios_aceptacion"]) != "null" || string(resp[1]["estimacion_sp"]) != "null" {
		t.Errorf("se esperaban criterios y estimación null: %s", rec.Body.String())
	}
}

func TestBacklogHandler_ConsultarBacklog_VacioEsArreglo(t *testing.T) {
	handler, _ := nuevoHandlerDePrueba()

	rec := consultarBacklog(handler, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
	if got := bytes.TrimSpace(rec.Body.Bytes()); string(got) != "[]" {
		t.Errorf("se esperaba [], se obtuvo %s", got)
	}
}

func TestBacklogHandler_ConsultarBacklog_ProyectoInexistente(t *testing.T) {
	handler, _ := nuevoHandlerDePrueba()

	rec := consultarBacklog(handler, "999")

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Mensaje != "proyecto no encontrado" {
		t.Errorf("mensaje inesperado: %q", resp.Mensaje)
	}
}

func TestBacklogHandler_ConsultarBacklog_IDInvalido(t *testing.T) {
	handler, _ := nuevoHandlerDePrueba()

	rec := consultarBacklog(handler, "abc")

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestBacklogHandler_ConsultarBacklog_ErrorDePersistencia(t *testing.T) {
	handler, repo := nuevoHandlerDePrueba()
	repo.errListar = errors.New("falla de persistencia")

	rec := consultarBacklog(handler, "1")

	if rec.Code != nethttp.StatusInternalServerError {
		t.Fatalf("se esperaba 500, se obtuvo %d", rec.Code)
	}
}
