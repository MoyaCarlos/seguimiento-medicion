---

description: "Tareas de HU-15: Consulta y Edición del Product Backlog (TDD con commits por fase)"
---

# Tasks: Consulta y Edición del Product Backlog (HU-15)

**Input**: Design documents from `/specs/009-hu-15-consulta-edicion-backlog/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica es RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT).

**Compilación en RED**: todos los commits compilan. En cada RED se agregan, junto al test, solo
las firmas nuevas que el test necesita con un cuerpo vacío (devolver valor cero / `nil`), de modo
que `go build ./...` pase y el test falle por aserción, no por símbolo indefinido.

**Organization**: por user story de la spec (US1 consultar, US2 editar).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2

---

## Phase 1: Setup

**Purpose**: dejar los escenarios BDD escritos antes de implementar (flujo de AGENTS.md).

- [ ] T001 Escribir features/consulta_edicion_backlog.feature en Gherkin español: tag `@HU-15`, Característica con el Como/Quiero/Para de la spec, Antecedentes `Dado que existe el proyecto con identificador 1`, y exactamente 7 escenarios, uno por Acceptance Scenario de spec.md (US1.1 solo historias sin Sprint ordenadas M, S, C, W y por creación; US1.2 backlog vacío; US1.3 proyecto inexistente → no encontrado; US2.1 edición válida conserva id y estado; US2.2 título en blanco rechazado sin modificar; US2.3 Story Points fuera de la escala Fibonacci rechazados sin modificar; US2.4 historia inexistente con datos inválidos → no encontrado). Indicar en un comentario del archivo qué escenario cubre cada categoría (normal, alternativo, límite, error). Sin pasos implementados: con el runner estricto, `go test ./features/` FALLA
- [ ] T002 COMMIT `RED: agregar escenarios BDD de consulta y edición del backlog (fallan)`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: ampliar la entidad y la persistencia que usan ambas historias.

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

- [ ] T003 Agregar `ErrHistoriaNoEncontrada = errors.New("historia no encontrada")` en internal/domain/errors.go y el campo `CriteriosAceptacion *string` a `BacklogItem` en internal/domain/backlog_item.go (`NewBacklogItem` lo deja en `nil`). Sin lógica
- [ ] T004 Declarar en internal/repository/backlog_repository.go los puertos `ProductBacklogRepository { ListarSinSprint(ctx, proyectoID int64) ([]domain.BacklogItem, error) }` y `HistoriaRepository { ObtenerPorID(ctx, id int64) (domain.BacklogItem, error); Actualizar(ctx, item domain.BacklogItem) error }`, documentando que `ObtenerPorID` devuelve `domain.ErrHistoriaNoEncontrada` y que `ListarSinSprint` ordena por `id`
- [ ] T005 COMMIT `Agregar criterios de aceptación a BacklogItem y los puertos de consulta y edición`
- [ ] T006 [P] RED: en internal/repository/sqlite_backlog_test.go, tests de `SQLiteBacklogRepository` contra `:memory:` con migraciones: (a) `ListarSinSprint` devuelve solo las historias del proyecto con `sprint_id IS NULL`, en orden de `id`, excluyendo las de otro proyecto y las que tienen Sprint (sembrar un Sprint y `UPDATE historias SET sprint_id=?, estado='En progreso'` directo en la base); (b) proyecto sin historias → slice vacío no nil; (c) `ObtenerPorID` inexistente → `errors.Is(err, domain.ErrHistoriaNoEncontrada)`; (d) `ObtenerPorID` lee todos los campos incluidos `criterios_aceptacion`, `estimacion_sp` y `sprint_id`; (e) `Actualizar` persiste `titulo`, `descripcion`, `criterios_aceptacion` y `estimacion_sp` (incluido volver a NULL) sin tocar `prioridad`, `estado`, `valor_negocio` ni `sprint_id`; (f) `Guardar` persiste `criterios_aceptacion` NULL. Stubs vacíos de los 3 métodos en internal/repository/sqlite_backlog.go — confirmar que FALLA
- [ ] T007 COMMIT `RED: agregar tests de ListarSinSprint, ObtenerPorID y Actualizar en SQLite (fallan)`
- [ ] T008 GREEN: implementar en internal/repository/sqlite_backlog.go `ListarSinSprint` (`WHERE proyecto_id = ? AND sprint_id IS NULL ORDER BY id`, slice inicializado con `make`), `ObtenerPorID` (`sql.ErrNoRows` → `domain.ErrHistoriaNoEncontrada`) y `Actualizar` (`UPDATE historias SET titulo=?, descripcion=?, criterios_aceptacion=?, estimacion_sp=? WHERE id=?`); incluir `criterios_aceptacion` en `Guardar` y en el `SELECT` de `ListarPorSprint`. NULL ↔ `nil` con un helper `punteroATextoSQL` análogo a `punteroAValorSQL` — `go test ./...` en verde
- [ ] T009 COMMIT `GREEN: implementar consulta y edición de historias en SQLite, tests en verde`
- [ ] T010 REFACTOR (solo si hay duplicación): extraer en internal/repository/sqlite_backlog.go un `escanearHistoria` compartido por `ListarPorSprint`, `ListarSinSprint` y `ObtenerPorID` + COMMIT `REFACTOR: escanear historias en un solo lugar`

**Checkpoint**: entidad y persistencia listas.

---

## Phase 3: User Story 1 - Consultar el Product Backlog de un proyecto (Priority: P1) 🎯 MVP

**Goal**: `GET /projects/{id}/backlog` devuelve las historias sin Sprint, ordenadas por MoSCoW y por creación.

**Independent Test**: sembrar un proyecto con historias de distintas prioridades (algunas con Sprint) y verificar que el listado trae solo las sin Sprint, en orden M, S, C, W y por id; proyecto sin historias → `[]`; proyecto inexistente → 404.

### Orden MoSCoW (dominio)

- [ ] T011 [P] [US1] RED: en internal/domain/product_backlog_test.go, tests table-driven de `Prioridad.Orden()` (M=0, S=1, C=2, W=3) y de `OrdenarBacklog(items []BacklogItem) []BacklogItem`: historias creadas en orden S, M, W, C, M terminan M(id menor), M, S, C, W; a igual prioridad gana el `ID` menor; lista vacía → vacía; no modifica el slice de entrada. Stubs vacíos en internal/domain/prioridad.go y internal/domain/product_backlog.go — confirmar que FALLA
- [ ] T012 [US1] COMMIT `RED: agregar tests del orden MoSCoW del Product Backlog (falla)`
- [ ] T013 [US1] GREEN: implementar `Prioridad.Orden()` en internal/domain/prioridad.go y `OrdenarBacklog` en internal/domain/product_backlog.go (copia + `sort.SliceStable` por `(Prioridad.Orden(), ID)`) — tests en verde
- [ ] T014 [US1] COMMIT `GREEN: implementar el orden MoSCoW del Product Backlog, test en verde`

### Caso de uso ConsultarProductBacklog

- [ ] T015 [US1] RED: en internal/service/consultar_product_backlog_test.go, tests con fakes (agregar `fakeProductBacklogRepository` en internal/service/fakes_test.go; reutilizar `fakeProjectRepository`): (a) proyecto inexistente → `domain.ErrProyectoNoEncontrado` y no se llama a `ListarSinSprint`; (b) devuelve las historias ordenadas por `OrdenarBacklog`; (c) sin historias → slice vacío no nil; (d) error del repositorio se propaga. Stub de `NewConsultarProductBacklog(proyectos repository.ProjectRepository, backlog repository.ProductBacklogRepository)` y `Ejecutar(ctx, proyectoID int64) ([]domain.BacklogItem, error)` en internal/service/consultar_product_backlog.go — confirmar que FALLA
- [ ] T016 [US1] COMMIT `RED: agregar tests de ConsultarProductBacklog (fallan)`
- [ ] T017 [US1] GREEN: implementar `ConsultarProductBacklog.Ejecutar` en internal/service/consultar_product_backlog.go: `ObtenerPorID` del proyecto → `ListarSinSprint` → `domain.OrdenarBacklog` — `go test ./...` en verde
- [ ] T018 [US1] COMMIT `GREEN: implementar ConsultarProductBacklog, tests en verde`

### Handler GET /projects/{id}/backlog

- [ ] T019 [US1] RED: en internal/http/backlog_handler_test.go, tests con `httptest` de `BacklogHandler.ConsultarBacklog`: 200 con el array en el orden recibido y cada historia con `id, proyecto_id, titulo, descripcion, criterios_aceptacion, prioridad, estado, valor_negocio, estimacion_sp` (FR-004; `criterios_aceptacion` y `estimacion_sp` como `null` si faltan); 200 con `[]` (no `null`) si está vacío; 404 `proyecto no encontrado`; 400 con id no numérico. Ajustar `NewBacklogHandler` para recibir también el service y adaptar `nuevoHandlerDePrueba`; agregar `CriteriosAceptacion *string \`json:"criterios_aceptacion"\`` a `historiaResponse` en internal/http/backlog_dto.go. Stub vacío del método — confirmar que FALLA
- [ ] T020 [US1] COMMIT `RED: agregar tests del handler GET /projects/{id}/backlog (fallan)`
- [ ] T021 [US1] GREEN: implementar `ConsultarBacklog` en internal/http/backlog_handler.go (`parsearID(r.PathValue("id"))`, mapear con `aRespuesta`, slice inicializado con `make`), completar `aRespuesta` con los criterios, y cablear `service.NewConsultarProductBacklog(proyectos, backlogRepo)` y `mux.HandleFunc("GET /projects/{id}/backlog", backlogHandler.ConsultarBacklog)` en cmd/api/main.go — `go test ./...` en verde
- [ ] T022 [US1] COMMIT `GREEN: implementar GET /projects/{id}/backlog y cablear la ruta, tests en verde`

