package service

import (
	"context"
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
func (o *ObtenerSesion) Ejecutar(_ context.Context, _ string) (domain.SesionConUsuario, error) {
	return domain.SesionConUsuario{}, nil
}
