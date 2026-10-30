# Implementation Plan: Registro y Seguimiento de Defectos (HU-08)

**Branch**: `feature/HU-08-registro-defectos` | **Date**: 2026-10-30 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/023-hu-08-registro-defectos/spec.md`

## Summary

Agregar al backend tres casos de uso sobre los defectos de un proyecto: **registrar** un
defecto (`POST /projects/{id}/defectos`), **marcarlo como resuelto**
(`POST /defectos/{id}/resolver`) y **listar** los defectos del proyecto
(`GET /projects/{id}/defectos`). La tabla `defectos` ya está diseñada en
`MODELO-DE-DATOS.md` §4.9 y entra con la migración reservada **005**. La entidad nueva
`Defecto` valida descripción y severidad y encapsula la transición única
`Pendiente → Resuelto`. Los services buscan el Sprint `Activo` del proyecto y lo asignan
como Sprint de detección o de resolución. El adaptador SQLite repite la condición
"ese Sprint sigue `Activo`" en el mismo `INSERT`/`UPDATE` para cubrir la carrera con el cierre
del Sprint y con dos resoluciones a la vez. El frontend queda fuera de este incremento
(sub-issue #68).

## Technical Context

**Language/Version**: Go (versión del `go.mod` del repo)

**Primary Dependencies**: `net/http` (routing por método y `PathValue`), `database/sql`, `encoding/json` (stdlib); `modernc.org/sqlite`; Godog para BDD. Sin dependencias nuevas.

**Storage**: SQLite embebido; tabla nueva `defectos` (migración `005_defectos.sql`, DDL de `MODELO-DE-DATOS.md` §4.9). Lee `proyectos`, `historias` y `sprints`.

**Testing**: `go test ./...` (dominio y service con fakes, repositorio y esquema contra SQLite `:memory:`, handler con `httptest`); Godog en `/features`.

**Target Platform**: Servidor local que expone REST/JSON en `:8080` para la SPA React.

**Project Type**: Aplicación web (backend Go hexagonal + frontend React/Vite).

**Performance Goals**: Mono-usuario local; respuesta percibida inmediata (< 1 s).

**Constraints**: Sin CGO; dominio sin dependencias de infraestructura; a lo sumo un Sprint `Activo` por proyecto (HU-05); los Sprints de un defecto nunca se recalculan.

**Scale/Scope**: Decenas de defectos por proyecto; sin filtros ni paginación.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitución v1.1.0 (`.specify/memory/constitution.md`).

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First (NON-NEGOTIABLE) | Validaciones de `Defecto` (descripción, severidad), transición `Resolver`, búsqueda del Sprint activo, orden de errores en los services, `INSERT`/`UPDATE` condicionales (carreras), listado y handlers: cada uno con commits RED → GREEN (→ REFACTOR si hay algo que limpiar) | PASS |
| II. SDD | spec → clarify (3 decisiones registradas) → plan; siguen tasks, analyze, BDD, implement | PASS |
| III. BDD | `features/registro_defectos.feature` (@HU-08) con los escenarios de la spec: normal (registro, resolución), alternativo (historia fuera del Sprint activo, resolución en otro Sprint), límite (2000 caracteres, mismo Sprint de detección y resolución) y de error (sin Sprint activo, ya resuelto, datos inválidos, inexistentes) | PASS |
| IV. Clean/Hexagonal | Reglas en `internal/domain`; puertos angostos en `internal/repository` (ISP); SQLite y HTTP como adaptadores; el handler solo decodifica el JSON | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un service por caso de uso (`RegistrarDefecto`, `ResolverDefecto`, `ListarDefectos`); se reutilizan `ObtenerProyecto`, `HistoriaRepository.ObtenerPorID` y `SprintRepository.ListarPorProyecto`; la búsqueda del Sprint activo vive en una sola función del dominio; sin edición, borrado, reapertura ni filtros | PASS |
| VI. Stack | Solo stdlib + `modernc.org/sqlite` + Godog | PASS |

**Re-evaluación post-diseño (Phase 1)**: todos los gates siguen en PASS. El esquema no cambia:
la migración 005 copia el DDL ya acordado, así que no hace falta modificar
`MODELO-DE-DATOS.md` en este PR. La ruta `POST /defectos/{id}/resolver` agrega el prefijo
`/defectos` a `PrefijosAPI` y al proxy de `web/vite.config.ts` (R1); es configuración, no una
dependencia nueva.

## Project Structure

### Documentation (this feature)

```text
specs/023-hu-08-registro-defectos/
├── plan.md              # Este archivo
├── research.md          # Fase 0: decisiones técnicas (R1–R10)
├── data-model.md        # Fase 1: Defecto, vista de listado, puertos y services
├── quickstart.md        # Fase 1: guía de validación
├── contracts/
│   └── openapi.yaml     # Fase 1: POST/GET /projects/{id}/defectos y POST /defectos/{id}/resolver
├── checklists/
│   └── requirements.md  # Calidad de la spec
└── tasks.md             # Fase 2 (/speckit.tasks)
```

### Source Code (repository root)

```text
cmd/api/main.go                              # wiring de los 3 services, el handler y las 3 rutas

