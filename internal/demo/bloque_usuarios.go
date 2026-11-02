package demo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

var nombresUsuarios = []string{"Beta Uno", "Beta Dos", "Beta Tres", "Beta Cuatro"}

// bloqueUsuarios asegura los usuarios predeterminados, sin rol. Es la primera
// escritura de la transacción (ver plan.md, decisión 3).
func bloqueUsuarios(ctx context.Context, tx *sql.Tx, c *carga) error {
	for _, nombre := range nombresUsuarios {
		usuario, err := domain.NewUser(nombre)
		if err != nil {
			return err
		}
		resultado, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES (?, ?, ?)`,
			usuario.Nombre, usuario.NombreNormalizado, c.creadoEn())
		if err != nil {
			return fmt.Errorf("crear usuario %s: %w", nombre, err)
		}
		creados, err := resultado.RowsAffected()
		if err != nil {
			return err
		}
		c.resumen.Usuarios += int(creados)

		var id int64
		if err := tx.QueryRowContext(ctx,
			`SELECT id FROM usuarios WHERE nombre_normalizado = ?`, usuario.NombreNormalizado).Scan(&id); err != nil {
			return fmt.Errorf("leer usuario %s: %w", nombre, err)
		}
		c.usuarios[nombre] = id
	}
	return nil
}
