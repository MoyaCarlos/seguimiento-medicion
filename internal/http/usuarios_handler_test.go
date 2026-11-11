package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type usuariosRegistradosDePrueba []domain.User

func (u usuariosRegistradosDePrueba) ListarRegistrados(context.Context) ([]domain.User, error) {
	return u, nil
}

// HU-21, escenario 8: GET /users responde id y nombre de cada usuario
// registrado, sin contraseña, huella ni datos de sesión.
func TestUsuariosHandler_Listar(t *testing.T) {
	registrados := usuariosRegistradosDePrueba{
		{ID: 1, Nombre: "ana", NombreNormalizado: "ana"},
		{ID: 2, Nombre: "Bruno", NombreNormalizado: "bruno"},
	}
	handler := NewUsuariosHandler(service.NewListarUsuarios(registrados))

	rec := httptest.NewRecorder()
	handler.Listar(rec, conSesionDe(httptest.NewRequest(http.MethodGet, "/users", nil), anaDePrueba))

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if body := strings.TrimSpace(rec.Body.String()); body != `[{"id":1,"nombre":"ana"},{"id":2,"nombre":"Bruno"}]` {
		t.Errorf("cuerpo inesperado: %s", body)
	}
}

func TestUsuariosHandler_ListarVacio(t *testing.T) {
	handler := NewUsuariosHandler(service.NewListarUsuarios(usuariosRegistradosDePrueba{}))

	rec := httptest.NewRecorder()
	handler.Listar(rec, conSesionDe(httptest.NewRequest(http.MethodGet, "/users", nil), anaDePrueba))

	if body := strings.TrimSpace(rec.Body.String()); rec.Code != http.StatusOK || body != "[]" {
		t.Errorf("se esperaba 200 [], se obtuvo %d %s", rec.Code, body)
	}
}
