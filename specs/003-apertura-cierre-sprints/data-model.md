# Data Model: Apertura y Cierre de Sprints

## Sprint (nueva entidad)

| Campo | Tipo | Regla |
|---|---|---|
| ID | int64 | autoincremental |
| ProyectoID | int64 | obligatorio, `> 0`, debe existir (vía `ProjectRepository.ObtenerPorID`; si no, `ErrProyectoNoEncontrado`) |
| SprintGoal | string | obligatorio al **iniciar** (no al crear) |
| FechaInicio | time.Time | obligatoria al iniciar |
| FechaFin | time.Time | obligatoria al iniciar; estrictamente posterior a FechaInicio (no se permite un Sprint de un día) |
| Estado | EstadoSprint | `Pendiente` (default al crear) → `Activo` → `Finalizado` |

**Invariantes** (validados en `NewSprint` / métodos de transición, patrón Factory):
- No se puede crear un Sprint para un `ProyectoID` inexistente (FR-002).
- No se puede crear un Sprint "Pendiente" si el proyecto ya tiene otro "Pendiente" (FR-011).
- No se puede iniciar un Sprint si el proyecto ya tiene otro "Activo" (FR-004).
- `FechaFin` > `FechaInicio` al iniciar (FR-005).
- Solo se puede cerrar un Sprint en estado "Activo" (FR-007).

**Transiciones de estado**:
```
Pendiente --(Iniciar: Goal + fechas válidas)--> Activo --(Cerrar)--> Finalizado
```
No hay transición inversa ni edición de un Sprint ya "Activo" o "Finalizado"
(fuera de alcance — YAGNI, no lo pide ningún criterio de aceptación).

## BacklogItem (extensión mínima — HU-01 ya existente)

Se agrega **un solo campo nuevo**, sin tocar nada de lo ya implementado en HU-01:

| Campo nuevo | Tipo | Regla |
|---|---|---|
| SprintID | *int64 | `nil` = sin Sprint asignado (está en el Product Backlog) |

**Importante — alcance acotado**: esta historia (HU-05) agrega el campo y
los dos métodos mínimos de `BacklogRepository` que necesita para poder
"soltar" las historias no completadas al cerrar un Sprint
(`ListarPorSprint`, `QuitarDeSprint`). **No implementa** la lógica de
asignar una historia a un Sprint por primera vez, ni sus validaciones (ej.
"no asignar una historia ya Terminada") — eso es el alcance completo de
**HU-06** (Movimiento de Historias al Sprint Backlog), que se construye
sobre este mismo campo más adelante.

## Relaciones

```
Project (1) ----< (N) Sprint
Sprint  (1) ----< (N) BacklogItem   [opcional: BacklogItem.SprintID puede ser nil]
```
