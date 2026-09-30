package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// BacklogRepository es el puerto de persistencia del Product Backlog.
type BacklogRepository interface {
	Guardar(ctx context.Context, item domain.BacklogItem) (domain.BacklogItem, error)
}
