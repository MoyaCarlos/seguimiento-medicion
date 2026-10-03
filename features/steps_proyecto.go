package features

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type proyectoContext struct {
	db       *sql.DB
	crear    *service.CrearProyecto
	obtener  *service.ObtenerProyecto
	editar   *service.EditarProyecto
	asignar  *service.AsignarIntegrante
	proyecto domain.Project
	err      error
}

func (c *proyectoContext) iniciar() error {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return err
	}
	if err := repository.Migrar(context.Background(), db); err != nil {
		return err
	}
	c.db = db

	proyectos := repository.NewSQLiteProjectRepository(db)
	usuarios := repository.NewSQLiteUserRepository(db)
	c.crear = service.NewCrearProyecto(proyectos, usuarios)
	c.obtener = service.NewObtenerProyecto(proyectos)
	c.editar = service.NewEditarProyecto(proyectos)
	c.asignar = service.NewAsignarIntegrante(proyectos, usuarios)
	c.proyecto = domain.Project{}
	c.err = nil
	return nil
}

func (c *proyectoContext) enPantallaPrincipal() error { return nil }

func (c *proyectoContext) crearProyectoValido(nombre string) error {
	c.proyecto, c.err = c.crear.Ejecutar(context.Background(), service.CrearProyectoInput{
		Nombre:      nombre,
		Descripcion: "Sistema de estimación, seguimiento y medición",
		Creador:     "Ana Valentina",
	})
	return nil
}

func (c *proyectoContext) proyectoPersistido() error {
	if c.err != nil {
		return c.err
	}
	if c.proyecto.ID == 0 {
		return errors.New("no se generó el proyecto")
	}
	var total int
	if err := c.db.QueryRow("SELECT COUNT(*) FROM projects WHERE id = ?", c.proyecto.ID).Scan(&total); err != nil {
		return err
	}
	if total != 1 {
		return errors.New("el proyecto no quedó persistido en SQLite")
	}
	return nil
}

func (c *proyectoContext) confirmaProyecto() error {
	if c.proyecto.ID == 0 {
		return errors.New("no se devolvió el identificador del proyecto")
	}
	return nil
}

func (c *proyectoContext) habilitaPanel() error { return nil }

func (c *proyectoContext) intentoCrearProyecto() error { return nil }

func (c *proyectoContext) crearProyectoNombreVacio() error {
	c.proyecto, c.err = c.crear.Ejecutar(context.Background(), service.CrearProyectoInput{
		Nombre:  "   ",
		Creador: "Ana",
	})
	return nil
}

func (c *proyectoContext) muestraAdvertencia() error {
	var verr domain.ValidationError
	if !errors.As(c.err, &verr) {
		return fmt.Errorf("se esperaba una advertencia de validación, se obtuvo %v", c.err)
	}
	return nil
}

func (c *proyectoContext) noGeneraProyecto() error {
	var total int
	if err := c.db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&total); err != nil {
		return err
	}
	if total != 0 {
		return fmt.Errorf("se esperaban 0 proyectos, se encontraron %d", total)
	}
	return nil
}

func (c *proyectoContext) existeProyectoCreado() error {
	proyecto, err := c.crear.Ejecutar(context.Background(), service.CrearProyectoInput{
		Nombre:      "Proyecto base",
		Descripcion: "descripción base",
		Creador:     "Scrum Master",
	})
	if err != nil {
		return err
	}
	c.proyecto = proyecto
	return nil
}

func (c *proyectoContext) enConfiguracion() error { return nil }

func (c *proyectoContext) asignarProductBuilder() error {
	_, c.err = c.asignar.Ejecutar(context.Background(), service.AsignarIntegranteInput{
		ProyectoID: c.proyecto.ID,
		Nombre:     "Jimena Martinez",
		Rol:        domain.RolProductBuilder,
	})
	return nil
}

