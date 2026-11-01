package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ListarDefectos devuelve los defectos de un proyecto en orden de registro (HU-08).
type ListarDefectos struct {
	proyectos repository.ProjectRepository
	defectos  repository.DefectoRepository
}

func NewListarDefectos(proyectos repository.ProjectRepository, defectos repository.DefectoRepository) *ListarDefectos {
	return &ListarDefectos{proyectos: proyectos, defectos: defectos}
}

func (s *ListarDefectos) Ejecutar(ctx context.Context, proyectoID int64) ([]domain.DefectoDetalle, error) {
	return nil, nil
}
