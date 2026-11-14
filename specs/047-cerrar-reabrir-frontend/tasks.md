---

description: "Tareas de #312: cerrar y reabrir un proyecto, frontend (TDD con commits por fase)"
---

# Tasks: Cerrar y reabrir un proyecto (frontend)

**Input**: Design documents from `specs/047-cerrar-reabrir-frontend/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/ui-contracts.md](contracts/ui-contracts.md), [quickstart.md](quickstart.md)

**Tests**: OBLIGATORIOS (Principio I, NON-NEGOTIABLE): un commit por fase del ciclo RED → GREEN → REFACTOR. Frontend: pruebas de componente con Vitest y Testing Library (`cd web && npm test`), con el escenario de la spec en el nombre de cada prueba (Principio III, constitución 1.1.0).

**Organization**: por user story. Las tres fases de entrega del clarify (F1, F2, F3) son independientes: se puede cortar después de cualquiera y mandar el resto en otro PR.

## Reglas que valen para TODAS las tareas

- **Solo frontend** (FR-012): no tocar el backend, `web/src/App.tsx`, `web/src/App.css`, `web/src/rutas.tsx` ni `web/src/compartido/router/LayoutProyecto.tsx`. Sin dependencias nuevas. Estilos solo en el CSS de cada feature.
- **`web/src/compartido/permisos/permisos.test.tsx` NO se modifica** (SC-006). Las pruebas nuevas van en archivos nuevos; en las pantallas se agregan casos a sus tests existentes sin cambiar los que hay.
- **Un solo pedido por pantalla** (FR-007a): cada pantalla llama UNA vez a `usePermisosProyecto` en su raíz y pasa `puedeGestionar` y `puedeEscribir` como props opcionales con valor por defecto `true`.
- **Mientras el rol o el cierre no se conocen** (cargando o error): todo habilitado; decide el servidor.
- **Una sola nota**: si el proyecto está cerrado se muestra `NotaCerrado`; si no, la `NotaGestion` de HU-21 cuando falta permiso. Nunca las dos.
- **Mensajes del servidor tal cual** (FR-003): «El proyecto ya está cerrado», «Cerrá el Sprint activo antes de cerrar el proyecto», «El proyecto no está cerrado», «El proyecto está cerrado».
- **Validación de cada fase**: `cd web && npm test && npm run lint && npm run build` en verde, más los pasos de [quickstart.md](quickstart.md) de esa fase.
- Cada COMMIT lleva el trailer de coautoría del entorno. Ninguna tarea hace `git push` ni merge.

## Mapa de acciones de escritura (cobertura de FR-007 y SC-004)

Cada acción de escritura de la interfaz, con la tarea (RED / GREEN) que la bloquea con el proyecto cerrado y su fase. Ninguna queda sin bloquear. Si al implementar aparece otra acción que escribe (otro `post` o `put` en un `api.ts` de `web/src/features`), se agrega aquí y en su pantalla antes de cerrar la fase.

| # | Acción de escritura | Pedido | Componente | Tareas | Fase |
|---|---|---|---|---|---|
| 1 | Guardar el proyecto | `PUT /projects/{id}` | `proyecto/Editar.tsx` | T024 / T026 | 1 |
| 2 | Cerrar el proyecto | `POST /projects/{id}/cerrar` | `proyecto/ZonaCierreProyecto.tsx` (se deshabilita por rol, y se reemplaza por «Reabrir» si ya está cerrado) | T019 / T021, T032 / T034 | 1 |
| 3 | Reabrir el proyecto | `POST /projects/{id}/reabrir` | `proyecto/ZonaCierreProyecto.tsx` (única acción que sigue habilitada con el proyecto cerrado) | T028 / T030, T032 / T034 | 1 |
| 4 | Crear un Sprint | `POST /sprints` | `sprints/index.tsx` | T037 / T039 | 2 |
| 5 | Iniciar un Sprint | `POST /sprints/{id}/iniciar` | `sprints/FormularioInicio.tsx` | T037 / T039 | 2 |
| 6 | Cerrar un Sprint | `POST /sprints/{id}/cerrar` | `sprints/BotonCerrar.tsx` | T037 / T039 | 2 |
| 7 | Cargar horas | `POST /backlog/{id}/horas` | `sprints/FormularioHoras.tsx` | T037 / T039 | 2 |
| 8 | Cambiar el estado de una historia | `PUT /backlog/{id}/estado` | `sprints/EstadoHistoria.tsx` | T037 / T039 | 2 |
| 9 | Crear una historia | `POST /backlog` | `backlog/crear.tsx` | T042 / T044 | 3 |
| 10 | Editar una historia | `PUT /backlog/{id}` | `backlog/Listado.tsx` (botón «Editar») y `backlog/Edicion.tsx` | T042 / T044 | 3 |
| 11 | Asignar una historia a un Sprint | `PUT /backlog/{id}/sprint` | `backlog/Listado.tsx` (botón «Asignar a Sprint») y `backlog/AsignarASprint.tsx` | T042 / T044 | 3 |
| 12 | Votar | `POST /backlog/{id}/votacion/votos` | `votacion/Cartas.tsx` vía `Ronda.tsx` | T046 / T048 | 3 |
| 13 | Iniciar la votación (gestión) | `POST /backlog/{id}/votacion` | `votacion/SinVotacion.tsx` | T046 / T048 | 3 |
| 14 | Revelar los votos (gestión) | `POST /backlog/{id}/votacion/revelar` | `votacion/AccionesScrumMaster.tsx` | T046 / T048 | 3 |
| 15 | Registrar la estimación (gestión) | `POST /backlog/{id}/votacion/estimacion` | `votacion/Consenso.tsx` | T046 / T048 | 3 |
| 16 | Abrir una nueva ronda (gestión) | `POST /backlog/{id}/votacion/rondas` | `votacion/Consenso.tsx` | T046 / T048 | 3 |
| 17 | Registrar un defecto | `POST /projects/{id}/defectos` | `defectos/FormularioDefecto.tsx` | T050 / T052 | 3 |
| 18 | Resolver un defecto | `POST /defectos/{id}/resolver` | `defectos/ListadoDefectos.tsx` (botón «Resolver») | T050 / T052 | 3 |
| 19 | Agregar un integrante | `POST /projects/{id}/members` | `proyecto/Configuracion.tsx` | T054 / T056 | 3 |

Fuera de la tabla, a propósito: la autenticación (`/auth/*`, fuera de alcance), la creación de un proyecto nuevo (`POST /projects`, no pertenece a un proyecto cerrado) y los enlaces de navegación («Votar», «Configuración», «Editar» del Panel), que no escriben.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: se puede hacer en paralelo (archivos distintos, sin dependencias pendientes)
- **[Story]**: user story de la spec (US1 a US4)

---

## Phase 1: Setup

- [ ] T001 Verificar la línea base: `cd web && npm test` en verde y confirmar que la rama contiene HU-21 (existen `web/src/compartido/permisos/useRolEnProyecto.ts` y `permisos.test.tsx`). Sin commit.
- [ ] T002 Verificar dónde está definido el estilo `.nota-gestion` (no aparece en los CSS de `web/src`): `grep -rn "nota-gestion" web/src`. Si vive en `App.css` basta reutilizar la clase; no se edita. Anotar el hallazgo en el mensaje del commit de T014 (el que agrega la nota y el aviso).

---

## Phase 2: Foundational (permisos compartidos) — bloquea a todas las user stories

**Goal**: `usePermisosProyecto`, aviso y nota, sin cambiar el comportamiento de HU-21.

### Tipo `fecha_cierre` y `useRolEnProyecto`

- [ ] T003 RED: pruebas nuevas en `web/src/compartido/permisos/useRolEnProyecto.cierre.test.tsx` — con `fetch` simulado: (a) `GET /projects/1` con `fecha_cierre: "2026-11-11"` → estado `listo` con `fechaCierre: "2026-11-11"`; (b) `fecha_cierre: null` → `fechaCierre: null`; (c) sin el campo → `fechaCierre: null`; (d) `fecha_cierre: 5` → estado `error` «No se pudo interpretar el rol en el proyecto». Correr y confirmar que FALLA.
- [ ] T004 COMMIT `RED: useRolEnProyecto devuelve la fecha de cierre (falla)` (solo el archivo de test)
- [ ] T005 GREEN: en `web/src/compartido/permisos/useRolEnProyecto.ts`, el estado `listo` pasa a `{ estado: 'listo'; rol: Rol; fechaCierre: string | null }`; leer `fecha_cierre` del mismo `GET /projects/{id}` (ausente o `null` → `null`; string → tal cual; otro tipo → `error`). Un solo pedido. `npm test` en verde, incluido `permisos.test.tsx` sin modificar.
- [ ] T006 COMMIT `GREEN: useRolEnProyecto devuelve la fecha de cierre, test en verde`

### `usePermisosProyecto` y `useEsScrumMaster`

- [ ] T007 RED: pruebas nuevas en `web/src/compartido/permisos/usePermisosProyecto.test.tsx` (con `renderHook` y `fetch` simulado; verificar siempre que `fetch` se llamó UNA vez). Casos: (a) Scrum Master y proyecto abierto → `{ rol: 'scrum_master', fechaCierre: null, cerrado: false, puedeGestionar: true, puedeEscribir: true }`; (b) Product Builder y abierto → `puedeGestionar: false`, `puedeEscribir: true`; (c) Scrum Master y cerrado → `cerrado: true`, `puedeGestionar: false`, `puedeEscribir: false`; (d) Product Builder y cerrado → ambos `false`; (e) cargando → ambos `true`, `rol: null`, `cerrado: false`, `fechaCierre: null`; (f) error de conexión o 500 → ambos `true` y `cerrado: false`; (g) `no_integrante` → `puedeGestionar: false` (igual que HU-21), `puedeEscribir: true` y `cerrado: false` (`puedeEscribir` depende solo del cierre); (g2) proyecto abierto (`fecha_cierre: null` o ausente) → `cerrado: false`: la nota «El proyecto está cerrado» se muestra SOLO con el cierre conocido, nunca cargando (e), con error (f) ni abierto (g2); (h) `aplicarCierre('2026-11-12')` pasa a cerrado sin otro pedido y `aplicarCierre(null)` lo reabre; (i) cambiar de `proyectoId` no arrastra el cierre del proyecto anterior. Confirmar que FALLA.
- [ ] T008 COMMIT `RED: usePermisosProyecto con rol, cierre y permisos (falla)`
- [ ] T009 GREEN: crear `web/src/compartido/permisos/usePermisosProyecto.ts` según [data-model.md](data-model.md): `cerrado = fechaCierre !== null`; `puedeGestionar = !cerrado && rol !== 'product_builder' && estado !== 'no_integrante'`; `puedeEscribir = !cerrado`; `cerrado` solo es `true` con el estado `listo` y `fecha_cierre` presente; desconocido → `true`; `aplicarCierre` guarda un reemplazo de la fecha atado al `proyectoId` (se descarta si cambia el proyecto). Reescribir `web/src/compartido/permisos/useEsScrumMaster.ts` como envoltorio fino que devuelve `puedeGestionar`. `npm test` en verde, incluido `permisos.test.tsx` sin modificar.
- [ ] T010 COMMIT `GREEN: usePermisosProyecto y useEsScrumMaster como envoltorio, test en verde`

### Aviso y nota (RED y GREEN separados, un commit por fase)

- [ ] T011 RED: escribir `web/src/compartido/permisos/AvisoProyectoCerrado.test.tsx` (con fecha muestra «Proyecto cerrado el AAAA-MM-DD» tal como llega, con `role="status"`; con `null` no renderiza nada) y `web/src/compartido/permisos/NotaCerrado.test.tsx` (muestra «El proyecto está cerrado» y usa el `id` recibido). Correr y confirmar que FALLA (los componentes no existen).
- [ ] T067 COMMIT `RED: aviso de proyecto cerrado y nota de cierre (falla)` (solo los archivos de test). *(ID agregado en el análisis; se ejecuta entre T011 y T012.)*
- [ ] T012 GREEN: crear `web/src/compartido/permisos/AvisoProyectoCerrado.tsx` (props `fechaCierre: string | null`; sin fecha → `null`; con fecha → `<div role="status" className="aviso-cerrado">`) y `web/src/compartido/permisos/NotaCerrado.tsx` (prop `id`; `<p id={id} className="nota-gestion">El proyecto está cerrado</p>`, sin tocar `NotaGestion.tsx`). Tests en verde.
- [ ] T013 Agregar los estilos de `.aviso-cerrado` en `web/src/features/proyecto/estilos.css` (nunca en `App.css`), usando las variables de color que ya usa ese archivo.
- [ ] T014 COMMIT `GREEN: aviso de proyecto cerrado y nota de cierre, test en verde` (incluye el hallazgo de T002 sobre dónde vive `.nota-gestion`)

**Checkpoint**: los permisos compartidos tienen el estado de cierre y HU-21 sigue igual.

---

## Phase 3: User Story 1 — El Scrum Master cierra un proyecto (P1) 🎯 MVP · Fase 1

**Goal**: la zona «Cerrar proyecto» con confirmación en línea dentro de Editar proyecto.

**Independent Test**: proyecto abierto sin Sprint Activo, como Scrum Master: Editar proyecto → cerrar → confirmar; aparece «Proyecto cerrado el AAAA-MM-DD» con estado «Finalizado».

### Cliente HTTP

- [ ] T015 [US1] RED: en `web/src/features/proyecto/api.cierre.test.ts` (archivo nuevo), con `fetch` simulado: `cerrarProyecto('7')` hace `POST /projects/7/cerrar` sin cuerpo y devuelve el proyecto con `fecha_cierre` y `estado: 'Finalizado'`; `reabrirProyecto('7')` hace `POST /projects/7/reabrir` sin cuerpo y devuelve `fecha_cierre: null`. Confirmar que FALLA.
- [ ] T016 [US1] COMMIT `RED: cerrarProyecto y reabrirProyecto del cliente (falla)`
- [ ] T017 [US1] GREEN: en `web/src/features/proyecto/api.ts` agregar `fecha_cierre: string | null` al tipo `Proyecto`, `export const cerrarProyecto = (id: string) => post<Proyecto & { estado: EstadoProyecto }>(...)` y `reabrirProyecto` (misma forma). Revisar que no se rompan los usos de `Proyecto` (`npm run build`). Tests en verde.
- [ ] T018 [US1] COMMIT `GREEN: cerrarProyecto y reabrirProyecto del cliente, test en verde`

### Zona de cierre: cerrar con confirmación

- [ ] T019 [US1] RED: en `web/src/features/proyecto/ZonaCierreProyecto.test.tsx` (archivo nuevo; props `proyectoId`, `rol`, `fechaCierre`, `onCambio`), con `fetch` simulado: (a) «Escenario: abierto y Scrum Master» → botón «Cerrar proyecto» habilitado y ningún pedido; (b) pulsarlo muestra «¿Cerrar el proyecto?» con «Confirmar» y «Cancelar», y el foco está en «Cancelar»; (c) «Cancelar» vuelve al reposo sin enviar nada (SC-002); (d) Escape cancela sin enviar nada; (e) «Confirmar» hace un `POST /projects/7/cerrar`, llama `onCambio('2026-11-11')` y muestra el estado «Finalizado»; (f) mientras el pedido está en curso, «Confirmar» y «Cancelar» están deshabilitados y un segundo clic no envía otro pedido; Escape en curso no cancela; (g) 409 «Cerrá el Sprint activo antes de cerrar el proyecto» → la confirmación se cierra, el mensaje se muestra tal cual junto a la zona y no se llama `onCambio`; (h) 409 «El proyecto ya está cerrado» → ídem; (i) 403 → `AvisoSinPermiso`; (j) 404 → aviso de proyecto inexistente; (k) conexión caída → mensaje de error de conexión y el estado no cambia. Confirmar que FALLA.
- [ ] T020 [US1] COMMIT `RED: zona de cierre con confirmación en línea (falla)`
- [ ] T021 [US1] GREEN: crear `web/src/features/proyecto/ZonaCierreProyecto.tsx` con los estados `reposo` / `confirmando` / `enviando` de [data-model.md](data-model.md): foco con ref y `useEffect` al pasar a `confirmando`; Escape (`onKeyDown` en la sección) cancela salvo en `enviando`; botones `disabled` en curso; errores con `mensajeDeError` y `ApiError.tipo` de `web/src/features/proyecto/errores.ts`, `ErrorMessage`, `AvisoSinPermiso` y `LoadingIndicator`; sección con `aria-labelledby` y título «Cerrar proyecto». El estado («Finalizado») se muestra solo a partir de la respuesta de la acción (spec FR-005). Estilos de la zona en `web/src/features/proyecto/estilos.css`. Tests de T019 en verde.
- [ ] T022 [US1] COMMIT `GREEN: zona de cierre con confirmación en línea, test en verde`
- [ ] T023 [US1] (solo si hay algo que limpiar) REFACTOR de `ZonaCierreProyecto.tsx` + COMMIT `REFACTOR: <qué se limpió>`

### Montar la zona en Editar proyecto

- [ ] T024 [US1] RED: casos nuevos en `web/src/features/proyecto/Editar.test.tsx` (sin cambiar los existentes): (a) «Escenario: Scrum Master abre Editar proyecto» → aparece la zona «Cerrar proyecto»; (b) cerrar y confirmar deja la pantalla con «Proyecto cerrado el AAAA-MM-DD» y estado «Finalizado» sin recargar, y «Guardar cambios» queda deshabilitado con la nota «El proyecto está cerrado» (SC-001, FR-005); (c) con el proyecto ya cerrado al abrir, se ve el aviso y la nota única (sin la nota de HU-21); (d) la nota «El proyecto está cerrado» NO aparece mientras el rol y el cierre se cargan, NI cuando el `GET /projects/{id}` del hook falla, NI en un proyecto abierto (`fecha_cierre: null`), y en esos tres casos «Guardar cambios» sigue habilitado. Confirmar que FALLA.
- [ ] T025 [US1] COMMIT `RED: Editar proyecto monta la zona de cierre y el aviso (falla)`
- [ ] T026 [US1] GREEN: en `web/src/features/proyecto/Editar.tsx`, cambio mínimo: reemplazar `useEsScrumMaster` por UNA llamada a `usePermisosProyecto(Number(id))` (después de la carga del proyecto, como hoy); montar `AvisoProyectoCerrado` y `ZonaCierreProyecto` (`onCambio={aplicarCierre}`); `Guardar cambios` usa `puedeGestionar` y muestra `NotaCerrado` si `cerrado`, si no la `NotaGestion` actual. Los casos previos de `Editar.test.tsx` siguen en verde.
- [ ] T027 [US1] COMMIT `GREEN: Editar proyecto con zona de cierre y aviso, test en verde`

**Checkpoint**: US1 funciona y se prueba sola.

---

## Phase 4: User Story 3 — El Scrum Master reabre un proyecto cerrado (P2) · Fase 1

**Goal**: «Reabrir proyecto» sin confirmación en la misma zona.

**Independent Test**: proyecto cerrado, como Scrum Master: reabrir; el aviso desaparece y Guardar se habilita.

- [ ] T028 [US3] RED: casos nuevos en `ZonaCierreProyecto.test.tsx`: (a) cerrado y Scrum Master → se ve «Proyecto cerrado el AAAA-MM-DD» y el botón «Reabrir proyecto» habilitado, sin botón «Cerrar proyecto»; (b) pulsarlo hace un `POST /projects/7/reabrir` SIN confirmación previa, llama `onCambio(null)` y muestra el estado devuelto; (c) botón deshabilitado mientras se envía; (d) 409 «El proyecto no está cerrado» → mensaje tal cual y no se llama `onCambio`; (e) 403 y 404 → avisos estándar; (f) conexión caída → error y sin cambio. Casos nuevos en `Editar.test.tsx`: (g) tras reabrir, el aviso desaparece y «Guardar cambios» se habilita sin recargar (SC-005). Confirmar que FALLA.
- [ ] T029 [US3] COMMIT `RED: reabrir proyecto desde la zona de cierre (falla)`
- [ ] T030 [US3] GREEN: completar `ZonaCierreProyecto.tsx` con la rama «cerrado» (reabrir sin confirmación, mismo manejo de errores) y verificar que `Editar.tsx` refleja la respuesta vía `aplicarCierre`. Tests en verde.
- [ ] T031 [US3] COMMIT `GREEN: reabrir proyecto desde la zona de cierre, test en verde`

**Checkpoint**: US1 y US3 funcionan.

---

## Phase 5: User Story 4 — El Product Builder no puede cerrar ni reabrir (P2) · Fase 1

**Goal**: la zona se ve deshabilitada con nota para el Product Builder (patrón de HU-21, FR-009).

**Independent Test**: como Product Builder, Editar proyecto muestra la zona sin poder cerrar ni reabrir.

- [ ] T032 [US4] RED: casos nuevos en `ZonaCierreProyecto.test.tsx`: (a) `rol: 'product_builder'` y abierto → «Cerrar proyecto» deshabilitado con `aria-describedby` hacia la nota «Solo el Scrum Master puede cerrar o reabrir el proyecto» y ningún pedido al hacer clic; (b) `product_builder` y cerrado → «Reabrir proyecto» deshabilitado con la misma nota; (c) `rol: null` (desconocido) → habilitado, decide el servidor; (d) un 403 del servidor muestra `AvisoSinPermiso`. Caso nuevo en `Editar.test.tsx`: (e) Product Builder ve la zona deshabilitada con la nota. Confirmar que FALLA.
- [ ] T033 [US4] COMMIT `RED: la zona de cierre se deshabilita para el Product Builder (falla)`
- [ ] T034 [US4] GREEN: en `ZonaCierreProyecto.tsx` deshabilitar los botones y mostrar `NotaGestion` con ese texto cuando `rol === 'product_builder'` (la zona decide por `rol`, no por `puedeGestionar`; ver research D3). Tests en verde.
- [ ] T035 [US4] COMMIT `GREEN: zona de cierre deshabilitada para el Product Builder, test en verde`

### Aviso en el Panel (Fase 1, US2.1; tareas agregadas en el análisis, IDs T063 a T066)

- [ ] T063 [US2] RED: casos nuevos en `web/src/features/proyecto/Panel.test.tsx` (sin cambiar los existentes): (a) proyecto cerrado → se ve «Proyecto cerrado el AAAA-MM-DD» y los datos del Panel siguen visibles; (b) proyecto abierto (`fecha_cierre: null`) → sin aviso; (c) mientras carga el rol y el cierre, o si ese pedido falla → sin aviso; (d) el Panel no deshabilita sus enlaces de navegación («Configuración», «Editar»). Confirmar que FALLA.
- [ ] T064 [US2] COMMIT `RED: el Panel muestra el aviso de proyecto cerrado (falla)`
- [ ] T065 [US2] GREEN: en `web/src/features/proyecto/Panel.tsx`, cambio mínimo: UNA llamada a `usePermisosProyecto(Number(id))` (un `GET /projects/{id}` igual al que ya hacen las demás pantallas; sin pedidos de otro tipo) y `<AvisoProyectoCerrado fechaCierre={fechaCierre} />` arriba. Sin otros cambios en el Panel. Tests en verde, incluidos los existentes.
- [ ] T066 [US2] COMMIT `GREEN: el Panel muestra el aviso de proyecto cerrado, test en verde`

**Checkpoint**: Fase 1 completa. Validación: T036.

- [ ] T036 Fase 1: `cd web && npm test && npm run lint && npm run build` en verde y pasos «Fase 1» de [quickstart.md](quickstart.md) contra el backend real. Marcar T036 como hecha en `tasks.md` y commitear junto con ese cambio: `Validar la Fase 1 de cerrar y reabrir (Refs #312)`.

---

## Phase 6: User Story 2 — Un proyecto cerrado se ve como solo lectura (P1) · Fases 2 y 3

**Goal**: aviso y acciones de escritura deshabilitadas con nota en las pantallas del proyecto. Las lecturas siguen funcionando (FR-008). El aviso en Editar ya quedó en la Fase 1 (T026).

**Independent Test**: abrir un proyecto cerrado y comprobar aviso, lecturas y acciones deshabilitadas con nota, pantalla por pantalla.

### Fase 2 — Sprints (crear, iniciar y cerrar Sprint, horas y estado de las historias)

- [ ] T037 [US2] RED: casos nuevos en `web/src/features/sprints/index.test.tsx` (con el `GET /projects/{id}` de la pantalla devolviendo `fecha_cierre`): (a) «Escenario: proyecto cerrado» → aviso «Proyecto cerrado el AAAA-MM-DD», «Crear Sprint» deshabilitado con `NotaCerrado`, y el listado de Sprints sigue visible; (b) proyecto abierto → sin aviso y todo como antes; (c) Scrum Master con proyecto cerrado → «Iniciar Sprint» y «Cerrar Sprint» deshabilitados con la nota única; (d) hay UNA sola llamada a `GET /projects/{id}` por la pantalla; (e) un 409 «El proyecto está cerrado» al crear se muestra tal cual; (f) con el cierre desconocido (cargando o error) todo queda habilitado y sin nota. Casos en `FormularioHoras.test.tsx`, `EstadoHistoria.test.tsx`, `FormularioInicio.test.tsx` y `BotonCerrar.test.tsx`: cada componente con `puedeEscribir={false}` (o `puedeGestionar={false}` por cierre) deshabilita su acción y no envía; con el valor por defecto se comporta como hoy. Confirmar que FALLA.
- [ ] T038 [US2] COMMIT `RED: Sprints en solo lectura con el proyecto cerrado (falla)`
- [ ] T039 [US2] GREEN: en `web/src/features/sprints/index.tsx` reemplazar `useEsScrumMaster` por UNA llamada a `usePermisosProyecto` (después de la carga de Sprints, como hoy); renderizar `AvisoProyectoCerrado`; «Crear Sprint» usa `puedeGestionar` y muestra `NotaCerrado` o la nota actual (nunca ambas); pasar `puedeGestionar` y `puedeEscribir` por `SeccionSprints.tsx` → `FormularioInicio.tsx`, `BotonCerrar.tsx`, `SprintBacklog.tsx` → `EstadoHistoria.tsx` y `EsfuerzoHistoria.tsx` → `FormularioHoras.tsx`, todas con default `true`. Cada componente deshabilita su acción y asocia la nota con `aria-describedby`. Los estilos nuevos, en `web/src/features/sprints/sprints.css`. Tests en verde, incluidos los existentes.
- [ ] T040 [US2] COMMIT `GREEN: Sprints en solo lectura con el proyecto cerrado, test en verde`
- [ ] T041 [US2] Fase 2: `cd web && npm test && npm run lint && npm run build` en verde y pasos «Fase 2» de [quickstart.md](quickstart.md). Marcar T037 a T041 en `tasks.md` y commitear junto con ese cambio: `Validar la Fase 2 de cerrar y reabrir (Refs #312)`.

### Fase 3 — Historias (crear, editar, asignar a Sprint)

- [ ] T042 [US2] RED: casos nuevos en `web/src/features/backlog/crear.test.tsx`: proyecto cerrado → aviso, formulario y «Crear historia» deshabilitados con `NotaCerrado`; abierto → igual que hoy; una sola llamada a `GET /projects/{id}`. Casos nuevos en `web/src/features/backlog/Listado.test.tsx` y `web/src/features/backlog/AsignarASprint.test.tsx` (no existe `Edicion.test.tsx`; `edicionHistoria.test.ts` es de lógica pura y no se toca; `Edicion` se prueba desde `Listado.test.tsx`): cerrado → aviso, el listado se sigue viendo, «Editar» y «Asignar a Sprint» deshabilitados con nota (el enlace «Votar» es navegación y no se deshabilita); `AsignarASprint` con `puedeEscribir={false}` no envía; un 409 «El proyecto está cerrado» al guardar o asignar se muestra tal cual; con el cierre desconocido (cargando o error) todo queda habilitado. Confirmar que FALLA.
- [ ] T043 [US2] COMMIT `RED: Historias en solo lectura con el proyecto cerrado (falla)`
- [ ] T044 [US2] GREEN: en `web/src/features/backlog/crear.tsx` y `web/src/features/backlog/Listado.tsx` (las dos pantallas raíz de `backlog/index.tsx`) llamar UNA vez a `usePermisosProyecto`, mostrar el aviso y la nota y pasar `puedeEscribir` a los botones «Editar» y «Asignar a Sprint» de `Listado.tsx` y a `Edicion.tsx` y `AsignarASprint.tsx` (default `true`). Estilos en `web/src/features/backlog/backlog.css`. Tests en verde.
- [ ] T045 [US2] COMMIT `GREEN: Historias en solo lectura con el proyecto cerrado, test en verde`

### Fase 3 — Votación

- [ ] T046 [US2] RED: casos nuevos en `web/src/features/votacion/index.test.tsx`: proyecto cerrado → aviso, la votación se sigue leyendo, y quedan deshabilitados iniciar votación, votar (cartas), revelar, registrar estimación y nueva ronda, con `NotaCerrado` (sin duplicar la nota de HU-21); abierto → igual que hoy; una sola llamada a `GET /projects/{id}`; un 409 «El proyecto está cerrado» al votar se muestra tal cual; con el cierre desconocido (cargando o error) todo queda habilitado. En `Cartas.test.tsx`: `deshabilitada` ya cubre el caso. Confirmar que FALLA.
- [ ] T047 [US2] COMMIT `RED: Votación en solo lectura con el proyecto cerrado (falla)`
- [ ] T048 [US2] GREEN: en `web/src/features/votacion/index.tsx` reemplazar `useEsScrumMaster` por UNA llamada a `usePermisosProyecto`; aviso arriba; `puedeGestionar` (que ya incluye el cierre) sigue a `SinVotacion`, `Consenso` y `AccionesScrumMaster`; `puedeEscribir` deshabilita `Cartas` vía `Ronda.tsx`; cuando `cerrado`, esos componentes muestran `NotaCerrado` en lugar de su nota de gestión. Estilos en `web/src/features/votacion/votacion.css`. Tests en verde.
- [ ] T049 [US2] COMMIT `GREEN: Votación en solo lectura con el proyecto cerrado, test en verde`

### Fase 3 — Defectos

- [ ] T050 [US2] RED: casos nuevos en `web/src/features/defectos/index.test.tsx`: proyecto cerrado → aviso, el listado de defectos se sigue viendo, y «Registrar defecto» y el botón «Resolver» de cada defecto Pendiente (`ListadoDefectos.tsx`, con Sprint activo) deshabilitados con `NotaCerrado`; abierto → igual que hoy; una sola llamada a `GET /projects/{id}`; un 409 «El proyecto está cerrado» al registrar o resolver se muestra tal cual; con el cierre desconocido todo queda habilitado. Confirmar que FALLA.
- [ ] T051 [US2] COMMIT `RED: Defectos en solo lectura con el proyecto cerrado (falla)`
- [ ] T052 [US2] GREEN: en `web/src/features/defectos/PantallaDefectos.tsx` llamar UNA vez a `usePermisosProyecto`, mostrar el aviso y pasar `puedeEscribir` a `FormularioDefecto.tsx` (deshabilita campos y botón) y a `ListadoDefectos.tsx` (deshabilita «Resolver»), ambas con default `true`. Estilos en `web/src/features/defectos/defectos.css`. Tests en verde.
- [ ] T053 [US2] COMMIT `GREEN: Defectos en solo lectura con el proyecto cerrado, test en verde`

### Fase 3 — Equipo

- [ ] T054 [US2] RED: casos nuevos en `web/src/features/proyecto/Configuracion.test.tsx`: proyecto cerrado → aviso, la lista de integrantes se sigue viendo, y «Asignar» deshabilitado con `NotaCerrado` (una sola nota); abierto → igual que hoy (incluida la nota de HU-21 para el Product Builder); un 409 «El proyecto está cerrado» al asignar se muestra tal cual; con el cierre desconocido todo queda habilitado. Confirmar que FALLA.
- [ ] T055 [US2] COMMIT `RED: Equipo en solo lectura con el proyecto cerrado (falla)`
- [ ] T056 [US2] GREEN: en `web/src/features/proyecto/Configuracion.tsx` reemplazar `useEsScrumMaster` por UNA llamada a `usePermisosProyecto`; aviso; la nota es `NotaCerrado` si `cerrado`, si no la actual. Tests en verde.
- [ ] T057 [US2] COMMIT `GREEN: Equipo en solo lectura con el proyecto cerrado, test en verde`
- [ ] T058 [US2] Fase 3: `cd web && npm test && npm run lint && npm run build` en verde y pasos «Fase 3» de [quickstart.md](quickstart.md). Marcar T042 a T058 en `tasks.md` y commitear junto con ese cambio: `Validar la Fase 3 de cerrar y reabrir (Refs #312)`.

**Checkpoint**: US2 completa en las tres fases.

---

## Phase 7: Polish & Cross-Cutting

- [ ] T059 [P] Cobertura SC-007 y mapa de escrituras: revisar que cada escenario de US1 a US4 y cada borde de la spec (conexión caída, doble clic, Escape, cargando, nota única, cambio de proyecto) tenga una prueba con su nombre en el título; agregar las que falten con su RED/GREEN. Verificar además que cada `post`/`put` de `web/src/features` (los de los `api.ts` y los que se llaman directo desde un componente, como `cerrarSprint` en `sprints/BotonCerrar.tsx`; salvo autenticación y `POST /projects`) esté en el mapa de acciones de escritura con su prueba: `grep -rnE "\b(post|put|patch|del)(<[^>]*>)?\(" web/src --include=*.ts --include=*.tsx`.
- [ ] T060 [P] Confirmar con `git diff main -- web/src/App.tsx web/src/App.css web/src/rutas.tsx web/src/compartido/router/LayoutProyecto.tsx web/src/compartido/permisos/permisos.test.tsx` que no hay cambios, y `git diff main --stat -- . ':!web' ':!specs'` sin cambios de backend (FR-012, SC-006).
- [ ] T061 Recorrido final de [quickstart.md](quickstart.md) contra el backend real, con un Scrum Master y un Product Builder, y marcar las tareas hechas.
- [ ] T062 COMMIT `Marcar las tareas de #312 como hechas (Refs #312)` (el commit incluye el tick de `tasks.md`, así no queda vacío)

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **US1 (Phase 3)** → **US3 (Phase 4)** → **US4 (Phase 5)** → **Aviso en el Panel (T063 a T066)** → **Validación F1 (T036)** → **US2 Fase 2 (T037 a T041)** → **US2 Fase 3 (T042 a T058)** → **Polish**.
- Fase 1 (Phases 2 a 5) es independiente de las Fases 2 y 3, que dependen solo de Phase 2 (hook, aviso y nota). Se puede cortar después de T036 o de T041.
- Dentro de las Fases 2 y 3, las pantallas son independientes entre sí; se pueden hacer en cualquier orden (o en paralelo por personas distintas, en archivos distintos).
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- US3 y US4 extienden `ZonaCierreProyecto.tsx` de US1: van después, en secuencia.

## Parallel Opportunities

- Dentro de Foundational, T003 a T010 (hooks) y T011 a T014 (aviso y nota) tocan archivos distintos y pueden ir en paralelo; T063 a T066 (Panel) necesitan T012.
- Tras T014: Phase 3 a 5 (Fase 1) y T037 (Fase 2) pueden avanzar en paralelo si hay más de una persona, porque tocan archivos distintos.
- Las pantallas de la Fase 3 (Historias, Votación, Defectos, Equipo): cuatro flujos independientes.
- T059 y T060.

## Implementation Strategy

1. **MVP**: Phase 1, Phase 2 y US1 (T001 a T027): el Scrum Master ya puede cerrar desde Editar proyecto y ve el aviso.
2. Completar Fase 1 con US3 y US4 y validar (T036). Entregable por sí solo.
3. Fase 2 (Sprints) y Fase 3 (Historias, Votación, Defectos, Equipo) de a una, cada una validada. Lo que no entre en este PR sigue mostrando el 409 del servidor.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- Los tests existentes de HU-21 y de cada pantalla no se modifican; solo se agregan casos.
