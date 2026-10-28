---

description: "Tareas de Navegar entre proyectos desde la interfaz (TDD con commits por fase)"
---

# Tasks: Navegar entre proyectos desde la interfaz

**Input**: Design documents from `/specs/022-navegar-entre-proyectos/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/ui.md`, `quickstart.md`

**Tests**: OBLIGATORIOS. La constitución del proyecto (Principio I, NON-NEGOTIABLE) exige TDD estricto con un commit por fase. No son opcionales en este proyecto. La cobertura BDD se cumple por equivalencia con pruebas de componente (Vitest + Testing Library); no se agregan escenarios Godog.

**Organización**: las fases siguen el orden pedido: preparar el fetch simulado, lógica pura, hook, selector, lista e integración, y al final eliminar `Inicio.tsx`. Cada comportamiento con lógica tiene su ciclo RED → COMMIT → GREEN → COMMIT. Los commits son tareas propias y no se marcan hasta ejecutar `git commit`.

**Rutas**: los comandos `npm` y las rutas `web/...` se ejecutan desde `web/`; las rutas de `specs/` y los comandos `git` se ejecutan desde la raíz del repo.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: puede ejecutarse en paralelo (archivos distintos, sin dependencias)
- **[Story]**: historia de usuario a la que pertenece la tarea (US1, US2, US3, US4)
- Todas las tareas indican archivos o comandos exactos

## Phase 1: Setup (Shared Infrastructure)

- [ ] T001 Agregar en `web/src/rutas.test.tsx` el `fetch` simulado con `beforeEach`/`afterEach`: una `Response` nueva por llamada, `[]` para `GET /projects`, un proyecto completo para `GET /projects/:id` y `[]` para `GET /projects/:id/members`; ejecutar `npm test` desde `web/` y confirmar que la suite existente sigue en verde.
- [ ] T002 COMMIT `test(rutas): simular GET /projects`; archivo: `web/src/rutas.test.tsx`.

**Checkpoint**: la suite existente queda verde antes y después de que el encabezado pida el listado. T001-T002 son infraestructura de prueba, no un ciclo RED; los ciclos RED/GREEN empiezan en T003.

---

## Phase 2: Foundational (Blocking Prerequisites)

**⚠️ CRITICAL**: ninguna historia puede empezar hasta completar la lógica pura y el hook.

- [ ] T003 [P] RED: crear `web/src/compartido/proyectos/ruta.test.ts` con los casos de `rutaDeSeccion` (`/projects/1/backlog`, `/projects/1/proyecto/edit`, `/projects/1`, sección desconocida y barra final) y `etiquetaProyecto` (`Alfa (#7)` y nombres repetidos); importar desde `./ruta`, ejecutar `npm test` desde `web/` y confirmar que FALLA.
- [ ] T004 COMMIT `RED: agregar tests de rutaDeSeccion y etiquetaProyecto (falla)`; archivo: `web/src/compartido/proyectos/ruta.test.ts`.
- [ ] T005 GREEN: implementar `web/src/compartido/proyectos/ruta.ts` con `rutaDeSeccion(pathname, idNuevo)` y `etiquetaProyecto({ id, nombre })` según `contracts/ui.md`; ejecutar `npm test` desde `web/` y confirmar que pasa.
- [ ] T006 COMMIT `GREEN: implementar rutaDeSeccion y etiquetaProyecto, tests en verde`; archivo: `web/src/compartido/proyectos/ruta.ts`.
- [ ] T007 [P] RED: crear `web/src/compartido/proyectos/useProyectos.test.tsx` con `renderHook` y `fetch` simulado: estado inicial `cargando`, éxito `listo` conservando el orden, error del cliente, respuesta `200` que no es arreglo y desmontaje antes de resolver; confirmar que FALLA.
- [ ] T008 COMMIT `RED: agregar tests de useProyectos (falla)`; archivo: `web/src/compartido/proyectos/useProyectos.test.tsx`.
- [ ] T009 GREEN: implementar `web/src/compartido/proyectos/api.ts` (`ProyectoResumen` con `id: number` y `nombre: string`, `listarProyectos()` con `get('/projects')`) y `web/src/compartido/proyectos/useProyectos.ts` (unión `cargando | listo | error`, bandera `vigente` y mensaje para respuestas no interpretables); ejecutar `npm test` desde `web/` y confirmar que pasa.
- [ ] T010 COMMIT `GREEN: implementar listarProyectos y useProyectos, tests en verde`; archivos: `web/src/compartido/proyectos/api.ts` y `web/src/compartido/proyectos/useProyectos.ts`.

