---
description: "Tareas de la Base del frontend (Issue #108), partiendo del código ya existente en la rama"
---

# Tasks: Base del frontend (Issue #108)

**Input**: Design documents from `/specs/007-frontend-base/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). El ciclo RED → GREEN → REFACTOR con commits separados aplica al **cliente de la API** (tiene lógica). Las piezas compartidas y el router no llevan commit RED por componente: sus tests van junto con una tarea de commit descriptiva.

**Punto de partida**: la rama ya tiene `compartido/api/client.ts`, las cuatro piezas, un router propio y Vitest (commits `ee4472b` a `ad991fa`). Estas tareas **solo cambian lo necesario**; las piezas, el proxy de Vite y Vitest se conservan.

**Rutas**: todo dentro de `web/`; se corre desde ahí (`cd web`). Verificación de cada fase: `npm run build`, `npm run lint` y `npm test` en verde.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: se puede hacer en paralelo (archivos distintos, sin dependencias)
- **[Story]**: historia de la spec a la que pertenece (US1 a US5)

## Phase 1: Setup (dependencias)

Solo esta base toca `package.json` (FR-022).

- [X] T001 En `web/` instalar `react-router` como dependencia y `@testing-library/react`, `@testing-library/user-event` y `jsdom` como devDependencies (`npm install react-router` y `npm install -D ...`); comprobar que la versión de `react-router` funciona con React 19 y que `npm run build` sigue pasando — `web/package.json`, `web/package-lock.json`
- [X] T002 COMMIT `Agregar react-router, Testing Library y jsdom a web`

---

## Phase 2: Foundational (entorno de pruebas)

**⚠️ CRITICAL**: bloquea las historias que usan pruebas de componentes.

- [X] T003 En `web/vitest.config.ts` (ya existe, solo tiene el plugin de React) agregar `test: { environment: 'jsdom' }`. No tocar `web/vite.config.ts` (proxy con `bypass` de `/backlog`, `/projects`, `/sprints`): cuando existen los dos archivos Vitest usa `vitest.config.ts` e ignora el bloque `test` de `vite.config.ts`. Agregar también `vitest.config.ts` al `include` de `web/tsconfig.node.json` (hoy solo lista `vite.config.ts`) para que `npm run build` revise sus tipos
- [X] T004 Correr `npm test`: los 7 tests actuales (3 del cliente + 4 de rutas) deben seguir en verde bajo jsdom. Si `Response`/`fetch` simulados fallan en jsdom, fijar `// @vitest-environment node` en `web/src/compartido/api/client.test.ts` (research §3)
- [X] T005 COMMIT `Configurar jsdom en Vitest`

**Checkpoint**: `npm test` en verde con el entorno nuevo.

---

## Phase 3: User Story 3 - Llamar a la API y mostrar errores de forma uniforme (Priority: P1)

**Goal**: cliente genérico GET/POST/PUT con `ApiError`, `ConnectionError` y respuesta no interpretable (FR-005 a FR-010; [contracts/cliente-api.md](contracts/cliente-api.md)).

**Independent Test**: `npm test` ejecuta los 7 casos de FR-010 y pasan.

