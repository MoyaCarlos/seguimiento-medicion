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

// Errores del registro y la resolución de defectos (HU-08).
var (
	// ErrDefectoNoEncontrado indica que el defecto solicitado no existe (HU-08).
	ErrDefectoNoEncontrado = errors.New("defecto no encontrado")
	// ErrSinSprintActivoParaDetectar indica un registro en un proyecto sin Sprint Activo (HU-08).
	ErrSinSprintActivoParaDetectar = ValidationError{Campo: "sprint_deteccion_id", Mensaje: "el proyecto no tiene un Sprint activo al que asociar el defecto como Sprint de detección"}
	// ErrSinSprintActivoParaResolver indica una resolución en un proyecto sin Sprint Activo (HU-08).
	ErrSinSprintActivoParaResolver = ValidationError{Campo: "sprint_resolucion_id", Mensaje: "el proyecto no tiene un Sprint activo al que asociar el defecto como Sprint de resolución"}
	// ErrDefectoYaResuelto indica una resolución de un defecto que ya está Resuelto (HU-08).
	ErrDefectoYaResuelto = ValidationError{Campo: "estado", Mensaje: "el defecto ya está resuelto"}
	// ErrHistoriaRelacionadaObligatoria indica un defecto sin historia relacionada (HU-08).
	ErrHistoriaRelacionadaObligatoria = ValidationError{Campo: "historia_id", Mensaje: "la historia relacionada es obligatoria"}
	// ErrHistoriaRelacionadaAjena indica una historia relacionada inexistente o de otro proyecto (HU-08).
	ErrHistoriaRelacionadaAjena = ValidationError{Campo: "historia_id", Mensaje: "la historia relacionada no existe en el proyecto"}
)

// Errores de la votación de Planning Poker.
var (
	// ErrVotacionNoIniciada indica que la historia no tiene una ronda de votación.
	ErrVotacionNoIniciada = errors.New("la historia no tiene votación iniciada")
	// ErrVotacionYaIniciada indica que la historia ya tiene una ronda de votación.
	ErrVotacionYaIniciada = errors.New("la historia ya tiene una votación")
	// ErrNoEsIntegrante indica que quien actúa no integra el proyecto de la historia.
	ErrNoEsIntegrante = errors.New("el usuario no es integrante del proyecto")
	// ErrNoEsScrumMaster indica que la acción es exclusiva del Scrum Master del proyecto.
	ErrNoEsScrumMaster = errors.New("solo el Scrum Master del proyecto puede hacerlo")

	ErrVotoInvalido       = ValidationError{Campo: "valor", Mensaje: "el voto debe ser una carta de la baraja: 1, 2, 3, 5, 8, 13, 21"}
	ErrHistoriaCompletada = ValidationError{Campo: "estado", Mensaje: "no se puede votar una historia completada"}
	ErrRondaRevelada      = ValidationError{Campo: "estado", Mensaje: "la ronda ya fue revelada"}
	ErrFaltanVotos        = ValidationError{Campo: "votos", Mensaje: "faltan votar integrantes del proyecto"}
)
