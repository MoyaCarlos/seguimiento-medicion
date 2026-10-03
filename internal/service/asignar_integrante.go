package service

import (
	"context"
	"errors"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// AsignarIntegranteInput son los datos de entrada para vincular a un integrante.
type AsignarIntegranteInput struct {
	ProyectoID int64
	Nombre     string
	Rol        domain.Role
}

// AsignarIntegrante vincula a un integrante (existente o nuevo) con un rol.
type AsignarIntegrante struct {
	proyectos repository.ProjectRepository
	usuarios  repository.UserRepository
}

func NewAsignarIntegrante(proyectos repository.ProjectRepository, usuarios repository.UserRepository) *AsignarIntegrante {
	return &AsignarIntegrante{proyectos: proyectos, usuarios: usuarios}
}

// Ejecutar verifica el proyecto, resuelve el integrante y persiste la
// vinculación (el rol se valida en el dominio), rechazando duplicados.
func (s *AsignarIntegrante) Ejecutar(ctx context.Context, input AsignarIntegranteInput) (domain.Member, error) {
	proyecto, err := s.proyectos.ObtenerPorID(ctx, input.ProyectoID)
	if err != nil {
		return domain.Member{}, err
	}

	candidato, err := domain.NewUser(input.Nombre)
	if err != nil {
		return domain.Member{}, err
	}
	usuario, err := resolverUsuario(ctx, s.usuarios, candidato)
	if err != nil {
		return domain.Member{}, err
	}

	membresia, err := domain.NewMembership(proyecto.ID, usuario.ID, input.Rol)
	if err != nil {
		return domain.Member{}, err
	}
	if err := s.proyectos.AgregarIntegrante(ctx, membresia); err != nil {
		if errors.Is(err, repository.ErrMiembroDuplicado) {
			return domain.Member{}, domain.ValidationError{Campo: "integrante", Mensaje: "el integrante ya pertenece al proyecto"}
		}
		return domain.Member{}, err
	}

	return domain.Member{
		UserID:    usuario.ID,
		ProjectID: proyecto.ID,
		Nombre:    usuario.Nombre,
		Role:      input.Rol,
		CreadoEn:  membresia.CreadoEn,
	}, nil
}