- [X] T006 [US3] RED: ampliar `web/src/compartido/api/client.test.ts` a los 7 casos de FR-010 usando `get`, `post` y `put`: éxito 200 con JSON; 400 con `campo`; 400 sin `campo` (campo `undefined`); 404; 409; error 500 con cuerpo que no es JSON (`ApiError` con ese estado y mensaje "La API respondió con estado 500"); `fetch` rechazado (`ConnectionError`). Agregar: `post` y `put` envían el cuerpo como JSON con `Content-Type: application/json`, y toda llamada envía `Accept: application/json`. Correr `npm test` y confirmar que FALLA
- [X] T007 [US3] COMMIT `RED: agregar tests del cliente de la API (falla)` (solo archivos de test)
- [X] T008 [US3] GREEN: en `web/src/compartido/api/client.ts` agregar `ConnectionError`, las operaciones `get`, `post` y `put` sobre una función base, y lectura tolerante del cuerpo: un cuerpo que no es JSON ya no lanza "respuesta inválida" en errores; si la respuesta es de error y no trae `mensaje` interpretable, usar el estado y el mensaje "La API respondió con estado N"; si es éxito con cuerpo no JSON, `ApiError` con el estado y mensaje genérico; cuerpo vacío devuelve `undefined`. Mantener `apiFetch` solo si algo lo usa; si no, retirarlo. `npm test` en verde
- [X] T009 [US3] COMMIT `GREEN: implementar cliente de la API con ConnectionError y respuesta no interpretable, test en verde`
- [X] T010 [US3] (solo si hay algo que limpiar) REFACTOR de `web/src/compartido/api/client.ts` + COMMIT `REFACTOR: [qué se limpió]`

**Checkpoint**: US3 completa y probada sola.

---

## Phase 4: User Story 2 - Sumar una pantalla sin pisar a nadie (Priority: P1) 🎯 MVP

**Goal**: React Router con dos registros de una línea por feature, placeholders dentro de cada feature, navegación con pantalla activa y "página no encontrada" (FR-011 a FR-015, FR-019, FR-020; [contracts/piezas-y-registros.md](contracts/piezas-y-registros.md)).

**Independent Test**: crear una pantalla de prueba siguiendo el README; fuera de su carpeta cambian solo 2 líneas (una en cada registro).

- [X] T011 [P] [US2] Crear los placeholders, uno por feature, exportando `Component` en `web/src/features/proyecto/index.tsx` (HU-04; cada feature se monta con `<carpeta>/*`, así que puede declarar sus propias sub-rutas sin tocar el registro, por ejemplo `settings` y `edit` en HU-04), `web/src/features/backlog/index.tsx` (HU-01), `web/src/features/sprints/index.tsx` (HU-05) y `web/src/features/estado/index.tsx` (HU-13); mover ahí el contenido y las clases CSS del `FeaturePlaceholder` actual de `web/src/compartido/router/AppRoutes.tsx`, cada uno con su título ("Panel del proyecto", "Product Backlog", "Sprints", "Estado del proyecto") y el proyecto actual leído de la URL (`useParams`). Ninguno importa a otro
- [X] T012 [P] [US2] Crear `web/src/navegacion.ts` con el registro de navegación, una línea por feature relativa al proyecto: `{ ruta: 'proyecto', texto: 'Proyecto' }`, `{ ruta: 'backlog', texto: 'Historias' }`, `{ ruta: 'sprints', texto: 'Sprints' }`, `{ ruta: 'status', texto: 'Estado' }`
- [X] T013 [US2] Crear `web/src/compartido/router/LayoutProyecto.tsx`: lee `:id` con `useParams`, muestra el encabezado, los enlaces de `navegacion.ts` como `NavLink` (la pantalla activa marcada con `aria-current="page"` y la clase `feature-navigation__item--active` ya existente; el enlace de `proyecto` queda activo también en sus sub-rutas) y un `<Outlet />`; y `web/src/compartido/router/NoEncontrada.tsx` con el aviso "Página no encontrada"
- [X] T014 [US2] Crear `web/src/rutas.tsx` con el registro de rutas: ruta `/` (guía mínima: "abrí `/projects/<id>`"), ruta `projects/:id` con `LayoutProyecto`, un `index` que redirige a `proyecto` (`<Navigate to="proyecto" replace />`, lo mantiene la base) y una línea hija por feature con carga diferida y comodín final (`{ path: 'proyecto/*', lazy: () => import('./features/proyecto') }`, `{ path: 'backlog/*', lazy: ... }`, `{ path: 'sprints/*', lazy: ... }`, `{ path: 'status/*', lazy: ... }`) y `{ path: '*', element: <NoEncontrada /> }` al final; el proyecto actual solo viaja en la URL
- [X] T015 [US2] Reemplazar el router propio: `web/src/main.tsx` monta `RouterProvider` con `createBrowserRouter(rutas)`; `web/src/App.tsx` queda como marco (encabezado con enlace a `/`, sin id fijo) que no nombra pantallas. Borrar `web/src/compartido/router/AppRoutes.tsx`, `routes.ts` y `routes.test.ts`
- [X] T016 [US2] Escribir `web/src/rutas.test.tsx` con `createMemoryRouter(rutas, { initialEntries })` y Testing Library: cada una de las cuatro rutas de proyecto muestra su placeholder con el id de la URL; el enlace de la pantalla actual tiene `aria-current="page"`; `/no-existe` y `/projects/1/otra` muestran "Página no encontrada"; el cambio de `id` en la URL cambia el proyecto mostrado (sin estado compartido); `/projects/1` redirige a `/projects/1/proyecto`; `/projects/1/proyecto/settings` sigue mostrando el placeholder de HU-04
- [X] T017 [US2] Impedir que una feature importe a otra (FR-020): primero probar una regla de oxlint (`no-restricted-imports` con `overrides` sobre `src/features/**`) en `web/.oxlintrc.json`; si exige nombrar cada feature, escribir `web/src/features/aislamiento.test.ts` que recorra `src/features/**/*.ts(x)` y falle si un import apunta a otra carpeta de `features/`. Comprobar a mano que agregar un import cruzado hace fallar `npm run lint` o `npm test`, y revertirlo
- [X] T018 [US2] `npm run build`, `npm run lint` y `npm test` en verde
- [X] T019 [US2] COMMIT `Reemplazar router propio por React Router con registros de una línea y placeholders por feature`

