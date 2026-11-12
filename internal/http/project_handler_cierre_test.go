package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// relojDeCierreHTTP fija el «hoy» de los tests de cierre (spec 044, #300).
func relojDeCierreHTTP() time.Time { return time.Date(2026, 11, 11, 15, 30, 0, 0, time.Local) }

const productBuilderHTTP int64 = 20

type proyectoConEstado struct {
	proyectoResponse
	Estado string `json:"estado"`
}

// proyectoDeAnaConBuilder crea un proyecto de ana (Scrum Master) con Jimena
// (id 20) como Product Builder.
func proyectoDeAnaConBuilder(t *testing.T, handler *ProjectHandler, proyectos *stubProjectRepository) string {
	t.Helper()
	id := crearProyectoViaHTTP(t, handler, "Tienda")
	proyectos.miembros = append(proyectos.miembros, domain.Membership{ProjectID: id, UserID: productBuilderHTTP, Role: domain.RolProductBuilder})
	return strconv.FormatInt(id, 10)
}

func cerrarProyectoHTTP(handler *ProjectHandler, idStr string, actor domain.User) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/projects/"+idStr+"/cerrar", nil)
	req.SetPathValue("id", idStr)
	handler.Cerrar(rec, conSesionDe(req, actor))
	return rec
}

func TestProjectHandler_Cerrar_ElScrumMasterCierra(t *testing.T) {
	handler, proyectos := nuevoProjectHandlerDePrueba()
	idStr := proyectoDeAnaConBuilder(t, handler, proyectos)

	rec := cerrarProyectoHTTP(handler, idStr, anaDePrueba)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp proyectoConEstado
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.FechaCierre == nil || *resp.FechaCierre != "2026-11-11" {
		t.Errorf("se esperaba fecha_cierre 2026-11-11, se obtuvo %v", resp.FechaCierre)
	}
	if resp.Estado != string(domain.ProyectoFinalizado) {
		t.Errorf("se esperaba el estado %q, se obtuvo %q", domain.ProyectoFinalizado, resp.Estado)
	}
}

func TestProjectHandler_Cerrar_Errores(t *testing.T) {
	casos := []struct {
		nombre          string
		yaCerrado       bool
		conSprintActivo bool
		id              func(idStr string) string
		actor           domain.User
		codigo          int
		mensaje         string
	}{
		{
			nombre: "ya cerrado", yaCerrado: true,
			actor: anaDePrueba, codigo: http.StatusConflict, mensaje: "El proyecto ya está cerrado",
		},
		{
			nombre: "con un Sprint activo", conSprintActivo: true,
			actor: anaDePrueba, codigo: http.StatusConflict, mensaje: "Cerrá el Sprint activo antes de cerrar el proyecto",
		},
		{
			nombre: "Product Builder",
			actor:  domain.User{ID: productBuilderHTTP}, codigo: http.StatusForbidden, mensaje: "Solo el Scrum Master puede realizar esta acción",
		},
		{
			nombre: "proyecto inexistente",
			id:     func(string) string { return "999" },
			actor:  anaDePrueba, codigo: http.StatusNotFound, mensaje: "proyecto no encontrado",
		},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			handler, proyectos, sprints := nuevoProjectHandlerDePruebaConSprints()
			idStr := proyectoDeAnaConBuilder(t, handler, proyectos)
			if c.conSprintActivo {
				id, _ := strconv.ParseInt(idStr, 10, 64)
				if _, err := sprints.Guardar(context.Background(), domain.Sprint{ProyectoID: id, Estado: domain.SprintActivo}); err != nil {
					t.Fatalf("no se pudo sembrar el Sprint: %v", err)
				}
			}
			if c.yaCerrado {
				if rec := cerrarProyectoHTTP(handler, idStr, anaDePrueba); rec.Code != http.StatusOK {
					t.Fatalf("no se pudo cerrar de antemano: %d", rec.Code)
				}
			}
			objetivo := idStr
			if c.id != nil {
				objetivo = c.id(idStr)
			}

			rec := cerrarProyectoHTTP(handler, objetivo, c.actor)

			if rec.Code != c.codigo {
				t.Fatalf("se esperaba %d, se obtuvo %d (%s)", c.codigo, rec.Code, rec.Body.String())
			}
			var resp errorResponse
			_ = json.Unmarshal(rec.Body.Bytes(), &resp)
			if resp.Mensaje != c.mensaje {
				t.Errorf("se esperaba el mensaje %q, se obtuvo %q", c.mensaje, resp.Mensaje)
			}
		})
	}
}

