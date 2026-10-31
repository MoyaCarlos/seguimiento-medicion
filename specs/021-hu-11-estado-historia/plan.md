# Implementation Plan: Actualización de Estado de una Historia en el Sprint (HU-11)

**Branch**: `feature/HU-11-estado-historia` | **Date**: 2026-10-29 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/021-hu-11-estado-historia/spec.md`

## Summary

Agregar un caso de uso sobre historias y Sprints que ya existen: **cambiar el estado de una
historia asignada a un Sprint activo** (`PUT /backlog/{id}/estado` con `{"estado": "..."}`). La
regla es una función pura del dominio, `ValidarCambioDeEstado(historia, sprint, nuevoEstado)`, que
devuelve el primer error según el orden fijado en los Edge Cases de la spec (estado inválido →
historia sin Sprint → Sprint no activo → estado repetido → retroceso → salto directo
"Nueva → Completada"). El servicio `CambiarEstadoHistoria` solo orquesta: obtiene la historia
(404), obtiene su Sprint, aplica la regla, calcula la fecha con un reloj inyectable y guarda con un
único `UPDATE historias SET estado = ?, fecha_finalizacion = ? WHERE id = ? AND estado = ?`; si no
afecta filas, devuelve un error de conflicto que el handler traduce a 409. Sin migraciones ni
cambios en `MODELO-DE-DATOS.md`: `historias.estado` y `historias.fecha_finalizacion` existen desde
la 001. La pantalla (US3) va en una fase propia, bloqueada hasta tener la pantalla de Sprints.

**Dos ajustes sobre las decisiones acordadas** (se avisan porque cambian firmas, no el
comportamiento):

1. `ValidarCambioDeEstado` recibe el Sprint como **puntero** (`*domain.Sprint`), no por valor: la
   historia puede no tener Sprint y entonces el servicio no tiene ningún Sprint que obtener. La
   regla de "historia sin Sprint" vive igual en el dominio y se evalúa mirando `SprintID`, no el
   puntero ([research.md](./research.md) R1).
2. `domain.BacklogItem` suma `FechaFinalizacion *time.Time` (y con eso `columnasHistoria` y
   `escanearHistoria` leen la columna) porque la respuesta aditiva `fecha_finalizacion` la
   necesita. Es lectura de una columna que ya existe, no un cambio de esquema.

El resto se implementa como se acordó: `EstadoEnProgreso` en `internal/domain/estado.go` (con
`Estado.EsValida()` siguiendo el patrón de `Prioridad.EsValida()`), el puerto chico, el service
con reloj inyectable y el 409 sin campo en `escribirError`.

## Technical Context

**Language/Version**: Go 1.26.5 (backend); TypeScript + React + Vite (frontend, US3)

**Primary Dependencies**: `net/http` (rutas con `PathValue`), `database/sql`, `encoding/json`;
`modernc.org/sqlite`; Godog para BDD. Frontend: las que ya usa `web/` (Vitest + Testing Library).
Sin dependencias nuevas.

**Storage**: SQLite embebido; columnas `historias.estado` y `historias.fecha_finalizacion` de
`001_esquema_base.sql`, con los `CHECK` que atan `Completada` a la fecha y prohíben
`En progreso`/`Completada` sin Sprint. Sin migración nueva.

**Testing**: `go test ./...` (dominio puro sin base, service con fakes, repositorio contra
`:memory:` migrada, handlers con `httptest`); Godog en `/features`; Vitest en `web/` para US3.

**Target Platform**: Servidor local que expone REST/JSON en `:8080` para la SPA React.

**Project Type**: Aplicación web (backend Go hexagonal + frontend React/Vite).

**Performance Goals**: Mono-usuario local; respuesta percibida inmediata (< 1 s, SC-001).

**Constraints**: Sin CGO; dominio sin dependencias de infraestructura; no tocar
`MODELO-DE-DATOS.md` ni agregar migraciones; cambio de a una historia; sin transiciones hacia
atrás ni salto directo.

**Scale/Scope**: Decenas de historias por Sprint; pocos Sprints por proyecto.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitución v1.1.0 (`.specify/memory/constitution.md`).

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First (NON-NEGOTIABLE) | `ValidarCambioDeEstado` (una prueba por regla y por orden de precedencia), `CambiarEstado` del repositorio (cambio, fecha, carrera), el service (reloj, 404, Sprint) y el handler (ruta, 409), con commits RED → GREEN (→ REFACTOR si hay algo que limpiar) | PASS |
| II. SDD | spec → clarify (5 Q resueltas) → plan; siguen tasks, analyze, BDD, implement | PASS |
| III. BDD | `features/cambiar_estado_historia.feature` (@HU-11) con los escenarios de US1 y US2 (normal, alternativo, límite, error); US3 se prueba con Vitest | PASS |
| IV. Clean/Hexagonal | Regla en `internal/domain` sin base; puerto en `internal/repository`; SQLite y HTTP como adaptadores | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un service por caso de uso; puerto chico y propio (ISP); una sola sentencia condicional; sin columnas nuevas ni reglas duplicadas | PASS |
| VI. Stack | Solo stdlib + `modernc.org/sqlite` + Godog; React/Vite en el frontend | PASS |

**Re-evaluación post-diseño (Phase 1)**: todos los gates siguen en PASS. Los dos ajustes de firma
están justificados en [research.md](./research.md) (R1, R3) y no agregan capas ni dependencias.

## Project Structure

### Documentation (this feature)

```text
specs/021-hu-11-estado-historia/
├── plan.md              # Este archivo
├── research.md          # Phase 0: decisiones y alternativas
├── data-model.md        # Phase 1: enlaza a MODELO-DE-DATOS.md §4.4/§5
├── quickstart.md        # Phase 1: validación de punta a punta
├── contracts/
│   └── openapi.yaml     # Phase 1: PUT /backlog/{id}/estado
├── checklists/
│   └── requirements.md
└── tasks.md             # Phase 2 (/speckit-tasks; no lo crea este comando)
```

### Source Code (repository root)

```text
internal/domain/
├── estado.go                    # + EstadoEnProgreso y Estado.EsValida()
├── estado_historia.go           # NUEVO: ValidarCambioDeEstado + errores de la transición
└── estado_historia_test.go      # NUEVO: pruebas puras, sin base

