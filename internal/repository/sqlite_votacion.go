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

// RondaActualRevelada indica si la ronda actual está revelada, sin cargar sus votos.
func (r *SQLiteVotacionRepository) RondaActualRevelada(ctx context.Context, historiaID int64) (bool, error) {
	return false, nil
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

// ListarRondasReveladas lee rondas y votos en una sola consulta: la conexión es
// única, así que no se puede consultar los votos con las filas de las rondas
// abiertas (research.md R9).
func (r *SQLiteVotacionRepository) ListarRondasReveladas(ctx context.Context, historiaID int64) ([]domain.RondaVotacion, error) {
	filas, err := r.db.QueryContext(ctx,
		`SELECT r.id, r.numero, r.revelada_en, v.usuario_id, v.valor
		 FROM rondas_votacion r LEFT JOIN votos v ON v.ronda_id = r.id
		 WHERE r.historia_id = ? AND r.revelada_en IS NOT NULL
		 ORDER BY r.numero, v.usuario_id`,
		historiaID,
	)
	if err != nil {
		return nil, fmt.Errorf("listar rondas de la historia %d: %w", historiaID, err)
	}
	defer filas.Close()

	rondas := make([]domain.RondaVotacion, 0)
	for filas.Next() {
		var (
			id        int64
			numero    int
			revelada  string
			usuarioID sql.NullInt64
			valor     sql.NullInt64
		)
		if err := filas.Scan(&id, &numero, &revelada, &usuarioID, &valor); err != nil {
			return nil, fmt.Errorf("leer ronda: %w", err)
		}
		if len(rondas) == 0 || rondas[len(rondas)-1].ID != id {
			en := parsearTiempo(revelada)
			rondas = append(rondas, domain.RondaVotacion{ID: id, HistoriaID: historiaID, Numero: numero, ReveladaEn: &en})
		}
		if usuarioID.Valid {
			actual := &rondas[len(rondas)-1]
			actual.Votos = append(actual.Votos, domain.Voto{UsuarioID: usuarioID.Int64, Valor: int(valor.Int64)})
		}
	}
	return rondas, filas.Err()
}
