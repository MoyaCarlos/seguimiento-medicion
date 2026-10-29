---

description: "Tareas de HU-06: Movimiento de Historias al Sprint Backlog (TDD con commits por fase)"
---

# Tasks: Movimiento de Historias al Sprint Backlog (HU-06)

**Input**: Design documents from `/specs/018-hu-06-asignar-historias-sprint/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica es RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT).

**Compilación en RED**: todos los commits compilan. En cada RED se agregan, junto al test, solo
las firmas nuevas que el test necesita con un cuerpo vacío (devolver valor cero / `nil`), de modo
que `go build ./...` pase y el test falle por aserción, no por símbolo indefinido.

**Organization**: por user story de la spec (US1 asignar, US2 consultar el Sprint Backlog, US3
pantalla). US1 queda completa, con su escenario BDD en verde, antes de empezar US2. US3 está
**bloqueada** hasta que estén en `main` el PR #149 y el #98.

**Sin migraciones**: no se toca `MODELO-DE-DATOS.md` ni `internal/repository/migraciones/`.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3

---

## Phase 1: Setup

**Purpose**: partir de una base en verde.

- [X] T001 Verificar en la rama `feature/HU-06-asignar-historias-sprint` que `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test ./...` están en verde antes de empezar. Sin cambios ni commit

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: exponer el Sprint de cada historia en la API, que usan la respuesta de la asignación (US1) y el Sprint Backlog (US2).

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

- [X] T002 RED: en internal/http/backlog_handler_test.go, test de que `GET /projects/{id}/backlog` devuelve en cada historia la clave `"sprint_id"` con valor `null` (decodificar a `map[string]any` y comprobar que la clave existe y es `nil`). Confirmar que FALLA (hoy `historiaResponse` no tiene el campo)
- [X] T003 COMMIT `RED: agregar test de sprint_id en la respuesta de historia (falla)`
- [X] T004 GREEN: en internal/http/backlog_dto.go, agregar a `historiaResponse` el campo `SprintID *int64 \`json:"sprint_id"\`` (sin `omitempty`: `null` sin Sprint, research.md R5) y cargarlo en `aRespuesta` desde `item.SprintID` — `go test ./...` en verde
- [X] T005 COMMIT `GREEN: exponer sprint_id en la respuesta de historia, test en verde`

**Checkpoint**: la API ya informa a qué Sprint pertenece cada historia.

---

## Phase 3: User Story 1 - Asignar una historia del Product Backlog a un Sprint (Priority: P1) 🎯 MVP

**Goal**: `PUT /backlog/{id}/sprint` con `{"sprint_id": N}` asigna una historia sin Sprint a un Sprint `Pendiente` o `Activo` de su proyecto, conservando estado `Nueva` y Story Points, con el orden de errores de la spec.

**Independent Test**: sembrar un proyecto con un Sprint `Pendiente` (y otro `Activo`) y una historia sin Sprint; asignarla y verificar que sale de `GET /projects/{id}/backlog`, tiene `sprint_id` cargado, estado `Nueva` y los mismos Story Points; provocar cada rechazo (otro proyecto, Finalizado, Completada, ya asignada, mismo Sprint) y los dos 404.

### Escenarios BDD de US1

- [X] T006 [US1] Escribir features/asignar_historias_sprint.feature en Gherkin español: tag `@HU-06`, Característica con el Como/Quiero/Para de la HU-06 ("Como Product Builder quiero seleccionar historias de usuario del Product Backlog y asignarlas a un Sprint específico para poder definir el compromiso de trabajo de la iteración actual"), Antecedentes `Dado que existe un proyecto con identificador 1` (paso existente de features/steps_creacion.go, no se redefine), y exactamente 9 escenarios, uno por Acceptance Scenario de US1 en spec.md: US1.1 asignar a Sprint "Pendiente" → sale del Product Backlog y queda en el Sprint Backlog; US1.2 ídem con Sprint "Activo"; US1.3 conserva estado "Nueva" y sus Story Points (usar una historia con 5 Story Points); US1.4 historia "Completada" → rechazo con el mensaje "solo se pueden planificar historias pendientes" y sin modificar; US1.5 historia que ya tiene Sprint, asignada **al mismo** Sprint → rechazo "la historia ya pertenece a un Sprint"; US1.6 Sprint "Finalizado" → rechazo "un Sprint finalizado no admite historias"; US1.7 Sprint de otro proyecto → rechazo "el Sprint no pertenece al proyecto de la historia"; US1.8 historia inexistente → "historia no encontrada"; US1.9 historia existente y Sprint inexistente → "sprint no encontrado". Comentario en el archivo con la categoría de cada escenario (normal: 1; alternativo: 2; límite: 3 y 5; error: 4, 6, 7, 8, 9). Sin pasos implementados: con el runner estricto, `go test ./features/` FALLA
- [X] T007 [US1] COMMIT `RED: agregar escenarios BDD de asignación de historias a un Sprint (fallan)`

