package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteUserRepository implementa UserRepository sobre SQLite.
type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// Create persiste el integrante, asignando ID y CreatedAt si vienen vacíos.
func (r *SQLiteUserRepository) Create(u *domain.User) error {
	if u.ID == "" {
		id, err := nuevoID()
		if err != nil {
			return err
		}
		u.ID = id
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	_, err := r.db.Exec(
		`INSERT INTO users (id, name, normalized_name, created_at) VALUES (?, ?, ?, ?)`,
		u.ID, u.Name, u.NormalizedName, u.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUsuarioDuplicado
		}
		return fmt.Errorf("crear usuario: %w", err)
	}
	return nil
}

// FindByNormalizedName devuelve el integrante cuyo nombre normalizado coincide.
func (r *SQLiteUserRepository) FindByNormalizedName(normalizedName string) (*domain.User, error) {
	row := r.db.QueryRow(
		`SELECT id, name, normalized_name, created_at FROM users WHERE normalized_name = ?`,
		normalizedName,
	)
	var (
		u         domain.User
		createdAt string
	)
	if err := row.Scan(&u.ID, &u.Name, &u.NormalizedName, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoEncontrado
		}
		return nil, fmt.Errorf("buscar usuario: %w", err)
	}
	u.CreatedAt = parsearTiempo(createdAt)
	return &u, nil
}

func parsearTiempo(valor string) time.Time {
	t, err := time.Parse(time.RFC3339, valor)
	if err != nil {
		return time.Time{}
	}
	return t
}
