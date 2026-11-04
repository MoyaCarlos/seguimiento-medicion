package repository

import (
	"context"
	"database/sql"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteMetricasRepository implementa DatosMetricasRepository sobre SQLite.
type SQLiteMetricasRepository struct {
	db *sql.DB
}

func NewSQLiteMetricasRepository(db *sql.DB) *SQLiteMetricasRepository {
	return &SQLiteMetricasRepository{db: db}
}

// ListarDatosSprintsFinalizados devuelve los datos de cada Sprint finalizado.
func (r *SQLiteMetricasRepository) ListarDatosSprintsFinalizados(ctx context.Context, proyectoID int64) ([]domain.DatosMetricasSprint, error) {
	return nil, nil
}
