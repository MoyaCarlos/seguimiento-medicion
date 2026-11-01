package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SprintRepository es el puerto de persistencia de los Sprints.
// ObtenerPorID devuelve domain.ErrSprintNoEncontrado si el Sprint no existe.
type SprintRepository interface {
	Guardar(ctx context.Context, s domain.Sprint) (domain.Sprint, error)
	ObtenerPorID(ctx context.Context, id int64) (domain.Sprint, error)
	ListarPorProyecto(ctx context.Context, proyectoID int64) ([]domain.Sprint, error)
	// Actualizar persiste el inicio de un Sprint. Devuelve
	// domain.ErrSprintNoPendiente si el Sprint guardado ya no está Pendiente.
	Actualizar(ctx context.Context, s domain.Sprint) error
	// CerrarConArrastre persiste el Sprint ya cerrado y, en una sola operación
	// atómica, registra las historias no completadas con la foto de sus Story
	// Points y las devuelve al Product Backlog en estado Nueva. Devuelve
	// domain.ErrSprintNoActivo si el Sprint guardado ya no está Activo.
	CerrarConArrastre(ctx context.Context, s domain.Sprint, noCompletadas []domain.BacklogItem) error
}

// HistoriasDeSprintRepository es el puerto que usa el cierre de un Sprint.
// Separado de BacklogRepository (ISP): crear historias no necesita esto.
type HistoriasDeSprintRepository interface {
	// ListarPorSprint devuelve las historias del Sprint ordenadas por prioridad
	// MoSCoW (M, S, C, W) y luego por id, el mismo orden que el Product Backlog.
	// El cierre de Sprint no depende de este orden.
	ListarPorSprint(ctx context.Context, sprintID int64) ([]domain.BacklogItem, error)
}
