# Research: Movimiento de Historias al Sprint Backlog (HU-06)

No quedaron `NEEDS CLARIFICATION` en el Technical Context. Este documento registra las
decisiones técnicas (las acordadas por el equipo y los ajustes que surgieron al revisar el
código) con sus alternativas.

## R1. La regla de la asignación es una función pura del dominio

- **Decision**: `domain.ValidarAsignacion(h BacklogItem, s Sprint) error` en
  `internal/domain/asignacion_sprint.go`. Evalúa en este orden y devuelve el primer error:
  1. `h.ProyectoID != s.ProyectoID` → `ErrSprintDeOtroProyecto`
  2. `s.Estado == SprintFinalizado` → `ErrSprintFinalizado`
  3. `h.Estado == EstadoCompletada` → `ErrHistoriaNoPendiente`
  4. `h.SprintID != nil` → `ErrHistoriaYaAsignada`

  Los cuatro son `ValidationError` (400 vía `escribirError`); los mensajes están en
  [data-model.md](./data-model.md).
- **Rationale**: el orden de precedencia es una regla de negocio (FR-009) y se prueba sin base,
  con una tabla de casos que combine varias violaciones a la vez. Las existencias (404) quedan
  fuera de la función porque necesitan el repositorio: las resuelve el servicio antes.
- **Alternatives considered**: un método `(*BacklogItem) AsignarA(s Sprint)` que además mute
  `SprintID`. Se descartó porque la escritura real es condicional en la base (R2) y el estado en
  memoria quedaría adelantado a lo persistido.

## R2. Escritura condicional en una sola sentencia

- **Decision**: método `AsignarASprint(ctx, historiaID, sprintID int64) error` en
  `SQLiteBacklogRepository`:
  `UPDATE historias SET sprint_id = ? WHERE id = ? AND sprint_id IS NULL`.
  - Error de la sentencia → `errorDeEscritura` (una violación de la FK compuesta llega como
    `ErrIntegridadReferencial`, 409, red de seguridad de `MODELO-DE-DATOS.md` §5).
  - `RowsAffected() == 0` → `domain.ErrHistoriaYaAsignada`.
- **Rationale**: si dos asignaciones de la misma historia llegan juntas, las dos pasan
  `ValidarAsignacion`, pero solo una sentencia encuentra `sprint_id IS NULL`; la otra recibe el
  mismo rechazo que el caso normal. No hace falta transacción explícita ni bloqueo.
- **Alternatives considered**: reutilizar `Actualizar` (no escribe `sprint_id` y no detecta la
  carrera); transacción con `SELECT … FOR UPDATE` (SQLite no lo tiene y es más código).
- **Riesgo aceptado**: si el Sprint se cierra entre la validación y el `UPDATE`, la historia
  queda en un Sprint recién `Finalizado`. Es una ventana mínima en una app mono-usuario; cubrirla
  requiere un `EXISTS` sobre `sprints` en el mismo `UPDATE` y se deja para un fix si aparece.

## R3. Puerto propio para la asignación (ajuste sobre lo acordado)

- **Decision**: puerto nuevo en `internal/repository/backlog_repository.go`:

  ```go
  type AsignacionHistoriaRepository interface {
      ObtenerPorID(ctx context.Context, id int64) (domain.BacklogItem, error)
      AsignarASprint(ctx context.Context, historiaID, sprintID int64) error
  }
  ```

  Lo implementa el mismo `SQLiteBacklogRepository`. **No** se agrega `AsignarASprint` a
  `HistoriaRepository`.
- **Rationale**: el repo ya separa los puertos por caso de uso (ISP): `ProductBacklogRepository`,
  `HistoriaRepository` (edición) y `HistoriasDeSprintRepository` (cierre, con el comentario
  "Separado de BacklogRepository (ISP)"). Ampliar `HistoriaRepository` obligaría a
  `EditarHistoria`, su fake (`fakeHistoriaRepository`) y los stubs de los handlers a conocer una
  operación que no usan. El cambio es solo de nombre de interfaz; el comportamiento acordado
  (un `UPDATE` condicional) no cambia.
- **Alternatives considered**: ampliar `HistoriaRepository` (lo acordado originalmente);
  descartado por lo anterior.

## R4. Orden del Sprint Backlog con un `ORDER BY` compartido (ajuste en las pruebas)

- **Decision**: extraer a una constante en `sqlite_backlog.go`
  `ordenBacklog = "ORDER BY CASE prioridad WHEN 'M' THEN 0 WHEN 'S' THEN 1 WHEN 'C' THEN 2 WHEN 'W' THEN 3 END, id"`
  y usarla en `ListarSinSprint` y `ListarPorSprint`. El RED es una prueba de
  `ListarPorSprint` contra SQLite con historias de prioridades mezcladas.
