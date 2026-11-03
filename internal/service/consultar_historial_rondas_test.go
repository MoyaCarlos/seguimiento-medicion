package service

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func votacionConRondaUnoReveladaYDosAbierta() *fakeVotacionRepository {
	revelada := time.Date(2026, 11, 3, 15, 4, 5, 0, time.UTC)
	return &fakeVotacionRepository{
		anteriores: []domain.RondaVotacion{{ID: 1, HistoriaID: 3, Numero: 1, ReveladaEn: &revelada, Votos: votosDe(3, 8)}},
		rondas:     map[int64]domain.RondaVotacion{3: {ID: 2, HistoriaID: 3, Numero: 2, Votos: votosDe(13)}},
	}
}

func TestConsultarHistorialRondas_SoloRondasReveladas(t *testing.T) {
	votaciones := votacionConRondaUnoReveladaYDosAbierta()
	servicio := NewConsultarHistorialRondas(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	historial, err := servicio.Ejecutar(context.Background(), 3)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if historial.HistoriaID != 3 || len(historial.Rondas) != 1 {
		t.Fatalf("se esperaba solo la ronda 1 de la historia 3, se obtuvo %+v", historial)
	}
	ronda := historial.Rondas[0]
	if ronda.Numero != 1 || !ronda.ReveladaEn.Equal(*votaciones.anteriores[0].ReveladaEn) {
		t.Errorf("se esperaba la ronda 1 con su instante de revelación, se obtuvo %+v", ronda)
	}
	esperados := []VotoRevelado{
		{UsuarioID: carlos, Nombre: "Carlos", Rol: domain.RolScrumMaster, Valor: 3},
		{UsuarioID: candela, Nombre: "Candela", Rol: domain.RolProductBuilder, Valor: 8},
	}
	if !slices.Equal(ronda.Votos, esperados) {
		t.Errorf("se esperaban los votos %+v, se obtuvo %+v", esperados, ronda.Votos)
	}
}

func TestConsultarHistorialRondas_SinRondasReveladas(t *testing.T) {
	servicio := NewConsultarHistorialRondas(historias(historiaDelBacklog()), equipoDelProyecto(), votacionAbierta(votosDe(3)...))

	historial, err := servicio.Ejecutar(context.Background(), 3)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if historial.HistoriaID != 3 || len(historial.Rondas) != 0 {
		t.Errorf("se esperaba un historial vacío de la historia 3, se obtuvo %+v", historial)
	}
}

func TestConsultarHistorialRondas_HistoriaInexistente(t *testing.T) {
	servicio := NewConsultarHistorialRondas(historias(), equipoDelProyecto(), votacionRevelada(votosDe(3)...))

	if _, err := servicio.Ejecutar(context.Background(), 999); !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Errorf("se esperaba ErrHistoriaNoEncontrada, se obtuvo %v", err)
	}
}

func TestConsultarHistorialRondas_NoEscribe(t *testing.T) {
	votaciones := votacionConRondaUnoReveladaYDosAbierta()
	antes := copiaDeRonda(votaciones.rondas[3])
	repoHistorias := historias(historiaDelBacklog())
	servicio := NewConsultarHistorialRondas(repoHistorias, equipoDelProyecto(), votaciones)

	if _, err := servicio.Ejecutar(context.Background(), 3); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if !reflect.DeepEqual(votaciones.rondas[3], antes) {
		t.Errorf("la consulta no debía cambiar la ronda actual: antes %+v, ahora %+v", antes, votaciones.rondas[3])
	}
	if len(repoHistorias.actualizaciones) != 0 {
		t.Error("la consulta no debía escribir en historias")
	}
}
