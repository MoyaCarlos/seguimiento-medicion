package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ObtenerHistoria es la lectura de una historia por su identificador.
type ObtenerHistoria struct {
	historias   repository.HistoriaRepository
	integrantes repository.IntegrantesRepository
}

func NewObtenerHistoria(historias repository.HistoriaRepository, integrantes repository.IntegrantesRepository) *ObtenerHistoria {
	return &ObtenerHistoria{historias: historias, integrantes: integrantes}
}

// ComoIntegrante devuelve la historia a quien integra su proyecto:
// domain.ErrHistoriaNoEncontrada si no existe y domain.ErrNoEsIntegrante si
// usuarioID no integra el proyecto de la historia.
func (s *ObtenerHistoria) ComoIntegrante(ctx context.Context, id, usuarioID int64) (domain.BacklogItem, error) {
	historia, err := s.historias.ObtenerPorID(ctx, id)
	if err != nil {
		return domain.BacklogItem{}, err
	}
	if err := exigirLector(ctx, s.integrantes, historia.ProyectoID, usuarioID); err != nil {
		return domain.BacklogItem{}, err
	}
	return historia, nil
}

// Ejecutar devuelve la historia o domain.ErrHistoriaNoEncontrada si no existe.
func (s *ObtenerHistoria) Ejecutar(ctx context.Context, id int64) (domain.BacklogItem, error) {
	return s.historias.ObtenerPorID(ctx, id)
}
