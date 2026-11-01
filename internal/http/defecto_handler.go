package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// DefectoHandler es el adaptador de entrada para los defectos de un proyecto (HU-08).
type DefectoHandler struct {
	obtenerProyecto *service.ObtenerProyecto
	registrar       *service.RegistrarDefecto
	resolver        *service.ResolverDefecto
	listar          *service.ListarDefectos
}

func NewDefectoHandler(
	obtenerProyecto *service.ObtenerProyecto,
	registrar *service.RegistrarDefecto,
	resolver *service.ResolverDefecto,
	listar *service.ListarDefectos,
) *DefectoHandler {
	return &DefectoHandler{obtenerProyecto: obtenerProyecto, registrar: registrar, resolver: resolver, listar: listar}
}

// Registrar atiende POST /projects/{id}/defectos.
func (h *DefectoHandler) Registrar(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}

	// El proyecto inexistente domina (404) sobre cualquier validación del cuerpo.
	if _, err := h.obtenerProyecto.Ejecutar(r.Context(), id); err != nil {
		escribirError(w, err)
		return
	}

	var req registrarDefectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}
	detalle, err := h.registrar.Ejecutar(r.Context(), service.RegistrarDefectoInput{
		ProyectoID:  id,
		HistoriaID:  req.HistoriaID,
		Descripcion: req.Descripcion,
		Severidad:   req.Severidad,
	})
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusCreated, aRespuestaDefecto(detalle))
}

// Resolver atiende POST /defectos/{id}/resolver. No tiene cuerpo.
func (h *DefectoHandler) Resolver(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}

	detalle, err := h.resolver.Ejecutar(r.Context(), id)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aRespuestaDefecto(detalle))
}

// Listar atiende GET /projects/{id}/defectos.
func (h *DefectoHandler) Listar(w nethttp.ResponseWriter, r *nethttp.Request) {
}