### Dominio: regla de la asignación

- [X] T008 [P] [US1] RED: crear internal/domain/asignacion_sprint_test.go con una tabla de casos de `ValidarAsignacion(h BacklogItem, s Sprint) error`, sin base de datos: (a) historia `Nueva` sin Sprint + Sprint `Pendiente` del mismo proyecto → `nil`; (b) ídem con Sprint `Activo` → `nil`; (c) historia sin Story Points → `nil`; (d) `h.ProyectoID != s.ProyectoID` → `ErrSprintDeOtroProyecto`; (e) Sprint `Finalizado` → `ErrSprintFinalizado`; (f) historia `Completada` (con `SprintID` cargado) → `ErrHistoriaNoPendiente`; (g) historia `Nueva` con `SprintID` cargado → `ErrHistoriaYaAsignada`; (h) historia `SprintID` igual al del Sprint destino → `ErrHistoriaYaAsignada`; precedencia (FR-009): (i) otro proyecto + Finalizado + Completada → `ErrSprintDeOtroProyecto`; (j) Finalizado + Completada → `ErrSprintFinalizado`; (k) Completada + ya asignada → `ErrHistoriaNoPendiente`; (l) historia con `Estado: domain.Estado("En progreso")` (sin constante en el dominio hasta HU-11) y `SprintID` cargado → `ErrHistoriaYaAsignada`. Comparar con `errors.Is`. Crear internal/domain/asignacion_sprint.go con los cuatro errores y `ValidarAsignacion` que devuelve `nil`. Valores exactos de data-model.md: `ErrSprintDeOtroProyecto = ValidationError{Campo: "sprint_id", Mensaje: "el Sprint no pertenece al proyecto de la historia"}`, `ErrSprintFinalizado = ValidationError{Campo: "sprint_id", Mensaje: "un Sprint finalizado no admite historias"}`, `ErrHistoriaNoPendiente = ValidationError{Campo: "estado", Mensaje: "solo se pueden planificar historias pendientes"}`, `ErrHistoriaYaAsignada = ValidationError{Campo: "sprint_id", Mensaje: "la historia ya pertenece a un Sprint"}`. Confirmar que FALLA
- [X] T009 [US1] COMMIT `RED: agregar tests de la regla de asignación de historias a un Sprint (falla)`
- [X] T010 [US1] GREEN: implementar `ValidarAsignacion` en internal/domain/asignacion_sprint.go evaluando en este orden y devolviendo el primer error: 1) `h.ProyectoID != s.ProyectoID`; 2) `s.Estado == SprintFinalizado`; 3) `h.Estado == EstadoCompletada`; 4) `h.SprintID != nil` — `go test ./...` en verde
- [X] T011 [US1] COMMIT `GREEN: implementar la regla de asignación de historias a un Sprint, test en verde`

### Persistencia: escritura condicional

