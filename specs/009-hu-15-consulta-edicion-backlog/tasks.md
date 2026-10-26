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

- [X] T001 Escribir features/consulta_edicion_backlog.feature en Gherkin español: tag `@HU-15`, Característica con el Como/Quiero/Para de la spec, Antecedentes `Dado que existe un proyecto con identificador 1` (mismo texto que el paso de HU-01 en features/steps_creacion.go, no se redefine), y exactamente 9 escenarios, uno por Acceptance Scenario de spec.md (US1.1 solo historias sin Sprint ordenadas M, S, C, W y por creación; US1.2 backlog vacío; US1.3 proyecto inexistente → no encontrado; US2.1 edición válida conserva id y estado; US2.2 título en blanco rechazado sin modificar; US2.3 Story Points fuera de la escala Fibonacci rechazados sin modificar; US2.4 historia inexistente con datos inválidos → no encontrado; US2.5 cambiar Story Points ya cargados de una historia de un Sprint "Activo" o "Finalizado" rechazado sin modificar; US2.6 cargar Story Points de una historia sin estimar de un Sprint "Activo" o "Finalizado" se guarda). Indicar en un comentario del archivo qué escenario cubre cada categoría (normal, alternativo, límite, error). Sin pasos implementados: con el runner estricto, `go test ./features/` FALLA
- [X] T002 COMMIT `RED: agregar escenarios BDD de consulta y edición del backlog (fallan)`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: ampliar la entidad y la persistencia que usan ambas historias.

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

- [X] T003 Agregar `ErrHistoriaNoEncontrada = errors.New("historia no encontrada")` en internal/domain/errors.go y el campo `CriteriosAceptacion *string` a `BacklogItem` en internal/domain/backlog_item.go (`NewBacklogItem` lo deja en `nil`). Sin lógica
- [X] T004 Declarar en internal/repository/backlog_repository.go los puertos `ProductBacklogRepository { ListarSinSprint(ctx, proyectoID int64) ([]domain.BacklogItem, error) }` y `HistoriaRepository { ObtenerPorID(ctx, id int64) (domain.BacklogItem, error); Actualizar(ctx, item domain.BacklogItem) error }`, documentando que `ObtenerPorID` devuelve `domain.ErrHistoriaNoEncontrada` y que `ListarSinSprint` devuelve las historias sin Sprint ordenadas por prioridad `M, S, C, W` y luego por `id`
- [X] T005 COMMIT `Agregar criterios de aceptación a BacklogItem y los puertos de consulta y edición`
- [X] T006 RED: en internal/repository/sqlite_backlog_test.go, tests de `SQLiteBacklogRepository` contra `:memory:` con migraciones: (a) `ListarSinSprint` devuelve solo las historias del proyecto con `sprint_id IS NULL`, excluyendo las de otro proyecto y las que tienen Sprint (sembrar un Sprint y `UPDATE historias SET sprint_id=?, estado='En progreso'` directo en la base); (b) orden MoSCoW: historias creadas con prioridades S, M, W, C, M salen M (id menor), M, S, C, W — el orden alfabético (C, M, S, W) haría fallar el test; (c) proyecto sin historias → slice vacío no nil; (d) `ObtenerPorID` inexistente → `errors.Is(err, domain.ErrHistoriaNoEncontrada)`; (e) `ObtenerPorID` lee todos los campos incluidos `criterios_aceptacion`, `estimacion_sp` y `sprint_id`; (f) `Actualizar` persiste `titulo`, `descripcion`, `criterios_aceptacion` y `estimacion_sp` (incluido volver a NULL) sin tocar `prioridad`, `estado`, `valor_negocio` ni `sprint_id`; (g) `Guardar` persiste `criterios_aceptacion` NULL. Stubs vacíos de los 3 métodos en internal/repository/sqlite_backlog.go — confirmar que FALLA
- [X] T007 COMMIT `RED: agregar tests de ListarSinSprint, ObtenerPorID y Actualizar en SQLite (fallan)`
- [X] T008 GREEN: implementar en internal/repository/sqlite_backlog.go `ListarSinSprint` (`WHERE proyecto_id = ? AND sprint_id IS NULL ORDER BY CASE prioridad WHEN 'M' THEN 0 WHEN 'S' THEN 1 WHEN 'C' THEN 2 WHEN 'W' THEN 3 END, id`, como fija MODELO-DE-DATOS.md §4.4; slice inicializado con `make`), `ObtenerPorID` (`sql.ErrNoRows` → `domain.ErrHistoriaNoEncontrada`) y `Actualizar` (`UPDATE historias SET titulo=?, descripcion=?, criterios_aceptacion=?, estimacion_sp=? WHERE id=?`); incluir `criterios_aceptacion` en `Guardar` y en el `SELECT` de `ListarPorSprint`. NULL ↔ `nil` con un helper `punteroATextoSQL` análogo a `punteroAValorSQL` — `go test ./...` en verde
- [X] T009 COMMIT `GREEN: implementar consulta y edición de historias en SQLite, tests en verde`
- [X] T010 REFACTOR (solo si hay duplicación): extraer en internal/repository/sqlite_backlog.go un `escanearHistoria` compartido por `ListarPorSprint`, `ListarSinSprint` y `ObtenerPorID` + COMMIT `REFACTOR: escanear historias en un solo lugar`

