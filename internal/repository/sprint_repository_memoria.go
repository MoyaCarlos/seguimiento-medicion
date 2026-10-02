package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SprintRepositoryEnMemoria implementa SprintRepository sin base de datos.
// Pensado para tests de cualquier historia que dependa de Sprint (ej. HU-13).
type SprintRepositoryEnMemoria struct {
	sprints []domain.Sprint
}

func NewSprintRepositoryEnMemoria() *SprintRepositoryEnMemoria {
	return &SprintRepositoryEnMemoria{}
}

func (r *SprintRepositoryEnMemoria) Guardar(_ context.Context, s domain.Sprint) (domain.Sprint, error) {
	s.ID = int64(len(r.sprints) + 1)
	r.sprints = append(r.sprints, s)
	return s, nil
}

func (r *SprintRepositoryEnMemoria) ObtenerPorID(_ context.Context, id int64) (domain.Sprint, error) {
	for _, s := range r.sprints {
		if s.ID == id {
			return s, nil
		}
	}
	return domain.Sprint{}, domain.ErrSprintNoEncontrado
}

func (r *SprintRepositoryEnMemoria) ListarPorProyecto(_ context.Context, proyectoID int64) ([]domain.Sprint, error) {
	var resultado []domain.Sprint
	for _, s := range r.sprints {
		if s.ProyectoID == proyectoID {
			resultado = append(resultado, s)
		}
	}
	return resultado, nil
}

func (r *SprintRepositoryEnMemoria) Actualizar(_ context.Context, s domain.Sprint) error {
	for i := range r.sprints {
		if r.sprints[i].ID == s.ID {
			r.sprints[i] = s
			return nil
		}
	}
	return domain.ErrSprintNoEncontrado
}
