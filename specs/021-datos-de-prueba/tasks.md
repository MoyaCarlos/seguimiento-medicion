---

description: "Tareas de la spec 021: cargar datos de prueba (TDD con commits por fase)"
---

# Tasks: Cargar datos de prueba (issue #144, versión 1)

**Input**: `specs/021-datos-de-prueba/` (spec, plan, research, data-model, contracts, quickstart)

**Tests**: OBLIGATORIOS (Principio I). Cada tarea RED va seguida de su COMMIT, luego GREEN y su COMMIT. Los COMMIT agregan archivos por nombre (nunca `git add -A`) y terminan con la línea `Co-Authored-By`. Ninguna tarea hace `git push`.

## Phase 1: Setup

- [ ] T001 RED: test de `cargarConfig` con `DEMO_HABILITADO` ("true" habilita; vacío, "false" o "TRUE" no) en cmd/api/config_test.go — confirmar que FALLA
- [ ] T002 COMMIT `RED: agregar test de DEMO_HABILITADO en la configuración (falla)`
- [ ] T003 GREEN: campo `DemoHabilitado bool` en cmd/api/config.go (`getenv("DEMO_HABILITADO") == "true"`) — tests en verde
- [ ] T004 COMMIT `GREEN: leer DEMO_HABILITADO en la configuración, test en verde`

## Phase 2: Foundational (esqueleto de la carga y bloque de usuarios)

**⚠️ Bloquea todas las historias.**

- [ ] T005 RED: tests del bloque de usuarios en internal/demo/bloque_usuarios_test.go sobre una base en memoria con las migraciones reales (helper de prueba que abre, migra y devuelve `*sql.DB`): crea Beta Uno a Beta Cuatro; si "Beta Dos" ya existe no se duplica ni falla y se reutiliza su id; el contador cuenta solo los creados — confirmar que FALLA
- [ ] T006 COMMIT `RED: agregar tests del bloque de usuarios de la carga (falla)`
- [ ] T007 GREEN: internal/demo/carga.go (tipo `carga` con día base, ids por nombre/clave y contadores; `type bloque func(ctx, *sql.Tx, *carga) error`) e internal/demo/bloque_usuarios.go (`NewUser`, `INSERT OR IGNORE` por `nombre_normalizado`, lectura de ids) — tests en verde
- [ ] T008 COMMIT `GREEN: implementar el bloque de usuarios de la carga, test en verde`

## Phase 3: User Story 1 - Cargar el proyecto de demostración (P1) 🎯 MVP

**Goal**: la carga crea el proyecto 25 completo.

**Independent Test**: base migrada + `Cargar` ⇒ el contenido de `data-model.md` y el resumen.

