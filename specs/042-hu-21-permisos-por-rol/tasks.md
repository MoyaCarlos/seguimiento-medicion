---

description: "Tareas de HU-21: Permisos por rol (TDD con commits por fase)"
---

# Tasks: Permisos por rol (HU-21)

**Input**: Design documents from `specs/042-hu-21-permisos-por-rol/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/api-permisos.md](contracts/api-permisos.md)

**Tests**: OBLIGATORIOS. La constitución (Principio I, NON-NEGOTIABLE) exige TDD estricto con
un commit por fase.

**Organization**: por user story de la spec. Orden de prioridad: US3 → US1 → US4 → US2 (P1)
→ US5 → US6 (P2) → US7 (frontend).

## Reglas que valen para TODAS las tareas

- **Sin esquema nuevo** (research R11). Si una tarea parece necesitar una tabla, columna o
  índice nuevos: frenar y consultar al Scrum Master (AGENTS.md).
- **Orden de errores** (FR-006, research R4): 401 (`conSesion`) → cargar el recurso (404) →
  `exigirIntegrante` / `exigirScrumMaster` (403) → validación de dominio y conflictos
  (400/409). Si un servicio hoy valida antes de cargar el recurso, se reordena en su GREEN.
- **Mensajes**: `ErrNoEsIntegrante` = «No sos integrante de este proyecto»;
  `ErrNoEsScrumMaster` = «Solo el Scrum Master puede realizar esta acción».
- **Los GREEN dejan `go test ./...` en verde, incluido `./features/`**: si un servicio cambia
  de firma (recibe `usuarioID`), el mismo GREEN ajusta `cmd/api/main.go`, los handlers y los
  `features/steps_*.go` que lo llaman, para que compile y siga en verde. Los steps de
  servicio fijan quién actúa con el paso «Dado que inicié sesión como "X"» (T011) en vez de
  leer `usuario_id`, `integrante_id` o `creador` de la tabla del escenario.
- **Sin atajos de prueba en producción** (FR-028).
- Ninguna tarea hace `git push` ni merge.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: se puede hacer en paralelo (archivos distintos, sin dependencias pendientes)
- **[Story]**: user story de la spec (US1…US7)

---

## Phase 1: Setup

- [X] T001 Verificar la base: `go build ./... && go vet ./... && go test ./...` y `cd web && npm test` en verde sobre la rama `feature/HU-21-permisos-por-rol` antes de tocar código (sin commit).

---

## Phase 2: Foundational (Blocking Prerequisites)

**⚠️ CRITICAL**: ninguna user story empieza antes de terminar esta fase.

### Reglas de pertenencia y rol (servicio)

- [X] T002 RED: tests en internal/service/integrantes_test.go (archivo nuevo) para `exigirIntegrante(ctx, proyectos, proyectoID, usuarioID)` → nil si integra, `domain.ErrNoEsIntegrante` si no; y `exigirScrumMaster` → nil si es Scrum Master, `ErrNoEsScrumMaster` si es Product Builder, **`ErrNoEsIntegrante` si no integra** (hoy devuelve `ErrNoEsScrumMaster`). Incluir aserciones de los mensajes: `ErrNoEsIntegrante.Error() == "No sos integrante de este proyecto"` y `ErrNoEsScrumMaster.Error() == "Solo el Scrum Master puede realizar esta acción"` — confirmar que FALLA
- [X] T003 COMMIT `RED: pertenencia y rol por proyecto con los mensajes de HU-21 (falla)`
- [X] T004 GREEN: implementar `exigirIntegrante` y el cambio de `exigirScrumMaster` en internal/service/integrantes.go; cambiar los mensajes en internal/domain/errors.go. Actualizar los tests y `.feature` existentes que esperan el mensaje viejo («solo el Scrum Master del proyecto puede hacerlo», `features/votacion_planning_poker.feature`, `features/consenso_estimacion.feature`, tests de internal/service y internal/http) y los que esperaban 403 de rol para un no integrante — `go test ./...` en verde
- [X] T005 COMMIT `GREEN: pertenencia y rol por proyecto con los mensajes de HU-21, test en verde`

### Sesión en la capa HTTP

- [X] T006 RED: tests en internal/http/sesion_middleware_test.go (nuevo) para `conSesion(obtener)(handler)`: sin cabecera, esquema distinto de Bearer, token mal formado, inexistente o vencido → 401 `{"mensaje":"Tu sesión no es válida"}` y el handler NO se ejecuta; con sesión vigente el handler se ejecuta y `usuarioDeSesion(r)` devuelve el `ID` y el `Nombre` del usuario. Usar un fake del puerto que usa `service.ObtenerSesion` — confirmar que FALLA
- [X] T007 COMMIT `RED: exigir sesión vigente en la capa HTTP (falla)`
- [X] T008 GREEN: implementar `conSesion` y `usuarioDeSesion` en internal/http/sesion_middleware.go, reutilizando `tokenBearer` (auth_handler.go) y `service.ObtenerSesion`; guardar el usuario en el `context` con una clave no exportada — tests en verde
- [X] T009 COMMIT `GREEN: exigir sesión vigente en la capa HTTP, test en verde`

### Registro de rutas compartido

- [X] T010 Mover el registro de rutas de cmd/api/main.go a `RegistrarRutas(mux *http.ServeMux, h Handlers, conSesion func(http.HandlerFunc) http.HandlerFunc)` en internal/http/rutas.go (archivo nuevo; `Handlers` agrupa los handlers ya construidos, incluidos `votacionHandler` y `demoHandler`). Las rutas se declaran en una **tabla exportada para tests** (`[]ruta{patron string; nivel nivelRuta; handler http.HandlerFunc}`, con `nivel` ∈ `publica`, `soloSesion`, `sesionIntegrante`, `sesionScrumMaster`, `lectura`) que `RegistrarRutas` recorre; las rutas de `VotacionHandler.Registrar` y `DemoHandler.Registrar` pasan a esa tabla. En este paso **ninguna ruta se envuelve todavía** (comportamiento idéntico); `main.go` solo construye y llama. `go test ./...` en verde
- [X] T011 Extraer el armado de repositorios, servicios y handlers de cmd/api/main.go a `app.NuevoMux(db, app.Opciones{DemoHabilitado, Ahora, Contrasenas})` en internal/app/app.go (raíz de composición compartida: `main` y Godog arman la misma API, sin atajos). Agregar en features/steps_permisos.go el contexto compartido `sesionDeEscenario` y el paso «Dado que inicié sesión como "X"», que fija quién actúa: `actor()` devuelve su **nombre** (cada step de servicio lo resuelve en su propia base, porque cada contexto de Godog arma la suya) y `peticion(metodo, ruta, cuerpo)` levanta la API con `app.NuevoMux` sobre una base en memoria solo si hace falta, registra o hace ingresar a la persona con `/auth/...` y agrega `Authorization: Bearer`. Inicializarlo en features/features_test.go antes de los demás `InitializeScenario*`. Sin escenarios nuevos todavía
- [X] T012 COMMIT `Extraer RegistrarRutas y la raíz de composición, y agregar el paso de sesión de Godog`

**Checkpoint**: reglas y sesión listas; todavía ninguna ruta exige sesión.

---

## Phase 3: User Story 3 - Sin sesión (Priority: P1) — Escenario 3

**Goal**: toda ruta protegida responde 401 sin sesión; las públicas siguen funcionando.

**Independent Test**: internal/http/rutas_permisos_test.go recorre la tabla «Qué se protege».

- [X] T013 [US3] RED: test de tabla en internal/app/rutas_permisos_test.go (nuevo; en `app` y no en `http` para armar la API real con `app.NuevoMux` sin ciclo de imports) que arma el mux real: cada ruta de los niveles «Solo sesión», «Sesión e integrante», «Sesión y Scrum Master» y **todas las lecturas `GET`** de la spec (incluidas `GET /projects/{id}/status`, `/members`, `/sprints`, `/metricas`, `/historial`, `/historial/{sprintId}`, `/avance`, `/defectos`, `/reporte`, `GET /sprints/{id}/reporte`, `/backlog`, `GET /backlog/{id}`, `/horas`, `/votacion/rondas`, `GET /users`, `GET /auth/yo`, `POST /auth/salir`) sin token → 401; `POST /auth/registro`, `POST /auth/ingreso`, `GET /health`, `GET /demo/estado`, `POST /demo/datos` sin token → distinto de 401. El test recorre la **tabla declarativa** de T010 y además compara su lista de patrones con una lista fija escrita en el test a partir de la spec: una ruta nueva sin clasificar, o clasificada en un nivel distinto al de la spec, hace fallar el test — confirmar que FALLA
- [X] T014 [US3] COMMIT `RED: rutas protegidas responden 401 sin sesión (falla)`
- [X] T015 [US3] GREEN: `RegistrarRutas(mux, h, conSesion)` en internal/http/rutas.go envuelve toda ruta que no es `Publica` (incluidas `GET /auth/yo` y `POST /auth/salir`: el envoltorio no cambia su respuesta) y `app.NuevoMux` le pasa `ConSesion(obtenerSesion)`. Se eliminan `VotacionHandler.Registrar` y `DemoHandler.Registrar` (repetían rutas de la tabla); sus tests y features/steps_datos_prueba.go registran desde `Rutas`. features/autenticacion.feature: el escenario de HU-20 «Las rutas existentes siguen respondiendo sin token» pasa a «Una ruta protegida sin token se rechaza», por la API real. Los demás steps de Godog que usan HTTP arman su mux con handlers sueltos y no cambian en esta fase; se ajustan cuando sus handlers empiecen a leer la sesión (US1, US4, US2) — `go test ./...` en verde
- [X] T016 [US3] COMMIT `GREEN: rutas protegidas responden 401 sin sesión, test en verde`

**Checkpoint**: la API exige sesión; las reglas por proyecto todavía no.

---

## Phase 4: User Story 1 - Identidad de la sesión (Priority: P1) — Escenario 7 🎯 MVP

**Goal**: votar, registrar horas, consultar la votación y crear proyectos quedan a nombre de la
sesión; `usuario_id`, `integrante_id` y `creador` se ignoran.

**Independent Test**: dos integrantes; uno vota; el otro consulta con `?usuario_id=<el primero>` y no ve su voto oculto.

### Crear proyecto

- [ ] T017 [P] [US1] RED: en internal/service/crear_proyecto_test.go, `CrearProyectoInput` recibe `CreadorID int64` (no `Creador string`): el proyecto queda con ese usuario como único Scrum Master y **no se crea ningún usuario**; `CreadorID <= 0` → error (no se persiste nada). En internal/http/project_handler_test.go: `POST /projects` con sesión de ana y cuerpo `{"nombre":"P","creador":"bruno"}` → 201, el Scrum Master es ana y la respuesta trae `"rol":"scrum_master"` — confirmar que FALLA
- [ ] T018 [US1] COMMIT `RED: el creador del proyecto sale de la sesión (falla)`
- [ ] T019 [US1] GREEN: internal/service/crear_proyecto.go usa `CreadorID` y `GuardarConScrumMaster`; quitar `Creador` de `crearProyectoRequest` en internal/http/project_dto.go; `ProjectHandler.Crear` pasa `usuarioDeSesion(r).ID` y agrega `rol` a la respuesta; ajustar features/steps_proyecto.go, steps_creacion.go y los demás steps que crean proyectos con `creador` (usar `actor()`), y features/creacion_proyecto_equipo.feature — `go test ./...` en verde
- [ ] T020 [US1] COMMIT `GREEN: el creador del proyecto sale de la sesión, test en verde`

### Votación

- [ ] T021 [P] [US1] RED: en internal/http/votacion_handler_test.go: `POST /backlog/{id}/votacion/votos` con cuerpo `{"valor":5,"usuario_id":<otro>}` registra el voto **a nombre de la sesión**; un cuerpo sin `usuario_id` ya no da 400 «el usuario es obligatorio»; `GET /backlog/{id}/votacion?usuario_id=<otro>` devuelve «mi voto» de la sesión y nunca el voto oculto del otro; iniciar/revelar/rondas/estimación usan la sesión como actor — confirmar que FALLA
- [ ] T022 [US1] COMMIT `RED: la votación toma la identidad de la sesión (falla)`
- [ ] T023 [US1] GREEN: quitar `UsuarioID` de los requests de internal/http/votacion_dto.go, eliminar `leerPedido` y la lectura de la query `usuario_id` en internal/http/votacion_handler.go; pasar `usuarioDeSesion(r).ID` a `Votar`, `ConsultarVotacion`, `IniciarVotacion`, `RevelarVotacion`, `AbrirNuevaRonda`, `RegistrarEstimacion`. Ajustar features/steps_votacion.go, steps_consenso.go, votacion_planning_poker.feature y consenso_estimacion.feature (quién actúa = sesión) — `go test ./...` en verde
- [ ] T024 [US1] COMMIT `GREEN: la votación toma la identidad de la sesión, test en verde`

### Horas

- [ ] T025 [P] [US1] RED: en internal/service/imputar_horas_test.go el input trae el `UsuarioID` de quien actúa (renombrar el campo si hoy es `IntegranteID`) y en internal/http/horas_handler_test.go `POST /backlog/{id}/horas` con `{"horas":2,"integrante_id":<otro>,...}` registra las horas a nombre de la sesión — confirmar que FALLA
- [ ] T026 [US1] COMMIT `RED: las horas quedan a nombre de la sesión (falla)`
- [ ] T027 [US1] GREEN: quitar `IntegranteID` del request en internal/http/horas_dto.go (se conserva en la respuesta si ya existe); `HorasHandler.Imputar` pasa `usuarioDeSesion(r).ID`; ajustar features/steps_horas.go e imputacion_horas.feature — `go test ./...` en verde
- [ ] T028 [US1] COMMIT `GREEN: las horas quedan a nombre de la sesión, test en verde`
- [ ] T029 [US1] (solo si hay algo que limpiar) REFACTOR: eliminar `resolverUsuario` de internal/service/usuarios.go si quedó sin uso (se termina de quitar en US6) + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: escenario 7 cumplido en la API.

---

## Phase 5: User Story 4 - Pertenencia al proyecto (Priority: P1) — Escenario 4

**Goal**: quien no integra el proyecto recibe 403 «No sos integrante de este proyecto» en el
panel, en la votación y en las escrituras de nivel «sesión e integrante».

**Independent Test**: un usuario con sesión, ajeno al proyecto, recibe 403 en cada una de esas rutas y nada cambia; puede leer historias, Sprints, horas y defectos (límite conocido).

- [ ] T030 [P] [US4] RED: tests de servicio con fakes en memoria — un no integrante recibe `ErrNoEsIntegrante` y **no se persiste nada** en: `ObtenerProyecto` (internal/service/obtener_proyecto_test.go), `ConsultarVotacion` (consultar_votacion_test.go), `CrearHistoriaBacklog` (crear_historia_test.go), `EditarHistoria` (editar_historia_test.go), `CambiarEstadoHistoria` (cambiar_estado_historia_test.go), `AsignarHistoriaASprint` (asignar_historia_sprint_test.go), `ImputarHoras` (imputar_horas_test.go), `RegistrarDefecto` (registrar_defecto_test.go), `ResolverDefecto` (resolver_defecto_test.go), `Votar` (votar_test.go). En la API (internal/http/rutas_permisos_test.go), un no integrante con sesión: `GET /projects/{id}` y `GET /backlog/{id}/votacion` → 403 `{"mensaje":"No sos integrante de este proyecto"}`; `GET /projects/{id}/backlog`, `/sprints`, `/defectos` y `GET /backlog/{id}/horas` → 200 (límite conocido). Por cada servicio: recurso inexistente → 404 antes que el 403, y un integrante con datos inválidos sigue recibiendo el 400 actual — confirmar que FALLA
- [ ] T031 [US4] COMMIT `RED: escrituras y lecturas protegidas exigen ser integrante (falla)`
- [ ] T032 [US4] GREEN: cada servicio recibe el `usuarioID` de quien actúa (parámetro o campo del input) y llama a `exigirIntegrante` sobre el proyecto del recurso **después de cargarlo y antes de validar**; inyectar `repository.ProjectRepository` donde falte. Handlers en internal/http (project, backlog, horas, defecto, votacion) pasan `usuarioDeSesion(r).ID`. Ajustar cmd/api/main.go y los steps de features/ (backlog, estado_historia, asignacion, horas, defectos, votacion, obtener_historia) para que el actor integre el proyecto. `ImputarHoras` reemplaza `ErrIntegranteAjeno` (400) por `ErrNoEsIntegrante` (403) y se elimina `ErrIntegranteAjeno` de internal/domain/errors.go; ajustar el escenario correspondiente de features/imputacion_horas.feature (pasa de 400 a 403 con «No sos integrante de este proyecto») — `go test ./...` en verde
- [ ] T033 [US4] COMMIT `GREEN: escrituras y lecturas protegidas exigen ser integrante, test en verde`
- [ ] T034 [US4] (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: escenario 4 cumplido.

---

## Phase 6: User Story 2 - Acciones de gestión solo para el Scrum Master (Priority: P1) — Escenarios 1 y 2

**Goal**: Sprints (crear, iniciar, cerrar), votación (iniciar, revelar, nueva ronda, estimación),
editar proyecto y configurar equipo: solo el Scrum Master.

**Independent Test**: cada acción con la sesión del Scrum Master se realiza; con la del Product Builder → 403 y nada cambia.

- [ ] T035 [P] [US2] RED: tests de servicio — Product Builder → `ErrNoEsScrumMaster`, no integrante → `ErrNoEsIntegrante`, Scrum Master → se realiza, y en los rechazos no cambia nada, para `CrearSprint` (internal/service/crear_sprint_test.go), `IniciarSprint` (iniciar_sprint_test.go), `CerrarSprint` (cerrar_sprint_test.go), `EditarProyecto` (editar_proyecto_test.go) y `AsignarIntegrante` (asignar_integrante_test.go). Un Product Builder con datos inválidos (por ejemplo fechas de inicio de Sprint inválidas, nombre de proyecto vacío) recibe 403, no 400. Los de votación (`IniciarVotacion`, `RevelarVotacion`, `AbrirNuevaRonda`, `RegistrarEstimacion`) ya exigen Scrum Master: agregar solo el caso «no integrante → `ErrNoEsIntegrante`» si falta — confirmar que FALLA
- [ ] T036 [US2] COMMIT `RED: acciones de gestión exclusivas del Scrum Master (falla)`
- [ ] T037 [US2] GREEN: agregar el `usuarioID` de quien actúa a `CrearSprint`, `IniciarSprint`, `CerrarSprint`, `EditarProyecto` y `AsignarIntegrante` y llamar a `exigirScrumMaster` después de cargar el recurso y antes de validar (reordenar si hoy validan antes). Handlers de internal/http/sprint_handler.go y project_handler.go pasan `usuarioDeSesion(r).ID`. Ajustar cmd/api/main.go y los steps de features/ (sprints, consultar_sprints, asignacion, estado_proyecto, historial_sprints, metricas, reporte, avance, defectos, proyecto, scrum_master_unico) para que el actor de las acciones de gestión sea el Scrum Master del proyecto — `go test ./...` en verde
- [ ] T038 [US2] COMMIT `GREEN: acciones de gestión exclusivas del Scrum Master, test en verde`
- [ ] T039 [US2] (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: escenarios 1 y 2 cumplidos en la API (la parte de pantalla es US7).

---

## Phase 7: User Story 5 - Mis proyectos (Priority: P2) — Escenarios 5 y 9

**Goal**: `GET /projects` devuelve solo los proyectos que integra la persona, con su rol; vacía si ninguno. `GET /projects/{id}` trae el `rol`.

**Independent Test**: dos personas con proyectos distintos ven cada una solo los suyos.

- [ ] T040 [P] [US5] RED: en internal/repository/sqlite_project_listar_test.go y project_repository_en_memoria_listar_test.go, `ListarPorIntegrante(ctx, usuarioID)` devuelve solo los proyectos donde el usuario está en `integrantes`, con el `rol` de esa fila, ordenados por `proyectos.id`; usuario sin proyectos → lista vacía (no nil) — confirmar que FALLA
- [ ] T041 [US5] COMMIT `RED: listar los proyectos de un integrante con su rol (falla)`
- [ ] T042 [US5] GREEN: agregar `ListarPorIntegrante` al puerto internal/repository/project_repository.go, a internal/repository/sqlite_project.go (`JOIN integrantes`, sin cambio de esquema) y a project_repository_en_memoria.go; el tipo de retorno lleva proyecto + rol (por ejemplo `domain.ProyectoConRol`) — tests en verde
- [ ] T043 [US5] COMMIT `GREEN: listar los proyectos de un integrante con su rol, test en verde`
- [ ] T044 [US5] RED: internal/service/listar_proyectos_test.go — `ListarProyectos.Ejecutar(ctx, usuarioID)` devuelve solo los suyos con rol; internal/http/project_handler_listar_test.go — `GET /projects` responde cada proyecto con `"rol"`, `[]` si no integra ninguno, nunca 403; `GET /projects/{id}` trae `"rol"` de quien consulta — confirmar que FALLA
- [ ] T045 [US5] COMMIT `RED: «Mis proyectos» con el rol de cada uno (falla)`
- [ ] T046 [US5] GREEN: internal/service/listar_proyectos.go usa `ListarPorIntegrante`; `ObtenerProyecto` devuelve también el rol; `proyectoResponse` en internal/http/project_dto.go suma `Rol string \`json:"rol,omitempty"\``; ajustar features/steps_listado_proyectos.go y listado_proyectos.feature — `go test ./...` en verde
- [ ] T047 [US5] COMMIT `GREEN: «Mis proyectos» con el rol de cada uno, test en verde`

