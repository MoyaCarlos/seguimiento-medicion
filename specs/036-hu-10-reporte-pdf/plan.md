# Implementation Plan: Generación de Reporte en PDF (HU-10)

**Branch**: `feature/HU-10-reporte-pdf` | **Date**: 2026-11-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/036-hu-10-reporte-pdf/spec.md`

## Summary

Generar a demanda dos reportes PDF con Maroto: el de un Sprint **finalizado** y el de un
**proyecto**. El servicio arma un agregado de solo lectura con los datos que ya exponen las
historias anteriores (proyecto, integrantes, Sprints, Product Backlog, Sprint Backlog,
historias arrastradas, horas, defectos y métricas de HU-03) y el paquete `internal/report`
lo convierte en PDF. El Sprint Activo se lista con sus historias, pero no suma a métricas,
esfuerzo ni defectos. Un proyecto sin Sprints finalizados genera el reporte con "sin datos"
en esas secciones. La generación se hace completa en memoria antes de responder: si falla, no
se descarga un archivo corrupto, vacío o parcial.

## Technical Context

**Language/Version**: Go 1.26.5 (backend); React + Vite en el frontend (fuera del alcance del Sprint 4).

**Primary Dependencies**: `net/http`, `context`, `database/sql`, `github.com/johnfercher/maroto/v2` v2.4.4 (Maroto, puro Go); sin CGO ni dependencias fuera del stack.

**Storage**: SQLite embebido, solo lectura. No hay tablas, columnas ni migraciones nuevas; se reutiliza el esquema de `MODELO-DE-DATOS.md` (incluida `historias_arrastradas`).

**Testing**: `go test ./...` (dominio puro, service con fakes, repositorio SQLite contra `:memory:`, handlers con `httptest` y un generador de PDF falso); Godog en `/features`; smoke test del PDF real (cabecera `%PDF-` y cierre `%%EOF`).

**Target Platform**: Servidor HTTP que expone REST/JSON y descargas PDF a la SPA React y a clientes directos.

**Project Type**: Aplicación web (backend Go hexagonal + frontend React/Vite).

**Performance Goals**: Equipo y proyecto académicos; respuesta < 1 s para decenas de historias y Sprints. No hay requisito de streaming.

**Constraints**: PDF con Maroto (constitución, Principio VI); solo lectura; totales del proyecto calculados únicamente con Sprints finalizados; UTF-8 conservado (Arial cp1252 cubre tildes, eñes y signos de apertura); sin autenticación.

**Scale/Scope**: Un proyecto por reporte; todos sus Sprints finalizados; sin filtros, plantillas ni historial de reportes.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitución v1.1.0 (`.specify/memory/constitution.md`). Re-evaluada después del diseño: sin cambios.

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First (NON-NEGOTIABLE) | El armado del reporte, la consulta de arrastradas, el generador PDF y los handlers se desarrollan con commits RED → GREEN. Maroto no tiene lógica de negocio propia | PASS |
| II. SDD | Spec y clarify ya hechos; este plan, research, data-model, contrato, quickstart y tasks; después analyze, BDD e implement | PASS |
| III. BDD | `features/reporte_pdf.feature` (@HU-10) con los escenarios de US1–US4: Sprint finalizado, proyecto completo, proyecto sin Sprints finalizados y error real. US5 (frontend, #78) queda para el Sprint 5 | PASS |
| IV. Clean/Hexagonal | El agregado y las reglas viven en `internal/domain`; el servicio depende de puertos de repositorio; Maroto y HTTP son adaptadores | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un servicio de caso de uso, un paquete de generación, reutilización de `domain.CalcularMetricas`; no se crea una consulta agregada nueva ni se persiste el reporte | PASS |
| VI. Stack | Maroto es la librería fijada por la constitución; no se agregan otras dependencias | PASS |

`MODELO-DE-DATOS.md`: sin cambios; `data-model.md` enlaza las secciones existentes.

## Project Structure

### Documentation (this feature)

```text
specs/036-hu-10-reporte-pdf/
├── plan.md              # Este archivo
├── research.md          # Fase 0
├── data-model.md        # Fase 1
├── quickstart.md        # Fase 1
├── contracts/
│   └── openapi.yaml     # Fase 1
├── checklists/
│   └── requirements.md
└── tasks.md             # Lo genera /speckit.tasks
```

### Source Code (repository root)

```text
internal/domain/
├── reporte.go                      # Agregado del reporte (no persistido)
├── reporte_test.go                 # Reglas de armado puras (si las hay)
└── errors.go                       # ErrSprintNoFinalizado y ErrReporteNoGenerado

internal/repository/
├── reporte_repository.go           # Puerto HistoriasArrastradasRepository
├── sqlite_backlog.go               # Implementa la lectura de arrastradas
└── sqlite_backlog_test.go          # Tests de la consulta

internal/service/
├── generar_reporte.go              # Armado del reporte de Sprint y de proyecto
├── generar_reporte_test.go
└── fakes_test.go                   # Fakes de métricas y arrastradas

internal/report/
├── doc.go
├── pdf.go                          # Generador con Maroto
└── pdf_test.go                     # PDF válido y no vacío

internal/http/
├── reporte_handler.go              # GET /projects/{id}/reporte y /sprints/{id}/reporte
├── reporte_handler_test.go
├── spa.go                          # Excepción de descarga en ServirSPA
├── spa_test.go
└── sprint_handler.go               # Mapeo de errores nuevos

cmd/api/main.go                     # Wiring y rutas

features/
├── reporte_pdf.feature             # @HU-10
├── steps_reporte.go
└── features_test.go                # Registra los pasos

go.mod / go.sum                     # Maroto v2.4.4
```

**Structure Decision**: Se mantiene la estructura hexagonal. El armado del reporte queda en
`internal/service`, los tipos sin persistencia en `internal/domain` y Maroto en
`internal/report`, que no conoce HTTP ni SQLite.

## Orden sugerido

1. Traer `origin/main` a la rama (hecho) y mergear la spec pendiente.
2. Escenarios BDD.
3. Agregado del dominio y puerto de arrastradas.
4. Service de armado (Sprint y proyecto).
5. Generador PDF con Maroto.
6. Handler, rutas, excepción de descarga en la SPA y cableado.
7. Pasos BDD, checks del CI y validación manual.

## Dependencias con otras historias

- **HU-03**: `domain.CalcularMetricas` es la única fuente de los valores de métricas; el reporte no recalcula nada.
- **HU-08**: los defectos y sus Sprints se leen con `DefectoRepository.ListarPorProyecto`.
- **HU-05/HU-06**: `ListarPorSprint` y `historias_arrastradas` aportan las historias planificadas y completadas.
- **HU-04/HU-13**: proyecto, integrantes y estado derivado.
- **HU-07**: los minutos reales por Sprint llegan ya agregados en `DatosMetricasSprint`.
- **Frontend #78**: queda para el Sprint 5; consumirá estas rutas con `fetch` y descarga por blob.

## Complexity Tracking

Sin violaciones de la constitución.
