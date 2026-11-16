package demo

import "github.com/MoyaCarlos/seguimiento-medicion/internal/domain"

func dia(n int) *int { return &n }

var (
	ana     = miembroDemo{usuario: "ana.iriarte", rol: domain.RolScrumMaster}
	carlos  = miembroDemo{usuario: "carlos.moya", rol: domain.RolProductBuilder}
	candela = miembroDemo{usuario: "candela.vulcano", rol: domain.RolProductBuilder}
	jimena  = miembroDemo{usuario: "jimena.martinez", rol: domain.RolProductBuilder}
)

var equipoCompleto = []miembroDemo{ana, carlos, candela, jimena}

func todos(historia string, sprint int, usuario string, d int, horas, actividad string) horasDemo {
	return horasDemo{historia, sprint, usuario, d, horas, actividad}
}

var (
	comp = domain.EstadoCompletada
	nva  = domain.EstadoNueva
	prog = estadoEnProgreso
)

// proyectosAdicionales: spec 052. El orden es el de la tabla de la spec.
var proyectosAdicionales = []proyectoDemo{
	plataformaTurnos,
	bibliotecaCerrada,
	tiendaFinalizada,
	appReservas,
	portalAlumnos,
	proyectoEnBlanco,
	campanaDonaciones,
}