- [X] T012 [US1] Declarar en internal/repository/backlog_repository.go el puerto `AsignacionHistoriaRepository { ObtenerPorID(ctx context.Context, id int64) (domain.BacklogItem, error); AsignarASprint(ctx context.Context, historiaID, sprintID int64) error }`, documentando que `AsignarASprint` solo asigna si la historia no tiene Sprint y devuelve `domain.ErrHistoriaYaAsignada` si no afectó ninguna fila. Separado de `HistoriaRepository` (ISP, research.md R3): **no** agregar el método a `HistoriaRepository`. Sin lógica
- [X] T013 [US1] COMMIT `Agregar puerto AsignacionHistoriaRepository`
- [X] T014 [US1] RED: en internal/repository/sqlite_backlog_test.go, tests de `SQLiteBacklogRepository.AsignarASprint` contra `nuevaBDMigrada(t)`, sembrando proyectos y Sprints con `NewSQLiteSprintRepository` y los helpers existentes (`guardarProyecto`, `guardarSprint`, `iniciarSprint`): (a) historia sin Sprint + Sprint del mismo proyecto → `nil`; un `ObtenerPorID` posterior tiene `SprintID` igual al Sprint, `Estado` `Nueva` y la misma `EstimacionSP` (sembrar 5 con `UPDATE` directo); (b) carrera simulada: dos llamadas seguidas sobre la misma historia con dos Sprints del proyecto, uno `Activo` (`guardarSprint` + `iniciarSprint`) y otro `Pendiente` (`guardarSprint`) — el índice único `sprints_un_pendiente_por_proyecto` de 001_esquema_base.sql impide dos `Pendiente` en el mismo proyecto (y `sprints_un_activo_por_proyecto`, dos `Activo`) → la primera `nil`, la segunda `errors.Is(err, domain.ErrHistoriaYaAsignada)`, y la historia queda en el primer Sprint; (c) Sprint de otro proyecto llamado directo (sin pasar por el dominio) → `errors.Is(err, ErrIntegridadReferencial)` y la historia sigue sin Sprint; (d) Sprint inexistente llamado directo → `errors.Is(err, ErrIntegridadReferencial)`; (e) opcional: una historia devuelta al Product Backlog por `CerrarConArrastre` (`sprint_id` NULL, estado `Nueva`) se puede asignar a otro Sprint del proyecto → `nil` y `SprintID` igual al nuevo Sprint. Stub vacío `AsignarASprint` que devuelve `nil` en internal/repository/sqlite_backlog.go. Confirmar que FALLA
- [X] T015 [US1] COMMIT `RED: agregar tests de AsignarASprint en SQLite (falla)`
- [X] T016 [US1] GREEN: implementar `AsignarASprint` en internal/repository/sqlite_backlog.go con una única sentencia `UPDATE historias SET sprint_id = ? WHERE id = ? AND sprint_id IS NULL`; error de `ExecContext` → `fmt.Errorf("asignar historia %d al sprint %d: %w", …, errorDeEscritura(err))`; `RowsAffected() == 0` → `domain.ErrHistoriaYaAsignada`. Sin transacción ni columnas nuevas — `go test ./...` en verde
- [X] T017 [US1] COMMIT `GREEN: implementar AsignarASprint en SQLite, test en verde`

### Caso de uso: AsignarHistoriaASprint

- [X] T018 [US1] RED: en internal/service/fakes_test.go agregar `fakeAsignacionRepository` (mapa de historias protegido por `sync.Mutex`; `ObtenerPorID` devuelve `domain.ErrHistoriaNoEncontrada` si no existe; `AsignarASprint` es atómico: bajo el lock, si la historia ya tiene `SprintID` devuelve `domain.ErrHistoriaYaAsignada`, si no lo carga; cuenta las llamadas a `AsignarASprint`). Crear internal/service/asignar_historia_sprint_test.go usando `repository.NewSprintRepositoryEnMemoria()` para los Sprints: (a) asignación válida (Sprint `Pendiente` y `Activo`) → devuelve la historia con `SprintID` cargado, `Estado` `Nueva` y la misma `EstimacionSP`; (b) historia inexistente → `ErrHistoriaNoEncontrada` aunque el Sprint tampoco exista (404 de historia primero, FR-008); (c) Sprint inexistente → `ErrSprintNoEncontrado`; (d) cada error de `ValidarAsignacion` (otro proyecto, Finalizado, Completada, ya asignada) se devuelve tal cual y `AsignarASprint` **no** se llama; (e) si `AsignarASprint` devuelve `ErrHistoriaYaAsignada` (carrera), el servicio lo propaga; (f) dos goroutines asignan la misma historia a la vez → exactamente una `nil` y una `ErrHistoriaYaAsignada` (correr con `go test -race ./internal/service/`). Crear internal/service/asignar_historia_sprint.go con `AsignarHistoriaASprint`, `NewAsignarHistoriaASprint(historias repository.AsignacionHistoriaRepository, sprints repository.SprintRepository)` y `Ejecutar(ctx, historiaID, sprintID int64) (domain.BacklogItem, error)` vacío. Confirmar que FALLA
- [X] T019 [US1] COMMIT `RED: agregar tests del caso de uso AsignarHistoriaASprint (falla)`
- [X] T020 [US1] GREEN: implementar `Ejecutar` en internal/service/asignar_historia_sprint.go: `historias.ObtenerPorID` → `sprints.ObtenerPorID` → `domain.ValidarAsignacion` → `historias.AsignarASprint` → devolver la historia con `SprintID = &sprintID`. El servicio solo orquesta, sin reglas propias — `go test ./...` en verde
- [X] T021 [US1] COMMIT `GREEN: implementar AsignarHistoriaASprint, test en verde`

