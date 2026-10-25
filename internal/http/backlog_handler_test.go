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
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type stubBacklogRepository struct {
	guardados       []domain.BacklogItem
	backlog         []domain.BacklogItem
	errListar       error
	historias       map[int64]domain.BacklogItem
	actualizaciones []domain.BacklogItem
	sprints         *repository.SprintRepositoryEnMemoria
}

func (s *stubBacklogRepository) Guardar(_ context.Context, item domain.BacklogItem) (domain.BacklogItem, error) {
	item.ID = int64(len(s.guardados) + 1)
	s.guardados = append(s.guardados, item)
	return item, nil
}

func (s *stubBacklogRepository) ListarSinSprint(_ context.Context, _ int64) ([]domain.BacklogItem, error) {
	return s.backlog, s.errListar
}

func (s *stubBacklogRepository) ObtenerPorID(_ context.Context, id int64) (domain.BacklogItem, error) {
	h, ok := s.historias[id]
	if !ok {
		return domain.BacklogItem{}, domain.ErrHistoriaNoEncontrada
	}
	return h, nil
}

func (s *stubBacklogRepository) Actualizar(_ context.Context, item domain.BacklogItem) error {
	s.actualizaciones = append(s.actualizaciones, item)
	s.historias[item.ID] = item
	return nil
}

func nuevoHandlerDePrueba() (*BacklogHandler, *stubBacklogRepository) {
	repo := &stubBacklogRepository{
		historias: map[int64]domain.BacklogItem{},
		sprints:   repository.NewSprintRepositoryEnMemoria(),
	}
	proyectos := &stubProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1}}}
	return NewBacklogHandler(
		service.NewCrearHistoriaBacklog(repo, proyectos),
		service.NewConsultarProductBacklog(proyectos, repo),
		service.NewObtenerHistoria(repo),
		service.NewEditarHistoria(repo, repo.sprints),
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

const edicionValidaJSON = `{"titulo":"Iniciar sesión con Google","descripcion":"Con OAuth","criterios_aceptacion":"Dado, cuando, entonces","estimacion_sp":8}`

func sembrarHistoria(repo *stubBacklogRepository) domain.BacklogItem {
	h := domain.BacklogItem{ID: 3, ProyectoID: 1, Titulo: "Original", Descripcion: "Descripción",
		Prioridad: domain.PrioridadMust, Estado: domain.EstadoNueva}
	repo.historias[h.ID] = h
	return h
}

func editarHistoria(handler *BacklogHandler, id, cuerpo string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodPut, "/backlog/"+id, bytes.NewBufferString(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	handler.Editar(rec, req)
	return rec
}

func leerError(t *testing.T, rec *httptest.ResponseRecorder) errorResponse {
	t.Helper()
	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v (%s)", err, rec.Body.String())
	}
	return resp
}

func TestBacklogHandler_Editar_Exitosa(t *testing.T) {
	handler, repo := nuevoHandlerDePrueba()
	sembrarHistoria(repo)

	rec := editarHistoria(handler, "3", edicionValidaJSON)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp historiaResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.ID != 3 || resp.Titulo != "Iniciar sesión con Google" || resp.Estado != "Nueva" || resp.Prioridad != "M" {
		t.Errorf("respuesta inesperada: %+v", resp)
	}
	if resp.CriteriosAceptacion == nil || *resp.CriteriosAceptacion != "Dado, cuando, entonces" ||
		resp.EstimacionSP == nil || *resp.EstimacionSP != 8 {
		t.Errorf("criterios o estimación inesperados: %+v", resp)
	}
	if len(repo.actualizaciones) != 1 {
		t.Errorf("se esperaba persistir la edición")
	}
}

func TestBacklogHandler_Editar_InexistenteDominaSobreElCuerpo(t *testing.T) {
	cuerpos := map[string]string{
		"titulo vacio":            `{"titulo":"","descripcion":"d","criterios_aceptacion":null,"estimacion_sp":null}`,
		"falta una clave":         `{"titulo":"t"}`,
		"story points no enteros": `{"titulo":"t","descripcion":"d","criterios_aceptacion":null,"estimacion_sp":2.5}`,
	}
	for nombre, cuerpo := range cuerpos {
		t.Run(nombre, func(t *testing.T) {
			handler, _ := nuevoHandlerDePrueba()

			rec := editarHistoria(handler, "999", cuerpo)

			if rec.Code != nethttp.StatusNotFound {
				t.Fatalf("se esperaba 404, se obtuvo %d (%s)", rec.Code, rec.Body.String())
			}
			if resp := leerError(t, rec); resp.Mensaje != "historia no encontrada" {
				t.Errorf("mensaje inesperado: %q", resp.Mensaje)
			}
		})
	}
}

func TestBacklogHandler_Editar_Validaciones(t *testing.T) {
	casos := []struct {
		nombre  string
		cuerpo  string
		campo   string
		mensaje string
	}{
		{"titulo en blanco", `{"titulo":"   ","descripcion":"d","criterios_aceptacion":null,"estimacion_sp":null}`, "titulo", "el título es obligatorio"},
		{"story points fuera de la escala", `{"titulo":"t","descripcion":"d","criterios_aceptacion":null,"estimacion_sp":4}`, "estimacion_sp", domain.ErrEstimacionInvalida.Mensaje},
		{"story points no enteros", `{"titulo":"t","descripcion":"d","criterios_aceptacion":null,"estimacion_sp":2.5}`, "estimacion_sp", domain.ErrEstimacionInvalida.Mensaje},
		{"falta titulo", `{"descripcion":"d","criterios_aceptacion":null,"estimacion_sp":null}`, "titulo", "el campo es obligatorio"},
		{"falta descripcion", `{"titulo":"t","criterios_aceptacion":null,"estimacion_sp":null}`, "descripcion", "el campo es obligatorio"},
		{"faltan criterios", `{"titulo":"t","descripcion":"d","estimacion_sp":null}`, "criterios_aceptacion", "el campo es obligatorio"},
		{"falta estimacion", `{"titulo":"t","descripcion":"d","criterios_aceptacion":null}`, "estimacion_sp", "el campo es obligatorio"},
	}
	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			handler, repo := nuevoHandlerDePrueba()
			sembrarHistoria(repo)

			rec := editarHistoria(handler, "3", tc.cuerpo)

			if rec.Code != nethttp.StatusBadRequest {
				t.Fatalf("se esperaba 400, se obtuvo %d (%s)", rec.Code, rec.Body.String())
			}
			if resp := leerError(t, rec); resp.Campo != tc.campo || resp.Mensaje != tc.mensaje {
				t.Errorf("se esperaba %q/%q, se obtuvo %q/%q", tc.campo, tc.mensaje, resp.Campo, resp.Mensaje)
			}
			if len(repo.actualizaciones) != 0 {
				t.Errorf("no se esperaba persistir una edición inválida")
			}
		})
	}
}

