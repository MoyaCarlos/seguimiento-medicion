---

description: "Task list for HU-01 frontend: pantalla de creación de historias de usuario"
---

# Tasks: Pantalla de creación de historias de usuario (HU-01, frontend)

**Input**: Design documents from `/specs/008-hu-01-frontend/`

**Prerequisites**: [plan.md](plan.md) (required), [spec.md](spec.md) (required), [research.md](research.md), [data-model.md](data-model.md), [contracts/api-backlog.md](contracts/api-backlog.md), [quickstart.md](quickstart.md)

**Tests**: OBLIGATORIOS. La constitución del proyecto (Principio I, NON-NEGOTIABLE) exige TDD estricto con un commit por fase. La spec además exige pruebas de componente con Testing Library + jsdom, con el escenario de aceptación en el nombre de cada prueba.

**Organización**: tareas por historia de usuario, en orden de prioridad. Todas las historias comparten la misma pantalla (`index.tsx`) y el mismo módulo de reglas (`historia.ts`), por lo que se implementan **en secuencia**: US1 → US2 → US3 → US4.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: puede ejecutarse en paralelo (archivos distintos, sin dependencias pendientes)
- **[Story]**: US1, US2, US3, US4
- Rutas exactas de archivo en cada tarea

## Phase 1: Setup (Shared Infrastructure)