**Checkpoint**: `ruta.ts`, `api.ts` y `useProyectos.ts` listos; el selector y la lista pueden consumirlos.

---

## Phase 3: User Story 1 - Cambiar de proyecto conservando la sección (Priority: P1) 🎯 MVP

**Goal**: el selector del encabezado cambia de proyecto y conserva la sección actual.

**Independent Test**: abrir `/projects/1/backlog` con dos proyectos cargados, elegir el 2 y comprobar `/projects/2/backlog`, el selector marcando el 2 y la pantalla del backlog reiniciada.

- [ ] T011 [US1] RED: crear `web/src/compartido/proyectos/SelectorProyecto.test.tsx` con `createMemoryRouter` y `fetch` simulado: combobox accesible "Proyecto", opciones `<nombre> (#<id>)` en orden, proyecto actual seleccionado por `:id`, fallback `Proyecto <id>` cuando el actual no está en la lista (incluida lista vacía), cambio a otro proyecto a la misma sección y mismo proyecto sin navegar; confirmar que FALLA.
- [ ] T012 [US1] COMMIT `RED: agregar tests de navegación del selector (falla)`; archivo: `web/src/compartido/proyectos/SelectorProyecto.test.tsx`.
- [ ] T013 [US1] GREEN: implementar `web/src/compartido/proyectos/SelectorProyecto.tsx` con el `<select>`, opciones, fallback y navegación mediante `rutaDeSeccion`; no agregar todavía el enlace "Nuevo proyecto" ni los estados de carga/error (US2/US4); ejecutar `npm test` desde `web/` y confirmar que pasa.
- [ ] T014 [US1] COMMIT `GREEN: implementar navegación del selector, test en verde`; archivo: `web/src/compartido/proyectos/SelectorProyecto.tsx`.
- [ ] T015 [US1] RED: agregar a `web/src/rutas.test.tsx` la integración del selector con dos proyectos: cambiar de `/projects/1/backlog` a `/projects/2/backlog`; además, cargar un título en el backlog de 1, cambiar a 2 y verificar que el formulario de 2 arranca vacío; confirmar que FALLA.
- [ ] T016 [US1] COMMIT `RED: agregar integración de navegación entre proyectos (falla)`; archivo: `web/src/rutas.test.tsx`.
- [ ] T017 [US1] GREEN: modificar `web/src/compartido/router/LayoutProyecto.tsx` para montar `<SelectorProyecto />` en el encabezado y `<Outlet key={id} />`; ejecutar `npm test` desde `web/` y confirmar que pasa.
- [ ] T018 [US1] COMMIT `GREEN: montar selector y reiniciar por :id, tests en verde`; archivo: `web/src/compartido/router/LayoutProyecto.tsx`.

**Checkpoint**: US1 funciona de forma independiente y el cambio de proyecto reinicia el estado de la pantalla.

---

## Phase 4: User Story 2 - Ir a crear un proyecto desde cualquier proyecto (Priority: P1)

**Goal**: el encabezado de cualquier pantalla de proyecto ofrece "Nuevo proyecto" hacia `/`.

**Independent Test**: desde `/projects/1/backlog`, pulsar "Nuevo proyecto" y llegar a `/` con el formulario de creación; el enlace sigue visible durante la carga y ante error.

