package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

const tamanoTokenSesion = 32

// EmisorSesiones crea sesiones nuevas y entrega el token en claro una sola vez (HU-20).
type EmisorSesiones struct {
	sesiones repository.CreadorSesiones
	ahora    func() time.Time
}

func NewEmisorSesiones(sesiones repository.CreadorSesiones, ahora func() time.Time) *EmisorSesiones {
	return &EmisorSesiones{sesiones: sesiones, ahora: ahora}
}

// Emitir crea una sesión para el usuario y devuelve el token en claro; solo se
// persiste su huella.
func (e *EmisorSesiones) Emitir(ctx context.Context, usuarioID int64) (string, domain.Sesion, error) {
	bytes := make([]byte, tamanoTokenSesion)
	if _, err := rand.Read(bytes); err != nil {
		return "", domain.Sesion{}, fmt.Errorf("generar token de sesión: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(bytes)

	ahora := e.ahora()
	sesion, err := e.sesiones.Crear(ctx, HuellaDeToken(token), domain.Sesion{
		UsuarioID: usuarioID,
		CreadoEn:  ahora,
		VenceEn:   ahora.Add(domain.DuracionSesion),
	})
	if err != nil {
		return "", domain.Sesion{}, fmt.Errorf("emitir sesión: %w", err)
	}
	return token, sesion, nil
}

// HuellaDeToken devuelve la huella SHA-256 hexadecimal con la que se guarda un token.
func HuellaDeToken(token string) string {
	suma := sha256.Sum256([]byte(token))
	return hex.EncodeToString(suma[:])
}
