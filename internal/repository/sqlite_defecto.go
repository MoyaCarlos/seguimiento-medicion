package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteDefectoRepository implementa DefectoRepository sobre SQLite.
type SQLiteDefectoRepository struct {
	db *sql.DB
}

func NewSQLiteDefectoRepository(db *sql.DB) *SQLiteDefectoRepository {
	return &SQLiteDefectoRepository{db: db}
}

const seleccionarDefecto = `SELECT d.id, d.proyecto_id, d.historia_id, d.descripcion, d.severidad, d.estado,
	       d.sprint_deteccion_id, d.sprint_resolucion_id, d.creado_en, h.titulo
	FROM defectos d
	JOIN historias h ON h.id = d.historia_id`

// Registrar inserta el defecto solo si su Sprint de detección sigue Activo, en
// la misma sentencia: si el Sprint se cerró entre la validación del service y
// el INSERT, el registro se rechaza (research.md R4).
func (r *SQLiteDefectoRepository) Registrar(ctx context.Context, d domain.Defecto) (domain.Defecto, error) {
	creadoEn := time.Now().UTC()
	resultado, err := r.db.ExecContext(ctx,
		`INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, creado_en)
		 SELECT ?, ?, ?, ?, 'Pendiente', ?, ?
		 WHERE EXISTS (SELECT 1 FROM sprints WHERE id = ? AND estado = 'Activo')`,
		d.ProyectoID, d.HistoriaID, d.Descripcion, string(d.Severidad), d.SprintDeteccionID, creadoEn.Format(time.RFC3339),
		d.SprintDeteccionID,
	)
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("registrar defecto: %w", errorDeEscritura(err))
	}
	insertadas, err := resultado.RowsAffected()
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("registrar defecto: %w", err)
	}
	if insertadas == 0 {
		return domain.Defecto{}, domain.ErrSinSprintActivoParaDetectar
	}
	id, err := resultado.LastInsertId()
	if err != nil {
		return domain.Defecto{}, fmt.Errorf("obtener id del defecto: %w", err)
	}
	d.ID, d.Estado, d.SprintResolucionID = id, domain.DefectoPendiente, nil
	d.CreadoEn = creadoEn.Truncate(time.Second)
	return d, nil
}

func (r *SQLiteDefectoRepository) ObtenerPorID(ctx context.Context, id int64) (domain.DefectoDetalle, error) {
	d, err := escanearDefecto(r.db.QueryRowContext(ctx, seleccionarDefecto+` WHERE d.id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.DefectoDetalle{}, domain.ErrDefectoNoEncontrado
	}
	if err != nil {
		return domain.DefectoDetalle{}, fmt.Errorf("obtener defecto %d: %w", id, err)
	}
	return d, nil
}

// Resolver marca el defecto como Resuelto solo si sigue Pendiente y su Sprint
// de resolución sigue Activo, en la misma sentencia: de dos resoluciones
// simultáneas gana una sola, y un Sprint cerrado a mitad de camino la rechaza
// (research.md R4).
func (r *SQLiteDefectoRepository) Resolver(ctx context.Context, d domain.Defecto) error {
	if d.SprintResolucionID == nil {
		return domain.ErrSinSprintActivoParaResolver
	}
	resultado, err := r.db.ExecContext(ctx,
		`UPDATE defectos SET estado = 'Resuelto', sprint_resolucion_id = ?
		 WHERE id = ? AND estado = 'Pendiente'
		   AND EXISTS (SELECT 1 FROM sprints WHERE id = ? AND estado = 'Activo')`,
		*d.SprintResolucionID, d.ID, *d.SprintResolucionID,
	)
	if err != nil {
		return fmt.Errorf("resolver defecto %d: %w", d.ID, errorDeEscritura(err))
	}
	actualizadas, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("resolver defecto %d: %w", d.ID, err)
	}
	if actualizadas > 0 {
		return nil
	}
	var estado string
	if err := r.db.QueryRowContext(ctx, `SELECT estado FROM defectos WHERE id = ?`, d.ID).Scan(&estado); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrDefectoNoEncontrado
		}
		return fmt.Errorf("releer defecto %d: %w", d.ID, err)
	}
	if domain.EstadoDefecto(estado) == domain.DefectoResuelto {
		return domain.ErrDefectoYaResuelto
	}
	return domain.ErrSinSprintActivoParaResolver
}

// ListarPorProyecto devuelve los defectos del proyecto con el título de su
// historia, en orden de registro.
func (r *SQLiteDefectoRepository) ListarPorProyecto(ctx context.Context, proyectoID int64) ([]domain.DefectoDetalle, error) {
	filas, err := r.db.QueryContext(ctx, seleccionarDefecto+` WHERE d.proyecto_id = ? ORDER BY d.id`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar defectos del proyecto %d: %w", proyectoID, err)
	}
	defer filas.Close()

	defectos := make([]domain.DefectoDetalle, 0)
	for filas.Next() {
		d, err := escanearDefecto(filas)
		if err != nil {
			return nil, fmt.Errorf("leer defecto: %w", err)
		}
		defectos = append(defectos, d)
	}
	if err := filas.Err(); err != nil {
		return nil, fmt.Errorf("listar defectos del proyecto %d: %w", proyectoID, err)
	}
	return defectos, nil
}

func escanearDefecto(fila escaneable) (domain.DefectoDetalle, error) {
	var (
		d                           domain.DefectoDetalle
		severidad, estado, creadoEn string
		resolucion                  sql.NullInt64
	)
	if err := fila.Scan(&d.ID, &d.ProyectoID, &d.HistoriaID, &d.Descripcion, &severidad, &estado,
		&d.SprintDeteccionID, &resolucion, &creadoEn, &d.TituloHistoria); err != nil {
		return domain.DefectoDetalle{}, err
	}
	d.Severidad, d.Estado = domain.Severidad(severidad), domain.EstadoDefecto(estado)
	if resolucion.Valid {
		d.SprintResolucionID = &resolucion.Int64
	}
	d.CreadoEn = parsearTiempo(creadoEn)
	return d, nil
}