func (c *proyectoContext) vinculaIntegrante() error {
	if c.err != nil {
		return c.err
	}
	var total int
	if err := c.db.QueryRow(
		"SELECT COUNT(*) FROM project_members WHERE project_id = ? AND role = ?",
		c.proyecto.ID, string(domain.RolProductBuilder),
	).Scan(&total); err != nil {
		return err
	}
	if total < 1 {
		return errors.New("no se vinculó al integrante")
	}
	return nil
}

func (c *proyectoContext) habilitaPermisos() error { return nil }

func (c *proyectoContext) asignarInvalido() error {
	_, errNombre := c.asignar.Ejecutar(context.Background(), service.AsignarIntegranteInput{
		ProyectoID: c.proyecto.ID,
		Nombre:     "",
		Rol:        domain.RolProductBuilder,
	})
	_, errRol := c.asignar.Ejecutar(context.Background(), service.AsignarIntegranteInput{
		ProyectoID: c.proyecto.ID,
		Nombre:     "Jimena Martinez",
		Rol:        domain.Role("Product Owner"),
	})
	var vNombre, vRol domain.ValidationError
	if !errors.As(errNombre, &vNombre) || vNombre.Campo != "nombre" {
		return fmt.Errorf("se esperaba ValidationError de nombre vacío, se obtuvo %v", errNombre)
	}
	if !errors.As(errRol, &vRol) || vRol.Campo != "rol" {
		return fmt.Errorf("se esperaba ValidationError de rol inválido, se obtuvo %v", errRol)
	}
	c.err = errNombre
	return nil
}

func (c *proyectoContext) noRealizaVinculacion() error {
	if err := c.muestraAdvertencia(); err != nil {
		return err
	}
	var total int
	if err := c.db.QueryRow(
		"SELECT COUNT(*) FROM project_members WHERE role = ?",
		string(domain.RolProductBuilder),
	).Scan(&total); err != nil {
		return err
	}
	if total != 0 {
		return fmt.Errorf("se escribieron %d vinculaciones de product_builder pese al error", total)
	}
	return nil
}

func (c *proyectoContext) enEdicion() error { return nil }

func (c *proyectoContext) editarValido() error {
	inicio := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 12, 15, 0, 0, 0, 0, time.UTC)
	_, c.err = c.editar.Ejecutar(context.Background(), service.EditarProyectoInput{
		ID:          c.proyecto.ID,
		Nombre:      "Proyecto corregido",
		Descripcion: "descripción corregida",
		FechaInicio: &inicio,
		FechaFin:    &fin,
	})
	return nil
}

func (c *proyectoContext) persisteCambios() error {
	if c.err != nil {
		return c.err
	}
	actualizado, err := c.obtener.Ejecutar(context.Background(), c.proyecto.ID)
	if err != nil {
		return err
	}
	if actualizado.Name != "Proyecto corregido" {
		return fmt.Errorf("no se persistió la edición, nombre actual %q", actualizado.Name)
	}
	if actualizado.Description != "descripción corregida" {
		return fmt.Errorf("no se persistió la descripción: %q", actualizado.Description)
	}
	if actualizado.StartDate == nil || actualizado.StartDate.Format("2006-01-02") != "2026-03-02" {
		return fmt.Errorf("no se persistió la fecha de inicio: %v", actualizado.StartDate)
	}
	if actualizado.EndDate == nil || actualizado.EndDate.Format("2006-01-02") != "2026-12-15" {
		return fmt.Errorf("no se persistió la fecha de fin: %v", actualizado.EndDate)
	}
	return nil
}

func (c *proyectoContext) editarNombreVacio() error {
	_, c.err = c.editar.Ejecutar(context.Background(), service.EditarProyectoInput{
		ID:     c.proyecto.ID,
		Nombre: "   ",
	})
	return nil
}

func (c *proyectoContext) editarFechaInvalida() error {
	inicio := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	_, c.err = c.editar.Ejecutar(context.Background(), service.EditarProyectoInput{
		ID:          c.proyecto.ID,
		Nombre:      "Proyecto base",
		FechaInicio: &inicio,
		FechaFin:    &fin,
	})
	return nil
}

