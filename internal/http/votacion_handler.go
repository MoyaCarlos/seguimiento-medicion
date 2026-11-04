package http

import (
	"encoding/json"
	nethttp "net/http"
	"strconv"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
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

// Registrar agrega las rutas de la votación al mux.
func (h *VotacionHandler) Registrar(mux *nethttp.ServeMux) {
	mux.HandleFunc("GET /backlog/{id}/votacion", h.Consultar)
	mux.HandleFunc("POST /backlog/{id}/votacion", h.Iniciar)
	mux.HandleFunc("POST /backlog/{id}/votacion/votos", h.Votar)
	mux.HandleFunc("POST /backlog/{id}/votacion/revelar", h.Revelar)
	mux.HandleFunc("POST /backlog/{id}/votacion/estimacion", h.RegistrarEstimacion)
	mux.HandleFunc("POST /backlog/{id}/votacion/rondas", h.AbrirRonda)
	mux.HandleFunc("GET /backlog/{id}/votacion/rondas", h.Historial)
}

// Iniciar atiende POST /backlog/{id}/votacion. Responde con la ronda 1 tal
// como la ve quien la inició.
func (h *VotacionHandler) Iniciar(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, ok := h.historiaDeRuta(w, r)
	if !ok {
		return
	}
	var req actorRequest
	if !leerPedido(w, r, &req) {
		return
	}

	if _, err := h.iniciar.Ejecutar(r.Context(), id, req.UsuarioID); err != nil {
		escribirError(w, err)
		return
	}
	h.responderVista(w, r, nethttp.StatusCreated, id, req.UsuarioID)
}

// Votar atiende POST /backlog/{id}/votacion/votos. Responde con la vista de
// quien votó.
func (h *VotacionHandler) Votar(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, ok := h.historiaDeRuta(w, r)
	if !ok {
		return
	}
	var req votoRequest
	if !leerPedido(w, r, &req) {
		return
	}

	if _, err := h.votar.Ejecutar(r.Context(), service.VotarInput{HistoriaID: id, UsuarioID: req.UsuarioID, Valor: req.valor()}); err != nil {
		escribirError(w, err)
		return
	}
	h.responderVista(w, r, nethttp.StatusOK, id, req.UsuarioID)
}

// Revelar atiende POST /backlog/{id}/votacion/revelar.
func (h *VotacionHandler) Revelar(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, ok := h.historiaDeRuta(w, r)
	if !ok {
		return
	}
	var req actorRequest
	if !leerPedido(w, r, &req) {
		return
	}

	if _, err := h.revelar.Ejecutar(r.Context(), id, req.UsuarioID, time.Now()); err != nil {
		escribirError(w, err)
		return
	}
	h.responderVista(w, r, nethttp.StatusOK, id, req.UsuarioID)
}

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
	h.responderVista(w, r, nethttp.StatusOK, id, consultanteID)
}

// RegistrarEstimacion atiende POST /backlog/{id}/votacion/estimacion. Responde
// con la historia.
func (h *VotacionHandler) RegistrarEstimacion(w nethttp.ResponseWriter, r *nethttp.Request) {
	historia, ok := h.cargarHistoria(w, r)
	if !ok {
		return
	}
	var req estimacionRequest
	if !leerPedido(w, r, &req) {
		return
	}

	actualizada, err := h.registrar.Ejecutar(r.Context(), historia, req.UsuarioID, req.valor())
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aRespuesta(actualizada))
}

// AbrirRonda atiende POST /backlog/{id}/votacion/rondas. Responde con la ronda
// nueva tal como la ve quien la abrió.
func (h *VotacionHandler) AbrirRonda(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, ok := h.historiaDeRuta(w, r)
	if !ok {
		return
	}
	var req actorRequest
	if !leerPedido(w, r, &req) {
		return
	}

	if _, err := h.abrir.Ejecutar(r.Context(), id, req.UsuarioID); err != nil {
		escribirError(w, err)
		return
	}
	h.responderVista(w, r, nethttp.StatusCreated, id, req.UsuarioID)
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

func (h *VotacionHandler) historiaDeRuta(w nethttp.ResponseWriter, r *nethttp.Request) (int64, bool) {
	historia, ok := h.cargarHistoria(w, r)
	return historia.ID, ok
}

// cargarHistoria valida el id de la ruta y que la historia exista: la historia
// inexistente domina (404) sobre cualquier validación del pedido (FR-011).
func (h *VotacionHandler) cargarHistoria(w nethttp.ResponseWriter, r *nethttp.Request) (domain.BacklogItem, bool) {
	id, ok := idDeRuta(w, r)
	if !ok {
		return domain.BacklogItem{}, false
	}
	historia, err := h.obtenerHistoria.Ejecutar(r.Context(), id)
	if err != nil {
		escribirError(w, err)
		return domain.BacklogItem{}, false
	}
	return historia, true
}

func (h *VotacionHandler) responderVista(w nethttp.ResponseWriter, r *nethttp.Request, status int, historiaID, consultanteID int64) {
	vista, err := h.consultar.Ejecutar(r.Context(), historiaID, consultanteID)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, status, aRespuestaVotacion(vista))
}

// leerPedido decodifica el cuerpo y exige quien actúa (usuario_id > 0).
func leerPedido(w nethttp.ResponseWriter, r *nethttp.Request, req interface{ usuario() int64 }) bool {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Mensaje: "JSON inválido"})
		return false
	}
	if req.usuario() <= 0 {
		escribirUsuarioInvalido(w)
		return false
	}
	return true
}

func escribirUsuarioInvalido(w nethttp.ResponseWriter) {
	escribirJSON(w, nethttp.StatusBadRequest, errorResponse{Campo: "usuario_id", Mensaje: "el usuario es obligatorio"})
}
