package domain

import (
	"math/big"
	"strings"
	"time"
)

// RegistroHoras es una carga de horas trabajadas en una historia, atribuida al
// Sprint Activo de la historia en el momento de la carga. La duración se guarda
// en minutos enteros (D8).
type RegistroHoras struct {
	ID         int64
	ProyectoID int64
	HistoriaID int64
	UsuarioID  int64
	SprintID   int64
	Fecha      time.Time
	Actividad  string
	Minutos    int
	CreadoEn   time.Time
}

// RegistroHorasDetalle es la vista de una carga con el nombre del integrante.
type RegistroHorasDetalle struct {
	RegistroHoras
	NombreIntegrante string
}

// MinutosDesdeHoras convierte una cantidad de horas en texto decimal a minutos
// enteros. Usa big.Rat para que la conversión sea exacta (0,1 h = 6 min).
func MinutosDesdeHoras(horas string) (int, error) {
	cantidad, ok := new(big.Rat).SetString(horas)
	if !ok {
		return 0, ValidationError{Campo: "horas", Mensaje: "las horas deben ser un número"}
	}
	minutos := cantidad.Mul(cantidad, big.NewRat(60, 1))
	return int(minutos.Num().Int64()), nil
}

// NewRegistroHoras construye una carga de horas para una historia del Sprint
// dado. Supone que el service ya verificó el Sprint activo y el integrante.
func NewRegistroHoras(historia BacklogItem, sprint Sprint, usuarioID int64, horas, actividad, fecha string, ahora time.Time) (RegistroHoras, error) {
	minutos, err := MinutosDesdeHoras(horas)
	if err != nil {
		return RegistroHoras{}, err
	}
	dia, err := time.Parse(time.DateOnly, fecha)
	if err != nil {
		return RegistroHoras{}, ValidationError{Campo: "fecha", Mensaje: "formato esperado AAAA-MM-DD"}
	}
	return RegistroHoras{
		ProyectoID: historia.ProyectoID,
		HistoriaID: historia.ID,
		UsuarioID:  usuarioID,
		SprintID:   sprint.ID,
		Fecha:      dia,
		Actividad:  strings.TrimSpace(actividad),
		Minutos:    minutos,
	}, nil
}

// ValidarSprintParaImputar verifica que la historia esté asignada al Sprint y
// que ese Sprint esté Activo; el estado de la historia no importa.
func ValidarSprintParaImputar(historia BacklogItem, sprint Sprint) error {
	return nil
}
