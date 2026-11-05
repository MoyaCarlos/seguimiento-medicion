package service

import (
	"context"

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

func (s *ConsultarHistorialSprints) Ejecutar(ctx context.Context, proyectoID int64) ([]domain.SprintHistorial, error) {
	return nil, nil
}
