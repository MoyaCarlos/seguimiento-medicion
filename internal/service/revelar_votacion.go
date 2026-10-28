package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// RevelarVotacion muestra los votos de la ronda actual de una historia cuando
// votaron todos los integrantes. Solo lo hace un Scrum Master del proyecto.
type RevelarVotacion struct {
	historias  repository.HistoriaRepository
	proyectos  repository.ProjectRepository
	votaciones repository.VotacionRepository
}

func NewRevelarVotacion(historias repository.HistoriaRepository, proyectos repository.ProjectRepository, votaciones repository.VotacionRepository) *RevelarVotacion {
	return &RevelarVotacion{historias: historias, proyectos: proyectos, votaciones: votaciones}
}

func (s *RevelarVotacion) Ejecutar(ctx context.Context, historiaID, usuarioID int64, ahora time.Time) (domain.RondaVotacion, error) {
	return domain.RondaVotacion{}, nil
}
