package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// BacklogRepository es el puerto de persistencia del Product Backlog.
type BacklogRepository interface {
	Guardar(ctx context.Context, item domain.BacklogItem) (domain.BacklogItem, error)
}

// ProductBacklogRepository es el puerto que usa la consulta del Product Backlog.
// ListarSinSprint devuelve las historias del proyecto sin Sprint asignado,
// ordenadas por prioridad M, S, C, W y luego por id.
type ProductBacklogRepository interface {
	ListarSinSprint(ctx context.Context, proyectoID int64) ([]domain.BacklogItem, error)
}

// HistoriaRepository es el puerto que usa la edición de una historia.
// ObtenerPorID devuelve domain.ErrHistoriaNoEncontrada si la historia no existe.
type HistoriaRepository interface {
	ObtenerPorID(ctx context.Context, id int64) (domain.BacklogItem, error)
	Actualizar(ctx context.Context, item domain.BacklogItem) error
}