### HTTP: PUT /backlog/{id}/sprint

- [X] T022 [US1] RED: en internal/http/backlog_handler_test.go, tests de `BacklogHandler.AsignarSprint` (ruta `PUT /backlog/{id}/sprint` en un `ServeMux` de prueba) contra los contratos de contracts/openapi.yaml: (a) 200 con la historia y `"sprint_id"` igual al enviado, `"estado": "Nueva"`; (b) historia inexistente con cuerpo `{}` y con JSON inválido → 404 `{"mensaje": "historia no encontrada"}` (la existencia domina, research.md R5); (c) JSON inválido → 400 `{"mensaje": "JSON inválido"}`; (d) `sprint_id` ausente o `null` → 400 `{"campo": "sprint_id", "mensaje": "el campo es obligatorio"}`; (e) Sprint inexistente → 404 `{"mensaje": "sprint no encontrado"}`; (f) cada rechazo de la regla → 400 con su `campo` y `mensaje` exactos (tabla de data-model.md); (g) id de ruta no numérico → 400. Agregar el parámetro `asignarSprint *service.AsignarHistoriaASprint` a `NewBacklogHandler` y actualizar su llamada en el helper de internal/http/backlog_handler_test.go y en cmd/api/main.go (pasar `nil` en main.go hasta T026 para que compile); método `AsignarSprint` vacío en internal/http/backlog_handler.go y `asignarSprintRequest { SprintID *int64 \`json:"sprint_id"\` }` en internal/http/backlog_dto.go. Confirmar que FALLA
- [X] T023 [US1] COMMIT `RED: agregar tests del handler PUT /backlog/{id}/sprint (falla)`
- [X] T024 [US1] GREEN: implementar `AsignarSprint` en internal/http/backlog_handler.go siguiendo el patrón de `Editar`: `parsearID` → `obtenerHistoria.Ejecutar` (404 antes que el cuerpo) → decodificar `asignarSprintRequest` (error → 400 "JSON inválido"; `SprintID == nil` → 400 campo `sprint_id` "el campo es obligatorio") → `asignarSprint.Ejecutar` → `escribirError` ante error / 200 con `aRespuesta`. Sin `case` nuevos en `escribirError` (los rechazos son `ValidationError`) — `go test ./...` en verde
- [X] T025 [US1] COMMIT `GREEN: implementar handler PUT /backlog/{id}/sprint, test en verde`
- [X] T026 [US1] Cablear en cmd/api/main.go: construir `service.NewAsignarHistoriaASprint(backlogRepo, sprintRepo)` con el `SQLiteBacklogRepository` y el `SQLiteSprintRepository` existentes, pasarlo a `NewBacklogHandler` (reemplazar el `nil` de T022) y registrar `mux.HandleFunc("PUT /backlog/{id}/sprint", backlogHandler.AsignarSprint)` junto a las rutas de `/backlog`
- [X] T027 [US1] COMMIT `Cablear PUT /backlog/{id}/sprint en cmd/api`

### BDD de US1 en verde

- [X] T028 [US1] GREEN: crear features/steps_asignacion.go con `InitializeScenarioAsignacion(ctx *godog.ScenarioContext, base *scenarioContext)`, que use `base.db` para construir `AsignarHistoriaASprint`, `ConsultarProductBacklog` y los services de Sprint de HU-05 (`CrearSprint`, `IniciarSprint`, `CerrarSprint`) para dejar cada Sprint en el estado que pide el escenario (sin sembrar el estado del Sprint a mano). La historia "Completada" y la "ya asignada" se siembran con SQL directo (`UPDATE historias SET sprint_id = ?, estado = 'Completada', fecha_finalizacion = '2026-10-20'`), como sqlite_sprint_test.go, porque HU-11 no existe todavía; el Sprint de otro proyecto se crea sobre un segundo proyecto. Registrar en features/features_test.go `InitializeScenarioAsignacion(ctx, sc) // HU-06`. No redefinir pasos existentes — `go test ./features/` en verde (los 9 escenarios de US1)
- [X] T029 [US1] COMMIT `GREEN: implementar los pasos BDD de asignación de historias a un Sprint, escenarios en verde`

**Checkpoint**: US1 funciona de punta a punta (dominio, SQLite, servicio, API y BDD). Es el MVP.

---

## Phase 4: User Story 2 - Consultar el Sprint Backlog de una iteración (Priority: P2)

