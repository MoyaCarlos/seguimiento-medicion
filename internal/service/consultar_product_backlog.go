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

// Ejecutar verifica que el proyecto exista y devuelve su Product Backlog en el
// orden que entrega el repositorio (prioridad y luego creación).
func (s *ConsultarProductBacklog) Ejecutar(ctx context.Context, proyectoID int64) ([]domain.BacklogItem, error) {
	if _, err := s.proyectos.ObtenerPorID(ctx, proyectoID); err != nil {
		return nil, err
	}
	historias, err := s.backlog.ListarSinSprint(ctx, proyectoID)
	if err != nil {
		return nil, err
	}
	if historias == nil {
		historias = make([]domain.BacklogItem, 0)
	}
	return historias, nil
}
