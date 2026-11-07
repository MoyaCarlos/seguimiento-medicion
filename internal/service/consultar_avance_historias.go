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

func (s *ConsultarAvanceHistorias) Ejecutar(ctx context.Context, proyectoID int64) (domain.AvanceHistorias, error) {
	return domain.AvanceHistorias{}, nil
}
