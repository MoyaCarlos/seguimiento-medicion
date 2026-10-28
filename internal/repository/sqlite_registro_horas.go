package repository

import (
	"context"
	"database/sql"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteRegistroHorasRepository implementa RegistroHorasRepository sobre SQLite.
type SQLiteRegistroHorasRepository struct {
	db *sql.DB
}

func NewSQLiteRegistroHorasRepository(db *sql.DB) *SQLiteRegistroHorasRepository {
	return &SQLiteRegistroHorasRepository{db: db}
}

func (r *SQLiteRegistroHorasRepository) Registrar(ctx context.Context, registro domain.RegistroHoras) (domain.RegistroHoras, error) {
	return domain.RegistroHoras{}, nil
}

func (r *SQLiteRegistroHorasRepository) ListarPorHistoria(ctx context.Context, historiaID int64) ([]domain.RegistroHorasDetalle, error) {
	return nil, nil
}
