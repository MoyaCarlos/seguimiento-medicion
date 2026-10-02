package service

import (
	"context"
	"errors"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ListarIntegrantes devuelve el equipo de un proyecto.
type ListarIntegrantes struct {
	proyectos repository.ProjectRepository
}

func NewListarIntegrantes(proyectos repository.ProjectRepository) *ListarIntegrantes {
	return &ListarIntegrantes{proyectos: proyectos}
}

// Ejecutar verifica que el proyecto exista y devuelve sus integrantes.
func (s *ListarIntegrantes) Ejecutar(_ context.Context, proyectoID string) ([]domain.Member, error) {
	_, err := s.proyectos.GetByID(proyectoID)
	if errors.Is(err, repository.ErrNoEncontrado) {
		return nil, domain.ErrNoEncontrado
	}
	if err != nil {
		return nil, err
	}
	return s.proyectos.ListMembers(proyectoID)
}
