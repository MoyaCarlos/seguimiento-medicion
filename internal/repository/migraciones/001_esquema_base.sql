-- Esquema base: MODELO-DE-DATOS.md §4.1–§4.5.

CREATE TABLE proyectos (
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre                TEXT NOT NULL,
    descripcion           TEXT NOT NULL DEFAULT '',
    fecha_inicio          TEXT,
    fecha_fin             TEXT,
    horas_por_story_point REAL,
    creado_en             TEXT NOT NULL,
    CHECK (horas_por_story_point IS NULL OR (typeof(horas_por_story_point) IN ('real','integer') AND horas_por_story_point > 0)),
    CHECK (date(fecha_inicio) IS fecha_inicio),
    CHECK (date(fecha_fin) IS fecha_fin),
    CHECK (fecha_inicio IS NULL OR fecha_fin IS NULL OR fecha_fin > fecha_inicio)
);

CREATE TABLE usuarios (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre             TEXT NOT NULL,
    nombre_normalizado TEXT NOT NULL UNIQUE,
    creado_en          TEXT NOT NULL
);

CREATE TABLE integrantes (
    proyecto_id INTEGER NOT NULL REFERENCES proyectos (id),
    usuario_id  INTEGER NOT NULL REFERENCES usuarios (id),
    rol         TEXT NOT NULL CHECK (rol IN ('scrum_master','product_builder')),
    creado_en   TEXT NOT NULL,
    PRIMARY KEY (proyecto_id, usuario_id)
);

CREATE TABLE sprints (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    proyecto_id  INTEGER NOT NULL REFERENCES proyectos (id),
    sprint_goal  TEXT,
    fecha_inicio TEXT,
    fecha_fin    TEXT,
    estado       TEXT NOT NULL CHECK (estado IN ('Pendiente','Activo','Finalizado')),
    UNIQUE (id, proyecto_id),
    CHECK ((estado = 'Pendiente' AND sprint_goal IS NULL AND fecha_inicio IS NULL AND fecha_fin IS NULL)
        OR (estado <> 'Pendiente' AND sprint_goal IS NOT NULL AND fecha_inicio IS NOT NULL AND fecha_fin IS NOT NULL AND fecha_fin > fecha_inicio)),
    CHECK (date(fecha_inicio) IS fecha_inicio),
    CHECK (date(fecha_fin) IS fecha_fin)
);

CREATE UNIQUE INDEX sprints_un_activo_por_proyecto ON sprints (proyecto_id) WHERE estado = 'Activo';
CREATE UNIQUE INDEX sprints_un_pendiente_por_proyecto ON sprints (proyecto_id) WHERE estado = 'Pendiente';

CREATE TABLE historias (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    proyecto_id          INTEGER NOT NULL REFERENCES proyectos (id),
    titulo               TEXT NOT NULL,
    descripcion          TEXT NOT NULL,
    criterios_aceptacion TEXT,
    prioridad            TEXT NOT NULL CHECK (prioridad IN ('M','S','C','W')),
    estado               TEXT NOT NULL CHECK (estado IN ('Nueva','En progreso','Completada')),
    valor_negocio        INTEGER,
    estimacion_sp        INTEGER,
    sprint_id            INTEGER,
    fecha_finalizacion   TEXT,
    FOREIGN KEY (sprint_id, proyecto_id) REFERENCES sprints (id, proyecto_id),
    UNIQUE (id, proyecto_id),
    CHECK (estado = 'Nueva' OR sprint_id IS NOT NULL),
    CHECK ((estado = 'Completada') = (fecha_finalizacion IS NOT NULL)),
    CHECK (date(fecha_finalizacion) IS fecha_finalizacion)
);
