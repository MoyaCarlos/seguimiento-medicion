package http

import (
	"encoding/json"
	nethttp "net/http"
	"time"

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
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}

	var req imputarHorasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}

	detalle, err := h.imputar.Ejecutar(r.Context(), service.ImputarHorasInput{
		HistoriaID:   id,
		IntegranteID: req.IntegranteID,
		Fecha:        req.Fecha,
		Actividad:    req.Actividad,
		Horas:        string(req.Horas),
	}, time.Now())
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusCreated, aRespuestaHoras(detalle))
}
