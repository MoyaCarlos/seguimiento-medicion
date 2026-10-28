---
description: "Tareas de la pantalla de Sprints (HU-05 Frontend, Issue #26): TDD con commits por fase"
---

# Tasks: Pantalla de Sprints (HU-05 Frontend, Issue #26)

**Input**: Design documents from `/specs/011-hu-05-frontend-sprints/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/pantalla-y-api.md, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). RED → COMMIT → GREEN → COMMIT donde hay lógica: mapeo de errores, agrupación por estado, estados de botones y avisos. Las piezas sin lógica (CSS, `SeccionSprints` presentacional, tipos, `api.ts`) llevan una tarea de commit descriptiva, sin prefijo RED/GREEN. BDD con Godog no aplica a frontend.

**Alcance**: solo `web/src/features/sprints/` y `specs/011-hu-05-frontend-sprints/`. Ninguna tarea toca `web/src/compartido/`, `web/package.json`, `web/README.md`, `web/src/rutas.tsx`, `web/src/rutas.test.tsx` ni otras features.

**Rutas**: `web/...` se corre desde `web/` (`cd web`). Las rutas de `specs/` y los `git` se corren desde la raíz del repo.

## Reglas de commit (todas las tareas COMMIT)

- `git add` **siempre por nombre de archivo**. Nunca `git add .`, `git add -A`, `git add docs/` ni `git add specs/`.
- No agregar: `"Proyecto - 2026.pdf"`, `docs/INFORME-REVISION-FIX-19.md`, `docs/INFORME-REVISION-FIX-100.md`, `docs/prototipo-sprint2/`, `docs/reviews/`, `specs/007-frontend-base/analisis.md`.
- Cada mensaje termina con la línea `Co-Authored-By: Claude Sonnet 5.5 <noreply@anthropic.com>`.
- Antes de cada commit, `git status --short` y confirmar que solo están los archivos esperados.
- **Sin `git push` ni apertura de PR en ninguna tarea.**

## Format: `[ID] [P?] [Story] Description`

- **[Story]**: US1 Crear (spec US1) · US2 Iniciar · US3 Cerrar · US4 Ver por estado

## Phase 1: Spec (análisis antes de implementar)

- [X] T001 Correr `/speckit-analyze` sobre `specs/011-hu-05-frontend-sprints/` (spec.md, plan.md, tasks.md) y guardar el informe completo en `specs/011-hu-05-frontend-sprints/analisis.md` — cubre todos los FR
- [X] T002 Revisar `specs/011-hu-05-frontend-sprints/analisis.md`: si hay hallazgos **CRITICAL o HIGH**, **frenar y consultar al Scrum Master** antes de seguir; corregir los artefactos que corresponda y repetir T001. Sin CRITICAL ni HIGH, continuar
- [X] T003 COMMIT de la spec 011, agregando por nombre: `specs/011-hu-05-frontend-sprints/spec.md`, `plan.md`, `research.md`, `data-model.md`, `quickstart.md`, `tasks.md`, `contracts/pantalla-y-api.md`, `checklists/requirements.md` y `analisis.md` (todos bajo `specs/011-hu-05-frontend-sprints/`). Mensaje: `Agregar spec, plan y tareas de HU-05 frontend (Issue #26)`

---

## Phase 2: Setup (tipos y api.ts)

**Purpose**: tipos y llamadas a la API; sin lógica propia (cubiertos por los tests de pantalla de la Phase 5).

- [X] T004 Crear `web/src/features/sprints/sprint.ts` solo con el tipo `Sprint` (`id: number`, `proyecto_id: number`, `sprint_goal?: string | null`, `fecha_inicio?: string | null`, `fecha_fin?: string | null`, `estado: 'Pendiente' | 'Activo' | 'Finalizado'`, según `specs/003-apertura-cierre-sprints/contracts/openapi.yaml`) y la constante `ESTADOS = ['Pendiente', 'Activo', 'Finalizado'] as const` — FR-008
- [X] T005 Crear `web/src/features/sprints/api.ts` con `crearSprint(proyectoId: number)` → `post<Sprint>('/sprints', { proyecto_id: proyectoId })` e `iniciarSprint(id: number, datos: { sprint_goal: string; fecha_inicio: string; fecha_fin: string })` → `post<Sprint>(`/sprints/${id}/iniciar`, datos)`, ambas con `post` de `../../compartido/api/client`. **No** incluir `cerrarSprint` (vive en `BotonCerrar.tsx`, FR-007) — FR-001, FR-003, FR-015
- [X] T006 Cierre de fase: en `web/` correr `npm test`, `npm run lint` y `npm run build`; los tres en verde
- [X] T007 COMMIT `Agregar tipo Sprint y llamadas de crear e iniciar a la API de sprints`; archivos: `web/src/features/sprints/sprint.ts`, `web/src/features/sprints/api.ts`

---

## Phase 3: Mapeo de errores y estado (RED/GREEN)

**Goal**: la lógica pura que decide dónde se ve un error y en qué sección está cada Sprint.

### Clasificación de errores — FR-011, FR-012, FR-013

- [X] T008 [US2] RED: crear `web/src/features/sprints/errores.test.ts` con `clasificarError(error, campos)` (aún inexistente): (a) `new ApiError(400, 'm', 'sprint_goal')` con campos `['sprint_goal','fecha_inicio','fecha_fin']` → `{ campo: 'sprint_goal', mensaje: 'm' }`; (b) `ApiError` con `field: 'estado'` → `{ aviso: 'm' }`; (c) `ApiError` sin `field` → `{ aviso }`; (d) `ApiError` con `field: 'proyecto_id'` → aviso; (e) `new ConnectionError()` → `{ aviso: 'No se pudo conectar con el servidor' }`; (f) `new Error('x')` → `{ aviso: 'x' }`; (g) valor no Error → `{ aviso: 'Ocurrió un error inesperado' }`. Correr `npm test` y confirmar que FALLA
- [X] T009 [US2] COMMIT `RED: agregar tests de clasificación de errores de Sprints (falla)`; archivo: `web/src/features/sprints/errores.test.ts`
- [X] T010 [US2] GREEN: crear `web/src/features/sprints/errores.ts` exportando el tipo `ErrorClasificado = { campo: string; mensaje: string } | { aviso: string }` y `clasificarError(error: unknown, campos: readonly string[])`, importando `ApiError` de `../../compartido/api/client`; `npm test` en verde
- [X] T011 [US2] COMMIT `GREEN: clasificar errores de Sprints en campo o aviso, test en verde`; archivo: `web/src/features/sprints/errores.ts`

### Agrupación y estado — FR-002, FR-008, FR-009, FR-010

- [X] T012 [US4] RED: crear `web/src/features/sprints/sprint.test.ts` para funciones aún inexistentes de `./sprint`: `agruparPorEstado(sprints)` devuelve siempre las tres claves `Pendiente`, `Activo`, `Finalizado` (vacías si no hay Sprints), cada Sprint solo en la lista de su estado y cada lista ordenada por `id` ascendente (entrada desordenada con ids 3, 1, 2); `guardarSprint(lista, sprint)` reemplaza el Sprint de igual `id` o lo agrega si es nuevo, sin mutar la lista original; `cargarSprintsIniciales()` devuelve `[]`. Correr `npm test` y confirmar que FALLA
- [X] T013 [US4] COMMIT `RED: agregar tests de agrupación y guardado de Sprints (falla)`; archivo: `web/src/features/sprints/sprint.test.ts`
- [X] T014 [US4] GREEN: en `web/src/features/sprints/sprint.ts` agregar `agruparPorEstado`, `guardarSprint` y `cargarSprintsIniciales` (único punto de carga inicial, hoy devuelve `[]`); `npm test` en verde
- [X] T015 [US4] COMMIT `GREEN: agrupar Sprints por estado y guardarlos por id, test en verde`; archivo: `web/src/features/sprints/sprint.ts`
- [X] T016 [US2] [US4] (solo si hay algo que limpiar) REFACTOR de `web/src/features/sprints/errores.ts` o `web/src/features/sprints/sprint.ts` + COMMIT `REFACTOR: [qué se limpió]` por nombre de archivo
- [X] T017 Cierre de fase: en `web/` correr `npm test`, `npm run lint` y `npm run build`; los tres en verde

**Checkpoint**: la lógica de errores y de agrupación está probada sin UI.

---

## Phase 4: Componentes de acción (RED/GREEN)

**Goal**: formulario de inicio y botón de cierre, cada uno con su estado de envío, su carga y su aviso propio (clarify: un aviso por Sprint).

### Formulario de inicio — US2: FR-003, FR-004, FR-012, FR-013, FR-014

- [X] T018 [US2] RED: crear `web/src/features/sprints/FormularioInicio.test.tsx` (Vitest + Testing Library, `fetch` simulado con `vi.stubGlobal`, el cliente real) para `FormularioInicio` con props `{ sprint, alActualizar }`: (a) hay campos "Sprint Goal", "Fecha de inicio" y "Fecha de fin" (estos dos con `type="date"`) asociados a su etiqueta y un botón "Iniciar Sprint"; (b) al enviar, hace `POST /sprints/{id}/iniciar` con `{ sprint_goal, fecha_inicio, fecha_fin }` y llama `alActualizar` con el Sprint devuelto; (c) durante el envío el botón está deshabilitado y se ve "Cargando…"; (d) 400 con `campo: 'fecha_fin'` muestra el mensaje junto a "Fecha de fin"; (e) 400 con `campo: 'estado'` muestra el mensaje en un aviso arriba del formulario; (f) fetch rechazado muestra el error de conexión en el aviso; (g) un nuevo envío borra los errores anteriores. Confirmar que FALLA
- [X] T019 [US2] COMMIT `RED: agregar tests del formulario de inicio de Sprint (falla)`; archivo: `web/src/features/sprints/FormularioInicio.test.tsx`
- [X] T020 [US2] GREEN: crear `web/src/features/sprints/FormularioInicio.tsx` con `Field`, `Button`, `ErrorMessage` y `LoadingIndicator` de `../../compartido/componentes/`, `iniciarSprint` de `./api` y `clasificarError` de `./errores` con campos `['sprint_goal','fecha_inicio','fecha_fin']`; los campos con `name` exactamente `sprint_goal`, `fecha_inicio` y `fecha_fin` (para que coincidan con el `campo` que devuelve la API) y un `id` propio por Sprint que incluya el id del Sprint (por ejemplo `sprint-${sprint.id}-sprint_goal`); estado local `enviando`, `aviso` y errores por campo; sin validar fechas ni Goal (FR-011). `npm test` en verde
- [X] T021 [US2] COMMIT `GREEN: implementar formulario de inicio de Sprint, test en verde`; archivo: `web/src/features/sprints/FormularioInicio.tsx`

### Botón Cerrar — US3: FR-005, FR-006, FR-007, FR-013, FR-014

- [X] T022 [US3] RED: crear `web/src/features/sprints/BotonCerrar.test.tsx` para `BotonCerrar` con props `{ sprint, alActualizar }`; `cerrarSprint` no se exporta, así que se prueba solo a través del componente: (a) el clic en "Cerrar Sprint" de un Sprint con `id` 7 hace `POST /sprints/7/cerrar` sin cuerpo y llama `alActualizar` con el Sprint devuelto; (b) durante el envío el botón está deshabilitado y se ve "Cargando…"; (c) un error de la API se muestra en un aviso arriba del botón y el botón vuelve a habilitarse; (d) un error de conexión se muestra en el aviso; (e) un nuevo clic borra el aviso anterior. Confirmar que FALLA
- [X] T023 [US3] COMMIT `RED: agregar tests del botón Cerrar Sprint (falla)`; archivo: `web/src/features/sprints/BotonCerrar.test.tsx`
- [X] T024 [US3] GREEN: crear `web/src/features/sprints/BotonCerrar.tsx` que exporta solo `BotonCerrar`; `cerrarSprint` (`post<Sprint>(`/sprints/${id}/cerrar`, undefined)`) vive sin exportar en el mismo archivo, porque exportar un componente y una función dispara la regla `react/only-export-components` de oxlint; FR-007 se mantiene: la llamada y el botón siguen aislados en un solo archivo, para que un cambio de #117 lo toque solo a él; errores vía `clasificarError(error, [])` (todo va al aviso). `npm test` en verde
- [X] T025 [US3] COMMIT `GREEN: implementar botón Cerrar Sprint aislado, test en verde`; archivo: `web/src/features/sprints/BotonCerrar.tsx`
- [X] T026 [US2] [US3] (solo si hay algo que limpiar) REFACTOR de `web/src/features/sprints/FormularioInicio.tsx` o `web/src/features/sprints/BotonCerrar.tsx` + COMMIT `REFACTOR: [qué se limpió]` por nombre de archivo
- [X] T027 Cierre de fase: en `web/` correr `npm test`, `npm run lint` y `npm run build`; los tres en verde

**Checkpoint**: iniciar y cerrar funcionan como componentes aislados.

---

## Phase 5: Secciones y pantalla (US1 y US4)

**Goal**: la pantalla completa: crear, tres secciones, "Sprint #<id>", texto de limitación y contrato con la base (FR-020).

### Secciones y estilos (sin lógica: sin RED) — FR-008, FR-019

- [X] T028 [P] [US4] Crear `web/src/features/sprints/SeccionSprints.tsx` (presentacional): props `{ estado, sprints, alActualizar }`; muestra el título de la sección (exactamente "Pendiente", "Activo" o "Finalizado"; ningún encabezado se llama "Sprints", que es el del `<h1>` de la pantalla), "No hay Sprints" si está vacía y, por cada Sprint, una tarjeta con "Sprint #<id>" (el prefijo y el id en el mismo texto), su Goal y fechas si existen, y su acción: `FormularioInicio` si está Pendiente, `BotonCerrar` si está Activo, ninguna si está Finalizado — FR-003, FR-005, FR-008, FR-019
- [X] T029 [P] [US4] Crear `web/src/features/sprints/sprints.css` con los estilos propios de la feature (secciones, tarjetas, formulario), sin tocar `web/src/index.css` ni `web/src/App.css` — FR-015
- [X] T030 Cierre de fase parcial: en `web/` correr `npm test`, `npm run lint` y `npm run build`; los tres en verde
- [X] T031 COMMIT `Agregar secciones y estilos de la pantalla de Sprints`; archivos: `web/src/features/sprints/SeccionSprints.tsx`, `web/src/features/sprints/sprints.css`

### Pantalla — US1 Crear y US4 Ver: FR-001, FR-002, FR-009, FR-010, FR-013, FR-014, FR-016, FR-018, FR-019, FR-020

- [X] T032 [US1] [US4] RED: crear `web/src/features/sprints/index.test.tsx` montando `Component` en un `MemoryRouter` con ruta `/projects/:id/sprints` (`/projects/42/sprints`), `fetch` simulado: (a) FR-020: un `<h1>` con el texto exacto "Sprints" y "Proyecto actual:" con el id en un único `<strong>`; ningún otro elemento muestra solo "42"; (b) las tres secciones vacías y el texto "Los Sprints existentes no se listan todavía. Los que ves son los de esta visita y se pierden al recargar." sin "docs/" ni "curl"; (c) "Crear Sprint" hace `POST /sprints` con `{ proyecto_id: 42 }` y muestra "Sprint #<id>" en Pendiente, sin Goal ni fechas; (d) durante la creación el botón está deshabilitado y se ve "Cargando…"; (e) 400 al crear muestra el mensaje en el aviso de la pantalla y no agrega ningún Sprint; (f) flujo completo: crear → iniciar (el Sprint pasa a Activo con Goal y fechas) → cerrar (pasa a Finalizado, sin acciones); (g) el error de un Sprint no aparece en otro; (h) con un id no numérico en la URL (`/projects/abc/sprints`), "Crear Sprint" igualmente llama a la API y se muestra en el aviso el error que devuelve (por ejemplo 400 con `campo: 'proyecto_id'`), sin que la pantalla valide el id. Confirmar que FALLA
- [X] T033 [US1] [US4] COMMIT `RED: agregar tests de la pantalla de Sprints (falla)`; archivo: `web/src/features/sprints/index.test.tsx`
- [X] T034 [US1] [US4] GREEN: reemplazar el placeholder de `web/src/features/sprints/index.tsx` por la pantalla real, **conservando el contrato de FR-020**: `export function Component()`, un `<h1>` con el texto exacto "Sprints" y una línea `Proyecto actual: <strong>{id}</strong>` con el id de `useParams()` en un único `<strong>`, sin ningún otro elemento que muestre solo ese número. Agregar: estado local `sprints` inicializado con `cargarSprintsIniciales()`, botón "Crear Sprint" con `crearSprint(Number(id))`, `LoadingIndicator`/botón deshabilitado, los errores de crear clasificados con `clasificarError(error, [])` y mostrados siempre en el aviso de la pantalla (arriba), el texto fijo de limitación (FR-018), las tres `SeccionSprints` con `agruparPorEstado` y `guardarSprint`, e `import './sprints.css'`. Las secciones se titulan "Pendiente", "Activo" y "Finalizado", y ningún otro encabezado se llama "Sprints" (`rutas.test.tsx` busca ese encabezado y falla si hay dos). Solo se dejan de usar las clases del placeholder en `index.tsx`; no se edita `web/src/App.css` ni `web/src/index.css`. `npm test` en verde
- [X] T035 [US1] [US4] Correr `web/src/rutas.test.tsx` completo (`npm test -- src/rutas.test.tsx`) y confirmar que sigue en verde con la pantalla nueva (en particular el caso de `/projects/42/sprints` con "Sprints" y el `<strong>` del id); luego `npm test` completo
- [X] T036 [US1] [US4] COMMIT `GREEN: implementar pantalla de Sprints (crear y ver por estado), test en verde`; archivo: `web/src/features/sprints/index.tsx`
- [X] T037 [US1] [US4] (solo si hay algo que limpiar) REFACTOR de `web/src/features/sprints/index.tsx` o `web/src/features/sprints/SeccionSprints.tsx` + COMMIT `REFACTOR: [qué se limpió]` por nombre de archivo
- [X] T038 Cierre de fase: en `web/` correr `npm test`, `npm run lint` y `npm run build`; los tres en verde. `web/src/features/aislamiento.test.ts` (parte de `npm test`) confirma que no se importa otra feature — FR-015

**Checkpoint**: las cuatro historias funcionan con tests.

---

## Phase 6: Verificación manual (quickstart.md)

- [X] T039 Levantar el backend (`go run ./cmd/api` desde la raíz) y el frontend (`npm run dev` en `web/`); crear un proyecto con el `curl` de `docs/curl-frontend-sprint-2.md` y anotar su id. No commitear archivos que genere el backend (base de datos); verificar con `git status --short` — referencia: `specs/011-hu-05-frontend-sprints/quickstart.md`
- [X] T040 Ejecutar los pasos 3, 4 y 5 de `specs/011-hu-05-frontend-sprints/quickstart.md` (crear/iniciar/cerrar; errores: segundo Pendiente, fecha de fin anterior, segundo Activo, backend apagado; recarga) y anotar para cada paso OK o el fallo observado — FR-001 a FR-019
- [X] T041 Ejecutar el paso 6 (aislamiento): `git diff --stat main...HEAD` solo muestra `web/src/features/sprints/` y `specs/011-*/`; ningún cambio en `compartido/`, `package.json`, README ni registros — FR-015
- [X] T042 Si algún paso falla, corregirlo con su propio RED/GREEN en `web/src/features/sprints/` (con sus tareas COMMIT) y repetir T040; si todo está OK, guardar el resumen del resultado para la descripción del PR (T046). Cierre de fase: `npm test`, `npm run lint` y `npm run build` en `web/` en verde
- [X] T042a [US4] RED: en `web/src/features/sprints/index.test.tsx`, agregar un test que crea un Sprint en `/projects/1/sprints`, navega con el router a `/projects/2/sprints` y espera las tres secciones vacías, sin "Sprint #<id>", un único `<h1>` "Sprints" y un único `<strong>` con el id (FR-009, FR-020). Confirmar que FALLA y commitear: `RED: agregar test de cambio de proyecto en la pantalla de Sprints (falla)`
- [X] T042b [US4] GREEN: en `web/src/features/sprints/index.tsx`, que `Component` solo lea el `:id` y renderice un componente interno `PantallaSprints` con `key={id}`, para que el estado local se reinicie al cambiar de proyecto sin desmontar la pantalla. `npm test` (incluye `rutas.test.tsx`), `npm run lint` y `npm run build` en verde y commitear: `GREEN: reiniciar los Sprints al cambiar de proyecto, test en verde`

---

## Phase 7: Cierre

- [X] T043 Marcar `[X]` las tareas hechas en `specs/011-hu-05-frontend-sprints/tasks.md` y, si cambió el estado de la spec, la línea `Status` de `specs/011-hu-05-frontend-sprints/spec.md`
- [X] T044 COMMIT `Cerrar tareas de la spec 011 (HU-05 frontend)`; archivos: `specs/011-hu-05-frontend-sprints/tasks.md` (y `specs/011-hu-05-frontend-sprints/spec.md` si se modificó)
- [X] T045 Verificación final en `web/`: `npm test`, `npm run lint`, `npm run build` en verde; `git status --short` solo muestra los archivos sin trackear que no van en este PR; `git log --oneline main..HEAD` muestra los commits `RED:`/`GREEN:` en orden
- [X] T046 Escribir la descripción del PR en `~/Escritorio/PR-26-frontend-sprints.md`: `Closes #26` (el padre #5 queda abierto hasta #117); resumen de la pantalla; resultado de la verificación manual del quickstart (T040 y T041); resultado de `npm test`, `npm run lint` y `npm run build`; recordatorio de merge con "Create a merge commit"; pedido de revisión y aprobación a otra integrante del equipo (**no Carlos**, que es quien escribe el PR); terminar con `🤖 Generated with [Claude Code](https://claude.com/claude-code)`. **No** hacer `git push` ni abrir el PR
- [X] T047 Review de #127 (Candela): sacar de la pantalla la referencia a `docs/curl-frontend-sprint-2.md` (FR-018) con ciclo `RED:`/`GREEN:` sobre `web/src/features/sprints/index.test.tsx` e `index.tsx`; actualizar `spec.md` (Clarifications y FR-018) y `plan.md` (punto 8)

---

## Dependencies & Execution Order

- Phase 1 (T001–T003) → Phase 2 → Phase 3 → Phase 4 → Phase 5 → Phase 6 → Phase 7. T002 frena el avance ante hallazgos CRITICAL o HIGH.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- T008–T011 (errores) y T012–T015 (estado) son independientes entre sí, pero se commitean por separado en orden. T028 y T029 son [P] (archivos distintos).
- Phase 4 depende de Phase 3 (`clasificarError`); Phase 5 depende de Phases 2, 3 y 4.
- Cada cierre de fase corre `npm test`, `npm run lint` y `npm run build` en `web/`.

## Implementation Strategy

1. Spec analizada y commiteada (Phase 1).
2. Tipos y llamadas (Phase 2) → lógica pura con RED/GREEN (Phase 3) → componentes de acción (Phase 4).
3. Pantalla completa conservando FR-020 (Phase 5) → verificación manual (Phase 6) → cierre y descripción del PR (Phase 7).
4. MVP: Phases 1 a 5 con US1 (crear) en verde; US2 y US3 se validan de forma independiente con los tests de sus componentes.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra. Merge a `main` con "Create a merge commit", nunca squash.
- Se agregan tres archivos de test que el plan no listaba (`FormularioInicio.test.tsx`, `BotonCerrar.test.tsx`), dentro de la misma carpeta, para probar los estados de botones y avisos de cada componente.
- Nunca `git push` ni apertura de PR desde las tareas.
