package repository

import (
	"context"
	"database/sql"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteDefectoRepository implementa DefectoRepository sobre SQLite.
type SQLiteDefectoRepository struct {
	db *sql.DB
}

func NewSQLiteDefectoRepository(db *sql.DB) *SQLiteDefectoRepository {
	return &SQLiteDefectoRepository{db: db}
}

func (r *SQLiteDefectoRepository) Registrar(ctx context.Context, d domain.Defecto) (domain.Defecto, error) {
	return domain.Defecto{}, nil
}

func (r *SQLiteDefectoRepository) ObtenerPorID(ctx context.Context, id int64) (domain.DefectoDetalle, error) {
	return domain.DefectoDetalle{}, nil
}

func (r *SQLiteDefectoRepository) Resolver(ctx context.Context, d domain.Defecto) error {
	return nil
}

func (r *SQLiteDefectoRepository) ListarPorProyecto(ctx context.Context, proyectoID int64) ([]domain.DefectoDetalle, error) {
	return nil, nil
}
