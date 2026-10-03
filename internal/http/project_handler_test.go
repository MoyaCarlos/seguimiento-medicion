package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

var errPersistencia = errors.New("fallo de persistencia")

type stubProjectRepository struct {
	proyectos     map[int64]domain.Project
	miembros      []domain.Membership
	listaMiembros []domain.Member
	err           error
	creates       int
}

func (s *stubProjectRepository) Create(p *domain.Project) error {
	if s.err != nil {
		return s.err
	}
	s.creates++
	if s.proyectos == nil {
		s.proyectos = map[int64]domain.Project{}
	}
	if p.ID == 0 {
		p.ID = int64(s.creates)
	}
	p.CreatedAt = time.Now().UTC()
	s.proyectos[p.ID] = *p
	return nil
}

func (s *stubProjectRepository) GetByID(id int64) (*domain.Project, error) {
	if s.err != nil {
		return nil, s.err
	}
	p, ok := s.proyectos[id]
	if !ok {
		return nil, repository.ErrNoEncontrado
	}
	return &p, nil
}

func (s *stubProjectRepository) Update(p *domain.Project) error {
	if _, ok := s.proyectos[p.ID]; !ok {
		return repository.ErrNoEncontrado
	}
	s.proyectos[p.ID] = *p
	return nil
}

func (s *stubProjectRepository) AddMember(m *domain.Membership) error {
	for _, e := range s.miembros {
		if e.ProjectID == m.ProjectID && e.UserID == m.UserID {
			return repository.ErrMiembroDuplicado
		}
	}
	s.miembros = append(s.miembros, *m)
	return nil
}

func (s *stubProjectRepository) ListMembers(int64) ([]domain.Member, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.listaMiembros, nil
}

type stubUserRepository struct {
	usuarios map[string]domain.User
	creates  int
}

func (s *stubUserRepository) FindByNormalizedName(normalized string) (*domain.User, error) {
	if u, ok := s.usuarios[normalized]; ok {
		return &u, nil
	}
	return nil, repository.ErrNoEncontrado
}

func (s *stubUserRepository) Create(u *domain.User) error {
	if s.usuarios == nil {
		s.usuarios = map[string]domain.User{}
	}
	if _, ok := s.usuarios[u.NormalizedName]; ok {
		return repository.ErrUsuarioDuplicado
	}
	s.creates++
	u.ID = int64(s.creates)
	u.CreatedAt = time.Now().UTC()
	s.usuarios[u.NormalizedName] = *u
	return nil
}

func nuevoProjectHandlerDePrueba() (*ProjectHandler, *stubProjectRepository) {
	proyectos := &stubProjectRepository{}
	usuarios := &stubUserRepository{}
	crear := service.NewCrearProyecto(proyectos, usuarios)
	obtener := service.NewObtenerProyecto(proyectos)
	editar := service.NewEditarProyecto(proyectos)
	asignar := service.NewAsignarIntegrante(proyectos, usuarios)
	listar := service.NewListarIntegrantes(proyectos)
	return NewProjectHandler(crear, obtener, editar, asignar, listar), proyectos
}

// crearProyectoViaHTTP usa el propio handler (y sus repos) para preparar un proyecto.
func crearProyectoViaHTTP(t *testing.T, handler *ProjectHandler, nombre string) int64 {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBufferString(
		`{"nombre":"`+nombre+`","creador":"Ana"}`,
	))
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

func TestProjectHandler_Crear_Exitosa(t *testing.T) {
	handler, repo := nuevoProjectHandlerDePrueba()
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBufferString(
		`{"nombre":"Software Metrics & Estimation","descripcion":"Sistema","creador":"Ana","fecha_inicio":"2026-03-01","fecha_fin":"2026-11-30"}`,
	))
	handler.Crear(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp proyectoResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.ID == 0 || resp.Nombre != "Software Metrics & Estimation" {
		t.Errorf("respuesta inesperada: %+v", resp)
	}
	if resp.FechaInicio == nil || *resp.FechaInicio != "2026-03-01" {
		t.Errorf("fecha_inicio inesperada: %v", resp.FechaInicio)
	}
	if len(repo.miembros) != 1 || repo.miembros[0].Role != domain.RolScrumMaster {
		t.Errorf("se esperaba el creador como scrum_master")
	}
}

func TestProjectHandler_Crear_NombreVacio(t *testing.T) {
	handler, repo := nuevoProjectHandlerDePrueba()
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBufferString(`{"nombre":"","creador":"Ana"}`))
	handler.Crear(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Campo != "nombre" {
		t.Errorf("se esperaba campo nombre, se obtuvo %q", resp.Campo)
	}
	if repo.creates != 0 {
		t.Error("no se debía persistir el proyecto")
	}
}

func TestProjectHandler_Crear_FechaInvalida(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBufferString(
		`{"nombre":"Proyecto","creador":"Ana","fecha_inicio":"2026-11-30","fecha_fin":"2026-03-01"}`,
	))
	handler.Crear(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Campo != "fecha_fin" {
		t.Errorf("se esperaba campo fecha_fin, se obtuvo %q", resp.Campo)
	}
}

func TestProjectHandler_Crear_JSONInvalido(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBufferString(`{no-es-json`))
	handler.Crear(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
}

func TestProjectHandler_Crear_ErrorPersistencia(t *testing.T) {
	handler, repo := nuevoProjectHandlerDePrueba()
	repo.err = errPersistencia
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBufferString(`{"nombre":"Proyecto","creador":"Ana"}`))
	handler.Crear(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("se esperaba 500, se obtuvo %d", rec.Code)
	}
}

