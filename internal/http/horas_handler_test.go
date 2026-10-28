package http

import (
	"bytes"
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type stubRegistroHorasRepository struct {
	registros []domain.RegistroHoras
}

func (s *stubRegistroHorasRepository) Registrar(_ context.Context, r domain.RegistroHoras) (domain.RegistroHoras, error) {
	r.ID = int64(len(s.registros) + 1)
	s.registros = append(s.registros, r)
	return r, nil
}

func (s *stubRegistroHorasRepository) ListarPorHistoria(_ context.Context, historiaID int64) ([]domain.RegistroHorasDetalle, error) {
	detalle := make([]domain.RegistroHorasDetalle, 0)
	for _, r := range s.registros {
		if r.HistoriaID == historiaID {
			detalle = append(detalle, domain.RegistroHorasDetalle{RegistroHoras: r, NombreIntegrante: "Ana"})
		}
	}
	return detalle, nil
}

// hoyCivil es el día local actual a las 00:00 UTC, como lo compara el dominio.
func hoyCivil() time.Time {
	ahora := time.Now()
	return time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 0, 0, 0, 0, time.UTC)
}

type escenarioHorasHTTP struct {
	handler *HorasHandler
	horas   *stubRegistroHorasRepository
	sprint  domain.Sprint
}

// nuevoEscenarioHorasHTTP siembra la historia 3 del proyecto 1 en un Sprint
// Activo que empezó hace 7 días y termina dentro de 7, y a Ana (usuario 2)
// como integrante. La historia 4 está en el Product Backlog.
func nuevoEscenarioHorasHTTP(t *testing.T) escenarioHorasHTTP {
	t.Helper()
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint, err := sprints.Guardar(context.Background(), domain.Sprint{
		ProyectoID: 1, SprintGoal: "g", Estado: domain.SprintActivo,
		FechaInicio: hoyCivil().AddDate(0, 0, -7), FechaFin: hoyCivil().AddDate(0, 0, 7),
	})
	if err != nil {
		t.Fatalf("no se pudo sembrar el Sprint: %v", err)
	}
	historias := &stubBacklogRepository{historias: map[int64]domain.BacklogItem{
		3: {ID: 3, ProyectoID: 1, Estado: domain.EstadoNueva, SprintID: &sprint.ID},
		4: {ID: 4, ProyectoID: 1, Estado: domain.EstadoNueva},
	}}
	integrantes := &stubProjectRepository{listaMiembros: []domain.Member{{UserID: 2, ProjectID: 1, Nombre: "Ana"}}}
	horas := &stubRegistroHorasRepository{}
	handler := NewHorasHandler(
		service.NewImputarHoras(historias, sprints, integrantes, horas),
		service.NewObtenerHistoria(historias),
	)
	return escenarioHorasHTTP{handler: handler, horas: horas, sprint: sprint}
}

