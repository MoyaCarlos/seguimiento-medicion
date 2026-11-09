package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CerrarSesion cierra la sesión dueña de un token (HU-20).
type CerrarSesion struct {
	obtener  *ObtenerSesion
	cerrador repository.CerradorSesiones
	ahora    func() time.Time
}

func NewCerrarSesion(obtener *ObtenerSesion, cerrador repository.CerradorSesiones, ahora func() time.Time) *CerrarSesion {
	return &CerrarSesion{obtener: obtener, cerrador: cerrador, ahora: ahora}
}

// Ejecutar cierra la sesión vigente del token o devuelve domain.ErrNoAutenticado.
func (c *CerrarSesion) Ejecutar(ctx context.Context, token string) error {
	sesion, err := c.obtener.Ejecutar(ctx, token)
	if err != nil {
		return err
	}
	if err := c.cerrador.Cerrar(ctx, sesion.ID, c.ahora()); err != nil {
		return fmt.Errorf("cerrar sesión: %w", err)
	}
	return nil
}
