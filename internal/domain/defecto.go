package domain

import "time"

// Severidad indica cuán grave es un defecto.
type Severidad string

const (
	SeveridadAlta  Severidad = "Alta"
	SeveridadMedia Severidad = "Media"
	SeveridadBaja  Severidad = "Baja"
)

// EsValida indica si la severidad es exactamente Alta, Media o Baja.
func (s Severidad) EsValida() bool {
	switch s {
	case SeveridadAlta, SeveridadMedia, SeveridadBaja:
		return true
	default:
		return false
	}
}

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
	descripcion, err := validarDescripcion(descripcion)
	if err != nil {
		return Defecto{}, err
	}
	if !Severidad(severidad).EsValida() {
		return Defecto{}, ValidationError{Campo: "severidad", Mensaje: "la severidad debe ser Alta, Media o Baja"}
	}
	return Defecto{
		ProyectoID:        proyectoID,
		HistoriaID:        historiaID,
		SprintDeteccionID: activo.ID,
		Descripcion:       descripcion,
		Severidad:         Severidad(severidad),
		Estado:            DefectoPendiente,
	}, nil
}

// Resolver pasa el defecto de Pendiente a Resuelto con el Sprint activo como
// Sprint de resolución. "Ya resuelto" va antes que la falta de Sprint activo
// (research.md R3); ante un error el defecto no cambia.
func (d *Defecto) Resolver(activo *Sprint) error {
	if d.Estado == DefectoResuelto {
		return ErrDefectoYaResuelto
	}
	if activo == nil {
		return ErrSinSprintActivoParaResolver
	}
	sprintID := activo.ID
	d.Estado, d.SprintResolucionID = DefectoResuelto, &sprintID
	return nil
}