func (e escenarioHorasHTTP) imputar(id, cuerpo string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(nethttp.MethodPost, "/backlog/"+id+"/horas", bytes.NewBufferString(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	rec := httptest.NewRecorder()
	e.handler.Imputar(rec, req)
	return rec
}

func TestHorasHandler_Imputar_Exitosa(t *testing.T) {
	e := nuevoEscenarioHorasHTTP(t)
	hoy := hoyCivil().Format(time.DateOnly)

	rec := e.imputar("3", `{"integrante_id": 2, "fecha": "`+hoy+`", "actividad": "Implementar", "horas": 4}`)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp registroHorasResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	esperada := registroHorasResponse{
		ID: 1, HistoriaID: 3, SprintID: e.sprint.ID, IntegranteID: 2, Integrante: "Ana",
		Fecha: hoy, Actividad: "Implementar", Minutos: 240, Horas: 4,
	}
	if resp != esperada {
		t.Errorf("respuesta:\n se esperaba %+v\n se obtuvo   %+v", esperada, resp)
	}
	if len(e.horas.registros) != 1 {
		t.Errorf("se esperaba 1 carga persistida, hay %d", len(e.horas.registros))
	}
}

func TestHorasHandler_Imputar_HorasConDecimales(t *testing.T) {
	e := nuevoEscenarioHorasHTTP(t)
	hoy := hoyCivil().Format(time.DateOnly)

	rec := e.imputar("3", `{"integrante_id": 2, "fecha": "`+hoy+`", "actividad": "Revisar", "horas": 1.5}`)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp registroHorasResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.Minutos != 90 || resp.Horas != 1.5 {
		t.Errorf("se esperaban 90 minutos y 1.5 horas, se obtuvo %d y %v", resp.Minutos, resp.Horas)
	}
}

func TestHorasHandler_Imputar_Errores(t *testing.T) {
	hoy := hoyCivil().Format(time.DateOnly)
	cuerpo := func(integranteID, actividad, horas string) string {
		return `{"integrante_id": ` + integranteID + `, "fecha": "` + hoy + `", "actividad": "` + actividad + `", "horas": ` + horas + `}`
	}
	casos := []struct {
		nombre, id, cuerpo string
		status             int
		campo, mensaje     string
	}{
		{"id no numérico", "abc", cuerpo("2", "x", "1"), nethttp.StatusBadRequest, "id", ""},
		{"id cero", "0", cuerpo("2", "x", "1"), nethttp.StatusBadRequest, "id", ""},
		{"historia inexistente con cuerpo inválido", "999", `no es json`, nethttp.StatusNotFound, "", "historia no encontrada"},
		{"cuerpo que no es JSON", "3", `no es json`, nethttp.StatusBadRequest, "", "JSON inválido"},
		{"horas como texto no numérico", "3", cuerpo("2", "x", `"abc"`), nethttp.StatusBadRequest, "horas", "las horas deben ser un número"},
		{"horas ausentes", "3", `{"integrante_id": 2, "fecha": "` + hoy + `", "actividad": "x"}`, nethttp.StatusBadRequest, "horas", "las horas son obligatorias"},
		{"horas nulas", "3", cuerpo("2", "x", "null"), nethttp.StatusBadRequest, "horas", "las horas son obligatorias"},
		{"horas negativas", "3", cuerpo("2", "x", "-2"), nethttp.StatusBadRequest, "horas", "las horas deben ser mayores que cero"},
		{"actividad vacía", "3", cuerpo("2", "", "1"), nethttp.StatusBadRequest, "actividad", "la actividad es obligatoria"},
		{"integrante ajeno", "3", cuerpo("99", "x", "1"), nethttp.StatusBadRequest, "integrante_id", "el integrante no pertenece al proyecto de la historia"},
		{"historia fuera del Sprint activo", "4", cuerpo("2", "x", "1"), nethttp.StatusBadRequest, "historia_id", "solo se imputan horas a historias del Sprint activo"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := nuevoEscenarioHorasHTTP(t)

			rec := e.imputar(c.id, c.cuerpo)

			if rec.Code != c.status {
				t.Fatalf("se esperaba %d, se obtuvo %d (%s)", c.status, rec.Code, rec.Body.String())
			}
			var resp errorResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("respuesta JSON inválida: %v", err)
			}
			if resp.Campo != c.campo || (c.mensaje != "" && resp.Mensaje != c.mensaje) {
				t.Errorf("se esperaba campo %q y mensaje %q, se obtuvo %+v", c.campo, c.mensaje, resp)
			}
			if len(e.horas.registros) != 0 {
				t.Errorf("no se debía registrar ninguna carga, hay %d", len(e.horas.registros))
			}
		})
	}
}

func TestHorasHandler_Imputar_HorasComoTextoNumerico(t *testing.T) {
	e := nuevoEscenarioHorasHTTP(t)
	hoy := hoyCivil().Format(time.DateOnly)

	rec := e.imputar("3", `{"integrante_id": 2, "fecha": "`+hoy+`", "actividad": "Implementar", "horas": "4"}`)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp registroHorasResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.Minutos != 240 {
		t.Errorf("se esperaban 240 minutos, se obtuvo %d", resp.Minutos)
	}
}
