package repository

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMigrar_EsquemaViejoIDsTexto(t *testing.T) {
	db, err := AbrirSQLite(":memory:")
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite: %v", err)
	}
	defer db.Close()
	ctx := context.Background()

	// Esquema viejo: projects.id es TEXT (previo a la migración a int64).
	if _, err := db.ExecContext(ctx, `CREATE TABLE projects (
		id          TEXT PRIMARY KEY,
		name        TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		start_date  TEXT,
		end_date    TEXT,
		created_at  TEXT NOT NULL
	)`); err != nil {
		t.Fatalf("no se pudo preparar el esquema viejo: %v", err)
	}

	err = Migrar(ctx, db)
	if err == nil {
		t.Fatal("se esperaba un error por esquema viejo (IDs TEXT)")
	}
	if !strings.Contains(err.Error(), "esquema viejo") {
		t.Fatalf("el error debe explicar el esquema viejo, se obtuvo %v", err)
	}
}

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
	if err := Migrar(ctx, db); err != nil {
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
				`INSERT INTO projects (name, description, created_at) VALUES (?, '', ?)`,
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
