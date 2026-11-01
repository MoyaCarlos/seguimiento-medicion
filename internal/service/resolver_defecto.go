package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ResolverDefecto marca un defecto Pendiente como Resuelto con el Sprint
// activo de su proyecto como Sprint de resolución (HU-08).
type ResolverDefecto struct {
	defectos repository.DefectoRepository
	sprints  repository.SprintRepository
}

func NewResolverDefecto(defectos repository.DefectoRepository, sprints repository.SprintRepository) *ResolverDefecto {
	return &ResolverDefecto{defectos: defectos, sprints: sprints}
}

func (s *ResolverDefecto) Ejecutar(ctx context.Context, defectoID int64) (domain.DefectoDetalle, error) {
	return domain.DefectoDetalle{}, nil
}