internal/domain/
├── defecto.go                               # Defecto, Severidad, EstadoDefecto, NewDefecto,
│                                            # (*Defecto).Resolver, DefectoDetalle
├── defecto_test.go
├── sprint.go                                # + SprintActivoDe(sprints) (Sprint, bool)
├── sprint_test.go
└── errors.go                                # + ErrDefectoNoEncontrado, ErrSinSprintActivoParaDetectar,
                                             #   ErrSinSprintActivoParaResolver, ErrDefectoYaResuelto,
                                             #   ErrHistoriaRelacionadaObligatoria, ErrHistoriaRelacionadaAjena

internal/service/
├── registrar_defecto.go                     # RegistrarDefecto (nuevo)
├── registrar_defecto_test.go
├── resolver_defecto.go                      # ResolverDefecto (nuevo)
├── resolver_defecto_test.go
├── listar_defectos.go                       # ListarDefectos (nuevo)
├── listar_defectos_test.go
└── fakes_test.go                            # + fake de DefectoRepository

internal/repository/
├── migraciones/005_defectos.sql             # DDL de MODELO-DE-DATOS.md §4.9
├── esquema_test.go                          # + casos de defectos (§9.6)
├── defecto_repository.go                    # puerto DefectoRepository
├── sqlite_defecto.go                        # Registrar (INSERT condicional), ObtenerPorID,
│                                            # Resolver (UPDATE condicional), ListarPorProyecto (JOIN historias)
└── sqlite_defecto_test.go

internal/http/
├── defecto_handler.go                       # Registrar, Resolver, Listar
├── defecto_dto.go                           # registrarDefectoRequest, defectoResponse
├── defecto_handler_test.go
├── sprint_handler.go                        # escribirError: + ErrDefectoNoEncontrado → 404
└── spa.go                                   # PrefijosAPI: + "/defectos"

web/vite.config.ts                           # proxy: + '/defectos': apiProxy

features/
├── registro_defectos.feature                # @HU-08 (nuevo)
├── steps_defectos.go                        # pasos Godog (nuevo)
└── features_test.go                         # + InitializeScenarioDefectos
```

**Structure Decision**: estructura hexagonal ya establecida en `AGENTS.md`. Los defectos
tienen su propio handler (`DefectoHandler`) en lugar de sumarse a `ProjectHandler`, que ya
tiene siete dependencias (SRP). `escribirError` suma un solo caso (`ErrDefectoNoEncontrado`
→ 404); el resto de los errores nuevos son `ValidationError` (400), ya mapeado.

## Complexity Tracking

> Sin violaciones a la constitución.

Notas de alcance (decisiones documentadas, no violaciones):

| Tema | Decisión | Justificación |
|------|----------|---------------|
| Frontend (formulario, listado, acción de resolver, mensaje sin Sprint activo) | Fuera de este incremento; sub-issue #68 | DoD: la sub-issue de Frontend no frena el cierre de la historia |
| Reabrir, editar o borrar defectos | Fuera de alcance | Clarification 1, FR-016, D9 |
| Quién registró o resolvió el defecto | No se guarda | La tabla no tiene esa columna y no hay autenticación (Assumptions) |
| Filtros y paginación del listado | No se implementan | Clarification 3, YAGNI |
| Métricas de defectos por Sprint | Fuera de alcance | HU-03; el listado y la tabla dejan los datos listos (SC-005) |
| Carrera con el cierre del Sprint y resolución doble | `INSERT`/`UPDATE` condicionales en una sola sentencia | R4; mismo criterio que HU-07 (R5) y el fix #145 |
| Prefijo `/defectos` nuevo | Se suma a `PrefijosAPI` y al proxy de Vite | R1; el test `TestPrefijosAPI_IncluyenLosQueProxeaVite` los mantiene alineados |
