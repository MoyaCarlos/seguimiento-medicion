package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	apphttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// listadoContext ejercita GET /projects («Mis proyectos», HU-21) sobre una base
// real en memoria, como la persona que inició sesión en el escenario.
type listadoContext struct {
	sesion    *sesionDeEscenario
	db        *sql.DB
	repo      *repository.SQLiteProjectRepository
	servicio  *service.ListarProyectos
	proyectos []domain.ProyectoConRol
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

// integroProyectosConIDs guarda dos proyectos con la persona de la sesión como
// integrante de ambos.
func (c *listadoContext) integroProyectosConIDs(primero, segundo int) error {
	usuario, err := c.usuarioDeLaSesion()
	if err != nil {
		return err
	}
	if _, err := c.db.Exec(
		`INSERT INTO proyectos (id, nombre, descripcion, fecha_inicio, fecha_fin, creado_en) VALUES
			(?, 'Nueve', '', NULL, NULL, '2026-10-23T00:00:00Z'),
			(?, 'Dos', 'Descripción', '2026-03-01', NULL, '2026-10-23T00:00:00Z')`,
		primero, segundo,
	); err != nil {
		return err
	}
	_, err = c.db.Exec(
		`INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES
			(?, ?, 'scrum_master', '2026-10-23T00:00:00Z'),
			(?, ?, 'product_builder', '2026-10-23T00:00:00Z')`,
		primero, usuario, segundo, usuario,
	)
	return err
}

// usuarioDeLaSesion devuelve el id de quien inició sesión, con su usuario en la
// base del contexto.
func (c *listadoContext) usuarioDeLaSesion() (int64, error) {
	usuario, err := domain.NewUser(c.sesion.actor())
	if err != nil {
		return 0, fmt.Errorf("el escenario no inició sesión: %w", err)
	}
	if _, err := c.db.Exec(`INSERT OR IGNORE INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES (?, ?, '2026-10-23T00:00:00Z')`,
		usuario.Nombre, usuario.NombreNormalizado); err != nil {
		return 0, err
	}
	var id int64
	err = c.db.QueryRow(`SELECT id FROM usuarios WHERE nombre_normalizado = ?`, usuario.NombreNormalizado).Scan(&id)
	return id, err
}

func (c *listadoContext) sinProyectos() error { return nil }

func (c *listadoContext) persistenciaFallida() error {
	c.falla = true
	return nil
}

func (c *listadoContext) consultoListado() error {
	if c.falla {
		return c.consultoListadoConPersistenciaFallida()
	}
	usuario, err := c.usuarioDeLaSesion()
	if err != nil {
		return err
	}
	c.proyectos, c.err = c.servicio.Ejecutar(context.Background(), usuario)
	return nil
}

func (c *listadoContext) proyectosConIDsEnOrden(primero, segundo int) error {
	if c.err != nil {
		return fmt.Errorf("no se esperaba error al listar: %w", c.err)
	}
	if len(c.proyectos) != 2 {
		return fmt.Errorf("se esperaban 2 proyectos, se obtuvieron %d", len(c.proyectos))
	}
	if c.proyectos[0].Proyecto.ID != int64(primero) || c.proyectos[1].Proyecto.ID != int64(segundo) {
		return fmt.Errorf(
			"se esperaba el orden [%d, %d], se obtuvo [%d, %d]",
			primero, segundo, c.proyectos[0].Proyecto.ID, c.proyectos[1].Proyecto.ID,
		)
	}
	return nil
}

func (c *listadoContext) mismosCamposQueConsultaIndividual() error {
	for _, conRol := range c.proyectos {
		listado := conRol.Proyecto
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

func (repositorioQueFalla) ListarPorIntegrante(context.Context, int64) ([]domain.ProyectoConRol, error) {
	return nil, errors.New("fallo de persistencia")
}

// consultoListadoConPersistenciaFallida pide GET /projects con una sesión real
// (emitida en la base del escenario) a un handler cuyo repositorio falla.
func (c *listadoContext) consultoListadoConPersistenciaFallida() error {
	usuario, err := c.usuarioDeLaSesion()
	if err != nil {
		return err
	}
	sesiones := repository.NewSQLiteSesionRepository(c.db)
	token, _, err := service.NewEmisorSesiones(sesiones, time.Now).Emitir(context.Background(), usuario)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	handler := apphttp.NewProjectHandler(nil, nil, service.NewListarProyectos(repositorioQueFalla{}), nil, nil, nil, nil, nil)
	apphttp.RegistrarRutas(mux, apphttp.Handlers{Proyectos: handler}, apphttp.ConSesion(service.NewObtenerSesion(sesiones, time.Now)))
	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	c.status = rec.Code
	return nil
}

// InitializeScenarioListado registra los pasos BDD del listado de proyectos.
func InitializeScenarioListado(ctx *godog.ScenarioContext, sesion *sesionDeEscenario) {
	c := &listadoContext{sesion: sesion}

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

	ctx.Step(`^que integro proyectos guardados con identificadores (\d+) y (\d+)$`, c.integroProyectosConIDs)
	ctx.Step(`^que no hay proyectos guardados$`, c.sinProyectos)
	ctx.Step(`^que la persistencia de proyectos falla$`, c.persistenciaFallida)
	ctx.Step(`^consulto el listado de proyectos$`, c.consultoListado)
	ctx.Step(`^recibo los proyectos con identificadores (\d+) y (\d+) en ese orden$`, c.proyectosConIDsEnOrden)
	ctx.Step(`^cada proyecto del listado tiene los mismos campos que su consulta individual$`, c.mismosCamposQueConsultaIndividual)
	ctx.Step(`^recibo un arreglo vacío$`, c.arregloVacio)
	ctx.Step(`^el sistema responde con un error interno (\d+)$`, c.respondeErrorInterno)
}
