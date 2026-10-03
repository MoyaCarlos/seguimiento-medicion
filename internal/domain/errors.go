package domain

import "errors"

// ValidationError describe una violación de una regla de negocio del dominio.
type ValidationError struct {
	Campo   string
	Mensaje string
}

func (e ValidationError) Error() string {
	return e.Campo + ": " + e.Mensaje
}

// ErrNoEncontrado indica que la entidad solicitada no existe.
var ErrNoEncontrado = errors.New("registro no encontrado")

// ErrProyectoNoEncontrado indica que el proyecto solicitado no existe.
var ErrProyectoNoEncontrado = errors.New("proyecto no encontrado")
