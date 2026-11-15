package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func registrarUsuarioDePrueba(t *testing.T, credenciales *credencialesEnMemoriaHTTP, nombre string) domain.User {
	t.Helper()
	u, err := credenciales.Registrar(context.Background(), domain.User{
		Nombre:            nombre,
		NombreNormalizado: domain.NormalizarNombre(nombre),
	}, "hash")
	if err != nil {
		t.Fatalf("preparar usuario %q: %v", nombre, err)
	}
	return u
}

func pedirCambiarNombre(h *AuthHandler, usuario domain.User, cuerpo string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(nethttp.MethodPut, "/auth/yo", strings.NewReader(cuerpo))
	if usuario.ID > 0 {
		req = req.WithContext(context.WithValue(req.Context(), claveUsuarioDeSesion{}, usuario))
	}
	rec := httptest.NewRecorder()
	h.CambiarNombre(rec, req)
	return rec
}

func TestAuthHandler_CambiarNombre_Devuelve200ConElNombreNuevo(t *testing.T) {
	h, credenciales := nuevoAuthHandlerDePrueba()
	ana := registrarUsuarioDePrueba(t, credenciales, "ana")

	rec := pedirCambiarNombre(h, ana, `{"nombre":"Ana.Perez"}`)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, se esperaba 200 (cuerpo: %s)", rec.Code, rec.Body.String())
	}
	var resp usuarioSesionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta no es JSON: %v", err)
	}
	if resp.ID != ana.ID || resp.Nombre != "Ana.Perez" {
		t.Errorf("respuesta = %+v, se esperaba el id y el nombre nuevo", resp)
	}
	guardada, _, err := credenciales.ObtenerPorNombreNormalizado(context.Background(), "ana.perez")
	if err != nil || guardada.ID != ana.ID || guardada.Nombre != "Ana.Perez" {
		t.Errorf("usuario no renombrado: %+v, %v", guardada, err)
	}
}

func TestAuthHandler_CambiarNombre_Devuelve401SinSesion(t *testing.T) {
	h, _ := nuevoAuthHandlerDePrueba()

	rec := pedirCambiarNombre(h, domain.User{}, `{"nombre":"ana.perez"}`)

	if rec.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status = %d, se esperaba 401", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Mensaje != domain.ErrNoAutenticado.Error() {
		t.Errorf("mensaje = %q, se esperaba %q", resp.Mensaje, domain.ErrNoAutenticado.Error())
	}
}

func TestAuthHandler_CambiarNombre_Devuelve400ConCuerpoONombreInvalido(t *testing.T) {
	casos := []struct {
		nombre string
		cuerpo string
		campo  string
	}{
		{"JSON inválido", `{no es json`, ""},
		{"nombre ausente", `{}`, "nombre"},
		{"nombre corto", `{"nombre":"ab"}`, "nombre"},
		{"nombre con espacio", `{"nombre":"ana perez"}`, "nombre"},
		{"nombre con arroba", `{"nombre":"ana@perez"}`, "nombre"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			h, credenciales := nuevoAuthHandlerDePrueba()
			ana := registrarUsuarioDePrueba(t, credenciales, "ana")

			rec := pedirCambiarNombre(h, ana, c.cuerpo)

			if rec.Code != nethttp.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400 (cuerpo: %s)", rec.Code, rec.Body.String())
			}
			var resp errorResponse
			_ = json.Unmarshal(rec.Body.Bytes(), &resp)
			if resp.Campo != c.campo {
				t.Errorf("campo = %q, se esperaba %q", resp.Campo, c.campo)
			}
			if resp.Mensaje == "" {
				t.Errorf("se esperaba un mensaje")
			}
		})
	}
}

func TestAuthHandler_CambiarNombre_Devuelve409ConNombreEnUso(t *testing.T) {
	h, credenciales := nuevoAuthHandlerDePrueba()
	ana := registrarUsuarioDePrueba(t, credenciales, "ana")
	registrarUsuarioDePrueba(t, credenciales, "bruno")

	rec := pedirCambiarNombre(h, ana, `{"nombre":"BRUNO"}`)

	if rec.Code != nethttp.StatusConflict {
		t.Fatalf("status = %d, se esperaba 409 (cuerpo: %s)", rec.Code, rec.Body.String())
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Mensaje != domain.ErrNombreUsuarioEnUso.Error() {
		t.Errorf("mensaje = %q, se esperaba %q", resp.Mensaje, domain.ErrNombreUsuarioEnUso.Error())
	}
	guardada, _, err := credenciales.ObtenerPorNombreNormalizado(context.Background(), "ana")
	if err != nil || guardada.ID != ana.ID || guardada.Nombre != "ana" {
		t.Errorf("Ana debía conservar su nombre: %+v, %v", guardada, err)
	}
}

func TestAuthHandler_CambiarNombre_IgnoraElUsuarioIDDelCuerpo(t *testing.T) {
	h, credenciales := nuevoAuthHandlerDePrueba()
	ana := registrarUsuarioDePrueba(t, credenciales, "ana")
	bruno := registrarUsuarioDePrueba(t, credenciales, "bruno")

	rec := pedirCambiarNombre(h, ana, `{"usuario_id":`+itoa(bruno.ID)+`,"nombre":"Ana.Perez"}`)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, se esperaba 200 (cuerpo: %s)", rec.Code, rec.Body.String())
	}
	guardada, _, err := credenciales.ObtenerPorNombreNormalizado(context.Background(), "ana.perez")
	if err != nil || guardada.ID != ana.ID {
		t.Errorf("se debía renombrar al usuario de la sesión: %+v, %v", guardada, err)
	}
	if otro, _, err := credenciales.ObtenerPorNombreNormalizado(context.Background(), "bruno"); err != nil || otro.ID != bruno.ID {
		t.Errorf("no se debía tocar al usuario del cuerpo: %+v, %v", otro, err)
	}
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	var digitos []byte
	for n > 0 {
		digitos = append([]byte{byte('0' + n%10)}, digitos...)
		n /= 10
	}
	return string(digitos)
}