// La interfaz sabe si el proyecto está en solo lectura por fecha_cierre (FR-008).
func TestProjectHandler_ConsultasIncluyenFechaDeCierre(t *testing.T) {
	handler, proyectos := nuevoProjectHandlerDePrueba()
	idStr := proyectoDeAnaConBuilder(t, handler, proyectos)

	obtener := func() map[string]any {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/projects/"+idStr, nil)
		req.SetPathValue("id", idStr)
		handler.Obtener(rec, conSesionDe(req, anaDePrueba))
		var resp map[string]any
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		return resp
	}
	listar := func() []map[string]any {
		rec := httptest.NewRecorder()
		handler.Listar(rec, conSesionDe(httptest.NewRequest(http.MethodGet, "/projects", nil), anaDePrueba))
		var resp []map[string]any
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		return resp
	}

	if valor, ok := obtener()["fecha_cierre"]; !ok || valor != nil {
		t.Errorf("abierto: se esperaba \"fecha_cierre\": null en el detalle, se obtuvo %v (presente: %v)", valor, ok)
	}
	if lista := listar(); len(lista) != 1 || lista[0]["fecha_cierre"] != nil {
		t.Errorf("abierto: se esperaba fecha_cierre null en Mis proyectos, se obtuvo %v", lista)
	} else if _, ok := lista[0]["fecha_cierre"]; !ok {
		t.Error("abierto: Mis proyectos debe incluir la clave fecha_cierre")
	}

	cerrarProyectoHTTP(handler, idStr, anaDePrueba)

	if valor := obtener()["fecha_cierre"]; valor != "2026-11-11" {
		t.Errorf("cerrado: se esperaba fecha_cierre 2026-11-11 en el detalle, se obtuvo %v", valor)
	}
	if lista := listar(); len(lista) != 1 || lista[0]["fecha_cierre"] != "2026-11-11" {
		t.Errorf("cerrado: se esperaba fecha_cierre 2026-11-11 en Mis proyectos, se obtuvo %v", lista)
	}
}

func reabrirProyectoHTTP(handler *ProjectHandler, idStr string, actor domain.User) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/projects/"+idStr+"/reabrir", nil)
	req.SetPathValue("id", idStr)
	handler.Reabrir(rec, conSesionDe(req, actor))
	return rec
}

func TestProjectHandler_Reabrir_ElScrumMasterReabre(t *testing.T) {
	handler, proyectos := nuevoProjectHandlerDePrueba()
	idStr := proyectoDeAnaConBuilder(t, handler, proyectos)
	cerrarProyectoHTTP(handler, idStr, anaDePrueba)

	rec := reabrirProyectoHTTP(handler, idStr, anaDePrueba)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if valor, ok := resp["fecha_cierre"]; !ok || valor != nil {
		t.Errorf("se esperaba \"fecha_cierre\": null, se obtuvo %v (presente: %v)", valor, ok)
	}
	if resp["estado"] != string(domain.ProyectoPlanificado) {
		t.Errorf("sin fechas ni Sprints se esperaba el estado %q, se obtuvo %v", domain.ProyectoPlanificado, resp["estado"])
	}
}

func TestProjectHandler_Reabrir_Errores(t *testing.T) {
	casos := []struct {
		nombre  string
		cerrar  bool
		id      string
		actor   domain.User
		codigo  int
		mensaje string
	}{
		{"abierto", false, "", anaDePrueba, http.StatusConflict, "El proyecto no está cerrado"},
		{"Product Builder", true, "", domain.User{ID: productBuilderHTTP}, http.StatusForbidden, "Solo el Scrum Master puede realizar esta acción"},
		{"no integrante", true, "", domain.User{ID: 99}, http.StatusForbidden, "No sos integrante de este proyecto"},
		{"proyecto inexistente", true, "999", anaDePrueba, http.StatusNotFound, "proyecto no encontrado"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			handler, proyectos := nuevoProjectHandlerDePrueba()
			idStr := proyectoDeAnaConBuilder(t, handler, proyectos)
			if c.cerrar {
				cerrarProyectoHTTP(handler, idStr, anaDePrueba)
			}
			objetivo := idStr
			if c.id != "" {
				objetivo = c.id
			}

			rec := reabrirProyectoHTTP(handler, objetivo, c.actor)

			if rec.Code != c.codigo {
				t.Fatalf("se esperaba %d, se obtuvo %d (%s)", c.codigo, rec.Code, rec.Body.String())
			}
			var resp errorResponse
			_ = json.Unmarshal(rec.Body.Bytes(), &resp)
			if resp.Mensaje != c.mensaje {
				t.Errorf("se esperaba el mensaje %q, se obtuvo %q", c.mensaje, resp.Mensaje)
			}
		})
	}
}
