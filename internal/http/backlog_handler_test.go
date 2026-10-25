package http

import (
	"bytes"
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type stubBacklogRepository struct {
	guardados []domain.BacklogItem
}

func (s *stubBacklogRepository) Guardar(_ context.Context, item domain.BacklogItem) (domain.BacklogItem, error) {
	item.ID = int64(len(s.guardados) + 1)
	s.guardados = append(s.guardados, item)
	return item, nil
}

func nuevoHandlerDePrueba() (*BacklogHandler, *stubBacklogRepository) {
	repo := &stubBacklogRepository{}
	proyectos := &stubProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1}}}
	return NewBacklogHandler(service.NewCrearHistoriaBacklog(repo, proyectos)), repo
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
