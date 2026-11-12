package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// CrearHistoriaInput son los datos de entrada del caso de uso.
type CrearHistoriaInput struct {
	ProyectoID          int64
	Titulo              string
	Descripcion         string
	Prioridad           domain.Prioridad
	ValorNegocio        *int
	CriteriosAceptacion *string
	// UsuarioID es quien crea la historia (la sesión, HU-21); debe integrar el proyecto.
	UsuarioID int64
}

// CrearHistoriaBacklog crea y persiste una historia del Product Backlog.
type CrearHistoriaBacklog struct {
	repo      repository.BacklogRepository
	proyectos repository.ProjectRepository
}

func NewCrearHistoriaBacklog(repo repository.BacklogRepository, proyectos repository.ProjectRepository) *CrearHistoriaBacklog {
	return &CrearHistoriaBacklog{repo: repo, proyectos: proyectos}
}

// Ejecutar verifica, en este orden, que el proyecto exista (404), que quien
// crea lo integre (403) y los datos de la historia (400), y la persiste. Sin
// proyecto en el pedido, la validación de la historia lo informa (400).
func (s *CrearHistoriaBacklog) Ejecutar(ctx context.Context, input CrearHistoriaInput) (domain.BacklogItem, error) {
	if input.ProyectoID > 0 {
		if _, err := s.proyectos.ObtenerPorID(ctx, input.ProyectoID); err != nil {
			return domain.BacklogItem{}, err
		}
		if err := exigirIntegrante(ctx, s.proyectos, input.ProyectoID, input.UsuarioID); err != nil {
			return domain.BacklogItem{}, err
		}
	}
	item, err := domain.NewBacklogItem(input.ProyectoID, input.Titulo, input.Descripcion, input.Prioridad, input.ValorNegocio)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	item, err = item.ConCriteriosAceptacion(input.CriteriosAceptacion)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	return s.repo.Guardar(ctx, item)
}
