package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CerrarProyecto deja el proyecto cerrado con la fecha de hoy (spec 044, #300).
type CerrarProyecto struct {
	proyectos repository.ProjectRepository
	sprints   repository.SprintRepository
	reloj     func() time.Time
}

func NewCerrarProyecto(proyectos repository.ProjectRepository, sprints repository.SprintRepository, reloj func() time.Time) *CerrarProyecto {
	return &CerrarProyecto{proyectos: proyectos, sprints: sprints, reloj: reloj}
}

// Ejecutar verifica, en este orden, que el proyecto exista (404), que quien
// cierra sea su Scrum Master (403) y que se pueda cerrar (409: ya cerrado o
// con un Sprint Activo). El permiso usa verificarScrumMaster y no el guardia de
// escritura, para que cerrar uno ya cerrado informe «ya está cerrado».
func (s *CerrarProyecto) Ejecutar(ctx context.Context, proyectoID, usuarioID int64) (domain.Project, error) {
	proyecto, err := s.proyectos.ObtenerPorID(ctx, proyectoID)
	if err != nil {
		return domain.Project{}, err
	}
	if err := verificarScrumMaster(ctx, s.proyectos, proyectoID, usuarioID); err != nil {
		return domain.Project{}, err
	}
	sprints, err := s.sprints.ListarPorProyecto(ctx, proyectoID)
	if err != nil {
		return domain.Project{}, err
	}
	cerrado, err := proyecto.Cerrar(s.reloj(), sprints)
	if err != nil {
		return domain.Project{}, err
	}
	if err := s.proyectos.Cerrar(ctx, proyectoID, *cerrado.FechaCierre); err != nil {
		return domain.Project{}, err
	}
	return cerrado, nil
}