internal/repository/
├── backlog_repository.go        # + puerto CambioEstadoHistoriaRepository
├── sqlite_backlog.go            # + CambiarEstado; columnasHistoria y escanearHistoria suman fecha_finalizacion
└── sqlite_backlog_test.go       # + cambio de estado, fecha, conflicto de carrera

internal/service/
├── cambiar_estado_historia.go       # NUEVO: CambiarEstadoHistoria (reloj inyectable)
├── cambiar_estado_historia_test.go  # NUEVO
└── fakes_test.go                    # + fake del puerto nuevo

internal/http/
├── backlog_handler.go           # + CambiarEstado (PUT /backlog/{id}/estado)
├── backlog_dto.go               # + cambiarEstadoRequest; historiaResponse suma fecha_finalizacion
├── sprint_handler.go            # escribirError (compartido) suma el case 409 sin campo
└── backlog_handler_test.go      # + tests del handler

cmd/api/main.go                  # wiring del service y la ruta

features/
├── cambiar_estado_historia.feature  # NUEVO (@HU-11)
├── steps_estado_historia.go         # NUEVO
└── features_test.go                 # registra los pasos nuevos

web/src/features/                 # US3 (bloqueada hasta tener la pantalla de Sprints)
└── sprints/                      # control de estado por historia del Sprint activo + cliente de la ruta
```

**Structure Decision**: se respeta la arquitectura hexagonal liviana de `AGENTS.md`. No hay
paquetes nuevos: cada capa suma archivos junto a los de HU-05, HU-06 y HU-15. La pantalla se
integra en la feature `sprints` existente de `web/`; los nombres exactos de componentes se fijan en
`tasks.md` cuando el código de la pantalla de Sprints esté en `main`.

## Etapas de implementación

1. **US1 — Avanzar el estado (P1, MVP)**: dominio (`EstadoEnProgreso` + `ValidarCambioDeEstado`) →
   repositorio (`CambiarEstado` condicional y `fecha_finalizacion` en la lectura) → service (reloj,
   404, Sprint) → handler, ruta y `escribirError` (409) → escenario BDD de las transiciones
   válidas y de los rechazos por transición. Todo en verde antes de seguir.
2. **US2 — Rechazo sin Sprint e historia inexistente (P2)**: cubre solo la historia sin Sprint
   (paso a "En progreso" y a "Completada") y la historia inexistente, que es lo que implementan
   T030–T033. Los demás rechazos (Sprint no activo, salto directo, retroceso y estado repetido)
   pertenecen a US1 y se cubren en su fase. La carrera (409) se cubre con un test del repositorio
   y del handler, donde es determinista.
3. **US3 — Pantalla (P3)**: **bloqueada** hasta que la pantalla de Sprints esté en `main`. Por cada
   historia del Sprint activo, un control para pasar a "En progreso" o a "Completada"; tras un 200
   se refresca la lista con el estado y la fecha de finalización; ante un rechazo se muestra el
   `mensaje` de la API (incluido el 409 de carrera, que pide recargar). Se prueba con Vitest y
   Testing Library, con el escenario en el nombre de cada prueba.

## Complexity Tracking

Sin violaciones de la constitución que justificar. Los dos ajustes de firma (Sprint como puntero y
`FechaFinalizacion` en el modelo de lectura) no agregan capas ni dependencias.
