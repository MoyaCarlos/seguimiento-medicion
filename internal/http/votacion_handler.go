package http

import (
	"encoding/json"
	nethttp "net/http"
	"strconv"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// VotacionHandler es el adaptador de entrada de la votación de Planning Poker.
type VotacionHandler struct {
	obtenerHistoria *service.ObtenerHistoria
	votar           *service.Votar
	consultar       *service.ConsultarVotacion
	revelar         *service.RevelarVotacion
}

func NewVotacionHandler(obtenerHistoria *service.ObtenerHistoria, votar *service.Votar, consultar *service.ConsultarVotacion, revelar *service.RevelarVotacion) *VotacionHandler {
	return &VotacionHandler{obtenerHistoria: obtenerHistoria, votar: votar, consultar: consultar, revelar: revelar}
}

// Votar atiende POST /backlog/{id}/votacion/votos. Responde con la vista de
// quien votó.
func (h *VotacionHandler) Votar(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, ok := h.historiaDeRuta(w, r)
	if !ok {
		return
	}
	var req votoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return
	}
	if req.UsuarioID <= 0 {
		escribirUsuarioInvalido(w)
		return
	}

	if _, err := h.votar.Ejecutar(r.Context(), service.VotarInput{HistoriaID: id, UsuarioID: req.UsuarioID, Valor: req.Valor}); err != nil {
		escribirError(w, err)
		return
	}
	h.responderVista(w, r, id, req.UsuarioID)
}

// Revelar atiende POST /backlog/{id}/votacion/revelar.
func (h *VotacionHandler) Revelar(w nethttp.ResponseWriter, r *nethttp.Request) {}

// Consultar atiende GET /backlog/{id}/votacion?usuario_id=N. Sin usuario_id
// todos los votos de una ronda oculta quedan ocultos.
func (h *VotacionHandler) Consultar(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, ok := h.historiaDeRuta(w, r)
	if !ok {
		return
	}
	var consultanteID int64
	if valor := r.URL.Query().Get("usuario_id"); valor != "" {
		n, err := strconv.ParseInt(valor, 10, 64)
		if err != nil || n <= 0 {
			escribirUsuarioInvalido(w)
			return
		}
		consultanteID = n
	}
	h.responderVista(w, r, id, consultanteID)
}

// historiaDeRuta valida el id de la ruta y que la historia exista: la historia
// inexistente domina (404) sobre cualquier validación del pedido (FR-011).
func (h *VotacionHandler) historiaDeRuta(w nethttp.ResponseWriter, r *nethttp.Request) (int64, bool) {
	id, ok := idDeRuta(w, r)
	if !ok {
		return 0, false
	}
	if _, err := h.obtenerHistoria.Ejecutar(r.Context(), id); err != nil {
		escribirError(w, err)
		return 0, false
	}
	return id, true
}

func (h *VotacionHandler) responderVista(w nethttp.ResponseWriter, r *nethttp.Request, historiaID, consultanteID int64) {
	vista, err := h.consultar.Ejecutar(r.Context(), historiaID, consultanteID)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aRespuestaVotacion(vista))
}

func escribirUsuarioInvalido(w nethttp.ResponseWriter) {
	escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Campo: "usuario_id", Mensaje: "el usuario es obligatorio"})
}
