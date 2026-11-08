package http

import "github.com/MoyaCarlos/seguimiento-medicion/internal/domain"

type metricasResponse struct {
	StoryPointsPlanificados        int      `json:"story_points_planificados"`
	StoryPointsCompletados         int      `json:"story_points_completados"`
	Velocidad                      float64  `json:"velocidad"`
	PorcentajeHistoriasCompletadas float64  `json:"porcentaje_historias_completadas"`
	HorasEstimadas                 *float64 `json:"horas_estimadas"`
	HorasReales                    float64  `json:"horas_reales"`
	DesviacionHoras                *float64 `json:"desviacion_horas"`
	DesviacionPorcentaje           *float64 `json:"desviacion_porcentaje"`
	DefectosDetectados             int      `json:"defectos_detectados"`
	DefectosResueltos              int      `json:"defectos_resueltos"`
}

type metricasDeSprintResponse struct {
	SprintID int64            `json:"sprint_id"`
	Numero   int              `json:"numero"`
	Metricas metricasResponse `json:"metricas"`
}

type metricasProyectoResponse struct {
	ProyectoID         int64                      `json:"proyecto_id"`
	HorasPorStoryPoint *float64                   `json:"horas_por_story_point"`
	Proyecto           metricasResponse           `json:"proyecto"`
	Sprints            []metricasDeSprintResponse `json:"sprints"`
}

func aMetricasResponse(m domain.Metricas) metricasResponse {
	return metricasResponse{
		StoryPointsPlanificados:        m.SPPlanificados,
		StoryPointsCompletados:         m.SPCompletados,
		Velocidad:                      m.Velocidad,
		PorcentajeHistoriasCompletadas: m.PorcentajeCompletadas,
		HorasEstimadas:                 m.HorasEstimadas,
		HorasReales:                    m.HorasReales,
		DesviacionHoras:                m.DesviacionHoras,
		DesviacionPorcentaje:           m.DesviacionPorcentaje,
		DefectosDetectados:             m.DefectosDetectados,
		DefectosResueltos:              m.DefectosResueltos,
	}
}

func aMetricasProyectoResponse(proyectoID int64, m domain.MetricasProyecto) metricasProyectoResponse {
	sprints := make([]metricasDeSprintResponse, 0, len(m.Sprints))
	for _, s := range m.Sprints {
		sprints = append(sprints, metricasDeSprintResponse{SprintID: s.SprintID, Numero: s.Numero, Metricas: aMetricasResponse(s.Metricas)})
	}
	return metricasProyectoResponse{
		ProyectoID:         proyectoID,
		HorasPorStoryPoint: m.HorasPorStoryPoint,
		Proyecto:           aMetricasResponse(m.Proyecto),
		Sprints:            sprints,
	}
}
