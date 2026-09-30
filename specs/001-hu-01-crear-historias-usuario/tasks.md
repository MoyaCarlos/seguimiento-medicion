---

description: "Task list for HU-01 — Creación de Historias de Usuario"
---

# Tasks: Creación de Historias de Usuario (HU-01)

**Input**: Design documents from `/specs/001-hu-01-crear-historias-usuario/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/openapi.yaml](./contracts/openapi.yaml)

**Tests**: SÍ son obligatorias. `AGENTS.md` exige TDD (RED→GREEN→REFACTOR, con cada fase en su propio commit) para dominio/servicio y BDD (Godog) por criterio de aceptación, así que las tareas de test se incluyen y se ejecutan antes de implementar.

**Organization**: Tareas agrupadas por historia de usuario (US1=P1, US2=P2) para permitir implementación y prueba independientes.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Puede correr en paralelo (archivo distinto, sin dependencias pendientes)
- **[Story]**: `[US1]` o `[US2]`, según `spec.md`
- Rutas de archivo exactas en cada tarea
- Convención de commits (AGENTS.md): `RED: ...`, `GREEN: ...`, `REFACTOR: ...`

## Path Conventions

Backend Go en la raíz: `cmd/`, `internal/`, `features/` (ver `plan.md`). Sin frontend en esta iteración.

## Estado canónico

`Estado` inicial = `"Nueva"` (enum `EstadoNueva`). Confirmado y normalizado en todos los artefactos y en `docs/PRODUCT-BACKLOG.md` (R9 resuelto).

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Dependencias y estructura base

- [X] T001 Agregar dependencias al módulo: `go get modernc.org/sqlite` y `go get github.com/cucumber/godog`, luego `go mod tidy` (actualiza `go.mod` y `go.sum`)
- [X] T002 [P] Ignorar la base SQLite local en `.gitignore` (agregar `*.db` y `*.db-journal`)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Tipos de dominio compartidos y puerto de persistencia que usan US1 y US2

**⚠️ CRITICAL**: Ninguna historia puede empezar hasta completar esta fase

- [X] T003 [P] Implementar el enum `Estado` con `EstadoNueva = "Nueva"` en `internal/domain/estado.go`
- [X] T004 [P] Implementar el enum `Prioridad` MoSCoW con valores `"M"`, `"S"`, `"C"`, `"W"` y función de validación `EsValida()` en `internal/domain/prioridad.go`
- [X] T005 [P] Implementar el tipo `ValidationError{Campo, Mensaje string}` con `Error() string` en `internal/domain/errors.go`
- [X] T006 Crear el struct `BacklogItem` en `internal/domain/backlog_item.go` con campos: `ID int64`, `ProyectoID int64`, `Titulo string`, `Descripcion string`, `Prioridad Prioridad`, `Estado Estado`, `ValorNegocio *int`, `EstimacionSP *int` (depende de T003, T004)
- [X] T007 [P] Definir el puerto `BacklogRepository` con `Guardar(ctx context.Context, item domain.BacklogItem) (domain.BacklogItem, error)` en `internal/repository/backlog_repository.go` (depende de T006)
- [X] T008 [P] Implementar apertura de SQLite (`database/sql` + driver `"sqlite"`) y migración idempotente `CREATE TABLE IF NOT EXISTS backlog_items (id INTEGER PRIMARY KEY AUTOINCREMENT, proyecto_id INTEGER NOT NULL, titulo TEXT NOT NULL, descripcion TEXT NOT NULL, prioridad TEXT NOT NULL, estado TEXT NOT NULL, valor_negocio INTEGER, estimacion_sp INTEGER)` en `internal/repository/migrate.go` (sin `FOREIGN KEY`, ver R5)

**Checkpoint**: Dominio base, puerto y esquema listos — US1 y US2 pueden comenzar

---

## Phase 3: User Story 1 - Crear una historia de usuario (Priority: P1) 🎯 MVP

**Goal**: Un Product Builder o Scrum Master crea una historia con título, descripción y prioridad; queda con estado `"Nueva"`, estimación `nil` y se persiste; título vacío devuelve advertencia y no registra.

**Independent Test**: `POST /backlog` con datos válidos → `201` con `estado: "Nueva"` y `estimacion_sp: null`, persistido tras reiniciar; `POST /backlog` con título vacío → `400 {"campo":"titulo",...}` y sin registro nuevo.

### Tests for User Story 1 (TDD — escribir y ver FALLAR antes de implementar) ⚠️

- [X] T009 [P] [US1] Test de dominio `NewBacklogItem` en `internal/domain/backlog_item_test.go`: caso válido; `proyectoID` 0 y negativo → `ValidationError{Campo:"proyecto_id"}`; `titulo` vacío y solo espacios → `ValidationError{Campo:"titulo"}`; `titulo` > 200 caracteres → error de `titulo`; `descripcion` vacía → error de `descripcion`; `descripcion` > 2000 → error de `descripcion`; `prioridad` fuera de `{M,S,C,W}` → error de `prioridad`; `Estado == "Nueva"`; `EstimacionSP == nil`
- [X] T011 [P] [US1] Test de servicio en `internal/service/crear_historia_test.go` con un fake en memoria de `BacklogRepository` (definido en el mismo archivo): creación válida persiste y devuelve el item con `ID` asignado; entrada inválida devuelve `ValidationError` y NO llama a `Guardar`
- [X] T013 [P] [US1] Test de integración del repositorio en `internal/repository/sqlite_backlog_test.go` (SQLite `:memory:`): `Guardar` asigna `ID`, y al releer se conserva `estado = "Nueva"` y `estimacion_sp = NULL`; dos historias con el mismo `titulo` se persisten ambas (FR-011)
- [X] T015 [P] [US1] Test del handler en `internal/http/backlog_handler_test.go` (con `httptest` y fake de servicio/repo): `201` con cuerpo creado; `400` con `{"campo":"titulo",...}` cuando el título está vacío; `400` con `{"campo":"proyecto_id",...}` cuando falta el proyecto (`proyecto_id` 0 o ausente); `400` con JSON malformado; verificar forma según `contracts/openapi.yaml`

### Implementation for User Story 1

- [X] T010 [US1] Implementar el constructor Factory `NewBacklogItem(proyectoID int64, titulo, descripcion string, prioridad Prioridad, valorNegocio *int) (BacklogItem, error)` en `internal/domain/backlog_item.go`: normalizar con `strings.TrimSpace`, validar en orden `proyectoID > 0` (`ValidationError{Campo:"proyecto_id"}`), título no vacío (≤ 200 runas), descripción no vacía (≤ 2000 runas), prioridad en MoSCoW; forzar `Estado = EstadoNueva`; `EstimacionSP = nil` (esta tarea NO valida aún `valorNegocio`, ver T020) — hace pasar T009
- [X] T012 [US1] Implementar `CrearHistoriaBacklog` con `CrearHistoriaInput{ProyectoID int64, Titulo, Descripcion string, Prioridad domain.Prioridad}` y `Ejecutar(ctx, input) (domain.BacklogItem, error)` en `internal/service/crear_historia.go`: construye con `domain.NewBacklogItem` y persiste vía `BacklogRepository.Guardar` (depende de T007, T010) — hace pasar T011
- [X] T014 [US1] Implementar `SQLiteBacklogRepository.Guardar` en `internal/repository/sqlite_backlog.go`: `INSERT` con `valor_negocio`/`estimacion_sp` a `NULL` cuando el puntero es `nil`; recuperar `LastInsertId` y devolver el item con `ID` (depende de T007, T008) — hace pasar T013
- [X] T016 [US1] Implementar DTOs (request/response) en `internal/http/backlog_dto.go` y el handler `POST /backlog` en `internal/http/backlog_handler.go`: decodificar JSON, llamar al service, mapear `domain.ValidationError` (vía `errors.As`) o JSON inválido a `400 {"campo","mensaje"}`, devolver `201` con la historia creada, `500` en errores de persistencia (depende de T012) — hace pasar T015
- [X] T017 [US1] Crear el arnés BDD Godog en `features/features_test.go` (TestMain + `InitializeScenario` con `scenarioContext` y registro de pasos) y el feature `features/crear_historia_backlog.feature` con los dos escenarios de US1 (creación exitosa → estado "Nueva" y persistencia; falta título → advertencia y sin registro) más los step definitions en `features/steps_creacion.go` (depende de T012, T016)
- [X] T018 [US1] Cablear dependencias en `cmd/api/main.go`: abrir SQLite, ejecutar migración, construir `SQLiteBacklogRepository` → `CrearHistoriaBacklog` → handler, y registrar `mux.HandleFunc("POST /backlog", ...)` (depende de T014, T016)

**Checkpoint**: US1 funcional y testeable de forma independiente (MVP del backlog)

---

## Phase 4: User Story 2 - Registrar el Valor de Negocio (Priority: P2)

**Goal**: Al crear una historia se puede indicar opcionalmente el Valor de Negocio (Fibonacci `{1,2,3,5,8,13,21}`); si se omite, la historia se registra igual.

**Independent Test**: `POST /backlog` con `valor_negocio: 13` → `201` con `valor_negocio: 13`; sin el campo → `201` con `valor_negocio: null`; `valor_negocio: 4` → `400` con `{"campo":"valor_negocio",...}`.

### Tests for User Story 2 (TDD — escribir y ver FALLAR antes de implementar) ⚠️

- [X] T019 [US2] Test de dominio en `internal/domain/backlog_item_test.go`: `valorNegocio` `nil` permitido; valores `1,2,3,5,8,13,21` permitidos; `4` (fuera de la escala) → `ValidationError{Campo:"valor_negocio"}`
- [X] T021 [US2] Test de servicio en `internal/service/crear_historia_test.go`: input con `ValorNegocio` lo persiste; sin `ValorNegocio` persiste `nil`
- [X] T023 [US2] Test del handler en `internal/http/backlog_handler_test.go`: `201` con `valor_negocio` en la respuesta; `400` cuando `valor_negocio` no pertenece a la escala

### Implementation for User Story 2

- [X] T020 [US2] Extender `NewBacklogItem` en `internal/domain/backlog_item.go` con la validación de `valorNegocio`: `nil` o perteneciente a `{1,2,3,5,8,13,21}`, si no `ValidationError{Campo:"valor_negocio"}` (depende de T010) — hace pasar T019
- [X] T022 [US2] Agregar `ValorNegocio *int` a `CrearHistoriaInput` y pasarlo a `NewBacklogItem` en `internal/service/crear_historia.go` (depende de T012) — hace pasar T021
- [X] T024 [US2] Agregar `valor_negocio` al request y a la respuesta en `internal/http/backlog_dto.go` (opcional, puntero a `int`), manteniendo el mapeo de `400` (depende de T016) — hace pasar T023
- [X] T025 [US2] Agregar a `features/crear_historia_backlog.feature` y `features/steps_creacion.go` el escenario de Valor de Negocio opcional (con valor válido y omitido) (depende de T017)

**Checkpoint**: US1 y US2 funcionan de forma independiente

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Calidad, consistencia y validación end-to-end

- [X] T026 [P] Ejecutar `gofmt -l .` y `go vet ./...` y corregir hallazgos
- [X] T027 Ejecutar la suite completa `go test ./...` y el BDD `go test ./features/...`; asegurar todo en verde y que los tests de US1/US2 sean independientes
- [X] T028 Ejecutar la validación manual de `quickstart.md` (ambos escenarios curl y los casos borde) y registrar evidencia
- [X] T029 [P] Verificar que las respuestas reales cumplen `contracts/openapi.yaml` (formas `201`, `400`, `500`; enums `M/S/C/W` y `Nueva`; `estimacion_sp: null` al crear)
- [X] T030 Verificar que no queden referencias residuales a `"Nuevo"` como estado y que `"Nueva"` sea el valor canónico en `spec.md` y `docs/PRODUCT-BACKLOG.md` (R9 resuelto)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: sin dependencias
- **Foundational (Phase 2)**: depende de Setup — BLOQUEA ambas historias
- **User Stories (Phase 3-4)**: dependen de Foundational
  - US1 puede empezar primero; US2 depende de la base de US1 (mismo entity/service/DTO)
- **Polish (Phase 5)**: depende de completar las historias deseadas

### User Story Dependencies

- **US1 (P1)**: única dependencia = Foundational. Entrega el MVP.
- **US2 (P2)**: depende de US1 (extiende `NewBacklogItem`, `CrearHistoriaInput`, DTOs y el feature file); no es independiente de US1 a nivel de código, aunque aporta un incremento verificable propio.

### Within Each User Story

- Tests (RED) antes de implementación (GREEN); REFACTOR solo si hay algo que limpiar
- Dominio → service → repository → handler → BDD → wiring
- Commit por fase con prefijo `RED:`/`GREEN:`/`REFACTOR:` (AGENTS.md)

### Parallel Opportunities

- T003, T004, T005 (archivos de dominio distintos) en paralelo
- T007 y T008 en paralelo tras T006
- T009, T011, T013, T015 (archivos de test distintos) en paralelo
- T026 y T029 en paralelo en la fase final

---

## Parallel Example: User Story 1

```text
# Tests RED de US1 en paralelo (archivos distintos):
Task: "T009 [US1] Test de NewBacklogItem en internal/domain/backlog_item_test.go"
Task: "T011 [US1] Test de CrearHistoriaBacklog en internal/service/crear_historia_test.go"
Task: "T013 [US1] Test SQLite en internal/repository/sqlite_backlog_test.go"
Task: "T015 [US1] Test del handler en internal/http/backlog_handler_test.go"

