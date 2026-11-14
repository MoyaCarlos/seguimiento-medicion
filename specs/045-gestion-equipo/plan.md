# Implementation Plan: Gestión del equipo del proyecto (HU-16)

**Branch**: `feature/175-gestion-equipo` | **Date**: 2026-11-12 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/045-gestion-equipo/spec.md`

## Summary

Quitar integrantes con **baja lógica** (`integrantes.fecha_baja`, D15, migración 009),
reactivarlos con el alta existente de HU-21 reutilizando la misma fila, y proteger al único
Scrum Master. «Integrante» pasa a significar «integrante activo» filtrando en el repositorio
(`ListarIntegrantes`, `ListarPorIntegrante`), así todas las reglas existentes (pertenencia,
rol, votación, horas, reporte) ignoran a quien fue dado de baja sin tocar cada servicio. El
historial de rondas y la vista de votación usan una lectura nueva con los dados de baja para
conservar nombres. La baja, la regla del Scrum Master y el descarte de los votos de rondas no
reveladas van en una sola transacción con `UPDATE` condicional (patrón de #145/#158). Endpoint
nuevo: `DELETE /projects/{id}/members/{usuario_id}` (204). Detalle en
[research.md](research.md).

## Technical Context

**Language/Version**: Go 1.26 (backend); TypeScript + React (Vite) en `web/`

**Primary Dependencies**: `net/http` estándar, `modernc.org/sqlite`, Godog; Vitest + Testing Library

**Storage**: SQLite embebido; migración nueva `009_baja_integrantes.sql` (después de la `008` de #300)

**Testing**: `testing` (dominio, servicio con fakes en memoria, repositorio SQLite, handlers), Godog (`features/gestion_equipo.feature`), Vitest (pantalla del equipo)

**Target Platform**: servidor Linux (Render) y navegador

**Project Type**: web-service + SPA

**Performance Goals**: sin requisitos nuevos (equipos de pocas personas)

**Constraints**: un solo PR con `MODELO-DE-DATOS.md` + migración + test (AGENTS.md); merge de `main` con #311 antes de la migración

**Scale/Scope**: 1 endpoint nuevo, 1 cambio de comportamiento en el alta, 2 lecturas filtradas, 1 lectura nueva, 1 columna

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Cumplimiento |
|---|---|
| I. Test-First | Reglas con RED/GREEN/REFACTOR en commits propios: `Member.ValidarBaja` (dominio), `QuitarIntegrante` (servicio), filtro de activos y reactivación (repositorio en memoria y SQLite), `DarDeBajaIntegrante` (UPDATE condicional y descarte de votos), handler. ✅ |
| II. SDD | specify → clarify (resuelto en la spec) → plan → tasks → analyze → BDD → implement. ✅ |
| III. BDD | `features/gestion_equipo.feature` con los escenarios 2 a 5 (normal, alternativo, límite, error); escenario 1 sigue en `permisos_por_rol.feature`. Frontend con Vitest. ✅ |
| IV. Clean/Hexagonal | Regla del Scrum Master en el dominio; servicio `QuitarIntegrante` depende de `ProjectRepository`; SQL solo en el adaptador. ✅ |
| V. SOLID/KISS/YAGNI/DRY | Un servicio por caso de uso; reactivación sin endpoint nuevo; sin endpoint de cambio de rol (Clarification Q2); filtro de activos en un solo lugar. ✅ |
| VI. Stack | Sin dependencias nuevas. ✅ |
| Modelo de datos (AGENTS.md) | D15 confirmada por el equipo (2026-11-12); `MODELO-DE-DATOS.md` y migración en el mismo PR. ✅ |

**Re-check post-design**: sin violaciones. Una observación: el adaptador SQLite de
`ProjectRepository` borra filas de `votos` dentro de la transacción de la baja (cruza
tablas de otro agregado). Se acepta por atomicidad (FR-014, FR-018), igual que el cierre de
Sprint escribe `historias` y `historias_arrastradas`; se documenta en `MODELO-DE-DATOS.md`.

## Project Structure

### Documentation (this feature)

```text
specs/045-gestion-equipo/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/api-equipo.md
├── checklists/requirements.md
└── tasks.md             # /speckit.tasks
```

### Source Code (repository root)

```text
internal/domain/
├── errors.go                      # ErrIntegranteNoEncontrado, ErrScrumMasterNoSeQuita
└── member.go                      # Member.ValidarBaja()
internal/repository/
├── migraciones/009_baja_integrantes.sql
├── project_repository.go          # ListarIntegrantesHistoricos, DarDeBajaIntegrante
├── project_repository_en_memoria.go
└── sqlite_project.go              # filtros de activos, ON CONFLICT de reactivación, baja transaccional
internal/service/
├── quitar_integrante.go           # nuevo caso de uso
├── consultar_historial_rondas.go  # usa ListarIntegrantesHistoricos
└── consultar_votacion.go          # activos + votantes de la ronda
internal/http/
├── project_handler.go             # QuitarIntegrante (DELETE)
├── rutas.go                       # DELETE /projects/{id}/members/{usuario_id}
└── sprint_handler.go              # escribirError: 404 y 409 nuevos
internal/app/app.go                # wiring
features/
├── gestion_equipo.feature
└── steps_gestion_equipo.go
web/src/features/proyecto/
├── api.ts                         # quitarIntegrante
├── Configuracion.tsx              # botón Quitar con confirmación
└── Configuracion.test.tsx
MODELO-DE-DATOS.md                 # §4.3, §4.7, §5, §7 (D15 cerrada), §8, migraciones
docs/PRODUCT-BACKLOG.md            # HU-16: estado y notas
```

**Structure Decision**: la estructura hexagonal existente; no se agregan paquetes.

## Orden de implementación

1. **Dominio**: errores nuevos y `Member.ValidarBaja` (TDD).
2. **Merge de `main`** con #311 (migración 008). Si no entró, se frena y se avisa.
3. **Esquema**: migración 009 + test de la migración (CHECK de `fecha_baja`) + `MODELO-DE-DATOS.md`.
4. **Repositorio**: filtros de activos, `ListarIntegrantesHistoricos`, reactivación en `AgregarIntegrante`, `DarDeBajaIntegrante` (memoria y SQLite, TDD).
5. **Servicios**: `QuitarIntegrante`; `ConsultarHistorialRondas` y `ConsultarVotacion` con los dados de baja (TDD).
6. **HTTP**: handler, ruta, mapeo de errores (TDD).
7. **BDD**: `gestion_equipo.feature` (se escribe antes de implementar, después de `/speckit.analyze`).
8. **Frontend**: `quitarIntegrante` y botón en `Configuracion.tsx` (Vitest).
9. **Docs**: backlog y guía de curl si corresponde.

## Complexity Tracking

Sin violaciones de la constitución que justificar.
