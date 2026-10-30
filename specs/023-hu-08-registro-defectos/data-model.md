# Data Model: Registro y Seguimiento de Defectos (HU-08)

**Fecha**: 2026-10-30 | **Spec**: [spec.md](./spec.md) | **Research**: [research.md](./research.md)

## Esquema (persistencia)

La tabla ya está diseñada; este documento **no la redefine**:

- [`MODELO-DE-DATOS.md` §4.9 `defectos`](../../MODELO-DE-DATOS.md#49-defectos--hu-08-tabla-nueva-migración-005):
  columnas, FK compuestas y `CHECK`.
- §2 (convenciones, instantes en RFC 3339 UTC), §4.10 (`CHECK` frente a `NULL`), §5 (dónde vive
  cada regla y traducción de errores), D6 (sin título), D9 (sin borrado), §9 (migraciones; la 005
  está reservada para HU-08).

Esta historia agrega `internal/repository/migraciones/005_defectos.sql` con ese DDL y no cambia el
documento.

## Entidades del dominio

### `Defecto` (nueva, `internal/domain/defecto.go`)

| Campo | Tipo Go | Origen | Regla |
|---|---|---|---|
| `ID` | `int64` | base | asignado al guardar; define el orden de registro |
| `ProyectoID` | `int64` | ruta | proyecto existente (handler) |
| `HistoriaID` | `int64` | cuerpo (`historia_id`) | obligatoria; historia del mismo proyecto, en cualquier Sprint o estado (service) |
| `Descripcion` | `string` | cuerpo | sin espacios externos; no vacía; ≤ 2000 runas |
| `Severidad` | `Severidad` | cuerpo | exactamente `Alta`, `Media` o `Baja` |
| `Estado` | `EstadoDefecto` | sistema | `Pendiente` al crear; `Resuelto` solo vía `Resolver` |
| `SprintDeteccionID` | `int64` | sistema | Sprint `Activo` del proyecto al registrar; no se recalcula |
| `SprintResolucionID` | `*int64` | sistema | `nil` mientras está `Pendiente`; Sprint `Activo` al resolver |
| `CreadoEn` | `time.Time` | base | fecha y hora de registro (RFC 3339 UTC) |

**Tipos**:

```text
Severidad      = "Alta" | "Media" | "Baja"
EstadoDefecto  = "Pendiente" | "Resuelto"
```

**Factory**:

```text
NewDefecto(proyectoID, historiaID int64, activo Sprint,
           descripcion, severidad string) (Defecto, error)
```

Valida, en este orden, descripción → severidad. Ante el primer error devuelve un
`ValidationError` y ningún defecto. Fija `Estado = Pendiente`,
`SprintDeteccionID = activo.ID` y `SprintResolucionID = nil`. Supone que el service ya verificó
que hay un Sprint activo y que la historia es del proyecto (research R2).

**Transición**:

```text
(*Defecto).Resolver(activo *Sprint) error
```

| Estado actual | `activo` | Resultado |
|---|---|---|
| `Resuelto` | cualquiera | `ErrDefectoYaResuelto`; no cambia nada |
| `Pendiente` | `nil` | `ErrSinSprintActivoParaResolver`; no cambia nada |
| `Pendiente` | Sprint `Activo` | `Estado = Resuelto`, `SprintResolucionID = &activo.ID` |

No hay transición `Resuelto → Pendiente` (Clarification 1, FR-016).

**Función del dominio** (`internal/domain/sprint.go`):

- `SprintActivoDe(sprints []Sprint) (Sprint, bool)`: el Sprint `Activo` de la lista, si hay
  (research R5).

**Errores nuevos** (`internal/domain/errors.go`):

| Error | Tipo / Campo | Mensaje |
|---|---|---|
| `ErrDefectoNoEncontrado` | `errors.New` (404) | "defecto no encontrado" |
| `ErrSinSprintActivoParaDetectar` | `sprint_deteccion_id` | "el proyecto no tiene un Sprint activo al que asociar el defecto como Sprint de detección" |
| `ErrSinSprintActivoParaResolver` | `sprint_resolucion_id` | "el proyecto no tiene un Sprint activo al que asociar el defecto como Sprint de resolución" |
| `ErrDefectoYaResuelto` | `estado` | "el defecto ya está resuelto" |
| `ErrHistoriaRelacionadaObligatoria` | `historia_id` | "la historia relacionada es obligatoria" |
| `ErrHistoriaRelacionadaAjena` | `historia_id` | "la historia relacionada no existe en el proyecto" |
| (descripción) | `descripcion` | "la descripción es obligatoria" / "la descripción no puede superar los 2000 caracteres" |
| (severidad) | `severidad` | "la severidad debe ser Alta, Media o Baja" |

Los textos son orientativos; los fija el test RED de cada caso.

### `DefectoDetalle` (vista de lectura)

`Defecto` + `TituloHistoria string`. Al listar y al obtener por id se arma con `JOIN historias`; al registrar, con el
`Titulo` de la historia que leyó el service. Igual que `RegistroHorasDetalle`, es una vista para
mostrar; no se construye por factory.

## Puertos (`internal/repository/defecto_repository.go`)

```text
DefectoRepository
  Registrar(ctx, d domain.Defecto) (domain.Defecto, error)
      INSERT condicional (R4): si el Sprint de detección ya no está Activo
      → domain.ErrSinSprintActivoParaDetectar. FK → ErrIntegridadReferencial.
  ObtenerPorID(ctx, id int64) (domain.DefectoDetalle, error)
      Con el título de la historia (JOIN historias). domain.ErrDefectoNoEncontrado si no existe.
  Resolver(ctx, d domain.Defecto) error
      UPDATE condicional (R4) con d.SprintResolucionID: si el defecto ya no está Pendiente
      → domain.ErrDefectoYaResuelto; si el Sprint ya no está Activo
      → domain.ErrSinSprintActivoParaResolver.
  ListarPorProyecto(ctx, proyectoID int64) ([]domain.DefectoDetalle, error)
      ORDER BY id. Lista vacía (no nil) si no hay defectos.
```

Lecturas que se reutilizan sin cambios: `ProjectRepository.ObtenerPorID` (vía el service
`ObtenerProyecto`), `HistoriaRepository.ObtenerPorID` (`SQLiteBacklogRepository`) y
`SprintRepository.ListarPorProyecto`. Cada service declara solo los métodos que usa (ISP).

## Casos de uso (`internal/service`)

| Service | Depende de | Hace |
|---|---|---|
| `RegistrarDefecto` | `SprintRepository` (`ListarPorProyecto`), `HistoriaRepository` (`ObtenerPorID`), `DefectoRepository` | Sprint activo → historia relacionada → `NewDefecto` → `Registrar`; devuelve `DefectoDetalle` con el título de la historia (respuesta 201) |
| `ResolverDefecto` | `DefectoRepository`, `SprintRepository` (`ListarPorProyecto`) | `ObtenerPorID` → Sprint activo del proyecto del defecto → `Resolver` del dominio → `Resolver` del repositorio; devuelve el `DefectoDetalle` resuelto (respuesta 200) |
| `ListarDefectos` | `ProjectRepository` (`ObtenerPorID`), `DefectoRepository` | proyecto existe → `ListarPorProyecto` |

## Estados

```text
           Registrar (con Sprint Activo)          Resolver (con Sprint Activo)
  (nada) ─────────────────────────────▶ Pendiente ────────────────────────────▶ Resuelto
                                       sprint_resolucion_id = NULL              sprint_resolucion_id = Sprint Activo
```

`Resuelto` es final. El estado de la historia relacionada y el Sprint en que esté no influyen en
ninguna transición (Clarification 2).
