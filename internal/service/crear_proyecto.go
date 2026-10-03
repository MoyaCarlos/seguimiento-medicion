package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CrearProyectoInput son los datos de entrada para crear un proyecto.
type CrearProyectoInput struct {
	Nombre      string
	Descripcion string
	FechaInicio *time.Time
	FechaFin    *time.Time
	Creador     string
}

// CrearProyecto crea un proyecto y vincula a su creador como Scrum Master.
type CrearProyecto struct {
	proyectos repository.ProjectRepository
	usuarios  repository.UserRepository
}

func NewCrearProyecto(proyectos repository.ProjectRepository, usuarios repository.UserRepository) *CrearProyecto {
	return &CrearProyecto{proyectos: proyectos, usuarios: usuarios}
}

// Ejecutar valida, persiste el proyecto y registra la vinculación del creador.
func (s *CrearProyecto) Ejecutar(_ context.Context, input CrearProyectoInput) (domain.Project, error) {
	proyecto, err := domain.NewProject(input.Nombre, input.Descripcion, input.FechaInicio, input.FechaFin)
	if err != nil {
		return domain.Project{}, err
	}
	creador, err := domain.NewUser(input.Creador)
	if err != nil {
		return domain.Project{}, err
	}

	usuario, err := resolverUsuario(s.usuarios, creador)
	if err != nil {
		return domain.Project{}, err
	}

	if err := s.proyectos.CreateWithScrumMaster(&proyecto, usuario.ID); err != nil {
		return domain.Project{}, err
	}
	return proyecto, nil
}
