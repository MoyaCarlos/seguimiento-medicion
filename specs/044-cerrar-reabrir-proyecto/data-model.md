# Data Model: Cerrar y reabrir un proyecto (spec 044)

La fuente de verdad del esquema es [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md). Esta
página resume el cambio y enlaza las secciones; no redefine tablas.

## Cambio de esquema (migración 008)

**Tabla `proyectos`** ([§4.1](../../MODELO-DE-DATOS.md#41-proyectos--hu-04-hu-13-hu-03-antes-projects)):
se agrega una columna.

| Columna | Tipo | Notas |
|---|---|---|
| `fecha_cierre` | TEXT? | **Nueva** (#300, cierra D16). `YYYY-MM-DD`. `NULL` ⇔ proyecto abierto. `CHECK (date(fecha_cierre) IS fecha_cierre)` (D12) |

- Los proyectos existentes quedan con `NULL` (abiertos).
- Sin relación con `fecha_inicio`/`fecha_fin`: se puede cerrar antes del inicio o después
  del fin (el cierre es una decisión del Scrum Master, no una fecha de plan).
- El estado del proyecto **sigue sin guardarse** ([§6](../../MODELO-DE-DATOS.md#6-qué-no-se-guarda-y-de-dónde-se-calcula)):
  su fuente pasa a ser `proyectos.fecha_cierre` + `proyectos.fecha_*` + `sprints.estado`.

## Entidad de dominio `Project`

| Campo | Tipo | Regla |
|---|---|---|
| `FechaCierre` | `*time.Time` | `nil` mientras el proyecto está abierto; día civil (00:00 UTC) del cierre |

### Transiciones

```text
            Cerrar(ahora, sprints)
  ABIERTO ─────────────────────────▶ CERRADO
 (FechaCierre=nil)                  (FechaCierre=hoy)
     ▲                                  │
     └──────────── Reabrir() ───────────┘
```

| Desde | Acción | Condición | Resultado / error (409) |
|---|---|---|---|
| Abierto | `Cerrar` | ningún Sprint `Activo` | Cerrado, `FechaCierre` = hoy |
| Abierto | `Cerrar` | hay un Sprint `Activo` | `ErrSprintActivoAlCerrarProyecto` «Cerrá el Sprint activo antes de cerrar el proyecto» |
| Cerrado | `Cerrar` | — | `ErrProyectoYaCerrado` «El proyecto ya está cerrado» |
| Cerrado | `Reabrir` | — | Abierto, `FechaCierre = nil` |
| Abierto | `Reabrir` | — | `ErrProyectoNoCerrado` «El proyecto no está cerrado» |
| Cerrado | cualquier escritura | — | `ErrProyectoCerrado` «El proyecto está cerrado» |

### Estado derivado (`CalcularEstadoProyecto`)

1. **Nuevo**: `FechaCierre != nil` ⇒ «Finalizado».
2. En adelante, sin cambios (HU-13): Sprint `Activo` ⇒ «En curso»; todos `Finalizado` ⇒
   «Finalizado»; algún `Finalizado` ⇒ «En curso»; si no, deciden las fechas.

## Dónde vive cada regla (para [§5](../../MODELO-DE-DATOS.md#5-dónde-vive-cada-regla))

| Regla | Base | Servicio / dominio |
|---|---|---|
| Formato de `fecha_cierre` | `CHECK` (D12) | El dominio la arma con el día civil |
| Cerrar solo sin Sprint Activo, y no dos veces | `UPDATE` condicional (cubre la carrera) | `Project.Cerrar` da el mensaje (409) |
| Reabrir solo un proyecto cerrado | `UPDATE` condicional | `Project.Reabrir` (409) |
| Proyecto cerrado = solo lectura | — | Guardia `exigirProyectoAbierto` en el control de acceso de HU-21 (409). Límite declarado: una escritura concurrente con el cierre puede completarse ([research.md §6](research.md)) |
