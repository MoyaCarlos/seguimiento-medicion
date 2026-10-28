package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// IniciarVotacion abre la ronda 1 de votación de una historia. Solo lo hace un
// Scrum Master del proyecto, y no modifica la historia (FR-013).
type IniciarVotacion struct {
	historias  repository.HistoriaRepository
	proyectos  repository.ProjectRepository
	votaciones repository.VotacionRepository
}

func NewIniciarVotacion(historias repository.HistoriaRepository, proyectos repository.ProjectRepository, votaciones repository.VotacionRepository) *IniciarVotacion {
	return &IniciarVotacion{historias: historias, proyectos: proyectos, votaciones: votaciones}
}

func (s *IniciarVotacion) Ejecutar(ctx context.Context, historiaID, usuarioID int64) (domain.RondaVotacion, error) {
	return domain.RondaVotacion{}, nil
}
