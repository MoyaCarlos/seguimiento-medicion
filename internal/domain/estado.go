package domain

// Estado representa el estado de una historia de usuario dentro del backlog.
type Estado string

const (
	// EstadoNueva es el estado inicial con el que se registra toda historia.
	EstadoNueva Estado = "Nueva"
	// EstadoCompletada lo necesita el cierre de Sprint (HU-05) para decidir
	// qué historias vuelven al backlog; la transición a este estado es de HU-11.
	EstadoCompletada Estado = "Completada"
)
