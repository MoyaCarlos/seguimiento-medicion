package http

import (
	nethttp "net/http"
)

// NivelRuta es lo que una ruta exige a quien la usa (tabla «Qué se protege»
// de la spec 042, HU-21). La sesión la exige la capa HTTP; la pertenencia al
// proyecto y el rol los comprueban los servicios.
type NivelRuta int

const (
	// Publica no exige sesión.
	Publica NivelRuta = iota
	// SoloSesion exige sesión, sin pertenencia ni rol.
	SoloSesion
	// Lectura es un GET con sesión, sin pertenencia: solo devuelve datos de la
	// propia persona (GET /auth/yo, GET /projects, que lista solo sus proyectos).
	Lectura
	// LecturaDeProyecto es un GET con sesión que devuelve datos de un proyecto:
	// exige ser integrante, también si el proyecto está cerrado (403 si no).
	LecturaDeProyecto
	// SesionIntegrante exige sesión y ser integrante del proyecto.
	SesionIntegrante
	// SesionScrumMaster exige sesión y ser el Scrum Master del proyecto.
	SesionScrumMaster
)

// Ruta es una entrada de la tabla de rutas de la API.
type Ruta struct {
	Patron  string
	Nivel   NivelRuta
	Handler nethttp.HandlerFunc
}

// Handlers agrupa los handlers ya construidos por el arranque.
type Handlers struct {
	Auth      *AuthHandler
	Proyectos *ProjectHandler
	Backlog   *BacklogHandler
	Horas     *HorasHandler
	Sprints   *SprintHandler
	Votacion  *VotacionHandler
	Defectos  *DefectoHandler
	Metricas  *MetricasHandler
	Historial *HistorialHandler
	Reporte   *ReporteHandler
	Avance    *AvanceHandler
	Demo      *DemoHandler
	Usuarios  *UsuariosHandler
}

// Rutas devuelve la tabla completa de rutas de la API con su nivel.
func Rutas(h Handlers) []Ruta {
	return []Ruta{
		{"GET /health", Publica, salud},
		{"POST /auth/registro", Publica, h.Auth.Registro},
		{"POST /auth/ingreso", Publica, h.Auth.Ingreso},
		{"GET /auth/yo", Lectura, h.Auth.Yo},
		{"PUT /auth/yo", SoloSesion, h.Auth.CambiarNombre},
		{"POST /auth/salir", SoloSesion, h.Auth.Salir},
		{"GET /demo/estado", Publica, h.Demo.soloSiHabilitado(h.Demo.Estado)},
		{"POST /demo/datos", Publica, h.Demo.soloSiHabilitado(h.Demo.Cargar)},
		{"/demo/", Publica, nethttp.NotFound},

		{"GET /users", SoloSesion, h.Usuarios.Listar},

		{"GET /projects", Lectura, h.Proyectos.Listar},
		{"POST /projects", SoloSesion, h.Proyectos.Crear},
		{"GET /projects/{id}", LecturaDeProyecto, h.Proyectos.Obtener},
		{"GET /projects/{id}/status", LecturaDeProyecto, h.Proyectos.ObtenerEstado},
		{"PUT /projects/{id}", SesionScrumMaster, h.Proyectos.Editar},
		{"POST /projects/{id}/members", SesionScrumMaster, h.Proyectos.AsignarIntegrante},
		{"DELETE /projects/{id}/members/{usuario_id}", SesionScrumMaster, h.Proyectos.QuitarIntegrante},
		{"POST /projects/{id}/cerrar", SesionScrumMaster, h.Proyectos.Cerrar},
		{"POST /projects/{id}/reabrir", SesionScrumMaster, h.Proyectos.Reabrir},
		{"GET /projects/{id}/members", LecturaDeProyecto, h.Proyectos.ListarIntegrantes},

		{"POST /backlog", SesionIntegrante, h.Backlog.Crear},
		{"GET /projects/{id}/backlog", LecturaDeProyecto, h.Backlog.ConsultarBacklog},
		{"PUT /backlog/{id}", SesionIntegrante, h.Backlog.Editar},
		{"GET /backlog/{id}", LecturaDeProyecto, h.Backlog.Obtener},
		{"PUT /backlog/{id}/sprint", SesionIntegrante, h.Backlog.AsignarSprint},
		{"PUT /backlog/{id}/estado", SesionIntegrante, h.Backlog.CambiarEstado},
		{"POST /backlog/{id}/horas", SesionIntegrante, h.Horas.Imputar},
		{"GET /backlog/{id}/horas", LecturaDeProyecto, h.Horas.ConsultarEsfuerzo},

		{"GET /backlog/{id}/votacion", LecturaDeProyecto, h.Votacion.Consultar},
		{"POST /backlog/{id}/votacion", SesionScrumMaster, h.Votacion.Iniciar},
		{"POST /backlog/{id}/votacion/votos", SesionIntegrante, h.Votacion.Votar},
		{"POST /backlog/{id}/votacion/revelar", SesionScrumMaster, h.Votacion.Revelar},
		{"POST /backlog/{id}/votacion/estimacion", SesionScrumMaster, h.Votacion.RegistrarEstimacion},
		{"POST /backlog/{id}/votacion/rondas", SesionScrumMaster, h.Votacion.AbrirRonda},
		{"GET /backlog/{id}/votacion/rondas", LecturaDeProyecto, h.Votacion.Historial},

		{"POST /sprints", SesionScrumMaster, h.Sprints.Crear},
		{"POST /sprints/{id}/iniciar", SesionScrumMaster, h.Sprints.Iniciar},
		{"POST /sprints/{id}/cerrar", SesionScrumMaster, h.Sprints.Cerrar},
		{"GET /sprints/{id}/backlog", LecturaDeProyecto, h.Sprints.ConsultarBacklog},
		{"GET /projects/{id}/sprints", LecturaDeProyecto, h.Sprints.Listar},

		{"POST /projects/{id}/defectos", SesionIntegrante, h.Defectos.Registrar},
		{"GET /projects/{id}/defectos", LecturaDeProyecto, h.Defectos.Listar},
		{"POST /defectos/{id}/resolver", SesionIntegrante, h.Defectos.Resolver},

		{"GET /projects/{id}/metricas", LecturaDeProyecto, h.Metricas.Consultar},
		{"GET /projects/{id}/historial", LecturaDeProyecto, h.Historial.Listar},
		{"GET /projects/{id}/historial/{sprintId}", LecturaDeProyecto, h.Historial.Detalle},
		{"GET /projects/{id}/avance", LecturaDeProyecto, h.Avance.Consultar},
		{"GET /sprints/{id}/reporte", LecturaDeProyecto, h.Reporte.Sprint},
		{"GET /projects/{id}/reporte", LecturaDeProyecto, h.Reporte.Proyecto},
	}
}

// RegistrarRutas registra en el mux todas las rutas de la tabla y envuelve con
// conSesion las que no son públicas.
func RegistrarRutas(mux *nethttp.ServeMux, h Handlers, conSesion func(nethttp.HandlerFunc) nethttp.HandlerFunc) {
	for _, r := range Rutas(h) {
		handler := r.Handler
		if r.Nivel != Publica {
			handler = conSesion(handler)
		}
		mux.HandleFunc(r.Patron, handler)
	}
}

func salud(w nethttp.ResponseWriter, _ *nethttp.Request) {
	_, _ = w.Write([]byte("ok"))
}
