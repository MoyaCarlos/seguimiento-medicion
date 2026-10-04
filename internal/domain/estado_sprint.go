package domain

// EstadoSprint representa la etapa del ciclo de vida de un Sprint.
type EstadoSprint string

const (
	SprintPendiente  EstadoSprint = "Pendiente"
	SprintActivo     EstadoSprint = "Activo"
	SprintFinalizado EstadoSprint = "Finalizado"
)
