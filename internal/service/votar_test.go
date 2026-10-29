package service

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// El equipo del proyecto 1 de historiaDelBacklog (historia 3).
const (
	carlos  int64 = 1
	candela int64 = 2
	jimena  int64 = 3
	ajeno   int64 = 9
)

func equipoDelProyecto() *fakeProjectRepository {
	return &fakeProjectRepository{listaMiembros: []domain.Member{
		{UserID: carlos, ProjectID: 1, Nombre: "Carlos", Role: domain.RolScrumMaster},
		{UserID: candela, ProjectID: 1, Nombre: "Candela", Role: domain.RolProductBuilder},
		{UserID: jimena, ProjectID: 1, Nombre: "Jimena", Role: domain.RolProductBuilder},
	}}
}

func historias(hs ...domain.BacklogItem) *fakeHistoriaRepository {
	repo := &fakeHistoriaRepository{historias: map[int64]domain.BacklogItem{}}
	for _, h := range hs {
		repo.historias[h.ID] = h
	}
	return repo
}

// votacionAbierta siembra la ronda 1, sin revelar, de la historia 3.
func votacionAbierta(votos ...domain.Voto) *fakeVotacionRepository {
	return &fakeVotacionRepository{rondas: map[int64]domain.RondaVotacion{
		3: {ID: 1, HistoriaID: 3, Numero: 1, Votos: votos},
	}}
}

func votacionRevelada(votos ...domain.Voto) *fakeVotacionRepository {
	repo := votacionAbierta(votos...)
	r := repo.rondas[3]
	revelada := time.Date(2026, 10, 27, 15, 4, 5, 0, time.UTC)
	r.ReveladaEn = &revelada
	repo.rondas[3] = r
	return repo
}

func votosGuardados(repo *fakeVotacionRepository) []domain.Voto {
	return repo.rondas[3].Votos
}

func TestVotar_HistoriaInexistente(t *testing.T) {
	votaciones := votacionAbierta()
	servicio := NewVotar(historias(), equipoDelProyecto(), votaciones)

	_, err := servicio.Ejecutar(context.Background(), VotarInput{HistoriaID: 999, UsuarioID: candela, Valor: 8})

	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Errorf("se esperaba ErrHistoriaNoEncontrada, se obtuvo %v", err)
	}
}

func TestVotar_NoEsIntegrante(t *testing.T) {
	votaciones := votacionAbierta()
	servicio := NewVotar(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	_, err := servicio.Ejecutar(context.Background(), VotarInput{HistoriaID: 3, UsuarioID: ajeno, Valor: 8})

	if !errors.Is(err, domain.ErrNoEsIntegrante) {
		t.Errorf("se esperaba ErrNoEsIntegrante, se obtuvo %v", err)
	}
	if len(votosGuardados(votaciones)) != 0 {
		t.Errorf("no se esperaba guardar el voto, quedó %+v", votosGuardados(votaciones))
	}
}

func TestVotar_SinVotacion(t *testing.T) {
	servicio := NewVotar(historias(historiaDelBacklog()), equipoDelProyecto(), &fakeVotacionRepository{})

	_, err := servicio.Ejecutar(context.Background(), VotarInput{HistoriaID: 3, UsuarioID: candela, Valor: 8})

	if !errors.Is(err, domain.ErrVotacionNoIniciada) {
		t.Errorf("se esperaba ErrVotacionNoIniciada, se obtuvo %v", err)
	}
}

func TestVotar_GuardaElVoto(t *testing.T) {
	votaciones := votacionAbierta()
	servicio := NewVotar(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	ronda, err := servicio.Ejecutar(context.Background(), VotarInput{HistoriaID: 3, UsuarioID: candela, Valor: 8})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	esperados := []domain.Voto{{UsuarioID: candela, Valor: 8}}
	if !slices.Equal(votosGuardados(votaciones), esperados) {
		t.Errorf("se esperaba guardar %+v, quedó %+v", esperados, votosGuardados(votaciones))
	}
	if ronda.ID != 1 || !slices.Equal(ronda.Votos, esperados) {
		t.Errorf("se esperaba la ronda 1 con %+v, se obtuvo %+v", esperados, ronda)
	}
}

func TestVotar_CambiaElVoto(t *testing.T) {
	votaciones := votacionAbierta(domain.Voto{UsuarioID: carlos, Valor: 3}, domain.Voto{UsuarioID: candela, Valor: 8})
	servicio := NewVotar(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	if _, err := servicio.Ejecutar(context.Background(), VotarInput{HistoriaID: 3, UsuarioID: candela, Valor: 13}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	esperados := []domain.Voto{{UsuarioID: carlos, Valor: 3}, {UsuarioID: candela, Valor: 13}}
	if !slices.Equal(votosGuardados(votaciones), esperados) {
		t.Errorf("se esperaba %+v, quedó %+v", esperados, votosGuardados(votaciones))
	}
}

func TestVotar_CartaFueraDeLaBaraja(t *testing.T) {
	votaciones := votacionAbierta()
	servicio := NewVotar(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	_, err := servicio.Ejecutar(context.Background(), VotarInput{HistoriaID: 3, UsuarioID: candela, Valor: 4})

	if !errors.Is(err, domain.ErrVotoInvalido) {
		t.Errorf("se esperaba ErrVotoInvalido, se obtuvo %v", err)
	}
	if len(votosGuardados(votaciones)) != 0 {
		t.Errorf("no se esperaba guardar el voto, quedó %+v", votosGuardados(votaciones))
	}
}

func TestVotar_RondaRevelada(t *testing.T) {
	votaciones := votacionRevelada(domain.Voto{UsuarioID: candela, Valor: 8})
	servicio := NewVotar(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	_, err := servicio.Ejecutar(context.Background(), VotarInput{HistoriaID: 3, UsuarioID: candela, Valor: 13})

	if !errors.Is(err, domain.ErrRondaRevelada) {
		t.Errorf("se esperaba ErrRondaRevelada, se obtuvo %v", err)
	}
	if !slices.Equal(votosGuardados(votaciones), []domain.Voto{{UsuarioID: candela, Valor: 8}}) {
		t.Errorf("los votos no debían cambiar, quedaron %+v", votosGuardados(votaciones))
	}
}

func TestVotar_ElScrumMasterTambienVota(t *testing.T) {
	votaciones := votacionAbierta()
	servicio := NewVotar(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	if _, err := servicio.Ejecutar(context.Background(), VotarInput{HistoriaID: 3, UsuarioID: carlos, Valor: 5}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if !slices.Equal(votosGuardados(votaciones), []domain.Voto{{UsuarioID: carlos, Valor: 5}}) {
		t.Errorf("se esperaba el voto 5 de Carlos, quedó %+v", votosGuardados(votaciones))
	}
}