- [ ] T009 [US1] RED: tests del bloque de proyecto y equipo en internal/demo/bloque_proyecto_test.go: id 25, nombre «Licencia de conducir (datos de prueba)», descripción, `horas_por_story_point = 4`, inicio `hoy-15`, fin `hoy+45` (reloj fijo), roles Beta Uno `scrum_master`, Beta Dos y Beta Tres `product_builder`, Beta Cuatro sin fila en `integrantes` — confirmar que FALLA
- [ ] T010 [US1] COMMIT `RED: agregar tests del bloque de proyecto y equipo (falla)`
- [ ] T011 [US1] GREEN: internal/demo/bloque_proyecto.go con `NewProject` y `NewMembership`; el `INSERT` con id 25 y `horas_por_story_point` — tests en verde
- [ ] T012 [US1] COMMIT `GREEN: implementar el bloque de proyecto y equipo, test en verde`
- [ ] T013 [US1] RED: tests del bloque de Sprints en internal/demo/bloque_sprints_test.go: #1 Finalizado (hoy-14 a hoy-8), #2 Activo (hoy-7 a hoy+6), #3 Pendiente con goal y fechas NULL; se respetan los índices de un Activo y un Pendiente — confirmar que FALLA
- [ ] T014 [US1] COMMIT `RED: agregar tests del bloque de Sprints (falla)`
- [ ] T015 [US1] GREEN: internal/demo/bloque_sprints.go con `NewSprint`, `Iniciar` y `Cerrar` del dominio — tests en verde
- [ ] T016 [US1] COMMIT `GREEN: implementar el bloque de Sprints, test en verde`
- [ ] T017 [US1] RED: tests del bloque de historias en internal/demo/bloque_historias_test.go: 10 historias con el reparto de `data-model.md` (3 Completadas en el Sprint 1 con `fecha_finalizacion` dentro del Sprint, 2 En progreso y 1 Nueva en el Sprint 2, 1 Nueva en el Sprint 3, 3 en el backlog); prioridades M, S, C y W presentes; estimaciones solo Fibonacci y al menos una historia del backlog sin estimar; criterios en algunas; la arrastrada en `historias_arrastradas` con Sprint 1 y `estimacion_sp = 8` — confirmar que FALLA
- [ ] T018 [US1] COMMIT `RED: agregar tests del bloque de historias (falla)`
- [ ] T019 [US1] GREEN: internal/demo/bloque_historias.go con `NewBacklogItem` + `ConDatosEditados`; `En progreso` como constante local; `INSERT` de estado, sprint y fecha de finalización — tests en verde
- [ ] T020 [US1] COMMIT `GREEN: implementar el bloque de historias, test en verde`
- [ ] T021 [US1] RED: tests del bloque de horas en internal/demo/bloque_horas_test.go: registros de `data-model.md` solo en historias de los Sprints 1 y 2; cada fecha dentro de su Sprint y ≤ hoy; los tres integrantes con rol aparecen; el total de cada Completada estimada queda entre el 50 % y el 150 % de SP × 4 h — confirmar que FALLA
- [ ] T022 [US1] COMMIT `RED: agregar tests del bloque de horas (falla)`
- [ ] T023 [US1] GREEN: internal/demo/bloque_horas.go con `NewRegistroHoras` (sin `ValidarSprintParaImputar`, ver research R4) — tests en verde
- [ ] T024 [US1] COMMIT `GREEN: implementar el bloque de horas, test en verde`
- [ ] T025 [US1] RED: tests del bloque de votación en internal/demo/bloque_votacion_test.go: una ronda 1 con `revelada_en` en la historia «Avisar por correo el vencimiento» (sin estimar, en el backlog) y tres votos (3, 5, 5) de Beta Uno, Beta Dos y Beta Tres; la historia sigue sin `estimacion_sp` — confirmar que FALLA
- [ ] T026 [US1] COMMIT `RED: agregar tests del bloque de votación (falla)`
- [ ] T027 [US1] GREEN: internal/demo/bloque_votacion.go con `NuevaRonda`, `Votar` y `Revelar` — tests en verde
- [ ] T028 [US1] COMMIT `GREEN: implementar el bloque de votación, test en verde`
- [ ] T029 [US1] RED: tests de `Cargador.Cargar` en internal/demo/carga_test.go: base vacía ⇒ `Resumen` {proyecto 25, 4 usuarios, 3 integrantes, 3 Sprints, 10 historias, 15 registros, 1 ronda, 3 votos}; ninguna otra tabla (defectos, métricas) recibe filas (FR-022); las fechas se calculan con el reloj inyectado (dos días distintos ⇒ fechas distintas); `Resumen` serializa a los nombres JSON de contracts/api.md — confirmar que FALLA
- [ ] T030 [US1] COMMIT `RED: agregar tests de la carga completa (falla)`
- [ ] T031 [US1] GREEN: `Cargador`, `Resumen` y la lista ordenada de bloques en internal/demo/carga.go (una sola transacción) — tests en verde
- [ ] T032 [US1] COMMIT `GREEN: implementar la carga completa en una transacción, test en verde`

## Phase 4: User Story 2 - Idempotencia y atomicidad (P1)

