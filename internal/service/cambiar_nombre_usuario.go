package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CambiarNombreUsuario cambia el nombre del propio usuario autenticado (HU-17).
type CambiarNombreUsuario struct {
	usuarios repository.ActualizadorNombreUsuario
}

func NewCambiarNombreUsuario(usuarios repository.ActualizadorNombreUsuario) *CambiarNombreUsuario {
	return &CambiarNombreUsuario{usuarios: usuarios}
}

// Ejecutar es el stub del RED de HU-17; el GREEN lo implementa.
func (s *CambiarNombreUsuario) Ejecutar(ctx context.Context, usuarioID int64, nombre string) (domain.User, error) {
	return domain.User{}, nil
}
