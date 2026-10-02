package http

import (
	"encoding/json"
	"errors"
	nethttp "net/http"
	"strconv"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// SprintHandler es el adaptador de entrada para el ciclo de vida de los Sprints.
type SprintHandler struct {
	iniciar *service.IniciarSprint
}

func NewSprintHandler(iniciar *service.IniciarSprint) *SprintHandler {
	return &SprintHandler{iniciar: iniciar}
}

// Iniciar atiende POST /sprints/{id}/iniciar.
func (h *SprintHandler) Iniciar(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, ok := idDeRuta(w, r)
	if !ok {
		return
	}
	var req iniciarSprintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}
	inicio, ok := fechaDeCampo(w, "fecha_inicio", req.FechaInicio)
	if !ok {
		return
	}
	fin, ok := fechaDeCampo(w, "fecha_fin", req.FechaFin)
	if !ok {
		return
	}

	sprint, err := h.iniciar.Ejecutar(r.Context(), service.IniciarSprintInput{
		SprintID:    id,
		SprintGoal:  req.SprintGoal,
		FechaInicio: inicio,
		FechaFin:    fin,
	})
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aRespuestaSprint(sprint))
}

func idDeRuta(w nethttp.ResponseWriter, r *nethttp.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Campo: "id", Mensaje: "identificador inválido"})
		return 0, false
	}
	return id, true
}

func fechaDeCampo(w nethttp.ResponseWriter, campo, valor string) (time.Time, bool) {
	fecha, err := time.Parse(time.DateOnly, valor)
	if err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Campo: campo, Mensaje: "formato esperado AAAA-MM-DD"})
		return time.Time{}, false
	}
	return fecha, true
}

// escribirError traduce los errores del dominio a respuestas HTTP.
func escribirError(w nethttp.ResponseWriter, err error) {
	var verr domain.ValidationError
	switch {
	case errors.As(err, &verr):
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Campo: verr.Campo, Mensaje: verr.Mensaje})
	case errors.Is(err, domain.ErrSprintNoEncontrado):
		escribirJSON(w, nethttp.StatusNotFound, errorResponse{Mensaje: err.Error()})
	default:
		escribirJSON(w, nethttp.StatusInternalServerError, errorResponse{Mensaje: "error interno"})
	}
}
