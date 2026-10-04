---
description: "Task list for HU-05: Apertura y Cierre de Sprints"
---

# Tasks: Apertura y Cierre de Sprints

**Input**: Design documents from `/specs/003-apertura-cierre-sprints/`
**Prerequisites**: plan.md, spec.md, data-model.md, research.md, contracts/, quickstart.md

**Tests**: la Constitución (Principio I, NON-NEGOTIABLE) exige TDD estricto
— los tasks de test NO son opcionales en este proyecto, van antes que su
implementación correspondiente en cada fase (RED antes que GREEN).

## Phase 1: Setup

- [X] T001 Agregar tabla `sprints` (id, proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) y columna `sprint_id` nullable en `backlog_items` a la migración en `internal/repository/migrate.go`

## Phase 2: Foundational (bloqueante para todas las user stories)

- [X] T002 [P] Crear tipo `EstadoSprint` (constantes `Pendiente`/`Activo`/`Finalizado` + método `EsValida()`) en `internal/domain/estado_sprint.go`
- [X] T003 [P] RED: test de `NewSprint` — rechaza `ProyectoID <= 0` — en `internal/domain/sprint_test.go`
- [X] T004 GREEN: implementar struct `Sprint{ID, ProyectoID, SprintGoal, FechaInicio, FechaFin, Estado}` + factory `NewSprint` en `internal/domain/sprint.go`
- [X] T005 [P] Definir puerto `SprintRepository` (`Crear`, `ObtenerPorID`, `ListarPorProyecto`) en `internal/repository/sprint_repository.go`
- [X] T006 [P] Implementar `SprintRepositoryEnMemoria` en `internal/repository/sprint_repository_memoria.go` — lo va a necesitar también HU-13 para su propio TDD

**Checkpoint**: dominio y puerto de Sprint listos — las 3 user stories pueden arrancar.

## Phase 3: User Story 1 - Iniciar un Sprint (Priority: P1) 🎯 MVP

**Goal**: pasar un Sprint de `Pendiente` a `Activo` definiendo Goal y fechas, respetando que no haya otro `Activo` en el mismo proyecto.

**Independent Test**: sembrar un Sprint `Pendiente` directo en `SprintRepositoryEnMemoria`, iniciarlo, verificar `Activo`; sembrar un segundo Sprint `Pendiente` en el mismo proyecto y verificar que iniciarlo falla porque ya hay un `Activo`.

- [X] T007 [P] [US1] RED: test `IniciarSprint` — caso exitoso (`Pendiente`→`Activo`) en `internal/service/iniciar_sprint_test.go`
- [X] T008 [US1] GREEN: implementar servicio `IniciarSprint` en `internal/service/iniciar_sprint.go`
- [X] T009 [P] [US1] RED: test `IniciarSprint` — rechaza si ya existe otro Sprint `Activo` en el proyecto (FR-004)
- [X] T010 [US1] GREEN: agregar la validación (consulta `SprintRepository.ListarPorProyecto`)
- [X] T011 [P] [US1] RED: test `IniciarSprint` — rechaza si `FechaFin` no es posterior a `FechaInicio` (FR-005)
- [X] T012 [US1] GREEN: agregar la validación de fechas
- [X] T013 [P] [US1] RED: test handler `POST /sprints/{id}/iniciar` — 200 caso éxito en `internal/http/sprint_handler_test.go`
- [X] T014 [P] [US1] RED: test handler — 400 en los 3 casos de validación fallida
- [X] T015 [US1] GREEN: implementar `SprintHandler.Iniciar` + DTOs en `internal/http/sprint_handler.go`, `internal/http/sprint_dto.go`
- [X] T016 [US1] Cablear ruta `POST /sprints/{id}/iniciar` en `cmd/api/main.go`

**Checkpoint**: US1 completa y demostrable de forma independiente.

## Phase 4: User Story 2 - Cerrar un Sprint con arrastre (Priority: P1)

**Goal**: cerrar un Sprint `Activo`, pasarlo a `Finalizado`, y devolver al Product Backlog las historias no completadas que tenía asignadas.

