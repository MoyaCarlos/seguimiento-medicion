package http

import (
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// HorasHandler es el adaptador de entrada para las horas imputadas a una historia.
type HorasHandler struct {
	imputar         *service.ImputarHoras
	obtenerHistoria *service.ObtenerHistoria
}

func NewHorasHandler(imputar *service.ImputarHoras, obtenerHistoria *service.ObtenerHistoria) *HorasHandler {
	return &HorasHandler{imputar: imputar, obtenerHistoria: obtenerHistoria}
}

// Imputar atiende POST /backlog/{id}/horas.
func (h *HorasHandler) Imputar(w nethttp.ResponseWriter, r *nethttp.Request) {
}