- [X] T001 [P] Crear `web/src/features/backlog/opciones.ts` con `PRIORIDADES` (`M`/Must have, `S`/Should have, `C`/Could have, `W`/Won't have), `VALORES_NEGOCIO` (`1, 2, 3, 5, 8, 13, 21`), los tipos `Prioridad` y `ValorNegocio`, y `etiquetaPrioridad(codigo)`. Es data/mapeo trivial: no lleva ciclo RED
- [X] T002 [P] Crear `web/src/features/backlog/backlog.css` con reglas acotadas a `.backlog-form textarea` y `.backlog-form select`, reutilizando las variables `--line`, `--ink` y `--lilac` de `web/src/App.css`. No tocar `App.css` ni `compartido/`
- [X] T003 COMMIT `Agregar opciones y estilos de la pantalla de backlog`

**Checkpoint**: opciones y estilos listos; todavía no hay pantalla.

---

## Phase 2: Foundational (Blocking Prerequisites)

**⚠️ CRITICAL**: sin estos tipos no se puede empezar ninguna historia.

- [X] T004 Crear `web/src/features/backlog/historia.ts` con los tipos `FormularioHistoria`, `CuerpoCrearHistoria`, `HistoriaCreada` y `ErrorDeCampo`, y la constante `ESTADO_INICIAL` (`titulo: ''`, `descripcion: ''`, `prioridad: 'M'`, `valorNegocio: ''`). Solo tipos/constante: sin funciones ni lógica
- [X] T005 COMMIT `Definir tipos del formulario de historia`

**Checkpoint**: Foundation ready; empieza la implementación por historia.

---

## Phase 3: User Story 1 - Crear una historia de usuario desde el Product Backlog (Priority: P1) 🎯 MVP

**Goal**: la pantalla del `backlog` del proyecto de la URL crea una historia con `POST /backlog`, muestra la confirmación (estado "Nueva", sin estimación) y deja el formulario listo para otra carga.

**Independent Test**: abrir `/projects/42/backlog`, completar título, descripción y prioridad, guardar; verificar una única llamada a `/backlog` con `proyecto_id: 42`, la confirmación con "Nueva" y `estimacion_sp: null`, el formulario restablecido y el proyecto conservado al recargar.

- [ ] T006 [US1] RED: en `web/src/features/backlog/historia.test.ts` escribir tests de `parsearProyectoId` (acepta `"42"` → `42`; rechaza `undefined`, `""`, `"0"`, `"-1"`, `"1.5"`, `"abc"`, `" 1"`, `"1e3"`, y los enteros fuera del rango seguro de JavaScript `"9007199254740993"` y `"99999999999999999999"`, que perderían precisión al serializarse) y de `armarCuerpo` (título y descripción con `trim`, prioridad `M/S/C/W`, `proyecto_id` numérico, `valor_negocio` omitido cuando es `''`). Correr `npm test` en `web/` y confirmar que FALLA
- [ ] T007 [US1] COMMIT `RED: agregar tests de parseo del proyecto y armado del cuerpo (falla)` (solo `historia.test.ts`)
- [ ] T008 [US1] GREEN: implementar `parsearProyectoId` y `armarCuerpo` en `web/src/features/backlog/historia.ts`; `npm test` en verde
- [ ] T009 [US1] COMMIT `GREEN: implementar parseo del proyecto y armado del cuerpo, test en verde`
- [ ] T010 [US1] RED: en `web/src/features/backlog/index.test.tsx` montar `Component` con `createMemoryRouter`/`MemoryRouter` en `/projects/42/backlog`, simular `fetch` y escribir un `it` por escenario de aceptación de US1 (los tres, con el escenario en el nombre) más los casos: id inválido (aviso general + formulario deshabilitado, sin llamada); FR-016 con un `fetch` pendiente (un segundo envío no genera otra llamada y el botón "Guardar" queda deshabilitado); dos historias con el mismo título (la segunda se envía y no se trata como error); las cuatro prioridades MoSCoW —escenarios 1 y 2 de US4, adelantados para que la fase RED de US4 no contenga tests ya en verde— (las cuatro etiquetas visibles, "Must have" por defecto, cada una envía su código `M/S/C/W` y la confirmación muestra la etiqueta correspondiente). Verificar una sola llamada, cuerpo con `proyecto_id: 42` y textos recortados, confirmación con la historia, "Nueva" y sin estimación, formulario restablecido, y que una recarga conserva el proyecto de la URL. Confirmar que FALLA
- [ ] T011 [US1] COMMIT `RED: agregar pruebas de componente de creación de historia (falla)` (solo `index.test.tsx`)
- [ ] T012 [US1] GREEN: implementar `web/src/features/backlog/api.ts` (`crearHistoria` sobre `post<HistoriaCreada>('/backlog', cuerpo)` del cliente compartido) y `web/src/features/backlog/index.tsx` (reemplaza el placeholder): formulario con título (`Field`), descripción (`textarea`), prioridad (`select` con las cuatro etiquetas y `M` por defecto), cada control con `label htmlFor`/`id` propios y su error asociado por `aria-describedby`, estado local, envío único, `LoadingIndicator` + `Button` compartido deshabilitado, confirmación con título, prioridad legible, estado "Nueva", "sin estimación" y valor de negocio si existe, formulario restablecido, proyecto leído de `useParams`, y `fieldset disabled` + aviso general si el id no es válido. Importar `backlog.css`. `npm test` en verde
- [ ] T013 [US1] COMMIT `GREEN: implementar pantalla de creación de historia, test en verde`

**Checkpoint**: US1 funcional y verificable por sí sola (MVP).

---

## Phase 4: User Story 2 - Evitar el envío de historias incompletas o demasiado largas (Priority: P1)

**Goal**: validar antes de llamar al servidor, con todas las advertencias juntas, sin registrar nada inválido.

**Independent Test**: intentar guardar sin título, solo espacios, con 201/2001 caracteres y con ambos campos inválidos; comprobar advertencias junto a cada campo y cero llamadas; con exactamente 200/2000, comprobar que se envía.

- [ ] T014 [US2] RED: en `web/src/features/backlog/historia.test.ts` escribir tests de `contarCaracteres` (cuenta puntos de código: `"😀"` → 1, `"ñ"` → 1, `"e\u0301"` → 2; no unidades UTF-16) y de `validarHistoria` con las reglas de `data-model.md` citadas textualmente: título "Obligatorio; se evalúa después de `trim`; máximo 200 puntos de código", descripción "Obligatoria; se evalúa después de `trim`; máximo 2000 puntos de código", límites inclusivos 200/2000, y devolución de **todos** los errores de campo del intento. Confirmar que FALLA
- [ ] T015 [US2] COMMIT `RED: agregar tests de validación del formulario (falla)`
- [ ] T016 [US2] GREEN: implementar `contarCaracteres` y `validarHistoria` en `web/src/features/backlog/historia.ts` (recortar antes de validar y contar con `[...texto].length`); `npm test` en verde
- [ ] T017 [US2] COMMIT `GREEN: implementar validación del formulario, test en verde`
- [ ] T018 [US2] RED: agregar en `web/src/features/backlog/index.test.tsx` un `it` por escenario de aceptación de US2 (cuatro) más el caso de varios errores simultáneos: advertencia junto al campo, cero llamadas a `fetch`, exactamente 200/2000 aceptados, y que un intento inválido reemplaza la confirmación anterior. Confirmar que FALLA
- [ ] T019 [US2] COMMIT `RED: agregar pruebas de componente de validación (falla)`
- [ ] T020 [US2] GREEN: integrar `validarHistoria` en `web/src/features/backlog/index.tsx`: mostrar todas las advertencias junto a cada control con `label` asociada por `htmlFor`/`id` y `aria-invalid`/`aria-describedby`, bloquear el envío si hay errores, limpiar confirmación/errores al iniciar el intento y no borrar lo escrito. `npm test` en verde
- [ ] T021 [US2] COMMIT `GREEN: integrar validación en la pantalla, test en verde`

**Checkpoint**: US1 y US2 funcionando; no se crean historias inválidas.

---

## Phase 5: User Story 3 - Entender por qué el servidor rechazó la historia (Priority: P2)

**Goal**: mostrar el rechazo del servidor junto al campo o en el aviso general, sin perder datos.

**Independent Test**: simular 400 con campo, 400 sin campo, 404 de proyecto, error de conexión/500 y comprobar el lugar del mensaje, que los datos se conservan, que la pantalla sigue usable y que un nuevo intento reemplaza el aviso.

- [ ] T022 [US3] RED: en `web/src/features/backlog/historia.test.ts` escribir tests de `mapearError`: `ApiError` con `campo` conocido (`titulo`, `descripcion`, `prioridad`, `valor_negocio`) → mensaje de ese campo; sin campo, campo desconocido o `proyecto_id` → aviso general; 404 "proyecto no encontrado" → general; `ConnectionError` → mensaje de conexión general. Confirmar que FALLA
- [ ] T023 [US3] COMMIT `RED: agregar tests de mapeo de errores del servidor (falla)`
- [ ] T024 [US3] GREEN: implementar `mapearError` en `web/src/features/backlog/historia.ts`; `npm test` en verde
- [ ] T025 [US3] COMMIT `GREEN: implementar mapeo de errores del servidor, test en verde`
- [ ] T026 [US3] RED: agregar en `web/src/features/backlog/index.test.tsx` un `it` por escenario de aceptación de US3 (cinco): error 400 con campo junto al campo y datos conservados; 400 sin campo en el aviso general; 404 de proyecto en el aviso general; backend caído/500 con aviso general y pantalla usable; y que un nuevo intento reemplaza el aviso anterior. Verificar además que, con un aviso visible, escribir en un campo no lo borra (FR-015). Confirmar que FALLA
- [ ] T027 [US3] COMMIT `RED: agregar pruebas de componente de errores del servidor (falla)`
- [ ] T028 [US3] GREEN: integrar `mapearError` y `ApiError`/`ConnectionError` en `web/src/features/backlog/index.tsx` (mensaje junto al campo o `ErrorMessage` arriba del formulario, sin borrar datos); `npm test` en verde
- [ ] T029 [US3] COMMIT `GREEN: integrar errores del servidor en la pantalla, test en verde`

**Checkpoint**: todos los rechazos se muestran en el lugar correcto.

---

## Phase 6: User Story 4 - Elegir la prioridad y el valor de negocio (Priority: P2)

**Goal**: completar la clasificación de la historia con el valor de negocio opcional de la escala; la prioridad MoSCoW ya se elige desde US1.

**Independent Test**: guardar con cada valor de negocio permitido y comprobar el valor enviado y la confirmación; guardar con "Sin valor" y comprobar que se crea sin `valor_negocio`. Los casos de prioridad se verifican en la prueba independiente de US1.

- [ ] T030 [US4] RED: agregar en `web/src/features/backlog/index.test.tsx` un `it` por escenario de aceptación de US4 que todavía no cubra T010: el valor de negocio ofrece solo "Sin valor", 1, 2, 3, 5, 8, 13 y 21; con "Sin valor" el campo se omite del cuerpo, la confirmación no lo muestra y la historia se crea igual; y el selector vuelve a "Sin valor" tras el éxito. Confirmar que FALLA (falta el selector de valor de negocio). Los casos de prioridad (etiquetas, "Must have" por defecto, código enviado y etiqueta en la confirmación) ya quedaron cubiertos en T010, para que esta fase RED no contenga tests que ya pasan.
- [ ] T031 [US4] COMMIT `RED: agregar pruebas de componente de valor de negocio (falla)`
- [ ] T032 [US4] GREEN: agregar en `web/src/features/backlog/index.tsx` el `select` de valor de negocio (opción inicial "Sin valor" y los valores de `opciones.ts`), enviar su número solo cuando corresponde, mostrarlo en la confirmación y restablecerlo; confirmar (sin agregar tests nuevos, ya cubiertos en T010) que la prioridad usa `PRIORIDADES` y su etiqueta. `npm test` en verde
- [ ] T033 [US4] COMMIT `GREEN: implementar valor de negocio en la pantalla, test en verde`

**Checkpoint**: las cuatro historias quedan funcionales.

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T034 Desde `web/` correr `npm test`, `npm run lint` y `npm run build`; corregir lo que falle sin tocar `web/src/compartido/`, `web/package.json` ni `web/src/rutas.tsx`
- [ ] T035 Ejecutar los 10 puntos de validación de [quickstart.md](quickstart.md) con backend y frontend levantados; confirmar que no hay cambios en `compartido/`, `package.json` ni en las demás pantallas, y que no quedan advertencias de lint ni tests en rojo

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **US1 (Phase 3)** → **US2 (Phase 4)** → **US3 (Phase 5)** → **US4 (Phase 6)** → **Polish (Phase 7)**.
- Las historias comparten `historia.ts`, `index.tsx` e `index.test.tsx`, por lo que **no se ejecutan en paralelo**; cada una deja la suite en verde.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT. Nunca agrupar test e implementación en el mismo commit.
- Los COMMIT de RED incluyen solo los archivos de test; los de GREEN, solo la implementación correspondiente.
- No se hace `git push` ni merge desde estas tareas: lo decide el equipo por PR (merge commit, nunca squash).
- Dependencias puntuales: T004 depende de T001 (tipos de prioridad); T006 depende de T004; T008 depende de T006; T010 depende de T008; T012 depende de T008 y T010; T014 depende de T012; T016 depende de T014; T018 depende de T016 y T012; T020 depende de T016 y T018; T022 depende de T012; T024 depende de T022; T026 depende de T024 y T020; T028 depende de T024 y T026; T030 depende de T028; T032 depende de T024 y T030; T034 depende de T032; T035 depende de T034.

## Parallel Opportunities

- T001 (`opciones.ts`) y T002 (`backlog.css`) pueden hacerse en paralelo: archivos distintos y sin dependencias entre sí.
- Fuera de eso, la implementación es secuencial porque casi todas las tareas tocan los mismos dos archivos (`historia.ts`/`historia.test.ts` o `index.tsx`/`index.test.tsx`).

## Implementation Strategy

1. Completar Setup y Foundational.
2. US1 (MVP): pantalla que crea y confirma; validar con su prueba independiente.
3. US2: validación local completa.
4. US3: errores del servidor.
5. US4: prioridades y valor de negocio.
6. Polish: suite, lint, build y quickstart.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra. Solo se agrega un REFACTOR si realmente hay algo que limpiar, con su propio commit.
- Todos los textos de la pantalla en español y los máximos de 200/2000 contados en puntos de código sobre el texto recortado, igual que el backend.
- No se agregan dependencias ni se modifica `web/src/compartido/`; la ruta `backlog/*` ya existe.
