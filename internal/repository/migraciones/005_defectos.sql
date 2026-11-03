-- Defectos: MODELO-DE-DATOS.md §4.9.

CREATE TABLE defectos (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    proyecto_id          INTEGER NOT NULL,
    historia_id          INTEGER NOT NULL,
    descripcion          TEXT NOT NULL,
    severidad            TEXT NOT NULL CHECK (severidad IN ('Alta', 'Media', 'Baja')),
    estado               TEXT NOT NULL DEFAULT 'Pendiente' CHECK (estado IN ('Pendiente', 'Resuelto')),
    sprint_deteccion_id  INTEGER NOT NULL,
    sprint_resolucion_id INTEGER,
    creado_en            TEXT NOT NULL,
    CHECK ((estado = 'Resuelto') = (sprint_resolucion_id IS NOT NULL)),
    FOREIGN KEY (historia_id, proyecto_id) REFERENCES historias (id, proyecto_id),
    FOREIGN KEY (sprint_deteccion_id, proyecto_id) REFERENCES sprints (id, proyecto_id),
    FOREIGN KEY (sprint_resolucion_id, proyecto_id) REFERENCES sprints (id, proyecto_id)
);
