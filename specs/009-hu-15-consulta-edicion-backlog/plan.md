# Implementation Plan: Consulta y Edición del Product Backlog (HU-15)

**Branch**: `feature/HU-15-consulta-edicion-backlog` | **Date**: 2026-10-25 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/009-hu-15-consulta-edicion-backlog/spec.md`

## Summary

Agregar al backend dos casos de uso sobre las historias de HU-01: **consultar el Product
Backlog** de un proyecto (`GET /projects/{id}/backlog`, solo historias sin Sprint, ordenadas
por MoSCoW y luego por id) y **editar una historia** (`PUT /backlog/{id}`: título, descripción,
criterios de aceptación y Story Points, todo junto). La entidad `BacklogItem` suma
`CriteriosAceptacion *string` (la columna ya existe desde #107) y la operación
`ConDatosEditados`, que reutiliza las validaciones de creación y valida los Story Points
contra la escala Fibonacci. El orden MoSCoW lo resuelve la consulta SQL con un `CASE`
explícito, como fija `MODELO-DE-DATOS.md` §4.4. `EditarHistoria` aplica además la regla de
servicio D2b: no se cambian los Story Points ya cargados de una historia de un Sprint `Activo`
o `Finalizado` (cargarlos desde vacío sí). La existencia (proyecto/historia) se verifica antes
de cualquier validación del cuerpo. El frontend queda fuera de este incremento (sub-issue
propia).

## Technical Context

**Language/Version**: Go 1.26.5

**Primary Dependencies**: `net/http` (routing por método y `PathValue`), `database/sql`, `encoding/json`; `modernc.org/sqlite`; Godog para BDD. Sin dependencias nuevas.

**Storage**: SQLite embebido; tabla `historias` de `001_esquema_base.sql`, sin migración nueva.

**Testing**: `go test ./...` (dominio y service con fakes, repositorio contra `:memory:`, handler con `httptest`); Godog en `/features`.

**Target Platform**: Servidor local que expone REST/JSON en `:8080` para la SPA React.

**Project Type**: Aplicación web (backend Go hexagonal + frontend React/Vite).

**Performance Goals**: Mono-usuario local; respuesta percibida inmediata (< 1 s).

**Constraints**: Sin CGO; dominio sin dependencias de infraestructura; sin paginación ni búsqueda (spec, Assumptions).

**Scale/Scope**: Decenas a cientos de historias por proyecto.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitución v1.0.0 (`.specify/memory/constitution.md`).

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First (NON-NEGOTIABLE) | Validación de edición, escala de SP, orden MoSCoW (test de integración del repositorio), regla D2b, services y handlers se hacen con commits RED → GREEN (→ REFACTOR al extraer las validaciones compartidas con `NewBacklogItem`) | PASS |
| II. SDD | spec → clarify (3 Q resueltas) → plan; siguen tasks, analyze, BDD, implement | PASS |
| III. BDD | `features/consulta_edicion_backlog.feature` con los 9 escenarios de la spec (normal, alternativo, límite, error) | PASS |
| IV. Clean/Hexagonal | Reglas en `internal/domain`; puertos nuevos en `internal/repository` (ISP); SQLite y HTTP como adaptadores | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un service por caso de uso (`ConsultarProductBacklog`, `EditarHistoria`); validaciones compartidas con HU-01; escala Fibonacci única; sin versionado/paginación | PASS |
| VI. Stack | Solo stdlib + `modernc.org/sqlite` + Godog | PASS |

**Re-evaluación post-diseño (Phase 1)**: todos los gates siguen en PASS. El parseo de
`estimacion_sp` como `json.Number` en el handler (R5) no duplica la regla: el mensaje y la
escala siguen en el dominio.

## Project Structure

### Documentation (this feature)

```text
specs/009-hu-15-consulta-edicion-backlog/
├── plan.md              # Este archivo
├── research.md          # Fase 0: decisiones técnicas (R1–R9)
├── data-model.md        # Fase 1: BacklogItem ampliado, vista Product Backlog, puertos
├── quickstart.md        # Fase 1: guía de validación
├── contracts/
│   └── openapi.yaml     # Fase 1: GET /projects/{id}/backlog, PUT /backlog/{id}
├── checklists/
│   └── requirements.md  # Calidad de la spec
└── tasks.md             # Fase 2 (/speckit.tasks)
```

### Source Code (repository root)

```text
cmd/api/main.go                          # wiring de los 2 services y las 2 rutas nuevas

internal/domain/
├── backlog_item.go                      # + CriteriosAceptacion, ConDatosEditados, helpers de validación compartidos
├── backlog_item_test.go                 # + tests de edición
├── valor_negocio.go                     # escala Fibonacci compartida + EsEstimacionValida, ErrEstimacionInvalida
└── errors.go                            # + ErrHistoriaNoEncontrada, ErrEstimacionCongelada

internal/service/
├── consultar_product_backlog.go         # ConsultarProductBacklog (nuevo)
├── consultar_product_backlog_test.go
├── editar_historia.go                   # EditarHistoria (nuevo), con la regla D2b (usa SprintRepository)
├── editar_historia_test.go
├── obtener_historia.go                  # ObtenerHistoria (nuevo): existencia antes de validar en el handler
├── obtener_historia_test.go
└── fakes_test.go                        # + fakeHistoriaRepository

internal/repository/
├── backlog_repository.go                # + ProductBacklogRepository, HistoriaRepository
├── sqlite_backlog.go                    # + ListarSinSprint (ORDER BY CASE prioridad, id), ObtenerPorID, Actualizar; criterios en Guardar/ListarPorSprint
└── sqlite_backlog_test.go

internal/http/
├── backlog_handler.go                   # + ConsultarBacklog, Editar
├── backlog_dto.go                       # + editarHistoriaRequest, criterios en historiaResponse
├── backlog_handler_test.go
└── sprint_handler.go                    # escribirError: ErrHistoriaNoEncontrada → 404

features/
├── consulta_edicion_backlog.feature     # @HU-15 (nuevo)
└── steps_backlog.go                     # pasos Godog (nuevo) + registro en features_test.go
```

**Structure Decision**: estructura hexagonal ya establecida en `AGENTS.md`. La consulta del
backlog se expone colgando del proyecto (`ProjectHandler` no se toca: el endpoint lo atiende
`BacklogHandler`, que ya es el adaptador de historias). Los fakes de service viven en
`fakes_test.go`, como el resto.

## Complexity Tracking

> Sin violaciones a la constitución.

Notas de alcance (decisiones documentadas, no violaciones):

| Tema | Decisión | Justificación |
|------|----------|---------------|
| Frontend (pantalla de backlog y edición) | Fuera de este incremento; sub-issue Frontend | DoD: la sub-issue de Frontend no frena el cierre de la historia |
| Asignar historias a un Sprint | No hay endpoint; los tests siembran `sprint_id` en la base o en el fake | Es HU-06 (Assumptions) |
| Editar prioridad / valor de negocio | Fuera de alcance | Clarification 3 |
| Concurrencia | Último guardado válido, sin versión | Edge case explícito de la spec |
| D2b sin endpoint para asignar Sprint | Se prueba sembrando `sprint_id` y el estado del Sprint en la base o en el fake | HU-06 todavía no existe |
