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

// ListarPorSprint devuelve las historias asignadas al Sprint, en orden de creación.
func (r *SQLiteBacklogRepository) ListarPorSprint(ctx context.Context, sprintID int64) ([]domain.BacklogItem, error) {
	filas, err := r.db.QueryContext(ctx,
		`SELECT id, proyecto_id, titulo, descripcion, prioridad, estado, valor_negocio, estimacion_sp, sprint_id
		 FROM backlog_items WHERE sprint_id = ? ORDER BY id`, sprintID)
	if err != nil {
		return nil, fmt.Errorf("listar historias del sprint %d: %w", sprintID, err)
	}
	defer filas.Close()

	var historias []domain.BacklogItem
	for filas.Next() {
		var (
			h                 domain.BacklogItem
			prioridad, estado string
		)
		if err := filas.Scan(&h.ID, &h.ProyectoID, &h.Titulo, &h.Descripcion, &prioridad, &estado,
			&h.ValorNegocio, &h.EstimacionSP, &h.SprintID); err != nil {
			return nil, fmt.Errorf("leer historia: %w", err)
		}
		h.Prioridad, h.Estado = domain.Prioridad(prioridad), domain.Estado(estado)
		historias = append(historias, h)
	}
	return historias, filas.Err()
}

// QuitarDeSprint devuelve la historia al Product Backlog (sin Sprint asignado).
func (r *SQLiteBacklogRepository) QuitarDeSprint(ctx context.Context, historiaID int64) error {
	if _, err := r.db.ExecContext(ctx, "UPDATE backlog_items SET sprint_id = NULL WHERE id = ?", historiaID); err != nil {
		return fmt.Errorf("quitar historia %d del sprint: %w", historiaID, err)
	}
	return nil
}

func punteroAValorSQL(v *int) any {
	if v == nil {
		return nil
	}
	return *v
}
