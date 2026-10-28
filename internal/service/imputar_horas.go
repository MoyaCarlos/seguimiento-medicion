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
	historia, err := s.historias.ObtenerPorID(ctx, in.HistoriaID)
	if err != nil {
		return domain.RegistroHorasDetalle{}, err
	}
	sprint, err := s.sprints.ObtenerPorID(ctx, *historia.SprintID)
	if err != nil {
		return domain.RegistroHorasDetalle{}, err
	}
	nombre, err := s.nombreDelIntegrante(ctx, historia.ProyectoID, in.IntegranteID)
	if err != nil {
		return domain.RegistroHorasDetalle{}, err
	}
	registro, err := domain.NewRegistroHoras(historia, sprint, in.IntegranteID, in.Horas, in.Actividad, in.Fecha, ahora)
	if err != nil {
		return domain.RegistroHorasDetalle{}, err
	}
	registro, err = s.horas.Registrar(ctx, registro)
	if err != nil {
		return domain.RegistroHorasDetalle{}, err
	}
	return domain.RegistroHorasDetalle{RegistroHoras: registro, NombreIntegrante: nombre}, nil
}

func (s *ImputarHoras) nombreDelIntegrante(ctx context.Context, proyectoID, usuarioID int64) (string, error) {
	integrantes, err := s.integrantes.ListarIntegrantes(ctx, proyectoID)
	if err != nil {
		return "", err
	}
	for _, m := range integrantes {
		if m.UserID == usuarioID {
			return m.Nombre, nil
		}
	}
	return "", nil
}
