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

// exigirIntegrante devuelve ErrNoEsIntegrante si el usuario no integra el proyecto.
func exigirIntegrante(ctx context.Context, proyectos repository.ProjectRepository, proyectoID, usuarioID int64) error {
	_, integra, err := rolEnProyecto(ctx, proyectos, proyectoID, usuarioID)
	if err != nil {
		return err
	}
	if !integra {
		return domain.ErrNoEsIntegrante
	}
	return nil
}

// exigirScrumMaster devuelve ErrNoEsIntegrante si el usuario no integra el
// proyecto y ErrNoEsScrumMaster si lo integra con otro rol.
func exigirScrumMaster(ctx context.Context, proyectos repository.ProjectRepository, proyectoID, usuarioID int64) error {
	rol, integra, err := rolEnProyecto(ctx, proyectos, proyectoID, usuarioID)
	if err != nil {
		return err
	}
	if !integra {
		return domain.ErrNoEsIntegrante
	}
	if rol != domain.RolScrumMaster {
		return domain.ErrNoEsScrumMaster
	}
	return nil
}
