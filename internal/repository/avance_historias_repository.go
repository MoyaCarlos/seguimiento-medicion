package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// AvanceHistoriasRepository es el puerto de lectura del avance de las
// historias (HU-18). ListarPorProyecto devuelve todas las historias del
// proyecto, con o sin Sprint, ordenadas por id; el slice vacío no es nil.
type AvanceHistoriasRepository interface {
	ListarPorProyecto(ctx context.Context, proyectoID int64) ([]domain.BacklogItem, error)
}
