package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CrearHistoriaInput son los datos de entrada del caso de uso.
type CrearHistoriaInput struct {
	ProyectoID   int64
	Titulo       string
	Descripcion  string
	Prioridad    domain.Prioridad
	ValorNegocio *int
}

// CrearHistoriaBacklog crea y persiste una historia del Product Backlog.
type CrearHistoriaBacklog struct {
	repo      repository.BacklogRepository
	proyectos repository.ProjectRepository
}

func NewCrearHistoriaBacklog(repo repository.BacklogRepository, proyectos repository.ProjectRepository) *CrearHistoriaBacklog {
	return &CrearHistoriaBacklog{repo: repo, proyectos: proyectos}
}

// Ejecutar construye la historia (devolviendo el error de validación si
// corresponde), verifica que el proyecto exista y la persiste a través del
// repositorio.
func (s *CrearHistoriaBacklog) Ejecutar(ctx context.Context, input CrearHistoriaInput) (domain.BacklogItem, error) {
	item, err := domain.NewBacklogItem(input.ProyectoID, input.Titulo, input.Descripcion, input.Prioridad, input.ValorNegocio)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	if _, err := s.proyectos.ObtenerPorID(ctx, input.ProyectoID); err != nil {
		return domain.BacklogItem{}, err
	}
	return s.repo.Guardar(ctx, item)
}
