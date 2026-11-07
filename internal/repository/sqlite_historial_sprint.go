package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// Las completadas toman los Story Points actuales y las arrastradas la foto
// guardada al cerrar el Sprint (D2c). El grupo ordena completadas primero.
const consultaHistoriasComprometidas = `
	SELECT id, titulo, estimacion_sp, 1 AS completada
	FROM historias
	WHERE sprint_id = ? AND estado = 'Completada'
	UNION ALL
	SELECT h.id, h.titulo, a.estimacion_sp, 0 AS completada
	FROM historias_arrastradas a
	JOIN historias h ON h.id = a.historia_id
	WHERE a.sprint_id = ?
	ORDER BY completada DESC, id`

const consultaMinutosPorHistoria = `
	SELECT historia_id, SUM(minutos)
	FROM registros_horas
	WHERE sprint_id = ?
	GROUP BY historia_id`

// SQLiteHistorialSprintRepository implementa HistorialSprintRepository sobre
// SQLite. Solo lee.
type SQLiteHistorialSprintRepository struct {
	db *sql.DB
}

func NewSQLiteHistorialSprintRepository(db *sql.DB) *SQLiteHistorialSprintRepository {
	return &SQLiteHistorialSprintRepository{db: db}
}

func (r *SQLiteHistorialSprintRepository) HistoriasComprometidas(ctx context.Context, sprintID int64) ([]domain.HistoriaComprometida, error) {
	filas, err := r.db.QueryContext(ctx, consultaHistoriasComprometidas, sprintID, sprintID)
	if err != nil {
		return nil, fmt.Errorf("listar historias comprometidas del sprint %d: %w", sprintID, err)
	}
	defer filas.Close()

	historias := make([]domain.HistoriaComprometida, 0)
	for filas.Next() {
		var (
			h          domain.HistoriaComprometida
			estimacion sql.NullInt64
		)
		if err := filas.Scan(&h.ID, &h.Titulo, &estimacion, &h.Completada); err != nil {
			return nil, fmt.Errorf("leer historia comprometida: %w", err)
		}
		if estimacion.Valid {
			sp := int(estimacion.Int64)
			h.EstimacionSP = &sp
		}
		historias = append(historias, h)
	}
	if err := filas.Err(); err != nil {
		return nil, fmt.Errorf("listar historias comprometidas del sprint %d: %w", sprintID, err)
	}
	return historias, nil
}

func (r *SQLiteHistorialSprintRepository) MinutosPorHistoria(ctx context.Context, sprintID int64) (map[int64]int, error) {
	filas, err := r.db.QueryContext(ctx, consultaMinutosPorHistoria, sprintID)
	if err != nil {
		return nil, fmt.Errorf("sumar minutos del sprint %d: %w", sprintID, err)
	}
	defer filas.Close()

	minutos := make(map[int64]int)
	for filas.Next() {
		var historiaID int64
		var total int
		if err := filas.Scan(&historiaID, &total); err != nil {
			return nil, fmt.Errorf("leer minutos por historia: %w", err)
		}
		minutos[historiaID] = total
	}
	if err := filas.Err(); err != nil {
		return nil, fmt.Errorf("sumar minutos del sprint %d: %w", sprintID, err)
	}
	return minutos, nil
}

// DefectosDeSprint reutiliza la lectura de defectos de HU-08.
func (r *SQLiteHistorialSprintRepository) DefectosDeSprint(ctx context.Context, sprintID int64) ([]domain.DefectoDetalle, error) {
	filas, err := r.db.QueryContext(ctx,
		seleccionarDefecto+` WHERE d.sprint_deteccion_id = ? OR d.sprint_resolucion_id = ? ORDER BY d.id`, sprintID, sprintID)
	if err != nil {
		return nil, fmt.Errorf("listar defectos del sprint %d: %w", sprintID, err)
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
		return nil, fmt.Errorf("listar defectos del sprint %d: %w", sprintID, err)
	}
	return defectos, nil
}