**Checkpoint**: entidad y persistencia listas (incluido el orden del backlog).

---

## Phase 3: User Story 1 - Consultar el Product Backlog de un proyecto (Priority: P1) 🎯 MVP

**Goal**: `GET /projects/{id}/backlog` devuelve las historias sin Sprint, ordenadas por MoSCoW y por creación.

**Independent Test**: sembrar un proyecto con historias de distintas prioridades (algunas con Sprint) y verificar que el listado trae solo las sin Sprint, en orden M, S, C, W y por id; proyecto sin historias → `[]`; proyecto inexistente → 404.

### Caso de uso ConsultarProductBacklog

- [X] T011 [US1] RED: en internal/service/consultar_product_backlog_test.go, tests con fakes (agregar `fakeProductBacklogRepository` en internal/service/fakes_test.go; reutilizar `fakeProjectRepository`): (a) proyecto inexistente → `domain.ErrProyectoNoEncontrado` y no se llama a `ListarSinSprint`; (b) devuelve las historias tal como las entrega el repositorio (el orden lo resuelve la consulta, T008); (c) sin historias → slice vacío no nil; (d) error del repositorio se propaga. Stub de `NewConsultarProductBacklog(proyectos repository.ProjectRepository, backlog repository.ProductBacklogRepository)` y `Ejecutar(ctx, proyectoID int64) ([]domain.BacklogItem, error)` en internal/service/consultar_product_backlog.go — confirmar que FALLA
- [X] T012 [US1] COMMIT `RED: agregar tests de ConsultarProductBacklog (fallan)`
- [X] T013 [US1] GREEN: implementar `ConsultarProductBacklog.Ejecutar` en internal/service/consultar_product_backlog.go: `ObtenerPorID` del proyecto → `ListarSinSprint` (slice no nil) — `go test ./...` en verde
- [X] T014 [US1] COMMIT `GREEN: implementar ConsultarProductBacklog, tests en verde`

### Handler GET /projects/{id}/backlog

