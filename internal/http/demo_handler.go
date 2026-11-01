package http

import (
	"context"
	"errors"
	"log"
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/demo"
)

// CargadorDatosPrueba carga el proyecto de demostración.
type CargadorDatosPrueba interface {
	Cargar(ctx context.Context) (demo.Resumen, error)
}

// DemoHandler es el adaptador de entrada de la carga de datos de prueba. Es el
// único lugar que decide si las rutas /demo/* existen: hoy lo define una
// variable de entorno; más adelante, un chequeo de perfil.
type DemoHandler struct {
	habilitado func() bool
	cargador   CargadorDatosPrueba
}

func NewDemoHandler(habilitado func() bool, cargador CargadorDatosPrueba) *DemoHandler {
	return &DemoHandler{habilitado: habilitado, cargador: cargador}
}

// Registrar agrega GET /demo/estado y POST /demo/datos al mux. Cualquier otra
// combinación de método y ruta bajo /demo/ responde 404, también con la
// función apagada.
func (h *DemoHandler) Registrar(mux *nethttp.ServeMux) {
	mux.HandleFunc("GET /demo/estado", h.soloSiHabilitado(h.Estado))
	mux.HandleFunc("POST /demo/datos", h.soloSiHabilitado(h.Cargar))
	mux.HandleFunc("/demo/", nethttp.NotFound)
}

func (h *DemoHandler) soloSiHabilitado(siguiente nethttp.HandlerFunc) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if !h.habilitado() {
			nethttp.NotFound(w, r)
			return
		}
		siguiente(w, r)
	}
}

// Estado atiende GET /demo/estado.
func (h *DemoHandler) Estado(w nethttp.ResponseWriter, _ *nethttp.Request) {
	escribirJSON(w, nethttp.StatusOK, map[string]bool{"habilitado": true})
}

// Cargar atiende POST /demo/datos.
func (h *DemoHandler) Cargar(w nethttp.ResponseWriter, r *nethttp.Request) {
	resumen, err := h.cargador.Cargar(r.Context())
	switch {
	case err == nil:
		escribirJSON(w, nethttp.StatusCreated, resumen)
	case errors.Is(err, demo.ErrYaCargado):
		escribirJSON(w, nethttp.StatusConflict, errorResponse{Campo: "ya_cargado", Mensaje: err.Error()})
	case errors.Is(err, demo.ErrIdOcupado):
		escribirJSON(w, nethttp.StatusConflict, errorResponse{Campo: "proyecto_id", Mensaje: err.Error()})
	default:
		log.Printf("cargar datos de prueba: %v", err)
		escribirJSON(w, nethttp.StatusInternalServerError, errorResponse{Mensaje: "no se pudieron cargar los datos de prueba"})
	}
}
