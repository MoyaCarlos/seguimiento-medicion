package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

type obtenedorSprint interface {
	ObtenerPorID(ctx context.Context, id int64) (domain.Sprint, error)
}

// ConsultarDetalleSprint devuelve el detalle de un Sprint finalizado: sus
// métricas de HU-03, sus historias comprometidas y sus defectos (HU-12). Solo
// lee.
type ConsultarDetalleSprint struct {
	proyectos buscadorProyectos
	sprints   obtenedorSprint
	metricas  calculadorMetricas
	historial repository.HistorialSprintRepository
}

func NewConsultarDetalleSprint(proyectos buscadorProyectos, sprints obtenedorSprint, metricas calculadorMetricas, historial repository.HistorialSprintRepository) *ConsultarDetalleSprint {
	return &ConsultarDetalleSprint{proyectos: proyectos, sprints: sprints, metricas: metricas, historial: historial}
}

func (s *ConsultarDetalleSprint) Ejecutar(ctx context.Context, proyectoID, sprintID int64) (domain.DetalleSprint, error) {
	return domain.DetalleSprint{}, nil
}
