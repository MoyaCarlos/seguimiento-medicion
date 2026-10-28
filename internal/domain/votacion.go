package domain

import "time"

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
	return RondaVotacion{}, nil
}

// Votar registra el voto del usuario o reemplaza el que ya tenía.
func (r *RondaVotacion) Votar(usuarioID int64, valor int) error {
	return nil
}

// Revelar muestra los votos cuando votaron todos los integrantes.
func (r *RondaVotacion) Revelar(integrantes []int64, ahora time.Time) error {
	return nil
}

// Revelada indica si los votos de la ronda ya son visibles.
func (r RondaVotacion) Revelada() bool {
	return false
}
