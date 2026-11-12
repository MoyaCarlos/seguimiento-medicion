package repository

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// ProjectRepository es el puerto de persistencia de proyectos y su equipo.
// Los métodos se renombraron a la convención en español (Create→Guardar,
// GetByID→ObtenerPorID, Update→Actualizar, AddMember→AgregarIntegrante,
// ListMembers→ListarIntegrantes); la interfaz es consumida por otras historias.
type ProjectRepository interface {
	Guardar(ctx context.Context, p domain.Project) (domain.Project, error)
	GuardarConScrumMaster(ctx context.Context, p domain.Project, creadorID int64) (domain.Project, error)
	ObtenerPorID(ctx context.Context, id int64) (domain.Project, error)
	// ListarPorIntegrante devuelve los proyectos que integra usuarioID, con su
	// rol en cada uno, ordenados por id; arreglo vacío si no integra ninguno.
	ListarPorIntegrante(ctx context.Context, usuarioID int64) ([]domain.ProyectoConRol, error)
	Actualizar(ctx context.Context, p domain.Project) error
	AgregarIntegrante(ctx context.Context, m domain.Membership) (domain.Membership, error)
	ListarIntegrantes(ctx context.Context, proyectoID int64) ([]domain.Member, error)
	// Cerrar guarda la fecha de cierre solo si el proyecto está abierto y sin
	// Sprint Activo (#300). Si no lo cierra, devuelve ErrProyectoNoEncontrado,
	// ErrProyectoYaCerrado o ErrSprintActivoAlCerrarProyecto.
	Cerrar(ctx context.Context, id int64, fecha time.Time) error
}
