package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ResolverDefecto marca un defecto Pendiente como Resuelto con el Sprint
// activo de su proyecto como Sprint de resolución (HU-08).
type ResolverDefecto struct {
	defectos    repository.DefectoRepository
	sprints     repository.SprintRepository
	integrantes repository.IntegrantesRepository
}

func NewResolverDefecto(defectos repository.DefectoRepository, sprints repository.SprintRepository, integrantes repository.IntegrantesRepository) *ResolverDefecto {
	return &ResolverDefecto{defectos: defectos, sprints: sprints, integrantes: integrantes}
}

// Ejecutar valida en el orden de la spec: defecto existente, quien resuelve
// (usuarioID, la sesión) integra el proyecto (HU-21), defecto no resuelto y
// proyecto con Sprint activo. Devuelve el defecto ya resuelto.
func (s *ResolverDefecto) Ejecutar(ctx context.Context, defectoID, usuarioID int64) (domain.DefectoDetalle, error) {
	detalle, err := s.defectos.ObtenerPorID(ctx, defectoID)
	if err != nil {
		return domain.DefectoDetalle{}, err
	}
	if err := exigirIntegrante(ctx, s.integrantes, detalle.ProyectoID, usuarioID); err != nil {
		return domain.DefectoDetalle{}, err
	}
	sprints, err := s.sprints.ListarPorProyecto(ctx, detalle.ProyectoID)
	if err != nil {
		return domain.DefectoDetalle{}, err
	}
	var activo *domain.Sprint
	if sprint, hay := domain.SprintActivoDe(sprints); hay {
		activo = &sprint
	}
	if err := detalle.Defecto.Resolver(activo); err != nil {
		return domain.DefectoDetalle{}, err
	}
	if err := s.defectos.Resolver(ctx, detalle.Defecto); err != nil {
		return domain.DefectoDetalle{}, err
	}
	return detalle, nil
}
