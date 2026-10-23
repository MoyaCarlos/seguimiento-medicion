package service

import (
	"context"
	"sort"

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

func (s *ListarProyectos) Ejecutar(ctx context.Context) ([]domain.Project, error) {
	proyectos, err := s.proyectos.Listar(ctx)
	if err != nil {
		return nil, err
	}
	if proyectos == nil {
		proyectos = []domain.Project{}
	}
	sort.Slice(proyectos, func(i, j int) bool { return proyectos[i].ID < proyectos[j].ID })
	return proyectos, nil
}