---

## Phase 8: User Story 6 - Configurar el equipo con usuarios registrados (Priority: P2) — Escenarios 6 y 8

**Goal**: `GET /users` lista los usuarios con cuenta; agregar un integrante solo acepta usuarios registrados y siempre como Product Builder.

**Independent Test**: el Scrum Master suma a un usuario de la lista (201, `product_builder`); un inexistente o sin contraseña → 404; `rol: scrum_master` → 409.

- [ ] T048 [P] [US6] RED: en internal/repository/sqlite_user_test.go, `ListarRegistrados(ctx)` devuelve `id` y `nombre` solo de usuarios con `contrasena_hash IS NOT NULL`, ordenados por `nombre_normalizado`; un usuario sin contraseña no aparece — confirmar que FALLA
- [ ] T049 [US6] COMMIT `RED: listar los usuarios registrados (falla)`
- [ ] T050 [US6] GREEN: agregar `ListarRegistrados` al puerto internal/repository/user_repository.go y a internal/repository/sqlite_user.go (y al fake de internal/service/fakes_test.go si implementa el puerto) — tests en verde
- [ ] T051 [US6] COMMIT `GREEN: listar los usuarios registrados, test en verde`
- [ ] T052 [US6] RED: internal/service/listar_usuarios_test.go (nuevo) y internal/http/usuarios_handler_test.go (nuevo): `GET /users` con sesión → 200 `[{"id":1,"nombre":"ana"}]`, sin campos de contraseña ni sesión; internal/service/asignar_integrante_test.go: usuario registrado sin `rol` o con `product_builder` → `product_builder`; `rol: scrum_master` → `ErrScrumMasterYaExiste` (409); usuario inexistente **o sin contraseña** → el mismo `domain.ErrUsuarioNoEncontrado` (404 «usuario no encontrado») sin crear usuario ni integrante; ya integrante → 400 de duplicado actual — confirmar que FALLA
- [ ] T053 [US6] COMMIT `RED: equipo solo con usuarios registrados, como Product Builder (falla)`
- [ ] T054 [US6] GREEN: servicio internal/service/listar_usuarios.go; handler internal/http/usuarios_handler.go registrado como `GET /users` (nivel «solo sesión») en internal/http/rutas.go y en la tabla de rutas_permisos_test.go; `domain.ErrUsuarioNoEncontrado` en internal/domain/errors.go mapeado a 404 en `escribirError`; internal/service/asignar_integrante.go deja de crear usuarios (busca por nombre normalizado y exige contraseña, vía un método de lectura del puerto de credenciales o de usuarios ya existente), `rol` vacío = `product_builder`; eliminar `resolverUsuario` de internal/service/usuarios.go; ajustar features/steps_proyecto.go, steps_scrum_master_unico.go, creacion_proyecto_equipo.feature y scrum_master_unico.feature (los integrantes se registran antes con «Dado que inicié sesión como "X"» o un paso equivalente que use `/auth/registro`) — `go test ./...` en verde
- [ ] T055 [US6] COMMIT `GREEN: equipo solo con usuarios registrados, como Product Builder, test en verde`

