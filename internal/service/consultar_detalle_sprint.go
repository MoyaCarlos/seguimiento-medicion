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

// Ejecutar valida en orden: proyecto existente, Sprint existente y del
// proyecto, Sprint finalizado. Recién entonces lee el detalle.
func (s *ConsultarDetalleSprint) Ejecutar(ctx context.Context, proyectoID, sprintID int64) (domain.DetalleSprint, error) {
	if _, err := s.proyectos.ObtenerPorID(ctx, proyectoID); err != nil {
		return domain.DetalleSprint{}, err
	}
	sprint, err := s.sprints.ObtenerPorID(ctx, sprintID)
	if err != nil {
		return domain.DetalleSprint{}, err
	}
	if sprint.ProyectoID != proyectoID {
		return domain.DetalleSprint{}, domain.ErrSprintNoEncontrado
	}
	if err := domain.ValidarSprintFinalizado(sprint); err != nil {
		return domain.DetalleSprint{}, err
	}

	metricasProyecto, err := s.metricas.Ejecutar(ctx, proyectoID)
	if err != nil {
		return domain.DetalleSprint{}, err
	}
	metricas, err := metricasDelSprint(metricasProyecto, sprintID)
	if err != nil {
		return domain.DetalleSprint{}, err
	}
	historias, err := s.historial.HistoriasComprometidas(ctx, sprintID)
	if err != nil {
		return domain.DetalleSprint{}, err
	}
	minutos, err := s.historial.MinutosPorHistoria(ctx, sprintID)
	if err != nil {
		return domain.DetalleSprint{}, err
	}
	defectos, err := s.historial.DefectosDeSprint(ctx, sprintID)
	if err != nil {
		return domain.DetalleSprint{}, err
	}
	base := domain.SprintHistorial{Sprint: sprint, Metricas: metricas}
	return domain.ArmarDetalleSprint(base, historias, minutos, defectos), nil
}
