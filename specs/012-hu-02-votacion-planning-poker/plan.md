# Implementation Plan: Votación de Esfuerzo en Planning Poker (HU-02)

**Branch**: `feature/HU-02-votacion-planning-poker` | **Date**: 2026-10-27 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/012-hu-02-votacion-planning-poker/spec.md`

## Summary

Agregar al backend la votación secreta de Planning Poker sobre una historia: el Scrum Master
**inicia** la votación (abre la ronda 1), cada integrante **vota** una carta de la baraja
Fibonacci (puede cambiarla hasta revelar), la **consulta** muestra quién votó y solo el voto
propio, y el Scrum Master **revela** cuando votaron todos, con el promedio exacto y la carta
sugerida. Las reglas de la ronda y el cálculo viven en un agregado nuevo `RondaVotacion` del
dominio; los permisos, la existencia de la historia y la ocultación, en cuatro casos de uso.
La persistencia usa las tablas `rondas_votacion` y `votos` de `MODELO-DE-DATOS.md` §4.7,
creadas por la migración 003, con escrituras condicionales para que nadie vote sobre una
ronda ya revelada. HU-14 (diferencias, nuevas rondas, estimación acordada) y el frontend
quedan fuera.

## Technical Context

**Language/Version**: Go 1.26.5

**Primary Dependencies**: `net/http` (routing por método y `PathValue`), `database/sql`, `encoding/json`, `math`; `modernc.org/sqlite`; Godog para BDD. Sin dependencias nuevas.

**Storage**: SQLite embebido; tablas nuevas `rondas_votacion` y `votos` (migración `003_votacion.sql`, §4.7). Lee `historias`, `integrantes` y `usuarios` de la 001.

**Testing**: `go test ./...` (dominio puro, services con fakes, repositorio y test de esquema contra `:memory:`, handler con `httptest`); Godog en `/features`.

**Target Platform**: Servidor local que expone REST/JSON en `:8080` para la SPA React.

**Project Type**: Aplicación web (backend Go hexagonal + frontend React/Vite).

**Performance Goals**: Equipo chico en local; respuesta percibida inmediata (< 1 s).

**Constraints**: Sin CGO; dominio sin infraestructura; sin autenticación (`usuario_id` en la petición, research R2); una sola conexión SQLite (`SetMaxOpenConns(1)`).

**Scale/Scope**: Equipos de 3 a 10 integrantes; una ronda por historia en HU-02.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitución v1.0.0 (`.specify/memory/constitution.md`).

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First (NON-NEGOTIABLE) | `RondaVotacion` (votar, reemplazar, revelar), `CalcularPromedio`, los 4 services, el repositorio SQLite (incluido el upsert condicional), el test de esquema de la 003 y el handler se hacen con commits RED → GREEN (→ REFACTOR si hay algo que limpiar) | PASS |
| II. SDD | spec → clarify (5 decisiones registradas) → plan; siguen tasks, analyze, BDD, implement | PASS |
| III. BDD | `features/votacion_planning_poker.feature` (@HU-02) con un escenario por criterio de la spec: normal, alternativo (cambio de voto), límite (empate en el promedio, integrante que falta) y error (carta inválida, permisos, revelada, completada) | PASS |
| IV. Clean/Hexagonal | Reglas en `internal/domain`; puerto `VotacionRepository` (ISP); SQLite y HTTP como adaptadores; el servicio no ve transacciones (R6) | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un service por caso de uso; la baraja reutiliza `EsEstimacionValida`; sin Strategy para el promedio hasta que HU-14 lo pida (R3); sin historial de rondas en la API todavía | PASS |
| VI. Stack | Solo stdlib + `modernc.org/sqlite` + Godog | PASS |

**Modelo de datos (AGENTS.md)**: la migración 003 implementa §4.7 sin cambios, con el número
reservado en §9; no hace falta tocar `MODELO-DE-DATOS.md`. Ninguna decisión del modelo quedó
abierta.

**Re-evaluación post-diseño (Phase 1)**: todos los gates siguen en PASS. El 403 es un código
nuevo en la API (R8), justificado por las reglas de rol de la spec; no agrega capas.

## Project Structure

### Documentation (this feature)

```text
specs/012-hu-02-votacion-planning-poker/
├── plan.md              # Este archivo
├── research.md          # Fase 0: decisiones técnicas (R1–R10)
├── data-model.md        # Fase 1: RondaVotacion, Promedio, vista de consulta, puerto
├── quickstart.md        # Fase 1: guía de validación
├── contracts/
│   └── openapi.yaml     # Fase 1: iniciar, votar, revelar, consultar
├── checklists/
│   └── requirements.md  # Calidad de la spec
└── tasks.md             # Fase 2 (/speckit.tasks)
```

### Source Code (repository root)

```text
cmd/api/main.go                              # wiring del repo, los 4 services, el handler y las 4 rutas

internal/domain/
├── votacion.go                              # RondaVotacion, Voto, Promedio, NuevaRonda, Votar, Revelar, CalcularPromedio (nuevo)
├── votacion_test.go
└── errors.go                                # + errores de votación (research R8)

internal/service/
├── iniciar_votacion.go                      # IniciarVotacion (nuevo)
├── iniciar_votacion_test.go
├── votar.go                                 # Votar (nuevo)
├── votar_test.go
├── revelar_votacion.go                      # RevelarVotacion (nuevo)
├── revelar_votacion_test.go
├── consultar_votacion.go                    # ConsultarVotacion + vista con ocultación (nuevo)
├── consultar_votacion_test.go
├── integrantes.go                           # rolEnProyecto, compartido por los 4 (nuevo)
└── fakes_test.go                            # + fakeVotacionRepository (simula las escrituras condicionales, R6)

internal/repository/
├── migraciones/003_votacion.sql             # rondas_votacion, votos (§4.7)
├── esquema_test.go                          # + estructura y filas de las tablas de la 003
├── votacion_repository.go                   # puerto VotacionRepository (nuevo)
├── sqlite_votacion.go                       # adaptador SQLite (nuevo)
└── sqlite_votacion_test.go

internal/http/
├── votacion_handler.go                      # VotacionHandler: Iniciar, Votar, Revelar, Consultar (nuevo)
├── votacion_dto.go                          # requests y votacionResponse (nuevo)
├── votacion_handler_test.go
└── sprint_handler.go                        # escribirError: + 403 y errores de votación

features/
├── votacion_planning_poker.feature          # @HU-02 (nuevo)
├── steps_votacion.go                        # pasos Godog sobre SQLite :memory: (nuevo)
└── features_test.go                         # + InitializeScenarioVotacion
```

**Structure Decision**: estructura hexagonal ya establecida en `AGENTS.md`. La votación tiene
handler propio (`VotacionHandler`) aunque las rutas cuelguen de `/backlog/{id}`, igual que
`SprintHandler` atiende su recurso: `BacklogHandler` no crece con un agregado que no es suyo.
`escribirError` sigue centralizado en `sprint_handler.go`, donde está hoy.

## Complexity Tracking

Sin violaciones de la constitución que justificar.
