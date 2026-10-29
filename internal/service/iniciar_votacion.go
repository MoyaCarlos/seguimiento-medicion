package service

import (
	"context"
	"errors"

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

// Ejecutar chequea antes que la historia no tenga ronda; si dos inicios
// simultáneos pasan el chequeo, el UNIQUE de la base frena al segundo.
func (s *IniciarVotacion) Ejecutar(ctx context.Context, historiaID, usuarioID int64) (domain.RondaVotacion, error) {
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
	if historia.Estado == domain.EstadoCompletada {
		return domain.RondaVotacion{}, domain.ErrHistoriaCompletada
	}

	_, err = s.votaciones.ObtenerRondaActual(ctx, historia.ID)
	if err == nil {
		return domain.RondaVotacion{}, domain.ErrVotacionYaIniciada
	}
	if !errors.Is(err, domain.ErrVotacionNoIniciada) {
		return domain.RondaVotacion{}, err
	}

	ronda, err := domain.NuevaRonda(historia.ID)
	if err != nil {
		return domain.RondaVotacion{}, err
	}
	return s.votaciones.CrearRonda(ctx, ronda)
}
