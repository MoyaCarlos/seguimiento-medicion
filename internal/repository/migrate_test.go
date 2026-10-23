package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"
)

func TestDSNConPragmas_SinParametros(t *testing.T) {
	got := dsnConPragmas("file:app.db")
	want := "file:app.db?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
	if got != want {
		t.Errorf("se esperaba %q, se obtuvo %q", want, got)
	}
}

func TestDSNConPragmas_ConParametros(t *testing.T) {
	got := dsnConPragmas("file:app.db?cache=shared")
	want := "file:app.db?cache=shared&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
	if got != want {
		t.Errorf("se esperaba %q, se obtuvo %q", want, got)
	}
}

func TestEsDSNEnMemoria(t *testing.T) {
	casos := map[string]bool{
		":memory:":                   true,
		"file::memory:":              true,
		"file::memory:?cache=shared": true,
		"file:app.db":                false,
		"app.db":                     false,
	}
	for dsn, esperado := range casos {
		if esDSNEnMemoria(dsn) != esperado {
			t.Errorf("esDSNEnMemoria(%q) se esperaba %v", dsn, esperado)
		}
	}
}

func TestAbrirSQLite_PragmasEnArchivoReal(t *testing.T) {
	archivo := filepath.Join(t.TempDir(), "prueba.db")
	db, err := AbrirSQLite(archivo)
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite: %v", err)
	}
	defer db.Close()

	var fk, busy int
	if err := db.QueryRow("PRAGMA foreign_keys").Scan(&fk); err != nil {
		t.Fatalf("no se pudo consultar foreign_keys: %v", err)
	}
	if err := db.QueryRow("PRAGMA busy_timeout").Scan(&busy); err != nil {
		t.Fatalf("no se pudo consultar busy_timeout: %v", err)
	}
	if fk != 1 {
		t.Errorf("se esperaba foreign_keys = 1, se obtuvo %d", fk)
	}
	if busy != 5000 {
		t.Errorf("se esperaba busy_timeout = 5000, se obtuvo %d", busy)
	}
}

func TestAbrirSQLite_DSNConParametrosPrevios(t *testing.T) {
	archivo := filepath.Join(t.TempDir(), "prueba.db")
	dsn := "file:" + filepath.ToSlash(archivo) + "?cache=shared"
	db, err := AbrirSQLite(dsn)
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite: %v", err)
	}
	defer db.Close()

	var fk int
	if err := db.QueryRow("PRAGMA foreign_keys").Scan(&fk); err != nil {
		t.Fatalf("no se pudo consultar foreign_keys: %v", err)
	}
	if fk != 1 {
		t.Errorf("se esperaba foreign_keys = 1, se obtuvo %d", fk)
	}
}

func TestAbrirSQLite_EscriturasConcurrentesSinLock(t *testing.T) {
	archivo := filepath.Join(t.TempDir(), "prueba.db")
	db, err := AbrirSQLite(archivo)
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite: %v", err)
	}
	defer db.Close()
	ctx := context.Background()
	if err := AplicarMigraciones(ctx, db); err != nil {
		t.Fatalf("no se pudo migrar: %v", err)
	}

	const escrituras = 30
	errc := make(chan error, escrituras)
	var wg sync.WaitGroup
	for i := 0; i < escrituras; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := db.ExecContext(ctx,
				`INSERT INTO proyectos (nombre, descripcion, creado_en) VALUES (?, '', ?)`,
				fmt.Sprintf("proyecto-%d", i), time.Now().UTC().Format(time.RFC3339),
			)
			if err != nil {
				errc <- err
			}
		}(i)
	}
	wg.Wait()
	close(errc)
	for err := range errc {
		if strings.Contains(err.Error(), "database is locked") {
			t.Fatalf("no se esperaba 'database is locked': %v", err)
		}
		t.Errorf("escritura concurrente falló: %v", err)
	}
}

