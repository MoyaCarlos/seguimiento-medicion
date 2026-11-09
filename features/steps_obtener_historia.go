package features

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"

	apphttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// obtenerHistoriaContext ejercita GET /backlog/{id} sobre una base real en memoria.
type obtenerHistoriaContext struct {
	db  *sql.DB
	mux *http.ServeMux
	rec *httptest.ResponseRecorder
}

func (c *obtenerHistoriaContext) iniciar() error {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(1)
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		return err
	}
	backlog := repository.NewSQLiteBacklogRepository(db)
	handler := apphttp.NewBacklogHandler(nil, nil, service.NewObtenerHistoria(backlog), nil, nil, nil)
	c.db = db
	c.mux = http.NewServeMux()
	c.mux.HandleFunc("GET /backlog/{id}", handler.Obtener)
	c.rec = nil
	return nil
}

func (c *obtenerHistoriaContext) cerrar() error {
	if c.db == nil {
		return nil
	}
	err := c.db.Close()
	c.db = nil
	return err
}

func (c *obtenerHistoriaContext) crearProyecto() error {
	_, err := c.db.Exec(
		`INSERT OR IGNORE INTO proyectos (id, nombre, descripcion, fecha_inicio, fecha_fin, creado_en)
		 VALUES (1, 'Proyecto', '', NULL, NULL, '2026-10-26T00:00:00Z')`)
	return err
}

func (c *obtenerHistoriaContext) crearSprint(id int) error {
	if err := c.crearProyecto(); err != nil {
		return err
	}
	_, err := c.db.Exec(
		`INSERT OR IGNORE INTO sprints (id, proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado)
		 VALUES (?, 1, NULL, NULL, NULL, 'Pendiente')`, id)
	return err
}

func (c *obtenerHistoriaContext) historiaEnBacklog(id int) error {
	if err := c.crearProyecto(); err != nil {
		return err
	}
	_, err := c.db.Exec(
		`INSERT INTO historias (id, proyecto_id, titulo, descripcion, prioridad, estado)
		 VALUES (?, 1, 'Historia', 'Descripción', 'M', 'Nueva')`, id)
	return err
}

func (c *obtenerHistoriaContext) historiaEnSprint(id, sprint int) error {
	if err := c.crearSprint(sprint); err != nil {
		return err
	}
	_, err := c.db.Exec(
		`INSERT INTO historias (id, proyecto_id, titulo, descripcion, prioridad, estado, sprint_id)
		 VALUES (?, 1, 'Historia', 'Descripción', 'M', 'En progreso', ?)`, id, sprint)
	return err
}

func (c *obtenerHistoriaContext) historiaCompletada(id int, fecha string, sprint int) error {
	if err := c.crearSprint(sprint); err != nil {
		return err
	}
	_, err := c.db.Exec(
		`INSERT INTO historias (id, proyecto_id, titulo, descripcion, prioridad, estado, sprint_id, fecha_finalizacion)
		 VALUES (?, 1, 'Historia', 'Descripción', 'M', 'Completada', ?, ?)`, id, sprint, fecha)
	return err
}

func (c *obtenerHistoriaContext) consultoHistoria(id string) error {
	c.rec = httptest.NewRecorder()
	c.mux.ServeHTTP(c.rec, httptest.NewRequest(http.MethodGet, "/backlog/"+id, nil))
	return nil
}

func (c *obtenerHistoriaContext) respondeEstado(esperado int) error {
	if c.rec.Code != esperado {
		return fmt.Errorf("se esperaba el estado %d, se obtuvo %d (%s)", esperado, c.rec.Code, c.rec.Body.String())
	}
	return nil
}

type historiaRecibida struct {
	ID                int64   `json:"id"`
	SprintID          *int64  `json:"sprint_id"`
	FechaFinalizacion *string `json:"fecha_finalizacion"`
}

func (c *obtenerHistoriaContext) decodificar() (historiaRecibida, error) {
	var h historiaRecibida
	if err := json.Unmarshal(c.rec.Body.Bytes(), &h); err != nil {
		return h, fmt.Errorf("respuesta JSON inválida: %w", err)
	}
	return h, nil
}

