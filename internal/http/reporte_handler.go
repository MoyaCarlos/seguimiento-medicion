package http

import (
	"fmt"
	nethttp "net/http"
	"time"

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
func (h *ReporteHandler) Sprint(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}
	reporte, err := h.generar.EjecutarSprint(r.Context(), id, time.Now())
	if err != nil {
		escribirError(w, err)
		return
	}
	contenido, err := h.pdf.GenerarSprint(reporte)
	if err != nil {
		escribirError(w, domain.ErrReporteNoGenerado)
		return
	}
	nombre := fmt.Sprintf("reporte-proyecto-%d-sprint-%d.pdf", reporte.Proyecto.ID, reporte.Sprint.Sprint.Numero)
	escribirPDF(w, nombre, contenido)
}

// Proyecto atiende GET /projects/{id}/reporte.
func (h *ReporteHandler) Proyecto(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, err := parsearID(r.PathValue("id"))
	if err != nil {
		escribirError(w, err)
		return
	}
	reporte, err := h.generar.EjecutarProyecto(r.Context(), id, time.Now())
	if err != nil {
		escribirError(w, err)
		return
	}
	contenido, err := h.pdf.GenerarProyecto(reporte)
	if err != nil {
		escribirError(w, domain.ErrReporteNoGenerado)
		return
	}
	nombre := fmt.Sprintf("reporte-proyecto-%d.pdf", reporte.Proyecto.ID)
	escribirPDF(w, nombre, contenido)
}

func escribirPDF(w nethttp.ResponseWriter, nombre string, contenido []byte) {
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", nombre))
	w.WriteHeader(nethttp.StatusOK)
	_, _ = w.Write(contenido)
}
