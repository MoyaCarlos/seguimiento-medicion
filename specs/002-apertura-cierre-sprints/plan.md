# Implementation Plan: Apertura y Cierre de Sprints

**Branch**: `002-apertura-cierre-sprints` | **Date**: 2026-10-01 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-apertura-cierre-sprints/spec.md`

## Summary

Permitir crear, iniciar y cerrar Sprints asociados a un Project, con la
regla de negocio de "un solo Sprint Activo (y un solo Pendiente) por
proyecto a la vez", y el arrastre automático de historias no completadas
al Product Backlog al cerrar. Se implementa en Go siguiendo la arquitectura
Clean/Hexagonal ya establecida (mismo patrón que HU-01), usando un fake en
memoria de `ProjectRepository` (ya commiteado) para no depender de la
implementación real de HU-04 todavía.

## Technical Context

**Language/Version**: Go 1.26

**Primary Dependencies**: `modernc.org/sqlite` (persistencia, ya en uso por
HU-01); stdlib `net/http` para el handler REST (framework HTTP definitivo
del proyecto sigue sin confirmar formalmente — se usa stdlib por default,
consistente con `cmd/api/main.go` actual).

**Storage**: SQLite embebido (`modernc.org/sqlite`), misma base que HU-01
(`seguimiento.db`), tabla nueva `sprints` + columna de asociación en
`backlog_items` para el Sprint asignado.

**Testing**: paquete `testing` estándar de Go (TDD, ciclo RED/GREEN/REFACTOR)
+ Godog (BDD) para el escenario Gherkin.

**Target Platform**: servidor backend (API REST), multiplataforma vía Go.

**Project Type**: web-service (backend del monorepo; el frontend React en
`/web` no forma parte del alcance de esta historia).

**Performance Goals**: no especificado por la consigna — prioridad en
correctitud de las reglas de negocio (un solo Activo/Pendiente, arrastre de
historias) por sobre cualquier meta de throughput.

**Constraints**: debe integrar con el contrato `Project`/`ProjectRepository`
ya commiteado (`internal/domain/project.go`,
`internal/repository/project_repository.go`); el dominio no puede depender
de SQLite ni HTTP (Principio IV de la constitución).

**Scale/Scope**: acotado a esta historia — entidad `Sprint` y sus 3
operaciones (crear, iniciar, cerrar). No incluye consulta de historial de
Sprints (eso es HU-12) ni cálculo de métricas (HU-03).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Cumplimiento |
|---|---|
| I. Test-First (RED/GREEN/REFACTOR, commits por fase) | PASS — se aplica igual que en HU-01, detallado en tasks.md |
| II. SDD (flujo Spec Kit completo) | PASS — este plan es parte del flujo |
| III. BDD (Godog, `/features`) | PASS — escenario se genera después de `/speckit.tasks`/`/speckit.analyze`, antes de `/speckit.implement` |
| IV. Clean/Hexagonal (domain sin dependencias externas) | PASS — `Sprint` en `internal/domain`, `SprintRepository` como puerto, sin SQLite/HTTP en el dominio |
| V. SOLID/KISS/YAGNI/DRY | PASS — un service por caso de uso (3 services separados: `CrearSprint`/`IniciarSprint`/`CerrarSprint`, no uno agrupado — corregido en `/speckit.analyze`, ver C1), reutiliza el patrón Factory ya usado en `Project`/`BacklogItem` |
| VI. Stack obligatorio | PASS — Go + `modernc.org/sqlite`, sin dependencias nuevas |

Sin violaciones — no hace falta llenar Complexity Tracking.

## Project Structure

### Documentation (this feature)

```text
specs/002-apertura-cierre-sprints/
├── plan.md              # Este archivo
├── research.md          # Fase 0
├── data-model.md        # Fase 1
├── quickstart.md        # Fase 1
├── contracts/           # Fase 1
└── tasks.md             # Fase 2 (/speckit.tasks, no este comando)
```

### Source Code (repository root)

```text
internal/domain/
├── sprint.go            # Entidad Sprint + Factory NewSprint
├── estado_sprint.go      # Enum EstadoSprint (Pendiente/Activo/Finalizado)
└── sprint_test.go        # Tests TDD del dominio

internal/service/
├── crear_sprint.go        # Caso de uso: Crear (un service por caso de uso, igual que HU-01)
├── crear_sprint_test.go
├── iniciar_sprint.go      # Caso de uso: Iniciar
├── iniciar_sprint_test.go
├── cerrar_sprint.go       # Caso de uso: Cerrar
└── cerrar_sprint_test.go

internal/repository/
├── sprint_repository.go       # Puerto SprintRepository
├── sqlite_sprint.go           # Adaptador SQLite
├── sqlite_sprint_test.go
└── sprint_repository_fake.go  # Fake en memoria (para tests propios y de quien consuma SprintRepository)

internal/http/
├── sprint_dto.go
├── sprint_handler.go     # POST /sprints, POST /sprints/{id}/iniciar, POST /sprints/{id}/cerrar
└── sprint_handler_test.go

features/
└── apertura_cierre_sprints.feature   # Generado después de /speckit.tasks + /speckit.analyze
```

**Structure Decision**: misma estructura de capas ya usada en HU-01, un
archivo por responsabilidad dentro de cada paquete. Se agrega
`sprint_repository_fake.go` versionado (no solo en tests) porque HU-13
también lo va a necesitar para su propio TDD, igual que ya se hizo con el
contrato de `Project`.

## Complexity Tracking

*Sin violaciones de la constitución — tabla no aplica.*
