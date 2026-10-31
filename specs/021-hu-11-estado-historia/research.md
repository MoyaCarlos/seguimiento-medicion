# Research: Actualización de Estado de una Historia en el Sprint (HU-11)

No quedaron `NEEDS CLARIFICATION` en el Technical Context. Este documento registra las decisiones
técnicas (las acordadas por el equipo y los ajustes que surgieron al revisar el código) con sus
alternativas.

## R1. La regla de la transición es una función pura del dominio

- **Decision**: `domain.ValidarCambioDeEstado(h BacklogItem, s *Sprint, nuevo Estado) error` en
  `internal/domain/estado_historia.go`. Evalúa en este orden y devuelve el primer error:

  | # | Condición | Error | Campo | Mensaje |
  |---|---|---|---|---|
  | 1 | `!nuevo.EsValida()` | `ErrEstadoInvalido` | `estado` | el estado debe ser Nueva, En progreso o Completada |
  | 2 | `h.SprintID == nil` y el destino no es `Nueva` | `ErrHistoriaSinSprint` | `estado` | primero debe asignarse a un Sprint activo |
  | 3 | `s == nil` o `s.Estado != SprintActivo` | `ErrSprintNoActivoParaEstado` | `estado` | solo se cambia el estado de historias de un Sprint activo |
  | 4 | `nuevo == h.Estado` | `ErrEstadoRepetido` | `estado` | la historia ya está en ese estado |
  | 5 | `h.Estado == Completada` o `nuevo == Nueva` | `ErrRetrocesoEstado` | `estado` | el estado de una historia solo avanza: no se puede volver a un estado anterior |
  | 6 | `h.Estado == Nueva` y `nuevo == Completada` | `ErrSaltoDirecto` | `estado` | primero debe pasar por En progreso |

  Una historia sin Sprint siempre está en `Nueva` (regla transversal del modelo), así que el
  destino `Nueva` sobre ella cae en el caso 4 (repetido) y no en el 2; es lo que fija la
  precedencia de la spec (FR-018). Después del caso 3, las combinaciones que quedan son solo
  `Nueva → En progreso`, `Nueva → Completada` y `En progreso → Completada`; por eso el retroceso
  se expresa con `h.Estado == Completada || nuevo == Nueva`.
- **Rationale**: el orden de precedencia es una regla de negocio (FR-018) y se prueba sin base con
  una tabla de casos que combine varias violaciones a la vez. Las existencias (404) quedan fuera
  porque necesitan el repositorio: las resuelve el servicio antes.
- **Ajuste sobre lo acordado**: el Sprint entra como `*Sprint` y no por valor, porque la historia
  puede no tener Sprint y el servicio no tiene qué obtener. La regla de "sin Sprint" sigue en el
  dominio y se decide con `SprintID`; el puntero solo alimenta el estado del Sprint.
- **Alternatives considered**: un método `(*BacklogItem) CambiarA(nuevo Estado)` que mute el
  estado. Se descartó por la misma razón que en HU-06: la escritura real es condicional en la base
  (R2) y mutar antes de persistir adelantaría el estado en memoria.

## R2. Escritura condicional en una sola sentencia

- **Decision**: método `CambiarEstado(ctx, id int64, desde, hacia domain.Estado, fechaFinalizacion *time.Time) error`
  en `SQLiteBacklogRepository`:

  ```sql
  UPDATE historias SET estado = ?, fecha_finalizacion = ? WHERE id = ? AND estado = ?
  ```

  - `fecha_finalizacion` va `NULL` salvo que `hacia == Completada`; se serializa con el helper
    `fechaASQL` que ya usa el repositorio de Sprints (formatea a `AAAA-MM-DD`, así el componente
    horario del reloj no llega a la base).
  - Error de la sentencia → `errorDeEscritura` (una violación de `CHECK` llega como error genérico,
    §5 de `MODELO-DE-DATOS.md`).
  - `RowsAffected() == 0` → `domain.ErrConflictoEstadoHistoria` (error de conflicto, 409).
