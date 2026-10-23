package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type listadorProyectos interface {
	Listar(ctx context.Context) ([]domain.Project, error)
}

// ListarProyectos devuelve todos los proyectos ordenados por id.
type ListarProyectos struct {
	proyectos listadorProyectos
}

func NewListarProyectos(proyectos listadorProyectos) *ListarProyectos {
	return &ListarProyectos{proyectos: proyectos}
}

func (s *ListarProyectos) Ejecutar(_ context.Context) ([]domain.Project, error) {
	return nil, nil
}
