package domain

import "time"

// Severidad indica cuán grave es un defecto.
type Severidad string

const (
	SeveridadAlta  Severidad = "Alta"
	SeveridadMedia Severidad = "Media"
	SeveridadBaja  Severidad = "Baja"
)

// EstadoDefecto es la etapa de un defecto: Pendiente hasta que se resuelve.
// Resuelto es definitivo (HU-08, Clarification 1).
type EstadoDefecto string

const (
	DefectoPendiente EstadoDefecto = "Pendiente"
	DefectoResuelto  EstadoDefecto = "Resuelto"
)

// Defecto es un error técnico detectado en el software, relacionado con una
// historia del proyecto y atribuido al Sprint Activo en que se registró. No
// tiene título (D6).
type Defecto struct {
	ID                int64
	ProyectoID        int64
	HistoriaID        int64
	SprintDeteccionID int64
	Descripcion       string
	Severidad         Severidad
	Estado            EstadoDefecto
	// SprintResolucionID es nil mientras el defecto está Pendiente.
	SprintResolucionID *int64
	CreadoEn           time.Time
}

// DefectoDetalle es la vista de un defecto con el título de su historia.
type DefectoDetalle struct {
	Defecto
	TituloHistoria string
}

// NewDefecto construye un defecto Pendiente detectado en el Sprint activo.
// Valida descripción y severidad, en ese orden. Supone que el service ya
// verificó el Sprint activo y la historia relacionada.
func NewDefecto(proyectoID, historiaID int64, activo Sprint, descripcion, severidad string) (Defecto, error) {
	return Defecto{}, nil
}
