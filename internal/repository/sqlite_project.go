package repository

import (
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

// Create persiste el proyecto, asignando ID y CreatedAt si vienen vacíos.
func (r *SQLiteProjectRepository) Create(p *domain.Project) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	res, err := r.db.Exec(
		`INSERT INTO projects (name, description, start_date, end_date, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		p.Name, p.Description,
		tiempoOpcionalAValor(p.StartDate), tiempoOpcionalAValor(p.EndDate),
		p.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("crear proyecto: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener id del proyecto: %w", err)
	}
	p.ID = id
	return nil
}

// CreateWithScrumMaster inserta el proyecto y su membresía de Scrum Master en
// una sola transacción, de modo que si la segunda escritura falla no queda un
// proyecto sin Scrum Master.
func (r *SQLiteProjectRepository) CreateWithScrumMaster(p *domain.Project, creatorID int64) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`INSERT INTO projects (name, description, start_date, end_date, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		p.Name, p.Description,
		tiempoOpcionalAValor(p.StartDate), tiempoOpcionalAValor(p.EndDate),
		p.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("crear proyecto: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("obtener id del proyecto: %w", err)
	}
	p.ID = id

	if _, err := tx.Exec(
		`INSERT INTO project_members (project_id, user_id, role, created_at) VALUES (?, ?, ?, ?)`,
		id, creatorID, string(domain.RolScrumMaster), time.Now().UTC().Format(time.RFC3339),
	); err != nil {
		return fmt.Errorf("vincular scrum master: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar transacción: %w", err)
	}
	return nil
}

// GetByID devuelve el proyecto con el identificador indicado.
func (r *SQLiteProjectRepository) GetByID(id int64) (*domain.Project, error) {
	row := r.db.QueryRow(
		`SELECT id, name, description, start_date, end_date, created_at FROM projects WHERE id = ?`,
		id,
	)
	return escanearProyecto(row)
}

// Update reemplaza los datos editables del proyecto.
func (r *SQLiteProjectRepository) Update(p *domain.Project) error {
	resultado, err := r.db.Exec(
		`UPDATE projects SET name = ?, description = ?, start_date = ?, end_date = ? WHERE id = ?`,
		p.Name, p.Description,
		tiempoOpcionalAValor(p.StartDate), tiempoOpcionalAValor(p.EndDate),
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
		return ErrNoEncontrado
	}
	return nil
}

// AddMember vincula un integrante con un rol al proyecto.
func (r *SQLiteProjectRepository) AddMember(m *domain.Membership) error {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	_, err := r.db.Exec(
		`INSERT INTO project_members (project_id, user_id, role, created_at) VALUES (?, ?, ?, ?)`,
		m.ProjectID, m.UserID, string(m.Role), m.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.Code() {
			case sqlite3.SQLITE_CONSTRAINT_UNIQUE, sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
				return ErrMiembroDuplicado
			case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY:
				return ErrIntegridadReferencial
			}
		}
		return fmt.Errorf("agregar integrante: %w", err)
	}
	return nil
}

// ListMembers devuelve los integrantes del proyecto con su nombre y rol.
func (r *SQLiteProjectRepository) ListMembers(projectID int64) ([]domain.Member, error) {
	filas, err := r.db.Query(
		`SELECT pm.user_id, pm.project_id, u.name, pm.role, pm.created_at
		 FROM project_members pm
		 JOIN users u ON u.id = pm.user_id
		 WHERE pm.project_id = ?
		 ORDER BY pm.created_at, u.name`,
		projectID,
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
		if err := filas.Scan(&m.UserID, &m.ProjectID, &m.Name, &m.Role, &createdAt); err != nil {
			return nil, fmt.Errorf("leer integrante: %w", err)
		}
		m.CreatedAt = parsearTiempo(createdAt)
		miembros = append(miembros, m)
	}
	if err := filas.Err(); err != nil {
		return nil, fmt.Errorf("listar integrantes: %w", err)
	}
	return miembros, nil
}

func escanearProyecto(row *sql.Row) (*domain.Project, error) {
	var (
		p                      domain.Project
		inicio, fin, createdAt sql.NullString
	)
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &inicio, &fin, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoEncontrado
		}
		return nil, fmt.Errorf("leer proyecto: %w", err)
	}
	p.StartDate = valorATiempoOpcional(inicio)
	p.EndDate = valorATiempoOpcional(fin)
	p.CreatedAt = parsearTiempo(createdAt.String)
	return &p, nil
}

func tiempoOpcionalAValor(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.RFC3339)
}

func valorATiempoOpcional(v sql.NullString) *time.Time {
	if !v.Valid {
		return nil
	}
	t, err := time.Parse(time.RFC3339, v.String)
	if err != nil {
		return nil
	}
	return &t
}
