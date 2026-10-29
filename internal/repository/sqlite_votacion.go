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

// SQLiteVotacionRepository implementa VotacionRepository sobre SQLite.
type SQLiteVotacionRepository struct {
	db *sql.DB
}

func NewSQLiteVotacionRepository(db *sql.DB) *SQLiteVotacionRepository {
	return &SQLiteVotacionRepository{db: db}
}

// CrearRonda persiste la ronda y devuelve una copia con el ID asignado. El
// UNIQUE (historia_id, numero) frena dos inicios simultáneos.
func (r *SQLiteVotacionRepository) CrearRonda(ctx context.Context, ronda domain.RondaVotacion) (domain.RondaVotacion, error) {
	resultado, err := r.db.ExecContext(ctx,
		`INSERT INTO rondas_votacion (historia_id, numero) VALUES (?, ?)`,
		ronda.HistoriaID, ronda.Numero,
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return domain.RondaVotacion{}, domain.ErrVotacionYaIniciada
		}
		return domain.RondaVotacion{}, fmt.Errorf("crear ronda: %w", errorDeEscritura(err))
	}
	id, err := resultado.LastInsertId()
	if err != nil {
		return domain.RondaVotacion{}, fmt.Errorf("obtener id de la ronda: %w", err)
	}
	ronda.ID = id
	return ronda, nil
}

// ObtenerRondaActual devuelve la ronda de mayor número de la historia con sus
// votos ordenados por usuario.
func (r *SQLiteVotacionRepository) ObtenerRondaActual(ctx context.Context, historiaID int64) (domain.RondaVotacion, error) {
	var (
		ronda    domain.RondaVotacion
		revelada sql.NullString
	)
	err := r.db.QueryRowContext(ctx,
		`SELECT id, historia_id, numero, revelada_en FROM rondas_votacion
		 WHERE historia_id = ? ORDER BY numero DESC LIMIT 1`,
		historiaID,
	).Scan(&ronda.ID, &ronda.HistoriaID, &ronda.Numero, &revelada)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.RondaVotacion{}, domain.ErrVotacionNoIniciada
	}
	if err != nil {
		return domain.RondaVotacion{}, fmt.Errorf("obtener ronda de la historia %d: %w", historiaID, err)
	}
	if revelada.Valid {
		en := parsearTiempo(revelada.String)
		ronda.ReveladaEn = &en
	}

	ronda.Votos, err = r.votosDe(ctx, ronda.ID)
	if err != nil {
		return domain.RondaVotacion{}, err
	}
	return ronda, nil
}

func (r *SQLiteVotacionRepository) votosDe(ctx context.Context, rondaID int64) ([]domain.Voto, error) {
	filas, err := r.db.QueryContext(ctx,
		"SELECT usuario_id, valor FROM votos WHERE ronda_id = ? ORDER BY usuario_id", rondaID)
	if err != nil {
		return nil, fmt.Errorf("listar votos de la ronda %d: %w", rondaID, err)
	}
	defer filas.Close()

	var votos []domain.Voto
	for filas.Next() {
		var v domain.Voto
		if err := filas.Scan(&v.UsuarioID, &v.Valor); err != nil {
			return nil, fmt.Errorf("leer voto: %w", err)
		}
		votos = append(votos, v)
	}
	return votos, filas.Err()
}

// GuardarVoto agrega el voto o reemplaza el del mismo usuario, solo si la
// ronda sigue sin revelar (research.md R6).
func (r *SQLiteVotacionRepository) GuardarVoto(ctx context.Context, rondaID int64, v domain.Voto) error {
	resultado, err := r.db.ExecContext(ctx,
		`INSERT INTO votos (ronda_id, usuario_id, valor)
		 SELECT ?, ?, ? WHERE EXISTS (SELECT 1 FROM rondas_votacion WHERE id = ? AND revelada_en IS NULL)
		 ON CONFLICT (ronda_id, usuario_id) DO UPDATE SET valor = excluded.valor`,
		rondaID, v.UsuarioID, v.Valor, rondaID,
	)
	if err != nil {
		return fmt.Errorf("guardar voto: %w", errorDeEscritura(err))
	}
	return ningunaFilaEsRondaRevelada(resultado)
}

// MarcarRevelada guarda el instante de revelación, solo si la ronda sigue sin
// revelar (research.md R6).
func (r *SQLiteVotacionRepository) MarcarRevelada(ctx context.Context, rondaID int64, en time.Time) error {
	resultado, err := r.db.ExecContext(ctx,
		`UPDATE rondas_votacion SET revelada_en = ? WHERE id = ? AND revelada_en IS NULL`,
		en.UTC().Format(time.RFC3339), rondaID,
	)
	if err != nil {
		return fmt.Errorf("revelar ronda %d: %w", rondaID, err)
	}
	return ningunaFilaEsRondaRevelada(resultado)
}

func ningunaFilaEsRondaRevelada(resultado sql.Result) error {
	n, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("contar filas escritas: %w", err)
	}
	if n == 0 {
		return domain.ErrRondaRevelada
	}
	return nil
}
