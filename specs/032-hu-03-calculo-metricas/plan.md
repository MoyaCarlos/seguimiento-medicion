# Implementation Plan: Cálculo de Métricas del Proyecto (HU-03)

**Branch**: `feature/HU-03-calculo-metricas` | **Date**: 2026-11-03 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/032-hu-03-calculo-metricas/spec.md`

## Summary

Calcular al consultar, con una sola petición, las métricas de cada Sprint finalizado y del proyecto:
Story Points planificados y completados, velocidad, horas estimadas y reales, desviación en horas y
porcentaje, porcentaje de historias completadas y defectos. El cálculo es una función pura del dominio,
alimentada por una consulta agregada de solo lectura. Además, el parámetro de horas por Story Point se
puede editar en el `PUT` del proyecto y en su pantalla, y la regla D2b pasa a prohibir la primera carga
de Story Points con Sprint `Finalizado`. No hay cambios de esquema.

## Technical Context

**Language/Version**: Go 1.26.5; TypeScript con React y Vite en el frontend.

**Primary Dependencies**: `net/http`, `database/sql`, `encoding/json`, `math`, `strconv`, `modernc.org/sqlite`; Godog para BDD; Vitest y Testing Library. Sin dependencias nuevas.

**Storage**: SQLite embebido, solo lectura en la consulta; escribe únicamente `proyectos.horas_por_story_point` (`MODELO-DE-DATOS.md` §4.1 y §6). Sin migración nueva.

**Testing**: `go test ./...` (dominio puro, services con fakes, repositorio contra `:memory:`, handler con `httptest`); Godog en `/features`; Vitest en `web/`.

**Target Platform**: Servidor que expone REST/JSON para la SPA React.

**Project Type**: Aplicación web (backend Go hexagonal + frontend React/Vite).

**Performance Goals**: Equipo chico; respuesta percibida inmediata (< 1 s). Una consulta por petición, sin N+1.

**Constraints**: Sin CGO; dominio sin infraestructura; sin autenticación; métricas no guardadas (§6); sumas en enteros y redondeo solo al final (D8).

**Scale/Scope**: Decenas de Sprints por proyecto como máximo; sin filtros ni paginación.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitución v1.1.0 (`.specify/memory/constitution.md`). Re-evaluada después del diseño: sin cambios.

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First (NON-NEGOTIABLE) | `CalcularMetricas`, `ConHorasPorStoryPoint`, la regla de D2b, el service, la consulta SQL y los handlers se hacen con commits RED → GREEN (→ REFACTOR si hay algo que limpiar). Los DTO sin lógica quedan fuera | PASS |
| II. SDD | spec → clarify (7 decisiones de la reunión y 4 definiciones confirmadas por el equipo) → plan; siguen tasks, analyze, BDD e implement | PASS |
| III. BDD | `features/calculo_metricas.feature` (@HU-03): normal, alternativo (sin parámetro, recálculo), límite (sin Sprints, Sprint sin historias, estimadas en cero, arrastrada en dos Sprints) y error (proyecto inexistente, valor inválido, carga rechazada). El frontend con pruebas de componente | PASS |
| IV. Clean/Hexagonal | Reglas en `internal/domain`; puerto `DatosMetricasRepository` (ISP); SQLite y HTTP como adaptadores | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un service por caso de uso; un único cálculo que consumen HU-09, HU-10 y HU-12; D2b sigue en una sola función (R9). No se usa Strategy: hay una sola forma de calcular | PASS |
| VI. Stack | Solo stdlib + `modernc.org/sqlite` + Godog + Vitest | PASS |

`MODELO-DE-DATOS.md`: sin cambios de tablas ni migraciones. Se actualizan sus secciones de reglas por
la decisión de la reunión que reabre D2b ([data-model.md](./data-model.md)).

## Project Structure

### Documentation (this feature)

```text
specs/032-hu-03-calculo-metricas/
├── plan.md              # Este archivo
├── research.md          # Fase 0
├── data-model.md        # Fase 1
├── quickstart.md        # Fase 1
├── contracts/
│   └── openapi.yaml     # Fase 1
├── checklists/
│   └── requirements.md
└── tasks.md             # Lo genera /speckit-tasks
```

### Source Code (repository root)

```text
internal/domain/
├── metricas.go                     # Nuevo: tipos y CalcularMetricas
├── metricas_test.go                # Nuevo
├── project.go                      # HorasPorStoryPoint y ConHorasPorStoryPoint
├── project_test.go
└── errors.go                       # ErrCargaEnSprintFinalizado