- [ ] T019 [US2] RED: agregar a `web/src/compartido/proyectos/SelectorProyecto.test.tsx` y a `web/src/rutas.test.tsx` las pruebas del enlace "Nuevo proyecto": link accesible, `href` `/`, visible con el listado cargando o en error y clic desde una pantalla de proyecto que abre `/` con "Crear proyecto"; confirmar que FALLA.
- [ ] T020 [US2] COMMIT `RED: agregar tests del enlace Nuevo proyecto (falla)`; archivos: `web/src/compartido/proyectos/SelectorProyecto.test.tsx` y `web/src/rutas.test.tsx`.
- [ ] T021 [US2] GREEN: agregar a `web/src/compartido/proyectos/SelectorProyecto.tsx` el enlace "Nuevo proyecto" a `/` como enlace real (no un `Button` ni una acción de botón), fuera de la rama de datos, para que quede visible y operativo en todos los estados; ejecutar `npm test` desde `web/` y confirmar que pasa.
- [ ] T022 [US2] COMMIT `GREEN: implementar enlace Nuevo proyecto, test en verde`; archivo: `web/src/compartido/proyectos/SelectorProyecto.tsx`.

**Checkpoint**: US2 funciona de forma independiente.

---

## Phase 5: User Story 3 - Ver y abrir los proyectos desde el inicio (Priority: P1)

**Goal**: `/` muestra el formulario de creación y la lista de proyectos con enlaces a cada panel.

**Independent Test**: abrir `/` con proyectos, comprobar la lista y que un enlace abre `/projects/<id>/proyecto`; con la lista vacía, ver la invitación a crear el primero.

**Trazabilidad**: el escenario 4 de US3 (la lista falla en el inicio y el formulario sigue disponible) no se prueba en esta fase: se verifica en las tareas de lista de US4 (T035-T038).

- [ ] T023 [US3] RED: crear `web/src/compartido/proyectos/ListaProyectos.test.tsx` con listas de enlaces a `/projects/<id>/proyecto`, etiquetas `<nombre> (#<id>)`, orden de la API e invitación a crear cuando está vacía; confirmar que FALLA.
- [ ] T024 [US3] COMMIT `RED: agregar tests de ListaProyectos (falla)`; archivo: `web/src/compartido/proyectos/ListaProyectos.test.tsx`.
- [ ] T025 [US3] GREEN: implementar `web/src/compartido/proyectos/ListaProyectos.tsx` con la lista y el estado vacío; no agregar todavía el indicador de carga ni el mensaje de error (US4); ejecutar `npm test` desde `web/` y confirmar que pasa.
- [ ] T026 [US3] COMMIT `GREEN: implementar ListaProyectos, test en verde`; archivo: `web/src/compartido/proyectos/ListaProyectos.tsx`.
- [ ] T027 [US3] RED: agregar a `web/src/rutas.test.tsx` la integración del inicio: `/` muestra "Crear proyecto" y la lista, y un enlace de la lista abre el panel; confirmar que FALLA.
- [ ] T028 [US3] COMMIT `RED: agregar integración del inicio con lista (falla)`; archivo: `web/src/rutas.test.tsx`.
- [ ] T029 [US3] GREEN: modificar `web/src/rutas.tsx` para definir `ComposicionInicio` (layout pathless con `<Outlet />` y `<ListaProyectos />`) y convertir la ruta índice en su hija; no tocar `web/src/features/proyecto/crear.tsx`; ejecutar `npm test` desde `web/` y confirmar que pasa.
- [ ] T030 [US3] COMMIT `GREEN: componer formulario y lista en el inicio, tests en verde`; archivo: `web/src/rutas.tsx`.

**Checkpoint**: US3 funciona de forma independiente.

---

## Phase 6: User Story 4 - Estados de la lista que no rompen la pantalla (Priority: P2)

**Goal**: selector y lista muestran carga y error sin impedir la navegación ni la creación.

**Independent Test**: simular `GET /projects` pendiente y fallido; ver los indicadores, el aviso de error, el `<select>` deshabilitado con la única opción `Proyecto <id>` del proyecto actual, la navegación de secciones y el formulario operativos. El estado vacío ya se cubre en US3.

