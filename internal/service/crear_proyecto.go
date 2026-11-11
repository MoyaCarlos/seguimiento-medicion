package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CrearProyectoInput son los datos de entrada para crear un proyecto. CreadorID
// es la persona con sesión iniciada (HU-21): queda como Scrum Master.
type CrearProyectoInput struct {
	Nombre      string
	Descripcion string
	FechaInicio *time.Time
	FechaFin    *time.Time
	CreadorID   int64
}

// CrearProyecto crea un proyecto y vincula a su creador como Scrum Master.
type CrearProyecto struct {
	proyectos repository.ProjectRepository
}

func NewCrearProyecto(proyectos repository.ProjectRepository) *CrearProyecto {
	return &CrearProyecto{proyectos: proyectos}
}

// Ejecutar valida, persiste el proyecto y registra la vinculación del creador.
func (s *CrearProyecto) Ejecutar(ctx context.Context, input CrearProyectoInput) (domain.Project, error) {
	if input.CreadorID <= 0 {
		return domain.Project{}, domain.ErrNoAutenticado
	}
	proyecto, err := domain.NewProject(input.Nombre, input.Descripcion, input.FechaInicio, input.FechaFin)
	if err != nil {
		return domain.Project{}, err
	}
	return s.proyectos.GuardarConScrumMaster(ctx, proyecto, input.CreadorID)
}
