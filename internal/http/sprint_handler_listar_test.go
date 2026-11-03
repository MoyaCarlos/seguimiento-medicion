package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func peticionListarSprints(id string) *nethttp.Request {
	req := httptest.NewRequest(nethttp.MethodGet, "/projects/"+id+"/sprints", nil)
	req.SetPathValue("id", id)
	return req
}

func TestSprintHandler_Listar_DevuelveLosSprintsDelProyecto(t *testing.T) {
	handler, sprints, proyectos := nuevoSprintHandlerConProyectos(t)
	proyectoID := sembrarProyectoParaSprints(t, proyectos)
	for range 2 {
		s, _ := domain.NewSprint(proyectoID)
		if _, err := sprints.Guardar(context.Background(), s); err != nil {
			t.Fatalf("no se pudo sembrar el Sprint: %v", err)
		}
	}
	rec := httptest.NewRecorder()

	handler.Listar(rec, peticionListarSprints(strconv.FormatInt(proyectoID, 10)))

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp []sprintResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if len(resp) != 2 || resp[0].ID >= resp[1].ID {
		t.Fatalf("se esperaban 2 Sprints ordenados por id, se obtuvo %+v", resp)
	}
	if resp[0].ProyectoID != proyectoID || resp[0].Estado != string(domain.SprintPendiente) || resp[0].SprintGoal != nil {
		t.Errorf("formato inesperado: %+v", resp[0])
	}
}

func TestSprintHandler_Listar_SinSprintsDevuelveArregloVacio(t *testing.T) {
	handler, _, proyectos := nuevoSprintHandlerConProyectos(t)
	proyectoID := sembrarProyectoParaSprints(t, proyectos)
	rec := httptest.NewRecorder()

	handler.Listar(rec, peticionListarSprints(strconv.FormatInt(proyectoID, 10)))

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if cuerpo := rec.Body.String(); cuerpo != "[]\n" {
		t.Fatalf("se esperaba [], se obtuvo %q", cuerpo)
	}
}

func TestSprintHandler_Listar_ProyectoInexistente(t *testing.T) {
	handler, _ := nuevoSprintHandlerDePrueba(t)
	rec := httptest.NewRecorder()

	handler.Listar(rec, peticionListarSprints("999"))

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
}

func TestSprintHandler_Listar_IDInvalido(t *testing.T) {
	handler, _ := nuevoSprintHandlerDePrueba(t)
	for _, id := range []string{"abc", "0", "-1"} {
		rec := httptest.NewRecorder()

		handler.Listar(rec, peticionListarSprints(id))

		if rec.Code != nethttp.StatusBadRequest {
			t.Errorf("id %q: se esperaba 400, se obtuvo %d", id, rec.Code)
		}
	}
}
