package http

import (
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type stubSprintPorID struct {
	sprints map[int64]domain.Sprint
}

func (s stubSprintPorID) ObtenerPorID(_ context.Context, id int64) (domain.Sprint, error) {
	sprint, ok := s.sprints[id]
	if !ok {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	return sprint, nil
}

type stubHistorialSprint struct {
	historias []domain.HistoriaComprometida
	minutos   map[int64]int
	defectos  []domain.DefectoDetalle
	err       error
}

func (s stubHistorialSprint) HistoriasComprometidas(context.Context, int64) ([]domain.HistoriaComprometida, error) {
	return s.historias, s.err
}

func (s stubHistorialSprint) MinutosPorHistoria(context.Context, int64) (map[int64]int, error) {
	return s.minutos, s.err
}

func (s stubHistorialSprint) DefectosDeSprint(context.Context, int64) ([]domain.DefectoDetalle, error) {
	return s.defectos, s.err
}

var sprintsDelDetalle = stubSprintPorID{sprints: map[int64]domain.Sprint{
	7: sprintFinalizadoDelHistorial(7, 1, "Primer incremento"),
	8: {ID: 8, ProyectoID: 1, Numero: 2, SprintGoal: "En curso", Estado: domain.SprintActivo},
	9: {ID: 9, ProyectoID: 2, Numero: 1, SprintGoal: "Ajeno", Estado: domain.SprintFinalizado},
}}

var datosDelDetalle = stubDatosMetricas{datos: []domain.DatosMetricasSprint{{
	SprintID: 7, Numero: 1, HistoriasCompletadas: 1, HistoriasArrastradas: 1, SPCompletados: 5, SPArrastrados: 3, MinutosReales: 140,
	DefectosDetectados: 1, DefectosResueltos: 1,
}}}

func consultarDetalle(t *testing.T, historial stubHistorialSprint, id, sprintID string) *httptest.ResponseRecorder {
	t.Helper()
	proyectos := repository.NewProjectRepositoryEnMemoria()
	horas := 4.0
	if _, err := proyectos.Guardar(context.Background(), domain.Project{Nombre: "Proyecto", HorasPorStoryPoint: &horas}); err != nil {
		t.Fatalf("no se pudo sembrar el proyecto: %v", err)
	}

	metricas := service.NewCalcularMetricas(proyectos, datosDelDetalle)
	detalle := service.NewConsultarDetalleSprint(proyectos, sprintsDelDetalle, metricas, historial)
	req := httptest.NewRequest(nethttp.MethodGet, "/projects/"+id+"/historial/"+sprintID, nil)
	req.SetPathValue("id", id)
	req.SetPathValue("sprintId", sprintID)
	rec := httptest.NewRecorder()
	NewHistorialHandler(nil, detalle).Detalle(rec, req)
	return rec
}

func TestHistorialHandler_Detalle_SprintFinalizado(t *testing.T) {
	cinco, tres := 5, 3
	resuelto := int64(7)
	historial := stubHistorialSprint{
		historias: []domain.HistoriaComprometida{
			{ID: 1, Titulo: "Login", EstimacionSP: &cinco, Completada: true},
			{ID: 2, Titulo: "Reportes", EstimacionSP: &tres},
			{ID: 3, Titulo: "Sin estimar"},
		},
		minutos: map[int64]int{1: 90, 2: 50},
		defectos: []domain.DefectoDetalle{{
			Defecto: domain.Defecto{ID: 4, ProyectoID: 1, HistoriaID: 1, SprintDeteccionID: 7, SprintResolucionID: &resuelto,
				Descripcion: "Error al guardar", Severidad: domain.Severidad("Alta"), Estado: domain.EstadoDefecto("Resuelto")},
			TituloHistoria: "Login",
		}},
	}

	rec := consultarDetalle(t, historial, "1", "7")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	cuerpo := cuerpoJSON(t, rec)
	for campo, valor := range map[string]any{
		"id": 7.0, "proyecto_id": 1.0, "numero": 1.0, "sprint_goal": "Primer incremento",
		"fecha_inicio": "2026-10-05", "fecha_fin": "2026-10-19",
	} {
		if cuerpo[campo] != valor {
			t.Errorf("%s: se esperaba %#v, se obtuvo %#v", campo, valor, cuerpo[campo])
		}
	}
	if metricas := cuerpo["metricas"].(map[string]any); metricas["horas_reales"] != 2.33 || metricas["story_points_planificados"] != 8.0 {
		t.Errorf("se esperaban las métricas de HU-03 del Sprint, se obtuvo %#v", metricas)
	}
	historias := []any{
		map[string]any{"id": 1.0, "titulo": "Login", "estimacion_sp": 5.0, "completada": true, "minutos": 90.0, "horas": 1.5},
		map[string]any{"id": 2.0, "titulo": "Reportes", "estimacion_sp": 3.0, "completada": false, "minutos": 50.0, "horas": 0.83},
		map[string]any{"id": 3.0, "titulo": "Sin estimar", "estimacion_sp": nil, "completada": false, "minutos": 0.0, "horas": 0.0},
	}
	if !reflect.DeepEqual(cuerpo["historias"], historias) {
		t.Errorf("historias inesperadas:\n se obtuvo %#v\n se esperaba %#v", cuerpo["historias"], historias)
	}
	defectos := []any{map[string]any{
		"id": 4.0, "historia_id": 1.0, "historia": "Login", "descripcion": "Error al guardar", "severidad": "Alta",
		"estado": "Resuelto", "detectado_en_sprint": true, "resuelto_en_sprint": true,
	}}
	if !reflect.DeepEqual(cuerpo["defectos"], defectos) {
		t.Errorf("defectos inesperados:\n se obtuvo %#v\n se esperaba %#v", cuerpo["defectos"], defectos)
	}
}

func TestHistorialHandler_Detalle_SinHistoriasNiDefectosDevuelveArreglosVacios(t *testing.T) {
	rec := consultarDetalle(t, stubHistorialSprint{}, "1", "7")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	cuerpo := cuerpoJSON(t, rec)
	for _, campo := range []string{"historias", "defectos"} {
		if lista, ok := cuerpo[campo].([]any); !ok || len(lista) != 0 {
			t.Errorf("%s: se esperaba [], se obtuvo %#v", campo, cuerpo[campo])
		}
	}
}

func TestHistorialHandler_Detalle_Errores(t *testing.T) {
	casos := []struct {
		nombre       string
		id, sprintID string
		historial    stubHistorialSprint
		codigo       int
		cuerpo       errorResponse
	}{
		{"id inválido", "abc", "7", stubHistorialSprint{}, nethttp.StatusBadRequest,
			errorResponse{Campo: "id", Mensaje: "el identificador debe ser un entero positivo"}},
		{"sprintId inválido", "1", "abc", stubHistorialSprint{}, nethttp.StatusBadRequest,
			errorResponse{Campo: "sprint_id", Mensaje: "el identificador debe ser un entero positivo"}},
		{"proyecto inexistente", "99", "7", stubHistorialSprint{}, nethttp.StatusNotFound,
			errorResponse{Mensaje: "proyecto no encontrado"}},
		{"Sprint inexistente", "1", "999", stubHistorialSprint{}, nethttp.StatusNotFound,
			errorResponse{Mensaje: "sprint no encontrado"}},
		{"Sprint de otro proyecto", "1", "9", stubHistorialSprint{}, nethttp.StatusNotFound,
			errorResponse{Mensaje: "sprint no encontrado"}},
		{"Sprint no finalizado", "1", "8", stubHistorialSprint{}, nethttp.StatusBadRequest,
			errorResponse{Campo: "estado", Mensaje: "el Sprint todavía no está finalizado"}},
		{"error de lectura", "1", "7", stubHistorialSprint{err: errors.New("disco roto")}, nethttp.StatusInternalServerError,
			errorResponse{Mensaje: "error interno"}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := consultarDetalle(t, c.historial, c.id, c.sprintID)

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
