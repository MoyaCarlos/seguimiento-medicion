---

description: "Tareas de la spec 041 (datos de prueba v2): TDD con commits por fase"
---

# Tasks: Datos de prueba v2 (issue #235)

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/api.md](contracts/api.md)

**Tests**: obligatorios (Principio I). Cada fase: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Los COMMIT son tareas propias y no se marcan hasta haber corrido `git commit`. Ninguna tarea hace `git push` ni merge.

**Nota sobre el RED**: cambiar la firma de `NewCargador` (fase A) hace que el paquete de pruebas no compile; ese fallo de compilación cuenta como RED. Las pruebas de la v1 que citan «Beta» o «3 Sprints» se actualizan en la fase que cambia ese comportamiento, así cada fase deja `go test ./...` en verde.

**Recortables si el tiempo aprieta**: las fases **E** y **F** (frontend y README). A–D son el mínimo para que Render tenga datos y se pueda iniciar sesión por API.

## Phase 0: Preparación

- [ ] T001 Confirmar con `go build ./... && go test ./...` y `cd web && npm test` que `main` parte en verde (sin commit)

---

## Phase A: Usuarios con contraseña y 409 por usuario existente (US1, US2) 🎯

**Goal**: cuatro usuarios reales con `demo1234` hasheado una sola vez, un Scrum Master, equipo de 4, y 409 `usuario_existente` con prioridad `ya_cargado` > usuario > id ocupado.

**Independent Test**: cargar con hasheador barato, iniciar sesión con los cuatro, y repetir la carga con un usuario preexistente.

- [X] T002 [US1] RED: pruebas en `internal/demo/bloque_usuarios_test.go` (4 usuarios de la lista con `contrasena_hash` verificable con `demo1234`, un mismo hash para los cuatro, nombres válidos según `ValidarNombreUsuario`, no existe ningún «Beta»; con un `Hasheador` falso, `Cargar` lo invoca **una sola vez** aunque haya cuatro usuarios, y también cuando la carga termina en 409) y en `bloque_proyecto_test.go` (4 integrantes, único Scrum Master `carlos.moya`); en `carga_test.go`, tras `Cargar`, `ServicioContrasenas.Verificar("demo1234", hash)` es verdadero para los cuatro y el hash guardado no contiene `demo1234`; adaptar `ayudas_test.go` (`cargador` y `ejecutar` reciben un hasheador `service.NewServicioContrasenas(1)` y su hash) — confirmar que FALLAN (no compilan)
- [X] T003 COMMIT `RED: agregar tests de usuarios con contraseña de la demo (falla)`
- [X] T004 [US1] GREEN: `internal/demo/carga.go` (`Hasheador`, `NewCargador(db, ahora, hasheador)`, hash único calculado en `Cargar` antes de la transacción, `carga.hashContrasena`, `Resumen.Integrantes` = 4), `bloque_usuarios.go` (lista `usuariosDemo` con rol, `ValidarNombreUsuario` + `NewUser`, `INSERT` estricto con el hash), `bloque_proyecto.go` (equipo desde `usuariosDemo`); renombrar las claves «Beta Uno/Dos/Tres» de `bloque_horas.go` y `bloque_votacion.go` a `carlos.moya`/`ana.iriarte`/`candela.vulcano` sin cambiar los datos; adaptar `carga_test.go`, `concurrencia_test.go` y `atomicidad_test.go` (usuarios 4, integrantes 4; una carga fallida no deja ningún usuario) y las demás pruebas que citan «Beta»; `cmd/api/main.go` (pasa el `hasheador` existente a `demo.NewCargador`) y `features/steps_datos_prueba.go` (hasheador barato) para que `go build ./...` no se rompa entre commits; `go build ./... && go test ./...` en verde
- [X] T005 COMMIT `GREEN: usuarios reales con contraseña hasheada una sola vez, test en verde`
- [X] T006 [US2] RED: pruebas en `internal/demo/idempotencia_test.go` — usuario existente ⇒ `ErrUsuarioExistente` con el nombre del primero en conflicto (también con otra capitalización) y base sin cambios; existe el proyecto y el usuario ⇒ `ErrYaCargado`; proyecto 25 de otro nombre sin conflicto de usuarios ⇒ `ErrIdOcupado`; usuario existente y proyecto 25 de otro nombre ⇒ `ErrUsuarioExistente` — confirmar que FALLAN
- [X] T007 [US2] RED: prueba en `internal/http/demo_handler_test.go` — `ErrUsuarioExistente` ⇒ 409 `{"campo":"usuario_existente","mensaje":"ya existe el usuario ana.iriarte"}` — confirmar que FALLA
- [X] T008 COMMIT `RED: agregar tests del conflicto por usuario existente (falla)`
- [X] T009 [US2] GREEN: `ErrUsuarioExistente` y clasificador de conflictos en `internal/demo/carga.go` (orden: `ya_cargado` > usuario > `id ocupado` > error original), mapeo en `internal/http/demo_handler.go` — tests en verde
- [X] T010 COMMIT `GREEN: 409 por usuario de la demo existente, test en verde`
- [X] T011 REFACTOR (solo si hay algo que limpiar): renombrar `bloquesV1` a `bloquesDemo` y limpiar nombres que quedaron de la v1; + COMMIT `REFACTOR: renombrar la lista de bloques de la carga`

