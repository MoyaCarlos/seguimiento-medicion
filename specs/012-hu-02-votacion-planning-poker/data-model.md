# Data Model: Votación de Esfuerzo en Planning Poker (HU-02)

**Fecha**: 2026-10-27 | **Spec**: [spec.md](./spec.md) | **Research**: [research.md](./research.md)

## Esquema

Las tablas `rondas_votacion` y `votos` están definidas en
[MODELO-DE-DATOS.md §4.7](../../MODELO-DE-DATOS.md#47-rondas_votacion-y-votos--hu-02-hu-14-tablas-nuevas-migración-003)
y se crean con la migración **003** (número reservado para HU-02, §9). Esta historia no
redefine ni cambia el esquema: la migración traduce §4.7 a DDL, con las convenciones de la
001 (sin `?` en §4.7 ⇒ `NOT NULL`, según la notación de §4).

Restricciones que la base garantiza y la historia aprovecha:

| Restricción | Origen | Uso en HU-02 |
|---|---|---|
| `UNIQUE (historia_id, numero)` | §4.7 | frena dos "iniciar" simultáneos → `ErrVotacionYaIniciada` (R6) |
| `PRIMARY KEY (ronda_id, usuario_id)` | §4.7 | un voto por integrante y ronda; el cambio de voto es un *upsert* |
| FK `historia_id → historias`, `ronda_id → rondas_votacion`, `usuario_id → usuarios` | §4.7 | `ErrIntegridadReferencial` (409) si algo se cuela |

Reglas que **no** están en la base (§5): escala Fibonacci del voto (dominio), que el votante
sea integrante del proyecto (servicio), la ocultación de valores (servicio/handler).

## Entidades de dominio (`internal/domain/votacion.go`)

### `RondaVotacion`

| Campo | Tipo | Notas |
|---|---|---|
| `ID` | `int64` | |
| `HistoriaID` | `int64` | |
| `Numero` | `int` | 1 en HU-02; HU-14 abre las siguientes |
| `ReveladaEn` | `*time.Time` | `nil` = votos ocultos |
| `Votos` | `[]Voto` | a lo sumo uno por `UsuarioID` |

- `NuevaRonda(historiaID int64) (RondaVotacion, error)` — ronda 1, sin votos ni revelar;
  `ValidationError{historia_id}` si `historiaID <= 0`.
- `(r *RondaVotacion) Votar(usuarioID int64, valor int) error`
  - revelada → `ErrRondaRevelada`;
  - `!EsEstimacionValida(valor)` → `ErrVotoInvalido`;
  - si el usuario ya votó, reemplaza el valor; si no, agrega el voto.
- `(r *RondaVotacion) Revelar(integrantes []int64, ahora time.Time) error`
  - revelada → `ErrRondaRevelada`;
  - algún integrante sin voto → `ErrFaltanVotos`;
  - fija `ReveladaEn = ahora`.
- `(r RondaVotacion) Revelada() bool`
- `(r RondaVotacion) Promedio() (Promedio, bool)` — `false` si no está revelada.

### `Voto`

| Campo | Tipo |
|---|---|
| `UsuarioID` | `int64` |
| `Valor` | `int` |

### `Promedio` y `CalcularPromedio`

| Campo | Tipo | Regla |
|---|---|---|
| `Exacto` | `float64` | media aritmética redondeada a 2 decimales |
| `CartaSugerida` | `int` | carta de la escala más cercana a `Exacto` (el valor que ve el equipo, así nunca se contradicen); empate → la mayor |

Ejemplos (casos de test): `[3,5,8]` → 5,33 / 5; `[3,5]` → 4 / 5 (empate); `[8,8,8]` →
8 / 8; `[1,21]` → 11 / 13 (8 dista 3, 13 dista 2); `[2,3]` → 2,5 / 3 (empate).

### Errores nuevos (`internal/domain/errors.go`)

Ver la tabla de [research.md R8](./research.md#r8-errores-y-códigos-http):
`ErrVotacionNoIniciada`, `ErrVotacionYaIniciada`, `ErrNoEsIntegrante`, `ErrNoEsScrumMaster`
(sentinels) y `ErrVotoInvalido`, `ErrHistoriaCompletada`, `ErrRondaRevelada`,
`ErrFaltanVotos` (`ValidationError`).

### Transiciones de una ronda

```text
            NuevaRonda                 Revelar (votaron todos)
 (sin ronda) ─────────▶ Abierta ───────────────────────────▶ Revelada
                         │  ▲                                  │
                         └──┘ Votar (agrega o reemplaza)       └─ Votar / Revelar → ErrRondaRevelada
```

En HU-02 una historia tiene como mucho una ronda. HU-14 agrega `Revelada → nueva Abierta
(numero + 1)` y conserva la anterior como historial.

## Vista de consulta (servicio)

`ConsultarVotacion` arma, a partir de la ronda actual y los integrantes del proyecto:

| Campo | Tipo | Regla |
|---|---|---|
| `HistoriaID`, `Numero` | | de la ronda |
| `ReveladaEn` | `*time.Time` | |
| `Integrantes[]` | `UsuarioID`, `Nombre`, `Rol`, `YaVoto bool`, `Valor *int` | `Valor` solo si la ronda está revelada **o** el integrante es quien consulta (FR-007, clarify Q4); si no, `nil` ("Oculto" si `YaVoto`) |
| `Promedio` | `*Promedio` | solo revelada |

Los integrantes salen de `ProjectRepository.ListarIntegrantes`, así "quién falta" incluye a
quien se sumó con la ronda abierta (Assumptions de la spec).

## Puertos (`internal/repository/votacion_repository.go`)

```go
// VotacionRepository es el puerto de persistencia de Planning Poker.
// ObtenerRondaActual devuelve domain.ErrVotacionNoIniciada si la historia no tiene rondas.
type VotacionRepository interface {
    CrearRonda(ctx context.Context, r domain.RondaVotacion) (domain.RondaVotacion, error)
    ObtenerRondaActual(ctx context.Context, historiaID int64) (domain.RondaVotacion, error)
    GuardarVoto(ctx context.Context, rondaID int64, v domain.Voto) error
    MarcarRevelada(ctx context.Context, rondaID int64, en time.Time) error
}
```

| Método | SQLite | Errores |
|---|---|---|
| `CrearRonda` | `INSERT INTO rondas_votacion` | 2067 → `ErrVotacionYaIniciada`; 787 → `ErrIntegridadReferencial` |
| `ObtenerRondaActual` | ronda de mayor `numero` + sus votos ordenados por `usuario_id` | sin filas → `ErrVotacionNoIniciada` |
| `GuardarVoto` | upsert condicional (R6) | 0 filas → `ErrRondaRevelada`; 787 → `ErrIntegridadReferencial` |
| `MarcarRevelada` | `UPDATE … WHERE revelada_en IS NULL` (R6) | 0 filas → `ErrRondaRevelada` |

Se reutilizan sin cambios: `HistoriaRepository.ObtenerPorID` (HU-15) y
`ProjectRepository.ListarIntegrantes` (HU-04).
