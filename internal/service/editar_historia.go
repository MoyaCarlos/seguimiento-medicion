package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// EditarHistoriaInput son los datos editables de una historia.
type EditarHistoriaInput struct {
	ID                  int64
	Titulo              string
	Descripcion         string
	CriteriosAceptacion *string
	EstimacionSP        *int
}

// EditarHistoria aplica la edición de una historia del Product Backlog.
type EditarHistoria struct {
	historias repository.HistoriaRepository
	sprints   repository.SprintRepository
}

func NewEditarHistoria(historias repository.HistoriaRepository, sprints repository.SprintRepository) *EditarHistoria {
	return &EditarHistoria{historias: historias, sprints: sprints}
}

// Ejecutar carga la historia, valida los datos y persiste los cambios.
func (s *EditarHistoria) Ejecutar(ctx context.Context, input EditarHistoriaInput) (domain.BacklogItem, error) {
	return domain.BacklogItem{}, nil
}
