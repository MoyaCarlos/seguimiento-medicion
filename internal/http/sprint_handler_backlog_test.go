package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type historiasDeSprintFake struct {
	historias []domain.BacklogItem
}

func (f historiasDeSprintFake) ListarPorSprint(context.Context, int64) ([]domain.BacklogItem, error) {
	return f.historias, nil
}

func nuevoSprintHandlerConBacklog(sprints repository.SprintRepository, historias repository.HistoriasDeSprintRepository) *SprintHandler {
	return NewSprintHandler(nil, nil, nil, nil, service.NewConsultarSprintBacklog(sprints, historias), sprints)
}

func peticionConsultarBacklog(id string) *nethttp.Request {
	req := httptest.NewRequest(nethttp.MethodGet, "/sprints/"+id+"/backlog", nil)
	req.SetPathValue("id", id)
	return req
}

func TestSprintHandler_ConsultarBacklog_DevuelveLasHistoriasDelSprint(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	pendiente, err := domain.NewSprint(1)
	if err != nil {
		t.Fatalf("no se pudo crear el Sprint: %v", err)
	}
	guardado, err := sprints.Guardar(context.Background(), pendiente)
	if err != nil {
		t.Fatalf("no se pudo sembrar el Sprint: %v", err)
	}
	sprintID := guardado.ID
	historias := historiasDeSprintFake{historias: []domain.BacklogItem{
		{ID: 2, ProyectoID: 1, Titulo: "Segunda", SprintID: &sprintID},
		{ID: 1, ProyectoID: 1, Titulo: "Primera", SprintID: &sprintID},
	}}
	handler := nuevoSprintHandlerConBacklog(sprints, historias)
	rec := httptest.NewRecorder()

	handler.ConsultarBacklog(rec, peticionConsultarBacklog(strconv.FormatInt(sprintID, 10)))

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp []historiaResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if len(resp) != 2 || resp[0].ID != 2 || resp[1].ID != 1 {
		t.Fatalf("se esperaba el orden recibido [2, 1], se obtuvo %+v", resp)
	}
	for _, historia := range resp {
		if historia.SprintID == nil || *historia.SprintID != sprintID {
			t.Errorf("sprint_id = %v, quería %d", historia.SprintID, sprintID)
		}
	}
}

func TestSprintHandler_ConsultarBacklog_SinHistoriasDevuelveArregloVacio(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	pendiente, _ := domain.NewSprint(1)
	guardado, err := sprints.Guardar(context.Background(), pendiente)
	if err != nil {
		t.Fatalf("no se pudo sembrar el Sprint: %v", err)
	}
	handler := nuevoSprintHandlerConBacklog(sprints, historiasDeSprintFake{})
	rec := httptest.NewRecorder()

	handler.ConsultarBacklog(rec, peticionConsultarBacklog(strconv.FormatInt(guardado.ID, 10)))

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if cuerpo := rec.Body.String(); cuerpo != "[]\n" {
		t.Fatalf("se esperaba [], se obtuvo %q", cuerpo)
	}
}

func TestSprintHandler_ConsultarBacklog_SprintInexistente(t *testing.T) {
	handler := nuevoSprintHandlerConBacklog(repository.NewSprintRepositoryEnMemoria(), historiasDeSprintFake{})
	rec := httptest.NewRecorder()

	handler.ConsultarBacklog(rec, peticionConsultarBacklog("999"))

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if resp := decodificarError(t, rec); resp.Mensaje != "sprint no encontrado" {
		t.Errorf("mensaje inesperado: %q", resp.Mensaje)
	}
}

func TestSprintHandler_ConsultarBacklog_IDInvalido(t *testing.T) {
	handler := nuevoSprintHandlerConBacklog(repository.NewSprintRepositoryEnMemoria(), historiasDeSprintFake{})
	rec := httptest.NewRecorder()

	handler.ConsultarBacklog(rec, peticionConsultarBacklog("abc"))

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	resp := decodificarError(t, rec)
	if resp.Campo != "id" {
		t.Errorf("campo = %q, quería %q", resp.Campo, "id")
	}
	if resp.Mensaje != "identificador inválido" {
		t.Errorf("mensaje = %q, quería %q", resp.Mensaje, "identificador inválido")
	}
}
