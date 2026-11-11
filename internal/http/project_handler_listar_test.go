package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestProjectHandler_Listar_ConProyectos(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	primero := crearProyectoViaHTTP(t, handler, "Primero")
	segundo := crearProyectoViaHTTPConFechas(t, handler, "Segundo", "2026-03-01", "2026-11-30")

	rec := listarComo(handler, anaDePrueba)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var respuesta []json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &respuesta); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if len(respuesta) != 2 {
		t.Fatalf("se esperaban 2 proyectos, se obtuvieron %d", len(respuesta))
	}

	var ids []int64
	for _, elemento := range respuesta {
		var proyecto proyectoResponse
		if err := json.Unmarshal(elemento, &proyecto); err != nil {
			t.Fatalf("elemento JSON inválido: %v", err)
		}
		ids = append(ids, proyecto.ID)
		if proyecto.Rol != "scrum_master" {
			t.Errorf("el proyecto %d debía traer el rol de quien consulta (scrum_master), se obtuvo %q", proyecto.ID, proyecto.Rol)
		}
	}
	if ids[0] != primero || ids[1] != segundo {
		t.Errorf("se esperaban los ids %d y %d, se obtuvo %v", primero, segundo, ids)
	}

	for _, elemento := range respuesta {
		var proyecto proyectoResponse
		if err := json.Unmarshal(elemento, &proyecto); err != nil {
			t.Fatalf("elemento JSON inválido: %v", err)
		}
		individual := obtenerProyectoViaHTTP(t, handler, proyecto.ID)
		if !mismoJSON(elemento, individual) {
			t.Errorf(
				"el proyecto %d del listado difiere de GET /projects/%d:\nlistado: %s\nindividual: %s",
				proyecto.ID, proyecto.ID, elemento, individual,
			)
		}
	}
}

func TestProjectHandler_Listar_Vacio(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()

	rec := listarComo(handler, anaDePrueba)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if body := strings.TrimSpace(rec.Body.String()); body != "[]" {
		t.Errorf("se esperaba un arreglo vacío, se obtuvo %q", body)
	}
}

func TestProjectHandler_Listar_ErrorDePersistencia(t *testing.T) {
	handler, repo := nuevoProjectHandlerDePrueba()
	repo.err = errPersistencia

	rec := listarComo(handler, anaDePrueba)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("se esperaba 500, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	verificarErrorSinCampo(t, rec)
}

// HU-21, escenarios 5 y 9: cada persona ve solo los proyectos que integra, con
// su rol en cada uno; quien no integra ninguno recibe [] y nunca 403.
func TestProjectHandler_Listar_SoloLosProyectosDeQuienConsultaConSuRol(t *testing.T) {
	handler, repo := nuevoProjectHandlerDePrueba()
	bruno := domain.User{ID: 8, Nombre: "bruno", NombreNormalizado: "bruno"}
	carla := domain.User{ID: 9, Nombre: "carla", NombreNormalizado: "carla"}
	deAna := crearProyectoViaHTTP(t, handler, "De Ana")
	deBruno := crearProyectoComo(t, handler, bruno, "De Bruno")
	repo.miembros = append(repo.miembros, domain.Membership{ProjectID: deAna, UserID: bruno.ID, Role: domain.RolProductBuilder})

	rec := listarComo(handler, anaDePrueba)
	if body := strings.TrimSpace(rec.Body.String()); rec.Code != http.StatusOK || strings.Count(body, `"id":`) != 1 ||
		!strings.Contains(body, `"nombre":"De Ana"`) || !strings.Contains(body, `"rol":"scrum_master"`) {
		t.Errorf("Ana debía ver solo su proyecto como scrum_master: %d %s", rec.Code, body)
	}

	rec = listarComo(handler, bruno)
	var deBrunoResp []proyectoResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &deBrunoResp); err != nil || rec.Code != http.StatusOK {
		t.Fatalf("respuesta inválida para Bruno: %d %s", rec.Code, rec.Body.String())
	}
	if len(deBrunoResp) != 2 ||
		deBrunoResp[0].ID != deAna || deBrunoResp[0].Rol != "product_builder" ||
		deBrunoResp[1].ID != deBruno || deBrunoResp[1].Rol != "scrum_master" {
		t.Errorf("Bruno debía ver %d como product_builder y %d como scrum_master, se obtuvo %+v", deAna, deBruno, deBrunoResp)
	}

	rec = listarComo(handler, carla)
	if body := strings.TrimSpace(rec.Body.String()); rec.Code != http.StatusOK || body != "[]" {
		t.Errorf("Carla no integra proyectos: se esperaba 200 [], se obtuvo %d %s", rec.Code, body)
	}
}

