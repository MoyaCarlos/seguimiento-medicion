# Implementation Plan: Movimiento de Historias al Sprint Backlog (HU-06)

**Branch**: `feature/HU-06-asignar-historias-sprint` | **Date**: 2026-10-28 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/018-hu-06-asignar-historias-sprint/spec.md`

## Summary

Agregar al backend dos casos de uso sobre historias y Sprints que ya existen: **asignar una
historia del Product Backlog a un Sprint** (`PUT /backlog/{id}/sprint`) y **consultar el Sprint
Backlog** de un Sprint (`GET /sprints/{id}/backlog`). La regla de la asignación es una función
pura del dominio, `ValidarAsignacion(historia, sprint)`, que devuelve el primer error según el
orden fijado en la spec (Sprint de otro proyecto → Sprint finalizado → historia completada →
historia que ya tiene Sprint). El servicio `AsignarHistoriaASprint` solo orquesta: obtiene la
historia (404), obtiene el Sprint (404), aplica la regla y guarda. La escritura es un único
`UPDATE historias SET sprint_id = ? WHERE id = ? AND sprint_id IS NULL`: si no afecta filas,
devuelve `ErrHistoriaYaAsignada`, así una carrera entre dos asignaciones se resuelve en una sola
sentencia. El Sprint Backlog reutiliza `ListarPorSprint`, cuyo `ORDER BY` pasa a ser el mismo
de `ListarSinSprint` (MoSCoW y luego `id`) mediante una constante compartida. Sin migraciones ni
cambios en `MODELO-DE-DATOS.md`. La pantalla (US3) va en una fase propia, bloqueada hasta que
estén en `main` el PR #149 y el #98.

## Technical Context

**Language/Version**: Go 1.26.5 (backend); TypeScript + React + Vite (frontend, US3)

**Primary Dependencies**: `net/http` (routing por método y `PathValue`), `database/sql`, `encoding/json`; `modernc.org/sqlite`; Godog para BDD. Frontend: las que ya usa `web/` (Vitest + Testing Library). Sin dependencias nuevas.

**Storage**: SQLite embebido; columna `historias.sprint_id` de `001_esquema_base.sql` con la FK compuesta `(sprint_id, proyecto_id)`. Sin migración nueva.

**Testing**: `go test ./...` (dominio puro sin base, service con fakes, repositorio contra `:memory:` migrada, handlers con `httptest`); Godog en `/features`; Vitest en `web/` para US3.

**Target Platform**: Servidor local que expone REST/JSON en `:8080` para la SPA React.

**Project Type**: Aplicación web (backend Go hexagonal + frontend React/Vite).

**Performance Goals**: Mono-usuario local; respuesta percibida inmediata (< 1 s, SC-001).

**Constraints**: Sin CGO; dominio sin dependencias de infraestructura; no tocar `MODELO-DE-DATOS.md` ni agregar migraciones; asignación de a una historia; sin mover entre Sprints (D5).

**Scale/Scope**: Decenas de historias por Sprint; pocos Sprints por proyecto.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitución v1.0.0 (`.specify/memory/constitution.md`).

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First (NON-NEGOTIABLE) | `ValidarAsignacion` (una prueba por regla y por orden de precedencia), `AsignarASprint` (asignación, carrera, FK), orden del Sprint Backlog, ambos services y ambos handlers, con commits RED → GREEN (→ REFACTOR al extraer el `ORDER BY` compartido) | PASS |
| II. SDD | spec → clarify (3 Q resueltas) → plan; siguen tasks, analyze, BDD, implement | PASS |
| III. BDD | `features/asignar_historias_sprint.feature` (@HU-06) con los escenarios de US1 y US2 (normal, alternativo, límite, error) | PASS |
| IV. Clean/Hexagonal | Regla en `internal/domain` sin base; puertos en `internal/repository`; SQLite y HTTP como adaptadores | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un service por caso de uso; puerto chico para la asignación (ISP, como `HistoriasDeSprintRepository`); un solo `ORDER BY` para ambos backlogs; sin asignación múltiple ni mover entre Sprints | PASS |
| VI. Stack | Solo stdlib + `modernc.org/sqlite` + Godog; React/Vite en el frontend | PASS |

**Re-evaluación post-diseño (Phase 1)**: todos los gates siguen en PASS. Las dos diferencias con
lo acordado (puerto propio en vez de ampliar `HistoriaRepository` y RED del orden solo en
SQLite) están justificadas en [research.md](./research.md) (R3 y R4) y no agregan complejidad.

## Project Structure

### Documentation (this feature)

```text
specs/018-hu-06-asignar-historias-sprint/
├── plan.md              # Este archivo
├── research.md          # Phase 0: decisiones y alternativas
├── data-model.md        # Phase 1: enlaza a MODELO-DE-DATOS.md §4.4/§4.5
├── quickstart.md        # Phase 1: validación de punta a punta
├── contracts/
│   └── openapi.yaml     # Phase 1: PUT /backlog/{id}/sprint, GET /sprints/{id}/backlog
├── checklists/
│   └── requirements.md
└── tasks.md             # Phase 2 (/speckit-tasks; no lo crea este comando)
```

### Source Code (repository root)

```text
internal/domain/
├── asignacion_sprint.go          # NUEVO: ValidarAsignacion + errores de la asignación
└── asignacion_sprint_test.go     # NUEVO: pruebas puras, sin base

