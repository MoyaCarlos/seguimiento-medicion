package http

import (
	"encoding/json"
	"io"
	nethttp "net/http"
	"strconv"
	"time"

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
	estado  *service.ObtenerEstadoProyecto
}

func NewProjectHandler(
	crear *service.CrearProyecto,
	obtener *service.ObtenerProyecto,
	editar *service.EditarProyecto,
	asignar *service.AsignarIntegrante,
	listar *service.ListarIntegrantes,
	estado *service.ObtenerEstadoProyecto,
) *ProjectHandler {
	return &ProjectHandler{crear: crear, obtener: obtener, editar: editar, asignar: asignar, listar: listar, estado: estado}
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
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}
	proyecto, err := h.obtener.Ejecutar(r.Context(), id)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aProyectoResponse(proyecto))
}

// ObtenerEstado atiende GET /projects/{id}/status.
func (h *ProjectHandler) ObtenerEstado(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}
	estado, err := h.estado.Ejecutar(r.Context(), id, time.Now())
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, estadoProyectoResponse{Estado: string(estado)})
}

// Editar atiende PUT /projects/{id}.
func (h *ProjectHandler) Editar(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}

	// El proyecto inexistente domina (404) sobre cualquier validación del cuerpo.
	if _, err := h.obtener.Ejecutar(r.Context(), id); err != nil {
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
	for _, campo := range []string{"nombre", "descripcion", "fecha_inicio", "fecha_fin"} {
		if _, ok := claves[campo]; !ok {
			escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Campo: campo, Mensaje: "el campo es obligatorio"})
			return
		}
	}

	var req editarProyectoRequest
	if err := json.Unmarshal(data, &req); err != nil {
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
		ID:          id,
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
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}

	var req asignarIntegranteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}

	miembro, err := h.asignar.Ejecutar(r.Context(), service.AsignarIntegranteInput{
		ProyectoID: id,
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
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}

	miembros, err := h.listar.Ejecutar(r.Context(), id)
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

func parsearID(valor string) (int64, error) {
	id, err := strconv.ParseInt(valor, 10, 64)
	if err != nil || id <= 0 {
		return 0, domain.ValidationError{Campo: "id", Mensaje: "el identificador debe ser un entero positivo"}
	}
	return id, nil
}
