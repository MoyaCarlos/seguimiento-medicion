# Implementation Plan: Consulta de Sprints Anteriores (HU-12)

**Branch**: `feature/HU-12-consulta-sprints-anteriores` | **Date**: 2026-11-03 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/033-hu-12-consulta-sprints-anteriores/spec.md`

## Summary

Agregar dos consultas de solo lectura: el **historial** de Sprints finalizados de un proyecto
(`GET /projects/{id}/historial`) y el **detalle** de uno de ellos
(`GET /projects/{id}/historial/{sprintId}`), con sus historias comprometidas, sus horas por
historia y sus defectos. Las métricas no se calculan acá: salen del cálculo único de HU-03. Sin
tablas ni migraciones: un adaptador SQLite nuevo lee `historias`, `historias_arrastradas`,
`registros_horas` y `defectos`. El frontend agrega la sección "Historial de Sprints" (sub-issue
#88).

**Bloqueo de implementación**: HU-12 necesita el cálculo de métricas de HU-03, que todavía no está
en `main` (spec 032 en borrador, sin plan). Este plan se puede completar ahora, pero
`/speckit-implement` tiene que esperar el merge de HU-03 y una revisión del contrato (R2).

## Technical Context

**Language/Version**: Go 1.26 (`go.mod`); TypeScript con React + Vite en `web/`

**Primary Dependencies**: `net/http` (routing por método y `PathValue`), `database/sql` (stdlib); `modernc.org/sqlite`; Godog; React Router, Vitest y Testing Library. Sin dependencias nuevas.

**Storage**: SQLite embebido, solo lectura. Lee `proyectos`, `sprints`, `historias`, `historias_arrastradas`, `registros_horas` y `defectos` (`MODELO-DE-DATOS.md` §4.4 a §4.6, §4.8 y §4.9). Sin migraciones (§8).

**Testing**: `go test ./...` (dominio y service con fakes, repositorio contra SQLite `:memory:`, handler con `httptest`); Godog en `/features`; Vitest y Testing Library en `web/`.

**Target Platform**: Servidor que expone REST/JSON y sirve la SPA React.

**Project Type**: Aplicación web (backend Go hexagonal + frontend React/Vite).

**Performance Goals**: Respuesta percibida inmediata (< 1 s) con el volumen esperado.

**Constraints**: Sin CGO; dominio sin dependencias de infraestructura; ninguna escritura; métricas únicamente desde HU-03 (FR-012).

**Scale/Scope**: Pocos Sprints por proyecto y decenas de historias y defectos por Sprint; sin filtros ni paginación.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitución v1.1.0 (`.specify/memory/constitution.md`).

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First (NON-NEGOTIABLE) | Filtro y orden de Sprints finalizados, validación de Sprint finalizado, armado del detalle (cero horas, marcas de defectos), orden de errores en los services, consultas SQL y handlers: cada uno con commits RED → GREEN (→ REFACTOR si hay algo que limpiar) | PASS |
| II. SDD | spec → clarify (3 decisiones registradas) → plan; siguen tasks, analyze, BDD, implement | PASS |
| III. BDD | `features/historial_sprints.feature` (@HU-12) cubre lo normal (lista y detalle), lo alternativo (arrastrada completada en otro Sprint, defecto entre dos Sprints), lo límite (proyecto y Sprint vacíos) y los errores (inexistentes, Sprint no finalizado). La interfaz se cubre con pruebas de componente | PASS |
| IV. Clean/Hexagonal | Reglas puras en `internal/domain`; puerto angosto en `internal/repository`; SQLite y HTTP como adaptadores | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un service por consulta; un handler propio; las métricas no se recalculan (R1); sin filtros ni transacción de lectura (R6) | PASS |
| VI. Stack | Sin dependencias nuevas | PASS |

**Re-evaluación post-diseño (Phase 1)**: todos los gates siguen en PASS. El esquema y
`MODELO-DE-DATOS.md` no cambian. La dependencia de HU-03 no viola ningún principio; es un bloqueo
de orden, registrado arriba y en R2.

## Project Structure

### Documentation (this feature)

```text
specs/033-hu-12-consulta-sprints-anteriores/
├── plan.md              # Este archivo
├── research.md          # Fase 0: decisiones técnicas (R1–R8)
├── data-model.md        # Fase 1: vistas de lectura, puerto y services
├── quickstart.md        # Fase 1: guía de validación
├── contracts/
│   └── openapi.yaml     # Fase 1: GET historial y GET detalle
├── checklists/
│   └── requirements.md  # Calidad de la spec
└── tasks.md             # Fase 2 (/speckit.tasks)
```

### Source Code (repository root)

```text
cmd/api/main.go                              # wiring de los 2 services, el handler y las 2 rutas

internal/domain/
├── historial_sprint.go                      # SprintsFinalizados, ValidarSprintFinalizado,
│                                            # HistoriaComprometida, DefectoDeSprint, ArmarDetalleSprint
├── historial_sprint_test.go
└── errors.go                                # + ErrSprintNoFinalizado

internal/service/
├── consultar_historial_sprints.go           # ConsultarHistorialSprints (nuevo)
├── consultar_historial_sprints_test.go
├── consultar_detalle_sprint.go              # ConsultarDetalleSprint (nuevo)
├── consultar_detalle_sprint_test.go
└── fakes_test.go                            # + fakes del puerto y del cálculo de métricas

internal/repository/
├── historial_sprint_repository.go           # puerto HistorialSprintRepository
├── sqlite_historial_sprint.go               # HistoriasComprometidas, MinutosPorHistoria, DefectosDeSprint
└── sqlite_historial_sprint_test.go

internal/http/
├── historial_handler.go                     # Listar, Detalle
├── historial_dto.go                         # respuestas de lista y detalle
└── historial_handler_test.go

features/
├── historial_sprints.feature                # @HU-12 (nuevo)
├── steps_historial_sprints.go               # pasos Godog (nuevo)
└── features_test.go                         # + registro de los pasos

web/src/
├── rutas.tsx                                # + historial/*
├── navegacion.ts                            # + "Historial de Sprints"
└── features/historial/                      # index, lista, detalle, api y tipos, con sus pruebas
```

**Structure Decision**: estructura hexagonal ya establecida en `AGENTS.md`. El historial tiene su
propio handler (`HistorialHandler`), porque `SprintHandler` atiende el ciclo de vida de los Sprints
(SRP). `escribirError` no cambia: `ErrSprintNoFinalizado` es un `ValidationError` (400) y
`ErrSprintNoEncontrado` ya responde 404. El frontend es una feature nueva y no importa de otras
(`web/src/features/aislamiento.test.ts`).

## Complexity Tracking

> Sin violaciones a la constitución.

Notas de alcance (decisiones documentadas, no violaciones):

| Tema | Decisión | Justificación |
|------|----------|---------------|
| Dependencia de HU-03 | La implementación espera su merge | R1 y R2: FR-012 prohíbe recalcular métricas |
| Frontend | Incremento posterior de esta misma spec; sub-issue #88 | DoD: la sub-issue de Frontend no frena el cierre de la historia |
| Registros de horas individuales | No se devuelven | Clarification 2; siguen disponibles en `GET /backlog/{id}/horas` |
| Filtros y paginación | No se implementan | Assumptions, YAGNI |
| Transacción de lectura | No se usa | R6; KISS |
