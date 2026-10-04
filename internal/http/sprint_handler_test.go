package http

import (
	"bytes"
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

const cuerpoIniciarValido = `{"sprint_goal":"Entregar el MVP","fecha_inicio":"2026-10-05","fecha_fin":"2026-10-19"}`

// sinHistorias simula un Sprint sin historias asignadas (el arrastre se
// prueba en el caso de uso, no acá).
type sinHistorias struct{}

func (sinHistorias) ListarPorSprint(context.Context, int64) ([]domain.BacklogItem, error) {
	return nil, nil
}
func (sinHistorias) QuitarDeSprint(context.Context, int64) error { return nil }

func nuevoSprintHandlerConProyectos(t *testing.T) (*SprintHandler, *repository.SprintRepositoryEnMemoria, *repository.ProjectRepositoryEnMemoria) {
	t.Helper()
	sprints := repository.NewSprintRepositoryEnMemoria()
	proyectos := repository.NewProjectRepositoryEnMemoria()
	handler := NewSprintHandler(
		service.NewCrearSprint(proyectos, sprints),
		service.NewIniciarSprint(sprints),
		service.NewCerrarSprint(sprints, sinHistorias{}),
		sprints,
	)
	return handler, sprints, proyectos
}

func nuevoSprintHandlerDePrueba(t *testing.T) (*SprintHandler, *repository.SprintRepositoryEnMemoria) {
	t.Helper()
	handler, sprints, _ := nuevoSprintHandlerConProyectos(t)
	return handler, sprints
}

func sembrarProyectoParaSprints(t *testing.T, repo *repository.ProjectRepositoryEnMemoria) int64 {
	t.Helper()
	p, _ := domain.NewProject("Software Metrics & Estimation", "", nil, nil)
	g, err := repo.Guardar(context.Background(), p)
	if err != nil {
		t.Fatalf("no se pudo sembrar el proyecto: %v", err)
	}
	return g.ID
}

func peticionCrearSprint(cuerpo string) *nethttp.Request {
	req := httptest.NewRequest(nethttp.MethodPost, "/sprints", bytes.NewBufferString(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func cuerpoCrearSprint(proyectoID int64) string {
	return `{"proyecto_id":` + strconv.FormatInt(proyectoID, 10) + `}`
}

func TestSprintHandler_Crear_Exitoso(t *testing.T) {
	handler, _, proyectos := nuevoSprintHandlerConProyectos(t)
	proyectoID := sembrarProyectoParaSprints(t, proyectos)
	rec := httptest.NewRecorder()

	handler.Crear(rec, peticionCrearSprint(cuerpoCrearSprint(proyectoID)))

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp sprintResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.ID == 0 || resp.ProyectoID != proyectoID || resp.Estado != "Pendiente" {
		t.Errorf("respuesta inesperada: %+v", resp)
	}
	if resp.SprintGoal != nil || resp.FechaInicio != nil || resp.FechaFin != nil {
		t.Errorf("un Sprint Pendiente no tiene Goal ni fechas: %+v", resp)
	}
}

func TestSprintHandler_Crear_ProyectoInexistente(t *testing.T) {
	handler, _ := nuevoSprintHandlerDePrueba(t)
	rec := httptest.NewRecorder()

	handler.Crear(rec, peticionCrearSprint(cuerpoCrearSprint(999)))

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
}

func TestSprintHandler_Crear_ErroresDeValidacion(t *testing.T) {
	casos := map[string]struct {
		cuerpo func(proyectoID int64) string
		campo  string
	}{
		"proyecto_id no positivo": {func(int64) string { return cuerpoCrearSprint(0) }, "proyecto_id"},
		"otro Sprint Pendiente":   {cuerpoCrearSprint, "estado"},
		"JSON inválido":           {func(int64) string { return `{no-es-json` }, ""},
	}
	for nombre, caso := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler, _, proyectos := nuevoSprintHandlerConProyectos(t)
			proyectoID := sembrarProyectoParaSprints(t, proyectos)
			handler.Crear(httptest.NewRecorder(), peticionCrearSprint(cuerpoCrearSprint(proyectoID)))
			rec := httptest.NewRecorder()

			handler.Crear(rec, peticionCrearSprint(caso.cuerpo(proyectoID)))

			if rec.Code != nethttp.StatusBadRequest {
				t.Fatalf("se esperaba 400, se obtuvo %d (%s)", rec.Code, rec.Body.String())
			}
			if resp := decodificarError(t, rec); resp.Campo != caso.campo {
				t.Errorf("se esperaba campo %q, se obtuvo %q", caso.campo, resp.Campo)
			}
		})
	}
}

func sembrarPendiente(t *testing.T, repo *repository.SprintRepositoryEnMemoria) int64 {
	t.Helper()
	s, _ := domain.NewSprint(1)
	g, err := repo.Guardar(context.Background(), s)
	if err != nil {
		t.Fatalf("no se pudo sembrar: %v", err)
	}
	return g.ID
}

func peticionSprint(id, accion, cuerpo string) *nethttp.Request {
	req := httptest.NewRequest(nethttp.MethodPost, "/sprints/"+id+"/"+accion, bytes.NewBufferString(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	return req
}

func decodificarError(t *testing.T, rec *httptest.ResponseRecorder) errorResponse {
	t.Helper()
	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v (%s)", err, rec.Body.String())
	}
	return resp
}

func TestSprintHandler_Iniciar_Exitoso(t *testing.T) {
	handler, repo := nuevoSprintHandlerDePrueba(t)
	id := sembrarPendiente(t, repo)
	rec := httptest.NewRecorder()

	handler.Iniciar(rec, peticionSprint(strconv.FormatInt(id, 10), "iniciar", cuerpoIniciarValido))

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp sprintResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.ID != id || resp.Estado != "Activo" {
		t.Errorf("respuesta inesperada: %+v", resp)
	}
	if resp.SprintGoal == nil || *resp.SprintGoal != "Entregar el MVP" ||
		resp.FechaInicio == nil || *resp.FechaInicio != "2026-10-05" ||
		resp.FechaFin == nil || *resp.FechaFin != "2026-10-19" {
		t.Errorf("no se devolvieron Goal y fechas: %+v", resp)
	}
}

func TestSprintHandler_Iniciar_ErroresDeValidacion(t *testing.T) {
	casos := map[string]struct {
		cuerpo string
		campo  string
	}{
		"fecha de fin no posterior": {`{"sprint_goal":"Goal","fecha_inicio":"2026-10-05","fecha_fin":"2026-10-05"}`, "fecha_fin"},
		"formato de fecha inválido": {`{"sprint_goal":"Goal","fecha_inicio":"05/10/2026","fecha_fin":"2026-10-19"}`, "fecha_inicio"},
		"goal vacío":                {`{"sprint_goal":"","fecha_inicio":"2026-10-05","fecha_fin":"2026-10-19"}`, "sprint_goal"},
	}
	for nombre, caso := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler, repo := nuevoSprintHandlerDePrueba(t)
			id := sembrarPendiente(t, repo)
			rec := httptest.NewRecorder()

			handler.Iniciar(rec, peticionSprint(strconv.FormatInt(id, 10), "iniciar", caso.cuerpo))

			if rec.Code != nethttp.StatusBadRequest {
				t.Fatalf("se esperaba 400, se obtuvo %d (%s)", rec.Code, rec.Body.String())
			}
			if resp := decodificarError(t, rec); resp.Campo != caso.campo {
				t.Errorf("se esperaba campo %q, se obtuvo %q", caso.campo, resp.Campo)
			}
		})
	}
}