**Goal**: `GET /sprints/{id}/backlog` devuelve las historias del Sprint, ordenadas por prioridad `M, S, C, W` y luego por `id`; `[]` si no tiene historias; 404 si el Sprint no existe.

**Independent Test**: sembrar un Sprint con historias de prioridades mezcladas y otro sin historias; la consulta del primero trae solo sus historias en orden MoSCoW y por creación, la del segundo `[]`, y la de un Sprint inexistente 404.

### Escenarios BDD de US2

- [X] T030 [US2] Agregar a features/asignar_historias_sprint.feature exactamente 3 escenarios, uno por Acceptance Scenario de US2 en spec.md: US2.1 Sprint con historias → solo las de ese Sprint, en orden M, S, C, W y por creación (tabla con título y prioridad, incluyendo historias de otro Sprint que no deben aparecer; el Sprint consultado es el `Pendiente` y el "otro Sprint" es el `Activo` del mismo proyecto, o uno de otro proyecto, porque el índice único `sprints_un_pendiente_por_proyecto` impide dos `Pendiente` en el mismo proyecto); US2.2 Sprint sin historias → Sprint Backlog vacío, sin error; US2.3 Sprint inexistente → "sprint no encontrado". Actualizar el comentario de categorías (normal: US2.1; límite: US2.2; error: US2.3). Con los pasos nuevos sin implementar, `go test ./features/` FALLA
- [X] T031 [US2] COMMIT `RED: agregar escenarios BDD de consulta del Sprint Backlog (fallan)`

### Persistencia: orden del Sprint Backlog

- [X] T032 [P] [US2] RED: en internal/repository/sqlite_backlog_test.go, test `TestSQLiteBacklogRepository_ListarPorSprint_OrdenMoSCoWYCreacion`: historias del mismo Sprint creadas con prioridades S, M, W, C, M salen M (id menor), M, S, C, W; una historia de otro Sprint no aparece (los dos Sprints son uno `Activo`, con `iniciarSprint`, y otro `Pendiente`, o el "otro Sprint" es de otro proyecto: el índice único `sprints_un_pendiente_por_proyecto` impide dos `Pendiente` en el mismo proyecto, y `sprints_un_activo_por_proyecto`, dos `Activo`); un Sprint sin historias → slice vacío no `nil`; un Sprint cerrado con `CerrarConArrastre` (una historia `Completada` y una no completada) → `ListarPorSprint` devuelve solo la completada. Con el `ORDER BY id` actual FALLA (research.md R4). Confirmar que `TestSQLiteBacklogRepository_HistoriasDeSprint` (cierre) sigue en verde
- [X] T033 [US2] COMMIT `RED: agregar test del orden MoSCoW del Sprint Backlog (falla)`
- [X] T034 [US2] GREEN: en internal/repository/sqlite_backlog.go, cambiar el `ORDER BY id` de `ListarPorSprint` por `ORDER BY CASE prioridad WHEN 'M' THEN 0 WHEN 'S' THEN 1 WHEN 'C' THEN 2 WHEN 'W' THEN 3 END, id` (MODELO-DE-DATOS.md §4.4) y actualizar su comentario ("en orden de creación" → "por prioridad M, S, C, W y luego por orden de creación") — `go test ./...` en verde
- [X] T035 [US2] COMMIT `GREEN: ordenar el Sprint Backlog por prioridad MoSCoW, test en verde`
- [X] T036 [US2] REFACTOR: extraer en internal/repository/sqlite_backlog.go la constante `ordenBacklog = "ORDER BY CASE prioridad WHEN 'M' THEN 0 WHEN 'S' THEN 1 WHEN 'C' THEN 2 WHEN 'W' THEN 3 END, id"` y usarla en `ListarSinSprint` y `ListarPorSprint` (DRY); mover a la constante el comentario sobre por qué no sirve el orden alfabético — `go test ./...` en verde
- [X] T037 [US2] COMMIT `REFACTOR: compartir el orden MoSCoW entre el Product Backlog y el Sprint Backlog`

### Caso de uso: ConsultarSprintBacklog

