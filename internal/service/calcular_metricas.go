package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CalcularMetricas calcula al consultar las métricas de un proyecto (HU-03).
// No guarda nada. HU-09, HU-10 y HU-12 toman los valores de este caso de uso.
type CalcularMetricas struct {
	proyectos repository.ProjectRepository
	datos     repository.DatosMetricasRepository
}

func NewCalcularMetricas(proyectos repository.ProjectRepository, datos repository.DatosMetricasRepository) *CalcularMetricas {
	return &CalcularMetricas{proyectos: proyectos, datos: datos}
}

func (s *CalcularMetricas) Ejecutar(ctx context.Context, proyectoID int64) (domain.MetricasProyecto, error) {
	proyecto, err := s.proyectos.ObtenerPorID(ctx, proyectoID)
	if err != nil {
		return domain.MetricasProyecto{}, err
	}
	datos, err := s.datos.ListarDatosSprintsFinalizados(ctx, proyectoID)
	if err != nil {
		return domain.MetricasProyecto{}, err
	}
	return domain.CalcularMetricas(datos, proyecto.HorasPorStoryPoint), nil
}