- **Rationale**: si dos cambios de la misma historia llegan juntos, los dos pasan
  `ValidarCambioDeEstado`, pero solo uno encuentra `estado = desde`; el otro recibe el conflicto.
  Es el mismo patrón que `AsignarASprint` (HU-06) y los `UPDATE` condicionales de los fixes #145 y
  #158: sin transacción explícita ni bloqueo.
- **Alternatives considered**: `SELECT` y después `UPDATE` (deja la carrera abierta);
  `BEGIN IMMEDIATE` con un lock (más código para una app mono-usuario).
- **Riesgo aceptado**: si la historia se borrara entre la lectura y el `UPDATE`, la respuesta sería
  409 en vez de 404. Hoy no existe endpoint de borrado de historias.

## R3. Puerto propio y `fecha_finalizacion` en el modelo de lectura (ajuste sobre lo acordado)

- **Decision**: puerto nuevo en `internal/repository/backlog_repository.go`:

  ```go
  type CambioEstadoHistoriaRepository interface {
      ObtenerPorID(ctx context.Context, id int64) (domain.BacklogItem, error)
      CambiarEstado(ctx context.Context, id int64, desde, hacia domain.Estado, fechaFinalizacion *time.Time) error
  }
  ```

  Lo implementa el mismo `SQLiteBacklogRepository`. Se agrega `FechaFinalizacion *time.Time` a
  `domain.BacklogItem`, `columnasHistoria` suma `fecha_finalizacion` y `escanearHistoria` la lee
  como `sql.NullString` y la parsea con `fechaDesdeSQL` (helpers ya existentes de
  `sqlite_sprint.go`). `Guardar` no la escribe: una historia nueva no tiene fecha.
- **Rationale**: ISP, igual que HU-06: `HistoriaRepository` es la edición y no conoce transiciones.
  `ObtenerPorID` ya aparece en otros puertos, pero es el mismo método del mismo tipo concreto, así
  que no hay código duplicado. El campo nuevo es necesario para devolver la fecha en la respuesta
  del cambio de estado (y en cualquier listado de historias que la incluya); la columna ya existe
  desde la 001.
- **Alternatives considered**: ampliar `HistoriaRepository` (obliga a `EditarHistoria` y sus fakes
  a conocer una operación que no usan); leer la fecha con una consulta aparte solo para la
  respuesta (dos consultas y un campo que igual falta en el modelo).

## R4. Servicio con reloj inyectable

- **Decision**: `CambiarEstadoHistoria` con `NewCambiarEstadoHistoria(historias, sprints, reloj func() time.Time)`:

  ```go
  func (s *CambiarEstadoHistoria) Ejecutar(ctx context.Context, id int64, nuevo domain.Estado) (domain.BacklogItem, error)
  ```

  1. `ObtenerPorID` (404 si no existe).
  2. Si `SprintID != nil`, `sprints.ObtenerPorID` para conocer su estado.
  3. `ValidarCambioDeEstado`.
  4. Si el destino es `Completada`, toma la fecha del reloj; si no, `nil`.
  5. `CambiarEstado` con el estado leído como `desde`.
  6. Devuelve la historia con el estado y la fecha ya actualizados.
- **Rationale**: el servicio solo orquesta (FR-001 a FR-003). El reloj inyectable hace
  determinista la fecha en los tests de service y en el escenario BDD (una fecha fija).
- **Nota de convención**: los otros services del proyecto reciben `ahora time.Time` como
  parámetro (`ObtenerEstadoProyecto`, `RevelarVotacion`, `ImputarHoras`). Acá se usa la inyección
  acordada por el equipo; para el handler no cambia nada, porque `main.go` entrega `time.Now`.
- **Alternatives considered**: pasar `time.Now()` desde el handler (convención actual); se
  descartó por la decisión del equipo, no por un problema técnico.

## R5. Ruta HTTP, DTO y conflicto 409

