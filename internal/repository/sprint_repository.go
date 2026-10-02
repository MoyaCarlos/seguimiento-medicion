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
	Actualizar(ctx context.Context, s domain.Sprint) error
}

// HistoriasDeSprintRepository es el puerto que usa el cierre de un Sprint.
// Separado de BacklogRepository (ISP): crear historias no necesita esto.
type HistoriasDeSprintRepository interface {
	ListarPorSprint(ctx context.Context, sprintID int64) ([]domain.BacklogItem, error)
	QuitarDeSprint(ctx context.Context, historiaID int64) error
}
