package http

import (
	"encoding/json"
	nethttp "net/http"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// AuthHandler es el adaptador de entrada de la autenticación (HU-20). No
// registra cuerpos ni cabeceras: llevan contraseñas y tokens.
type AuthHandler struct {
	registrar *service.RegistrarUsuario
	ingresar  *service.Ingresar
	obtener   *service.ObtenerSesion
}

func NewAuthHandler(registrar *service.RegistrarUsuario, ingresar *service.Ingresar, obtener *service.ObtenerSesion) *AuthHandler {
	return &AuthHandler{registrar: registrar, ingresar: ingresar, obtener: obtener}
}

// Yo atiende GET /auth/yo.
func (h *AuthHandler) Yo(_ nethttp.ResponseWriter, _ *nethttp.Request) {}

// Registro atiende POST /auth/registro.
func (h *AuthHandler) Registro(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req registroRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
