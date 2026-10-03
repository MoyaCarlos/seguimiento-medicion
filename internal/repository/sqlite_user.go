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

// SQLiteUserRepository implementa UserRepository sobre SQLite.
type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// Guardar persiste el integrante y devuelve una copia con el ID asignado.
func (r *SQLiteUserRepository) Guardar(ctx context.Context, u domain.User) (domain.User, error) {
	if u.CreadoEn.IsZero() {
		u.CreadoEn = time.Now().UTC()
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (name, normalized_name, created_at) VALUES (?, ?, ?)`,
		u.Nombre, u.NombreNormalizado, u.CreadoEn.Format(time.RFC3339),
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return domain.User{}, ErrUsuarioDuplicado
		}
		return domain.User{}, fmt.Errorf("crear usuario: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, fmt.Errorf("obtener id del usuario: %w", err)
	}
	u.ID = id
	return u, nil
}

// ObtenerPorNombreNormalizado devuelve el integrante cuyo nombre normalizado coincide.
func (r *SQLiteUserRepository) ObtenerPorNombreNormalizado(ctx context.Context, normalized string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, normalized_name, created_at FROM users WHERE normalized_name = ?`,
		normalized,
	)
	var (
		u         domain.User
		createdAt string
	)
	if err := row.Scan(&u.ID, &u.Nombre, &u.NombreNormalizado, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoEncontrado
		}
		return nil, fmt.Errorf("buscar usuario: %w", err)
	}
	u.CreadoEn = parsearTiempo(createdAt)
	return &u, nil
}

func parsearTiempo(valor string) time.Time {
	t, err := time.Parse(time.RFC3339, valor)
	if err != nil {
		return time.Time{}
	}
	return t
}
