package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ImputarHoras registra horas trabajadas en una historia del Sprint activo.
type ImputarHoras struct {
	historias   repository.HistoriaRepository
	sprints     repository.SprintRepository
	integrantes repository.IntegrantesRepository
	horas       repository.RegistroHorasRepository
}

// ImputarHorasInput son los datos de una carga. Horas y Fecha llegan como texto
// para que el dominio las valide en el orden de la spec.
type ImputarHorasInput struct {
	HistoriaID   int64
	IntegranteID int64
	Fecha        string
	Actividad    string
	Horas        string
}

func NewImputarHoras(historias repository.HistoriaRepository, sprints repository.SprintRepository, integrantes repository.IntegrantesRepository, horas repository.RegistroHorasRepository) *ImputarHoras {
	return &ImputarHoras{historias: historias, sprints: sprints, integrantes: integrantes, horas: horas}
}

func (s *ImputarHoras) Ejecutar(ctx context.Context, in ImputarHorasInput, ahora time.Time) (domain.RegistroHorasDetalle, error) {
	return domain.RegistroHorasDetalle{}, nil
}
