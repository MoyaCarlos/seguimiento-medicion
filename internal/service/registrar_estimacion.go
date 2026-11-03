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

// Ejecutar no escribe si el valor ya es el cargado.
func (s *RegistrarEstimacion) Ejecutar(ctx context.Context, historiaID, usuarioID int64, valor int) (domain.BacklogItem, error) {
	historia, err := s.historias.ObtenerPorID(ctx, historiaID)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	rol, _, err := rolEnProyecto(ctx, s.proyectos, historia.ProyectoID, usuarioID)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	if rol != domain.RolScrumMaster {
		return domain.BacklogItem{}, domain.ErrNoEsScrumMaster
	}

	conEstimacion, err := historia.ConEstimacion(valor)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	ronda, err := s.votaciones.ObtenerRondaActual(ctx, historia.ID)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	if !ronda.Revelada() {
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
