---

description: "Tareas de HU-02: Votación de Esfuerzo en Planning Poker (TDD con commits por fase)"
---

# Tasks: Votación de Esfuerzo en Planning Poker (HU-02)

**Input**: Design documents from `/specs/012-hu-02-votacion-planning-poker/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica es RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT).

**Compilación en RED**: todos los commits compilan. En cada RED se agregan, junto al test, solo
las firmas nuevas que el test necesita con un cuerpo vacío (devolver valor cero / `nil`), de modo
que `go build ./...` pase y el test falle por aserción, no por símbolo indefinido.

**Organization**: por user story de la spec (US1 votar en secreto, US2 revelar, US3 iniciar).
US1 y US2 son P1; US3 es P2 pero es el prerrequisito funcional de las otras dos en la API: sus
tests de service y handler siembran la ronda directo en el fake o en la base, así cada historia
se prueba sola.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3

---

## Phase 1: Setup

**Purpose**: dejar los escenarios BDD escritos antes de implementar (flujo de AGENTS.md).

- [ ] T001 Escribir features/votacion_planning_poker.feature en Gherkin español: tag `@HU-02`, Característica con el Como/Quiero/Para de la Descripción de la spec, Antecedentes con un proyecto que tiene a Ana (`scrum_master`), Bruno y Carla (`product_builder`) y una historia "Nueva" del proyecto, con pasos propios de HU-02 para el equipo y la historia (el paso `que existe un proyecto con identificador (\d+)` de features/steps_creacion.go crea el proyecto sin Bruno ni Carla, así que no alcanza; no redefinir pasos existentes ni repetir sus textos), y exactamente 13 escenarios, uno por Acceptance Scenario de spec.md: US1.1 Bruno vota 8 → consultado por Ana figura votó con valor oculto, consultado por Bruno ve 8; US1.2 Bruno cambia 8 por 13 antes de revelar → un solo voto suyo, 13, oculto para los demás; US1.3 voto con 4 rechazado y la ronda no cambia; US1.4 voto de quien no es integrante rechazado; US2.1 votos 3, 5 y 8 → al revelar se ven los tres valores, promedio exacto 5,33 y carta sugerida 5; US2.2 votar o cambiar el voto en una ronda revelada rechazado sin cambios; US2.3 revelar con un integrante sin votar rechazado y los votos siguen ocultos; US2.4 Bruno (Product Builder) intenta revelar → rechazado; US3.1 Ana inicia → ronda 1 sin votos ni revelar; US3.2 iniciar de nuevo una historia con votación → rechazado; US3.3 historia inexistente → "historia no encontrada"; US3.4 Bruno intenta iniciar → rechazado y sin rondas; US3.5 historia "Completada" → rechazado y sin rondas. Comentario en el archivo indicando qué escenario cubre cada categoría (normal: US1.1, US2.1, US3.1; alternativo: US1.2; límite: US2.3; error: el resto). Sin pasos implementados: con el runner estricto, `go test ./features/` FALLA
- [ ] T002 COMMIT `RED: agregar escenarios BDD de la votación de Planning Poker (fallan)`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: esquema, entidad de dominio, puerto y adaptador SQLite que usan las tres historias.

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

### Migración 003

- [ ] T003 RED: en internal/repository/esquema_test.go, sumar a `esquemaEsperado` (y a la lista de tablas) `rondas_votacion` y `votos` según MODELO-DE-DATOS.md §4.7: columnas `rondas_votacion(id INTEGER PK, historia_id INTEGER NOT NULL, numero INTEGER NOT NULL, revelada_en TEXT NULL)` y `votos(ronda_id INTEGER NOT NULL, usuario_id INTEGER NOT NULL, valor INTEGER NOT NULL)`; FK `rondas_votacion.historia_id → historias(id)`, `votos.ronda_id → rondas_votacion(id)`, `votos.usuario_id → usuarios(id)`; `UNIQUE (historia_id, numero)`; PK compuesta `(ronda_id, usuario_id)`. En `TestEsquema_Filas`, casos: ronda válida con `revelada_en` NULL y con valor; segunda ronda con el mismo `(historia_id, numero)` → 2067; ronda de historia inexistente → 787; `numero` NULL → 1299; voto válido; segundo voto del mismo `(ronda_id, usuario_id)` → 1555; voto con `ronda_id` o `usuario_id` inexistente → 787; voto con `valor` NULL, `ronda_id` NULL y `usuario_id` NULL → 1299 (la PK compuesta de SQLite admitiría NULL sin el `NOT NULL`). Confirmar que FALLA
- [ ] T004 COMMIT `RED: agregar al test de esquema las tablas de votación (falla)`
- [ ] T005 GREEN: crear internal/repository/migraciones/003_votacion.sql con el comentario `-- Votación de Planning Poker: MODELO-DE-DATOS.md §4.7.` y el DDL: `CREATE TABLE rondas_votacion (id INTEGER PRIMARY KEY AUTOINCREMENT, historia_id INTEGER NOT NULL REFERENCES historias (id), numero INTEGER NOT NULL, revelada_en TEXT, UNIQUE (historia_id, numero));` y `CREATE TABLE votos (ronda_id INTEGER NOT NULL REFERENCES rondas_votacion (id), usuario_id INTEGER NOT NULL REFERENCES usuarios (id), valor INTEGER NOT NULL, PRIMARY KEY (ronda_id, usuario_id));` (sin `CHECK` de escala: D11, la valida el dominio). No tocar MODELO-DE-DATOS.md: el esquema no cambia — `go test ./...` en verde
- [ ] T006 COMMIT `GREEN: agregar la migración 003 de votación, test de esquema en verde`

### Entidad de dominio

- [ ] T007 Agregar en internal/domain/errors.go los errores de research.md R8: sentinels `ErrVotacionNoIniciada` ("la historia no tiene votación iniciada"), `ErrVotacionYaIniciada` ("la historia ya tiene una votación"), `ErrNoEsIntegrante` ("el usuario no es integrante del proyecto"), `ErrNoEsScrumMaster` ("solo el Scrum Master del proyecto puede hacerlo"); y `ValidationError`: `ErrVotoInvalido{Campo: "valor", Mensaje: "el voto debe ser una carta de la baraja: 1, 2, 3, 5, 8, 13, 21"}`, `ErrHistoriaCompletada{Campo: "estado", Mensaje: "no se puede votar una historia completada"}`, `ErrRondaRevelada{Campo: "estado", Mensaje: "la ronda ya fue revelada"}`, `ErrFaltanVotos{Campo: "votos", Mensaje: "faltan votar integrantes del proyecto"}`. Sin lógica
- [ ] T008 COMMIT `Agregar los errores de la votación de Planning Poker`
- [ ] T009 RED: en internal/domain/votacion_test.go, tests de `RondaVotacion` (data-model.md): (a) `NuevaRonda(historiaID)` → `Numero` 1, sin votos, `ReveladaEn` nil; `historiaID <= 0` → `ValidationError` campo `historia_id`; (b) `Votar` agrega el voto; (c) segundo `Votar` del mismo usuario reemplaza el valor y sigue habiendo un voto suyo (FR-005); (d) `Votar` con 0, 4, 22 → `ErrVotoInvalido` sin cambiar los votos; (e) `Votar` sobre ronda revelada → `ErrRondaRevelada` (aunque el valor sea inválido, la revelada se chequea primero); (f) `Revelar(integrantes, ahora)` con todos votados fija `ReveladaEn = ahora` y `Revelada()` es true; (g) `Revelar` con un integrante sin voto → `ErrFaltanVotos` y `ReveladaEn` sigue nil; (h) segundo `Revelar` → `ErrRondaRevelada`. Stubs en internal/domain/votacion.go con los tipos `RondaVotacion{ID, HistoriaID int64; Numero int; ReveladaEn *time.Time; Votos []Voto}` y `Voto{UsuarioID int64; Valor int}` — confirmar que FALLA
- [ ] T010 COMMIT `RED: agregar tests de la ronda de votación (fallan)`
- [ ] T011 GREEN: implementar `NuevaRonda`, `Votar`, `Revelar` y `Revelada` en internal/domain/votacion.go; la baraja reutiliza `EsEstimacionValida` (internal/domain/valor_negocio.go), no se duplica la escala — `go test ./...` en verde
- [ ] T012 COMMIT `GREEN: implementar la ronda de votación, tests en verde`
- [ ] T013 RED: en internal/domain/votacion_test.go, tests de `CalcularPromedio(valores []int) Promedio` (`Promedio{Exacto float64; CartaSugerida int}`) con la tabla de data-model.md: `[3,5,8]` → 5.33 / 5; `[3,5]` → 4 / 5 (empate, la mayor); `[8,8,8]` → 8 / 8; `[1,21]` → 11 / 13; `[2,3]` → 2.5 / 3 (empate); y `(RondaVotacion).Promedio()` → `(Promedio, false)` si no está revelada, `(CalcularPromedio(valores), true)` si lo está. Regla: "`Exacto` = media aritmética redondeada a 2 decimales; `CartaSugerida` = carta de la escala más cercana a `Exacto`; empate → la mayor". Stubs — confirmar que FALLA
- [ ] T014 COMMIT `RED: agregar tests del promedio sugerido (fallan)`
- [ ] T015 GREEN: implementar `CalcularPromedio` y `Promedio()` en internal/domain/votacion.go (`math.Round(media*100)/100`; recorrer la escala en orden ascendente para resolver el empate hacia la mayor; la escala es un `map`, así que ordenar sus claves) — `go test ./...` en verde
- [ ] T016 COMMIT `GREEN: implementar el promedio sugerido, tests en verde`
- [ ] T017 (solo si hay algo que limpiar) REFACTOR en internal/domain/votacion.go o internal/domain/valor_negocio.go (por ejemplo, exponer la escala ordenada en un solo lugar si `CalcularPromedio` la reconstruye) + COMMIT `REFACTOR: <qué se limpió>`

### Puerto y adaptador SQLite

- [ ] T018 Declarar en internal/repository/votacion_repository.go el puerto `VotacionRepository { CrearRonda(ctx, r domain.RondaVotacion) (domain.RondaVotacion, error); ObtenerRondaActual(ctx, historiaID int64) (domain.RondaVotacion, error); GuardarVoto(ctx, rondaID int64, v domain.Voto) error; MarcarRevelada(ctx, rondaID int64, en time.Time) error }`, documentando: la ronda actual es la de mayor `numero` (§4.7); `ObtenerRondaActual` devuelve `domain.ErrVotacionNoIniciada` si no hay rondas; `CrearRonda` devuelve `domain.ErrVotacionYaIniciada` si choca el `UNIQUE`; `GuardarVoto` y `MarcarRevelada` devuelven `domain.ErrRondaRevelada` si la ronda ya fue revelada (escritura condicional, research.md R6)
- [ ] T019 COMMIT `Agregar el puerto VotacionRepository`
- [ ] T020 RED: en internal/repository/sqlite_votacion_test.go, tests de `SQLiteVotacionRepository` contra `:memory:` con migraciones (sembrar proyecto, usuarios, integrantes e historia con los repos SQLite existentes): (a) `CrearRonda` asigna id y persiste `numero` 1 con `revelada_en` NULL; (b) segunda `CrearRonda` con la misma historia y número → `errors.Is(err, domain.ErrVotacionYaIniciada)`; (c) `CrearRonda` de historia inexistente → `ErrIntegridadReferencial`; (d) `ObtenerRondaActual` sin rondas → `domain.ErrVotacionNoIniciada`; (e) con rondas 1 y 2 sembradas por SQL devuelve la 2 (mayor `numero`) con sus votos ordenados por `usuario_id`; (f) `GuardarVoto` inserta y un segundo `GuardarVoto` del mismo usuario reemplaza el valor (sigue habiendo una fila); (g) `MarcarRevelada` guarda `revelada_en` en RFC 3339 UTC y `ObtenerRondaActual` la devuelve como `*time.Time` igual al instante; (h) `GuardarVoto` y `MarcarRevelada` sobre una ronda revelada → `domain.ErrRondaRevelada` y la base no cambia; (i) `GuardarVoto` con usuario inexistente → `ErrIntegridadReferencial`. Stubs en internal/repository/sqlite_votacion.go — confirmar que FALLA
- [ ] T021 COMMIT `RED: agregar tests del repositorio SQLite de votación (fallan)`
- [ ] T022 GREEN: implementar internal/repository/sqlite_votacion.go: `CrearRonda` (`INSERT`; código 2067 → `domain.ErrVotacionYaIniciada`; 787 → `ErrIntegridadReferencial` vía `errorDeEscritura`); `ObtenerRondaActual` (`SELECT ... FROM rondas_votacion WHERE historia_id = ? ORDER BY numero DESC LIMIT 1` + `SELECT usuario_id, valor FROM votos WHERE ronda_id = ? ORDER BY usuario_id`; `sql.ErrNoRows` → `domain.ErrVotacionNoIniciada`); `GuardarVoto` (`INSERT INTO votos (ronda_id, usuario_id, valor) SELECT ?, ?, ? WHERE EXISTS (SELECT 1 FROM rondas_votacion WHERE id = ? AND revelada_en IS NULL) ON CONFLICT (ronda_id, usuario_id) DO UPDATE SET valor = excluded.valor`; 0 filas afectadas → `domain.ErrRondaRevelada`); `MarcarRevelada` (`UPDATE rondas_votacion SET revelada_en = ? WHERE id = ? AND revelada_en IS NULL`, instante en `time.RFC3339` UTC; 0 filas → `domain.ErrRondaRevelada`) — `go test ./...` en verde
- [ ] T023 COMMIT `GREEN: implementar el repositorio SQLite de votación, tests en verde`

### Soporte de los services

- [ ] T024 Agregar en internal/service/fakes_test.go un `fakeVotacionRepository` (ronda en memoria por historia, con `err` configurable) que simule las escrituras condicionales de R6: `GuardarVoto` y `MarcarRevelada` devuelven `domain.ErrRondaRevelada` si la ronda ya está revelada, y `CrearRonda` devuelve `domain.ErrVotacionYaIniciada` si la historia ya tiene ronda. Sin lógica de negocio propia
- [ ] T025 COMMIT `Agregar el fake de VotacionRepository para los tests de service`

**Checkpoint**: esquema, dominio y persistencia listos.

---

## Phase 3: User Story 1 - Votar en secreto (Priority: P1) 🎯 MVP

**Goal**: un integrante vota una carta y la consulta muestra quién votó, con el valor oculto para todos salvo para quien consulta.

**Independent Test**: con una ronda 1 abierta sembrada en el repositorio, Bruno vota 8; la consulta de Ana muestra a Bruno con `ya_voto: true, valor: null` y la de Bruno muestra `valor: 8`.

### Casos de uso Votar y ConsultarVotacion

- [ ] T026 [US1] RED: en internal/service/votar_test.go, tests con fakes (`fakeHistoriaRepository`, `fakeProjectRepository` con `listaMiembros`, `fakeVotacionRepository`) de `Votar.Ejecutar(ctx, VotarInput{HistoriaID, UsuarioID int64; Valor int}) (domain.RondaVotacion, error)`: (a) historia inexistente → `domain.ErrHistoriaNoEncontrada`; (b) usuario que no es integrante del proyecto de la historia → `domain.ErrNoEsIntegrante` sin guardar; (c) sin ronda → `domain.ErrVotacionNoIniciada`; (d) voto válido se guarda en la ronda actual; (e) cambio de voto reemplaza; (f) valor 4 → `domain.ErrVotoInvalido` sin guardar; (g) ronda revelada → `domain.ErrRondaRevelada`; (h) un Scrum Master también puede votar. Incluir en internal/service/integrantes.go el stub de `rolEnProyecto(ctx, proyectos repository.ProjectRepository, proyectoID, usuarioID int64) (domain.Role, bool, error)`. Stubs de `NewVotar(historias repository.HistoriaRepository, proyectos repository.ProjectRepository, votaciones repository.VotacionRepository)` en internal/service/votar.go — confirmar que FALLA
- [ ] T027 [US1] COMMIT `RED: agregar tests de Votar (fallan)`
- [ ] T028 [US1] GREEN: implementar `rolEnProyecto` en internal/service/integrantes.go (busca al usuario en `ListarIntegrantes` del proyecto) y `Votar.Ejecutar` en internal/service/votar.go: historia → integrante (403) → ronda actual → `ronda.Votar` (reglas del dominio) → `GuardarVoto` — `go test ./...` en verde
- [ ] T029 [US1] COMMIT `GREEN: implementar Votar, tests en verde`
- [ ] T030 [US1] RED: en internal/service/consultar_votacion_test.go, tests de `ConsultarVotacion.Ejecutar(ctx, historiaID, consultanteID int64) (VistaVotacion, error)` (data-model.md, "Vista de consulta"): (a) historia inexistente → `ErrHistoriaNoEncontrada`; (b) sin ronda → `ErrVotacionNoIniciada`; (c) ronda abierta: lista a todos los integrantes ordenados por `usuario_id` con `YaVoto`; `Valor` es nil para los demás y el propio para quien consulta (FR-007); (d) `consultanteID` 0 o de alguien ajeno al proyecto → todos los `Valor` nil; (e) integrante sumado con la ronda abierta aparece con `YaVoto` false; (f) ronda revelada → todos los valores y `Promedio` no nil; abierta → `Promedio` nil; (g) no escribe en ningún repositorio. Stubs de `VistaVotacion{HistoriaID int64; Numero int; ReveladaEn *time.Time; Integrantes []VotoDeIntegrante; Promedio *domain.Promedio}`, `VotoDeIntegrante{UsuarioID int64; Nombre string; Rol domain.Role; YaVoto bool; Valor *int}` y `NewConsultarVotacion(historias, proyectos, votaciones)` en internal/service/consultar_votacion.go — confirmar que FALLA
- [ ] T031 [US1] COMMIT `RED: agregar tests de ConsultarVotacion (fallan)`
- [ ] T032 [US1] GREEN: implementar `ConsultarVotacion.Ejecutar` en internal/service/consultar_votacion.go — `go test ./...` en verde
- [ ] T033 [US1] COMMIT `GREEN: implementar ConsultarVotacion con el voto propio visible, tests en verde`

### Handlers POST /backlog/{id}/votacion/votos y GET /backlog/{id}/votacion

- [ ] T034 [US1] RED: en internal/http/votacion_handler_test.go, tests con `httptest` de `VotacionHandler.Votar` y `VotacionHandler.Consultar` según contracts/openapi.yaml: `POST /backlog/{id}/votacion/votos` → 200 con la vista de quien votó (`historia_id, ronda, revelada_en, integrantes[{usuario_id, nombre, rol, ya_voto, valor}], promedio`), 400 id no numérico, 404 historia inexistente aunque el JSON sea inválido (la existencia domina, FR-011), 400 JSON inválido, 400 `campo: usuario_id` si falta o es ≤ 0, 403 no integrante, 404 sin votación, 400 `campo: valor`, 400 `campo: estado` con ronda revelada; `GET /backlog/{id}/votacion?usuario_id=N` → 200 con `valor: null` para los demás, 400 si `usuario_id` no es numérico, 404 `historia no encontrada` si la historia no existe (FR-011), 404 sin votación. Crear internal/http/votacion_dto.go con `votoRequest{UsuarioID int64; Valor int}` y `votacionResponse`; constructor `NewVotacionHandler(obtenerHistoria *service.ObtenerHistoria, votar *service.Votar, consultar *service.ConsultarVotacion)`. Stubs vacíos — confirmar que FALLA
- [ ] T035 [US1] COMMIT `RED: agregar tests de los handlers de voto y consulta de la votación (fallan)`
- [ ] T036 [US1] GREEN: implementar `Votar` y `Consultar` en internal/http/votacion_handler.go (validar la historia con `obtenerHistoria` antes de decodificar el cuerpo, como `SprintHandler.Iniciar`; `ahora` no se usa acá) y ampliar `escribirError` en internal/http/sprint_handler.go solo con lo que ejercitan los tests de T034: `ErrNoEsIntegrante` → 403 y `ErrVotacionNoIniciada` → 404 (`ErrNoEsScrumMaster` entra en T046 y `ErrVotacionYaIniciada` en T054, cada uno con su RED) — `go test ./...` en verde
- [ ] T037 [US1] COMMIT `GREEN: implementar los handlers de voto y consulta de la votación, tests en verde`
- [ ] T038 [US1] Cablear en cmd/api/main.go `repository.NewSQLiteVotacionRepository(db)`, los services `NewVotar` y `NewConsultarVotacion`, `NewVotacionHandler` y las rutas `POST /backlog/{id}/votacion/votos` y `GET /backlog/{id}/votacion`
- [ ] T039 [US1] COMMIT `Cablear el voto y la consulta de la votación en cmd/api`

**Checkpoint**: US1 funciona con una ronda sembrada (quickstart pasos 4–8 tras iniciar por SQL o tras US3).

---

## Phase 4: User Story 2 - Revelar los votos (Priority: P1)

**Goal**: el Scrum Master revela cuando votaron todos; se ven los valores y el promedio.

**Independent Test**: con una ronda sembrada donde votaron los tres integrantes (3, 5, 8), Ana revela y la respuesta trae los tres valores y `promedio: {exacto: 5.33, carta_sugerida: 5}`.

- [ ] T040 [US2] RED: en internal/service/revelar_votacion_test.go, tests de `RevelarVotacion.Ejecutar(ctx, historiaID, usuarioID int64, ahora time.Time) (domain.RondaVotacion, error)`: (a) historia inexistente → `ErrHistoriaNoEncontrada`; (b) usuario Product Builder o ajeno → `ErrNoEsScrumMaster` sin revelar; (c) sin ronda → `ErrVotacionNoIniciada`; (d) falta un voto → `ErrFaltanVotos` y la ronda sigue sin revelar; (e) todos votaron → `MarcarRevelada` con `ahora` y la ronda devuelta revelada con promedio; (f) ya revelada → `ErrRondaRevelada`; (g) con dos Scrum Master en el proyecto, cualquiera revela (Assumptions). Stub de `NewRevelarVotacion(historias, proyectos, votaciones)` en internal/service/revelar_votacion.go — confirmar que FALLA
- [ ] T041 [US2] COMMIT `RED: agregar tests de RevelarVotacion (fallan)`
- [ ] T042 [US2] GREEN: implementar `RevelarVotacion.Ejecutar`: historia → Scrum Master (`rolEnProyecto`) → ronda actual → `ronda.Revelar(ids de ListarIntegrantes, ahora)` → `MarcarRevelada` — `go test ./...` en verde
- [ ] T043 [US2] COMMIT `GREEN: implementar RevelarVotacion, tests en verde`
- [ ] T044 [US2] RED: en internal/http/votacion_handler_test.go, tests de `VotacionHandler.Revelar` (`POST /backlog/{id}/votacion/revelar` con `{"usuario_id": N}`): 200 con todos los valores, `revelada_en` RFC 3339 y `promedio`; 404 historia inexistente aunque el JSON sea inválido; 400 `campo: usuario_id`; 403 no Scrum Master; 404 sin votación; 400 `campo: votos` si falta votar; 400 `campo: estado` si ya estaba revelada. Agregar `revelar *service.RevelarVotacion` a `NewVotacionHandler` (adaptar el helper de prueba y cmd/api/main.go en este mismo commit; la ruta se registra en T046). Stub vacío — confirmar que FALLA
- [ ] T045 [US2] COMMIT `RED: agregar tests del handler POST /backlog/{id}/votacion/revelar (fallan)`
- [ ] T046 [US2] GREEN: implementar `Revelar` en internal/http/votacion_handler.go (pasa `time.Now()` al service, research.md R7), mapear `ErrNoEsScrumMaster` → 403 en `escribirError` (internal/http/sprint_handler.go) y registrar `POST /backlog/{id}/votacion/revelar` en cmd/api/main.go — `go test ./...` en verde
- [ ] T047 [US2] COMMIT `GREEN: implementar POST /backlog/{id}/votacion/revelar y cablear la ruta, tests en verde`

**Checkpoint**: US1 + US2 funcionan juntas sobre una ronda sembrada.

---

## Phase 5: User Story 3 - Iniciar la votación de una historia (Priority: P2)

**Goal**: el Scrum Master abre la ronda 1 de una historia que no esté Completada ni tenga votación.

**Independent Test**: Ana inicia la votación de una historia "Nueva" → 201 con `ronda: 1`, los integrantes con `ya_voto: false` y `revelada_en: null`; un segundo inicio → 409.

- [ ] T048 [US3] RED: en internal/service/iniciar_votacion_test.go, tests de `IniciarVotacion.Ejecutar(ctx, historiaID, usuarioID int64) (domain.RondaVotacion, error)`: (a) historia inexistente → `ErrHistoriaNoEncontrada`; (b) Product Builder o ajeno → `ErrNoEsScrumMaster` sin crear ronda; (c) historia `Completada` → `ErrHistoriaCompletada` sin crear ronda (FR-017); (d) historia "Nueva" y "En progreso" (con Sprint) → crea la ronda 1 sin votos; (e) historia con ronda → `ErrVotacionYaIniciada` (chequeo previo con `ObtenerRondaActual`) y también si `CrearRonda` lo devuelve (carrera); (f) no modifica la historia: `Actualizar` del `fakeHistoriaRepository` no se llama (FR-013). Stub de `NewIniciarVotacion(historias, proyectos, votaciones)` en internal/service/iniciar_votacion.go — confirmar que FALLA
- [ ] T049 [US3] COMMIT `RED: agregar tests de IniciarVotacion (fallan)`
- [ ] T050 [US3] GREEN: implementar `IniciarVotacion.Ejecutar`: historia → Scrum Master → no `Completada` → sin ronda previa → `domain.NuevaRonda` → `CrearRonda` — `go test ./...` en verde
- [ ] T051 [US3] COMMIT `GREEN: implementar IniciarVotacion, tests en verde`
- [ ] T052 [US3] RED: en internal/http/votacion_handler_test.go, tests de `VotacionHandler.Iniciar` (`POST /backlog/{id}/votacion` con `{"usuario_id": N}`): 201 con la vista de la ronda 1 (integrantes con `ya_voto: false`, `promedio: null`); 404 historia inexistente con cuerpo `{}` (no 400); 400 `campo: usuario_id`; 403 no Scrum Master; 400 `campo: estado` con historia Completada; 409 si ya tiene votación. La respuesta se arma con `ConsultarVotacion` para quien inició (una sola forma de la vista). Agregar `iniciar *service.IniciarVotacion` a `NewVotacionHandler` (adaptar helper de prueba y cmd/api/main.go en este commit). Stub vacío — confirmar que FALLA
- [ ] T053 [US3] COMMIT `RED: agregar tests del handler POST /backlog/{id}/votacion (fallan)`
- [ ] T054 [US3] GREEN: implementar `Iniciar` en internal/http/votacion_handler.go, mapear `ErrVotacionYaIniciada` → 409 en `escribirError` (internal/http/sprint_handler.go) y registrar `POST /backlog/{id}/votacion` en cmd/api/main.go — `go test ./...` en verde
- [ ] T055 [US3] COMMIT `GREEN: implementar POST /backlog/{id}/votacion y cablear la ruta, tests en verde`
- [ ] T056 [US3] (solo si hay algo que limpiar) REFACTOR de internal/http/votacion_handler.go (por ejemplo, un helper compartido para "id de ruta → historia existe → decodificar `usuario_id`" usado por `Iniciar`, `Votar` y `Revelar`) + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: las tres historias funcionan de punta a punta por la API.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T057 Implementar features/steps_votacion.go con los pasos de los 13 escenarios de T001 sobre SQLite `:memory:` (repos SQLite reales y los 4 services, como features/steps_backlog.go; reusar `scenarioContext` de features/steps_creacion.go), y registrar `InitializeScenarioVotacion(ctx, sc)` en features/features_test.go con el comentario `// HU-02`. Para la historia "Completada" (US3.5) sembrar el estado por SQL (hoy no hay endpoint, HU-11), cumpliendo los `CHECK` de §4.4 (Sprint asignado y `fecha_finalizacion`) — `go test ./features/` en verde
- [ ] T058 COMMIT `GREEN: implementar los pasos BDD de la votación de Planning Poker, escenarios en verde`
- [ ] T059 `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test ./...` sin errores
- [ ] T060 Validación manual con specs/012-hu-02-votacion-planning-poker/quickstart.md (pasos 1–14) y registrar el resultado en ese archivo
- [ ] T061 Actualizar **Status** de specs/012-hu-02-votacion-planning-poker/spec.md a "Implementada, pendiente de revisión por PR" y marcar [X] las tareas hechas en este archivo
- [ ] T062 COMMIT `docs: registrar la validación de HU-02 y marcar las tareas completadas`

