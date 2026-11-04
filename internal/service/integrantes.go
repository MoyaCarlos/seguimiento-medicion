package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// rolEnProyecto devuelve el rol del usuario en el proyecto; false si no lo integra.
func rolEnProyecto(ctx context.Context, proyectos repository.ProjectRepository, proyectoID, usuarioID int64) (domain.Role, bool, error) {
	integrantes, err := proyectos.ListarIntegrantes(ctx, proyectoID)
	if err != nil {
		return "", false, err
	}
	for _, m := range integrantes {
		if m.UserID == usuarioID {
			return m.Role, true, nil
		}
	}
	return "", false, nil
}

// exigirScrumMaster devuelve ErrNoEsScrumMaster si el usuario no tiene ese rol
// en el proyecto, también si no lo integra.
func exigirScrumMaster(ctx context.Context, proyectos repository.ProjectRepository, proyectoID, usuarioID int64) error {
	rol, _, err := rolEnProyecto(ctx, proyectos, proyectoID, usuarioID)
	if err != nil {
		return err
	}
	if rol != domain.RolScrumMaster {
		return domain.ErrNoEsScrumMaster
	}
	return nil
}
