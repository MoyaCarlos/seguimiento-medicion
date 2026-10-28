package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// AsignarHistoriaASprint mueve una historia del Product Backlog a un Sprint.
type AsignarHistoriaASprint struct {
	historias repository.AsignacionHistoriaRepository
	sprints   repository.SprintRepository
}

func NewAsignarHistoriaASprint(historias repository.AsignacionHistoriaRepository, sprints repository.SprintRepository) *AsignarHistoriaASprint {
	return &AsignarHistoriaASprint{historias: historias, sprints: sprints}
}

// Ejecutar asigna la historia al Sprint y devuelve la historia con su SprintID.
func (s *AsignarHistoriaASprint) Ejecutar(ctx context.Context, historiaID, sprintID int64) (domain.BacklogItem, error) {
	historia, err := s.historias.ObtenerPorID(ctx, historiaID)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	sprint, err := s.sprints.ObtenerPorID(ctx, sprintID)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	if err := domain.ValidarAsignacion(historia, sprint); err != nil {
		return domain.BacklogItem{}, err
	}
	if err := s.historias.AsignarASprint(ctx, historiaID, sprintID); err != nil {
		return domain.BacklogItem{}, err
	}
	historia.SprintID = &sprintID
	return historia, nil
}