---

## Phase 9: BDD de HU-21 (escenarios de la historia)

El texto de la feature se redacta en el paso BDD del flujo (después de `/speckit.analyze`, antes de `/speckit.implement`, con el prompt de referencia de AGENTS.md) y se revisa con el equipo; se commitea como RED recién en esta fase porque Godog corre en modo `Strict` y unos pasos sin definir romperían el check `backend` de las fases anteriores (mismo criterio que HU-20, spec 034).

- [ ] T056 RED: agregar features/permisos_por_rol.feature (`@HU-21`, español) con la Característica Como/Quiero/Para de la spec y **un Escenario por criterio de aceptación (1 a 9)**, sin inventar casos; Antecedentes con «Dado que inicié sesión como "ana"» y un proyecto de ana con «bruno» como Product Builder; cubrir caso normal (1, 5), alternativo (8, 9), límite (6, 7) y error (2, 3, 4). Registrar en features/steps_permisos.go solo un stub de los pasos nuevos (sin implementar), de modo que el paquete compile; `go test ./features/` debe FALLAR por pasos sin definir (Strict), no por compilación
- [ ] T057 COMMIT `RED: agregar escenarios BDD de permisos por rol (fallan)`
- [ ] T058 GREEN: implementar los pasos en features/steps_permisos.go usando solo la API real (`RegistrarRutas` + `peticion()` con token; usuarios con `/auth/registro`, sin atajos de prueba); `go test ./features/` en verde
- [ ] T059 COMMIT `GREEN: agregar pasos BDD de permisos por rol, escenarios en verde`

