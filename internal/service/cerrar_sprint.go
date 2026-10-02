package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CerrarSprint finaliza un Sprint Activo y devuelve al Product Backlog las
// historias que no se completaron.
type CerrarSprint struct {
	sprints   repository.SprintRepository
	historias repository.HistoriasDeSprintRepository
}

func NewCerrarSprint(sprints repository.SprintRepository, historias repository.HistoriasDeSprintRepository) *CerrarSprint {
	return &CerrarSprint{sprints: sprints, historias: historias}
}

func (s *CerrarSprint) Ejecutar(ctx context.Context, sprintID int64) (domain.Sprint, error) {
	sprint, err := s.sprints.ObtenerPorID(ctx, sprintID)
	if err != nil {
		return domain.Sprint{}, err
	}
	if err := sprint.Cerrar(); err != nil {
		return domain.Sprint{}, err
	}

	historias, err := s.historias.ListarPorSprint(ctx, sprint.ID)
	if err != nil {
		return domain.Sprint{}, err
	}
	// ponytail: sin transacción. Primero se liberan las historias y al final
	// se marca el Sprint: si algo falla a mitad, el Sprint sigue Activo y
	// reintentar el cierre completa lo que faltó. Transacción si hiciera
	// falta atomicidad estricta.
	for _, h := range historias {
		if h.Estado == domain.EstadoCompletada {
			continue
		}
		if err := s.historias.QuitarDeSprint(ctx, h.ID); err != nil {
			return domain.Sprint{}, err
		}
	}

	if err := s.sprints.Actualizar(ctx, sprint); err != nil {
		return domain.Sprint{}, err
	}
	return sprint, nil
}
