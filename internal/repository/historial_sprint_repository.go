package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// HistorialSprintRepository es el puerto de lectura del detalle de un Sprint
// en el historial (HU-12). Solo lee; los slices y el mapa vacíos no son nil.
type HistorialSprintRepository interface {
	// HistoriasComprometidas devuelve las historias completadas en el Sprint y
	// las arrastradas al cerrarlo, con la foto de sus Story Points: primero las
	// completadas y luego las arrastradas, cada grupo por id.
	HistoriasComprometidas(ctx context.Context, sprintID int64) ([]domain.HistoriaComprometida, error)
	// MinutosPorHistoria suma los minutos imputados en el Sprint, por historia.
	MinutosPorHistoria(ctx context.Context, sprintID int64) (map[int64]int, error)
	// DefectosDeSprint devuelve los defectos detectados o resueltos en el
	// Sprint, una sola vez cada uno, ordenados por id.
	DefectosDeSprint(ctx context.Context, sprintID int64) ([]domain.DefectoDetalle, error)
}
