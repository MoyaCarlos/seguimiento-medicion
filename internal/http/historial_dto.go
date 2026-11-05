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
