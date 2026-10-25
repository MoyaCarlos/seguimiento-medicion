package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ConsultarProductBacklog devuelve las historias sin Sprint de un proyecto.
type ConsultarProductBacklog struct {
	proyectos repository.ProjectRepository
	backlog   repository.ProductBacklogRepository
}

func NewConsultarProductBacklog(proyectos repository.ProjectRepository, backlog repository.ProductBacklogRepository) *ConsultarProductBacklog {
	return &ConsultarProductBacklog{proyectos: proyectos, backlog: backlog}
}

// Ejecutar devuelve el Product Backlog del proyecto.
func (s *ConsultarProductBacklog) Ejecutar(ctx context.Context, proyectoID int64) ([]domain.BacklogItem, error) {
	return nil, nil
}
