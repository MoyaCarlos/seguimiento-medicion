# Implementation Plan: Imputación de horas en pantalla (HU-07 Frontend)

**Branch**: `feature/HU-07-frontend` | **Date**: 2026-11-02 | **Spec**: [spec.md](./spec.md)

**Issue**: #63 (parte de #7)

**Input**: Feature specification from `specs/029-hu-07-frontend/spec.md`

## Summary

Agregar a cada historia de la pantalla de Sprints su esfuerzo total (con el detalle plegado) y, en las del Sprint activo, un botón "Registrar horas" que abre el formulario ahí mismo. La pantalla consume `GET` y `POST /backlog/{id}/horas` y `GET /projects/{id}/members`, sin cambiar la API. Las reglas de validación siguen solo en el backend: la pantalla muestra sus mensajes junto al campo y agrega apenas dos cosas propias (exigir que se elija un integrante y traducir la coma decimal).

## Technical Context

**Language/Version**: TypeScript 6, React 19, react-router 8 (Vite 8).

**Primary Dependencies**: ninguna nueva. Se reutilizan `get`/`post` del cliente compartido, `clasificarError` de Sprints y los componentes `Button`, `Field`, `ErrorMessage` y `LoadingIndicator`.

**Storage**: N/A en el frontend. Las cargas viven en `registros_horas` (`MODELO-DE-DATOS.md` §4.8), sin cambios.

**Testing**: Vitest 4 + Testing Library + jsdom. Los escenarios Given-When-Then se automatizan como pruebas de componente, con el escenario en el nombre de cada prueba (`AGENTS.md`). Las funciones puras llevan TDD con commits RED/GREEN.

**Target Platform**: navegador (SPA servida por el backend Go).

**Project Type**: aplicación web; este cambio toca solo `web/`.

**Performance Goals**: una consulta `GET /backlog/{id}/horas` por historia visible. Aceptable para el tamaño de un Sprint de este proyecto; no se justifica un endpoint agregado (ver [research.md](./research.md), R2).

**Constraints**: sin cambios de backend ni de esquema; una feature no importa a otra (`web/src/features/aislamiento.test.ts`).

**Scale/Scope**: un botón de formulario, un total y un detalle por historia. Sin edición ni borrado de cargas, sin pantalla ni ruta nueva.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Cumplimiento |
|---|---|
| I. Test-First | Cumple. Las funciones puras (`horas.ts`) y la validación del integrante van con RED → GREEN en commits propios. El resto de la pantalla se cubre con pruebas de componente escritas antes del código. |
| II. SDD | Cumple. Flujo specify → clarify → plan → tasks → analyze antes de implementar. |
| III. BDD | Cumple por la vía de frontend de la v1.1.0: pruebas de componente con el escenario en el nombre. Los comportamientos de la API ya están cubiertos con Godog (spec 019). |
| IV. Clean/Hexagonal | No aplica al dominio: no se toca `internal/`. |
| V. Clean Code, KISS, YAGNI, DRY | Cumple. Las reglas de horas, actividad y fecha no se duplican en el cliente; el total lo calcula el backend y se reconsulta. |
| VI. Stack | Cumple. Sin dependencias nuevas. |

Resultado: sin violaciones, antes y después del diseño.

## Project Structure

### Documentation (this feature)

```text
specs/029-hu-07-frontend/
├── spec.md
├── plan.md              # este archivo
├── research.md          # decisiones técnicas
├── data-model.md        # tipos que consume la pantalla
├── contracts/api.md     # qué usa la pantalla del contrato existente
├── quickstart.md        # validación manual
├── checklists/requirements.md
└── tasks.md             # lo genera /speckit-tasks
```

### Source Code (repository root)

```text
web/src/features/sprints/
├── horas.ts                    # NUEVO: tipos y funciones puras (normalizarHoras, formatearHoras, fechaPropuesta)
├── horas.test.ts               # NUEVO
├── EsfuerzoHistoria.tsx        # NUEVO: total, detalle plegado y botón "Registrar horas"
├── EsfuerzoHistoria.test.tsx   # NUEVO
├── FormularioHoras.tsx         # NUEVO: formulario, integrantes y errores
├── FormularioHoras.test.tsx    # NUEVO
├── api.ts                      # + consultarEsfuerzo, imputarHoras, listarIntegrantes
├── SprintBacklog.tsx           # + renderiza EsfuerzoHistoria en cada historia
├── SprintBacklog.test.tsx      # se ajusta el fetch simulado por URL
├── SeccionSprints.tsx          # + pasa proyecto y fecha de fin cuando el Sprint está Activo
├── index.test.tsx              # + respuesta del esfuerzo y pruebas del botón según el estado del Sprint
└── sprints.css                 # + estilos del esfuerzo y el detalle
```

**Structure Decision**: todo vive en `features/sprints`, porque el formulario se muestra dentro de la tarjeta del Sprint y `aislamiento.test.ts` prohíbe importar otra feature. Por eso `listarIntegrantes` se declara también en `sprints/api.ts` en lugar de importarse de `features/proyecto` (ver [research.md](./research.md), R1). No cambian `rutas.tsx` ni `navegacion.ts`.

**Posible conflicto**: el PR abierto #216 (fix #129) cambia una línea de `SeccionSprints.tsx` y `sprint.ts`. Quien mergee segundo lo resuelve; el cambio no se superpone en la lógica.

## Complexity Tracking

Sin violaciones que justificar.
