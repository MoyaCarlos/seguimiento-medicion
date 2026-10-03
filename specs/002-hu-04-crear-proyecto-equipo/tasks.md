---

description: "Task list for HU-04 — Creación de Proyecto y Asignación de Equipo"
---

# Tasks: Creación de Proyecto y Asignación de Equipo (HU-04)

**Input**: Design documents from `/specs/002-hu-04-crear-proyecto-equipo/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/openapi.yaml](./contracts/openapi.yaml), [quickstart.md](./quickstart.md)

**Tests**: SÍ son obligatorias. `AGENTS.md` exige TDD (RED→GREEN→REFACTOR, con cada fase en su propio commit) para dominio/servicio y BDD (Godog) por criterio de aceptación, así que las tareas de test se incluyen y se ejecutan antes de implementar.

**Organization**: Tareas agrupadas por historia de usuario (US1=P1, US2=P2, US3=P3) para permitir implementación y prueba independientes.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Puede correr en paralelo (archivo distinto, sin dependencias pendientes)
- **[Story]**: `[US1]`, `[US2]` o `[US3]`, según `spec.md`
- Rutas de archivo exactas en cada tarea
- Convención de commits (AGENTS.md): `RED: ...`, `GREEN: ...`, `REFACTOR: ...`

## Path Conventions

Backend Go en la raíz: `cmd/`, `internal/`, `features/` (ver `plan.md`). Sin frontend en esta iteración (R12): la redirección y las vistas del SPA quedan diferidas.

## Restricciones de `data-model.md` (citas literales)

- `Project.Name`: no vacío (ni solo espacios), **≤ 100 caracteres**; nombre repetido permitido.
- `Project.Description`: **≤ 2000 caracteres**, opcional.
- Fechas `StartDate`/`EndDate`: `*time.Time` opcionales; si ambas existen, `EndDate` **no puede ser anterior** a `StartDate`.
- `User.Name`: no vacío (ni solo espacios), **≤ 200 caracteres**; `NormalizedName = ToLower(TrimSpace(Name))` único.
- `Role`: `scrum_master` o `product_builder` (único por integrante y proyecto).
- Firmas de HU-05 que NO se cambian: `Create(p *domain.Project) error`, `GetByID(id string) (*domain.Project, error)`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirmar dependencias y estructura base (no se agregan dependencias nuevas)

- [X] T001 Confirmar que `modernc.org/sqlite` y `github.com/cucumber/godog` ya están en `go.mod` y ejecutar `go mod tidy` (actualiza `go.mod`/`go.sum` si hace falta)
- [X] T002 [P] Verificar que `.gitignore` ignora la base local (`*.db`, `*.db-journal`); agregar la regla si falta

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Tipos de dominio compartidos (Role/User/Membership), extensión de `Project`, puertos, adaptadores SQLite, generación de IDs y esquema que usan US1, US2 y US3

**⚠️ CRITICAL**: Ninguna historia puede empezar hasta completar esta fase

- [X] T003 [P] Test de dominio (RED) de `Role`, `User` y `Membership` en `internal/domain/member_test.go`: rol solo acepta `scrum_master`/`product_builder`; `User.Name` vacío o solo espacios → `ValidationError{Campo:"nombre"}`; `User.Name` > 200 runas → error de `nombre`; `NormalizedName` = minúsculas + sin espacios externos; `Membership` con `projectID`/`userID` vacíos → error del campo correspondiente
- [X] T004 Implementar el enum `Role` con `RolScrumMaster = "scrum_master"`, `RolProductBuilder = "product_builder"` y `Valido() bool` en `internal/domain/member.go` (GREEN, depende de T003)
- [X] T005 Implementar el struct `User` y el constructor `NewUser(name string) (User, error)` (recorta espacios, valida ≤ 200 caracteres, calcula `NormalizedName = ToLower(TrimSpace(Name))`) en `internal/domain/member.go` (depende de T004) — hace pasar la parte de `User` de T003
- [X] T006 Implementar el struct `Membership` y el constructor `NewMembership(projectID, userID string, role Role) (Membership, error)` en `internal/domain/member.go` (depende de T004) — hace pasar la parte de `Membership` de T003
- [X] T007 [P] Extender el struct `Project` en `internal/domain/project.go` con `StartDate *time.Time` y `EndDate *time.Time`, preservando `ID`, `Name`, `Description`, `CreatedAt` y su orden
- [X] T008 Extender la interfaz `ProjectRepository` en `internal/repository/project_repository.go` con `Update(p *domain.Project) error`, `AddMember(m *domain.Membership) error` y `ListMembers(projectID string) ([]domain.Membership, error)`, manteniendo sin cambios `Create(p *domain.Project) error` y `GetByID(id string) (*domain.Project, error)`
- [X] T009 [P] Definir el puerto `UserRepository` con `FindByNormalizedName(normalizedName string) (*domain.User, error)` y `Create(u *domain.User) error` en `internal/repository/user_repository.go`
- [X] T010 [P] Implementar el generador de identificadores UUID v4 con `crypto/rand` (`nuevoID() string`) en `internal/repository/id.go`
- [X] T011 Extender `Migrar` en `internal/repository/migrate.go` con `CREATE TABLE IF NOT EXISTS` para `projects` (`id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT NOT NULL DEFAULT '', start_date TEXT, end_date TEXT, created_at TEXT NOT NULL`), `users` (`id TEXT PRIMARY KEY, name TEXT NOT NULL, normalized_name TEXT NOT NULL UNIQUE, created_at TEXT NOT NULL`) y `project_members` (`project_id TEXT NOT NULL, user_id TEXT NOT NULL, role TEXT NOT NULL, created_at TEXT NOT NULL, PRIMARY KEY (project_id, user_id), FOREIGN KEY (project_id) REFERENCES projects(id), FOREIGN KEY (user_id) REFERENCES users(id)`), y activar `PRAGMA foreign_keys = ON` al abrir la conexión
- [X] T012 [P] Test de integración (RED) de `SQLiteUserRepository` en `internal/repository/sqlite_user_test.go` (SQLite `:memory:`): `Create` asigna `ID` y persiste `NormalizedName`; `FindByNormalizedName` encuentra por minúsculas+trim; nombre normalizado duplicado no crea un segundo usuario
- [X] T013 [P] Test de integración (RED) de `SQLiteProjectRepository` en `internal/repository/sqlite_project_test.go` (SQLite `:memory:`): `Create` asigna `ID`; `GetByID` relee con fechas `NULL`/RFC3339; `Update` persiste cambios; `AddMember` + `ListMembers` devuelven la relación Proyecto–Usuario–Rol; `AddMember` duplicado falla por la PK; dos proyectos con el **mismo nombre** se persisten ambos (FR-007)
- [X] T014 Implementar `SQLiteUserRepository` en `internal/repository/sqlite_user.go` (GREEN, depende de T009, T010, T011) — hace pasar T012
- [X] T015 Implementar `SQLiteProjectRepository` en `internal/repository/sqlite_project.go` (GREEN, depende de T008, T010, T011) — hace pasar T013

**Checkpoint**: Dominio base, puertos, adaptadores y esquema listos — US1, US2 y US3 pueden comenzar

---

## Phase 3: User Story 1 - Crear el entorno del proyecto (Priority: P1) 🎯 MVP

**Goal**: Un Scrum Master crea un proyecto con nombre obligatorio, descripción y fechas opcionales; el proyecto se persiste, el creador queda vinculado como `scrum_master` y la API confirma la creación devolviendo el proyecto (la redirección del SPA queda diferida).

**Independent Test**: `POST /projects` con datos válidos → `201` con `id` (UUID), fechas y `GET /projects/{id}/members` mostrando al creador como `scrum_master`; `GET /projects/{id}` devuelve el proyecto; nombre vacío → `400 {"campo":"nombre",...}` y sin registro.

### Tests for User Story 1 (TDD — escribir y ver FALLAR antes de implementar) ⚠️

- [X] T016 [P] [US1] Test de dominio (RED) de `NewProject` en `internal/domain/project_test.go`: válido con y sin fechas; `name` vacío y solo espacios → `ValidationError{Campo:"nombre"}`; `name` > 100 runas → error de `nombre`; `description` > 2000 runas → error de `descripcion`; `end` anterior a `start` → `ValidationError{Campo:"fecha_fin"}`; `Name`/`Description` recortados; dos proyectos con el mismo nombre son válidos (FR-007)
- [X] T017 [P] [US1] Test de servicio (RED) de `CrearProyecto` en `internal/service/crear_proyecto_test.go` con fakes en memoria de `ProjectRepository` y `UserRepository` (definidos en el mismo archivo): crea y persiste el proyecto y vincula al `creador` como `scrum_master`; entrada inválida devuelve `ValidationError` y NO llama a `Create` ni `AddMember`
- [X] T018 [P] [US1] Test del handler (RED) en `internal/http/project_handler_test.go` (`httptest` + fakes): `POST /projects` → `201` con el proyecto creado; `400 {"campo":"nombre",...}` con nombre vacío; `400 {"campo":"fecha_fin",...}` con fin < inicio; `400` con JSON malformado; `500` ante fallo de persistencia; `GET /projects/{id}` → `200` con el proyecto y `404` si no existe
- [X] T022 [P] [US1] Test de servicio (RED) de `ObtenerProyecto` en `internal/service/obtener_proyecto_test.go`: devuelve el proyecto existente por `id`; `id` inexistente → error no encontrado

### Implementation for User Story 1

- [X] T019 [US1] Implementar `NewProject(name, description string, start, end *time.Time) (Project, error)` en `internal/domain/project.go` con validaciones en orden: nombre no vacío, nombre ≤ 100 runas, descripción ≤ 2000 runas, y `end` no anterior a `start` cuando ambas existen (depende de T007, T016) — hace pasar T016
- [X] T020 [US1] Implementar `CrearProyecto` en `internal/service/crear_proyecto.go` con `CrearProyectoInput{Nombre, Descripcion string, FechaInicio, FechaFin *time.Time, Creador string}` y `Ejecutar(ctx, input)`: construye con `domain.NewProject`, persiste con `ProjectRepository.Create`, resuelve/crea el `User` del `Creador` vía `UserRepository` y crea la `Membership` con rol `scrum_master` (FR-008; depende de T006, T008, T009, T019) — hace pasar T017
- [X] T021 [US1] Implementar los DTOs de creación y el handler `POST /projects` en `internal/http/project_dto.go` y `internal/http/project_handler.go`: parsear fechas `YYYY-MM-DD`, mapear `domain.ValidationError` (vía `errors.As`) y JSON inválido a `400 {"campo","mensaje"}`, devolver `201` con el proyecto y `500` en errores de persistencia (depende de T020) — hace pasar la parte POST de T018
- [X] T023 [US1] Implementar `ObtenerProyecto` en `internal/service/obtener_proyecto.go` y el handler `GET /projects/{id}` en `internal/http/project_handler.go` (usar `r.PathValue("id")`), devolviendo `200` con el proyecto y `404` si no existe (depende de T008, T022) — hace pasar T022 y la parte GET de T018
- [X] T024 [US1] Registrar los step definitions de proyecto en `features/features_test.go` (extender `InitializeScenario`) y crear `features/creacion_proyecto_equipo.feature` con los escenarios de US1 (creación válida con creador Scrum Master; nombre vacío → advertencia y sin registro) más su implementación en `features/steps_proyecto.go` (depende de T021, T023)
- [X] T025 [US1] Cablear en `cmd/api/main.go`: construir `SQLiteUserRepository` y `SQLiteProjectRepository`, `CrearProyecto`/`ObtenerProyecto` → handlers y registrar `mux.HandleFunc("POST /projects", ...)` y `mux.HandleFunc("GET /projects/{id}", ...)` (depende de T014, T015, T021, T023)

**Checkpoint**: US1 funcional y testeable de forma independiente (MVP de proyectos)

---

## Phase 4: User Story 2 - Asignar integrantes y roles (Priority: P2)

**Goal**: Desde la configuración de un proyecto existente se vincula a un integrante indicando nombre y rol (`scrum_master` o `product_builder`); el usuario se reutiliza por nombre normalizado y se rechaza el duplicado en el mismo proyecto.

**Independent Test**: `POST /projects/{id}/members` con `{"nombre":"Jimena","rol":"product_builder"}` → `201`; el mismo nombre con distintas mayúsculas/espacios → `400` por duplicado; rol inválido → `400 {"campo":"rol",...}`; proyecto inexistente → `404`; `GET /projects/{id}/members` lista el equipo.

### Tests for User Story 2 (TDD — escribir y ver FALLAR antes de implementar) ⚠️

- [X] T026 [P] [US2] Test de servicio (RED) de `AsignarIntegrante` en `internal/service/asignar_integrante_test.go`: vincula un integrante nuevo con rol válido; reutiliza el `User` existente cuando el nombre difiere solo en mayúsculas/espacios; rechaza duplicado en el mismo proyecto con `ValidationError{Campo:"integrante"}`; rol inválido devuelve `ValidationError{Campo:"rol"}`; proyecto inexistente no persiste; dos integrantes distintos con rol `scrum_master` en el mismo proyecto se permiten (FR-017)
- [X] T027 [P] [US2] Test del handler (RED) de `POST /projects/{id}/members` y `GET /projects/{id}/members` en `internal/http/project_handler_test.go`: `201` con `{id, proyecto_id, nombre, rol}`; `400` con rol inválido y con integrante duplicado; `404` con proyecto inexistente; `200` con el listado

### Implementation for User Story 2

- [X] T028 [US2] Implementar `AsignarIntegrante` en `internal/service/asignar_integrante.go` con `AsignarIntegranteInput{ProyectoID, Nombre string, Rol domain.Role}` y `Ejecutar(ctx, input)`: verificar el proyecto con `GetByID`, resolver/crear el `User` por nombre normalizado vía `UserRepository`, y persistir la `Membership` con `AddMember`; mapear la violación de la PK (proyecto, usuario) a `ValidationError{Campo:"integrante", Mensaje:"el integrante ya pertenece al proyecto"}` (FR-009..FR-016; depende de T006, T008, T009) — hace pasar T026
- [X] T029 [US2] Implementar los DTOs y los handlers `POST /projects/{id}/members` y `GET /projects/{id}/members` en `internal/http/project_dto.go` y `internal/http/project_handler.go` (usar `r.PathValue("id")`), con mapeo de `404` a proyecto inexistente (depende de T028) — hace pasar T027
- [X] T030 [US2] Agregar a `features/creacion_proyecto_equipo.feature` y `features/steps_proyecto.go` los escenarios de asignación (alta con rol, rol inválido, duplicado) (depende de T029)
- [X] T031 [US2] Cablear `POST /projects/{id}/members` y `GET /projects/{id}/members` en `cmd/api/main.go` (depende de T029)

**Checkpoint**: US1 y US2 funcionan de forma independiente

---

## Phase 5: User Story 3 - Editar un proyecto existente (Priority: P3)

**Goal**: Un Scrum Master corrige nombre, descripción y fechas de un proyecto existente reutilizando las validaciones del alta, sin versionado ni historial; los datos inválidos se rechazan conservando los valores previos.

**Independent Test**: `PUT /projects/{id}` con datos válidos → `200` y `GET /projects/{id}` refleja los cambios; nombre vacío o fin < inicio → `400 {"campo",...}` sin modificar el proyecto; proyecto inexistente → `404`.

### Tests for User Story 3 (TDD — escribir y ver FALLAR antes de implementar) ⚠️

- [X] T032 [P] [US3] Test de dominio (RED) de `ConDatosEditados` en `internal/domain/project_test.go`: actualiza nombre/descripción/fechas válidos y conserva `ID` y `CreatedAt`; reutiliza las mismas validaciones del alta (nombre vacío, > 100, fin < inicio, descripción > 2000); descripción vacía permitida
- [X] T033 [P] [US3] Test de servicio (RED) de `EditarProyecto` en `internal/service/editar_proyecto_test.go`: carga por `id`, valida y persiste con `Update`; proyecto inexistente devuelve error no encontrado; entrada inválida NO llama a `Update`
- [X] T034 [P] [US3] Test del handler (RED) de `PUT /projects/{id}` en `internal/http/project_handler_test.go`: `200` con el proyecto actualizado; `400 {"campo","mensaje"}` con nombre vacío o fin < inicio; `404` con proyecto inexistente

### Implementation for User Story 3

- [X] T035 [US3] Implementar `(p Project) ConDatosEditados(name, description string, start, end *time.Time) (Project, error)` en `internal/domain/project.go`, reutilizando la misma validación de `NewProject` y conservando `ID` y `CreatedAt` (depende de T019, T032) — hace pasar T032
- [X] T036 [US3] Implementar `EditarProyecto` en `internal/service/editar_proyecto.go` con `EditarProyectoInput{ID, Nombre, Descripcion string, FechaInicio, FechaFin *time.Time}` y `Ejecutar(ctx, input)`: `GetByID` + `ConDatosEditados` + `Update`; error no encontrado si no existe (depende de T008, T035) — hace pasar T033
- [X] T037 [US3] Implementar el DTO `EditarProyectoRequest` y el handler `PUT /projects/{id}` en `internal/http/project_dto.go` y `internal/http/project_handler.go` (depende de T036) — hace pasar T034
- [X] T038 [US3] Agregar a `features/creacion_proyecto_equipo.feature` y `features/steps_proyecto.go` los escenarios de edición (edición válida; nombre vacío; fin < inicio; cancelar sin cambios) (depende de T037)
- [X] T039 [US3] Cablear `PUT /projects/{id}` en `cmd/api/main.go` (depende de T037)

**Checkpoint**: US1, US2 y US3 funcionan de forma independiente

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Calidad, consistencia y validación end-to-end

- [X] T040 [P] Ejecutar `gofmt -l .` y `go vet ./...` y corregir hallazgos
- [X] T041 Ejecutar la suite completa `go test ./...` y el BDD `go test ./features/...`; asegurar todo en verde y que US1/US2/US3 sean independientemente testeables
- [X] T042 Ejecutar la validación manual de `quickstart.md` (los escenarios curl de creación, asignación y edición, más los casos borde) y registrar evidencia
- [X] T043 [P] Verificar que las respuestas reales cumplen `contracts/openapi.yaml` (formas `201`/`200`/`400`/`404`/`500`; enums `scrum_master`/`product_builder`; fechas `null` cuando corresponden)
- [X] T044 Verificar que no se agregaron dependencias nuevas, que las firmas `ProjectRepository.Create(p *domain.Project) error` y `GetByID(id string) (*domain.Project, error)` siguen idénticas (contrato con HU-05) y que no existe ningún endpoint de versiones/historial (FR-023)

---

## Phase 7: Remediación del code review

**Input**: `Informe-remediacion-HU-04.md` (revisión externa de calidad).

**Reglas**: un commit por fase (`RED:`/`GREEN:`/`REFACTOR:`) y, tras cada GREEN/REFACTOR,
`go build ./...`, `go vet ./...`, `gofmt -l .` (vacío) y `go test ./...` en verde. Ningún
commit de esta fase puede romper `go build ./...`.

**Desvío de proceso justificado (aceptado por el Agile Enabler)**: cambiar las firmas
`string→int64` no admite un test RED que compile contra el código viejo. Por eso la
migración de IDs se hace como **commit mecánico por capa** (`REFACTOR: pasar IDs a int64
en <capa>`) que solo adapta tipos y tests existentes, sin comportamiento nuevo; el
RED/GREEN real se reserva para lo observable nuevo (ID `> 0` por `LastInsertId`, `400`
por id no numérico o `<= 0`, IDs numéricos en el JSON).

**Decisiones de diseño (propuestas en el plan, sin objeción del Agile Enabler)**:
- Método transaccional (B2): `GuardarConScrumMaster(ctx, p domain.Project, creadorID int64) (domain.Project, error)`.
- Orden de validación (mejora): **404 primero** (proyecto inexistente domina sobre rol inválido), documentado en `spec.md` (casos límite).
- Usuario huérfano (B2): `resolverUsuario` (find-or-create del creador) queda fuera de la transacción; un usuario sin membresía es inocuo (no viola SC-006/SC-007 y se reutiliza luego).

### 7.1 IDs `int64` (sección 1 del informe)

- [ ] R001 REFACTOR dominio: `Project.ID`, `User.ID`, `Membership.ProjectID`, `Membership.UserID` → `int64`; agregar `ErrProyectoNoEncontrado` en `internal/domain/errors.go`; adaptar `internal/domain/project_test.go` y `member_test.go`. Verde.
- [ ] R002 COMMIT `REFACTOR: pasar IDs a int64 en dominio`
- [ ] R003 REFACTOR repositorio: `ProjectRepository` → `Guardar(ctx, Project) (Project, error)`, `ObtenerPorID(ctx, int64) (Project, error)`, `Actualizar`, `AgregarIntegrante`, `ListarIntegrantes`, `GuardarConScrumMaster(ctx, Project, creadorID int64) (Project, error)` (ctx primero); `UserRepository` int64; borrar `internal/repository/id.go` y todo `crypto/rand`; migración con `projects.id`/`users.id` `INTEGER PRIMARY KEY AUTOINCREMENT` y `project_members.project_id`/`user_id` `INTEGER` con FK; adaptar `sqlite_project.go`, `sqlite_user.go` y sus tests (ID por `LastInsertId()`). Verde.
- [ ] R004 COMMIT `REFACTOR: pasar IDs a int64 en repositorio`
- [ ] R005 RED repo: test en `internal/repository/sqlite_project_test.go` de que `Guardar` asigna `ID > 0` y `ObtenerPorID(id int64)` relee por ese entero — confirmar FALLA.
- [ ] R006 COMMIT `RED: test de ID numérico y lectura por ID (falla)`
- [ ] R007 GREEN repo: implementar `Guardar`/`ObtenerPorID` int64 con `LastInsertId()` — verde.
- [ ] R008 COMMIT `GREEN: ID numérico asignado por LastInsertId, test en verde`
- [ ] R009 REFACTOR service: adaptar `CrearProyecto`, `ObtenerProyecto`, `EditarProyecto`, `AsignarIntegrante`, `ListarIntegrantes`, `usuarios.go` y `fakes_test.go` + tests al puerto int64. Verde.
- [ ] R010 COMMIT `REFACTOR: pasar IDs a int64 en service`
- [ ] R011 REFACTOR http: `{id}` de ruta como `int64`; DTOs con IDs `int64`; `400 {"campo":"id",...}` si no numérico o `<= 0`; adaptar `internal/http/project_dto.go`, `project_handler.go` y tests. Verde.
- [ ] R012 COMMIT `REFACTOR: pasar IDs a int64 en http`
- [ ] R013 RED http: test `GET/PUT/POST /projects/{id}...` con id no numérico y con id `<= 0` → `400 {"campo":"id",...}` — confirmar FALLA.
- [ ] R014 COMMIT `RED: 400 por id no numérico o <= 0 (falla)`
- [ ] R015 GREEN http: parsear `{id}` a int64 con 400 — verde.
- [ ] R016 COMMIT `GREEN: parsear id a int64 con 400, test en verde`
- [ ] R017 Agregar `NewProjectRepositoryEnMemoria()` en `internal/repository` (archivo no-test, implementa la interfaz completa) para HU-05/HU-13.
- [ ] R018 COMMIT `Agregar NewProjectRepositoryEnMemoria`
- [ ] R019 Actualizar `spec.md`, `plan.md`, `data-model.md`, `quickstart.md`, `research.md` y `contracts/openapi.yaml` (IDs `integer`, contrato int64). Borrar `seguimiento.db` local (esquema TEXT→INTEGER).
- [ ] R020 COMMIT `Actualizar specs y OpenAPI a IDs int64`

### 7.2 B2 — creación transaccional proyecto + Scrum Master

- [ ] R021 RED integración: test en `internal/repository/sqlite_project_test.go` sobre `:memory:` que fuerza fallo en la 2ª escritura (`GuardarConScrumMaster` con `creadorID` inexistente y FK activas) y verifica `SELECT COUNT(*) FROM projects` == 0 — confirmar FALLA.
- [ ] R022 COMMIT `RED: proyecto + Scrum Master en una sola transacción (falla)`
- [ ] R023 GREEN: `GuardarConScrumMaster` con `BeginTx`/`Commit`/`Rollback` en `SQLiteProjectRepository`; `CrearProyecto` usa solo ese método (creador resuelto vía `resolverUsuario`, huérfano inocuo justificado). Verde.
- [ ] R024 COMMIT `GREEN: crear proyecto + Scrum Master transaccional, test en verde`

### 7.3 B1 — runner Godog único (11 escenarios)

- [ ] R025 Unificar runner en `features/features_test.go`: `Paths:["."]`, `Name:"BDD"`, registrando HU-01 y HU-04 (cada una en su `inicializarPasosXxx`); desambiguar el regex colisionado de HU-04 (`^el sistema muestra una advertencia de validación$` → `... del proyecto`) y actualizar `features/creacion_proyecto_equipo.feature` sin salir del spec. `go test ./features/ -v` → **11 escenarios**.
- [ ] R026 COMMIT `Unificar runner BDD con los pasos de HU-01 y HU-04`
- [ ] R027 Documentar en `features/README.md` cómo se agrega una historia nueva.
- [ ] R028 COMMIT `Documentar agregado de historias BDD en features/README.md`

### 7.4 C1 — pasos BDD que verifican lo que dicen

- [ ] R029 Mejorar `features/steps_proyecto.go`: `noRealizaVinculacion` comprueba que no se escribió fila nueva en `project_members` (p.ej. 0 filas con rol `product_builder`); `asignarInvalido` ejercita nombre vacío **y** rol inválido; `editarValido` manda fechas y `persisteCambios` las verifica persistidas. Demostrar que cada `Entonces` falla rompiendo temporalmente una regla (sin commitear).
- [ ] R030 COMMIT `Reforzar pasos BDD de HU-04 para verificar lo que dicen`

### 7.5 C2 — PUT completo con presencia de claves

- [ ] R031 RED http: `PUT /projects/{id}` con `{"nombre":"X"}` (sin `descripcion`/`fecha_inicio`/`fecha_fin`) sobre proyecto con descripción y fechas → `400` y proyecto intacto — confirmar FALLA.
- [ ] R032 COMMIT `RED: PUT sin claves obligatorias devuelve 400 (falla)`
- [ ] R033 GREEN http: DTO de edición con detección de presencia de clave (punteros/`json.RawMessage`); `nombre`, `descripcion`, `fecha_inicio`, `fecha_fin` obligatorios como claves; `null` explícito para limpiar fechas; clave ausente → `400` con el campo. Actualizar `contracts/openapi.yaml`. Verde.
- [ ] R034 COMMIT `GREEN: PUT completo con presencia de claves, test en verde`

### 7.6 C3 — campo del error del creador

- [ ] R035 RED service: `CrearProyecto` con `creador` vacío → `ValidationError{Campo:"creador"}` (hoy dice `"nombre"`) — confirmar FALLA.
- [ ] R036 COMMIT `RED: error del creador con campo correcto (falla)`
- [ ] R037 GREEN service: mapear a `"creador"` solo en el flujo de creación (sin romper `AsignarIntegrante`, que sigue con `"nombre"`). Verde.
- [ ] R038 COMMIT `GREEN: campo creador en error del creador, test en verde`

### 7.7 C4 — foreign_keys por DSN

- [ ] R039 RED repo: test que abre dos conexiones del pool sobre un archivo en `t.TempDir()` y verifica `PRAGMA foreign_keys` == 1 en ambas — confirmar FALLA.
- [ ] R040 COMMIT `RED: foreign_keys activas en todas las conexiones (falla)`
- [ ] R041 GREEN repo: activar por DSN `?_pragma=foreign_keys(1)` en `AbrirSQLite`; mantener las pruebas `:memory:` funcionando (cap de 1 conexión solo para memoria). Verde.
- [ ] R042 COMMIT `GREEN: foreign_keys por DSN, test en verde`

### 7.8 Mejoras requeridas

- [ ] R043 RED repo: `AddMember` con `user_id` inexistente (FK) → error de integridad, NO `ErrProyectoNoEncontrado` (404) — confirmar FALLA.
- [ ] R044 COMMIT `RED: falla de FK no se confunde con proyecto no encontrado (falla)`
- [ ] R045 GREEN repo: reemplazar `strings.Contains` por el código de error de modernc en `sqlite_project.go` y `sqlite_user.go`, distinguiendo UNIQUE / FK / no-encontrado. Verde.
- [ ] R046 COMMIT `GREEN: errores SQLite por código (FK ≠ 404), test en verde`
- [ ] R047 REFACTOR: un solo `ErrNoEncontrado` de dominio (eliminar el de `repository` y la traducción repetida en los 4 services). Verde.
- [ ] R048 COMMIT `REFACTOR: un solo error de dominio para no-encontrado`
- [ ] R049 REFACTOR: validar el rol una sola vez (eliminar el duplicado entre `AsignarIntegrante` y `NewMembership`). Verde.
- [ ] R050 COMMIT `REFACTOR: validación de rol en un único lugar`
- [ ] R051 RED service: proyecto inexistente + rol inválido → `404` (orden: proyecto primero) — confirmar FALLA.
- [ ] R052 COMMIT `RED: orden de validación 404 antes que 400 (falla)`
- [ ] R053 GREEN service + spec: aplicar el orden 404→400 y documentarlo como caso límite en `spec.md`. Verde.
- [ ] R054 COMMIT `GREEN: orden de validación 404 antes que 400, test en verde`
- [ ] R055 RED repo: `ListMembers` con 3 altas en el mismo segundo devuelve el orden real de alta (no arbitrario) — confirmar FALLA.
- [ ] R056 COMMIT `RED: ListMembers ordena por orden real de alta (falla)`
- [ ] R057 GREEN repo: ordenar por el orden real de alta (p.ej. `ORDER BY pm.rowid`). Verde.
- [ ] R058 COMMIT `GREEN: ListMembers por orden de alta, test en verde`

### 7.9 Renombres a convención español (REFACTOR, tests en verde)

- [ ] R059 REFACTOR dominio: `Name→Nombre`, `Description→Descripcion`, `StartDate→FechaInicio`, `EndDate→FechaFin`, `CreatedAt→CreadoEn`; `User.Name→Nombre`, `NormalizedName→NombreNormalizado`, `CreatedAt→CreadoEn`. Adaptar services/repos/http/tests. Verde.
- [ ] R060 COMMIT `REFACTOR: campos de dominio a español (Project/User)`
- [ ] R061 REFACTOR puertos: `Create→Guardar`, `GetByID→ObtenerPorID`, `Update→Actualizar`, `AddMember→AgregarIntegrante`, `ListMembers→ListarIntegrantes`, `FindByNormalizedName→ObtenerPorNombreNormalizado`. Adaptar services/tests. Verde.
- [ ] R062 COMMIT `REFACTOR: métodos de puerto a español`

### 7.10 Spec en EARS, Branch y /speckit.analyze

- [ ] R063 Reescribir reglas de negocio, casos límite y condiciones de error de `spec.md` en formato EARS; corregir el campo `Branch` con el nombre real de la rama.
- [ ] R064 COMMIT `Reescribir spec HU-04 en EARS y corregir Branch`
- [ ] R065 Ejecutar `/speckit.analyze` y registrar resultado y resolución en `research.md` (criterio: 0 CRITICAL).
- [ ] R066 COMMIT `Registrar resultado de /speckit.analyze en research.md`

---

## Excepción de proceso

Los **13 commits originales** de esta rama (previos a la remediación) no llevan prefijos
`RED:`/`GREEN:`/`REFACTOR:` y 8 de ellos no compilan por separado, porque la plantilla de
`/speckit.tasks` vigente en ese momento no exigía un commit por fase (corregida luego en el
PR #15). **No se reescribe el historial**: reconstruir un proceso que no ocurrió así es peor
que una excepción documentada. La evidencia del ciclo TDD queda cubierta de aquí en adelante
por los commits de esta fase de remediación.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: sin dependencias
- **Foundational (Phase 2)**: depende de Setup — BLOQUEA las tres historias
- **User Stories (Phase 3-5)**: dependen de Foundational. US1 puede empezar primero; US2 y US3 son incrementos independientes a nivel de caso de uso, aunque comparten `project_handler.go`/`project_dto.go` y el feature file (coordinar por archivo)
- **Polish (Phase 6)**: depende de completar las historias deseadas

### User Story Dependencies

- **US1 (P1)**: depende solo de Foundational. Entrega el MVP.
- **US2 (P2)**: depende de Foundational (User/Role/Membership y `UserRepository`); no depende del caso de uso de US1, pero comparte DTOs/handler.
- **US3 (P3)**: depende de Foundational y de `NewProject` de US1 (reutiliza su validación vía `ConDatosEditados`).

### Within Each User Story

- Tests (RED) antes de implementación (GREEN); REFACTOR solo si hay algo que limpiar
- Dominio → service → repository (foundational) → handler → BDD → wiring
- Commit por fase con prefijo `RED:`/`GREEN:`/`REFACTOR:` (AGENTS.md)

### Parallel Opportunities

- T002, T007, T009, T010 (archivos distintos) en paralelo
- T012 y T013 (tests de repositorios distintos) en paralelo
- T016, T017, T018 y T022 (tests de US1 en archivos distintos) en paralelo
- T026 y T027 (tests de US2) en paralelo
- T032, T033 y T034 (tests de US3) en paralelo
- T040 y T043 en paralelo en la fase final

---

## Parallel Example: User Story 1

```text
# Tests RED de US1 en paralelo (archivos distintos):
Task: "T016 [US1] Test de NewProject en internal/domain/project_test.go"
Task: "T017 [US1] Test de CrearProyecto en internal/service/crear_proyecto_test.go"
Task: "T018 [US1] Test de handlers POST/GET en internal/http/project_handler_test.go"
Task: "T022 [US1] Test de ObtenerProyecto en internal/service/obtener_proyecto_test.go"

