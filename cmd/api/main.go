package main

import (
	"context"
	"log"
	"net/http"

	apihttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

func main() {
	ctx := context.Background()

	db, err := repository.AbrirSQLite("seguimiento.db")
	if err != nil {
		log.Fatalf("abrir base de datos: %v", err)
	}
	defer db.Close()

	if err := repository.Migrar(ctx, db); err != nil {
		log.Fatalf("migrar base de datos: %v", err)
	}

	backlogRepo := repository.NewSQLiteBacklogRepository(db)
	crearHistoria := service.NewCrearHistoriaBacklog(backlogRepo)
	backlogHandler := apihttp.NewBacklogHandler(crearHistoria)

	proyectos := repository.NewSQLiteProjectRepository(db)
	usuarios := repository.NewSQLiteUserRepository(db)
	sprintRepo := repository.NewSQLiteSprintRepository(db)
	crearProyecto := service.NewCrearProyecto(proyectos, usuarios)
	obtenerProyecto := service.NewObtenerProyecto(proyectos)
	editarProyecto := service.NewEditarProyecto(proyectos)
	asignarIntegrante := service.NewAsignarIntegrante(proyectos, usuarios)
	listarIntegrantes := service.NewListarIntegrantes(proyectos)
	obtenerEstadoProyecto := service.NewObtenerEstadoProyecto(proyectos, sprintRepo)
	projectHandler := apihttp.NewProjectHandler(crearProyecto, obtenerProyecto, editarProyecto, asignarIntegrante, listarIntegrantes, obtenerEstadoProyecto)

	sprintHandler := apihttp.NewSprintHandler(
		service.NewCrearSprint(proyectos, sprintRepo),
		service.NewIniciarSprint(sprintRepo),
		service.NewCerrarSprint(sprintRepo, backlogRepo),
		sprintRepo,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /backlog", backlogHandler.Crear)
	mux.HandleFunc("POST /projects", projectHandler.Crear)
	mux.HandleFunc("GET /projects/{id}", projectHandler.Obtener)
	mux.HandleFunc("GET /projects/{id}/status", projectHandler.ObtenerEstado)
	mux.HandleFunc("PUT /projects/{id}", projectHandler.Editar)
	mux.HandleFunc("POST /projects/{id}/members", projectHandler.AsignarIntegrante)
	mux.HandleFunc("GET /projects/{id}/members", projectHandler.ListarIntegrantes)
	mux.HandleFunc("POST /sprints", sprintHandler.Crear)
	mux.HandleFunc("POST /sprints/{id}/iniciar", sprintHandler.Iniciar)
	mux.HandleFunc("POST /sprints/{id}/cerrar", sprintHandler.Cerrar)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
