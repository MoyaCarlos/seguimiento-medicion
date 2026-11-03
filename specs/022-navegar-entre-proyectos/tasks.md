---

description: "Tareas de Navegar entre proyectos desde la interfaz (TDD con commits por fase)"
---

# Tasks: Navegar entre proyectos desde la interfaz

**Input**: Design documents from `/specs/022-navegar-entre-proyectos/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/ui.md`, `quickstart.md`

**Tests**: OBLIGATORIOS. La constitución del proyecto (Principio I, NON-NEGOTIABLE) exige TDD estricto con un commit por fase. No son opcionales en este proyecto. La cobertura BDD se cumple por equivalencia con pruebas de componente (Vitest + Testing Library); no se agregan escenarios Godog.

**Organización**: las fases siguen el orden pedido: preparar el fetch simulado, lógica pura, hook, selector, lista e integración, y al final el nombre accesible del `<nav>` y la verificación. Cada comportamiento con lógica tiene su ciclo RED → COMMIT → GREEN → COMMIT. Los commits son tareas propias y no se marcan hasta ejecutar `git commit`.

**Rutas**: los comandos `npm` y las rutas `web/...` se ejecutan desde `web/`; las rutas de `specs/` y los comandos `git` se ejecutan desde la raíz del repo.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: puede ejecutarse en paralelo (archivos distintos, sin dependencias)
- **[Story]**: historia de usuario a la que pertenece la tarea (US1, US2, US3, US4)
- Todas las tareas indican archivos o comandos exactos

## Phase 1: Setup (Shared Infrastructure)

- [x] T001 Agregar en `web/src/rutas.test.tsx` el `fetch` simulado con `beforeEach`/`afterEach`: una `Response` nueva por llamada, `[]` para `GET /projects`, un proyecto completo para `GET /projects/:id` y `[]` para `GET /projects/:id/members`; ejecutar `npm test` desde `web/` y confirmar que la suite existente sigue en verde.
- [x] T002 COMMIT `test(rutas): simular GET /projects`; archivo: `web/src/rutas.test.tsx`.

**Checkpoint**: la suite existente queda verde antes y después de que el encabezado pida el listado. T001-T002 son infraestructura de prueba, no un ciclo RED; los ciclos RED/GREEN empiezan en T003.

---

## Phase 2: Foundational (Blocking Prerequisites)

**⚠️ CRITICAL**: ninguna historia puede empezar hasta completar la lógica pura y el hook.

- [x] T003 [P] RED: crear `web/src/compartido/proyectos/ruta.test.ts` con los casos de `rutaDeSeccion` (`/projects/1/backlog`, `/projects/1/proyecto/edit`, `/projects/1`, `/projects/1/otra` y `/projects/1/votacion/3` → `/projects/2/proyecto`, y barra final) y `etiquetaProyecto` (`Alfa (#7)` y nombres repetidos); crear antes `web/src/compartido/proyectos/ruta.ts` con un stub mínimo (funciones que devuelven un valor vacío) para que el RED falle por aserción y no por import inexistente; ejecutar `npm test` desde `web/` y confirmar que FALLA.
- [x] T004 COMMIT `RED: agregar tests de rutaDeSeccion y etiquetaProyecto (falla)`; archivos: `web/src/compartido/proyectos/ruta.test.ts` y el stub `web/src/compartido/proyectos/ruta.ts`.
- [x] T005 GREEN: implementar `web/src/compartido/proyectos/ruta.ts` con `rutaDeSeccion(pathname, idNuevo)` (conserva el primer tramo solo si está en `navegacion` de `web/src/navegacion.ts`; si no, `proyecto`) y `etiquetaProyecto({ id, nombre })` según `contracts/ui.md`; ejecutar `npm test` desde `web/` y confirmar que pasa.
- [x] T006 COMMIT `GREEN: implementar rutaDeSeccion y etiquetaProyecto, tests en verde`; archivo: `web/src/compartido/proyectos/ruta.ts`.
- [x] T007 [P] RED: crear `web/src/compartido/proyectos/useProyectos.test.tsx` con `renderHook` y `fetch` simulado: estado inicial `cargando`, éxito `listo` conservando el orden, error del cliente, respuesta `200` que no es arreglo y desmontaje antes de resolver; crear antes `api.ts` y `useProyectos.ts` con stubs mínimos (el hook devuelve siempre `cargando`) para que el RED falle por aserción; confirmar que FALLA.
- [x] T008 COMMIT `RED: agregar tests de useProyectos (falla)`; archivos: `useProyectos.test.tsx` y los stubs `api.ts` y `useProyectos.ts` en `web/src/compartido/proyectos/`.
- [x] T009 GREEN: implementar `web/src/compartido/proyectos/api.ts` (`ProyectoResumen` con `id: number` y `nombre: string`, `listarProyectos()` con `get('/projects')`) y `web/src/compartido/proyectos/useProyectos.ts` (unión `cargando | listo | error`, bandera `vigente` y mensaje para respuestas no interpretables); ejecutar `npm test` desde `web/` y confirmar que pasa.
- [x] T010 COMMIT `GREEN: implementar listarProyectos y useProyectos, tests en verde`; archivos: `web/src/compartido/proyectos/api.ts` y `web/src/compartido/proyectos/useProyectos.ts`.

