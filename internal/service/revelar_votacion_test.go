package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

var momentoDeRevelar = time.Date(2026, 10, 28, 10, 0, 0, 0, time.UTC)

func votosDeTodos() []domain.Voto {
	return []domain.Voto{{UsuarioID: carlos, Valor: 3}, {UsuarioID: candela, Valor: 5}, {UsuarioID: jimena, Valor: 8}}
}

func esperarSinRevelar(t *testing.T, votaciones *fakeVotacionRepository) {
	t.Helper()
	if r := votaciones.rondas[3]; r.ReveladaEn != nil {
		t.Errorf("la ronda no debía quedar revelada, quedó en %v", r.ReveladaEn)
	}
}

func TestRevelarVotacion_HistoriaInexistente(t *testing.T) {
	servicio := NewRevelarVotacion(historias(), equipoDelProyecto(), votacionAbierta(votosDeTodos()...))

	_, err := servicio.Ejecutar(context.Background(), 999, carlos, momentoDeRevelar)

	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Errorf("se esperaba ErrHistoriaNoEncontrada, se obtuvo %v", err)
	}
}

func TestRevelarVotacion_NoEsScrumMaster(t *testing.T) {
	for _, usuarioID := range []int64{candela, ajeno} {
		votaciones := votacionAbierta(votosDeTodos()...)
		servicio := NewRevelarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

		_, err := servicio.Ejecutar(context.Background(), 3, usuarioID, momentoDeRevelar)

		if !errors.Is(err, domain.ErrNoEsScrumMaster) {
			t.Errorf("usuario %d: se esperaba ErrNoEsScrumMaster, se obtuvo %v", usuarioID, err)
		}
		esperarSinRevelar(t, votaciones)
	}
}

func TestRevelarVotacion_SinVotacion(t *testing.T) {
	servicio := NewRevelarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), &fakeVotacionRepository{})

	_, err := servicio.Ejecutar(context.Background(), 3, carlos, momentoDeRevelar)

	if !errors.Is(err, domain.ErrVotacionNoIniciada) {
		t.Errorf("se esperaba ErrVotacionNoIniciada, se obtuvo %v", err)
	}
}

func TestRevelarVotacion_FaltaUnVoto(t *testing.T) {
	votaciones := votacionAbierta(domain.Voto{UsuarioID: carlos, Valor: 3}, domain.Voto{UsuarioID: candela, Valor: 5})
	servicio := NewRevelarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	_, err := servicio.Ejecutar(context.Background(), 3, carlos, momentoDeRevelar)

	if !errors.Is(err, domain.ErrFaltanVotos) {
		t.Errorf("se esperaba ErrFaltanVotos, se obtuvo %v", err)
	}
	esperarSinRevelar(t, votaciones)
}

func TestRevelarVotacion_VotaronTodos(t *testing.T) {
	votaciones := votacionAbierta(votosDeTodos()...)
	servicio := NewRevelarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	ronda, err := servicio.Ejecutar(context.Background(), 3, carlos, momentoDeRevelar)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if guardada := votaciones.rondas[3].ReveladaEn; guardada == nil || !guardada.Equal(momentoDeRevelar) {
		t.Errorf("se esperaba guardar la revelación en %v, quedó %v", momentoDeRevelar, guardada)
	}
	promedio, ok := ronda.Promedio()
	if !ronda.Revelada() || !ok || promedio != domain.CalcularPromedio([]int{3, 5, 8}) {
		t.Errorf("se esperaba la ronda revelada con el promedio de 3, 5 y 8, se obtuvo %+v (%+v)", ronda, promedio)
	}
}

func TestRevelarVotacion_YaRevelada(t *testing.T) {
	votaciones := votacionRevelada(votosDeTodos()...)
	antes := *votaciones.rondas[3].ReveladaEn
	servicio := NewRevelarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	_, err := servicio.Ejecutar(context.Background(), 3, carlos, momentoDeRevelar)

	if !errors.Is(err, domain.ErrRondaRevelada) {
		t.Errorf("se esperaba ErrRondaRevelada, se obtuvo %v", err)
	}
	if !votaciones.rondas[3].ReveladaEn.Equal(antes) {
		t.Errorf("no debía cambiar el instante de revelación, quedó %v", votaciones.rondas[3].ReveladaEn)
	}
}

func TestRevelarVotacion_CualquierScrumMasterRevela(t *testing.T) {
	equipo := equipoDelProyecto()
	equipo.listaMiembros[2].Role = domain.RolScrumMaster
	votaciones := votacionAbierta(votosDeTodos()...)
	servicio := NewRevelarVotacion(historias(historiaDelBacklog()), equipo, votaciones)

	if _, err := servicio.Ejecutar(context.Background(), 3, jimena, momentoDeRevelar); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if votaciones.rondas[3].ReveladaEn == nil {
		t.Error("la ronda debía quedar revelada")
	}
}
