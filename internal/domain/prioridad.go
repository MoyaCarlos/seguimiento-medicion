package domain

// Prioridad representa la prioridad de una historia según la escala MoSCoW.
type Prioridad string

const (
	PrioridadMust   Prioridad = "M"
	PrioridadShould Prioridad = "S"
	PrioridadCould  Prioridad = "C"
	PrioridadWont   Prioridad = "W"
)

// EsValida indica si la prioridad pertenece a la escala MoSCoW.
func (p Prioridad) EsValida() bool {
	switch p {
	case PrioridadMust, PrioridadShould, PrioridadCould, PrioridadWont:
		return true
	default:
		return false
	}
}
