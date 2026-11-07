---

description: "Tareas de HU-18: Avance de las Historias del Proyecto (TDD con commits por fase)"
---

# Tasks: Avance de las Historias del Proyecto (HU-18)

**Input**: Design documents from `/specs/037-hu-18-avance-historias/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica es RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT).

**Compilación en RED**: todos los commits compilan. En cada RED se agregan, junto al test, solo las firmas nuevas que el test necesita con un cuerpo vacío (valor cero / `nil`), de modo que `go build ./...` pase y el test falle por aserción, no por símbolo indefinido.

**Tests en verde**: `go test ./features/` falla de T002 a la tarea de pasos BDD (T044) porque los pasos de `avance_historias.feature` todavía no están definidos (modo estricto). Por eso las tareas GREEN del backend comprueban `go test ./internal/... ./cmd/...`, y `go test ./...` se exige desde T044.

**Organization**: por user story de la spec. El backend sirve a las tres historias con un solo endpoint: US1 arma los grupos y las sumas, US2 agrega el Sprint de cada historia y US3 los casos de proyecto vacío, inexistente e id inválido. El frontend es la sub-issue de Frontend (Fase 6), que por la Definition of Done no frena el cierre de la historia.

**Firma de `AgruparAvance`**: `AgruparAvance(historias []BacklogItem, sprints []Sprint) (AvanceHistorias, error)` (data-model.md) desde US1. En US1 los tests le pasan una lista de Sprints vacía donde no importa y el service la llama con `nil`; resolver el Sprint de cada historia es el comportamiento nuevo de US2.

**Constructor del service**: `NewConsultarAvanceHistorias(proyectos, historias, sprints)` recibe sus tres dependencias desde US1 (`buscadorProyectos`, `repository.AvanceHistoriasRepository`, `listadorSprints`), para no cambiar el cableado de `cmd/api/main.go` en cada historia.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3

---

## Phase 1: Setup

**Purpose**: dejar los escenarios BDD escritos antes de implementar (flujo de AGENTS.md).

- [X] T001 Escribir features/avance_historias.feature en Gherkin español: tag `@HU-18`, Característica con el Como/Quiero/Para de la historia, Antecedentes si hay una precondición común y un Escenario por cada criterio de aceptación de US1, US2 y US3 de la spec (incluido el 400 con campo `id` por identificador inválido). Comentario inicial con la cobertura por categoría (normal, alternativo, límite, error), como en features/historial_sprints.feature. La pantalla (frontend) se cubre con pruebas de componente
- [X] T002 COMMIT `RED: agregar escenarios BDD del avance de las historias (fallan)`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: tipos de dominio, puerto y lectura SQL que usan las tres historias.

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

- [X] T003 Agregar en internal/domain/avance_historias.go los tipos de data-model.md sin lógica: `SprintRef{ID int64; Numero int}`, `HistoriaAvance{ID int64; Titulo string; EstimacionSP *int; Estado Estado; Sprint *SprintRef}` (`EstimacionSP` en `nil` = "sin estimar"; `Sprint` en `nil` = Product Backlog), `GrupoAvance{Estado Estado; Cantidad int; StoryPoints int; Historias []HistoriaAvance}` y `AvanceHistorias{Grupos []GrupoAvance; Total TotalesAvance}` con `TotalesAvance{Historias int; StoryPoints int}`. Definir el puerto `AvanceHistoriasRepository` con `ListarPorProyecto(ctx, proyectoID int64) ([]domain.BacklogItem, error)` ("todas las historias del proyecto, ORDER BY id; vacío ≠ nil") en internal/repository/avance_historias_repository.go
- [X] T004 COMMIT `Agregar los tipos del avance de historias y su puerto`
- [X] T005 RED: en internal/repository/sqlite_backlog_test.go, `ListarPorProyecto` contra SQLite `:memory:`: devuelve todas las historias del proyecto (Nueva con y sin Sprint, En progreso, Completada) ordenadas por `id`; excluye las de otros proyectos; conserva `EstimacionSP` en `nil` cuando no está estimada y `SprintID`; proyecto sin historias → slice vacío, nunca `nil`. Stub en internal/repository/sqlite_backlog.go — confirmar que FALLA
- [X] T006 COMMIT `RED: agregar tests de ListarPorProyecto de historias (fallan)`
- [X] T007 GREEN: implementar `SQLiteBacklogRepository.ListarPorProyecto` en internal/repository/sqlite_backlog.go con `listarHistorias` y `columnasHistoria`: `WHERE proyecto_id = ? ORDER BY id` — `go test ./internal/... ./cmd/...` en verde
- [X] T008 COMMIT `GREEN: implementar ListarPorProyecto de historias, test en verde`

**Checkpoint**: las historias pueden empezar.

---

## Phase 3: User Story 1 - Ver el avance de las historias por estado (Priority: P1) 🎯 MVP

**Goal**: `GET /projects/{id}/avance` devuelve las historias agrupadas en Nueva, En progreso y Completada, con la cantidad y los Story Points de cada grupo, y responde 404 y 400 en los casos de error.

**Independent Test**: con un proyecto con historias en los tres estados y distintos Story Points, la respuesta trae tres grupos con la cantidad y la suma esperadas, y las sin estimar cuentan y suman 0; un proyecto inexistente responde 404 y un id inválido, 400 con campo `id`.

- [X] T009 [US1] RED: en internal/domain/avance_historias_test.go, tabla de `AgruparAvance(historias []BacklogItem, sprints []Sprint) (AvanceHistorias, error)` pasando una lista de Sprints vacía: siempre tres grupos en el orden Nueva, En progreso, Completada, aunque estén vacíos (US1.4); cada historia en exactamente un grupo; `Cantidad` y `StoryPoints` por grupo (US1.2); historia sin Story Points cuenta en `Cantidad` y suma 0 con `EstimacionSP` en `nil` (US1.3); `Total` = suma de los tres grupos; conserva el orden recibido dentro de cada grupo; slices de `Historias` nunca `nil`; estado desconocido → error. Stub con la firma completa y cuerpo vacío en internal/domain/avance_historias.go — confirmar que FALLA
- [X] T010 [US1] COMMIT `RED: agregar tests de la agrupación del avance de historias (fallan)`
- [X] T011 [US1] GREEN: implementar la agrupación y las sumas en `AgruparAvance` (internal/domain/avance_historias.go); todavía no resuelve el Sprint de cada historia (eso es US2) — `go test ./internal/... ./cmd/...` en verde
- [X] T012 [US1] COMMIT `GREEN: implementar la agrupación del avance de historias, test en verde`
- [ ] T013 [US1] (solo si hay algo que limpiar) REFACTOR de `AgruparAvance` + COMMIT `REFACTOR: <qué se limpió>`
- [ ] T014 [US1] RED: en internal/service/consultar_avance_historias_test.go, con fakes nuevos en internal/service/fakes_test.go (puerto `AvanceHistoriasRepository`): `Ejecutar(ctx, proyectoID)` devuelve el avance de las historias que entrega el repositorio, agrupado por el dominio (historias sin Sprint); proyecto inexistente → `ErrProyectoNoEncontrado` sin leer historias ni Sprints (research.md R5); un error del repositorio de proyectos o del de historias se propaga sin datos parciales. Stub en internal/service/consultar_avance_historias.go con el constructor `NewConsultarAvanceHistorias(proyectos, historias, sprints)` — confirmar que FALLA
- [ ] T015 [US1] COMMIT `RED: agregar tests del caso de uso ConsultarAvanceHistorias (fallan)`
- [ ] T016 [US1] GREEN: implementar `ConsultarAvanceHistorias.Ejecutar(ctx, proyectoID) (domain.AvanceHistorias, error)` en internal/service/consultar_avance_historias.go: verifica el proyecto con `buscadorProyectos.ObtenerPorID`, lee las historias y llama a `AgruparAvance(historias, nil)` — `go test ./internal/... ./cmd/...` en verde
- [ ] T017 [US1] COMMIT `GREEN: implementar el caso de uso ConsultarAvanceHistorias, test en verde`
- [ ] T018 [US1] RED: en internal/http/avance_handler_test.go, `GET /projects/{id}/avance`: 200 con los nombres exactos de contracts/openapi.yaml (`total{historias,story_points}`, `grupos[{estado,cantidad,story_points,historias[]}]`, cada historia con `id`, `titulo`, `estimacion_sp` en `null` si no está estimada, `estado`, `sprint` en `null`); los tres grupos presentes y `historias` como `[]`, nunca `null`; id inválido (`abc`, `0`, `-1`) → 400 con `campo: id` y el mensaje de `parsearID`, sin consultar al repositorio; proyecto inexistente → 404 "proyecto no encontrado"; error de lectura → 500 "error interno" sin cuerpo parcial. Stubs con cuerpo vacío en internal/http/avance_handler.go e internal/http/avance_dto.go — confirmar que FALLA (todos los casos, por aserción)
- [ ] T019 [US1] COMMIT `RED: agregar tests del handler GET /projects/{id}/avance (fallan)`
- [ ] T020 [US1] GREEN: implementar `AvanceHandler.Consultar` con el orden de evaluación id inválido (`parsearID`, 400) → proyecto inexistente (404) → lectura, y el DTO que solo copia, en internal/http/avance_handler.go e internal/http/avance_dto.go — `go test ./internal/... ./cmd/...` en verde
- [ ] T021 [US1] COMMIT `GREEN: implementar el handler GET /projects/{id}/avance, test en verde`
- [ ] T022 [US1] Cablear `ConsultarAvanceHistorias` (con `proyectos`, `backlogRepo` y `sprintRepo` ya construidos), `AvanceHandler` y `GET /projects/{id}/avance` en cmd/api/main.go
- [ ] T023 [US1] COMMIT `Cablear GET /projects/{id}/avance en cmd/api`

**Checkpoint**: los grupos, los totales y los errores 400 y 404 funcionan por sí solos (con `sprint: null` hasta US2).

---

## Phase 4: User Story 2 - Ver el Sprint de cada historia (Priority: P1)

**Goal**: cada historia informa el Sprint al que pertenece (id y número) o `null` si está en el Product Backlog.

**Independent Test**: con una historia asignada a un Sprint y otra sin Sprint, la respuesta informa `{id, numero}` para la primera y `null` para la segunda.

- [ ] T024 [US2] RED: en internal/domain/avance_historias_test.go, resolución del Sprint en `AgruparAvance` con una lista de Sprints: la historia con `SprintID` recibe `Sprint` = `{ID, Numero}` del Sprint de la lista (US2.1); sin `SprintID` queda `Sprint` en `nil` (US2.2); una historia Nueva asignada a un Sprint no iniciado conserva ese Sprint (Clarification 2); una historia que referencia un Sprint ausente de la lista → error (fuentes inconsistentes, research.md R3) — confirmar que FALLA
- [ ] T025 [US2] COMMIT `RED: agregar tests del Sprint de cada historia en el dominio (fallan)`
- [ ] T026 [US2] GREEN: resolver `SprintID → SprintRef` en `AgruparAvance` (internal/domain/avance_historias.go) — `go test ./internal/... ./cmd/...` en verde (los tests de service y handler de US2 todavía no existen)
- [ ] T027 [US2] COMMIT `GREEN: informar el Sprint de cada historia en el dominio, test en verde`
- [ ] T028 [US2] RED: en internal/service/consultar_avance_historias_test.go, con un fake de Sprints con número: el service lista los Sprints del proyecto y los pasa al dominio, de modo que cada historia trae su Sprint o `nil`; un error del listado de Sprints se propaga; una historia que apunta a un Sprint ausente del listado devuelve error interno (no `ErrSprintNoEncontrado`, que respondería 404) — confirmar que FALLA (el service de US1 pasa `nil`)
- [ ] T029 [US2] COMMIT `RED: agregar tests del listado de Sprints en ConsultarAvanceHistorias (fallan)`
- [ ] T030 [US2] RED: en internal/http/avance_handler_test.go, con el service real y fakes: `sprint` serializado como `{id, numero}` para una historia con Sprint y como `null` para una sin Sprint (contracts/openapi.yaml) — confirmar que FALLA (el service todavía no entrega el Sprint)
- [ ] T031 [US2] COMMIT `RED: agregar tests del Sprint de cada historia en el handler de avance (fallan)`
- [ ] T032 [US2] GREEN: en `ConsultarAvanceHistorias` (internal/service/consultar_avance_historias.go) listar los Sprints con `listadorSprints.ListarPorProyecto` y pasarlos a `AgruparAvance`, propagando los errores sin envolverlos como 404; el DTO ya copia `HistoriaAvance.Sprint` — `go test ./internal/... ./cmd/...` en verde
- [ ] T033 [US2] COMMIT `GREEN: informar el Sprint de cada historia en service y handler, test en verde`

**Checkpoint**: la vista dice en qué Sprint está cada historia.

---

## Phase 5: User Story 3 - Proyecto sin historias (Priority: P2)

**Goal**: un proyecto sin historias responde 200 con los tres grupos y los totales en cero. El 404 y el 400 ya se implementaron en US1 (T014-T021).

**Independent Test**: se consulta un proyecto recién creado: 200, tres grupos con `cantidad: 0`, `story_points: 0` y `historias: []`, y `total` en cero.

- [ ] T034 [US3] Tests de caracterización (no son RED: pasan de entrada porque US1.4 y T009 ya cubren los grupos vacíos): en internal/service/consultar_avance_historias_test.go y internal/http/avance_handler_test.go, proyecto existente sin historias → avance con los tres grupos y `total` en cero (service) y 200 con ese cuerpo y `historias: []` (handler). Si alguno falla, es un defecto de US1: se corrige con el ciclo RED/GREEN completo en su propio commit
- [ ] T035 [US3] COMMIT `Agregar tests de proyecto sin historias en el avance (caracterización)`

**Checkpoint**: el backend cubre los cuatro escenarios del backlog.

---

## Phase 6: Frontend — pantalla "Estado" (sub-issue de Frontend; US1, US2 y US3)

**Purpose**: reemplazar el placeholder de web/src/features/estado/index.tsx. No toca rutas.tsx ni navegacion.ts (`status/*` ya apunta a `./features/estado`) y no importa de otras features (web/src/features/aislamiento.test.ts).

- [ ] T036 [US1] [US2] RED: en web/src/features/estado/avance.test.ts y web/src/features/estado/api.test.ts, tipos y funciones puras: `consultarAvance(proyectoId)` pide `GET /projects/{id}/avance` y rechaza una respuesta sin `grupos` (como `consultarMetricas`); `textoDeSprint(sprint)` devuelve "Sprint <número>" o "Product Backlog" (US2); `textoDeEstimacion(sp)` devuelve "sin estimar" para `null` y nunca "0" (US1.3) — stubs en web/src/features/estado/avance.ts y web/src/features/estado/api.ts — confirmar que FALLA
- [ ] T037 [US1] [US2] COMMIT `RED: agregar tests de la API y las funciones del avance (fallan)`
- [ ] T038 [US1] [US2] GREEN: implementar web/src/features/estado/api.ts y web/src/features/estado/avance.ts — `npm test` en verde
- [ ] T039 [US1] [US2] COMMIT `GREEN: implementar la API y las funciones del avance, test en verde`
- [ ] T040 [US1] [US2] [US3] RED: en web/src/features/estado/index.test.tsx y GrupoHistorias.test.tsx, pruebas de componente con el escenario en el nombre: tres grupos con cantidad y Story Points (US1.1-US1.2) e historia "sin estimar" (US1.3); Sprint de cada historia o "Product Backlog" (US2.1-US2.2); proyecto sin historias → mensaje y totales en cero (US3.1); id de la ruta inválido → mensaje de error sin pedir al servidor; carga, error de API y error de conexión como estados distintos, sin datos parciales — confirmar que FALLA
- [ ] T041 [US1] [US2] [US3] COMMIT `RED: agregar pruebas de componente de la pantalla Estado (fallan)`
- [ ] T042 [US1] [US2] [US3] GREEN: reemplazar web/src/features/estado/index.tsx (validar el id con `/^[1-9]\d*$/` antes de pedir, `LoadingIndicator`, `ErrorMessage`), crear web/src/features/estado/GrupoHistorias.tsx y web/src/features/estado/estado.css con el diseño actual y los componentes de web/src/compartido/componentes, sin librerías nuevas, sin desborde en 320 px — `npm run lint && npm test && npm run build` en verde
- [ ] T043 [US1] [US2] [US3] COMMIT `GREEN: implementar la pantalla Estado, test en verde`

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T044 Step definitions de Godog para features/avance_historias.feature en features/steps_avance_historias.go y registro (`InitializeScenarioAvance(ctx)`) en features/features_test.go — `go test ./features/` en verde
- [ ] T045 COMMIT `Agregar pasos BDD del avance de las historias`
- [ ] T046 Verificación final: `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test ./...` sin errores; `cd web && npm run lint && npm test && npm run build`. Comprobar **FR-007** (solo lectura): `AvanceHistoriasRepository` expone únicamente `ListarPorProyecto`, la única ruta nueva es `GET /projects/{id}/avance`, `git diff origin/main --stat` no toca `MODELO-DE-DATOS.md` ni `internal/repository/migraciones/`, y `GET /projects/{id}/backlog` (HU-15) sigue devolviendo solo las historias sin Sprint (sus pruebas en verde). Comprobar **FR-008**: `git diff origin/main --stat` no toca el estado Planificado / En curso / Finalizado del panel (`features/estado_proyecto.feature`, `features/steps_estado_proyecto.go`, `web/src/features/proyecto/`) y esas pruebas siguen en verde
- [ ] T047 Recorrer quickstart.md (API y pantalla) y corregirlo si algo no coincide; COMMIT `Ajustar el quickstart de la spec 037` solo si hubo cambios

---

## Dependencies & Execution Order

- **Setup (Fase 1)** → **Foundational (Fase 2)** → **US1 (Fase 3)** → **US2 (Fase 4)** → **US3 (Fase 5)** → **Polish (Fase 7)**.
- **Frontend (Fase 6)** depende solo del contrato (contracts/openapi.yaml): puede hacerse en paralelo al backend desde el fin de la Fase 2, pero su prueba manual (T047) necesita el backend completo.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Cada RED (T005, T009, T014, T018, T024, T028, T030, T036, T040) va antes del GREEN que lo satisface (T007, T011, T016, T020, T026, T032, T038, T042). Nunca saltear el COMMIT de RED.
- T032 satisface a la vez los RED de T028 y T030, que quedan commiteados antes.
- T034 son tests de caracterización (pasan de entrada), no un ciclo RED/GREEN: US3 no agrega comportamiento nuevo al backend.
- T044 es el único punto donde `go test ./features/` vuelve a verde.

## Implementation Strategy

1. Setup + Foundational → fundación lista.
2. US1 → validar de forma independiente (MVP: grupos, totales y errores).
3. US2 y US3 de a una, cada una dejando el sistema funcionando.
4. Frontend en paralelo o al final; no frena el cierre del issue padre.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- Sin cambios de esquema ni de MODELO-DE-DATOS.md.
- [P] no se usó: las tareas comparten archivos (`avance_historias.go`, `avance_handler_test.go`) y el orden TDD es secuencial.
