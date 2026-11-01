package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// DefectoRepository es el puerto de persistencia de los defectos (HU-08).
//   - Registrar inserta el defecto solo si su Sprint de detección sigue Activo;
//     si no, devuelve domain.ErrSinSprintActivoParaDetectar. Una FK violada
//     devuelve ErrIntegridadReferencial.
//   - ObtenerPorID devuelve el defecto con el título de su historia, o
//     domain.ErrDefectoNoEncontrado si no existe.
//   - Resolver persiste la resolución solo si el defecto sigue Pendiente y su
//     Sprint de resolución sigue Activo; si no, devuelve
//     domain.ErrDefectoYaResuelto o domain.ErrSinSprintActivoParaResolver.
//   - ListarPorProyecto ordena por orden de registro y devuelve un slice vacío
//     (no nil) si el proyecto no tiene defectos.
type DefectoRepository interface {
	Registrar(ctx context.Context, d domain.Defecto) (domain.Defecto, error)
	ObtenerPorID(ctx context.Context, id int64) (domain.DefectoDetalle, error)
	Resolver(ctx context.Context, d domain.Defecto) error
	ListarPorProyecto(ctx context.Context, proyectoID int64) ([]domain.DefectoDetalle, error)
}
