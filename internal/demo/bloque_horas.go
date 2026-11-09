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
	sprint    int // 1 a 3
	usuario   string
	dia       int // relativo a hoy
	horas     string
	actividad string
}

var horasImputadas = []horasDemo{
	{"tipos", 1, "ana.iriarte", -27, "5", "Relevamiento de los tipos de licencia"},
	{"datos", 1, "candela.vulcano", -26, "5", "Formulario de datos personales"},
	{"datos", 1, "ana.iriarte", -25, "6", "Validación del DNI y del domicilio"},
	{"datos", 1, "jimena.martinez", -25, "3", "Revisión de los textos del formulario"},
	{"apto", 1, "ana.iriarte", -24, "10", "Lectura y carga del certificado"},
	{"apto", 1, "carlos.moya", -23, "8", "Reglas de vigencia del certificado"},
	{"apto", 1, "candela.vulcano", -22, "12", "Pruebas con certificados vencidos"},
	{"turno", 2, "ana.iriarte", -13, "6", "Diseño de la agenda de turnos"},
	{"turno", 2, "candela.vulcano", -12, "8", "Implementación del cupo por turno"},
	{"turno", 2, "carlos.moya", -11, "4", "Revisión y pruebas de aceptación"},
	{"comprobante", 2, "candela.vulcano", -9, "4", "Plantilla del comprobante"},
	{"comprobante", 2, "ana.iriarte", -8, "3", "Generación del código QR"},
	{"tasa", 2, "carlos.moya", -10, "3", "Análisis de la pasarela de pagos"},
	{"tasa", 2, "candela.vulcano", -9, "4", "Maqueta de la pantalla de pago"},
	{"tasa", 2, "jimena.martinez", -9, "3", "Pruebas de la maqueta con usuarios"},
	{"examen", 3, "ana.iriarte", -6, "8", "Banco de preguntas y sorteo"},
	{"examen", 3, "candela.vulcano", -4, "7", "Pantalla del examen con temporizador"},
	{"examen", 3, "carlos.moya", -2, "5", "Corrección automática del examen"},
	{"reprogramar", 3, "jimena.martinez", -3, "4", "Pantalla para cambiar el turno"},
	{"reprogramar", 3, "candela.vulcano", -1, "3", "Liberar el cupo del turno anterior"},
}

// bloqueHoras imputa horas históricas en los Sprints 1 a 3. No exige que el
// Sprint esté Activo (regla de la API de HU-07): son datos de demostración.
func bloqueHoras(ctx context.Context, tx *sql.Tx, c *carga) error {
	sprints := make([]domain.Sprint, 0, 3)
	for i, definicion := range sprintsDemo[:3] {
		sprints = append(sprints, domain.Sprint{ID: c.sprints[i], ProyectoID: proyectoID,
			FechaInicio: c.dia(definicion.desde), FechaFin: c.dia(definicion.hasta)})
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
