# Implementation Plan: Imputación de Horas Trabajadas (HU-07)

**Branch**: `feature/HU-07-imputacion-horas` | **Date**: 2026-10-28 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/019-hu-07-imputacion-horas/spec.md`

## Summary

Agregar al backend dos casos de uso sobre las historias: **imputar horas** a una historia del
Sprint activo (`POST /backlog/{id}/horas`) y **consultar su esfuerzo real**
(`GET /backlog/{id}/horas`: total y detalle de cargas). La tabla `registros_horas` ya está
diseñada en `MODELO-DE-DATOS.md` §4.8 y entra con la migración reservada **004**. La entidad
nueva `RegistroHoras` valida horas, actividad y fecha; las horas llegan como texto decimal y
se convierten a minutos enteros con `math/big.Rat`, sin `float64` (D8). El service
`ImputarHoras` verifica, en el orden de la spec, la historia, el Sprint activo y el
integrante, y el adaptador SQLite inserta con un `INSERT ... WHERE EXISTS` que repite la
condición de Sprint activo para cubrir la carrera con el cierre del Sprint. El frontend queda
fuera de este incremento (sub-issue #63).

## Technical Context

**Language/Version**: Go (versión del `go.mod` del repo)

**Primary Dependencies**: `net/http` (routing por método y `PathValue`), `database/sql`, `encoding/json`, `math/big` (stdlib); `modernc.org/sqlite`; Godog para BDD. Sin dependencias nuevas.

**Storage**: SQLite embebido; tabla nueva `registros_horas` (migración `004_registros_horas.sql`, DDL de `MODELO-DE-DATOS.md` §4.8). Lee `historias`, `sprints`, `integrantes` y `usuarios`.

**Testing**: `go test ./...` (dominio y service con fakes, repositorio y esquema contra SQLite `:memory:`, handler con `httptest`); Godog en `/features`.

**Target Platform**: Servidor local que expone REST/JSON en `:8080` para la SPA React.

**Project Type**: Aplicación web (backend Go hexagonal + frontend React/Vite).

**Performance Goals**: Mono-usuario local; respuesta percibida inmediata (< 1 s).

**Constraints**: Sin CGO; dominio sin dependencias de infraestructura; minutos enteros y sumas exactas (D8); "hoy" = día civil local del servidor (R7).

**Scale/Scope**: Decenas de cargas por historia; sin paginación.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitución v1.0.0 (`.specify/memory/constitution.md`).

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First (NON-NEGOTIABLE) | Conversión horas→minutos, validaciones de `RegistroHoras`, rango de fechas, regla de Sprint activo, integrante, orden de errores, INSERT condicional (carrera), total y handlers: cada uno con commits RED → GREEN (→ REFACTOR si hay algo que limpiar) | PASS |
| II. SDD | spec → clarify (6 decisiones registradas) → plan; siguen tasks, analyze, BDD, implement | PASS |
| III. BDD | `features/imputacion_horas.feature` (@HU-07) con los escenarios de la spec: normal, alternativo (acumulado, cambio de Sprint), límite (24 h, bordes de fecha) y de error | PASS |
| IV. Clean/Hexagonal | Reglas en `internal/domain`; puertos angostos en `internal/repository` (ISP); SQLite y HTTP como adaptadores; el handler solo extrae texto del JSON | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Un service por caso de uso (`ImputarHoras`, `ConsultarEsfuerzoReal`); se reutilizan `ListarIntegrantes`, `ObtenerPorID` y el helper de día civil de HU-13; sin edición ni borrado de cargas ni tope diario | PASS |
| VI. Stack | Solo stdlib + `modernc.org/sqlite` + Godog | PASS |

**Re-evaluación post-diseño (Phase 1)**: todos los gates siguen en PASS. Recibir `horas` y
`fecha` como texto en el dominio (R2–R4) no duplica reglas en el handler: el handler no
valida, solo extrae el texto del JSON. El esquema no cambia: la migración 004 copia el DDL ya
acordado, así que no hace falta modificar `MODELO-DE-DATOS.md` en este PR.

## Project Structure

### Documentation (this feature)

```text
specs/019-hu-07-imputacion-horas/
├── plan.md              # Este archivo
├── research.md          # Fase 0: decisiones técnicas (R1–R11)
├── data-model.md        # Fase 1: RegistroHoras, vistas, puertos y services
├── quickstart.md        # Fase 1: guía de validación
├── contracts/
│   └── openapi.yaml     # Fase 1: POST y GET /backlog/{id}/horas
├── checklists/
│   └── requirements.md  # Calidad de la spec
└── tasks.md             # Fase 2 (/speckit.tasks)
```

### Source Code (repository root)

```text
cmd/api/main.go                              # wiring de los 2 services, el handler y las 2 rutas

