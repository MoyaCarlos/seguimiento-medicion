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

// ErrScrumMasterYaExiste indica que el proyecto ya tiene un Scrum Master y se
// quiso agregar otro (spec 039, #236).
var ErrScrumMasterYaExiste = errors.New("el proyecto ya tiene un Scrum Master")

// ErrNoEncontrado indica que la entidad solicitada no existe.
var ErrNoEncontrado = errors.New("registro no encontrado")

// ErrProyectoNoEncontrado indica que el proyecto solicitado no existe.
var ErrProyectoNoEncontrado = errors.New("proyecto no encontrado")

// ErrHistoriaNoEncontrada indica que la historia solicitada no existe.
var ErrHistoriaNoEncontrada = errors.New("historia no encontrada")

// ErrSprintNoFinalizado indica que se pidió en el historial un Sprint que
// todavía está Pendiente o Activo (HU-12).
var ErrSprintNoFinalizado = ValidationError{Campo: "estado", Mensaje: "el Sprint todavía no está finalizado"}

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
	// ErrRondaYaAbierta indica que otra petición abrió esa misma ronda mientras tanto.
	ErrRondaYaAbierta = errors.New("otra persona abrió esta ronda mientras tanto: recargá la votación")
	// ErrNoEsIntegrante indica que quien actúa no integra el proyecto de la historia.
	ErrNoEsIntegrante = errors.New("el usuario no es integrante del proyecto")
	// ErrNoEsScrumMaster indica que la acción es exclusiva del Scrum Master del proyecto.
	ErrNoEsScrumMaster = errors.New("solo el Scrum Master del proyecto puede hacerlo")

	ErrVotoInvalido       = ValidationError{Campo: "valor", Mensaje: "el voto debe ser una carta de la baraja: 1, 2, 3, 5, 8, 13, 21"}
	ErrHistoriaCompletada = ValidationError{Campo: "estado", Mensaje: "no se puede votar una historia completada"}
	ErrRondaRevelada      = ValidationError{Campo: "estado", Mensaje: "la ronda ya fue revelada"}
	ErrFaltanVotos        = ValidationError{Campo: "votos", Mensaje: "faltan votar integrantes del proyecto"}
	ErrRondaNoRevelada    = ValidationError{Campo: "estado", Mensaje: "la ronda actual todavía no fue revelada"}

	ErrRondaEnHistoriaCompletada = ValidationError{Campo: "estado", Mensaje: "no se puede abrir una ronda en una historia completada"}
)

// Errores de autenticación (HU-20).
// Errores de autenticación (HU-20).
var (
	// ErrCredencialesInvalidas indica un ingreso con usuario o contraseña que no coinciden.
	ErrCredencialesInvalidas = errors.New("Usuario o contraseña incorrectos")
	// ErrNoAutenticado indica un pedido sin sesión vigente.
	ErrNoAutenticado = errors.New("Tu sesión no es válida")
	// ErrNombreUsuarioEnUso indica un registro con un nombre de usuario ya existente.
	ErrNombreUsuarioEnUso = errors.New("Ese nombre de usuario ya existe")
)

// ErrHorasPorStoryPointInvalidas indica un parámetro que no es un número mayor que cero (HU-03).
var ErrHorasPorStoryPointInvalidas = ValidationError{Campo: "horas_por_story_point", Mensaje: "las horas por Story Point deben ser un número mayor que cero"}

// ErrCargaEnSprintFinalizado indica una carga de Story Points en una historia
// sin estimar de un Sprint Finalizado (decisión D2b, reabierta en HU-03).
var ErrCargaEnSprintFinalizado = ValidationError{Campo: "estimacion_sp", Mensaje: "no se pueden cargar Story Points en un Sprint Finalizado"}
