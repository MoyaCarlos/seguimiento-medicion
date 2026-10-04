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
- [X] T006 [P] Implementar `SprintRepositoryFake` (en memoria) en `internal/repository/sprint_repository_fake.go` — lo va a necesitar también HU-13 para su propio TDD

**Checkpoint**: dominio y puerto de Sprint listos — las 3 user stories pueden arrancar.

## Phase 3: User Story 1 - Iniciar un Sprint (Priority: P1) 🎯 MVP

**Goal**: pasar un Sprint de `Pendiente` a `Activo` definiendo Goal y fechas, respetando que no haya otro `Activo` en el mismo proyecto.

**Independent Test**: sembrar un Sprint `Pendiente` directo en `SprintRepositoryFake`, iniciarlo, verificar `Activo`; sembrar un segundo Sprint `Pendiente` en el mismo proyecto y verificar que iniciarlo falla porque ya hay un `Activo`.

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

**Independent Test**: crear un Sprint para un `Project` existente (vía `ProjectRepositoryFake`) → queda `Pendiente`; crear para un `ProyectoID` inexistente → error; crear un segundo `Pendiente` en el mismo proyecto → error.

- [ ] T029 [P] [US3] RED: test `CrearSprint` — caso exitoso en `internal/service/crear_sprint_test.go`
- [ ] T030 [US3] GREEN: implementar servicio `CrearSprint` (usa `ProjectRepository.GetByID`) en `internal/service/crear_sprint.go`
- [ ] T031 [P] [US3] RED: test `CrearSprint` — rechaza `ProyectoID` inexistente (FR-002)
- [ ] T032 [US3] GREEN: agregar esa validación
- [ ] T033 [P] [US3] RED: test `CrearSprint` — rechaza si ya existe otro Sprint `Pendiente` en el proyecto (FR-011)
- [ ] T034 [US3] GREEN: agregar esa validación
- [ ] T035 [P] [US3] RED: test handler `POST /sprints` en `internal/http/sprint_handler_test.go`
- [ ] T036 [US3] GREEN: implementar `SprintHandler.Crear`
- [ ] T037 [US3] Cablear ruta `POST /sprints` en `cmd/api/main.go`

**Checkpoint**: las 3 user stories completas — HU-05 cerrada funcionalmente.

## Phase 6: Polish

- [ ] T038 Escenario BDD en `features/apertura_cierre_sprints.feature` + step definitions (usar el prompt de referencia de `AGENTS.md`, después de este `/speckit.tasks` + `/speckit.analyze`)
- [ ] T039 Confirmar `go build ./...` y `go vet ./...` sin errores
- [ ] T040 Marcar HU-05 como implementada en `docs/PRODUCT-BACKLOG.md` si cambió algo del alcance original

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
- **US3 en pausa**: el contrato `Project`/`ProjectRepository` usa
  `ID string` y el resto del sistema `int64` — se resuelve con quien hace
  HU-04 antes de implementar `CrearSprint`.

## Dependencies & Execution Order

- **Phase 1 → Phase 2**: bloqueante (Setup antes que Foundational).
- **Phase 2 → Fases 3/4/5**: bloqueante (dominio/puerto de Sprint antes que cualquier user story).
- **US1 (P1) y US2 (P1)**: no dependen una de la otra para ser *testeadas* (cada una siembra su propio estado de Sprint directo en el fake), pero **US2 depende de los tasks T017-T019** (campo `SprintID` + métodos de `BacklogRepository`) antes de poder implementar el arrastre.
- **US3 (P2)**: depende de Phase 2 únicamente; no depende de US1/US2. Se ordena después por prioridad de valor, no por dependencia técnica real — se podría adelantar si conviene por reparto de trabajo.
- **Phase 6**: después de que las 3 user stories estén implementadas.

## Parallel Example

Dentro de una misma user story, los tasks marcados `[P]` tocan archivos distintos y pueden hacerse en paralelo (ej. T007 y T009 son tests en el mismo archivo pero casos distintos — en la práctica se escriben uno a la vez siguiendo RED→GREEN, `[P]` acá indica que no hay dependencia de otro task no completado, no que deban ejecutarse literalmente a la vez con Go).

## Implementation Strategy

**MVP = Phase 1 + Phase 2 + Phase 3 (US1)**: con eso ya se puede demostrar "iniciar un Sprint" de punta a punta. US2 y US3 se agregan incrementalmente después, cada una dejando el sistema en un estado funcional y demostrable.
