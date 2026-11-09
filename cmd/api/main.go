package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/demo"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	apihttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/report"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

func main() {
	ctx := context.Background()
	cfg := cargarConfig(os.Getenv)

	db, err := repository.AbrirSQLite(cfg.RutaBase)
	if err != nil {
		log.Fatalf("abrir base de datos: %v", err)
	}
	defer db.Close()

	if err := repository.AplicarMigraciones(ctx, db); err != nil {
		log.Fatalf("migrar base de datos: %v", err)
	}

	backlogRepo := repository.NewSQLiteBacklogRepository(db)

	proyectos := repository.NewSQLiteProjectRepository(db)
	crearHistoria := service.NewCrearHistoriaBacklog(backlogRepo, proyectos)
	usuarios := repository.NewSQLiteUserRepository(db)
	sprintRepo := repository.NewSQLiteSprintRepository(db)
	backlogHandler := apihttp.NewBacklogHandler(
		crearHistoria,
		service.NewConsultarProductBacklog(proyectos, backlogRepo),
		service.NewObtenerHistoria(backlogRepo),
		service.NewEditarHistoria(backlogRepo, sprintRepo),
		service.NewAsignarHistoriaASprint(backlogRepo, sprintRepo),
		service.NewCambiarEstadoHistoria(backlogRepo, sprintRepo, time.Now),
	)
	crearProyecto := service.NewCrearProyecto(proyectos, usuarios)
	obtenerProyecto := service.NewObtenerProyecto(proyectos)
	listarProyectos := service.NewListarProyectos(proyectos)
	editarProyecto := service.NewEditarProyecto(proyectos)
	asignarIntegrante := service.NewAsignarIntegrante(proyectos, usuarios)
	listarIntegrantes := service.NewListarIntegrantes(proyectos)
	obtenerEstadoProyecto := service.NewObtenerEstadoProyecto(proyectos, sprintRepo)
	projectHandler := apihttp.NewProjectHandler(crearProyecto, obtenerProyecto, listarProyectos, editarProyecto, asignarIntegrante, listarIntegrantes, obtenerEstadoProyecto)

	sprintHandler := apihttp.NewSprintHandler(
		service.NewCrearSprint(proyectos, sprintRepo),
		service.NewIniciarSprint(sprintRepo),
		service.NewCerrarSprint(sprintRepo, backlogRepo),
		service.NewConsultarSprints(proyectos, sprintRepo),
		service.NewConsultarSprintBacklog(sprintRepo, backlogRepo),
		sprintRepo,
	)

	horasRepo := repository.NewSQLiteRegistroHorasRepository(db)
	imputarHoras := service.NewImputarHoras(backlogRepo, sprintRepo, proyectos, horasRepo)
	horasHandler := apihttp.NewHorasHandler(
		imputarHoras,
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
		service.NewRegistrarDefecto(sprintRepo, backlogRepo, defectosRepo),
		service.NewResolverDefecto(defectosRepo, sprintRepo),
		service.NewListarDefectos(proyectos, defectosRepo),
	)

	metricasRepo := repository.NewSQLiteMetricasRepository(db)
	credenciales := repository.NewSQLiteCredencialRepository(db)
	hasheador := service.NewServicioContrasenasProduccion()
	sesionesRepo := repository.NewSQLiteSesionRepository(db)
	emisorSesiones := service.NewEmisorSesiones(sesionesRepo, time.Now)
	obtenerSesion := service.NewObtenerSesion(sesionesRepo, time.Now)
	authHandler := apihttp.NewAuthHandler(
		service.NewRegistrarUsuario(credenciales, hasheador, emisorSesiones),
		service.NewIngresar(credenciales, hasheador, emisorSesiones),
		obtenerSesion,
		service.NewCerrarSesion(obtenerSesion, sesionesRepo, time.Now),
	)

	calcularMetricas := service.NewCalcularMetricas(proyectos, metricasRepo)
	metricasHandler := apihttp.NewMetricasHandler(calcularMetricas)
	historialHandler := apihttp.NewHistorialHandler(
		service.NewConsultarHistorialSprints(calcularMetricas, sprintRepo),
		service.NewConsultarDetalleSprint(proyectos, sprintRepo, calcularMetricas, repository.NewSQLiteHistorialSprintRepository(db)),
	)
	reporteHandler := apihttp.NewReporteHandler(
		service.NewGenerarReporte(proyectos, sprintRepo, backlogRepo, backlogRepo, backlogRepo, metricasRepo, defectosRepo),
		report.Generador{},
	)
	avanceHandler := apihttp.NewAvanceHandler(service.NewConsultarAvanceHistorias(proyectos, backlogRepo, sprintRepo))

	demoHandler := apihttp.NewDemoHandler(
		func() bool { return cfg.DemoHabilitado },
		demo.NewCargador(db, time.Now),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /auth/registro", authHandler.Registro)
	mux.HandleFunc("POST /auth/ingreso", authHandler.Ingreso)
	mux.HandleFunc("GET /auth/yo", authHandler.Yo)
	mux.HandleFunc("POST /auth/salir", authHandler.Salir)
	mux.HandleFunc("POST /backlog", backlogHandler.Crear)
	mux.HandleFunc("GET /projects/{id}/backlog", backlogHandler.ConsultarBacklog)
	mux.HandleFunc("PUT /backlog/{id}", backlogHandler.Editar)
	mux.HandleFunc("GET /backlog/{id}", backlogHandler.Obtener)
	mux.HandleFunc("POST /backlog/{id}/horas", horasHandler.Imputar)
	mux.HandleFunc("GET /backlog/{id}/horas", horasHandler.ConsultarEsfuerzo)
	mux.HandleFunc("PUT /backlog/{id}/sprint", backlogHandler.AsignarSprint)
	mux.HandleFunc("PUT /backlog/{id}/estado", backlogHandler.CambiarEstado)
	mux.HandleFunc("GET /projects", projectHandler.Listar)
	mux.HandleFunc("POST /projects", projectHandler.Crear)
	mux.HandleFunc("GET /projects/{id}", projectHandler.Obtener)
	mux.HandleFunc("GET /projects/{id}/status", projectHandler.ObtenerEstado)
	mux.HandleFunc("PUT /projects/{id}", projectHandler.Editar)
	mux.HandleFunc("POST /projects/{id}/members", projectHandler.AsignarIntegrante)
	mux.HandleFunc("GET /projects/{id}/members", projectHandler.ListarIntegrantes)
	mux.HandleFunc("GET /projects/{id}/sprints", sprintHandler.Listar)
	mux.HandleFunc("GET /projects/{id}/metricas", metricasHandler.Consultar)
	mux.HandleFunc("GET /projects/{id}/historial", historialHandler.Listar)
	mux.HandleFunc("GET /projects/{id}/historial/{sprintId}", historialHandler.Detalle)
	mux.HandleFunc("GET /projects/{id}/avance", avanceHandler.Consultar)
	mux.HandleFunc("POST /projects/{id}/defectos", defectoHandler.Registrar)
	mux.HandleFunc("GET /projects/{id}/defectos", defectoHandler.Listar)
	mux.HandleFunc("POST /defectos/{id}/resolver", defectoHandler.Resolver)
	mux.HandleFunc("GET /sprints/{id}/reporte", reporteHandler.Sprint)
	mux.HandleFunc("GET /projects/{id}/reporte", reporteHandler.Proyecto)
	mux.HandleFunc("POST /sprints", sprintHandler.Crear)
	mux.HandleFunc("POST /sprints/{id}/iniciar", sprintHandler.Iniciar)
	mux.HandleFunc("POST /sprints/{id}/cerrar", sprintHandler.Cerrar)
	mux.HandleFunc("GET /sprints/{id}/backlog", sprintHandler.ConsultarBacklog)

	votacionHandler.Registrar(mux)
	demoHandler.Registrar(mux)

	log.Printf("listening on :%s", cfg.Puerto)
	if err := http.ListenAndServe(":"+cfg.Puerto, apihttp.ServirSPA(mux, cfg.DirWeb)); err != nil {
		log.Fatal(err)
	}
}
