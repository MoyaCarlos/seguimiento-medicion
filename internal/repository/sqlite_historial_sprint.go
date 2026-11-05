package repository

import (
	"context"
	"database/sql"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteHistorialSprintRepository implementa HistorialSprintRepository sobre
// SQLite. Solo lee.
type SQLiteHistorialSprintRepository struct {
	db *sql.DB
}

func NewSQLiteHistorialSprintRepository(db *sql.DB) *SQLiteHistorialSprintRepository {
	return &SQLiteHistorialSprintRepository{db: db}
}

func (r *SQLiteHistorialSprintRepository) HistoriasComprometidas(ctx context.Context, sprintID int64) ([]domain.HistoriaComprometida, error) {
	return nil, nil
}

func (r *SQLiteHistorialSprintRepository) MinutosPorHistoria(ctx context.Context, sprintID int64) (map[int64]int, error) {
	return nil, nil
}

func (r *SQLiteHistorialSprintRepository) DefectosDeSprint(ctx context.Context, sprintID int64) ([]domain.DefectoDetalle, error) {
	return nil, nil
}
