# Data Model: Fechas y Estado del Proyecto (HU-13)

## Sin cambios en la persistencia

Esta historia **no agrega tablas ni columnas ni campos**. Se reutilizan las entidades y
tablas existentes:

- `projects` (HU-04): `id, name, description, start_date, end_date, created_at`.
- `sprints` (HU-05): `id, proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado`.

`internal/repository/migrate.go` y los adaptadores `sqlite_project.go` / `sqlite_sprint.go`
**no se modifican**. La entidad `domain.Project` **no recibe ningún campo `Estado`**.

## Estado del Proyecto (valor derivado, no almacenado)

| Valor | Significado |
|---|---|
| `Planificado` | El proyecto aún no inició ejecución: sin Sprints iniciados y sin alcanzar la fecha de inicio (o sin fecha de inicio cargada). |
| `En curso` | Hay un Sprint `Activo`, o hay un Sprint `Finalizado` junto a un `Pendiente`, o (sin Sprints iniciados) la fecha actual está dentro del rango del proyecto. |
| `Finalizado` | Todos los Sprints existentes están `Finalizado`, o (sin Sprints iniciados) la fecha actual superó la fecha de fin. |

**Representación**: `domain.EstadoProyecto string` con constantes
(`ProyectoPlanificado`/`ProyectoEnCurso`/`ProyectoFinalizado`) y método `EsValida()`, mismo
patrón que `domain.Estado`, `domain.Prioridad` y `domain.EstadoSprint`. **No se persiste**.

## Regla de cálculo (función pura)

`CalcularEstadoProyecto(p domain.Project, sprints []domain.Sprint, ahora time.Time) domain.EstadoProyecto`

Entradas: el `Project` (con `FechaInicio`/`FechaFin` opcionales), los `Sprint` del proyecto
(con `Estado` ∈ {`Pendiente`, `Activo`, `Finalizado`}) y el instante `ahora`.

Algoritmo (precedencia: Sprints sobre fechas):

1. ¿Algún Sprint con `Estado == Activo`? → `En curso`.
2. ¿Hay ≥ 1 Sprint y todos están `Finalizado`? → `Finalizado`.
3. ¿Hay al menos un Sprint `Finalizado` (sin `Activo` y sin ser todos `Finalizado`, es
   decir, conviviendo con al menos un `Pendiente`)? → `En curso`.
4. Si no (sin ejecución iniciada), decide la fecha:
   - `FechaInicio == nil` o `ahora` < `FechaInicio` → `Planificado`.
   - `FechaFin != nil` y `ahora` > `FechaFin` → `Finalizado`.
   - en otro caso (dentro del rango, bordes inclusivos) → `En curso`.

Reglas de borde:
- Los límites `[FechaInicio, FechaFin]` son **inclusivos**: el día de inicio y el día de fin
  cuentan como `En curso`.
- Un Sprint `Pendiente` no es evidencia de ejecución.
- Un proyecto sin fechas y sin Sprints iniciados ⇒ `Planificado`.

## Consumo por el caso de uso

`ObtenerEstadoProyecto.Ejecutar(ctx, proyectoID, ahora)`:
1. `projects.ObtenerPorID(ctx, proyectoID)` → si no existe, `domain.ErrProyectoNoEncontrado`.
2. `sprints.ListarPorProyecto(ctx, proyectoID)`.
3. `domain.CalcularEstadoProyecto(proyecto, sprints, ahora)` → devuelve el `EstadoProyecto`.

Solo lectura: no invoca `Guardar`/`Actualizar` en ningún puerto.

## Relaciones (sin cambios)

```text
Project (1) ----< (N) Sprint
```

El estado del proyecto es una **proyección calculada** sobre esa relación; no se materializa.
