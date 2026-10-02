package service

import (
	"context"
	"errors"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// AsignarIntegranteInput son los datos de entrada para vincular a un integrante.
type AsignarIntegranteInput struct {
	ProyectoID string
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

// Ejecutar valida el rol y el proyecto, resuelve el integrante y persiste la
// vinculación, rechazando duplicados.
func (s *AsignarIntegrante) Ejecutar(_ context.Context, input AsignarIntegranteInput) (domain.Member, error) {
	if !input.Rol.Valido() {
		return domain.Member{}, domain.ValidationError{Campo: "rol", Mensaje: "el rol debe ser scrum_master o product_builder"}
	}

	proyecto, err := s.proyectos.GetByID(input.ProyectoID)
	if errors.Is(err, repository.ErrNoEncontrado) {
		return domain.Member{}, domain.ErrNoEncontrado
	}
	if err != nil {
		return domain.Member{}, err
	}

	candidato, err := domain.NewUser(input.Nombre)
	if err != nil {
		return domain.Member{}, err
	}
	usuario, err := resolverUsuario(s.usuarios, candidato)
	if err != nil {
		return domain.Member{}, err
	}

	membresia, err := domain.NewMembership(proyecto.ID, usuario.ID, input.Rol)
	if err != nil {
		return domain.Member{}, err
	}
	if err := s.proyectos.AddMember(&membresia); err != nil {
		if errors.Is(err, repository.ErrMiembroDuplicado) {
			return domain.Member{}, domain.ValidationError{Campo: "integrante", Mensaje: "el integrante ya pertenece al proyecto"}
		}
		if errors.Is(err, repository.ErrNoEncontrado) {
			return domain.Member{}, domain.ErrNoEncontrado
		}
		return domain.Member{}, err
	}

	return domain.Member{
		UserID:    usuario.ID,
		ProjectID: proyecto.ID,
		Name:      usuario.Name,
		Role:      input.Rol,
		CreatedAt: membresia.CreatedAt,
	}, nil
}