---

## Phase 10: User Story 7 - Pantallas según el rol (Priority: P2, frontend)

**Goal**: la interfaz refleja los permisos (FR-014, FR-022 a FR-026).

**Independent Test**: pruebas de componente con el escenario en el nombre (Vitest + Testing Library).

### Cliente y permisos compartidos

- [ ] T060 [P] [US7] RED: web/src/compartido/api/client.test.ts — un 403 con mensaje «No sos integrante de este proyecto» se distingue (por ejemplo `error.tipo === 'no_integrante'`) de cualquier otro 403 (`'sin_permiso'`, conserva el mensaje del servidor); web/src/compartido/permisos/permisos.test.tsx (nuevo) — `NoIntegrante` muestra «No sos integrante de este proyecto» con enlace a «Mis proyectos»; `AvisoSinPermiso` muestra «No tenés permiso» y el mensaje; `useEsScrumMaster(proyecto)` según `proyecto.rol` — confirmar que FALLA
- [ ] T061 [US7] COMMIT `RED: 403 de pertenencia y de rol en la interfaz (falla)`
- [ ] T062 [US7] GREEN: web/src/compartido/api/client.ts y web/src/compartido/permisos/ (NoIntegrante.tsx, AvisoSinPermiso.tsx, useEsScrumMaster.ts); el panel (web/src/features/proyecto/Panel.tsx) muestra `NoIntegrante` ante ese 403 — `npm test` en verde
- [ ] T063 [US7] COMMIT `GREEN: 403 de pertenencia y de rol en la interfaz, test en verde`