**Checkpoint**: se puede iniciar sesión con los cuatro usuarios; los 409 funcionan.

---

## Phase B: Cuatro Sprints y reubicación de historias, horas y votaciones (US3)

**Goal**: dos Sprints finalizados, uno activo y uno pendiente; 12 historias; 20 registros de horas; ronda revelada y ronda abierta.

**Independent Test**: tras la carga, los conteos y el reparto de [data-model.md](data-model.md).

- [X] T012 [US3] RED: `internal/demo/bloque_sprints_test.go` — 4 Sprints (Finalizado −28/−15, Finalizado −14/−8, Activo −7/+6, Pendiente sin goal ni fechas) y proyecto con inicio hoy−30 en `bloque_proyecto_test.go` — confirmar que FALLAN
- [X] T013 [US3] RED: `internal/demo/bloque_historias_test.go` — 12 historias con el reparto 3/2/3/1/3, velocidades 13 y 8, fechas de finalización dentro de su Sprint, arrastrada `tasa` en el Sprint 2 con 8 SP y en estado Nueva, al menos una sin estimar en el backlog — confirmar que FALLAN
- [X] T014 [US3] RED: `bloque_horas_test.go` (20 registros en los Sprints 1 a 3, fechas dentro del Sprint y ≤ hoy, los cuatro integrantes con horas, Completadas entre el 50 % y el 150 % de SP × 4 h) y `bloque_votacion_test.go` (ronda revelada con 4 votos en `aviso`, ronda abierta con 2 votos y `revelada_en NULL` en `digital`; resumen 2 rondas y 6 votos) — confirmar que FALLAN
- [X] T015 COMMIT `RED: agregar tests de cuatro Sprints, historias, horas y votaciones (falla)`
- [X] T016 [US3] GREEN: `bloque_proyecto.go` (inicio `c.dia(-30)`), `bloque_sprints.go` (lista de 4 definiciones con el mismo bucle y funciones de la v1), `bloque_historias.go` (12 filas, `sprint` 0 a 4, arrastrada en `c.sprints[1]`), `bloque_horas.go` (20 filas, tres Sprints con sus fechas), `bloque_votacion.go` (dos rondas) — adaptar `carga_test.go` y `concurrencia_test.go` (sprints 4, historias 12, horas 20, rondas 2, votos 6) — `go test ./internal/demo/...` en verde
- [X] T017 COMMIT `GREEN: cuatro Sprints, historias reubicadas, horas y rondas de votación, test en verde`
- [X] T018 REFACTOR (solo si hay algo que limpiar, p. ej. duplicación entre las dos rondas) + COMMIT `REFACTOR: <qué se limpió>` (sin cambios: no había nada que limpiar)

**Checkpoint**: historial con dos Sprints comparables; `go test ./...` en verde.

---

## Phase C: Bloque de defectos y conteo en el resumen (US4)

**Goal**: 5 defectos (3 Resueltos, 2 Pendientes; Alta/Media/Baja) ligados a historias de los Sprints 1 a 3.

**Independent Test**: tras la carga, `SELECT` sobre `defectos` y el campo `defectos` del resumen.

- [X] T019 [US4] RED: `internal/demo/bloque_defectos_test.go` — 5 filas con las tres severidades; 3 Resueltos con `sprint_resolucion_id` ≥ el de detección y 2 Pendientes sin él; historia y Sprints del proyecto 25 (FK compuestas); `resumen.Defectos` = 5; `carga_test.go` — JSON del resumen incluye `"defectos":5`; `idempotencia_test.go` — `defectos` en `tablasDeLaCarga` — confirmar que FALLAN
- [X] T020 COMMIT `RED: agregar tests del bloque de defectos (falla)`
- [X] T021 [US4] GREEN: `internal/demo/bloque_defectos.go` (`NewDefecto`, `Resolver`, INSERT con `creado_en` relativo), `Resumen.Defectos` y `bloqueDefectos` al final de `bloquesDemo` en `carga.go` — adaptar `concurrencia_test.go` (defectos 5) — tests en verde
- [X] T022 COMMIT `GREEN: bloque de defectos y conteo en el resumen, test en verde`
- [X] T023 REFACTOR (solo si hay algo que limpiar) + COMMIT `REFACTOR: <qué se limpió>` (sin cambios: no había nada que limpiar)

