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

type stubDatosMetricas struct {
	datos []domain.DatosMetricasSprint
	err   error
}

func (s stubDatosMetricas) ListarDatosSprintsFinalizados(context.Context, int64) ([]domain.DatosMetricasSprint, error) {
	return s.datos, s.err
}

func consultarMetricas(t *testing.T, horas *float64, datos stubDatosMetricas, id string) *httptest.ResponseRecorder {
	t.Helper()
	proyectos := repository.NewProjectRepositoryEnMemoria()
	proyecto, err := proyectos.Guardar(context.Background(), domain.Project{Nombre: "Proyecto", HorasPorStoryPoint: horas})
	if err != nil {
		t.Fatalf("no se pudo sembrar el proyecto: %v", err)
	}
	if proyecto.ID != 1 {
		t.Fatalf("se esperaba el proyecto 1, se obtuvo %d", proyecto.ID)
	}

	req := httptest.NewRequest(nethttp.MethodGet, "/projects/"+id+"/metricas", nil)
	req.SetPathValue("id", id)
	rec := httptest.NewRecorder()
	NewMetricasHandler(service.NewCalcularMetricas(proyectos, datos)).Consultar(rec, req)
	return rec
}

func cuerpoJSON(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var cuerpo map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &cuerpo); err != nil {
		t.Fatalf("la respuesta no es JSON: %v (%s)", err, rec.Body.String())
	}
	return cuerpo
}

func TestMetricasHandler_Consultar_ConUnSprint(t *testing.T) {
	horas := 4.0
	datos := stubDatosMetricas{datos: []domain.DatosMetricasSprint{{
		SprintID: 7, Numero: 2, HistoriasCompletadas: 1, HistoriasArrastradas: 1, SPCompletados: 5, SPArrastrados: 3,
		MinutosReales: 1800, DefectosDetectados: 2, DefectosResueltos: 1,
	}}}

	rec := consultarMetricas(t, &horas, datos, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	metricas := map[string]any{
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
	}
	esperado := map[string]any{
		"proyecto_id":           1.0,
		"horas_por_story_point": 4.0,
		"proyecto":              metricas,
		"sprints":               []any{map[string]any{"sprint_id": 7.0, "numero": 2.0, "metricas": metricas}},
	}
	if cuerpo := cuerpoJSON(t, rec); !reflect.DeepEqual(cuerpo, esperado) {
		t.Errorf("respuesta inesperada:\n%v\nse esperaba:\n%v", cuerpo, esperado)
	}
}

func TestMetricasHandler_Consultar_SinHorasPorStoryPoint(t *testing.T) {
	datos := stubDatosMetricas{datos: []domain.DatosMetricasSprint{{SprintID: 7, Numero: 1, HistoriasCompletadas: 1, SPCompletados: 5, MinutosReales: 60}}}

	rec := consultarMetricas(t, nil, datos, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	cuerpo := cuerpoJSON(t, rec)
	if valor, existe := cuerpo["horas_por_story_point"]; !existe || valor != nil {
		t.Errorf("horas_por_story_point: se esperaba null, se obtuvo %v", valor)
	}
	proyecto, _ := cuerpo["proyecto"].(map[string]any)
	for _, campo := range []string{"horas_estimadas", "desviacion_horas", "desviacion_porcentaje"} {
		if valor, existe := proyecto[campo]; !existe || valor != nil {
			t.Errorf("%s: se esperaba null, se obtuvo %v (existe: %v)", campo, valor, existe)
		}
	}
	if proyecto["horas_reales"] != 1.0 {
		t.Errorf("horas_reales: se esperaba 1, se obtuvo %v", proyecto["horas_reales"])
	}
}

func TestMetricasHandler_Consultar_SinSprintsFinalizados(t *testing.T) {
	horas := 4.0

	rec := consultarMetricas(t, &horas, stubDatosMetricas{datos: []domain.DatosMetricasSprint{}}, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	cuerpo := cuerpoJSON(t, rec)
	if sprints, esLista := cuerpo["sprints"].([]any); !esLista || len(sprints) != 0 {
		t.Errorf("se esperaba una lista vacía de Sprints, se obtuvo %v", cuerpo["sprints"])
	}
	proyecto, _ := cuerpo["proyecto"].(map[string]any)
	if proyecto["velocidad"] != 0.0 || proyecto["horas_estimadas"] != 0.0 {
		t.Errorf("se esperaban velocidad y horas estimadas en cero, se obtuvo %v", proyecto)
	}
	if valor, existe := proyecto["desviacion_porcentaje"]; !existe || valor != nil {
		t.Errorf("desviacion_porcentaje: se esperaba null, se obtuvo %v", valor)
	}
}

func TestMetricasHandler_Consultar_Errores(t *testing.T) {
	casos := []struct {
		nombre  string
		id      string
		datos   stubDatosMetricas
		status  int
		mensaje string
	}{
		{nombre: "id inválido", id: "abc", status: nethttp.StatusBadRequest},
		{nombre: "proyecto inexistente", id: "99", status: nethttp.StatusNotFound, mensaje: domain.ErrProyectoNoEncontrado.Error()},
		{nombre: "falla del repositorio", id: "1", datos: stubDatosMetricas{err: errors.New("falla")}, status: nethttp.StatusInternalServerError, mensaje: "error interno"},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := consultarMetricas(t, nil, c.datos, c.id)

			if rec.Code != c.status {
				t.Fatalf("se esperaba %d, se obtuvo %d (%s)", c.status, rec.Code, rec.Body.String())
			}
			if c.mensaje != "" && leerError(t, rec).Mensaje != c.mensaje {
				t.Errorf("se esperaba el mensaje %q, se obtuvo %s", c.mensaje, rec.Body.String())
			}
		})
	}
}

func TestMetricasHandler_Consultar_ConHorasExtremasRespondeJSONValido(t *testing.T) {
	horas := 1e308
	datos := stubDatosMetricas{datos: []domain.DatosMetricasSprint{{SprintID: 7, Numero: 1, HistoriasCompletadas: 1, SPCompletados: 16, MinutosReales: 60}}}

	rec := consultarMetricas(t, &horas, datos, "1")

	proyecto, _ := cuerpoJSON(t, rec)["proyecto"].(map[string]any)
	if valor, existe := proyecto["horas_estimadas"]; !existe || valor != nil {
		t.Errorf("horas_estimadas: se esperaba null, se obtuvo %v", valor)
	}
	if proyecto["horas_reales"] != 1.0 {
		t.Errorf("horas_reales: se esperaba 1, se obtuvo %v", proyecto["horas_reales"])
	}
}
