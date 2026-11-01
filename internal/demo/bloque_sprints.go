package demo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// bloqueSprints crea el Sprint 1 Finalizado, el 2 Activo y el 3 Pendiente. El
// Pendiente no lleva Sprint Goal ni fechas: lo exige el CHECK de la tabla.
func bloqueSprints(ctx context.Context, tx *sql.Tx, c *carga) error {
	finalizado, err := sprintIniciado("Gestión inicial del trámite", c.dia(-14), c.dia(-8))
	if err != nil {
		return err
	}
	if err := finalizado.Cerrar(); err != nil {
		return err
	}
	activo, err := sprintIniciado("Examen teórico y comprobantes", c.dia(-7), c.dia(6))
	if err != nil {
		return err
	}
	pendiente, err := domain.NewSprint(proyectoID)
	if err != nil {
		return err
	}

	for _, sprint := range []domain.Sprint{finalizado, activo, pendiente} {
		id, err := insertarSprint(ctx, tx, sprint)
		if err != nil {
			return err
		}
		c.sprints = append(c.sprints, id)
		c.resumen.Sprints++
	}
	return nil
}

func sprintIniciado(goal string, inicio, fin time.Time) (domain.Sprint, error) {
	sprint, err := domain.NewSprint(proyectoID)
	if err != nil {
		return domain.Sprint{}, err
	}
	return sprint, sprint.Iniciar(goal, inicio, fin)
}

func insertarSprint(ctx context.Context, tx *sql.Tx, s domain.Sprint) (int64, error) {
	var goal, inicio, fin any
	if s.Estado != domain.SprintPendiente {
		goal, inicio, fin = s.SprintGoal, s.FechaInicio.Format(time.DateOnly), s.FechaFin.Format(time.DateOnly)
	}
	resultado, err := tx.ExecContext(ctx,
		`INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (?, ?, ?, ?, ?)`,
		s.ProyectoID, goal, inicio, fin, string(s.Estado))
	if err != nil {
		return 0, fmt.Errorf("crear Sprint %s: %w", s.Estado, err)
	}
	return resultado.LastInsertId()
}