**Checkpoint**: `ruta.ts`, `api.ts` y `useProyectos.ts` listos; el selector y la lista pueden consumirlos.

---

## Phase 3: User Story 1 - Cambiar de proyecto conservando la sección (Priority: P1) 🎯 MVP

**Goal**: el selector del encabezado cambia de proyecto y conserva la sección actual.

**Independent Test**: abrir `/projects/1/backlog` con dos proyectos cargados, elegir el 2 y comprobar `/projects/2/backlog`, el selector marcando el 2 y la pantalla del backlog reiniciada.

- [x] T011 [US1] RED: crear `web/src/compartido/proyectos/SelectorProyecto.test.tsx` con `createMemoryRouter` y `fetch` simulado: combobox accesible con nombre "Proyecto" que se enfoca con Tab y se puede operar con el teclado (`userEvent.tab()` y `selectOptions`), opciones `<nombre> (#<id>)` en orden, proyecto actual seleccionado por `:id`, fallback `Proyecto <id>` cuando el actual no está en la lista (incluida lista vacía), cambio a otro proyecto a la misma sección, desde una ruta fuera del menú (`/projects/1/otra`) a `/projects/2/proyecto`, mismo proyecto sin navegar y dos proyectos con el mismo nombre que abren el elegido por su `#id`; crear antes `SelectorProyecto.tsx` con un stub mínimo (renderiza vacío) para que falle por aserción; confirmar que FALLA.
- [x] T012 [US1] COMMIT `RED: agregar tests de navegación del selector (falla)`; archivos: `SelectorProyecto.test.tsx` y el stub `SelectorProyecto.tsx` en `web/src/compartido/proyectos/`.
- [x] T013 [US1] GREEN: implementar `web/src/compartido/proyectos/SelectorProyecto.tsx` con el `<select>`, opciones, fallback y navegación mediante `rutaDeSeccion`; no agregar todavía el enlace "Nuevo proyecto" ni los estados de carga/error (US2/US4); ejecutar `npm test` desde `web/` y confirmar que pasa.
- [x] T014 [US1] COMMIT `GREEN: implementar navegación del selector, test en verde`; archivo: `web/src/compartido/proyectos/SelectorProyecto.tsx`.
- [x] T015 [US1] RED: agregar a `web/src/rutas.test.tsx` la integración del selector con dos proyectos: cambiar de `/projects/1/backlog` a `/projects/2/backlog`; además, contar las llamadas al `fetch` simulado: `GET /projects` se pide una sola vez al montar el layout y no se vuelve a pedir al cambiar de sección dentro del mismo proyecto; además, cargar un título en el backlog de 1, cambiar a 2 y verificar que el formulario de 2 arranca vacío; confirmar que FALLA.
- [x] T016 [US1] COMMIT `RED: agregar integración de navegación entre proyectos (falla)`; archivo: `web/src/rutas.test.tsx`.
- [x] T017 [US1] GREEN: modificar `web/src/compartido/router/LayoutProyecto.tsx` para montar `<SelectorProyecto />` en el encabezado y `<Outlet key={id} />`; ejecutar `npm test` desde `web/` y confirmar que pasa.
- [x] T018 [US1] COMMIT `GREEN: montar selector y reiniciar por :id, tests en verde`; archivo: `web/src/compartido/router/LayoutProyecto.tsx`.

**Checkpoint**: US1 funciona de forma independiente y el cambio de proyecto reinicia el estado de la pantalla.

---

## Phase 4: User Story 2 - Ir a crear un proyecto desde cualquier proyecto (Priority: P1)

