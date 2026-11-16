package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ConsultarEsfuerzoReal devuelve el total acumulado y el detalle de las cargas
// de horas de una historia.
type ConsultarEsfuerzoReal struct {
	historias   repository.HistoriaRepository
	horas       repository.RegistroHorasRepository
	integrantes repository.IntegrantesRepository
}

// EsfuerzoReal es el resultado de la consulta; no se guarda.
type EsfuerzoReal struct {
	HistoriaID   int64
	TotalMinutos int
	Registros    []domain.RegistroHorasDetalle
}

func NewConsultarEsfuerzoReal(historias repository.HistoriaRepository, horas repository.RegistroHorasRepository, integrantes repository.IntegrantesRepository) *ConsultarEsfuerzoReal {
	return &ConsultarEsfuerzoReal{historias: historias, horas: horas, integrantes: integrantes}
}

func (s *ConsultarEsfuerzoReal) Ejecutar(ctx context.Context, historiaID, consultanteID int64) (EsfuerzoReal, error) {
	historia, err := s.historias.ObtenerPorID(ctx, historiaID)
	if err != nil {
		return EsfuerzoReal{}, err
	}
	if err := exigirLector(ctx, s.integrantes, historia.ProyectoID, consultanteID); err != nil {
		return EsfuerzoReal{}, err
	}
	registros, err := s.horas.ListarPorHistoria(ctx, historiaID)
	if err != nil {
		return EsfuerzoReal{}, err
	}
	return EsfuerzoReal{HistoriaID: historiaID, TotalMinutos: domain.TotalMinutos(registros), Registros: registros}, nil
}
