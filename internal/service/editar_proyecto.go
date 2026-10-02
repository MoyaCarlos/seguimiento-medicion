package service

import (
	"context"
	"errors"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// EditarProyectoInput son los datos editables de un proyecto.
type EditarProyectoInput struct {
	ID          string
	Nombre      string
	Descripcion string
	FechaInicio *time.Time
	FechaFin    *time.Time
}

// EditarProyecto aplica una edición simple reutilizando las validaciones del alta.
type EditarProyecto struct {
	proyectos repository.ProjectRepository
}

func NewEditarProyecto(proyectos repository.ProjectRepository) *EditarProyecto {
	return &EditarProyecto{proyectos: proyectos}
}

// Ejecutar carga el proyecto, valida los datos y persiste los cambios.
func (s *EditarProyecto) Ejecutar(_ context.Context, input EditarProyectoInput) (domain.Project, error) {
	actual, err := s.proyectos.GetByID(input.ID)
	if errors.Is(err, repository.ErrNoEncontrado) {
		return domain.Project{}, domain.ErrNoEncontrado
	}
	if err != nil {
		return domain.Project{}, err
	}

	editado, err := actual.ConDatosEditados(input.Nombre, input.Descripcion, input.FechaInicio, input.FechaFin)
	if err != nil {
		return domain.Project{}, err
	}

	if err := s.proyectos.Update(&editado); err != nil {
		if errors.Is(err, repository.ErrNoEncontrado) {
			return domain.Project{}, domain.ErrNoEncontrado
		}
		return domain.Project{}, err
	}
	return editado, nil
}
