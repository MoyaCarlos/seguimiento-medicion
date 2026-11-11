package demo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

const contrasenaDeLaDemo = "demo1234"

type usuarioDemo struct {
	nombre string
	rol    domain.Role
}

// usuariosDemo va en orden: el conflicto se informa por el primero que ya existe.
var usuariosDemo = []usuarioDemo{
	{"carlos.moya", domain.RolScrumMaster},
	{"ana.iriarte", domain.RolProductBuilder},
	{"candela.vulcano", domain.RolProductBuilder},
	{"jimena.martinez", domain.RolProductBuilder},
}

// bloqueUsuarios crea los usuarios de la demo con el hash único de la carga.
// Es la primera escritura de la transacción (ver plan.md, decisión 1) y no
// reutiliza usuarios: si uno existe, la clave única aborta la carga.
func bloqueUsuarios(ctx context.Context, tx *sql.Tx, c *carga) error {
	for _, demo := range usuariosDemo {
		if err := domain.ValidarNombreUsuario(demo.nombre); err != nil {
			return err
		}
		usuario, err := domain.NewUser(demo.nombre)
		if err != nil {
			return err
		}
		resultado, err := tx.ExecContext(ctx,
			`INSERT INTO usuarios (nombre, nombre_normalizado, contrasena_hash, creado_en) VALUES (?, ?, ?, ?)`,
			usuario.Nombre, usuario.NombreNormalizado, c.hashContrasena, c.creadoEn())
		if err != nil {
			return fmt.Errorf("crear usuario %s: %w", demo.nombre, err)
		}
		id, err := resultado.LastInsertId()
		if err != nil {
			return err
		}
		c.usuarios[demo.nombre] = id
		c.resumen.Usuarios++
	}
	return nil
}
