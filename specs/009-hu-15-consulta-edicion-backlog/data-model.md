# Data Model: Consulta y Edición del Product Backlog (HU-15)

**Fecha**: 2026-10-25 | **Spec**: [spec.md](./spec.md) | **Research**: [research.md](./research.md)

## Entidad: BacklogItem (Historia de Usuario) — ampliada

Definida en HU-01 (`internal/domain/backlog_item.go`). HU-15 agrega un campo y una operación.

| Campo | Tipo Go | Columna SQLite | Editable en HU-15 | Regla |
|-------|---------|----------------|-------------------|-------|
| ID | `int64` | `id` (PK autoincrement) | No | Orden de creación |
| ProyectoID | `int64` | `proyecto_id` | No | |
| Titulo | `string` | `titulo` | **Sí** | Obligatorio, trim, ≤ 200 runas |
| Descripcion | `string` | `descripcion` | **Sí** | Obligatoria, trim, ≤ 2000 runas |
| **CriteriosAceptacion** (nuevo) | `*string` | `criterios_aceptacion` (ya existe) | **Sí** | Opcional; trim; vacío → `nil`; ≤ 2000 runas |
| Prioridad | `Prioridad` | `prioridad` | No | M/S/C/W |
| Estado | `Estado` | `estado` | No | |
| ValorNegocio | `*int` | `valor_negocio` | No | |
| EstimacionSP | `*int` | `estimacion_sp` | **Sí** | `nil` o ∈ {1, 2, 3, 5, 8, 13, 21} |
| SprintID | `*int64` | `sprint_id` | No | `nil` ⇒ está en el Product Backlog |

### Operación: `ConDatosEditados`

```text
(b BacklogItem) ConDatosEditados(titulo, descripcion string, criterios *string, estimacionSP *int) (BacklogItem, error)
```

Validaciones, en este orden (el primer error se devuelve; nunca se aplica una parte):

| # | Condición | `ValidationError.Campo` | Mensaje |
|---|-----------|-------------------------|---------|
| 1 | título vacío / solo espacios | `titulo` | el título es obligatorio |
| 2 | título > 200 | `titulo` | el título no puede superar los 200 caracteres |
| 3 | descripción vacía / solo espacios | `descripcion` | la descripción es obligatoria |
| 4 | descripción > 2000 | `descripcion` | la descripción no puede superar los 2000 caracteres |
| 5 | criterios > 2000 (tras trim) | `criterios_aceptacion` | los criterios de aceptación no pueden superar los 2000 caracteres |
| 6 | estimación no `nil` y fuera de la escala | `estimacion_sp` | los Story Points deben ser uno de 1, 2, 3, 5, 8, 13, 21 |

Las reglas 1–4 son las mismas de `NewBacklogItem` (se comparten vía helpers). La regla 6 se
expone como `domain.ErrEstimacionInvalida` para que el handler la reutilice ante un número no
entero (R5).

Resultado: copia con los cuatro campos reemplazados y todos los demás intactos (FR-011).

## Vista: Product Backlog (no persistida)

- **Origen**: historias con `proyecto_id = ?` y `sprint_id IS NULL`, cualquiera sea su estado.
- **Orden**: `Prioridad.Orden()` ascendente (M=0, S=1, C=2, W=3) y, a igual prioridad, `ID`
  ascendente. Lo aplica `domain.OrdenarBacklog` (orden estable).
- **Vacía**: lista vacía (`[]`), nunca `null`, cuando el proyecto existe y no hay historias.

## Errores de dominio

| Error | Uso | HTTP |
|-------|-----|------|
| `domain.ErrProyectoNoEncontrado` (existe) | consulta de backlog de proyecto inexistente | 404 |
| `domain.ErrHistoriaNoEncontrada` (nuevo) | edición de historia inexistente | 404 |
| `domain.ValidationError` | edición inválida | 400 con `campo` |

## Puertos (internal/repository)

```text
ProductBacklogRepository
  ListarSinSprint(ctx, proyectoID int64) ([]domain.BacklogItem, error)   // ORDER BY id

HistoriaRepository
  ObtenerPorID(ctx, id int64) (domain.BacklogItem, error)               // ErrHistoriaNoEncontrada
  Actualizar(ctx, item domain.BacklogItem) error                         // titulo, descripcion, criterios, estimacion_sp
```

`SQLiteBacklogRepository` implementa ambos, además de `BacklogRepository` y
`HistoriasDeSprintRepository`. `Guardar` y `ListarPorSprint` pasan a leer/escribir
`criterios_aceptacion`.

## Esquema

Sin cambios: `criterios_aceptacion TEXT` ya está en `001_esquema_base.sql`.
