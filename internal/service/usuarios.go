package service

import (
	"context"
	"errors"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// resolverUsuario busca al integrante por su nombre normalizado y lo crea si no
// existe. Ante una carrera de creación (UNIQUE), re-consulta y devuelve el
// existente.
func resolverUsuario(ctx context.Context, usuarios repository.UserRepository, candidato domain.User) (domain.User, error) {
	existente, err := usuarios.ObtenerPorNombreNormalizado(ctx, candidato.NombreNormalizado)
	if err == nil {
		return *existente, nil
	}
	if !errors.Is(err, domain.ErrNoEncontrado) {
		return domain.User{}, err
	}

	guardado, err := usuarios.Guardar(ctx, candidato)
	if err != nil {
		if errors.Is(err, repository.ErrUsuarioDuplicado) {
			existente, errRelectura := usuarios.ObtenerPorNombreNormalizado(ctx, candidato.NombreNormalizado)
			if errRelectura != nil {
				return domain.User{}, errRelectura
			}
			return *existente, nil
		}
		return domain.User{}, err
	}
	return guardado, nil
}
