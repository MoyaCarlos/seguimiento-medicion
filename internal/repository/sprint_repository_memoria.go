package repository

import (
	"context"
	"slices"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SprintRepositoryEnMemoria implementa SprintRepository sin base de datos.
// Pensado para tests de cualquier historia que dependa de Sprint (ej. HU-13).
type SprintRepositoryEnMemoria struct {
	sprints     []domain.Sprint
	arrastradas map[int64][]domain.BacklogItem
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

// CerrarConArrastre valida todo antes de modificar nada, así un fallo deja el
// Sprint y las arrastradas como estaban, igual que la transacción de SQLite.
// No guarda historias: solo registra la foto de las arrastradas.
func (r *SprintRepositoryEnMemoria) CerrarConArrastre(_ context.Context, s domain.Sprint, noCompletadas []domain.BacklogItem) error {
	i := slices.IndexFunc(r.sprints, func(guardado domain.Sprint) bool { return guardado.ID == s.ID })
	if i < 0 {
		return domain.ErrSprintNoEncontrado
	}
	for _, h := range noCompletadas {
		if h.ProyectoID != s.ProyectoID {
			return ErrIntegridadReferencial
		}
		if h.SprintID == nil || *h.SprintID != s.ID {
			return ErrHistoriaFueraDelSprint
		}
	}

	r.sprints[i] = s
	if r.arrastradas == nil {
		r.arrastradas = map[int64][]domain.BacklogItem{}
	}
	for _, h := range noCompletadas {
		if h.EstimacionSP != nil {
			foto := *h.EstimacionSP
			h.EstimacionSP = &foto
		}
		r.arrastradas[s.ID] = append(r.arrastradas[s.ID], h)
	}
	return nil
}

// Arrastradas devuelve las historias registradas al cerrar el Sprint.
func (r *SprintRepositoryEnMemoria) Arrastradas(sprintID int64) []domain.BacklogItem {
	return r.arrastradas[sprintID]
}