func (c *obtenerHistoriaContext) historiaConSprintNulo(id int) error {
	h, err := c.decodificar()
	if err != nil {
		return err
	}
	if h.ID != int64(id) || h.SprintID != nil {
		return fmt.Errorf("se esperaba la historia %d con sprint_id nulo, se obtuvo %s", id, c.rec.Body.String())
	}
	return nil
}

func (c *obtenerHistoriaContext) historiaConSprint(id, sprint int) error {
	h, err := c.decodificar()
	if err != nil {
		return err
	}
	if h.ID != int64(id) || h.SprintID == nil || *h.SprintID != int64(sprint) {
		return fmt.Errorf("se esperaba la historia %d con sprint_id %d, se obtuvo %s", id, sprint, c.rec.Body.String())
	}
	return nil
}

func (c *obtenerHistoriaContext) historiaConFecha(fecha string) error {
	h, err := c.decodificar()
	if err != nil {
		return err
	}
	if h.FechaFinalizacion == nil || *h.FechaFinalizacion != fecha {
		return fmt.Errorf("se esperaba fecha_finalizacion %q, se obtuvo %s", fecha, c.rec.Body.String())
	}
	return nil
}

func (c *obtenerHistoriaContext) mensajeEs(mensaje string) error {
	var e struct {
		Mensaje string `json:"mensaje"`
	}
	if err := json.Unmarshal(c.rec.Body.Bytes(), &e); err != nil {
		return fmt.Errorf("respuesta JSON inválida: %w", err)
	}
	if e.Mensaje != mensaje {
		return fmt.Errorf("se esperaba el mensaje %q, se obtuvo %q", mensaje, e.Mensaje)
	}
	return nil
}

func (c *obtenerHistoriaContext) errorEs(campo, mensaje string) error {
	var e struct {
		Campo   string `json:"campo"`
		Mensaje string `json:"mensaje"`
	}
	if err := json.Unmarshal(c.rec.Body.Bytes(), &e); err != nil {
		return fmt.Errorf("respuesta JSON inválida: %w", err)
	}
	if e.Campo != campo || e.Mensaje != mensaje {
		return fmt.Errorf("se esperaba campo %q y mensaje %q, se obtuvo %+v", campo, mensaje, e)
	}
	return nil
}

// InitializeScenarioObtenerHistoria registra los pasos BDD de GET /backlog/{id}.
func InitializeScenarioObtenerHistoria(ctx *godog.ScenarioContext) {
	c := &obtenerHistoriaContext{}

	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		_ = c.cerrar()
		return goctx, c.iniciar()
	})
	ctx.After(func(goctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return goctx, c.cerrar()
	})

	ctx.Step(`^que existe la historia (\d+) en el Product Backlog$`, c.historiaEnBacklog)
	ctx.Step(`^que existe la historia (\d+) asignada al Sprint (\d+)$`, c.historiaEnSprint)
	ctx.Step(`^que existe la historia (\d+) completada el "([^"]*)" en el Sprint (\d+)$`, c.historiaCompletada)
	ctx.Step(`^consulto la historia "([^"]*)"$`, c.consultoHistoria)
	ctx.Step(`^la consulta de la historia responde (\d+)$`, c.respondeEstado)
	ctx.Step(`^la historia recibida es la (\d+) y tiene sprint_id nulo$`, c.historiaConSprintNulo)
	ctx.Step(`^la historia recibida es la (\d+) y tiene sprint_id (\d+)$`, c.historiaConSprint)
	ctx.Step(`^la historia recibida tiene fecha_finalizacion "([^"]*)"$`, c.historiaConFecha)
	ctx.Step(`^el mensaje de la consulta de la historia es "([^"]*)"$`, c.mensajeEs)
	ctx.Step(`^el error de la consulta de la historia es del campo "([^"]*)" con mensaje "([^"]*)"$`, c.errorEs)
}
