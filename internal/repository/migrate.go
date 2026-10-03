package repository

import (
	"context"
	"database/sql"
	"fmt"

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

// AbrirSQLite abre una conexión a SQLite (driver puro Go, sin CGO) y activa
// las claves foráneas. Se limita a una conexión para que los PRAGMA y las
// bases :memory: se comporten de forma predecible.
func AbrirSQLite(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.ExecContext(context.Background(), "PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("activar foreign keys: %w", err)
	}
	return db, nil
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
	return nil
}
