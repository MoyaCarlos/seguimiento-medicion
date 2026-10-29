-- Registros de horas: MODELO-DE-DATOS.md §4.8.

CREATE TABLE registros_horas (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    proyecto_id INTEGER NOT NULL,
    historia_id INTEGER NOT NULL,
    usuario_id  INTEGER NOT NULL,
    sprint_id   INTEGER NOT NULL,
    fecha       TEXT NOT NULL CHECK (date(fecha) IS fecha),
    actividad   TEXT NOT NULL,
    minutos     INTEGER NOT NULL CHECK (typeof(minutos) = 'integer' AND minutos > 0),
    creado_en   TEXT NOT NULL,
    FOREIGN KEY (historia_id, proyecto_id) REFERENCES historias (id, proyecto_id),
    FOREIGN KEY (sprint_id, proyecto_id) REFERENCES sprints (id, proyecto_id),
    FOREIGN KEY (proyecto_id, usuario_id) REFERENCES integrantes (proyecto_id, usuario_id)
);
