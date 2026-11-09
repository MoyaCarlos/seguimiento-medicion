# Data Model: Avance de las Historias del Proyecto (HU-18)

**Sin cambios de esquema.** Las tablas están en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md): `historias` (§4.4) y `sprints` (§4.5). HU-18 es de solo lectura (§8). Nada de lo siguiente se guarda.

## Invariante que se aprovecha

`CHECK (estado = 'Nueva' OR sprint_id IS NOT NULL)` (§4.4): En progreso y Completada siempre tienen Sprint; solo Nueva puede estar en el Product Backlog.

## Vistas de lectura (`internal/domain`)

### HistoriaAvance

| Campo | Tipo | Origen |
|---|---|---|
| `ID` | `int64` | `historias.id` |
| `Titulo` | `string` | `historias.titulo` |
| `EstimacionSP` | `*int` | `historias.estimacion_sp`; `nil` = "sin estimar" |
| `Estado` | `Estado` | `historias.estado` |
| `Sprint` | `*SprintRef` | `nil` si `sprint_id` es NULL (Product Backlog) |

`SprintRef` = `{ID int64, Numero int}`; el `Numero` es el derivado de `Sprint.Numero` (#129).

### GrupoAvance

| Campo | Regla |
|---|---|
| `Estado` | Nueva, En progreso o Completada |
| `Cantidad` | cantidad de historias del grupo (incluye las sin estimar) |
| `StoryPoints` | suma de `EstimacionSP` del grupo; las `nil` suman 0 |
| `Historias` | en el orden recibido (creación); nunca `nil` |

### AvanceHistorias

`Grupos` con exactamente tres elementos, siempre en el orden Nueva, En progreso, Completada, más `Total` (`Cantidad` y `StoryPoints` de todo el proyecto). Se construye con `AgruparAvance(historias []BacklogItem, sprints []Sprint) (AvanceHistorias, error)`; devuelve error si una historia referencia un Sprint ausente de `sprints` o tiene un estado desconocido. Esa es la firma desde la primera tarea de implementación (`tasks.md`, US1): con `sprints` vacío o `nil` ninguna historia con Sprint se puede resolver, y la resolución del Sprint se agrega en US2.

## Puerto (`internal/repository`)

```text
AvanceHistoriasRepository
  ListarPorProyecto(ctx, proyectoID) ([]domain.BacklogItem, error)   // todas las historias, ORDER BY id; vacío ≠ nil
```

## Service (`internal/service`)

`ConsultarAvanceHistorias` depende de `buscadorProyectos` (existencia), `AvanceHistoriasRepository` y `listadorSprints`. `Ejecutar(ctx, proyectoID) (domain.AvanceHistorias, error)`.

## Reglas de negocio → dónde viven

| Regla (spec) | Dónde |
|---|---|
| Tres grupos siempre (FR-002) | `AgruparAvance` |
| Sin estimar cuenta y suma 0 (FR-004) | `AgruparAvance` |
| Sprint o Product Backlog (FR-003) | `AgruparAvance` |
| 404 (FR-006) | service (`ErrProyectoNoEncontrado`), mapeado por `escribirError` |
| 400 con campo `id` (FR-006) | `parsearID` en el handler, antes de llamar al service |
| Sin escritura (FR-007) | el puerto solo expone lectura |