func TestProjectHandler_Obtener_Existente(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	id := crearProyectoViaHTTP(t, handler, "Proyecto")
	idStr := strconv.FormatInt(id, 10)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects/"+idStr, nil)
	req.SetPathValue("id", idStr)
	handler.Obtener(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
	var resp proyectoResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.ID != id {
		t.Errorf("se esperaba el proyecto %d, se obtuvo %d", id, resp.ID)
	}
}

func TestProjectHandler_Obtener_Inexistente(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects/no-existe", nil)
	req.SetPathValue("id", "no-existe")
	handler.Obtener(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestProjectHandler_Obtener_IDNoNumerico(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects/abc", nil)
	req.SetPathValue("id", "abc")
	handler.Obtener(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Campo != "id" {
		t.Errorf("se esperaba campo id, se obtuvo %q", resp.Campo)
	}
}

func TestProjectHandler_Obtener_IDNoPositivo(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects/0", nil)
	req.SetPathValue("id", "0")
	handler.Obtener(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Campo != "id" {
		t.Errorf("se esperaba campo id, se obtuvo %q", resp.Campo)
	}
}

func TestProjectHandler_Editar_Exitosa(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	id := crearProyectoViaHTTP(t, handler, "Original")
	idStr := strconv.FormatInt(id, 10)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/projects/"+idStr, bytes.NewBufferString(`{"nombre":"Corregido","descripcion":"nueva"}`))
	req.SetPathValue("id", idStr)
	handler.Editar(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp proyectoResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Nombre != "Corregido" {
		t.Errorf("no se actualizó el nombre: %+v", resp)
	}
}

func TestProjectHandler_Editar_NombreVacio(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	id := crearProyectoViaHTTP(t, handler, "Original")
	idStr := strconv.FormatInt(id, 10)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/projects/"+idStr, bytes.NewBufferString(`{"nombre":""}`))
	req.SetPathValue("id", idStr)
	handler.Editar(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Campo != "nombre" {
		t.Errorf("se esperaba campo nombre, se obtuvo %q", resp.Campo)
	}
}

func TestProjectHandler_Editar_Inexistente(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/projects/no-existe", bytes.NewBufferString(`{"nombre":"X"}`))
	req.SetPathValue("id", "no-existe")
	handler.Editar(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestProjectHandler_AsignarIntegrante_Exitosa(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	id := crearProyectoViaHTTP(t, handler, "Proyecto")
	idStr := strconv.FormatInt(id, 10)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/projects/"+idStr+"/members", bytes.NewBufferString(`{"nombre":"Jimena","rol":"product_builder"}`))
	req.SetPathValue("id", idStr)
	handler.AsignarIntegrante(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp integranteResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Nombre != "Jimena" || resp.Rol != "product_builder" {
		t.Errorf("integrante inesperado: %+v", resp)
	}
}

func TestProjectHandler_AsignarIntegrante_RolInvalido(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	id := crearProyectoViaHTTP(t, handler, "Proyecto")
	idStr := strconv.FormatInt(id, 10)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/projects/"+idStr+"/members", bytes.NewBufferString(`{"nombre":"Jimena","rol":"Product Owner"}`))
	req.SetPathValue("id", idStr)
	handler.AsignarIntegrante(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Campo != "rol" {
		t.Errorf("se esperaba campo rol, se obtuvo %q", resp.Campo)
	}
}

func TestProjectHandler_AsignarIntegrante_Duplicado(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	id := crearProyectoViaHTTP(t, handler, "Proyecto")
	idStr := strconv.FormatInt(id, 10)

	cuerpo := `{"nombre":"Jimena","rol":"product_builder"}`
	req1 := httptest.NewRequest(http.MethodPost, "/projects/"+idStr+"/members", bytes.NewBufferString(cuerpo))
	req1.SetPathValue("id", idStr)
	handler.AsignarIntegrante(httptest.NewRecorder(), req1)

	rec := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/projects/"+idStr+"/members", bytes.NewBufferString(cuerpo))
	req2.SetPathValue("id", idStr)
	handler.AsignarIntegrante(rec, req2)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400, se obtuvo %d", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Campo != "integrante" {
		t.Errorf("se esperaba campo integrante, se obtuvo %q", resp.Campo)
	}
}

func TestProjectHandler_AsignarIntegrante_ProyectoInexistente(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/projects/no-existe/members", bytes.NewBufferString(`{"nombre":"Jimena","rol":"product_builder"}`))
	req.SetPathValue("id", "no-existe")
	handler.AsignarIntegrante(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d", rec.Code)
	}
}

func TestProjectHandler_ListarIntegrantes(t *testing.T) {
	handler, repo := nuevoProjectHandlerDePrueba()
	id := crearProyectoViaHTTP(t, handler, "Proyecto")
	idStr := strconv.FormatInt(id, 10)
	repo.listaMiembros = []domain.Member{{UserID: 1, ProjectID: id, Name: "Ana", Role: domain.RolScrumMaster}}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects/"+idStr+"/members", nil)
	req.SetPathValue("id", idStr)
	handler.ListarIntegrantes(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
	}
	var resp []integranteResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if len(resp) != 1 || resp[0].Nombre != "Ana" {
		t.Errorf("listado inesperado: %+v", resp)
	}
}