- **Ajuste**: lo acordado pedía también un RED "en el fake". No existe un fake de repositorio de
  historias en `internal/repository`: los únicos que implementan `ListarPorSprint` son dobles
  locales de prueba (`historiasEnMemoria` en `cerrar_sprint_test.go`, `sinHistorias` en
  `sprint_handler_test.go`). El orden es responsabilidad del adaptador (igual que en
  `ConsultarProductBacklog`, que devuelve "el orden que entrega el repositorio"), así que un RED
  sobre un doble de prueba no probaría código de producción. El RED del orden queda solo en
  SQLite; el fake del servicio devuelve las historias tal como se siembran.
- **Rationale**: DRY (una sola definición del orden MoSCoW, `MODELO-DE-DATOS.md` §4.4) y FR-011.
  El cierre de Sprint no depende del orden: `TestSQLiteBacklogRepository_HistoriasDeSprint` usa
  historias de la misma prioridad, así que sigue en verde.
- **Alternatives considered**: ordenar en el servicio con `slices.SortFunc` (duplica la regla
  que ya vive en SQL); una consulta nueva `ListarSprintBacklog` (dos consultas casi iguales).

## R5. Rutas HTTP y precedencia de errores

- **Decision**:
  - `PUT /backlog/{id}/sprint` con cuerpo `{"sprint_id": N}` en `BacklogHandler.AsignarSprint`.
    Siguiendo el patrón de `Editar`, el handler verifica primero que la historia exista
    (`ObtenerHistoria`, 404) y recién después decodifica el cuerpo: así una historia inexistente
    da 404 aunque el JSON sea inválido o falte `sprint_id`. `sprint_id` ausente o `null` → 400
    `{campo: "sprint_id", mensaje: "el campo es obligatorio"}`; un valor numérico se pasa al
    servicio, que da 404 si el Sprint no existe. Respuesta 200 con la historia
    (mismo `historiaResponse` que el Product Backlog), con `sprint_id` cargado.
  - **`historiaResponse` suma `sprint_id` (`*int64`, `null` sin Sprint).** Hoy el DTO no lo
    expone, y sin él ni la respuesta de la asignación ni el Sprint Backlog dicen a qué Sprint
    pertenece la historia. Es un campo nuevo en la respuesta (no cambia ninguno existente): en
    `GET /projects/{id}/backlog` siempre vale `null` y los clientes actuales lo ignoran.
  - `GET /sprints/{id}/backlog` en `SprintHandler.ConsultarBacklog`: 404 si el Sprint no existe,
    200 con `[]` si no tiene historias.
  - Todos los errores pasan por `escribirError`; no hace falta ningún `case` nuevo porque los
    errores de la asignación son `ValidationError`.
- **Rationale**: la historia en la ruta hace que "historia inexistente" domine de forma natural
  (FR-008); `PUT` porque asignar es idempotente en su efecto sobre el recurso (el segundo intento
  se rechaza sin cambiar nada). Go 1.22+ distingue `PUT /backlog/{id}` de
  `PUT /backlog/{id}/sprint` sin conflicto.
- **Alternatives considered**: `POST /sprints/{id}/historias` con `historia_id` en el cuerpo
  (haría que el 404 de Sprint saliera antes que el de historia, al revés de la spec);
  `PATCH /backlog/{id}` genérico (mezcla con la edición de HU-15).

## R6. Escenario BDD

- **Decision**: `features/asignar_historias_sprint.feature` con `@HU-06`, un escenario por
  criterio de aceptación de US1 (9) y US2 (3), sobre SQLite en memoria como el resto de los
  pasos (`steps_backlog.go`). Los Sprints se preparan con los servicios de HU-05
  (crear, iniciar, cerrar) para no sembrar estados a mano; la historia "Completada" se siembra
  con SQL directo, como ya hace `sqlite_sprint_test.go`, porque HU-11 todavía no existe.
- **Rationale**: Principio III; las 4 categorías (normal, alternativo, límite, error) están
  cubiertas por los escenarios de la spec. US3 no lleva escenario Godog: se prueba con Vitest.

## R7. Pantalla (US3)

- **Decision**: fase propia, bloqueada hasta que estén en `main` el PR #149 y el #98. En cada
  historia del Product Backlog, una acción "Asignar a Sprint" con un selector que solo ofrece los
  Sprints `Pendiente` y `Activo` del proyecto (de `GET /projects/{id}/sprints`). El botón de
  confirmar queda deshabilitado mientras la petición está en curso (como FR-016 de HU-01, evita
  envíos duplicados). Tras un 200 se refrescan el Product Backlog y el Sprint Backlog afectado;
  ante un error se muestra el `mensaje` de la respuesta sin reemplazarlo (FR-015). Cada Sprint
  muestra su Sprint Backlog con `GET /sprints/{id}/backlog`.
- **Rationale**: reutiliza el cliente de `web/src/compartido/api` y las features existentes; el
  arrastrar y soltar es opcional (spec, Assumptions).