var plataformaTurnos = proyectoDemo{
	nombre:      "Plataforma de turnos médicos (datos de prueba)",
	descripcion: "Gestión de turnos para una red de consultorios: agenda de médicos, reserva y cancelación, recordatorios e historial clínico. Proyecto de demostración con datos de ejemplo.",
	inicio:      -45, fin: 30, horasPorSP: 4,
	equipo: equipoCompleto,
	sprints: []sprintDemo{
		{"Registro de pacientes y agenda de médicos", -42, -29, domain.SprintFinalizado},
		{"Reserva y cancelación de turnos", -28, -15, domain.SprintFinalizado},
		{"Recordatorios e historial clínico", -7, 6, domain.SprintActivo},
		{estado: domain.SprintPendiente},
	},
	arrastreDesde: 2,
	historias: []historiaDemo{
		{clave: "registro", titulo: "Registrar pacientes con su obra social",
			descripcion: "Como **recepcionista** quiero cargar los datos del paciente para:\n\n- identificarlo en cada turno\n- facturar a su obra social",
			criterios:   "- El DNI no admite puntos\n- La obra social es opcional",
			prioridad:   domain.PrioridadMust, valor: 8, puntos: 5, sprint: 1, estado: comp, finalizada: -33},
		{clave: "medicos", titulo: "Dar de alta médicos y especialidades",
			descripcion: "Como administrador quiero registrar médicos con su especialidad para ofrecerlos en la agenda.",
			prioridad:   domain.PrioridadMust, valor: 8, puntos: 3, sprint: 1, estado: comp, finalizada: -32},
		{clave: "login", titulo: "Ingresar con usuario y contraseña",
			descripcion: "Como paciente quiero ingresar con mi usuario para ver mis turnos.",
			criterios:   "- Tras 5 intentos fallidos se bloquea 15 minutos",
			prioridad:   domain.PrioridadMust, valor: 13, puntos: 3, sprint: 1, estado: comp, finalizada: -31},
		{clave: "agenda", titulo: "Definir la agenda semanal de cada médico",
			descripcion: "## Contexto\n\nCada médico atiende en *franjas* fijas por día de la semana.\n\n1. Elegir días y horarios\n2. Definir la duración del turno\n3. Publicar la agenda",
			criterios:   "- Las franjas no pueden superponerse\n- La duración mínima del turno es de 10 minutos",
			prioridad:   domain.PrioridadMust, valor: 13, puntos: 8, sprint: 1, estado: comp, finalizada: -30},
		{clave: "consultorios", titulo: "Administrar los consultorios",
			descripcion: "Como administrador quiero asignar consultorios a los médicos para evitar choques de ocupación.",
			prioridad:   domain.PrioridadShould, valor: 5, puntos: 2, sprint: 1, estado: comp, finalizada: -30},
		{clave: "buscar", titulo: "Buscar médicos por especialidad",
			descripcion: "Como paciente quiero filtrar médicos por especialidad para elegir con quién atenderme.",
			prioridad:   domain.PrioridadShould, valor: 8, puntos: 3, sprint: 1, estado: comp, finalizada: -29},

		{clave: "reservar", titulo: "Reservar un turno con un médico",
			descripcion: "Como paciente quiero elegir día y hora disponibles. Ver [la guía de turnos](https://ejemplo.org/guia-turnos).\n\n> Un turno no puede superponerse con otro del mismo médico.",
			criterios:   "- Solo se ofrecen horarios libres\n- Al confirmar se bloquea el horario",
			prioridad:   domain.PrioridadMust, valor: 21, puntos: 8, sprint: 2, estado: comp, finalizada: -22},
		{clave: "cancelar", titulo: "Cancelar un turno reservado",
			descripcion: "Como paciente quiero cancelar mi turno para liberar el horario.",
			criterios:   "- Se puede cancelar hasta 2 horas antes\n- El horario vuelve a quedar libre",
			prioridad:   domain.PrioridadMust, valor: 13, puntos: 5, sprint: 2, estado: comp, finalizada: -20},
		{clave: "confirmacion", titulo: "Enviar el comprobante del turno por correo",
			descripcion: "Como paciente quiero recibir un correo con los datos del turno para no olvidarlo.",
			prioridad:   domain.PrioridadShould, valor: 8, puntos: 3, sprint: 2, estado: comp, finalizada: -19},
		{clave: "sobreturnos", titulo: "Permitir sobreturnos con autorización del médico",
			descripcion: "Como recepcionista quiero agendar un sobreturno cuando el médico lo autoriza.",
			prioridad:   domain.PrioridadCould, valor: 3, puntos: 5, sprint: 2, estado: comp, finalizada: -17},
		{clave: "feriados", titulo: "Bloquear feriados y licencias de los médicos",
			descripcion: "Como administrador quiero cargar feriados y licencias para no ofrecer turnos esos días.",
			prioridad:   domain.PrioridadShould, valor: 8, puntos: 3, sprint: 2, estado: comp, finalizada: -16},
		{clave: "reprogramar", titulo: "Reprogramar un turno a otra fecha",
			descripcion: "Como paciente quiero mover mi turno a otro horario sin cancelarlo.",
			prioridad:   domain.PrioridadShould, valor: 8, puntos: 5, sprint: 2, estado: comp, finalizada: -15},

		{clave: "recordatorios", titulo: "Enviar recordatorios de turno por mensaje",
			descripcion: "Como paciente quiero recibir un recordatorio 24 horas antes para no faltar.",
			criterios:   "- Se envía un solo recordatorio por turno\n- Se puede desactivar desde el perfil",
			prioridad:   domain.PrioridadMust, valor: 13, puntos: 5, sprint: 3, estado: comp, finalizada: -3},
		{clave: "pagos", titulo: "Cobrar el coseguro en línea",
			descripcion: "Como paciente quiero pagar el coseguro al reservar para no hacerlo en el consultorio.",
			prioridad:   domain.PrioridadMust, valor: 13, puntos: 8, sprint: 3, estado: comp, finalizada: -2},
		{clave: "historial", titulo: "Historial clínico del paciente con consultas, estudios y recetas ordenados por fecha de atención",
			descripcion: "Guardar cada consulta con `fecha`, `motivo` y `diagnostico`.\n\n```\nGET /pacientes/{id}/historial\n```\n\nSolo lo ve el médico tratante.",
			criterios:   "- Las consultas se ordenan de la más nueva a la más vieja\n- El paciente puede ver su propio historial",
			prioridad:   domain.PrioridadMust, valor: 21, puntos: 13, sprint: 3, estado: prog},
		{clave: "ausentismo", titulo: "Registrar el ausentismo del paciente",
			descripcion: "Como médico quiero marcar si el paciente faltó para medir el ausentismo.",
			prioridad:   domain.PrioridadShould, valor: 5, puntos: 3, sprint: 3, estado: prog},
		{clave: "notificaciones", titulo: "Avisar al paciente los cambios de turno",
			descripcion: "Como paciente quiero que me avisen si el médico mueve mi turno.",
			prioridad:   domain.PrioridadShould, valor: 8, sprint: 3, estado: nva},
		{clave: "recetas", titulo: "Adjuntar recetas al historial",
			descripcion: "Como médico quiero adjuntar la receta a la consulta para que el paciente la descargue.",
			prioridad:   domain.PrioridadCould, valor: 5, sprint: 3, estado: nva},

		{clave: "estadisticas", titulo: "Estadísticas de ocupación por consultorio",
			descripcion: "Como administrador quiero ver qué porcentaje de la agenda se ocupa para redistribuir consultorios.",
			prioridad:   domain.PrioridadShould, valor: 8, puntos: 5, sprint: 4, estado: nva},
		{clave: "teleconsulta", titulo: "Teleconsulta por video",
			descripcion: "Como paciente quiero atenderme por video cuando no puedo viajar.",
			prioridad:   domain.PrioridadCould, valor: 8, sprint: 4, estado: nva},
		{clave: "exportar", titulo: "Exportar la agenda a un calendario externo",
			descripcion: "Como médico quiero exportar mi agenda para verla en mi celular.",
			prioridad:   domain.PrioridadCould, valor: 3, sprint: 4, estado: nva},

		{clave: "listaespera", titulo: "Lista de espera para turnos cancelados",
			descripcion: "Como paciente quiero anotarme en una lista de espera para tomar un turno que se libere.",
			prioridad:   domain.PrioridadCould, valor: 5, puntos: 8, estado: nva, arrastrada: true},
		{clave: "app", titulo: "Aplicación móvil para pacientes",
			descripcion: "Como paciente quiero gestionar mis turnos desde el celular.",
			prioridad:   domain.PrioridadShould, valor: 13, estado: nva},
		{clave: "informes", titulo: "Informe mensual para obras sociales",
			descripcion: "Como administrador quiero un informe de prestaciones por obra social.",
			prioridad:   domain.PrioridadCould, valor: 5, estado: nva},
		{clave: "mapa", titulo: "Mostrar el mapa del consultorio",
			descripcion: "Como paciente quiero ver cómo llegar al consultorio.",
			prioridad:   domain.PrioridadWont, estado: nva},
	},
	horas: []horasDemo{
		todos("registro", 1, "ana.iriarte", -41, "4", "Formulario de pacientes"),
		todos("registro", 1, "candela.vulcano", -40, "1:30", "Validación del DNI"),
		todos("medicos", 1, "carlos.moya", -39, "3", "Alta de médicos"),
		todos("login", 1, "jimena.martinez", -38, "2.5", "Pantalla de ingreso"),
		todos("agenda", 1, "ana.iriarte", -36, "6", "Modelo de franjas semanales"),
		todos("agenda", 1, "candela.vulcano", -35, "8", "Pantalla de la agenda"),
		todos("agenda", 1, "carlos.moya", -33, "2:15", "Revisión con el cliente"),
		todos("consultorios", 1, "jimena.martinez", -32, "0:45", "ABM de consultorios"),
		todos("buscar", 1, "candela.vulcano", -30, "3", "Filtro por especialidad"),
		todos("reservar", 2, "ana.iriarte", -26, "7", "Reglas de disponibilidad"),
		todos("reservar", 2, "candela.vulcano", -24, "9", "Pantalla de reserva"),
		todos("reservar", 2, "jimena.martinez", -22, "1:30", "Pruebas de aceptación"),
		todos("cancelar", 2, "carlos.moya", -21, "4", "Regla de cancelación"),
		todos("cancelar", 2, "ana.iriarte", -20, "2.5", "Liberar el horario"),
		todos("confirmacion", 2, "candela.vulcano", -19, "3", "Plantilla del correo"),
		todos("sobreturnos", 2, "jimena.martinez", -18, "5", "Autorización del médico"),
		todos("feriados", 2, "carlos.moya", -17, "2", "Carga de feriados"),
		todos("reprogramar", 2, "ana.iriarte", -16, "4:30", "Reprogramación de turnos"),
		todos("recordatorios", 3, "candela.vulcano", -6, "5", "Envío de recordatorios"),
		todos("recordatorios", 3, "carlos.moya", -4, "1:30", "Pruebas del recordatorio"),
		todos("pagos", 3, "ana.iriarte", -5, "6", "Integración con la pasarela"),
		todos("pagos", 3, "jimena.martinez", -3, "3", "Pantalla de pago"),
		todos("historial", 3, "candela.vulcano", -2, "8", "Modelo del historial"),
		todos("historial", 3, "ana.iriarte", -1, "0:45", "Revisión del diseño"),
		todos("ausentismo", 3, "jimena.martinez", -1, "1.5", "Marca de ausencia"),
	},
	rondas: []rondaAdicional{
		{"estadisticas", 1, []votoDemo{{"carlos.moya", 5}, {"ana.iriarte", 5}, {"candela.vulcano", 5}, {"jimena.martinez", 5}}, true},
		{"teleconsulta", 1, []votoDemo{{"carlos.moya", 3}, {"ana.iriarte", 5}, {"candela.vulcano", 8}, {"jimena.martinez", 13}}, true},
		{"teleconsulta", 2, []votoDemo{{"ana.iriarte", 8}, {"jimena.martinez", 8}}, false},
		{"notificaciones", 1, []votoDemo{{"ana.iriarte", 8}, {"candela.vulcano", 5}, {"jimena.martinez", 8}}, false},
	},
	defectos: []defectoDemo{
		{"registro", "El DNI escrito con puntos no se reconoce al buscar al paciente", "Media", 1, 2, -33},
		{"agenda", "Dos médicos con el mismo consultorio quedan con turnos superpuestos", "Alta", 1, 2, -29},
		{"reservar", "Se pueden reservar dos turnos en el mismo horario desde dos pestañas", "Alta", 2, 2, -21},
		{"cancelar", "El mensaje de cancelación no indica el nombre del médico", "Baja", 2, 3, -14},
		{"recordatorios", "El recordatorio se envía dos veces cuando el turno se reprograma", "Media", 3, 0, -3},
		{"historial", "El historial no pagina y tarda con más de 200 consultas", "Alta", 3, 0, -2},
	},
}