**Goal**: segunda carga, id ocupado, fallo y concurrencia no dejan datos parciales.

- [ ] T033 [US2] RED: tests en internal/demo/idempotencia_test.go: segunda carga ⇒ `ErrYaCargado` y la base queda idéntica (conteos de todas las tablas); proyecto 25 con otro nombre ⇒ `ErrIdOcupado` y nada cambia, ni siquiera los usuarios Beta — confirmar que FALLA
- [ ] T034 [US2] COMMIT `RED: agregar tests de idempotencia e id ocupado (falla)`
- [ ] T035 [US2] GREEN: traducir la violación de clave primaria del proyecto 25 en `Cargar`: rollback, lectura del nombre y clasificación en `ErrYaCargado`/`ErrIdOcupado` (mensajes de contracts/api.md) — tests en verde
- [ ] T036 [US2] COMMIT `GREEN: rechazar la carga repetida o con el id ocupado, test en verde`
- [ ] T037 [US2] RED: test de fallo inyectado: `Cargador` con un bloque extra que falla al final ⇒ error y ninguna tabla con filas (incluidos los usuarios nuevos) — confirmar que FALLA
- [ ] T038 [US2] COMMIT `RED: agregar test de carga fallida que no deja datos (falla)`
- [ ] T039 [US2] GREEN: `Cargador` recibe la lista de bloques (constructor con la lista por defecto y variante para pruebas) y hace rollback ante cualquier error — tests en verde
- [ ] T040 [US2] COMMIT `GREEN: deshacer toda la carga ante un fallo, test en verde`
- [ ] T041 [US2] RED: test de concurrencia en internal/demo/concurrencia_test.go sobre una base en archivo: (a) 8 cargas simultáneas con un `*sql.DB` ⇒ exactamente una OK y las demás `ErrYaCargado`; (b) dos `*sql.DB` sobre el mismo archivo ⇒ una OK y una `ErrYaCargado`; en ambos casos 1 proyecto, 10 historias, sin duplicados — confirmar si FALLA o ya pasa (si ya pasa, registrarlo en el mensaje del commit)
- [ ] T042 [US2] COMMIT `RED: agregar test de cargas simultáneas (falla)` (o `test: ...` si ya pasa)
- [ ] T043 [US2] GREEN (solo si falla): ajustar `Cargar` para que la primera sentencia de la transacción sea una escritura (plan, decisión 3) — tests en verde
- [ ] T044 [US2] COMMIT `GREEN: resolver las cargas simultáneas, test en verde` (omitir si T043 no hizo falta)

## Phase 5: User Story 3 - Habilitación controlada y API (P1)

- [ ] T045 [US3] RED: tests del handler en internal/http/demo_handler_test.go con un `Cargador` falso: deshabilitada ⇒ GET y POST responden 404 con el cuerpo estándar y no llaman al cargador; habilitada ⇒ GET 200 `{"habilitado":true}`; POST 201 con el resumen; `ErrYaCargado` ⇒ 409 con `campo: "ya_cargado"`; `ErrIdOcupado` ⇒ 409 con `campo: "proyecto_id"`; otro error ⇒ 500 con el mensaje genérico — confirmar que FALLA
- [ ] T046 [US3] COMMIT `RED: agregar tests del handler de datos de prueba (falla)`
- [ ] T047 [US3] GREEN: internal/http/demo_handler.go con la interfaz `CargadorDatosPrueba`, `soloSiHabilitado` (único punto de control, recibe `func() bool`) y los dos handlers — tests en verde
- [ ] T048 [US3] COMMIT `GREEN: implementar el handler de datos de prueba, test en verde`
- [ ] T049 [US3] Agregar `/demo` a `PrefijosAPI` en internal/http/spa.go y `'/demo': apiProxy` en web/vite.config.ts (el test internal/http/spa_prefijos_test.go existente exige que coincidan); registrar `GET /demo/estado` y `POST /demo/datos` en cmd/api/main.go con `cfg.DemoHabilitado` y `demo.NewCargador(db, ...)`
- [ ] T050 [US3] COMMIT `Cablear las rutas /demo y el proxy de Vite`

