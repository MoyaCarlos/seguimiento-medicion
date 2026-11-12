package demo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

const (
	proyectoID     int64 = 25
	nombreProyecto       = "Licencia de conducir (datos de prueba)"
	horasPorSP           = 4
)

// bloqueProyecto crea el proyecto con id fijo y su equipo. El INSERT con id
// fijo es la fuente de la idempotencia: falla por clave primaria si ya existe.
func bloqueProyecto(ctx context.Context, tx *sql.Tx, c *carga) error {
	inicio, fin := c.dia(-30), c.dia(45)
	proyecto, err := domain.NewProject(nombreProyecto,
		"Trámite de licencia de conducir: turno, examen teórico, pago y emisión. Proyecto de demostración con datos de ejemplo.",
		&inicio, &fin)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO proyectos (id, nombre, descripcion, fecha_inicio, fecha_fin, horas_por_story_point, creado_en)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		proyectoID, proyecto.Nombre, proyecto.Descripcion,
		inicio.Format(time.DateOnly), fin.Format(time.DateOnly), horasPorSP, c.creadoEn()); err != nil {
		return fmt.Errorf("crear proyecto %d: %w", proyectoID, err)
	}
	c.resumen.ProyectoID = proyectoID

	for _, integrante := range usuariosDemo {
		m, err := domain.NewMembership(proyectoID, c.usuarios[integrante.nombre], integrante.rol)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, ?, ?)`,
			m.ProjectID, m.UserID, string(m.Role), c.creadoEn()); err != nil {
			return fmt.Errorf("agregar integrante %s: %w", integrante.nombre, err)
		}
		c.resumen.Integrantes++
	}
	return nil
}