### Mis proyectos

- [ ] T064 [P] [US7] RED: web/src/compartido/proyectos/ListaProyectos.test.tsx — «Escenario 5: veo solo mis proyectos con mi rol» (muestra «Scrum Master» / «Product Builder»); «Escenario 9: sin proyectos veo el mensaje y puedo crear uno» — confirmar que FALLA
- [ ] T065 [US7] COMMIT `RED: «Mis proyectos» con rol y estado vacío (falla)`
- [ ] T066 [US7] GREEN: tipo de proyecto con `rol` en web/src/compartido/proyectos/api.ts; web/src/compartido/proyectos/ListaProyectos.tsx muestra el rol y el estado vacío con el botón de crear — `npm test` en verde
- [ ] T067 [US7] COMMIT `GREEN: «Mis proyectos» con rol y estado vacío, test en verde`

### Identidad sin selectores

- [ ] T068 [P] [US7] RED: web/src/features/votacion/index.test.tsx — «Escenario 7: voto sin elegir quién actúa» (no hay selector «Quién actúa», el pedido de voto no lleva `usuario_id`, la consulta no lleva `?usuario_id`); web/src/features/sprints/FormularioHoras.test.tsx — sin selector de integrante, el pedido no lleva `integrante_id`; web/src/features/proyecto/crear.test.tsx — el pedido no lleva `creador` y no hay campo «Creador» — confirmar que FALLA
- [ ] T069 [US7] COMMIT `RED: la interfaz no elige quién actúa (falla)`
- [ ] T070 [US7] GREEN: eliminar web/src/features/votacion/SelectorActor.tsx y su uso (index.tsx, useVotacion.ts, api.ts, votacion.ts, Integrantes.tsx, Ronda.tsx); quitar el integrante de web/src/features/sprints/FormularioHoras.tsx y horas.ts; quitar `creador` de web/src/features/proyecto/crear.tsx, api.ts y validacion.ts — `npm test` en verde
- [ ] T071 [US7] COMMIT `GREEN: la interfaz no elige quién actúa, test en verde`

