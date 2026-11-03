package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// validarEstimacionCongelada impide cambiar los Story Points ya cargados de una
// historia de un Sprint Activo o Finalizado, para que las métricas históricas
// se puedan reproducir (D2b). Cargarlos por primera vez sí se permite.
func validarEstimacionCongelada(ctx context.Context, sprints repository.SprintRepository, actual domain.BacklogItem, nueva *int) error {
	if actual.SprintID == nil || actual.EstimacionSP == nil {
		return nil
	}
	if nueva != nil && *nueva == *actual.EstimacionSP {
		return nil
	}

	sprint, err := sprints.ObtenerPorID(ctx, *actual.SprintID)
	if err != nil {
		return err
	}
	if sprint.Estado == domain.SprintActivo || sprint.Estado == domain.SprintFinalizado {
		return domain.ErrEstimacionCongelada
	}
	return nil
}