var bibliotecaCerrada = proyectoDemo{
	nombre:      "Sistema de biblioteca (datos de prueba)",
	descripcion: "Catálogo, socios y préstamos de una biblioteca barrial. Proyecto cerrado por el Scrum Master.",
	inicio:      -100, fin: 15, cierre: dia(-4), horasPorSP: 3,
	equipo: equipoCompleto,
	sprints: []sprintDemo{
		{"Catálogo, socios y préstamos", -98, -85, domain.SprintFinalizado},
		{"Devoluciones, multas y reservas", -84, -71, domain.SprintFinalizado},
	},
	historias: []historiaDemo{
		{clave: "catalogo", titulo: "Consultar el catálogo de libros", descripcion: "Como socio quiero buscar libros por título o autor.",
			prioridad: domain.PrioridadMust, valor: 8, puntos: 5, sprint: 1, estado: comp, finalizada: -90},
		{clave: "socios", titulo: "Registrar socios", descripcion: "Como bibliotecario quiero dar de alta socios con su DNI.",
			prioridad: domain.PrioridadMust, valor: 8, puntos: 3, sprint: 1, estado: comp, finalizada: -88},
		{clave: "prestamo", titulo: "Registrar un préstamo", descripcion: "Como bibliotecario quiero prestar un libro a un socio y fijar la fecha de devolución.",
			prioridad: domain.PrioridadMust, valor: 13, puntos: 8, sprint: 1, estado: comp, finalizada: -86},
		{clave: "devolucion", titulo: "Registrar una devolución", descripcion: "Como bibliotecario quiero marcar el libro como devuelto.",
			prioridad: domain.PrioridadMust, valor: 13, puntos: 5, sprint: 2, estado: comp, finalizada: -76},
		{clave: "multas", titulo: "Calcular multas por demora", descripcion: "Como bibliotecario quiero que el sistema calcule la multa por días de atraso.",
			prioridad: domain.PrioridadShould, valor: 5, puntos: 3, sprint: 2, estado: comp, finalizada: -73},
		{clave: "reservas", titulo: "Reservar un libro prestado", descripcion: "Como socio quiero reservar un libro que está prestado.",
			prioridad: domain.PrioridadCould, valor: 3, puntos: 2, sprint: 2, estado: comp, finalizada: -72},
	},
	horas: []horasDemo{
		todos("catalogo", 1, "candela.vulcano", -96, "6", "Pantalla del catálogo"),
		todos("socios", 1, "ana.iriarte", -94, "3", "Alta de socios"),
		todos("prestamo", 1, "carlos.moya", -91, "8", "Regla de préstamo"),
		todos("prestamo", 1, "jimena.martinez", -90, "1:30", "Pruebas de préstamo"),
		todos("devolucion", 2, "candela.vulcano", -80, "4", "Devolución de libros"),
		todos("multas", 2, "ana.iriarte", -77, "2.5", "Cálculo de la multa"),
		todos("reservas", 2, "jimena.martinez", -74, "2", "Reservas de libros"),
	},
	defectos: []defectoDemo{
		{"prestamo", "Un socio con multa pendiente puede retirar otro libro", "Media", 1, 2, -88},
	},
}

