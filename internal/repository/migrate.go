package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

const esquemaBacklogItems = `
CREATE TABLE IF NOT EXISTS backlog_items (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    proyecto_id   INTEGER NOT NULL,
    titulo        TEXT    NOT NULL,
    descripcion   TEXT    NOT NULL,
    prioridad     TEXT    NOT NULL,
    estado        TEXT    NOT NULL,
    valor_negocio INTEGER,
    estimacion_sp INTEGER
);`

const esquemaProjects = `
CREATE TABLE IF NOT EXISTS projects (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    start_date  TEXT,
    end_date    TEXT,
    created_at  TEXT NOT NULL
);`

const esquemaUsers = `
CREATE TABLE IF NOT EXISTS users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    normalized_name TEXT NOT NULL UNIQUE,
    created_at      TEXT NOT NULL
);`

const esquemaProjectMembers = `
CREATE TABLE IF NOT EXISTS project_members (
    project_id INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    role       TEXT NOT NULL,
    created_at TEXT NOT NULL,
    PRIMARY KEY (project_id, user_id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (user_id)    REFERENCES users(id)
);`

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

// Migrar crea el esquema de forma idempotente.
func Migrar(ctx context.Context, db *sql.DB) error {
	esquemas := []struct {
		nombre string
		ddl    string
	}{
		{"backlog_items", esquemaBacklogItems},
		{"projects", esquemaProjects},
		{"users", esquemaUsers},
		{"project_members", esquemaProjectMembers},
	}
	for _, e := range esquemas {
		if _, err := db.ExecContext(ctx, e.ddl); err != nil {
			return fmt.Errorf("migrar %s: %w", e.nombre, err)
		}
	}
	if err := verificarEsquemaViejo(ctx, db); err != nil {
		return err
	}
	return nil
}

// verificarEsquemaViejo detecta una base creada con IDs TEXT (esquema previo a
// int64) y devuelve un error que indica cómo resolverlo. CREATE TABLE IF NOT
// EXISTS no migra una tabla ya existente, así que ese caso debe señalarse.
func verificarEsquemaViejo(ctx context.Context, db *sql.DB) error {
	filas, err := db.QueryContext(ctx, "PRAGMA table_info(projects)")
	if err != nil {
		return fmt.Errorf("inspeccionar esquema de projects: %w", err)
	}
	defer filas.Close()

	for filas.Next() {
		var (
			cid     int
			nombre  string
			tipo    string
			notnull int
			dflt    sql.NullString
			pk      int
		)
		if err := filas.Scan(&cid, &nombre, &tipo, &notnull, &dflt, &pk); err != nil {
			return fmt.Errorf("leer esquema de projects: %w", err)
		}
		if nombre == "id" && strings.EqualFold(tipo, "TEXT") {
			return ErrEsquemaViejo
		}
	}
	return filas.Err()
}
