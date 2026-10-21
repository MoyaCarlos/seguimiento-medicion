# Implementation Plan: Base del frontend (Issue #108)

**Branch**: `feature/frontend-base` | **Date**: 2026-10-20 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/007-frontend-base/spec.md`

**Punto de partida**: la rama ya tiene una implementación parcial en `web/` (commits `RED: ee4472b`, `GREEN: 6b9f83c`, `REFACTOR: cf0b183` y `ad991fa`). Compila, el lint está limpio y pasan 7 tests. Este plan **conserva lo que cumple la spec** y lista solo los cambios necesarios (sección "Estado actual vs. objetivo").

## Summary

Dejar `web/` lista para que HU-04, HU-01, HU-05 y HU-13 se hagan en paralelo. Ya están: las cuatro piezas compartidas, el reenvío de llamadas a la API en `:8080` y Vitest. Falta: reemplazar el router propio por React Router con dos registros de una línea por pantalla, un placeholder dentro de cada carpeta de feature, navegación con enlaces que marquen la pantalla activa, un cliente de la API genérico (GET/POST/PUT) con error de conexión propio y manejo de respuestas no interpretables, los 7 tests del cliente, tests de las piezas con Testing Library + jsdom, una regla que impida que una feature importe a otra, y la versión bash de los `curl`.

## Technical Context

**Language/Version**: TypeScript ~6.0, React 19 (ya instalados en `web/`)

**Primary Dependencies**: `react-router` (nueva); `@testing-library/react`, `@testing-library/user-event` y `jsdom` como devDependencies (nuevas). Aprobadas en `/speckit-clarify` (FR-026). Vitest ya está instalado (`^4.0.0`). Solo esta base toca `package.json`

**Storage**: N/A (el proyecto actual vive en la URL; los datos, en la API)

**Testing**: Vitest. Cliente de la API con `fetch` simulado; piezas y rutas con Testing Library en entorno jsdom

**Target Platform**: Navegador moderno; desarrollo local con el servidor de Vite

**Project Type**: SPA que consume la API REST existente

**Performance Goals**: N/A

**Constraints**: Sin tocar Go; `compartido/` y `package.json` solo los modifica esta base; ninguna feature importa a otra; una línea por feature en cada registro

**Scale/Scope**: 4 pantallas del Sprint 2

## Constitution Check

| Principio | Estado | Nota |
|---|---|---|
| I. Test-First | Cumple | El cambio del cliente de la API sigue RED→GREEN→REFACTOR con commits separados. Las piezas y el router no requieren commit RED por componente (código sin lógica de negocio; se prueban con tests propios) |
| II. SDD | Desvío regularizado | La implementación parcial se hizo antes de `plan` y `tasks`. Se regulariza acá: el plan y las tareas parten de ese código y el `/speckit-analyze` debe chequearlo antes de seguir implementando |
| III. BDD (Godog) | No aplica | La base no agrega comportamiento de negocio; Godog cubre la API en Go. Los criterios se verifican con [quickstart.md](quickstart.md) |
| IV. Clean/Hexagonal | No aplica | Solo frontend; no toca `internal/` |
| V. Clean Code/KISS/YAGNI/DRY | Cumple | Sin estado global, sin librería de datos, sin sistema de diseño; la lectura de errores vive en un solo lugar |
| VI. Stack | Cumple con decisión explícita | React Router y Testing Library/jsdom son dependencias fuera del stack, aprobadas por el equipo el 2026-10-19 |

Sin violaciones; no hay Complexity Tracking.

## Estado actual vs. objetivo

| # | Elemento | Hoy | Acción |
|---|---|---|---|
| — | Piezas `Button`, `Field`, `ErrorMessage`, `LoadingIndicator` | Hechas en `compartido/componentes/` | **Conservar tal cual** (se agregan tests) |
| — | Proxy de Vite con `bypass` | Hecho en `vite.config.ts` | **Conservar** |
| — | Vitest `^4.0.0` | Instalado | **Conservar** |
| 1 | Router | Propio sobre History API (`routes.ts`, `AppRoutes.tsx`, con nombres de pantallas escritos a mano) | **Reemplazar** por React Router y dos registros planos |
| 2 | Placeholders | `FeaturePlaceholder` único en `compartido/` | **Mover**: uno por feature en `features/<nombre>/` |
| 3 | Navegación | Botones con `pushState` | **Reemplazar** por `NavLink` que marca la pantalla activa |
| 4 | Cliente de la API | `apiFetch` genérico; cuerpo no JSON lanza "respuesta inválida" aun en errores | **Ajustar**: GET/POST/PUT, `ConnectionError`, error con estado y mensaje genérico si el cuerpo no se interpreta |
| 5 | Tests del cliente | 3 | **Completar** a los 7 casos de FR-010 |
| 6 | Tests de piezas | Ninguno | **Agregar** (FR-008, FR-017, FR-018) con Testing Library + jsdom |
| 7 | Regla FR-020 | No existe | **Evaluar** regla de lint; si no alcanza, test en `npm test` (ver [research.md](research.md)) |
| 8 | `docs/curl-frontend-sprint-2.md` | Solo sintaxis `cmd` (`^`) | **Agregar** versión bash |
| — | `web/README.md` | Describe el router propio | **Actualizar** (React Router, cómo agregar pantalla) |

## Project Structure

### Documentation (this feature)

```text
specs/007-frontend-base/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── cliente-api.md
│   └── piezas-y-registros.md
├── checklists/requirements.md
└── tasks.md          # lo genera /speckit-tasks
```

### Source Code (repository root)

```text
web/
├── vite.config.ts            # proxy con bypass (se conserva, sin cambios)
├── vitest.config.ts          # + bloque `test` con jsdom
├── package.json              # + react-router, @testing-library/*, jsdom
├── README.md                 # actualizar
└── src/
    ├── main.tsx              # monta <RouterProvider>
    ├── App.tsx               # marco: encabezado + <Outlet>; no conoce pantallas
    ├── rutas.tsx             # registro de rutas: una línea por feature
    ├── navegacion.ts         # registro de navegación: una línea por feature
    ├── compartido/           # solo lo modifica esta base
    │   ├── api/
    │   │   ├── client.ts         # ajustar: GET/POST/PUT + ConnectionError
    │   │   └── client.test.ts    # completar a 7 casos
    │   ├── componentes/          # se conservan; + Button/Field/ErrorMessage/LoadingIndicator tests
    │   └── router/               # nuevo contenido: LayoutProyecto (navegación con NavLink) y NoEncontrada
    └── features/
        ├── proyecto/index.tsx    # placeholder HU-04 (se monta en proyecto/*)
        ├── backlog/index.tsx     # placeholder HU-01 (backlog/*)
        ├── sprints/index.tsx     # placeholder HU-05 (sprints/*)
        └── estado/index.tsx      # placeholder HU-13 (status/*)

