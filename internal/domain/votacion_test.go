package domain

import (
	"errors"
	"slices"
	"testing"
	"time"
)

var momentoRevelacion = time.Date(2026, 10, 27, 15, 4, 5, 0, time.UTC)

func rondaAbierta(t *testing.T) RondaVotacion {
	t.Helper()
	r, err := NuevaRonda(100)
	if err != nil {
		t.Fatalf("no se esperaba error al crear la ronda: %v", err)
	}
	return r
}

func votar(t *testing.T, r *RondaVotacion, usuarioID int64, valor int) {
	t.Helper()
	if err := r.Votar(usuarioID, valor); err != nil {
		t.Fatalf("no se esperaba error al votar %d: %v", valor, err)
	}
}

func rondaRevelada(t *testing.T) RondaVotacion {
	t.Helper()
	r := rondaAbierta(t)
	votar(t, &r, 1, 8)
	if err := r.Revelar([]int64{1}, momentoRevelacion); err != nil {
		t.Fatalf("no se esperaba error al revelar: %v", err)
	}
	return r
}

func TestNuevaRonda_EsLaPrimeraSinVotosNiRevelar(t *testing.T) {
	r, err := NuevaRonda(100)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if r.HistoriaID != 100 || r.Numero != 1 {
		t.Errorf("se esperaba la ronda 1 de la historia 100, se obtuvo %+v", r)
	}
	if len(r.Votos) != 0 || r.ReveladaEn != nil || r.Revelada() {
		t.Errorf("una ronda nueva no debe tener votos ni estar revelada: %+v", r)
	}
}

func TestNuevaRonda_HistoriaObligatoria(t *testing.T) {
	for _, id := range []int64{0, -1} {
		_, err := NuevaRonda(id)
		esperarErrorDeValidacion(t, err, "historia_id")
	}
}

func TestRondaVotacion_Votar_AgregaElVoto(t *testing.T) {
	r := rondaAbierta(t)

	votar(t, &r, 2, 8)

	if !slices.Equal(r.Votos, []Voto{{UsuarioID: 2, Valor: 8}}) {
		t.Errorf("se esperaba el voto 8 del usuario 2, se obtuvo %+v", r.Votos)
	}
}

func TestRondaVotacion_Votar_ReemplazaElVotoPropio(t *testing.T) {
	r := rondaAbierta(t)
	votar(t, &r, 1, 3)
	votar(t, &r, 2, 8)

	votar(t, &r, 2, 13)

	esperados := []Voto{{UsuarioID: 1, Valor: 3}, {UsuarioID: 2, Valor: 13}}
	if !slices.Equal(r.Votos, esperados) {
		t.Errorf("se esperaba %+v, se obtuvo %+v", esperados, r.Votos)
	}
}

func TestRondaVotacion_Votar_CartaFueraDeLaBaraja(t *testing.T) {
	for _, valor := range []int{0, 4, 22} {
		r := rondaAbierta(t)
		votar(t, &r, 1, 5)

		err := r.Votar(2, valor)

		if !errors.Is(err, ErrVotoInvalido) {
			t.Errorf("valor %d: se esperaba ErrVotoInvalido, se obtuvo %v", valor, err)
		}
		if !slices.Equal(r.Votos, []Voto{{UsuarioID: 1, Valor: 5}}) {
			t.Errorf("valor %d: los votos no debían cambiar, quedaron %+v", valor, r.Votos)
		}
	}
}

func TestRondaVotacion_Votar_RondaRevelada(t *testing.T) {
	for _, valor := range []int{13, 4} {
		r := rondaRevelada(t)

		err := r.Votar(1, valor)

		if !errors.Is(err, ErrRondaRevelada) {
			t.Errorf("valor %d: se esperaba ErrRondaRevelada, se obtuvo %v", valor, err)
		}
		if !slices.Equal(r.Votos, []Voto{{UsuarioID: 1, Valor: 8}}) {
			t.Errorf("valor %d: los votos no debían cambiar, quedaron %+v", valor, r.Votos)
		}
	}
}

func TestRondaVotacion_Revelar_ConTodosLosVotos(t *testing.T) {
	r := rondaAbierta(t)
	votar(t, &r, 1, 3)
	votar(t, &r, 2, 8)

	if err := r.Revelar([]int64{1, 2}, momentoRevelacion); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if r.ReveladaEn == nil || !r.ReveladaEn.Equal(momentoRevelacion) {
		t.Errorf("se esperaba revelada en %v, se obtuvo %v", momentoRevelacion, r.ReveladaEn)
	}
	if !r.Revelada() {
		t.Error("la ronda debía quedar revelada")
	}
}

func TestRondaVotacion_Revelar_FaltaUnVoto(t *testing.T) {
	r := rondaAbierta(t)
	votar(t, &r, 1, 3)

	err := r.Revelar([]int64{1, 2}, momentoRevelacion)

	if !errors.Is(err, ErrFaltanVotos) {
		t.Errorf("se esperaba ErrFaltanVotos, se obtuvo %v", err)
	}
	if r.ReveladaEn != nil || r.Revelada() {
		t.Errorf("la ronda no debía quedar revelada: %+v", r)
	}
}

func TestRondaVotacion_Revelar_YaRevelada(t *testing.T) {
	r := rondaRevelada(t)

	err := r.Revelar([]int64{1}, momentoRevelacion.Add(time.Hour))

	if !errors.Is(err, ErrRondaRevelada) {
		t.Errorf("se esperaba ErrRondaRevelada, se obtuvo %v", err)
	}
	if r.ReveladaEn == nil || !r.ReveladaEn.Equal(momentoRevelacion) {
		t.Errorf("no debía cambiar el instante de revelación, quedó %v", r.ReveladaEn)
	}
}
