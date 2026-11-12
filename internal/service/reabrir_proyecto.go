package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ReabrirProyecto deja vacía la fecha de cierre del proyecto (spec 044, #300).
type ReabrirProyecto struct {
	proyectos repository.ProjectRepository
}

func NewReabrirProyecto(proyectos repository.ProjectRepository) *ReabrirProyecto {
	return &ReabrirProyecto{proyectos: proyectos}
}

// Ejecutar verifica, en este orden, que el proyecto exista (404), que quien
// reabre sea su Scrum Master (403) y que el proyecto esté cerrado (409). Es la
// única escritura permitida sobre un proyecto cerrado, por eso usa
// verificarScrumMaster y no el guardia de escritura.
func (s *ReabrirProyecto) Ejecutar(ctx context.Context, proyectoID, usuarioID int64) (domain.Project, error) {
	proyecto, err := s.proyectos.ObtenerPorID(ctx, proyectoID)
	if err != nil {
		return domain.Project{}, err
	}
	if err := verificarScrumMaster(ctx, s.proyectos, proyectoID, usuarioID); err != nil {
		return domain.Project{}, err
	}
	abierto, err := proyecto.Reabrir()
	if err != nil {
		return domain.Project{}, err
	}
	if err := s.proyectos.Reabrir(ctx, proyectoID); err != nil {
		return domain.Project{}, err
	}
	return abierto, nil
}
