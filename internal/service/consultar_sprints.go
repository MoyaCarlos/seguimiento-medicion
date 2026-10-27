package service

import (
	"context"

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
	return nil, nil
}
