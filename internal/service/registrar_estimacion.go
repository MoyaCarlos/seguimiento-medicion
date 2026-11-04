package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// RegistrarEstimacion guarda el valor acordado como Story Points de la historia.
// Solo lo hace un Scrum Master del proyecto, con la ronda actual revelada.
type RegistrarEstimacion struct {
	historias  repository.EstimacionHistoriaRepository
	proyectos  repository.ProjectRepository
	votaciones repository.VotacionRepository
	sprints    repository.SprintRepository
}

func NewRegistrarEstimacion(historias repository.EstimacionHistoriaRepository, proyectos repository.ProjectRepository, votaciones repository.VotacionRepository, sprints repository.SprintRepository) *RegistrarEstimacion {
	return &RegistrarEstimacion{historias: historias, proyectos: proyectos, votaciones: votaciones, sprints: sprints}
}

// Ejecutar recibe la historia ya leída por quien la llama, que tiene que
// comprobar que existe antes de validar el pedido (FR-011). No escribe si el
// valor ya es el cargado.
func (s *RegistrarEstimacion) Ejecutar(ctx context.Context, historia domain.BacklogItem, usuarioID int64, valor int) (domain.BacklogItem, error) {
	if err := exigirScrumMaster(ctx, s.proyectos, historia.ProyectoID, usuarioID); err != nil {
		return domain.BacklogItem{}, err
	}

	conEstimacion, err := historia.ConEstimacion(valor)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	revelada, err := s.votaciones.RondaActualRevelada(ctx, historia.ID)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	if !revelada {
		return domain.BacklogItem{}, domain.ErrRondaNoRevelada
	}
	if err := validarEstimacionCongelada(ctx, s.sprints, historia, conEstimacion.EstimacionSP); err != nil {
		return domain.BacklogItem{}, err
	}

	if historia.EstimacionSP != nil && *historia.EstimacionSP == valor {
		return historia, nil
	}
	if err := s.historias.GuardarEstimacion(ctx, historia.ID, valor); err != nil {
		return domain.BacklogItem{}, err
	}
	return conEstimacion, nil
}
