package repository

import (
	"context"
	"database/sql"
	"io/fs"
	"testing"
	"testing/fstest"
)

type clave struct{ proyecto, usuario int }

// baseAntesDeLa007 migra una base en memoria con todas las migraciones
// embebidas salvo la 007, para sembrar datos con el esquema anterior.
func baseAntesDeLa007(t *testing.T) *sql.DB {
	t.Helper()
	return baseSinLaMigracion(t, "007_")
}

// baseSinLaMigracion migra una base en memoria con todas las migraciones
// embebidas salvo la que empieza con prefijo y las posteriores.
func baseSinLaMigracion(t *testing.T, prefijo string) *sql.DB {
	t.Helper()
	db := nuevaBDVacia(t)
	previas := fstest.MapFS{}
	entradas, err := fs.ReadDir(migracionesFS, "migraciones")
	if err != nil {
		t.Fatalf("no se pudieron listar las migraciones: %v", err)
	}
	for _, e := range entradas {
		if e.Name() >= prefijo {
			continue
		}
		datos, err := migracionesFS.ReadFile("migraciones/" + e.Name())
		if err != nil {
			t.Fatalf("no se pudo leer %s: %v", e.Name(), err)
		}
		previas[e.Name()] = &fstest.MapFile{Data: datos}
	}
	if err := aplicarMigraciones(context.Background(), db, previas); err != nil {
		t.Fatalf("no se pudieron aplicar las migraciones previas: %v", err)
	}
	return db
}

// Proyecto 1: tres Scrum Master; el más antiguo (usuario 3) se inserta último.
// Proyecto 2: empate de creado_en entre los usuarios 4 y 2 (se inserta primero el 4).
// Proyecto 3: sin Scrum Master. Proyecto 4: un solo Scrum Master.
const integrantesConVariosScrumMaster = `
INSERT INTO proyectos (id, nombre, creado_en) VALUES (1, 'P1', '2026-10-01T00:00:00Z'), (2, 'P2', '2026-10-01T00:00:00Z'), (3, 'P3', '2026-10-01T00:00:00Z'), (4, 'P4', '2026-10-01T00:00:00Z');
INSERT INTO usuarios (id, nombre, nombre_normalizado, creado_en) VALUES
    (1, 'a', 'a', '2026-10-01T00:00:00Z'), (2, 'b', 'b', '2026-10-01T00:00:00Z'),
    (3, 'c', 'c', '2026-10-01T00:00:00Z'), (4, 'd', 'd', '2026-10-01T00:00:00Z');
INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES
    (1, 1, 'scrum_master', '2026-10-02T10:00:00Z'),
    (1, 2, 'scrum_master', '2026-10-03T10:00:00Z'),
    (1, 3, 'scrum_master', '2026-10-01T10:00:00Z'),
    (2, 4, 'scrum_master', '2026-10-01T10:00:00Z'),
    (2, 2, 'scrum_master', '2026-10-01T10:00:00Z'),
    (2, 3, 'product_builder', '2026-10-01T10:00:00Z'),
    (3, 1, 'product_builder', '2026-10-01T10:00:00Z'),
    (4, 1, 'scrum_master', '2026-10-01T10:00:00Z');
`

func rolesDeIntegrantes(t *testing.T, db *sql.DB) map[clave]string {
	t.Helper()
	filas, err := db.Query("SELECT proyecto_id, usuario_id, rol FROM integrantes")
	if err != nil {
		t.Fatalf("no se pudieron leer los integrantes: %v", err)
	}
	defer filas.Close()
	roles := map[clave]string{}
	for filas.Next() {
		var k clave
		var rol string
		if err := filas.Scan(&k.proyecto, &k.usuario, &rol); err != nil {
			t.Fatalf("no se pudo leer el integrante: %v", err)
		}
		roles[k] = rol
	}
	return roles
}

func TestMigracion007_CorrigeProyectosConVariosScrumMaster(t *testing.T) {
	db := baseAntesDeLa007(t)
	if _, err := db.Exec(integrantesConVariosScrumMaster); err != nil {
		t.Fatalf("no se pudieron sembrar los integrantes: %v", err)
	}

	if err := AplicarMigraciones(context.Background(), db); err != nil {
		t.Fatalf("la migración 007 no debe fallar sobre una base con duplicados: %v", err)
	}

	esperados := map[clave]string{
		{1, 1}: "product_builder",
		{1, 2}: "product_builder",
		{1, 3}: "scrum_master",
		{2, 4}: "product_builder",
		{2, 2}: "scrum_master",
		{2, 3}: "product_builder",
		{3, 1}: "product_builder",
		{4, 1}: "scrum_master",
	}
	obtenidos := rolesDeIntegrantes(t, db)
	if len(obtenidos) != len(esperados) {
		t.Errorf("la migración no debe borrar integrantes: se esperaban %d filas y hay %d", len(esperados), len(obtenidos))
	}
	for k, rol := range esperados {
		if obtenidos[k] != rol {
			t.Errorf("proyecto %d, usuario %d: se esperaba %q y se obtuvo %q", k.proyecto, k.usuario, rol, obtenidos[k])
		}
	}
}

func TestMigracion007_IndiceRechazaUnSegundoScrumMaster(t *testing.T) {
	db := baseAntesDeLa007(t)
	if _, err := db.Exec(integrantesConVariosScrumMaster); err != nil {
		t.Fatalf("no se pudieron sembrar los integrantes: %v", err)
	}
	if err := AplicarMigraciones(context.Background(), db); err != nil {
		t.Fatalf("no se pudo aplicar la migración 007: %v", err)
	}

	_, err := db.Exec(`INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (4, 2, 'scrum_master', '2026-10-05T00:00:00Z')`)

	if got := codigoDeViolacion(t, err); got != violaUnique {
		t.Fatalf("se esperaba el código %d (violación del índice), se obtuvo %d (%v)", violaUnique, got, err)
	}
}