- [X] T015 [US1] RED: en internal/http/backlog_handler_test.go, tests con `httptest` de `BacklogHandler.ConsultarBacklog`: 200 con el array en el orden recibido y cada historia con `id, proyecto_id, titulo, descripcion, criterios_aceptacion, prioridad, estado, valor_negocio, estimacion_sp` (FR-004; `criterios_aceptacion` y `estimacion_sp` como `null` si faltan); 200 con `[]` (no `null`) si está vacío; 404 `proyecto no encontrado`; 400 con id no numérico. Ajustar `NewBacklogHandler` para recibir también el service y adaptar `nuevoHandlerDePrueba`; agregar `CriteriosAceptacion *string \`json:"criterios_aceptacion"\`` a `historiaResponse` en internal/http/backlog_dto.go. En el test existente de creación (internal/http/backlog_handler_test.go), verificar que la respuesta de `POST /backlog` incluye `criterios_aceptacion: null` (falla hasta T017). En el mismo commit RED, adaptar cmd/api/main.go a la nueva firma de `NewBacklogHandler` pasando `service.NewConsultarProductBacklog(proyectos, backlogRepo)`; la ruta se registra recién en T017. Stub vacío del método — confirmar que FALLA
- [X] T016 [US1] COMMIT `RED: agregar tests del handler GET /projects/{id}/backlog (fallan)`
- [X] T017 [US1] GREEN: implementar `ConsultarBacklog` en internal/http/backlog_handler.go (`parsearID(r.PathValue("id"))`, mapear con `aRespuesta`, slice inicializado con `make`), completar `aRespuesta` con los criterios, y registrar `mux.HandleFunc("GET /projects/{id}/backlog", backlogHandler.ConsultarBacklog)` en cmd/api/main.go — `go test ./...` en verde, incluidos los tests de HU-01
- [X] T018 [US1] COMMIT `GREEN: implementar GET /projects/{id}/backlog y cablear la ruta, tests en verde`

**Checkpoint**: US1 funciona y se prueba sola (quickstart pasos 1–2).

---

## Phase 4: User Story 2 - Editar una historia del Product Backlog (Priority: P1)

**Goal**: `PUT /backlog/{id}` reemplaza título, descripción, criterios y Story Points, conservando el resto y respetando D2b.

**Independent Test**: editar una historia existente con datos válidos y verificar persistencia con id, estado, prioridad, valor de negocio y Sprint intactos; datos inválidos → 400 con `campo` y la historia sin cambios; historia inexistente → 404 aunque los datos sean inválidos; Story Points ya cargados de una historia de un Sprint `Activo` o `Finalizado` → 400 sin cambios, y cargarlos desde vacío se permite.

### Escala de Story Points (dominio)

- [X] T019 [P] [US2] RED: en internal/domain/backlog_item_test.go, test table-driven de `EsEstimacionValida(v int) bool`: verdadero para 1, 2, 3, 5, 8, 13, 21; falso para 0, -1, 4, 22, 34. Stub en internal/domain/valor_negocio.go — confirmar que FALLA
- [X] T020 [US2] COMMIT `RED: agregar test de la escala de Story Points (falla)`
- [X] T021 [US2] GREEN: implementar `EsEstimacionValida` sobre el mismo mapa que `EsValorNegocioValido` (renombrar `valoresNegocioValidos` a `escalaFibonacci`) y exportar `ErrEstimacionInvalida = ValidationError{Campo: "estimacion_sp", Mensaje: "los Story Points deben ser uno de 1, 2, 3, 5, 8, 13, 21"}` en internal/domain/valor_negocio.go — tests en verde
- [X] T022 [US2] COMMIT `GREEN: implementar la escala de Story Points, test en verde`

### Edición de la historia (dominio)