internal/domain/
├── registro_horas.go                        # RegistroHoras, NewRegistroHoras, MinutosDesdeHoras,
│                                            # ValidarSprintParaImputar, RegistroHorasDetalle, TotalMinutos
├── registro_horas_test.go
├── estado_proyecto.go                       # inicioDelDiaCivil se reutiliza (mismo paquete)
└── errors.go                                # + ErrHistoriaFueraDelSprintActivo, ErrIntegranteAjeno

internal/service/
├── imputar_horas.go                         # ImputarHoras (nuevo)
├── imputar_horas_test.go
├── consultar_esfuerzo_real.go               # ConsultarEsfuerzoReal (nuevo)
├── consultar_esfuerzo_real_test.go
└── fakes_test.go                            # + fakes de RegistroHorasRepository e IntegrantesRepository

internal/repository/
├── migraciones/004_registros_horas.sql      # DDL de MODELO-DE-DATOS.md §4.8
├── esquema_test.go                          # + casos de registros_horas (§9.6)
├── registro_horas_repository.go             # puertos RegistroHorasRepository, IntegrantesRepository
├── sqlite_registro_horas.go                 # Registrar (INSERT condicional), ListarPorHistoria (JOIN usuarios)
└── sqlite_registro_horas_test.go

internal/http/
├── horas_handler.go                         # Imputar, ConsultarEsfuerzo
├── horas_dto.go                             # imputarHorasRequest (horas como json.RawMessage), respuestas
└── horas_handler_test.go

features/
├── imputacion_horas.feature                 # @HU-07 (nuevo)
├── steps_horas.go                           # pasos Godog (nuevo)
└── features_test.go                         # + InitializeScenarioHoras
```

**Structure Decision**: estructura hexagonal ya establecida en `AGENTS.md`. Las horas tienen
su propio handler (`HorasHandler`) en lugar de sumarse a `BacklogHandler`, que ya tiene cuatro
dependencias (SRP). `escribirError` no cambia: los errores nuevos son `ValidationError` (400)
y `ErrHistoriaNoEncontrada` (404) ya está mapeado.

## Complexity Tracking

> Sin violaciones a la constitución.

Notas de alcance (decisiones documentadas, no violaciones):

| Tema | Decisión | Justificación |
|------|----------|---------------|
| Frontend (formulario, errores, acumulado) | Fuera de este incremento; sub-issue #63 | DoD: la sub-issue de Frontend no frena el cierre de la historia |
| Asignar historias a un Sprint | Los tests siembran `sprint_id` en la base o en el fake | HU-06 no estaba en `main` al planificar; entró después con #156 (Assumptions) |
| Editar o borrar cargas | Fuera de alcance | FR-015, D9 |
| Tope diario por integrante | No se implementa | Clarification: solo tope por carga |
| Métricas (desviación, horas por Sprint) | Fuera de alcance | HU-03 |
| Carrera con el cierre del Sprint | INSERT condicional en una sola sentencia | R5; mismo criterio que el fix #145 |
| Migración 003 ausente | La 004 se aplica igual | El runner aplica las que faltan aunque haya huecos (§9.1); la 003 es de HU-02 |
