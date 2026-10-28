---

description: "Tareas de HU-07: Imputación de Horas Trabajadas (TDD con commits por fase)"
---

# Tasks: Imputación de Horas Trabajadas (HU-07)

**Input**: Design documents from `/specs/019-hu-07-imputacion-horas/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica es RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT).

**Compilación en RED**: todos los commits compilan. En cada RED se agregan, junto al test, solo
las firmas nuevas que el test necesita con un cuerpo vacío (devolver valor cero / `nil`), de modo
que `go build ./...` pase y el test falle por aserción, no por símbolo indefinido. Si el
stub cambia una firma que ya se usa (`NewHorasHandler` en T044), el mismo RED actualiza la
llamada en cmd/api/main.go; es la única excepción a "solo archivos de test y firmas".

**Organization**: por user story de la spec (US1 imputar, US2 rechazar cargas inválidas, US3 consultar el esfuerzo real).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3

---

## Phase 1: Setup

**Purpose**: dejar los escenarios BDD escritos antes de implementar (flujo de AGENTS.md).

- [ ] T001 Escribir features/imputacion_horas.feature en Gherkin español: tag `@HU-07`, Característica con el Como/Quiero/Para de la spec ("Como integrante del equipo / Quiero registrar las horas reales que invertí trabajando en una Historia de Usuario / Para poder proveer datos al sistema que permitan comparar el esfuerzo estimado contra el esfuerzo real"), Antecedentes que arrancan con el paso existente `Dado que existe un proyecto con identificador 1` (features/steps_creacion.go, no se redefine) y siguen con un integrante, un Sprint "Activo" y una historia asignada a ese Sprint, y exactamente 12 escenarios, uno por Acceptance Scenario de spec.md: US1.1 imputar 4 horas con la fecha de hoy registra integrante, fecha, actividad, duración y Sprint, y el acumulado sube 4 horas; US1.2 cargas de dos integrantes se suman en el acumulado; US1.3 las horas siguen atribuidas al Sprint en que se cargaron después de que la historia cambia de Sprint (cerrar el Sprint y planificarla en otro); US2.1 horas negativas (-2) o cero rechazadas con error sobre las horas; US2.2 horas no numéricas rechazadas con error sobre las horas; US2.3 actividad vacía o solo espacios rechazada con error sobre la actividad; US2.4 persona que no es integrante del proyecto rechazada con error sobre el integrante; US2.5 historia fuera del Sprint "Activo" (en el Product Backlog, en un Sprint "Pendiente" o en uno "Finalizado", como Esquema del escenario) rechazada con "solo se imputan horas a historias del Sprint activo"; US2.6 historia inexistente → "historia no encontrada"; US3.1 consultar el esfuerzo de una historia con varias cargas muestra el total y el detalle; US3.2 historia sin cargas → total 0 y sin detalle; US3.3 consultar una historia inexistente → "historia no encontrada". En los rechazos, verificar también que el acumulado no cambia. Indicar en un comentario del archivo qué escenario cubre cada categoría: normal (US1.1, US3.1), alternativo (US1.2, US1.3, US3.2), límite (US2.1 con cero; los bordes de 24 h y de las fechas no son criterios de aceptación y se cubren con tests unitarios de dominio, T013 y T027), error (US2.x, US3.3). No inventar escenarios fuera de la spec. Sin pasos implementados: con el runner estricto, `go test ./features/` FALLA
- [ ] T002 COMMIT `RED: agregar escenarios BDD de imputación de horas (fallan)`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: esquema, entidad, puertos y persistencia que usan las tres historias.

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

- [ ] T003 RED: en internal/repository/esquema_test.go, (a) agregar `registros_horas` al mapa de estructura esperada de `TestEsquema_Estructura` con las columnas de MODELO-DE-DATOS.md §4.8 (`id` INTEGER pk; `proyecto_id`, `historia_id`, `usuario_id`, `sprint_id` INTEGER NOT NULL; `fecha`, `actividad` TEXT NOT NULL; `minutos` INTEGER NOT NULL; `creado_en` TEXT NOT NULL) y las FK `historia_id,proyecto_id -> historias(id,proyecto_id)`, `sprint_id,proyecto_id -> sprints(id,proyecto_id)`, `proyecto_id,usuario_id -> integrantes(proyecto_id,usuario_id)`, listadas en el orden en que las devuelve `pragma_foreign_key_list` (`fksDe` ordena por `id`, que es el inverso de la declaración: ver `historias_arrastradas`); actualizar la aserción final `versionesAplicadas(...) == []int{1, 2}` a `[]int{1, 2, 4}` (internal/repository/esquema_test.go, línea 178); (b) agregar a `TestEsquema_Filas` una constante `insertarHoras` y los casos de §9.6 sobre `datosBase` (historia 101 en el Sprint 10 del proyecto 1, usuario 1 integrante del 1): `minutos` 60 → aceptada; 0, -30, NULL, 1.5 y 'abc' → `violaCheck` (NULL → `violaNotNull`); `fecha` '2026-10-9', 'basura' y '2026-02-30' → `violaCheck`, '2026-10-09' → aceptada; historia de otro proyecto (el caso inserta su propia historia del proyecto 2 en la misma sentencia, porque `datosBase` no tiene ninguna), Sprint de otro proyecto (21), usuario que no es integrante (agregar el usuario 2 en el caso) → `violaFK`. Confirmar que FALLA (la tabla no existe)
- [ ] T004 COMMIT `RED: agregar tests de esquema de registros_horas (fallan)`
- [ ] T005 GREEN: crear internal/repository/migraciones/004_registros_horas.sql con el DDL de MODELO-DE-DATOS.md §4.8 tal cual: `id INTEGER PRIMARY KEY AUTOINCREMENT`, columnas `NOT NULL`, `CHECK (typeof(minutos) = 'integer' AND minutos > 0)`, `CHECK (date(fecha) IS fecha)`, las tres FK compuestas, sin `ON DELETE` (D9); comentario de cabecera `-- Registros de horas: MODELO-DE-DATOS.md §4.8.` como la 002 — `go test ./...` en verde
- [ ] T006 COMMIT `GREEN: agregar la migración 004 de registros_horas, tests en verde`
- [ ] T007 Agregar en internal/domain/errors.go `ErrHistoriaFueraDelSprintActivo = ValidationError{Campo: "historia_id", Mensaje: "solo se imputan horas a historias del Sprint activo"}` y `ErrIntegranteAjeno = ValidationError{Campo: "integrante_id", Mensaje: "el integrante no pertenece al proyecto de la historia"}`; crear internal/domain/registro_horas.go con los tipos `RegistroHoras` (`ID, ProyectoID, HistoriaID, UsuarioID, SprintID int64; Fecha time.Time; Actividad string; Minutos int; CreadoEn time.Time`) y `RegistroHorasDetalle` (`RegistroHoras` embebido + `NombreIntegrante string`); crear internal/repository/registro_horas_repository.go con los puertos `RegistroHorasRepository { Registrar(ctx, r domain.RegistroHoras) (domain.RegistroHoras, error); ListarPorHistoria(ctx, historiaID int64) ([]domain.RegistroHorasDetalle, error) }` e `IntegrantesRepository { ListarIntegrantes(ctx, proyectoID int64) ([]domain.Member, error) }`, documentando que `Registrar` devuelve `domain.ErrHistoriaFueraDelSprintActivo` si la historia ya no está en ese Sprint Activo y que `ListarPorHistoria` ordena por `fecha, id` y devuelve slice vacío no nil. Sin lógica
- [ ] T008 COMMIT `Agregar la entidad RegistroHoras, sus errores y los puertos de horas`
- [ ] T009 RED: en internal/repository/sqlite_registro_horas_test.go, tests de `SQLiteRegistroHorasRepository` contra `:memory:` con migraciones (sembrar proyecto, usuarios, integrantes, Sprint Activo y historia con `UPDATE historias SET sprint_id = ?` directo en la base): (a) `Registrar` persiste todos los campos (`fecha` como `AAAA-MM-DD`, `creado_en` RFC 3339 UTC) y devuelve el ID asignado; (b) `Registrar` con la historia ya fuera de ese Sprint (`sprint_id` NULL) → `errors.Is(err, domain.ErrHistoriaFueraDelSprintActivo)` y no inserta ninguna fila; (c) `Registrar` con el Sprint ya `Finalizado` → mismo error, sin filas; (d) `Registrar` con un usuario que no es integrante → `errors.Is(err, ErrIntegridadReferencial)`; (e) `ListarPorHistoria` devuelve solo las cargas de esa historia, con `NombreIntegrante`, ordenadas por `fecha` y, a igual fecha, por `id`; (f) historia sin cargas → slice vacío no nil; (g) dos `Registrar` seguidos sobre la misma historia y el mismo integrante quedan las dos cargas, con IDs distintos (edge case de cargas simultáneas). Stubs vacíos en internal/repository/sqlite_registro_horas.go — confirmar que FALLA
- [ ] T010 COMMIT `RED: agregar tests del repositorio SQLite de horas (fallan)`
- [ ] T011 GREEN: implementar en internal/repository/sqlite_registro_horas.go `NewSQLiteRegistroHorasRepository(db)`, `Registrar` con el INSERT condicional de research.md R5 (`INSERT INTO registros_horas (...) SELECT ?, ... WHERE EXISTS (SELECT 1 FROM historias h JOIN sprints s ON s.id = h.sprint_id WHERE h.id = ? AND h.sprint_id = ? AND s.estado = 'Activo')`; 0 filas afectadas → `domain.ErrHistoriaFueraDelSprintActivo`; errores de escritura con `errorDeEscritura`; fecha con `fechaASQL`; `creado_en` con `time.Now().UTC().Format(time.RFC3339)`, como internal/repository/sqlite_project.go) y `ListarPorHistoria` (`JOIN usuarios`, `ORDER BY fecha, id`, slice con `make`) — `go test ./...` en verde
- [ ] T012 COMMIT `GREEN: implementar el repositorio SQLite de horas, tests en verde`

**Checkpoint**: esquema, entidad, puertos y persistencia listos (incluida la carrera con el cierre del Sprint).

---

## Phase 3: User Story 1 - Imputar horas trabajadas en una historia (Priority: P1) 🎯 MVP

**Goal**: `POST /backlog/{id}/horas` registra una carga válida atribuida al Sprint activo de la historia.

**Independent Test**: sembrar un proyecto con un integrante, un Sprint "Activo" y una historia asignada; imputar 4 horas y verificar `201` con `minutos: 240`, `horas: 4` y el `sprint_id` del Sprint activo.

- [ ] T013 [US1] RED: en internal/domain/registro_horas_test.go, casos válidos: `MinutosDesdeHoras("4")` = 240, `"1.5"` = 90, `"0.25"` = 15, `"24"` = 1440 (borde del tope, aceptado), `"0.5"` = 30; `NewRegistroHoras(historia, sprint, usuarioID, "4", "  Implementar  ", fecha, ahora)` devuelve `ProyectoID` y `HistoriaID` de la historia, `SprintID` del Sprint, `Minutos` 240, `Actividad` sin espacios externos y `Fecha` a las 00:00 UTC; acepta una fecha igual a hoy y una igual al inicio del Sprint (bordes inclusivos), con `ahora` fijo (no `time.Now()`). Stubs vacíos de `MinutosDesdeHoras` y `NewRegistroHoras` en internal/domain/registro_horas.go — confirmar que FALLA
- [ ] T014 [US1] COMMIT `RED: agregar tests de RegistroHoras válido y conversión de horas a minutos (falla)`
- [ ] T015 [US1] GREEN: implementar en internal/domain/registro_horas.go `MinutosDesdeHoras(texto string) (int, error)` con `math/big.Rat` (solo `SetString` y multiplicar por 60, lo mínimo para los casos válidos de T013; el formato, el entero, el cero y el tope llegan con sus tests en T027/T029) y `NewRegistroHoras(historia BacklogItem, sprint Sprint, usuarioID int64, horas, actividad, fecha string, ahora time.Time) (RegistroHoras, error)` (copia los campos, recorta la actividad y parsea la fecha con `time.DateOnly`; `ahora` todavía no se usa: el rango de fechas llega en T029) — `go test ./...` en verde
- [ ] T016 [US1] COMMIT `GREEN: implementar RegistroHoras y la conversión exacta de horas a minutos, test en verde`
- [ ] T017 [US1] RED: en internal/service/imputar_horas_test.go, con fakes nuevos en internal/service/fakes_test.go (`fakeRegistroHorasRepository` que guarda las cargas y asigna IDs; `fakeIntegrantesRepository` con una lista de `domain.Member`; reutilizar `fakeHistoriaRepository` y `repository.NewSprintRepositoryEnMemoria()`, como internal/service/cerrar_sprint_test.go): `ImputarHoras.Ejecutar(ctx, ImputarHorasInput{HistoriaID, IntegranteID, Fecha, Actividad, Horas string}, ahora)` con historia en un Sprint Activo e integrante del proyecto registra una carga con el `SprintID` de la historia y devuelve un `domain.RegistroHorasDetalle` con el ID asignado y `NombreIntegrante` igual al `Nombre` del `domain.Member` del fake (la respuesta 201 del contrato lo muestra); dos cargas sobre la misma historia quedan las dos registradas. Stub vacío en internal/service/imputar_horas.go — confirmar que FALLA
- [ ] T018 [US1] COMMIT `RED: agregar tests de ImputarHoras con una carga válida (falla)`
- [ ] T019 [US1] GREEN: implementar `ImputarHoras` en internal/service/imputar_horas.go (`NewImputarHoras(historias repository.HistoriaRepository, sprints repository.SprintRepository, integrantes repository.IntegrantesRepository, horas repository.RegistroHorasRepository)`): obtener historia → obtener su Sprint → buscar el integrante en `ListarIntegrantes(historia.ProyectoID)` para tomar su nombre (sin rechazar todavía al ajeno: eso es T031/T033) → `NewRegistroHoras` → `Registrar` → devolver `RegistroHorasDetalle` — `go test ./...` en verde
- [ ] T020 [US1] COMMIT `GREEN: implementar ImputarHoras, test en verde`
- [ ] T021 [US1] RED: en internal/http/horas_handler_test.go, `POST /backlog/{id}/horas` con `{"integrante_id": 2, "fecha": "<hoy>", "actividad": "Implementar", "horas": 4}` → `201` con el cuerpo `RegistroHoras` del contrato (`id`, `historia_id`, `sprint_id`, `integrante_id`, `integrante`, `fecha`, `actividad`, `minutos: 240`, `horas: 4`); `"horas": 1.5` → `minutos: 90`. Como el handler pasa `time.Now()`, el test arma `<hoy>` y las fechas del Sprint relativas al día actual. Los services se construyen con stubs propios del paquete `http` (los fakes de internal/service/fakes_test.go no son visibles desde otro paquete), como `stubBacklogRepository` en internal/http/backlog_handler_test.go, y `repository.NewSprintRepositoryEnMemoria()`. Stubs de `HorasHandler`, `NewHorasHandler(imputar *service.ImputarHoras, obtenerHistoria *service.ObtenerHistoria)` (la segunda dependencia se usa recién en T037, pero va desde ahora para que la firma no cambie) e `Imputar` en internal/http/horas_handler.go y DTOs en internal/http/horas_dto.go — confirmar que FALLA
- [ ] T022 [US1] COMMIT `RED: agregar tests del handler POST /backlog/{id}/horas (falla)`
- [ ] T023 [US1] GREEN: implementar `Imputar` en internal/http/horas_handler.go (id con `parsearID`; `horas` como `json.RawMessage` cuyo literal numérico se pasa como texto, solo el caso número de T021; los strings y `null` de research.md R3 llegan con sus tests en T035/T037; `fecha` y `actividad` como string; pasar `time.Now()` al service; responder `201`) y `imputarHorasRequest` / `registroHorasResponse` en internal/http/horas_dto.go (`horas = float64(minutos) / 60`) — tests en verde
- [ ] T024 [US1] COMMIT `GREEN: implementar handler POST /backlog/{id}/horas, test en verde`
- [ ] T025 [US1] Cablear en cmd/api/main.go `repository.NewSQLiteRegistroHorasRepository(db)`, `service.NewImputarHoras(backlogRepo, sprintRepo, proyectos, horasRepo)`, `apihttp.NewHorasHandler(imputarHoras, service.NewObtenerHistoria(backlogRepo))` y `mux.HandleFunc("POST /backlog/{id}/horas", horasHandler.Imputar)`
- [ ] T026 [US1] COMMIT `Cablear POST /backlog/{id}/horas en cmd/api`

**Checkpoint**: se pueden imputar horas válidas de punta a punta.

---

## Phase 4: User Story 2 - Rechazar cargas inválidas (Priority: P1)

**Goal**: toda carga inválida se rechaza con un error que identifica el campo, en el orden de la spec, sin registrar nada.

**Independent Test**: intentar imputar con cada dato inválido por separado y verificar `400` con el `campo` esperado (o `404` si la historia no existe) y que el repositorio no recibió ninguna carga.

- [ ] T027 [US2] RED: en internal/domain/registro_horas_test.go, tabla de rechazos de `NewRegistroHoras` con `ValidationError` y campo esperado: horas `""` ("las horas son obligatorias"), `"abc"`, `"1/2"`, `"1e1"` ("las horas deben ser un número"), `"-2"`, `"0"` ("las horas deben ser mayores que cero"), `"24.01"` y `"25"` ("las horas no pueden superar 24 por carga"), `"0.01"` ("las horas deben equivaler a minutos enteros") → campo `horas`; actividad `""`, `"   "` ("la actividad es obligatoria"), 501 runas ("la actividad no puede superar los 500 caracteres"; con 500 se acepta) → campo `actividad`; fecha `""`, `"2026-10-9"`, `"basura"`, `"2026-02-30"` ("formato esperado AAAA-MM-DD"), mañana ("la fecha no puede ser posterior a hoy"), el día anterior al inicio del Sprint ("la fecha no puede ser anterior al inicio del Sprint"), y con un Sprint vencido (`fin` < hoy) el día siguiente a su fin ("la fecha no puede ser posterior al fin del Sprint"; el día del fin se acepta), y con ese mismo Sprint vencido una fecha posterior a hoy (también posterior al fin) → "la fecha no puede ser posterior a hoy" (hoy se evalúa primero), y con un Sprint cuyo inicio es mañana la fecha de hoy → "la fecha no puede ser anterior al inicio del Sprint" y la fecha de inicio → "la fecha no puede ser posterior a hoy" (ninguna fecha entra hasta ese día) → campo `fecha`; orden: horas inválidas + actividad vacía + fecha inválida → error de `horas`; actividad vacía + fecha inválida → error de `actividad`. Además `ValidarSprintParaImputar`: historia sin Sprint, con Sprint `Pendiente` o `Finalizado`, o asignada a otro Sprint → `ErrHistoriaFueraDelSprintActivo`; Sprint `Activo` de la historia → `nil`; historia "Completada" en el Sprint Activo → `nil` (Clarifications). Stub vacío de `ValidarSprintParaImputar` — confirmar que FALLA
- [ ] T028 [US2] COMMIT `RED: agregar tests de validación de RegistroHoras (falla)`
- [ ] T029 [US2] GREEN: completar las validaciones en internal/domain/registro_horas.go (horas → actividad → fecha; en `MinutosDesdeHoras`, texto vacío, formato `^-?\d+(\.\d+)?$` antes de `big.Rat`, mayor que cero, minutos enteros y `minutos <= 1440`; actividad con `strings.TrimSpace` y `len([]rune(...)) > 500`; fecha con `time.Parse(time.DateOnly, ...)` y rango `inicio ≤ fecha ≤ min(hoy, fin)`, con "hoy" calculado con `inicioDelDiaCivil(ahora)` de internal/domain/estado_proyecto.go y comparando primero contra hoy, después contra el fin y por último contra el inicio) e implementar `ValidarSprintParaImputar(historia BacklogItem, sprint Sprint) error` — `go test ./...` en verde
- [ ] T030 [US2] COMMIT `GREEN: implementar las validaciones de RegistroHoras, test en verde`
- [ ] T031 [US2] RED: en internal/service/imputar_horas_test.go, rechazos de `ImputarHoras` sin llamar a `Registrar`: historia inexistente → `domain.ErrHistoriaNoEncontrada`; historia sin Sprint, sin consultar el Sprint (verificarlo con un fake en internal/service/fakes_test.go que envuelva `repository.NewSprintRepositoryEnMemoria()` y cuente las llamadas a `ObtenerPorID`), con Sprint `Pendiente` o `Finalizado` → `ErrHistoriaFueraDelSprintActivo`; integrante de otro proyecto o inexistente → `ErrIntegranteAjeno`; datos inválidos → el `ValidationError` del dominio; orden: historia fuera del Sprint + integrante ajeno + horas inválidas → `ErrHistoriaFueraDelSprintActivo`; integrante ajeno + horas inválidas → `ErrIntegranteAjeno`; `Registrar` que devuelve `ErrHistoriaFueraDelSprintActivo` (carrera con el cierre) → se propaga tal cual — confirmar que FALLA (el caso "historia sin Sprint" puede fallar por *panic* al desreferenciar `SprintID` en el código de T019: cuenta como falla, pero no se commitea código que lo esquive)
- [ ] T032 [US2] COMMIT `RED: agregar tests de rechazos y orden de errores de ImputarHoras (falla)`
- [ ] T033 [US2] GREEN: en internal/service/imputar_horas.go, si la historia no tiene Sprint devolver `ErrHistoriaFueraDelSprintActivo` sin consultar el Sprint; llamar a `ValidarSprintParaImputar` antes de buscar el integrante, y devolver `ErrIntegranteAjeno` cuando la búsqueda de T019 en `ListarIntegrantes(historia.ProyectoID)` (por `UserID`) no lo encuentra, antes de `NewRegistroHoras` — `go test ./...` en verde
- [ ] T034 [US2] COMMIT `GREEN: implementar rechazos de ImputarHoras en el orden de la spec, test en verde`
- [ ] T035 [US2] RED: en internal/http/horas_handler_test.go, errores de `POST /backlog/{id}/horas`: id `abc` o `0` → `400` campo `id`; historia inexistente con cuerpo inválido → `404` "historia no encontrada" (la existencia domina, como `PUT /backlog/{id}`); cuerpo que no es JSON → `400` "JSON inválido"; `"horas": "abc"` → `400` campo `horas` (no "JSON inválido"); `horas` ausente o `null` → `400` campo `horas`; `"horas": -2` → `400` campo `horas`; actividad vacía → `400` campo `actividad`; integrante ajeno → `400` campo `integrante_id`; historia fuera del Sprint activo → `400` campo `historia_id` — confirmar que FALLA
- [ ] T036 [US2] COMMIT `RED: agregar tests de errores del handler POST /backlog/{id}/horas (falla)`
- [ ] T037 [US2] GREEN: en internal/http/horas_handler.go consultar la historia antes de leer el cuerpo (reutilizar `service.ObtenerHistoria`); completar la extracción de `horas` de research.md R3: un string JSON pasa su contenido y `horas` ausente o `null` pasa texto vacío; los errores siguen saliendo por `escribirError` sin cambiarlo — tests en verde
- [ ] T038 [US2] COMMIT `GREEN: implementar errores del handler POST /backlog/{id}/horas, test en verde`
- [ ] T039 [US2] (solo si hay algo que limpiar) REFACTOR de las validaciones de internal/domain/registro_horas.go en helpers `validarHoras` / `validarActividad` / `validarFecha`, como `validarTitulo` en backlog_item.go + COMMIT `REFACTOR: separar las validaciones de RegistroHoras`

**Checkpoint**: US1 y US2 completas; la carga queda protegida contra datos inválidos y contra la carrera con el cierre.

---

## Phase 5: User Story 3 - Consultar el esfuerzo real de una historia (Priority: P2)

**Goal**: `GET /backlog/{id}/horas` devuelve el total acumulado y el detalle de las cargas.

**Independent Test**: sembrar varias cargas en una historia (de distintos integrantes y Sprints) y verificar el total exacto y el detalle; historia sin cargas → total 0 y `registros: []`; historia inexistente → 404.

- [ ] T040 [US3] RED: en internal/domain/registro_horas_test.go, `TotalMinutos` de cargas de 30, 90 y 120 minutos = 240 (0,5 h + 1,5 h + 2 h = 4 h exactas, SC-003) y de una lista vacía = 0; en internal/service/consultar_esfuerzo_real_test.go, `ConsultarEsfuerzoReal.Ejecutar(ctx, historiaID)` devuelve el detalle del repositorio y su total; historia sin cargas → total 0 y slice vacío; historia inexistente → `domain.ErrHistoriaNoEncontrada` sin consultar las cargas (que `fakeRegistroHorasRepository` cuente las llamadas a `ListarPorHistoria`). Stubs vacíos de `TotalMinutos` en internal/domain/registro_horas.go y del service en internal/service/consultar_esfuerzo_real.go — confirmar que FALLA
- [ ] T041 [US3] COMMIT `RED: agregar tests del esfuerzo real de una historia (falla)`
- [ ] T042 [US3] GREEN: implementar `TotalMinutos` en internal/domain/registro_horas.go y `ConsultarEsfuerzoReal` en internal/service/consultar_esfuerzo_real.go (`NewConsultarEsfuerzoReal(historias repository.HistoriaRepository, horas repository.RegistroHorasRepository)`, devuelve `service.EsfuerzoReal{HistoriaID, TotalMinutos, Registros}`, definido en ese mismo archivo) — `go test ./...` en verde
- [ ] T043 [US3] COMMIT `GREEN: implementar ConsultarEsfuerzoReal, test en verde`
- [ ] T044 [US3] RED: en internal/http/horas_handler_test.go, `GET /backlog/{id}/horas` → `200` con el cuerpo `EsfuerzoReal` del contrato (`historia_id`, `total_minutos`, `total_horas`, `registros` con `minutos` y `horas` por carga); historia sin cargas → `registros: []` (no `null`) y totales 0; historia inexistente → `404`; id inválido → `400` campo `id`. Stub de `ConsultarEsfuerzo` y tercer parámetro `consultar *service.ConsultarEsfuerzoReal` en `NewHorasHandler`; en el mismo commit, actualizar la llamada en cmd/api/main.go con `service.NewConsultarEsfuerzoReal(backlogRepo, horasRepo)` para que `go build ./...` siga pasando (la ruta se agrega en T048) — confirmar que FALLA
- [ ] T045 [US3] COMMIT `RED: agregar tests del handler GET /backlog/{id}/horas (falla)`
- [ ] T046 [US3] GREEN: implementar `ConsultarEsfuerzo` en internal/http/horas_handler.go y `esfuerzoRealResponse` en internal/http/horas_dto.go (slice inicializado con `make`) — tests en verde
- [ ] T047 [US3] COMMIT `GREEN: implementar handler GET /backlog/{id}/horas, test en verde`
- [ ] T048 [US3] Cablear en cmd/api/main.go `mux.HandleFunc("GET /backlog/{id}/horas", horasHandler.ConsultarEsfuerzo)` (el service ya se inyecta desde T044)
- [ ] T049 [US3] COMMIT `Cablear GET /backlog/{id}/horas en cmd/api`

**Checkpoint**: las tres historias funcionan de punta a punta.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T050 GREEN: en features/steps_horas.go, step definitions de Godog para features/imputacion_horas.feature sobre SQLite `:memory:` con migraciones y los services `ImputarHoras` / `ConsultarEsfuerzoReal` / `CerrarSprint`. Sembrar el Sprint y la historia con `UPDATE historias SET sprint_id = ?` (HU-06 no está en `main`), con fechas relativas a `hoyCivil()` como en features/steps_estado_proyecto.go. Definir `InitializeScenarioHoras(ctx *godog.ScenarioContext, base *scenarioContext)` que use `base.db` y registrarla en features/features_test.go como `InitializeScenarioHoras(ctx, sc) // HU-07`. No redefinir pasos existentes — `go test ./features/` en verde
- [ ] T051 COMMIT `GREEN: implementar los pasos BDD de imputación de horas, escenarios en verde`
- [ ] T052 Correr `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test ./...` sin errores (lo mismo que el check `backend` del CI)
- [ ] T053 Validar a mano con specs/019-hu-07-imputacion-horas/quickstart.md (pasos 1 a 8) y anotar el resultado
- [ ] T054 Marcar como `[X]` en este tasks.md las tareas terminadas y, con el OK de la usuaria, actualizar las casillas de las sub-issues #60, #61 y #62 + COMMIT `docs: actualizar el avance de tasks de HU-07`

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **US1 (Phase 3)** → **US2 (Phase 4)** → **US3 (Phase 5)** → **Polish (Phase 6)**.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- US2 depende de US1: endurece el mismo dominio, service y handler. US3 depende solo de Phase 2 (repositorio) y podría hacerse antes que US2.
- T050 (pasos BDD) va al final: los pasos llaman a código que tiene que existir para no romper el build.

### Parallel Opportunities

Hay una sola persona en la rama y cada RED/GREEN termina en un commit, así que el orden
secuencial es el recomendado. Si se reparte el trabajo:

- T003–T006 (esquema) y T007–T008 (entidad y puertos) tocan archivos distintos.
- En US3, T040 (dominio + service) y T044 (handler) usan archivos distintos, pero el handler necesita el service: conviene hacerlas en orden.

## Implementation Strategy

1. Setup + Foundational → esquema y persistencia listos.
2. US1 → imputar horas válidas de punta a punta (MVP).
3. US2 → rechazos en el orden de la spec. Juntas, US1 y US2 cubren los dos criterios de aceptación del backlog.
4. US3 → consulta del acumulado (la necesita el frontend, #63).
5. Polish → BDD en verde, checks del CI y quickstart.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- Los commits del `.feature` (T002) y de sus pasos (T051) usan `RED:` / `GREEN:` siguiendo el precedente de las specs 009 y 010, aunque la regla 2 de la plantilla pide commit sin prefijo para escenarios.
- La migración (T005) entra con un ciclo RED/GREEN porque tiene un test de esquema que la prueba (§9.6), no como tarea sin lógica.
