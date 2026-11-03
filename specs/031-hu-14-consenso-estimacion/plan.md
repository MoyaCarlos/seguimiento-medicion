# Implementation Plan: Consenso y Registro de la Estimación Acordada (HU-14)

**Branch**: `feature/HU-14-consenso-estimacion` | **Date**: 2026-11-03 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/031-hu-14-consenso-estimacion/spec.md`

## Summary

Completar en el backend el Planning Poker de HU-02: la consulta de una ronda revelada **indica las
diferencias** (si hay, con mínimo y máximo), el Scrum Master **abre nuevas rondas** conservando las
anteriores, cualquiera **consulta el historial** de rondas reveladas, y el Scrum Master **registra la
estimación acordada** en los Story Points de la historia, respetando D2b. Las reglas puras se agregan
al agregado `RondaVotacion`, con una `ReglaConsenso` (Strategy) para detectar diferencias; los
permisos y D2b, en tres casos de uso nuevos. No hay cambios de esquema. El frontend queda en la
sub-issue #93.

## Technical Context

**Language/Version**: Go 1.26.5

**Primary Dependencies**: `net/http`, `database/sql`, `encoding/json`, `modernc.org/sqlite`; Godog para BDD. Sin dependencias nuevas.

**Storage**: SQLite embebido; usa `rondas_votacion`, `votos` e `historias.estimacion_sp` (`MODELO-DE-DATOS.md` §4.4 y §4.7). Sin migración nueva.

**Testing**: `go test ./...` (dominio puro, services con fakes, repositorio contra `:memory:`, handler con `httptest`); Godog en `/features`.

**Target Platform**: Servidor que expone REST/JSON para la SPA React.

**Project Type**: Aplicación web (backend Go hexagonal + frontend React/Vite).

**Performance Goals**: Equipo chico; respuesta percibida inmediata (< 1 s).

**Constraints**: Sin CGO; dominio sin infraestructura; sin autenticación (`usuario_id` en la petición); una sola conexión SQLite (`SetMaxOpenConns(1)`, R9).

**Scale/Scope**: Equipos de 3 a 10 integrantes; pocas rondas por historia.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitución v1.1.0 (`.specify/memory/constitution.md`).

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First (NON-NEGOTIABLE) | `CoincidenciaExacta`, `Diferencias`, `Siguiente`, `ConEstimacion`, los services, el repositorio y el handler se hacen con commits RED → GREEN (→ REFACTOR si hay algo que limpiar). Extraer D2b es un commit `REFACTOR:` | PASS |
| II. SDD | spec → clarify (6 decisiones registradas) → plan; siguen tasks, analyze, BDD, implement. Decisiones validadas en la reunión de definición | PASS |
| III. BDD | `features/consenso_estimacion.feature` (@HU-14) con un escenario por criterio de la spec (20): normal, alternativo (sin diferencias, carta no votada), límite (primera carga con Sprint `Activo`) y error (permisos, baraja, sin revelar, D2b, Completada) | PASS |
| IV. Clean/Hexagonal | Reglas en `internal/domain`; puerto `EstimacionHistoriaRepository` (ISP); SQLite y HTTP como adaptadores; Strategy para el consenso, como nombra este principio | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un service por caso de uso; D2b en un solo lugar (R6). Strategy con una sola regla: ver Complexity Tracking | PASS |
| VI. Stack | Solo stdlib + `modernc.org/sqlite` + Godog | PASS |

**Modelo de datos (AGENTS.md)**: sin cambios en el esquema ni en `MODELO-DE-DATOS.md`. Ninguna decisión
del modelo quedó abierta: §4.7 define el historial y dónde va la acordada, y §5 asigna las reglas a HU-14.

**Re-evaluación post-diseño (Phase 1)**: todos los gates siguen en PASS. Los cambios en la consulta son
aditivos y no rompen a quien ya la consume, como la pantalla de HU-02 (PR #229).

## Project Structure

### Documentation (this feature)

```text
specs/031-hu-14-consenso-estimacion/
├── plan.md              # Este archivo
├── research.md          # Fase 0: decisiones técnicas (R1–R11)
├── data-model.md        # Fase 1: métodos nuevos, ReglaConsenso, vistas, puertos
├── quickstart.md        # Fase 1: guía de validación
├── contracts/
│   └── openapi.yaml     # Fase 1: nueva ronda, historial, estimación, campo diferencias
├── checklists/
│   └── requirements.md  # Calidad de la spec
└── tasks.md             # Fase 2 (/speckit.tasks)
```

### Source Code (repository root)

```text
cmd/api/main.go                              # wiring: regla de consenso, 3 services y las rutas de la votación

internal/domain/
├── votacion.go                              # + Siguiente, Diferencias, ReglaConsenso, CoincidenciaExacta
├── votacion_test.go
├── backlog_item.go                          # + ConEstimacion
├── backlog_item_test.go
└── errors.go                                # + ErrRondaNoRevelada

internal/service/
├── abrir_nueva_ronda.go                     # AbrirNuevaRonda (nuevo)
├── abrir_nueva_ronda_test.go
├── registrar_estimacion.go                  # RegistrarEstimacion (nuevo)
├── registrar_estimacion_test.go
├── consultar_historial_rondas.go            # ConsultarHistorialRondas (nuevo)
├── consultar_historial_rondas_test.go
├── estimacion_congelada.go                  # D2b, compartida con EditarHistoria (nuevo)
├── editar_historia.go                       # usa la función compartida
├── consultar_votacion.go                    # + Diferencias en la vista, recibe la regla
├── consultar_votacion_test.go
└── fakes_test.go                            # fakeVotacionRepository conserva las rondas anteriores

internal/repository/
├── votacion_repository.go                   # + ListarRondasReveladas
├── sqlite_votacion.go
├── sqlite_votacion_test.go
├── backlog_repository.go                    # + puerto EstimacionHistoriaRepository
├── sqlite_backlog.go                        # + GuardarEstimacion
└── sqlite_backlog_test.go

internal/http/
├── votacion_handler.go                      # + AbrirRonda, Historial, RegistrarEstimacion
├── votacion_dto.go                          # + diferencias, historial, estimacionRequest
└── votacion_handler_test.go

features/
├── consenso_estimacion.feature              # @HU-14 (nuevo)
├── steps_consenso.go                        # pasos Godog sobre SQLite :memory: (nuevo)
├── steps_votacion.go                        # wiring de la regla de consenso
└── features_test.go                         # + InitializeScenarioConsenso
```

**Structure Decision**: estructura hexagonal ya establecida en `AGENTS.md`. Los endpoints nuevos se suman
a `VotacionHandler`, porque son parte de la votación de la historia. `VotacionHandler.Registrar(mux)` agrupa
sus rutas, como `DemoHandler`, para poder probarlas. `escribirError` suma el 409 de una apertura simultánea (R10).

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Strategy (`ReglaConsenso`) con una sola implementación, en tensión con YAGNI | La constitución (IV), `AGENTS.md` y la sub-issue #91 lo piden para el consenso, y R3 de la spec 012 difirió la extracción a esta historia | Una función suelta alcanza hoy, pero no cumple esas indicaciones del equipo. Se limita a una interfaz de un método (R2) |
