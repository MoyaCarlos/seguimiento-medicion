package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	apphttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// listadoContext ejercita GET /projects sobre una base real en memoria.
type listadoContext struct {
	db        *sql.DB
	repo      *repository.SQLiteProjectRepository
	servicio  *service.ListarProyectos
	proyectos []domain.Project
	status    int
	falla     bool
	err       error
}

func (c *listadoContext) iniciar() error {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(1)
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		return err
	}
	c.db = db
	c.repo = repository.NewSQLiteProjectRepository(db)
	c.servicio = service.NewListarProyectos(c.repo)
	c.proyectos = nil
	c.status = 0
	c.falla = false
	c.err = nil
	return nil
}

func (c *listadoContext) cerrar() error {
	if c.db == nil {
		return nil
	}
	err := c.db.Close()
	c.db = nil
	return err
}

func (c *listadoContext) existenProyectosConIDs(primero, segundo int) error {
	_, err := c.db.Exec(
		`INSERT INTO proyectos (id, nombre, descripcion, fecha_inicio, fecha_fin, creado_en) VALUES
			(?, 'Nueve', '', NULL, NULL, '2026-10-23T00:00:00Z'),
			(?, 'Dos', 'Descripción', '2026-03-01', NULL, '2026-10-23T00:00:00Z')`,
		primero, segundo,
	)
	return err
}

func (c *listadoContext) sinProyectos() error { return nil }

func (c *listadoContext) persistenciaFallida() error {
	c.falla = true
	return nil
}

func (c *listadoContext) consultoListado() error {
	if c.falla {
		handler := apphttp.NewProjectHandler(nil, nil, service.NewListarProyectos(repositorioQueFalla{}), nil, nil, nil, nil)
		rec := httptest.NewRecorder()
		handler.Listar(rec, httptest.NewRequest(http.MethodGet, "/projects", nil))
		c.status = rec.Code
		return nil
	}
	c.proyectos, c.err = c.servicio.Ejecutar(context.Background())
	return nil
}

func (c *listadoContext) proyectosConIDsEnOrden(primero, segundo int) error {
	if c.err != nil {
		return fmt.Errorf("no se esperaba error al listar: %w", c.err)
	}
	if len(c.proyectos) != 2 {
		return fmt.Errorf("se esperaban 2 proyectos, se obtuvieron %d", len(c.proyectos))
	}
	if c.proyectos[0].ID != int64(primero) || c.proyectos[1].ID != int64(segundo) {
		return fmt.Errorf(
			"se esperaba el orden [%d, %d], se obtuvo [%d, %d]",
			primero, segundo, c.proyectos[0].ID, c.proyectos[1].ID,
		)
	}
	return nil
}

func (c *listadoContext) mismosCamposQueConsultaIndividual() error {
	for _, listado := range c.proyectos {
		individual, err := c.repo.ObtenerPorID(context.Background(), listado.ID)
		if err != nil {
			return fmt.Errorf("no se pudo consultar el proyecto %d: %w", listado.ID, err)
		}
		if !reflect.DeepEqual(listado, individual) {
			return fmt.Errorf(
				"el proyecto %d del listado difiere de su consulta individual: %+v vs %+v",
				listado.ID, listado, individual,
			)
		}
	}
	return nil
}

func (c *listadoContext) arregloVacio() error {
	if c.err != nil {
		return fmt.Errorf("no se esperaba error: %w", c.err)
	}
	if c.proyectos == nil {
		return errors.New("se esperaba un arreglo vacío, se obtuvo nil")
	}
	if len(c.proyectos) != 0 {
		return fmt.Errorf("se esperaba un arreglo vacío, se obtuvo %+v", c.proyectos)
	}
	return nil
}

func (c *listadoContext) respondeErrorInterno(esperado int) error {
	if c.status != esperado {
		return fmt.Errorf("se esperaba el estado %d, se obtuvo %d", esperado, c.status)
	}
	return nil
}

type repositorioQueFalla struct{}

func (repositorioQueFalla) Listar(context.Context) ([]domain.Project, error) {
	return nil, errors.New("fallo de persistencia")
}

// InitializeScenarioListado registra los pasos BDD del listado de proyectos.
func InitializeScenarioListado(ctx *godog.ScenarioContext) {
	c := &listadoContext{}

	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		if c.db != nil {
			_ = c.db.Close()
			c.db = nil
		}
		return goctx, c.iniciar()
	})
	ctx.After(func(goctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		return goctx, c.cerrar()
	})

	ctx.Step(`^que existen proyectos guardados con identificadores (\d+) y (\d+)$`, c.existenProyectosConIDs)
	ctx.Step(`^que no hay proyectos guardados$`, c.sinProyectos)
	ctx.Step(`^que la persistencia de proyectos falla$`, c.persistenciaFallida)
	ctx.Step(`^consulto el listado de proyectos$`, c.consultoListado)
	ctx.Step(`^recibo los proyectos con identificadores (\d+) y (\d+) en ese orden$`, c.proyectosConIDsEnOrden)
	ctx.Step(`^cada proyecto del listado tiene los mismos campos que su consulta individual$`, c.mismosCamposQueConsultaIndividual)
	ctx.Step(`^recibo un arreglo vacío$`, c.arregloVacio)
	ctx.Step(`^el sistema responde con un error interno (\d+)$`, c.respondeErrorInterno)
}