**Independent Test**: sembrar un Sprint `Activo` con 2 historias asignadas (una completada, una no) directo en los fakes de `SprintRepository`/`BacklogRepository`; cerrar el Sprint; verificar que la no completada queda con `SprintID == nil` y la completada conserva el suyo.

- [X] T017 [P] [US2] Agregar campo `SprintID *int64` a `domain.BacklogItem` en `internal/domain/backlog_item.go` (alcance mínimo para esta historia — la asignación completa es de HU-06)
- [X] T018 [P] [US2] Agregar métodos `ListarPorSprint` y `QuitarDeSprint` al puerto `BacklogRepository` en `internal/repository/backlog_repository.go`
- [X] T019 [US2] Implementar esos métodos en `SQLiteBacklogRepository` (`internal/repository/sqlite_backlog.go`)
- [X] T020 [P] [US2] RED: test `CerrarSprint` — caso exitoso (`Activo`→`Finalizado`) en `internal/service/cerrar_sprint_test.go`
- [X] T021 [US2] GREEN: implementar servicio `CerrarSprint` (cambio de estado) en `internal/service/cerrar_sprint.go`
- [X] T022 [P] [US2] RED: test `CerrarSprint` — rechaza si el Sprint no está `Activo` (FR-007)
- [X] T023 [US2] GREEN: agregar esa validación
- [X] T024 [P] [US2] RED: test `CerrarSprint` — historias no completadas quedan con `SprintID nil`, las completadas conservan el suyo (FR-008, FR-009)
- [X] T025 [US2] GREEN: implementar el arrastre usando `BacklogRepository.ListarPorSprint`/`QuitarDeSprint`
- [X] T026 [P] [US2] RED: test handler `POST /sprints/{id}/cerrar` en `internal/http/sprint_handler_test.go`
- [X] T027 [US2] GREEN: implementar `SprintHandler.Cerrar`
- [X] T028 [US2] Cablear ruta `POST /sprints/{id}/cerrar` en `cmd/api/main.go`

**Checkpoint**: US1 + US2 completas — el ciclo de vida completo de un Sprint ya es demostrable.

## Phase 5: User Story 3 - Crear un Sprint (Priority: P2)

**Goal**: crear un Sprint nuevo en estado `Pendiente` para un Project existente, respetando que no haya ya otro `Pendiente` en ese proyecto.

**Independent Test**: crear un Sprint para un `Project` guardado en `ProjectRepositoryEnMemoria` → queda `Pendiente`; crear para un `ProyectoID` inexistente → `ErrProyectoNoEncontrado`; crear un segundo `Pendiente` en el mismo proyecto → error de validación.

