package http

import (
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// AvanceHandler es el adaptador de entrada del avance de las historias (HU-18).
type AvanceHandler struct {
	consultar *service.ConsultarAvanceHistorias
}

func NewAvanceHandler(consultar *service.ConsultarAvanceHistorias) *AvanceHandler {
	return &AvanceHandler{consultar: consultar}
}

// Consultar atiende GET /projects/{id}/avance.
func (h *AvanceHandler) Consultar(w nethttp.ResponseWriter, r *nethttp.Request) {
}
