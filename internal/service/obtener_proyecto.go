package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ObtenerProyecto es la lectura de un proyecto por su identificador.
type ObtenerProyecto struct {
	proyectos repository.ProjectRepository
}

func NewObtenerProyecto(proyectos repository.ProjectRepository) *ObtenerProyecto {
	return &ObtenerProyecto{proyectos: proyectos}
}

// Ejecutar devuelve el proyecto o domain.ErrProyectoNoEncontrado si no existe.
func (s *ObtenerProyecto) Ejecutar(ctx context.Context, id int64) (domain.Project, error) {
	return s.proyectos.ObtenerPorID(ctx, id)
}

// ComoIntegrante devuelve el proyecto y el rol de usuarioID en él (panel del
// proyecto, HU-21): domain.ErrProyectoNoEncontrado si no existe y
// domain.ErrNoEsIntegrante si usuarioID no lo integra.
func (s *ObtenerProyecto) ComoIntegrante(ctx context.Context, id, usuarioID int64) (domain.Project, domain.Role, error) {
	proyecto, err := s.proyectos.ObtenerPorID(ctx, id)
	if err != nil {
		return domain.Project{}, "", err
	}
	rol, integra, err := rolEnProyecto(ctx, s.proyectos, id, usuarioID)
	if err != nil {
		return domain.Project{}, "", err
	}
	if !integra {
		return domain.Project{}, "", domain.ErrNoEsIntegrante
	}
	return proyecto, rol, nil
}
