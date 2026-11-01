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

type stubDefectoRepository struct {
	defectos []domain.Defecto
}

func (s *stubDefectoRepository) Registrar(_ context.Context, d domain.Defecto) (domain.Defecto, error) {
	d.ID = int64(len(s.defectos) + 1)
	d.CreadoEn = time.Date(2026, 10, 30, 14, 5, 0, 0, time.UTC)
	s.defectos = append(s.defectos, d)
	return d, nil
}

func (s *stubDefectoRepository) ObtenerPorID(_ context.Context, id int64) (domain.DefectoDetalle, error) {
	for _, d := range s.defectos {
		if d.ID == id {
			return domain.DefectoDetalle{Defecto: d, TituloHistoria: "Carrito de compras"}, nil
		}
	}
	return domain.DefectoDetalle{}, domain.ErrDefectoNoEncontrado
}

func (s *stubDefectoRepository) Resolver(_ context.Context, d domain.Defecto) error {
	for i := range s.defectos {
		if s.defectos[i].ID == d.ID {
			s.defectos[i] = d
			return nil
		}
	}
	return domain.ErrDefectoNoEncontrado
}

func (s *stubDefectoRepository) ListarPorProyecto(_ context.Context, proyectoID int64) ([]domain.DefectoDetalle, error) {
	detalle := make([]domain.DefectoDetalle, 0)
	for _, d := range s.defectos {
		if d.ProyectoID == proyectoID {
			detalle = append(detalle, domain.DefectoDetalle{Defecto: d, TituloHistoria: "Carrito de compras"})
		}
	}
	return detalle, nil
}

type escenarioDefectosHTTP struct {
	handler  *DefectoHandler
	sprints  *repository.SprintRepositoryEnMemoria
	defectos *stubDefectoRepository
}

// nuevoEscenarioDefectosHTTP siembra los proyectos 1 y 2 y la historia 3
// "Carrito de compras" del proyecto 1. Los Sprints los agrega cada test.
func nuevoEscenarioDefectosHTTP() escenarioDefectosHTTP {
	sprints := repository.NewSprintRepositoryEnMemoria()
	proyectos := &stubProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1}, 2: {ID: 2}}}
	historias := &stubBacklogRepository{historias: map[int64]domain.BacklogItem{
		3: {ID: 3, ProyectoID: 1, Titulo: "Carrito de compras", Estado: domain.EstadoNueva},
	}}
	defectos := &stubDefectoRepository{}
	handler := NewDefectoHandler(
		service.NewObtenerProyecto(proyectos),
		service.NewRegistrarDefecto(sprints, historias, defectos),
	)
	return escenarioDefectosHTTP{handler: handler, sprints: sprints, defectos: defectos}
}

func (e escenarioDefectosHTTP) sprint(t *testing.T, proyectoID int64, estado domain.EstadoSprint) domain.Sprint {
	t.Helper()
	s, err := e.sprints.Guardar(context.Background(), domain.Sprint{ProyectoID: proyectoID, SprintGoal: "g", Estado: estado})
	if err != nil {
		t.Fatalf("no se pudo sembrar el Sprint: %v", err)
	}
	return s
}

func (e escenarioDefectosHTTP) registrar(id, cuerpo string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(nethttp.MethodPost, "/projects/"+id+"/defectos", bytes.NewBufferString(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	rec := httptest.NewRecorder()
	e.handler.Registrar(rec, req)
	return rec
}

func leerDefecto(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v (%s)", err, rec.Body.String())
	}
	return resp
}

const cuerpoDefecto = `{"historia_id": 3, "descripcion": "El total no se actualiza", "severidad": "Alta"}`

func TestDefectoHandler_Registrar_Exitoso(t *testing.T) {
	e := nuevoEscenarioDefectosHTTP()
	activo := e.sprint(t, 1, domain.SprintActivo)

	rec := e.registrar("1", cuerpoDefecto)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	esperada := map[string]any{
		"id": 1.0, "proyecto_id": 1.0, "historia_id": 3.0, "historia": "Carrito de compras",
		"descripcion": "El total no se actualiza", "severidad": "Alta", "estado": "Pendiente",
		"sprint_deteccion_id": float64(activo.ID), "sprint_resolucion_id": nil, "creado_en": "2026-10-30T14:05:00Z",
	}
	resp := leerDefecto(t, rec)
	if len(resp) != len(esperada) {
		t.Errorf("se esperaban los campos %v, se obtuvo %v", esperada, resp)
	}
	for campo, valor := range esperada {
		if got, ok := resp[campo]; !ok || got != valor {
			t.Errorf("%s: se esperaba %v, se obtuvo %v (presente: %v)", campo, valor, got, ok)
		}
	}
}

func TestDefectoHandler_Registrar_IgnoraEstadoYSprints(t *testing.T) {
	e := nuevoEscenarioDefectosHTTP()
	activo := e.sprint(t, 1, domain.SprintActivo)

	rec := e.registrar("1", `{"historia_id": 3, "descripcion": "d", "severidad": "Baja", "estado": "Resuelto", "sprint_deteccion_id": 99, "sprint_resolucion_id": 99}`)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	resp := leerDefecto(t, rec)
	if resp["estado"] != "Pendiente" || resp["sprint_deteccion_id"] != float64(activo.ID) || resp["sprint_resolucion_id"] != nil {
		t.Errorf("se esperaba Pendiente en el Sprint %d sin resolución, se obtuvo %v", activo.ID, resp)
	}
}

func TestDefectoHandler_Registrar_Errores(t *testing.T) {
	casos := []struct {
		nombre, id, cuerpo string
		conActivo          bool
		status             int
		campo, mensaje     string
	}{
		{"id no numérico", "abc", cuerpoDefecto, true, nethttp.StatusBadRequest, "id", ""},
		{"id cero", "0", cuerpoDefecto, true, nethttp.StatusBadRequest, "id", ""},
		{"proyecto inexistente con cuerpo inválido", "999", `no es json`, true, nethttp.StatusNotFound, "", "proyecto no encontrado"},
		{"cuerpo que no es JSON", "1", `no es json`, true, nethttp.StatusBadRequest, "", "JSON inválido"},
		{"sin Sprint activo", "1", cuerpoDefecto, false, nethttp.StatusBadRequest, "sprint_deteccion_id", ""},
		{"sin historia", "1", `{"descripcion": "d", "severidad": "Alta"}`, true, nethttp.StatusBadRequest, "historia_id", ""},
		{"descripción vacía", "1", `{"historia_id": 3, "descripcion": "", "severidad": "Alta"}`, true, nethttp.StatusBadRequest, "descripcion", ""},
		{"severidad en minúsculas", "1", `{"historia_id": 3, "descripcion": "d", "severidad": "alta"}`, true, nethttp.StatusBadRequest, "severidad", ""},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := nuevoEscenarioDefectosHTTP()
			if c.conActivo {
				e.sprint(t, 1, domain.SprintActivo)
			}

			rec := e.registrar(c.id, c.cuerpo)

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
			if len(e.defectos.defectos) != 0 {
				t.Errorf("no se esperaba ningún defecto registrado, hay %d", len(e.defectos.defectos))
			}
		})
	}
}
