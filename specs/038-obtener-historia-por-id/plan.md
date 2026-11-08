# Implementation Plan: Consultar una historia por su identificador

**Branch**: `feature/222-obtener-historia` | **Date**: 2026-11-08 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/038-obtener-historia-por-id/spec.md` (issue #222)

## Summary

Agregar `GET /backlog/{id}`, que devuelve una sola historia con la misma forma que el resto de la API de backlog. Es solo un adaptador HTTP: el caso de uso `ObtenerHistoria` ya existe y ya está inyectado en `BacklogHandler`.

Decisión del punto pendiente de la spec (texto del 400): el handler nuevo usa el helper `idDeRuta` de `internal/http/sprint_handler.go`, que ya responde 400 con campo `id` y mensaje `identificador inválido`. No se toca `parsearID` ni el texto de ningún endpoint existente. Ver [research.md](research.md).

## Technical Context

**Language/Version**: Go (módulo `github.com/MoyaCarlos/seguimiento-medicion`), `net/http` con `ServeMux` y patrones `GET /backlog/{id}`

**Primary Dependencies**: solo las existentes; Godog para BDD. Sin dependencias nuevas.

**Storage**: SQLite (`modernc.org/sqlite`). Sin migración y sin cambios en `MODELO-DE-DATOS.md`.

**Testing**: `testing` estándar (tests de handler con `stubBacklogRepository`) y Godog en `/features`

**Target Platform**: API REST local (`http://localhost:8080`)

**Project Type**: web-service (backend Go; el frontend no se toca)

**Performance Goals**: N/A (lectura por clave)

**Constraints**: solo lectura; sin sesión (HU-21 en Sprint 5); no cambiar el comportamiento de endpoints existentes

**Scale/Scope**: 1 método de handler, 1 ruta, 1 `.feature` con 5 escenarios, 1 sección en la guía de curl

## Constitution Check

| Principio | Evaluación |
|-----------|------------|
| I. Test-First | PASA. Ciclo RED → GREEN por commit. El RED compila con un método `Obtener` vacío y falla por aserción (no por error de compilación). |
| II. SDD | PASA. Flujo specify → plan → tasks → implement, una historia a la vez. |
| III. BDD | PASA. Cinco escenarios Godog en un `.feature` nuevo: historia del Product Backlog, en Sprint, completada, inexistente e identificador inválido. |
| IV. Hexagonal | PASA. El handler solo adapta HTTP; la lógica queda en `service.ObtenerHistoria`. `domain` no cambia. |
| V. Clean/DRY/YAGNI | PASA. Se reutilizan `ObtenerHistoria`, `idDeRuta`, `escribirError` y `aRespuesta`; sin `case` nuevos ni helper nuevo. |
| VI. Stack | PASA. Sin dependencias nuevas. |

Resultado: sin violaciones; no se completa Complexity Tracking. Se reevalúa tras el diseño: sin cambios.

## Project Structure

### Documentation (this feature)

```text
specs/038-obtener-historia-por-id/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── obtener-historia.md
└── tasks.md             # lo genera /speckit.tasks
```

### Source Code (repository root)

```text
cmd/api/main.go                          # + mux.HandleFunc("GET /backlog/{id}", backlogHandler.Obtener), junto a PUT /backlog/{id}
internal/http/backlog_handler.go         # + método Obtener
internal/http/backlog_handler_test.go    # + tests de Obtener (RED/GREEN)
features/obtener_historia.feature        # nuevo, 5 escenarios
features/steps_obtener_historia.go       # nuevo, pasos Godog
features/features_test.go                # + registro del nuevo InitializeScenario...
docs/curl-frontend-sprint-2.md           # + sección en versión bash y en versión cmd
```

**Structure Decision**: estructura hexagonal existente. No se agregan paquetes, servicios ni DTOs: la respuesta es `historiaResponse` vía `aRespuesta`.

## Diseño

### Handler

`Obtener` en `BacklogHandler`:

1. `id, ok := idDeRuta(w, r)`; si `!ok`, ya respondió 400 (campo `id`, `identificador inválido`) y retorna. Cubre `abc`, `0`, `-3`, `1.5` y un número que desborda `int64` (`ParseInt` devuelve error), y `007` se parsea como 7.
2. `h.obtenerHistoria.Ejecutar(ctx, id)`; ante error, `escribirError(w, err)`. `ErrHistoriaNoEncontrada` ya mapea a 404 con `historia no encontrada` (sin `case` nuevos).
3. `escribirJSON(w, 200, aRespuesta(item))`.

La validación del id va antes de consultar (FR-007). `idDeRuta` está en el mismo paquete `http`, así que no hay import nuevo.

### Ruta

`mux.HandleFunc("GET /backlog/{id}", backlogHandler.Obtener)` en `cmd/api/main.go`, junto a `PUT /backlog/{id}`. No colisiona con `GET /backlog/{id}/horas` (patrón de más segmentos).

### Pruebas del handler (RED/GREEN)

- **RED**: método `Obtener` vacío (compila, responde 200 con cuerpo vacío por defecto) → los tests fallan por aserción.
- Casos: 200 sin `sprint_id` (`null`); 200 con `sprint_id`; 200 completada con `fecha_finalizacion`; 404 con `historia no encontrada`; 400 con campo `id` y `identificador inválido` para `abc`, `0`, `-3`, `1.5` y un número fuera de rango (p. ej. `99999999999999999999`); `007` se interpreta como 7.
- **GREEN**: implementación mínima de arriba.
- Un test adicional verifica que el 400 no consulta el repositorio (FR-007).

### BDD (Godog)

`features/obtener_historia.feature` con los cinco escenarios de la spec y `features/steps_obtener_historia.go` con sus pasos, sobre SQLite en memoria como `steps_consultar_sprints.go`. Se registra en `features_test.go`. Las suites existentes no cambian.

### Documentación

Sección "Consultar una historia" en `docs/curl-frontend-sprint-2.md`, en las dos versiones (bash y cmd), con éxito, 404 y 400.

## Complexity Tracking

Sin violaciones.