**Checkpoint**: la carga completa coincide con [data-model.md](data-model.md).

---

## Phase D: Escenarios Godog de `features/datos_prueba.feature`

- [X] T024 Actualizar `features/datos_prueba.feature` (en español, tag `@HU-…` existente) a partir de los criterios de la spec 041: el escenario de éxito pasa a «12 historias, 4 Sprints y 4 integrantes» y suma defectos; nuevos escenarios: inicio de sesión con un usuario de la demo, usuario existente ⇒ 409 con su nombre y sin cambios, `ya_cargado` aunque haya usuarios; revisar las 4 categorías (normal, alternativo, límite, error) y decir cuáles no aplican
- [X] T025 RED: pasos pendientes en `features/steps_datos_prueba.go` — `go test ./features/` falla por pasos sin definir (nota: no hubo RED; el comportamiento ya existía, ver commit 4899cbe)
- [X] T026 COMMIT `RED: agregar escenarios BDD de datos de prueba v2 (falla)`
- [X] T027 GREEN: definir los pasos (hasheador barato, login con `ServicioContrasenas`, conteo de defectos y Sprints) — `go test ./features/` en verde
- [X] T028 COMMIT `GREEN: pasos BDD de datos de prueba v2, test en verde`

---

## Phase E: Frontend `/ingresar` (US5) — recortable

- [ ] T029 [US5] RED: `web/src/compartido/demo/BotonDatosDePrueba.test.tsx` — con el nombre del escenario en cada prueba: 201 ⇒ aviso + tabla con los cuatro usuarios, `demo1234` y roles; 409 `ya_cargado` ⇒ aviso + tabla; 409 `usuario_existente` ⇒ aviso con el nombre y sin tabla; 409 `proyecto_id` y 500 ⇒ error; botón deshabilitado mientras carga; sin botón si `/demo/estado` no es 200 — confirmar que FALLAN
- [ ] T030 [US5] RED: `web/src/features/autenticacion/Ingresar.test.tsx` — con `/demo/estado` habilitado aparece «Cargar datos de prueba»; con 404 no aparece; el `fetch` simulado distingue por URL para no romper las pruebas de ingreso existentes — confirmar que FALLAN
- [ ] T031 COMMIT `RED: agregar pruebas de la pantalla de ingreso con datos de prueba (falla)`
- [ ] T032 [US5] GREEN: `web/src/compartido/demo/CredencialesDemo.tsx` (tabla constante), `BotonDatosDePrueba.tsx` (tabla tras 201 y `ya_cargado`, aviso para `usuario_existente`), `web/src/features/autenticacion/Ingresar.tsx` (monta el botón bajo el formulario); `npm test`, `npm run lint` y `npm run build` en verde
- [ ] T033 COMMIT `GREEN: botón y tabla de credenciales de la demo en /ingresar, test en verde`
- [ ] T034 REFACTOR (solo si hay algo que limpiar, p. ej. estilos) + COMMIT `REFACTOR: <qué se limpió>`

---

## Phase F: README — recortable

- [ ] T035 Documentar en `README.md` (sección de Render y datos de prueba) los cuatro usuarios, la contraseña `demo1234`, sus roles y el 409 por usuario existente
- [ ] T036 COMMIT `Documentar los usuarios y la contraseña de la demo en el README`

---

## Phase G: Cierre

- [ ] T037 `gofmt -l .`, `go build ./...`, `go vet ./...` y `go test ./...` sin errores; `cd web && npm run lint && npm test && npm run build`
- [ ] T038 Recorrer [quickstart.md](quickstart.md) a mano y anotar el resultado en el PR (la carga contra Render la hace el equipo)
- [ ] T039 Revisar `git log` de la rama: cada ciclo con sus commits `RED:`/`GREEN:`; confirmar que no hay push ni PR hechos por el agente

## Dependencias

- T001 → A → B → C → D (los pasos de D ven los datos finales). E y F dependen solo del contrato de A (409 `usuario_existente`) y pueden hacerse en paralelo con B–D por otra persona (archivos distintos), pero se commitean aparte.
- Dentro de cada fase: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT).
- Si el tiempo aprieta, se mergean A–D y se dejan E y F como PR aparte o recortados.
