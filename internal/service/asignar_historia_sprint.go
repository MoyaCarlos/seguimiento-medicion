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
	return domain.BacklogItem{}, nil
}