var tiendaFinalizada = proyectoDemo{
	nombre:      "Tienda de barrio (datos de prueba)",
	descripcion: "Ventas y stock de un comercio de barrio. Terminó por fecha de fin, sin cierre manual.",
	inicio:      -80, fin: -3, horasPorSP: 4,
	equipo: []miembroDemo{ana, carlos, candela},
	sprints: []sprintDemo{
		{"Productos y stock", -78, -65, domain.SprintFinalizado},
		{"Ventas y caja", -64, -51, domain.SprintFinalizado},
	},
	historias: []historiaDemo{
		{clave: "productos", titulo: "Cargar productos con su precio", descripcion: "Como comerciante quiero cargar mis productos.",
			prioridad: domain.PrioridadMust, valor: 8, puntos: 3, sprint: 1, estado: comp, finalizada: -70},
		{clave: "stock", titulo: "Controlar el stock", descripcion: "Como comerciante quiero ver cuánto me queda de cada producto.",
			prioridad: domain.PrioridadMust, valor: 8, puntos: 5, sprint: 1, estado: comp, finalizada: -67},
		{clave: "venta", titulo: "Registrar una venta", descripcion: "Como vendedor quiero registrar la venta y descontar el stock.",
			prioridad: domain.PrioridadMust, valor: 13, puntos: 8, sprint: 2, estado: comp, finalizada: -58},
		{clave: "caja", titulo: "Cerrar la caja del día", descripcion: "Como comerciante quiero un resumen de lo vendido al cerrar.",
			prioridad: domain.PrioridadShould, valor: 5, puntos: 3, sprint: 2, estado: comp, finalizada: -53},
	},
	horas: []horasDemo{
		todos("productos", 1, "candela.vulcano", -75, "3", "ABM de productos"),
		todos("stock", 1, "carlos.moya", -70, "5", "Control de stock"),
		todos("venta", 2, "ana.iriarte", -60, "6", "Registro de ventas"),
		todos("venta", 2, "candela.vulcano", -58, "2:30", "Descuento de stock"),
		todos("caja", 2, "carlos.moya", -54, "2", "Resumen de caja"),
	},
	defectos: []defectoDemo{
		{"venta", "La venta de un producto sin stock no muestra advertencia", "Media", 2, 2, -59},
	},
}

