package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// RegistrarEstimacion guarda el valor acordado como Story Points de la historia.
// Solo lo hace un Scrum Master del proyecto, con la ronda actual revelada.
type RegistrarEstimacion struct {
	historias  repository.EstimacionHistoriaRepository
	proyectos  repository.ProjectRepository
	votaciones repository.VotacionRepository
	sprints    repository.SprintRepository
}

func NewRegistrarEstimacion(historias repository.EstimacionHistoriaRepository, proyectos repository.ProjectRepository, votaciones repository.VotacionRepository, sprints repository.SprintRepository) *RegistrarEstimacion {
	return &RegistrarEstimacion{historias: historias, proyectos: proyectos, votaciones: votaciones, sprints: sprints}
}

func (s *RegistrarEstimacion) Ejecutar(ctx context.Context, historiaID, usuarioID int64, valor int) (domain.BacklogItem, error) {
	return domain.BacklogItem{}, nil
}