**Checkpoint**: US2 completa; `/projects/42/backlog` recargado conserva el proyecto.

---

## Phase 5: User Story 4 - Piezas compartidas (Priority: P2)

**Goal**: tests de las piezas, que se conservan sin cambios (FR-008, FR-016 a FR-018).

**Independent Test**: `npm test` ejecuta los tests de las cuatro piezas y pasan.

- [X] T020 [US4] Escribir `web/src/compartido/componentes/componentes.test.tsx` con Testing Library: (a) `Field` con `label`: al hacer clic en la etiqueta el campo recibe el foco, y con `error` el mensaje aparece junto al campo y el campo queda con `aria-invalid="true"`; (b) un formulario de prueba con una operación en curso muestra `LoadingIndicator` (`role="status"`) y su `Button` queda deshabilitado; (c) `ErrorMessage` con `role="alert"` renderizado arriba del formulario cuando el error no trae `campo`. No modificar las piezas; si un test falla por un defecto real, corregirlo en una tarea aparte
- [X] T021 [US4] `npm test` en verde
- [X] T022 [US4] COMMIT `Agregar tests de las piezas compartidas`

**Checkpoint**: US4 probada sola.

---

## Phase 6: User Story 5 - Documentación y datos de ejemplo (Priority: P2)

**Goal**: README y `curl` usables por quien no conoce Go ni React (FR-024, FR-025).

**Independent Test**: una persona que no escribió la base sigue el README y los `curl` desde cero.

