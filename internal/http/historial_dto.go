package http

import (
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type sprintHistorialResponse struct {
	ID          int64            `json:"id"`
	ProyectoID  int64            `json:"proyecto_id"`
	Numero      int              `json:"numero"`
	SprintGoal  string           `json:"sprint_goal"`
	FechaInicio string           `json:"fecha_inicio"`
	FechaFin    string           `json:"fecha_fin"`
	Metricas    metricasResponse `json:"metricas"`
}

func aSprintHistorialResponse(h domain.SprintHistorial) sprintHistorialResponse {
	return sprintHistorialResponse{
		ID:          h.Sprint.ID,
		ProyectoID:  h.Sprint.ProyectoID,
		Numero:      h.Sprint.Numero,
		SprintGoal:  h.Sprint.SprintGoal,
		FechaInicio: h.Sprint.FechaInicio.Format(time.DateOnly),
		FechaFin:    h.Sprint.FechaFin.Format(time.DateOnly),
		Metricas:    aMetricasResponse(h.Metricas),
	}
}

func aHistorialResponse(historial []domain.SprintHistorial) []sprintHistorialResponse {
	respuesta := make([]sprintHistorialResponse, 0, len(historial))
	for _, h := range historial {
		respuesta = append(respuesta, aSprintHistorialResponse(h))
	}
	return respuesta
}

type historiaComprometidaResponse struct {
	ID           int64   `json:"id"`
	Titulo       string  `json:"titulo"`
	EstimacionSP *int    `json:"estimacion_sp"`
	Completada   bool    `json:"completada"`
	Minutos      int     `json:"minutos"`
	Horas        float64 `json:"horas"`
}

type defectoDeSprintResponse struct {
	ID                int64  `json:"id"`
	HistoriaID        int64  `json:"historia_id"`
	Historia          string `json:"historia"`
	Descripcion       string `json:"descripcion"`
	Severidad         string `json:"severidad"`
	Estado            string `json:"estado"`
	DetectadoEnSprint bool   `json:"detectado_en_sprint"`
	ResueltoEnSprint  bool   `json:"resuelto_en_sprint"`
}

type detalleSprintResponse struct {
	sprintHistorialResponse
	Historias []historiaComprometidaResponse `json:"historias"`
	Defectos  []defectoDeSprintResponse      `json:"defectos"`
}

// Las horas ya vienen redondeadas del dominio: el DTO no las recalcula.
func aDetalleSprintResponse(d domain.DetalleSprint) detalleSprintResponse {
	historias := make([]historiaComprometidaResponse, 0, len(d.Historias))
	for _, h := range d.Historias {
		historias = append(historias, historiaComprometidaResponse{
			ID: h.ID, Titulo: h.Titulo, EstimacionSP: h.EstimacionSP, Completada: h.Completada, Minutos: h.Minutos, Horas: h.Horas,
		})
	}
	defectos := make([]defectoDeSprintResponse, 0, len(d.Defectos))
	for _, def := range d.Defectos {
		defectos = append(defectos, defectoDeSprintResponse{
			ID: def.ID, HistoriaID: def.HistoriaID, Historia: def.TituloHistoria, Descripcion: def.Descripcion,
			Severidad: string(def.Severidad), Estado: string(def.Estado),
			DetectadoEnSprint: def.DetectadoEnSprint, ResueltoEnSprint: def.ResueltoEnSprint,
		})
	}
	return detalleSprintResponse{
		sprintHistorialResponse: aSprintHistorialResponse(d.SprintHistorial),
		Historias:               historias,
		Defectos:                defectos,
	}
}
