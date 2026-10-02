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

const esquemaSprints = `
CREATE TABLE IF NOT EXISTS sprints (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    proyecto_id  INTEGER NOT NULL,
    sprint_goal  TEXT,
    fecha_inicio TEXT,
    fecha_fin    TEXT,
    estado       TEXT    NOT NULL
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
	if _, err := db.ExecContext(ctx, esquemaSprints); err != nil {
		return fmt.Errorf("migrar sprints: %w", err)
	}
	return asegurarColumna(ctx, db, "backlog_items", "sprint_id", "INTEGER")
}

// asegurarColumna agrega la columna si todavía no existe: SQLite no soporta
// "ADD COLUMN IF NOT EXISTS" y las bases locales creadas por HU-01 no la tienen.
func asegurarColumna(ctx context.Context, db *sql.DB, tabla, columna, tipo string) error {
	var existe int
	err := db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = ?", tabla, columna,
	).Scan(&existe)
	if err != nil {
		return fmt.Errorf("inspeccionar %s: %w", tabla, err)
	}
	if existe > 0 {
		return nil
	}
	if _, err := db.ExecContext(ctx, fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tabla, columna, tipo)); err != nil {
		return fmt.Errorf("agregar %s.%s: %w", tabla, columna, err)
	}
	return nil
}
