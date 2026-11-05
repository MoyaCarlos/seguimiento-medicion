package http

import (
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// HistorialHandler es el adaptador de entrada del historial de Sprints (HU-12).
type HistorialHandler struct {
	listar  *service.ConsultarHistorialSprints
	detalle *service.ConsultarDetalleSprint
}

func NewHistorialHandler(listar *service.ConsultarHistorialSprints, detalle *service.ConsultarDetalleSprint) *HistorialHandler {
	return &HistorialHandler{listar: listar, detalle: detalle}
}

// Listar atiende GET /projects/{id}/historial.
func (h *HistorialHandler) Listar(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}
	historial, err := h.listar.Ejecutar(r.Context(), id)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aHistorialResponse(historial))
}

// Detalle atiende GET /projects/{id}/historial/{sprintId}.
func (h *HistorialHandler) Detalle(w nethttp.ResponseWriter, r *nethttp.Request) {
}