- [X] T038 [P] [US2] RED: crear internal/service/consultar_sprint_backlog_test.go con `repository.NewSprintRepositoryEnMemoria()` y un doble local de `repository.HistoriasDeSprintRepository` que devuelve las historias sembradas tal cual: (a) Sprint existente → devuelve las historias en el orden que entrega el repositorio; (b) Sprint sin historias (el doble devuelve `nil`) → slice vacío no `nil`; (c) Sprint inexistente → `ErrSprintNoEncontrado` y el listado no se consulta. Crear internal/service/consultar_sprint_backlog.go con `ConsultarSprintBacklog`, `NewConsultarSprintBacklog(sprints repository.SprintRepository, historias repository.HistoriasDeSprintRepository)` y `Ejecutar(ctx, sprintID int64) ([]domain.BacklogItem, error)` vacío. Confirmar que FALLA
- [X] T039 [US2] COMMIT `RED: agregar tests del caso de uso ConsultarSprintBacklog (falla)`
- [X] T040 [US2] GREEN: implementar `Ejecutar` en internal/service/consultar_sprint_backlog.go: `sprints.ObtenerPorID` (404) → `historias.ListarPorSprint` → `nil` se convierte en `make([]domain.BacklogItem, 0)`, igual que `ConsultarProductBacklog` — `go test ./...` en verde
- [X] T041 [US2] COMMIT `GREEN: implementar ConsultarSprintBacklog, test en verde`

### HTTP: GET /sprints/{id}/backlog

- [X] T042 [US2] RED: en internal/http/sprint_handler_test.go, tests de `SprintHandler.ConsultarBacklog` (ruta `GET /sprints/{id}/backlog`): (a) 200 con un arreglo de historias, cada una con `"sprint_id"` igual al id del Sprint, en el orden recibido; (b) Sprint sin historias → 200 `[]` (no `null`); (c) Sprint inexistente → 404 `{"mensaje": "sprint no encontrado"}`; (d) id no numérico → 400. Agregar el parámetro `consultarBacklog *service.ConsultarSprintBacklog` a `NewSprintHandler` y actualizar sus llamadas en internal/http/sprint_handler_test.go y cmd/api/main.go (`nil` en main.go hasta T046); método `ConsultarBacklog` vacío en internal/http/sprint_handler.go. Confirmar que FALLA
- [X] T043 [US2] COMMIT `RED: agregar tests del handler GET /sprints/{id}/backlog (falla)`
- [X] T044 [US2] GREEN: implementar `ConsultarBacklog` en internal/http/sprint_handler.go: `idDeRuta` → `consultarBacklog.Ejecutar` → `escribirError` / 200 con las historias mapeadas con `aRespuesta` (de internal/http/backlog_dto.go) — `go test ./...` en verde
- [X] T045 [US2] COMMIT `GREEN: implementar handler GET /sprints/{id}/backlog, test en verde`
- [ ] T046 [US2] Cablear en cmd/api/main.go: construir `service.NewConsultarSprintBacklog(sprintRepo, backlogRepo)`, pasarlo a `NewSprintHandler` (reemplazar el `nil` de T042) y registrar `mux.HandleFunc("GET /sprints/{id}/backlog", sprintHandler.ConsultarBacklog)` junto a las rutas de `/sprints`
- [ ] T047 [US2] COMMIT `Cablear GET /sprints/{id}/backlog en cmd/api`

### BDD de US2 en verde

- [ ] T048 [US2] GREEN: en features/steps_asignacion.go, agregar los pasos de los 3 escenarios de US2 usando `ConsultarSprintBacklog` sobre `base.db`; las historias de US2.1 se asignan con `AsignarHistoriaASprint` (no con SQL directo). No redefinir pasos existentes — `go test ./features/` en verde (12 escenarios @HU-06)
- [ ] T049 [US2] COMMIT `GREEN: implementar los pasos BDD de consulta del Sprint Backlog, escenarios en verde`

**Checkpoint**: US1 y US2 funcionan de punta a punta; el backend de HU-06 está completo.

---

## Phase 5: User Story 3 - Planificar el Sprint desde la pantalla (Priority: P3) ⛔ BLOQUEADA

**Bloqueo**: no empezar hasta que estén mergeados en `main` el PR #149 (`GET /projects/{id}/sprints`) y el #98 (listado del backlog en pantalla). Los nombres exactos de componentes y archivos de esta fase se ajustan al código que traigan esos PRs; lo que no cambia es el comportamiento pedido.

**Goal**: desde la pantalla del proyecto, elegir una historia del Product Backlog y un Sprint `Pendiente` o `Activo`, confirmar la asignación y ver el Sprint Backlog de cada Sprint, con el mensaje de la API ante un rechazo.

**Independent Test**: con la API real, asignar una historia: desaparece del Product Backlog y aparece en el Sprint Backlog del Sprint elegido sin recargar a mano; provocar un rechazo (Sprint cerrado mientras tanto) y ver el `mensaje` de la API con los listados sin cambios.

