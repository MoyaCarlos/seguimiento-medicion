package http

import (
	"encoding/json"
	nethttp "net/http"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// VotacionHandler es el adaptador de entrada de la votación de Planning Poker.
type VotacionHandler struct {
	obtenerHistoria *service.ObtenerHistoria
	votar           *service.Votar
	consultar       *service.ConsultarVotacion
	revelar         *service.RevelarVotacion
	iniciar         *service.IniciarVotacion
	registrar       *service.RegistrarEstimacion
	abrir           *service.AbrirNuevaRonda
	historial       *service.ConsultarHistorialRondas
}

func NewVotacionHandler(obtenerHistoria *service.ObtenerHistoria, votar *service.Votar, consultar *service.ConsultarVotacion, revelar *service.RevelarVotacion, iniciar *service.IniciarVotacion, registrar *service.RegistrarEstimacion, abrir *service.AbrirNuevaRonda, historial *service.ConsultarHistorialRondas) *VotacionHandler {
	return &VotacionHandler{obtenerHistoria: obtenerHistoria, votar: votar, consultar: consultar, revelar: revelar, iniciar: iniciar, registrar: registrar, abrir: abrir, historial: historial}
}

// Iniciar atiende POST /backlog/{id}/votacion. Responde con la ronda 1 tal
// como la ve quien la inició.
func (h *VotacionHandler) Iniciar(w nethttp.ResponseWriter, r *nethttp.Request) {
	actor, id, ok := h.actorEHistoria(w, r)
	if !ok {
		return
	}
	if _, err := h.iniciar.Ejecutar(r.Context(), id, actor); err != nil {
		escribirError(w, err)
		return
	}
	h.responderVista(w, r, nethttp.StatusCreated, id, actor)
}

// Votar atiende POST /backlog/{id}/votacion/votos. El voto queda a nombre de la
// sesión y la respuesta es la vista de quien votó.
func (h *VotacionHandler) Votar(w nethttp.ResponseWriter, r *nethttp.Request) {
	actor, id, ok := h.actorEHistoria(w, r)
	if !ok {
		return
	}
	var req votoRequest
	if !leerCuerpo(w, r, &req) {
		return
	}
	if _, err := h.votar.Ejecutar(r.Context(), service.VotarInput{HistoriaID: id, UsuarioID: actor, Valor: req.valor()}); err != nil {
		escribirError(w, err)
		return
	}
	h.responderVista(w, r, nethttp.StatusOK, id, actor)
}

// Revelar atiende POST /backlog/{id}/votacion/revelar.
func (h *VotacionHandler) Revelar(w nethttp.ResponseWriter, r *nethttp.Request) {
	actor, id, ok := h.actorEHistoria(w, r)
	if !ok {
		return
	}
	if _, err := h.revelar.Ejecutar(r.Context(), id, actor, time.Now()); err != nil {
		escribirError(w, err)
		return
	}
	h.responderVista(w, r, nethttp.StatusOK, id, actor)
}

// Consultar atiende GET /backlog/{id}/votacion: la vista de quien tiene la
// sesión. Un usuario_id en la query se ignora (escenario 7 de HU-21).
func (h *VotacionHandler) Consultar(w nethttp.ResponseWriter, r *nethttp.Request) {
	actor, id, ok := h.actorEHistoria(w, r)
	if !ok {
		return
	}
	h.responderVista(w, r, nethttp.StatusOK, id, actor)
}

// RegistrarEstimacion atiende POST /backlog/{id}/votacion/estimacion. Responde
// con la historia.
func (h *VotacionHandler) RegistrarEstimacion(w nethttp.ResponseWriter, r *nethttp.Request) {
	actor, id, ok := h.actorEHistoria(w, r)
	if !ok {
		return
	}
	var req estimacionRequest
	if !leerCuerpo(w, r, &req) {
		return
	}
	historia, err := h.registrar.Ejecutar(r.Context(), id, actor, req.valor())
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aRespuesta(historia))
}

// AbrirRonda atiende POST /backlog/{id}/votacion/rondas. Responde con la ronda
// nueva tal como la ve quien la abrió.
func (h *VotacionHandler) AbrirRonda(w nethttp.ResponseWriter, r *nethttp.Request) {
	actor, id, ok := h.actorEHistoria(w, r)
	if !ok {
		return
	}
	if _, err := h.abrir.Ejecutar(r.Context(), id, actor); err != nil {
		escribirError(w, err)
		return
	}
	h.responderVista(w, r, nethttp.StatusCreated, id, actor)
}

// Historial atiende GET /backlog/{id}/votacion/rondas.
func (h *VotacionHandler) Historial(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, ok := h.historiaDeRuta(w, r)
	if !ok {
		return
	}
	historial, err := h.historial.Ejecutar(r.Context(), id)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aRespuestaHistorial(historial))
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

func (h *VotacionHandler) responderVista(w nethttp.ResponseWriter, r *nethttp.Request, status int, historiaID, consultanteID int64) {
	vista, err := h.consultar.Ejecutar(r.Context(), historiaID, consultanteID)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, status, aRespuestaVotacion(vista))
}

// actorEHistoria resuelve quién actúa (401 sin sesión) y la historia de la ruta.
func (h *VotacionHandler) actorEHistoria(w nethttp.ResponseWriter, r *nethttp.Request) (int64, int64, bool) {
	actor, ok := actorDeSesion(w, r)
	if !ok {
		return 0, 0, false
	}
	id, ok := h.historiaDeRuta(w, r)
	return actor, id, ok
}

// leerCuerpo decodifica el cuerpo JSON; si no es válido responde 400.
func leerCuerpo(w nethttp.ResponseWriter, r *nethttp.Request, req any) bool {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return false
	}
	return true
}
