package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/app"
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

	mux := app.NuevoMux(db, app.Opciones{
		DemoHabilitado: func() bool { return cfg.DemoHabilitado },
		Ahora:          time.Now,
		Contrasenas:    service.NewServicioContrasenasProduccion(),
	})

	log.Printf("listening on :%s", cfg.Puerto)
	if err := http.ListenAndServe(":"+cfg.Puerto, apihttp.ServirSPA(mux, cfg.DirWeb)); err != nil {
		log.Fatal(err)
	}
}
