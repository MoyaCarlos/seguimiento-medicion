package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// AbrirNuevaRonda abre la ronda siguiente de una votación con la ronda actual
// revelada, haya o no diferencias. Solo lo hace un Scrum Master del proyecto.
type AbrirNuevaRonda struct {
	historias  repository.HistoriaRepository
	proyectos  repository.ProjectRepository
	votaciones repository.VotacionRepository
}

func NewAbrirNuevaRonda(historias repository.HistoriaRepository, proyectos repository.ProjectRepository, votaciones repository.VotacionRepository) *AbrirNuevaRonda {
	return &AbrirNuevaRonda{historias: historias, proyectos: proyectos, votaciones: votaciones}
}

func (s *AbrirNuevaRonda) Ejecutar(ctx context.Context, historiaID, usuarioID int64) (domain.RondaVotacion, error) {
	return domain.RondaVotacion{}, nil
}
