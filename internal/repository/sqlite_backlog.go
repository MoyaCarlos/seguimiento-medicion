package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteBacklogRepository implementa BacklogRepository sobre SQLite.
type SQLiteBacklogRepository struct {
	db *sql.DB
}

func NewSQLiteBacklogRepository(db *sql.DB) *SQLiteBacklogRepository {
	return &SQLiteBacklogRepository{db: db}
}

// Guardar persiste la historia y devuelve una copia con el ID asignado.
func (r *SQLiteBacklogRepository) Guardar(ctx context.Context, item domain.BacklogItem) (domain.BacklogItem, error) {
	resultado, err := r.db.ExecContext(ctx,
		`INSERT INTO backlog_items (proyecto_id, titulo, descripcion, prioridad, estado, valor_negocio, estimacion_sp)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		item.ProyectoID,
		item.Titulo,
		item.Descripcion,
		string(item.Prioridad),
		string(item.Estado),
		punteroAValorSQL(item.ValorNegocio),
		punteroAValorSQL(item.EstimacionSP),
	)
	if err != nil {
		return domain.BacklogItem{}, fmt.Errorf("guardar backlog item: %w", err)
	}

	id, err := resultado.LastInsertId()
	if err != nil {
		return domain.BacklogItem{}, fmt.Errorf("obtener id del backlog item: %w", err)
	}
	item.ID = id
	return item, nil
}

func punteroAValorSQL(v *int) any {
	if v == nil {
		return nil
	}
	return *v
}