internal/repository/
├── backlog_repository.go         # + puerto AsignacionHistoriaRepository
├── sqlite_backlog.go             # + AsignarASprint; ORDER BY compartido en ListarPorSprint/ListarSinSprint
└── sqlite_backlog_test.go        # + asignación, carrera, FK, orden del Sprint Backlog

internal/service/
├── asignar_historia_sprint.go        # NUEVO: AsignarHistoriaASprint
├── asignar_historia_sprint_test.go   # NUEVO
├── consultar_sprint_backlog.go       # NUEVO: ConsultarSprintBacklog
├── consultar_sprint_backlog_test.go  # NUEVO
└── fakes_test.go                     # + fake atómico (mutex) de la asignación

internal/http/
├── backlog_handler.go            # + AsignarSprint (PUT /backlog/{id}/sprint)
├── backlog_dto.go                # + asignarSprintRequest; historiaResponse suma sprint_id
├── backlog_handler_test.go       # + tests del handler
├── sprint_handler.go             # + ConsultarBacklog (GET /sprints/{id}/backlog)
└── sprint_handler_test.go        # + tests del handler

cmd/api/main.go                   # wiring de los dos services y las dos rutas

features/
├── asignar_historias_sprint.feature  # NUEVO (@HU-06)
├── steps_asignacion.go               # NUEVO
└── features_test.go                  # registra los pasos nuevos

web/src/features/                 # US3 (bloqueada hasta #149 y #98)
├── backlog/                      # acción "Asignar a Sprint" en cada historia del backlog
└── sprints/                      # Sprint Backlog de cada Sprint + cliente de las dos rutas
```

**Structure Decision**: se respeta la arquitectura hexagonal liviana de `AGENTS.md`. No hay
paquetes nuevos: cada capa suma archivos junto a los de HU-05 y HU-15. La pantalla se integra
en las features `backlog` y `sprints` existentes de `web/`; los nombres exactos de componentes
se fijan en `tasks.md` cuando el código de #149 y #98 esté en `main`.

## Etapas de implementación

1. **US1 — Asignar (P1, MVP)**: dominio → repositorio (`AsignarASprint`) → service → handler y
   ruta → escenario BDD de asignación. Todo en verde antes de seguir.
2. **US2 — Consultar el Sprint Backlog (P2)**: RED del orden en `ListarPorSprint` → REFACTOR del
   `ORDER BY` compartido → service → handler y ruta → escenarios BDD de consulta.
3. **US3 — Pantalla (P3)**: **bloqueada** hasta que estén en `main` el PR #149
   (`GET /projects/{id}/sprints`) y el #98 (listado del backlog en pantalla). Elegir historia y
   Sprint (solo `Pendiente`/`Activo`), confirmar con el botón deshabilitado mientras la petición
   está en curso (como FR-016 de HU-01), refrescar ambos listados y mostrar el `mensaje` de la
   API ante un rechazo.

## Complexity Tracking

Sin violaciones de la constitución que justificar.
