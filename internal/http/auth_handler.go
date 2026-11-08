package http

import (
	"encoding/json"
	nethttp "net/http"
	"strings"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// AuthHandler es el adaptador de entrada de la autenticación (HU-20). No
// registra cuerpos ni cabeceras: llevan contraseñas y tokens.
type AuthHandler struct {
	registrar *service.RegistrarUsuario
	ingresar  *service.Ingresar
	obtener   *service.ObtenerSesion
	cerrar    *service.CerrarSesion
}

func NewAuthHandler(registrar *service.RegistrarUsuario, ingresar *service.Ingresar, obtener *service.ObtenerSesion, cerrar *service.CerrarSesion) *AuthHandler {
	return &AuthHandler{registrar: registrar, ingresar: ingresar, obtener: obtener, cerrar: cerrar}
}

// Yo atiende GET /auth/yo.
func (h *AuthHandler) Yo(w nethttp.ResponseWriter, r *nethttp.Request) {
	token, ok := tokenBearer(r)
	if !ok {
		escribirError(w, domain.ErrNoAutenticado)
		return
	}
	sesion, err := h.obtener.Ejecutar(r.Context(), token)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, usuarioSesionResponse{ID: sesion.Usuario.ID, Nombre: sesion.Usuario.Nombre})
}

// Salir atiende POST /auth/salir.
func (h *AuthHandler) Salir(w nethttp.ResponseWriter, r *nethttp.Request) {
	token, ok := tokenBearer(r)
	if !ok {
		escribirError(w, domain.ErrNoAutenticado)
		return
	}
	if err := h.cerrar.Ejecutar(r.Context(), token); err != nil {
		escribirError(w, err)
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}

// tokenBearer lee el token de "Authorization: Bearer <token>". El esquema se
// compara sin distinguir mayúsculas; sin cabecera, con otro esquema o sin valor
// devuelve false.
func tokenBearer(r *nethttp.Request) (string, bool) {
	esquema, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	if !ok || !strings.EqualFold(esquema, "Bearer") || token == "" {
		return "", false
	}
	return token, true
}

// maxCuerpoAuth es el tamaño máximo del cuerpo de /auth/registro e /ingreso.
const maxCuerpoAuth = 64 * 1024

// Registro atiende POST /auth/registro.
func (h *AuthHandler) Registro(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req registroRequest
	if err := json.NewDecoder(nethttp.MaxBytesReader(w, r.Body, maxCuerpoAuth)).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}
	resultado, err := h.registrar.Ejecutar(r.Context(), req.Nombre, req.Contrasena)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusCreated, aRespuestaSesion(resultado))
}

// Ingreso atiende POST /auth/ingreso.
func (h *AuthHandler) Ingreso(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req registroRequest
	if err := json.NewDecoder(nethttp.MaxBytesReader(w, r.Body, maxCuerpoAuth)).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}
	resultado, err := h.ingresar.Ejecutar(r.Context(), req.Nombre, req.Contrasena)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aRespuestaSesion(resultado))
}

func aRespuestaSesion(r service.ResultadoSesion) sesionResponse {
	return sesionResponse{
		Usuario: usuarioSesionResponse{ID: r.Usuario.ID, Nombre: r.Usuario.Nombre},
		Token:   r.Token,
		VenceEn: r.VenceEn.UTC().Format(time.RFC3339),
	}
}
