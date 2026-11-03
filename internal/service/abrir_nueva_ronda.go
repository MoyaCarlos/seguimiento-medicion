package service

import (
	"context"
	"errors"

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

// Ejecutar deja que el UNIQUE de la base frene a la segunda de dos aperturas
// simultáneas de la misma ronda.
func (s *AbrirNuevaRonda) Ejecutar(ctx context.Context, historiaID, usuarioID int64) (domain.RondaVotacion, error) {
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

	actual, err := s.votaciones.ObtenerRondaActual(ctx, historia.ID)
	if err != nil {
		return domain.RondaVotacion{}, err
	}
	siguiente, err := actual.Siguiente()
	if err != nil {
		return domain.RondaVotacion{}, err
	}
	ronda, err := s.votaciones.CrearRonda(ctx, siguiente)
	if errors.Is(err, domain.ErrVotacionYaIniciada) {
		return domain.RondaVotacion{}, domain.ErrRondaYaAbierta
	}
	return ronda, err
}
