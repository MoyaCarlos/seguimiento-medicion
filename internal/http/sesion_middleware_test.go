package http

import (
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// handlerQueAnota registra si se ejecutó y qué usuario vio en la sesión.
type handlerQueAnota struct {
	ejecutado bool
	usuario   domain.User
}

func (h *handlerQueAnota) atender(w nethttp.ResponseWriter, r *nethttp.Request) {
	h.ejecutado = true
	h.usuario = usuarioDeSesion(r)
	w.WriteHeader(nethttp.StatusNoContent)
}

func pedirConAutorizacion(conSesion func(nethttp.HandlerFunc) nethttp.HandlerFunc, h *handlerQueAnota, autorizacion string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(nethttp.MethodGet, "/protegida", nil)
	if autorizacion != "" {
		req.Header.Set("Authorization", autorizacion)
	}
	rec := httptest.NewRecorder()
	conSesion(h.atender)(rec, req)
	return rec
}

func TestConSesion_ConSesionVigenteEjecutaElHandlerConElUsuario(t *testing.T) {
	auth, _, _, _ := nuevoAuthHandlerConSesiones()
	token := registrarYObtenerToken(t, auth)
	h := &handlerQueAnota{}

	rec := pedirConAutorizacion(ConSesion(auth.obtener), h, "Bearer "+token)

	if rec.Code != nethttp.StatusNoContent {
		t.Fatalf("status = %d, se esperaba 204 (cuerpo: %s)", rec.Code, rec.Body.String())
	}
	if !h.ejecutado {
		t.Fatal("el handler no se ejecutó")
	}
	if h.usuario.ID == 0 || h.usuario.Nombre != "Ana.Perez" {
		t.Errorf("usuario de la sesión = %+v, se esperaba Ana.Perez con ID", h.usuario)
	}
}

func TestConSesion_SinSesionValidaResponde401SinEjecutarElHandler(t *testing.T) {
	casos := []struct {
		nombre       string
		autorizacion func(token string, sesiones *sesionesEnMemoriaHTTP, ahora *time.Time) string
	}{
		{"sin cabecera", func(string, *sesionesEnMemoriaHTTP, *time.Time) string { return "" }},
		{"esquema Basic", func(token string, _ *sesionesEnMemoriaHTTP, _ *time.Time) string { return "Basic " + token }},
		{"Bearer sin valor", func(string, *sesionesEnMemoriaHTTP, *time.Time) string { return "Bearer " }},
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
			auth, _, sesiones, ahora := nuevoAuthHandlerConSesiones()
			token := registrarYObtenerToken(t, auth)
			h := &handlerQueAnota{}

			rec := pedirConAutorizacion(ConSesion(auth.obtener), h, c.autorizacion(token, sesiones, ahora))

			if rec.Code != nethttp.StatusUnauthorized {
				t.Fatalf("status = %d, se esperaba 401", rec.Code)
			}
			if rec.Body.String() != esperado {
				t.Errorf("cuerpo = %q, se esperaba %q", rec.Body.String(), esperado)
			}
			if h.ejecutado {
				t.Error("el handler se ejecutó sin sesión válida")
			}
		})
	}
}