## Phase 6: BDD (Godog)

- [ ] T051 Escribir features/datos_prueba.feature en español con tag `@Issue-144` (no es una HU, por eso no `@HU-XX`), Característica con Como/Quiero/Para y cuatro escenarios: carga exitosa, ya cargado, id ocupado, deshabilitado; revisar las cuatro categorías (normal, alternativo, límite, error) e informar cuáles no aplican — sin steps todavía
- [ ] T052 COMMIT `Agregar escenarios BDD de la carga de datos de prueba`
- [ ] T053 RED: confirmar que `go test ./features/...` FALLA por steps indefinidos (Strict) — escribir features/steps_datos_prueba.go con los steps sobre el handler real y una base en memoria migrada; registrar `InitializeScenarioDatosPrueba` en features/features_test.go solo en T055
- [ ] T054 COMMIT `RED: agregar pasos de los escenarios de datos de prueba (falla)`
- [ ] T055 GREEN: registrar el inicializador en features/features_test.go — `go test ./features/...` en verde
- [ ] T056 COMMIT `GREEN: automatizar los escenarios de datos de prueba, test en verde`

## Phase 7: User Story 4 - Botón en el frontend (P2)

- [ ] T057 [US4] RED: pruebas en web/src/compartido/demo/BotonDatosDePrueba.test.tsx con `fetch` simulado y el escenario Given-When-Then en el nombre de cada prueba: estado 200 ⇒ botón visible; estado 404 ⇒ sin botón; estado con error de red ⇒ sin botón; 201 ⇒ aviso y enlace «Abrir el proyecto de prueba» a `/projects/25`; 409 `ya_cargado` ⇒ aviso (rol `status`, no `alert`) con el enlace; 409 `proyecto_id` ⇒ mensaje de error; 500 ⇒ mensaje de error; botón deshabilitado mientras carga — confirmar que FALLA
- [ ] T058 [US4] COMMIT `RED: agregar pruebas del botón de datos de prueba (falla)`
- [ ] T059 [US4] GREEN: web/src/compartido/demo/api.ts (`datosDePruebaHabilitados`, `cargarDatosDePrueba` sobre `compartido/api/client`) y web/src/compartido/demo/BotonDatosDePrueba.tsx reutilizando `Button`, `ErrorMessage` y `LoadingIndicator` — pruebas en verde
- [ ] T060 [US4] COMMIT `GREEN: implementar el botón de datos de prueba, test en verde`
- [ ] T061 [US4] Montar `<BotonDatosDePrueba />` con una línea (más su import) en web/src/features/proyecto/crear.tsx; verificar que `crear.test.tsx` y `aislamiento.test.ts` siguen en verde
- [ ] T062 [US4] COMMIT `Montar el botón de datos de prueba en el inicio`

## Phase 8: Polish

- [ ] T063 Verificar contra quickstart.md: `go build ./...`, `go vet ./...`, `gofmt -l .`, `go test ./...`, y en web/: `npm run lint`, `npm test`, `npm run build`; corregir lo que falle (cada corrección con su propio commit)
- [ ] T064 Confirmar que la sección «Alcance de la versión 1 y evolución» de spec.md coincide con lo implementado y actualizar el estado de la spec; marcar estas tareas como hechas
- [ ] T065 COMMIT `Cerrar la spec 021 y marcar las tareas como hechas`

## Dependencias

Phase 1 y 2 → Phase 3 (US1) → Phase 4 (US2) → Phase 5 (US3) → Phase 6 (BDD). Phase 7 (US4) depende solo del contrato (contracts/api.md) y puede hacerse en paralelo con las fases 4–6. Polish al final.

## Estrategia

MVP = Phases 1–3 (carga funcional por función Go). Después atomicidad/idempotencia, API, BDD y botón. Una fase a la vez, con revisión humana del diff antes de seguir.
