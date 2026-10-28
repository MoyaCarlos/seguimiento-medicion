-- Historias arrastradas al cerrar un Sprint: MODELO-DE-DATOS.md §4.6.

CREATE TABLE historias_arrastradas (
    sprint_id     INTEGER NOT NULL,
    historia_id   INTEGER NOT NULL,
    proyecto_id   INTEGER NOT NULL,
    estimacion_sp INTEGER,
    PRIMARY KEY (sprint_id, historia_id),
    FOREIGN KEY (sprint_id, proyecto_id) REFERENCES sprints (id, proyecto_id),
    FOREIGN KEY (historia_id, proyecto_id) REFERENCES historias (id, proyecto_id)
);
