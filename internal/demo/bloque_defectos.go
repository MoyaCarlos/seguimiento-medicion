package demo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type defectoDemo struct {
	historia    string
	descripcion string
	severidad   string
	deteccion   int // Sprint 1 a 4, el de la detección
	resolucion  int // Sprint de la resolución; 0 si sigue Pendiente
	dia         int // día relativo a hoy en que se registró
}

var defectosDemo = []defectoDemo{
	{"apto", "El certificado médico que vence el mismo día del trámite se acepta como vigente", "Alta", 1, 2, -24},
	{"turno", "Al cancelar un turno el cupo no se libera y el horario queda como ocupado", "Media", 2, 2, -12},
	{"comprobante", "El código QR del comprobante sale desalineado en la versión para impresora", "Baja", 2, 0, -9},
	{"examen", "El temporizador del examen no se detiene al enviar las respuestas y lo marca como vencido", "Alta", 3, 0, -3},
	{"reprogramar", "El mensaje de confirmación al reprogramar muestra la fecha del turno anterior", "Baja", 3, 3, -2},
}

// bloqueDefectos crea defectos de ejemplo con las reglas del dominio de HU-08.
// Los Resueltos en un Sprint ya finalizado son datos históricos: no se exige
// Sprint activo, igual que con las horas.
func bloqueDefectos(ctx context.Context, tx *sql.Tx, c *carga) error {
	for _, d := range defectosDemo {
		defecto, err := armarDefecto(c, d)
		if err != nil {
			return err
		}
		if err := insertarDefecto(ctx, tx, c, d, defecto); err != nil {
			return err
		}
		c.resumen.Defectos++
	}
	return nil
}

func armarDefecto(c *carga, d defectoDemo) (domain.Defecto, error) {
	deteccion := domain.Sprint{ID: c.sprints[d.deteccion-1]}
	defecto, err := domain.NewDefecto(proyectoID, c.historias[d.historia], deteccion, d.descripcion, d.severidad)
	if err != nil {
		return domain.Defecto{}, err
	}
	if d.resolucion > 0 {
		resolucion := domain.Sprint{ID: c.sprints[d.resolucion-1]}
		if err := defecto.Resolver(&resolucion); err != nil {
			return domain.Defecto{}, err
		}
	}
	return defecto, nil
}

func insertarDefecto(ctx context.Context, tx *sql.Tx, c *carga, d defectoDemo, defecto domain.Defecto) error {
	var resolucion any
	if defecto.SprintResolucionID != nil {
		resolucion = *defecto.SprintResolucionID
	}
	creado := c.dia(d.dia).Add(10 * time.Hour).Format(time.RFC3339)
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id, creado_en)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		defecto.ProyectoID, defecto.HistoriaID, defecto.Descripcion, string(defecto.Severidad), string(defecto.Estado),
		defecto.SprintDeteccionID, resolucion, creado); err != nil {
		return fmt.Errorf("crear defecto de %s: %w", d.historia, err)
	}
	return nil
}