**Checkpoint**: US1 funciona y se prueba sola (quickstart pasos 1–2).

---

## Phase 4: User Story 2 - Editar una historia del Product Backlog (Priority: P1)

**Goal**: `PUT /backlog/{id}` reemplaza título, descripción, criterios y Story Points, conservando el resto.

**Independent Test**: editar una historia existente con datos válidos y verificar persistencia con id, estado, prioridad, valor de negocio y Sprint intactos; datos inválidos → 400 con `campo` y la historia sin cambios; historia inexistente → 404 aunque los datos sean inválidos.

### Escala de Story Points (dominio)

- [ ] T023 [P] [US2] RED: en internal/domain/backlog_item_test.go, test table-driven de `EsEstimacionValida(v int) bool`: verdadero para 1, 2, 3, 5, 8, 13, 21; falso para 0, -1, 4, 22, 34. Stub en internal/domain/valor_negocio.go — confirmar que FALLA
- [ ] T024 [US2] COMMIT `RED: agregar test de la escala de Story Points (falla)`
- [ ] T025 [US2] GREEN: implementar `EsEstimacionValida` sobre el mismo mapa que `EsValorNegocioValido` (renombrar `valoresNegocioValidos` a `escalaFibonacci`) y exportar `ErrEstimacionInvalida = ValidationError{Campo: "estimacion_sp", Mensaje: "los Story Points deben ser uno de 1, 2, 3, 5, 8, 13, 21"}` en internal/domain/valor_negocio.go — tests en verde
- [ ] T026 [US2] COMMIT `GREEN: implementar la escala de Story Points, test en verde`