- [X] T023 [US2] RED: en internal/domain/backlog_item_test.go, tests table-driven de `(b BacklogItem) ConDatosEditados(titulo, descripcion string, criterios *string, estimacionSP *int) (BacklogItem, error)` partiendo de una historia con ID, ProyectoID, Prioridad, Estado "En progreso", ValorNegocio y SprintID cargados: (a) válida → aplica trim y los 4 campos nuevos, conserva ID, ProyectoID, Prioridad, Estado, ValorNegocio y SprintID (FR-011); (b) título vacío o solo espacios → `ValidationError{Campo:"titulo"}` "el título es obligatorio"; (c) título de 201 runas → "el título no puede superar los 200 caracteres", 200 runas acepta; (d) descripción vacía → `descripcion` "la descripción es obligatoria"; (e) descripción de 2001 runas → "la descripción no puede superar los 2000 caracteres"; (f) criterios `nil`, `""` o `"   "` → `CriteriosAceptacion == nil`; (g) criterios de 2001 runas tras trim → `Campo:"criterios_aceptacion"` "los criterios de aceptación no pueden superar los 2000 caracteres", 2000 acepta; (h) estimación `nil` → sin estimación; (i) estimación 4 o 0 → `ErrEstimacionInvalida`; (j) varios campos inválidos → error y el receptor `b` queda sin cambios. Stub que devuelva `b, nil` en internal/domain/backlog_item.go — confirmar que FALLA
- [X] T024 [US2] COMMIT `RED: agregar tests de edición de historias (falla)`
- [X] T025 [US2] GREEN: implementar `ConDatosEditados` en internal/domain/backlog_item.go con el orden de validación de data-model.md (título, descripción, criterios, estimación) y la constante `longitudMaximaCriterios = 2000` — tests en verde
- [X] T026 [US2] COMMIT `GREEN: implementar la edición de historias, test en verde`
- [X] T027 [US2] REFACTOR: extraer `validarTitulo(string) (string, error)` y `validarDescripcion(string) (string, error)` en internal/domain/backlog_item.go y usarlos desde `NewBacklogItem` y `ConDatosEditados` (DRY con HU-01); tests de HU-01 y HU-15 en verde
- [X] T028 [US2] COMMIT `REFACTOR: compartir la validación de título y descripción entre alta y edición`

### Casos de uso ObtenerHistoria y EditarHistoria (incluye D2b)

- [X] T029 [US2] RED: tests con fakes (agregar `fakeHistoriaRepository` en internal/service/fakes_test.go, con mapa por id y registro de llamadas a `Actualizar`; para el Sprint usar `repository.NewSprintRepositoryEnMemoria()`, que ya existe). En internal/service/editar_historia_test.go: (a) historia inexistente con título vacío → `domain.ErrHistoriaNoEncontrada` (existencia antes que validación, FR-012); (b) válida sin Sprint → llama a `Actualizar` con la historia editada y la devuelve; (c) inválida → `ValidationError` y `Actualizar` no se llama; (d) error de `Actualizar` se propaga; (e) D2b, table-driven sobre el estado del Sprint asignado: `Pendiente` con SP 3 → 5 se guarda; `Activo` y `Finalizado` con SP 3 → 5 y con SP 3 → `nil` → `domain.ErrEstimacionCongelada` sin llamar a `Actualizar`; `Activo` y `Finalizado` con SP `nil` → 5 se guarda (carga inicial); `Activo` con SP 3 → 3 se guarda (no cambia) y se pueden editar título, descripción y criterios. En internal/service/obtener_historia_test.go: (f) `ObtenerHistoria.Ejecutar(ctx, id)` devuelve la historia o `domain.ErrHistoriaNoEncontrada` (lo usa el handler para que la existencia domine, igual que `ObtenerProyecto` en `ProjectHandler.Editar`). Stubs de `ErrEstimacionCongelada = ValidationError{Campo: "estimacion_sp", Mensaje: "los Story Points ya cargados no se pueden cambiar en un Sprint Activo o Finalizado"}` en internal/domain/errors.go, `NewObtenerHistoria(historias repository.HistoriaRepository)` en internal/service/obtener_historia.go y de `EditarHistoriaInput{ID int64; Titulo, Descripcion string; CriteriosAceptacion *string; EstimacionSP *int}`, `NewEditarHistoria(historias repository.HistoriaRepository, sprints repository.SprintRepository)` y `Ejecutar(ctx, input) (domain.BacklogItem, error)` en internal/service/editar_historia.go — confirmar que FALLA
- [X] T030 [US2] COMMIT `RED: agregar tests de ObtenerHistoria y EditarHistoria con la regla D2b (fallan)`
- [X] T031 [US2] GREEN: implementar `ObtenerHistoria.Ejecutar` en internal/service/obtener_historia.go y `EditarHistoria.Ejecutar` en internal/service/editar_historia.go: `ObtenerPorID` → `ConDatosEditados` → si `actual.SprintID != nil`, `sprints.ObtenerPorID(*actual.SprintID)` y, con estado `domain.SprintActivo` o `domain.SprintFinalizado`, `actual.EstimacionSP != nil` y la nueva distinta (otro valor o `nil`) → `domain.ErrEstimacionCongelada` (regla de servicio D2b, MODELO-DE-DATOS.md §5) → `Actualizar`, igual que `ObtenerProyecto` / `EditarProyecto` — `go test ./...` en verde
- [X] T032 [US2] COMMIT `GREEN: implementar ObtenerHistoria y EditarHistoria con la regla D2b, tests en verde`

