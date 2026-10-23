package repository

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// AbrirSQLite abre una conexión a SQLite (driver puro Go, sin CGO) limitada a
// una única conexión, activa las claves foráneas y configura un busy_timeout
// para las bases en archivo (vía _pragma en el DSN). Las bases en memoria
// (también en su forma URI) activan las FK por sentencia.
func AbrirSQLite(dsn string) (*sql.DB, error) {
	enMemoria := esDSNEnMemoria(dsn)

	abierto := dsn
	if !enMemoria {
		abierto = dsnConPragmas(dsn)
	}

	db, err := sql.Open("sqlite", abierto)
	if err != nil {
		return nil, fmt.Errorf("abrir sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)

	if enMemoria {
		if _, err := db.ExecContext(context.Background(), "PRAGMA foreign_keys = ON"); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("activar foreign keys: %w", err)
		}
	}
	return db, nil
}

// dsnConPragmas agrega foreign_keys y busy_timeout respetando los parámetros que
// el DSN ya pudiera traer (? para el primero, & para los siguientes).
func dsnConPragmas(dsn string) string {
	separador := "?"
	if strings.Contains(dsn, "?") {
		separador = "&"
	}
	return dsn + separador + "_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
}

// esDSNEnMemoria reconoce tanto la forma simple ":memory:" como las URI de
// memoria (p.ej. "file::memory:?cache=shared").
func esDSNEnMemoria(dsn string) bool {
	base := dsn
	if i := strings.IndexByte(dsn, '?'); i >= 0 {
		base = dsn[:i]
	}
	return base == ":memory:" || base == "file::memory:"
}

// verificarBaseVieja detecta una base creada con el esquema previo a la
// migración 001 (tablas en inglés). sprints no cuenta: conserva su nombre.
func verificarBaseVieja(ctx context.Context, db *sql.DB) error {
	var viejas int
	if err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name IN ('projects','users','project_members','backlog_items')`,
	).Scan(&viejas); err != nil {
		return fmt.Errorf("inspeccionar base: %w", err)
	}
	if viejas > 0 {
		return ErrEsquemaViejo
	}
	return nil
}

var nombreDeMigracion = regexp.MustCompile(`^(\d{3})_.+\.sql$`)

// aplicarMigraciones aplica, en orden de nombre, las migraciones NNN_nombre.sql
// de fsys que no figuran en migraciones_aplicadas, cada una en su transacción.
func aplicarMigraciones(ctx context.Context, db *sql.DB, fsys fs.FS) error {
	if err := verificarBaseVieja(ctx, db); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS migraciones_aplicadas (version INTEGER PRIMARY KEY, aplicada_en TEXT NOT NULL)`,
	); err != nil {
		return fmt.Errorf("crear migraciones_aplicadas: %w", err)
	}
	archivos, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return fmt.Errorf("listar migraciones: %w", err)
	}
	for _, archivo := range archivos {
		nombre := archivo.Name()
		partes := nombreDeMigracion.FindStringSubmatch(nombre)
		if partes == nil {
			return fmt.Errorf("nombre de migración inválido %q: se espera NNN_nombre.sql", nombre)
		}
		version, _ := strconv.Atoi(partes[1])
		if err := aplicarMigracion(ctx, db, fsys, nombre, version); err != nil {
			return fmt.Errorf("aplicar migración %s: %w", nombre, err)
		}
	}
	return nil
}

func aplicarMigracion(ctx context.Context, db *sql.DB, fsys fs.FS, nombre string, version int) error {
	var aplicada int
	if err := db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM migraciones_aplicadas WHERE version = ?", version,
	).Scan(&aplicada); err != nil {
		return err
	}
	if aplicada > 0 {
		return nil
	}
	contenido, err := fs.ReadFile(fsys, nombre)
	if err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, string(contenido)); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx,
		"INSERT INTO migraciones_aplicadas (version, aplicada_en) VALUES (?, ?)",
		version, time.Now().UTC().Format(time.RFC3339),
	); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

//go:embed migraciones/*.sql
var migracionesFS embed.FS

// AplicarMigraciones aplica las migraciones embebidas en el binario.
func AplicarMigraciones(ctx context.Context, db *sql.DB) error {
	fsys, err := fs.Sub(migracionesFS, "migraciones")
	if err != nil {
		return fmt.Errorf("abrir migraciones embebidas: %w", err)
	}
	return aplicarMigraciones(ctx, db, fsys)
}
