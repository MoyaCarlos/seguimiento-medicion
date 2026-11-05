package http

import (
	"bytes"
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

func (c *credencialesEnMemoriaHTTP) ObtenerPorNombreNormalizado(_ context.Context, normalizado string) (domain.User, string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	u, ok := c.usuarios[normalizado]
	if !ok {
		return domain.User{}, "", domain.ErrNoEncontrado
	}
	return u, c.hashes[normalizado], nil
}

type sesionesEnMemoriaHTTP struct {
	mu           sync.Mutex
	sesiones     map[string]domain.Sesion
	credenciales *credencialesEnMemoriaHTTP
}

func (s *sesionesEnMemoriaHTTP) ObtenerPorHuella(_ context.Context, huella string) (domain.SesionConUsuario, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sesion, ok := s.sesiones[huella]
	if !ok {
		return domain.SesionConUsuario{}, domain.ErrNoEncontrado
	}
	resultado := domain.SesionConUsuario{Sesion: sesion}
	for _, u := range s.credenciales.usuarios {
		if u.ID == sesion.UsuarioID {
			resultado.Usuario = u
		}
	}
	return resultado, nil
}

func (s *sesionesEnMemoriaHTTP) cerrar(huella string, en time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sesion := s.sesiones[huella]
	sesion.CerradaEn = &en
	s.sesiones[huella] = sesion
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
	h, credenciales, _, _ := nuevoAuthHandlerConSesiones()
	return h, credenciales
}

// nuevoAuthHandlerConSesiones también devuelve las sesiones y el reloj, que el
// test puede adelantar para vencer una sesión.
func nuevoAuthHandlerConSesiones() (*AuthHandler, *credencialesEnMemoriaHTTP, *sesionesEnMemoriaHTTP, *time.Time) {
	credenciales := &credencialesEnMemoriaHTTP{}
	sesiones := &sesionesEnMemoriaHTTP{credenciales: credenciales}
	ahora := relojAuthHTTP()
	reloj := func() time.Time { return ahora }
	emisor := service.NewEmisorSesiones(sesiones, reloj)
	hasheador := service.NewServicioContrasenas(iteracionesDePruebaAuth)
	registrar := service.NewRegistrarUsuario(credenciales, hasheador, emisor)
	ingresar := service.NewIngresar(credenciales, hasheador, emisor)
	obtener := service.NewObtenerSesion(sesiones, reloj)
	return NewAuthHandler(registrar, ingresar, obtener), credenciales, sesiones, &ahora
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

func postIngreso(h *AuthHandler, cuerpo string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(nethttp.MethodPost, "/auth/ingreso", strings.NewReader(cuerpo))
	rec := httptest.NewRecorder()
	h.Ingreso(rec, req)
	return rec
}

func TestAuthHandler_Ingreso_Devuelve200ConElMismoCuerpoQueElRegistro(t *testing.T) {
	h, _ := nuevoAuthHandlerDePrueba()
	postRegistro(h, `{"nombre":"Ana.Perez","contrasena":"una-clave-larga"}`)

	rec := postIngreso(h, `{"nombre":"ANA.PEREZ","contrasena":"una-clave-larga"}`)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, se esperaba 200 (cuerpo: %s)", rec.Code, rec.Body.String())
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
		t.Errorf("usuario = %+v, se esperaba el nombre con el que se registró", resp.Usuario)
	}
	if len(resp.Token) != 43 {
		t.Errorf("token de %d caracteres, se esperaban 43", len(resp.Token))
	}
	if resp.VenceEn != "2026-11-11T12:00:00Z" {
		t.Errorf("vence_en = %q, se esperaba %q", resp.VenceEn, "2026-11-11T12:00:00Z")
	}
}

func TestAuthHandler_Ingreso_LaRespuestaNoContieneContrasenaNiHash(t *testing.T) {
	h, credenciales := nuevoAuthHandlerDePrueba()
	postRegistro(h, `{"nombre":"ana.perez","contrasena":"una-clave-larga"}`)

	rec := postIngreso(h, `{"nombre":"ana.perez","contrasena":"una-clave-larga"}`)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, se esperaba 200", rec.Code)
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

func TestAuthHandler_Ingreso_Los401TienenElMismoCuerpoByteAByte(t *testing.T) {
	h, credenciales := nuevoAuthHandlerDePrueba()
	postRegistro(h, `{"nombre":"ana.perez","contrasena":"una-clave-larga"}`)
	if _, err := credenciales.Registrar(context.Background(), domain.User{Nombre: "previo", NombreNormalizado: "previo"}, ""); err != nil {
		t.Fatalf("preparar usuario previo: %v", err)
	}
	casos := []struct {
		nombre string
		cuerpo string
	}{
		{"contraseña incorrecta", `{"nombre":"ana.perez","contrasena":"otra-clave-larga"}`},
		{"usuario inexistente", `{"nombre":"fantasma","contrasena":"una-clave-larga"}`},
		{"usuario sin contraseña", `{"nombre":"previo","contrasena":"una-clave-larga"}`},
	}
	const esperado = `{"mensaje":"Usuario o contraseña incorrectos"}`

	var primero []byte
	for _, c := range casos {
		rec := postIngreso(h, c.cuerpo)

		if rec.Code != nethttp.StatusUnauthorized {
			t.Fatalf("%s: status = %d, se esperaba 401", c.nombre, rec.Code)
		}
		cuerpo := rec.Body.Bytes()
		var resp errorResponse
		if err := json.Unmarshal(cuerpo, &resp); err != nil || resp.Mensaje != "Usuario o contraseña incorrectos" || resp.Campo != "" {
			t.Errorf("%s: cuerpo = %s, se esperaba %s", c.nombre, cuerpo, esperado)
		}
		if primero == nil {
			primero = cuerpo
		} else if !bytes.Equal(primero, cuerpo) {
			t.Errorf("%s: cuerpo %q distinto de %q", c.nombre, cuerpo, primero)
		}
	}
}

func TestAuthHandler_Ingreso_NoValidaFormatoDeNombreNiLargoDeContrasena(t *testing.T) {
	h, _ := nuevoAuthHandlerDePrueba()
	postRegistro(h, `{"nombre":"ana.perez","contrasena":"una-clave-larga"}`)
	esperado := postIngreso(h, `{"nombre":"fantasma","contrasena":"una-clave-larga"}`)

	for nombre, cuerpo := range map[string]string{
		"cuerpo vacío":       `{}`,
		"nombre muy corto":   `{"nombre":"ab","contrasena":"una-clave-larga"}`,
		"contraseña corta":   `{"nombre":"ana.perez","contrasena":"x"}`,
		"contraseña ausente": `{"nombre":"ana.perez"}`,
	} {
		rec := postIngreso(h, cuerpo)

		if rec.Code != nethttp.StatusUnauthorized {
			t.Errorf("%s: status = %d, se esperaba 401", nombre, rec.Code)
		}
		if !bytes.Equal(rec.Body.Bytes(), esperado.Body.Bytes()) {
			t.Errorf("%s: cuerpo %q distinto del 401 común %q", nombre, rec.Body.Bytes(), esperado.Body.Bytes())
		}
	}
}

func TestAuthHandler_Ingreso_Devuelve400ConJSONInvalido(t *testing.T) {
	h, _ := nuevoAuthHandlerDePrueba()

	rec := postIngreso(h, `{no es json`)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("status = %d, se esperaba 400", rec.Code)
	}
	var resp errorResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Mensaje != "JSON inválido" || resp.Campo != "" {
		t.Errorf("respuesta = %+v, se esperaba solo el mensaje JSON inválido", resp)
	}
}

func getYo(h *AuthHandler, autorizacion string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(nethttp.MethodGet, "/auth/yo", nil)
	if autorizacion != "" {
		req.Header.Set("Authorization", autorizacion)
	}
	rec := httptest.NewRecorder()
	h.Yo(rec, req)
	return rec
}

func registrarYObtenerToken(t *testing.T, h *AuthHandler) string {
	t.Helper()
	rec := postRegistro(h, `{"nombre":"Ana.Perez","contrasena":"una-clave-larga"}`)
	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil || resp.Token == "" {
		t.Fatalf("preparar sesión: %v (cuerpo: %s)", err, rec.Body.String())
	}
	return resp.Token
}

func TestAuthHandler_Yo_Devuelve200SoloConIdYNombre(t *testing.T) {
	h, credenciales := nuevoAuthHandlerDePrueba()
	token := registrarYObtenerToken(t, h)

	rec := getYo(h, "Bearer "+token)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, se esperaba 200 (cuerpo: %s)", rec.Code, rec.Body.String())
	}
	var campos map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &campos); err != nil {
		t.Fatalf("respuesta no es JSON: %v", err)
	}
	if len(campos) != 2 || campos["id"] == nil || campos["nombre"] != "Ana.Perez" {
		t.Errorf("campos = %v, se esperaba solo id y nombre %q", campos, "Ana.Perez")
	}
	cuerpo := rec.Body.String()
	if strings.Contains(cuerpo, "una-clave-larga") {
		t.Errorf("la respuesta contiene la contraseña")
	}
	if hash := credenciales.hashes["ana.perez"]; hash == "" || strings.Contains(cuerpo, hash) {
		t.Errorf("la respuesta contiene el hash o no se guardó uno")
	}
	if strings.Contains(cuerpo, token) {
		t.Errorf("la respuesta repite el token")
	}
}