### Edición de la historia (dominio)

- [ ] T027 [US2] RED: en internal/domain/backlog_item_test.go, tests table-driven de `(b BacklogItem) ConDatosEditados(titulo, descripcion string, criterios *string, estimacionSP *int) (BacklogItem, error)` partiendo de una historia con ID, ProyectoID, Prioridad, Estado "En progreso", ValorNegocio y SprintID cargados: (a) válida → aplica trim y los 4 campos nuevos, conserva ID, ProyectoID, Prioridad, Estado, ValorNegocio y SprintID (FR-011); (b) título vacío o solo espacios → `ValidationError{Campo:"titulo"}` "el título es obligatorio"; (c) título de 201 runas → "el título no puede superar los 200 caracteres", 200 runas acepta; (d) descripción vacía → `descripcion` "la descripción es obligatoria"; (e) descripción de 2001 runas → "la descripción no puede superar los 2000 caracteres"; (f) criterios `nil`, `""` o `"   "` → `CriteriosAceptacion == nil`; (g) criterios de 2001 runas tras trim → `Campo:"criterios_aceptacion"` "los criterios de aceptación no pueden superar los 2000 caracteres", 2000 acepta; (h) estimación `nil` → sin estimación; (i) estimación 4 o 0 → `ErrEstimacionInvalida`; (j) varios campos inválidos → error y el receptor `b` queda sin cambios. Stub que devuelva `b, nil` en internal/domain/backlog_item.go — confirmar que FALLA
- [ ] T028 [US2] COMMIT `RED: agregar tests de edición de historias (falla)`
- [ ] T029 [US2] GREEN: implementar `ConDatosEditados` en internal/domain/backlog_item.go con el orden de validación de data-model.md (título, descripción, criterios, estimación) y la constante `longitudMaximaCriterios = 2000` — tests en verde
- [ ] T030 [US2] COMMIT `GREEN: implementar la edición de historias, test en verde`
- [ ] T031 [US2] REFACTOR: extraer `validarTitulo(string) (string, error)` y `validarDescripcion(string) (string, error)` en internal/domain/backlog_item.go y usarlos desde `NewBacklogItem` y `ConDatosEditados` (DRY con HU-01); tests de HU-01 y HU-15 en verde
- [ ] T032 [US2] COMMIT `REFACTOR: compartir la validación de título y descripción entre alta y edición`