### Acciones de gestión deshabilitadas

- [ ] T072 [P] [US7] RED: pruebas «Escenario 2: como Product Builder veo la acción deshabilitada con la nota "Solo el Scrum Master puede…"» y «Escenario 1: como Scrum Master la acción está habilitada» en web/src/features/sprints/index.test.tsx (crear, iniciar — FormularioInicio — y cerrar — BotonCerrar), web/src/features/votacion/index.test.tsx (AccionesScrumMaster: iniciar, revelar, nueva ronda, registrar estimación), web/src/features/proyecto/Editar.test.tsx y Configuracion.test.tsx; y «un 403 de rol inesperado muestra "No tenés permiso" sin perder lo escrito» en Editar.test.tsx — confirmar que FALLA
- [ ] T073 [US7] COMMIT `RED: acciones de gestión deshabilitadas para el Product Builder (falla)`
- [ ] T074 [US7] GREEN: usar `useEsScrumMaster` y `AvisoSinPermiso` en web/src/features/sprints/ (SeccionSprints.tsx, FormularioInicio.tsx, BotonCerrar.tsx), web/src/features/votacion/AccionesScrumMaster.tsx, web/src/features/proyecto/Editar.tsx y Configuracion.tsx — `npm test` en verde
- [ ] T075 [US7] COMMIT `GREEN: acciones de gestión deshabilitadas para el Product Builder, test en verde`

