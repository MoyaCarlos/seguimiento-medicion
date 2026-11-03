# Data Model: Movimiento de Historias al Sprint Backlog (HU-06)

**Sin cambios de esquema.** La fuente única es [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md);
este documento enlaza a sus secciones y no redefine tablas.

## Tablas que usa

| Tabla | Sección | Uso en HU-06 |
|---|---|---|
| `historias` | [§4.4](../../MODELO-DE-DATOS.md#44-historias--hu-01-hu-15-hu-06-hu-11-hu-14-existe-hoy-como-backlog_items) | Se escribe **solo** `sprint_id` (de `NULL` al Sprint destino). Se lee para el Sprint Backlog (`WHERE sprint_id = ?`). |
| `sprints` | [§4.5](../../MODELO-DE-DATOS.md#45-sprints--hu-05-existe-hoy-mismo-nombre) | Solo lectura: `proyecto_id` y `estado` del Sprint destino. |

Mapeo de HU-06 en el documento: [§8](../../MODELO-DE-DATOS.md#8-mapa-historia--tablas), fila "HU-06 Mover al Sprint" → `historias.sprint_id`,
migración 001.

## Restricciones de base que intervienen

- `FOREIGN KEY (sprint_id, proyecto_id) REFERENCES sprints (id, proyecto_id)`: impide asociar una
  historia a un Sprint de otro proyecto. El servicio lo valida antes (`ErrSprintDeOtroProyecto`,
  400); si una violación llega igual a la base, se traduce a `ErrIntegridadReferencial` (409),
  como fija §5.
- `CHECK (estado = 'Nueva' OR sprint_id IS NOT NULL)`: toda historia sin Sprint está en `Nueva`.
  La asignación no cambia `estado`, así que nunca lo viola.

## Entidades de dominio (sin campos nuevos)

- `domain.BacklogItem`: se usan `ID`, `ProyectoID`, `Estado` y `SprintID`. Ningún otro campo
  cambia al asignar (FR-003, D2b: los Story Points no se tocan).
- `domain.Sprint`: se usan `ID`, `ProyectoID` y `Estado`.

## Regla de asignación (`domain.ValidarAsignacion`)

Se evalúa después de comprobar que la historia y el Sprint existen (404 primero, FR-008). Devuelve
el primer error en este orden (FR-009):

| # | Condición | Error | Campo | Mensaje |
|---|---|---|---|---|
| 1 | `historia.ProyectoID != sprint.ProyectoID` | `ErrSprintDeOtroProyecto` | `sprint_id` | el Sprint no pertenece al proyecto de la historia |
| 2 | `sprint.Estado == Finalizado` | `ErrSprintFinalizado` | `sprint_id` | un Sprint finalizado no admite historias |
| 3 | `historia.Estado == Completada` | `ErrHistoriaNoPendiente` | `estado` | solo se pueden planificar historias pendientes |
| 4 | `historia.SprintID != nil` | `ErrHistoriaYaAsignada` | `sprint_id` | la historia ya pertenece a un Sprint |

`ErrHistoriaYaAsignada` lo devuelve también el repositorio cuando el `UPDATE` condicional no
afecta filas (carrera, [research.md](./research.md) R2).

## Transición de estado

```text
Historia en Product Backlog            Historia en Sprint Backlog
(sprint_id = NULL, estado = Nueva) ──► (sprint_id = S, estado = Nueva)
                                       S en estado Pendiente o Activo
```

No hay transición inversa en esta historia (D5). Volver al Product Backlog solo ocurre al cerrar
el Sprint con la historia no completada (HU-05, §4.6).

## Vistas (no se guardan)

- **Product Backlog** (HU-15): `historias WHERE proyecto_id = ? AND sprint_id IS NULL`.
- **Sprint Backlog** (HU-06): `historias WHERE sprint_id = ?`. Ambas con el mismo orden: prioridad
  `M, S, C, W` con `CASE` explícito y luego `id` (§4.4).

## Métricas (sin cambios en esta historia)

Una historia asignada con el Sprint `Activo` cuenta como planificada al cerrarlo (D2c, §4.6). Lo
resuelven el cierre (HU-05) y las métricas (HU-03); HU-06 no agrega cálculo.
