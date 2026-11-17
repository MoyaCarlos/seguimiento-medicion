package features

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/app"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// Las lecturas de datos de un proyecto exigen sesión y ser integrante: los
// escenarios que consultan la API arman un lector con sesión y lo suman al
// proyecto que consultan.

// apiDeEscenario arma la API completa sobre db, igual que producción.
func apiDeEscenario(db *sql.DB) *http.ServeMux {
	return app.NuevoMux(db, app.Opciones{
		DemoHabilitado: func() bool { return false },
		Ahora:          time.Now,
		Contrasenas:    service.NewServicioContrasenas(1),
	})
}

// lectorDeEscenario es la persona con sesión que consulta la API en un escenario.
type lectorDeEscenario struct {
	db    *sql.DB
	id    int64
	token string
}

func nuevoLector(db *sql.DB) (lectorDeEscenario, error) {
	if _, err := db.Exec(`INSERT OR IGNORE INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES ('lector', 'lector', '2026-10-26T00:00:00Z')`); err != nil {
		return lectorDeEscenario{}, err
	}
	var id int64
	if err := db.QueryRow(`SELECT id FROM usuarios WHERE nombre_normalizado = 'lector'`).Scan(&id); err != nil {
		return lectorDeEscenario{}, err
	}
	emisor := service.NewEmisorSesiones(repository.NewSQLiteSesionRepository(db), time.Now)
	token, _, err := emisor.Emitir(context.Background(), id)
	return lectorDeEscenario{db: db, id: id, token: token}, err
}

// integra suma al lector como Product Builder del proyecto.
func (l lectorDeEscenario) integra(proyectoID int64) error {
	_, err := l.db.Exec(
		`INSERT OR IGNORE INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, 'product_builder', '2026-10-26T00:00:00Z')`,
		proyectoID, l.id)
	return err
}

// pedir agrega la sesión del lector a la petición.
func (l lectorDeEscenario) pedir(req *http.Request) *http.Request {
	req.Header.Set("Authorization", "Bearer "+l.token)
	return req
}