**Goal**: el encabezado de cualquier pantalla de proyecto ofrece "Nuevo proyecto" hacia `/`.

**Independent Test**: desde `/projects/1/backlog`, pulsar "Nuevo proyecto" y llegar a `/` con el formulario de creación; el enlace sigue visible durante la carga y ante error.

- [x] T019 [US2] RED: agregar a `web/src/compartido/proyectos/SelectorProyecto.test.tsx` y a `web/src/rutas.test.tsx` las pruebas del enlace "Nuevo proyecto": link accesible, `href` `/`, visible con el listado cargando o en error y clic desde una pantalla de proyecto que abre `/` con "Crear proyecto"; confirmar que FALLA.
- [x] T020 [US2] COMMIT `RED: agregar tests del enlace Nuevo proyecto (falla)`; archivos: `web/src/compartido/proyectos/SelectorProyecto.test.tsx` y `web/src/rutas.test.tsx`.
- [x] T021 [US2] GREEN: agregar a `web/src/compartido/proyectos/SelectorProyecto.tsx` el enlace "Nuevo proyecto" a `/` como enlace real (no un `Button` ni una acción de botón), fuera de la rama de datos, para que quede visible y operativo en todos los estados; ejecutar `npm test` desde `web/` y confirmar que pasa.
- [x] T022 [US2] COMMIT `GREEN: implementar enlace Nuevo proyecto, test en verde`; archivo: `web/src/compartido/proyectos/SelectorProyecto.tsx`.

**Checkpoint**: US2 funciona de forma independiente.

---

## Phase 5: User Story 3 - Ver y abrir los proyectos desde el inicio (Priority: P1)

**Goal**: `/` muestra el formulario de creación y la lista de proyectos con enlaces a cada panel.

**Independent Test**: abrir `/` con proyectos, comprobar la lista y que un enlace abre `/projects/<id>/proyecto`; con la lista vacía, ver la invitación a crear el primero.

**Trazabilidad**: el escenario 4 de US3 (la lista falla en el inicio y el formulario sigue disponible) no se prueba en esta fase: se verifica en las tareas de lista de US4 (T035-T038).

- [x] T023 [US3] RED: crear `web/src/compartido/proyectos/ListaProyectos.test.tsx` con listas de enlaces a `/projects/<id>/proyecto`, etiquetas `<nombre> (#<id>)`, orden de la API e invitación a crear cuando está vacía; crear antes `ListaProyectos.tsx` con un stub mínimo (renderiza vacío) para que falle por aserción; confirmar que FALLA.
- [x] T024 [US3] COMMIT `RED: agregar tests de ListaProyectos (falla)`; archivos: `ListaProyectos.test.tsx` y el stub `ListaProyectos.tsx` en `web/src/compartido/proyectos/`.
- [x] T025 [US3] GREEN: implementar `web/src/compartido/proyectos/ListaProyectos.tsx` con la lista y el estado vacío; no agregar todavía el indicador de carga ni el mensaje de error (US4); ejecutar `npm test` desde `web/` y confirmar que pasa.
- [x] T026 [US3] COMMIT `GREEN: implementar ListaProyectos, test en verde`; archivo: `web/src/compartido/proyectos/ListaProyectos.tsx`.
- [x] T027 [US3] RED: agregar a `web/src/rutas.test.tsx` la integración del inicio: `/` muestra "Crear proyecto" y la lista, y un enlace de la lista abre el panel; confirmar que FALLA.
- [x] T028 [US3] COMMIT `RED: agregar integración del inicio con lista (falla)`; archivo: `web/src/rutas.test.tsx`.
- [x] T029 [US3] GREEN: crear `web/src/compartido/proyectos/ComposicionInicio.tsx` (layout pathless con `<Outlet />` y `<ListaProyectos />`) y modificar `web/src/rutas.tsx` para convertir la ruta índice en su hija; no tocar `web/src/features/proyecto/crear.tsx`; ejecutar `npm test` desde `web/` y confirmar que pasa.
- [x] T030 [US3] COMMIT `GREEN: componer formulario y lista en el inicio, tests en verde`; archivos: `web/src/rutas.tsx` y `web/src/compartido/proyectos/ComposicionInicio.tsx`.

**Checkpoint**: US3 funciona de forma independiente.

---

## Phase 6: User Story 4 - Estados de la lista que no rompen la pantalla (Priority: P2)

**Goal**: selector y lista muestran carga y error sin impedir la navegación ni la creación.