# Foundation en paralelo:
Task: "T007 Project con fechas en internal/domain/project.go"
Task: "T009 UserRepository en internal/repository/user_repository.go"
Task: "T010 nuevoID en internal/repository/id.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1 Setup → Phase 2 Foundational
2. Phase 3 US1 (T016–T025)
3. **STOP y VALIDAR**: `POST /projects` válido → `201` + creador `scrum_master`; `GET /projects/{id}` → `200`; nombre vacío → `400`
4. Demo/deploy del MVP de proyectos

### Incremental Delivery

1. Setup + Foundational → base lista (dominio, puertos, adaptadores, esquema)
2. US1 → MVP (crear proyecto, vincular creador y consultar proyecto)
3. US2 → equipo y roles (asignación, reutilización, duplicados)
4. US3 → edición simple (reutiliza validaciones del alta, sin historial)
5. Polish → gofmt/vet, suite completa, conformidad de contrato, quickstart

---

## Notes

- [P] = archivos distintos, sin dependencias pendientes
- La etiqueta [Story] da trazabilidad HU → Spec → BDD → Tests → Código
- Restricciones del `data-model.md` citadas literalmente (nombre ≤ 100, descripción ≤ 2000, integrante ≤ 200, `scrum_master`/`product_builder`, fin ≥ inicio)
- FR-008 (creador Scrum Master) se implementa en US1 con un campo `creador` en el request, ya que la autenticación está fuera de alcance
- `GET /projects/{id}` (T023) y `GET /projects/{id}/members` (T029) se incluyen como lectura de soporte para verificar persistencia (SC-002, SC-009); el resto del frontend se difiere (R12)
- La aplicación efectiva de permisos por rol y el estado de ciclo de vida del proyecto (HU-13) quedan fuera de alcance
- Evitar: tareas vagas, conflictos de archivo simultáneos (`project_handler.go`, `project_dto.go`, `creacion_proyecto_equipo.feature`), dependencias cruzadas que rompan la independencia