### Configurar el equipo

- [ ] T076 [P] [US7] RED: web/src/features/proyecto/Configuracion.test.tsx — «Escenario 8: elijo de la lista un usuario registrado y entra como Product Builder» (carga `GET /users`, no hay selector de rol, el pedido manda `{"nombre"}`); muestra el error del servidor si el usuario no existe — confirmar que FALLA
- [ ] T077 [US7] COMMIT `RED: configurar el equipo eligiendo de la lista de usuarios (falla)`
- [ ] T078 [US7] GREEN: web/src/features/proyecto/Configuracion.tsx y api.ts (`listarUsuarios`) — `npm test` en verde
- [ ] T079 [US7] COMMIT `GREEN: configurar el equipo eligiendo de la lista de usuarios, test en verde`
- [ ] T080 [US7] (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

---

## Phase 11: Polish & Cross-Cutting Concerns

- [ ] T081 Actualizar MODELO-DE-DATOS.md **solo texto, sin esquema**: §4.2 (`contrasena_hash` `NULL` = usuario previo; desde HU-21 ninguna ruta crea usuarios sin contraseña) y §8 (fila «HU-21 Permisos por rol | lee `integrantes`, `usuarios`, `sesiones` | — | sin migración»)
- [ ] T082 Agregar en docs/PRODUCT-BACKLOG.md una nota en las secciones de HU-02, HU-07 y HU-14 (en «Notas», con el mismo estilo de la nota que ya tiene HU-07 sobre su escenario 1): desde HU-21, un no integrante recibe 403 «No sos integrante de este proyecto» antes que el error de rol, y en las horas el rechazo pasa de 400 sobre `integrante_id` a ese 403 (cambio confirmado por el equipo el 2026-11-11, spec 042, sección Clarifications)
- [ ] T083 Declarar el **límite conocido** (con sesión, un no integrante puede leer historias, Sprints, horas, defectos, métricas, historial y reportes de un proyecto) y que el traspaso del Scrum Master queda fuera de alcance, en README.md o en la documentación de la API que corresponda; actualizar los `curl` de docs/ que manden `usuario_id`, `integrante_id` o `creador`
- [ ] T084 COMMIT `docs: permisos por rol en el modelo de datos, el backlog y el límite conocido de lectura (Refs #234)`
- [ ] T085 `gofmt -l .` vacío, `go build ./...`, `go vet ./...`, `go test ./...`; `cd web && npm run lint && npm test && npm run build` — todo en verde
- [ ] T086 Recorrer [quickstart.md](quickstart.md) a mano (curl e interfaz) y anotar el resultado en el PR

---

## Dependencies & Execution Order

- **Phase 1 → Phase 2** (bloqueante) → **US3** (T013–T016: todas las rutas pasan a exigir sesión) → **US1** → **US4** → **US2** → **US5** → **US6** → **BDD HU-21** → **US7** → **Polish**.
- US1, US4 y US2 tocan los mismos servicios y steps de Godog: hacerlas en ese orden, no en paralelo.
- US5 y US6 son independientes entre sí (repositorios distintos) una vez terminada US2.
- La fase BDD (T056–T059) necesita US1 a US6 terminadas (sus escenarios cubren los 9 criterios).
- US7 (frontend) depende solo del contrato ([contracts/api-permisos.md](contracts/api-permisos.md)); se puede empezar en paralelo por otra persona después de la Phase 2, con mocks del cliente.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.

### Parallel Opportunities

- Las RED marcadas [P] de una misma story tocan archivos de test distintos y se pueden escribir a la vez, pero sus GREEN se integran en orden.
- T040 (repositorio de proyectos) y T048 (repositorio de usuarios).
- Frontend: T060, T064, T068, T072 y T076 tocan pruebas distintas.

## Implementation Strategy

1. Phase 1 + 2 → reglas y sesión listas, sin cambios visibles.
2. US3 + US1 → MVP de seguridad: nada sin sesión y nadie actúa por otro (escenarios 3 y 7).
3. US4 + US2 → segregación de funciones (escenarios 1, 2 y 4). **Plan B del backlog**: estas
   fases (y US3/US1) son las que deben estar listas el 13-nov; si no, se declara.
4. US5 + US6 + BDD → escenarios 5, 6, 8 y 9, y la feature de Godog completa.
5. US7 → pantallas; Polish → documentación y validación.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- Mensajes de commit con `(Refs #234)` cuando corresponda.
