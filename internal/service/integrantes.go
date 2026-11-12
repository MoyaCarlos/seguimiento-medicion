package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// rolEnProyecto devuelve el rol del usuario en el proyecto; false si no lo integra.
func rolEnProyecto(ctx context.Context, integrantes repository.IntegrantesRepository, proyectoID, usuarioID int64) (domain.Role, bool, error) {
	miembros, err := integrantes.ListarIntegrantes(ctx, proyectoID)
	if err != nil {
		return "", false, err
	}
	for _, m := range miembros {
		if m.UserID == usuarioID {
			return m.Role, true, nil
		}
	}
	return "", false, nil
}

// exigirIntegrante es el guardia de las escrituras de un integrante (HU-21):
// exige integrar el proyecto (403) y que el proyecto esté abierto (409, #300).
func exigirIntegrante(ctx context.Context, integrantes repository.IntegrantesRepository, proyectoID, usuarioID int64) error {
	if err := verificarIntegrante(ctx, integrantes, proyectoID, usuarioID); err != nil {
		return err
	}
	return exigirProyectoAbierto(ctx, integrantes, proyectoID)
}

// exigirScrumMaster es el guardia de las acciones de gestión (HU-21): exige ser
// el Scrum Master del proyecto (403) y que el proyecto esté abierto (409, #300).
func exigirScrumMaster(ctx context.Context, integrantes repository.IntegrantesRepository, proyectoID, usuarioID int64) error {
	if err := verificarScrumMaster(ctx, integrantes, proyectoID, usuarioID); err != nil {
		return err
	}
	return exigirProyectoAbierto(ctx, integrantes, proyectoID)
}

// exigirProyectoAbierto es el único lugar donde se impone que un proyecto
// cerrado sea de solo lectura (spec 044, FR-005).
func exigirProyectoAbierto(ctx context.Context, integrantes repository.IntegrantesRepository, proyectoID int64) error {
	cerrado, err := integrantes.ProyectoCerrado(ctx, proyectoID)
	if err != nil {
		return err
	}
	if cerrado {
		return domain.ErrProyectoCerrado
	}
	return nil
}

// verificarIntegrante devuelve ErrNoEsIntegrante si el usuario no integra el proyecto.
func verificarIntegrante(ctx context.Context, integrantes repository.IntegrantesRepository, proyectoID, usuarioID int64) error {
	_, integra, err := rolEnProyecto(ctx, integrantes, proyectoID, usuarioID)
	if err != nil {
		return err
	}
	if !integra {
		return domain.ErrNoEsIntegrante
	}
	return nil
}

// verificarScrumMaster devuelve ErrNoEsIntegrante si el usuario no integra el
// proyecto y ErrNoEsScrumMaster si lo integra con otro rol.
func verificarScrumMaster(ctx context.Context, integrantes repository.IntegrantesRepository, proyectoID, usuarioID int64) error {
	rol, integra, err := rolEnProyecto(ctx, integrantes, proyectoID, usuarioID)
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