**Independent Test**: simular `GET /projects` pendiente y fallido; ver los indicadores, el aviso de error, el `<select>` deshabilitado con la única opción `Proyecto <id>` del proyecto actual, la navegación de secciones y el formulario operativos. El estado vacío ya se cubre en US3.

- [x] T031 [US4] RED: agregar a `web/src/compartido/proyectos/SelectorProyecto.test.tsx` los estados `LoadingIndicator` y `ErrorMessage` (con el enlace "Nuevo proyecto" visible) verificando que en ambos el `<select>` está deshabilitado y muestra como única opción `Proyecto <id>` del proyecto actual, y a `web/src/rutas.test.tsx` una prueba con `GET /projects` en error que verifique el aviso en el encabezado y que los enlaces de sección siguen navegando; confirmar que FALLA.
- [x] T032 [US4] COMMIT `RED: agregar tests de estados del selector (falla)`; archivos: `web/src/compartido/proyectos/SelectorProyecto.test.tsx` y `web/src/rutas.test.tsx`.
- [x] T033 [US4] GREEN: actualizar `web/src/compartido/proyectos/SelectorProyecto.tsx` para mostrar `LoadingIndicator` al cargar y `ErrorMessage` ante error, manteniendo el `<select>` visible y deshabilitado con la única opción `Proyecto <id>` del proyecto actual, y habilitándolo con la lista completa al resolverse la carga; conservar el enlace "Nuevo proyecto" y los enlaces de sección operativos; ejecutar `npm test` desde `web/` y confirmar que pasa.
- [x] T034 [US4] COMMIT `GREEN: implementar estados del selector, tests en verde`; archivo: `web/src/compartido/proyectos/SelectorProyecto.tsx`.
- [x] T035 [US4] RED: agregar a `web/src/compartido/proyectos/ListaProyectos.test.tsx` los estados `LoadingIndicator` y `ErrorMessage`, y a `web/src/rutas.test.tsx` una prueba de `/` con `GET /projects` en error que verifique el aviso de la lista y el formulario de creación disponible, y que completarlo y enviarlo (con `POST /projects` simulado) navega al panel del proyecto creado igual que antes (US4 escenario 4); esta prueba cubre el escenario 4 de US3 (la lista falla en el inicio y el formulario sigue disponible), que se verifica aquí y no en US3; confirmar que FALLA.
- [x] T036 [US4] COMMIT `RED: agregar tests de estados de la lista (falla)`; archivos: `web/src/compartido/proyectos/ListaProyectos.test.tsx` y `web/src/rutas.test.tsx`.
- [x] T037 [US4] GREEN: actualizar `web/src/compartido/proyectos/ListaProyectos.tsx` para mostrar `LoadingIndicator` al cargar y `ErrorMessage` ante error, sin ocultar ni deshabilitar el formulario de creación de `/` (escenario 4 de US3); ejecutar `npm test` desde `web/` y confirmar que pasa.
- [x] T038 [US4] COMMIT `GREEN: implementar estados de la lista, tests en verde`; archivo: `web/src/compartido/proyectos/ListaProyectos.tsx`.

**Nota de ejecución**: T032 y T036 se commitearon juntos (`RED: agregar tests de estados del selector y de la lista (falla)`) porque comparten `web/src/rutas.test.tsx`; T034 y T038 también (`GREEN: implementar estados del selector y de la lista, tests en verde`).

**Checkpoint**: US4 completa los estados de carga y error sin romper la pantalla.

---

## Phase 6b: User Story 5 - Tarjetas y buscador en el inicio (Priority: P2)

**Goal**: la lista de `/` muestra tarjetas con descripción y fechas, un contador y un buscador por nombre; el inicio va en dos columnas.

**Independent Test**: abrir `/` con varios proyectos, ver las tarjetas y el contador, buscar por nombre sin acentos ni mayúsculas y ver "N de M".

**Orden**: esta fase se ejecuta antes de los cierres T043 y T044 (se numera después para no renumerar). Los RED fallan por aserción: si hace falta, el módulo se crea con un stub mínimo. Los tests de la lista que antes buscaban el enlace por su texto se ajustan en el RED de las tarjetas (el nombre accesible sigue siendo `<nombre> (#<id>)`).

