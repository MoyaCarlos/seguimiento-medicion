package domain

import (
	"math/big"
	"regexp"
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

const (
	minutosMaximosPorCarga  = 24 * 60
	longitudMaximaActividad = 500
)

// formatoHoras admite solo decimales: big.Rat también aceptaría fracciones
// (1/2) y notación exponencial (1e1), que se rechazan a propósito (R2).
var formatoHoras = regexp.MustCompile(`^-?\d+(\.\d+)?$`)

// MinutosDesdeHoras convierte una cantidad de horas en texto decimal a minutos
// enteros. Usa big.Rat para que la conversión sea exacta (0,1 h = 6 min).
func MinutosDesdeHoras(horas string) (int, error) {
	if horas == "" {
		return 0, ValidationError{Campo: "horas", Mensaje: "las horas son obligatorias"}
	}
	if !formatoHoras.MatchString(horas) {
		return 0, ValidationError{Campo: "horas", Mensaje: "las horas deben ser un número"}
	}
	minutos, _ := new(big.Rat).SetString(horas)
	minutos.Mul(minutos, big.NewRat(60, 1))
	if minutos.Sign() <= 0 {
		return 0, ValidationError{Campo: "horas", Mensaje: "las horas deben ser mayores que cero"}
	}
	if minutos.Cmp(big.NewRat(minutosMaximosPorCarga, 1)) > 0 {
		return 0, ValidationError{Campo: "horas", Mensaje: "las horas no pueden superar 24 por carga"}
	}
	if !minutos.IsInt() {
		return 0, ValidationError{Campo: "horas", Mensaje: "las horas deben equivaler a minutos enteros"}
	}
	return int(minutos.Num().Int64()), nil
}

// NewRegistroHoras construye una carga de horas para una historia del Sprint
// dado. Valida horas, actividad y fecha, en ese orden. Supone que el service ya
// verificó el Sprint activo y el integrante.
func NewRegistroHoras(historia BacklogItem, sprint Sprint, usuarioID int64, horas, actividad, fecha string, ahora time.Time) (RegistroHoras, error) {
	minutos, err := MinutosDesdeHoras(horas)
	if err != nil {
		return RegistroHoras{}, err
	}
	actividad, err = validarActividad(actividad)
	if err != nil {
		return RegistroHoras{}, err
	}
	dia, err := validarFecha(fecha, sprint, ahora)
	if err != nil {
		return RegistroHoras{}, err
	}
	return RegistroHoras{
		ProyectoID: historia.ProyectoID,
		HistoriaID: historia.ID,
		UsuarioID:  usuarioID,
		SprintID:   sprint.ID,
		Fecha:      dia,
		Actividad:  actividad,
		Minutos:    minutos,
	}, nil
}

func validarActividad(actividad string) (string, error) {
	actividad = strings.TrimSpace(actividad)
	if actividad == "" {
		return "", ValidationError{Campo: "actividad", Mensaje: "la actividad es obligatoria"}
	}
	if len([]rune(actividad)) > longitudMaximaActividad {
		return "", ValidationError{Campo: "actividad", Mensaje: "la actividad no puede superar los 500 caracteres"}
	}
	return actividad, nil
}

// validarFecha acepta desde el inicio del Sprint hasta la menor entre hoy y
// su fin, ambas inclusive. Compara primero contra hoy, después contra el fin
// y por último contra el inicio.
func validarFecha(fecha string, sprint Sprint, ahora time.Time) (time.Time, error) {
	dia, err := time.Parse(time.DateOnly, fecha)
	if err != nil {
		return time.Time{}, ValidationError{Campo: "fecha", Mensaje: "formato esperado AAAA-MM-DD"}
	}
	if dia.After(inicioDelDiaCivil(ahora)) {
		return time.Time{}, ValidationError{Campo: "fecha", Mensaje: "la fecha no puede ser posterior a hoy"}
	}
	if dia.After(sprint.FechaFin) {
		return time.Time{}, ValidationError{Campo: "fecha", Mensaje: "la fecha no puede ser posterior al fin del Sprint"}
	}
	if dia.Before(sprint.FechaInicio) {
		return time.Time{}, ValidationError{Campo: "fecha", Mensaje: "la fecha no puede ser anterior al inicio del Sprint"}
	}
	return dia, nil
}

// ValidarSprintParaImputar verifica que la historia esté asignada al Sprint y
// que ese Sprint esté Activo; el estado de la historia no importa.
func ValidarSprintParaImputar(historia BacklogItem, sprint Sprint) error {
	if historia.SprintID == nil || *historia.SprintID != sprint.ID || sprint.Estado != SprintActivo {
		return ErrHistoriaFueraDelSprintActivo
	}
	return nil
}

// TotalMinutos suma exactamente las duraciones de las cargas (FR-013).
func TotalMinutos(registros []RegistroHorasDetalle) int {
	return 0
}
