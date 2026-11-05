---

description: "Tareas de HU-12: Consulta de Sprints Anteriores (TDD con commits por fase)"
---

# Tasks: Consulta de Sprints Anteriores (HU-12)

**Input**: Design documents from `/specs/033-hu-12-consulta-sprints-anteriores/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica es RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT).

**Compilación en RED**: todos los commits compilan. En cada RED se agregan, junto al test, solo las
firmas nuevas que el test necesita con un cuerpo vacío (devolver valor cero / `nil`), de modo que
`go build ./...` pase y el test falle por aserción, no por símbolo indefinido.

**Tests en verde**: `go test ./features/` falla de T002 a T039 porque los pasos de
`historial_sprints.feature` todavía no están definidos (modo estricto). Por eso las tareas GREEN
comprueban `go test ./internal/... ./cmd/...`, y `go test ./...` se exige desde T041.

**Organization**: por user story de la spec. US1 (lista) y US2 (detalle) comparten la fase
fundacional. US3 (proyecto y Sprint vacíos) son casos límite de US1 y US2: se prueban dentro de sus
tests y se verifican en su propia fase. US4 es el frontend (sub-issue #88), que no frena el cierre
de la historia.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3 / US4

---

## Phase 1: Setup

**Purpose**: dejar los escenarios BDD escritos antes de implementar (flujo de AGENTS.md, sub-issue #85).

- [X] T001 Escribir features/historial_sprints.feature en Gherkin español: tag `@HU-12`, Característica con el Como/Quiero/Para de la historia, Antecedentes si hay una precondición común y un Escenario por cada criterio de aceptación de US1, US2 y US3 de la spec. Comentario inicial con la cobertura por categoría (normal, alternativo, límite, error), como en features/calculo_metricas.feature. US4 se cubre con pruebas de componente
- [X] T002 COMMIT `RED: agregar escenarios BDD del historial de Sprints (fallan)`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: reglas puras y puerto que usan US1 y US2.

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

- [X] T003 RED: en internal/domain/historial_sprint_test.go, tabla de `SprintsFinalizados(sprints []Sprint) []Sprint`: devuelve solo los `Finalizado`, excluye `Activo` y `Pendiente`, ordena por `Numero` ascendente aunque lleguen desordenados, y con entrada vacía o sin finalizados devuelve un slice vacío "nunca `nil`"; y de `ValidarSprintFinalizado(sprint Sprint) error`: `nil` con `Finalizado`, `ErrSprintNoFinalizado` con `Activo` y con `Pendiente`. Agregar en internal/domain/errors.go `ErrSprintNoFinalizado = ValidationError{Campo: "estado", Mensaje: "el Sprint todavía no está finalizado"}` (research.md R4) y en internal/domain/historial_sprint.go las firmas con cuerpo vacío — confirmar que FALLA
- [X] T004 COMMIT `RED: agregar tests del filtro de Sprints finalizados (fallan)`
- [X] T005 GREEN: implementar `SprintsFinalizados` y `ValidarSprintFinalizado` en internal/domain/historial_sprint.go — `go test ./internal/... ./cmd/...` en verde
- [X] T006 COMMIT `GREEN: implementar el filtro de Sprints finalizados, tests en verde`
- [X] T007 Agregar en internal/domain/historial_sprint.go los tipos de data-model.md sin lógica: `HistoriaComprometida{ID int64; Titulo string; EstimacionSP *int; Completada bool; Minutos int; Horas float64}` (`EstimacionSP` en `nil` = "sin estimar"), `DefectoDeSprint{DefectoDetalle; DetectadoEnSprint bool; ResueltoEnSprint bool}`, `SprintHistorial{Sprint Sprint; Metricas Metricas}` y `DetalleSprint{SprintHistorial; Historias []HistoriaComprometida; Defectos []DefectoDeSprint}`. Definir el puerto `HistorialSprintRepository` con `HistoriasComprometidas(ctx, sprintID int64) ([]domain.HistoriaComprometida, error)`, `MinutosPorHistoria(ctx, sprintID int64) (map[int64]int, error)` y `DefectosDeSprint(ctx, sprintID int64) ([]domain.DefectoDetalle, error)` en internal/repository/historial_sprint_repository.go
- [X] T008 COMMIT `Agregar los tipos del historial de Sprints y su puerto`

**Checkpoint**: las historias pueden empezar.

---

## Phase 3: User Story 1 - Ver la lista de Sprints cerrados (Priority: P1) 🎯 MVP

**Goal**: `GET /projects/{id}/historial` devuelve los Sprints finalizados con su Sprint Goal, fechas y métricas de HU-03.

**Independent Test**: con un proyecto con dos Sprints finalizados, uno `Activo` y uno `Pendiente`, el historial devuelve solo los dos finalizados, ordenados por número, con las mismas métricas que `GET /projects/{id}/metricas`.

- [X] T009 [US1] RED: en internal/service/consultar_historial_sprints_test.go, con fakes (los nuevos, en internal/service/fakes_test.go; dependencias angostas declaradas en el service: `calculadorMetricas` con la firma de `CalcularMetricas.Ejecutar(ctx, proyectoID) (domain.MetricasProyecto, error)` y `listadorSprints` ya existente): une cada Sprint finalizado con las métricas de su `SprintID` (research.md R1 y R2); excluye `Activo` y `Pendiente`; ordena por `Numero`; proyecto inexistente → `ErrProyectoNoEncontrado` (lo devuelve el cálculo) sin listar Sprints; proyecto sin Sprints finalizados → slice vacío no `nil` (US3.1); un Sprint finalizado sin métricas → error interno (data-model.md, inconsistencia entre fuentes); un error del cálculo o del listado se propaga. Stub en internal/service/consultar_historial_sprints.go — confirmar que FALLA
- [X] T010 [US1] COMMIT `RED: agregar tests del caso de uso ConsultarHistorialSprints (fallan)`
- [ ] T011 [US1] GREEN: implementar `ConsultarHistorialSprints.Ejecutar(ctx, proyectoID) ([]domain.SprintHistorial, error)` en internal/service/consultar_historial_sprints.go — `go test ./internal/... ./cmd/...` en verde
- [ ] T012 [US1] COMMIT `GREEN: implementar el caso de uso ConsultarHistorialSprints, tests en verde`
- [ ] T013 [US1] RED: en internal/http/historial_handler_test.go, `GET /projects/{id}/historial`: 200 con los nombres exactos de `SprintHistorial` en contracts/openapi.yaml (`id`, `proyecto_id`, `numero`, `sprint_goal`, `fecha_inicio` y `fecha_fin` en `AAAA-MM-DD`, `metricas` con el mismo objeto que HU-03, armado con `aMetricasResponse`); `horas_estimadas` y `desviacion_*` en `null` sin horas por Story Point (US1.5); sin Sprints → `[]`, nunca `null`; id inválido → 400 `campo: id`; proyecto inexistente → 404 "proyecto no encontrado"; error de lectura → 500 "error interno". Stubs en internal/http/historial_handler.go e internal/http/historial_dto.go — confirmar que FALLA
- [ ] T014 [US1] COMMIT `RED: agregar tests del handler GET /projects/{id}/historial (fallan)`
- [ ] T015 [US1] GREEN: implementar `HistorialHandler.Listar` y el DTO en internal/http/historial_handler.go e internal/http/historial_dto.go — `go test ./internal/... ./cmd/...` en verde
- [ ] T016 [US1] COMMIT `GREEN: implementar el handler GET /projects/{id}/historial, tests en verde`
- [ ] T017 [US1] Cablear `ConsultarHistorialSprints` (con el `CalcularMetricas` ya construido) y `GET /projects/{id}/historial` en cmd/api/main.go
- [ ] T018 [US1] COMMIT `Cablear GET /projects/{id}/historial en cmd/api`

**Checkpoint**: la lista del historial funciona por sí sola.

---

## Phase 4: User Story 2 - Ver el detalle de un Sprint cerrado (Priority: P1)

**Goal**: `GET /projects/{id}/historial/{sprintId}` devuelve el Sprint con sus métricas, sus historias comprometidas, su esfuerzo y sus defectos, sin modificar nada.

**Independent Test**: un Sprint finalizado con una historia completada de 5 SP, una arrastrada de 3, horas imputadas y dos defectos devuelve ambas historias con su resultado, minutos por historia cuya suma coincide con `horas_reales`, y los defectos marcados; un Sprint `Activo` del mismo proyecto no cambia.

- [ ] T019 [US2] RED: en internal/domain/historial_sprint_test.go, tabla de `ArmarDetalleSprint(base SprintHistorial, historias []HistoriaComprometida, minutos map[int64]int, defectos []DefectoDetalle) DetalleSprint`: asigna a cada historia sus minutos y `0` si no tiene; calcula `Horas` desde los minutos con dos decimales reutilizando `redondear(minutosAHoras(...))` de internal/domain/metricas.go (50 min → 0.83, 90 min → 1.5, 0 min → 0; research.md R7); conserva el orden de las historias; marca `DetectadoEnSprint` si `SprintDeteccionID` es el Sprint consultado y `ResueltoEnSprint` si `SprintResolucionID` no es `nil` y es el Sprint consultado; un defecto detectado y resuelto en el mismo Sprint aparece una sola vez con ambas marcas (US2.6); sin historias ni defectos → slices vacíos no `nil` (US3.2). Firma con cuerpo vacío — confirmar que FALLA
- [ ] T020 [US2] COMMIT `RED: agregar tests del armado del detalle de un Sprint (fallan)`
- [ ] T021 [US2] GREEN: implementar `ArmarDetalleSprint` en internal/domain/historial_sprint.go — `go test ./internal/... ./cmd/...` en verde
- [ ] T022 [US2] COMMIT `GREEN: implementar el armado del detalle de un Sprint, tests en verde`
- [ ] T023 [US2] RED: en internal/repository/sqlite_historial_sprint_test.go, contra `:memory:` con datos sembrados por los repositorios existentes o por SQL: `HistoriasComprometidas` devuelve las `Completada` con `sprint_id = S` (con su `estimacion_sp` actual) y las de `historias_arrastradas` con `sprint_id = S` (con la foto de `estimacion_sp`, aunque la historia se haya editado después), completadas primero y luego arrastradas, cada grupo por `id` (research.md R5); una historia arrastrada en S1 y completada en S2 figura no completada en S1 y completada en S2 (US2.3); `estimacion_sp` `NULL` → `nil`; no incluye historias de otro Sprint. `MinutosPorHistoria` suma `registros_horas.minutos` del Sprint por historia, sin las horas de otros Sprints de la misma historia. `DefectosDeSprint` devuelve los defectos con `sprint_deteccion_id = S` o `sprint_resolucion_id = S`, una sola vez, ordenados por `id`, con el título de su historia. Sprint sin datos → slices y mapa vacíos no `nil`. Stub en internal/repository/sqlite_historial_sprint.go — confirmar que FALLA
- [ ] T024 [US2] COMMIT `RED: agregar tests de la lectura del detalle de un Sprint (fallan)`
- [ ] T025 [US2] GREEN: implementar `SQLiteHistorialSprintRepository` con sus tres consultas de solo lectura en internal/repository/sqlite_historial_sprint.go (data-model.md, sección Puerto) — `go test ./internal/... ./cmd/...` en verde
- [ ] T026 [US2] COMMIT `GREEN: leer el detalle de un Sprint, tests en verde`
- [ ] T027 [US2] RED: en internal/service/consultar_detalle_sprint_test.go, con fakes (el del puerto `HistorialSprintRepository`, en internal/service/fakes_test.go), en el orden de research.md R4: proyecto inexistente → `ErrProyectoNoEncontrado`; Sprint inexistente → `ErrSprintNoEncontrado`; Sprint de otro proyecto → `ErrSprintNoEncontrado`; Sprint `Activo` o `Pendiente` → `ErrSprintNoFinalizado` sin leer historias, horas ni defectos; con un Sprint finalizado arma el detalle con las métricas de su `SprintID` tomadas del cálculo de HU-03; un error de cualquier lectura se propaga sin datos parciales (R6); no llama a ninguna escritura (FR-014). Stub en internal/service/consultar_detalle_sprint.go — confirmar que FALLA
- [ ] T028 [US2] COMMIT `RED: agregar tests del caso de uso ConsultarDetalleSprint (fallan)`
- [ ] T029 [US2] GREEN: implementar `ConsultarDetalleSprint.Ejecutar(ctx, proyectoID, sprintID) (domain.DetalleSprint, error)` en internal/service/consultar_detalle_sprint.go — `go test ./internal/... ./cmd/...` en verde
- [ ] T030 [US2] COMMIT `GREEN: implementar el caso de uso ConsultarDetalleSprint, tests en verde`
- [ ] T031 [US2] RED: en internal/http/historial_handler_test.go, `GET /projects/{id}/historial/{sprintId}`: 200 con los nombres exactos de `DetalleSprint` en contracts/openapi.yaml; cada historia con `estimacion_sp` (`null` = "sin estimar"), `completada`, `minutos` y `horas` (el valor que ya calculó `ArmarDetalleSprint`; el DTO no convierte ni redondea, research.md R7); cada defecto con `historia_id`, `historia`, `descripcion`, `severidad`, `estado`, `detectado_en_sprint` y `resuelto_en_sprint`; sin historias ni defectos → `[]`, nunca `null`; id inválido → 400 `campo: id`; sprintId inválido → 400 `campo: sprint_id`; proyecto inexistente → 404; Sprint inexistente o de otro proyecto → 404 "sprint no encontrado"; Sprint no finalizado → 400 `campo: estado` "el Sprint todavía no está finalizado"; error de lectura → 500 "error interno" — confirmar que FALLA
- [ ] T032 [US2] COMMIT `RED: agregar tests del handler GET /projects/{id}/historial/{sprintId} (fallan)`
- [ ] T033 [US2] GREEN: implementar `HistorialHandler.Detalle` y su DTO en internal/http/historial_handler.go e internal/http/historial_dto.go; el DTO copia `Horas` del dominio sin recalcularlo — `go test ./internal/... ./cmd/...` en verde
- [ ] T034 [US2] COMMIT `GREEN: implementar el handler GET /projects/{id}/historial/{sprintId}, tests en verde`
- [ ] T035 [US2] (solo si hay algo que limpiar) REFACTOR de lo repetido entre la lista y el detalle (por ejemplo, la búsqueda de las métricas de un Sprint en `MetricasProyecto.Sprints`) + COMMIT `REFACTOR: <qué se limpió>`
- [ ] T036 [US2] Cablear `SQLiteHistorialSprintRepository`, `ConsultarDetalleSprint` y `GET /projects/{id}/historial/{sprintId}` en cmd/api/main.go
- [ ] T037 [US2] COMMIT `Cablear GET /projects/{id}/historial/{sprintId} en cmd/api`

**Checkpoint**: lista y detalle funcionan; el backend de la historia está completo.

---

## Phase 5: User Story 3 - Consultar un proyecto sin Sprints cerrados (Priority: P2)

**Goal**: un proyecto o un Sprint sin datos responde con éxito y listas vacías.

**Independent Test**: un proyecto recién creado y otro con un único Sprint `Activo` responden `200` con `[]`; un Sprint finalizado sin historias, horas ni defectos responde `200` con listas vacías y las métricas de HU-03 para ese caso.

- [ ] T038 [US3] Verificar que los casos de US3 están cubiertos sin tests nuevos: US3.1 en T009 y T013, US3.2 en T019, T023 y T031, y ambos en features/historial_sprints.feature. Si falta alguno, agregarlo con su ciclo RED → COMMIT → GREEN → COMMIT; si no falta nada, esta tarea no lleva commit

**Checkpoint**: el estado inicial del proyecto está cubierto.

---

## Phase 6: Polish (backend)

- [ ] T039 Step definitions de Godog en features/steps_historial_sprints.go (sembrar con los services reales sobre SQLite `:memory:`, como features/steps_metricas.go) y registrar `InitializeScenarioHistorial(ctx)` en features/features_test.go. Incluir un paso que compare las métricas del historial con las de `CalcularMetricas` (US1.4, FR-010) y otro que compruebe que el Sprint `Activo` no cambió (US2.7)
- [ ] T040 COMMIT `Agregar pasos BDD del historial de Sprints`
- [ ] T041 `go build ./...`, `go vet ./...`, `gofmt -l .` sin salida y `go test ./...` en verde
- [ ] T042 Prueba manual de la API según quickstart.md (pasos 1 a 8)

**Checkpoint**: sub-issues #84, #85, #86 y #87 cerradas; la historia se puede cerrar por la Definition of Done.

---

## Phase 7: User Story 4 - Navegar el historial desde la interfaz (Priority: P3)

**Goal**: sección "Historial de Sprints" con la lista y el detalle (sub-issue #88).

**Independent Test**: con pruebas de componente (Vitest y Testing Library, con el escenario en el nombre de cada prueba) y con la aplicación levantada: se abre la sección, se ve la lista, se selecciona un Sprint y se ve su detalle; sin Sprints cerrados se ve un estado vacío informativo.

- [ ] T043 [P] [US4] RED: en web/src/features/historial/api.test.ts, `listarHistorial(proyectoId)` y `obtenerDetalleSprint(proyectoId, sprintId)` piden `GET /projects/{id}/historial` y `GET /projects/{id}/historial/{sprintId}` y devuelven los tipos de contracts/openapi.yaml; una respuesta de error lanza con el mensaje del backend. Tipos en web/src/features/historial/tipos.ts y stubs en web/src/features/historial/api.ts — confirmar que FALLA
- [ ] T044 [US4] COMMIT `RED: agregar tests del cliente del historial de Sprints (fallan)`
- [ ] T045 [US4] GREEN: implementar web/src/features/historial/api.ts — `npm test` en verde
- [ ] T046 [US4] COMMIT `GREEN: implementar el cliente del historial de Sprints, tests en verde`
- [ ] T047 [US4] RED: en web/src/features/historial/Lista.test.tsx, escenarios US4.1, US4.3, US4.4 y US4.5: cada Sprint como "Sprint <número>" con su Sprint Goal, fechas y las métricas de FR-003; sin Sprints → mensaje de que todavía no hay Sprints cerrados; `null` → "sin datos", nunca cero; error → mensaje sin datos parciales. Stub en web/src/features/historial/Lista.tsx — confirmar que FALLA
- [ ] T048 [US4] COMMIT `RED: agregar pruebas de la lista del historial (fallan)`
- [ ] T049 [US4] GREEN: implementar web/src/features/historial/Lista.tsx — `npm test` en verde
- [ ] T050 [US4] COMMIT `GREEN: implementar la lista del historial, tests en verde`
- [ ] T051 [US4] RED: en web/src/features/historial/Detalle.test.tsx, escenarios US4.2, US4.4 y US4.5: muestra todas las métricas, las historias (completada o no, "sin estimar" si no tiene Story Points, horas por historia), los defectos (detectado / resuelto en el Sprint, estado actual) y un enlace para volver a la lista; error → mensaje sin datos parciales. Stub en web/src/features/historial/Detalle.tsx — confirmar que FALLA
- [ ] T052 [US4] COMMIT `RED: agregar pruebas del detalle del historial (fallan)`
- [ ] T053 [US4] GREEN: implementar web/src/features/historial/Detalle.tsx — `npm test` en verde
- [ ] T054 [US4] COMMIT `GREEN: implementar el detalle del historial, tests en verde`
- [ ] T055 [US4] Agregar web/src/features/historial/index.tsx (rutas internas: lista en `historial` y detalle en `historial/:sprintId`, recargable con F5), el estilo en web/src/features/historial/historial.css, la ruta `historial/*` en web/src/rutas.tsx y `{ ruta: 'historial', texto: 'Historial de Sprints' }` en web/src/navegacion.ts; ajustar web/src/rutas.test.tsx si enumera las secciones. La feature no importa de otras (web/src/features/aislamiento.test.ts)
- [ ] T056 [US4] COMMIT `Agregar la sección Historial de Sprints a la navegación`
- [ ] T057 [US4] `npm run lint`, `npm test` y `npm run build` en web/ sin errores; prueba manual de la interfaz según quickstart.md

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **US1 (Phase 3)** y **US2 (Phase 4)** → **US3 (Phase 5)** → **Polish backend (Phase 6)** → **US4 (Phase 7)**.
- US1 y US2 dependen solo de la Phase 2; se pueden hacer en cualquier orden, pero ambas tocan internal/http/historial_handler.go y cmd/api/main.go, así que no en paralelo.
- US3 depende de US1 y US2 (sus casos están dentro de esos tests).
- US4 depende del contrato, no del código: se puede empezar después de la Phase 2, pero la prueba manual (T057) necesita US1 y US2 cableadas.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.

## Parallel Opportunities

- La Phase 2 va en secuencia: T003–T007 tocan internal/domain/historial_sprint.go.
- T043 (cliente del frontend) puede ir en paralelo con cualquier tarea de backend desde la Phase 2.
- Dentro de US2, T019–T022 (dominio) y T023–T026 (repositorio) son independientes entre sí.

## Implementation Strategy

1. Setup + Foundational → fundación lista.
2. US1 → `GET /projects/{id}/historial` validado de forma independiente (MVP).
3. US2 → detalle; con esto el Escenario 2 del backlog queda cubierto.
4. US3 + pasos BDD + checks → PR de la historia (`Closes #12`), con #84 a #87 cerradas.
5. US4 → frontend, en el mismo PR o en uno posterior (la sub-issue #88 no frena el cierre).

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Sin cambios de esquema ni de `MODELO-DE-DATOS.md`: si alguna consulta parece necesitar una columna nueva, frenar y preguntar al equipo.
- Nunca `git push` ni squash desde las tareas.
- [P] tasks = different files, no dependencies.
