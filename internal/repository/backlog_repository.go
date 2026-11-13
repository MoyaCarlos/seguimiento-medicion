package repository

import (
	"context"
	"time"

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
// Actualizar aplica D2b en el mismo UPDATE: si el Sprint cambió a un estado que
// congela la estimación entre la validación y la escritura, devuelve
// domain.ErrEdicionDesactualizada (#249).
type HistoriaRepository interface {
	ObtenerPorID(ctx context.Context, id int64) (domain.BacklogItem, error)
	Actualizar(ctx context.Context, item domain.BacklogItem) error
}

// EstimacionHistoriaRepository es el puerto que usa el registro de la estimación
// acordada. Separado de HistoriaRepository (ISP): escribe solo los Story Points,
// así no pisa una edición simultánea. GuardarEstimacion aplica en el mismo
// UPDATE la regla D2b y la ronda revelada; si no afecta filas devuelve
// domain.ErrRegistroEstimacionDesactualizado (#249) o
// domain.ErrHistoriaNoEncontrada si la historia no existe.
type EstimacionHistoriaRepository interface {
	ObtenerPorID(ctx context.Context, id int64) (domain.BacklogItem, error)
	GuardarEstimacion(ctx context.Context, id int64, sp int, rondaID int64) error
}

// AsignacionHistoriaRepository es el puerto que usa la asignación de una
// historia a un Sprint. Separado de HistoriaRepository (ISP): la edición no
// asigna Sprints. AsignarASprint solo asigna si la historia no tiene Sprint y
// devuelve domain.ErrHistoriaYaAsignada si no afectó ninguna fila.
type AsignacionHistoriaRepository interface {
	ObtenerPorID(ctx context.Context, id int64) (domain.BacklogItem, error)
	AsignarASprint(ctx context.Context, historiaID, sprintID int64) error
}

// CambioEstadoHistoriaRepository es el puerto que usa el cambio de estado de
// una historia. Separado de HistoriaRepository (ISP): la edición no cambia
// estados. CambiarEstado solo aplica si la historia sigue en "desde" y
// devuelve domain.ErrConflictoEstadoHistoria si no afectó ninguna fila.
type CambioEstadoHistoriaRepository interface {
	ObtenerPorID(ctx context.Context, id int64) (domain.BacklogItem, error)
	CambiarEstado(ctx context.Context, id int64, desde, hacia domain.Estado, fechaFinalizacion *time.Time) error
}
