package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteVotacionRepository implementa VotacionRepository sobre SQLite.
type SQLiteVotacionRepository struct {
	db *sql.DB
}

func NewSQLiteVotacionRepository(db *sql.DB) *SQLiteVotacionRepository {
	return &SQLiteVotacionRepository{db: db}
}

func (r *SQLiteVotacionRepository) CrearRonda(ctx context.Context, ronda domain.RondaVotacion) (domain.RondaVotacion, error) {
	return ronda, nil
}

func (r *SQLiteVotacionRepository) ObtenerRondaActual(ctx context.Context, historiaID int64) (domain.RondaVotacion, error) {
	return domain.RondaVotacion{}, nil
}

func (r *SQLiteVotacionRepository) GuardarVoto(ctx context.Context, rondaID int64, v domain.Voto) error {
	return nil
}

func (r *SQLiteVotacionRepository) MarcarRevelada(ctx context.Context, rondaID int64, en time.Time) error {
	return nil
}
