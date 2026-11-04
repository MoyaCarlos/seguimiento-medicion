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

// SQLiteCredencialRepository implementa RegistroCredenciales y BuscadorCredenciales sobre SQLite.
type SQLiteCredencialRepository struct {
	db *sql.DB
}

func NewSQLiteCredencialRepository(db *sql.DB) *SQLiteCredencialRepository {
	return &SQLiteCredencialRepository{db: db}
}

// Registrar persiste el usuario con el hash de su contraseña y devuelve una
// copia con el ID asignado. La unicidad del nombre normalizado la garantiza la
// base, así dos registros simultáneos del mismo nombre no pueden crear ambos.
func (r *SQLiteCredencialRepository) Registrar(ctx context.Context, u domain.User, hash string) (domain.User, error) {
	if u.CreadoEn.IsZero() {
		u.CreadoEn = time.Now().UTC()
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO usuarios (nombre, nombre_normalizado, contrasena_hash, creado_en) VALUES (?, ?, ?, ?)`,
		u.Nombre, u.NombreNormalizado, hash, u.CreadoEn.UTC().Format(time.RFC3339),
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return domain.User{}, ErrUsuarioDuplicado
		}
		return domain.User{}, fmt.Errorf("registrar usuario: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, fmt.Errorf("obtener id del usuario: %w", err)
	}
	u.ID = id
	return u, nil
}

// ObtenerPorNombreNormalizado devuelve el usuario y su hash, vacío si el
// usuario no tiene contraseña (NULL).
func (r *SQLiteCredencialRepository) ObtenerPorNombreNormalizado(ctx context.Context, normalizado string) (domain.User, string, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, nombre, nombre_normalizado, contrasena_hash, creado_en FROM usuarios WHERE nombre_normalizado = ?`,
		normalizado,
	)
	var (
		u        domain.User
		hash     sql.NullString
		creadoEn string
	)
	if err := row.Scan(&u.ID, &u.Nombre, &u.NombreNormalizado, &hash, &creadoEn); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, "", domain.ErrNoEncontrado
		}
		return domain.User{}, "", fmt.Errorf("buscar credenciales: %w", err)
	}
	u.CreadoEn = parsearTiempo(creadoEn)
	return u, hash.String, nil
}
