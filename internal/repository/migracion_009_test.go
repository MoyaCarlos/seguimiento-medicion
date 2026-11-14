package repository

import (
	"context"
	"database/sql"
	"testing"
)

// La migración 009 (spec 045, #175) agrega integrantes.fecha_baja: los
// integrantes que ya existían quedan activos (fecha_baja NULL).
func TestMigracion009_LosIntegrantesPreviosQuedanActivos(t *testing.T) {
	db := baseSinLaMigracion(t, "009_")
	if _, err := db.Exec(`
INSERT INTO proyectos (id, nombre, creado_en) VALUES (1, 'P1', '2026-10-01T00:00:00Z');
INSERT INTO usuarios (id, nombre, nombre_normalizado, creado_en) VALUES
    (1, 'Ana', 'ana', '2026-10-01T00:00:00Z'), (2, 'Beto', 'beto', '2026-10-01T00:00:00Z');
INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES
    (1, 1, 'scrum_master', '2026-10-01T10:00:00Z'),
    (1, 2, 'product_builder', '2026-10-02T10:00:00Z');`); err != nil {
		t.Fatalf("no se pudieron sembrar los integrantes: %v", err)
	}

	if err := AplicarMigraciones(context.Background(), db); err != nil {
		t.Fatalf("la migración 009 no debe fallar sobre una base con integrantes: %v", err)
	}

	filas, err := db.Query("SELECT usuario_id, rol, fecha_baja FROM integrantes ORDER BY usuario_id")
	if err != nil {
		t.Fatalf("no se pudieron leer los integrantes: %v", err)
	}
	defer filas.Close()
	cantidad := 0
	for filas.Next() {
		var (
			usuario int64
			rol     string
			baja    sql.NullString
		)
		if err := filas.Scan(&usuario, &rol, &baja); err != nil {
			t.Fatalf("no se pudo leer el integrante: %v", err)
		}
		if baja.Valid {
			t.Errorf("usuario %d (%s): se esperaba fecha_baja NULL, se obtuvo %q", usuario, rol, baja.String)
		}
		cantidad++
	}
	if cantidad != 2 {
		t.Errorf("la migración no debe borrar integrantes: se esperaban 2 y hay %d", cantidad)
	}
}
