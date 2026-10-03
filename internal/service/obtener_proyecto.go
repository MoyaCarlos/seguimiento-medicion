package service

import (
	"context"
	"errors"

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

// Ejecutar devuelve el proyecto o domain.ErrNoEncontrado si no existe.
func (s *ObtenerProyecto) Ejecutar(_ context.Context, id int64) (domain.Project, error) {
	proyecto, err := s.proyectos.GetByID(id)
	if errors.Is(err, repository.ErrNoEncontrado) {
		return domain.Project{}, domain.ErrNoEncontrado
	}
	if err != nil {
		return domain.Project{}, err
	}
	return *proyecto, nil
}