var appReservas = proyectoDemo{
	nombre:      "App de reservas (datos de prueba)",
	descripcion: "Reservas de canchas deportivas. Todavía no empezó: sirve para arrancar un Sprint y estimar con Planning Poker desde cero.",
	inicio:      10, fin: 80,
	equipo:  equipoCompleto,
	sprints: []sprintDemo{{estado: domain.SprintPendiente}},
	historias: []historiaDemo{
		{clave: "canchas", titulo: "Ver las canchas disponibles", descripcion: "Como jugador quiero ver qué canchas están libres.",
			prioridad: domain.PrioridadMust, valor: 13, estado: nva},
		{clave: "reservar", titulo: "Reservar una cancha", descripcion: "Como jugador quiero reservar en un horario libre.",
			prioridad: domain.PrioridadMust, valor: 21, estado: nva},
		{clave: "pago", titulo: "Pagar la seña", descripcion: "Como jugador quiero pagar una seña para confirmar.",
			prioridad: domain.PrioridadShould, valor: 8, estado: nva},
		{clave: "equipos", titulo: "Armar equipos", descripcion: "Como jugador quiero invitar a mis amigos al partido.",
			prioridad: domain.PrioridadCould, valor: 5, estado: nva},
		{clave: "resenas", titulo: "Dejar reseñas de las canchas", descripcion: "Como jugador quiero opinar sobre la cancha.",
			prioridad: domain.PrioridadCould, valor: 3, estado: nva},
		{clave: "torneos", titulo: "Organizar torneos", descripcion: "Como organizador quiero armar un torneo con fixture.",
			prioridad: domain.PrioridadWont, estado: nva},
	},
}

