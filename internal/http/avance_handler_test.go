package http

import (
	"context"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type stubAvanceHistorias struct {
	historias []domain.BacklogItem
	err       error
	llamadas  int
}

func (s *stubAvanceHistorias) ListarPorProyecto(context.Context, int64) ([]domain.BacklogItem, error) {
	s.llamadas++
	return s.historias, s.err
}

func consultarAvance(t *testing.T, historias *stubAvanceHistorias, id string) *httptest.ResponseRecorder {
	t.Helper()
	return consultarAvanceConSprints(t, historias, stubSprintsDelHistorial{}, id)
}

func consultarAvanceConSprints(t *testing.T, historias *stubAvanceHistorias, sprints stubSprintsDelHistorial, id string) *httptest.ResponseRecorder {
	t.Helper()
	proyectos := repository.NewProjectRepositoryEnMemoria()
	proyecto, err := proyectos.Guardar(context.Background(), domain.Project{Nombre: "Proyecto"})
	if err != nil {
		t.Fatalf("no se pudo sembrar el proyecto: %v", err)
	}
	if proyecto.ID != 1 {
		t.Fatalf("se esperaba el proyecto 1, se obtuvo %d", proyecto.ID)
	}

	req := httptest.NewRequest(nethttp.MethodGet, "/projects/"+id+"/avance", nil)
	req.SetPathValue("id", id)
	rec := httptest.NewRecorder()
	NewAvanceHandler(service.NewConsultarAvanceHistorias(proyectos, historias, sprints)).Consultar(rec, req)
	return rec
}

func historiaDelAvance(id int64, titulo string, estado domain.Estado, sp *int) domain.BacklogItem {
	return domain.BacklogItem{ID: id, ProyectoID: 1, Titulo: titulo, Estado: estado, EstimacionSP: sp}
}

func TestAvanceHandler_Consultar_Devuelve200ConLosTresGruposDelContrato(t *testing.T) {
	cinco, ocho := 5, 8
	repo := &stubAvanceHistorias{historias: []domain.BacklogItem{
		historiaDelAvance(1, "H1", domain.EstadoNueva, nil),
		historiaDelAvance(2, "H2", domain.EstadoEnProgreso, &cinco),
		historiaDelAvance(3, "H3", domain.EstadoCompletada, &ocho),
	}}

	rec := consultarAvance(t, repo, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	esperado := map[string]any{
		"total": map[string]any{"historias": 3.0, "story_points": 13.0},
		"grupos": []any{
			map[string]any{"estado": "Nueva", "cantidad": 1.0, "story_points": 0.0, "historias": []any{
				map[string]any{"id": 1.0, "titulo": "H1", "estimacion_sp": nil, "estado": "Nueva", "sprint": nil},
			}},
			map[string]any{"estado": "En progreso", "cantidad": 1.0, "story_points": 5.0, "historias": []any{
				map[string]any{"id": 2.0, "titulo": "H2", "estimacion_sp": 5.0, "estado": "En progreso", "sprint": nil},
			}},
			map[string]any{"estado": "Completada", "cantidad": 1.0, "story_points": 8.0, "historias": []any{
				map[string]any{"id": 3.0, "titulo": "H3", "estimacion_sp": 8.0, "estado": "Completada", "sprint": nil},
			}},
		},
	}
	if got := cuerpoJSON(t, rec); !reflect.DeepEqual(got, esperado) {
		t.Errorf("cuerpo distinto del contrato:\n got: %v\nwant: %v", got, esperado)
	}
	if tipo := rec.Header().Get("Content-Type"); tipo != "application/json" {
		t.Errorf("se esperaba Content-Type application/json, se obtuvo %q", tipo)
	}
}

func TestAvanceHandler_Consultar_LosGruposVaciosTienenHistoriasComoArregloVacio(t *testing.T) {
	rec := consultarAvance(t, &stubAvanceHistorias{}, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	grupos, ok := cuerpoJSON(t, rec)["grupos"].([]any)
	if !ok || len(grupos) != 3 {
		t.Fatalf("se esperaban 3 grupos, se obtuvo %v", cuerpoJSON(t, rec)["grupos"])
	}
	for _, g := range grupos {
		historias, ok := g.(map[string]any)["historias"].([]any)
		if !ok || len(historias) != 0 {
			t.Errorf("se esperaba historias como [] (nunca null), se obtuvo %v", g)
		}
	}
}

func TestAvanceHandler_Consultar_IdInvalidoEs400ConCampoIdYNoLeeHistorias(t *testing.T) {
	for _, id := range []string{"abc", "0", "-1"} {
		t.Run(id, func(t *testing.T) {
			repo := &stubAvanceHistorias{}

			rec := consultarAvance(t, repo, id)

			if rec.Code != nethttp.StatusBadRequest {
				t.Fatalf("se esperaba 400, se obtuvo %d (%s)", rec.Code, rec.Body.String())
			}
			cuerpo := cuerpoJSON(t, rec)
			if cuerpo["campo"] != "id" || cuerpo["mensaje"] != "el identificador debe ser un entero positivo" {
				t.Errorf("se esperaba el error de parsearID sobre el campo id, se obtuvo %v", cuerpo)
			}
			if repo.llamadas != 0 {
				t.Errorf("no se esperaba leer historias con un id inválido")
			}
		})
	}
}

func TestAvanceHandler_Consultar_ProyectoInexistenteEs404(t *testing.T) {
	repo := &stubAvanceHistorias{}

	rec := consultarAvance(t, repo, "999")

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("se esperaba 404, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if mensaje := cuerpoJSON(t, rec)["mensaje"]; mensaje != "proyecto no encontrado" {
		t.Errorf("se esperaba %q, se obtuvo %v", "proyecto no encontrado", mensaje)
	}
	if repo.llamadas != 0 {
		t.Errorf("no se esperaba leer historias de un proyecto inexistente")
	}
}

func TestAvanceHandler_Consultar_ErrorDeLecturaEs500SinDatosParciales(t *testing.T) {
	rec := consultarAvance(t, &stubAvanceHistorias{err: errors.New("falló la lectura")}, "1")

	if rec.Code != nethttp.StatusInternalServerError {
		t.Fatalf("se esperaba 500, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	cuerpo := cuerpoJSON(t, rec)
	if cuerpo["mensaje"] != "error interno" {
		t.Errorf("se esperaba %q, se obtuvo %v", "error interno", cuerpo["mensaje"])
	}
	if _, hay := cuerpo["grupos"]; hay {
		t.Errorf("no se esperaban datos parciales, se obtuvo %v", cuerpo)
	}
}

func TestAvanceHandler_Consultar_InformaElSprintDeCadaHistoriaONullSiEstaEnElProductBacklog(t *testing.T) {
	sprintID := int64(20)
	conSprint := historiaDelAvance(2, "H2", domain.EstadoEnProgreso, nil)
	conSprint.SprintID = &sprintID
	repo := &stubAvanceHistorias{historias: []domain.BacklogItem{
		historiaDelAvance(1, "H1", domain.EstadoNueva, nil),
		conSprint,
	}}
	sprints := stubSprintsDelHistorial{sprints: []domain.Sprint{{ID: 20, ProyectoID: 1, Numero: 2, Estado: domain.SprintActivo}}}

	rec := consultarAvanceConSprints(t, repo, sprints, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	grupos := cuerpoJSON(t, rec)["grupos"].([]any)
	nueva := grupos[0].(map[string]any)["historias"].([]any)[0].(map[string]any)
	enProgreso := grupos[1].(map[string]any)["historias"].([]any)[0].(map[string]any)
	if nueva["sprint"] != nil {
		t.Errorf("la historia sin Sprint debe tener sprint null, se obtuvo %v", nueva["sprint"])
	}
	if esperado := map[string]any{"id": 20.0, "numero": 2.0}; !reflect.DeepEqual(enProgreso["sprint"], esperado) {
		t.Errorf("se esperaba sprint %v, se obtuvo %v", esperado, enProgreso["sprint"])
	}
}

func TestAvanceHandler_Consultar_ProyectoSinHistoriasResponde200ConTotalesEnCero(t *testing.T) {
	rec := consultarAvance(t, &stubAvanceHistorias{}, "1")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	cuerpo := cuerpoJSON(t, rec)
	if esperado := map[string]any{"historias": 0.0, "story_points": 0.0}; !reflect.DeepEqual(cuerpo["total"], esperado) {
		t.Errorf("se esperaba total %v, se obtuvo %v", esperado, cuerpo["total"])
	}
	for i, estado := range []string{"Nueva", "En progreso", "Completada"} {
		esperado := map[string]any{"estado": estado, "cantidad": 0.0, "story_points": 0.0, "historias": []any{}}
		if got := cuerpo["grupos"].([]any)[i]; !reflect.DeepEqual(got, esperado) {
			t.Errorf("grupo %d: se esperaba %v, se obtuvo %v", i, esperado, got)
		}
	}
}
