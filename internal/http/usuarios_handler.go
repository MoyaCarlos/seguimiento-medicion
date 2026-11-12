package http

import (
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// UsuariosHandler expone los usuarios registrados (HU-21).
type UsuariosHandler struct {
	listar *service.ListarUsuarios
}

func NewUsuariosHandler(listar *service.ListarUsuarios) *UsuariosHandler {
	return &UsuariosHandler{listar: listar}
}

// usuarioResponse es lo único que se expone de un usuario: nunca su contraseña,
// huella ni sesiones.
type usuarioResponse struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
}

// Listar atiende GET /users.
func (h *UsuariosHandler) Listar(w nethttp.ResponseWriter, r *nethttp.Request) {
	usuarios, err := h.listar.Ejecutar(r.Context())
	if err != nil {
		escribirError(w, err)
		return
	}
	respuesta := make([]usuarioResponse, 0, len(usuarios))
	for _, u := range usuarios {
		respuesta = append(respuesta, usuarioResponse{ID: u.ID, Nombre: u.Nombre})
	}
	escribirJSON(w, nethttp.StatusOK, respuesta)
}
