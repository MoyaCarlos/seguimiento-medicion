package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ObtenerHistoria es la lectura de una historia por su identificador.
type ObtenerHistoria struct {
	historias repository.HistoriaRepository
}

func NewObtenerHistoria(historias repository.HistoriaRepository) *ObtenerHistoria {
	return &ObtenerHistoria{historias: historias}
}

// Ejecutar devuelve la historia o domain.ErrHistoriaNoEncontrada si no existe.
func (s *ObtenerHistoria) Ejecutar(ctx context.Context, id int64) (domain.BacklogItem, error) {
	return s.historias.ObtenerPorID(ctx, id)
}
