-- Votación de Planning Poker: MODELO-DE-DATOS.md §4.7.

CREATE TABLE rondas_votacion (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    historia_id INTEGER NOT NULL REFERENCES historias (id),
    numero      INTEGER NOT NULL,
    revelada_en TEXT,
    UNIQUE (historia_id, numero)
);

CREATE TABLE votos (
    ronda_id   INTEGER NOT NULL REFERENCES rondas_votacion (id),
    usuario_id INTEGER NOT NULL REFERENCES usuarios (id),
    valor      INTEGER NOT NULL,
    PRIMARY KEY (ronda_id, usuario_id)
);