docs/curl-frontend-sprint-2.md    # + versión bash
```

Se eliminan `compartido/router/routes.ts`, `routes.test.ts` y `AppRoutes.tsx` (los reemplaza React Router).

**Structure Decision**: dos únicas zonas de código de UI: `compartido/` (de la base) y `features/<nombre>/` (de cada historia). Los nombres de carpeta son descriptivos (`proyecto`, `backlog`, `sprints`, `estado`) y cada historia puede renombrar el suyo antes de mergear la base, no después. Cada feature expone **una sola** pantalla de entrada (`index.tsx`); lo que necesite adentro lo organiza como quiera.

## Design decisions

1. **Proxy con `bypass`, sin prefijo.** Se conserva el reenvío que ya existe: las llamadas con `Accept: application/json` van al backend y las que piden HTML (recargar `/projects/42/backlog`) se resuelven con `index.html`. El cliente de la API envía siempre `Accept: application/json`, que es lo que hace funcionar el `bypass`. Ver [research.md](research.md).
2. **Registros por carga diferida.** Una línea de `rutas.tsx` por feature, por ejemplo `{ path: 'backlog/*', lazy: () => import('./features/backlog') }`, anidada bajo la ruta `projects/:id`. `rutas.tsx` es lo único que "conoce" a las features.
3. **Navegación relativa al proyecto.** `navegacion.ts` lista `{ ruta: 'backlog', texto: 'Historias' }`. El `LayoutProyecto` (en `compartido/router/`) lee el `:id` de la URL, arma los `NavLink` y marca la activa. Así el proyecto actual nunca se guarda en estado compartido (FR-013).
4. **Una línea por feature, aunque tenga varias vistas.** Todas las features se montan en `<carpeta>/*` (`proyecto/*`, `backlog/*`, `sprints/*`, `status/*`), tengan hoy una vista o varias, para que ninguna tenga que editar el registro si más adelante necesita sub-pantallas. Cada una reparte sus sub-pantallas (HU-04: panel, `settings`, `edit`) con rutas anidadas dentro de su carpeta, con un `*` interno para su "no encontrada". La base agrega un `index` que redirige a `proyecto`, así `/projects/:id` abre el panel. React Router prefiere la ruta más específica, así que el comodín de una feature no pisa a las otras ni al "no encontrada" del proyecto.
5. **Entrada sin proyecto.** Como la API no lista proyectos, la ruta `/` muestra una guía mínima ("abrí `/projects/<id>`") y el enlace del encabezado ya no apunta a un id fijo.
6. **Errores como clases.** `ApiError` (estado, mensaje, campo opcional) y `ConnectionError`; la pantalla decide con `instanceof`. El aviso general sin campo sale de `ErrorMessage`, arriba del formulario (clarify).
7. **Endpoints fuera de `compartido/`.** El cliente no conoce rutas de la API (FR-005); cada feature crea sus llamadas sobre `get/post/put` en su carpeta. Los endpoints disponibles se documentan en [contracts/cliente-api.md](contracts/cliente-api.md).
