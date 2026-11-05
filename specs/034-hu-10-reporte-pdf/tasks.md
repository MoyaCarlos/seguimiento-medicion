---
description: "Tareas de HU-10: Generación de Reporte en PDF (TDD con commits por fase)"
---

# Tasks: Generación de Reporte en PDF (HU-10)

**Input**: Design documents from `/specs/034-hu-10-reporte-pdf/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica es
RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT).

**Compilación en RED**: todos los commits compilan. En cada RED se agregan, junto al test,
solo las firmas nuevas que el test necesita con un cuerpo vacío (devolver valor cero / `nil`),
de modo que `go build ./...` pase y el test falle por aserción, no por símbolo indefinido.

**Tests en verde durante la implementación**: `go test ./features/` falla desde T002 hasta
T038 porque los pasos del `.feature` todavía no están definidos (Godog estricto). Por eso las
tareas GREEN comprueban `go test ./internal/... ./cmd/...`, y `go test ./...` se exige desde
T040.

**Organization**: por user story de la spec (US1 Sprint, US2 proyecto, US3 proyecto sin
Sprints finalizados, US4 errores). US5 (frontend, #78) queda fuera del Sprint 4.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3 / US4

---

## Phase 1: Setup

**Purpose**: dejar los escenarios BDD escritos antes de implementar (flujo de AGENTS.md).

- [ ] T001 Escribir features/reporte_pdf.feature en Gherkin español: tag `@HU-10`,
  Característica con el Como/Quiero/Para de la spec y un escenario por cada criterio de
  aceptación de US1–US4. US1: reporte de un Sprint finalizado con historias planificadas y
  completadas, estimaciones, esfuerzo, métricas y defectos; Sprint finalizado sin horas,
  defectos o historias de una sección; Sprint Pendiente/Activo rechazado; Sprint inexistente;
  historia arrastrada planificada y no completada; defecto detectado en un Sprint y resuelto
  en otro. US2: proyecto con Sprints finalizados, Product Backlog y Sprint Activo (el Activo
  no suma); proyecto sin Sprint Activo; proyecto sin historias en el Product Backlog; proyecto
  inexistente; Sprint finalizado sin horas ni defectos; defecto detectado en un finalizado y
  resuelto en otro; historias del Product Backlog vs. Sprint Activo; varios Sprints
  finalizados en orden. US3: proyecto sin Sprints finalizados con y sin Sprint Activo
  ("sin datos"); sin historias; proyecto inexistente. US4: fallo real de generación sin
  archivo; generación exitosa con PDF válido y no vacío; fallo después de empezar sin archivo
  parcial. Indicar en un comentario qué escenario cubre cada categoría (normal, alternativo,
  límite, error). No inventar escenarios fuera de la spec. Sin pasos implementados, con el
  runner estricto `go test ./features/` FALLA
- [ ] T002 COMMIT `RED: agregar escenarios BDD de reporte en PDF (fallan)`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: persistencia de las historias arrastradas, tipos del reporte y puertos.

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

- [ ] T003 RED: en internal/repository/sqlite_backlog_test.go, `ListarArrastradasPorSprint`
  contra `:memory:`: devuelve solo las historias de `historias_arrastradas` del Sprint, con el
  título actual de la historia y la estimación congelada en la tabla (aunque la historia se
  haya editado después), ordenadas por prioridad M, S, C, W y luego por id; Sprint sin
  arrastradas → slice vacío no `nil`; no incluye historias de otro Sprint. Stub del método en
  internal/repository/sqlite_backlog.go — confirmar que FALLA
- [ ] T004 COMMIT `RED: agregar tests de las historias arrastradas de un Sprint (fallan)`
- [ ] T005 GREEN: implementar `ListarArrastradasPorSprint` en
  internal/repository/sqlite_backlog.go con `JOIN historias` y la foto de
  `historias_arrastradas.estimacion_sp` — `go test ./internal/... ./cmd/...` en verde
- [ ] T006 COMMIT `GREEN: leer las historias arrastradas de un Sprint, tests en verde`
- [ ] T007 Crear internal/domain/reporte.go con los tipos de data-model.md
  (`HistoriaReporte`, `DefectoReporte`, `SprintReporte`, `SprintActivoReporte`,
  `ReporteSprint`, `ReporteProyecto`); agregar en internal/domain/errors.go
  `ErrSprintNoFinalizado = ValidationError{Campo: "estado", Mensaje: "solo se puede generar el
  reporte de un Sprint finalizado"}` y
  `ErrReporteNoGenerado = errors.New("no se pudo generar el reporte")`, con comentario
  `(HU-10)`; crear internal/repository/reporte_repository.go con el puerto
  `HistoriasArrastradasRepository { ListarArrastradasPorSprint(ctx, sprintID int64)
  ([]domain.BacklogItem, error) }` y documentar que devuelve un slice vacío (no nil). Sin lógica
- [ ] T008 COMMIT `Agregar los tipos del reporte, sus errores y el puerto de arrastradas`

**Checkpoint**: tipos, puerto y lectura de arrastradas listos.

---

## Phase 3: User Story 1 - Reporte de un Sprint finalizado (Priority: P1) — MVP

**Goal**: `GET /sprints/{id}/reporte` devuelve un PDF de un Sprint finalizado con su detalle.

**Independent Test**: sembrar un Sprint finalizado con una historia completada, una arrastrada,
horas, métricas y defectos; pedir el reporte y verificar el agregado y el PDF; un Sprint no
finalizado da `400` y uno inexistente `404`.

- [ ] T009 [US1] RED: en internal/service/generar_reporte_test.go, con fakes de proyecto,
  Sprint, backlog, arrastradas, métricas y defectos: `EjecutarSprint` con Sprint inexistente →
  `domain.ErrSprintNoEncontrado` sin leer historias; Sprint Pendiente o Activo →
  `domain.ErrSprintNoFinalizado` sin generar; Sprint finalizado con una historia completada y
  una arrastrada → `Planificadas` con ambas, la arrastrada `Completada=false` y su estimación
  congelada; métricas del Sprint tomadas de `domain.CalcularMetricas`; defectos detectados y
  resueltos filtrados para ese Sprint; un defecto detectado en otro Sprint y resuelto en este
  solo aparece en `DefectosResueltos`; Sprint sin horas ni defectos → métricas en cero y
  slices vacíos no nil; error de un repositorio se propaga. Stub de `GenerarReporte` en
  internal/service/generar_reporte.go — confirmar que FALLA
- [ ] T010 [US1] COMMIT `RED: agregar tests del armado del reporte de Sprint (fallan)`
- [ ] T011 [US1] GREEN: implementar `GenerarReporte.EjecutarSprint` en
  internal/service/generar_reporte.go con las dependencias de repositorio de data-model.md y
  `domain.CalcularMetricas` — `go test ./internal/... ./cmd/...` en verde
- [ ] T012 [US1] COMMIT `GREEN: implementar el armado del reporte de Sprint, tests en verde`
- [ ] T013 [US1] RED: en internal/report/pdf_test.go, `GenerarSprint(domain.ReporteSprint)`
  devuelve bytes no vacíos que empiezan con `%PDF-` y contienen `%%EOF`, con datos mínimos y con
  una lista larga que fuerza paginación; `GenerarProyecto` idem con un proyecto sin Sprints
  finalizados. Stub en internal/report/pdf.go — confirmar que FALLA
- [ ] T014 [US1] COMMIT `RED: agregar tests de generación del PDF (fallan)`
- [ ] T015 [US1] GREEN: agregar `github.com/johnfercher/maroto/v2` v2.4.4 a go.mod e
  implementar en internal/report/pdf.go los generadores de Sprint y proyecto con Maroto
  (secciones de data-model.md, tablas fila por fila, paginación automática) — tests en verde
- [ ] T016 [US1] COMMIT `GREEN: generar los reportes PDF con Maroto, tests en verde`
- [ ] T017 [US1] RED: en internal/http/reporte_handler_test.go, `GET /sprints/{id}/reporte`
  con un service real sobre fakes y un generador inyectado: `200`, `Content-Type:
  application/pdf`, `Content-Disposition` con `reporte-proyecto-<proyecto>-sprint-<numero>.pdf`
  y cuerpo `%PDF-`; Sprint no finalizado → `400` campo `estado` sin cabecera PDF; Sprint
  inexistente → `404` "sprint no encontrado"; id `abc` o `0` → `400` campo `id`; error del
  generador → `500` "no se pudo generar el reporte" y sin cabecera ni cuerpo PDF. Stubs de
  `ReporteHandler` y `NewReporteHandler` en internal/http/reporte_handler.go — confirmar que
  FALLA
- [ ] T018 [US1] COMMIT `RED: agregar tests del handler GET /sprints/{id}/reporte (falla)`
- [ ] T019 [US1] GREEN: implementar `ReporteHandler.Sprint` y el helper `escribirPDF` en
  internal/http/reporte_handler.go (generar los bytes antes de escribir cabeceras); mapear
  `ErrSprintNoFinalizado` a `400` y `ErrReporteNoGenerado` a `500` en
  internal/http/sprint_handler.go — tests en verde
- [ ] T020 [US1] COMMIT `GREEN: implementar handler GET /sprints/{id}/reporte, test en verde`
- [ ] T021 [US1] Cablear en cmd/api/main.go el generador de reportes, `NewGenerarReporte`,
  `NewReporteHandler` y `mux.HandleFunc("GET /sprints/{id}/reporte", ...)`
- [ ] T022 [US1] COMMIT `Cablear GET /sprints/{id}/reporte en cmd/api`

**Checkpoint**: el reporte de un Sprint finalizado funciona de punta a punta.

---

## Phase 4: User Story 2 - Reporte del proyecto (Priority: P1)

**Goal**: `GET /projects/{id}/reporte` devuelve un PDF con los datos del proyecto, su Product
Backlog, todos sus Sprints finalizados y el Sprint Activo si existe.

**Independent Test**: un proyecto con integrantes, Product Backlog, dos Sprints finalizados y
un Sprint Activo genera un PDF con todas las secciones; los totales no incluyen el Activo; sin
Sprints finalizados, métricas/esfuerzo/defectos dicen "sin datos".

- [ ] T023 [US2] RED: en internal/service/generar_reporte_test.go, `EjecutarProyecto` con
  proyecto inexistente → `domain.ErrProyectoNoEncontrado`; proyecto con integrantes, Product
  Backlog, dos Sprints finalizados y un Sprint Activo → cada sección presente, Sprints en orden
  de `id`, Activo con sus historias y sin sumar a las métricas, esfuerzo ni defectos;
  `Estado` calculado con `domain.CalcularEstadoProyecto`; proyecto sin Sprint Activo →
  `SprintActivo=nil`; sin historias en el Product Backlog → slice vacío no nil; sin Sprints
  finalizados → `SinSprintsFinalizados=true`, métricas en cero y slices vacíos; defecto
  detectado en un finalizado y resuelto en otro aparece en cada uno; error de repositorio se
  propaga. Stub de `EjecutarProyecto` — confirmar que FALLA
- [ ] T024 [US2] COMMIT `RED: agregar tests del armado del reporte de proyecto (fallan)`
- [ ] T025 [US2] GREEN: implementar `GenerarReporte.EjecutarProyecto` reutilizando el armado de
  cada Sprint finalizado — `go test ./internal/... ./cmd/...` en verde
- [ ] T026 [US2] COMMIT `GREEN: implementar el armado del reporte de proyecto, tests en verde`
- [ ] T027 [US2] RED: en internal/http/reporte_handler_test.go, `GET /projects/{id}/reporte`:
  `200` con `Content-Type: application/pdf`, `Content-Disposition` con
  `reporte-proyecto-<id>.pdf` y cuerpo `%PDF-`; proyecto inexistente → `404` "proyecto no
  encontrado"; id inválido → `400` campo `id`; error del generador → `500` sin archivo. Stub de
  `Proyecto` en el handler — confirmar que FALLA
- [ ] T028 [US2] COMMIT `RED: agregar tests del handler GET /projects/{id}/reporte (falla)`
- [ ] T029 [US2] GREEN: implementar `ReporteHandler.Proyecto` en
  internal/http/reporte_handler.go — tests en verde
- [ ] T030 [US2] COMMIT `GREEN: implementar handler GET /projects/{id}/reporte, test en verde`
- [ ] T031 [US3] RED: en internal/http/spa_test.go, un `GET /projects/1/reporte` y
  `/sprints/1/reporte` con `Accept: text/html` llega a la API y no devuelve `index.html` —
  confirmar que FALLA
- [ ] T032 [US3] COMMIT `RED: agregar tests de la descarga de reportes a través de la SPA (fallan)`
- [ ] T033 [US3] GREEN: excluir en `esPedidoDePagina` (internal/http/spa.go) los caminos que
  terminan en `/reporte` — tests en verde
- [ ] T034 [US3] COMMIT `GREEN: permitir la descarga directa de reportes, tests en verde`
- [ ] T035 [US2] Cablear en cmd/api/main.go `mux.HandleFunc("GET /projects/{id}/reporte", ...)`
- [ ] T036 [US2] COMMIT `Cablear GET /projects/{id}/reporte en cmd/api`

**Checkpoint**: las secciones de US1, US2 y US3 funcionan de punta a punta.

---

## Phase 5: User Story 4 - Informar un error real de generación (Priority: P2)

**Goal**: un fallo real devuelve error y nunca un archivo corrupto, vacío o parcial.

- [ ] T037 [US4] Verificar con los tests de T017/T027 y T013 que: el generador se ejecuta
  completo antes de escribir cabeceras; un error del generador deja la respuesta `500` sin
  `Content-Type: application/pdf` ni cuerpo PDF; un error de lectura de datos no llega a
  escribir la respuesta. Si falta un caso, agregarlo en el RED correspondiente antes de seguir
  (no se admite un commit sin test)
- [ ] T038 [US4] (solo si hay algo que limpiar) REFACTOR de helpers repetidos de secciones o
  formato de números + COMMIT `REFACTOR: unificar el formato de las secciones del reporte`

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T039 GREEN: en features/steps_reporte.go, steps de Godog para
  features/reporte_pdf.feature sobre una base SQLite `:memory:` con migraciones, usando los
  services reales, `GenerarReporte` y los generadores de `internal/report`; para el escenario
  de error real inyectar un repositorio que falla y verificar que no se generan bytes; definir
  `InitializeScenarioReporte(ctx, ...)` y registrarla en features/features_test.go como
  `InitializeScenarioReporte(ctx) // HU-10`. No redefinir pasos existentes — `go test
  ./features/` en verde
- [ ] T040 COMMIT `GREEN: implementar los pasos BDD de reporte en PDF, escenarios en verde`
- [ ] T041 Correr `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test ./...`,
  y en web/ `npm run lint` y `npm test`, sin errores (los checks `backend` y `frontend` del CI)
- [ ] T042 Validar a mano con specs/034-hu-10-reporte-pdf/quickstart.md (pasos 1 a 9) y anotar
  el resultado al final de ese archivo
- [ ] T043 Marcar como `[X]` en este tasks.md las tareas terminadas y actualizar las casillas de
  las sub-issues #74, #75, #76 y #77 + COMMIT `docs: actualizar el avance de tasks de HU-10`

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **US1 (Phase 3)** → **US2/US3 (Phase 4)**
  → **US4 (Phase 5)** → **Polish (Phase 6)**.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT).
- US2 reutiliza el armado de Sprint de US1; por eso va después.
- US4 no agrega código nuevo: se cubre con los handlers y el generador de US1/US2.
- T039 va al final: los pasos llaman a código que tiene que existir para no romper el build.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia
  evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- `MODELO-DE-DATOS.md` no cambia: no hay tablas ni columnas nuevas.
- La dependencia nueva es solo Maroto v2.4.4, la librería que fija la constitución.