- **Decision**:
  - `PUT /backlog/{id}/estado` en `BacklogHandler.CambiarEstado`, con cuerpo `{"estado": "..."}`.
    Siguiendo el patrón de `Editar` y `AsignarSprint`, el handler verifica primero que la historia
    exista (`ObtenerHistoria`, 404) y recién después decodifica el cuerpo: una historia inexistente
    da 404 aunque el JSON sea inválido o falte `estado`.
  - `estado` ausente o vacío → 400 `{campo: "estado", mensaje: "el campo es obligatorio"}`. Un valor
    que no es uno de los tres estados lo rechaza el dominio (R1, caso 1) con el mismo campo.
  - Respuesta 200 con el mismo `historiaResponse` del Product Backlog, que **suma el campo aditivo
    `fecha_finalizacion`** (`*string`, `null` sin fecha), formateado con `time.DateOnly`.
  - `escribirError` suma un `case` para `domain.ErrConflictoEstadoHistoria` → **409 sin campo**,
    con el mensaje "la historia cambió de estado mientras tanto: recargá y volvé a intentar". El
    error se define como error plano (no `ValidationError`) para que no lo capture el primer
    `case`.
  - El constructor de `BacklogHandler` suma el service nuevo; hoy tiene un solo sitio de
    construcción en tests además de `main.go`.
- **Rationale**: la historia en la ruta hace que "historia inexistente" domine de forma natural
  (FR-007); `PUT` porque el cambio de estado reemplaza un atributo del recurso, igual que
  `PUT /backlog/{id}/sprint`. El 409 sigue la convención de `MODELO-DE-DATOS.md` §5 para
  conflictos de concurrencia.
- **Alternatives considered**: `PATCH /backlog/{id}` genérico (se mezcla con la edición de HU-15);
  `POST /backlog/{id}/estado` (la transición no crea un recurso).

## R6. Escenario BDD

- **Decision**: `features/cambiar_estado_historia.feature` con `@HU-11`, un escenario por criterio
  de aceptación de US1 y US2, sobre SQLite en memoria como el resto de los pasos
  (`steps_asignacion.go`, `steps_backlog.go`). Comparte el `scenarioContext` de HU-01 y usa los
  services de HU-05 (crear, iniciar, cerrar) y de HU-06 (asignar) para preparar cada caso; el
  estado inicial "En progreso" se llega con el propio service de HU-11, sin SQL a mano. La fecha
  del reloj se fija en el paso para poder afirmarla.
- **Rationale**: Principio III. Las cuatro categorías (normal, alternativo, límite, error) están
  cubiertas por la spec: normal `Nueva → En progreso`; alternativo `En progreso → Completada`;
  límite persistencia y conservación de Story Points; error sin Sprint, Sprint Pendiente,
  Sprint Finalizado, repetido, retroceso, salto directo, cuerpo inválido e historia inexistente.
  La carrera entre dos cambios simultáneos no se puede reproducir de forma determinista con
  Godog: se cubre con una prueba del repositorio contra SQLite, como en HU-06.
- **Alternatives considered**: sembrar la fecha con SQL (no probaría el reloj); dejar la carrera
  solo en el escenario (no es determinista).

## R7. Pantalla (US3)

- **Decision**: fase propia, **bloqueada** hasta que la pantalla de Sprints esté en `main`. Por
  cada historia del Sprint activo, un control para pasarla a "En progreso" o a "Completada". Tras
  un 200 se refresca la lista con el estado y la fecha de finalización; ante un rechazo se muestra
  el `mensaje` de la API, incluido el 409 de carrera (que pide recargar). Se prueba con Vitest y
  Testing Library, con el escenario en el nombre de cada prueba.
- **Rationale**: reutiliza el cliente de `web/src/compartido/api` y la feature `sprints`; el
  estado "Completada" no se ofrece como acción directa desde "Nueva" porque el backend la rechaza
  (Q1), así que la pantalla muestra solo la transición que corresponde al estado actual. La lista
  de historias del Sprint sale de la misma pantalla de Sprints que habilita esta fase; si al
  implementarla faltara una consulta para listarlas, se resuelve en esa fase y no cambia el backend
  de HU-11.
- **Alternatives considered**: dos botones siempre visibles (el usuario recibiría rechazos
  evitables); arrastrar y soltar (no se pide).
