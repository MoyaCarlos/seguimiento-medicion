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

// Ejecutar devuelve el Sprint Backlog en el orden que entrega el repositorio.
func (s *ConsultarSprintBacklog) Ejecutar(_ context.Context, _ int64) ([]domain.BacklogItem, error) {
	return nil, nil
}