func nuevaBDVacia(t *testing.T) *sql.DB {
	t.Helper()
	db, err := AbrirSQLite(":memory:")
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func migracion(sql string) *fstest.MapFile {
	return &fstest.MapFile{Data: []byte(sql)}
}

func versionesAplicadas(t *testing.T, db *sql.DB) []int {
	t.Helper()
	filas, err := db.Query("SELECT version FROM migraciones_aplicadas ORDER BY version")
	if err != nil {
		t.Fatalf("no se pudo leer migraciones_aplicadas: %v", err)
	}
	defer filas.Close()
	var versiones []int
	for filas.Next() {
		var v int
		if err := filas.Scan(&v); err != nil {
			t.Fatalf("no se pudo leer la versión: %v", err)
		}
		versiones = append(versiones, v)
	}
	return versiones
}

func existeTabla(t *testing.T, db *sql.DB, nombre string) bool {
	t.Helper()
	var total int
	if err := db.QueryRow(
		"SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?", nombre,
	).Scan(&total); err != nil {
		t.Fatalf("no se pudo consultar sqlite_master: %v", err)
	}
	return total > 0
}

func TestAplicarMigraciones_EnOrdenYRegistra(t *testing.T) {
	db := nuevaBDVacia(t)
	fsys := fstest.MapFS{
		"001_crear.sql":    migracion("CREATE TABLE t1 (x INTEGER);"),
		"002_insertar.sql": migracion("INSERT INTO t1 (x) VALUES (1);"),
	}

	if err := aplicarMigraciones(context.Background(), db, fsys); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if got := versionesAplicadas(t, db); !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("se esperaban las versiones [1 2], se obtuvo %v", got)
	}
	filas, err := db.Query("SELECT aplicada_en FROM migraciones_aplicadas")
	if err != nil {
		t.Fatalf("no se pudo leer aplicada_en: %v", err)
	}
	defer filas.Close()
	for filas.Next() {
		var aplicadaEn string
		if err := filas.Scan(&aplicadaEn); err != nil {
			t.Fatalf("no se pudo leer aplicada_en: %v", err)
		}
		if _, err := time.Parse(time.RFC3339, aplicadaEn); err != nil {
			t.Errorf("aplicada_en %q no es RFC 3339: %v", aplicadaEn, err)
		}
	}
}

func TestAplicarMigraciones_Idempotente(t *testing.T) {
	db := nuevaBDVacia(t)
	ctx := context.Background()
	fsys := fstest.MapFS{
		"001_crear.sql":    migracion("CREATE TABLE t1 (x INTEGER);"),
		"002_insertar.sql": migracion("INSERT INTO t1 (x) VALUES (1);"),
	}

	for i := 0; i < 2; i++ {
		if err := aplicarMigraciones(ctx, db, fsys); err != nil {
			t.Fatalf("corrida %d: no se esperaba error: %v", i+1, err)
		}
	}

	var total int
	if err := db.QueryRow("SELECT COUNT(*) FROM t1").Scan(&total); err != nil {
		t.Fatalf("no se pudo contar t1: %v", err)
	}
	if total != 1 {
		t.Fatalf("se esperaba una sola fila en t1, se obtuvo %d", total)
	}
}

func TestAplicarMigraciones_ToleraHuecos(t *testing.T) {
	db := nuevaBDVacia(t)
	fsys := fstest.MapFS{
		"001_crear.sql": migracion("CREATE TABLE t1 (x INTEGER);"),
		"003_otra.sql":  migracion("CREATE TABLE t3 (x INTEGER);"),
	}

	if err := aplicarMigraciones(context.Background(), db, fsys); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if got := versionesAplicadas(t, db); !slices.Equal(got, []int{1, 3}) {
		t.Fatalf("se esperaban las versiones [1 3], se obtuvo %v", got)
	}
}

func TestAplicarMigraciones_SoloLasNuevas(t *testing.T) {
	db := nuevaBDVacia(t)
	ctx := context.Background()
	primera := migracion("CREATE TABLE t1 (x INTEGER);")

	if err := aplicarMigraciones(ctx, db, fstest.MapFS{"001_crear.sql": primera}); err != nil {
		t.Fatalf("primera corrida: no se esperaba error: %v", err)
	}
	err := aplicarMigraciones(ctx, db, fstest.MapFS{
		"001_crear.sql": primera,
		"003_otra.sql":  migracion("CREATE TABLE t3 (x INTEGER);"),
	})
	if err != nil {
		t.Fatalf("segunda corrida: no se esperaba error (¿se re-ejecutó la 001?): %v", err)
	}

	if !existeTabla(t, db, "t3") {
		t.Fatal("se esperaba que la 003 creara t3")
	}
	if got := versionesAplicadas(t, db); !slices.Equal(got, []int{1, 3}) {
		t.Fatalf("se esperaban las versiones [1 3], se obtuvo %v", got)
	}
}

