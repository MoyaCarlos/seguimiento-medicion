package http

import (
	"context"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type stubReporteBacklog struct {
	backlog     []domain.BacklogItem
	porSprint   map[int64][]domain.BacklogItem
	arrastradas map[int64][]domain.BacklogItem
}

func (s *stubReporteBacklog) ListarSinSprint(_ context.Context, _ int64) ([]domain.BacklogItem, error) {
	if s.backlog == nil {
		return []domain.BacklogItem{}, nil
	}
	return s.backlog, nil
}

func (s *stubReporteBacklog) ListarPorSprint(_ context.Context, sprintID int64) ([]domain.BacklogItem, error) {
	historias := s.porSprint[sprintID]
	if historias == nil {
		return []domain.BacklogItem{}, nil
	}
	return historias, nil
}

func (s *stubReporteBacklog) ListarArrastradasPorSprint(_ context.Context, sprintID int64) ([]domain.BacklogItem, error) {
	historias := s.arrastradas[sprintID]
	if historias == nil {
		return []domain.BacklogItem{}, nil
	}
	return historias, nil
}

type generadorReporteStub struct {
	contenidoSprint   []byte
	contenidoProyecto []byte
	err               error
	llamadasSprint    int
	llamadasProyecto  int
}

func (g *generadorReporteStub) GenerarSprint(domain.ReporteSprint) ([]byte, error) {
	g.llamadasSprint++
	return g.contenidoSprint, g.err
}

func (g *generadorReporteStub) GenerarProyecto(domain.ReporteProyecto) ([]byte, error) {
	g.llamadasProyecto++
	return g.contenidoProyecto, g.err
}

type escenarioReporteHTTP struct {
	handler *ReporteHandler
	pdf     *generadorReporteStub
	sprints *repository.SprintRepositoryEnMemoria
}

func nuevoEscenarioReporteHTTP() escenarioReporteHTTP {
	sprints := repository.NewSprintRepositoryEnMemoria()
	proyectos := &stubProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1, Nombre: "Proyecto"}}}
	backlog := &stubReporteBacklog{}
	generar := service.NewGenerarReporte(proyectos, sprints, backlog, backlog, backlog, stubDatosMetricas{}, &stubDefectoRepository{})
	pdf := &generadorReporteStub{contenidoSprint: []byte("%PDF-1.4 contenido"), contenidoProyecto: []byte("%PDF-1.4 proyecto")}
	return escenarioReporteHTTP{handler: NewReporteHandler(generar, pdf), pdf: pdf, sprints: sprints}
}

func (e escenarioReporteHTTP) sprint(t *testing.T, estado domain.EstadoSprint) domain.Sprint {
	t.Helper()
	sprint := domain.Sprint{ProyectoID: 1, Estado: estado}
	if estado != domain.SprintPendiente {
		sprint.SprintGoal = "Objetivo"
		sprint.FechaInicio = time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
		sprint.FechaFin = time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC)
	}
	guardado, err := e.sprints.Guardar(context.Background(), sprint)
	if err != nil {
		t.Fatalf("no se pudo sembrar el Sprint: %v", err)
	}
	return guardado
}

func (e escenarioReporteHTTP) pedirSprint(id string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(nethttp.MethodGet, "/sprints/"+id+"/reporte", nil)
	req.SetPathValue("id", id)
	rec := httptest.NewRecorder()
	e.handler.Sprint(rec, req)
	return rec
}

func (e escenarioReporteHTTP) pedirProyecto(id string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(nethttp.MethodGet, "/projects/"+id+"/reporte", nil)
	req.SetPathValue("id", id)
	rec := httptest.NewRecorder()
	e.handler.Proyecto(rec, req)
	return rec
}

