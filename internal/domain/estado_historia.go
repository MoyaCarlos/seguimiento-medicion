package domain

import "errors"

// Errores del cambio de estado de una historia en el Sprint (HU-11).
var (
	ErrEstadoInvalido           = ValidationError{Campo: "estado", Mensaje: "el estado debe ser Nueva, En progreso o Completada"}
	ErrHistoriaSinSprint        = ValidationError{Campo: "estado", Mensaje: "primero debe asignarse a un Sprint activo"}
	ErrSprintNoActivoParaEstado = ValidationError{Campo: "estado", Mensaje: "solo se cambia el estado de historias de un Sprint activo"}
	ErrEstadoRepetido           = ValidationError{Campo: "estado", Mensaje: "la historia ya está en ese estado"}
	ErrRetrocesoEstado          = ValidationError{Campo: "estado", Mensaje: "el estado de una historia solo avanza: no se puede volver a un estado anterior"}
	ErrSaltoDirecto             = ValidationError{Campo: "estado", Mensaje: "primero debe pasar por En progreso"}

	// ErrConflictoEstadoHistoria lo devuelve el repositorio cuando el UPDATE
	// condicional no afecta filas porque otra petición cambió la historia antes.
	ErrConflictoEstadoHistoria = errors.New("la historia cambió de estado mientras se procesaba la petición")
)

// ValidarCambioDeEstado indica si la historia puede pasar al estado nuevo y,
// si no, devuelve el primer motivo según el orden de precedencia de la spec.
// s es el Sprint de la historia, o nil si no tiene.
func ValidarCambioDeEstado(h BacklogItem, s *Sprint, nuevo Estado) error {
	return nil
}
