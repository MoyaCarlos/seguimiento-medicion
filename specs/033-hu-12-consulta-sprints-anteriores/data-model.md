# Data Model: Consulta de Sprints Anteriores (HU-12)

**Sin cambios de esquema.** Las tablas están definidas en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md): `sprints` (§4.5), `historias_arrastradas` (§4.6), `historias` (§4.4), `registros_horas` (§4.8) y `defectos` (§4.9). HU-12 es de solo lectura (§8). Las métricas no se guardan (§6).

## Vistas de lectura (`internal/domain`)

Todas son valores calculados al consultar. No tienen tabla propia.

### HistoriaComprometida

| Campo | Tipo | Origen |
|---|---|---|
| `ID` | `int64` | `historias.id` |
| `Titulo` | `string` | `historias.titulo` (valor actual) |
| `EstimacionSP` | `*int` | Completada: `historias.estimacion_sp`. Arrastrada: `historias_arrastradas.estimacion_sp` (foto). `nil` = "sin estimar" |
| `Completada` | `bool` | `true` si viene de `historias`; `false` si viene de `historias_arrastradas` |
| `Minutos` | `int` | Minutos imputados en este Sprint; `0` si no tiene |

### DefectoDeSprint

`DefectoDetalle` (HU-08) más dos marcas:

| Campo | Regla |
|---|---|
| `DetectadoEnSprint` | `SprintDeteccionID == sprint consultado` |
| `ResueltoEnSprint` | `SprintResolucionID` no es `nil` e igual al Sprint consultado |

Al menos una de las dos es verdadera. Un defecto aparece una sola vez aunque tenga ambas. `Estado` es el actual: un defecto de un Sprint cerrado figura como `Resuelto` si se resolvió después.

### SprintHistorial y DetalleSprint

- **SprintHistorial**: `Sprint` + métricas de HU-03.
- **DetalleSprint**: `SprintHistorial` + `[]HistoriaComprometida` + `[]DefectoDeSprint`.

El tipo de las métricas es el de HU-03 (ver R2).

## Reglas puras (`internal/domain`)

| Función | Regla |
|---|---|
| `SprintsFinalizados(sprints)` | Devuelve solo los `Finalizado`, ordenados por `Numero` ascendente; nunca `nil` |
| `ValidarSprintFinalizado(sprint)` | Devuelve `ErrSprintNoFinalizado` si no está `Finalizado` |
| `ArmarDetalleSprint(...)` | Une las horas con las historias (cero si no tiene) y marca los defectos |

## Puerto (`internal/repository`)

```go
type HistorialSprintRepository interface {
	HistoriasComprometidas(ctx context.Context, sprintID int64) ([]domain.HistoriaComprometida, error)
	MinutosPorHistoria(ctx context.Context, sprintID int64) (map[int64]int, error)
	DefectosDeSprint(ctx context.Context, sprintID int64) ([]domain.DefectoDetalle, error)
}
```

| Método | Fuente | Orden |
|---|---|---|
| `HistoriasComprometidas` | `historias` con `sprint_id = S` y `estado = 'Completada'` ∪ `historias_arrastradas` con `sprint_id = S`, unida a `historias` para el título | Completadas, luego arrastradas; por `id` |
| `MinutosPorHistoria` | `SUM(minutos)` de `registros_horas` con `sprint_id = S`, agrupado por historia | — |
| `DefectosDeSprint` | `defectos` con `sprint_deteccion_id = S` o `sprint_resolucion_id = S`, unida a `historias` para el título | Por `id` |

## Services (`internal/service`)

Cada uno declara interfaces angostas para sus dependencias, como `ConsultarSprints`.

### ConsultarHistorialSprints

1. Pide las métricas del proyecto a HU-03, que responde `ErrProyectoNoEncontrado` si no existe.
2. Lista los Sprints del proyecto y aplica `SprintsFinalizados`.
3. Une cada Sprint con sus métricas. Si falta alguna, devuelve un error interno, porque es una inconsistencia entre fuentes.

### ConsultarDetalleSprint

1. Verifica el proyecto: `ErrProyectoNoEncontrado`.
2. Obtiene el Sprint: `ErrSprintNoEncontrado` si no existe o si es de otro proyecto.
3. Aplica `ValidarSprintFinalizado`.
4. Lee las métricas, las historias comprometidas, los minutos y los defectos, y aplica `ArmarDetalleSprint`.

Ningún service escribe. Cualquier error de lectura se devuelve sin datos parciales.