### Handler PUT /backlog/{id}

- [X] T033 [US2] RED: en internal/http/backlog_handler_test.go, tests con `httptest` de `BacklogHandler.Editar`: (a) 200 con la historia completa; (b) 404 `historia no encontrada` aunque el cuerpo tenga título vacío, falte una clave o `estimacion_sp` sea 2.5 (existencia primero, R6); (c) 400 `campo:"titulo"` con título en blanco; (d) 400 `campo:"estimacion_sp"` con 4 y con 2.5 (mensaje de `domain.ErrEstimacionInvalida`); (e) 400 `campo:"<clave>"` "el campo es obligatorio" si falta cualquiera de `titulo`, `descripcion`, `criterios_aceptacion`, `estimacion_sp`; (f) `estimacion_sp: null` y `criterios_aceptacion: null` → 200 con ambos `null`; (g) 400 `campo:"estimacion_sp"` con el mensaje de `domain.ErrEstimacionCongelada` al cambiar los SP de una historia de un Sprint `Activo`; (h) JSON inválido → 400 "JSON inválido"; (i) id no numérico → 400. Agregar `editarHistoriaRequest{Titulo, Descripcion string; CriteriosAceptacion *string; EstimacionSP *json.Number}` en internal/http/backlog_dto.go. `NewBacklogHandler` recibe también `*service.ObtenerHistoria` y `*service.EditarHistoria`: en el mismo commit RED, adaptar cmd/api/main.go y `nuevoHandlerDePrueba` pasando `service.NewObtenerHistoria(backlogRepo)` y `service.NewEditarHistoria(backlogRepo, sprintRepo)`; la ruta se registra recién en T035. Stub vacío del método — confirmar que FALLA
- [X] T034 [US2] COMMIT `RED: agregar tests del handler PUT /backlog/{id} (fallan)`
- [X] T035 [US2] GREEN: implementar `Editar` en internal/http/backlog_handler.go siguiendo `ProjectHandler.Editar`: parsear id → verificar existencia con `ObtenerHistoria.Ejecutar` (comentario: "La historia inexistente domina (404) sobre cualquier validación del cuerpo.") → leer cuerpo y exigir las 4 claves → convertir `estimacion_sp` (`json.Number` no entero → `domain.ErrEstimacionInvalida`) → `EditarHistoria.Ejecutar` → 200 `aRespuesta`. Mapear `ErrHistoriaNoEncontrada` a 404 en `escribirError` (internal/http/sprint_handler.go). Registrar `mux.HandleFunc("PUT /backlog/{id}", backlogHandler.Editar)` en cmd/api/main.go — `go test ./...` en verde
- [X] T036 [US2] COMMIT `GREEN: implementar PUT /backlog/{id} y cablear la ruta, tests en verde`

**Checkpoint**: US1 y US2 funcionan (quickstart pasos 3–9).

---

