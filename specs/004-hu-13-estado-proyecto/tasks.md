---
description: "Task list for HU-13: Fechas y Estado del Proyecto"
---

# Tasks: Fechas y Estado del Proyecto (HU-13)

**Input**: Design documents from `/specs/004-hu-13-estado-proyecto/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: la Constitución (Principio I, NON-NEGOTIABLE) exige TDD estricto — los tasks de
test NO son opcionales: van antes que su implementación con commits separados
`RED:` / `GREEN:` / `REFACTOR:`. Cada COMMIT es una tarea propia con su propio ID, y se marca
`[X]` recién después de correr `git commit`.

**Recordatorio de alcance (compliance)**: esta historia **no agrega tablas, columnas ni
campos**. `internal/repository/migrate.go`, `internal/domain/project.go` y
`internal/domain/sprint.go` NO se tocan (salvo el comentario de T021). El estado se calcula
al consultar. Un solo caso de uso (`ObtenerEstadoProyecto`) y los services dependen de
interfaces de repositorio, nunca de SQLite.

## Phase 1: Setup (Shared Infrastructure)

**No aplica — sin cambios de esquema.** No hay tareas de migración ni de infraestructura:
el `Project` (HU-04) ya tiene `FechaInicio`/`FechaFin` y el `Sprint` (HU-05) ya tiene
`Estado`; los puertos `ProjectRepository` y `SprintRepository` y sus fakes en memoria
(`internal/repository/project_repository_en_memoria.go`,
`internal/repository/sprint_repository_memoria.go`) ya están en `main`. Esta ausencia es
intencional y evidencia el requisito de "sin tablas ni campos nuevos".

---

## Phase 2: Foundational (Blocking Prerequisites)

**⚠️ CRITICAL**: sin la regla de dominio no puede empezar ninguna user story.

- [X] T001 [P] RED: tests de `CalcularEstadoProyecto` con la tabla de decisión completa en `internal/domain/estado_proyecto_test.go` — Sprint `Activo` ⇒ "En curso" (aun con `FechaFin` vencida); ninguno `Activo` y todos `Finalizado` ⇒ "Finalizado"; sin Sprints iniciados ⇒ fechas deciden (`Planificado` si `ahora < FechaInicio`, "En curso" si `FechaInicio ≤ ahora ≤ FechaFin` bordes inclusivos, `Finalizado` si `ahora > FechaFin`); sin Sprints ni `FechaInicio` ⇒ "Planificado" — confirmar que FALLA
- [X] T002 COMMIT `RED: agregar tests de EstadoProyecto (falla)` — stagear **solo** `internal/domain/estado_proyecto_test.go`
- [X] T003 GREEN: implementar `type EstadoProyecto string` con constantes `ProyectoPlanificado`/`ProyectoEnCurso`/`ProyectoFinalizado` y función pura `CalcularEstadoProyecto(p domain.Project, sprints []domain.Sprint, ahora time.Time) domain.EstadoProyecto` en `internal/domain/estado_proyecto.go` (sin `EsValida()`: nadie lo usa — YAGNI, igual que `EstadoSprint` de HU-05) — `go test ./internal/domain/` en verde
- [X] T004 COMMIT `GREEN: implementar EstadoProyecto, test en verde` — stagear `internal/domain/estado_proyecto.go`
- [ ] T005 (solo si hay algo que limpiar) REFACTOR de `CalcularEstadoProyecto` en `internal/domain/estado_proyecto.go` — tests en verde
- [ ] T006 COMMIT `REFACTOR: <qué se limpió>` (solo si se ejecutó T005) — stagear `internal/domain/estado_proyecto.go`

**Checkpoint**: regla de derivación lista y testeada; las user stories pueden arrancar.

---

## Phase 3: User Story 1 - Consultar el estado general de un proyecto (Priority: P1) 🎯 MVP

**Goal**: exponer `GET /projects/{id}/status` que devuelve `{ "estado": "..." }` calculado
on-demand a partir del proyecto, sus Sprints y la fecha actual (solo lectura).

**Independent Test**: sembrar un `Project` con fechas y Sprints en los fakes en memoria,
llamar a `ObtenerEstadoProyecto.Ejecutar(ctx, proyectoID, ahoraFijo)` y verificar la etiqueta;
`GET /projects/{id}/status` devuelve 200/404/400 según corresponda. No depende de US2.

- [X] T007 [P] [US1] RED: test `ObtenerEstadoProyecto.Ejecutar` — 200 con la tabla de decisión del dominio, bordes inclusivos y proyecto inexistente ⇒ `domain.ErrProyectoNoEncontrado` (FR-010) — en `internal/service/obtener_estado_proyecto_test.go`, usando `ProjectRepositoryEnMemoria` + `SprintRepositoryEnMemoria` y un `ahora` fijo — confirmar que FALLA
- [X] T008 [US1] COMMIT `RED: agregar tests de ObtenerEstadoProyecto (falla)` — stagear **solo** `internal/service/obtener_estado_proyecto_test.go`
- [X] T009 [US1] GREEN: implementar `ObtenerEstadoProyecto` en `internal/service/obtener_estado_proyecto.go` — struct con los puertos `repository.ProjectRepository` y `repository.SprintRepository` (interfaces, no SQLite); `Ejecutar(ctx, proyectoID int64, ahora time.Time)` = `ObtenerPorID` → `ListarPorProyecto` → `domain.CalcularEstadoProyecto`; sin llamadas de escritura — `go test ./...` en verde
- [X] T010 [US1] COMMIT `GREEN: implementar ObtenerEstadoProyecto, test en verde` — stagear `internal/service/obtener_estado_proyecto.go`
- [X] T011 [P] [US1] RED: test del handler `GET /projects/{id}/status` en `internal/http/project_handler_test.go` — 200 `{ "estado": "En curso" }`; 404 proyecto inexistente; 400 `id` no entero positivo — confirmar que FALLA
- [X] T012 [US1] COMMIT `RED: agregar tests del handler GET /projects/{id}/status (falla)` — stagear **solo** `internal/http/project_handler_test.go`
- [X] T013 [US1] GREEN: implementar `ProjectHandler.ObtenerEstado` en `internal/http/project_handler.go` y `estadoProyectoResponse{ Estado string `json:"estado"` }` en `internal/http/project_dto.go` (pasa `time.Now()` al service; el 404 ya lo traduce `escribirError`) — tests en verde
- [X] T014 [US1] COMMIT `GREEN: implementar handler GET /projects/{id}/status, test en verde` — stagear `internal/http/project_handler.go` y `internal/http/project_dto.go`
- [X] T015 [US1] Cablear `service.NewObtenerEstadoProyecto(proyectos, sprintRepo)` y la ruta `GET /projects/{id}/status` en `cmd/api/main.go`
- [X] T016 [US1] COMMIT `Cablear GET /projects/{id}/status en cmd/api` — stagear `cmd/api/main.go`

**Checkpoint**: US1 completa y demostrable de forma independiente (MVP).

---

## Phase 4: User Story 2 - Ver el estado actualizado sin recálculo manual (Priority: P2)

**Goal**: garantizar que el estado se recalcula en cada consulta (on-demand) y que la
consulta es de solo lectura.

**Independent Test**: consultar el estado de un proyecto, cambiar el estado de un Sprint en
el fake (o el `ahora`) y volver a consultar: la etiqueta cambia sin ninguna acción de guardado;
y verificar que el service no invoca `Guardar`/`Actualizar`.

> **Nota de TDD**: US2 no introduce código de producción nuevo (es una propiedad del diseño
> on-demand de US1), por lo que no hay fase RED/GREEN: son tests de caracterización/regresión
> que quedan en verde apenas se escriben. Se documenta el desvío explícitamente, como se hizo
> en HU-05.

- [X] T017 [P] [US2] Agregar tests de caracterización en `internal/service/obtener_estado_proyecto_test.go` — (a) mismo proyecto consultado dos veces con `ahora` distinto devuelve estados distintos; (b) tras pasar un Sprint a `Finalizado` en el fake, la re-consulta refleja `Finalizado`; (c) un fake "espía" confirma que `Ejecutar` no llama `Guardar`/`Actualizar` (FR-011, FR-012)
- [X] T018 [US2] COMMIT `Agregar tests de estado on-demand y solo lectura de HU-13` — stagear **solo** `internal/service/obtener_estado_proyecto_test.go`

**Checkpoint**: US1 + US2 — la consulta queda verificada como determinista, on-demand y read-only.

---

## Phase 5: Polish & Cross-Cutting Concerns

> El `.feature` se commitea **antes** de la implementación (a propósito): sus pasos quedan
> `undefined` y, con Godog en modo no estricto, no rompen la suite. La definición de pasos se
> restaura recién al final, cuando el dominio y el servicio ya existen (T027-T029).

- [X] T019 Escenario BDD en `features/estado_proyecto.feature` (tag `@HU-13`) — 7 escenarios de Criterios de Aceptación + 2 de límite (`hoy == fecha_inicio` y `hoy == fecha_fin`, ambos "En curso" según los Edge Cases del spec). Commiteado antes de implementar; pasos `undefined` sin romper la suite
- [X] T020 COMMIT `Agregar escenario BDD de estado del proyecto (HU-13)` — stagear **solo** `features/estado_proyecto.feature`
- [X] T021 Corregir el comentario de `internal/domain/project.go` (hoy sugiere campos de "estado" en la struct) para aclarar que el estado de HU-13 se calcula y NO se persiste
- [X] T022 COMMIT `Aclarar que el estado del proyecto es derivado, no un campo persistido` — stagear `internal/domain/project.go`
- [X] T023 Confirmar `go build ./...`, `go vet ./...` y `go test ./...`
- [X] T024 Ejecutar la validación de `quickstart.md`: tabla de decisión, endpoint 200/404/400 y verificación de solo lectura
- [X] T025 Actualizar `spec.md` (Status: Implementada) y `docs/PRODUCT-BACKLOG.md` si cambió el alcance de HU-13
- [X] T026 COMMIT `Marcar HU-13 como implementada` — stagear `specs/004-hu-13-estado-proyecto/spec.md` (y `docs/PRODUCT-BACKLOG.md` si se tocó)
- [X] T027 **[final] Restaurar los pasos BDD**: copiar `specs/004-hu-13-estado-proyecto/_steps_estado_proyecto.go` a `features/steps_estado_proyecto.go` (renombrar sin el guion bajo) y registrarlos — `InitializeScenario` en `features/steps_creacion.go` debe devolver `*scenarioContext`, y en `features/features_test.go` usar `sc := InitializeScenario(ctx)` + `InitializeScenarioEstado(ctx, sc)`
- [X] T028 Verificar `go build ./...`, `go vet ./...`, `go test ./...` y `go test ./features/ -v` en verde (0 pasos `undefined`)
- [X] T029 COMMIT `Agregar pasos BDD de estado del proyecto (HU-13)` — stagear `features/steps_estado_proyecto.go`, `features/steps_creacion.go` y `features/features_test.go`

---

## Dependencies & Execution Order

- **Phase 2 → Phase 3/4**: bloqueante (la regla de dominio es prerequisito de todo).
- **US1 (P1)**: depende solo de Phase 2. Es el MVP.
- **US2 (P2)**: depende de US1 (reutiliza el service ya implementado); agrega solo tests.
- **Phase 5**: T027-T029 van **al final**, después de toda la implementación (el `.feature` ya quedó commiteado antes; recién ahí se restauran los pasos). T023-T026 son de cierre.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT (nunca test e implementación en el mismo commit). El COMMIT de REFACTOR (T006) es opcional, solo si hubo limpieza.

## Parallel Example

- T001 (test de dominio) es `[P]`: archivo nuevo, sin dependencias.
- T007/T009 (service) y T011/T013 (handler) tocan archivos distintos, pero el handler depende
  del service ya implementado, así que se ejecutan en secuencia por fase.
- Los `[P]` indican "sin dependencia de otro task incompleto", no ejecución literal simultánea.

## Implementation Strategy

**MVP = Phase 2 + Phase 3 (US1)**: con eso ya se demuestra "consultar el estado de un
proyecto" de punta a punta. US2 (Phase 4) agrega la verificación de on-demand/read-only y
Phase 5 cierra con BDD y control de calidad.

## Notes

- Ningún task hace `git push` ni merge: eso se decide por PR (merge commit, nunca squash).
- **Cero cambios de esquema**: si una tarea sintiera la tentación de agregar un campo `Estado`
  o una columna, es una violación del plan — el estado se calcula en `CalcularEstadoProyecto`.
- Un commit por fase TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada
  por la cátedra. Los COMMIT del RED stagean solo el/los archivo(s) de test.
- El instante `ahora` se inyecta como parámetro (no `time.Now()` dentro del dominio) para que
  la regla sea pura y testeable.
