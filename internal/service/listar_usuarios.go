package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type listadorUsuarios interface {
	ListarRegistrados(ctx context.Context) ([]domain.User, error)
}

// ListarUsuarios devuelve los usuarios registrados, entre los que el Scrum
// Master elige a quién sumar al equipo (HU-21).
type ListarUsuarios struct {
	usuarios listadorUsuarios
}

func NewListarUsuarios(usuarios listadorUsuarios) *ListarUsuarios {
	return &ListarUsuarios{usuarios: usuarios}
}

// Ejecutar devuelve los usuarios con cuenta; un arreglo vacío si no hay.
func (s *ListarUsuarios) Ejecutar(ctx context.Context) ([]domain.User, error) {
	usuarios, err := s.usuarios.ListarRegistrados(ctx)
	if err != nil {
		return nil, err
	}
	if usuarios == nil {
		usuarios = []domain.User{}
	}
	return usuarios, nil
}
