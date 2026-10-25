# Implementation Plan: Validar proyecto existente al crear historia de backlog (Issue #19)

**Branch**: `fix/validar-proyecto-crear-historia` | **Date**: 2026-10-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/005-validar-proyecto-crear-historia/spec.md`

## Summary

`CrearHistoriaBacklog` hoy construye el `BacklogItem` y lo persiste sin verificar que el
proyecto referenciado exista, lo que genera historias huérfanas. La corrección reutiliza
exactamente el patrón ya vigente en `CrearSprint` (`internal/service/crear_sprint.go:25`):
agregar una dependencia `repository.ProjectRepository` al service y llamar
`proyectos.ObtenerPorID(ctx, input.ProyectoID)` **después** de `domain.NewBacklogItem`
(y antes de `repo.Guardar`), propagando el error tal cual (FR-006). No se agrega lógica
nueva: la validación de existencia ya existe en el repositorio (`ObtenerPorID` devuelve
`domain.ErrProyectoNoEncontrado`), así que el cambio es de cableado y orden de llamadas.

**Hallazgo que reduce alcance (FR-007/US2)**: el mapeo a HTTP 404 **ya existe** en
`internal/http/sprint_handler.go:118` (`escribirError` traduce `domain.ErrProyectoNoEncontrado`
a `StatusNotFound`). Por lo tanto, una vez que el service devuelva ese error, el endpoint
`POST /backlog` responderá 404 automáticamente, sin tocar el handler. El "500 genérico
actual" del issue ya no aplica. Igual se respeta la consigna: el mapeo queda **fuera de
alcance** de esta historia y no se agregan tests de handler para el 404; solo se deja
constancia (ver research.md).

## Technical Context

**Language/Version**: Go 1.26.5 (módulo `github.com/MoyaCarlos/seguimiento-medicion`).

**Primary Dependencies**: stdlib; `modernc.org/sqlite` solo detrás de los adaptadores
existentes. No se agregan dependencias nuevas.

**Storage**: **Sin cambios de esquema.** No se tocan `internal/repository/migrate.go` ni las
tablas. Las claves foráneas (`FOREIGN KEY`) quedan explícitamente fuera de alcance (FR-008):
la integridad se resuelve en la capa de aplicación, reutilizando `ProjectRepository.ObtenerPorID`.

**Testing**: paquete `testing` estándar de Go (TDD, RED/GREEN/REFACTOR) + Godog (BDD).
Los tests de `internal/service` reutilizan el fake ya versionado `fakeProjectRepository`
(`internal/service/fakes_test.go`), no se crea uno nuevo. El test de handler reutiliza el
`stubProjectRepository` ya existente en `internal/http/project_handler_test.go` (mismo paquete
`http`).

**Target Platform**: servidor backend (API REST), multiplataforma vía Go.

**Project Type**: web-service (backend del monorepo; el frontend React en `/web` queda fuera).

**Performance Goals**: no aplica; la validación agrega una lectura por identificador
(`ObtenerPorID`), O(1) indexada por PK. Sin metas de latencia específicas.

**Constraints**:
- El service debe seguir dependiendo de la **interfaz** `repository.ProjectRepository`, nunca
  de `*sql.DB` ni de `sqlite_project.go` (Dependency Inversion, Principio V).
- Se mantiene el **orden de validación actual**: primero las invariantes de `BacklogItem`
  (`domain.NewBacklogItem`), luego la existencia del proyecto (FR-005). No reordenar.
- El error de `ObtenerPorID` se **propaga tal cual** (`return`, sin envolver ni traducir),
  de modo que cubre `ErrProyectoNoEncontrado` y cualquier otro error (FR-006).
- Sin `FOREIGN KEY` en el esquema (FR-008).

**Scale/Scope**: acotado a un caso de uso y su cableado. Un cambio de firma de constructor
que impacta a **todos** sus llamadores (servicio, tests de servicio, tests de handler, wiring
y steps BDD). No incluye el mapeo HTTP ni el esquema.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Cumplimiento |
|---|---|
| I. Test-First (RED/GREEN/REFACTOR, commits por fase) | PASS — la validación es regla de negocio: el test de proyecto inexistente se escribe primero (RED) en `internal/service`, luego el GREEN de cableado/llamada. Commits por fase según la plantilla de tasks. |
| II. SDD (flujo Spec Kit completo) | PASS — este plan es parte del flujo de la corrección. |
| III. BDD (Godog, `/features`) | PASS — se agrega el escenario "proyecto inexistente" al `.feature` de creación de historias antes de `/speckit.implement`. |
| IV. Clean/Hexagonal (dependencia hacia adentro) | PASS — el dominio no cambia; el service depende de la interfaz `ProjectRepository`; SQLite solo se cablea en `main.go`. |
| V. SOLID/KISS/YAGNI/DRY | PASS — se reutiliza `ObtenerPorID` y el patrón de `CrearSprint` (DRY, sin duplicar validación); no se crea fake nuevo (KISS); no se agregan abstracciones (YAGNI). SRP intacto: un service por caso de uso. |
| VI. Stack obligatorio | PASS — Go + `modernc.org/sqlite`, sin dependencias nuevas. |

Sin violaciones — no hace falta llenar Complexity Tracking.

## Project Structure

### Documentation (this feature)

```text
specs/005-validar-proyecto-crear-historia/
├── plan.md              # Este archivo
├── research.md          # Fase 0
├── data-model.md        # Fase 1
├── quickstart.md        # Fase 1
├── contracts/
│   └── openapi.yaml     # Fase 1 (delta sobre el contrato de HU-01)
├── checklists/
│   └── requirements.md  # Ya generado en /speckit.specify
└── tasks.md             # Fase 2 (/speckit.tasks, no este comando)
```

### Source Code (repository root)

Pedido explícito:

```text
internal/service/
├── crear_historia.go        # EDITAR: nuevo campo proyectos + firma de constructor
│                            #   + validación ObtenerPorID antes de Guardar
└── crear_historia_test.go   # EDITAR: precargar proyecto en el fake + 3 casos nuevos

