package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ConsultarSprintBacklog devuelve las historias asignadas a un Sprint.
type ConsultarSprintBacklog struct {
	sprints     repository.SprintRepository
	historias   repository.HistoriasDeSprintRepository
	integrantes repository.IntegrantesRepository
}

func NewConsultarSprintBacklog(sprints repository.SprintRepository, historias repository.HistoriasDeSprintRepository, integrantes repository.IntegrantesRepository) *ConsultarSprintBacklog {
	return &ConsultarSprintBacklog{sprints: sprints, historias: historias, integrantes: integrantes}
}

// Ejecutar verifica que el Sprint exista y que consultanteID integre su
// proyecto, y devuelve su Sprint Backlog en el orden que entrega el repositorio.
func (s *ConsultarSprintBacklog) Ejecutar(ctx context.Context, sprintID, consultanteID int64) ([]domain.BacklogItem, error) {
	sprint, err := s.sprints.ObtenerPorID(ctx, sprintID)
	if err != nil {
		return nil, err
	}
	if err := exigirLector(ctx, s.integrantes, sprint.ProyectoID, consultanteID); err != nil {
		return nil, err
	}
	historias, err := s.historias.ListarPorSprint(ctx, sprintID)
	if err != nil {
		return nil, err
	}
	if historias == nil {
		historias = make([]domain.BacklogItem, 0)
	}
	return historias, nil
}
