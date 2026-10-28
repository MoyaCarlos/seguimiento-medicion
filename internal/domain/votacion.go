package domain

import (
	"slices"
	"time"
)

// RondaVotacion es una ronda de Planning Poker sobre una historia. Mientras
// ReveladaEn es nil los votos están ocultos.
type RondaVotacion struct {
	ID         int64
	HistoriaID int64
	Numero     int
	ReveladaEn *time.Time
	Votos      []Voto
}

// Voto es la carta que eligió un integrante en una ronda.
type Voto struct {
	UsuarioID int64
	Valor     int
}

// NuevaRonda crea la ronda 1 de una historia, sin votos ni revelar.
func NuevaRonda(historiaID int64) (RondaVotacion, error) {
	if historiaID <= 0 {
		return RondaVotacion{}, ValidationError{Campo: "historia_id", Mensaje: "la historia es obligatoria"}
	}
	return RondaVotacion{HistoriaID: historiaID, Numero: 1}, nil
}

// Votar registra el voto del usuario o reemplaza el que ya tenía.
func (r *RondaVotacion) Votar(usuarioID int64, valor int) error {
	if r.Revelada() {
		return ErrRondaRevelada
	}
	if !EsEstimacionValida(valor) {
		return ErrVotoInvalido
	}
	if i := r.indiceDelVoto(usuarioID); i >= 0 {
		r.Votos[i].Valor = valor
		return nil
	}
	r.Votos = append(r.Votos, Voto{UsuarioID: usuarioID, Valor: valor})
	return nil
}

// Revelar muestra los votos cuando votaron todos los integrantes.
func (r *RondaVotacion) Revelar(integrantes []int64, ahora time.Time) error {
	if r.Revelada() {
		return ErrRondaRevelada
	}
	for _, id := range integrantes {
		if r.indiceDelVoto(id) < 0 {
			return ErrFaltanVotos
		}
	}
	r.ReveladaEn = &ahora
	return nil
}

// Revelada indica si los votos de la ronda ya son visibles.
func (r RondaVotacion) Revelada() bool {
	return r.ReveladaEn != nil
}

func (r RondaVotacion) indiceDelVoto(usuarioID int64) int {
	return slices.IndexFunc(r.Votos, func(v Voto) bool { return v.UsuarioID == usuarioID })
}
