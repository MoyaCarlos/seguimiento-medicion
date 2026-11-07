package http

import (
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type stubSprintsDelHistorial struct {
	sprints []domain.Sprint
	err     error
}

func (s stubSprintsDelHistorial) ListarPorProyecto(context.Context, int64) ([]domain.Sprint, error) {
	return s.sprints, s.err
}

type fuentesDelHistorial struct {
	horasPorSP *float64
	datos      stubDatosMetricas
	sprints    stubSprintsDelHistorial
}

func sprintFinalizadoDelHistorial(id int64, numero int, goal string) domain.Sprint {
	return domain.Sprint{
		ID: id, ProyectoID: 1, Numero: numero, SprintGoal: goal, Estado: domain.SprintFinalizado,
		FechaInicio: time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC),
		FechaFin:    time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC),
	}
}

func consultarHistorial(t *testing.T, f fuentesDelHistorial, id string) *httptest.ResponseRecorder {
	t.Helper()
	proyectos := repository.NewProjectRepositoryEnMemoria()
	proyecto, err := proyectos.Guardar(context.Background(), domain.Project{Nombre: "Proyecto", HorasPorStoryPoint: f.horasPorSP})
	if err != nil {
		t.Fatalf("no se pudo sembrar el proyecto: %v", err)
	}
	if proyecto.ID != 1 {
		t.Fatalf("se esperaba el proyecto 1, se obtuvo %d", proyecto.ID)
	}

	metricas := service.NewCalcularMetricas(proyectos, f.datos)
	req := httptest.NewRequest(nethttp.MethodGet, "/projects/"+id+"/historial", nil)
	req.SetPathValue("id", id)
	rec := httptest.NewRecorder()
	NewHistorialHandler(service.NewConsultarHistorialSprints(metricas, f.sprints), nil).Listar(rec, req)
	return rec
}

func cuerpoJSONLista(t *testing.T, rec *httptest.ResponseRecorder) []any {
	t.Helper()
	var cuerpo []any
	if err := json.Unmarshal(rec.Body.Bytes(), &cuerpo); err != nil {
		t.Fatalf("la respuesta no es un arreglo JSON: %v (%s)", err, rec.Body.String())
	}
	return cuerpo
}

func TestHistorialHandler_Listar_SprintFinalizadoConSusMetricas(t *testing.T) {
	horas := 4.0
	f := fuentesDelHistorial{
		horasPorSP: &horas,
		datos: stubDatosMetricas{datos: []domain.DatosMetricasSprint{{
			SprintID: 7, Numero: 1, HistoriasCompletadas: 1, HistoriasArrastradas: 1, SPCompletados: 5, SPArrastrados: 3,
			MinutosReales: 1800, DefectosDetectados: 2, DefectosResueltos: 1,
		}}},
		sprints: stubSprintsDelHistorial{sprints: []domain.Sprint{sprintFinalizadoDelHistorial(7, 1, "Primer incremento")}},
	}

	rec := consultarHistorial(t, f, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	esperado := []any{map[string]any{
		"id":           7.0,
		"proyecto_id":  1.0,
		"numero":       1.0,
		"sprint_goal":  "Primer incremento",
		"fecha_inicio": "2026-10-05",
		"fecha_fin":    "2026-10-19",
		"metricas": map[string]any{
			"story_points_planificados":        8.0,
			"story_points_completados":         5.0,
			"velocidad":                        5.0,
			"porcentaje_historias_completadas": 50.0,
			"horas_estimadas":                  32.0,
			"horas_reales":                     30.0,
			"desviacion_horas":                 -2.0,
			"desviacion_porcentaje":            -6.25,
			"defectos_detectados":              2.0,
			"defectos_resueltos":               1.0,
		},
	}}
	if cuerpo := cuerpoJSONLista(t, rec); !reflect.DeepEqual(cuerpo, esperado) {
		t.Errorf("respuesta inesperada:\n se obtuvo %#v\n se esperaba %#v", cuerpo, esperado)
	}
}

func TestHistorialHandler_Listar_SinHorasPorStoryPointInformaSinDatos(t *testing.T) {
	f := fuentesDelHistorial{
		datos:   stubDatosMetricas{datos: []domain.DatosMetricasSprint{{SprintID: 7, Numero: 1, HistoriasCompletadas: 1, SPCompletados: 5}}},
		sprints: stubSprintsDelHistorial{sprints: []domain.Sprint{sprintFinalizadoDelHistorial(7, 1, "Primer incremento")}},
	}

	rec := consultarHistorial(t, f, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	sprint := cuerpoJSONLista(t, rec)[0].(map[string]any)
	metricas := sprint["metricas"].(map[string]any)
	for _, campo := range []string{"horas_estimadas", "desviacion_horas", "desviacion_porcentaje"} {
		valor, presente := metricas[campo]
		if !presente || valor != nil {
			t.Errorf("%s: se esperaba null (sin datos), se obtuvo %#v (presente: %v)", campo, valor, presente)
		}
	}
}

func TestHistorialHandler_Listar_SinSprintsFinalizadosDevuelveArregloVacio(t *testing.T) {
	rec := consultarHistorial(t, fuentesDelHistorial{}, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if cuerpo := rec.Body.String(); cuerpo != "[]\n" && cuerpo != "[]" {
		t.Errorf("se esperaba [], se obtuvo %q", cuerpo)
	}
}

func TestHistorialHandler_Listar_Errores(t *testing.T) {
	casos := []struct {
		nombre  string
		id      string
		fuentes fuentesDelHistorial
		codigo  int
		cuerpo  errorResponse
	}{
		{"id inválido", "abc", fuentesDelHistorial{}, nethttp.StatusBadRequest,
			errorResponse{Campo: "id", Mensaje: "el identificador debe ser un entero positivo"}},
		{"proyecto inexistente", "99", fuentesDelHistorial{}, nethttp.StatusNotFound,
			errorResponse{Mensaje: "proyecto no encontrado"}},
		{"error de lectura", "1", fuentesDelHistorial{sprints: stubSprintsDelHistorial{err: errors.New("disco roto")}}, nethttp.StatusInternalServerError,
			errorResponse{Mensaje: "error interno"}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := consultarHistorial(t, c.fuentes, c.id)

			if rec.Code != c.codigo {
				t.Fatalf("se esperaba %d, se obtuvo %d (%s)", c.codigo, rec.Code, rec.Body.String())
			}
			var cuerpo errorResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &cuerpo); err != nil {
				t.Fatalf("la respuesta no es JSON: %v", err)
			}
			if cuerpo != c.cuerpo {
				t.Errorf("se esperaba %+v, se obtuvo %+v", c.cuerpo, cuerpo)
			}
		})
	}
}
