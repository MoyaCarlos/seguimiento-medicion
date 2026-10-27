package service

import (
	"context"
	"sort"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type buscadorProyectos interface {
	ObtenerPorID(ctx context.Context, id int64) (domain.Project, error)
}

type listadorSprints interface {
	ListarPorProyecto(ctx context.Context, proyectoID int64) ([]domain.Sprint, error)
}

// ConsultarSprints devuelve los Sprints de un proyecto ordenados por id.
type ConsultarSprints struct {
	proyectos buscadorProyectos
	sprints   listadorSprints
}

func NewConsultarSprints(proyectos buscadorProyectos, sprints listadorSprints) *ConsultarSprints {
	return &ConsultarSprints{proyectos: proyectos, sprints: sprints}
}

func (s *ConsultarSprints) Ejecutar(ctx context.Context, proyectoID int64) ([]domain.Sprint, error) {
	if _, err := s.proyectos.ObtenerPorID(ctx, proyectoID); err != nil {
		return nil, err
	}
	sprints, err := s.sprints.ListarPorProyecto(ctx, proyectoID)
	if err != nil {
		return nil, err
	}
	if sprints == nil {
		sprints = []domain.Sprint{}
	}
	sort.Slice(sprints, func(i, j int) bool { return sprints[i].ID < sprints[j].ID })
	return sprints, nil
}
