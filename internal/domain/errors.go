package domain

import "errors"

// ErrSprintNoEncontrado indica que el Sprint pedido no existe.
var ErrSprintNoEncontrado = errors.New("sprint no encontrado")

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

// ErrHistoriaNoEncontrada indica que la historia solicitada no existe.
var ErrHistoriaNoEncontrada = errors.New("historia no encontrada")

// ErrEstimacionCongelada indica un cambio de los Story Points ya cargados de una
// historia de un Sprint Activo o Finalizado (decisión D2b).
var ErrEstimacionCongelada = ValidationError{Campo: "estimacion_sp", Mensaje: "los Story Points ya cargados no se pueden cambiar en un Sprint Activo o Finalizado"}

// ErrHistoriaFueraDelSprintActivo indica una carga de horas a una historia que
// no está en un Sprint Activo (HU-07).
var ErrHistoriaFueraDelSprintActivo = ValidationError{Campo: "historia_id", Mensaje: "solo se imputan horas a historias del Sprint activo"}

// ErrIntegranteAjeno indica una carga de horas de alguien que no es integrante
// del proyecto de la historia (HU-07).
var ErrIntegranteAjeno = ValidationError{Campo: "integrante_id", Mensaje: "el integrante no pertenece al proyecto de la historia"}