> Fases 5 a 7 renumeradas el 2026-10-04 (antes T029-T040, ninguna empezada) para agregar las tareas COMMIT de la plantilla vigente (PR #15) y lo que salió de `/speckit.clarify` (sesión 2026-10-04).

- [X] T029 [US3] RED: test `CrearSprint` — caso exitoso (`Pendiente`, persistido) en `internal/service/crear_sprint_test.go`, con `ProjectRepositoryEnMemoria` + `SprintRepositoryEnMemoria`
- [X] T030 [US3] COMMIT `RED: agregar test de CrearSprint exitoso (falla)`
- [X] T031 [US3] GREEN: implementar servicio `CrearSprint` (`NewSprint` + `SprintRepository.Guardar`) en `internal/service/crear_sprint.go`
- [X] T032 [US3] COMMIT `GREEN: implementar CrearSprint, test en verde`
- [X] T033 [US3] RED: test `CrearSprint` — rechaza `ProyectoID` inexistente con `domain.ErrProyectoNoEncontrado` y no guarda nada (FR-002)
- [X] T034 [US3] COMMIT `RED: agregar test de CrearSprint con proyecto inexistente (falla)`
- [X] T035 [US3] GREEN: verificar el proyecto con `ProjectRepository.ObtenerPorID` (puerto de HU-04, sin puertos nuevos)
- [X] T036 [US3] COMMIT `GREEN: rechazar proyecto inexistente en CrearSprint, test en verde`
- [X] T037 [US3] RED: test `CrearSprint` — rechaza si ya existe otro Sprint `Pendiente` en el proyecto (FR-011, US3 escenario 3)
- [X] T038 [US3] COMMIT `RED: agregar test de CrearSprint con otro Pendiente (falla)`
- [X] T039 [US3] GREEN: agregar esa validación (consulta `SprintRepository.ListarPorProyecto`)
- [X] T040 [US3] COMMIT `GREEN: rechazar segundo Sprint Pendiente, test en verde`
- [X] T041 [US3] RED: test handler `POST /sprints` en `internal/http/sprint_handler_test.go` — 201 éxito; 404 proyecto inexistente; 400 `proyecto_id <= 0`, otro `Pendiente` y JSON inválido
- [X] T042 [US3] COMMIT `RED: agregar tests del handler POST /sprints (falla)`
- [X] T043 [US3] GREEN: implementar `SprintHandler.Crear` + DTOs en `internal/http/sprint_handler.go`, `internal/http/sprint_dto.go`
- [X] T044 [US3] COMMIT `GREEN: implementar handler POST /sprints, test en verde`
- [X] T045 [US3] Cablear ruta `POST /sprints` en `cmd/api/main.go` (reusar `proyectos` de HU-04)
- [X] T046 [US3] COMMIT `Cablear POST /sprints en cmd/api`

**Checkpoint**: las 3 user stories completas — HU-05 cerrada funcionalmente.

## Phase 6: 404 antes que 400 al iniciar (clarify 2026-10-04)

**Goal**: `POST /sprints/{id}/iniciar` sobre un Sprint inexistente responde 404 aunque el cuerpo sea inválido (fecha mal formada o JSON malformado), mismo criterio que HU-04. Hoy el handler valida el cuerpo antes de buscar el Sprint.

- [X] T047 RED: tests handler — `POST /sprints/999/iniciar` con fecha mal formada → 404, y con JSON malformado → 404
- [X] T048 COMMIT `RED: agregar tests de 404 antes que 400 en iniciar Sprint (falla)`
- [X] T049 GREEN: `NewSprintHandler` recibe el puerto `repository.SprintRepository`; `SprintHandler.Iniciar` llama a `ObtenerPorID` antes de decodificar el cuerpo. Actualizar `cmd/api/main.go` y el helper de `sprint_handler_test.go`; T013/T014 siguen en verde
- [X] T050 COMMIT `GREEN: responder 404 antes que 400 en iniciar Sprint, test en verde`

## Phase 7: BDD y cierre

- [X] T051 Escenarios BDD de US1/US2 en `features/apertura_cierre_sprints.feature` + pasos en `features/steps_sprints.go` (5 de 7 escenarios definidos)
- [X] T052 Agregar el escenario de FR-011 (US3 escenario 3) a `features/apertura_cierre_sprints.feature`
- [X] T053 Definir los pasos de los 3 escenarios de creación. El paso compartido "existe un proyecto con identificador 1" (HU-01, `steps_creacion.go`) solo asigna un número: para HU-05 tiene que existir un proyecto real en el repositorio que usan los pasos de Sprint, sin romper los escenarios de HU-01
- [X] T054 COMMIT `Agregar pasos BDD de creación de Sprint (HU-05)`
- [X] T055 Confirmar `go build ./...`, `go vet ./...`, `gofmt -l .` (vacío), `go test ./...` y Godog sin escenarios `undefined`
- [X] T056 Quitar "NO TERMINADA, NO MERGEAR" del Status de `spec.md` y marcar HU-05 en `docs/PRODUCT-BACKLOG.md` si cambió el alcance
- [X] T057 COMMIT `Marcar HU-05 como implementada`

## Notas de implementación (desvíos respecto al plan)

- **T002**: `EstadoSprint` sin método `EsValida()` — ningún código lo
  necesita (los estados solo los asigna el dominio). YAGNI.
- **T005**: el puerto usa `Guardar` (mismo verbo que `BacklogRepository`)
  y suma `Actualizar`, que faltaba en el plan para persistir transiciones.
- **T018**: `ListarPorSprint`/`QuitarDeSprint` van en un puerto aparte
  (`HistoriasDeSprintRepository`) en vez de agrandar `BacklogRepository`,
  para no obligar a los fakes de HU-01 a implementar métodos que no usan.
- **Reglas de transición en el dominio**: `Sprint.Iniciar`/`Sprint.Cerrar`
  validan estado, Goal y fechas (T011/T012 quedaron en el dominio, no en
  el service); los services solo manejan reglas entre Sprints/historias.
- **Recortes de dominio para US2**: `BacklogItem.SprintID` (asignación
  completa → HU-06) y `EstadoCompletada` (transición → HU-11).
- **E1/F1 de `/speckit.analyze` resueltos**: "Sprint no encontrado" → 404
  con tests en service y handler; `ListarPorProyecto` con test directo.
- **REFACTOR**: `BacklogHandler` (HU-01) pasó a usar el mismo mapeo de
  errores que `SprintHandler`.
- **404 antes que 400 al iniciar (2026-10-04, decisión de Carlos)**: el
  handler consulta el puerto `SprintRepository` directamente, sin un service
  `ObtenerSprint` aparte. Difiere de `ProjectHandler.Editar` (HU-04), que usa
  el service `ObtenerProyecto`. Sigue dependiendo de una interfaz, no de SQLite.
- **T041**: el handler responde 400 por `proyecto_id <= 0` antes de buscar
  el proyecto (es un dato mal formado, no un recurso inexistente).
- **T037**: se sumó un test de guarda (ya en verde en el RED) que permite
  preparar un Sprint Pendiente mientras otro está Activo.
- **T053**: el paso compartido de HU-01 "existe un proyecto con
  identificador N" ahora crea el proyecto real (repositorio de HU-04) en vez
  de solo guardar el número; HU-01 sigue en verde.
- **T052 + T053** quedaron en un solo commit (T054): el escenario y sus pasos.
- **T056**: `PRODUCT-BACKLOG.md` sin cambios — el alcance no cambió y el
  estado de las historias se sigue en el tablero de GitHub Projects.
- **US3 destrabada (2026-10-04)**: HU-04 entró a `main` con IDs `int64`;
  `CrearSprint` usa su `ProjectRepository.ObtenerPorID`.

## Dependencies & Execution Order

- **Phase 1 → Phase 2**: bloqueante (Setup antes que Foundational).
- **Phase 2 → Fases 3/4/5**: bloqueante (dominio/puerto de Sprint antes que cualquier user story).
- **US1 (P1) y US2 (P1)**: no dependen una de la otra para ser *testeadas* (cada una siembra su propio estado de Sprint directo en el fake), pero **US2 depende de los tasks T017-T019** (campo `SprintID` + métodos de `BacklogRepository`) antes de poder implementar el arrastre.
- **US3 (P2)**: depende de Phase 2 únicamente; no depende de US1/US2. Se ordena después por prioridad de valor, no por dependencia técnica real — se podría adelantar si conviene por reparto de trabajo.
- **Phase 6** (404 antes que 400): solo depende de US1; puede ir antes o después de US3.
- **Phase 7**: T052-T054 después de US3 (los pasos usan `CrearSprint`); T055-T057 al final de todo.

## Parallel Example

Dentro de una misma user story, los tasks marcados `[P]` tocan archivos distintos y pueden hacerse en paralelo (ej. T007 y T009 son tests en el mismo archivo pero casos distintos — en la práctica se escriben uno a la vez siguiendo RED→GREEN, `[P]` acá indica que no hay dependencia de otro task no completado, no que deban ejecutarse literalmente a la vez con Go).

## Implementation Strategy

**MVP = Phase 1 + Phase 2 + Phase 3 (US1)**: con eso ya se puede demostrar "iniciar un Sprint" de punta a punta. US2 y US3 se agregan incrementalmente después, cada una dejando el sistema en un estado funcional y demostrable.
