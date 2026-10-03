package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// UserRepository abstrae la persistencia de los integrantes del equipo.
type UserRepository interface {
	ObtenerPorNombreNormalizado(ctx context.Context, normalized string) (*domain.User, error)
	Guardar(ctx context.Context, u domain.User) (domain.User, error)
}
