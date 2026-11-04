package repository

import (
	"context"
	"database/sql"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteCredencialRepository implementa RegistroCredenciales y BuscadorCredenciales sobre SQLite.
type SQLiteCredencialRepository struct {
	db *sql.DB
}

func NewSQLiteCredencialRepository(db *sql.DB) *SQLiteCredencialRepository {
	return &SQLiteCredencialRepository{db: db}
}

// Registrar persiste el usuario con el hash de su contraseña.
func (r *SQLiteCredencialRepository) Registrar(ctx context.Context, u domain.User, hash string) (domain.User, error) {
	return domain.User{}, nil
}

// ObtenerPorNombreNormalizado devuelve el usuario y su hash (vacío si no tiene contraseña).
func (r *SQLiteCredencialRepository) ObtenerPorNombreNormalizado(ctx context.Context, normalizado string) (domain.User, string, error) {
	return domain.User{}, "", nil
}
