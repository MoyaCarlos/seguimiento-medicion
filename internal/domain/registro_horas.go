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
