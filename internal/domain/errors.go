package domain

// ValidationError describe una violación de una regla de negocio del dominio.
type ValidationError struct {
	Campo   string
	Mensaje string
}

func (e ValidationError) Error() string {
	return e.Campo + ": " + e.Mensaje
}