func TestSprintHandler_Iniciar_OtroSprintActivo(t *testing.T) {
	handler, repo := nuevoSprintHandlerDePrueba(t)
	primero, segundo := sembrarPendiente(t, repo), sembrarPendiente(t, repo)
	handler.Iniciar(httptest.NewRecorder(), peticionSprint(strconv.FormatInt(primero, 10), "iniciar", cuerpoIniciarValido))
	rec := httptest.NewRecorder()

	handler.Iniciar(rec, peticionSprint(strconv.FormatInt(segundo, 10), "iniciar", cuerpoIniciarValido))

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	if resp := decodificarError(t, rec); resp.Campo != "estado" {
		t.Errorf("se esperaba campo estado, se obtuvo %q", resp.Campo)
	}
}

func TestSprintHandler_Iniciar_NoEncontrado(t *testing.T) {
	handler, _ := nuevoSprintHandlerDePrueba(t)
	rec := httptest.NewRecorder()

	handler.Iniciar(rec, peticionSprint("999", "iniciar", cuerpoIniciarValido))

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestSprintHandler_Iniciar_NoEncontradoAntesQueCuerpoInvalido(t *testing.T) {
	casos := map[string]string{
		"fecha mal formada": `{"sprint_goal":"Goal","fecha_inicio":"05/10/2026","fecha_fin":"2026-10-19"}`,
		"JSON malformado":   `{no-es-json`,
	}
	for nombre, cuerpo := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler, _ := nuevoSprintHandlerDePrueba(t)
			rec := httptest.NewRecorder()

			handler.Iniciar(rec, peticionSprint("999", "iniciar", cuerpo))

			if rec.Code != nethttp.StatusNotFound {
				t.Fatalf("se esperaba 404, se obtuvo %d (%s)", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestSprintHandler_Iniciar_PeticionInvalida(t *testing.T) {
	casos := map[string]struct{ id, cuerpo string }{
		"id no numérico": {"abc", cuerpoIniciarValido},
		"JSON inválido":  {"1", `{no-es-json`},
	}
	for nombre, caso := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler, repo := nuevoSprintHandlerDePrueba(t)
			sembrarPendiente(t, repo)
			rec := httptest.NewRecorder()

			handler.Iniciar(rec, peticionSprint(caso.id, "iniciar", caso.cuerpo))

			if rec.Code != nethttp.StatusBadRequest {
				t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
			}
		})
	}
}

func TestSprintHandler_Cerrar_Exitoso(t *testing.T) {
	handler, repo := nuevoSprintHandlerDePrueba(t)
	id := strconv.FormatInt(sembrarPendiente(t, repo), 10)
	handler.Iniciar(httptest.NewRecorder(), peticionSprint(id, "iniciar", cuerpoIniciarValido))
	rec := httptest.NewRecorder()

	handler.Cerrar(rec, peticionSprint(id, "cerrar", ""))

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp sprintResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.Estado != "Finalizado" {
		t.Errorf("se esperaba estado Finalizado, se obtuvo %q", resp.Estado)
	}
}

func TestSprintHandler_Cerrar_NoActivo(t *testing.T) {
	handler, repo := nuevoSprintHandlerDePrueba(t)
	id := strconv.FormatInt(sembrarPendiente(t, repo), 10)
	rec := httptest.NewRecorder()

	handler.Cerrar(rec, peticionSprint(id, "cerrar", ""))

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	if resp := decodificarError(t, rec); resp.Campo != "estado" {
		t.Errorf("se esperaba campo estado, se obtuvo %q", resp.Campo)
	}
}

func TestSprintHandler_Cerrar_NoEncontradoEIdInvalido(t *testing.T) {
	casos := map[string]struct {
		id     string
		codigo int
	}{
		"inexistente":    {"999", nethttp.StatusNotFound},
		"id no numérico": {"abc", nethttp.StatusBadRequest},
	}
	for nombre, caso := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler, _ := nuevoSprintHandlerDePrueba(t)
			rec := httptest.NewRecorder()

			handler.Cerrar(rec, peticionSprint(caso.id, "cerrar", ""))

			if rec.Code != caso.codigo {
				t.Fatalf("se esperaba %d, se obtuvo %d", caso.codigo, rec.Code)
			}
		})
	}
}
