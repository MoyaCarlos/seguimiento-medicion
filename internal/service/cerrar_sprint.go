package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CerrarSprint finaliza un Sprint Activo y, en la misma operación, registra
// como arrastradas las historias que no se completaron y las devuelve al
// Product Backlog.
type CerrarSprint struct {
	sprints     repository.SprintRepository
	historias   repository.HistoriasDeSprintRepository
	integrantes repository.IntegrantesRepository
}

func NewCerrarSprint(sprints repository.SprintRepository, historias repository.HistoriasDeSprintRepository, integrantes repository.IntegrantesRepository) *CerrarSprint {
	return &CerrarSprint{sprints: sprints, historias: historias, integrantes: integrantes}
}

// Ejecutar cierra el Sprint; quien lo cierra (usuarioID, la sesión) debe ser el
// Scrum Master del proyecto (HU-21).
func (s *CerrarSprint) Ejecutar(ctx context.Context, sprintID, usuarioID int64) (domain.Sprint, error) {
	sprint, err := s.sprints.ObtenerPorID(ctx, sprintID)
	if err != nil {
		return domain.Sprint{}, err
	}
	if err := exigirScrumMaster(ctx, s.integrantes, sprint.ProyectoID, usuarioID); err != nil {
		return domain.Sprint{}, err
	}
	if err := sprint.Cerrar(); err != nil {
		return domain.Sprint{}, err
	}

	historias, err := s.historias.ListarPorSprint(ctx, sprint.ID)
	if err != nil {
		return domain.Sprint{}, err
	}
	var noCompletadas []domain.BacklogItem
	for _, h := range historias {
		if h.Estado != domain.EstadoCompletada {
			noCompletadas = append(noCompletadas, h)
		}
	}

	if err := s.sprints.CerrarConArrastre(ctx, sprint, noCompletadas); err != nil {
		return domain.Sprint{}, err
	}
	return sprint, nil
}
