package domain

// Estado representa el estado de una historia de usuario dentro del backlog.
type Estado string

const (
	// EstadoNueva es el estado inicial con el que se registra toda historia.
	EstadoNueva Estado = "Nueva"
)