### Caso de uso EditarHistoria

- [ ] T033 [US2] RED: en internal/service/editar_historia_test.go, tests con `fakeHistoriaRepository` (agregar en internal/service/fakes_test.go, con mapa por id y registro de llamadas a `Actualizar`): (a) historia inexistente con título vacío → `domain.ErrHistoriaNoEncontrada` (existencia antes que validación, FR-012); (b) válida → llama a `Actualizar` con la historia editada y la devuelve; (c) inválida → `ValidationError` y `Actualizar` no se llama; (d) error de `Actualizar` se propaga; (e) `ObtenerHistoria.Ejecutar(ctx, id)` devuelve la historia o `domain.ErrHistoriaNoEncontrada` (lo usa el handler para que la existencia domine, igual que `ObtenerProyecto` en `ProjectHandler.Editar`), en internal/service/obtener_historia_test.go. Stubs de `NewObtenerHistoria(historias repository.HistoriaRepository)` en internal/service/obtener_historia.go y de `EditarHistoriaInput{ID int64; Titulo, Descripcion string; CriteriosAceptacion *string; EstimacionSP *int}`, `NewEditarHistoria(historias repository.HistoriaRepository)` y `Ejecutar(ctx, input) (domain.BacklogItem, error)` en internal/service/editar_historia.go — confirmar que FALLA
- [ ] T034 [US2] COMMIT `RED: agregar tests de EditarHistoria (fallan)`
- [ ] T035 [US2] GREEN: implementar `ObtenerHistoria.Ejecutar` en internal/service/obtener_historia.go y `EditarHistoria.Ejecutar` en internal/service/editar_historia.go (`ObtenerPorID` → `ConDatosEditados` → `Actualizar`), igual que `ObtenerProyecto` / `EditarProyecto` — `go test ./...` en verde
- [ ] T036 [US2] COMMIT `GREEN: implementar ObtenerHistoria y EditarHistoria, tests en verde`

### Handler PUT /backlog/{id}

