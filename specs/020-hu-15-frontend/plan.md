# Implementation Plan: Product Backlog en pantalla (HU-15 Frontend)

**Branch**: `feature/HU-15-frontend` | **Date**: 2026-10-29 | **Spec**: [spec.md](./spec.md)

**Issue**: #98 (parte de #43)

## Summary

Convertir la ruta `backlog` en la pantalla del Product Backlog de HU-15: listar las
historias que devuelve `GET /projects/{id}/backlog` y editarlas con un formulario
controlado que envía `PUT /backlog/{id}`. La pantalla de creación de HU-01 se conserva
en la sub-ruta `backlog/nueva`, así no se pierde la funcionalidad entregada en el
Sprint 2. No se agregan dependencias ni se cambia la API.

## Technical Context

**Lenguaje/versión**: TypeScript 6, React 19, react-router 8 (Vite 8).

**Pruebas**: Vitest 4 + Testing Library + jsdom. Los escenarios Given-When-Then se
automatizan como pruebas de componente, con el escenario en el nombre de cada prueba
(convención BDD del frontend en `AGENTS.md`).

**Dependencias**: ninguna nueva. Se reutilizan `get`/`put` del cliente compartido y los
componentes `Button`, `Field`, `ErrorMessage` y `LoadingIndicator`.

**Alcance**: listado, edición inline, validación de cliente y mapeo de errores de la API.
Sin paginación, búsqueda, Markdown ni edición de prioridad/valor/estado/Sprint.

## Estructura

```text
web/src/
├── rutas.tsx                              # sin cambios: la feature ya cuelga de backlog/*
├── navegacion.ts                          # sin cambios: "Historias" apunta a backlog
└── features/backlog/
    ├── index.tsx                          # NUEVO: shell de rutas (index → Listado, nueva → Crear)
    ├── crear.tsx                          # MOVIDO desde index.tsx: pantalla de HU-01
    ├── crear.test.tsx                     # MOVIDO desde index.test.tsx
    ├── Listado.tsx                        # NUEVO: consulta, estados y listado
    ├── Listado.test.tsx                   # NUEVO: pruebas de listado y edición
    ├── Edicion.tsx                        # NUEVO: formulario controlado de edición
    ├── edicion.ts                         # NUEVO: tipos, formulario, cuerpo y validación
    ├── edicion.test.ts                    # NUEVO: unitarios de validación y mapeo
    ├── api.ts                             # + listarBacklog, editarHistoria
    ├── historia.ts                        # + tipos de historia y edición
    ├── opciones.ts                        # + escala de Story Points
    └── backlog.css                        # + estilos del listado y la edición
```

`rutas.test.tsx` se ajusta para simular `GET /projects/{id}/backlog` y para cubrir la
sub-ruta `backlog/nueva`.

## Decisiones

| Tema | Decisión | Motivo |
|---|---|---|
| Ruta del listado | `backlog` index | Es la pantalla que esperan HU-06 y la navegación "Historias" |
| Creación de HU-01 | Se mueve a `backlog/nueva` | No se rompe la funcionalidad ya entregada; el listado tiene prioridad |
| Orden | Se respeta el orden de la API | Evita duplicar la regla MoSCoW en el frontend (DRY) |
| Edición | Formulario inline, una historia por vez | Menos navegación y menos estado global; patrón de las otras pantallas |
| Validación | Cliente para avisar antes de llamar; API como fuente de verdad | Mejor experiencia sin duplicar el mensaje final de la API |
| Story Points | Selector con "Sin estimar" + escala Fibonacci | Evita valores inválidos y refleja la regla del dominio |
| Criterios vacíos | Se envían como `null` | Contrato del backend |
| Errores | `ApiError.field` conocido → junto al campo; resto → aviso general | Mismo patrón que HU-01 y Sprints |
| Sin Story Points | "Sin estimar" y `null` | La estimación es opcional (HU-15) |

## Constitución

| Principio | Cumplimiento |
|---|---|
| I. Test-First | RED de pruebas de componente y unitarios antes de la implementación; GREEN después |
| II. SDD | spec → plan → tasks; BDD expresado en nombres de pruebas y trazabilidad en la spec |
| III. BDD | Escenarios de la spec en `Listado.test.tsx` (Vitest/Testing Library) |
| IV. Clean/Hexagonal | La feature consume el adaptador HTTP compartido; la regla vive en la API |
| V. Clean Code/SOLID/KISS/YAGNI | Componentes chicos (`Listado`, `Edicion`), validación pura en `edicion.ts`, sin dependencias nuevas |
| VI. Stack | React + Vite + Vitest ya existentes |

## Verificación

En `web/`:

```text
npm test
npm run lint
npm run build
```

Además, `npm test` debe seguir en verde para las pruebas de HU-01 y de rutas.
