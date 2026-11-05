# Implementation Plan: Dashboard de Métricas (HU-09)

**Branch**: `feature/HU-09-dashboard-metricas` | **Date**: 2026-11-04 | **Spec**: [spec.md](./spec.md)

**Issue**: #69 (parte de #9)

**Input**: Feature specification from `specs/035-hu-09-dashboard-metricas/spec.md`

## Summary

Agregar la sección "Dashboard" a cada proyecto: un gráfico de barras con los Story Points comprometidos y completados por Sprint finalizado, indicadores con los valores del proyecto en conjunto, un estado vacío y los mensajes de error. La pantalla consume `GET /projects/{id}/metricas` (HU-03) sin cambiarlo; no hay cambios de backend ni de esquema. El gráfico usa Recharts, que solo se descarga al entrar a la sección.

## Technical Context

**Language/Version**: TypeScript 6, React 19, react-router 8 (Vite 8).

**Primary Dependencies**: `recharts` 3.x y su peer `react-is` (nuevas, decididas por el equipo en la reunión de definición). Se reutilizan `get` del cliente compartido y los componentes `ErrorMessage` y `LoadingIndicator`.

**Storage**: N/A. Las métricas se calculan al consultar y no se guardan (`MODELO-DE-DATOS.md` §6).

**Testing**: Vitest 4 + Testing Library + jsdom. Los escenarios Given-When-Then se automatizan como pruebas de componente, con el escenario en el nombre de cada prueba (`AGENTS.md`). Las funciones puras llevan TDD con commits RED/GREEN.

**Target Platform**: navegador (SPA servida por el backend Go).

**Project Type**: aplicación web; este cambio toca solo `web/`.

**Performance Goals**: una sola consulta al entrar a la sección. Recharts no debe formar parte de la carga inicial de la aplicación.

**Constraints**: sin cambios de backend ni de esquema; una feature no importa a otra (`web/src/features/aislamiento.test.ts`); sin desborde de página en 320 px (FR-018); valores de HU-03 sin recalcular (FR-017).

**Scale/Scope**: una ruta, una entrada de navegación, un gráfico y siete indicadores, para las decenas de Sprints que puede tener un proyecto como máximo.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Cumplimiento |
|---|---|
| I. Test-First | Cumple. Las funciones puras (`metricas.ts`) van con RED → GREEN en commits propios. El resto de la pantalla se cubre con pruebas de componente escritas antes del código. |
| II. SDD | Cumple. Flujo specify → clarify → plan → tasks → analyze antes de implementar. |
| III. BDD | Cumple por la vía de frontend: pruebas de componente con el escenario en el nombre. Los comportamientos de la API ya están en `features/calculo_metricas.feature` (incluye proyecto sin Sprints finalizados, sin parámetro e inexistente). |
| IV. Clean/Hexagonal | No aplica al dominio: no se toca `internal/`. |
| V. Clean Code, KISS, YAGNI, DRY | Cumple. El cálculo sigue en un único lugar (HU-03); la pantalla solo da formato y escala. Sin endpoint nuevo. |
| VI. Stack | Cumple. La constitución pide que el equipo decida explícitamente toda dependencia nueva, y Recharts se confirmó en la reunión de definición. |

Resultado: sin violaciones, antes y después del diseño.
## Project Structure

### Documentation (this feature)

```text
specs/035-hu-09-dashboard-metricas/
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
web/
├── package.json, package-lock.json   # recharts, react-is
└── src/
    ├── navegacion.ts                 # entrada "Dashboard"
    ├── rutas.tsx, rutas.test.tsx     # ruta dashboard/*, con lazy
    └── features/dashboard/
        ├── index.tsx, index.test.tsx # carga, estado vacío y errores
        ├── api.ts                    # consultarMetricas
        ├── metricas.ts, metricas.test.ts       # tipos y funciones puras
        ├── GraficoVelocidad.tsx, GraficoVelocidad.test.tsx
        ├── Indicadores.tsx, Indicadores.test.tsx
        └── dashboard.css
```

**Structure Decision**: una feature nueva, igual que las demás de `web/src/features/`, con su ruta diferida. No toca `internal/`, `features/` ni el esquema.

## Riesgos

- **PR #231 abierto** (HU-19) cambia `rutas.tsx`, `rutas.test.tsx`, `LayoutProyecto.tsx` y `App.css`. Si se mergea antes, hay que actualizar la rama y resolver el conflicto en las rutas. Por eso los estilos del Dashboard van en `dashboard.css`, no en `App.css`.
- **Recharts en jsdom**: no se verificó todavía. jsdom no mide el ancho, así que es probable que haya que simular `ResizeObserver` o fijar las dimensiones del gráfico. La primera tarea lo comprueba antes de seguir (R3).
- **Punto abierto**: los gráficos de horas y defectos por Sprint no se trataron en la reunión (ver Clarifications). Si se suman, son componentes nuevos en la misma feature.

## Complexity Tracking

Sin violaciones que justificar.