- [x] T045 [US5] RED: crear `web/src/compartido/proyectos/busqueda.test.ts` con los casos de `filtrarProyectos` (coincidencia parcial, mayúsculas, acentos, espacios al inicio y al final, texto vacío o de espacios, solo por nombre, orden conservado, sin coincidencias); crear antes `busqueda.ts` con un stub que devuelva `[]`; confirmar que FALLA por aserción.
- [x] T046 [US5] COMMIT `RED: agregar tests de filtrarProyectos (falla)`; archivos: `busqueda.test.ts` y el stub `busqueda.ts` en `web/src/compartido/proyectos/`.
- [x] T047 [US5] GREEN: implementar `filtrarProyectos` en `web/src/compartido/proyectos/busqueda.ts` (recortar, normalizar NFD sin marcas diacríticas, minúsculas, `includes` sobre el nombre); confirmar que pasa.
- [x] T048 [US5] COMMIT `GREEN: implementar filtrarProyectos, tests en verde`; archivo: `web/src/compartido/proyectos/busqueda.ts`.
- [x] T049 [US5] RED: en `web/src/compartido/proyectos/ListaProyectos.test.tsx` agregar (escenarios US3.5 a US3.8): tarjeta con nombre, `#<id>`, descripción y fechas; la tarjeta entera es el enlace con nombre accesible `<nombre> (#<id>)`; "Sin descripción" y "Sin definir"; `title` con el nombre completo; contador con `N`; ajustar los tests de la lista que miraban el texto del enlace; confirmar que FALLA.
- [x] T050 [US5] COMMIT `RED: agregar tests de tarjetas y contador de ListaProyectos (falla)`; archivo: `web/src/compartido/proyectos/ListaProyectos.test.tsx`.
- [x] T051 [US5] GREEN: ampliar `ProyectoResumen` en `web/src/compartido/proyectos/api.ts` con `descripcion`, `fecha_inicio` y `fecha_fin`, y renderizar tarjetas y contador en `ListaProyectos.tsx` (enlace con `aria-label`, `title`, descripción y fechas con los textos de ausencia del panel); confirmar que pasa.
- [x] T052 [US5] COMMIT `GREEN: mostrar tarjetas y contador en ListaProyectos, tests en verde`; archivos: `web/src/compartido/proyectos/api.ts` y `web/src/compartido/proyectos/ListaProyectos.tsx`.
- [x] T053 [US5] RED: en `ListaProyectos.test.tsx` agregar (escenarios US5.1 a US5.10): buscador "Buscar proyecto" que filtra por nombre, ignora mayúsculas y acentos, recorta espacios, vacío muestra todos, contador "N de M", orden conservado, no filtra por descripción ni id, sin coincidencias ("Ningún proyecto coincide con la búsqueda", "0 de M", campo visible), lista vacía sin buscador, cargando y error sin buscador, sin nuevas llamadas al `fetch`; confirmar que FALLA.
- [x] T054 [US5] COMMIT `RED: agregar tests del buscador de proyectos (falla)`; archivo: `web/src/compartido/proyectos/ListaProyectos.test.tsx`.
- [x] T055 [US5] GREEN: agregar el campo de búsqueda y el estado del filtro a `ListaProyectos.tsx` usando `filtrarProyectos`, contador "N de M" y el mensaje de sin coincidencias; confirmar que pasa.
- [x] T056 [US5] COMMIT `GREEN: agregar buscador por nombre a ListaProyectos, tests en verde`; archivo: `web/src/compartido/proyectos/ListaProyectos.tsx`.
- [ ] T057 [US5] Estilo y composición (FR-028, sin prueba automatizada): dos columnas en `/`, tarjetas, buscador y encabezado de la lista en `web/src/App.css` y `web/src/compartido/proyectos/ComposicionInicio.tsx`; se revisa en el navegador (escritorio y 390 px) y se commitea solo con el OK del Scrum Master.
- [ ] T058 [US5] COMMIT `Estilo: ...` del inicio (tarjetas, buscador y dos columnas) y otro `Estilo: ...` del encabezado del proyecto (selector en píldora y "Nuevo proyecto"); archivos: `web/src/App.css`, `ComposicionInicio.tsx`, `SelectorProyecto.tsx`, `LayoutProyecto.tsx`; actualizar el alcance del estilo en `spec.md` y `research.md`.

---

## Phase 7: Polish & Cross-Cutting Concerns

