package demo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type horasDemo struct {
	historia  string
	sprint    int // 1 o 2
	usuario   string
	dia       int // relativo a hoy
	horas     string
	actividad string
}

var horasImputadas = []horasDemo{
	{"turno", 1, "Beta Dos", -13, "6", "Diseño de la agenda de turnos"},
	{"turno", 1, "Beta Tres", -12, "8", "Implementación del cupo por turno"},
	{"turno", 1, "Beta Uno", -11, "4", "Revisión y pruebas de aceptación"},
	{"datos", 1, "Beta Tres", -13, "5", "Formulario de datos personales"},
	{"datos", 1, "Beta Dos", -12, "6", "Validación del DNI y del domicilio"},
	{"apto", 1, "Beta Dos", -11, "10", "Lectura y carga del certificado"},
	{"apto", 1, "Beta Uno", -10, "8", "Reglas de vigencia del certificado"},
	{"apto", 1, "Beta Tres", -9, "12", "Pruebas con certificados vencidos"},
	{"tasa", 1, "Beta Uno", -10, "3", "Análisis de la pasarela de pagos"},
	{"tasa", 1, "Beta Tres", -9, "4", "Maqueta de la pantalla de pago"},
	{"examen", 2, "Beta Dos", -6, "8", "Banco de preguntas y sorteo"},
	{"examen", 2, "Beta Tres", -4, "7", "Pantalla del examen con temporizador"},
	{"examen", 2, "Beta Uno", -2, "5", "Corrección automática del examen"},
	{"comprobante", 2, "Beta Tres", -3, "3", "Plantilla del comprobante"},
	{"comprobante", 2, "Beta Dos", -1, "2", "Generación del código QR"},
}

// bloqueHoras imputa horas históricas en los Sprints 1 y 2. No exige que el
// Sprint esté Activo (regla de la API de HU-07): son datos de demostración.
func bloqueHoras(ctx context.Context, tx *sql.Tx, c *carga) error {
	sprints := []domain.Sprint{
		{ID: c.sprints[0], ProyectoID: proyectoID, FechaInicio: c.dia(-14), FechaFin: c.dia(-8)},
		{ID: c.sprints[1], ProyectoID: proyectoID, FechaInicio: c.dia(-7), FechaFin: c.dia(6)},
	}
	for _, h := range horasImputadas {
		historia := domain.BacklogItem{ID: c.historias[h.historia], ProyectoID: proyectoID}
		registro, err := domain.NewRegistroHoras(historia, sprints[h.sprint-1], c.usuarios[h.usuario],
			h.horas, h.actividad, c.dia(h.dia).Format(time.DateOnly), c.ahora)
		if err != nil {
			return fmt.Errorf("horas de %s: %w", h.historia, err)
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO registros_horas (proyecto_id, historia_id, usuario_id, sprint_id, fecha, actividad, minutos, creado_en)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			registro.ProyectoID, registro.HistoriaID, registro.UsuarioID, registro.SprintID,
			registro.Fecha.Format(time.DateOnly), registro.Actividad, registro.Minutos, c.creadoEn()); err != nil {
			return fmt.Errorf("imputar horas de %s: %w", h.historia, err)
		}
		c.resumen.RegistrosHoras++
	}
	return nil
}