func TestAuthHandler_Yo_Los401TienenElMismoCuerpoByteAByte(t *testing.T) {
	casos := []struct {
		nombre       string
		autorizacion func(token string, sesiones *sesionesEnMemoriaHTTP, ahora *time.Time) string
	}{
		{"sin cabecera", func(string, *sesionesEnMemoriaHTTP, *time.Time) string { return "" }},
		{"esquema Basic", func(token string, _ *sesionesEnMemoriaHTTP, _ *time.Time) string { return "Basic " + token }},
		{"Bearer sin valor", func(string, *sesionesEnMemoriaHTTP, *time.Time) string { return "Bearer " }},
		{"Bearer sin espacio ni valor", func(string, *sesionesEnMemoriaHTTP, *time.Time) string { return "Bearer" }},
		{"token mal formado", func(string, *sesionesEnMemoriaHTTP, *time.Time) string { return "Bearer no-es-un-token" }},
		{"token inexistente", func(string, *sesionesEnMemoriaHTTP, *time.Time) string {
			return "Bearer " + strings.Repeat("A", 43)
		}},
		{"sesión vencida", func(token string, _ *sesionesEnMemoriaHTTP, ahora *time.Time) string {
			*ahora = ahora.Add(domain.DuracionSesion)
			return "Bearer " + token
		}},
		{"sesión cerrada", func(token string, sesiones *sesionesEnMemoriaHTTP, ahora *time.Time) string {
			sesiones.cerrar(service.HuellaDeToken(token), *ahora)
			return "Bearer " + token
		}},
	}
	const esperado = "{\"mensaje\":\"Tu sesión no es válida\"}\n"

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			h, _, sesiones, ahora := nuevoAuthHandlerConSesiones()
			token := registrarYObtenerToken(t, h)

			rec := getYo(h, c.autorizacion(token, sesiones, ahora))

			if rec.Code != nethttp.StatusUnauthorized {
				t.Fatalf("status = %d, se esperaba 401", rec.Code)
			}
			if rec.Body.String() != esperado {
				t.Errorf("cuerpo = %q, se esperaba %q", rec.Body.String(), esperado)
			}
		})
	}
}
