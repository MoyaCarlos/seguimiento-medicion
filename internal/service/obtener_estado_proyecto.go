package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ObtenerEstadoProyecto calcula, en cada consulta, el estado general de un
// proyecto a partir de sus Sprints y sus fechas. Es una operación de solo
// lectura: nunca escribe.
type ObtenerEstadoProyecto struct {
	proyectos repository.ProjectRepository
	sprints   repository.SprintRepository
}

func NewObtenerEstadoProyecto(proyectos repository.ProjectRepository, sprints repository.SprintRepository) *ObtenerEstadoProyecto {
	return &ObtenerEstadoProyecto{proyectos: proyectos, sprints: sprints}
}

// Ejecutar devuelve el estado derivado, domain.ErrProyectoNoEncontrado si el
// proyecto no existe y domain.ErrNoEsIntegrante si consultanteID no lo integra.
func (s *ObtenerEstadoProyecto) Ejecutar(ctx context.Context, proyectoID, consultanteID int64, ahora time.Time) (domain.EstadoProyecto, error) {
	proyecto, err := s.proyectos.ObtenerPorID(ctx, proyectoID)
	if err != nil {
		return "", err
	}
	if err := exigirLector(ctx, s.proyectos, proyectoID, consultanteID); err != nil {
		return "", err
	}
	sprints, err := s.sprints.ListarPorProyecto(ctx, proyectoID)
	if err != nil {
		return "", err
	}
	return domain.CalcularEstadoProyecto(proyecto, sprints, ahora), nil
}
