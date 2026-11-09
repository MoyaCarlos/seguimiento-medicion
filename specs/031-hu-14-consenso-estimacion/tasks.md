---

description: "Tareas de HU-14: Consenso y Registro de la Estimación Acordada (TDD con commits por fase)"
---

# Tasks: Consenso y Registro de la Estimación Acordada (HU-14)

**Input**: Design documents from `/specs/031-hu-14-consenso-estimacion/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica es RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT).

**Compilación en RED**: todos los commits compilan. En cada RED se agregan, junto al test, solo las
firmas nuevas que el test necesita con un cuerpo vacío (devolver valor cero / `nil`), de modo que
`go build ./...` pase y el test falle por aserción, no por símbolo indefinido. El RED de un service
o handler que cambia un constructor actualiza también a quienes lo llaman.

**Organization**: por user story de la spec (US1 diferencias, US2 estimación acordada, US3 nueva ronda
e historial). Los tests de cada historia siembran las rondas directo en el fake o en la base, así cada
una se prueba sola.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3

---

## Phase 1: Setup

**Purpose**: dejar los escenarios BDD escritos antes de implementar (flujo de AGENTS.md).

- [X] T001 Escribir features/consenso_estimacion.feature en Gherkin español: tag `@HU-14`, Característica con el Como/Quiero/Para de la Descripción de la spec y Antecedentes reutilizando los pasos de features/steps_votacion.go (`que el proyecto de la votación tiene a "Carlos" como Scrum Master y a "Candela" y "Jimena" como Product Builders` y `que el proyecto de la votación tiene una historia en estado "Nueva"`), sin redefinir pasos existentes ni repetir sus textos. Exactamente 20 escenarios, uno por Acceptance Scenario de spec.md: US1.1 a US1.3, US2.1 a US2.9 y US3.1 a US3.8. Comentario en el archivo con la cobertura por categoría, como en features/votacion_planning_poker.feature (normal: US1.1, US2.1, US3.1; alternativo: US1.2, US2.2, US2.8, US3.2, US3.7, US3.8; límite: US1.3, US2.3, US2.9; error: US2.4 a US2.7 y US3.3 a US3.6). Sin pasos implementados: con el runner estricto, `go test ./features/` FALLA
- [X] T002 COMMIT `RED: agregar escenarios BDD de consenso y estimación acordada (fallan)`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: error compartido por US2 y US3.

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

- [X] T003 Agregar en internal/domain/errors.go, en el bloque de errores de la votación, `ErrRondaNoRevelada = ValidationError{Campo: "estado", Mensaje: "la ronda actual todavía no fue revelada"}` (research.md R10). Sin lógica
- [X] T004 COMMIT `Agregar el error de ronda sin revelar`

**Checkpoint**: las historias pueden empezar.

---

## Phase 3: User Story 1 - Detectar diferencias entre los votos (Priority: P1) 🎯 MVP

**Goal**: la consulta de una ronda revelada indica si hay diferencias y, si las hay, el mínimo y el máximo.

**Independent Test**: con una ronda revelada con votos 3, 5 y 8 sembrada en la base, `GET /backlog/{id}/votacion` devuelve `diferencias: {hay_diferencias: true, minimo: 3, maximo: 8}`; con la ronda oculta, `diferencias: null`.

- [X] T005 [US1] RED: en internal/domain/votacion_test.go, tabla de `CoincidenciaExacta{}.HayDiferencias(valores []int) bool`: `[3,5,8]` y `[3,5]` → `true`; `[5,5,5]`, `[8]` y `[]` → `false`. Stubs en internal/domain/votacion.go de la interfaz `ReglaConsenso` y de `CoincidenciaExacta` — confirmar que FALLA
- [X] T006 [US1] COMMIT `RED: agregar tests de la regla de consenso (fallan)`
- [X] T007 [US1] GREEN: implementar `CoincidenciaExacta.HayDiferencias` en internal/domain/votacion.go — `go test ./...` en verde
- [X] T008 [US1] COMMIT `GREEN: implementar la regla de consenso, tests en verde`
- [X] T009 [US1] RED: en internal/domain/votacion_test.go, tests de `(RondaVotacion).Diferencias(regla ReglaConsenso) (Diferencias, bool)`: sin revelar → `false`; revelada con votos 3, 5 y 8 → `{Hay: true, Minimo: 3, Maximo: 8}`, `true`; revelada con votos 5, 5 y 5 → `Hay: false`, `true`; una regla de prueba que siempre responde `true` cambia el resultado (prueba que la regla es intercambiable); revelada sin votos, con esa misma regla → `Hay: false`, `true`, sin entrar en pánico. Stubs del tipo `Diferencias{Hay bool; Minimo, Maximo int}` y del método — confirmar que FALLA
- [X] T010 [US1] COMMIT `RED: agregar tests de las diferencias de una ronda (fallan)`
- [X] T011 [US1] GREEN: implementar `Diferencias` en internal/domain/votacion.go, con los extremos calculados sobre `r.Votos`; sin votos no hay diferencias, sin consultar la regla — `go test ./...` en verde
- [X] T012 [US1] COMMIT `GREEN: implementar las diferencias de una ronda, tests en verde`
- [X] T013 [US1] RED: en internal/service/consultar_votacion_test.go, tests de que `VistaVotacion.Diferencias` (`*domain.Diferencias`) es `nil` con la ronda oculta, también para quien ya votó (FR-002), y se completa con la ronda revelada. Cambiar la firma a `NewConsultarVotacion(historias, proyectos, votaciones, regla domain.ReglaConsenso)` y actualizar a quienes lo llaman, pasando `domain.CoincidenciaExacta{}`: internal/service/consultar_votacion_test.go, internal/http/votacion_handler_test.go, features/steps_votacion.go y cmd/api/main.go — confirmar que FALLA
- [X] T014 [US1] COMMIT `RED: agregar a la consulta de la votación las diferencias (falla)`
- [X] T015 [US1] GREEN: completar `Diferencias` en `ConsultarVotacion.Ejecutar` en internal/service/consultar_votacion.go — `go test ./...` en verde
- [X] T016 [US1] COMMIT `GREEN: informar las diferencias en la consulta de la votación, test en verde`
- [X] T017 [US1] RED: en internal/http/votacion_handler_test.go, tests de `GET /backlog/{id}/votacion` con los tres casos de research.md R3: ronda oculta → `"diferencias": null`; revelada con votos 3 y 8 (`rondaRevelada` y `votosDeCarlosYCandela`) → `{"hay_diferencias": true, "minimo": 3, "maximo": 8}`; revelada con votos iguales → `{"hay_diferencias": false, "minimo": null, "maximo": null}`. Confirmar que FALLA
- [X] T018 [US1] COMMIT `RED: agregar a la respuesta de la votación las diferencias (falla)`
- [X] T019 [US1] GREEN: agregar `Diferencias *diferenciasResponse` con `json:"diferencias"` a `votacionResponse` en internal/http/votacion_dto.go, con `Minimo` y `Maximo` como `*int` — `go test ./...` en verde
- [X] T020 [US1] COMMIT `GREEN: devolver las diferencias en la respuesta de la votación, test en verde`

**Checkpoint**: la consulta informa las diferencias, sin cambiar los campos de HU-02.

---

## Phase 4: User Story 2 - Registrar la estimación acordada (Priority: P1)

**Goal**: el Scrum Master registra el valor acordado como Story Points de la historia.

**Independent Test**: con una ronda revelada sembrada, `POST /backlog/{id}/votacion/estimacion` con `{"usuario_id":1,"estimacion_sp":5}` responde 200 con `estimacion_sp: 5`, y la historia queda con ese valor.

- [X] T021 [US2] RED: en internal/domain/backlog_item_test.go, tests de `(BacklogItem).ConEstimacion(sp int) (BacklogItem, error)`: 5 → copia con `EstimacionSP` 5, sin cambiar la original; 0, 4 y 22 → `ErrEstimacionInvalida`. Stub en internal/domain/backlog_item.go — confirmar que FALLA
- [X] T022 [US2] COMMIT `RED: agregar tests de la estimación de una historia (fallan)`
- [X] T023 [US2] GREEN: implementar `ConEstimacion` con `EsEstimacionValida` en internal/domain/backlog_item.go — `go test ./...` en verde
- [X] T024 [US2] COMMIT `GREEN: implementar la estimación de una historia, tests en verde`
- [X] T025 [US2] REFACTOR: mover `validarEstimacionCongelada` de método de `EditarHistoria` a la función `validarEstimacionCongelada(ctx, sprints repository.SprintRepository, actual domain.BacklogItem, nueva *int) error` en internal/service/estimacion_congelada.go, y usarla desde internal/service/editar_historia.go, sin cambiar su comportamiento (research.md R6) — `go test ./...` en verde
- [X] T026 [US2] COMMIT `REFACTOR: compartir la regla D2b entre las escrituras de Story Points`
- [X] T027 [US2] Definir en internal/repository/backlog_repository.go el puerto `EstimacionHistoriaRepository` con `ObtenerPorID` y `GuardarEstimacion(ctx, id int64, sp int) error`, documentando `domain.ErrHistoriaNoEncontrada` si no afecta ninguna fila (data-model.md). Sin lógica
- [X] T028 [US2] COMMIT `Agregar el puerto de la estimación de una historia`
- [X] T029 [US2] RED: en internal/repository/sqlite_backlog_test.go, tests de `GuardarEstimacion`: escribe `estimacion_sp` y deja título, descripción y criterios como estaban (R7); id inexistente → `ErrHistoriaNoEncontrada`. Stub en internal/repository/sqlite_backlog.go — confirmar que FALLA
- [X] T030 [US2] COMMIT `RED: agregar tests de guardar la estimación en SQLite (fallan)`
- [X] T031 [US2] GREEN: implementar `GuardarEstimacion` con `UPDATE historias SET estimacion_sp = ? WHERE id = ?` en internal/repository/sqlite_backlog.go — `go test ./...` en verde
- [X] T032 [US2] COMMIT `GREEN: guardar la estimación en SQLite, tests en verde`
- [X] T033 [US2] RED: en internal/service/registrar_estimacion_test.go, tests de `RegistrarEstimacion.Ejecutar(ctx, historiaID, usuarioID, valor)` con un `fakeEstimacionRepository` nuevo en internal/service/fakes_test.go que cuenta las escrituras: US2.1 historia sin estimación → queda con 5; US2.2 con 3 cargados, sin Sprint y con Sprint `Pendiente` → queda con 8; US2.3 sin estimación con Sprint `Activo` → queda con 5; US2.8 valor que nadie votó → se acepta; US2.9 historia `Completada` sin estimación → queda con 5; mismo valor con Sprint `Activo` → sin error ni escritura; US2.4 valor 4 → `ErrEstimacionInvalida`; US2.5 ronda oculta → `ErrRondaNoRevelada`; sin votación → `ErrVotacionNoIniciada`; US2.6 valor distinto con Sprint `Activo` y con `Finalizado` → `ErrEstimacionCongelada`; US2.7 Product Builder → `ErrNoEsScrumMaster`; quien no integra el proyecto → `ErrNoEsScrumMaster`, igual que iniciar y revelar en HU-02; historia inexistente → `ErrHistoriaNoEncontrada`. En cada rechazo, ninguna escritura. Stub en internal/service/registrar_estimacion.go — confirmar que FALLA
- [X] T034 [US2] COMMIT `RED: agregar tests de registrar la estimación acordada (fallan)`
- [X] T035 [US2] GREEN: implementar `RegistrarEstimacion` en internal/service/registrar_estimacion.go, en el orden de research.md R10: historia → permiso → valor → ronda actual revelada → D2b → escritura — `go test ./...` en verde
- [X] T036 [US2] COMMIT `GREEN: implementar registrar la estimación acordada, tests en verde`
- [X] T037 [US2] RED: en internal/http/votacion_handler_test.go, tests de `RegistrarEstimacion` (`POST /backlog/{id}/votacion/estimacion`): éxito → 200 con la historia y `estimacion_sp`; `estimacion_sp` 4 → 400 campo `estimacion_sp`; ronda oculta → 400 campo `estado`; sin `usuario_id` → 400 campo `usuario_id`; Product Builder → 403; sin votación → 404; historia inexistente con cuerpo `{}` → 404, no 400 (FR-012). Agregar `estimacionRequest` en internal/http/votacion_dto.go, el parámetro `registrar *service.RegistrarEstimacion` a `NewVotacionHandler` y actualizar a quienes lo llaman: `nuevoVotacionHandlerDePrueba` y cmd/api/main.go. Stub del método — confirmar que FALLA
- [X] T038 [US2] COMMIT `RED: agregar tests del handler de la estimación acordada (fallan)`
- [X] T039 [US2] GREEN: implementar el método en internal/http/votacion_handler.go con `historiaDeRuta`, `leerPedido` y `aRespuesta` — `go test ./...` en verde
- [X] T040 [US2] COMMIT `GREEN: implementar el handler de la estimación acordada, tests en verde`
- [X] T041 [US2] Cablear `POST /backlog/{id}/votacion/estimacion` en cmd/api/main.go
- [X] T042 [US2] COMMIT `Cablear el registro de la estimación acordada en cmd/api`

**Checkpoint**: la estimación acordada se ve en el Product Backlog (SC-005).

---

## Phase 5: User Story 3 - Abrir una nueva ronda y consultar el historial (Priority: P2)

**Goal**: el Scrum Master abre rondas sobre una revelada, y las rondas reveladas se pueden consultar.

**Independent Test**: con la ronda 1 revelada, `POST /backlog/{id}/votacion/rondas` responde 201 con `ronda: 2` sin votos, y `GET /backlog/{id}/votacion/rondas` devuelve solo la ronda 1 con sus votos.

### Nueva ronda

- [X] T043 [US3] RED: en internal/domain/votacion_test.go, tests de `(RondaVotacion).Siguiente() (RondaVotacion, error)`: revelada → misma historia, `Numero + 1`, sin votos, `ReveladaEn` nil y sin cambiar la original; sin revelar → `ErrRondaNoRevelada`. Stub — confirmar que FALLA
- [X] T044 [US3] COMMIT `RED: agregar tests de la ronda siguiente (fallan)`
- [X] T045 [US3] GREEN: implementar `Siguiente` en internal/domain/votacion.go — `go test ./...` en verde
- [X] T046 [US3] COMMIT `GREEN: implementar la ronda siguiente, tests en verde`
- [X] T047 [US3] Ampliar `fakeVotacionRepository` en internal/service/fakes_test.go (data-model.md): conservar las rondas anteriores, responder `ErrVotacionYaIniciada` cuando el número no es mayor que el de la ronda actual, y asignar ids únicos. Los tests existentes no cambian — `go test ./...` en verde
- [X] T048 [US3] COMMIT `Conservar las rondas anteriores en el fake de votación`
- [X] T049 [US3] RED: en internal/service/abrir_nueva_ronda_test.go, tests de `AbrirNuevaRonda.Ejecutar(ctx, historiaID, usuarioID)`: US3.1 ronda 1 revelada con diferencias → ronda 2 abierta sin votos y la ronda 1 se conserva; US3.2 votar y revelar la ronda 2 calcula diferencias y promedio solo con sus votos; US3.7 revelada sin diferencias → también abre; después de registrar una estimación → también abre (clarify); US3.3 ronda oculta → `ErrRondaNoRevelada`; US3.4 Product Builder → `ErrNoEsScrumMaster`; quien no integra el proyecto → `ErrNoEsScrumMaster`; US3.5 sin votación → `ErrVotacionNoIniciada`; US3.6 `Completada` → `ErrHistoriaCompletada`; historia inexistente → `ErrHistoriaNoEncontrada`. En cada rechazo no se crea ninguna ronda. Stub en internal/service/abrir_nueva_ronda.go — confirmar que FALLA
- [X] T050 [US3] COMMIT `RED: agregar tests de abrir una nueva ronda (fallan)`
- [X] T051 [US3] GREEN: implementar `AbrirNuevaRonda` en internal/service/abrir_nueva_ronda.go: historia → permiso → no `Completada` → ronda actual → `Siguiente` → `CrearRonda` — `go test ./...` en verde
- [X] T052 [US3] COMMIT `GREEN: implementar abrir una nueva ronda, tests en verde`
- [X] T053 [US3] Tests de caracterización en internal/repository/sqlite_votacion_test.go: `CrearRonda` guarda la ronda 2 de una historia con la ronda 1, `ObtenerRondaActual` pasa a devolverla sin votos, y repetir la ronda 2 → `ErrVotacionYaIniciada`. `CrearRonda` ya existe y no cambia, así que se espera que pasen de entrada; si alguno falla, frenar y avisar, porque es un defecto del adaptador — `go test ./...` en verde
- [X] T054 [US3] COMMIT `Agregar tests de caracterización de la segunda ronda en SQLite`
- [X] T055 [US3] RED: en internal/http/votacion_handler_test.go, tests de `AbrirRonda` (`POST /backlog/{id}/votacion/rondas`): éxito → 201 con `ronda: 2`, `diferencias: null` y nadie votó; ronda oculta → 400 campo `estado`; sin `usuario_id` → 400; Product Builder → 403; sin votación → 404; historia inexistente con cuerpo `{}` → 404. Agregar el parámetro `abrir *service.AbrirNuevaRonda` a `NewVotacionHandler` y actualizar a quienes lo llaman. Stub — confirmar que FALLA
- [X] T056 [US3] COMMIT `RED: agregar tests del handler de nueva ronda (fallan)`
- [X] T057 [US3] GREEN: implementar `AbrirRonda` en internal/http/votacion_handler.go, respondiendo con `responderVista` y `StatusCreated` — `go test ./...` en verde
- [X] T058 [US3] COMMIT `GREEN: implementar el handler de nueva ronda, tests en verde`

### Historial

- [X] T059 [US3] RED: en internal/repository/sqlite_votacion_test.go, tests de `ListarRondasReveladas`: rondas 1 y 2 reveladas y 3 abierta → devuelve 1 y 2 ordenadas, con sus votos por `usuario_id`, sin la 3; ronda revelada sin votos → aparece con `Votos` vacío; historia sin rondas → lista vacía sin error. Agregar el método a `VotacionRepository` en internal/repository/votacion_repository.go, con stub en internal/repository/sqlite_votacion.go y en `fakeVotacionRepository` — confirmar que FALLA
- [X] T060 [US3] COMMIT `RED: agregar tests de listar las rondas reveladas en SQLite (fallan)`
- [X] T061 [US3] GREEN: implementar `ListarRondasReveladas` en internal/repository/sqlite_votacion.go con **una sola consulta** `rondas_votacion LEFT JOIN votos`, sin llamar a `votosDe` con filas abiertas (la conexión es única, R9); completar el método del fake — `go test ./...` en verde
- [X] T062 [US3] COMMIT `GREEN: listar las rondas reveladas en SQLite, tests en verde`
- [X] T063 [US3] RED: en internal/service/consultar_historial_rondas_test.go, tests de `ConsultarHistorialRondas.Ejecutar(ctx, historiaID)`: US3.8 ronda 1 revelada y 2 abierta → solo la 1, con número, momento y el voto de cada integrante con `Nombre` y `Rol`; sin rondas reveladas → `Rondas` vacío sin error; historia inexistente → `ErrHistoriaNoEncontrada`; consultar deja la ronda actual y sus votos iguales que antes (solo lectura). Stub en internal/service/consultar_historial_rondas.go — confirmar que FALLA
- [X] T064 [US3] COMMIT `RED: agregar tests de consultar el historial de rondas (fallan)`
- [X] T065 [US3] GREEN: implementar `ConsultarHistorialRondas` en internal/service/consultar_historial_rondas.go, con los nombres y roles de `ListarIntegrantes` — `go test ./...` en verde
- [X] T066 [US3] COMMIT `GREEN: implementar consultar el historial de rondas, tests en verde`
- [X] T067 [US3] RED: en internal/http/votacion_handler_test.go, tests de `Historial` (`GET /backlog/{id}/votacion/rondas`): 200 con la forma de `HistorialRondas` de contracts/openapi.yaml, `revelada_en` en RFC 3339 UTC; sin rondas reveladas → `"rondas": []`, no `null`; id inválido → 400; historia inexistente → 404. Agregar los DTO en internal/http/votacion_dto.go, el parámetro `historial *service.ConsultarHistorialRondas` a `NewVotacionHandler` y actualizar a quienes lo llaman. Stub — confirmar que FALLA
- [X] T068 [US3] COMMIT `RED: agregar tests del handler del historial de rondas (fallan)`
- [X] T069 [US3] GREEN: implementar `Historial` en internal/http/votacion_handler.go — `go test ./...` en verde
- [X] T070 [US3] COMMIT `GREEN: implementar el handler del historial de rondas, tests en verde`
- [X] T071 [US3] Cablear `POST /backlog/{id}/votacion/rondas` y `GET /backlog/{id}/votacion/rondas` en cmd/api/main.go
- [X] T072 [US3] COMMIT `Cablear la nueva ronda y el historial en cmd/api`

**Checkpoint**: las tres historias funcionan por la API.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T073 Implementar los pasos Godog en features/steps_consenso.go, con `InitializeScenarioConsenso` registrado en features/features_test.go. Reutilizar el estado del equipo y la historia de HU-02 haciendo que `InitializeScenarioVotacion` devuelva sus `*pasosVotacion`, como `InitializeScenario` devuelve `sc` — `go test ./features/` en verde, incluidos los 20 escenarios `@HU-14`
- [X] T074 COMMIT `Agregar pasos BDD de consenso y estimación acordada`
- [X] T075 `go build ./...`, `go vet ./...`, `gofmt -l .` sin salida y `go test ./...` en verde
- [X] T076 Correr la validación de specs/031-hu-14-consenso-estimacion/quickstart.md contra un servidor con base nueva y anotar el resultado al final de ese archivo
- [X] T077 COMMIT `Registrar la validación de HU-14 en el quickstart`

---

## Phase 7: Correcciones de la revisión del PR

**Purpose**: cerrar los huecos que encontró la revisión con errores introducidos a propósito, y aclarar el mensaje de una apertura simultánea.

- [X] T078 Test de caracterización en internal/repository/sqlite_backlog_test.go: `GuardarEstimacion` no cambia los Story Points de otra historia. Se espera que pase de entrada; verificar que detecta un `WHERE` mal escrito — `go test ./...` en verde
- [X] T079 COMMIT `Agregar test de que guardar la estimación no toca otras historias`
- [X] T080 RED: en internal/http/votacion_handler_test.go, test de que `(VotacionHandler).Registrar(mux)` atiende las 7 rutas de la votación, cada una con su resultado propio. Stub en internal/http/votacion_handler.go que no registra nada — confirmar que FALLA
- [X] T081 COMMIT `RED: agregar test de las rutas de la votación (falla)`
- [X] T082 GREEN: implementar `Registrar` y usarlo desde cmd/api/main.go en lugar de las 7 líneas sueltas, como `demoHandler.Registrar(mux)` — `go test ./...` en verde
- [X] T083 COMMIT `GREEN: registrar las rutas de la votación desde su handler, test en verde`
- [X] T084 RED: en internal/service/abrir_nueva_ronda_test.go, test con un fake que pierde la carrera en `CrearRonda`: `AbrirNuevaRonda` devuelve `ErrRondaYaAbierta`. Agregar el error en internal/domain/errors.go — confirmar que FALLA
- [X] T085 COMMIT `RED: agregar test del error de una apertura simultánea (falla)`
- [X] T086 GREEN: traducir `ErrVotacionYaIniciada` a `ErrRondaYaAbierta` en internal/service/abrir_nueva_ronda.go — `go test ./...` en verde
- [X] T087 COMMIT `GREEN: aclarar el error de una apertura simultánea, test en verde`
- [X] T088 RED: en internal/http/votacion_handler_test.go, test de que `escribirError` responde 409 con el mensaje de `ErrRondaYaAbierta` — confirmar que FALLA
- [X] T089 COMMIT `RED: agregar test del 409 de una apertura simultánea (falla)`
- [X] T090 GREEN: agregar el caso en `escribirError` en internal/http/sprint_handler.go — `go test ./...` en verde
- [X] T091 COMMIT `GREEN: responder 409 a una apertura simultánea, test en verde`
- [X] T092 Actualizar research.md (R8, R10), data-model.md, plan.md y contracts/openapi.yaml, registrando que el último registro simultáneo gana. Correr `go test -race` y repetir las comprobaciones con el servidor real
- [X] T093 COMMIT `docs: registrar las correcciones de la revisión de HU-14`

---

## Phase 8: Correcciones de la revisión independiente

**Purpose**: respetar el contrato con valores que no son enteros, proteger con tests la regla de la ronda actual, y aclarar el mensaje de abrir una ronda en una historia completada.

- [X] T094 RED: en internal/http/votacion_handler_test.go, tests de `estimacion_sp` con `2.5`, `"5"`, `null` y sin el campo → 400 campo `estimacion_sp`; Product Builder con `2.5` → 403. Confirmar que FALLA
- [X] T095 COMMIT `RED: agregar tests de una estimación que no es entera (fallan)`
- [X] T096 GREEN: en internal/http/votacion_dto.go, leer `estimacion_sp` como `json.RawMessage`; lo que no sea un entero llega al service como 0, que lo rechaza después de validar el permiso (research.md R10) — `go test ./...` en verde
- [X] T097 COMMIT `GREEN: rechazar una estimación que no es entera con su campo, tests en verde`
- [X] T098 Tests de caracterización: en internal/service/registrar_estimacion_test.go, ronda 1 revelada y 2 abierta → `ErrRondaNoRevelada`; en internal/http/votacion_handler_test.go, valor distinto con Sprint Activo → 400 campo `estimacion_sp`. Verificar que detectan un chequeo sobre "alguna ronda revelada" — `go test ./...` en verde
- [X] T099 COMMIT `Agregar tests de la ronda actual y de D2b en el handler`
- [X] T100 RED: en internal/service/abrir_nueva_ronda_test.go, `Completada` → `ErrRondaEnHistoriaCompletada`, y en internal/http/votacion_handler_test.go, 400 campo `estado` con su mensaje. Agregar el error en internal/domain/errors.go — confirmar que FALLA
- [X] T101 COMMIT `RED: agregar tests del error de abrir una ronda en una historia completada (fallan)`
- [X] T102 GREEN: devolver `ErrRondaEnHistoriaCompletada` en internal/service/abrir_nueva_ronda.go — `go test ./...` en verde
- [X] T103 COMMIT `GREEN: aclarar el error de abrir una ronda en una historia completada, tests en verde`
- [X] T104 Actualizar research.md (R10), data-model.md y contracts/openapi.yaml; repetir las comprobaciones con el servidor real
- [X] T105 COMMIT `docs: registrar las correcciones de la revisión independiente de HU-14`

---

## Phase 9: Frontend (sub-issue #93)

**Goal**: en la pantalla de Planning Poker (web/src/features/votacion), con la ronda revelada se ven las diferencias con el mínimo y el máximo, y el Scrum Master puede abrir una nueva ronda o registrar la estimación acordada. El historial de rondas no entra en la pantalla (Clarifications, 2026-11-03).

**Independent Test**: con pruebas de componente (Vitest y Testing Library, con el escenario de la spec en el nombre de cada prueba) y con la aplicación levantada: se revela una ronda con votos 3, 5 y 8, se ve "Hay diferencias" con mínimo 3 y máximo 8, se registra el valor 5 y la historia queda con 5 SP; se abre una nueva ronda y se vuelve a votar.

- [X] T106 [US1] [US2] [US3] RED: en web/src/features/votacion/api.test.ts, `abrirNuevaRonda(historiaId, usuarioId)` hace `POST /backlog/{id}/votacion/rondas` con `usuario_id` y devuelve la `Votacion`; `registrarEstimacion(historiaId, usuarioId, valor)` hace `POST /backlog/{id}/votacion/estimacion` con `usuario_id` y `estimacion_sp` y devuelve la historia. Sumar `diferencias` al tipo `Votacion` (web/src/features/votacion/votacion.ts) y stubs en web/src/features/votacion/api.ts — confirmar que FALLA
- [X] T107 COMMIT `RED: agregar tests del cliente de consenso y estimación acordada (fallan)`
- [X] T108 GREEN: implementar las dos funciones en web/src/features/votacion/api.ts — `npm test` en verde
- [X] T109 COMMIT `GREEN: implementar el cliente de consenso y estimación acordada, tests en verde`
- [X] T110 [US1] RED: en web/src/features/votacion/votacion.test.ts, `textoDeDiferencias` ("Hay diferencias: mínimo 3, máximo 8" / "No hay diferencias: todos votaron lo mismo" / vacío si es null), y en web/src/features/votacion/index.test.tsx, escenarios US1.1 a US1.3 de la spec. Stub de la función — confirmar que FALLA
- [X] T111 [US1] COMMIT `RED: agregar pruebas de las diferencias entre los votos (fallan)`
- [X] T112 [US1] GREEN: mostrar las diferencias en web/src/features/votacion/Resultado.tsx — `npm test` en verde
- [X] T113 [US1] COMMIT `GREEN: mostrar las diferencias entre los votos, tests en verde`
- [X] T114 [US2] RED: en web/src/features/votacion/index.test.tsx, escenarios US2.1, US2.2, US2.4 a US2.8: con la ronda revelada el Scrum Master elige una carta de la baraja (por defecto la sugerida, aunque nadie la haya votado) y registra; la historia muestra los nuevos SP y un aviso; un rechazo del servidor muestra su mensaje y la historia mantiene su estimación; un Product Builder no ve el formulario; un doble envío registra una sola vez. Reemplaza la aclaración "se hace en otra historia" — confirmar que FALLA
- [X] T115 [US2] COMMIT `RED: agregar pruebas del registro de la estimación acordada (fallan)`
- [X] T116 [US2] GREEN: formulario en web/src/features/votacion/Consenso.tsx y acción en web/src/features/votacion/useVotacion.ts — `npm test` en verde
- [X] T117 [US2] COMMIT `GREEN: registrar la estimación acordada desde la votación, tests en verde`
- [X] T118 [US3] RED: en web/src/features/votacion/index.test.tsx, escenarios US3.1, US3.3, US3.4 y US3.7: con la ronda revelada, haya o no diferencias, el Scrum Master pulsa "Nueva ronda" y se ve la ronda 2 sin votos y con las cartas; sin revelar o siendo Product Builder no hay botón; un rechazo del servidor muestra su mensaje — confirmar que FALLA
- [X] T119 [US3] COMMIT `RED: agregar pruebas de la nueva ronda (fallan)`
- [X] T120 [US3] GREEN: botón "Nueva ronda" en web/src/features/votacion/Consenso.tsx y acción en web/src/features/votacion/useVotacion.ts — `npm test` en verde
- [X] T121 [US3] COMMIT `GREEN: abrir una nueva ronda desde la votación, tests en verde`
- [X] T122 Estilos en web/src/features/votacion/votacion.css con los tokens de web/src/App.css; `npm run lint`, `npm test` y `npm run build` en web/ sin errores; prueba manual con la aplicación levantada
- [X] T123 COMMIT `Ajustar el estilo del consenso en la votación`

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **User Stories (Phase 3 a 5)** → **Polish (Phase 6)**.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- **US1** no depende de las otras. **US2** no depende de US1 ni de US3. **US3** no depende de US1 ni de US2, salvo T055, cuya respuesta incluye `diferencias`: si US3 va primero, esa aserción se agrega al implementar US1.
- Las historias tocan los mismos archivos (`votacion.go`, `votacion_handler.go`, `votacion_dto.go`, `fakes_test.go`): se hacen de a una. No hay tareas `[P]`, porque el historial RED/GREEN es lineal y cada tarea sigue a su commit.
- T073 va al final porque los pasos llaman a código que debe existir (sub-issue #90).

## Implementation Strategy

1. Setup + Foundational → escenarios y error listos.
2. US1 → validar la consulta (MVP: el equipo ve las diferencias).
3. US2 → la historia ya produce Story Points acordados.
4. US3 → nuevas rondas e historial.
5. Polish → pasos BDD, validación final y PR contra `main` a cargo del equipo.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- Sin migración ni cambios en `MODELO-DE-DATOS.md` (research.md R11). Si aparece una contradicción con el documento, se frena y se avisa.
