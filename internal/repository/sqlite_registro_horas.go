package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteRegistroHorasRepository implementa RegistroHorasRepository sobre SQLite.
type SQLiteRegistroHorasRepository struct {
	db *sql.DB
}

func NewSQLiteRegistroHorasRepository(db *sql.DB) *SQLiteRegistroHorasRepository {
	return &SQLiteRegistroHorasRepository{db: db}
}

// Registrar inserta la carga solo si la historia sigue en ese Sprint y el
// Sprint sigue Activo, en la misma sentencia: si el Sprint se cerró entre la
// validación del service y el INSERT, la carga se rechaza (research.md R5).
func (r *SQLiteRegistroHorasRepository) Registrar(ctx context.Context, registro domain.RegistroHoras) (domain.RegistroHoras, error) {
	creadoEn := time.Now().UTC()
	resultado, err := r.db.ExecContext(ctx,
		`INSERT INTO registros_horas (proyecto_id, historia_id, usuario_id, sprint_id, fecha, actividad, minutos, creado_en)
		 SELECT ?, ?, ?, ?, ?, ?, ?, ?
		 WHERE EXISTS (SELECT 1 FROM historias h JOIN sprints s ON s.id = h.sprint_id
		               WHERE h.id = ? AND h.sprint_id = ? AND s.estado = 'Activo')`,
		registro.ProyectoID, registro.HistoriaID, registro.UsuarioID, registro.SprintID,
		fechaASQL(registro.Fecha), registro.Actividad, registro.Minutos, creadoEn.Format(time.RFC3339),
		registro.HistoriaID, registro.SprintID,
	)
	if err != nil {
		return domain.RegistroHoras{}, fmt.Errorf("registrar horas: %w", errorDeEscritura(err))
	}
	insertadas, err := resultado.RowsAffected()
	if err != nil {
		return domain.RegistroHoras{}, fmt.Errorf("registrar horas: %w", err)
	}
	if insertadas == 0 {
		return domain.RegistroHoras{}, domain.ErrHistoriaFueraDelSprintActivo
	}
	id, err := resultado.LastInsertId()
	if err != nil {
		return domain.RegistroHoras{}, fmt.Errorf("obtener id de la carga: %w", err)
	}
	registro.ID = id
	registro.CreadoEn = creadoEn
	return registro, nil
}

// ListarPorHistoria devuelve las cargas de la historia con el nombre del
// integrante, ordenadas por fecha trabajada y luego por orden de carga.
func (r *SQLiteRegistroHorasRepository) ListarPorHistoria(ctx context.Context, historiaID int64) ([]domain.RegistroHorasDetalle, error) {
	filas, err := r.db.QueryContext(ctx,
		`SELECT rh.id, rh.proyecto_id, rh.historia_id, rh.usuario_id, rh.sprint_id, rh.fecha, rh.actividad, rh.minutos, rh.creado_en, u.nombre
		 FROM registros_horas rh
		 JOIN usuarios u ON u.id = rh.usuario_id
		 WHERE rh.historia_id = ?
		 ORDER BY rh.fecha, rh.id`,
		historiaID,
	)
	if err != nil {
		return nil, fmt.Errorf("listar horas de la historia %d: %w", historiaID, err)
	}
	defer filas.Close()

	registros := make([]domain.RegistroHorasDetalle, 0)
	for filas.Next() {
		var (
			d               domain.RegistroHorasDetalle
			fecha, creadoEn string
		)
		if err := filas.Scan(&d.ID, &d.ProyectoID, &d.HistoriaID, &d.UsuarioID, &d.SprintID,
			&fecha, &d.Actividad, &d.Minutos, &creadoEn, &d.NombreIntegrante); err != nil {
			return nil, fmt.Errorf("leer carga de horas: %w", err)
		}
		if d.Fecha, err = time.Parse(time.DateOnly, fecha); err != nil {
			return nil, fmt.Errorf("leer fecha de la carga %d: %w", d.ID, err)
		}
		d.CreadoEn = parsearTiempo(creadoEn)
		registros = append(registros, d)
	}
	if err := filas.Err(); err != nil {
		return nil, fmt.Errorf("listar horas de la historia %d: %w", historiaID, err)
	}
	return registros, nil
}
