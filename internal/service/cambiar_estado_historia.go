package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CambiarEstadoHistoria avanza el estado de una historia dentro de su Sprint.
type CambiarEstadoHistoria struct {
	historias repository.CambioEstadoHistoriaRepository
	sprints   repository.SprintRepository
	reloj     func() time.Time
}

func NewCambiarEstadoHistoria(historias repository.CambioEstadoHistoriaRepository, sprints repository.SprintRepository, reloj func() time.Time) *CambiarEstadoHistoria {
	return &CambiarEstadoHistoria{historias: historias, sprints: sprints, reloj: reloj}
}

// Ejecutar cambia el estado de la historia y la devuelve actualizada.
func (s *CambiarEstadoHistoria) Ejecutar(ctx context.Context, id int64, nuevo domain.Estado) (domain.BacklogItem, error) {
	historia, err := s.historias.ObtenerPorID(ctx, id)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	var sprint *domain.Sprint
	if historia.SprintID != nil {
		encontrado, err := s.sprints.ObtenerPorID(ctx, *historia.SprintID)
		if err != nil {
			return domain.BacklogItem{}, err
		}
		sprint = &encontrado
	}
	if err := domain.ValidarCambioDeEstado(historia, sprint, nuevo); err != nil {
		return domain.BacklogItem{}, err
	}
	var fecha *time.Time
	if nuevo == domain.EstadoCompletada {
		ahora := s.reloj()
		fecha = &ahora
	}
	if err := s.historias.CambiarEstado(ctx, id, historia.Estado, nuevo, fecha); err != nil {
		return domain.BacklogItem{}, err
	}
	historia.Estado = nuevo
	historia.FechaFinalizacion = fecha
	return historia, nil
}
