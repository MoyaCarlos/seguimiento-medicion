package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// ProjectRepository es el puerto de persistencia de proyectos y su equipo.
// El contrato de Guardar/ObtenerPorID es consumido por otras historias y no se
// altera.
type ProjectRepository interface {
	Guardar(ctx context.Context, p domain.Project) (domain.Project, error)
	GuardarConScrumMaster(ctx context.Context, p domain.Project, creadorID int64) (domain.Project, error)
	ObtenerPorID(ctx context.Context, id int64) (domain.Project, error)
	Actualizar(ctx context.Context, p domain.Project) error
	AgregarIntegrante(ctx context.Context, m domain.Membership) (domain.Membership, error)
	ListarIntegrantes(ctx context.Context, proyectoID int64) ([]domain.Member, error)
}
