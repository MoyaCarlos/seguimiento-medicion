---

description: "Tareas de la feature 038: GET /backlog/{id} (issue #222), TDD con commits por fase"
---

# Tasks: Consultar una historia por su identificador

**Input**: `specs/038-obtener-historia-por-id/` (plan.md, spec.md, research.md, data-model.md, contracts/obtener-historia.md, quickstart.md)

**Tests**: OBLIGATORIOS (Principio I). Un commit por fase RED / GREEN / REFACTOR. Ninguna tarea hace `git push` ni merge.

**Alcance**: solo `internal/http`, `cmd/api`, `features/` y `docs/`. Sin migración, sin cambios en `internal/domain`, `internal/service` ni `MODELO-DE-DATOS.md`. No se modifica `parsearID` ni `escribirError`.

**Archivos de test**: los tests del handler van en `internal/http/backlog_handler_test.go`, reutilizando `nuevoHandlerDePrueba()` y `stubBacklogRepository`. Para inyectar `{id}` se usa `req.SetPathValue("id", ...)` (ver cómo lo hace el helper `consultarBacklog` de ese archivo).

## Phase 1: Setup

Sin tareas: no hay dependencias, migraciones ni puertos nuevos.

---

## Phase 2: Foundational

**Propósito**: que el RED compile y falle por aserción, no por error de compilación.

- [X] T001 Agregar el método vacío `func (h *BacklogHandler) Obtener(w nethttp.ResponseWriter, r *nethttp.Request) {}` con el comentario `// Obtener atiende GET /backlog/{id}.` en internal/http/backlog_handler.go (sin lógica); confirmar `go build ./...`
- [X] T002 COMMIT `Agregar método Obtener vacío en BacklogHandler`

**Checkpoint**: el paquete compila; el método todavía no responde nada (el recorder queda con 200 y cuerpo vacío).

---

## Phase 3: User Story 1 - Consultar una historia del Product Backlog (Priority: P1) 🎯 MVP

**Goal**: `GET /backlog/{id}` devuelve 200 con la historia completa y `sprint_id` en `null`.

**Independent Test**: historia sin Sprint en el stub; la respuesta trae todos los campos y `sprint_id: null`.

- [X] T003 [US1] RED: agregar `consultarHistoria(handler, id string)` (helper que arma `GET /backlog/{id}` con `SetPathValue` y llama a `handler.Obtener`) y `TestBacklogHandler_Obtener_SinSprint` en internal/http/backlog_handler_test.go: espera 200, `Content-Type: application/json` y los 11 campos de `historiaResponse` con los valores cargados, `sprint_id` y `fecha_finalizacion` en `null`; además, comprobar (FR-003, SC-001) que el cuerpo del 200 es idéntico, campo por campo y formato por formato, al elemento que `ConsultarBacklog` (`GET /projects/{id}/backlog`, con la misma historia en `repo.backlog`) devuelve para esa historia — decodificar ambos JSON a `map[string]any` y compararlos con `reflect.DeepEqual`, y comparar también los bytes del objeto con `json.Marshal(aRespuesta(item))` — correr `go test ./internal/http/ -run Obtener` y confirmar que FALLA por aserción (no por compilación)
- [X] T004 [US1] COMMIT `RED: agregar test del handler GET /backlog/{id} sin Sprint (falla)` (solo el archivo de test)
- [X] T005 [US1] GREEN: implementar `Obtener` en internal/http/backlog_handler.go con lo mínimo: `strconv.ParseInt` del `PathValue("id")`, `h.obtenerHistoria.Ejecutar` y `escribirJSON(w, nethttp.StatusOK, aRespuesta(item))` — `go test ./...` en verde
- [X] T006 [US1] COMMIT `GREEN: implementar handler GET /backlog/{id}, test en verde`
- [X] T007 [US1] Cablear `mux.HandleFunc("GET /backlog/{id}", backlogHandler.Obtener)` en cmd/api/main.go, justo después de `PUT /backlog/{id}`; confirmar `go build ./...`
- [X] T008 [US1] COMMIT `Cablear GET /backlog/{id} en cmd/api`

**Checkpoint**: la historia sin Sprint se consulta por id.

---

## Phase 4: User Story 2 - Historia asignada a un Sprint o completada (Priority: P2)

