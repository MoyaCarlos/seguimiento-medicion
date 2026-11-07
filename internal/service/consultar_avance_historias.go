package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ConsultarAvanceHistorias devuelve las historias de un proyecto agrupadas por
// estado (HU-18). Solo lee.
type ConsultarAvanceHistorias struct {
	proyectos buscadorProyectos
	historias repository.AvanceHistoriasRepository
	sprints   listadorSprints
}

func NewConsultarAvanceHistorias(proyectos buscadorProyectos, historias repository.AvanceHistoriasRepository, sprints listadorSprints) *ConsultarAvanceHistorias {
	return &ConsultarAvanceHistorias{proyectos: proyectos, historias: historias, sprints: sprints}
}

// Ejecutar responde domain.ErrProyectoNoEncontrado si el proyecto no existe, sin
// leer historias ni Sprints. Una historia que apunta a un Sprint ausente del
// listado es un error interno, no un 404.
func (s *ConsultarAvanceHistorias) Ejecutar(ctx context.Context, proyectoID int64) (domain.AvanceHistorias, error) {
	if _, err := s.proyectos.ObtenerPorID(ctx, proyectoID); err != nil {
		return domain.AvanceHistorias{}, err
	}
	historias, err := s.historias.ListarPorProyecto(ctx, proyectoID)
	if err != nil {
		return domain.AvanceHistorias{}, err
	}
	sprints, err := s.sprints.ListarPorProyecto(ctx, proyectoID)
	if err != nil {
		return domain.AvanceHistorias{}, err
	}
	return domain.AgruparAvance(historias, sprints)
}
