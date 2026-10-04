# Research: Apertura y Cierre de Sprints

Sin `NEEDS CLARIFICATION` pendientes en el Technical Context del plan — las
decisiones técnicas ya estaban fijadas a nivel de proyecto (`AGENTS.md`,
constitución) o resueltas en `spec.md`. Este documento registra el
razonamiento, no investigación externa nueva.

## Decisión: Modelado del estado del Sprint

**Decisión**: tipo `EstadoSprint string` con constantes
(`Pendiente`/`Activo`/`Finalizado`) y método `EsValida()`, igual patrón que
`domain.Prioridad` y `domain.Estado` ya existentes en HU-01.

**Rationale**: consistencia de estilo dentro del mismo paquete `domain`;
evita introducir un enfoque distinto (ej. `iota` numérico) para el mismo
tipo de problema ya resuelto.

**Alternatives considered**: `iota` con tipo numérico — descartado, ya hay
un patrón establecido con strings legibles directamente en SQLite sin
tabla de traducción.

## Decisión: Validación de "un solo Activo/Pendiente por proyecto"

**Decisión**: la validación vive en la capa `service` (caso de uso
`IniciarSprint`/`CrearSprint`), consultando `SprintRepository` antes de
persistir — no es una constraint a nivel de base de datos (ej. UNIQUE
INDEX).

**Rationale**: la regla de negocio debe poder testearse con TDD puro sobre
el dominio/servicio usando el fake en memoria, sin depender de que SQLite
esté levantado. Una constraint de base de datos la volvería invisible para
esos tests y correspondería a `internal/repository`, no al servicio.

**Alternatives considered**: UNIQUE INDEX en SQLite sobre
`(proyecto_id, estado)` para `estado IN ('Pendiente','Activo')` — se
descarta como mecanismo único porque el error de SQLite sería genérico (no
distinguible de otros conflictos) y no directamente testeable con TDD de
dominio; podría agregarse después como salvaguarda adicional, pero no
reemplaza la validación en `service` (YAGNI: no se agrega ahora sin una
necesidad concreta de esa capa extra de seguridad).

## Decisión: Arrastre de historias no completadas al cerrar

**Decisión**: `CerrarSprint` (service) recibe el `SprintID`, pide a
`BacklogRepository` las historias asignadas a ese Sprint, y por cada una
que no esté en estado completada, actualiza su asignación a "sin Sprint".

**Rationale**: reutiliza `BacklogRepository` (de HU-01) en vez de crear un
mecanismo de "eventos" o callbacks entre historias — más simple y
suficiente para el volumen de datos de un TP (YAGNI: no se agrega un bus de
eventos para un caso de uso que se resuelve con una consulta directa).

**Alternatives considered**: trigger a nivel de SQLite — descartado, misma
razón que el punto anterior (no testeable con TDD de dominio/servicio en
memoria).

**Nota de dependencia cruzada**: esto requiere que `BacklogRepository` (ya
existente, HU-01) tenga un método para listar/actualizar historias por
Sprint — se agrega en `internal/repository/backlog_repository.go` como
parte de las tasks de esta historia, sin tocar la lógica ya implementada de
HU-01.
