package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestObtenerHistoria_Existente(t *testing.T) {
	historias := &fakeHistoriaRepository{historias: map[int64]domain.BacklogItem{3: {ID: 3, Titulo: "Historia"}}}

	h, err := NewObtenerHistoria(historias).Ejecutar(context.Background(), 3)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if h.ID != 3 || h.Titulo != "Historia" {
		t.Errorf("historia inesperada: %+v", h)
	}
}

func TestObtenerHistoria_Inexistente(t *testing.T) {
	historias := &fakeHistoriaRepository{historias: map[int64]domain.BacklogItem{}}

	_, err := NewObtenerHistoria(historias).Ejecutar(context.Background(), 999)

	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Fatalf("se esperaba domain.ErrHistoriaNoEncontrada, se obtuvo %v", err)
	}
}