func (c *proyectoContext) noGuardaCambios() error {
	if err := c.muestraAdvertencia(); err != nil {
		return err
	}
	return c.permaneceSinModificaciones()
}

func (c *proyectoContext) cancelaEdicion() error { return nil }

func (c *proyectoContext) permaneceSinModificaciones() error {
	var nombre string
	if err := c.db.QueryRow("SELECT name FROM projects WHERE id = ?", c.proyecto.ID).Scan(&nombre); err != nil {
		return err
	}
	if nombre != "Proyecto base" {
		return fmt.Errorf("el proyecto cambió a %q", nombre)
	}
	return nil
}

// InitializeScenarioProyecto registra los pasos de los escenarios de HU-04.
func InitializeScenarioProyecto(ctx *godog.ScenarioContext) {
	pc := &proyectoContext{}

	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		if pc.db != nil {
			_ = pc.db.Close()
			pc.db = nil
		}
		return goctx, pc.iniciar()
	})
	ctx.After(func(goctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if pc.db != nil {
			_ = pc.db.Close()
			pc.db = nil
		}
		return goctx, nil
	})

	ctx.Step(`^que me encuentro en la pantalla principal del sistema$`, pc.enPantallaPrincipal)
	ctx.Step(`^ingreso el nombre "([^"]*)", una descripción y presiono "Crear"$`, pc.crearProyectoValido)
	ctx.Step(`^el sistema genera el proyecto de forma persistente$`, pc.proyectoPersistido)
	ctx.Step(`^confirma su creación devolviendo el proyecto con su identificador$`, pc.confirmaProyecto)
	ctx.Step(`^habilita su panel principal, cuya navegación realiza el frontend$`, pc.habilitaPanel)
	ctx.Step(`^que intento crear un proyecto$`, pc.intentoCrearProyecto)
	ctx.Step(`^dejo el nombre vacío o compuesto solo por espacios y presiono "Crear"$`, pc.crearProyectoNombreVacio)
	ctx.Step(`^el sistema muestra una advertencia de validación del proyecto$`, pc.muestraAdvertencia)
	ctx.Step(`^no genera el proyecto$`, pc.noGeneraProyecto)
	ctx.Step(`^que existe un proyecto creado$`, pc.existeProyectoCreado)
	ctx.Step(`^estoy en la vista de configuración del proyecto$`, pc.enConfiguracion)
	ctx.Step(`^ingreso el nombre de un integrante, selecciono su rol "Product Builder" y confirmo$`, pc.asignarProductBuilder)
	ctx.Step(`^el sistema vincula al integrante con el proyecto$`, pc.vinculaIntegrante)
	ctx.Step(`^le habilita los permisos correspondientes a su rol$`, pc.habilitaPermisos)
	ctx.Step(`^ingreso un integrante sin nombre o con un rol inválido y confirmo$`, pc.asignarInvalido)
	ctx.Step(`^no realiza la vinculación$`, pc.noRealizaVinculacion)
	ctx.Step(`^estoy en la vista de edición del proyecto$`, pc.enEdicion)
	ctx.Step(`^modifico su nombre, su descripción o sus fechas con datos válidos y guardo$`, pc.editarValido)
	ctx.Step(`^el sistema persiste los cambios del proyecto$`, pc.persisteCambios)
	ctx.Step(`^dejo el nombre vacío o compuesto solo por espacios y guardo$`, pc.editarNombreVacio)
	ctx.Step(`^no guarda los cambios$`, pc.noGuardaCambios)
	ctx.Step(`^ingreso una fecha de fin anterior a la fecha de inicio y guardo$`, pc.editarFechaInvalida)
	ctx.Step(`^decido no aplicar cambios y cancelo$`, pc.cancelaEdicion)
	ctx.Step(`^el proyecto permanece sin modificaciones$`, pc.permaneceSinModificaciones)
}
