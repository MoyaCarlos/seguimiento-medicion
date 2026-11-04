package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// EmisorSesiones crea sesiones nuevas y entrega el token en claro una sola vez (HU-20).
type EmisorSesiones struct {
	sesiones repository.CreadorSesiones
	ahora    func() time.Time
}

func NewEmisorSesiones(sesiones repository.CreadorSesiones, ahora func() time.Time) *EmisorSesiones {
	return &EmisorSesiones{sesiones: sesiones, ahora: ahora}
}

func (e *EmisorSesiones) Emitir(ctx context.Context, usuarioID int64) (string, domain.Sesion, error) {
	return "", domain.Sesion{}, nil
}

// HuellaDeToken devuelve la huella SHA-256 hexadecimal con la que se guarda un token.
func HuellaDeToken(token string) string {
	return ""
}
