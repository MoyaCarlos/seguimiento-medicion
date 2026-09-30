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
	repo repository.BacklogRepository
}

func NewCrearHistoriaBacklog(repo repository.BacklogRepository) *CrearHistoriaBacklog {
	return &CrearHistoriaBacklog{repo: repo}
}

// Ejecutar construye la historia (devolviendo el error de validación si
// corresponde) y la persiste a través del repositorio.
func (s *CrearHistoriaBacklog) Ejecutar(ctx context.Context, input CrearHistoriaInput) (domain.BacklogItem, error) {
	item, err := domain.NewBacklogItem(input.ProyectoID, input.Titulo, input.Descripcion, input.Prioridad, input.ValorNegocio)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	return s.repo.Guardar(ctx, item)
}
