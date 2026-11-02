---

description: "Tareas de HU-02 Frontend: pantalla de Planning Poker (TDD con commits por fase)"
---

# Tasks: Pantalla de Planning Poker (HU-02, Frontend)

**Input**: Design documents from `/specs/029-hu-02-frontend-planning-poker/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/api-votacion.md, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I y III v1.1.0). Cada comportamiento es RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT). Las pruebas son de Vitest + Testing Library y llevan el escenario Given-When-Then (o el id `US1.1`) en el nombre.

**Convenciones**: todos los comandos de `npm` se corren dentro de `web/` con Node 24 (`nvm use`). Todos los commits terminan con `Refs #48`. En cada RED, junto al test se agregan solo los exports que el test necesita con cuerpo mínimo, de modo que `npm run build` compile y el test falle por aserción. Se agregan archivos por nombre, nunca `git add -A`. Sin dependencias nuevas, sin tocar el backend.

**Organization**: por user story de la spec. Se implementan en este orden: US1 (votar), US3 (revelar), US2 (iniciar), US4 (acceso y actualización); las dos primeras son P1.

## Format: `[ID] [P?] [Story] Description`

---

## Phase 1: Setup

- [ ] T001 Verificar en la rama `feature/HU-02-frontend-planning-poker` que `npm ci`, `npm run lint`, `npm test` y `npm run build` están en verde antes de empezar. Sin cambios ni commit

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: tipos, reglas puras, cliente de la API y mensajes de error que usan todas las historias.

- [ ] T002 RED: crear web/src/features/votacion/votacion.test.ts con tests de las reglas puras: `CARTAS` es exactamente `[1, 2, 3, 5, 8, 13, 21]`; `esScrumMaster` solo para `scrum_master`; `estaRevelada`; `faltanPorVotar` lista a los de `ya_voto: false`; `puedeRevelar` es `true` solo si el actor es Scrum Master, la ronda no está revelada y no falta nadie; `textoDeVoto`: «Falta» si no votó, valor propio si el integrante es el actor, «Oculto» si votó y es ajeno sin revelar, valor si está revelada — confirmar que FALLA
- [ ] T003 COMMIT `RED: agregar tests de las reglas de la votación (falla)`
- [ ] T004 GREEN: implementar web/src/features/votacion/votacion.ts (tipos `Rol`, `Carta`, `Integrante`, `VotoDeIntegrante`, `Votacion` de [data-model.md](./data-model.md) y las funciones puras) — `npm test` en verde
- [ ] T005 COMMIT `GREEN: implementar las reglas de la votación, test en verde`
- [ ] T006 RED: crear web/src/features/votacion/api.test.ts con tests del cliente (con `fetch` simulado): `consultarVotacion(historiaId, usuarioId?)` hace `GET /backlog/{id}/votacion` con `?usuario_id=N` solo si hay usuario; `iniciarVotacion` → `POST /backlog/{id}/votacion` con `{usuario_id}`; `votar` → `POST /backlog/{id}/votacion/votos` con `{usuario_id, valor}`; `revelarVotacion` → `POST /backlog/{id}/votacion/revelar` con `{usuario_id}`; `listarIntegrantes(proyectoId)` → `GET /projects/{id}/members` y devuelve `{id, nombre, rol}` — confirmar que FALLA
- [ ] T007 COMMIT `RED: agregar tests del cliente de la API de votación (falla)`
- [ ] T008 GREEN: implementar web/src/features/votacion/api.ts con `get` y `post` de `../../compartido/api/client` — `npm test` en verde
- [ ] T009 COMMIT `GREEN: implementar el cliente de la API de votación, test en verde`
- [ ] T010 RED: crear web/src/features/votacion/errores.test.ts: `mensajeDeError` devuelve el mensaje de `ApiError` y de `ConnectionError`, y «Ocurrió un error inesperado» para cualquier otro valor; `esVotacionNoIniciada` es `true` solo para `ApiError` 404 con el mensaje exacto «la historia no tiene votación iniciada» (R1 de [research.md](./research.md)) — confirmar que FALLA
- [ ] T011 COMMIT `RED: agregar tests de los mensajes de error de la votación (falla)`
- [ ] T012 GREEN: implementar web/src/features/votacion/errores.ts — `npm test` en verde
- [ ] T013 COMMIT `GREEN: implementar los mensajes de error de la votación, test en verde`

