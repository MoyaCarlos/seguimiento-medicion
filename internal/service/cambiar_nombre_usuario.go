package service

import (
	"context"
	"errors"
	"fmt"

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

// Ejecutar valida el nombre con las reglas de HU-20, normaliza la identidad y
// actualiza al usuario de la sesión. El nombre visible se conserva tal como se
// recibió.
func (s *CambiarNombreUsuario) Ejecutar(ctx context.Context, usuarioID int64, nombre string) (domain.User, error) {
	if err := domain.ValidarNombreUsuario(nombre); err != nil {
		return domain.User{}, err
	}
	normalizado := domain.NormalizarNombre(nombre)
	if err := s.usuarios.CambiarNombre(ctx, usuarioID, nombre, normalizado); err != nil {
		if errors.Is(err, repository.ErrUsuarioDuplicado) {
			return domain.User{}, domain.ErrNombreUsuarioEnUso
		}
		return domain.User{}, fmt.Errorf("cambiar nombre del usuario %d: %w", usuarioID, err)
	}
	return domain.User{ID: usuarioID, Nombre: nombre, NombreNormalizado: normalizado}, nil
}
