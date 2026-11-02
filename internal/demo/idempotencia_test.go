package demo

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
)

var tablasDeLaCarga = []string{"usuarios", "proyectos", "integrantes", "sprints", "historias",
	"historias_arrastradas", "registros_horas", "rondas_votacion", "votos"}

func fotoDeLaBase(t *testing.T, db *sql.DB) map[string]int {
	t.Helper()
	foto := map[string]int{}
	for _, tabla := range tablasDeLaCarga {
		foto[tabla] = contar(t, db, `SELECT COUNT(*) FROM `+tabla)
	}
	return foto
}

func cargador(db *sql.DB) *Cargador {
	return NewCargador(db, func() time.Time { return ahoraFijo })
}

func TestCargar_UnaSegundaCargaFallaConYaCargadoYNoCambiaNada(t *testing.T) {
	db := baseMigrada(t)
	if _, err := cargador(db).Cargar(context.Background()); err != nil {
		t.Fatal(err)
	}
	antes := fotoDeLaBase(t, db)

	_, err := cargador(db).Cargar(context.Background())

	if !errors.Is(err, ErrYaCargado) {
		t.Fatalf("error %v, se esperaba ErrYaCargado", err)
	}
	if err.Error() != "los datos de prueba ya están cargados" {
		t.Errorf("mensaje: %q", err)
	}
	for tabla, filas := range fotoDeLaBase(t, db) {
		if filas != antes[tabla] {
			t.Errorf("%s cambió: %d -> %d", tabla, antes[tabla], filas)
		}
	}
}

func TestCargar_ConElId25DeOtroProyectoFallaConIdOcupadoYNoCambiaNada(t *testing.T) {
	db := baseMigrada(t)
	if _, err := db.Exec(`INSERT INTO proyectos (id, nombre, descripcion, creado_en) VALUES (25, 'Proyecto real', '', '2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	antes := fotoDeLaBase(t, db)

	_, err := cargador(db).Cargar(context.Background())

	if !errors.Is(err, ErrIdOcupado) {
		t.Fatalf("error %v, se esperaba ErrIdOcupado", err)
	}
	if err.Error() != "el id 25 ya lo usa otro proyecto" {
		t.Errorf("mensaje: %q", err)
	}
	for tabla, filas := range fotoDeLaBase(t, db) {
		if filas != antes[tabla] {
			t.Errorf("%s cambió: %d -> %d (ni los usuarios Beta deben quedar)", tabla, antes[tabla], filas)
		}
	}
}
