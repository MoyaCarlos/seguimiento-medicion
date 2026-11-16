# Implementation Plan: Actualización automática de la votación (#225)

**Branch**: `feature/225-actualizacion-automatica-votacion` | **Date**: 2026-11-15 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/051-actualizacion-automatica-votacion/spec.md`

## Summary

Se agrega polling en `useVotacion`: mientras la ronda está abierta, la pestaña visible y no hay
una acción en curso, cada 5 segundos se reconsulta `GET /backlog/{id}/votacion` de forma
silenciosa. La consulta automática no toca el estado de carga, no deshabilita la pantalla y
descarta su respuesta si una acción empezó mientras estaba en curso. Con la pestaña oculta, la
ronda revelada o una acción pendiente el intervalo se detiene. El botón «Actualizar» sigue
igual. Sin cambios de API, esquema ni dependencias.

## Technical Context

**Language/Version**: TypeScript + React 19 + Vite, Vitest + Testing Library (jsdom).

**Primary Dependencies**: las existentes; no se agregan dependencias.

**Storage**: no aplica.

**Testing**: pruebas del hook con `renderHook` y temporizadores falsos, con el escenario
Given-When-Then en el nombre; se mantienen en verde las pruebas actuales de la pantalla.

**Target Platform**: navegador.

**Project Type**: SPA React del monorepo.

**Constraints**: no tocar la API, el esquema ni `MODELO-DE-DATOS.md`; no usar SSE; conservar
«Actualizar» y los mensajes de error actuales.

**Scale/Scope**: `useVotacion.ts` y una prueba nueva; no se tocan los componentes visuales.

## Constitution Check

| Principio | Estado |
|---|---|
| I. Test-First | RED/GREEN con commits por fase en [tasks.md](tasks.md) |
| II. SDD | Spec, plan, research, data-model, contrato y quickstart antes de implementar |
| III. BDD | Escenarios de componente/hook con Given-When-Then en el nombre |
| IV. Clean/Hexagonal | Solo frontend; la consulta sigue usando el mismo cliente de API |
| V. SOLID/KISS/YAGNI | Polling simple, sin SSE ni store global; un efecto acotado al hook |
| VI. Stack | Sin dependencias nuevas |

Sin violaciones.

## Project Structure

### Documentation (this feature)

```text
specs/051-actualizacion-automatica-votacion/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── polling.md
└── tasks.md
```

### Source Code (repository root)

```text
web/src/features/votacion/
├── useVotacion.ts            # + efecto de polling silencioso
└── useVotacion.test.ts       # NUEVO: escenarios de polling con temporizadores falsos
```

**Structure Decision**: el sondeo vive en el hook que ya concentra la consulta y el estado de
la votación; no se agrega un componente ni un servicio nuevo.

## Decisiones técnicas

| # | Decisión | Motivo |
|---|---|---|
| 1 | Polling cada 5 segundos | Simple, suficiente para el TP y sin cambios de backend. |
| 2 | Consulta silenciosa, sin `recarga`/`cargadaPara` | Evita que el indicador de carga parpadee cada 5 segundos. |
| 3 | Pausar con `document.visibilityState !== 'visible'` | No gasta pedidos con la pestaña oculta y deja dormir a Render. |
| 4 | Pausar con la ronda revelada o una acción en curso | No hay nada nuevo que traer y evita carreras con las acciones. |
| 5 | Descartar respuestas viejas con `accionesIniciadas` | Impide que una consulta automática pise el resultado de una acción. |
| 6 | Conservar «Actualizar» | Sigue sirviendo como refresco manual y no rompe la spec 030. |

## Complexity Tracking

Sin violaciones de la constitución que justificar.
