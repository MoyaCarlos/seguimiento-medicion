package repository

import (
	"context"
	"database/sql"
	"testing"
)

// La migración 008 (spec 044, #300) agrega proyectos.fecha_cierre: los
// proyectos que ya existían quedan abiertos (fecha_cierre NULL).
func TestMigracion008_LosProyectosPreviosQuedanAbiertos(t *testing.T) {
	db := baseSinLaMigracion(t, "008_")
	if _, err := db.Exec(`INSERT INTO proyectos (id, nombre, fecha_inicio, fecha_fin, creado_en) VALUES
		(1, 'P1', '2026-10-01', '2026-12-15', '2026-10-01T00:00:00Z'),
		(2, 'P2', NULL, NULL, '2026-10-01T00:00:00Z')`); err != nil {
		t.Fatalf("no se pudieron sembrar los proyectos: %v", err)
	}

	if err := AplicarMigraciones(context.Background(), db); err != nil {
		t.Fatalf("la migración 008 no debe fallar sobre una base con proyectos: %v", err)
	}

	filas, err := db.Query("SELECT id, nombre, fecha_cierre FROM proyectos ORDER BY id")
	if err != nil {
		t.Fatalf("no se pudieron leer los proyectos: %v", err)
	}
	defer filas.Close()
	var ids []int64
	for filas.Next() {
		var (
			id     int64
			nombre string
			cierre sql.NullString
		)
		if err := filas.Scan(&id, &nombre, &cierre); err != nil {
			t.Fatalf("no se pudo leer el proyecto: %v", err)
		}
		if cierre.Valid {
			t.Errorf("proyecto %d (%s): se esperaba fecha_cierre NULL, se obtuvo %q", id, nombre, cierre.String)
		}
		ids = append(ids, id)
	}
	if len(ids) != 2 {
		t.Errorf("la migración no debe borrar proyectos: se esperaban 2 y hay %d", len(ids))
	}
}