- [ ] T031 [US4] RED: agregar a `web/src/compartido/proyectos/SelectorProyecto.test.tsx` los estados `LoadingIndicator` y `ErrorMessage` (con el enlace "Nuevo proyecto" visible) verificando que en ambos el `<select>` está deshabilitado y muestra como única opción `Proyecto <id>` del proyecto actual, y a `web/src/rutas.test.tsx` una prueba con `GET /projects` en error que verifique el aviso en el encabezado y que los enlaces de sección siguen navegando; confirmar que FALLA.
- [ ] T032 [US4] COMMIT `RED: agregar tests de estados del selector (falla)`; archivos: `web/src/compartido/proyectos/SelectorProyecto.test.tsx` y `web/src/rutas.test.tsx`.
- [ ] T033 [US4] GREEN: actualizar `web/src/compartido/proyectos/SelectorProyecto.tsx` para mostrar `LoadingIndicator` al cargar y `ErrorMessage` ante error, manteniendo el `<select>` visible y deshabilitado con la única opción `Proyecto <id>` del proyecto actual, y habilitándolo con la lista completa al resolverse la carga; conservar el enlace "Nuevo proyecto" y los enlaces de sección operativos; ejecutar `npm test` desde `web/` y confirmar que pasa.
- [ ] T034 [US4] COMMIT `GREEN: implementar estados del selector, tests en verde`; archivo: `web/src/compartido/proyectos/SelectorProyecto.tsx`.
- [ ] T035 [US4] RED: agregar a `web/src/compartido/proyectos/ListaProyectos.test.tsx` los estados `LoadingIndicator` y `ErrorMessage`, y a `web/src/rutas.test.tsx` una prueba de `/` con `GET /projects` en error que verifique el aviso de la lista y el formulario de creación disponible; esta prueba cubre el escenario 4 de US3 (la lista falla en el inicio y el formulario sigue disponible), que se verifica aquí y no en US3; confirmar que FALLA.
- [ ] T036 [US4] COMMIT `RED: agregar tests de estados de la lista (falla)`; archivos: `web/src/compartido/proyectos/ListaProyectos.test.tsx` y `web/src/rutas.test.tsx`.
- [ ] T037 [US4] GREEN: actualizar `web/src/compartido/proyectos/ListaProyectos.tsx` para mostrar `LoadingIndicator` al cargar y `ErrorMessage` ante error, sin ocultar ni deshabilitar el formulario de creación de `/` (escenario 4 de US3); ejecutar `npm test` desde `web/` y confirmar que pasa.
- [ ] T038 [US4] COMMIT `GREEN: implementar estados de la lista, tests en verde`; archivo: `web/src/compartido/proyectos/ListaProyectos.tsx`.

**Checkpoint**: US4 completa los estados de carga y error sin romper la pantalla.

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T039 Eliminar `web/src/compartido/router/Inicio.tsx` y verificar que no queden importaciones ni referencias; no modificar `web/src/features/proyecto/crear.tsx` ni `web/src/features/proyecto/crear.test.tsx`.
- [ ] T040 COMMIT `chore: eliminar Inicio.tsx sin uso`; archivo eliminado: `web/src/compartido/router/Inicio.tsx`.
- [ ] T041 Ejecutar desde `web/`: `npm test`, `npm run lint` y `npm run build`; verificar que no haya dependencias nuevas ni cambios en `web/src/features/proyecto/crear.tsx`, `web/src/features/proyecto/crear.test.tsx`, la API o el esquema, y que no queden referencias a `Inicio.tsx`.
- [ ] T042 Validar `specs/022-navegar-entre-proyectos/quickstart.md` con el backend: lista en el inicio, selector conservando la sección, enlace "Nuevo proyecto", estados de carga/error y fallback del proyecto actual; confirmar que el PR cerrará #148 y #141.

**Checkpoint**: `Inicio.tsx` eliminado, suite/lint/build verdes y quickstart validado.

---

## Dependencies & Execution Order

- **Setup (T001-T002)** → **Foundational (T003-T010)** → **US1 (T011-T018)** → **US2 (T019-T022)** → **US3 (T023-T030)** → **US4 (T031-T038)** → **Polish (T039-T042)**.
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

- Tests obligatorios por constitución; no hay tareas de REFACTOR planificadas porque no hay limpieza prevista. Si durante la implementación aparece algo que limpiar, se agrega una tarea REFACTOR con su commit y los tests en verde.
- No se agregan dependencias, endpoints ni cambios de esquema.
- No se agregan escenarios Godog: la cobertura BDD es por equivalencia con Vitest + Testing Library.
- `web/README.md` menciona el `Inicio.tsx` eliminado; queda fuera de los archivos permitidos y se registra como deuda documental, sin bloquear esta historia.
- El PR debe incluir `Closes #148` y `Closes #141` en la descripción.
