package repository

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// CreadorSesiones persiste una sesión nueva identificada por la huella de su token (HU-20).
type CreadorSesiones interface {
	Crear(ctx context.Context, huella string, s domain.Sesion) (domain.Sesion, error)
}

// BuscadorSesiones busca una sesión, con su usuario, por la huella del token (HU-20).
type BuscadorSesiones interface {
	// ObtenerPorHuella devuelve domain.ErrNoEncontrado si no existe; devuelve
	// también las sesiones vencidas o cerradas, la vigencia la decide el dominio.
	ObtenerPorHuella(ctx context.Context, huella string) (domain.SesionConUsuario, error)
}

// CerradorSesiones marca una sesión como cerrada sin borrarla (HU-20).
type CerradorSesiones interface {
	Cerrar(ctx context.Context, sesionID int64, ahora time.Time) error
}
