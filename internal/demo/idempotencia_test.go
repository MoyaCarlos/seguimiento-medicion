package demo

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"
)

var tablasDeLaCarga = []string{"usuarios", "proyectos", "integrantes", "sprints", "historias",
	"historias_arrastradas", "registros_horas", "rondas_votacion", "votos", "defectos"}

func fotoDeLaBase(t *testing.T, db *sql.DB) map[string]int {
	t.Helper()
	foto := map[string]int{}
	for _, tabla := range tablasDeLaCarga {
		foto[tabla] = contar(t, db, `SELECT COUNT(*) FROM `+tabla)
	}
	return foto
}

func cargador(db *sql.DB) *Cargador {
	return NewCargador(db, func() time.Time { return ahoraFijo }, hasheadorBarato)
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
			t.Errorf("%s cambió: %d -> %d (ningún usuario de la demo debe quedar)", tabla, antes[tabla], filas)
		}
	}
}

func insertarUsuario(t *testing.T, db *sql.DB, nombre string) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES (?, ?, '2026-01-01T00:00:00Z')`,
		nombre, strings.ToLower(nombre)); err != nil {
		t.Fatal(err)
	}
}

func cargarYEsperarUsuarioExistente(t *testing.T, db *sql.DB, nombreEsperado string) {
	t.Helper()
	antes := fotoDeLaBase(t, db)

	_, err := cargador(db).Cargar(context.Background())

	var existente ErrUsuarioExistente
	if !errors.As(err, &existente) {
		t.Fatalf("error %v, se esperaba ErrUsuarioExistente", err)
	}
	if existente.Nombre != nombreEsperado {
		t.Errorf("usuario en conflicto %q, se esperaba %q", existente.Nombre, nombreEsperado)
	}
	if err.Error() != "ya existe el usuario "+nombreEsperado {
		t.Errorf("mensaje: %q", err)
	}
	for tabla, filas := range fotoDeLaBase(t, db) {
		if filas != antes[tabla] {
			t.Errorf("%s cambió: %d -> %d", tabla, antes[tabla], filas)
		}
	}
}

func TestCargar_UnUsuarioDeLaDemoExistenteFallaNombrandoloYNoCambiaNada(t *testing.T) {
	db := baseMigrada(t)
	insertarUsuario(t, db, "ana.iriarte")

	cargarYEsperarUsuarioExistente(t, db, "ana.iriarte")
}

func TestCargar_ElConflictoDeUsuarioIgnoraMayusculas(t *testing.T) {
	db := baseMigrada(t)
	insertarUsuario(t, db, "Ana.Iriarte")

	cargarYEsperarUsuarioExistente(t, db, "ana.iriarte")
}

func TestCargar_ConVariosUsuariosExistentesNombraAlPrimeroDeLaLista(t *testing.T) {
	db := baseMigrada(t)
	insertarUsuario(t, db, "jimena.martinez")
	insertarUsuario(t, db, "candela.vulcano")

	cargarYEsperarUsuarioExistente(t, db, "candela.vulcano")
}

func TestCargar_UsuarioExistenteGanaAlIdOcupadoPorOtroProyecto(t *testing.T) {
	db := baseMigrada(t)
	insertarUsuario(t, db, "carlos.moya")
	if _, err := db.Exec(`INSERT INTO proyectos (id, nombre, descripcion, creado_en) VALUES (25, 'Proyecto real', '', '2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}

	cargarYEsperarUsuarioExistente(t, db, "carlos.moya")
}

func TestCargar_YaCargadoGanaAlUsuarioExistente(t *testing.T) {
	db := baseMigrada(t)
	if _, err := cargador(db).Cargar(context.Background()); err != nil {
		t.Fatal(err)
	}

	_, err := cargador(db).Cargar(context.Background())

	if !errors.Is(err, ErrYaCargado) {
		t.Fatalf("error %v, se esperaba ErrYaCargado aunque los usuarios existan", err)
	}
}
