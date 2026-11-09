package features

import (
	"context"
	"errors"
	"fmt"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// scrumMasterUnicoContext reutiliza la base y el proyecto del contexto de HU-04.
type scrumMasterUnicoContext struct {
	pc               *proyectoContext
	integrantesAntes int
	errInsercion     error
}

func (c *scrumMasterUnicoContext) scrumMasters() ([]string, error) {
	filas, err := c.pc.db.Query(
		`SELECT u.nombre FROM integrantes i JOIN usuarios u ON u.id = i.usuario_id
		 WHERE i.proyecto_id = ? AND i.rol = ? ORDER BY u.nombre`,
		c.pc.proyecto.ID, string(domain.RolScrumMaster))
	if err != nil {
		return nil, err
	}
	defer filas.Close()
	var nombres []string
	for filas.Next() {
		var nombre string
		if err := filas.Scan(&nombre); err != nil {
			return nil, err
		}
		nombres = append(nombres, nombre)
	}
	return nombres, filas.Err()
}

func (c *scrumMasterUnicoContext) contarIntegrantes() (int, error) {
	var total int
	err := c.pc.db.QueryRow("SELECT COUNT(*) FROM integrantes WHERE proyecto_id = ?", c.pc.proyecto.ID).Scan(&total)
	return total, err
}

func (c *scrumMasterUnicoContext) unicoScrumMaster() (string, error) {
	nombres, err := c.scrumMasters()
	if err != nil {
		return "", err
	}
	if len(nombres) != 1 {
		return "", fmt.Errorf("se esperaba un único Scrum Master, se encontraron %d: %v", len(nombres), nombres)
	}
	return nombres[0], nil
}

func (c *scrumMasterUnicoContext) tieneUnicoScrumMaster() error {
	_, err := c.unicoScrumMaster()
	return err
}

func (c *scrumMasterUnicoContext) esElCreador() error {
	nombre, err := c.unicoScrumMaster()
	if err != nil {
		return err
	}
	if nombre != "Ana Valentina" {
		return fmt.Errorf("se esperaba como Scrum Master al creador, se obtuvo %q", nombre)
	}
	return nil
}

func (c *scrumMasterUnicoContext) intentoAgregarScrumMaster() error {
	antes, err := c.contarIntegrantes()
	if err != nil {
		return err
	}
	c.integrantesAntes = antes
	_, c.pc.err = c.pc.asignar.Ejecutar(context.Background(), service.AsignarIntegranteInput{
		ProyectoID: c.pc.proyecto.ID,
		Nombre:     "Candela Vulcano",
		Rol:        domain.RolScrumMaster,
	})
	return nil
}

func (c *scrumMasterUnicoContext) rechazaConConflicto() error {
	if !errors.Is(c.pc.err, domain.ErrScrumMasterYaExiste) {
		return fmt.Errorf("se esperaba domain.ErrScrumMasterYaExiste, se obtuvo %v", c.pc.err)
	}
	return nil
}

func (c *scrumMasterUnicoContext) listaSinCambios() error {
	despues, err := c.contarIntegrantes()
	if err != nil {
		return err
	}
	if despues != c.integrantesAntes {
		return fmt.Errorf("la lista de integrantes cambió de %d a %d", c.integrantesAntes, despues)
	}
	return nil
}

func (c *scrumMasterUnicoContext) scrumMasterNoCambia() error {
	return c.tieneUnicoScrumMaster()
}

func (c *scrumMasterUnicoContext) insertarSegundoScrumMaster() error {
	ctx := context.Background()
	usuario, err := domain.NewUser("Candela Vulcano")
	if err != nil {
		return err
	}
	guardado, err := repository.NewSQLiteUserRepository(c.pc.db).Guardar(ctx, usuario)
	if err != nil {
		return err
	}
	_, c.errInsercion = c.pc.db.ExecContext(ctx,
		"INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, ?, ?)",
		c.pc.proyecto.ID, guardado.ID, string(domain.RolScrumMaster), "2026-11-08T00:00:00Z")
	return nil
}

func (c *scrumMasterUnicoContext) baseRechaza() error {
	if c.errInsercion == nil {
		return errors.New("la base aceptó un segundo Scrum Master")
	}
	return nil
}

// InitializeScenarioScrumMasterUnico registra los pasos propios de la spec 039 (#236).
func InitializeScenarioScrumMasterUnico(ctx *godog.ScenarioContext, pc *proyectoContext) {
	c := &scrumMasterUnicoContext{pc: pc}

	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		c.errInsercion = nil
		return goctx, nil
	})

	ctx.Step(`^el proyecto tiene un único Scrum Master$`, c.tieneUnicoScrumMaster)
	ctx.Step(`^ese Scrum Master es quien lo creó$`, c.esElCreador)
	ctx.Step(`^intento agregar a otra persona con el rol "Scrum Master"$`, c.intentoAgregarScrumMaster)
	ctx.Step(`^el sistema rechaza la operación con un conflicto que indica que el proyecto ya tiene Scrum Master$`, c.rechazaConConflicto)
	ctx.Step(`^la lista de integrantes del proyecto no cambia$`, c.listaSinCambios)
	ctx.Step(`^el Scrum Master del proyecto no cambia$`, c.scrumMasterNoCambia)
	ctx.Step(`^se inserta directamente en la base un segundo Scrum Master para ese proyecto$`, c.insertarSegundoScrumMaster)
	ctx.Step(`^la base rechaza la inserción$`, c.baseRechaza)
	ctx.Step(`^el proyecto sigue con su único Scrum Master$`, c.tieneUnicoScrumMaster)
}