cmd/api/
└── main.go                  # EDITAR: pasar el ProjectRepository de SQLite a NewCrearHistoriaBacklog
```

Impactos **derivados del cambio de firma del constructor** (no listados en el pedido, pero
obligatorios para compilar y para no romper la suite; ver research.md D3):

```text
internal/http/
└── backlog_handler_test.go  # EDITAR: nuevoHandlerDePrueba() debe pasar un ProjectRepository;
                             #   reutilizar stubProjectRepository (ya en project_handler_test.go)
                             #   con el proyecto 1 precargado para el caso exitoso (201)

features/
├── steps_creacion.go        # EDITAR: iniciarBD() construye NewCrearHistoriaBacklog(...);
│                            #   pasarle el ProjectRepository SQLite de la misma db
└── crear_historia_backlog.feature  # EDITAR: nuevo escenario "proyecto inexistente"
```

**Archivos NO tocados** (evidencia de alcance): `internal/domain/backlog_item.go`,
`internal/domain/errors.go`, `internal/repository/project_repository.go`,
`internal/repository/sqlite_project.go`, `internal/repository/backlog_repository.go`,
`internal/repository/migrate.go`, `internal/http/backlog_handler.go` (el 404 ya lo cubre
`escribirError`), `internal/service/crear_sprint.go`.

**Structure Decision**: se mantiene un archivo por responsabilidad. La validación vive en el
caso de uso (`internal/service`), consistente con `CrearSprint`; el dominio y la persistencia
no cambian. El único motivo por el que el fix toca `internal/http` y `features` es que ambos
instancian el service: es un impacto de compilación, no un cambio de comportamiento HTTP.

## Complexity Tracking

*Sin violaciones de la constitución — tabla no aplica.*
