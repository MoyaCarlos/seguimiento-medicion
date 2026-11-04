package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// DatosMetricasRepository es el puerto de lectura del cálculo de métricas (HU-03).
// ListarDatosSprintsFinalizados devuelve una fila por Sprint finalizado del
// proyecto, ordenada por Sprint, y un slice vacío (no nil) si no hay ninguno.
type DatosMetricasRepository interface {
	ListarDatosSprintsFinalizados(ctx context.Context, proyectoID int64) ([]domain.DatosMetricasSprint, error)
}
