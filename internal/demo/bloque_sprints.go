package demo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// sprintDemo define un Sprint con sus días relativos a hoy. El Pendiente no
// lleva Sprint Goal ni fechas: lo exige el CHECK de la tabla.
type sprintDemo struct {
	goal         string
	desde, hasta int
	estado       domain.EstadoSprint
}

var sprintsDemo = []sprintDemo{
	{"Alta del trámite y validación médica", -28, -15, domain.SprintFinalizado},
	{"Turno y comprobante", -14, -8, domain.SprintFinalizado},
	{"Examen teórico y reprogramación de turnos", -7, 6, domain.SprintActivo},
	{estado: domain.SprintPendiente},
}

// bloqueSprints crea los dos Sprints Finalizados, el Activo y el Pendiente.
func bloqueSprints(ctx context.Context, tx *sql.Tx, c *carga) error {
	for _, definicion := range sprintsDemo {
		sprint, err := armarSprint(c, definicion)
		if err != nil {
			return err
		}
		id, err := insertarSprint(ctx, tx, sprint)
		if err != nil {
			return err
		}
		c.sprints = append(c.sprints, id)
		c.resumen.Sprints++
	}
	return nil
}

func armarSprint(c *carga, definicion sprintDemo) (domain.Sprint, error) {
	if definicion.estado == domain.SprintPendiente {
		return domain.NewSprint(proyectoID)
	}
	sprint, err := sprintIniciado(definicion.goal, c.dia(definicion.desde), c.dia(definicion.hasta))
	if err != nil {
		return domain.Sprint{}, err
	}
	if definicion.estado == domain.SprintFinalizado {
		return sprint, sprint.Cerrar()
	}
	return sprint, nil
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
