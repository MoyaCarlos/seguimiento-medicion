# Data Model: HU-01 — Creación de Historias de Usuario

**Feature**: [spec.md](./spec.md) | [plan.md](./plan.md) | **Fecha**: 2026-09-30

## Entidad de dominio: `BacklogItem`

Representa una historia de usuario del Product Backlog.

| Campo | Tipo Go | Obligatorio | Descripción |
|-------|---------|-------------|-------------|
| `ID` | `int64` | — (asignado al persistir) | Identificador autogenerado por SQLite (`AUTOINCREMENT`); preserva el orden de creación. |
| `ProyectoID` | `int64` | Sí | Proyecto al que pertenece (FR-012). Entero positivo (`> 0`); la existencia del proyecto no se valida aún (HU-04). |
| `Titulo` | `string` | Sí | No vacío (ni solo espacios), ≤ 200 caracteres (FR-004). |
| `Descripcion` | `string` | Sí | No vacía (ni solo espacios), ≤ 2000 caracteres (FR-005). |
| `Prioridad` | `Prioridad` (enum string) | Sí | Escala MoSCoW: `"M"`, `"S"`, `"C"`, `"W"` (FR-010). |
| `Estado` | `Estado` (enum string) | Sí (automático) | Siempre `"Nueva"` al crear (FR-002). |
| `ValorNegocio` | `*int` | No | Escala Fibonacci `{1,2,3,5,8,13,21}` (FR-008/FR-009). `nil` si no se informa. |
| `EstimacionSP` | `*int` | No | Story Points. Siempre `nil` al crear; se completa en Planning Poker (HU-02) (FR-013). |

**Orden y listado**: el orden de creación se preserva mediante el `ID` incremental
(FR-003). La vista de listado del Product Backlog y su actualización en vivo quedan
fuera del alcance de esta historia (ver `plan.md`, Complexity Tracking).

### Enumeraciones

- `Prioridad`: `PrioridadMust = "M"`, `PrioridadShould = "S"`, `PrioridadCould = "C"`, `PrioridadWont = "W"`. Se mapean a Alta/Media/Baja en el tablero (convención del Product Backlog).
- `Estado`: `EstadoNueva = "Nueva"`.
- Escala de Valor de Negocio (Fibonacci): `1, 2, 3, 5, 8, 13, 21`.

## Constructor Factory: `NewBacklogItem`

```go
func NewBacklogItem(
    proyectoID int64,
    titulo string,
    descripcion string,
    prioridad Prioridad,
    valorNegocio *int,
) (BacklogItem, error)
```

Invariantes validados (en este orden; se devuelve el primer `ValidationError`):

1. `proyectoID > 0` → si falla: `ValidationError{Campo: "proyecto_id", Mensaje: "el proyecto es obligatorio"}`.
2. `titulo` no vacío tras `strings.TrimSpace` → si falla: `ValidationError{Campo: "titulo", Mensaje: "el título es obligatorio"}`.
3. `len([]rune(titulo)) <= 200` → si falla: `ValidationError{Campo: "titulo", Mensaje: "el título no puede superar los 200 caracteres"}`.
4. `descripcion` no vacía tras `strings.TrimSpace` → si falla: `ValidationError{Campo: "descripcion", Mensaje: "la descripción es obligatoria"}`.
5. `len([]rune(descripcion)) <= 2000` → si falla: `ValidationError{Campo: "descripcion", Mensaje: "la descripción no puede superar los 2000 caracteres"}`.
6. `prioridad` pertenece al enum MoSCoW → si falla: `ValidationError{Campo: "prioridad", Mensaje: "la prioridad debe ser M, S, C o W"}`.
7. `valorNegocio` es `nil` o pertenece a `{1,2,3,5,8,13,21}` → si falla: `ValidationError{Campo: "valor_negocio", Mensaje: "el valor de negocio debe ser uno de 1, 2, 3, 5, 8, 13, 21"}`.

Comportamiento de inicialización:

- `Estado` se fuerza a `EstadoNueva` (nunca se recibe como entrada).
- `EstimacionSP` se inicializa en `nil` (nunca se recibe como entrada).
- `ID` queda en `0` hasta que el repositorio lo asigne.
- Los campos `string` se normalizan con `TrimSpace` antes de asignarse.

### Tipo de error

```go
type ValidationError struct {
    Campo   string
    Mensaje string
}

func (e ValidationError) Error() string
```

`errors.As(err, &domain.ValidationError{})` permite al adaptador HTTP mapear a `400`.

## Tabla SQLite: `backlog_items`

```sql
CREATE TABLE IF NOT EXISTS backlog_items (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    proyecto_id   INTEGER NOT NULL,
    titulo        TEXT    NOT NULL,
    descripcion   TEXT    NOT NULL,
    prioridad     TEXT    NOT NULL,
    estado        TEXT    NOT NULL,
    valor_negocio INTEGER,
    estimacion_sp INTEGER
);
```

- `valor_negocio` y `estimacion_sp` son `NULL` cuando el dominio tiene `nil`.
- `proyecto_id` es obligatorio y debe ser `> 0` (validado en el dominio, invariante 1); su existencia no se valida (HU-04).
- No se declara `FOREIGN KEY (proyecto_id) REFERENCES projects(id)` porque la tabla `projects` se crea en HU-04.

## Transiciones de estado

| Estado inicial | Evento | Estado resultante |
|----------------|--------|-------------------|
| — | Crear historia (esta HU) | `"Nueva"` |
| `"Nueva"` | Planificación / Sprint (HU-06, HU-11) | fuera de alcance |

En esta historia el estado es siempre `"Nueva"`; no hay transiciones administradas aquí.

## Mapeo Input → Entidad

`service.CrearHistoriaInput` (entrada del caso de uso):

| Campo | Tipo | Origen |
|-------|------|--------|
| `ProyectoID` | `int64` | body JSON |
| `Titulo` | `string` | body JSON |
| `Descripcion` | `string` | body JSON |
| `Prioridad` | `domain.Prioridad` | body JSON |
| `ValorNegocio` | `*int` | body JSON (opcional) |

El servicio construye el `BacklogItem` con `domain.NewBacklogItem` y delega la
persistencia en `BacklogRepository.Guardar`.
