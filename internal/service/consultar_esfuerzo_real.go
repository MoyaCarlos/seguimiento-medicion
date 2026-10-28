package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ConsultarEsfuerzoReal devuelve el total acumulado y el detalle de las cargas
// de horas de una historia.
type ConsultarEsfuerzoReal struct {
	historias repository.HistoriaRepository
	horas     repository.RegistroHorasRepository
}

// EsfuerzoReal es el resultado de la consulta; no se guarda.
type EsfuerzoReal struct {
	HistoriaID   int64
	TotalMinutos int
	Registros    []domain.RegistroHorasDetalle
}

func NewConsultarEsfuerzoReal(historias repository.HistoriaRepository, horas repository.RegistroHorasRepository) *ConsultarEsfuerzoReal {
	return &ConsultarEsfuerzoReal{historias: historias, horas: horas}
}

func (s *ConsultarEsfuerzoReal) Ejecutar(ctx context.Context, historiaID int64) (EsfuerzoReal, error) {
	return EsfuerzoReal{}, nil
}
