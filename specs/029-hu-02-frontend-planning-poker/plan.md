# Implementation Plan: Pantalla de Planning Poker (HU-02, Frontend)

**Branch**: `feature/HU-02-frontend-planning-poker` | **Date**: 2026-11-02 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/029-hu-02-frontend-planning-poker/spec.md`

## Summary

Una feature nueva del frontend, `web/src/features/votacion/`, con la pantalla de Planning Poker de una historia en la ruta `/projects/:id/votacion/:historiaId`. Consume la API ya implementada en la spec 012 (`/backlog/{id}/votacion...`) y la lista de integrantes (`GET /projects/{id}/members`). La pantalla: selector de quién actúa, cartas Fibonacci, estado «ya votó / falta» por integrante, «Oculto» para votos ajenos, «Iniciar votación» y «Revelar» solo para el `scrum_master`, resultado con promedio exacto y carta sugerida, botón «Actualizar» y mensajes de error del servidor. El Product Backlog suma un botón «Votar» que es un enlace (`Link`) a esa ruta, sin importar la nueva feature. No se toca el backend ni se agregan dependencias ni tablas.

## Technical Context

**Language/Version**: TypeScript + React 19 + Vite (frontend en `web/`), Node 24

**Primary Dependencies**: las que ya usa `web/` (react-router, Testing Library, Vitest). Sin dependencias nuevas.

**Storage**: N/A (el estado vive en el servidor; ver [data-model.md](./data-model.md))

**Testing**: Vitest + Testing Library (jsdom); escenarios Given-When-Then en el nombre de cada prueba (Principio III, v1.1.0). El comportamiento de la API ya está cubierto por Godog (spec 012).

**Target Platform**: SPA en navegador, API REST/JSON del mismo origen (proxy de Vite en desarrollo).

**Project Type**: Aplicación web (solo parte frontend).

**Performance Goals**: respuesta percibida inmediata para una pantalla de pocos integrantes.

**Constraints**: aislamiento entre features (test `aislamiento.test.ts`); sin tocar backend, esquema ni `MODELO-DE-DATOS.md`; un pedido por acción a la vez.

**Scale/Scope**: una pantalla, 3 a 10 integrantes por proyecto.

## Constitution Check

Constitución v1.1.0 (`.specify/memory/constitution.md`).

| Principio | Cómo se cumple | Estado |
|-----------|----------------|--------|
| I. Test-First | Cada fase (cliente de la API, reglas de la vista, componentes, enlace en el backlog) con commit `RED:` y `GREEN:`; `REFACTOR:` solo si hay algo que limpiar | PASS |
| II. SDD | specify → clarify → plan; siguen tasks, analyze, implement | PASS |
| III. BDD | Pruebas de componente Vitest con el escenario Given-When-Then en el nombre; la API sigue cubierta por Godog (spec 012) | PASS |
| IV. Clean/Hexagonal | No aplica al backend (no se toca); en el frontend, lógica de la vista en funciones puras separadas de los componentes | PASS |
| V. Clean Code/SOLID/KISS/YAGNI/DRY | Sin polling, sin store global, sin persistir la elección; el promedio y la carta sugerida los calcula el backend (no se duplica la regla) | PASS |
| VI. Stack | React + Vite + Vitest + Testing Library, sin dependencias nuevas | PASS |

**Re-evaluación post-diseño**: todos los gates siguen en PASS.

## Project Structure

### Documentation (this feature)

```text
specs/029-hu-02-frontend-planning-poker/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── api-votacion.md
├── checklists/requirements.md
└── tasks.md            # lo genera /speckit-tasks
```

### Source Code

```text
web/src/
├── rutas.tsx                         # + ruta lazy 'votacion/:historiaId'
└── features/
    ├── backlog/Listado.tsx           # + enlace «Votar» (Link, sin importar votacion)
    └── votacion/
        ├── index.tsx                 # Component: pantalla, carga, acciones
        ├── api.ts                    # cliente: consultar, iniciar, votar, revelar, integrantes
        ├── votacion.ts               # tipos, CARTAS, reglas puras (puedeRevelar, esScrumMaster, ...)
        ├── errores.ts                # mensajeDeError
        ├── Cartas.tsx                # baraja Fibonacci
        ├── Integrantes.tsx           # lista ya votó / falta / Oculto / valor
        ├── Resultado.tsx             # promedio exacto y carta sugerida
        ├── votacion.css
        └── *.test.ts(x)              # pruebas con el escenario en el nombre
```

**Structure Decision**: feature propia por el aislamiento. `api.ts` declara su propia copia mínima del tipo de integrante (no se importa de `proyecto`), igual que `backlog` declara `SprintParaAsignar`.
