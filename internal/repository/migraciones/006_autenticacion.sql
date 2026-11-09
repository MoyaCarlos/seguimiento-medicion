-- Autenticación: MODELO-DE-DATOS.md §4.2 y §4.11.

ALTER TABLE usuarios ADD COLUMN contrasena_hash TEXT;

CREATE TABLE sesiones (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    usuario_id INTEGER NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    creado_en  TEXT NOT NULL,
    vence_en   TEXT NOT NULL,
    cerrada_en TEXT,
    FOREIGN KEY (usuario_id) REFERENCES usuarios (id)
);