func TestBacklogHandler_Editar_CamposOpcionalesNull(t *testing.T) {
	handler, repo := nuevoHandlerDePrueba()
	sembrarHistoria(repo)

	rec := editarHistoria(handler, "3", `{"titulo":"t","descripcion":"d","criterios_aceptacion":null,"estimacion_sp":null}`)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var crudo map[string]json.RawMessage
	_ = json.Unmarshal(rec.Body.Bytes(), &crudo)
	if string(crudo["criterios_aceptacion"]) != "null" || string(crudo["estimacion_sp"]) != "null" {
		t.Errorf("se esperaban criterios y estimación null: %s", rec.Body.String())
	}
}

func TestBacklogHandler_Editar_StoryPointsCongelados(t *testing.T) {
	handler, repo := nuevoHandlerDePrueba()
	sprint, _ := repo.sprints.Guardar(context.Background(), domain.Sprint{ProyectoID: 1, Estado: domain.SprintActivo})
	h := sembrarHistoria(repo)
	sp := 3
	h.SprintID, h.EstimacionSP = &sprint.ID, &sp
	repo.historias[h.ID] = h

	rec := editarHistoria(handler, "3", `{"titulo":"t","descripcion":"d","criterios_aceptacion":null,"estimacion_sp":5}`)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if resp := leerError(t, rec); resp.Campo != "estimacion_sp" || resp.Mensaje != domain.ErrEstimacionCongelada.Mensaje {
		t.Errorf("error inesperado: %+v", resp)
	}
	if len(repo.actualizaciones) != 0 {
		t.Errorf("no se esperaba persistir la edición")
	}
}

func TestBacklogHandler_Editar_JSONInvalido(t *testing.T) {
	handler, repo := nuevoHandlerDePrueba()
	sembrarHistoria(repo)

	rec := editarHistoria(handler, "3", `{no-es-json`)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	if resp := leerError(t, rec); resp.Mensaje != "JSON inválido" {
		t.Errorf("mensaje inesperado: %q", resp.Mensaje)
	}
}

func TestBacklogHandler_Editar_IDInvalido(t *testing.T) {
	handler, _ := nuevoHandlerDePrueba()

	rec := editarHistoria(handler, "abc", edicionValidaJSON)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}
