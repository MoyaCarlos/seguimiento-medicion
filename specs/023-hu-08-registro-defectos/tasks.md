---

description: "Tareas de HU-08: Registro y Seguimiento de Defectos (TDD con commits por fase)"
---

# Tasks: Registro y Seguimiento de Defectos (HU-08)

**Input**: Design documents from `/specs/023-hu-08-registro-defectos/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica es RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT).

**Compilación en RED**: todos los commits compilan. En cada RED se agregan, junto al test, solo
las firmas nuevas que el test necesita con un cuerpo vacío (devolver valor cero / `nil`), de modo
que `go build ./...` pase y el test falle por aserción, no por símbolo indefinido. Si el stub
cambia una firma que ya se usa (`NewDefectoHandler` en T042 y T052), el mismo RED actualiza la
llamada en cmd/api/main.go; es la única excepción a "solo archivos de test y firmas".

**Organization**: por user story de la spec (US1 registrar, US2 resolver, US3 consultar).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3

---

## Phase 1: Setup

**Purpose**: dejar los escenarios BDD escritos antes de implementar (flujo de AGENTS.md).

- [X] T001 Escribir features/registro_defectos.feature en Gherkin español: tag `@HU-08`, Característica con el Como/Quiero/Para de la spec ("Como integrante del equipo / Quiero registrar un defecto técnico indicando su severidad, la historia relacionada y el sprint de detección / Para poder hacer un seguimiento de la calidad y asegurar que se resuelva antes del despliegue final"), Antecedentes que arrancan con el paso existente `Dado que existe un proyecto con identificador 1` (features/steps_creacion.go, no se redefine) y siguen con una historia del proyecto, y exactamente 14 escenarios, uno por Acceptance Scenario de spec.md: US1.1 con Sprint "Activo", registrar con severidad "Alta" deja el defecto asociado a la historia, con ese Sprint de detección, estado "Pendiente" y sin Sprint de resolución; US1.2 historia fuera del Sprint activo (completada en un Sprint anterior o en el Product Backlog, como Esquema del escenario) se registra igual; US1.3 sin Sprint "Activo" (sin Sprints, solo "Pendiente" o solo "Finalizado", como Esquema del escenario) → rechazo "no hay un Sprint activo al que asociarlo como Sprint de detección"; US1.4 descripción vacía, severidad que no es Alta/Media/Baja o sin historia relacionada (Esquema del escenario) → error sobre el campo; US1.5 historia inexistente o de otro proyecto → error sobre la historia; US1.6 proyecto inexistente → "proyecto no encontrado"; US2.1 resolver un "Pendiente" con Sprint "Activo" → "Resuelto" con ese Sprint de resolución y el de detección sin cambios; US2.2 detectado en un Sprint finalizado y resuelto en el activo actual → Sprints distintos; US2.3 sin Sprint "Activo" → rechazo "no hay un Sprint activo al que asociarlo como Sprint de resolución" y sigue "Pendiente"; US2.4 resolver uno ya "Resuelto" → "el defecto ya está resuelto" y su Sprint de resolución no cambia; US2.5 defecto inexistente → "defecto no encontrado"; US3.1 listar un proyecto con defectos pendientes y resueltos → todos, en orden de registro, con sus datos; US3.2 proyecto sin defectos → lista vacía; US3.3 listar un proyecto inexistente → "proyecto no encontrado". En los rechazos de registro, verificar también que el listado del proyecto no cambia. Indicar en un comentario del archivo qué escenario cubre cada categoría: normal (US1.1, US2.1, US3.1), alternativo (US1.2, US2.2, US3.2), límite (ningún criterio de aceptación es de borde: el largo de 2000 caracteres y el mismo Sprint de detección y resolución se cubren con tests unitarios, T017 y T032), error (US1.3–US1.6, US2.3–US2.5, US3.3). No inventar escenarios fuera de la spec. Sin pasos implementados: con el runner estricto, `go test ./features/` FALLA
- [X] T002 COMMIT `RED: agregar escenarios BDD de registro de defectos (fallan)`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: esquema, entidad, búsqueda del Sprint activo, puerto y persistencia que usan las tres historias.

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

- [X] T003 RED: en internal/repository/esquema_test.go, (a) agregar `defectos` al mapa de estructura esperada de `TestEsquema_Estructura` con las columnas de MODELO-DE-DATOS.md §4.9 (`id` INTEGER pk; `proyecto_id`, `historia_id`, `sprint_deteccion_id` INTEGER NOT NULL; `descripcion`, `severidad` TEXT NOT NULL; `estado` TEXT NOT NULL con defecto `'Pendiente'`; `sprint_resolucion_id` INTEGER que admite NULL; `creado_en` TEXT NOT NULL) y las FK `historia_id,proyecto_id -> historias(id,proyecto_id)`, `sprint_deteccion_id,proyecto_id -> sprints(id,proyecto_id)`, `sprint_resolucion_id,proyecto_id -> sprints(id,proyecto_id)`, en orden alfabético (`fksDe` ordena con `slices.Sort`); actualizar la aserción `versionesAplicadas(...)` de `[]int{1, 2, 3, 4}` a `[]int{1, 2, 3, 4, 5}` (internal/repository/esquema_test.go, línea 220); (b) agregar a `TestEsquema_Filas` una constante `insertarDefecto` y los casos de §9.6 sobre `datosBase` (historia 100 del proyecto 1, Sprints 10 Activo y 11 Finalizado del proyecto 1, Sprint 21 del proyecto 2): `Pendiente` con `sprint_resolucion_id NULL` → aceptada; `Resuelto` con `sprint_resolucion_id` 10 → aceptada; `Resuelto` sin `sprint_resolucion_id` y `Pendiente` con él → `violaCheck`; severidad `'Crítica'` y estado `'Abierto'` → `violaCheck`; severidad NULL y `historia_id` NULL → `violaNotNull`; sin columna `estado` en el INSERT → queda `'Pendiente'` (aceptada); historia de otro proyecto (insertar una historia del proyecto 2 en la misma sentencia), Sprint de detección de otro proyecto (21) y Sprint de resolución de otro proyecto (21) → `violaFK`. Confirmar que FALLA (la tabla no existe)
- [X] T004 COMMIT `RED: agregar tests de esquema de defectos (fallan)`
- [X] T005 GREEN: crear internal/repository/migraciones/005_defectos.sql con el DDL de MODELO-DE-DATOS.md §4.9 tal cual: `id INTEGER PRIMARY KEY AUTOINCREMENT`, `estado TEXT NOT NULL DEFAULT 'Pendiente'`, `CHECK (severidad IN ('Alta','Media','Baja'))`, `CHECK (estado IN ('Pendiente','Resuelto'))`, `CHECK ((estado = 'Resuelto') = (sprint_resolucion_id IS NOT NULL))`, las tres FK compuestas, sin `ON DELETE` (D9); comentario de cabecera `-- Defectos: MODELO-DE-DATOS.md §4.9.` como la 004 — `go test ./...` en verde
- [X] T006 COMMIT `GREEN: agregar la migración 005 de defectos, tests en verde`
- [X] T007 Crear internal/domain/defecto.go con los tipos `Severidad` (`SeveridadAlta = "Alta"`, `SeveridadMedia = "Media"`, `SeveridadBaja = "Baja"`), `EstadoDefecto` (`DefectoPendiente = "Pendiente"`, `DefectoResuelto = "Resuelto"`), `Defecto` (`ID, ProyectoID, HistoriaID, SprintDeteccionID int64; Descripcion string; Severidad Severidad; Estado EstadoDefecto; SprintResolucionID *int64; CreadoEn time.Time`) y `DefectoDetalle` (`Defecto` embebido + `TituloHistoria string`); agregar en internal/domain/errors.go `ErrDefectoNoEncontrado = errors.New("defecto no encontrado")`, `ErrSinSprintActivoParaDetectar = ValidationError{Campo: "sprint_deteccion_id", Mensaje: "el proyecto no tiene un Sprint activo al que asociar el defecto como Sprint de detección"}`, `ErrSinSprintActivoParaResolver = ValidationError{Campo: "sprint_resolucion_id", Mensaje: "el proyecto no tiene un Sprint activo al que asociar el defecto como Sprint de resolución"}`, `ErrDefectoYaResuelto = ValidationError{Campo: "estado", Mensaje: "el defecto ya está resuelto"}`, `ErrHistoriaRelacionadaObligatoria = ValidationError{Campo: "historia_id", Mensaje: "la historia relacionada es obligatoria"}` y `ErrHistoriaRelacionadaAjena = ValidationError{Campo: "historia_id", Mensaje: "la historia relacionada no existe en el proyecto"}`, cada uno con su comentario `(HU-08)`; crear internal/repository/defecto_repository.go con el puerto `DefectoRepository { Registrar(ctx, d domain.Defecto) (domain.Defecto, error); ObtenerPorID(ctx, id int64) (domain.DefectoDetalle, error); Resolver(ctx, d domain.Defecto) error; ListarPorProyecto(ctx, proyectoID int64) ([]domain.DefectoDetalle, error) }`, documentando los errores de cada método según data-model.md (Puertos). Sin lógica
- [X] T008 COMMIT `Agregar la entidad Defecto, sus errores y el puerto DefectoRepository`
- [X] T009 [P] RED: en internal/domain/sprint_test.go, `SprintActivoDe(sprints []Sprint) (Sprint, bool)`: lista vacía → `false`; solo `Pendiente` y `Finalizado` → `false`; un `Activo` entre otros → ese Sprint y `true`. Stub vacío en internal/domain/sprint.go — confirmar que FALLA
- [X] T010 COMMIT `RED: agregar tests de SprintActivoDe (falla)`
- [X] T011 GREEN: implementar `SprintActivoDe` en internal/domain/sprint.go (recorre la lista y devuelve el primero con `Estado == SprintActivo`) — `go test ./...` en verde
- [X] T012 COMMIT `GREEN: implementar SprintActivoDe, test en verde`
- [X] T013 RED: en internal/repository/sqlite_defecto_test.go, tests de `SQLiteDefectoRepository` contra `:memory:` con migraciones (sembrar proyectos, historias y Sprints con INSERT directos, como internal/repository/sqlite_registro_horas_test.go): (a) `Registrar` persiste todos los campos (`estado` 'Pendiente', `sprint_resolucion_id` NULL, `creado_en` RFC 3339 UTC) y devuelve el ID asignado; (b) `Registrar` con el Sprint de detección ya `Finalizado` → `errors.Is(err, domain.ErrSinSprintActivoParaDetectar)` y no inserta ninguna fila; (c) `Registrar` con una historia de otro proyecto → `errors.Is(err, ErrIntegridadReferencial)`; (d) `ObtenerPorID` devuelve el `DefectoDetalle` con `TituloHistoria`; id inexistente → `domain.ErrDefectoNoEncontrado`; (e) `Resolver` de un `Pendiente` con el Sprint `Activo` persiste `estado` 'Resuelto' y `sprint_resolucion_id`; (f) `Resolver` de uno ya `Resuelto` → `domain.ErrDefectoYaResuelto` y conserva el `sprint_resolucion_id` original (resolución doble: alcanza con dos llamadas seguidas, porque la atomicidad la da el `UPDATE` condicional de una sola sentencia, research.md R4); (g) `Resolver` con el Sprint ya `Finalizado` → `domain.ErrSinSprintActivoParaResolver` y el defecto sigue `Pendiente`; (h) `ListarPorProyecto` devuelve solo los del proyecto, con `TituloHistoria`, ordenados por `id`; (i) proyecto sin defectos → slice vacío no nil. Stubs vacíos en internal/repository/sqlite_defecto.go — confirmar que FALLA
- [X] T014 COMMIT `RED: agregar tests del repositorio SQLite de defectos (fallan)`
- [X] T015 GREEN: implementar en internal/repository/sqlite_defecto.go `NewSQLiteDefectoRepository(db)`; `Registrar` con el INSERT condicional de research.md R4 (`INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, creado_en) SELECT ?, ?, ?, ?, 'Pendiente', ?, ? WHERE EXISTS (SELECT 1 FROM sprints WHERE id = ? AND estado = 'Activo')`; 0 filas → `domain.ErrSinSprintActivoParaDetectar`; errores con `errorDeEscritura`; `creado_en` con `time.Now().UTC().Format(time.RFC3339)`); `ObtenerPorID` con `JOIN historias` (`sql.ErrNoRows` → `domain.ErrDefectoNoEncontrado`); `Resolver` con el UPDATE condicional de R4 (`UPDATE defectos SET estado = 'Resuelto', sprint_resolucion_id = ? WHERE id = ? AND estado = 'Pendiente' AND EXISTS (SELECT 1 FROM sprints WHERE id = ? AND estado = 'Activo')`; 0 filas → releer el estado: `Resuelto` → `domain.ErrDefectoYaResuelto`, si no → `domain.ErrSinSprintActivoParaResolver`); `ListarPorProyecto` con `JOIN historias`, `ORDER BY d.id` y slice con `make` — `go test ./...` en verde
- [X] T016 COMMIT `GREEN: implementar el repositorio SQLite de defectos, tests en verde`

**Checkpoint**: esquema, entidad, búsqueda del Sprint activo y persistencia listos (incluidas las carreras con el cierre del Sprint y la resolución doble).

---

## Phase 3: User Story 1 - Registrar un defecto (Priority: P1) 🎯 MVP

**Goal**: `POST /projects/{id}/defectos` registra un defecto "Pendiente" con el Sprint activo como Sprint de detección y rechaza los inválidos en el orden de la spec.

**Independent Test**: sembrar un proyecto con un Sprint "Activo" y una historia; registrar un defecto y verificar `201` con `estado: Pendiente`, el `sprint_deteccion_id` del Sprint activo y `sprint_resolucion_id: null`; cada dato inválido da `400` con su `campo` y no registra nada.

- [X] T017 [US1] RED: en internal/domain/defecto_test.go, `NewDefecto(proyectoID, historiaID int64, activo Sprint, descripcion, severidad string) (Defecto, error)`: válido con `"  El total no se actualiza  "` y `"Alta"` → `Descripcion` sin espacios externos, `Severidad` Alta, `Estado` Pendiente, `SprintDeteccionID` = `activo.ID`, `SprintResolucionID` nil, `ProyectoID` e `HistoriaID` copiados; `"Media"` y `"Baja"` aceptadas; descripción de exactamente 2000 runas aceptada. Rechazos con `ValidationError`: descripción `""` y `"   "` ("la descripción es obligatoria"), 2001 runas ("la descripción no puede superar los 2000 caracteres") → campo `descripcion`; severidad `""`, `"alta"`, `" Alta"`, `"Crítica"` ("la severidad debe ser Alta, Media o Baja") → campo `severidad`; orden: descripción vacía + severidad inválida → error de `descripcion`. Stub vacío de `NewDefecto` en internal/domain/defecto.go — confirmar que FALLA
- [X] T018 [US1] COMMIT `RED: agregar tests de NewDefecto (falla)`
- [X] T019 [US1] GREEN: implementar `NewDefecto` en internal/domain/defecto.go (descripción con `strings.TrimSpace` y `len([]rune(...)) > 2000`, como `validarDescripcion` en internal/domain/backlog_item.go; severidad comparada exacta contra las tres constantes; descripción → severidad) — `go test ./...` en verde
- [X] T020 [US1] COMMIT `GREEN: implementar NewDefecto, test en verde`
- [X] T021 [US1] RED: en internal/service/registrar_defecto_test.go, con un fake nuevo `fakeDefectoRepository` en internal/service/fakes_test.go (guarda los defectos, asigna IDs y cuenta las llamadas a `Registrar`) y reutilizando `fakeHistoriaRepository` y `repository.NewSprintRepositoryEnMemoria()`: `RegistrarDefecto.Ejecutar(ctx, RegistrarDefectoInput{ProyectoID, HistoriaID int64; Descripcion, Severidad string})` con Sprint `Activo` y una historia del proyecto en el Product Backlog registra el defecto con el `SprintDeteccionID` del activo y devuelve `domain.DefectoDetalle` con el ID asignado y `TituloHistoria`; historia en un Sprint `Finalizado` → también se registra (Clarification 2). Rechazos sin llamar a `Registrar`: proyecto sin Sprints, con solo `Pendiente` o con solo `Finalizado` → `ErrSinSprintActivoParaDetectar`; `HistoriaID` 0 → `ErrHistoriaRelacionadaObligatoria`; historia inexistente o de otro proyecto → `ErrHistoriaRelacionadaAjena`; datos inválidos → el `ValidationError` del dominio; orden: sin Sprint activo + historia ajena + descripción vacía → `ErrSinSprintActivoParaDetectar`; historia ajena + descripción vacía → `ErrHistoriaRelacionadaAjena`; `Registrar` que devuelve `ErrSinSprintActivoParaDetectar` (carrera con el cierre) → se propaga tal cual. Stub vacío en internal/service/registrar_defecto.go — confirmar que FALLA
- [X] T022 [US1] COMMIT `RED: agregar tests de RegistrarDefecto (falla)`
- [X] T023 [US1] GREEN: implementar `RegistrarDefecto` en internal/service/registrar_defecto.go (`NewRegistrarDefecto(sprints repository.SprintRepository, historias repository.HistoriaRepository, defectos repository.DefectoRepository)`): `ListarPorProyecto` + `domain.SprintActivoDe` → `HistoriaID == 0` → `ObtenerPorID` de la historia (`ErrHistoriaNoEncontrada` o `historia.ProyectoID != input.ProyectoID` → `ErrHistoriaRelacionadaAjena`) → `NewDefecto` → `Registrar` → devolver `DefectoDetalle` con `historia.Titulo` — `go test ./...` en verde
- [X] T024 [US1] COMMIT `GREEN: implementar RegistrarDefecto, test en verde`
- [X] T025 [US1] RED: en internal/http/defecto_handler_test.go, `POST /projects/{id}/defectos` con `{"historia_id": 3, "descripcion": "El total no se actualiza", "severidad": "Alta"}` → `201` con el cuerpo `Defecto` del contrato (`id`, `proyecto_id`, `historia_id`, `historia`, `descripcion`, `severidad`, `estado: Pendiente`, `sprint_deteccion_id`, `sprint_resolucion_id: null`, `creado_en`); un cuerpo con `"estado": "Resuelto"` o `"sprint_deteccion_id": 99` se ignora; errores: id `abc` o `0` → `400` campo `id`; proyecto inexistente con cuerpo inválido → `404` "proyecto no encontrado" (la existencia domina); cuerpo que no es JSON → `400` "JSON inválido"; sin Sprint activo → `400` campo `sprint_deteccion_id`; sin `historia_id` → `400` campo `historia_id`; descripción vacía → `400` campo `descripcion`; severidad `"alta"` → `400` campo `severidad`. Los services se construyen con stubs propios del paquete `http` (los fakes de internal/service no son visibles), como en internal/http/horas_handler_test.go. Stubs de `DefectoHandler`, `NewDefectoHandler(obtenerProyecto *service.ObtenerProyecto, registrar *service.RegistrarDefecto)` y `Registrar` en internal/http/defecto_handler.go y DTOs en internal/http/defecto_dto.go — confirmar que FALLA
- [X] T026 [US1] COMMIT `RED: agregar tests del handler POST /projects/{id}/defectos (falla)`
- [X] T027 [US1] GREEN: implementar `Registrar` en internal/http/defecto_handler.go (id con `parsearID`; `obtenerProyecto.Ejecutar` antes de leer el cuerpo; decodificar `registrarDefectoRequest{HistoriaID int64; Descripcion, Severidad string}` sin campos de estado ni Sprints; responder `201`) y `defectoResponse` en internal/http/defecto_dto.go (`sprint_resolucion_id` como `*int64` para serializar `null`; `creado_en` en RFC 3339 UTC) — tests en verde
- [X] T028 [US1] COMMIT `GREEN: implementar handler POST /projects/{id}/defectos, test en verde`
- [X] T029 [US1] (no hizo falta: `NewDefecto` reutiliza `validarDescripcion` y `Severidad.EsValida`) (solo si hay algo que limpiar) REFACTOR de las validaciones de internal/domain/defecto.go en helpers `validarDescripcionDefecto` / `validarSeveridad`, como `validarTitulo` en backlog_item.go + COMMIT `REFACTOR: separar las validaciones de Defecto`
- [X] T030 [US1] Cablear en cmd/api/main.go `defectosRepo := repository.NewSQLiteDefectoRepository(db)`, `service.NewRegistrarDefecto(sprintRepo, backlogRepo, defectosRepo)`, `apihttp.NewDefectoHandler(obtenerProyecto, registrarDefecto)` y `mux.HandleFunc("POST /projects/{id}/defectos", defectoHandler.Registrar)`
- [X] T031 [US1] COMMIT `Cablear POST /projects/{id}/defectos en cmd/api`

**Checkpoint**: se pueden registrar defectos de punta a punta, con todos sus rechazos.

---

## Phase 4: User Story 2 - Marcar un defecto como resuelto (Priority: P1)

**Goal**: `POST /defectos/{id}/resolver` pasa un defecto "Pendiente" a "Resuelto" con el Sprint activo como Sprint de resolución.

**Independent Test**: sembrar un defecto "Pendiente" y un Sprint "Activo"; resolverlo y verificar `200` con `estado: Resuelto`, el `sprint_resolucion_id` del activo y el `sprint_deteccion_id` sin cambios; resolverlo otra vez → `400` campo `estado`.

- [X] T032 [US2] RED: en internal/domain/defecto_test.go, `(*Defecto).Resolver(activo *Sprint) error`: `Pendiente` con un Sprint `Activo` → `Estado` Resuelto, `SprintResolucionID` = `activo.ID`, `SprintDeteccionID` sin cambios; mismo Sprint que el de detección → aceptado; `Pendiente` con `activo` nil → `ErrSinSprintActivoParaResolver` y no cambia nada; `Resuelto` con un Sprint `Activo` → `ErrDefectoYaResuelto` y conserva su `SprintResolucionID`; `Resuelto` con `activo` nil → `ErrDefectoYaResuelto` (va antes que la falta de Sprint, research.md R3). Stub vacío del método en internal/domain/defecto.go — confirmar que FALLA
- [X] T033 [US2] COMMIT `RED: agregar tests de la resolución de un Defecto (falla)`
- [X] T034 [US2] GREEN: implementar `Resolver` en internal/domain/defecto.go (primero el estado, después `activo == nil`; copiar el ID antes de tomar su dirección) — `go test ./...` en verde
- [X] T035 [US2] COMMIT `GREEN: implementar la resolución de un Defecto, test en verde`
- [X] T036 [US2] RED: en internal/service/resolver_defecto_test.go, extendiendo `fakeDefectoRepository` (que cuente las llamadas a `Resolver`): `ResolverDefecto.Ejecutar(ctx, defectoID)` con un `Pendiente` y un Sprint `Activo` en su proyecto llama a `Resolver` del repositorio con `SprintResolucionID` = el activo y devuelve el `DefectoDetalle` resuelto; detectado en un Sprint `Finalizado` y resuelto en el activo → Sprints distintos; defecto inexistente → `domain.ErrDefectoNoEncontrado` sin consultar Sprints; ya `Resuelto` → `ErrDefectoYaResuelto` sin llamar a `Resolver` del repositorio; sin Sprint `Activo` → `ErrSinSprintActivoParaResolver` sin llamar a `Resolver`; el Sprint activo se busca en el proyecto del defecto, no en otro; `Resolver` del repositorio que devuelve `ErrDefectoYaResuelto` o `ErrSinSprintActivoParaResolver` (carreras) → se propaga tal cual. Stub vacío en internal/service/resolver_defecto.go — confirmar que FALLA
- [X] T037 [US2] COMMIT `RED: agregar tests de ResolverDefecto (falla)`
- [X] T038 [US2] GREEN: implementar `ResolverDefecto` en internal/service/resolver_defecto.go (`NewResolverDefecto(defectos repository.DefectoRepository, sprints repository.SprintRepository)`): `ObtenerPorID` → `ListarPorProyecto(detalle.ProyectoID)` + `domain.SprintActivoDe` → `detalle.Defecto.Resolver(activo o nil)` → `Resolver` del repositorio → devolver el detalle actualizado — `go test ./...` en verde
- [X] T039 [US2] COMMIT `GREEN: implementar ResolverDefecto, test en verde`
- [X] T040 [US2] Agregar `"/defectos"` a `PrefijosAPI` en internal/http/spa.go y `'/defectos': apiProxy` al `proxy` de web/vite.config.ts (research.md R1); `TestPrefijosAPI_IncluyenLosQueProxeaVite` y `cd web && npm test` siguen en verde
- [X] T041 [US2] COMMIT `Agregar el prefijo /defectos a PrefijosAPI y al proxy de Vite`
- [X] T042 [US2] RED: en internal/http/defecto_handler_test.go, `POST /defectos/{id}/resolver` sin cuerpo → `200` con el `Defecto` resuelto (`estado: Resuelto`, `sprint_resolucion_id` del activo); ya resuelto → `400` campo `estado`; sin Sprint activo → `400` campo `sprint_resolucion_id`; defecto inexistente → `404` "defecto no encontrado"; id `abc` o `0` → `400` campo `id`. Stub de `Resolver` y segundo servicio `resolver *service.ResolverDefecto` en `NewDefectoHandler`; en el mismo commit, actualizar la llamada en cmd/api/main.go con `service.NewResolverDefecto(defectosRepo, sprintRepo)` para que `go build ./...` siga pasando (la ruta se agrega en T046) — confirmar que FALLA
- [X] T043 [US2] COMMIT `RED: agregar tests del handler POST /defectos/{id}/resolver (falla)`
- [X] T044 [US2] GREEN: implementar `Resolver` en internal/http/defecto_handler.go (responder `200` con `defectoResponse`) y sumar `errors.Is(err, domain.ErrDefectoNoEncontrado)` al caso 404 de `escribirError` en internal/http/sprint_handler.go — tests en verde
- [X] T045 [US2] COMMIT `GREEN: implementar handler POST /defectos/{id}/resolver, test en verde`
- [X] T046 [US2] Cablear en cmd/api/main.go `mux.HandleFunc("POST /defectos/{id}/resolver", defectoHandler.Resolver)` (el service ya se inyecta desde T042)
- [X] T047 [US2] COMMIT `Cablear POST /defectos/{id}/resolver en cmd/api`

**Checkpoint**: US1 y US2 completas; cubren los tres criterios de aceptación del backlog.

---

## Phase 5: User Story 3 - Consultar los defectos de un proyecto (Priority: P2)

**Goal**: `GET /projects/{id}/defectos` devuelve los defectos del proyecto en orden de registro.

**Independent Test**: sembrar defectos pendientes y resueltos en un proyecto y uno en otro; el listado devuelve solo los del proyecto, en orden de registro; proyecto sin defectos → `[]`; proyecto inexistente → `404`.

- [X] T048 [US3] RED: en internal/service/listar_defectos_test.go, con `fakeProjectRepository` y `fakeDefectoRepository` (que cuente las llamadas a `ListarPorProyecto`): `ListarDefectos.Ejecutar(ctx, proyectoID)` devuelve los `DefectoDetalle` del repositorio; proyecto sin defectos → slice vacío no nil; proyecto inexistente → `domain.ErrProyectoNoEncontrado` sin consultar los defectos. Stub vacío en internal/service/listar_defectos.go — confirmar que FALLA
- [X] T049 [US3] COMMIT `RED: agregar tests de ListarDefectos (falla)`
- [X] T050 [US3] GREEN: implementar `ListarDefectos` en internal/service/listar_defectos.go (`NewListarDefectos(proyectos repository.ProjectRepository, defectos repository.DefectoRepository)`: `ObtenerPorID` del proyecto → `ListarPorProyecto`) — `go test ./...` en verde
- [X] T051 [US3] COMMIT `GREEN: implementar ListarDefectos, test en verde`
- [X] T052 [US3] RED: en internal/http/defecto_handler_test.go, `GET /projects/{id}/defectos` → `200` con un arreglo de `Defecto` del contrato en orden de registro (un `Pendiente` con `sprint_resolucion_id: null` y un `Resuelto` con su Sprint); proyecto sin defectos → `[]` (no `null`); proyecto inexistente → `404`; id inválido → `400` campo `id`. Stub de `Listar` y tercer servicio `listar *service.ListarDefectos` en `NewDefectoHandler`; en el mismo commit, actualizar la llamada en cmd/api/main.go con `service.NewListarDefectos(proyectos, defectosRepo)` (la ruta se agrega en T056) — confirmar que FALLA
- [X] T053 [US3] COMMIT `RED: agregar tests del handler GET /projects/{id}/defectos (falla)`
- [X] T054 [US3] GREEN: implementar `Listar` en internal/http/defecto_handler.go (slice de respuesta inicializado con `make`) — tests en verde
- [X] T055 [US3] COMMIT `GREEN: implementar handler GET /projects/{id}/defectos, test en verde`
- [X] T056 [US3] Cablear en cmd/api/main.go `mux.HandleFunc("GET /projects/{id}/defectos", defectoHandler.Listar)`
- [X] T057 [US3] COMMIT `Cablear GET /projects/{id}/defectos en cmd/api`

**Checkpoint**: las tres historias funcionan de punta a punta.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T058 GREEN: en features/steps_defectos.go, step definitions de Godog para features/registro_defectos.feature sobre la base de `base.db` (SQLite `:memory:` con migraciones) y los services `RegistrarDefecto` / `ResolverDefecto` / `ListarDefectos`, más `CrearSprint` / `IniciarSprint` / `CerrarSprint` de HU-05 para armar los Sprints "Pendiente", "Activo" y "Finalizado", con fechas relativas a `hoyCivil()` como en features/steps_horas.go. Para US1.2, la historia "completada en un Sprint anterior" se siembra con `UPDATE historias SET estado = 'Completada', fecha_finalizacion = ?, sprint_id = ? WHERE id = ?` directo en `base.db` sobre un Sprint ya `Activo`, y recién después se cierra ese Sprint con `CerrarSprint` (que deja en el Sprint las completadas) y se inicia otro; así la historia queda en un Sprint `Finalizado` sin pasar por la edición de HU-15. Definir `InitializeScenarioDefectos(ctx *godog.ScenarioContext, base *scenarioContext)` y registrarla en features/features_test.go como `InitializeScenarioDefectos(ctx, sc) // HU-08`. No redefinir pasos existentes — `go test ./features/` en verde
- [X] T059 COMMIT `GREEN: implementar los pasos BDD de registro de defectos, escenarios en verde`
- [X] T060 Correr `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test ./...`, y en web/ `npm run lint` y `npm test`, sin errores (los checks `backend` y `frontend` del CI)
- [X] T061 Validar a mano con specs/023-hu-08-registro-defectos/quickstart.md (pasos 1 a 9) y anotar el resultado al final de ese archivo
- [X] T062 Marcar como `[X]` en este tasks.md las tareas terminadas y, con el OK de la usuaria, actualizar las casillas de las sub-issues #65, #66 y #67 + COMMIT `docs: actualizar el avance de tasks de HU-08`

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **US1 (Phase 3)** → **US2 (Phase 4)** → **US3 (Phase 5)** → **Polish (Phase 6)**.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- US2 depende de US1 solo por `DefectoHandler` (se crea en T025) y por tener defectos que resolver en las pruebas de punta a punta; su dominio y su service dependen solo de la Phase 2.
- US3 depende de la Phase 2 (repositorio) y de `DefectoHandler`; podría hacerse antes que US2.
- T058 (pasos BDD) va al final: los pasos llaman a código que tiene que existir para no romper el build.

