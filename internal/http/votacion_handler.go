package http

import (
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// VotacionHandler es el adaptador de entrada de la votación de Planning Poker.
type VotacionHandler struct {
	obtenerHistoria *service.ObtenerHistoria
	votar           *service.Votar
	consultar       *service.ConsultarVotacion
}

func NewVotacionHandler(obtenerHistoria *service.ObtenerHistoria, votar *service.Votar, consultar *service.ConsultarVotacion) *VotacionHandler {
	return &VotacionHandler{obtenerHistoria: obtenerHistoria, votar: votar, consultar: consultar}
}

// Votar atiende POST /backlog/{id}/votacion/votos.
func (h *VotacionHandler) Votar(w nethttp.ResponseWriter, r *nethttp.Request) {}

// Consultar atiende GET /backlog/{id}/votacion.
func (h *VotacionHandler) Consultar(w nethttp.ResponseWriter, r *nethttp.Request) {}
