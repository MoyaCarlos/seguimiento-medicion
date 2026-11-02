# Tasks: Imputación de horas en pantalla (HU-07 Frontend)

**Input**: Design documents from `specs/029-hu-07-frontend/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/api.md, quickstart.md

**Feature Branch**: `feature/HU-07-frontend`

**Issue**: #63 (parte de #7)

**Tests**: OBLIGATORIOS (Principio I). Cada comportamiento se prueba primero, con el escenario en el nombre de la prueba (`US1.2: ...`). Un commit por fase; el test nunca va en el mismo commit que su implementación.

**Comandos**: todos desde `web/`. `npx vitest run src/features/sprints` para el ciclo; `npm run lint`, `npm test` y `npm run build` para la verificación final.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: se puede hacer en paralelo (archivos distintos, sin dependencias)
- **[Story]**: historia de usuario a la que pertenece

## Phase 1: Setup

- [X] T001 COMMIT de `specs/029-hu-07-frontend/`: `docs: agregar la especificación del frontend de HU-07 (Issue #63)`

---

## Phase 2: Foundational

**Goal**: tipos, funciones puras y llamadas a la API que usan las tres historias.

- [X] T002 RED: en `web/src/features/sprints/horas.test.ts`, probar `normalizarHoras` (`1,5` → `1.5`; `1.5` y `abc` sin cambios; recorta espacios; vacío sigue vacío), `formatearHoras` (`5.5` → `5,5 h`; `4` → `4 h`; `0` → `0 h`; `61/60` → `1,02 h`) y `fechaPropuesta` (hoy si no pasó el fin; la fecha de fin si ya pasó; usa la fecha **local**, incluso a las 23:30 en UTC-3; hoy si no hay fecha de fin) — confirmar que FALLA
- [X] T003 COMMIT `RED: agregar tests de las funciones de horas del frontend (falla)`
- [X] T004 GREEN: implementar las tres funciones y los tipos `RegistroHoras`, `EsfuerzoReal`, `IntegranteDeProyecto` y `DatosImputacion` (ver data-model.md) en `web/src/features/sprints/horas.ts` — tests en verde
- [X] T005 COMMIT `GREEN: implementar las funciones de horas del frontend, test en verde`
- [X] T006 RED: en `web/src/features/sprints/api.test.ts`, probar `consultarEsfuerzo` (`GET /backlog/{id}/horas`), `imputarHoras` (`POST /backlog/{id}/horas` con el cuerpo de `DatosImputacion`) y `listarIntegrantes` (`GET /projects/{id}/members`) — confirmar que FALLA
- [X] T007 COMMIT `RED: agregar tests de las llamadas de horas a la API (falla)`
- [X] T008 GREEN: agregar las tres llamadas a `web/src/features/sprints/api.ts` sin importar de otra feature (research R1) — tests en verde
- [X] T009 COMMIT `GREEN: implementar las llamadas de horas a la API, test en verde`

**Checkpoint**: la base está lista para las historias.

---

## Phase 3: User Story 1 - Registrar horas en una historia del Sprint activo (Priority: P1) 🎯 MVP

**Goal**: abrir el formulario desde una historia del Sprint activo y guardar una carga.

**Independent Test**: con `fetch` simulado, abrir el formulario, completar los cuatro datos, guardar y verificar el cuerpo enviado y la confirmación. Manual: pasos 2 y 8 de quickstart.md.

- [X] T010 [US1] RED: en `web/src/features/sprints/FormularioHoras.test.tsx`, escribir US1.1 (consulta integrantes al abrir, indicador mientras espera, lista sin nadie seleccionado, campos de fecha, actividad y horas, y fecha propuesta), US1.2 (envía los datos elegidos, muestra la confirmación, avisa con `alRegistrar`, vacía actividad y horas, conserva integrante y fecha), US1.3 (`1,5` se envía como `1.5`), US1.5 (sin doble envío: botón deshabilitado, indicador de envío y una sola llamada) y que dos formularios de historias distintas no repiten `id` de campos — confirmar que FALLA
- [X] T011 [US1] COMMIT `RED: agregar pruebas del formulario de horas (falla)`
- [X] T012 [US1] GREEN: implementar `web/src/features/sprints/FormularioHoras.tsx` con props `historiaId`, `proyectoId`, `fechaFin` y `alRegistrar`; integrante con `<select>` accesible (mismo patrón que `AsignarASprint.tsx`, opción vacía "Elegí un integrante"); fecha con `Field type="date"`; actividad con `Field`; horas con `Field type="text" inputMode="decimal"` (research R3); `id` de campos con el formato `historia-{id}-{campo}`, como `FormularioInicio.tsx`; `LoadingIndicator` mientras envía; confirmación con `role="status"` — tests en verde
- [X] T013 [US1] COMMIT `GREEN: implementar el formulario de horas, test en verde`
- [X] T014 [US1] RED: en `web/src/features/sprints/EsfuerzoHistoria.test.tsx`, probar que con `imputacion` hay un botón "Registrar horas" que abre y cierra el formulario, que al reabrir la lista vuelve sin nadie seleccionado, y en `web/src/features/sprints/SprintBacklog.test.tsx` escribir US1.4 (sin `imputacion`, no se ofrece registrar) — confirmar que FALLA
- [X] T015 [US1] COMMIT `RED: agregar pruebas de abrir el registro de horas desde la historia (falla)`
- [X] T016 [US1] GREEN: implementar `web/src/features/sprints/EsfuerzoHistoria.tsx` (botón con `aria-expanded` que monta y desmonta `FormularioHoras`, research R7), renderizarlo por historia en `web/src/features/sprints/SprintBacklog.tsx` con la nueva prop opcional `imputacion: { proyectoId, fechaFin }`, y pasarla desde `web/src/features/sprints/SeccionSprints.tsx` solo si el Sprint está "Activo" — tests en verde
- [X] T017 [US1] COMMIT `GREEN: abrir el registro de horas desde la historia, test en verde`

**Checkpoint**: se pueden registrar horas desde la pantalla.

---

## Phase 4: User Story 2 - Ver qué dato corregir cuando la carga se rechaza (Priority: P1)

**Goal**: cada rechazo se ve junto a su campo, conservando lo escrito.

**Independent Test**: con `fetch` simulado, responder cada error y verificar dónde aparece. Manual: pasos 3 a 7 y 11 de quickstart.md.

- [X] T018 [US2] RED: en `web/src/features/sprints/FormularioHoras.test.tsx`, escribir US2.1 (sin integrante: mensaje "Elegí un integrante" y ninguna llamada `POST`), US2.2, US2.3, US2.4 y US2.5 (errores con `campo` `horas`, `actividad`, `fecha` e `integrante_id` junto a su campo, con `aria-invalid`), US2.6 (`campo: historia_id` como aviso general), US2.7 (lo escrito se conserva y `alRegistrar` no se llama), US2.8 (el error desaparece al aceptarse), más el aviso de conexión y el fallo al consultar integrantes (muestra el error y no permite guardar) — confirmar que FALLA
- [X] T019 [US2] COMMIT `RED: agregar pruebas de los errores del formulario de horas (falla)`
- [X] T020 [US2] GREEN: implementar en `web/src/features/sprints/FormularioHoras.tsx` la validación del integrante y el mapeo con `clasificarError` sobre `['integrante_id', 'fecha', 'actividad', 'horas']`, sin repetir reglas del backend — tests en verde
- [X] T021 [US2] COMMIT `GREEN: mostrar los errores del formulario de horas, test en verde`

**Checkpoint**: US1 y US2 funcionan juntas.

---

## Phase 5: User Story 3 - Ver el esfuerzo acumulado de cada historia (Priority: P2)

**Goal**: total y detalle plegado en cada historia de la pantalla de Sprints.

**Independent Test**: con `fetch` simulado, verificar totales con y sin cargas, el detalle y la actualización al guardar. Manual: pasos 1, 8, 9, 10 y 12 de quickstart.md.

- [X] T022 [US3] RED: en `web/src/features/sprints/EsfuerzoHistoria.test.tsx`, escribir US3.1 (`5,5 h`), US3.2 (`0 h`, sin error), US3.3 (se reconsulta y actualiza al guardar, sin recargar), US3.4 (aviso en lugar del total), US3.5 (detalle plegado que se despliega con integrante, fecha, actividad y horas, en el orden de la API), US3.6 (sin cargas, sin "Ver detalle") y US3.7 (la carga nueva aparece en el detalle desplegado) — confirmar que FALLA
- [X] T023 [US3] RED: ajustar `web/src/features/sprints/SprintBacklog.test.tsx` para simular `fetch` según la URL (research R8) y agregar: que el fallo del esfuerzo de una historia no oculta las demás (FR-010). En `web/src/features/sprints/index.test.tsx`, agregar FR-001 (el botón "Registrar horas" solo aparece en el Sprint "Activo", no en "Pendiente" ni "Finalizado") y que al cerrar el Sprint desaparece el botón y sigue el total — confirmar que FALLA lo que depende del esfuerzo
- [X] T024 [US3] COMMIT `RED: agregar pruebas del esfuerzo acumulado por historia (falla)`
- [X] T025 [US3] GREEN: implementar en `web/src/features/sprints/EsfuerzoHistoria.tsx` la consulta al montarse, la reconsulta sin parpadeo (research R2), el total, el aviso de error y el detalle plegado; agregar estilos en `web/src/features/sprints/sprints.css` — tests en verde
- [X] T026 [US3] COMMIT `GREEN: mostrar el esfuerzo acumulado de cada historia, test en verde`

**Checkpoint**: las tres historias funcionan.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T027 (solo si hay algo que limpiar) REFACTOR de `web/src/features/sprints/` + COMMIT `REFACTOR: <qué se limpió>`, con los tests en verde
- [X] T028 Correr en `web/` `npm run lint`, `npm test` (incluye `aislamiento.test.ts`) y `npm run build`, sin errores
- [ ] T029 Validar en el navegador la tabla de `specs/029-hu-07-frontend/quickstart.md`, confirmar que el comando de arranque es correcto, cronometrar una carga (SC-001: menos de 1 minuto) y registrar el resultado. Hecho, salvo cronometrar la carga con una persona
- [X] T030 Marcar el avance y actualizar el Status de `specs/029-hu-07-frontend/spec.md`
- [X] T031 COMMIT `docs: registrar la validación del frontend de HU-07 (Issue #63)`
- [X] T032 [US1] Arreglo de la revisión: deshabilitar "Registrar horas" mientras se envía una carga en `web/src/features/sprints/EsfuerzoHistoria.tsx`, para no duplicarla al cerrar y reabrir el formulario (FR-009), con su RED y su GREEN
- [X] T033 Arreglo de la revisión: limitar el ancho de la lista de integrantes y darle estilo de foco en `web/src/features/sprints/sprints.css`
- [X] T034 [US3] Agregar la prueba de la reconsulta fallida tras guardar en `web/src/features/sprints/EsfuerzoHistoria.test.tsx`
- [ ] T035 Registrar como issue `bug` la inestabilidad de las pruebas de `backlog` y de `rutas.test.tsx` con la máquina cargada (también en `main`)

---

## Dependencies & Execution Order

- **Setup** → **Foundational** → **US1** → **US2** → **US3** → **Polish**.
- US2 depende de `FormularioHoras` (US1). US3 extiende `EsfuerzoHistoria` (US1).
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT. Nunca saltear el COMMIT de RED.

### Parallel Opportunities

Las fases comparten archivos y cada commit sigue al anterior, así que se ejecutan en orden, sin tareas [P]. T002 y T006 no dependen entre sí, pero cada una lleva su propio commit.

## Implementation Strategy

1. Foundational + US1 → MVP: se pueden registrar horas.
2. US2 → los errores se ven junto al campo.
3. US3 → el resultado se ve en el total y el detalle.

## Notes

- Ninguna tarea hace `git push` ni merge. El PR lleva `Closes #63`, base `main`, y se mergea con "Create a merge commit" después de 1 aprobación.
- Esta historia no agrega `.feature`: el frontend cumple BDD con pruebas de componente (Principio III) y la API ya tiene los escenarios de Godog de la spec 019.
