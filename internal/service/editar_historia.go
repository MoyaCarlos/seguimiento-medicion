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
	// EstimacionLeida es la estimación que tenía la historia cuando se abrió el
	// formulario. Si no coincide con la actual, la edición se rechaza (#248).
	EstimacionLeida *int
	// UsuarioID es quien edita (la sesión, HU-21); debe integrar el proyecto.
	UsuarioID int64
}

// EditarHistoria aplica la edición de una historia del Product Backlog.
type EditarHistoria struct {
	historias   repository.HistoriaRepository
	sprints     repository.SprintRepository
	integrantes repository.IntegrantesRepository
}

func NewEditarHistoria(historias repository.HistoriaRepository, sprints repository.SprintRepository, integrantes repository.IntegrantesRepository) *EditarHistoria {
	return &EditarHistoria{historias: historias, sprints: sprints, integrantes: integrantes}
}

// Ejecutar carga la historia, valida los datos, aplica la regla D2b sobre los
// Story Points y persiste los cambios.
func (s *EditarHistoria) Ejecutar(ctx context.Context, input EditarHistoriaInput) (domain.BacklogItem, error) {
	actual, err := s.historias.ObtenerPorID(ctx, input.ID)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	if err := exigirIntegrante(ctx, s.integrantes, actual.ProyectoID, input.UsuarioID); err != nil {
		return domain.BacklogItem{}, err
	}

	editada, err := actual.ConDatosEditados(input.Titulo, input.Descripcion, input.CriteriosAceptacion, input.EstimacionSP)
	if err != nil {
		return domain.BacklogItem{}, err
	}

	if !mismaEstimacion(actual.EstimacionSP, input.EstimacionLeida) {
		return domain.BacklogItem{}, domain.ErrEstimacionDesactualizada
	}

	if err := validarEstimacionCongelada(ctx, s.sprints, actual, editada.EstimacionSP); err != nil {
		return domain.BacklogItem{}, err
	}

	if err := s.historias.ActualizarSiEstimacionSigueSiendo(ctx, editada, input.EstimacionLeida); err != nil {
		return domain.BacklogItem{}, err
	}
	return editada, nil
}
