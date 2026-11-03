package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestConsultarEsfuerzoReal_DevuelveDetalleYTotal(t *testing.T) {
	historias := &fakeHistoriaRepository{historias: map[int64]domain.BacklogItem{3: {ID: 3, ProyectoID: 1}}}
	horas := &fakeRegistroHorasRepository{registros: []domain.RegistroHoras{
		{ID: 1, HistoriaID: 3, Minutos: 30, Fecha: time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC)},
		{ID: 2, HistoriaID: 4, Minutos: 45},
		{ID: 3, HistoriaID: 3, Minutos: 90, Fecha: time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC)},
		{ID: 4, HistoriaID: 3, Minutos: 120, Fecha: time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC)},
	}}

	esfuerzo, err := NewConsultarEsfuerzoReal(historias, horas).Ejecutar(context.Background(), 3)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if esfuerzo.HistoriaID != 3 || esfuerzo.TotalMinutos != 240 {
		t.Errorf("se esperaba la historia 3 con 240 minutos, se obtuvo %d con %d", esfuerzo.HistoriaID, esfuerzo.TotalMinutos)
	}
	var ids []int64
	for _, r := range esfuerzo.Registros {
		ids = append(ids, r.ID)
	}
	if len(ids) != 3 || ids[0] != 1 || ids[1] != 3 || ids[2] != 4 {
		t.Errorf("se esperaban las cargas 1, 3 y 4 en ese orden, se obtuvo %v", ids)
	}
}

func TestConsultarEsfuerzoReal_HistoriaSinCargas(t *testing.T) {
	historias := &fakeHistoriaRepository{historias: map[int64]domain.BacklogItem{3: {ID: 3, ProyectoID: 1}}}

	esfuerzo, err := NewConsultarEsfuerzoReal(historias, &fakeRegistroHorasRepository{}).Ejecutar(context.Background(), 3)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if esfuerzo.HistoriaID != 3 || esfuerzo.TotalMinutos != 0 || esfuerzo.Registros == nil || len(esfuerzo.Registros) != 0 {
		t.Errorf("se esperaba la historia 3 con total 0 y un detalle vacío no nil, se obtuvo %+v", esfuerzo)
	}
}

func TestConsultarEsfuerzoReal_HistoriaInexistente(t *testing.T) {
	historias := &fakeHistoriaRepository{historias: map[int64]domain.BacklogItem{}}
	horas := &fakeRegistroHorasRepository{}

	_, err := NewConsultarEsfuerzoReal(historias, horas).Ejecutar(context.Background(), 99)

	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Errorf("se esperaba ErrHistoriaNoEncontrada, se obtuvo %v", err)
	}
	if horas.listados != 0 {
		t.Errorf("no se debían consultar las cargas, se consultaron %d veces", horas.listados)
	}
}
