package repository

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// RegistroHorasRepository es el puerto de persistencia de las cargas de horas.
// Registrar devuelve domain.ErrHistoriaFueraDelSprintActivo si, al insertar, la
// historia ya no está en ese Sprint o el Sprint ya no está Activo.
// ListarPorHistoria ordena por fecha y luego por id, y devuelve un slice vacío
// (no nil) si la historia no tiene cargas.
type RegistroHorasRepository interface {
	Registrar(ctx context.Context, r domain.RegistroHoras) (domain.RegistroHoras, error)
	ListarPorHistoria(ctx context.Context, historiaID int64) ([]domain.RegistroHorasDetalle, error)
}

// IntegrantesRepository es el puerto que usa la imputación de horas para
// verificar que quien carga es integrante del proyecto.
type IntegrantesRepository interface {
	ListarIntegrantes(ctx context.Context, proyectoID int64) ([]domain.Member, error)
}