// HU-21, escenario 5: el panel trae el rol de quien consulta.
func TestProjectHandler_Obtener_TraeElRolDeQuienConsulta(t *testing.T) {
	handler, repo := nuevoProjectHandlerDePrueba()
	bruno := domain.User{ID: 8, Nombre: "bruno", NombreNormalizado: "bruno"}
	id := crearProyectoViaHTTP(t, handler, "De Ana")
	repo.miembros = append(repo.miembros, domain.Membership{ProjectID: id, UserID: bruno.ID, Role: domain.RolProductBuilder})

	for _, c := range []struct {
		quien domain.User
		rol   string
	}{{anaDePrueba, "scrum_master"}, {bruno, "product_builder"}} {
		idStr := strconv.FormatInt(id, 10)
		req := httptest.NewRequest(http.MethodGet, "/projects/"+idStr, nil)
		req.SetPathValue("id", idStr)
		rec := httptest.NewRecorder()
		handler.Obtener(rec, conSesionDe(req, c.quien))
		var resp proyectoResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil || rec.Code != http.StatusOK {
			t.Fatalf("respuesta inválida para %s: %d %s", c.quien.Nombre, rec.Code, rec.Body.String())
		}
		if resp.Rol != c.rol {
			t.Errorf("%s: se esperaba el rol %q, se obtuvo %q", c.quien.Nombre, c.rol, resp.Rol)
		}
	}
}

func listarComo(handler *ProjectHandler, u domain.User) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	handler.Listar(rec, conSesionDe(httptest.NewRequest(http.MethodGet, "/projects", nil), u))
	return rec
}

func crearProyectoComo(t *testing.T, handler *ProjectHandler, u domain.User, nombre string) int64 {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(`{"nombre":"`+nombre+`"}`))
	handler.Crear(rec, conSesionDe(req, u))
	var resp proyectoResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil || rec.Code != http.StatusCreated {
		t.Fatalf("no se pudo preparar el proyecto de %s: %d %s", u.Nombre, rec.Code, rec.Body.String())
	}
	return resp.ID
}

func crearProyectoViaHTTPConFechas(t *testing.T, handler *ProjectHandler, nombre, inicio, fin string) int64 {
	t.Helper()
	rec := httptest.NewRecorder()
	cuerpo := `{"nombre":"` + nombre + `","fecha_inicio":"` + inicio + `","fecha_fin":"` + fin + `"}`
	req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(cuerpo))
	handler.Crear(rec, conSesionDe(req, anaDePrueba))
	if rec.Code != http.StatusCreated {
		t.Fatalf("no se pudo preparar el proyecto: %d (%s)", rec.Code, rec.Body.String())
	}
	var resp proyectoResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta inválida al preparar el proyecto: %v", err)
	}
	return resp.ID
}

func obtenerProyectoViaHTTP(t *testing.T, handler *ProjectHandler, id int64) json.RawMessage {
	t.Helper()
	idStr := strconv.FormatInt(id, 10)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/projects/"+idStr, nil)
	req.SetPathValue("id", idStr)
	handler.Obtener(rec, conSesionDe(req, anaDePrueba))
	if rec.Code != http.StatusOK {
		t.Fatalf("no se pudo consultar el proyecto %d: %d (%s)", id, rec.Code, rec.Body.String())
	}
	var resp json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta inválida al consultar el proyecto %d: %v", id, err)
	}
	return resp
}

func mismoJSON(a, b json.RawMessage) bool {
	var izquierda, derecha any
	if err := json.Unmarshal(a, &izquierda); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &derecha); err != nil {
		return false
	}
	return reflect.DeepEqual(izquierda, derecha)
}
