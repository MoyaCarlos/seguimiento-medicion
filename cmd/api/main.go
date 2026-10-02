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

	repo := repository.NewSQLiteBacklogRepository(db)
	crearHistoria := service.NewCrearHistoriaBacklog(repo)
	backlogHandler := apihttp.NewBacklogHandler(crearHistoria)

	sprintRepo := repository.NewSQLiteSprintRepository(db)
	sprintHandler := apihttp.NewSprintHandler(
		service.NewIniciarSprint(sprintRepo),
		service.NewCerrarSprint(sprintRepo, repo),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /backlog", backlogHandler.Crear)
	mux.HandleFunc("POST /sprints/{id}/iniciar", sprintHandler.Iniciar)
	mux.HandleFunc("POST /sprints/{id}/cerrar", sprintHandler.Cerrar)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