- [x] T039 RED: agregar a `web/src/rutas.test.tsx` la prueba de que el `<nav>` de secciones de un proyecto se encuentra por rol `navigation` con nombre "Secciones del proyecto" y que ya no existe uno llamado "Historias del Sprint 2"; confirmar que FALLA por aserción (no por import).
- [x] T040 COMMIT `RED: esperar el nombre accesible "Secciones del proyecto" en la navegación (falla)`; archivo: `web/src/rutas.test.tsx`.
- [x] T041 GREEN: cambiar el `aria-label` del `<nav>` en `web/src/compartido/router/LayoutProyecto.tsx` a "Secciones del proyecto"; ejecutar `npm test` desde `web/` y confirmar que pasa.
- [x] T042 COMMIT `GREEN: nombrar la navegación "Secciones del proyecto", test en verde`; archivo: `web/src/compartido/router/LayoutProyecto.tsx`.
- [ ] T043 Ejecutar desde `web/`: `npm test`, `npm run lint` y `npm run build`; verificar que no haya dependencias nuevas ni cambios en `web/src/features/proyecto/crear.tsx`, `web/src/features/proyecto/crear.test.tsx`, `web/src/navegacion.ts`, la API o el esquema.
- [ ] T044 Validar `specs/022-navegar-entre-proyectos/quickstart.md` con el backend: lista en el inicio, selector conservando la sección, enlace "Nuevo proyecto", estados de carga/error y fallback del proyecto actual; confirmar que el PR cerrará #148.

**Checkpoint**: nombre accesible corregido, suite/lint/build verdes y quickstart validado.

---

## Dependencies & Execution Order

- **Setup (T001-T002)** → **Foundational (T003-T010)** → **US1 (T011-T018)** → **US2 (T019-T022)** → **US3 (T023-T030)** → **US4 (T031-T038)** → **Polish (T039-T044)**.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT. Nunca se agrupan test e implementación en el mismo commit.
- US1 depende de Foundational y deja el selector montado.
- US2 depende de US1 porque agrega el enlace al selector ya montado.
- US3 depende de Foundational; su componente puede trabajarse en paralelo con US1/US2, pero la integración en `web/src/rutas.test.tsx` debe esperar a que terminen las ediciones de US1/US2 sobre ese mismo archivo.
- US4 depende de US1, US2 y US3 porque completa los estados del selector y de la lista.
- Polish depende de US4.

### Parallel opportunities

- **T003 y T007** están marcadas `[P]`: son archivos de prueba distintos (`ruta.test.ts` y `useProyectos.test.tsx`) y no comparten código. Los commits T004 y T008 deben seguir siendo separados.
- Después de Foundational, **T023-T026** (`ListaProyectos`) son candidatas a paralelizarse con US1/US2 porque usan archivos distintos; no se marcan `[P]` para respetar el orden pedido. Las tareas de integración **T027-T030** deben esperar para no editar `web/src/rutas.test.tsx` al mismo tiempo.
- Fuera de esos casos no se marcan tareas `[P]`: US2, US4 y las integraciones comparten `SelectorProyecto.tsx`, `ListaProyectos.tsx` o `web/src/rutas.test.tsx`, y el ciclo TDD exige orden RED → GREEN antes de avanzar.

## Implementation Strategy

1. **MVP técnico**: US1 (selector, navegación conservando sección y reinicio por `:id`).
2. **Alcance P1 del issue**: US1 + US2 + US3 (navegación, acceso a crear y lista del inicio).
3. **Calidad visual y de error**: US4 (carga, error y no bloqueo), y luego Polish.
4. No se hace push ni merge desde estas tareas; el PR lo gestiona el equipo con merge commit.

## Notes

- Todo RED falla por aserción, nunca por un import inexistente: el módulo se crea con un stub mínimo en el mismo commit RED.
- Tests obligatorios por constitución; no hay tareas de REFACTOR planificadas porque no hay limpieza prevista. Si durante la implementación aparece algo que limpiar, se agrega una tarea REFACTOR con su commit y los tests en verde.
- No se agregan dependencias, endpoints ni cambios de esquema.
- No se agregan escenarios Godog: la cobertura BDD es por equivalencia con Vitest + Testing Library.
- `Inicio.tsx` ya fue eliminado por el PR #208 (issue #141 cerrado); esta historia no lo toca.
- El PR debe incluir `Closes #148` en la descripción (sin `Closes #141`).
- Frontend: los escenarios Given-When-Then de la spec se automatizan como pruebas de componente (Vitest + Testing Library) con el escenario en el nombre de cada prueba; no hay `.feature` de Godog porque la API no cambia.