internal/service/
├── calcular_metricas.go            # Nuevo
├── calcular_metricas_test.go       # Nuevo
├── editar_proyecto.go              # Aplica el parámetro después de los demás datos
├── editar_proyecto_test.go
├── estimacion_congelada.go         # Regla de R9
├── editar_historia_test.go         # Cambia "finalizado permite la carga inicial"
└── registrar_estimacion_test.go    # Cambia "primera carga con un Sprint Finalizado"

internal/repository/
├── metricas_repository.go          # Nuevo: puerto
├── sqlite_metricas.go              # Nuevo: consulta agregada
├── sqlite_metricas_test.go         # Nuevo
├── sqlite_project.go               # Lee y actualiza el parámetro
└── sqlite_project_test.go

internal/http/
├── metricas_dto.go                 # Nuevo
├── metricas_handler.go             # Nuevo: GET /projects/{id}/metricas
├── metricas_handler_test.go        # Nuevo
├── project_dto.go                  # Campo en la respuesta
├── project_handler.go              # Tres estados del parámetro (R7)
└── project_handler_test.go

cmd/api/main.go                     # Wiring y ruta

features/
├── calculo_metricas.feature        # Nuevo (@HU-03)
├── steps_metricas.go               # Nuevo
├── features_test.go                # Registra los pasos
└── consulta_edicion_backlog.feature # Saca el ejemplo Finalizado de la carga inicial

web/src/features/proyecto/          # Sub-issue #53
├── api.ts
├── validacion.ts
├── validacion.test.ts
├── Editar.tsx
└── Editar.test.tsx

MODELO-DE-DATOS.md, docs/PRODUCT-BACKLOG.md   # Regla D2b, §6 y nota de HU-03
specs/002-…/contracts/openapi.yaml, specs/009-…/spec.md, specs/031-…/spec.md
```

**Structure Decision**: Se mantiene la estructura hexagonal existente. Las métricas agregan un archivo
por capa, y los demás cambios extienden el código de HU-04, HU-14 y HU-15.

## Orden sugerido

1. Traer `origin/main` a la rama: está 102 commits atrás y no tiene HU-14 ([R12](./research.md)).
2. Escenarios BDD.
3. User Story 1 y 2: dominio → consulta SQL → service → handler.
4. User Story 3: parámetro en dominio → repositorio → service → handler.
5. User Story 5: regla de D2b y sus documentos. Es independiente del cálculo.
6. User Story 4: pantalla (#53). No frena el cierre de #3.

## Dependencias con otras historias

- **HU-12 (spec 033)**: espera este merge. Usa los mismos nombres de `MetricasSprint`, así que su contrato no cambia. Debe llamar a `service.CalcularMetricas` y unir por `sprint_id`, y corregir su supuesto sobre D2b (línea de "Reproducibilidad").
- **HU-10 (spec 034)**: consume el mismo service; su FR-028 ya prevé la integración.
- **HU-09**: consume `GET /projects/{id}/metricas`.
- **PR #268**: modifica `consenso_estimacion.feature`, que este plan no toca (R10).
- **PR #250**: modifica `docs/PRODUCT-BACKLOG.md`; posible conflicto menor.

## Complexity Tracking

Sin violaciones de la constitución.
