package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
		`INSERT INTO historias (proyecto_id, titulo, descripcion, criterios_aceptacion, prioridad, estado, valor_negocio, estimacion_sp)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ProyectoID,
		item.Titulo,
		item.Descripcion,
		punteroATextoSQL(item.CriteriosAceptacion),
		string(item.Prioridad),
		string(item.Estado),
		punteroAValorSQL(item.ValorNegocio),
		punteroAValorSQL(item.EstimacionSP),
	)
	if err != nil {
		return domain.BacklogItem{}, fmt.Errorf("guardar backlog item: %w", errorDeEscritura(err))
	}

	id, err := resultado.LastInsertId()
	if err != nil {
		return domain.BacklogItem{}, fmt.Errorf("obtener id del backlog item: %w", err)
	}
	item.ID = id
	return item, nil
}

const columnasHistoria = `id, proyecto_id, titulo, descripcion, criterios_aceptacion, prioridad, estado, valor_negocio, estimacion_sp, sprint_id, fecha_finalizacion`

// ListarPorSprint devuelve las historias asignadas al Sprint, en orden de creación.
func (r *SQLiteBacklogRepository) ListarPorSprint(ctx context.Context, sprintID int64) ([]domain.BacklogItem, error) {
	historias, err := r.listarHistorias(ctx,
		`SELECT `+columnasHistoria+` FROM historias WHERE sprint_id = ? ORDER BY id`, sprintID)
	if err != nil {
		return nil, fmt.Errorf("listar historias del sprint %d: %w", sprintID, err)
	}
	return historias, nil
}

// ListarSinSprint devuelve el Product Backlog del proyecto: las historias sin
// Sprint, por prioridad M, S, C, W y luego por orden de creación. El orden
// alfabético de la prioridad no sirve, por eso el CASE.
func (r *SQLiteBacklogRepository) ListarSinSprint(ctx context.Context, proyectoID int64) ([]domain.BacklogItem, error) {
	historias, err := r.listarHistorias(ctx,
		`SELECT `+columnasHistoria+` FROM historias WHERE proyecto_id = ? AND sprint_id IS NULL
		 ORDER BY CASE prioridad WHEN 'M' THEN 0 WHEN 'S' THEN 1 WHEN 'C' THEN 2 WHEN 'W' THEN 3 END, id`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar product backlog del proyecto %d: %w", proyectoID, err)
	}
	return historias, nil
}

// ObtenerPorID devuelve la historia con el identificador indicado.
func (r *SQLiteBacklogRepository) ObtenerPorID(ctx context.Context, id int64) (domain.BacklogItem, error) {
	return escanearHistoria(r.db.QueryRowContext(ctx,
		`SELECT `+columnasHistoria+` FROM historias WHERE id = ?`, id))
}

func (r *SQLiteBacklogRepository) listarHistorias(ctx context.Context, consulta string, args ...any) ([]domain.BacklogItem, error) {
	filas, err := r.db.QueryContext(ctx, consulta, args...)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	historias := make([]domain.BacklogItem, 0)
	for filas.Next() {
		h, err := escanearHistoria(filas)
		if err != nil {
			return nil, err
		}
		historias = append(historias, h)
	}
	return historias, filas.Err()
}

func escanearHistoria(fila escaneable) (domain.BacklogItem, error) {
	var (
		h                 domain.BacklogItem
		prioridad, estado string
		finalizacion      sql.NullString
	)
	if err := fila.Scan(&h.ID, &h.ProyectoID, &h.Titulo, &h.Descripcion, &h.CriteriosAceptacion, &prioridad, &estado,
		&h.ValorNegocio, &h.EstimacionSP, &h.SprintID, &finalizacion); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.BacklogItem{}, domain.ErrHistoriaNoEncontrada
		}
		return domain.BacklogItem{}, fmt.Errorf("leer historia: %w", err)
	}
	h.Prioridad, h.Estado = domain.Prioridad(prioridad), domain.Estado(estado)
	fecha, err := fechaDesdeSQL(finalizacion)
	if err != nil {
		return domain.BacklogItem{}, fmt.Errorf("leer fecha de finalización de la historia %d: %w", h.ID, err)
	}
	if !fecha.IsZero() {
		h.FechaFinalizacion = &fecha
	}
	return h, nil
}

// Actualizar persiste los campos editables de la historia: título, descripción,
// criterios de aceptación y Story Points.
func (r *SQLiteBacklogRepository) Actualizar(ctx context.Context, item domain.BacklogItem) error {
	if _, err := r.db.ExecContext(ctx,
		"UPDATE historias SET titulo = ?, descripcion = ?, criterios_aceptacion = ?, estimacion_sp = ? WHERE id = ?",
		item.Titulo, item.Descripcion, punteroATextoSQL(item.CriteriosAceptacion), punteroAValorSQL(item.EstimacionSP), item.ID,
	); err != nil {
		return fmt.Errorf("actualizar historia %d: %w", item.ID, err)
	}
	return nil
}

// AsignarASprint asigna la historia al Sprint solo si todavía no tiene uno. La
// condición va en el mismo UPDATE: si dos asignaciones compiten, solo una
// encuentra sprint_id IS NULL y la otra recibe domain.ErrHistoriaYaAsignada.
func (r *SQLiteBacklogRepository) AsignarASprint(ctx context.Context, historiaID, sprintID int64) error {
	resultado, err := r.db.ExecContext(ctx,
		"UPDATE historias SET sprint_id = ? WHERE id = ? AND sprint_id IS NULL", sprintID, historiaID)
	if err != nil {
		return fmt.Errorf("asignar historia %d al sprint %d: %w", historiaID, sprintID, errorDeEscritura(err))
	}
	filas, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("asignar historia %d al sprint %d: %w", historiaID, sprintID, err)
	}
	if filas == 0 {
		return domain.ErrHistoriaYaAsignada
	}
	return nil
}

// CambiarEstado pasa la historia de "desde" a "hacia" solo si todavía está en
// "desde". La condición va en el mismo UPDATE: si dos cambios compiten, solo uno
// la encuentra en ese estado y el otro recibe domain.ErrConflictoEstadoHistoria.
// Sin fecha, fecha_finalizacion queda NULL.
func (r *SQLiteBacklogRepository) CambiarEstado(ctx context.Context, id int64, desde, hacia domain.Estado, fechaFinalizacion *time.Time) error {
	var fecha time.Time
	if fechaFinalizacion != nil {
		fecha = *fechaFinalizacion
	}
	resultado, err := r.db.ExecContext(ctx,
		"UPDATE historias SET estado = ?, fecha_finalizacion = ? WHERE id = ? AND estado = ?",
		string(hacia), fechaASQL(fecha), id, string(desde))
	if err != nil {
		return fmt.Errorf("cambiar estado de la historia %d: %w", id, errorDeEscritura(err))
	}
	filas, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("cambiar estado de la historia %d: %w", id, err)
	}
	if filas == 0 {
		return domain.ErrConflictoEstadoHistoria
	}
	return nil
}

func punteroATextoSQL(v *string) any {
	if v == nil {
		return nil
	}
	return *v
}

func punteroAValorSQL(v *int) any {
	if v == nil {
		return nil
	}
	return *v
}
