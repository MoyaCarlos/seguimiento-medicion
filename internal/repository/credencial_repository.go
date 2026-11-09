package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// RegistroCredenciales persiste un usuario nuevo junto con el hash de su contraseña (HU-20).
type RegistroCredenciales interface {
	// Registrar devuelve ErrUsuarioDuplicado si el nombre normalizado ya existe.
	Registrar(ctx context.Context, u domain.User, hash string) (domain.User, error)
}

// BuscadorCredenciales busca un usuario por nombre normalizado junto con su hash (HU-20).
type BuscadorCredenciales interface {
	// ObtenerPorNombreNormalizado devuelve el hash vacío si el usuario no tiene
	// contraseña (NULL) y domain.ErrNoEncontrado si no existe.
	ObtenerPorNombreNormalizado(ctx context.Context, normalizado string) (domain.User, string, error)
}
