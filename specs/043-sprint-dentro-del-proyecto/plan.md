# Implementation Plan: Sprints dentro del plazo del proyecto (#168)

**Branch**: `feature/168-sprint-dentro-del-proyecto` | **Date**: 2026-11-11 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/043-sprint-dentro-del-proyecto/spec.md`

## Summary

Al iniciar un Sprint, rechazar con 400 las fechas que queden fuera del plazo del proyecto: inicio antes de la fecha de inicio del proyecto (`fecha_inicio`) o fin después de su fecha de fin (`fecha_fin`), solo contra las fechas que el proyecto tenga y con bordes incluidos. La regla es un método de dominio, `Project.AdmiteSprint(inicio, fin)`, que `IniciarSprint` llama después de las validaciones propias del Sprint y antes de persistir; para eso el servicio recibe el `ProjectRepository`. Sin cambio de esquema; solo se documenta la regla en `MODELO-DE-DATOS.md` §5.

## Technical Context

**Language/Version**: Go (versión de `go.mod`)

**Primary Dependencies**: `net/http` estándar, Godog; sin dependencias nuevas

**Storage**: SQLite (`modernc.org/sqlite`); sin cambio de esquema ni migración

**Testing**: `testing` estándar (dominio y servicio, con repositorios en memoria), Godog en `/features`

**Target Platform**: servidor Linux (Render) y desarrollo en Windows

**Project Type**: web-service (API REST en Go + SPA React, esta historia toca solo el backend)

**Performance Goals**: N/A (una lectura extra del proyecto por inicio de Sprint)

**Constraints**: mantener el orden de errores existente (FR-006); un inicio rechazado no persiste nada (FR-007)

**Scale/Scope**: 1 método de dominio, 1 servicio, 1 `.feature`, ajuste de llamadores de `NewIniciarSprint` (cableado, 5 archivos de pasos, tests de servicio y handler)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Cumplimiento |
|---|---|
| I. Test-First | Regla de validación: ciclo RED → GREEN (→ REFACTOR si hace falta), cada fase en su commit. Primero dominio (`AdmiteSprint`), después servicio (`IniciarSprint`). |
| II. SDD | specify → clarify (decisiones del equipo en la spec) → plan → tasks → analyze → BDD → implement. Spec 043 con `--number` a mano. |
| III. BDD | `features/sprint_dentro_del_proyecto.feature` con caso normal (bordes), alternativo (proyecto sin fechas), límite (bordes exactos) y error (inicio antes, fin después). |
| IV. Clean/Hexagonal | La regla vive en `internal/domain`; el servicio depende de la interfaz `ProjectRepository`, no de SQLite. |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un método en la entidad dueña del plazo; sin puerto nuevo ni normalización innecesaria de fechas ([research.md §2, §3, §5](research.md)). |
| VI. Stack | Sin cambios de stack. |

**Resultado**: pasa. Sin violaciones que justificar. Re-chequeo después del diseño: pasa (el diseño de Fase 1 no agrega capas ni dependencias).

## Project Structure

### Documentation (this feature)

```text
specs/043-sprint-dentro-del-proyecto/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── iniciar-sprint-400.yaml
├── checklists/
│   └── requirements.md
└── tasks.md             # /speckit-tasks
```

### Source Code (repository root)

```text
internal/domain/
├── project.go                     # + AdmiteSprint(inicio, fin) error
└── project_test.go                # + tests de la regla (tabla de casos)
internal/service/
├── iniciar_sprint.go              # + proyectos ProjectRepository; llama AdmiteSprint
├── iniciar_sprint_test.go         # + casos del plazo; siembra el proyecto
├── cerrar_sprint_test.go          # ajuste de NewIniciarSprint
└── crear_sprint_test.go           # ajuste de NewIniciarSprint
internal/http/
└── sprint_handler_test.go         # ajuste de NewIniciarSprint (+ caso 400 por plazo)
cmd/api/main.go                    # cableado: NewIniciarSprint(sprintRepo, projectRepo)
features/
├── sprint_dentro_del_proyecto.feature   # nuevo, @HU-05 @168
├── steps_sprint_dentro_del_proyecto.go  # nuevo
└── steps_{asignacion,defectos,estado_historia,estado_proyecto,sprints}.go  # ajuste de NewIniciarSprint
MODELO-DE-DATOS.md                 # §5: fila de la regla nueva
```

**Structure Decision**: estructura existente del backend (Clean/Hexagonal liviana). Sin archivos de frontend: el formulario de inicio ya muestra el error por campo ([research.md §8](research.md)).

## Fases de implementación (insumo para `/speckit-tasks`)

1. **Dominio** (RED/GREEN/REFACTOR): `Project.AdmiteSprint` con tabla de casos: sin fechas; solo inicio (antes, igual, después); solo fin (antes, igual, después); ambas (dentro, bordes, inicio afuera, fin afuera, las dos afuera → `fecha_inicio`); mensaje con la fecha límite.
2. **Servicio** (RED/GREEN): `IniciarSprint` rechaza fuera del plazo y deja el Sprint Pendiente; acepta en bordes; el orden de errores se mantiene (Goal vacío con fechas afuera → `sprint_goal`; otro Activo → 409). El cambio de firma de `NewIniciarSprint` y el ajuste de todos sus llamadores (cableado, pasos BDD, tests) van en el mismo commit `GREEN:` para que el árbol compile; el `RED:` previo solo agrega los tests nuevos (que no compilan o fallan, como en las historias anteriores).
3. **BDD**: `.feature` + pasos; escenarios de la spec (US1–US3).
4. **Documentación**: fila de §5 en `MODELO-DE-DATOS.md`; marcar tareas hechas.
5. **Cierre**: `go build`, `go vet`, `gofmt -l .`, `go test ./...` en verde; `/speckit-analyze` registrado.

## Riesgos

- **Conflicto con #301 (HU-21)**: ambos cambian `NewIniciarSprint`. Resolución acordada en [research.md §5](research.md): firma `NewIniciarSprint(sprints, proyectos ProjectRepository)` y un solo campo usado para permisos y para el plazo. Si #301 se mergea primero, se trae `main` a esta rama antes del PR y el cableado pasa a `internal/app/app.go`.

## Complexity Tracking

Sin violaciones de la constitución; no aplica.
