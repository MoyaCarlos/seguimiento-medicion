# Implementation Plan: Fechas y Estado del Proyecto (HU-13)

**Branch**: `feature/HU-13-fechas-estado-proyecto` | **Date**: 2026-10-04 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/004-hu-13-estado-proyecto/spec.md`

## Summary

Exponer la **consulta del estado general de un proyecto** ("Planificado", "En curso" o
"Finalizado") a partir de la fecha actual, de las fechas de inicio/fin ya existentes del
`Project` (HU-04) y del estado de sus `Sprint` (HU-05). El estado es un valor **derivado y
no persistido**: se calcula en cada consulta con una función pura del dominio. No se agregan
tablas, columnas ni campos nuevos a la base de datos ni a las entidades existentes. Se
implementa con el mismo patrón Clean/Hexagonal de HU-04/HU-05: regla en `internal/domain`,
un caso de uso en `internal/service` que depende únicamente de los puertos
`ProjectRepository` y `SprintRepository`, y un adaptador HTTP de lectura.

Precedencia acordada (clarificación 2026-10-04): **los Sprints mandan, las fechas son
respaldo**. Un Sprint `Activo` ⇒ "En curso" aunque la fecha de fin haya pasado; si no hay
`Activo` y todos los Sprints están `Finalizado` ⇒ "Finalizado"; en el resto deciden las
fechas del proyecto.

## Technical Context

**Language/Version**: Go 1.26 (mismo módulo `github.com/MoyaCarlos/seguimiento-medicion`)

**Primary Dependencies**: stdlib `net/http` para el handler REST (patrón ya usado en
`internal/http/*`); `modernc.org/sqlite` sólo detrás de los adaptadores existentes.
No se agregan dependencias nuevas.

**Storage**: **Sin cambios de esquema.** El estado no se guarda. Se reutilizan las tablas
`projects` (con `start_date`/`end_date`) y `sprints` (con `estado`) ya migradas por HU-04 y
HU-05. `internal/repository/migrate.go` no se toca.

**Testing**: paquete `testing` estándar de Go (TDD, ciclo RED/GREEN/REFACTOR) + Godog (BDD).
Los tests de `internal/service` usan los fakes en memoria ya versionados
(`ProjectRepositoryEnMemoria`, `SprintRepositoryEnMemoria`).

**Target Platform**: servidor backend (API REST), multiplataforma vía Go.

**Project Type**: web-service (backend del monorepo; el frontend React en `/web` queda fuera
de alcance).

**Performance Goals**: no especificado por la consigna; prioridad en correctitud
determinista de la regla de derivación. La consulta es O(cantidad de Sprints del proyecto).

**Constraints**: el dominio no puede depender de SQLite, HTTP ni del reloj del sistema: la
función de cálculo recibe el instante `ahora time.Time` como parámetro (testeable sin tocar
`time.Now()` dentro del dominio). Los servicios dependen de interfaces de repositorio, nunca
de `*sql.DB` ni de `sqlite_*.go`.

**Scale/Scope**: acotado a esta historia — un tipo de estado derivado en el dominio, un caso
de uso de consulta y un endpoint de lectura. No incluye registrar/editar fechas (HU-04) ni
el frontend de visualización.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Cumplimiento |
|---|---|
| I. Test-First (RED/GREEN/REFACTOR, commits por fase) | PASS — el cálculo es regla de negocio: se testea primero en `internal/domain` y luego en el service; commits por fase según la plantilla de tasks |
| II. SDD (flujo Spec Kit completo) | PASS — este plan es parte del flujo |
| III. BDD (Godog, `/features`) | PASS — escenario Gherkin de consulta de estado generado antes de `/speckit.implement` |
| IV. Clean/Hexagonal (dominio sin dependencias externas) | PASS — `EstadoProyecto` + función pura en `internal/domain`; `ObtenerEstadoProyecto` depende de los puertos `ProjectRepository`/`SprintRepository`; SQLite sólo se cablea en `cmd/api/main.go` |
| V. SOLID/KISS/YAGNI/DRY | PASS — un service para el único caso de uso; no se persiste el estado (YAGNI); la regla vive en un solo lugar del dominio (DRY) y no se duplica en el handler |
| VI. Stack obligatorio | PASS — Go + `modernc.org/sqlite`, sin dependencias nuevas |

Sin violaciones — no hace falta llenar Complexity Tracking.

**Verificación explícita de los puntos pedidos**:

- **Arquitectura por capas**: dominio (regla) → servicio (caso de uso) → repositorios (puertos
  ya existentes) → HTTP (adaptador de lectura). La dependencia apunta hacia adentro.
- **Un caso de uso por service**: `ObtenerEstadoProyecto` es el único caso de uso de esta
  historia; no se agrupa con crear/editar proyecto.
- **Services dependen de interfaces**: `ObtenerEstadoProyecto` guarda
  `repository.ProjectRepository` y `repository.SprintRepository` (interfaces), no `*sql.DB`.
- **Sin tablas ni campos nuevos**: `migrate.go` y las structs `Project`/`Sprint` no se
  modifican; el estado se calcula al consultar.

## Project Structure

### Documentation (this feature)

```text
specs/004-hu-13-estado-proyecto/
├── plan.md              # Este archivo
├── research.md          # Fase 0
├── data-model.md        # Fase 1
├── quickstart.md        # Fase 1
├── contracts/           # Fase 1
│   └── openapi.yaml
├── checklists/
│   └── requirements.md
└── tasks.md             # Fase 2 (/speckit.tasks, no este comando)
```

### Source Code (repository root)

```text
internal/domain/
├── estado_proyecto.go        # NUEVO: tipo EstadoProyecto (Planificado/En curso/Finalizado),
│                             #   EsValida() y función pura CalcularEstadoProyecto(p, sprints, ahora)
└── estado_proyecto_test.go   # NUEVO: tests TDD de la regla (RED/GREEN)

internal/service/
├── obtener_estado_proyecto.go      # NUEVO: caso de uso único (depende de los 2 puertos)
└── obtener_estado_proyecto_test.go # NUEVO: tests con fakes en memoria

internal/http/
├── project_dto.go            # EDITAR: agregar estadoProyectoResponse
├── project_handler.go        # EDITAR: método ObtenerEstado (GET /projects/{id}/status)
└── project_handler_test.go   # EDITAR: tests del endpoint

cmd/api/
└── main.go                   # EDITAR: wiring de ObtenerEstadoProyecto + ruta

features/
└── estado_proyecto.feature   # NUEVO (se genera antes de /speckit.implement)
```

**Estructura de archivos NO tocados** (evidencia de "sin cambios de esquema"):
`internal/repository/migrate.go`, `internal/domain/project.go`, `internal/domain/sprint.go`,
`internal/repository/sqlite_project.go`, `internal/repository/sqlite_sprint.go`.

**Structure Decision**: un archivo por responsabilidad dentro de cada paquete, igual que
HU-01/HU-04/HU-05. El estado es un concepto de dominio (regla), no un DTO de persistencia, por
eso vive en `internal/domain` y no en `internal/repository`.

## Complexity Tracking

*Sin violaciones de la constitución — tabla no aplica.*