## Phase 5: Polish & Cross-Cutting Concerns

- [X] T037 GREEN: en features/steps_backlog.go, step definitions de Godog para features/consulta_edicion_backlog.feature sobre SQLite `:memory:` con migraciones y los services `ConsultarProductBacklog` / `EditarHistoria` (sembrar los Sprints de US1.1, US2.5 y US2.6 con el estado que pida el escenario y `UPDATE historias SET sprint_id=?, estado='En progreso'`). Definir `InitializeScenarioBacklog(ctx *godog.ScenarioContext, base *scenarioContext)` que use `base.db` (la base que prepara el paso de proyecto de HU-01) para construir los services, igual que `InitializeScenarioEstado`, y registrarla en features/features_test.go como `InitializeScenarioBacklog(ctx, sc) // HU-15`. No redefinir pasos existentes — `go test ./features/` en verde
- [X] T038 COMMIT `GREEN: implementar los pasos BDD de consulta y edición del backlog, escenarios en verde`
- [X] T039 [P] Agregar al contrato de HU-01 (specs/001-hu-01-crear-historias-usuario/contracts/openapi.yaml) el campo `criterios_aceptacion` (nullable) en la respuesta de `POST /backlog`, con referencia a specs/009-hu-15-consulta-edicion-backlog/contracts/openapi.yaml
- [X] T040 COMMIT `docs: sumar criterios de aceptación a la respuesta de POST /backlog (HU-15)`
- [X] T041 `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test ./...` en verde; correr la validación manual de specs/009-hu-15-consulta-edicion-backlog/quickstart.md. *Hecho el 2026-10-25: verificación automática en verde (46 escenarios BDD) y validación manual del quickstart con la API real (`go run ./cmd/api` sobre una base nueva): los pasos 1–9 dieron lo esperado, incluida la persistencia tras reiniciar. Para probar textos con tildes desde Git Bash en Windows, enviar el cuerpo con `curl --data-binary @archivo.json`: con `-d` en la línea de comandos la codificación de los argumentos rompe los caracteres antes de llegar a la API.*
- [X] T042 Marcar las tareas completadas en specs/009-hu-15-consulta-edicion-backlog/tasks.md y pasar el Status de spec.md a `Implemented`
- [X] T043 COMMIT `docs: marcar las tareas completadas de HU-15`

---

## Dependencies & Execution Order

- **Setup (T001–T002)** → **Foundational (T003–T010)** → **US1 (T011–T018)** → **US2 (T019–T036)** → **Polish (T037–T043)**.
- US2 no depende de US1 a nivel de código (solo comparten la fase Foundational y el `BacklogHandler`); se hace después para no pisar `backlog_handler.go` y `main.go` a la vez.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- T037 necesita US1 y US2 completas.

## Parallel Opportunities

- T019 (escala de SP) es independiente de toda la US1 y podría adelantarse.
- T039 (contrato HU-01) es documentación y puede hacerse en cualquier momento después de T008.
- Los commits igual son secuenciales: cada COMMIT incluye solo los archivos de su fase.

## Implementation Strategy

1. Setup + Foundational → escenarios escritos (en rojo) y persistencia lista.
2. US1 → consultar el backlog: MVP de esta historia, demostrable con el quickstart.
3. US2 → edición, con la regla D2b.
4. Polish → BDD en verde, contrato y verificación final; después PR con `Closes #43` (merge commit, nunca squash).

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Ninguna tarea hace `git push` ni merge.
- Los commits del `.feature` (T002) y de sus pasos (T038) usan `RED:` / `GREEN:` siguiendo el precedente de la spec 010 (a0eceee, 6d445f9), aunque la regla 2 de la plantilla pide commit sin prefijo para escenarios.
- Las reglas de orden del backlog (CASE en SQL) y D2b (servicio) siguen MODELO-DE-DATOS.md §4.4 y §5.
- Frontend fuera de alcance (sub-issue #98).