func TestReporteHandler_Sprint_DescargaElPDF(t *testing.T) {
	e := nuevoEscenarioReporteHTTP()
	sprint := e.sprint(t, domain.SprintFinalizado)

	rec := e.pedirSprint("1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/pdf" {
		t.Errorf("Content-Type: se esperaba application/pdf, se obtuvo %q", got)
	}
	esperado := `attachment; filename="reporte-proyecto-1-sprint-` + "1" + `.pdf"`
	if got := rec.Header().Get("Content-Disposition"); got != esperado {
		t.Errorf("Content-Disposition: se esperaba %q, se obtuvo %q", esperado, got)
	}
	if got := rec.Body.String(); got != string(e.pdf.contenidoSprint) {
		t.Errorf("cuerpo inesperado: %q", got)
	}
	if e.pdf.llamadasSprint != 1 {
		t.Errorf("se esperaba una generación, se hicieron %d", e.pdf.llamadasSprint)
	}
	if sprint.Numero != 1 {
		t.Errorf("se esperaba el Sprint número 1, se obtuvo %d", sprint.Numero)
	}
}

func TestReporteHandler_Sprint_NoFinalizado(t *testing.T) {
	for _, estado := range []domain.EstadoSprint{domain.SprintPendiente, domain.SprintActivo} {
		t.Run(string(estado), func(t *testing.T) {
			e := nuevoEscenarioReporteHTTP()
			e.sprint(t, estado)

			rec := e.pedirSprint("1")

			if rec.Code != nethttp.StatusBadRequest {
				t.Fatalf("se esperaba 400, se obtuvo %d (%s)", rec.Code, rec.Body.String())
			}
			cuerpo := cuerpoJSON(t, rec)
			if cuerpo["campo"] != "estado" {
				t.Errorf("se esperaba el campo estado, se obtuvo %v", cuerpo)
			}
			if e.pdf.llamadasSprint != 0 {
				t.Errorf("no se debía generar el PDF")
			}
		})
	}
}

func TestReporteHandler_Sprint_Inexistente(t *testing.T) {
	e := nuevoEscenarioReporteHTTP()

	rec := e.pedirSprint("999")

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if cuerpo := cuerpoJSON(t, rec); cuerpo["mensaje"] != "sprint no encontrado" {
		t.Errorf("mensaje inesperado: %v", cuerpo)
	}
	if e.pdf.llamadasSprint != 0 {
		t.Errorf("no se debía generar el PDF")
	}
}

func TestReporteHandler_Sprint_IdInvalido(t *testing.T) {
	for _, id := range []string{"abc", "0"} {
		t.Run(id, func(t *testing.T) {
			e := nuevoEscenarioReporteHTTP()

			rec := e.pedirSprint(id)

			if rec.Code != nethttp.StatusBadRequest {
				t.Fatalf("se esperaba 400, se obtuvo %d (%s)", rec.Code, rec.Body.String())
			}
			if cuerpo := cuerpoJSON(t, rec); cuerpo["campo"] != "id" {
				t.Errorf("se esperaba el campo id, se obtuvo %v", cuerpo)
			}
		})
	}
}

func TestReporteHandler_Sprint_ErrorDeGeneracion(t *testing.T) {
	e := nuevoEscenarioReporteHTTP()
	e.pdf.err = errors.New("falla de Maroto")
	e.sprint(t, domain.SprintFinalizado)

	rec := e.pedirSprint("1")

	if rec.Code != nethttp.StatusInternalServerError {
		t.Fatalf("se esperaba 500, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got == "application/pdf" {
		t.Errorf("no se debía responder como PDF")
	}
	if cuerpoJSON(t, rec)["mensaje"] != "no se pudo generar el reporte" {
		t.Errorf("mensaje inesperado: %s", rec.Body.String())
	}
	if rec.Body.Len() == 0 || (len(rec.Body.Bytes()) >= 4 && string(rec.Body.Bytes()[:4]) == "%PDF") {
		t.Errorf("no se debía entregar un archivo: %q", rec.Body.String())
	}
}

func TestReporteHandler_Proyecto_DescargaElPDF(t *testing.T) {
	e := nuevoEscenarioReporteHTTP()

	rec := e.pedirProyecto("1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/pdf" {
		t.Errorf("Content-Type: se esperaba application/pdf, se obtuvo %q", got)
	}
	esperado := `attachment; filename="reporte-proyecto-1.pdf"`
	if got := rec.Header().Get("Content-Disposition"); got != esperado {
		t.Errorf("Content-Disposition: se esperaba %q, se obtuvo %q", esperado, got)
	}
	if got := rec.Body.String(); got != string(e.pdf.contenidoProyecto) {
		t.Errorf("cuerpo inesperado: %q", got)
	}
	if e.pdf.llamadasProyecto != 1 {
		t.Errorf("se esperaba una generación, se hicieron %d", e.pdf.llamadasProyecto)
	}
}

func TestReporteHandler_Proyecto_Inexistente(t *testing.T) {
	e := nuevoEscenarioReporteHTTP()

	rec := e.pedirProyecto("999")

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if cuerpo := cuerpoJSON(t, rec); cuerpo["mensaje"] != "proyecto no encontrado" {
		t.Errorf("mensaje inesperado: %v", cuerpo)
	}
	if e.pdf.llamadasProyecto != 0 {
		t.Errorf("no se debía generar el PDF")
	}
}

func TestReporteHandler_Proyecto_IdInvalido(t *testing.T) {
	for _, id := range []string{"abc", "0"} {
		t.Run(id, func(t *testing.T) {
			e := nuevoEscenarioReporteHTTP()

			rec := e.pedirProyecto(id)

			if rec.Code != nethttp.StatusBadRequest {
				t.Fatalf("se esperaba 400, se obtuvo %d (%s)", rec.Code, rec.Body.String())
			}
			if cuerpo := cuerpoJSON(t, rec); cuerpo["campo"] != "id" {
				t.Errorf("se esperaba el campo id, se obtuvo %v", cuerpo)
			}
		})
	}
}

func TestReporteHandler_ErrorDeLecturaNoEscribeRespuesta(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	proyectos := &stubProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1, Nombre: "Proyecto"}}}
	backlog := &stubReporteBacklog{}
	falla := errors.New("falla de la base")
	pdf := &generadorReporteStub{contenidoSprint: []byte("%PDF-1.4"), contenidoProyecto: []byte("%PDF-1.4")}
	generar := service.NewGenerarReporte(proyectos, sprints, backlog, backlog, backlog, stubDatosMetricas{err: falla}, &stubDefectoRepository{})
	handler := NewReporteHandler(generar, pdf)

	finalizado, err := sprints.Guardar(context.Background(), domain.Sprint{ProyectoID: 1, Estado: domain.SprintFinalizado})
	if err != nil {
		t.Fatalf("no se pudo sembrar el Sprint: %v", err)
	}
	reqSprint := httptest.NewRequest(nethttp.MethodGet, "/sprints/1/reporte", nil)
	reqSprint.SetPathValue("id", "1")
	recSprint := httptest.NewRecorder()
	handler.Sprint(recSprint, reqSprint)

	reqProyecto := httptest.NewRequest(nethttp.MethodGet, "/projects/1/reporte", nil)
	reqProyecto.SetPathValue("id", "1")
	recProyecto := httptest.NewRecorder()
	handler.Proyecto(recProyecto, reqProyecto)

	if finalizado.ID == 0 {
		t.Fatal("el Sprint debía tener ID")
	}
	for nombre, rec := range map[string]*httptest.ResponseRecorder{"sprint": recSprint, "proyecto": recProyecto} {
		if rec.Code != nethttp.StatusInternalServerError {
			t.Errorf("%s: se esperaba 500, se obtuvo %d", nombre, rec.Code)
		}
		if got := rec.Header().Get("Content-Type"); got == "application/pdf" {
			t.Errorf("%s: no se debía responder como PDF", nombre)
		}
		if rec.Body.Len() > 0 && string(rec.Body.Bytes()[:4]) == "%PDF" {
			t.Errorf("%s: no se debía entregar un archivo", nombre)
		}
	}
	if pdf.llamadasSprint != 0 || pdf.llamadasProyecto != 0 {
		t.Errorf("no se debía generar el PDF: sprint=%d proyecto=%d", pdf.llamadasSprint, pdf.llamadasProyecto)
	}
}

func TestReporteHandler_Proyecto_ErrorDeGeneracion(t *testing.T) {
	e := nuevoEscenarioReporteHTTP()
	e.pdf.err = errors.New("falla de Maroto")

	rec := e.pedirProyecto("1")

	if rec.Code != nethttp.StatusInternalServerError {
		t.Fatalf("se esperaba 500, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got == "application/pdf" {
		t.Errorf("no se debía responder como PDF")
	}
	if cuerpoJSON(t, rec)["mensaje"] != "no se pudo generar el reporte" {
		t.Errorf("mensaje inesperado: %s", rec.Body.String())
	}
	if rec.Body.Len() == 0 || (len(rec.Body.Bytes()) >= 4 && string(rec.Body.Bytes()[:4]) == "%PDF") {
		t.Errorf("no se debía entregar un archivo: %q", rec.Body.String())
	}
}
