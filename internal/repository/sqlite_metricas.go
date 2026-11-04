package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// El número se cuenta sobre todos los Sprints del proyecto, no solo los
// finalizados (#129): debe coincidir con el del resto de la API.
const consultaDatosMetricas = `
	SELECT
		s.id,
		(SELECT COUNT(*) FROM sprints previos WHERE previos.proyecto_id = s.proyecto_id AND previos.id <= s.id),
		(SELECT COUNT(*) FROM historias h WHERE h.sprint_id = s.id AND h.estado = 'Completada'),
		(SELECT COUNT(*) FROM historias_arrastradas a WHERE a.sprint_id = s.id),
		(SELECT COALESCE(SUM(h.estimacion_sp), 0) FROM historias h WHERE h.sprint_id = s.id AND h.estado = 'Completada'),
		(SELECT COALESCE(SUM(a.estimacion_sp), 0) FROM historias_arrastradas a WHERE a.sprint_id = s.id),
		(SELECT COALESCE(SUM(r.minutos), 0) FROM registros_horas r WHERE r.sprint_id = s.id),
		(SELECT COUNT(*) FROM defectos d WHERE d.sprint_deteccion_id = s.id),
		(SELECT COUNT(*) FROM defectos d WHERE d.sprint_resolucion_id = s.id)
	FROM sprints s
	WHERE s.proyecto_id = ? AND s.estado = 'Finalizado'
	ORDER BY s.id`

// SQLiteMetricasRepository implementa DatosMetricasRepository sobre SQLite.
type SQLiteMetricasRepository struct {
	db *sql.DB
}

func NewSQLiteMetricasRepository(db *sql.DB) *SQLiteMetricasRepository {
	return &SQLiteMetricasRepository{db: db}
}

// ListarDatosSprintsFinalizados devuelve los datos de cada Sprint finalizado.
func (r *SQLiteMetricasRepository) ListarDatosSprintsFinalizados(ctx context.Context, proyectoID int64) ([]domain.DatosMetricasSprint, error) {
	filas, err := r.db.QueryContext(ctx, consultaDatosMetricas, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar datos de métricas: %w", err)
	}
	defer filas.Close()

	datos := make([]domain.DatosMetricasSprint, 0)
	for filas.Next() {
		var d domain.DatosMetricasSprint
		if err := filas.Scan(&d.SprintID, &d.Numero, &d.HistoriasCompletadas, &d.HistoriasArrastradas,
			&d.SPCompletados, &d.SPArrastrados, &d.MinutosReales, &d.DefectosDetectados, &d.DefectosResueltos); err != nil {
			return nil, fmt.Errorf("leer datos de métricas: %w", err)
		}
		datos = append(datos, d)
	}
	if err := filas.Err(); err != nil {
		return nil, fmt.Errorf("listar datos de métricas: %w", err)
	}
	return datos, nil
}
