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

func TestCalcularPromedio(t *testing.T) {
	casos := []struct {
		valores  []int
		esperado Promedio
	}{
		{[]int{3, 5, 8}, Promedio{Exacto: 5.33, CartaSugerida: 5}},
		{[]int{3, 5}, Promedio{Exacto: 4, CartaSugerida: 5}},
		{[]int{8, 8, 8}, Promedio{Exacto: 8, CartaSugerida: 8}},
		{[]int{1, 21}, Promedio{Exacto: 11, CartaSugerida: 13}},
		{[]int{2, 3}, Promedio{Exacto: 2.5, CartaSugerida: 3}},
	}

	for _, c := range casos {
		if got := CalcularPromedio(c.valores); got != c.esperado {
			t.Errorf("%v: se esperaba %+v, se obtuvo %+v", c.valores, c.esperado, got)
		}
	}
}

func TestRondaVotacion_Promedio_SinRevelar(t *testing.T) {
	r := rondaAbierta(t)
	votar(t, &r, 1, 8)

	if p, ok := r.Promedio(); ok {
		t.Errorf("una ronda sin revelar no debe tener promedio, se obtuvo %+v", p)
	}
}

func TestRondaVotacion_Promedio_Revelada(t *testing.T) {
	r := rondaAbierta(t)
	votar(t, &r, 1, 3)
	votar(t, &r, 2, 5)
	votar(t, &r, 3, 8)
	if err := r.Revelar([]int64{1, 2, 3}, momentoRevelacion); err != nil {
		t.Fatalf("no se esperaba error al revelar: %v", err)
	}

	p, ok := r.Promedio()

	if !ok || p != CalcularPromedio([]int{3, 5, 8}) {
		t.Errorf("se esperaba el promedio de 3, 5 y 8, se obtuvo %+v (%v)", p, ok)
	}
}

func TestCoincidenciaExacta_HayDiferencias(t *testing.T) {
	casos := []struct {
		nombre   string
		valores  []int
		esperado bool
	}{
		{"cartas distintas", []int{3, 5, 8}, true},
		{"cartas vecinas", []int{3, 5}, true},
		{"todas iguales", []int{5, 5, 5}, false},
		{"un solo voto", []int{8}, false},
		{"sin votos", []int{}, false},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := (CoincidenciaExacta{}).HayDiferencias(c.valores); got != c.esperado {
				t.Errorf("HayDiferencias(%v) = %v, se esperaba %v", c.valores, got, c.esperado)
			}
		})
	}
}

// siempreDiferencias prueba que la regla de consenso es intercambiable.
type siempreDiferencias struct{}

func (siempreDiferencias) HayDiferencias([]int) bool { return true }

func rondaConVotos(revelada bool, valores ...int) RondaVotacion {
	r := RondaVotacion{ID: 1, HistoriaID: 100, Numero: 1}
	for i, v := range valores {
		r.Votos = append(r.Votos, Voto{UsuarioID: int64(i + 1), Valor: v})
	}
	if revelada {
		en := momentoRevelacion
		r.ReveladaEn = &en
	}
	return r
}

func TestRondaVotacion_Diferencias_SinRevelar(t *testing.T) {
	if _, ok := rondaConVotos(false, 3, 8).Diferencias(CoincidenciaExacta{}); ok {
		t.Error("una ronda sin revelar no debía informar diferencias")
	}
}

func TestRondaVotacion_Diferencias_Revelada(t *testing.T) {
	casos := []struct {
		nombre   string
		ronda    RondaVotacion
		regla    ReglaConsenso
		esperado Diferencias
	}{
		{"con diferencias", rondaConVotos(true, 3, 5, 8), CoincidenciaExacta{}, Diferencias{Hay: true, Minimo: 3, Maximo: 8}},
		{"sin importar el orden", rondaConVotos(true, 8, 3, 5), CoincidenciaExacta{}, Diferencias{Hay: true, Minimo: 3, Maximo: 8}},
		{"sin diferencias", rondaConVotos(true, 5, 5, 5), CoincidenciaExacta{}, Diferencias{Hay: false}},
		{"otra regla cambia el resultado", rondaConVotos(true, 5, 5, 5), siempreDiferencias{}, Diferencias{Hay: true, Minimo: 5, Maximo: 5}},
		{"sin votos no consulta la regla", rondaConVotos(true), siempreDiferencias{}, Diferencias{Hay: false}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got, ok := c.ronda.Diferencias(c.regla)
			if !ok {
				t.Fatal("una ronda revelada debía informar diferencias")
			}
			if got != c.esperado {
				t.Errorf("se esperaba %+v, se obtuvo %+v", c.esperado, got)
			}
		})
	}
}