# Foundation en paralelo:
Task: "T003 Estado en internal/domain/estado.go"
Task: "T004 Prioridad en internal/domain/prioridad.go"
Task: "T005 ValidationError en internal/domain/errors.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1 Setup → Phase 2 Foundational
2. Phase 3 US1 (T009–T018)
3. **STOP y VALIDAR**: `POST /backlog` válido → `201` con `"Nueva"`; título vacío → `400`
4. Demo/deploy del MVP de backlog

### Incremental Delivery

1. Setup + Foundational → base lista
2. US1 → MVP (creación con campos obligatorios)
3. US2 → Valor de Negocio opcional (incremento)
4. Polish → gofmt/vet, suite completa, conformidad de contrato, quickstart

---

## Notes

- [P] = archivos distintos, sin dependencias pendientes
- La etiqueta [Story] da trazabilidad HU → Spec → BDD → Tests → Código
- Restricciones del `data-model.md` citadas literalmente en las tareas (longitudes 200/2000, MoSCoW `M/S/C/W`, Fibonacci `1,2,3,5,8,13,21`, estado `"Nueva"`, `EstimacionSP` nulo al crear)
- `valor_negocio` (US2) y `proyecto_id` (US1) provienen de FR-008/FR-009 y FR-012; la validación de existencia del proyecto y la autorización por rol quedan diferidas a HU-04 (R5, R8)
- El listado/vista del Product Backlog y su actualización en vivo (FR-003, ex SC-006) quedan fuera de alcance y se difieren al incremento de frontend (R11); no se agrega `GET /backlog`. El orden de creación se preserva con `id` incremental
- Evitar: tareas vagas, conflictos de archivo simultáneos, dependencias cruzadas que rompan la independencia
