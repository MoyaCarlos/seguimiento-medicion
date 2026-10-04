package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteSprintRepository implementa SprintRepository sobre SQLite.
type SQLiteSprintRepository struct {
	db *sql.DB
}

func NewSQLiteSprintRepository(db *sql.DB) *SQLiteSprintRepository {
	return &SQLiteSprintRepository{db: db}
}

const columnasSprint = "id, proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado"

// Guardar persiste un Sprint nuevo y devuelve una copia con el ID asignado.
func (r *SQLiteSprintRepository) Guardar(ctx context.Context, s domain.Sprint) (domain.Sprint, error) {
	resultado, err := r.db.ExecContext(ctx,
		`INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (?, ?, ?, ?, ?)`,
		s.ProyectoID, textoONulo(s.SprintGoal), fechaASQL(s.FechaInicio), fechaASQL(s.FechaFin), string(s.Estado),
	)
	if err != nil {
		return domain.Sprint{}, fmt.Errorf("guardar sprint: %w", err)
	}
	id, err := resultado.LastInsertId()
	if err != nil {
		return domain.Sprint{}, fmt.Errorf("obtener id del sprint: %w", err)
	}
	s.ID = id
	return s, nil
}

func (r *SQLiteSprintRepository) ObtenerPorID(ctx context.Context, id int64) (domain.Sprint, error) {
	s, err := escanearSprint(r.db.QueryRowContext(ctx, "SELECT "+columnasSprint+" FROM sprints WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	if err != nil {
		return domain.Sprint{}, fmt.Errorf("obtener sprint %d: %w", id, err)
	}
	return s, nil
}

func (r *SQLiteSprintRepository) ListarPorProyecto(ctx context.Context, proyectoID int64) ([]domain.Sprint, error) {
	filas, err := r.db.QueryContext(ctx, "SELECT "+columnasSprint+" FROM sprints WHERE proyecto_id = ? ORDER BY id", proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar sprints del proyecto %d: %w", proyectoID, err)
	}
	defer filas.Close()

	var sprints []domain.Sprint
	for filas.Next() {
		s, err := escanearSprint(filas)
		if err != nil {
			return nil, fmt.Errorf("leer sprint: %w", err)
		}
		sprints = append(sprints, s)
	}
	return sprints, filas.Err()
}

func (r *SQLiteSprintRepository) Actualizar(ctx context.Context, s domain.Sprint) error {
	resultado, err := r.db.ExecContext(ctx,
		`UPDATE sprints SET sprint_goal = ?, fecha_inicio = ?, fecha_fin = ?, estado = ? WHERE id = ?`,
		textoONulo(s.SprintGoal), fechaASQL(s.FechaInicio), fechaASQL(s.FechaFin), string(s.Estado), s.ID,
	)
	if err != nil {
		return fmt.Errorf("actualizar sprint %d: %w", s.ID, err)
	}
	if n, err := resultado.RowsAffected(); err == nil && n == 0 {
		return domain.ErrSprintNoEncontrado
	}
	return nil
}

type escaneable interface {
	Scan(dest ...any) error
}

func escanearSprint(fila escaneable) (domain.Sprint, error) {
	var (
		s                 domain.Sprint
		goal, inicio, fin sql.NullString
		estado            string
	)
	if err := fila.Scan(&s.ID, &s.ProyectoID, &goal, &inicio, &fin, &estado); err != nil {
		return domain.Sprint{}, err
	}
	s.SprintGoal = goal.String
	s.Estado = domain.EstadoSprint(estado)
	var err error
	if s.FechaInicio, err = fechaDesdeSQL(inicio); err != nil {
		return domain.Sprint{}, err
	}
	if s.FechaFin, err = fechaDesdeSQL(fin); err != nil {
		return domain.Sprint{}, err
	}
	return s, nil
}

// Las fechas de un Sprint son días, sin hora: se guardan como AAAA-MM-DD.
func fechaASQL(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t.Format(time.DateOnly)
}

func fechaDesdeSQL(v sql.NullString) (time.Time, error) {
	if !v.Valid {
		return time.Time{}, nil
	}
	return time.Parse(time.DateOnly, v.String)
}

func textoONulo(s string) any {
	if s == "" {
		return nil
	}
	return s
}
