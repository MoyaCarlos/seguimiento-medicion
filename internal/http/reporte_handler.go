package http

import (
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// generadorReportes convierte el agregado del reporte en bytes PDF.
type generadorReportes interface {
	GenerarSprint(domain.ReporteSprint) ([]byte, error)
	GenerarProyecto(domain.ReporteProyecto) ([]byte, error)
}

// ReporteHandler es el adaptador de entrada de los reportes en PDF (HU-10).
type ReporteHandler struct {
	generar *service.GenerarReporte
	pdf     generadorReportes
}

func NewReporteHandler(generar *service.GenerarReporte, pdf generadorReportes) *ReporteHandler {
	return &ReporteHandler{generar: generar, pdf: pdf}
}

// Sprint atiende GET /sprints/{id}/reporte.
func (h *ReporteHandler) Sprint(_ nethttp.ResponseWriter, _ *nethttp.Request) {
}

// Proyecto atiende GET /projects/{id}/reporte.
func (h *ReporteHandler) Proyecto(_ nethttp.ResponseWriter, _ *nethttp.Request) {
}
