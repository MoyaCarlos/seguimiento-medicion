# Implementation Plan: Avance de las Historias del Proyecto (HU-18)

**Branch**: `feature/HU-18-avance-historias` | **Date**: 2026-11-07 | **Spec**: [spec.md](./spec.md)

**Issue**: #177

**Input**: Feature specification from `specs/037-hu-18-avance-historias/spec.md`

## Summary

Reemplazar el placeholder de la pestaña "Estado" (`/projects/:id/status`) por una pantalla de solo lectura que muestra las historias del proyecto agrupadas en Nueva, En progreso y Completada, con la cantidad y los Story Points de cada grupo y el Sprint de cada historia. El backend suma **un endpoint de lectura**, `GET /projects/{id}/avance`, cuya agrupación y sumas viven en el dominio (única fuente de verdad); el DTO solo copia. Sin tablas ni migraciones: lee `historias` y `sprints`. El frontend consume el endpoint.

## Technical Context

**Language/Version**: Go 1.26 (`go.mod`); TypeScript 6, React 19, react-router 8 (Vite 8).

**Primary Dependencies**: `net/http`, `database/sql`, `modernc.org/sqlite`, Godog; Vitest y Testing Library. **Sin dependencias nuevas.**

**Storage**: SQLite embebido, solo lectura: `proyectos`, `historias` y `sprints` (`MODELO-DE-DATOS.md` §4.4 y §4.5). Sin migraciones (§8).

**Testing**: `go test ./...` (dominio y service con fakes, repositorio contra SQLite `:memory:`, handler con `httptest`); Godog en `/features`; Vitest + Testing Library en `web/` con el escenario en el nombre de cada prueba.

**Target Platform**: servidor que expone REST/JSON y sirve la SPA.

**Project Type**: aplicación web (backend Go hexagonal + frontend React/Vite).

**Performance Goals**: una consulta por apertura de la pantalla; respuesta percibida inmediata con decenas de historias.

**Constraints**: ninguna escritura; sin recalcular métricas de HU-03; el dominio no conoce SQLite ni HTTP; una feature de `web/` no importa de otra (`web/src/features/aislamiento.test.ts`).

**Scale/Scope**: un endpoint, un service, una función de dominio, un método de repositorio, una pantalla.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.* Constitución v1.1.0 (`.specify/memory/constitution.md`).

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First | Agrupación y sumas (dominio), orden de errores y unión con Sprints (service), consulta SQL (repositorio) y handler: cada uno con RED → GREEN (→ REFACTOR si hay algo que limpiar). Las funciones puras del frontend, igual. | PASS |
| II. SDD | specify → clarify → plan; siguen tasks, analyze, BDD, implement. | PASS |
| III. BDD | `features/avance_historias.feature` (@HU-18) para la API; la pantalla, con pruebas de componente con el escenario en el nombre. | PASS |
| IV. Clean/Hexagonal | Reglas puras en `internal/domain`; puerto angosto en `internal/repository`; SQLite y HTTP como adaptadores. | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un service por caso de uso; la agrupación vive solo en el dominio; sin filtros, paginación ni transacción de lectura. | PASS |
| VI. Stack | Sin dependencias nuevas. | PASS |

**Re-evaluación post-diseño**: sin cambios, todo en PASS. `MODELO-DE-DATOS.md` no se toca.

## Project Structure

### Documentation (this feature)

```text
specs/037-hu-18-avance-historias/
├── spec.md
├── plan.md              # este archivo
├── research.md          # decisiones técnicas (R1–R7)
├── data-model.md        # vistas de lectura, puerto y service
├── contracts/openapi.yaml
├── quickstart.md
├── checklists/requirements.md
└── tasks.md             # lo genera /speckit-tasks
```

### Source Code (repository root)

```text
cmd/api/main.go                       # wiring del service y del handler + ruta GET /projects/{id}/avance

internal/domain/
├── avance_historias.go               # AvanceHistorias, GrupoAvance, HistoriaAvance, AgruparAvance
└── avance_historias_test.go

internal/service/
├── consultar_avance_historias.go     # ConsultarAvanceHistorias
├── consultar_avance_historias_test.go
└── fakes_test.go                     # + fake del puerto (si no alcanza el existente)

internal/repository/
├── avance_historias_repository.go    # puerto AvanceHistoriasRepository (ListarPorProyecto)
├── sqlite_backlog.go                 # + ListarPorProyecto (ORDER BY id)
└── sqlite_backlog_test.go

internal/http/
├── avance_handler.go                 # AvanceHandler.Consultar
├── avance_dto.go
└── avance_handler_test.go

features/
├── avance_historias.feature          # @HU-18
├── steps_avance_historias.go
└── features_test.go                  # + registro de pasos

web/src/features/estado/
├── index.tsx, index.test.tsx         # reemplaza el placeholder
├── api.ts                            # consultarAvance
├── avance.ts, avance.test.ts         # tipos y funciones puras de presentación
├── GrupoHistorias.tsx, GrupoHistorias.test.tsx
└── estado.css
```

**Structure Decision**: estructura hexagonal ya establecida. `AvanceHandler` propio (SRP: `BacklogHandler` atiende el Product Backlog, que sigue mostrando solo las historias sin Sprint). `escribirError` no cambia: `ErrProyectoNoEncontrado` ya responde 404 y el `ValidationError` de `parsearID`, 400. El frontend reutiliza la feature `estado` existente y sus rutas (`status/*` ya apunta a `./features/estado`); no se toca `rutas.tsx` ni `navegacion.ts`.

## Riesgos

- **Conflictos en `cmd/api/main.go`**: otras ramas agregan rutas y wiring en el mismo bloque; se resuelven al actualizar la rama.
- **Puerto sobre el adaptador de historias**: `SQLiteBacklogRepository` ya cumple varios puertos; el método nuevo no cambia los existentes.

## Complexity Tracking

Sin violaciones que justificar.
