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

// AsignacionHistoriaRepository es el puerto que usa la asignación de una
// historia a un Sprint. Separado de HistoriaRepository (ISP): la edición no
// asigna Sprints. AsignarASprint solo asigna si la historia no tiene Sprint y
// devuelve domain.ErrHistoriaYaAsignada si no afectó ninguna fila.
type AsignacionHistoriaRepository interface {
	ObtenerPorID(ctx context.Context, id int64) (domain.BacklogItem, error)
	AsignarASprint(ctx context.Context, historiaID, sprintID int64) error
}
