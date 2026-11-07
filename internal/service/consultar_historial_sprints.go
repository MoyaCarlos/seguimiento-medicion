package service

import (
	"context"
	"fmt"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// calculadorMetricas es el cálculo único de métricas de HU-03.
type calculadorMetricas interface {
	Ejecutar(ctx context.Context, proyectoID int64) (domain.MetricasProyecto, error)
}

// ConsultarHistorialSprints devuelve los Sprints finalizados de un proyecto
// con sus métricas de HU-03 (HU-12). Solo lee.
type ConsultarHistorialSprints struct {
	metricas calculadorMetricas
	sprints  listadorSprints
}

func NewConsultarHistorialSprints(metricas calculadorMetricas, sprints listadorSprints) *ConsultarHistorialSprints {
	return &ConsultarHistorialSprints{metricas: metricas, sprints: sprints}
}

// Ejecutar responde domain.ErrProyectoNoEncontrado si el proyecto no existe:
// lo verifica el cálculo de métricas.
func (s *ConsultarHistorialSprints) Ejecutar(ctx context.Context, proyectoID int64) ([]domain.SprintHistorial, error) {
	metricas, err := s.metricas.Ejecutar(ctx, proyectoID)
	if err != nil {
		return nil, err
	}
	sprints, err := s.sprints.ListarPorProyecto(ctx, proyectoID)
	if err != nil {
		return nil, err
	}

	finalizados := domain.SprintsFinalizados(sprints)
	historial := make([]domain.SprintHistorial, 0, len(finalizados))
	for _, sprint := range finalizados {
		m, err := metricasDelSprint(metricas, sprint.ID)
		if err != nil {
			return nil, err
		}
		historial = append(historial, domain.SprintHistorial{Sprint: sprint, Metricas: m})
	}
	return historial, nil
}

// metricasDelSprint busca las métricas de un Sprint finalizado. Si faltan, las
// fuentes son inconsistentes: es un error interno.
func metricasDelSprint(metricas domain.MetricasProyecto, sprintID int64) (domain.Metricas, error) {
	for _, m := range metricas.Sprints {
		if m.SprintID == sprintID {
			return m.Metricas, nil
		}
	}
	return domain.Metricas{}, fmt.Errorf("faltan las métricas del Sprint finalizado %d", sprintID)
}
