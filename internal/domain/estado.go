package domain

// Estado representa el estado de una historia de usuario dentro del backlog.
type Estado string

const (
	// EstadoNueva es el estado inicial con el que se registra toda historia.
	EstadoNueva Estado = "Nueva"
	// EstadoEnProgreso indica que la historia está siendo trabajada en el Sprint.
	EstadoEnProgreso Estado = "En progreso"
	// EstadoCompletada lo necesita el cierre de Sprint (HU-05) para decidir
	// qué historias vuelven al backlog; la transición a este estado es de HU-11.
	EstadoCompletada Estado = "Completada"
)

// EsValida indica si el estado es uno de los que admite el modelo de datos.
func (e Estado) EsValida() bool {
	switch e {
	case EstadoNueva, EstadoEnProgreso, EstadoCompletada:
		return true
	default:
		return false
	}
}