- [X] T023 [P] [US5] Actualizar `web/README.md`: instalar (`npm install`), correr (`npm run dev`), compilar, lint y testear; cómo funciona el proxy (las llamadas con `Accept: application/json` van al backend, las de navegador devuelven `index.html`; el cliente compartido siempre lo envía); React Router en lugar del router propio; pasos para agregar una pantalla (carpeta en `features/<nombre>/` con `index.tsx` que exporta `Component`, una línea en `rutas.tsx`, una línea en `navegacion.ts`, usar solo `compartido/`, no importar otra feature, no tocar `package.json`); regla de revisión: el panel del proyecto (`features/proyecto`) es el único punto de contacto entre HU-04 y HU-13; no se enlazan de otra forma; cómo armar llamadas a un endpoint sobre `get`/`post`/`put`, qué hace cada error y dónde mostrar el aviso sin campo. Escrito sin jerga
- [X] T024 [P] [US5] En `docs/curl-frontend-sprint-2.md` agregar la versión bash (continuación de línea con `\` y JSON entre comillas simples) de cada comando, dejando la versión `cmd` (`^`) como segunda opción claramente rotulada; mantener la indicación de dónde anotar los ids devueltos
- [X] T025 [US5] COMMIT `Actualizar README del frontend y agregar curl en bash`

**Checkpoint**: US5 revisada leyendo los dos documentos de principio a fin.

---

## Phase 7: User Story 1 - Arrancar `web/` limpio y conectado (Priority: P1)

**Goal**: comprobar que el conjunto cumple FR-001 a FR-004 (la mayor parte ya está hecha).

**Independent Test**: [quickstart.md](quickstart.md) pasos 1 a 4.

- [X] T026 [US1] Buscar restos de la plantilla de Vite en `web/` (contador, `react.svg`, `vite.svg`, `hero.png`, textos de demostración) y borrar los archivos y referencias que no se usen
- [X] T027 [US1] Con el backend levantado (`go run ./cmd/api`) y los `curl` de `docs/curl-frontend-sprint-2.md`, correr `npm run dev` y comprobar el quickstart pasos 3 y 4: `curl -H "Accept: application/json" http://localhost:5173/projects/<id>` devuelve JSON, y recargar `http://localhost:5173/projects/<id>/backlog` muestra el placeholder
- [X] T028 [US1] COMMIT `Quitar restos de la plantilla de Vite` (solo si T026 cambió archivos)

---

## Phase 8: Polish & Cross-Cutting Concerns

- [ ] T029 Correr los pasos 1, 2 y 5 a 7 del [quickstart.md](quickstart.md) (los 3 y 4 ya se hicieron en T027), incluida la pantalla de ejemplo que toque solo 2 líneas fuera de su carpeta (descartarla después). Pedir a otra persona del equipo que siga solo `web/README.md` con cronómetro para verificar SC-002 (menos de 15 minutos)
- [X] T030 Confirmar con `git diff main --stat -- web/package.json web/src/compartido` que los cambios en dependencias y en código compartido son solo de esta base, y que ninguna feature importa a otra
- [x] T031 Si el plan o las tareas cambiaron después de la primera corrida, volver a correr `/speckit-analyze` y resolver lo que marque antes de abrir el PR

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → historias → **Polish**.
- **US3 (Phase 3)** y **US2 (Phase 4)** dependen de Phase 2 y son independientes entre sí; US3 va primero porque es el único con ciclo RED/GREEN.
- **US4 (Phase 5)** depende de Phase 2 (jsdom y Testing Library) y de nada más.
- **US5 (Phase 6)**: el README (T023) conviene después de US2 y US3 para describir el código final.
- **US1 (Phase 7)** va después de las demás: verifica el conjunto.
- Dentro de US3: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.

### Parallel Opportunities

- T011 y T012 (placeholders y registro de navegación, archivos distintos).
- T023 y T024 (README y `curl`).
- US3 y US2 pueden hacerse en paralelo por personas distintas tras Phase 2; aun así los commits van en el orden de este archivo.

## Implementation Strategy

1. Setup + Foundational: dependencias y jsdom.
2. US3 (cliente con TDD) y US2 (router) → MVP de la base: lo que las cuatro pantallas necesitan.
3. US4 (tests de piezas), US5 (documentación) y US1 (verificación final).
4. Polish: quickstart completo y `/speckit-analyze`.

## Notes

- Un commit por fase del ciclo TDD en el cliente: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra. El historial actual ya tiene un ciclo previo; los commits nuevos se suman encima.
- Nunca `git push` ni squash desde las tareas; el PR se mergea con "Create a merge commit".
- Las piezas (`Button`, `Field`, `ErrorMessage`, `LoadingIndicator`) y el proxy de Vite no se modifican.
