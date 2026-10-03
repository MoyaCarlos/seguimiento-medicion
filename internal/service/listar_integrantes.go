package service

import (
	"context"

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
func (s *ListarIntegrantes) Ejecutar(ctx context.Context, proyectoID int64) ([]domain.Member, error) {
	if _, err := s.proyectos.ObtenerPorID(ctx, proyectoID); err != nil {
		return nil, err
	}
	return s.proyectos.ListarIntegrantes(ctx, proyectoID)
}