### Parallel Opportunities

Hay una sola persona en la rama y cada RED/GREEN termina en un commit, así que el orden
secuencial es el recomendado. Si se reparte el trabajo:

- T009–T012 (`SprintActivoDe`, internal/domain/sprint.go) no comparte archivos con T003–T008 ni con T013–T016.
- T040–T041 (prefijo `/defectos`) toca solo internal/http/spa.go y web/vite.config.ts y se puede hacer en cualquier momento antes de T046.
- En cada historia, el dominio y el service (archivos distintos) van antes que el handler, que los necesita.

## Implementation Strategy

1. Setup + Foundational → esquema y persistencia listos, con las carreras cubiertas.
2. US1 → registrar defectos de punta a punta (MVP).
3. US2 → resolverlos. Juntas, US1 y US2 cubren los tres criterios de aceptación del backlog.
4. US3 → listado (lo necesita el frontend, #68).
5. Polish → BDD en verde, checks del CI y quickstart.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- Los commits del `.feature` (T002) y de sus pasos (T059) usan `RED:` / `GREEN:` siguiendo el precedente de las specs 009, 010 y 019, aunque la regla 2 de la plantilla pide commit sin prefijo para escenarios.
- La migración (T005) entra con un ciclo RED/GREEN porque tiene un test de esquema que la prueba (§9.6), no como tarea sin lógica.
- `MODELO-DE-DATOS.md` no cambia: la migración 005 copia el DDL de §4.9.