- [ ] T050 [US3] Verificar que el PR #149 y el #98 están en `main` y traer `main` a la rama con `git merge origin/main` (merge local, sin push); `cd web && npm test` y `go test ./...` en verde. Si alguno falta, detener la fase y avisar
- [ ] T051 [P] [US3] RED: crear web/src/features/sprints/api.test.ts con tests de `asignarHistoriaASprint(historiaId, sprintId)` (hace `PUT /backlog/{historiaId}/sprint` con cuerpo `{"sprint_id": sprintId}` y devuelve la historia) y `consultarSprintBacklog(sprintId)` (hace `GET /sprints/{sprintId}/backlog`), y de que un 400/404 rechaza con el `mensaje` de la respuesta, usando el cliente de web/src/compartido/api/client.ts. Exportaciones vacías en web/src/features/sprints/api.ts. Confirmar que FALLA
- [ ] T052 [US3] COMMIT `RED: agregar tests del cliente de asignación y Sprint Backlog (falla)`
- [ ] T053 [US3] GREEN: implementar `asignarHistoriaASprint` y `consultarSprintBacklog` en web/src/features/sprints/api.ts — `npm test` en verde
- [ ] T054 [US3] COMMIT `GREEN: implementar el cliente de asignación y Sprint Backlog, test en verde`
- [ ] T055 [P] [US3] RED: crear web/src/features/backlog/AsignarASprint.test.tsx para el componente `AsignarASprint` (props: la historia, los Sprints del proyecto y un callback al asignar): (a) el selector ofrece solo los Sprints `Pendiente` y `Activo` (FR-014); (b) al confirmar llama a `asignarHistoriaASprint` con la historia y el Sprint elegidos; (c) el botón de confirmar queda deshabilitado mientras la petición está en curso y no envía dos veces (como FR-016 de HU-01); (d) tras un 200 llama al callback; (e) ante un error muestra el `mensaje` de la API tal cual (FR-015) y no llama al callback; (f) sin Sprints `Pendiente`/`Activo` no permite confirmar. Componente vacío en web/src/features/backlog/AsignarASprint.tsx. Confirmar que FALLA
- [ ] T056 [US3] COMMIT `RED: agregar tests del componente AsignarASprint (falla)`
- [ ] T057 [US3] GREEN: implementar `AsignarASprint` en web/src/features/backlog/AsignarASprint.tsx — `npm test` en verde
- [ ] T058 [US3] COMMIT `GREEN: implementar el componente AsignarASprint, test en verde`
- [ ] T059 [P] [US3] RED: crear web/src/features/sprints/SprintBacklog.test.tsx para el componente `SprintBacklog` (prop: el Sprint): (a) muestra las historias de `consultarSprintBacklog` en el orden recibido; (b) Sprint sin historias → mensaje de Sprint Backlog vacío, sin error; (c) error de la API → muestra su `mensaje`; (d) vuelve a consultar cuando cambia una señal de recarga (prop o key) tras una asignación. Componente vacío en web/src/features/sprints/SprintBacklog.tsx. Confirmar que FALLA
- [ ] T060 [US3] COMMIT `RED: agregar tests del componente SprintBacklog (falla)`
- [ ] T061 [US3] GREEN: implementar `SprintBacklog` en web/src/features/sprints/SprintBacklog.tsx — `npm test` en verde
- [ ] T062 [US3] COMMIT `GREEN: implementar el componente SprintBacklog, test en verde`
- [ ] T063 [US3] RED: en web/src/features/backlog/index.test.tsx (y/o web/src/features/sprints/index.test.tsx, según cómo quede la pantalla tras #149 y #98), test de integración de la pantalla del proyecto: cada historia del Product Backlog tiene la acción "Asignar a Sprint"; cada Sprint muestra su `SprintBacklog`; tras asignar, la historia deja de verse en el Product Backlog y se ve en el Sprint Backlog del Sprint elegido; ante un rechazo, ambos listados quedan como estaban. Confirmar que FALLA
- [ ] T064 [US3] COMMIT `RED: agregar test de integración de la planificación del Sprint en pantalla (falla)`
- [ ] T065 [US3] GREEN: integrar `AsignarASprint` en el listado de web/src/features/backlog/index.tsx y `SprintBacklog` en web/src/features/sprints/SeccionSprints.tsx (o donde #149 muestre los Sprints), refrescando el Product Backlog y el Sprint Backlog afectado tras una asignación correcta; estilos en los `.css` existentes de cada feature — `npm test` en verde
- [ ] T066 [US3] COMMIT `GREEN: integrar la planificación del Sprint en la pantalla del proyecto, test en verde`
- [ ] T067 [US3] `cd web && npm run lint && npm test && npm run build` sin errores y validación manual del paso 3 de specs/018-hu-06-asignar-historias-sprint/quickstart.md

**Checkpoint**: la sub-issue de Frontend de HU-06 queda lista para cerrarse.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T068 `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida), `go test ./...` y `go test -race ./internal/...` en verde
- [ ] T069 Validación manual de los pasos 1 y 2 de specs/018-hu-06-asignar-historias-sprint/quickstart.md con la API real (`go run ./cmd/api` sobre una base nueva), anotando el resultado en esta tarea. Para cuerpos con tildes desde Git Bash en Windows, usar `curl --data-binary @archivo.json`
- [ ] T070 Revisar que contracts/openapi.yaml coincide con las respuestas reales (códigos, `campo`, `mensaje` y `sprint_id`); si hubo diferencias, corregir el contrato + COMMIT `docs: alinear el contrato de HU-06 con la implementación`

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **US1 (Phase 3)** → **US2 (Phase 4)** → **Polish (Phase 6)**.
- **US3 (Phase 5)** depende de US1 y US2 **y** de que estén en `main` el PR #149 y el #98. Si el backend se entrega antes, Polish se corre sobre US1 + US2 y la Phase 5 queda abierta (la sub-issue de Frontend no frena el cierre del issue padre, AGENTS.md).
- Dentro de US1: T006–T007 (BDD) → dominio T008–T011 → puerto T012–T013 → SQLite T014–T017 → servicio T018–T021 → handler T022–T025 → main T026–T027 → pasos BDD T028–T029.
- Dentro de US2: T030–T031 (BDD) → SQLite T032–T037 → servicio T038–T041 → handler T042–T045 → main T046–T047 → pasos BDD T048–T049.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.

## Parallel Opportunities

- **US1**: T008 (dominio) no depende de nada de US1 salvo el `.feature`; se puede escribir en paralelo con T012 (puerto). El servicio (T018) necesita el dominio (T010) y el puerto (T012); el handler (T022), el servicio.
- **US2**: T032 (orden en SQLite) y T038 (servicio, con un doble local) tocan archivos distintos y no dependen entre sí; se pueden hacer en paralelo antes de T042.
- **US3**: T051 (cliente), T055 (`AsignarASprint`) y T059 (`SprintBacklog`) son archivos distintos; los componentes pueden escribirse contra un mock del cliente en paralelo. T063 necesita los tres.
- Cada commit sigue siendo propio de su fase aunque las tareas se escriban en paralelo.

## Implementation Strategy

1. Setup + Foundational (T001–T005): la API expone `sprint_id`.
2. **MVP = US1** (T006–T029): asignar con todas las reglas y su BDD en verde. Se puede demostrar con curl (quickstart, pasos 2–12).
3. US2 (T030–T049): Sprint Backlog ordenado; completa el backend de HU-06.
4. Polish (T068–T070) sobre el backend → PR con `Closes #6` (merge commit, nunca squash), con las sub-issues de Especificación SDD, Escenarios BDD, Backend y Tests cerradas.
5. US3 (T050–T067) cuando estén #149 y #98: en la misma rama si el PR sigue abierto, o en un PR siguiente de la misma spec; cierra la sub-issue de Frontend.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Los commits del `.feature` (T007, T031) y de sus pasos (T029, T049) usan `RED:` / `GREEN:` siguiendo el precedente de las specs 009 y 010, aunque la regla 2 de la plantilla pide commit sin prefijo para escenarios.
- Nunca `git push` ni squash desde las tareas.
- [P] tasks = different files, no dependencies.

## Excepción de proceso

Entre T007 y T029 el paquete `features` falló a propósito: los 9 escenarios del `.feature` no tenían pasos (RED del BDD) y los GREEN de dominio, repositorio, servicio y handler se verificaron por paquete con `go test ./internal/... ./cmd/...`. T029 puso los 9 escenarios en verde. Los pasos siembran con SQL directo la historia "Completada" y la "ya asignada", porque HU-11 todavía no existe; se reemplazan por el servicio cuando exista. El Sprint a planificar se crea de forma diferida porque los índices únicos de `001_esquema_base.sql` impiden dos Sprints Pendiente (o dos Activo) en un proyecto. `go test -race` no se pudo correr localmente (compilador de C de 32 bits); la prueba de las dos goroutines usa un mutex.
