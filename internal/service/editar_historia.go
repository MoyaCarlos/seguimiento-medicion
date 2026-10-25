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

// Ejecutar carga la historia, valida los datos, aplica la regla D2b sobre los
// Story Points y persiste los cambios.
func (s *EditarHistoria) Ejecutar(ctx context.Context, input EditarHistoriaInput) (domain.BacklogItem, error) {
	actual, err := s.historias.ObtenerPorID(ctx, input.ID)
	if err != nil {
		return domain.BacklogItem{}, err
	}

	editada, err := actual.ConDatosEditados(input.Titulo, input.Descripcion, input.CriteriosAceptacion, input.EstimacionSP)
	if err != nil {
		return domain.BacklogItem{}, err
	}

	if err := s.validarEstimacionCongelada(ctx, actual, editada.EstimacionSP); err != nil {
		return domain.BacklogItem{}, err
	}

	if err := s.historias.Actualizar(ctx, editada); err != nil {
		return domain.BacklogItem{}, err
	}
	return editada, nil
}

// validarEstimacionCongelada impide cambiar los Story Points ya cargados de una
// historia de un Sprint Activo o Finalizado, para que las métricas históricas
// se puedan reproducir (D2b). Cargarlos por primera vez sí se permite.
func (s *EditarHistoria) validarEstimacionCongelada(ctx context.Context, actual domain.BacklogItem, nueva *int) error {
	if actual.SprintID == nil || actual.EstimacionSP == nil {
		return nil
	}
	if nueva != nil && *nueva == *actual.EstimacionSP {
		return nil
	}

	sprint, err := s.sprints.ObtenerPorID(ctx, *actual.SprintID)
	if err != nil {
		return err
	}
	if sprint.Estado == domain.SprintActivo || sprint.Estado == domain.SprintFinalizado {
		return domain.ErrEstimacionCongelada
	}
	return nil
}
