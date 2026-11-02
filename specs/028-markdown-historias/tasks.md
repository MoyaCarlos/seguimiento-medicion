---

description: "Tareas de Markdown seguro en el Product Backlog (Issue #147): frontend, TDD con commits por fase"

---

# Tasks: Markdown seguro en el Product Backlog

**Input**: Design documents from `/specs/028-markdown-historias/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: OBLIGATORIOS. Cada comportamiento se implementa con RED → COMMIT → GREEN →
COMMIT. La vista previa opcional queda fuera de alcance.

## Phase 1: Setup

- [X] T001 Agregar `react-markdown` a `web/package.json` y actualizar `web/package-lock.json` con `npm install react-markdown` (sin `rehype-raw`)
- [X] T002 COMMIT `Agregar react-markdown para renderizar Markdown seguro`

---

## Phase 2: Foundational — componente Markdown

**Goal**: componente de presentación que interpreta CommonMark sin HTML crudo.

- [X] T003 RED: crear `web/src/features/backlog/Markdown.test.tsx` con: (a) texto plano visible; (b) negritas, listas, código y enlaces comunes; (c) `<script>alert(1)</script>` y `<img src=x onerror="alert(2)">` no insertan nodos ni ejecutan código; (d) enlace `[x](javascript:alert(1))` no queda con `href` ejecutable. Correr `npx vitest run src/features/backlog/Markdown.test.tsx` y confirmar que FALLA
- [X] T004 COMMIT `RED: agregar tests de Markdown seguro (fallan)`
- [X] T005 GREEN: crear `web/src/features/backlog/Markdown.tsx` con `react-markdown`, sin `rehype-raw` ni `dangerouslySetInnerHTML`; correr el test y confirmar que pasa
- [X] T006 COMMIT `GREEN: implementar componente Markdown seguro, tests en verde`
- [X] T007 REFACTOR: no aplica, no quedó nada que limpiar

---

## Phase 3: User Story 1 — Markdown en el listado del backlog (P1)

**Goal**: la descripción y los criterios de aceptación se ven formateados en el listado.

**Independent Test**: `npx vitest run src/features/backlog/Listado.test.tsx` con una
historia que use Markdown en descripción y criterios.

- [X] T008 RED: agregar a `web/src/features/backlog/Listado.test.tsx` un caso con descripción `**Importante:** ...` y criterios con lista, que espere un `<strong>` y un `<li>`; confirmar que FALLA
- [X] T009 COMMIT `RED: agregar test de Markdown en el listado del backlog (falla)`
- [X] T010 GREEN: usar `Markdown` en `web/src/features/backlog/Listado.tsx` para `descripcion` y `criterios_aceptacion`, conservando el mensaje de criterios vacíos; correr `npm test`
- [X] T011 COMMIT `GREEN: mostrar descripción y criterios con Markdown, tests en verde`
- [X] T012 REFACTOR: no aplica, no quedó nada que limpiar

---

## Phase 4: Polish

- [X] T013 Agregar estilos acotados a `.backlog` para el contenido Markdown en `web/src/features/backlog/backlog.css` (párrafos, listas, código, pre, enlaces e imágenes sin desbordar)
- [X] T014 COMMIT `Estilizar el contenido Markdown del backlog`
- [X] T015 Correr `npm test`, `npm run lint` y `npm run build` en `web/`, sin errores
- [X] T016 Actualizar `quickstart.md`/`spec.md` con lo verificado y marcar estas tareas

## Dependencies & Execution Order

- Setup → componente → integración → estilos.
- Cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT).
- No hay backend, `.feature` de Godog ni cambios de esquema: el escenario se automatiza
  como prueba de componente, con el escenario en el nombre de la prueba.

## Resultado de la verificación

- `npx vitest run src/features/backlog/Markdown.test.tsx`: 4/4 en verde.
- `npx vitest run src/features/backlog/Listado.test.tsx`: 17/17 en verde.
- `npm test`: 227/228 en la corrida local; el único fallo fue el timeout de carga de
  `rutas.test.tsx` en Windows (un flake ya documentado). Ese archivo pasa solo (13/13) y
  el CI del PR corre la suite en Linux.
- `npm run lint` y `npm run build`: en verde.
- Commits de la feature: docs de spec, dependencia, RED/GREEN del componente, RED/GREEN del listado, estilos y este cierre de tareas.
