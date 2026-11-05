package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ObtenerSesion resuelve quién es el dueño de un token de sesión (HU-20).
type ObtenerSesion struct {
	sesiones repository.BuscadorSesiones
	ahora    func() time.Time
}

func NewObtenerSesion(sesiones repository.BuscadorSesiones, ahora func() time.Time) *ObtenerSesion {
	return &ObtenerSesion{sesiones: sesiones, ahora: ahora}
}

// Ejecutar devuelve la sesión vigente del token o domain.ErrNoAutenticado.
func (o *ObtenerSesion) Ejecutar(ctx context.Context, token string) (domain.SesionConUsuario, error) {
	if !tokenBienFormado(token) {
		return domain.SesionConUsuario{}, domain.ErrNoAutenticado
	}

	encontrada, err := o.sesiones.ObtenerPorHuella(ctx, HuellaDeToken(token))
	if errors.Is(err, domain.ErrNoEncontrado) {
		return domain.SesionConUsuario{}, domain.ErrNoAutenticado
	}
	if err != nil {
		return domain.SesionConUsuario{}, fmt.Errorf("obtener sesión: %w", err)
	}
	if !encontrada.Vigente(o.ahora()) {
		return domain.SesionConUsuario{}, domain.ErrNoAutenticado
	}
	return encontrada, nil
}

// tokenBienFormado indica si el token tiene la forma que emite EmisorSesiones:
// 32 bytes en base64url sin relleno, es decir 43 caracteres del alfabeto.
func tokenBienFormado(token string) bool {
	if len(token) != base64.RawURLEncoding.EncodedLen(tamanoTokenSesion) {
		return false
	}
	for i := 0; i < len(token); i++ {
		c := token[i]
		letra := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
		if !letra && !(c >= '0' && c <= '9') && c != '-' && c != '_' {
			return false
		}
	}
	return true
}
