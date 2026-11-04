# Data Model: Cálculo de Métricas del Proyecto (HU-03)

**Sin cambios de esquema ni migraciones.** Las tablas y restricciones son las de
[`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md); este documento enlaza y no las redefine.

| Dato | Fuente | Sección |
|---|---|---|
| Horas por Story Point | `proyectos.horas_por_story_point` | [§4.1](../../MODELO-DE-DATOS.md#41-proyectos--hu-04-hu-13-hu-03-antes-projects) |
| Completadas | `historias` con `sprint_id = S` y `estado = 'Completada'` | §4.4 |
| Arrastradas y su foto | `historias_arrastradas` con `sprint_id = S` | §4.6 |
| Horas reales | `registros_horas.minutos` por `sprint_id` | §4.8 |
| Defectos | `defectos` por `sprint_deteccion_id` y `sprint_resolucion_id` | §4.9 |
| Número del Sprint | Derivado: posición entre todos los Sprints del proyecto, por `id` | spec 024 |

## Cambio en una entidad existente

### `domain.Project`

| Campo nuevo | Tipo | Regla |
|---|---|---|
| `HorasPorStoryPoint` | `*float64` | `nil` = sin definir. Si está definido: finito y mayor que cero |

`ConDatosEditados` no cambia y conserva el campo. Método nuevo:

`ConHorasPorStoryPoint(valor *string) (Project, error)`: `nil` lo deja sin definir; un texto que no
sea un número finito mayor que cero devuelve `ValidationError` con campo `horas_por_story_point`.

## Tipos nuevos (`internal/domain/metricas.go`)

No se guardan.

### DatosMetricasSprint (entrada)

Una fila por Sprint finalizado, leída por el repositorio.

| Campo | Tipo | Origen |
|---|---|---|
| `SprintID`, `Numero` | `int64`, `int` | `sprints` |
| `HistoriasCompletadas` | `int` | `COUNT` de las completadas |
| `HistoriasArrastradas` | `int` | `COUNT` de las arrastradas |
| `SPCompletados` | `int` | `COALESCE(SUM(estimacion_sp), 0)` de las completadas |
| `SPArrastrados` | `int` | `COALESCE(SUM(estimacion_sp), 0)` de la foto |
| `MinutosReales` | `int` | `COALESCE(SUM(minutos), 0)` |
| `DefectosDetectados`, `DefectosResueltos` | `int` | `COUNT` |

### Metricas (salida)

Misma forma para un Sprint y para el proyecto.

| Campo | Tipo | Sprint | Proyecto |
|---|---|---|---|
| `SPPlanificados` | `int` | completados + arrastrados | suma |
| `SPCompletados` | `int` | completados | suma |
| `Velocidad` | `float64` | `SPCompletados` | promedio por Sprint finalizado; 0 sin Sprints |
| `PorcentajeCompletadas` | `float64` | completadas / comprometidas × 100; 0 sin comprometidas | sobre los totales |
| `HorasEstimadas` | `*float64` | `SPPlanificados` × parámetro; `nil` sin parámetro | sobre los totales |
| `HorasReales` | `float64` | minutos / 60 | total de minutos / 60 |
| `DesviacionHoras` | `*float64` | reales − estimadas; `nil` sin parámetro | sobre los totales |
| `DesviacionPorcentaje` | `*float64` | desviación / estimadas × 100; `nil` sin parámetro o con estimadas en 0 | sobre los totales |
| `DefectosDetectados`, `DefectosResueltos` | `int` | conteo | suma |

Los decimales se redondean a dos solo al final, y nunca se devuelve `-0` ([R4](./research.md)).

### MetricasDeSprint y MetricasProyecto

- **MetricasDeSprint**: `SprintID`, `Numero`, `Metricas`.
- **MetricasProyecto**: `HorasPorStoryPoint *float64`, `Proyecto Metricas`, `Sprints []MetricasDeSprint` (vacío, no `nil`, y ordenado por número).

## Función pura

`CalcularMetricas(datos []DatosMetricasSprint, horasPorSP *float64) MetricasProyecto`

## Puerto y service

`repository.DatosMetricasRepository`:
`ListarDatosSprintsFinalizados(ctx, proyectoID) ([]domain.DatosMetricasSprint, error)`, ordenado por
`sprints.id` y con un slice vacío si no hay Sprints finalizados.

`service.CalcularMetricas.Ejecutar(ctx, proyectoID) (domain.MetricasProyecto, error)`:

1. `proyectos.ObtenerPorID`: `ErrProyectoNoEncontrado` si no existe (FR-014).
2. `ListarDatosSprintsFinalizados`.
3. `domain.CalcularMetricas` con `proyecto.HorasPorStoryPoint`.

No escribe nada (FR-015). HU-12 y HU-10 consumen este service.

## Regla de Story Points

| Sprint de la historia | Sin estimar → carga | Con estimación → otro valor o vacío | Sin cambio de valor |
|---|---|---|---|
| Ninguno (Product Backlog) | Permitido | Permitido | Permitido |
| `Pendiente` | Permitido | Permitido | Permitido |
| `Activo` | Permitido | `ErrEstimacionCongelada` | Permitido |
| `Finalizado` | **`ErrCargaEnSprintFinalizado` (nuevo)** | `ErrEstimacionCongelada` | Permitido |

## Documentos a actualizar en esta historia

| Documento | Qué cambia |
|---|---|
| `MODELO-DE-DATOS.md` §4.4, §4.6, §5 y D2b | Se prohíbe también la primera carga con Sprint `Finalizado` |
| `MODELO-DE-DATOS.md` §6 | Horas estimadas: SP planificados (ya no "lo define la spec") |
| `docs/PRODUCT-BACKLOG.md`, decisiones transversales, HU-14 Esc. 4 y HU-15 Esc. 3 | Regla nueva |
| `docs/PRODUCT-BACKLOG.md`, nota de HU-03 | Excepción: cambiar el parámetro recalcula los Sprints cerrados |
| `specs/009-…/spec.md` y `specs/031-…/spec.md` | Regla nueva (Clarifications, escenarios, edge cases y FR) |
| `specs/002-…/contracts/openapi.yaml` | Campo opcional en la edición y en la respuesta del proyecto |

La spec 033 (HU-12) asume la regla anterior y se ajusta en su propia rama.
