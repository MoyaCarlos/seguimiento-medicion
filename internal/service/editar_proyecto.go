package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// EditarProyectoInput son los datos editables de un proyecto.
type EditarProyectoInput struct {
	ID          int64
	Nombre      string
	Descripcion string
	FechaInicio *time.Time
	FechaFin    *time.Time
	// HorasPorSPEnviadas distingue no mencionar el parámetro (lo conserva) de
	// enviarlo sin valor (HorasPorSP nil, lo borra).
	HorasPorSPEnviadas bool
	HorasPorSP         *string
	// UsuarioID es quien edita (la sesión, HU-21); debe ser el Scrum Master.
	UsuarioID int64
}

// EditarProyecto aplica una edición simple reutilizando las validaciones del alta.
type EditarProyecto struct {
	proyectos repository.ProjectRepository
}

func NewEditarProyecto(proyectos repository.ProjectRepository) *EditarProyecto {
	return &EditarProyecto{proyectos: proyectos}
}

// Autorizar devuelve ErrProyectoNoEncontrado (404) o el error de permisos (403)
// sin tocar nada: el handler lo llama antes de leer el cuerpo, así un Product
// Builder recibe 403 aunque mande datos inválidos (HU-21, FR-006).
func (s *EditarProyecto) Autorizar(ctx context.Context, proyectoID, usuarioID int64) error {
	_, err := s.proyectoDelScrumMaster(ctx, proyectoID, usuarioID)
	return err
}

func (s *EditarProyecto) proyectoDelScrumMaster(ctx context.Context, proyectoID, usuarioID int64) (domain.Project, error) {
	proyecto, err := s.proyectos.ObtenerPorID(ctx, proyectoID)
	if err != nil {
		return domain.Project{}, err
	}
	if err := exigirScrumMaster(ctx, s.proyectos, proyectoID, usuarioID); err != nil {
		return domain.Project{}, err
	}
	return proyecto, nil
}

// Ejecutar carga el proyecto, exige que quien edita sea su Scrum Master, valida
// los datos y persiste los cambios.
func (s *EditarProyecto) Ejecutar(ctx context.Context, input EditarProyectoInput) (domain.Project, error) {
	actual, err := s.proyectoDelScrumMaster(ctx, input.ID, input.UsuarioID)
	if err != nil {
		return domain.Project{}, err
	}

	editado, err := actual.ConDatosEditados(input.Nombre, input.Descripcion, input.FechaInicio, input.FechaFin)
	if err != nil {
		return domain.Project{}, err
	}

	if input.HorasPorSPEnviadas {
		editado, err = editado.ConHorasPorStoryPoint(input.HorasPorSP)
		if err != nil {
			return domain.Project{}, err
		}
	}

	if err := s.proyectos.Actualizar(ctx, editado); err != nil {
		return domain.Project{}, err
	}
	return editado, nil
}
