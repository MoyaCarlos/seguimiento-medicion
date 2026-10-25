package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// BacklogHandler es el adaptador de entrada para las historias del backlog.
type BacklogHandler struct {
	crearHistoria    *service.CrearHistoriaBacklog
	consultarBacklog *service.ConsultarProductBacklog
	obtenerHistoria  *service.ObtenerHistoria
	editarHistoria   *service.EditarHistoria
}

func NewBacklogHandler(
	crearHistoria *service.CrearHistoriaBacklog,
	consultarBacklog *service.ConsultarProductBacklog,
	obtenerHistoria *service.ObtenerHistoria,
	editarHistoria *service.EditarHistoria,
) *BacklogHandler {
	return &BacklogHandler{
		crearHistoria:    crearHistoria,
		consultarBacklog: consultarBacklog,
		obtenerHistoria:  obtenerHistoria,
		editarHistoria:   editarHistoria,
	}
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

// ConsultarBacklog atiende GET /projects/{id}/backlog.
func (h *BacklogHandler) ConsultarBacklog(w nethttp.ResponseWriter, r *nethttp.Request) {
	proyectoID, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}

	historias, err := h.consultarBacklog.Ejecutar(r.Context(), proyectoID)
	if err != nil {
		escribirError(w, err)
		return
	}

	respuesta := make([]historiaResponse, 0, len(historias))
	for _, item := range historias {
		respuesta = append(respuesta, aRespuesta(item))
	}
	escribirJSON(w, nethttp.StatusOK, respuesta)
}

// Editar atiende PUT /backlog/{id}.
func (h *BacklogHandler) Editar(w nethttp.ResponseWriter, r *nethttp.Request) {
}

func escribirJSON(w nethttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
