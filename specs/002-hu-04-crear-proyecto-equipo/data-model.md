# Data Model: HU-04 — Creación de Proyecto y Asignación de Equipo

**Feature**: [spec.md](./spec.md) | [plan.md](./plan.md) | **Fecha**: 2026-10-01

## Entidad de dominio: `Project`

Representa el espacio de trabajo que agrupa backlog, sprints, defectos y métricas.
La struct ya existe (`internal/domain/project.go`); se extiende preservando los campos
actuales y el orden, sin renombrar.

| Campo | Tipo Go | Obligatorio | Descripción |
|-------|---------|-------------|-------------|
| `ID` | `string` | — (asignado al persistir) | UUID v4 generado por el repositorio; se mantiene `string` por el contrato con HU-05. |
| `Name` | `string` | Sí | No vacío (ni solo espacios), ≤ 100 caracteres (FR-005). Nombre repetido permitido (FR-007). |
| `Description` | `string` | No | ≤ 2000 caracteres. Puede quedar vacía (FR-001, FR-019). |
| `StartDate` | `*time.Time` | No | Fecha de inicio; `nil` si no se informa. |
| `EndDate` | `*time.Time` | No | Fecha de fin; `nil` si no se informa. Si ambas existen, no puede ser anterior a `StartDate` (FR-020). |
| `CreatedAt` | `time.Time` | Sí (automático) | Timestamp de creación. No editable. |

> HU-13 (`Should have`) es quien agrega la consulta del estado derivado
> (Planificado/En curso/Finalizado). HU-04 solo registra y edita las fechas.

### Constructor y edición (validaciones reutilizadas)

```go
func NewProject(name, description string, start, end *time.Time) (Project, error)
func (p Project) ConDatosEditados(name, description string, start, end *time.Time) (Project, error)
```

Invariantes validados (en este orden; se devuelve el primer `ValidationError`).
`ConDatosEditados` reutiliza exactamente la misma validación que `NewProject` (DRY, FR-019)
y conserva `ID` y `CreatedAt`:

1. `name` no vacío tras `strings.TrimSpace` → si falla: `ValidationError{Campo: "nombre", Mensaje: "el nombre es obligatorio"}`.
2. `len([]rune(name)) <= 100` → si falla: `ValidationError{Campo: "nombre", Mensaje: "el nombre no puede superar los 100 caracteres"}`.
3. `len([]rune(description)) <= 2000` → si falla: `ValidationError{Campo: "descripcion", Mensaje: "la descripción no puede superar los 2000 caracteres"}`.
4. `start == nil || end == nil || !end.Before(*start)` → si falla: `ValidationError{Campo: "fecha_fin", Mensaje: "la fecha de fin no puede ser anterior a la fecha de inicio"}`.

Comportamiento de inicialización: `Name` y `Description` se normalizan con `TrimSpace`;
`ID` queda vacío hasta que el repositorio lo asigne; `CreatedAt` se fija al persistir.

## Entidad de dominio: `User` (integrante)

| Campo | Tipo Go | Obligatorio | Descripción |
|-------|---------|-------------|-------------|
| `ID` | `string` | — (asignado al persistir) | UUID v4. |
| `Name` | `string` | Sí | Nombre tal como se ingresó, recortado. No vacío, ≤ 200 caracteres (FR-013). |
| `NormalizedName` | `string` | Sí (derivado) | `ToLower(TrimSpace(Name))`; clave de unicidad para reutilización (FR-016, R3). |
| `CreatedAt` | `time.Time` | Sí (automático) | Timestamp de creación. |

### Constructor

```go
func NewUser(name string) (User, error)
```

1. `name` no vacío tras `TrimSpace` → si falla: `ValidationError{Campo: "nombre", Mensaje: "el nombre del integrante es obligatorio"}`.
2. `len([]rune(name)) <= 200` → si falla: `ValidationError{Campo: "nombre", Mensaje: "el nombre del integrante no puede superar los 200 caracteres"}`.

`NormalizedName` se calcula siempre con `ToLower(TrimSpace(name))` (FR-016).

## Enum de dominio: `Role`

