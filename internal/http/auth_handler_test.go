package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

const iteracionesDePruebaAuth = 1

type credencialesEnMemoriaHTTP struct {
	mu       sync.Mutex
	usuarios map[string]domain.User
	hashes   map[string]string
}

func (c *credencialesEnMemoriaHTTP) Registrar(_ context.Context, u domain.User, hash string) (domain.User, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.usuarios == nil {
		c.usuarios = map[string]domain.User{}
		c.hashes = map[string]string{}
	}
	if _, ok := c.usuarios[u.NombreNormalizado]; ok {
		return domain.User{}, repository.ErrUsuarioDuplicado
	}
	u.ID = int64(len(c.usuarios) + 1)
	c.usuarios[u.NombreNormalizado] = u
	c.hashes[u.NombreNormalizado] = hash
	return u, nil
}

type sesionesEnMemoriaHTTP struct {
	mu       sync.Mutex
	sesiones map[string]domain.Sesion
}

func (s *sesionesEnMemoriaHTTP) Crear(_ context.Context, huella string, sesion domain.Sesion) (domain.Sesion, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.sesiones == nil {
		s.sesiones = map[string]domain.Sesion{}
	}
	sesion.ID = int64(len(s.sesiones) + 1)
	s.sesiones[huella] = sesion
	return sesion, nil
}

var relojAuthHTTP = func() time.Time { return time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC) }

func nuevoAuthHandlerDePrueba() (*AuthHandler, *credencialesEnMemoriaHTTP) {
	credenciales := &credencialesEnMemoriaHTTP{}
	emisor := service.NewEmisorSesiones(&sesionesEnMemoriaHTTP{}, relojAuthHTTP)
	registrar := service.NewRegistrarUsuario(credenciales, service.NewServicioContrasenas(iteracionesDePruebaAuth), emisor)
	return NewAuthHandler(registrar), credenciales
}

func postRegistro(h *AuthHandler, cuerpo string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(nethttp.MethodPost, "/auth/registro", strings.NewReader(cuerpo))
	rec := httptest.NewRecorder()
	h.Registro(rec, req)
	return rec
}

func TestAuthHandler_Registro_Devuelve201ConUsuarioTokenYVencimiento(t *testing.T) {
	h, _ := nuevoAuthHandlerDePrueba()

	rec := postRegistro(h, `{"nombre":"Ana.Perez","contrasena":"una-clave-larga"}`)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201 (cuerpo: %s)", rec.Code, rec.Body.String())
	}
	var resp struct {
		Usuario struct {
			ID     int64  `json:"id"`
			Nombre string `json:"nombre"`
		} `json:"usuario"`
		Token   string `json:"token"`
		VenceEn string `json:"vence_en"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta no es JSON: %v", err)
	}
	if resp.Usuario.ID == 0 || resp.Usuario.Nombre != "Ana.Perez" {
		t.Errorf("usuario = %+v, se esperaba id y nombre tal como se escribió", resp.Usuario)
	}
	if len(resp.Token) != 43 {
		t.Errorf("token de %d caracteres, se esperaban 43", len(resp.Token))
	}
	if resp.VenceEn != "2026-11-11T12:00:00Z" {
		t.Errorf("vence_en = %q, se esperaba %q", resp.VenceEn, "2026-11-11T12:00:00Z")
	}
}

func TestAuthHandler_Registro_LaRespuestaNoContieneContrasenaNiHash(t *testing.T) {
	h, credenciales := nuevoAuthHandlerDePrueba()

	rec := postRegistro(h, `{"nombre":"ana.perez","contrasena":"una-clave-larga"}`)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("status = %d, se esperaba 201", rec.Code)
	}
	cuerpo := rec.Body.String()
	if strings.Contains(cuerpo, "una-clave-larga") {
		t.Errorf("la respuesta contiene la contraseña")
	}
	if hash := credenciales.hashes["ana.perez"]; hash == "" || strings.Contains(cuerpo, hash) {
		t.Errorf("la respuesta contiene el hash o no se guardó uno")
	}
	for _, prohibido := range []string{"contrasena", "hash"} {
		if strings.Contains(cuerpo, prohibido) {
			t.Errorf("la respuesta menciona %q: %s", prohibido, cuerpo)
		}
	}
}

func TestAuthHandler_Registro_Devuelve400PorCampoInvalido(t *testing.T) {
	casos := []struct {
		nombre string
		cuerpo string
		campo  string
	}{
		{"nombre fuera de reglas", `{"nombre":"ab","contrasena":"una-clave-larga"}`, "nombre"},
		{"contraseña corta", `{"nombre":"ana.perez","contrasena":"corta"}`, "contrasena"},
		{"contraseña ausente", `{"nombre":"ana.perez"}`, "contrasena"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			h, _ := nuevoAuthHandlerDePrueba()

			rec := postRegistro(h, c.cuerpo)

			if rec.Code != nethttp.StatusBadRequest {
				t.Fatalf("status = %d, se esperaba 400", rec.Code)
			}
			var resp errorResponse
			_ = json.Unmarshal(rec.Body.Bytes(), &resp)
			if resp.Campo != c.campo || resp.Mensaje == "" {
				t.Errorf("respuesta = %+v, se esperaba campo %q con mensaje", resp, c.campo)
			}
		})
	}
}

func TestAuthHandler_Registro_Devuelve400ConJSONInvalido(t *testing.T) {
	h, _ := nuevoAuthHandlerDePrueba()

	rec := postRegistro(h, `{no es json`)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Mensaje != "JSON inválido" || resp.Campo != "" {
		t.Errorf("respuesta = %+v, se esperaba solo el mensaje JSON inválido", resp)
	}
}

func TestAuthHandler_Registro_Devuelve409ConNombreEnUso(t *testing.T) {
	h, _ := nuevoAuthHandlerDePrueba()
	postRegistro(h, `{"nombre":"Ana.Perez","contrasena":"una-clave-larga"}`)

	rec := postRegistro(h, `{"nombre":"ana.perez","contrasena":"otra-clave-larga"}`)

	if rec.Code != nethttp.StatusConflict {
		t.Fatalf("status = %d, se esperaba 409", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Mensaje != "Ese nombre de usuario ya existe" {
		t.Errorf("mensaje = %q", resp.Mensaje)
	}
}
