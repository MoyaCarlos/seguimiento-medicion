# Data Model: Generación de Reporte en PDF (HU-10)

El reporte es **derivado y no persistido**: no agrega tablas ni columnas. Toma los datos de
las entidades y tablas ya definidas en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md):

- Proyecto e integrantes: §4.1 y §4.2 (HU-04).
- Sprints: §4.3 (HU-05) y `historias_arrastradas`: §4.6 (HU-05/HU-06).
- Historias: §4.4 (HU-01/HU-15).
- Registros de horas: §4.8 (HU-07); llegan agregados por Sprint en `DatosMetricasSprint`.
- Defectos: §4.9 (HU-08).
- Métricas: no se persisten (§6); se calculan con `domain.CalcularMetricas` (HU-03).

## Tipos del dominio (nuevos, en `internal/domain/reporte.go`)

### `HistoriaReporte`

| Campo | Tipo | Regla |
|-------|------|-------|
| `ID` | `int64` | Identificador de la historia |
| `Titulo` | `string` | Título actual de la historia |
| `Prioridad` | `Prioridad` | M, S, C o W |
| `EstimacionSP` | `*int` | `nil` = "sin estimar"; suma 0 a los totales |
| `Estado` | `Estado` | Estado actual de la historia |
| `Completada` | `bool` | `true` si la historia estaba completada en ese Sprint |

### `DefectoReporte`

| Campo | Tipo | Regla |
|-------|------|-------|
| `Descripcion` | `string` | Descripción del defecto |
| `Severidad` | `Severidad` | Alta, Media o Baja |
| `Estado` | `EstadoDefecto` | Pendiente o Resuelto |
| `Historia` | `string` | Título de la historia relacionada |
| `SprintDeteccion` | `int` | Número del Sprint de detección |
| `SprintResolucion` | `*int` | Número del Sprint de resolución; `nil` si está Pendiente |

### `SprintReporte`

| Campo | Tipo | Regla |
|-------|------|-------|
| `Sprint` | `Sprint` | Goal, fechas, número y estado del Sprint |
| `Planificadas` | `[]HistoriaReporte` | Completadas + arrastradas al cierre |
| `Metricas` | `Metricas` | Valores de HU-03 para ese Sprint |
| `DefectosDetectados` | `[]DefectoReporte` | Defectos con `sprint_deteccion_id` = Sprint |
| `DefectosResueltos` | `[]DefectoReporte` | Defectos con `sprint_resolucion_id` = Sprint |

### `SprintActivoReporte`

| Campo | Tipo | Regla |
|-------|------|-------|
| `Sprint` | `Sprint` | Sprint en estado Activo |
| `Historias` | `[]HistoriaReporte` | Historias actualmente asignadas |

### `ReporteSprint`

| Campo | Tipo | Regla |
|-------|------|-------|
| `Proyecto` | `Project` | Datos del proyecto del Sprint |
| `Sprint` | `SprintReporte` | Sprint finalizado solicitado |

### `ReporteProyecto`

| Campo | Tipo | Regla |
|-------|------|-------|
| `Proyecto` | `Project` | Nombre, descripción, fechas |
| `Estado` | `EstadoProyecto` | Estado derivado (HU-13) |
| `Integrantes` | `[]Member` | Equipo del proyecto; puede estar vacío |
| `ProductBacklog` | `[]HistoriaReporte` | Historias sin Sprint, en orden de HU-15 |
| `Sprints` | `[]SprintReporte` | Sprints finalizados, en orden de `id` |
| `SprintActivo` | `*SprintActivoReporte` | `nil` si el proyecto no tiene Sprint Activo |
| `Metricas` | `Metricas` | Totales calculados solo con Sprints finalizados |
| `HorasPorStoryPoint` | `*float64` | Parámetro con el que se calcularon las métricas |
| `SinSprintsFinalizados` | `bool` | `true` cuando no hay Sprints finalizados |

## Reglas de armado

1. El reporte de Sprint solo se arma si el Sprint existe y está `Finalizado`; si no existe,
   `ErrSprintNoEncontrado`; si no está finalizado, `ErrSprintNoFinalizado`.
2. El reporte de proyecto solo se arma si el proyecto existe; si no,
   `ErrProyectoNoEncontrado`.
3. `Planificadas` = historias completadas del Sprint (`historias.sprint_id = sprint_id` y
   estado `Completada`) + historias de `historias_arrastradas` para ese Sprint. Las
   arrastradas usan la estimación congelada en `historias_arrastradas.estimacion_sp`.
4. Los defectos se filtran por proyecto y luego por Sprint de detección o de resolución. Un
   defecto puede aparecer en el reporte de un Sprint por detección y en otro por resolución.
5. Las métricas de cada Sprint y del proyecto salen de `domain.CalcularMetricas`; el reporte
   no recalcula valores.
6. Los totales del proyecto usan únicamente los Sprints finalizados. El Sprint Activo no
   aporta métricas, esfuerzo ni defectos a los totales.
7. Si no hay Sprints finalizados, `SinSprintsFinalizados` es `true` y el generador escribe
   "sin datos" en métricas, esfuerzo y defectos.
8. `Metricas` y `HorasPorStoryPoint` no se guardan: se recalculan en cada solicitud.
