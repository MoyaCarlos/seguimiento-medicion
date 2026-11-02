package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ConsultarSprintBacklog devuelve las historias asignadas a un Sprint.
type ConsultarSprintBacklog struct {
	sprints   repository.SprintRepository
	historias repository.HistoriasDeSprintRepository
}

func NewConsultarSprintBacklog(sprints repository.SprintRepository, historias repository.HistoriasDeSprintRepository) *ConsultarSprintBacklog {
	return &ConsultarSprintBacklog{sprints: sprints, historias: historias}
}

// Ejecutar verifica que el Sprint exista y devuelve su Sprint Backlog en el
// orden que entrega el repositorio.
func (s *ConsultarSprintBacklog) Ejecutar(ctx context.Context, sprintID int64) ([]domain.BacklogItem, error) {
	if _, err := s.sprints.ObtenerPorID(ctx, sprintID); err != nil {
		return nil, err
	}
	historias, err := s.historias.ListarPorSprint(ctx, sprintID)
	if err != nil {
		return nil, err
	}
	if historias == nil {
		historias = make([]domain.BacklogItem, 0)
	}
	return historias, nil
}
