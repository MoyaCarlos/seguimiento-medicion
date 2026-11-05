package http

import (
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// HistorialHandler es el adaptador de entrada del historial de Sprints (HU-12).
type HistorialHandler struct {
	listar *service.ConsultarHistorialSprints
}

func NewHistorialHandler(listar *service.ConsultarHistorialSprints) *HistorialHandler {
	return &HistorialHandler{listar: listar}
}

// Listar atiende GET /projects/{id}/historial.
func (h *HistorialHandler) Listar(w nethttp.ResponseWriter, r *nethttp.Request) {
}