**Goal**: la respuesta refleja `sprint_id` y `fecha_finalizacion` actuales (FR-004, FR-005).

**Independent Test**: una historia con Sprint y otra completada con fecha.

- [X] T009 [US2] Agregar en internal/http/backlog_handler_test.go `TestBacklogHandler_Obtener_ConSprint` (200, `sprint_id` = id del Sprint) y `TestBacklogHandler_Obtener_Completada` (200, `estado` de completada y `fecha_finalizacion` con formato `AAAA-MM-DD`, el de `time.DateOnly` que usa `aRespuesta`) y `TestBacklogHandler_Obtener_CompletadaDeOtroProyecto` (historia completada con `proyecto_id` distinto del único proyecto del stub, p. ej. 2: responde 200 igual, porque la consulta no filtra por proyecto); en los tres casos comprobar que el cuerpo es idéntico, campo por campo y formato por formato, a `json.Marshal(aRespuesta(historia))`, la misma función que arman el listado del backlog y `GET /sprints/{id}/backlog` (FR-003, SC-001) — correr `go test ./internal/http/ -run Obtener`
- [X] T010 [US2] COMMIT `Agregar tests del handler GET /backlog/{id} con Sprint y completada` (**nota**: `aRespuesta` ya cubre estos casos tras T005, así que probablemente pasan de entrada; son tests de caracterización y por eso el commit no lleva prefijo `RED:`. Si alguno falla, tratarlo como RED, usar el prefijo `RED:` y corregir en una tarea GREEN con su commit antes de seguir)

**Checkpoint**: US1 y US2 en verde.

---

## Phase 5: User Story 3 - Error claro ante id inexistente o inválido (Priority: P2)

**Goal**: 404 `historia no encontrada` para inexistente; 400 con campo `id` y `identificador inválido` para id no válido (FR-006, FR-007).

**Independent Test**: id inexistente y los ids `abc`, `0`, `-3`, `1.5`, fuera de rango.

