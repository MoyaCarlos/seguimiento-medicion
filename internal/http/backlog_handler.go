package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// BacklogHandler es el adaptador de entrada para las historias del backlog.
type BacklogHandler struct {
	crearHistoria *service.CrearHistoriaBacklog
}

func NewBacklogHandler(crearHistoria *service.CrearHistoriaBacklog) *BacklogHandler {
	return &BacklogHandler{crearHistoria: crearHistoria}
}

// Crear atiende POST /backlog.
func (h *BacklogHandler) Crear(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req crearHistoriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}

	item, err := h.crearHistoria.Ejecutar(r.Context(), service.CrearHistoriaInput{
		ProyectoID:   req.ProyectoID,
		Titulo:       req.Titulo,
		Descripcion:  req.Descripcion,
		Prioridad:    domain.Prioridad(req.Prioridad),
		ValorNegocio: req.ValorNegocio,
	})
	if err != nil {
		escribirError(w, err)
		return
	}

	escribirJSON(w, nethttp.StatusCreated, aRespuesta(item))
}

func escribirJSON(w nethttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