var portalAlumnos = proyectoDemo{
	nombre:      "Portal de alumnos (datos de prueba)",
	descripcion: "Inscripciones y notas de una escuela. Un integrante fue dado de baja: sus horas y votos ya revelados se conservan.",
	inicio:      -40, fin: 40, horasPorSP: 4,
	equipo: []miembroDemo{ana, carlos, candela, {usuario: "jimena.martinez", rol: domain.RolProductBuilder, baja: dia(-10)}},
	sprints: []sprintDemo{
		{"Inscripción a materias", -35, -22, domain.SprintFinalizado},
		{"Notas y asistencia", -7, 6, domain.SprintActivo},
	},
	historias: []historiaDemo{
		{clave: "inscripcion", titulo: "Inscribirse a una materia", descripcion: "Como alumno quiero inscribirme a una materia del cuatrimestre.",
			prioridad: domain.PrioridadMust, valor: 13, puntos: 5, sprint: 1, estado: comp, finalizada: -26},
		{clave: "materias", titulo: "Ver las materias disponibles", descripcion: "Como alumno quiero ver las materias con sus correlatividades.",
			prioridad: domain.PrioridadShould, valor: 8, puntos: 3, sprint: 1, estado: comp, finalizada: -24},
		{clave: "notas", titulo: "Consultar las notas", descripcion: "Como alumno quiero ver mis notas por materia.",
			prioridad: domain.PrioridadMust, valor: 13, puntos: 5, sprint: 2, estado: prog},
		{clave: "asistencia", titulo: "Registrar la asistencia", descripcion: "Como docente quiero tomar lista de mis alumnos.",
			prioridad: domain.PrioridadShould, valor: 5, puntos: 3, sprint: 2, estado: nva},
		{clave: "boletin", titulo: "Descargar el boletín", descripcion: "Como alumno quiero descargar mi boletín en PDF.",
			prioridad: domain.PrioridadCould, valor: 5, estado: nva},
		{clave: "avisos", titulo: "Avisos de la escuela", descripcion: "Como alumno quiero ver los comunicados de la escuela.",
			prioridad: domain.PrioridadCould, valor: 3, estado: nva},
	},
	horas: []horasDemo{
		todos("inscripcion", 1, "candela.vulcano", -33, "6", "Pantalla de inscripción"),
		todos("inscripcion", 1, "jimena.martinez", -31, "3", "Pruebas con alumnos"),
		todos("materias", 1, "jimena.martinez", -28, "2.5", "Listado de materias"),
		todos("materias", 1, "ana.iriarte", -27, "1:30", "Correlatividades"),
		todos("inscripcion", 1, "carlos.moya", -25, "4", "Revisión de la inscripción"),
		todos("notas", 2, "candela.vulcano", -5, "5", "Consulta de notas"),
		todos("notas", 2, "ana.iriarte", -3, "2", "Diseño del boletín"),
	},
	rondas: []rondaAdicional{
		{"inscripcion", 1, []votoDemo{{"ana.iriarte", 5}, {"carlos.moya", 5}, {"candela.vulcano", 5}, {"jimena.martinez", 5}}, true},
		{"notas", 1, []votoDemo{{"ana.iriarte", 3}, {"carlos.moya", 5}, {"candela.vulcano", 5}, {"jimena.martinez", 8}}, true},
	},
	defectos: []defectoDemo{
		{"inscripcion", "Un alumno puede inscribirse dos veces a la misma materia", "Media", 1, 0, -25},
	},
}

