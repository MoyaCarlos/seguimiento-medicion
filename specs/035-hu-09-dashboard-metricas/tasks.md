# Tasks: Dashboard de Métricas (HU-09)

**Input**: Design documents from `specs/035-hu-09-dashboard-metricas/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/api.md, quickstart.md

**Feature Branch**: `feature/HU-09-dashboard-metricas`

**Issue**: #69 (parte de #9)

**Tests**: OBLIGATORIOS (Principio I). Cada comportamiento se prueba primero, con el escenario en el nombre de la prueba (`US1.2: ...`). Un commit por fase; el test nunca va en el mismo commit que su implementación.

**Comandos**: todos desde `web/`. `npx vitest run src/features/dashboard` para el ciclo; `npm run lint`, `npm test` y `npm run build` para la verificación final.

**Reglas**: sin `git push` ni merge desde las tareas (merge commit, nunca squash); archivos agregados por nombre; no se toca `internal/`, `features/` ni el esquema.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: se puede hacer en paralelo (archivos distintos, sin dependencias)
- **[Story]**: historia de usuario a la que pertenece

## Phase 1: Setup

El equipo confirmó las decisiones en la reunión de definición, incluida la dependencia nueva que exige decidir el Principio VI.

- [X] T001 COMMIT de `specs/035-hu-09-dashboard-metricas/`: `docs: agregar la especificación de HU-09 (Issue #69)`
- [X] T002 Registrar en `specs/035-hu-09-dashboard-metricas/spec.md` las decisiones confirmadas en la reunión de definición (sin backend nuevo, Recharts, alcance del panel) y el punto abierto de horas y defectos por Sprint, que no se trató. COMMIT `docs: registrar las decisiones de la reunión de HU-09`
- [ ] T003 Instalar `recharts` y `react-is` con `npm install` en `web/`; se modifican `web/package.json` y `web/package-lock.json`
- [ ] T004 Verificar Recharts en jsdom con una prueba mínima, temporal, en `web/src/features/dashboard/GraficoVelocidad.test.tsx`: montar un gráfico de barras, resolver lo que haga falta (`ResizeObserver` o dimensiones fijas) y anotar el resultado en R3 de `specs/035-hu-09-dashboard-metricas/research.md`
- [ ] T005 COMMIT de `web/package.json`, `web/package-lock.json` y `research.md`, sin la prueba temporal: `build: agregar Recharts para el Dashboard (HU-09)`

---

## Phase 2: Foundational

**Goal**: tipos, funciones puras y consulta a la API que usan todas las historias.

- [ ] T006 RED: en `web/src/features/dashboard/metricas.test.ts`, probar `datosDelGrafico` (una fila por Sprint, conserva el orden, valores en cero incluidos; `[]` → `[]`), `formatoNumero` (`-6.25` → `-6,25`; `4` → `4`; `0` → `0`; `null` → `sin datos`), `textoDesviacion` (positiva → más de lo estimado; negativa → menos; `0` → sin desviación; `null` → `null`) y `proporcion` (`5, 10` → `50`; máximo `0` → `0`; valor `null` → `0`; nunca más de `100`) — confirmar que FALLA
- [ ] T007 COMMIT `RED: agregar tests de las funciones del Dashboard (falla)`
- [ ] T008 GREEN: implementar las funciones y los tipos `Metricas` y `MetricasProyecto` (ver data-model.md, `null` significa "sin datos") en `web/src/features/dashboard/metricas.ts` — tests en verde
- [ ] T009 COMMIT `GREEN: implementar las funciones del Dashboard, test en verde`
- [ ] T010 RED: en `web/src/features/dashboard/api.test.ts`, probar que `consultarMetricas('7')` hace `GET /projects/7/metricas` y devuelve la respuesta — confirmar que FALLA
- [ ] T011 COMMIT `RED: agregar test de la consulta de métricas (falla)`
- [ ] T012 GREEN: implementar `consultarMetricas` en `web/src/features/dashboard/api.ts` con `get` del cliente compartido — tests en verde
- [ ] T013 COMMIT `GREEN: implementar la consulta de métricas, test en verde`

**Checkpoint**: la base está lista para las historias.

---

## Phase 3: User Story 1 - Comparar lo comprometido y lo completado (Priority: P1) 🎯 MVP

**Goal**: la sección "Dashboard" muestra el gráfico de velocidad por Sprint finalizado.

**Independent Test**: con una respuesta simulada de tres Sprints, se entra a la ruta y se verifican los Sprints en orden y sus dos valores.

- [ ] T014 [US1] RED: en `web/src/features/dashboard/GraficoVelocidad.test.tsx`, reemplazar la prueba temporal de T004 por los escenarios: US1.1 (el gráfico se monta y hay leyenda "Comprometidos" y "Completados"), US1.2 (8 y 5 en la tabla "Ver valores por Sprint"), US1.3 (orden por número), US1.4 (un Sprint con ceros sigue en la tabla) y US1.5 (la pantalla muestra todos los Sprints que recibe, sin filtrar por estado: que solo lleguen los finalizados lo garantiza HU-03, con el escenario Godog "Un Sprint activo no aporta a las métricas") — confirmar que FALLA
- [ ] T015 [US1] COMMIT `RED: agregar tests del gráfico de velocidad (falla)`
- [ ] T016 [US1] GREEN: implementar `web/src/features/dashboard/GraficoVelocidad.tsx` y sus estilos en `web/src/features/dashboard/dashboard.css`, siguiendo la skill `dataviz` para los colores y R4 de research (series distinguibles sin color, tabla plegada, desplazamiento dentro del contenedor, eje Y con enteros) — tests en verde
- [ ] T017 [US1] COMMIT `GREEN: implementar el gráfico de velocidad, test en verde`
- [ ] T018 [US1] RED: en `web/src/features/dashboard/index.test.tsx`, probar US4.1 anticipado (indicador de carga durante la consulta), que con Sprints finalizados se muestra el gráfico, que entrar hace una sola consulta y solo `GET` (FR-002, FR-003), y que al salir y volver se consulta de nuevo y se ven los valores nuevos (caso límite); en `web/src/rutas.test.tsx`, agregar una prueba propia para `/projects/42/dashboard` con `fetch` simulado que devuelve una respuesta de métricas (el stub global devuelve `[]`, así que no sirve la tabla de rutas), y probar que la navegación ofrece "Dashboard" — confirmar que FALLA
- [ ] T019 [US1] COMMIT `RED: agregar tests de la sección Dashboard (falla)`
- [ ] T020 [US1] GREEN: implementar `web/src/features/dashboard/index.tsx` (exporta `Component`, con `key` por proyecto y el patrón `vigente` de Sprints); agregar la entrada a `web/src/navegacion.ts` y la ruta `dashboard/*` con `lazy` a `web/src/rutas.tsx` — tests en verde
- [ ] T021 [US1] COMMIT `GREEN: implementar la sección Dashboard, test en verde`
- [ ] T022 [US1] (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: el gráfico funciona de forma independiente.

---

## Phase 4: User Story 2 - Ver el resto de las métricas en indicadores (Priority: P1)

**Goal**: un indicador visual por cada métrica restante, con los valores del proyecto.

**Independent Test**: con valores conocidos, cada indicador muestra el valor de la respuesta; sin horas por Story Point, muestra "sin datos".

- [ ] T023 [US2] RED: en `web/src/features/dashboard/Indicadores.test.tsx`, probar US2.1 (los siete indicadores), US2.2 (50% con una barra de proporción accesible), US2.3 (32 y 30 horas en una misma escala), US2.4 (-2 h, -6,25% y "menos de lo estimado"), US2.5 (2 y 1 defectos en una misma escala), US2.6 (`null` → "sin datos", sin barra, horas reales visibles y enlace a `/projects/:id/proyecto/edit`) y los casos límite: horas en valor con porcentaje `null`; defectos y horas en cero se muestran como `0`, sin barras rotas — confirmar que FALLA
- [ ] T024 [US2] COMMIT `RED: agregar tests de los indicadores del Dashboard (falla)`
- [ ] T025 [US2] GREEN: implementar `web/src/features/dashboard/Indicadores.tsx` y sus estilos en `web/src/features/dashboard/dashboard.css`, con HTML y CSS, sin Recharts (research R5) — tests en verde
- [ ] T026 [US2] COMMIT `GREEN: implementar los indicadores del Dashboard, test en verde`
- [ ] T027 [US2] RED: en `web/src/features/dashboard/index.test.tsx`, probar que con Sprints finalizados los indicadores usan los valores de `proyecto`, no los de un Sprint — confirmar que FALLA
- [ ] T028 [US2] COMMIT `RED: agregar test de los indicadores en la sección (falla)`
- [ ] T029 [US2] GREEN: mostrar `Indicadores` en `web/src/features/dashboard/index.tsx` — tests en verde
- [ ] T030 [US2] COMMIT `GREEN: mostrar los indicadores en el Dashboard, test en verde`

**Checkpoint**: el Escenario 1 de la historia está completo.

---

## Phase 5: User Story 3 - Entender por qué todavía no hay métricas (Priority: P1)

**Goal**: estado vacío informativo sin Sprints finalizados.

**Independent Test**: con `sprints: []`, aparece el estado vacío, sin gráfico, indicadores ni error.

- [ ] T031 [US3] RED: en `web/src/features/dashboard/index.test.tsx`, probar US3.1 (mensaje, sin gráfico ni indicadores), US3.2 (`sprints: []` con valores del proyecto en cero sigue vacío) y US3.3 (enlace a `/projects/:id/sprints`) — confirmar que FALLA
- [ ] T032 [US3] COMMIT `RED: agregar tests del estado vacío del Dashboard (falla)`
- [ ] T033 [US3] GREEN: implementar el estado vacío en `web/src/features/dashboard/index.tsx`, decidido por `sprints.length`, con sus estilos en `web/src/features/dashboard/dashboard.css` — tests en verde
- [ ] T034 [US3] COMMIT `GREEN: implementar el estado vacío del Dashboard, test en verde`

**Checkpoint**: el Escenario 2 de la historia está completo.

---

## Phase 6: User Story 4 - Saber qué pasó si las métricas no se pueden obtener (Priority: P2)

**Goal**: los errores se distinguen de la falta de datos.

**Independent Test**: se simula cada falla y se verifica el mensaje, sin gráfico, indicadores ni estado vacío.

- [ ] T035 [US4] RED: en `web/src/features/dashboard/index.test.tsx`, probar US4.2 (404 → "proyecto no encontrado"), US4.3 (`ConnectionError` → su mensaje) y un 500 con su mensaje; en todos, sin gráfico, sin indicadores y sin estado vacío — confirmar que FALLA
- [ ] T036 [US4] COMMIT `RED: agregar tests de los errores del Dashboard (falla)`
- [ ] T037 [US4] GREEN: mostrar el `message` del error con `ErrorMessage` en `web/src/features/dashboard/index.tsx` (research R6) — tests en verde
- [ ] T038 [US4] COMMIT `GREEN: implementar los errores del Dashboard, test en verde`

**Checkpoint**: todas las historias funcionan.

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T039 Verificar en `web/dist` tras `npm run build` que Recharts queda en un chunk aparte de la entrada, y anotar el tamaño en R2 de `specs/035-hu-09-dashboard-metricas/research.md`
- [ ] T040 COMMIT `docs: registrar la carga diferida de Recharts en HU-09` con `research.md`
- [ ] T041 `npm run lint`, `npm test` y `npm run build` desde `web/` sin errores
- [ ] T042 Ejecutar la validación manual de `specs/035-hu-09-dashboard-metricas/quickstart.md`, incluida la de 320 px, y anotar el resultado en ese archivo
- [ ] T043 Actualizar el estado en `specs/035-hu-09-dashboard-metricas/spec.md` y hacer COMMIT `docs: registrar la validación de HU-09`

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **US1** → **US2** → **US3** → **US4** → **Polish**.
- US1 crea `index.tsx`, la ruta y la navegación. US2, US3 y US4 agregan estados a ese archivo, así que van de a una.
- `GraficoVelocidad.tsx` (T014 a T017) e `Indicadores.tsx` (T023 a T026) son independientes entre sí.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.

## Parallel Opportunities

- Después de la Phase 2, una persona puede tomar `GraficoVelocidad` (T014 a T017) y otra `Indicadores` (T023 a T026).
- T006 a T009 y T010 a T013 tocan archivos distintos.

## Implementation Strategy

1. Setup + Foundational → Recharts verificada y base lista.
2. US1 → el gráfico se puede mostrar (MVP).
3. US2 + US3 → los dos escenarios de la historia completos; desde acá se puede entregar.
4. US4 y Polish → errores y validación final.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- Las sub-issues #70, #71 y #72 quedaron ajustadas a la decisión de research R1. Los escenarios de #70 están mapeados en `quickstart.md`.
- `web/src/compartido/router/LayoutProyecto.tsx` no cambia: la navegación sale de `navegacion.ts`. Si #231 se mergea antes, actualizar la rama antes de T018.
