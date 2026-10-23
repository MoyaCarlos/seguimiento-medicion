package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// SQLiteProjectRepository implementa ProjectRepository sobre SQLite.
type SQLiteProjectRepository struct {
	db *sql.DB
}

func NewSQLiteProjectRepository(db *sql.DB) *SQLiteProjectRepository {
	return &SQLiteProjectRepository{db: db}
}

// Guardar persiste el proyecto y devuelve una copia con el ID asignado.
func (r *SQLiteProjectRepository) Guardar(ctx context.Context, p domain.Project) (domain.Project, error) {
	return insertarProyecto(ctx, r.db, p)
}

// GuardarConScrumMaster inserta el proyecto y su membresía de Scrum Master en
// una sola transacción, de modo que si la segunda escritura falla no queda un
// proyecto sin Scrum Master.
func (r *SQLiteProjectRepository) GuardarConScrumMaster(ctx context.Context, p domain.Project, creadorID int64) (domain.Project, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Project{}, fmt.Errorf("iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	p, err = insertarProyecto(ctx, tx, p)
	if err != nil {
		return domain.Project{}, err
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, ?, ?)`,
		p.ID, creadorID, string(domain.RolScrumMaster), time.Now().UTC().Format(time.RFC3339),
	); err != nil {
		return domain.Project{}, fmt.Errorf("vincular scrum master: %w", errorDeEscritura(err))
	}

	if err := tx.Commit(); err != nil {
		return domain.Project{}, fmt.Errorf("confirmar transacción: %w", err)
	}
	return p, nil
}

// ejecutorSQL abstrae *sql.DB y *sql.Tx para compartir el INSERT de proyectos.
type ejecutorSQL interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// insertarProyecto inserta el proyecto en el ejecutor dado y asigna el ID.
func insertarProyecto(ctx context.Context, ex ejecutorSQL, p domain.Project) (domain.Project, error) {
	if p.CreadoEn.IsZero() {
		p.CreadoEn = time.Now().UTC()
	}
	res, err := ex.ExecContext(ctx,
		`INSERT INTO proyectos (nombre, descripcion, fecha_inicio, fecha_fin, creado_en)
		 VALUES (?, ?, ?, ?, ?)`,
		p.Nombre, p.Descripcion,
		tiempoOpcionalAValor(p.FechaInicio), tiempoOpcionalAValor(p.FechaFin),
		p.CreadoEn.Format(time.RFC3339),
	)
	if err != nil {
		return domain.Project{}, fmt.Errorf("crear proyecto: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Project{}, fmt.Errorf("obtener id del proyecto: %w", err)
	}
	p.ID = id
	return p, nil
}

// ObtenerPorID devuelve el proyecto con el identificador indicado.
func (r *SQLiteProjectRepository) ObtenerPorID(ctx context.Context, id int64) (domain.Project, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, nombre, descripcion, fecha_inicio, fecha_fin, creado_en FROM proyectos WHERE id = ?`,
		id,
	)
	return escanearProyecto(row)
}

// Actualizar reemplaza los datos editables del proyecto.
func (r *SQLiteProjectRepository) Actualizar(ctx context.Context, p domain.Project) error {
	resultado, err := r.db.ExecContext(ctx,
		`UPDATE proyectos SET nombre = ?, descripcion = ?, fecha_inicio = ?, fecha_fin = ? WHERE id = ?`,
		p.Nombre, p.Descripcion,
		tiempoOpcionalAValor(p.FechaInicio), tiempoOpcionalAValor(p.FechaFin),
		p.ID,
	)
	if err != nil {
		return fmt.Errorf("actualizar proyecto: %w", err)
	}
	afectadas, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("actualizar proyecto: %w", err)
	}
	if afectadas == 0 {
		return domain.ErrProyectoNoEncontrado
	}
	return nil
}

// AgregarIntegrante vincula un integrante con un rol al proyecto y devuelve la
// membresía guardada (con su CreadoEn asignado).
func (r *SQLiteProjectRepository) AgregarIntegrante(ctx context.Context, m domain.Membership) (domain.Membership, error) {
	if m.CreadoEn.IsZero() {
		m.CreadoEn = time.Now().UTC()
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, ?, ?)`,
		m.ProjectID, m.UserID, string(m.Role), m.CreadoEn.Format(time.RFC3339),
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) &&
			(sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE || sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY) {
			return domain.Membership{}, ErrMiembroDuplicado
		}
		return domain.Membership{}, fmt.Errorf("agregar integrante: %w", errorDeEscritura(err))
	}
	return m, nil
}

// ListarIntegrantes devuelve los integrantes del proyecto con su nombre y rol.
func (r *SQLiteProjectRepository) ListarIntegrantes(ctx context.Context, proyectoID int64) ([]domain.Member, error) {
	filas, err := r.db.QueryContext(ctx,
		`SELECT i.usuario_id, i.proyecto_id, u.nombre, i.rol, i.creado_en
		 FROM integrantes i
		 JOIN usuarios u ON u.id = i.usuario_id
		 WHERE i.proyecto_id = ?
		 ORDER BY i.rowid`,
		proyectoID,
	)
	if err != nil {
		return nil, fmt.Errorf("listar integrantes: %w", err)
	}
	defer filas.Close()

	var miembros []domain.Member
	for filas.Next() {
		var (
			m         domain.Member
			createdAt string
		)
		if err := filas.Scan(&m.UserID, &m.ProjectID, &m.Nombre, &m.Role, &createdAt); err != nil {
			return nil, fmt.Errorf("leer integrante: %w", err)
		}
		m.CreadoEn = parsearTiempo(createdAt)
		miembros = append(miembros, m)
	}
	if err := filas.Err(); err != nil {
		return nil, fmt.Errorf("listar integrantes: %w", err)
	}
	return miembros, nil
}

func escanearProyecto(row *sql.Row) (domain.Project, error) {
	var (
		p                      domain.Project
		inicio, fin, createdAt sql.NullString
	)
	if err := row.Scan(&p.ID, &p.Nombre, &p.Descripcion, &inicio, &fin, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, domain.ErrProyectoNoEncontrado
		}
		return domain.Project{}, fmt.Errorf("leer proyecto: %w", err)
	}
	p.FechaInicio = valorATiempoOpcional(inicio)
	p.FechaFin = valorATiempoOpcional(fin)
	p.CreadoEn = parsearTiempo(createdAt.String)
	return p, nil
}

func tiempoOpcionalAValor(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.DateOnly)
}

func valorATiempoOpcional(v sql.NullString) *time.Time {
	if !v.Valid {
		return nil
	}
	t, err := time.Parse(time.DateOnly, v.String)
	if err != nil {
		return nil
	}
	return &t
}
