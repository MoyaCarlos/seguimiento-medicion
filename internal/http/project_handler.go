package http

import (
	"encoding/json"
	"errors"
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// ProjectHandler es el adaptador de entrada de los proyectos y su equipo.
type ProjectHandler struct {
	crear   *service.CrearProyecto
	obtener *service.ObtenerProyecto
	editar  *service.EditarProyecto
	asignar *service.AsignarIntegrante
	listar  *service.ListarIntegrantes
}

func NewProjectHandler(
	crear *service.CrearProyecto,
	obtener *service.ObtenerProyecto,
	editar *service.EditarProyecto,
	asignar *service.AsignarIntegrante,
	listar *service.ListarIntegrantes,
) *ProjectHandler {
	return &ProjectHandler{crear: crear, obtener: obtener, editar: editar, asignar: asignar, listar: listar}
}

// Crear atiende POST /projects.
func (h *ProjectHandler) Crear(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req crearProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}

	inicio, err := parsearFecha(req.FechaInicio, "fecha_inicio")
	if err != nil {
		escribirError(w, err)
		return
	}
	fin, err := parsearFecha(req.FechaFin, "fecha_fin")
	if err != nil {
		escribirError(w, err)
		return
	}

	proyecto, err := h.crear.Ejecutar(r.Context(), service.CrearProyectoInput{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		FechaInicio: inicio,
		FechaFin:    fin,
		Creador:     req.Creador,
	})
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusCreated, aProyectoResponse(proyecto))
}

// Obtener atiende GET /projects/{id}.
func (h *ProjectHandler) Obtener(w nethttp.ResponseWriter, r *nethttp.Request) {
	proyecto, err := h.obtener.Ejecutar(r.Context(), r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aProyectoResponse(proyecto))
}

// Editar atiende PUT /projects/{id}.
func (h *ProjectHandler) Editar(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req editarProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}

	inicio, err := parsearFecha(req.FechaInicio, "fecha_inicio")
	if err != nil {
		escribirError(w, err)
		return
	}
	fin, err := parsearFecha(req.FechaFin, "fecha_fin")
	if err != nil {
		escribirError(w, err)
		return
	}

	proyecto, err := h.editar.Ejecutar(r.Context(), service.EditarProyectoInput{
		ID:          r.PathValue("id"),
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		FechaInicio: inicio,
		FechaFin:    fin,
	})
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aProyectoResponse(proyecto))
}

// AsignarIntegrante atiende POST /projects/{id}/members.
func (h *ProjectHandler) AsignarIntegrante(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req asignarIntegranteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}

	miembro, err := h.asignar.Ejecutar(r.Context(), service.AsignarIntegranteInput{
		ProyectoID: r.PathValue("id"),
		Nombre:     req.Nombre,
		Rol:        domain.Role(req.Rol),
	})
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusCreated, aIntegranteResponse(miembro))
}

// ListarIntegrantes atiende GET /projects/{id}/members.
func (h *ProjectHandler) ListarIntegrantes(w nethttp.ResponseWriter, r *nethttp.Request) {
	miembros, err := h.listar.Ejecutar(r.Context(), r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}
	respuesta := make([]integranteResponse, 0, len(miembros))
	for _, m := range miembros {
		respuesta = append(respuesta, aIntegranteResponse(m))
	}
	escribirJSON(w, nethttp.StatusOK, respuesta)
}

// escribirError traduce los errores de dominio a respuestas HTTP.
func escribirError(w nethttp.ResponseWriter, err error) {
	var verr domain.ValidationError
	if errors.As(err, &verr) {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Campo: verr.Campo, Mensaje: verr.Mensaje})
		return
	}
	if errors.Is(err, domain.ErrNoEncontrado) {
		escribirJSON(w, nethttp.StatusNotFound, errorResponse{Mensaje: "proyecto no encontrado"})
		return
	}
	escribirJSON(w, nethttp.StatusInternalServerError, errorResponse{Mensaje: "error interno"})
}