**Checkpoint**: tipos, reglas, cliente y errores listos.

---

## Phase 3: User Story 1 - Votar en secreto con cartas (Priority: P1) 🎯 MVP

**Goal**: elegir quién actúa, ver quién votó y quién falta con los votos ajenos «Oculto», y votar o cambiar el voto con las cartas.

**Independent Test**: con una votación abierta simulada, elegir un integrante y la carta 8: aparece «ya votó» con su 8 y «Oculto» en los demás.

- [ ] T014 [US1] RED: crear web/src/features/votacion/Cartas.test.tsx: `Given una ronda abierta, When se muestra la baraja, Then hay exactamente un botón por carta 1, 2, 3, 5, 8, 13 y 21`; `When se pulsa una carta Then se llama a alElegir con su valor`; `Given deshabilitada Then ninguna carta responde` y la carta seleccionada queda marcada (`aria-pressed`) — confirmar que FALLA
- [ ] T015 [US1] COMMIT `RED: agregar tests de la baraja de cartas (falla)`
- [ ] T016 [US1] GREEN: implementar web/src/features/votacion/Cartas.tsx (props `seleccionada`, `deshabilitada`, `alElegir`) — `npm test` en verde
- [ ] T017 [US1] COMMIT `GREEN: implementar la baraja de cartas, test en verde`
- [ ] T018 [US1] RED: crear web/src/features/votacion/Integrantes.test.tsx: `Given una ronda sin revelar, When se muestran los integrantes, Then quien votó figura «Ya votó» con «Oculto» y quien no figura «Falta»`; `Given el integrante elegido votó 8, Then se ve su 8`; `Given una ronda revelada, Then se ve el valor de cada uno` — confirmar que FALLA
- [ ] T019 [US1] COMMIT `RED: agregar tests de la lista de integrantes (falla)`
- [ ] T020 [US1] GREEN: implementar web/src/features/votacion/Integrantes.tsx usando `textoDeVoto` — `npm test` en verde
- [ ] T021 [US1] COMMIT `GREEN: implementar la lista de integrantes, test en verde`
- [ ] T022 [US1] RED: crear web/src/features/votacion/index.test.tsx (pantalla, con `./api` simulado y un router en memoria en `/projects/25/votacion/3`): `US1.0 Given integrantes del proyecto y votación abierta, When se abre la pantalla, Then se pide elegir quién actúa y no hay cartas`; `US1.1 When elijo un integrante y la carta 8, Then se llama a votar con (3, usuario, 8) y se muestra su voto`; `US1.2 Given ya voté, When elijo otra carta, Then se muestra la nueva`; `US1.3 When cambio de integrante, Then se consulta con el nuevo usuario_id y se muestra su voto`; `US1.4 Given un voto en curso, When pulso otra carta, Then votar se llamó una sola vez y las cartas están deshabilitadas`; `US1.5 Given el servidor rechaza el voto, Then se muestra su mensaje, se vuelve a consultar y la pantalla sigue operable`; proyecto sin integrantes muestra el aviso; `Given el servidor no responde, Then se muestra el aviso de conexión`; `Given la respuesta de un integrante anterior llega tarde, Then se descarta` — confirmar que FALLA
- [ ] T023 [US1] COMMIT `RED: agregar tests de la pantalla de votación, voto secreto (falla)`
- [ ] T024 [US1] GREEN: implementar web/src/features/votacion/index.tsx (`Component`: lee `id` y `historiaId`, carga integrantes y votación con marca de vigencia, selector, `Cartas`, `Integrantes`, estado `enviando` único) y web/src/features/votacion/votacion.css — `npm test` en verde
- [ ] T025 [US1] COMMIT `GREEN: implementar la pantalla de votación con voto secreto, test en verde`
- [ ] T026 [US1] (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`
- [ ] T027 [US1] RED: agregar en web/src/rutas.test.tsx una prueba que resuelve `/projects/25/votacion/3` hacia la pantalla de votación (no hacia `NoEncontrada`) — confirmar que FALLA
- [ ] T028 [US1] COMMIT `RED: agregar test de la ruta de la votación (falla)`
- [ ] T029 [US1] GREEN: registrar la ruta `votacion/:historiaId` con `lazy` en web/src/rutas.tsx — `npm test` en verde
- [ ] T030 [US1] COMMIT `GREEN: registrar la ruta de la votación, test en verde`

**Checkpoint**: US1 funciona y se prueba sola.

---

## Phase 4: User Story 3 - Revelar y ver el resultado (Priority: P1)

**Goal**: el Scrum Master revela cuando votaron todos y se ven valores, promedio exacto y carta sugerida.

**Independent Test**: con todos los votos cargados, «Revelar» muestra 3, 5, 8, promedio 5,33 y carta 5.

- [ ] T031 [US3] RED: crear web/src/features/votacion/Resultado.test.tsx: `Given una ronda revelada con promedio 5.33 y carta 5, When se muestra el resultado, Then se ven "5,33" y la carta sugerida 5`; sin promedio no muestra nada — confirmar que FALLA
- [ ] T032 [US3] COMMIT `RED: agregar tests del resultado de la votación (falla)`
- [ ] T033 [US3] GREEN: implementar web/src/features/votacion/Resultado.tsx — `npm test` en verde
- [ ] T034 [US3] COMMIT `GREEN: implementar el resultado de la votación, test en verde`
- [ ] T035 [US3] RED: agregar en index.test.tsx: `US3.1 Given votaron todos y soy Scrum Master, When pulso Revelar, Then se llama a revelar y se ven valores, promedio y carta sugerida sin cartas`; `US3.2 Given falta alguien, Then Revelar está deshabilitado y se indica quién falta`; `US3.3 Given soy Product Builder, Then no hay botón Revelar`; `US3.4 Given revelar en curso, When pulso otra vez, Then se llamó una sola vez`; `US3.5 Given el servidor rechaza revelar, Then se muestra su mensaje`; `Given ronda revelada Then se aclara que registrar la estimación acordada es otra historia` — confirmar que FALLA
- [ ] T036 [US3] COMMIT `RED: agregar tests de revelar los votos en la pantalla (falla)`
- [ ] T037 [US3] GREEN: agregar «Revelar» y `Resultado` en web/src/features/votacion/index.tsx — `npm test` en verde
- [ ] T038 [US3] COMMIT `GREEN: implementar la revelación de votos en la pantalla, test en verde`
- [ ] T039 [US3] (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: US1 y US3 funcionan juntas.

---

## Phase 5: User Story 2 - Iniciar la votación (Priority: P2)

**Goal**: el Scrum Master abre la ronda 1 de una historia sin votación.

**Independent Test**: historia sin votación + Scrum Master elegido: «Iniciar votación» muestra la ronda 1 con todos «Falta».

- [ ] T040 [US2] RED: agregar en index.test.tsx: `US2.0 Given la API responde 404 "la historia no tiene votación iniciada", Then no hay cartas y se ofrece iniciar`; `US2.1 Given Scrum Master elegido, When pulso Iniciar votación, Then se llama a iniciar y aparece la ronda 1 con todos "Falta"`; `US2.2 Given Product Builder elegido, Then no hay botón Iniciar y se indica que solo el Scrum Master puede iniciar`; `US2.4 Given inicio en curso, When pulso otra vez, Then se llamó una sola vez`; `US2.5 Given el servidor rechaza (409 o historia Completada), Then se muestra su mensaje`; `Given 404 "historia no encontrada" Then se muestra el mensaje y no se ofrece iniciar` — confirmar que FALLA
- [ ] T041 [US2] COMMIT `RED: agregar tests de iniciar la votación en la pantalla (falla)`
- [ ] T042 [US2] GREEN: agregar «Iniciar votación» y el estado «sin votación» en web/src/features/votacion/index.tsx usando `esVotacionNoIniciada` — `npm test` en verde
- [ ] T043 [US2] COMMIT `GREEN: implementar el inicio de la votación en la pantalla, test en verde`
- [ ] T044 [US2] (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

---

## Phase 6: User Story 4 - Llegar a la pantalla y mantenerla al día (Priority: P2)

**Goal**: botón «Votar» en el Product Backlog y botón «Actualizar» en la pantalla.

**Independent Test**: desde el backlog pulsar «Votar» y llegar a la pantalla; «Actualizar» refleja un voto nuevo.

- [ ] T045 [US4] RED: agregar en web/src/features/votacion/index.test.tsx: `US4.2 When pulso Actualizar, Then se vuelve a consultar y se reflejan los votos nuevos`; `US4.4 Given falló la consulta por falta de conexión, When pulso Actualizar y el servidor responde, Then se muestra la votación` — confirmar que FALLA
- [ ] T046 [US4] COMMIT `RED: agregar tests de actualizar la votación (falla)`
- [ ] T047 [US4] GREEN: agregar «Actualizar» en web/src/features/votacion/index.tsx — `npm test` en verde
- [ ] T048 [US4] COMMIT `GREEN: implementar Actualizar en la pantalla de votación, test en verde`
- [ ] T049 [US4] RED: agregar en web/src/features/backlog/Listado.test.tsx: `US4.1 Given historias en el backlog del proyecto 42, When se muestra, Then cada una tiene un enlace "Votar" a /projects/42/votacion/{id}` — confirmar que FALLA
- [ ] T050 [US4] COMMIT `RED: agregar test del enlace Votar en el backlog (falla)`
- [ ] T051 [US4] GREEN: agregar el `Link` «Votar» en web/src/features/backlog/Listado.tsx, sin importar la feature `votacion` — `npm test` (incluye `aislamiento.test.ts`) en verde
- [ ] T052 [US4] COMMIT `GREEN: agregar el enlace Votar al Product Backlog, test en verde`

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T053 `npm run lint`, `npm test` y `npm run build` sin errores (Node 24, dentro de `web/`)
- [ ] T054 Validación manual de quickstart.md con el backend desde este worktree (base temporal y `/demo/datos`); apagar backend y frontend al terminar
- [ ] T055 Actualizar el estado de la spec y marcar las tareas hechas en specs/029-hu-02-frontend-planning-poker/ y COMMIT `Marcar las tareas de la spec 029 como hechas (Refs #48)`

---

## Dependencies & Execution Order

- Setup → Foundational → US1 → US3 → US2 → US4 → Polish.
- US3, US2 y US4 extienden `index.tsx` de US1; no son paralelas.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- Oportunidad en paralelo: T014 a T021 (Cartas e Integrantes) son archivos distintos, pero se commitean uno a uno.

## Implementation Strategy

1. Setup + Foundational.
2. US1 (MVP): votar con voto secreto y ruta accesible por URL.
3. US3, US2 y US4 de a una, dejando todo en verde en cada commit.

## Notes

- El historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas; el PR lo abre el equipo con `Refs #48` o `Closes #48`.
- No se commitean `docs/reviews/`, `INFORME-*.md`, `Proyecto - 2026.pdf`, `docs/prototipo-sprint2/` ni `specs/007-frontend-base/analisis.md`.
