# Implementation Plan: Tarjetas plegables y confirmación de estado (#219)

**Branch**: `feature/219-sprint-tarjetas-confirmacion` | **Date**: 2026-11-15 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/050-sprint-tarjetas-plegables/spec.md`

## Summary

Mejora solo de frontend sobre el Sprint Backlog de HU-11. Cada historia pasa a una tarjeta
con `<details>`/`<summary>`: el resumen muestra título, estado y Story Points, y el detalle
muestra prioridad, fecha de finalización, control de estado y esfuerzo. El cambio de estado
agrega una confirmación en línea con «Confirmar» y «Cancelar», foco en «Cancelar», cierre con
Escape y protección contra el doble envío. Sin dependencias nuevas ni cambios de API o esquema.

## Technical Context

**Language/Version**: TypeScript 6 + React 19, Vite 8, Vitest + Testing Library (jsdom).

**Primary Dependencies**: las existentes; no se agregan dependencias.

**Storage**: no aplica.

**Testing**: pruebas de componente Vitest con el escenario Given-When-Then en el nombre; se
mantienen en verde las pruebas de HU-11 y del esfuerzo.

**Target Platform**: navegador.

**Project Type**: SPA React del monorepo.

**Constraints**: no tocar la API, el esquema ni las transiciones de HU-11; no usar
`window.confirm`; conservar la nota y el bloqueo por proyecto cerrado.

**Scale/Scope**: `SprintBacklog.tsx`, `EstadoHistoria.tsx`, `sprints.css` y sus pruebas.

## Constitution Check

| Principio | Estado |
|---|---|
| I. Test-First | RED/GREEN por comportamiento y commits por fase en [tasks.md](tasks.md) |
| II. SDD | Spec, plan, research, contrato, data-model y quickstart antes de implementar |
| III. BDD | Escenarios de componente con Given-When-Then en el nombre; sin cambios de API |
| IV. Clean/Hexagonal | Solo frontend; las reglas de estado siguen en el backend de HU-11 |
| V. SOLID/KISS/YAGNI | `<details>` nativo y reutilización del patrón de confirmación existente |
| VI. Stack | Sin dependencias nuevas |

Sin violaciones.

## Project Structure

### Documentation (this feature)

```text
specs/050-sprint-tarjetas-plegables/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── ui.md
└── tasks.md
```

### Source Code (repository root)

```text
web/src/features/sprints/
├── SprintBacklog.tsx        # tarjetas <details>/<summary>
├── SprintBacklog.test.tsx   # + plegado y acciones desplegadas
├── EstadoHistoria.tsx       # confirmación en línea y doble envío
├── EstadoHistoria.test.tsx  # + confirmar, cancelar, Escape y clic fuera
└── sprints.css              # estilos de tarjeta, resumen, marca y confirmación
```

**Structure Decision**: la mejora queda dentro de `features/sprints`, sin componentes nuevos:
`SprintBacklog` arma la tarjeta y `EstadoHistoria` conserva la responsabilidad del cambio de
estado, ahora en dos pasos.

## Decisiones técnicas

| # | Decisión | Motivo |
|---|---|---|
| 1 | `<details>`/`<summary>` nativos | Accesibles con teclado, sin estado controlado ni dependencias. |
| 2 | Resumen con título, marca de estado y SP | Permite leer el avance sin desplegar y reduce la altura. |
| 3 | Confirmación en línea en `EstadoHistoria` | Se prueba con Vitest y sigue el patrón de `ZonaCierreProyecto`. |
| 4 | Foco inicial en «Cancelar» y Escape | Evita confirmar por accidente con el teclado. |
| 5 | Cerrar la confirmación ante un error de la API | Muestra el mensaje de HU-11 y deja reintentar desde la acción. |
| 6 | Sin `window.confirm` | No es accesible ni testeable como componente. |

## Complexity Tracking

Sin violaciones de la constitución que justificar.
