package http

import (
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// AuthHandler es el adaptador de entrada de la autenticación (HU-20).
type AuthHandler struct {
	registrar *service.RegistrarUsuario
}

func NewAuthHandler(registrar *service.RegistrarUsuario) *AuthHandler {
	return &AuthHandler{registrar: registrar}
}

// Registro atiende POST /auth/registro.
func (h *AuthHandler) Registro(w nethttp.ResponseWriter, r *nethttp.Request) {
}
