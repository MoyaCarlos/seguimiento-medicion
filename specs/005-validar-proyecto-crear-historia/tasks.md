---
description: "Task list for Issue #19: validar proyecto existente al crear historia de backlog"
---

# Tasks: Validar proyecto existente al crear historia de backlog (Issue #19)

**Input**: Design documents from `/specs/005-validar-proyecto-crear-historia/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: la Constitución (Principio I, NON-NEGOTIABLE) exige TDD estricto — los tasks de
test NO son opcionales: van antes que su implementación con commits separados
`RED:` / `GREEN:` / `REFACTOR:`. Cada COMMIT es una tarea propia con su propio ID, y se marca
`[X]` recién después de correr `git commit`.

**Recordatorio de alcance (compliance)**: esta historia **no agrega tablas, columnas, campos
ni `FOREIGN KEY`**. `internal/repository/migrate.go`, `internal/domain/backlog_item.go` y
`internal/repository/backlog_repository.go` NO se tocan. **El mapeo HTTP 404 NO se implementa
acá** (FR-007/US2 fuera de alcance): ya lo cubre el traductor compartido `escribirError`
(`internal/http/sprint_handler.go:118`), por lo que no hay tareas de handler para el 404.

**Impactos por cambio de firma del constructor** (descubiertos en plan/research, no estaban en
el pedido original): además de `main.go` y de los tests del service, el cambio de firma de
`NewCrearHistoriaBacklog` obliga a adaptar `internal/http/backlog_handler_test.go:27` y
`features/steps_creacion.go:44`, o el módulo no compila. Se reflejan como tareas propias.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2)
- Include exact file paths in descriptions

## Phase 1: Gate de consistencia y remediación (previo a implementar)

**Objetivo**: cumplir el paso estándar de `AGENTS.md` (`/speckit.analyze` **antes** de
`/speckit.implement`) y dejar aplicadas las correcciones que surgieron del análisis. Resultado:
**0 CRITICAL**. Detalle en `research.md` → "/speckit.analyze (remediación, 2026-10-05)".

- [X] T001 Correr `/speckit.analyze` sobre `specs/005-validar-proyecto-crear-historia/` (spec, plan, tasks, data-model, contracts, quickstart) — resultado **0 CRITICAL**, sin violaciones de la constitución
- [X] T002 (I1) Corregir en este `tasks.md` la descripción del test RED del service: no referenciar `errPersistencia` (vive en el paquete `http`); usar un error local del paquete service (`errFalloPersistencia := errors.New("fallo de persistencia")` en `crear_historia_test.go`)
- [X] T003 (I2) Reformular `FR-007` en `spec.md` como criterio de aceptación relacionado ya cubierto por `escribirError`, y alinear US2 y el último supuesto (ya no hablan de una "tarea de mapeo de errores" pendiente)
- [X] T004 (U1) Aclarar en este `tasks.md` que el paso BDD `Dado … que no existe` solo fija `s.proyectoID`; la ejecución la hace el `Cuando` (`ingresoValido`)
- [X] T005 Registrar el resultado del análisis y la resolución de I1/I2/U1/C1/C2/D1 en `research.md`
- [X] T006 COMMIT `Aplicar remediación del análisis de consistencia (Issue #19)` — stagear `specs/005-validar-proyecto-crear-historia/spec.md`, `specs/005-validar-proyecto-crear-historia/tasks.md` y `specs/005-validar-proyecto-crear-historia/research.md`

**Checkpoint**: spec/plan/tasks consistentes; se puede implementar.

---

## Phase 2: Setup (Shared Infrastructure)

**No aplica — sin cambios de esquema ni de infraestructura.** No hay migraciones, tablas ni
dependencias nuevas. Los puertos `ProjectRepository` y `BacklogRepository` ya existen, igual
que el adaptador SQLite y los fakes. Esta ausencia es intencional (FR-008: sin `FOREIGN KEY`).

---

## Phase 3: Foundational (Blocking Prerequisites)

**No aplica — no hay prerequisito bloqueante.** `repository.ProjectRepository` ya está definido
(`internal/repository/project_repository.go`), `ProjectRepositoryEnMemoria`/`fakeProjectRepository`
ya existen (`internal/service/fakes_test.go` y `internal/repository/project_repository_en_memoria.go`),
y `domain.ErrProyectoNoEncontrado` ya existe (`internal/domain/errors.go`). No hay que crear
nada antes de la user story.

---

## Phase 4: User Story 1 - Rechazar la creación si el proyecto no existe (Priority: P1) 🎯 MVP

**Goal**: que `CrearHistoriaBacklog.Ejecutar` valide la existencia del proyecto
(`proyectos.ObtenerPorID`) después de construir el `BacklogItem` (invariantes primero) y antes
de `repo.Guardar`, propagando el error tal cual y sin persistir la historia.

**Independent Test**: ejecutar el caso de uso con un `ProyectoID` sin proyecto en el
`fakeProjectRepository` y verificar `errors.Is(err, domain.ErrProyectoNoEncontrado)` con
`len(repo.guardados) == 0`; y que con el proyecto precargado la creación sigue funcionando.

- [ ] T007 [P] [US1] RED: en `internal/service/crear_historia_test.go` (a) actualizar los tests existentes (`Exitosa`, `NoPersisteSiEsInvalida`, `ValorNegocio`) para construir el service con `NewCrearHistoriaBacklog(repo, proyectos)` usando `fakeProjectRepository` con el proyecto `1` precargado (`proyectos: map[int64]domain.Project{1: {ID: 1}}`); (b) agregar 3 tests nuevos: proyecto inexistente ⇒ `errors.Is(err, domain.ErrProyectoNoEncontrado)` y `len(repo.guardados) == 0`; error genérico de `ObtenerPorID` propagado y `len(repo.guardados) == 0` (usar un error local del paquete service, p. ej. `errFalloPersistencia := errors.New("fallo de persistencia")` definido en `crear_historia_test.go`, y `fakeProjectRepository{err: errFalloPersistencia}`; NO reusar `errPersistencia`, que vive en `internal/http/project_handler_test.go`); historia inválida (título `"   "`) con `ProyectoID` inexistente ⇒ `domain.ValidationError` con `Campo == "titulo"` — confirmar que FALLA (no compila contra la firma actual de 1 argumento)
- [ ] T008 [US1] COMMIT `RED: agregar tests de validación de proyecto en CrearHistoriaBacklog (falla)` — stagear **solo** `internal/service/crear_historia_test.go`
- [ ] T009 [US1] GREEN: en `internal/service/crear_historia.go` agregar el campo `proyectos repository.ProjectRepository` a `CrearHistoriaBacklog`, actualizar `NewCrearHistoriaBacklog(repo repository.BacklogRepository, proyectos repository.ProjectRepository)` y, en `Ejecutar`, tras `domain.NewBacklogItem` y antes de `s.repo.Guardar`, llamar `s.proyectos.ObtenerPorID(ctx, input.ProyectoID)`; si devuelve error, `return domain.BacklogItem{}, err` (propagación tal cual, cubre `ErrProyectoNoEncontrado` y cualquier otro, FR-006) — `go test ./internal/service/...` en verde
- [ ] T010 [US1] COMMIT `GREEN: validar proyecto existente en CrearHistoriaBacklog, test en verde` — stagear `internal/service/crear_historia.go`
- [ ] T011 [P] [US1] Cablear en `cmd/api/main.go`: mover la construcción de `proyectos := repository.NewSQLiteProjectRepository(db)` **antes** de `service.NewCrearHistoriaBacklog(...)` y pasarle `proyectos` (hoy está en la línea 27, antes de que `proyectos` exista)
- [ ] T012 [US1] COMMIT `Cablear ProjectRepository en NewCrearHistoriaBacklog (cmd/api)` — stagear `cmd/api/main.go`
- [ ] T013 [P] [US1] Adaptar `internal/http/backlog_handler_test.go`: en `nuevoHandlerDePrueba()` pasar también un `*stubProjectRepository` (ya definido en `internal/http/project_handler_test.go`, no crear otro) con el proyecto `1` precargado, para que `TestBacklogHandler_Crear_Exitosa` siga dando `201`; NO agregar test de 404 para el proyecto inexistente (FR-007 fuera de alcance, ya cubierto por `escribirError`) — `go test ./internal/http/...` en verde
- [ ] T014 [US1] COMMIT `Adaptar tests del handler de backlog a la nueva firma del service` — stagear **solo** `internal/http/backlog_handler_test.go`
- [ ] T015 [P] [US1] Adaptar `features/steps_creacion.go`: en `iniciarBD()` pasar `repository.NewSQLiteProjectRepository(s.db)` a `service.NewCrearHistoriaBacklog(...)`; agregar el paso `que intento crear una historia para el proyecto con identificador (\d+), que no existe` (fija `s.proyectoID` al identificador indicado; la ejecución del caso de uso la hace el paso `Cuando …`, `ingresoValido`) y el paso `el sistema responde "proyecto no encontrado"` (`errors.Is(s.err, domain.ErrProyectoNoEncontrado)`) — `go test ./features/...` compila y en verde
- [ ] T016 [US1] COMMIT `Cablear ProjectRepository y agregar pasos BDD de proyecto inexistente (Issue #19)` — stagear `features/steps_creacion.go`
- [ ] T017 [P] [US1] Agregar en `features/crear_historia_backlog.feature` el Escenario `Rechazo de creación por proyecto inexistente` (`Dado que intento crear una historia para el proyecto con identificador 999, que no existe` / `Cuando ingreso un título, una descripción y una prioridad "M" válidos y presiono "Guardar"` / `Entonces el sistema responde "proyecto no encontrado"` / `Y no registra la historia en el Product Backlog`) — `go test ./features/... -run TestFeatures` en verde
- [ ] T018 [US1] COMMIT `Agregar escenario BDD de rechazo por proyecto inexistente (Issue #19)` — stagear **solo** `features/crear_historia_backlog.feature`
- [ ] T019 [US1] (solo si hay algo que limpiar) REFACTOR de `internal/service/crear_historia.go` — tests en verde
- [ ] T020 [US1] COMMIT `REFACTOR: <qué se limpió>` (solo si se ejecutó T019) — stagear `internal/service/crear_historia.go`

**Checkpoint**: US1 completa y demostrable de forma independiente (MVP): el caso de uso rechaza
historias huérfanas y el camino feliz sigue funcionando.

---

## Phase 5: User Story 2 - Responder 404 en POST /backlog (Priority: P2)

**Fuera de alcance — sin tareas de implementación.**

El mapeo de `domain.ErrProyectoNoEncontrado` a `404` ya existe en el traductor compartido
`escribirError` (`internal/http/sprint_handler.go:118`), que `BacklogHandler.Crear` ya invoca
(`internal/http/backlog_handler.go:36`). Por lo tanto, apenas US1 devuelve el error, el
endpoint responde `404` sin tocar el handler. No se agregan tests de handler en esta historia
(FR-007) porque el comportamiento se verificará de punta a punta en `quickstart.md` (T022). Si
en T022 el `404` no aparece, **no** se resuelve acá: se reporta como hallazgo al paso de mapeo
de errores.

**Independent Test**: `POST /backlog` con `proyecto_id` inexistente responde `404` con
`{"mensaje":"proyecto no encontrado"}` (verificado manualmente en T022).

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T021 Confirmar `go build ./...`, `go vet ./...` y `go test ./...` en verde (incluye `internal/service`, `internal/http` y `features`)
- [ ] T022 Ejecutar la validación de `quickstart.md`: tests de service (4 casos), tests de handler (201/400) y verificación manual del `404` con `curl` contra `POST /backlog` con `proyecto_id: 999` (sin filas nuevas en `backlog_items`)
- [ ] T023 Actualizar `specs/005-validar-proyecto-crear-historia/spec.md` (Status) y `docs/PRODUCT-BACKLOG.md` si el issue #19 lo amerita
- [ ] T024 COMMIT `Marcar Issue #19 como implementado` — stagear `specs/005-validar-proyecto-crear-historia/spec.md` (y `docs/PRODUCT-BACKLOG.md` si se tocó)

---

## Dependencies & Execution Order

- **Phase 1 (gate + remediación)** → antes de cualquier implementación (T001–T005 ya ejecutados; queda T006, el commit).
- **Phase 4 (US1)**: es el MVP; no depende de otra user story.
  - Dentro del comportamiento: RED (T007) → COMMIT (T008) → GREEN (T009) → COMMIT (T010).
  - T011–T018 son cableado/adaptación mecánica y BDD: dependen de T009 (la firma nueva ya debe existir).
- **Phase 5 (US2)**: sin tareas; se satisface por infraestructura existente al terminar US1.
- **Phase 6**: al final, con toda la implementación hecha.

## Parallel Example

- T007 (RED, test de service), T013 (adaptar test de handler), T015 (adaptar steps) y T017
  (feature BDD) tocan archivos distintos; pueden prepararse en paralelo, pero **T011–T018
  dependen de T009** (la firma nueva del constructor), así que se ejecutan después del GREEN.
- Los `[P]` indican "sin dependencia de otro task incompleto", no ejecución literal simultánea.

## Implementation Strategy

**MVP = Phase 4 (US1)**: valida la existencia del proyecto al crear la historia y deja el caso
de uso testeado de forma independiente. US2 no agrega código (ya cubierto por `escribirError`);
el cierre (T021–T024) verifica build, suite completa y quickstart.

## Notes

- Ningún task hace `git push` ni merge: eso se decide por PR (merge commit, nunca squash).
- **Cero cambios de esquema**: si una tarea sintiera la tentación de agregar `FOREIGN KEY`,
  es una violación del plan (FR-008, paso posterior acordado).
- El orden de validación se mantiene: `NewBacklogItem` (invariantes) → `ObtenerPorID` →
  `Guardar`. No reordenar.
- Un commit por fase TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada
  por la cátedra. El COMMIT del RED stagea solo el/los archivo(s) de test.
- El GREEN de T009 se valida con `go test ./internal/service/...` porque, hasta que T011–T015
  adapten los llamadores, el módulo completo no compila (consecuencia normal del cambio de
  firma del constructor en Go). T021 cierra con `go build ./...` y `go test ./...`.
- La remediación del análisis (I1/I2/U1) quedó aplicada en `spec.md`, `research.md` y este
  `tasks.md`; su commit es T006.