func TestAplicarMigraciones_FallaALaMitad(t *testing.T) {
	db := nuevaBDVacia(t)
	fsys := fstest.MapFS{
		"001_crear.sql": migracion("CREATE TABLE t1 (x INTEGER);"),
		"002_rota.sql":  migracion("CREATE TABLE b (x); CREATE TABLE b (x);"),
	}

	err := aplicarMigraciones(context.Background(), db, fsys)

	if err == nil || !strings.Contains(err.Error(), "002_rota.sql") {
		t.Fatalf("se esperaba un error que nombre 002_rota.sql, se obtuvo %v", err)
	}
	if existeTabla(t, db, "b") {
		t.Error("la tabla b no debía quedar creada")
	}
	if got := versionesAplicadas(t, db); !slices.Equal(got, []int{1}) {
		t.Fatalf("se esperaba solo la versión 1 registrada, se obtuvo %v", got)
	}
}

func TestAplicarMigraciones_NombreInvalido(t *testing.T) {
	db := nuevaBDVacia(t)
	fsys := fstest.MapFS{"1_x.sql": migracion("CREATE TABLE t1 (x INTEGER);")}

	if err := aplicarMigraciones(context.Background(), db, fsys); err == nil {
		t.Fatal("se esperaba un error por nombre de migración inválido")
	}
}

func TestAplicarMigraciones_VersionRepetida(t *testing.T) {
	db := nuevaBDVacia(t)
	fsys := fstest.MapFS{
		"002_a.sql": migracion("CREATE TABLE t1 (x INTEGER);"),
		"002_b.sql": migracion("CREATE TABLE t2 (x INTEGER);"),
	}

	err := aplicarMigraciones(context.Background(), db, fsys)

	if err == nil {
		t.Fatal("se esperaba un error por versión repetida")
	}
	if !strings.Contains(err.Error(), "002_a.sql") || !strings.Contains(err.Error(), "002_b.sql") {
		t.Errorf("el error debe nombrar ambos archivos, got %q", err)
	}
}

func TestAplicarMigraciones_IgnoraVersionDesconocida(t *testing.T) {
	db := nuevaBDVacia(t)
	ctx := context.Background()
	if _, err := db.Exec(`CREATE TABLE migraciones_aplicadas (version INTEGER PRIMARY KEY, aplicada_en TEXT NOT NULL);
		INSERT INTO migraciones_aplicadas (version, aplicada_en) VALUES (9, '2026-10-23T00:00:00Z');`); err != nil {
		t.Fatalf("no se pudo preparar la tabla de control: %v", err)
	}
	fsys := fstest.MapFS{
		"001_crear.sql":    migracion("CREATE TABLE t1 (x INTEGER);"),
		"002_insertar.sql": migracion("INSERT INTO t1 (x) VALUES (1);"),
	}

	if err := aplicarMigraciones(ctx, db, fsys); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if got := versionesAplicadas(t, db); !slices.Equal(got, []int{1, 2, 9}) {
		t.Fatalf("se esperaban las versiones [1 2 9], se obtuvo %v", got)
	}
}

func TestAplicarMigraciones_BaseVieja(t *testing.T) {
	casos := map[string]string{
		"projects":        `CREATE TABLE projects (id INTEGER PRIMARY KEY)`,
		"users":           `CREATE TABLE users (id INTEGER PRIMARY KEY)`,
		"project_members": `CREATE TABLE project_members (project_id INTEGER, user_id INTEGER)`,
		"backlog_items":   `CREATE TABLE backlog_items (id INTEGER PRIMARY KEY)`,
	}
	for tabla, ddl := range casos {
		t.Run(tabla, func(t *testing.T) {
			db := nuevaBDVacia(t)
			ctx := context.Background()
			if _, err := db.ExecContext(ctx, ddl); err != nil {
				t.Fatalf("no se pudo preparar la base vieja: %v", err)
			}

			err := AplicarMigraciones(ctx, db)

			if !errors.Is(err, ErrEsquemaViejo) {
				t.Fatalf("se esperaba ErrEsquemaViejo, se obtuvo %v", err)
			}
			if !strings.Contains(err.Error(), "seguimiento.db") {
				t.Errorf("el error debe nombrar seguimiento.db, se obtuvo %v", err)
			}
			if existeTabla(t, db, "migraciones_aplicadas") {
				t.Error("no se esperaba migraciones_aplicadas en una base vieja")
			}
		})
	}
}

func TestVerificarBaseVieja_SoloSprints(t *testing.T) {
	db := nuevaBDVacia(t)
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `CREATE TABLE sprints (id INTEGER PRIMARY KEY)`); err != nil {
		t.Fatalf("no se pudo preparar la base: %v", err)
	}

	if err := verificarBaseVieja(ctx, db); err != nil {
		t.Fatalf("una base con solo sprints no es vieja, se obtuvo %v", err)
	}
}
