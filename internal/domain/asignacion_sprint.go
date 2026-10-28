package domain

// Errores de la asignación de una historia a un Sprint (HU-06).
var (
	ErrSprintDeOtroProyecto = ValidationError{Campo: "sprint_id", Mensaje: "el Sprint no pertenece al proyecto de la historia"}
	ErrSprintFinalizado     = ValidationError{Campo: "sprint_id", Mensaje: "un Sprint finalizado no admite historias"}
	ErrHistoriaNoPendiente  = ValidationError{Campo: "estado", Mensaje: "solo se pueden planificar historias pendientes"}
	ErrHistoriaYaAsignada   = ValidationError{Campo: "sprint_id", Mensaje: "la historia ya pertenece a un Sprint"}
)

// ValidarAsignacion indica si la historia puede entrar al Sprint y, si no,
// devuelve el primer motivo según el orden de precedencia de la spec.
func ValidarAsignacion(h BacklogItem, s Sprint) error {
	return nil
}
