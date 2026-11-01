package main

import (
	"context"
	"log"
	"net/http"
	"os"

	apihttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
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
		service.NewConsultarVotacion(backlogRepo, proyectos, votaciones),
		service.NewRevelarVotacion(backlogRepo, proyectos, votaciones),
		service.NewIniciarVotacion(backlogRepo, proyectos, votaciones),
	)

	defectosRepo := repository.NewSQLiteDefectoRepository(db)
	defectoHandler := apihttp.NewDefectoHandler(
		obtenerProyecto,
		service.NewRegistrarDefecto(sprintRepo, backlogRepo, defectosRepo),
		service.NewResolverDefecto(defectosRepo, sprintRepo),
		service.NewListarDefectos(proyectos, defectosRepo),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /backlog", backlogHandler.Crear)
	mux.HandleFunc("GET /projects/{id}/backlog", backlogHandler.ConsultarBacklog)
	mux.HandleFunc("PUT /backlog/{id}", backlogHandler.Editar)
	mux.HandleFunc("POST /backlog/{id}/horas", horasHandler.Imputar)
	mux.HandleFunc("GET /backlog/{id}/horas", horasHandler.ConsultarEsfuerzo)
	mux.HandleFunc("GET /backlog/{id}/votacion", votacionHandler.Consultar)
	mux.HandleFunc("POST /backlog/{id}/votacion", votacionHandler.Iniciar)
	mux.HandleFunc("POST /backlog/{id}/votacion/votos", votacionHandler.Votar)
	mux.HandleFunc("POST /backlog/{id}/votacion/revelar", votacionHandler.Revelar)
	mux.HandleFunc("PUT /backlog/{id}/sprint", backlogHandler.AsignarSprint)
	mux.HandleFunc("GET /projects", projectHandler.Listar)
	mux.HandleFunc("POST /projects", projectHandler.Crear)
	mux.HandleFunc("GET /projects/{id}", projectHandler.Obtener)
	mux.HandleFunc("GET /projects/{id}/status", projectHandler.ObtenerEstado)
	mux.HandleFunc("PUT /projects/{id}", projectHandler.Editar)
	mux.HandleFunc("POST /projects/{id}/members", projectHandler.AsignarIntegrante)
	mux.HandleFunc("GET /projects/{id}/members", projectHandler.ListarIntegrantes)
	mux.HandleFunc("GET /projects/{id}/sprints", sprintHandler.Listar)
	mux.HandleFunc("POST /projects/{id}/defectos", defectoHandler.Registrar)
	mux.HandleFunc("GET /projects/{id}/defectos", defectoHandler.Listar)
	mux.HandleFunc("POST /defectos/{id}/resolver", defectoHandler.Resolver)
	mux.HandleFunc("POST /sprints", sprintHandler.Crear)
	mux.HandleFunc("POST /sprints/{id}/iniciar", sprintHandler.Iniciar)
	mux.HandleFunc("POST /sprints/{id}/cerrar", sprintHandler.Cerrar)

	log.Printf("listening on :%s", cfg.Puerto)
	if err := http.ListenAndServe(":"+cfg.Puerto, apihttp.ServirSPA(mux, cfg.DirWeb)); err != nil {
		log.Fatal(err)
	}
}
