---

description: "Tareas de HU-03: Cálculo de Métricas del Proyecto (TDD con commits por fase)"
---

# Tasks: Cálculo de Métricas del Proyecto (HU-03)

**Input**: Design documents from `/specs/032-hu-03-calculo-metricas/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica es RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT).

**Compilación en RED**: todos los commits compilan. En cada RED se agregan, junto al test, solo las
firmas nuevas que el test necesita con un cuerpo vacío (devolver valor cero / `nil`), de modo que
`go build ./...` pase y el test falle por aserción, no por símbolo indefinido. El RED de un service
o handler que cambia un constructor actualiza también a quienes lo llaman.

**Tests en verde**: `go test ./features/` falla de T002 a T062 porque los pasos de
`calculo_metricas.feature` todavía no están definidos (modo estricto). Por eso las tareas GREEN
comprueban `go test ./internal/... ./cmd/...`, y `go test ./...` se exige desde T064.

**Organization**: por user story de la spec. US1 y US2 comparten el cálculo; US3, US4 y US5 no dependen de él.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3 / US4 / US5

---

## Phase 1: Setup

**Purpose**: dejar los escenarios BDD escritos antes de implementar (flujo de AGENTS.md).

- [X] T001 Escribir features/calculo_metricas.feature en Gherkin español: tag `@HU-03`, Característica con el Como/Quiero/Para de la historia y un Escenario por cada criterio de aceptación de US1, US2, US3 y US5 de la spec. US4 se cubre con pruebas de componente
- [X] T002 COMMIT `RED: agregar escenarios BDD de cálculo de métricas (fallan)`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: leer y escribir el parámetro; lo usan US1, US2 y US3.

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

- [X] T003 RED: en internal/repository/sqlite_project_test.go, con el valor sembrado por SQL: `ObtenerPorID` y `Listar` devuelven `HorasPorStoryPoint` 2.5, y `nil` si la columna es `NULL`; `Actualizar` guarda un valor y con `nil` lo borra; leer el proyecto, cambiar el nombre y actualizar conserva el valor (research.md R8). Agregar en internal/domain/project.go el campo `HorasPorStoryPoint *float64` — confirmar que FALLA
- [X] T004 COMMIT `RED: agregar tests de la persistencia de las horas por Story Point (fallan)`
- [X] T005 GREEN: en internal/repository/sqlite_project.go, agregar `horas_por_story_point` a los dos `SELECT`, a `escanearProyecto` (con `sql.NullFloat64`) y al `UPDATE` de `Actualizar`; el `INSERT` no cambia — `go test ./internal/... ./cmd/...` en verde
- [X] T006 COMMIT `GREEN: persistir las horas por Story Point, tests en verde`

**Checkpoint**: las historias pueden empezar.

---

## Phase 3: User Story 1 - Consultar las métricas de cada Sprint cerrado y del proyecto (Priority: P1) 🎯 MVP

**Goal**: una sola consulta devuelve las métricas de cada Sprint finalizado y del proyecto.

**Independent Test**: con un proyecto de 4 horas por Story Point y un Sprint finalizado con una historia completada de 5 SP, una arrastrada de 3, 30 horas y dos defectos (uno resuelto), `GET /projects/{id}/metricas` informa 8 y 5 SP, velocidad 5, 50%, 32 y 30 horas, −2 horas, −6,25%, 2 detectados y 1 resuelto.

- [X] T007 [US1] RED: en internal/domain/metricas_test.go, tabla de `CalcularMetricas` para un Sprint: el caso del Independent Test; historias sin Story Points (cuentan en el porcentaje y suman cero); Sprint sin historias comprometidas (ceros y 0%); sin parámetro (`HorasEstimadas`, `DesviacionHoras` y `DesviacionPorcentaje` en `nil`, las demás con su valor); con parámetro y 0 SP planificados (estimadas 0, desviación igual a las reales y porcentaje `nil`); sin horas (reales 0). Agregar en internal/domain/metricas.go los tipos de data-model.md y la firma con cuerpo vacío — confirmar que FALLA
- [X] T008 [US1] COMMIT `RED: agregar tests de las métricas de un Sprint (fallan)`
- [X] T009 [US1] GREEN: implementar el cálculo por Sprint en internal/domain/metricas.go, con la aritmética de research.md R4 — `go test ./internal/... ./cmd/...` en verde
- [X] T010 [US1] COMMIT `GREEN: implementar las métricas de un Sprint, tests en verde`
- [X] T011 [US1] RED: en internal/domain/metricas_test.go, valores del proyecto: dos Sprints con 5 y 8 SP completados → 13 y velocidad 6,5; porcentajes recalculados sobre los totales; historia de 3 SP arrastrada en un Sprint y completada en otro → 6 planificados, 3 completados, 24 horas y 50% (US1.9); `Sprints` conserva el orden recibido; redondeo a dos decimales calculado desde los totales enteros y no desde los valores redondeados de cada Sprint; una desviación de −0,001 horas se informa como `0` y no como `-0` — confirmar que FALLA
- [X] T012 [US1] COMMIT `RED: agregar tests de las métricas del proyecto (fallan)`
- [X] T013 [US1] GREEN: implementar los valores del proyecto en internal/domain/metricas.go — `go test ./internal/... ./cmd/...` en verde
- [X] T014 [US1] COMMIT `GREEN: implementar las métricas del proyecto, tests en verde`
- [X] T015 [US1] Definir el puerto `DatosMetricasRepository` con `ListarDatosSprintsFinalizados(ctx, proyectoID) ([]domain.DatosMetricasSprint, error)` en internal/repository/metricas_repository.go
- [X] T016 [US1] COMMIT `Agregar el puerto de los datos de métricas`
- [X] T017 [US1] RED: en internal/repository/sqlite_metricas_test.go, contra `:memory:`: solo aparecen los Sprints `Finalizado` del proyecto, y el `Activo` y el `Pendiente` no aportan horas ni defectos; las completadas suman su `estimacion_sp` y las arrastradas su foto, aunque después se edite la historia; `NULL` suma cero y cuenta como historia; las horas cuentan en el Sprint al que se imputaron aunque la historia ya no esté ahí; un defecto detectado en un Sprint y resuelto en otro cuenta en cada uno; un Sprint sin datos da ceros; `Numero` es la posición entre todos los Sprints del proyecto, con Sprints de otro proyecto intercalados y un Sprint no finalizado de `id` menor (research.md R3); proyecto sin Sprints finalizados → slice vacío no `nil`. Stub en internal/repository/sqlite_metricas.go — confirmar que FALLA
- [X] T018 [US1] COMMIT `RED: agregar tests de la lectura de los datos de métricas (fallan)`
- [X] T019 [US1] GREEN: implementar la consulta única con subconsultas correlacionadas en internal/repository/sqlite_metricas.go, ordenada por `sprints.id` (research.md R2 y R3) — `go test ./internal/... ./cmd/...` en verde
- [X] T020 [US1] COMMIT `GREEN: leer los datos de métricas, tests en verde`
- [X] T021 [US1] RED: en internal/service/calcular_metricas_test.go, con fakes: proyecto inexistente → `ErrProyectoNoEncontrado` sin consultar los datos; usa `HorasPorStoryPoint` del proyecto; un error del repositorio se propaga; no llama a ninguna escritura (FR-015). Stub en internal/service/calcular_metricas.go — confirmar que FALLA
- [X] T022 [US1] COMMIT `RED: agregar tests del caso de uso CalcularMetricas (fallan)`
- [X] T023 [US1] GREEN: implementar `CalcularMetricas.Ejecutar` en internal/service/calcular_metricas.go — `go test ./internal/... ./cmd/...` en verde
- [X] T024 [US1] COMMIT `GREEN: implementar el caso de uso CalcularMetricas, tests en verde`
- [X] T025 [US1] RED: en internal/http/metricas_handler_test.go, `GET /projects/{id}/metricas`: 200 con los nombres exactos de contracts/openapi.yaml; `horas_estimadas`, `desviacion_horas` y `desviacion_porcentaje` en `null` sin parámetro; id inválido → 400; proyecto inexistente → 404 "proyecto no encontrado"; error del repositorio → 500 "error interno". Stubs en internal/http/metricas_handler.go e internal/http/metricas_dto.go — confirmar que FALLA
- [X] T026 [US1] COMMIT `RED: agregar tests del handler GET /projects/{id}/metricas (fallan)`
- [X] T027 [US1] GREEN: implementar `MetricasHandler.Consultar` y sus DTO en internal/http/metricas_handler.go e internal/http/metricas_dto.go — `go test ./internal/... ./cmd/...` en verde
- [X] T028 [US1] COMMIT `GREEN: implementar el handler GET /projects/{id}/metricas, tests en verde`
- [X] T029 [US1] Cablear el repositorio, el service, el handler y la ruta `GET /projects/{id}/metricas` en cmd/api/main.go
- [X] T030 [US1] COMMIT `Cablear GET /projects/{id}/metricas en cmd/api`

**Checkpoint**: la consulta funciona con la aplicación levantada (quickstart.md, paso 1).

---

## Phase 4: User Story 2 - Consultar un proyecto sin Sprints cerrados (Priority: P2)

**Goal**: la consulta no falla y devuelve un estado vacío.

**Independent Test**: un proyecto recién creado responde 200, con `sprints: []` y los valores en cero; `desviacion_porcentaje` es `null` y, sin parámetro, también `horas_estimadas` y `desviacion_horas`.

- [ ] T031 [US2] RED: en internal/domain/metricas_test.go, `CalcularMetricas` sin datos: con parámetro, ceros y `DesviacionPorcentaje` `nil`; sin parámetro, además `HorasEstimadas` y `DesviacionHoras` `nil`; `Velocidad` es 0 y no `NaN`; `Sprints` es vacío y no `nil`. En internal/http/metricas_handler_test.go, el cuerpo tiene `"sprints":[]` — confirmar que FALLA
- [ ] T032 [US2] COMMIT `RED: agregar tests de las métricas de un proyecto sin Sprints cerrados (fallan)`
- [ ] T033 [US2] GREEN: cubrir el caso sin Sprints en internal/domain/metricas.go — `go test ./internal/... ./cmd/...` en verde
- [ ] T034 [US2] COMMIT `GREEN: informar las métricas de un proyecto sin Sprints cerrados, tests en verde`

**Checkpoint**: US1 y US2 completan el cálculo (MVP).

---

## Phase 5: User Story 3 - Definir las horas por Story Point del proyecto (Priority: P2)

**Goal**: editar el proyecto define, cambia o borra el parámetro, y las métricas lo usan.

**Independent Test**: un `PUT` con 4 guarda el valor; sin el campo lo conserva; con `null` lo borra; con 0, −1, `"abc"` o `""` responde 400 sobre el campo y el proyecto no cambia.

- [ ] T035 [US3] RED: en internal/domain/project_test.go, tabla de `ConHorasPorStoryPoint(valor *string)`: `nil` → sin definir; `4`, `2.5` y `1e2` → válidos; `0`, `-1`, `"abc"`, `""`, `"5"` (con comillas), `true`, `[]`, `1e999`, `Inf` y `NaN` → `ErrHorasPorStoryPointInvalidas`, `ValidationError{Campo: "horas_por_story_point", Mensaje: "las horas por Story Point deben ser un número mayor que cero"}`. Agregar el error en internal/domain/errors.go y el stub en internal/domain/project.go — confirmar que FALLA
- [ ] T036 [US3] COMMIT `RED: agregar tests de la validación de las horas por Story Point (fallan)`
- [ ] T037 [US3] GREEN: implementar `ConHorasPorStoryPoint` en internal/domain/project.go (research.md R7) — `go test ./internal/... ./cmd/...` en verde
- [ ] T038 [US3] COMMIT `GREEN: validar las horas por Story Point, tests en verde`
- [ ] T039 [US3] RED: en internal/service/editar_proyecto_test.go, con `HorasPorSPEnviadas` y `HorasPorSP` en `EditarProyectoInput`: sin enviar conserva; enviado `nil` borra; enviado `4` guarda; un valor inválido no llama a `Actualizar`; con nombre e horas inválidos el error es el del nombre (edge case del orden). Agregar los campos en internal/service/editar_proyecto.go — confirmar que FALLA
- [ ] T040 [US3] COMMIT `RED: agregar tests de la edición de las horas por Story Point (fallan)`
- [ ] T041 [US3] GREEN: aplicar el parámetro después de `ConDatosEditados` en internal/service/editar_proyecto.go — `go test ./internal/... ./cmd/...` en verde
- [ ] T042 [US3] COMMIT `GREEN: editar las horas por Story Point, tests en verde`
- [ ] T043 [US3] RED: en internal/http/project_handler_test.go, `PUT /projects/{id}`: sin el campo conserva y responde 200, aunque los otros cuatro siguen obligatorios; `null` borra; `4` y `2.5` guardan; `0`, `-1`, `"abc"` y `""` → 400 con `campo: horas_por_story_point`, nunca "JSON inválido"; con nombre vacío y horas inválidas el campo es `nombre`. `GET /projects/{id}`, `GET /projects` y `POST /projects` incluyen `horas_por_story_point` (`null` si no está definido) — confirmar que FALLA
- [ ] T044 [US3] COMMIT `RED: agregar tests del parámetro en el PUT y la consulta del proyecto (fallan)`
- [ ] T045 [US3] GREEN: en internal/http/project_handler.go, pasar el valor desde `claves` sin interpretarlo y sin agregarlo a los campos obligatorios; `null` llega como `nil`. En internal/http/project_dto.go, agregar el campo a `proyectoResponse` — `go test ./internal/... ./cmd/...` en verde
- [ ] T046 [US3] COMMIT `GREEN: aceptar y devolver las horas por Story Point, tests en verde`
- [ ] T047 [P] [US3] Actualizar specs/002-hu-04-crear-proyecto-equipo/contracts/openapi.yaml con el campo opcional de la edición y el de la respuesta
- [ ] T048 [US3] COMMIT `docs: agregar las horas por Story Point al contrato de HU-04`

**Checkpoint**: el parámetro cambia las métricas, incluidos los Sprints cerrados (quickstart.md, pasos 2 a 4).

---

## Phase 6: User Story 5 - Mantener fijos los Story Points de un Sprint finalizado (Priority: P2)

**Goal**: no se pueden cargar Story Points en una historia de un Sprint finalizado.

**Independent Test**: con una historia completada y sin estimar de un Sprint finalizado, cargarle Story Points por edición o por votación responde 400 sobre `estimacion_sp`; editar su título funciona y las métricas no cambian.

- [ ] T049 [US5] RED: en internal/service/editar_historia_test.go, cambiar "finalizado permite la carga inicial" para que se rechace con `ErrCargaEnSprintFinalizado`, y agregar "finalizado permite editar sin cargar" (`nil` → `nil`), "finalizado permite reenviar el mismo valor", "pendiente permite la carga inicial" y "sin Sprint permite la carga inicial" (FR-025). En internal/service/registrar_estimacion_test.go, sacar "primera carga con un Sprint Finalizado" de los casos válidos y agregarla a los rechazados, también con la historia `Completada`. Agregar `ErrCargaEnSprintFinalizado = ValidationError{Campo: "estimacion_sp", Mensaje: "no se pueden cargar Story Points en un Sprint Finalizado"}` en internal/domain/errors.go — confirmar que FALLA
- [ ] T050 [US5] RED: en internal/http/backlog_handler_test.go y internal/http/votacion_handler_test.go, 400 con `campo: estimacion_sp` y ese mensaje. En features/consulta_edicion_backlog.feature, sacar el ejemplo `Finalizado` del esquema "Carga inicial de Story Points" y actualizar el comentario de cobertura del encabezado — confirmar que FALLA
- [ ] T051 [US5] COMMIT `RED: agregar tests del rechazo de la carga con Sprint Finalizado (fallan)`
- [ ] T052 [US5] GREEN: reordenar `validarEstimacionCongelada` en internal/service/estimacion_congelada.go según research.md R9 — `go test ./internal/... ./cmd/...` en verde
- [ ] T053 [US5] COMMIT `GREEN: rechazar la carga de Story Points con Sprint Finalizado, tests en verde`
- [ ] T054 [P] [US5] Actualizar MODELO-DE-DATOS.md: §4.4, §4.6, §5 y D2b con la regla nueva y la fecha de la reunión, y §6 con que las horas estimadas salen de los SP planificados
- [ ] T055 [P] [US5] Actualizar docs/PRODUCT-BACKLOG.md: decisiones transversales, HU-14 Escenario 4, HU-15 Escenario 3 y la nota de HU-03 (la regla nueva y la excepción por cambiar el parámetro)
- [ ] T056 [P] [US5] Actualizar specs/009-hu-15-consulta-edicion-backlog/spec.md (Clarifications, escenario 6, edge case y FR) y specs/031-hu-14-consenso-estimacion/spec.md (FR-010)
- [ ] T057 [US5] COMMIT `docs: registrar la regla de la carga de Story Points con Sprint Finalizado (D2b)`

**Checkpoint**: los documentos y el código coinciden en la regla nueva.

---

## Phase 7: User Story 4 - Definir el parámetro desde la pantalla de edición del proyecto (Priority: P3)

**Goal**: la pantalla muestra y modifica el parámetro (sub-issue #53). No frena el cierre de #3.

**Independent Test**: se abre la edición, se carga un valor, se guarda y al volver a abrir el valor sigue ahí; vaciar el campo lo deja sin definir; un valor inválido muestra el mensaje junto al campo.

- [ ] T058 [US4] RED: en web/src/features/proyecto/validacion.test.ts, `validarEdicion` acepta vacío, `4` y `2.5`, y rechaza `0` y `-1`; agregar `horas_por_story_point` a los datos de ese test, para que el tipado siga pasando. En web/src/features/proyecto/Editar.test.tsx, con el escenario en el nombre de cada prueba: US4.1 muestra el valor actual; US4.2 guarda y envía el número; US4.3 muestra junto al campo el error `horas_por_story_point` del backend; US4.4 vacío envía `null` — `npm test` confirma que FALLA
- [ ] T059 [US4] COMMIT `RED: agregar pruebas del campo de horas por Story Point (fallan)`
- [ ] T060 [US4] GREEN: según research.md R11, agregar el campo en web/src/features/proyecto/validacion.ts, web/src/features/proyecto/api.ts, web/src/features/proyecto/Editar.tsx y web/src/features/proyecto/crear.tsx (solo el valor inicial) — `npm test` en verde
- [ ] T061 [US4] COMMIT `GREEN: editar las horas por Story Point desde la pantalla, pruebas en verde`

**Checkpoint**: las cinco historias funcionan.

---

## Phase 8: Polish & Cross-Cutting Concerns

- [ ] T062 Implementar features/steps_metricas.go con `InitializeScenarioMetricas`, reutilizando pasos existentes sin colisiones de regex, y registrarlo en features/features_test.go — `go test ./features/` en verde
- [ ] T063 COMMIT `Agregar pasos BDD de cálculo de métricas`
- [ ] T064 Correr `go build ./...`, `go vet ./...`, `gofmt -l .` y `go test ./...`, y en web/ `npm run lint`, `npm test` y `npm run build`, sin errores; validar la sintaxis de specs/032-hu-03-calculo-metricas/contracts/openapi.yaml
- [ ] T065 Validar specs/032-hu-03-calculo-metricas/quickstart.md con la aplicación levantada; si algún valor no coincide, se frena y se avisa
- [ ] T066 (solo si hay algo que limpiar) REFACTOR de lo que haya quedado duplicado + COMMIT `REFACTOR: <qué se limpió>`

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **User Stories (Phase 3 a 7)** → **Polish (Phase 8)**.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- **US2** depende de US1 (mismo cálculo y handler). **US3** solo depende de Phase 2. **US5** no depende de ninguna historia. **US4** depende del backend de US3.
- T001 falla en modo estricto hasta T062, porque los pasos llaman a código que debe existir, igual que en HU-14.

### Parallel Opportunities

El historial RED/GREEN es lineal, así que dentro de un mismo clon el código se hace en orden. Solo son paralelas las tareas de documentos:

- US5: T054, T055 y T056 (archivos distintos).
- US3: T047, respecto de cualquier otra tarea de documentos.
- Entre personas: US3, US4 y US5 no comparten archivos con el cálculo de US1, salvo internal/domain/errors.go.

## Implementation Strategy

1. Setup + Foundational → escenarios y persistencia listos.
2. US1 + US2 → MVP: destraba HU-09, HU-10 y HU-12.
3. US3 → las horas estimadas se pueden definir desde la API.
4. US5 → las métricas de los Sprints cerrados quedan fijas.
5. US4 → pantalla (#53).
6. Polish → pasos BDD, validación final y PR contra `main` a cargo del equipo (`Closes #3`).

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- Sin migración ni cambios de tablas. Si aparece una contradicción con `MODELO-DE-DATOS.md`, se frena y se avisa.
- El PR abierto #268 modifica `specs/031-…/spec.md` y `consenso_estimacion.feature`: T056 puede requerir resolver un conflicto. El #250 modifica `docs/PRODUCT-BACKLOG.md` (T055).
- Spec 033 (HU-12): el ajuste de su supuesto sobre D2b se hace en su propia rama.