- [X] T011 [US3] RED: agregar en internal/http/backlog_handler_test.go `TestBacklogHandler_Obtener_NoExiste` (404, `{"mensaje":"historia no encontrada"}`), `TestBacklogHandler_Obtener_IDInvalido` (tabla con `abc`, `0`, `-3`, `1.5` y `99999999999999999999`, que supera el máximo de `int64` (9223372036854775807): 400 con `campo: "id"` y `mensaje: "identificador inválido"`, nunca 404 ni 500; para cada fila verificar que NO se consultó el repositorio: agregar al `stubBacklogRepository` de ese archivo un campo `llamadasObtenerPorID int` que `ObtenerPorID` incrementa, y asertar que queda en 0 tras la petición con id inválido) y `TestBacklogHandler_Obtener_CerosAIzquierda` (con la historia 7 cargada en el stub, `GET /backlog/007` responde 200 con `"id": 7`, igual que `/backlog/7`, y `llamadasObtenerPorID` queda en 1; `007` no es inválido ni se trata como octal) — correr `go test ./internal/http/ -run Obtener` y confirmar que FALLA por aserción
- [X] T012 [US3] COMMIT `RED: agregar tests de error del handler GET /backlog/{id} (falla)` (solo el archivo de test)
- [X] T013 [US3] GREEN: en `Obtener` (internal/http/backlog_handler.go) reemplazar el `ParseInt` por `id, ok := idDeRuta(w, r); if !ok { return }` y enviar el error del caso de uso con `escribirError(w, err)`; sin `case` nuevos ni cambios en `parsearID`/`escribirError` — `go test ./...` en verde
- [X] T014 [US3] COMMIT `GREEN: validar id y traducir errores en handler GET /backlog/{id}, test en verde`
- [X] T015 [US3] (solo si hay algo que limpiar en `Obtener` o en los tests nuevos) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`, con los tests en verde — sin cambios que limpiar (no hubo REFACTOR, no se hizo commit)

**Checkpoint**: las tres historias en verde.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T016 [P] Crear features/obtener_historia.feature (`# language: es` y tag `@Issue-222`, en el estilo de features/consultar_sprints.feature) con los cinco escenarios de la spec: historia del Product Backlog (200, `sprint_id` null), historia en un Sprint (200, `sprint_id` del Sprint), historia completada (200, `fecha_finalizacion`), historia inexistente (404 `historia no encontrada`) e identificador inválido (400, campo `id`, `identificador inválido`)
- [X] T017 [P] Crear features/steps_obtener_historia.go con `InitializeScenarioObtenerHistoria(ctx *godog.ScenarioContext)`: SQLite en memoria y migraciones como en features/steps_consultar_sprints.go, `NewBacklogHandler` con `service.NewObtenerHistoria` (resto `nil`) y `mux.HandleFunc("GET /backlog/{id}", handler.Obtener)`; **frases propias** para no chocar con los pasos ya registrados (`el sistema responde ...` ya existe): p. ej. `^consulto la historia "([^"]*)"$`, `^la consulta de la historia responde (\d+)$`, `^la historia recibida tiene sprint_id nulo$`, etc.; los datos se insertan por SQL y el JSON recibido se compara campo a campo
- [X] T018 Registrar `InitializeScenarioObtenerHistoria(ctx) // Issue #222` en features/features_test.go; correr `go test ./features/` y confirmar los 5 escenarios nuevos en verde y los existentes sin cambios (si Godog reporta un paso ambiguo, renombrar la frase nueva)
- [X] T019 COMMIT `Agregar escenarios y pasos BDD de GET /backlog/{id}`
- [X] T020 [P] Agregar la sección "Consultar una historia" en docs/curl-frontend-sprint-2.md, **versión bash**, antes de "Cambiar el estado de una historia": `curl http://localhost:8080/backlog/1` con la descripción de `sprint_id`/`fecha_finalizacion`, el 404 `{"mensaje":"historia no encontrada"}` y el 400 `{"campo":"id","mensaje":"identificador inválido"}`
- [X] T021 Agregar la misma sección en la **versión cmd** de docs/curl-frontend-sprint-2.md (bloques ```` ```bat ````), con los mismos ejemplos
- [X] T022 COMMIT `Documentar GET /backlog/{id} en la guía de curl`
- [X] T023 `go build ./...`, `go vet ./...` y `go test ./...` sin errores
- [X] T024 Ejecutar la validación de quickstart.md (API local con `go run ./cmd/api` y los `curl` de la guía) y comprobar los resultados esperados — Validado a mano con la API real: GET /backlog/2 (sin Sprint) devuelve 200 con sprint_id y fecha_finalizacion nulos; GET /backlog/1 (en un Sprint, completada) devuelve 200 con sprint_id 1 y la fecha de finalización; /backlog/9999 da 404 historia no encontrada; /backlog/abc y /backlog/0 dan 400 con campo id e 'identificador inválido'; la forma coincide con la del listado del backlog

---

## Dependencies & Execution Order

- Foundational (T001–T002) → US1 (T003–T008) → US2 (T009–T010) → US3 (T011–T015) → Polish.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- US2 y US3 dependen de que `Obtener` exista (T005) y se escriben sobre el mismo archivo de test, por lo que van en secuencia.
- Polish: T016 y T017 son [P] (archivos distintos) pero T018 necesita ambos; T020 es [P] respecto de los de BDD; T021 sigue a T020 (mismo archivo).

### Parallel Example

```text
T016 features/obtener_historia.feature   ||  T017 features/steps_obtener_historia.go  ||  T020 docs/curl-frontend-sprint-2.md
```

## Implementation Strategy

1. **MVP = US1** (T001–T008): endpoint cableado que devuelve la historia.
2. US2: confirma `sprint_id` y `fecha_finalizacion`.
3. US3: agrega el 404 y el 400 con `idDeRuta` y `escribirError`.
4. Polish: BDD, guía de curl y verificación final.

## Notes

- Cobertura de la spec: FR-001–FR-003 → T003/T005; FR-004/FR-005 → T009; FR-006/FR-007 → T011/T013; FR-008 y FR-009 por construcción (lectura sin sesión); FR-010 → T005/T013 reutilizan `ObtenerHistoria`; SC-003 → T016–T018; SC-005 → T020–T021.
- Los COMMIT no se marcan hasta correr `git commit`. Mensajes de commit con la atribución `Co-Authored-By` indicada por la sesión.
