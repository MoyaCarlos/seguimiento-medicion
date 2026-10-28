package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// rolEnProyecto devuelve el rol del usuario en el proyecto; false si no lo integra.
func rolEnProyecto(ctx context.Context, proyectos repository.ProjectRepository, proyectoID, usuarioID int64) (domain.Role, bool, error) {
	return "", false, nil
}
