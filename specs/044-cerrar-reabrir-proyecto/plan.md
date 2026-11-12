# Implementation Plan: Cerrar y reabrir un proyecto (#300)

**Branch**: `feature/300-cerrar-reabrir-proyecto` | **Date**: 2026-11-11 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/044-cerrar-reabrir-proyecto/spec.md`

## Summary

El Scrum Master cierra un proyecto (`POST /projects/{id}/cerrar`) o lo reabre
(`POST /projects/{id}/reabrir`). El cierre se guarda en una columna nueva,
`proyectos.fecha_cierre` (migración 008, cierra D16). En el dominio, `Project` suma
`FechaCierre` y los métodos `Cerrar`/`Reabrir`, y `CalcularEstadoProyecto` agrega una
primera regla: si el proyecto está cerrado, es «Finalizado». La solo lectura se impone en un
único guardia, `exigirProyectoAbierto`, que llaman los guardias de escritura de HU-21
(`exigirIntegrante`/`exigirScrumMaster`). Así responde 409 «El proyecto está cerrado»
después de los 403 y antes de las validaciones. Los `UPDATE` condicionales cubren los
cierres simultáneos. Solo backend.

## Technical Context

**Language/Version**: Go (versión de `go.mod`)

**Primary Dependencies**: `net/http` estándar, Godog; sin dependencias nuevas

**Storage**: SQLite (`modernc.org/sqlite`); **cambio de esquema**: migración 008
(`ALTER TABLE proyectos ADD COLUMN fecha_cierre`) + `MODELO-DE-DATOS.md` en el mismo PR

**Testing**: `testing` estándar (dominio, repositorio SQLite, servicios con fakes, integración en `internal/app`), Godog en `/features`

**Target Platform**: servidor Linux (Render) y desarrollo en Windows

**Project Type**: web-service (API REST en Go + SPA React; esta historia toca solo el backend)

**Performance Goals**: N/A (una consulta extra, `SELECT fecha_cierre`, por escritura)

**Constraints**: mantener el orden de errores de HU-21 (FR-006 de la spec 042) con el 409
nuevo en su lugar; una escritura rechazada no cambia datos; `Actualizar` del proyecto no
toca `fecha_cierre`

**Scale/Scope**: 1 columna, 2 rutas, 2 servicios, 4 errores de dominio, 1 guardia que
cubre 17 servicios de escritura (17 rutas), 1 `.feature`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Cumplimiento |
|---|---|
| I. Test-First | Reglas de negocio con RED → GREEN (→ REFACTOR), un commit por fase: `CalcularEstadoProyecto` con cierre, `Project.Cerrar`/`Reabrir`, repositorio (`Cerrar`/`Reabrir` condicionales, `ProyectoCerrado`), guardia `exigirProyectoAbierto`, servicios `CerrarProyecto`/`ReabrirProyecto`, test exhaustivo de rutas de escritura. |
| II. SDD | specify (044, `--number` a mano) → clarify (decisiones confirmadas en la spec) → plan → tasks → analyze → BDD → implement. |
| III. BDD | `features/cerrar_reabrir_proyecto.feature` (`@300`): normal (cerrar, reabrir), alternativo (reabrir con fechas vencidas), límite (cerrar sin Sprints; reabrir y volver a escribir), error (Sprint Activo, dos veces, reabrir abierto, 403/401/404, escrituras con 409). |
| IV. Clean/Hexagonal | Reglas en `internal/domain`; el servicio depende de los puertos `ProjectRepository`/`IntegrantesRepository`; la base solo agrega la columna y el `UPDATE` condicional. |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un servicio por caso de uso (SRP). El chequeo vive en un solo lugar (DRY, FR-005). Se extiende un puerto existente en vez de crear uno (KISS, [research §2](research.md)). No se guarda quién cerró (YAGNI). |
| VI. Stack | Sin cambios. La migración usa el runner `go:embed` existente. |

**Gate de modelo de datos (AGENTS.md)**: el cambio de esquema fue confirmado por la usuaria
el 2026-11-11 (spec, Clarifications) y va en un solo PR con `MODELO-DE-DATOS.md`.

**Resultado**: pasa. Sin violaciones que justificar. Re-chequeo después del diseño: pasa;
el diseño de la Fase 1 no agrega capas, puertos ni dependencias.

## Project Structure

### Documentation (this feature)

```text
specs/044-cerrar-reabrir-proyecto/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── cerrar-reabrir-proyecto.yaml
├── checklists/
│   └── requirements.md
└── tasks.md             # /speckit-tasks
```

### Source Code (repository root)

```text
internal/repository/migraciones/
└── 008_cierre_proyecto.sql          # nuevo: ADD COLUMN fecha_cierre + CHECK (D12)
internal/repository/
├── esquema_test.go                  # + casos D12 de fecha_cierre
├── migracion_008_test.go            # nuevo: base previa migra y queda abierta
├── registro_horas_repository.go     # IntegrantesRepository + ProyectoCerrado
├── project_repository.go            # ProjectRepository + Cerrar, Reabrir
├── sqlite_project.go                # lee/escribe fecha_cierre; Cerrar/Reabrir condicionales; ProyectoCerrado
├── sqlite_project_test.go           # + tests de Cerrar/Reabrir/ProyectoCerrado y carreras
└── project_repository_en_memoria.go # ídem en memoria
internal/domain/
├── project.go                       # + FechaCierre, Cerrar(ahora, sprints), Reabrir()
├── project_test.go                  # + tests de transiciones
├── estado_proyecto.go               # + regla 1: cerrado ⇒ Finalizado
├── estado_proyecto_test.go          # + casos con cierre
└── errors.go                        # + ErrProyectoCerrado, ErrProyectoYaCerrado, ErrProyectoNoCerrado, ErrSprintActivoAlCerrarProyecto
internal/service/
├── integrantes.go                   # verificarIntegrante/ScrumMaster + exigirProyectoAbierto
├── cerrar_proyecto.go               # nuevo (+ _test)
├── reabrir_proyecto.go              # nuevo (+ _test)
├── consultar_votacion.go            # usa verificarIntegrante (lectura)
├── votar.go, imputar_horas.go       # pasan a exigirIntegrante
├── proyecto_cerrado_test.go         # nuevo: cada servicio de escritura devuelve ErrProyectoCerrado
└── fakes_test.go                    # fakes con ProyectoCerrado
internal/http/
├── project_dto.go                   # proyectoResponse + fecha_cierre; respuesta con estado
├── project_handler.go               # handlers Cerrar y Reabrir (+ _test); NewProjectHandler recibe los 2 servicios
├── sprint_handler.go                # escribirError: 4 errores nuevos → 409
└── rutas.go                         # 2 rutas, nivel SesionScrumMaster
internal/app/
├── app.go                           # cableado de CerrarProyecto/ReabrirProyecto
├── rutas_permisos_test.go           # nivelesDeLaSpec + 2 rutas
└── proyecto_cerrado_test.go         # nuevo: toda ruta de escritura → 409 sin cambios; tabla vs. Rutas()
features/
├── cerrar_reabrir_proyecto.feature  # nuevo, @300
├── steps_cerrar_reabrir_proyecto.go # nuevo
└── features_test.go                 # registro de los pasos
MODELO-DE-DATOS.md                   # §2, §3, §4.1, §4.10, §5, §6, §7 (D16), §8, §9
```

**Structure Decision**: estructura existente del backend (Clean/Hexagonal liviana). Sin
archivos de frontend (lo hace otra integrante, a partir de `fecha_cierre` en las respuestas
del proyecto y del contrato).

## Fases de implementación (insumo para `/speckit-tasks`)

1. **Esquema** (RED/GREEN): test de esquema con los casos D12 de `fecha_cierre` y test de la
   migración 008 → `008_cierre_proyecto.sql`. `MODELO-DE-DATOS.md` en la misma fase
   ([research §7](research.md)).
2. **Dominio — estado** (RED/GREEN): `CalcularEstadoProyecto` con `FechaCierre` ⇒
   «Finalizado» aunque haya un Sprint Activo, fechas futuras o ningún Sprint; sin cierre,
   los casos de HU-13 siguen pasando.
3. **Dominio — transiciones** (RED/GREEN/REFACTOR): `Project.Cerrar` (abierto → cerrado con
   el día civil; ya cerrado → `ErrProyectoYaCerrado`; Sprint Activo →
   `ErrSprintActivoAlCerrarProyecto`; Pendientes o Finalizados no impiden) y
   `Project.Reabrir` (cerrado → abierto; abierto → `ErrProyectoNoCerrado`).
4. **Repositorio** (RED/GREEN): lectura de `fecha_cierre` en `ObtenerPorID` y
   `ListarPorIntegrante`; `Actualizar` no la pisa; `Cerrar`/`Reabrir` condicionales con sus
   errores al afectar 0 filas; `ProyectoCerrado`. Lo mismo en el repositorio en memoria.
5. **Guardia** (RED/GREEN/REFACTOR): `exigirProyectoAbierto` en `exigirIntegrante` y
   `exigirScrumMaster`; `verificar*` para la lectura de la votación; `Votar` e
   `ImputarHoras` pasan a `exigirIntegrante`. El `RED:` incluye el test de cada servicio de
   escritura contra un proyecto cerrado y el ajuste de los fakes (método nuevo del puerto);
   el `GREEN:`, solo producción.
6. **Servicios** (RED/GREEN): `CerrarProyecto` (404 → 403 → 409s → `Cerrar` del repo) y
   `ReabrirProyecto`, con reloj inyectado para «hoy».
7. **HTTP** (RED/GREEN): rutas, handlers, `fecha_cierre` en `proyectoResponse`, `estado` en la
   respuesta de cerrar/reabrir y los 409 en `escribirError`; `nivelesDeLaSpec` y test de
   `internal/app` con **todas** las rutas de escritura contra un proyecto cerrado
   ([research §9](research.md)).
8. **BDD**: `.feature` + pasos, escenarios de US1 a US4.
9. **Cierre**: `go build`, `go vet`, `gofmt -l .`, `go test ./...` en verde;
   `/speckit-analyze`; marcar tareas.

## Riesgos

- **Conflicto con #304 (spec 043, abierto)**: comparten `internal/domain/project.go`,
  `MODELO-DE-DATOS.md` §5 y `features/steps_estado_proyecto.go`. Son cambios de texto, no
  de lógica. Si #304 se mergea primero, se trae `main` antes de abrir el PR
  ([research §10](research.md)).
- **Escritura concurrente con el cierre**: límite declarado, no se resuelve
  ([research §6](research.md)).
- **Escenarios existentes**: ninguno usa proyectos cerrados (la columna nace en `NULL`), así
  que no deberían cambiar. Si alguno falla, se frena y se revisa antes de ajustarlo.

## Complexity Tracking

Sin violaciones de la constitución; no aplica.
