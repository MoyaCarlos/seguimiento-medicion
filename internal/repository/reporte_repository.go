package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// HistoriasArrastradasRepository es el puerto de lectura de las historias que
// se arrastraron al cerrar un Sprint (HU-10). ListarArrastradasPorSprint
// devuelve las historias con la foto de sus Story Points y un slice vacío (no
// nil) si el Sprint no tiene arrastradas.
type HistoriasArrastradasRepository interface {
	ListarArrastradasPorSprint(ctx context.Context, sprintID int64) ([]domain.BacklogItem, error)
}
