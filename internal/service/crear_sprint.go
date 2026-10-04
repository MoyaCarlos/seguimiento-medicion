package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CrearSprint crea un Sprint Pendiente para un proyecto existente.
type CrearSprint struct {
	proyectos repository.ProjectRepository
	sprints   repository.SprintRepository
}

func NewCrearSprint(proyectos repository.ProjectRepository, sprints repository.SprintRepository) *CrearSprint {
	return &CrearSprint{proyectos: proyectos, sprints: sprints}
}

func (s *CrearSprint) Ejecutar(ctx context.Context, proyectoID int64) (domain.Sprint, error) {
	sprint, err := domain.NewSprint(proyectoID)
	if err != nil {
		return domain.Sprint{}, err
	}
	if _, err := s.proyectos.ObtenerPorID(ctx, proyectoID); err != nil {
		return domain.Sprint{}, err
	}
	return s.sprints.Guardar(ctx, sprint)
}
