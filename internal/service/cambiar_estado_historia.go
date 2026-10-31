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
	return domain.BacklogItem{}, nil
}