- [ ] T037 [US2] RED: en internal/http/backlog_handler_test.go, tests con `httptest` de `BacklogHandler.Editar`: (a) 200 con la historia completa; (b) 404 `historia no encontrada` aunque el cuerpo tenga título vacío, falte una clave o `estimacion_sp` sea 2.5 (existencia primero, R6); (c) 400 `campo:"titulo"` con título en blanco; (d) 400 `campo:"estimacion_sp"` con 4 y con 2.5 (mensaje de `domain.ErrEstimacionInvalida`); (e) 400 `campo:"<clave>"` "el campo es obligatorio" si falta cualquiera de `titulo`, `descripcion`, `criterios_aceptacion`, `estimacion_sp`; (f) `estimacion_sp: null` y `criterios_aceptacion: null` → 200 con ambos `null`; (g) JSON inválido → 400 "JSON inválido"; (h) id no numérico → 400. Agregar `editarHistoriaRequest{Titulo, Descripcion string; CriteriosAceptacion *string; EstimacionSP *json.Number}` en internal/http/backlog_dto.go y `ErrHistoriaNoEncontrada` al caso 404 de `escribirError` en internal/http/sprint_handler.go (cuando se necesite para el GREEN, no antes). Stub vacío del método — confirmar que FALLA
- [ ] T038 [US2] COMMIT `RED: agregar tests del handler PUT /backlog/{id} (fallan)`
- [ ] T039 [US2] GREEN: implementar `Editar` en internal/http/backlog_handler.go siguiendo `ProjectHandler.Editar`: parsear id → verificar existencia con `ObtenerHistoria.Ejecutar` (comentario: "La historia inexistente domina (404) sobre cualquier validación del cuerpo.") → leer cuerpo y exigir las 4 claves → convertir `estimacion_sp` (`json.Number` no entero → `domain.ErrEstimacionInvalida`) → `EditarHistoria.Ejecutar` → 200 `aRespuesta`. Mapear `ErrHistoriaNoEncontrada` a 404 en `escribirError`. `NewBacklogHandler` recibe también `*service.ObtenerHistoria` y `*service.EditarHistoria`. Cablear `service.NewObtenerHistoria(backlogRepo)`, `service.NewEditarHistoria(backlogRepo)` y `mux.HandleFunc("PUT /backlog/{id}", backlogHandler.Editar)` en cmd/api/main.go — `go test ./...` en verde
- [ ] T040 [US2] COMMIT `GREEN: implementar PUT /backlog/{id} y cablear la ruta, tests en verde`

**Checkpoint**: US1 y US2 funcionan (quickstart pasos 3–9).

---

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T041 GREEN: en features/steps_backlog.go, step definitions de Godog para features/consulta_edicion_backlog.feature sobre SQLite `:memory:` con migraciones y los services `ConsultarProductBacklog` / `EditarHistoria` (sembrar el Sprint de US1.1 con `UPDATE historias SET sprint_id=?, estado='En progreso'`); registrarlas en features/features_test.go — `go test ./features/` en verde
- [ ] T042 COMMIT `GREEN: implementar los pasos BDD de consulta y edición del backlog, escenarios en verde`
- [ ] T043 [P] Agregar al contrato de HU-01 (specs/001-hu-01-crear-historias-usuario/contracts/openapi.yaml) el campo `criterios_aceptacion` (nullable) en la respuesta de `POST /backlog`, con referencia a specs/009-hu-15-consulta-edicion-backlog/contracts/openapi.yaml
- [ ] T044 COMMIT `docs: sumar criterios de aceptación a la respuesta de POST /backlog (HU-15)`
- [ ] T045 `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test ./...` en verde; correr la validación manual de specs/009-hu-15-consulta-edicion-backlog/quickstart.md
- [ ] T046 Marcar las tareas completadas en specs/009-hu-15-consulta-edicion-backlog/tasks.md y pasar el Status de spec.md a `Implemented`
- [ ] T047 COMMIT `docs: marcar las tareas completadas de HU-15`

---

## Dependencies & Execution Order

- **Setup (T001–T002)** → **Foundational (T003–T010)** → **US1 (T011–T022)** → **US2 (T023–T040)** → **Polish (T041–T047)**.
- US2 no depende de US1 a nivel de código (solo comparten la fase Foundational y el `BacklogHandler`); se hace después para no pisar `backlog_handler.go` y `main.go` a la vez.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- T041 necesita US1 y US2 completas.

## Parallel Opportunities

- T006 (tests de repositorio) puede escribirse en paralelo con T011 (tests de dominio del orden): archivos distintos, sin dependencia.
- T023 (escala de SP) es independiente de toda la US1 y podría adelantarse.
- T043 (contrato HU-01) es documentación y puede hacerse en cualquier momento después de T008.
- Los commits igual son secuenciales: cada COMMIT incluye solo los archivos de su fase.

## Implementation Strategy

1. Setup + Foundational → escenarios escritos (en rojo) y persistencia lista.
2. US1 → consultar el backlog: MVP de esta historia, demostrable con el quickstart.
3. US2 → edición.
4. Polish → BDD en verde, contrato y verificación final; después PR con `Closes #<issue HU-15>` (merge commit, nunca squash).

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Ninguna tarea hace `git push` ni merge.
- Frontend fuera de alcance (sub-issue propia).
