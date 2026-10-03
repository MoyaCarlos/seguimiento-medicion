package service

import (
	"errors"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// resolverUsuario busca al integrante por su nombre normalizado y lo crea si no
// existe. Ante una carrera de creación (UNIQUE), re-consulta y devuelve el
// existente.
func resolverUsuario(usuarios repository.UserRepository, candidato domain.User) (domain.User, error) {
	existente, err := usuarios.FindByNormalizedName(candidato.NormalizedName)
	if err == nil {
		return *existente, nil
	}
	if !errors.Is(err, domain.ErrNoEncontrado) {
		return domain.User{}, err
	}

	if err := usuarios.Create(&candidato); err != nil {
		if errors.Is(err, repository.ErrUsuarioDuplicado) {
			existente, errRelectura := usuarios.FindByNormalizedName(candidato.NormalizedName)
			if errRelectura != nil {
				return domain.User{}, errRelectura
			}
			return *existente, nil
		}
		return domain.User{}, err
	}
	return candidato, nil
}
