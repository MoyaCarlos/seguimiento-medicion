package repository

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// fechaDeBaja es el 12 de noviembre en UTC: el repositorio guarda el día UTC.
var fechaDeBaja = time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)

// Proyecto 1: Ana (1) Scrum Master, Beto (2) y Carla (3) Product Builders.
// Proyecto 2: Beto (2) Product Builder. Sprint 10 del proyecto 1 Finalizado,
// con la historia 100 completada (5 SP) y 90 minutos de Beto.
// Rondas: 1000 (historia 100, revelada) y 1001 (historia 101, abierta) del
// proyecto 1; 2000 (historia 200, abierta) del proyecto 2.
const equipoConVotosYHoras = `
INSERT INTO usuarios (id, nombre, nombre_normalizado, creado_en) VALUES
    (1, 'Ana', 'ana', '2026-10-01T00:00:00Z'),
    (2, 'Beto', 'beto', '2026-10-01T00:00:00Z'),
    (3, 'Carla', 'carla', '2026-10-01T00:00:00Z');
INSERT INTO proyectos (id, nombre, creado_en) VALUES
    (1, 'Uno', '2026-10-01T00:00:00Z'), (2, 'Dos', '2026-10-01T00:00:00Z');
INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES
    (1, 1, 'scrum_master', '2026-10-01T00:00:00Z'),
    (1, 2, 'product_builder', '2026-10-02T00:00:00Z'),
    (1, 3, 'product_builder', '2026-10-03T00:00:00Z'),
    (2, 2, 'product_builder', '2026-10-02T00:00:00Z');
INSERT INTO sprints (id, proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES
    (10, 1, 'g', '2026-10-05', '2026-10-19', 'Finalizado');
INSERT INTO historias (id, proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp, sprint_id, fecha_finalizacion) VALUES
    (100, 1, 'h100', 'd', 'M', 'Completada', 5, 10, '2026-10-18'),
    (101, 1, 'h101', 'd', 'M', 'Nueva', NULL, NULL, NULL),
    (200, 2, 'h200', 'd', 'M', 'Nueva', NULL, NULL, NULL);
INSERT INTO registros_horas (proyecto_id, historia_id, usuario_id, sprint_id, fecha, actividad, minutos, creado_en) VALUES
    (1, 100, 2, 10, '2026-10-10', 'desarrollo', 90, '2026-10-10T00:00:00Z');
INSERT INTO rondas_votacion (id, historia_id, numero, revelada_en) VALUES
    (1000, 100, 1, '2026-10-06T00:00:00Z'),
    (1001, 101, 1, NULL),
    (2000, 200, 1, NULL);
INSERT INTO votos (ronda_id, usuario_id, valor) VALUES
    (1000, 2, 5), (1000, 3, 3),
    (1001, 1, 8), (1001, 2, 8), (1001, 3, 3),
    (2000, 2, 13);
`

type votoGuardado struct {
	ronda, usuario int64
	valor          int
}

func votosGuardados(t *testing.T, db *sql.DB) []votoGuardado {
	t.Helper()
	filas, err := db.Query("SELECT ronda_id, usuario_id, valor FROM votos ORDER BY ronda_id, usuario_id")
	if err != nil {
		t.Fatalf("no se pudieron leer los votos: %v", err)
	}
	defer filas.Close()
	var votos []votoGuardado
	for filas.Next() {
		var v votoGuardado
		if err := filas.Scan(&v.ronda, &v.usuario, &v.valor); err != nil {
			t.Fatalf("no se pudo leer el voto: %v", err)
		}
		votos = append(votos, v)
	}
	return votos
}

func fechaBajaGuardada(t *testing.T, db *sql.DB, proyectoID, usuarioID int64) sql.NullString {
	t.Helper()
	var baja sql.NullString
	if err := db.QueryRow("SELECT fecha_baja FROM integrantes WHERE proyecto_id = ? AND usuario_id = ?", proyectoID, usuarioID).Scan(&baja); err != nil {
		t.Fatalf("no se pudo leer la fila de integrantes (%d, %d): %v", proyectoID, usuarioID, err)
	}
	return baja
}