---

## Dependencies & Execution Order

- **Setup (T001–T002)** → **Foundational (T003–T025)** → **US1 (T026–T039)** → **US2 (T040–T047)** → **US3 (T048–T056)** → **Polish (T057–T062)**.
- Dentro de Foundational: migración (T003–T006) antes del adaptador SQLite (T020–T023); errores (T007) antes de la entidad (T009) y del puerto (T018); el fake (T024) después del puerto.
- `rolEnProyecto` nace en US1 (T028) y lo reutilizan US2 y US3.
- Los handlers de US2 y US3 amplían `NewVotacionHandler`: cada RED adapta el helper de prueba y cmd/api/main.go para que todo compile.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.

## Parallel Opportunities

- T003 (test de esquema) y T009/T013 (dominio) tocan archivos distintos y no dependen entre sí: pueden escribirse en paralelo, pero cada uno con su propio commit en el orden de la lista.
- En US1, los tests de `Votar` (T026) y `ConsultarVotacion` (T030) son archivos distintos; se escriben en paralelo si hay dos personas, respetando un commit por fase.
- El resto es secuencial: comparte `votacion_handler.go`, `sprint_handler.go`, `fakes_test.go` o `cmd/api/main.go`.

## Implementation Strategy

1. Setup + Foundational → esquema, dominio y persistencia listos.
2. US1 (votar en secreto) → MVP: se prueba con una ronda sembrada.
3. US2 (revelar) → completa el valor de la historia (votos + promedio).
4. US3 (iniciar) → la API queda usable de punta a punta sin sembrar nada.
5. Polish: BDD en verde, quickstart, PR con `Closes #2` cuando estén cerradas las sub-issues de spec, BDD, backend y tests.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- Fuera de alcance (HU-14): detectar diferencias, abrir la ronda 2 y registrar la estimación acordada. Frontend: sub-issue propia.
