package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// VotarInput son los datos de entrada del caso de uso.
type VotarInput struct {
	HistoriaID int64
	UsuarioID  int64
	Valor      int
}

// Votar registra en secreto el voto de un integrante en la ronda actual de una
// historia, o reemplaza el que ya tenía.
type Votar struct {
	historias  repository.HistoriaRepository
	proyectos  repository.ProjectRepository
	votaciones repository.VotacionRepository
}

func NewVotar(historias repository.HistoriaRepository, proyectos repository.ProjectRepository, votaciones repository.VotacionRepository) *Votar {
	return &Votar{historias: historias, proyectos: proyectos, votaciones: votaciones}
}

func (s *Votar) Ejecutar(ctx context.Context, input VotarInput) (domain.RondaVotacion, error) {
	historia, err := s.historias.ObtenerPorID(ctx, input.HistoriaID)
	if err != nil {
		return domain.RondaVotacion{}, err
	}
	_, integra, err := rolEnProyecto(ctx, s.proyectos, historia.ProyectoID, input.UsuarioID)
	if err != nil {
		return domain.RondaVotacion{}, err
	}
	if !integra {
		return domain.RondaVotacion{}, domain.ErrNoEsIntegrante
	}

	ronda, err := s.votaciones.ObtenerRondaActual(ctx, historia.ID)
	if err != nil {
		return domain.RondaVotacion{}, err
	}
	if err := ronda.Votar(input.UsuarioID, input.Valor); err != nil {
		return domain.RondaVotacion{}, err
	}
	if err := s.votaciones.GuardarVoto(ctx, ronda.ID, domain.Voto{UsuarioID: input.UsuarioID, Valor: input.Valor}); err != nil {
		return domain.RondaVotacion{}, err
	}
	return ronda, nil
}
