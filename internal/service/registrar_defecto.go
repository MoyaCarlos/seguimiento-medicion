package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// RegistrarDefecto registra un defecto Pendiente con el Sprint activo del
// proyecto como Sprint de detección (HU-08).
type RegistrarDefecto struct {
	sprints   repository.SprintRepository
	historias repository.HistoriaRepository
	defectos  repository.DefectoRepository
}

// RegistrarDefectoInput son los datos de un defecto. El estado y los Sprints
// los asigna el sistema.
type RegistrarDefectoInput struct {
	ProyectoID  int64
	HistoriaID  int64
	Descripcion string
	Severidad   string
}

func NewRegistrarDefecto(sprints repository.SprintRepository, historias repository.HistoriaRepository, defectos repository.DefectoRepository) *RegistrarDefecto {
	return &RegistrarDefecto{sprints: sprints, historias: historias, defectos: defectos}
}

func (s *RegistrarDefecto) Ejecutar(ctx context.Context, in RegistrarDefectoInput) (domain.DefectoDetalle, error) {
	return domain.DefectoDetalle{}, nil
}
