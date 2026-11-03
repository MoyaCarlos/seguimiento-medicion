package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// SQLiteSprintRepository implementa SprintRepository sobre SQLite.
type SQLiteSprintRepository struct {
	db *sql.DB
}

func NewSQLiteSprintRepository(db *sql.DB) *SQLiteSprintRepository {
	return &SQLiteSprintRepository{db: db}
}

// El número del Sprint se deriva de los de su proyecto, sin columna propia
// (#129). Es una subconsulta y no ROW_NUMBER() porque esa función numeraría
// solo las filas que deja el WHERE de cada consulta.
const columnasSprint = `id, proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado,
	(SELECT COUNT(*) FROM sprints previos WHERE previos.proyecto_id = sprints.proyecto_id AND previos.id <= sprints.id)`

// Guardar persiste un Sprint nuevo y lo devuelve con el ID y el número asignados.
func (r *SQLiteSprintRepository) Guardar(ctx context.Context, s domain.Sprint) (domain.Sprint, error) {
	resultado, err := r.db.ExecContext(ctx,
		`INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (?, ?, ?, ?, ?)`,
		s.ProyectoID, textoONulo(s.SprintGoal), fechaASQL(s.FechaInicio), fechaASQL(s.FechaFin), string(s.Estado),
	)
	if err != nil {
		return domain.Sprint{}, fmt.Errorf("guardar sprint: %w", errorDeEscrituraDeSprint(s, err))
	}
	id, err := resultado.LastInsertId()
	if err != nil {
		return domain.Sprint{}, fmt.Errorf("obtener id del sprint: %w", err)
	}
	return r.ObtenerPorID(ctx, id)
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

// Actualizar persiste el inicio de un Sprint. Solo lo aplica si en la base
// sigue Pendiente: de dos inicios simultáneos, el segundo recibe
// domain.ErrSprintNoPendiente y no pisa el Goal ni las fechas (#158).
func (r *SQLiteSprintRepository) Actualizar(ctx context.Context, s domain.Sprint) error {
	resultado, err := r.db.ExecContext(ctx,
		`UPDATE sprints SET sprint_goal = ?, fecha_inicio = ?, fecha_fin = ?, estado = ? WHERE id = ? AND estado = ?`,
		textoONulo(s.SprintGoal), fechaASQL(s.FechaInicio), fechaASQL(s.FechaFin), string(s.Estado), s.ID,
		string(domain.SprintPendiente),
	)
	if err != nil {
		return fmt.Errorf("actualizar sprint %d: %w", s.ID, errorDeEscrituraDeSprint(s, err))
	}
	if n, err := resultado.RowsAffected(); err == nil && n == 0 {
		return sprintSinActualizar(ctx, r.db, s.ID, domain.ErrSprintNoPendiente)
	}
	return nil
}

// CerrarConArrastre persiste el Sprint cerrado y, en la misma transacción,
// registra cada historia no completada en historias_arrastradas con la foto
// de sus Story Points y la devuelve al Product Backlog en estado Nueva
// (MODELO-DE-DATOS.md §4.6). Si algo falla, no se aplica nada.
// Solo cierra el Sprint si en la base sigue Activo: de dos cierres
// simultáneos, el segundo recibe domain.ErrSprintNoActivo (#145).
func (r *SQLiteSprintRepository) CerrarConArrastre(ctx context.Context, s domain.Sprint, noCompletadas []domain.BacklogItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cerrar sprint %d: %w", s.ID, err)
	}
	defer func() { _ = tx.Rollback() }()

	resultado, err := tx.ExecContext(ctx, `UPDATE sprints SET estado = ? WHERE id = ? AND estado = ?`,
		string(s.Estado), s.ID, string(domain.SprintActivo))
	if err != nil {
		return fmt.Errorf("cerrar sprint %d: %w", s.ID, errorDeEscrituraDeSprint(s, err))
	}
	if n, err := resultado.RowsAffected(); err == nil && n == 0 {
		return sprintSinActualizar(ctx, tx, s.ID, domain.ErrSprintNoActivo)
	}

	if len(noCompletadas) > 0 {
		if err := arrastrar(ctx, tx, s, noCompletadas); err != nil {
			return fmt.Errorf("cerrar sprint %d: %w", s.ID, err)
		}
	}
	return tx.Commit()
}

type consultable interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// sprintSinActualizar explica por qué un UPDATE condicional sobre el estado
// no cambió nada: el Sprint no existe o ya no está en el estado esperado, y
// en ese caso devuelve errEstado.
func sprintSinActualizar(ctx context.Context, db consultable, id int64, errEstado error) error {
	var existe bool
	if err := db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM sprints WHERE id = ?)", id).Scan(&existe); err != nil {
		return fmt.Errorf("sprint %d: %w", id, err)
	}
	if !existe {
		return domain.ErrSprintNoEncontrado
	}
	return errEstado
}

func arrastrar(ctx context.Context, tx *sql.Tx, s domain.Sprint, historias []domain.BacklogItem) error {
	ids := make([]any, 0, len(historias)+2)
	ids = append(ids, string(domain.EstadoNueva), s.ID)
	for _, h := range historias {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO historias_arrastradas (sprint_id, historia_id, proyecto_id, estimacion_sp) VALUES (?, ?, ?, ?)`,
			s.ID, h.ID, s.ProyectoID, punteroAValorSQL(h.EstimacionSP),
		); err != nil {
			return errorDeEscritura(err)
		}
		ids = append(ids, h.ID)
	}

	marcadores := strings.Repeat(", ?", len(historias))[2:]
	resultado, err := tx.ExecContext(ctx,
		`UPDATE historias SET sprint_id = NULL, estado = ? WHERE sprint_id = ? AND id IN (`+marcadores+`)`, ids...)
	if err != nil {
		return err
	}
	n, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if n != int64(len(historias)) {
		return ErrHistoriaFueraDelSprint
	}
	return nil
}

// errorDeEscrituraDeSprint traduce además la violación de los índices
// parciales (un Sprint activo y uno pendiente por proyecto) al mismo
// ValidationError que devuelve la validación previa del servicio.
func errorDeEscrituraDeSprint(s domain.Sprint, err error) error {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		switch s.Estado {
		case domain.SprintActivo:
			return domain.ErrYaHaySprintActivo
		case domain.SprintPendiente:
			return domain.ErrYaHaySprintPendiente
		}
	}
	return errorDeEscritura(err)
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
	if err := fila.Scan(&s.ID, &s.ProyectoID, &goal, &inicio, &fin, &estado, &s.Numero); err != nil {
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
