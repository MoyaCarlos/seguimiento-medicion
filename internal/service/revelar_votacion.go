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
	historia, err := s.historias.ObtenerPorID(ctx, historiaID)
	if err != nil {
		return domain.RondaVotacion{}, err
	}
	rol, _, err := rolEnProyecto(ctx, s.proyectos, historia.ProyectoID, usuarioID)
	if err != nil {
		return domain.RondaVotacion{}, err
	}
	if rol != domain.RolScrumMaster {
		return domain.RondaVotacion{}, domain.ErrNoEsScrumMaster
	}

	ronda, err := s.votaciones.ObtenerRondaActual(ctx, historia.ID)
	if err != nil {
		return domain.RondaVotacion{}, err
	}
	integrantes, err := s.proyectos.ListarIntegrantes(ctx, historia.ProyectoID)
	if err != nil {
		return domain.RondaVotacion{}, err
	}
	ids := make([]int64, len(integrantes))
	for i, m := range integrantes {
		ids[i] = m.UserID
	}
	if err := ronda.Revelar(ids, ahora); err != nil {
		return domain.RondaVotacion{}, err
	}
	if err := s.votaciones.MarcarRevelada(ctx, ronda.ID, ahora); err != nil {
		return domain.RondaVotacion{}, err
	}
	return ronda, nil
}