| Constante | Valor canónico | Etiqueta visible |
|-----------|----------------|------------------|
| `RolScrumMaster` | `"scrum_master"` | Scrum Master |
| `RolProductBuilder` | `"product_builder"` | Product Builder |

`func (r Role) Valido() bool` acepta únicamente los dos valores anteriores (FR-010).
Un rol inválido produce `ValidationError{Campo: "rol", Mensaje: "el rol debe ser scrum_master o product_builder"}`.

## Entidad de dominio: `Membership` (relación Proyecto–Usuario–Rol)

| Campo | Tipo Go | Obligatorio | Descripción |
|-------|---------|-------------|-------------|
| `ProjectID` | `string` | Sí | Proyecto existente (FR-009). |
| `UserID` | `string` | Sí | Integrante (existente o recién creado). |
| `Role` | `Role` | Sí | Un único rol por integrante y proyecto (FR-011, SC-007). |
| `CreatedAt` | `time.Time` | Sí (automático) | Timestamp de vinculación. |

### Constructor

```go
func NewMembership(projectID, userID string, role Role) (Membership, error)
```

1. `projectID` no vacío → si falla: `ValidationError{Campo: "proyecto_id", Mensaje: "el proyecto es obligatorio"}`.
2. `userID` no vacío → si falla: `ValidationError{Campo: "integrante_id", Mensaje: "el integrante es obligatorio"}`.
3. `role.Valido()` → si falla: el `ValidationError` de `Role`.

## Tipo de error

Se reutiliza el `domain.ValidationError{Campo, Mensaje}` existente
(`internal/domain/errors.go`); el adaptador HTTP lo mapea a `400`.

## Tablas SQLite

```sql
CREATE TABLE IF NOT EXISTS projects (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    start_date  TEXT,              -- RFC3339 o NULL
    end_date    TEXT,              -- RFC3339 o NULL
    created_at  TEXT NOT NULL      -- RFC3339
);

CREATE TABLE IF NOT EXISTS users (
    id              TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    normalized_name TEXT NOT NULL UNIQUE,
    created_at      TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS project_members (
    project_id TEXT NOT NULL,
    user_id    TEXT NOT NULL,
    role       TEXT NOT NULL,
    created_at TEXT NOT NULL,
    PRIMARY KEY (project_id, user_id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (user_id)    REFERENCES users(id)
);
```

- `projects.start_date` / `end_date` son `NULL` cuando el dominio tiene `nil`.
- `users.normalized_name` es `UNIQUE`: materializa la reutilización de integrantes (R3).
- `project_members.PRIMARY KEY (project_id, user_id)` rechaza la vinculación duplicada (FR-015).
- `PRAGMA foreign_keys = ON` se activa al abrir la conexión para que las FK se apliquen.

## Transiciones de estado

El proyecto **no** tiene estados de ciclo de vida (clarificado en `/speckit.clarify`);
se crea operativo. La relación integrante–proyecto–rol tampoco cambia de estado en esta
historia (la baja o el cambio de rol quedan fuera de alcance).

## Mapeo Casos de uso → Entidades

| Caso de uso | Entradas | Efecto |
|-------------|----------|--------|
| `CrearProyecto` | `nombre`, `descripcion`, `fecha_inicio?`, `fecha_fin?`, `creador` (nombre) | Crea y persiste el `Project`; resuelve/crea el `User` del creador y persiste `Membership` con rol `scrum_master` (FR-008). |
| `ObtenerProyecto` | `id` | Carga el `Project` por `GetByID`; soporta la lectura para verificar persistencia (FR-004, SC-002/SC-009). |
| `EditarProyecto` | `id`, `nombre`, `descripcion`, `fecha_inicio?`, `fecha_fin?` | Carga el `Project`, valida con las reglas del alta y persiste el `UPDATE` (FR-018..FR-023). |
| `AsignarIntegrante` | `proyecto_id`, `nombre`, `rol` | Valida el proyecto, resuelve/crea el `User` por nombre normalizado y persiste la `Membership`; rechaza duplicados (FR-009..FR-016). |
