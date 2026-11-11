// Package app es la raíz de composición: arma repositorios, servicios y
// handlers sobre una base ya migrada y devuelve el mux con todas las rutas. Lo
// usan cmd/api y los escenarios de Godog, así ambos ejercitan la misma API.
package app

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/demo"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	apihttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/report"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// Opciones son las dependencias que cambian entre producción y pruebas.
type Opciones struct {
	DemoHabilitado func() bool
	Ahora          func() time.Time
	Contrasenas    *service.ServicioContrasenas
}

// NuevoMux arma la API completa sobre db.
func NuevoMux(db *sql.DB, op Opciones) *http.ServeMux {
	backlogRepo := repository.NewSQLiteBacklogRepository(db)
	proyectos := repository.NewSQLiteProjectRepository(db)
	usuarios := repository.NewSQLiteUserRepository(db)
	credenciales := repository.NewSQLiteCredencialRepository(db)
	sprintRepo := repository.NewSQLiteSprintRepository(db)

	backlogHandler := apihttp.NewBacklogHandler(
		service.NewCrearHistoriaBacklog(backlogRepo, proyectos),
		service.NewConsultarProductBacklog(proyectos, backlogRepo),
		service.NewObtenerHistoria(backlogRepo),
		service.NewEditarHistoria(backlogRepo, sprintRepo, proyectos),
		service.NewAsignarHistoriaASprint(backlogRepo, sprintRepo, proyectos),
		service.NewCambiarEstadoHistoria(backlogRepo, sprintRepo, op.Ahora, proyectos),
	)
	obtenerProyecto := service.NewObtenerProyecto(proyectos)
	projectHandler := apihttp.NewProjectHandler(
		service.NewCrearProyecto(proyectos),
		obtenerProyecto,
		service.NewListarProyectos(proyectos),
		service.NewEditarProyecto(proyectos),
		service.NewAsignarIntegrante(proyectos, credenciales),
		service.NewListarIntegrantes(proyectos),
		service.NewObtenerEstadoProyecto(proyectos, sprintRepo),
	)

	sprintHandler := apihttp.NewSprintHandler(
		service.NewCrearSprint(proyectos, sprintRepo),
		service.NewIniciarSprint(sprintRepo, proyectos),
		service.NewCerrarSprint(sprintRepo, backlogRepo, proyectos),
		service.NewConsultarSprints(proyectos, sprintRepo),
		service.NewConsultarSprintBacklog(sprintRepo, backlogRepo),
		sprintRepo,
	)

	horasRepo := repository.NewSQLiteRegistroHorasRepository(db)
	horasHandler := apihttp.NewHorasHandler(
		service.NewImputarHoras(backlogRepo, sprintRepo, proyectos, horasRepo),
		service.NewObtenerHistoria(backlogRepo),
		service.NewConsultarEsfuerzoReal(backlogRepo, horasRepo),
	)

	votaciones := repository.NewSQLiteVotacionRepository(db)
	votacionHandler := apihttp.NewVotacionHandler(
		service.NewObtenerHistoria(backlogRepo),
		service.NewVotar(backlogRepo, proyectos, votaciones),
		service.NewConsultarVotacion(backlogRepo, proyectos, votaciones, domain.CoincidenciaExacta{}),
		service.NewRevelarVotacion(backlogRepo, proyectos, votaciones),
		service.NewIniciarVotacion(backlogRepo, proyectos, votaciones),
		service.NewRegistrarEstimacion(backlogRepo, proyectos, votaciones, sprintRepo),
		service.NewAbrirNuevaRonda(backlogRepo, proyectos, votaciones),
		service.NewConsultarHistorialRondas(backlogRepo, proyectos, votaciones),
	)

	defectosRepo := repository.NewSQLiteDefectoRepository(db)
	defectoHandler := apihttp.NewDefectoHandler(
		obtenerProyecto,
		service.NewRegistrarDefecto(sprintRepo, backlogRepo, defectosRepo, proyectos),
		service.NewResolverDefecto(defectosRepo, sprintRepo, proyectos),
		service.NewListarDefectos(proyectos, defectosRepo),
	)

	metricasRepo := repository.NewSQLiteMetricasRepository(db)
	sesionesRepo := repository.NewSQLiteSesionRepository(db)
	emisorSesiones := service.NewEmisorSesiones(sesionesRepo, op.Ahora)
	obtenerSesion := service.NewObtenerSesion(sesionesRepo, op.Ahora)
	authHandler := apihttp.NewAuthHandler(
		service.NewRegistrarUsuario(credenciales, op.Contrasenas, emisorSesiones),
		service.NewIngresar(credenciales, op.Contrasenas, emisorSesiones),
		obtenerSesion,
		service.NewCerrarSesion(obtenerSesion, sesionesRepo, op.Ahora),
	)

	calcularMetricas := service.NewCalcularMetricas(proyectos, metricasRepo)
	historialHandler := apihttp.NewHistorialHandler(
		service.NewConsultarHistorialSprints(calcularMetricas, sprintRepo),
		service.NewConsultarDetalleSprint(proyectos, sprintRepo, calcularMetricas, repository.NewSQLiteHistorialSprintRepository(db)),
	)
	reporteHandler := apihttp.NewReporteHandler(
		service.NewGenerarReporte(proyectos, sprintRepo, backlogRepo, backlogRepo, backlogRepo, metricasRepo, defectosRepo),
		report.Generador{},
	)

	mux := http.NewServeMux()
	apihttp.RegistrarRutas(mux, apihttp.Handlers{
		Auth:      authHandler,
		Proyectos: projectHandler,
		Backlog:   backlogHandler,
		Horas:     horasHandler,
		Sprints:   sprintHandler,
		Votacion:  votacionHandler,
		Defectos:  defectoHandler,
		Metricas:  apihttp.NewMetricasHandler(calcularMetricas),
		Historial: historialHandler,
		Reporte:   reporteHandler,
		Avance:    apihttp.NewAvanceHandler(service.NewConsultarAvanceHistorias(proyectos, backlogRepo, sprintRepo)),
		Demo:      apihttp.NewDemoHandler(op.DemoHabilitado, demo.NewCargador(db, op.Ahora, op.Contrasenas)),
		Usuarios:  apihttp.NewUsuariosHandler(service.NewListarUsuarios(usuarios)),
	}, apihttp.ConSesion(obtenerSesion))
	return mux
}
