package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// HU-16: DELETE /projects/{id}/members/{usuario_id} quita a un integrante.

func quitarIntegranteHTTP(handler *ProjectHandler, idStr, usuarioStr string, actor domain.User) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/projects/"+idStr+"/members/"+usuarioStr, nil)
	req.SetPathValue("id", idStr)
	req.SetPathValue("usuario_id", usuarioStr)
	handler.QuitarIntegrante(rec, conSesionDe(req, actor))
	return rec
}

func TestProjectHandler_QuitarIntegrante_ElScrumMasterQuitaAUnProductBuilder(t *testing.T) {
	handler, proyectos := nuevoProjectHandlerDePrueba()
	idStr := proyectoDeAnaConBuilder(t, handler, proyectos)

	rec := quitarIntegranteHTTP(handler, idStr, "20", anaDePrueba)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("se esperaba 204, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("se esperaba una respuesta sin cuerpo, se obtuvo %q", rec.Body.String())
	}
	for _, m := range proyectos.miembros {
		if m.UserID == productBuilderHTTP {
			t.Errorf("el Product Builder no debía seguir en el equipo: %+v", proyectos.miembros)
		}
	}
}

func TestProjectHandler_QuitarIntegrante_Errores(t *testing.T) {
	casos := []struct {
		nombre  string
		id      string
		usuario string
		actor   domain.User
		cerrado bool
		codigo  int
		campo   string
		mensaje string
	}{
		{nombre: "id no numérico", id: "abc", usuario: "20", actor: anaDePrueba, codigo: http.StatusBadRequest, campo: "id", mensaje: "el identificador debe ser un entero positivo"},
		{nombre: "usuario_id no numérico", usuario: "abc", actor: anaDePrueba, codigo: http.StatusBadRequest, campo: "usuario_id", mensaje: "el identificador debe ser un entero positivo"},
		{nombre: "usuario_id cero", usuario: "0", actor: anaDePrueba, codigo: http.StatusBadRequest, campo: "usuario_id", mensaje: "el identificador debe ser un entero positivo"},
		{nombre: "proyecto inexistente", id: "999", usuario: "20", actor: anaDePrueba, codigo: http.StatusNotFound, mensaje: "proyecto no encontrado"},
		{nombre: "actor no integrante", usuario: "20", actor: domain.User{ID: 21}, codigo: http.StatusForbidden, mensaje: "No sos integrante de este proyecto"},
		{nombre: "actor Product Builder", usuario: "20", actor: domain.User{ID: productBuilderHTTP}, codigo: http.StatusForbidden, mensaje: "Solo el Scrum Master puede realizar esta acción"},
		{nombre: "proyecto cerrado", usuario: "20", actor: anaDePrueba, cerrado: true, codigo: http.StatusConflict, mensaje: "El proyecto está cerrado"},
		{nombre: "integrante inexistente", usuario: "21", actor: anaDePrueba, codigo: http.StatusNotFound, mensaje: "integrante no encontrado"},
		{nombre: "el Scrum Master", usuario: "7", actor: anaDePrueba, codigo: http.StatusConflict, mensaje: "No se puede quitar al Scrum Master del proyecto"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			handler, proyectos := nuevoProjectHandlerDePrueba()
			idStr := proyectoDeAnaConBuilder(t, handler, proyectos)
			if c.cerrado {
				if rec := cerrarProyectoHTTP(handler, idStr, anaDePrueba); rec.Code != http.StatusOK {
					t.Fatalf("no se pudo cerrar de antemano: %d", rec.Code)
				}
			}
			objetivo := idStr
			if c.id != "" {
				objetivo = c.id
			}
			antes := len(proyectos.miembros)

			rec := quitarIntegranteHTTP(handler, objetivo, c.usuario, c.actor)

			if rec.Code != c.codigo {
				t.Fatalf("se esperaba %d, se obtuvo %d (%s)", c.codigo, rec.Code, rec.Body.String())
			}
			var resp errorResponse
			_ = json.Unmarshal(rec.Body.Bytes(), &resp)
			if resp.Mensaje != c.mensaje || resp.Campo != c.campo {
				t.Errorf("se esperaba {campo: %q, mensaje: %q}, se obtuvo %+v", c.campo, c.mensaje, resp)
			}
			if len(proyectos.miembros) != antes {
				t.Errorf("el equipo no debía cambiar: %+v", proyectos.miembros)
			}
		})
	}
}

func TestProjectHandler_QuitarIntegrante_SinSesionResponde401(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/projects/1/members/20", nil)
	req.SetPathValue("id", "1")
	req.SetPathValue("usuario_id", "20")

	handler.QuitarIntegrante(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
}