var proyectoEnBlanco = proyectoDemo{
	nombre:      "Proyecto en blanco (datos de prueba)",
	descripcion: "Solo el equipo: sin Sprints, historias ni defectos, para mostrar los estados vacíos.",
	inicio:      5, fin: 60,
	equipo: equipoCompleto,
}

var campanaDonaciones = proyectoDemo{
	nombre:      "Campaña de donaciones (datos de prueba)",
	descripcion: "Todos los Sprints terminaron pero la fecha de fin es futura: sigue «En curso» hasta que el Scrum Master lo cierra.",
	inicio:      -50, fin: 15, horasPorSP: 4,
	equipo: equipoCompleto,
	sprints: []sprintDemo{
		{"Alta de campañas y donantes", -48, -35, domain.SprintFinalizado},
		{"Donaciones y comprobantes", -34, -21, domain.SprintFinalizado},
	},
	historias: []historiaDemo{
		{clave: "campanas", titulo: "Crear una campaña", descripcion: "Como organizador quiero crear una campaña con su meta.",
			prioridad: domain.PrioridadMust, valor: 8, puntos: 3, sprint: 1, estado: comp, finalizada: -40},
		{clave: "donantes", titulo: "Registrar donantes", descripcion: "Como organizador quiero registrar a quienes donan.",
			prioridad: domain.PrioridadShould, valor: 5, puntos: 3, sprint: 1, estado: comp, finalizada: -37},
		{clave: "donar", titulo: "Registrar una donación", descripcion: "Como donante quiero aportar un monto a una campaña.",
			prioridad: domain.PrioridadMust, valor: 13, puntos: 8, sprint: 2, estado: comp, finalizada: -28},
		{clave: "comprobante", titulo: "Emitir el comprobante de la donación", descripcion: "Como donante quiero un comprobante de mi aporte.",
			prioridad: domain.PrioridadShould, valor: 8, puntos: 3, sprint: 2, estado: comp, finalizada: -24},
		{clave: "meta", titulo: "Mostrar el avance hacia la meta", descripcion: "Como organizador quiero ver cuánto falta para la meta.",
			prioridad: domain.PrioridadCould, valor: 5, puntos: 2, sprint: 2, estado: comp, finalizada: -22},
	},
	horas: []horasDemo{
		todos("campanas", 1, "ana.iriarte", -46, "3", "Alta de campañas"),
		todos("donantes", 1, "candela.vulcano", -41, "2:30", "Registro de donantes"),
		todos("donar", 2, "carlos.moya", -31, "6", "Registro de donaciones"),
		todos("donar", 2, "jimena.martinez", -29, "1.5", "Pruebas de donación"),
		todos("comprobante", 2, "candela.vulcano", -25, "3", "Comprobante en PDF"),
		todos("meta", 2, "ana.iriarte", -23, "0:45", "Barra de avance"),
	},
	defectos: []defectoDemo{
		{"donar", "Una donación en cero se acepta como válida", "Baja", 2, 2, -30},
	},
}