func TestSQLiteProjectRepository_DarDeBajaIntegrante(t *testing.T) {
	ctx := context.Background()
	preparar := func(t *testing.T) (*sql.DB, *SQLiteProjectRepository) {
		t.Helper()
		db := nuevaBDMigrada(t)
		mustExec(t, db, equipoConVotosYHoras)
		return db, NewSQLiteProjectRepository(db)
	}

	t.Run("Product Builder activo: fecha de baja y votos no revelados descartados", func(t *testing.T) {
		db, repo := preparar(t)
		metricas := NewSQLiteMetricasRepository(db)
		antes, err := metricas.ListarDatosSprintsFinalizados(ctx, 1)
		if err != nil {
			t.Fatalf("métricas antes de la baja: %v", err)
		}

		if err := repo.DarDeBajaIntegrante(ctx, 1, 2, fechaDeBaja); err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}

		if baja := fechaBajaGuardada(t, db, 1, 2); !baja.Valid || baja.String != "2026-11-12" {
			t.Errorf("se esperaba fecha_baja '2026-11-12' en la misma fila, se obtuvo %+v", baja)
		}
		if baja := fechaBajaGuardada(t, db, 2, 2); baja.Valid {
			t.Errorf("la baja no debe tocar el vínculo con otro proyecto, se obtuvo %q", baja.String)
		}
		activos, _ := repo.ListarIntegrantes(ctx, 1)
		if got := nombresDe(activos); slices.Contains(got, "Beto:product_builder") {
			t.Errorf("Beto no debe figurar entre los activos, se obtuvo %v", got)
		}
		historicos, _ := repo.ListarIntegrantesHistoricos(ctx, 1)
		if got := nombresDe(historicos); !slices.Contains(got, "Beto:product_builder") {
			t.Errorf("Beto debe seguir en el histórico, se obtuvo %v", got)
		}

		// El voto de Beto en la ronda abierta 1001 se descarta; el de la
		// revelada 1000 y el de la ronda abierta de otro proyecto (2000) quedan.
		esperados := []votoGuardado{
			{1000, 2, 5}, {1000, 3, 3},
			{1001, 1, 8}, {1001, 3, 3},
			{2000, 2, 13},
		}
		if got := votosGuardados(t, db); !reflect.DeepEqual(got, esperados) {
			t.Errorf("votos después de la baja:\n se esperaba %v\n se obtuvo   %v", esperados, got)
		}

		var horas int
		if err := db.QueryRow("SELECT COUNT(*) FROM registros_horas WHERE usuario_id = 2").Scan(&horas); err != nil || horas != 1 {
			t.Errorf("las horas de Beto deben seguir (1 registro), se obtuvo %d (%v)", horas, err)
		}
		despues, err := metricas.ListarDatosSprintsFinalizados(ctx, 1)
		if err != nil {
			t.Fatalf("métricas después de la baja: %v", err)
		}
		if !reflect.DeepEqual(antes, despues) {
			t.Errorf("las métricas no deben cambiar con la baja (SC-002):\n antes   %+v\n después %+v", antes, despues)
		}
	})

	t.Run("segunda baja del mismo", func(t *testing.T) {
		_, repo := preparar(t)
		if err := repo.DarDeBajaIntegrante(ctx, 1, 2, fechaDeBaja); err != nil {
			t.Fatalf("primera baja: %v", err)
		}
		if err := repo.DarDeBajaIntegrante(ctx, 1, 2, fechaDeBaja); !errors.Is(err, domain.ErrIntegranteNoEncontrado) {
			t.Errorf("se esperaba ErrIntegranteNoEncontrado, se obtuvo %v", err)
		}
	})

	t.Run("usuario que nunca integró", func(t *testing.T) {
		db, repo := preparar(t)
		if err := repo.DarDeBajaIntegrante(ctx, 2, 3, fechaDeBaja); !errors.Is(err, domain.ErrIntegranteNoEncontrado) {
			t.Errorf("se esperaba ErrIntegranteNoEncontrado, se obtuvo %v", err)
		}
		if got := votosGuardados(t, db); len(got) != 6 {
			t.Errorf("no debe borrar votos, quedaron %v", got)
		}
	})

	t.Run("Scrum Master: no se quita ni pierde sus votos de rondas abiertas", func(t *testing.T) {
		db, repo := preparar(t)
		if err := repo.DarDeBajaIntegrante(ctx, 1, 1, fechaDeBaja); !errors.Is(err, domain.ErrScrumMasterNoSeQuita) {
			t.Errorf("se esperaba ErrScrumMasterNoSeQuita, se obtuvo %v", err)
		}
		if baja := fechaBajaGuardada(t, db, 1, 1); baja.Valid {
			t.Errorf("el Scrum Master debe seguir activo, se obtuvo fecha_baja %q", baja.String)
		}
		if got := votosGuardados(t, db); !slices.Contains(got, votoGuardado{1001, 1, 8}) {
			t.Errorf("el voto abierto del Scrum Master debe seguir, quedaron %v", got)
		}
	})
}
