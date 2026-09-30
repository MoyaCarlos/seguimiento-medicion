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

// AbrirSQLite abre una conexión a SQLite (driver puro Go, sin CGO).
func AbrirSQLite(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir sqlite: %w", err)
	}
	return db, nil
}

// Migrar crea el esquema de forma idempotente.
func Migrar(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, esquemaBacklogItems); err != nil {
		return fmt.Errorf("migrar backlog_items: %w", err)
	}
	return nil
}
