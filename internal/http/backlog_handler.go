package http

import (
	"encoding/json"
	"io"
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
	asignarSprint    *service.AsignarHistoriaASprint
}

func NewBacklogHandler(
	crearHistoria *service.CrearHistoriaBacklog,
	consultarBacklog *service.ConsultarProductBacklog,
	obtenerHistoria *service.ObtenerHistoria,
	editarHistoria *service.EditarHistoria,
	asignarSprint *service.AsignarHistoriaASprint,
) *BacklogHandler {
	return &BacklogHandler{
		crearHistoria:    crearHistoria,
		consultarBacklog: consultarBacklog,
		obtenerHistoria:  obtenerHistoria,
		editarHistoria:   editarHistoria,
		asignarSprint:    asignarSprint,
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
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}

	// La historia inexistente domina (404) sobre cualquier validación del cuerpo.
	if _, err := h.obtenerHistoria.Ejecutar(r.Context(), id); err != nil {
		escribirError(w, err)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}
	var claves map[string]json.RawMessage
	if err := json.Unmarshal(data, &claves); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}
	for _, campo := range []string{"titulo", "descripcion", "criterios_aceptacion", "estimacion_sp"} {
		if _, ok := claves[campo]; !ok {
			escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Campo: campo, Mensaje: "el campo es obligatorio"})
			return
		}
	}

	var req editarHistoriaRequest
	if err := json.Unmarshal(data, &req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}
	estimacion, err := parsearEstimacion(req.EstimacionSP)
	if err != nil {
		escribirError(w, err)
		return
	}

	item, err := h.editarHistoria.Ejecutar(r.Context(), service.EditarHistoriaInput{
		ID:                  id,
		Titulo:              req.Titulo,
		Descripcion:         req.Descripcion,
		CriteriosAceptacion: req.CriteriosAceptacion,
		EstimacionSP:        estimacion,
	})
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aRespuesta(item))
}

// AsignarSprint atiende PUT /backlog/{id}/sprint.
func (h *BacklogHandler) AsignarSprint(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}

	// La historia inexistente domina (404) sobre cualquier validación del cuerpo.
	if _, err := h.obtenerHistoria.Ejecutar(r.Context(), id); err != nil {
		escribirError(w, err)
		return
	}

	var req asignarSprintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}
	if req.SprintID == nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Campo: "sprint_id", Mensaje: "el campo es obligatorio"})
		return
	}

	item, err := h.asignarSprint.Ejecutar(r.Context(), id, *req.SprintID)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aRespuesta(item))
}

// parsearEstimacion convierte los Story Points del cuerpo a entero; un número no
// entero se rechaza con el mismo error que un valor fuera de la escala.
func parsearEstimacion(valor *json.Number) (*int, error) {
	if valor == nil {
		return nil, nil
	}
	n, err := valor.Int64()
	if err != nil {
		return nil, domain.ErrEstimacionInvalida
	}
	estimacion := int(n)
	return &estimacion, nil
}

func escribirJSON(w nethttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
