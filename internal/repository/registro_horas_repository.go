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

// IntegrantesRepository es el puerto del control de acceso de HU-21: quiénes
// integran el proyecto, con qué rol, y si el proyecto está cerrado (#300).
type IntegrantesRepository interface {
	ListarIntegrantes(ctx context.Context, proyectoID int64) ([]domain.Member, error)
	// ProyectoCerrado informa si el proyecto tiene fecha de cierre; false si no
	// existe (el 404 lo resuelve el servicio antes).
	ProyectoCerrado(ctx context.Context, proyectoID int64) (bool, error)
}
