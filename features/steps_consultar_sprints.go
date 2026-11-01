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

// consultaSprintsContext ejercita GET /projects/{id}/sprints sobre una base real en memoria.
type consultaSprintsContext struct {
	db  *sql.DB
	mux *http.ServeMux
	rec *httptest.ResponseRecorder
}

func (c *consultaSprintsContext) iniciar() error {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(1)
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		return err
	}
	proyectos := repository.NewSQLiteProjectRepository(db)
	sprints := repository.NewSQLiteSprintRepository(db)
	handler := apphttp.NewSprintHandler(nil, nil, nil, service.NewConsultarSprints(proyectos, sprints), nil, sprints)
	c.db = db
	c.mux = http.NewServeMux()
	c.mux.HandleFunc("GET /projects/{id}/sprints", handler.Listar)
	c.rec = nil
	return nil
}

func (c *consultaSprintsContext) cerrar() error {
	if c.db == nil {
		return nil
	}
	err := c.db.Close()
	c.db = nil
	return err
}

func (c *consultaSprintsContext) crearProyecto(id int) error {
	_, err := c.db.Exec(
		`INSERT OR IGNORE INTO proyectos (id, nombre, descripcion, fecha_inicio, fecha_fin, creado_en)
		 VALUES (?, 'Proyecto', '', NULL, NULL, '2026-10-26T00:00:00Z')`, id)
	return err
}

func (c *consultaSprintsContext) proyectoConSprints(proyecto, primero, segundo int) error {
	if err := c.crearProyecto(proyecto); err != nil {
		return err
	}
	_, err := c.db.Exec(
		`INSERT INTO sprints (id, proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES
			(?, ?, NULL, NULL, NULL, 'Pendiente'),
			(?, ?, 'Entregar el MVP', '2026-10-05', '2026-10-19', 'Finalizado')`,
		primero, proyecto, segundo, proyecto)
	return err
}

func (c *consultaSprintsContext) proyectoSinSprints(proyecto int) error {
	return c.crearProyecto(proyecto)
}

func (c *consultaSprintsContext) consultoSprintsDe(id string) error {
	c.rec = httptest.NewRecorder()
	c.mux.ServeHTTP(c.rec, httptest.NewRequest(http.MethodGet, "/projects/"+id+"/sprints", nil))
	return nil
}

func (c *consultaSprintsContext) consultoSprintsDeNumero(id int) error {
	return c.consultoSprintsDe(fmt.Sprint(id))
}

func (c *consultaSprintsContext) respondeEstado(esperado int) error {
	if c.rec.Code != esperado {
		return fmt.Errorf("se esperaba el estado %d, se obtuvo %d (%s)", esperado, c.rec.Code, c.rec.Body.String())
	}
	return nil
}

func (c *consultaSprintsContext) recibeSprintsEnOrden(primero, segundo int) error {
	var sprints []struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(c.rec.Body.Bytes(), &sprints); err != nil {
		return fmt.Errorf("respuesta JSON inválida: %w", err)
	}
	if len(sprints) != 2 || sprints[0].ID != primero || sprints[1].ID != segundo {
		return fmt.Errorf("se esperaba el orden [%d, %d], se obtuvo %+v", primero, segundo, sprints)
	}
	return nil
}

func (c *consultaSprintsContext) recibeArregloVacio() error {
	if cuerpo := c.rec.Body.String(); cuerpo != "[]\n" {
		return fmt.Errorf("se esperaba [], se obtuvo %q", cuerpo)
	}
	return nil
}

// InitializeScenarioConsultarSprints registra los pasos BDD de la consulta de Sprints.
func InitializeScenarioConsultarSprints(ctx *godog.ScenarioContext) {
	c := &consultaSprintsContext{}

	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		_ = c.cerrar()
		return goctx, c.iniciar()
	})
	ctx.After(func(goctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return goctx, c.cerrar()
	})

	ctx.Step(`^que el proyecto (\d+) tiene los Sprints (\d+) y (\d+)$`, c.proyectoConSprints)
	ctx.Step(`^que el proyecto (\d+) no tiene Sprints$`, c.proyectoSinSprints)
	ctx.Step(`^consulto los Sprints del proyecto (\d+)$`, c.consultoSprintsDeNumero)
	ctx.Step(`^consulto los Sprints del proyecto "([^"]*)"$`, c.consultoSprintsDe)
	ctx.Step(`^el sistema responde (\d+)$`, c.respondeEstado)
	ctx.Step(`^recibo los Sprints (\d+) y (\d+) en ese orden$`, c.recibeSprintsEnOrden)
	ctx.Step(`^recibo un arreglo de Sprints vacío$`, c.recibeArregloVacio)
}
