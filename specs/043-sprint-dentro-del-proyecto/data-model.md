# Data Model: Sprints dentro del plazo del proyecto (#168)

**Sin cambio de esquema ni migración.** Las tablas, columnas y restricciones son las de [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md): `proyectos` en §4.1 y `sprints` en §4.5. Este documento solo describe la regla nueva sobre esas entidades.

## Entidades involucradas

| Entidad | Campos que usa la regla | Fuente |
|---|---|---|
| Proyecto (`domain.Project`) | `FechaInicio *time.Time`, `FechaFin *time.Time` (opcionales, HU-04) | `MODELO-DE-DATOS.md` §4.1 |
| Sprint (`domain.Sprint`) | `ProyectoID`, `FechaInicio`, `FechaFin` (se cargan al iniciar) | `MODELO-DE-DATOS.md` §4.5 |

## Regla nueva

`Project.AdmiteSprint(inicio, fin)`:

| Condición del proyecto | Se exige | Si no se cumple |
|---|---|---|
| Tiene `FechaInicio` | `inicio` ≥ `FechaInicio` | `ValidationError{Campo: "fecha_inicio"}` |
| Tiene `FechaFin` | `fin` ≤ `FechaFin` | `ValidationError{Campo: "fecha_fin"}` |
| Sin la fecha correspondiente | nada | — |

Se evalúa primero el inicio. Bordes incluidos. Comparación por día (todas las fechas a las 00:00 UTC, ver [research.md §3](research.md)).

## Transición de estado afectada

`Pendiente → Activo` (`IniciarSprint`) gana una precondición más: el rango del Sprint cabe en el plazo del proyecto. Si no cabe, el Sprint queda Pendiente, sin Goal ni fechas. Las demás transiciones (`Activo → Finalizado`) no cambian.

## Cambio en `MODELO-DE-DATOS.md`

Solo documentación: una fila nueva en la tabla de §5 ("Dónde vive cada regla"):

| Regla | Base | Servicio / dominio |
|---|---|---|
| Las fechas de un Sprint, al iniciarlo, caen dentro del plazo del proyecto (bordes incluidos; solo contra las fechas que el proyecto tenga) | — | Dominio (`Project.AdmiteSprint`), llamado desde `IniciarSprint`: `ValidationError` sobre `fecha_inicio` o `fecha_fin` (#168). Editar después las fechas del proyecto no revisa los Sprints existentes (#302) |

La base no puede hacerla cumplir con un `CHECK` (la regla cruza dos tablas), así que, por la regla de oro de §5, vive en el dominio.
