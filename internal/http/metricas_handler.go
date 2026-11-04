package http

import (
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// MetricasHandler es el adaptador de entrada de las métricas de un proyecto.
type MetricasHandler struct {
	calcular *service.CalcularMetricas
}

func NewMetricasHandler(calcular *service.CalcularMetricas) *MetricasHandler {
	return &MetricasHandler{calcular: calcular}
}

// Consultar atiende GET /projects/{id}/metricas.
func (h *MetricasHandler) Consultar(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}
	metricas, err := h.calcular.Ejecutar(r.Context(), id)
	if err != nil {
		escribirError(w, err)
		return
	}
	escribirJSON(w, nethttp.StatusOK, aMetricasProyectoResponse(id, metricas))
}
