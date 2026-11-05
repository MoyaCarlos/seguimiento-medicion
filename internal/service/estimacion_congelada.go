package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// validarEstimacionCongelada mantiene fijos los Story Points que usan las
// métricas de un Sprint (D2b): con un Sprint Finalizado no se cargan ni se
// cambian, y con uno Activo solo se permite la primera carga.
func validarEstimacionCongelada(ctx context.Context, sprints repository.SprintRepository, actual domain.BacklogItem, nueva *int) error {
	if actual.SprintID == nil || mismaEstimacion(actual.EstimacionSP, nueva) {
		return nil
	}

	sprint, err := sprints.ObtenerPorID(ctx, *actual.SprintID)
	if err != nil {
		return err
	}
	switch {
	case actual.EstimacionSP == nil && sprint.Estado == domain.SprintFinalizado:
		return domain.ErrCargaEnSprintFinalizado
	case actual.EstimacionSP != nil && (sprint.Estado == domain.SprintActivo || sprint.Estado == domain.SprintFinalizado):
		return domain.ErrEstimacionCongelada
	}
	return nil
}

func mismaEstimacion(actual, nueva *int) bool {
	if actual == nil || nueva == nil {
		return actual == nueva
	}
	return *actual == *nueva
}
