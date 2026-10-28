package domain

import "time"

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
// enteros.
func MinutosDesdeHoras(horas string) (int, error) {
	return 0, nil
}

// NewRegistroHoras construye una carga de horas para una historia del Sprint
// dado. Supone que el service ya verificó el Sprint activo y el integrante.
func NewRegistroHoras(historia BacklogItem, sprint Sprint, usuarioID int64, horas, actividad, fecha string, ahora time.Time) (RegistroHoras, error) {
	return RegistroHoras{}, nil
}
