# Implementation Plan: Creación de Historias de Usuario (HU-01)

**Branch**: `feature/HU-01-crear-historias-usuario` | **Date**: 2026-09-30 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-hu-01-crear-historias-usuario/spec.md`

## Summary

Implementar la creación de historias de usuario del Product Backlog de punta a punta
con arquitectura hexagonal liviana (Ports & Adapters): entidad de dominio `BacklogItem`
con invariantes validadas por un constructor Factory, caso de uso `CrearHistoriaBacklog`
que persiste a través del puerto `BacklogRepository`, adaptador de persistencia SQLite
(`modernc.org/sqlite`) y adaptador de entrada HTTP (`POST /backlog`). La estimación en
Story Points se modela como puntero opcional (`*int`) que queda en `nil` hasta cargarse
vía Planning Poker (HU-02). No se implementan edición ni borrado (YAGNI).

**Alcance de esta iteración (backend)**: creación, persistencia y respuesta de la historia.
La vista/listado del Product Backlog y su actualización en vivo, la autorización por rol
(FR-014) y la validación de existencia del proyecto (FR-012) quedan diferidas (ver
Complexity Tracking). El orden de creación se preserva vía `id` incremental.

## Technical Context

**Language/Version**: Go 1.26.5

**Primary Dependencies**: Librería estándar `net/http` (routing por método, Go ≥ 1.22) y `database/sql`; `modernc.org/sqlite` como driver SQLite puro Go (sin CGO). Testing: paquete `testing` estándar (TDD) y Godog para BDD.

**Storage**: SQLite embebido vía `modernc.org/sqlite` (driver `"sqlite"`).

**Testing**: `go test ./...` para dominio, servicio, repositorio y handler; escenario BDD con Godog en `/features`.

**Target Platform**: Servidor local (Windows/Linux) que expone una API REST/JSON. El frontend React consume el endpoint; su UI queda fuera del alcance de esta iteración backend.

**Project Type**: Aplicación web (backend Go + frontend React), con estructura hexagonal.

**Performance Goals**: Aplicación local mono-usuario; la creación debe percibirse inmediata (< 1 s de extremo a extremo). Sin objetivos de throughput.

**Constraints**: Sin CGO; la dependencia apunta siempre hacia el dominio; no incorporar dependencias más allá de `modernc.org/sqlite` y Godog (KISS/YAGNI).

**Scale/Scope**: MVP de Product Backlog; decenas a cientos de historias por proyecto.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

El archivo `.specify/memory/constitution.md` sigue siendo la plantilla sin ratificar de Spec Kit (solo placeholders, sin principios formales), por lo que no hay gates constitucionales formales. Se usan como gates los principios rectores documentados en `AGENTS.md`.

| Gate | Fuente | Estado |
|------|--------|--------|
| Dependencia apunta hacia adentro (domain no conoce SQLite/HTTP) | Arquitectura hexagonal / DIP | PASS |
| Un service por caso de uso; services dependen de interfaces | SRP / DIP | PASS |
| Sin dependencias innecesarias (API REST con stdlib) | KISS / YAGNI | PASS |
| Lógica de negocio y validaciones en un único lugar (dominio) | DRY | PASS |
| Ciclo TDD RED→GREEN→REFACTOR en dominio/servicio | TDD obligatorio | PASS |
| Escenario BDD por criterio de aceptación | BDD | PASS |
| Sin edición/borrado ni otros alcances fuera de la spec | YAGNI | PASS |
| Constitución formal ratificada | Gobernanza | WARN — `/speckit.constitution` no ejecutado; `constitution.md` es plantilla vacía |

**Re-evaluación post-diseño (Phase 1)**: los gates siguen en PASS. La única alerta (WARN) es administrativa y no bloquea; se recomienda ratificar la constitución con `/speckit.constitution` para dejar de apoyarse en `AGENTS.md`.

## Project Structure

### Documentation (this feature)

```text
specs/001-hu-01-crear-historias-usuario/
├── plan.md              # Este archivo
├── research.md          # Fase 0: decisiones técnicas
├── data-model.md        # Fase 1: entidad BacklogItem y reglas
├── quickstart.md        # Fase 1: guía de validación end-to-end
├── contracts/
│   └── openapi.yaml     # Fase 1: contrato REST de POST /backlog
├── checklists/
│   └── requirements.md  # Calidad de la spec (ya existe)
└── tasks.md             # Fase 2: tareas (generado por /speckit.tasks)
```

### Source Code (repository root)

```text
cmd/
└── api/
    └── main.go                      # wiring: abre SQLite, migra, repo→service→handler, registra POST /backlog

internal/
├── domain/
│   ├── backlog_item.go              # BacklogItem + NewBacklogItem (Factory + invariantes)
│   ├── prioridad.go                 # enum MoSCoW (M/S/C/W) + validación
│   ├── estado.go                    # enum Estado (valor inicial "Nueva")
│   ├── valor_negocio.go             # escala Fibonacci permitida (1,2,3,5,8,13,21)
│   ├── errors.go                    # ValidationError{Campo, Mensaje}
│   └── backlog_item_test.go         # tests unitarios TDD
├── service/
│   ├── crear_historia.go            # CrearHistoriaBacklog (caso de uso) + input
│   └── crear_historia_test.go       # tests con fake en memoria
├── repository/
│   ├── backlog_repository.go        # interfaz BacklogRepository (puerto)
│   ├── sqlite_backlog.go            # adaptador SQLite (modernc.org/sqlite)
│   ├── sqlite_backlog_test.go       # test de integración en :memory:
│   └── migrate.go                   # CREATE TABLE IF NOT EXISTS backlog_items
└── http/
    ├── backlog_handler.go           # handler POST /backlog
    ├── backlog_dto.go               # DTOs request/response
    └── backlog_handler_test.go      # tests con httptest

features/
└── crear_historia_backlog.feature   # escenario BDD (Godog)
```

**Structure Decision**: Se respeta la estructura ya definida en `AGENTS.md`. El paquete `internal/http` conserva su nombre actual (documentado en `internal/http/doc.go`); dentro de él la librería estándar se referencia por su nombre de paquete `http` sin conflicto, ya que el nombre del paquete contenedor no es un identificador en su propio scope. El fake en memoria del repositorio vive como test double no exportado dentro de `internal/service/crear_historia_test.go` (no en producción), según YAGNI.

## Complexity Tracking

> Sin violaciones a los principios del proyecto.

Notas de alcance (decisiones documentadas, no violaciones):

| Tema | Decisión | Justificación |
|------|----------|---------------|
| Listado/vista del backlog (FR-003, SC-002) | Esta iteración solo persiste y preserva el orden de creación vía `id` incremental; la vista/orden visual se difiere a un incremento de frontend | Backend-only; la UI no forma parte de esta historia |
| Actualización en vivo sin recargar (ex SC-006) | Diferida al incremento de frontend | Criterio de interfaz, no verificable por backend |
| Asociación a proyecto (FR-012) | Se persiste `proyecto_id` como referencia; la validación de existencia del proyecto se difiere a HU-04 | La entidad `Project` y su repositorio aún no existen; no se declara FK a una tabla inexistente |
| Autorización por rol (FR-014) | Diferida a HU-04 | No hay autenticación/identidad todavía; forzarla ahora sería YAGNI y agregaría alcance no solicitado |
| Valor de Negocio | Se incluye como campo opcional (`*int`) con validación Fibonacci | FR-008/FR-009 de la spec (confirmado en `/speckit.clarify`) |
