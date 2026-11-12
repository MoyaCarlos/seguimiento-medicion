---

description: "Tareas de la spec 044: Cerrar y reabrir un proyecto (#300)"
---

# Tasks: Cerrar y reabrir un proyecto (#300)

**Input**: Design documents from `/specs/044-cerrar-reabrir-proyecto/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/)

**Tests**: OBLIGATORIOS. La constitución del proyecto (Principio I, NON-NEGOTIABLE) exige TDD estricto con un commit por fase.

**Organization**: por historia de usuario, en el orden en que se pueden probar:

- US1 (cerrar) va primero porque las demás necesitan un proyecto cerrado.
- US2 (solo lectura) y US3 (reabrir) vienen después.
- US4 (solo el Scrum Master) **no tiene fase propia**: el permiso lo da el guardia de HU-21 que ya existe, así que sus casos (403, 401 y 404) se prueban dentro de los servicios y handlers de US1 y US3. Una fase propia tendría tests que pasan apenas se escriben, sin `RED:` real.
- Godog va al final, con el `.feature` y sus pasos juntos: `TestFeatures` corre con `Strict: true`, y un `.feature` sin pasos rompe el CI (el mismo criterio que en las specs 040 y 043).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 cerrar, US2 solo lectura, US3 reabrir, US4 permisos

---

## Phase 1: Setup

Sin tareas: la rama `feature/300-cerrar-reabrir-proyecto` ya existe y sale de `main` (con HU-21, PR #301), y no hay dependencias nuevas.

---

## Phase 2: Foundational — columna `fecha_cierre` (migración 008) y `Project.FechaCierre`

**Goal**: la base guarda la fecha de cierre y el dominio la conoce. Todavía nadie cierra ([data-model.md](data-model.md), [research.md §7](research.md)).

**Independent Test**: `go test ./internal/repository/ -run 'Esquema|Migracion008|Proyecto' -v`.

- [ ] T001 RED: esquema y migración 008.
  - En `internal/repository/esquema_test.go`, `TestEsquema_Filas`: sumar los casos D12 de `fecha_cierre` junto a los de `fecha_inicio`/`fecha_fin`:
    - `'2026-10-9'`, `'basura'` y `'2026-02-30'` → `violaCheck`;
    - `'2026-10-09'` y `NULL` → `aceptada`.
  - En `TestEsquema_Estructura`, la columna `fecha_cierre` en `proyectos`.
  - Nuevo `internal/repository/migracion_008_test.go`, con el patrón de `migracion_007_test.go`:
    - helper `baseAntesDeLa008` (todas las migraciones embebidas salvo `008_`);
    - siembra dos proyectos, aplica todas y verifica que siguen ahí con `fecha_cierre IS NULL`.
  - Correr `go test ./internal/repository/` y confirmar que FALLA (no existe la columna).
- [ ] T002 COMMIT `RED: agregar tests de esquema de proyectos.fecha_cierre y migración 008 (falla)` (solo `_test.go`)
- [ ] T003 GREEN: crear `internal/repository/migraciones/008_cierre_proyecto.sql`.
  - Contenido: comentario de cabecera (`-- Cierre manual del proyecto: MODELO-DE-DATOS.md §4.1 (D16, spec 044, #300).`) y `ALTER TABLE proyectos ADD COLUMN fecha_cierre TEXT CHECK (date(fecha_cierre) IS fecha_cierre);`.
  - `go test ./...` en verde.
- [ ] T004 COMMIT `GREEN: migración 008 con proyectos.fecha_cierre, test en verde`
- [ ] T005 Documentar el esquema en `MODELO-DE-DATOS.md`, en el mismo PR que la migración (regla de AGENTS.md):
  - §2: `fecha_cierre` en las fechas de calendario y en el `CHECK` D12.
  - §3: `text fecha_cierre` en `proyectos` del diagrama.
  - §4.1: fila `fecha_cierre | TEXT? | Nueva. #300. YYYY-MM-DD; NULL ⇔ abierto` y el `CHECK`.
  - §4.10: fila de `date(col) IS col` con `fecha_cierre`.
  - §6: estado del proyecto ← `proyectos.fecha_cierre` + `proyectos.fecha_*` + `sprints.estado`.
  - §7: D16 **cerrada el 2026-11-11**, con la decisión del issue (cerrado ⇒ Finalizado con prioridad; abierto ⇒ HU-13 sin cambios; la propuesta «sin cierre manual no hay Finalizado» se descarta).
  - §8: fila de #300 → `proyectos.fecha_cierre`, migración 008.
  - §9: fila `008 | proyectos.fecha_cierre | #300 (spec 044)` en la tabla de números reservados, y los casos de `fecha_cierre` en el punto 6.
  - Las filas de §5 van en T036, cuando existan las reglas.
- [ ] T006 COMMIT `docs: proyectos.fecha_cierre en MODELO-DE-DATOS.md, cierra D16 (Refs #300)`
- [ ] T007 RED: lectura y escritura de `fecha_cierre` en `internal/repository/sqlite_project_test.go`.
  - Con un proyecto al que se le pone `fecha_cierre = '2026-11-11'` por SQL directo, verificar:
    - `ObtenerPorID` devuelve `FechaCierre` = 2026-11-11 00:00 UTC;
    - `ListarPorIntegrante` también la devuelve;
    - con `NULL`, `FechaCierre == nil`;
    - `Actualizar` (editar nombre) **no la pisa**.
  - Confirmar que FALLA (no compila: `domain.Project` no tiene `FechaCierre`).
- [ ] T008 COMMIT `RED: agregar tests de lectura de fecha_cierre en el repositorio de proyectos (falla)`
- [ ] T009 GREEN: implementar la lectura.
  - `internal/domain/project.go`: campo `FechaCierre *time.Time`, con un comentario corto: nil mientras el proyecto está abierto (#300).
  - `internal/repository/sqlite_project.go`:
    - `fecha_cierre` en los `SELECT` de `ObtenerPorID` y `ListarPorIntegrante`;
    - `escanearProyecto` la lee con `valorATiempoOpcional`.
  - `Actualizar` **no** la incluye en su `UPDATE`.
  - `go test ./...` en verde.
- [ ] T010 COMMIT `GREEN: leer proyectos.fecha_cierre en el repositorio, test en verde`

**Checkpoint**: la columna existe, se lee y está documentada; nadie la escribe todavía.

---

## Phase 3: User Story 1 — Cerrar un proyecto (P1) 🎯 MVP (+ US4)

**Goal**: `POST /projects/{id}/cerrar` deja el proyecto «Finalizado», con la fecha de cierre de hoy ([contracts/cerrar-reabrir-proyecto.yaml](contracts/cerrar-reabrir-proyecto.yaml)).

**Independent Test**: `go test ./internal/domain/ ./internal/repository/ ./internal/service/ ./internal/http/ -run 'Cerrar|EstadoProyecto'`.

- [ ] T011 [US1] RED: estado con cierre.
  - En `internal/domain/estado_proyecto_test.go`, `TestCalcularEstadoProyecto_ProyectoCerrado` (tabla): un proyecto con `FechaCierre` da «Finalizado»:
    - con un Sprint `Activo`;
    - con fechas futuras (Planificado sin cierre);
    - en curso por fechas;
    - sin Sprints ni fechas.
  - Los casos de `TestCalcularEstadoProyecto_TablaDeDecision` (sin cierre) siguen igual.
  - Confirmar que FALLA.
- [ ] T012 [US1] COMMIT `RED: agregar tests del estado de un proyecto cerrado (falla)`
- [ ] T013 [US1] GREEN: en `internal/domain/estado_proyecto.go`, primera línea de `CalcularEstadoProyecto`: `if p.FechaCierre != nil { return ProyectoFinalizado }`. Actualizar el comentario: el cierre manual (#300) tiene prioridad. `go test ./...` en verde.
- [ ] T014 [US1] COMMIT `GREEN: un proyecto cerrado es Finalizado, test en verde`
- [ ] T015 [P] [US1] RED: en `internal/domain/project_test.go`, `TestProject_Cerrar` (tabla) sobre `func (p Project) Cerrar(ahora time.Time, sprints []Sprint) (Project, error)`:
  - abierto, sin Sprints → `FechaCierre` = día civil de `ahora` (con `ahora` 2026-11-11 15:30 local ⇒ 2026-11-11 00:00 UTC, el criterio de `inicioDelDiaCivil`);
  - abierto con Sprints `Pendiente` y `Finalizado` → se cierra;
  - abierto con un Sprint `Activo` → `ErrSprintActivoAlCerrarProyecto`, mensaje exacto `Cerrá el Sprint activo antes de cerrar el proyecto`;
  - ya cerrado → `ErrProyectoYaCerrado`, `El proyecto ya está cerrado`, y la fecha original no cambia;
  - ya cerrado con un Sprint `Activo` → `ErrProyectoYaCerrado` (va primero);
  - el receptor no se modifica (devuelve copia).
  - Confirmar que FALLA (no compila).
- [ ] T016 [US1] COMMIT `RED: agregar tests de Project.Cerrar (falla)`
- [ ] T017 [US1] GREEN: implementar las reglas del cierre.
  - `internal/domain/errors.go`: `ErrProyectoYaCerrado = errors.New("El proyecto ya está cerrado")` y `ErrSprintActivoAlCerrarProyecto = errors.New("Cerrá el Sprint activo antes de cerrar el proyecto")`, con comentario de doc.
  - `internal/domain/project.go`: `Cerrar`, con `inicioDelDiaCivil(ahora)`.
  - `go test ./...` en verde.
- [ ] T018 [US1] COMMIT `GREEN: implementar Project.Cerrar, test en verde`
- [ ] T019 [US1] RED: en `internal/repository/sqlite_project_test.go`, tests de `Cerrar(ctx, id int64, fecha time.Time) error`:
  - abierto sin Sprint `Activo` → guarda `'2026-11-11'` (verificar con `ObtenerPorID`);
  - ya cerrado → `domain.ErrProyectoYaCerrado`, y la fecha previa queda;
  - con un Sprint `Activo` (sembrado por SQL) → `domain.ErrSprintActivoAlCerrarProyecto`, sin cerrar;
  - inexistente → `domain.ErrProyectoNoEncontrado`;
  - carrera: dos `Cerrar` seguidos sobre el mismo proyecto → el segundo `ErrProyectoYaCerrado` (FR-017).
  - Lo mismo en `internal/repository/project_repository_en_memoria_test.go`, salvo el caso del Sprint Activo: el repositorio en memoria no conoce los Sprints.
  - Fakes de test que implementan `ProjectRepository`: sumar `Cerrar` a `fakeProjectRepository` (`internal/service/fakes_test.go`) y a `stubProjectRepository` (`internal/http/project_handler_test.go`). Un método de más compila aunque el puerto todavía no lo tenga, así que van en este RED y el GREEN queda solo con producción.
  - Confirmar que FALLA.
- [ ] T020 [US1] COMMIT `RED: agregar tests de Cerrar en los repositorios de proyectos (falla)`
- [ ] T021 [US1] GREEN: implementar `Cerrar` en los repositorios.
  - `internal/repository/project_repository.go`: `Cerrar` en `ProjectRepository`, con comentario (UPDATE condicional, errores al afectar 0 filas).
  - `internal/repository/sqlite_project.go`: `UPDATE proyectos SET fecha_cierre = ? WHERE id = ? AND fecha_cierre IS NULL AND NOT EXISTS (SELECT 1 FROM sprints WHERE proyecto_id = proyectos.id AND estado = 'Activo')` con `tiempoOpcionalAValor(&fecha)`. Con `RowsAffected` 0, releer:
    - no existe → `ErrProyectoNoEncontrado`;
    - `fecha_cierre` no nula → `ErrProyectoYaCerrado`;
    - si no → `ErrSprintActivoAlCerrarProyecto` ([research.md §5](research.md)).
  - `internal/repository/project_repository_en_memoria.go`: lo mismo, bajo su mutex. El repositorio en memoria no conoce los Sprints: verifica solo «ya cerrado» y la regla del Sprint Activo la da el dominio (documentarlo en el comentario).
  - `go test ./...` en verde.
- [ ] T022 [US1] COMMIT `GREEN: Cerrar en los repositorios de proyectos con UPDATE condicional, test en verde`
- [ ] T023 [US1] REFACTOR previo (sin cambio de comportamiento): en `internal/service/integrantes.go`, separar el chequeo de rol del guardia de escritura. Renombrar el cuerpo actual de `exigirIntegrante` a `verificarIntegrante` y el de `exigirScrumMaster` a `verificarScrumMaster`, y dejar `exigirIntegrante`/`exigirScrumMaster` llamándolas. `go test ./...` en verde antes y después ([research.md §1](research.md)).
- [ ] T024 [US1] COMMIT `REFACTOR: separar verificarIntegrante/ScrumMaster de los guardias exigir*`
- [ ] T025 [US1] [US4] RED: en `internal/service/cerrar_proyecto_test.go`, `NewCerrarProyecto(proyectos repository.ProjectRepository, sprints repository.SprintRepository, reloj func() time.Time)` y `Ejecutar(ctx, proyectoID, usuarioID int64) (domain.Project, error)`, con `fakeProjectRepository` y `repository.NewSprintRepositoryEnMemoria()`:
  - Scrum Master, proyecto abierto sin Sprints → devuelve el proyecto con `FechaCierre` = hoy del reloj, persistido;
  - con un Sprint `Activo` → `ErrSprintActivoAlCerrarProyecto` y el proyecto sigue abierto;
  - ya cerrado → `ErrProyectoYaCerrado`;
  - US4: Product Builder → `ErrNoEsScrumMaster`;
  - US4: no integrante → `ErrNoEsIntegrante`;
  - US4: proyecto inexistente → `ErrProyectoNoEncontrado`;
  - orden: inexistente gana sobre no integrante; Product Builder sobre un proyecto ya cerrado → `ErrNoEsScrumMaster` (403 antes que 409).
  - Confirmar que FALLA.
- [ ] T026 [US1] COMMIT `RED: agregar tests de CerrarProyecto (falla)`
- [ ] T027 [US1] GREEN: crear `internal/service/cerrar_proyecto.go`.
  - Flujo de `Ejecutar`:
    1. `ObtenerPorID` (404);
    2. `verificarScrumMaster` (403, el chequeo de rol **sin** el de proyecto abierto);
    3. listar los Sprints del proyecto;
    4. `p.Cerrar(s.reloj(), sprints)` (409);
    5. `s.proyectos.Cerrar(ctx, id, *cerrado.FechaCierre)`;
    6. devolver `cerrado`.
  - `verificarScrumMaster` ya existe desde el REFACTOR de T023.
  - `go test ./...` en verde.
- [ ] T028 [US1] COMMIT `GREEN: implementar CerrarProyecto, test en verde`
- [ ] T029 [US1] [US4] RED: handler y rutas de cerrar.
  - En `internal/http/project_handler_test.go`, tests de `POST /projects/{id}/cerrar` (actor por sesión, como los demás tests del handler):
    - 200 con cuerpo `{..., "fecha_cierre":"2026-11-11", "estado":"Finalizado"}`;
    - 409 `{"mensaje":"El proyecto ya está cerrado"}`;
    - 409 `{"mensaje":"Cerrá el Sprint activo antes de cerrar el proyecto"}`;
    - 403 de Product Builder;
    - 404 de proyecto inexistente.
  - `GET /projects/{id}` y `GET /projects` incluyen `"fecha_cierre": null` en un proyecto abierto.
  - En `internal/app/rutas_permisos_test.go`, `nivelesDeLaSpec`: `"POST /projects/{id}/cerrar": apihttp.SesionScrumMaster` (y `"POST /projects/{id}/reabrir"`, que se registra en US3: agregarlo recién en T053, para que este RED no dependa de US3).
  - Ajustar los llamadores de `NewProjectHandler` en test (`internal/http/project_handler_test.go`, `features/steps_listado_proyectos.go`) a la firma nueva, que suma solo `*service.CerrarProyecto`. `ReabrirProyecto` todavía no existe: se agrega en US3 (T055).
  - Confirmar que FALLA.
- [ ] T030 [US1] COMMIT `RED: agregar tests del handler POST /projects/{id}/cerrar (falla)`
- [ ] T031 [US1] GREEN: implementar el handler y la ruta de cerrar.
  - `internal/http/project_dto.go`:
    - `FechaCierre *string \`json:"fecha_cierre"\`` en `proyectoResponse`, con `aProyectoResponse` formateando `YYYY-MM-DD`;
    - tipo `proyectoConEstadoResponse { proyectoResponse; Estado string \`json:"estado"\` }`.
  - `internal/http/project_handler.go`:
    - `NewProjectHandler` recibe `cerrar *service.CerrarProyecto`;
    - handler `Cerrar`: `actorDeSesion`, id de la ruta, `Ejecutar`, `escribirError` ante error y 200 con el proyecto y `domain.CalcularEstadoProyecto(p, nil, time.Now())`. Un proyecto cerrado siempre da «Finalizado», así que no hacen falta los Sprints.
  - `internal/http/sprint_handler.go`, `escribirError`: `ErrProyectoYaCerrado` y `ErrSprintActivoAlCerrarProyecto` en el `case` de 409.
  - `internal/http/rutas.go`: `{"POST /projects/{id}/cerrar", SesionScrumMaster, h.Proyectos.Cerrar}`.
  - `internal/app/app.go`: `service.NewCerrarProyecto(proyectos, sprintRepo, op.Ahora)` en `NewProjectHandler`.
  - `gofmt -l .` sin salida; `go test ./...` en verde.
- [ ] T032 [US1] COMMIT `GREEN: POST /projects/{id}/cerrar, test en verde`
- [ ] T033 [US1] (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: el Scrum Master puede cerrar; el proyecto se ve «Finalizado». Todavía se puede escribir sobre él (US2).

---

## Phase 4: User Story 2 — Proyecto cerrado en solo lectura (P1)

**Goal**: toda escritura sobre un proyecto cerrado responde 409 «El proyecto está cerrado» sin cambiar datos, y las lecturas siguen igual. El chequeo vive en un solo lugar ([research.md §1, §2, §9](research.md)).

**Independent Test**: `go test ./internal/service/ -run 'Integrantes|Votar|ImputarHoras' && go test ./internal/app/ -run ProyectoCerrado`.

- [ ] T034 [US2] RED: el dato «cerrado» en los repositorios.
  - En `internal/repository/sqlite_project_test.go`: `ProyectoCerrado(ctx, id) (bool, error)` da `true` con `fecha_cierre` y `false` sin ella, y también `false` para un id inexistente (el 404 lo resuelve el servicio antes).
  - El mismo test para `ProjectRepositoryEnMemoria`.
  - Fakes de test que implementarán `IntegrantesRepository`, con `ProyectoCerrado(ctx, id) (bool, error)` (en este RED, como en T019):
    - `fakeProjectRepository` (campo `cerrados map[int64]bool`) y `fakeIntegrantesRepository` (ídem) en `internal/service/fakes_test.go`;
    - `integrantesConActor` y `scrumMasterConActor` → `false`;
    - `stubProjectRepository` de `internal/http/project_handler_test.go` → `false`.
  - Confirmar que FALLA.
- [ ] T035 [US2] COMMIT `RED: agregar tests de ProyectoCerrado en los repositorios (falla)`
- [ ] T036 [US2] GREEN: implementar `ProyectoCerrado`.
  - Implementarlo en `internal/repository/sqlite_project.go` (`SELECT fecha_cierre IS NOT NULL FROM proyectos WHERE id = ?`; `sql.ErrNoRows` → `false, nil`) y en `internal/repository/project_repository_en_memoria.go`.
  - Agregarlo a `IntegrantesRepository` (`internal/repository/registro_horas_repository.go`), con el comentario del puerto actualizado: «control de acceso de HU-21: integrantes y si el proyecto está cerrado (#300)».
  - En `MODELO-DE-DATOS.md` §5, dos filas:
    - «Cerrar solo sin Sprint Activo, y no dos veces»: base = `UPDATE` condicional; servicio = `Project.Cerrar` (409).
    - «Proyecto cerrado = solo lectura»: servicio = guardia `exigirProyectoAbierto` (409), con el límite de la escritura concurrente ([research.md §6](research.md)).
  - `go test ./...` en verde.
- [ ] T037 [US2] COMMIT `GREEN: ProyectoCerrado en el puerto de control de acceso, test en verde`
- [ ] T038 [US2] RED: el guardia.
  - En `internal/service/integrantes_test.go`, con `fakeIntegrantesRepository` y un proyecto cerrado:
    - `exigirIntegrante` y `exigirScrumMaster` → `domain.ErrProyectoCerrado` (mensaje exacto `El proyecto está cerrado`) para un integrante o Scrum Master;
    - un no integrante sigue recibiendo `ErrNoEsIntegrante`, y un Product Builder en `exigirScrumMaster`, `ErrNoEsScrumMaster` (403 antes que 409);
    - `verificarIntegrante` y `verificarScrumMaster` **no** devuelven error por el cierre;
    - con el proyecto abierto, todo como antes.
  - En `internal/service/votar_test.go` e `internal/service/imputar_horas_test.go`: un integrante que vota, o registra horas, en un proyecto cerrado → `ErrProyectoCerrado`, sin voto ni registro guardado.
  - En `internal/service/consultar_votacion_test.go` (o el test existente de la consulta): con el proyecto cerrado, un integrante **sí** consulta.
  - Confirmar que FALLA (no compila: `ErrProyectoCerrado` no existe).
- [ ] T039 [US2] COMMIT `RED: agregar tests del guardia de proyecto cerrado (falla)`
- [ ] T040 [US2] GREEN: implementar el guardia.
  - `internal/domain/errors.go`: `ErrProyectoCerrado = errors.New("El proyecto está cerrado")`.
  - `internal/service/integrantes.go`: `exigirProyectoAbierto(ctx, integrantes, proyectoID)` (único lugar del chequeo, FR-005), llamada al final de `exigirIntegrante` y de `exigirScrumMaster`.
  - `internal/service/consultar_votacion.go`: pasa a `verificarIntegrante`.
  - `internal/service/votar.go`: reemplazar el chequeo con `rolEnProyecto` por `exigirIntegrante`, conservando el error actual del no integrante.
  - `internal/service/imputar_horas.go`: `exigirIntegrante` antes de `nombreDelIntegrante`.
  - `internal/http/sprint_handler.go`, `escribirError`: `ErrProyectoCerrado` → 409.
  - `go test ./...` en verde, incluido Godog sin cambios (ningún escenario existente cierra proyectos).
- [ ] T041 [US2] COMMIT `GREEN: un proyecto cerrado rechaza las escrituras con 409, test en verde`
- [ ] T042 [US2] RED/verificación: nuevo `internal/app/proyecto_cerrado_test.go`, sobre `clienteDePrueba`.
  - Arrange, vía la API, con `ana` como Scrum Master del proyecto 1:
    - historias 1 y 2;
    - Sprint 1 iniciado, con la historia 1 asignada y un defecto registrado;
    - Sprint 1 cerrado y Sprint 2 creado (Pendiente);
    - votación iniciada en la historia 2;
    - `POST /projects/1/cerrar` → 200.
  - `TestProyectoCerrado_TodaEscrituraResponde409SinCambios`: una tabla `escriturasDeProyecto` (ruta concreta + cuerpo válido) con **cada** patrón de nivel `SesionIntegrante` y `SesionScrumMaster` salvo `POST /projects/{id}/cerrar` y `/reabrir` (17 rutas). Para armar cuerpos válidos, ver los `*_dto.go` y los tests de cada handler. Cada una → 409 `{"mensaje":"El proyecto está cerrado"}`.
  - Para «sin cambios», comparar el cuerpo de `GET /projects/1`, `/projects/1/backlog`, `/projects/1/sprints`, `/projects/1/defectos`, `/projects/1/members` y `/backlog/2/votacion` antes y después.
  - `TestProyectoCerrado_LaTablaCubreTodasLasEscrituras`: compara los patrones de `escriturasDeProyecto` con los de `apihttp.Rutas(apihttp.Handlers{})` de esos dos niveles (menos cerrar y reabrir), con el patrón de `TestRutas_CadaRutaTieneElNivelDeLaSpec`.
  - `TestProyectoCerrado_LasLecturasSiguen`: cada `GET` de `nivelesDeLaSpec` sobre el proyecto 1 responde 200, y `GET /projects/1/status` responde `{"estado":"Finalizado"}`.
  - Si los tres pasan de entrada (el guardia ya está), es la verificación de US2 y no hay `RED:` real. Commitear como `test: verificar que un proyecto cerrado rechaza todas las escrituras (Refs #300)`.
  - Si alguna ruta **no** da 409, ese caso es el `RED:`: commitear `RED: <ruta> escribe sobre un proyecto cerrado (falla)`, después un `GREEN:` que la pase por el guardia y avisar en la revisión.
- [ ] T043 [US2] COMMIT según T042
- [ ] T044 [US2] (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: un proyecto cerrado es de solo lectura en toda la API.

---

## Phase 5: User Story 3 — Reabrir un proyecto (P1) (+ US4)

**Goal**: `POST /projects/{id}/reabrir` deja vacía la fecha de cierre y el proyecto vuelve a aceptar escrituras.

**Independent Test**: `go test ./internal/domain/ ./internal/repository/ ./internal/service/ ./internal/http/ -run Reabrir`.

- [ ] T045 [P] [US3] RED: en `internal/domain/project_test.go`, `TestProject_Reabrir`:
  - cerrado → copia con `FechaCierre == nil`, y el receptor no cambia;
  - abierto → `ErrProyectoNoCerrado`, mensaje exacto `El proyecto no está cerrado`;
  - un proyecto cerrado y reabierto con la fecha de fin vencida da «Finalizado» por HU-13 (ver el borde en la spec).
  - Confirmar que FALLA.
- [ ] T046 [US3] COMMIT `RED: agregar tests de Project.Reabrir (falla)`
- [ ] T047 [US3] GREEN: `ErrProyectoNoCerrado` en `internal/domain/errors.go` y `Reabrir` en `internal/domain/project.go`. `go test ./...` en verde.
- [ ] T048 [US3] COMMIT `GREEN: implementar Project.Reabrir, test en verde`
- [ ] T049 [US3] RED: repositorios y servicio de reabrir.
  - En `internal/repository/sqlite_project_test.go` (y en memoria), `Reabrir(ctx, id) error`:
    - cerrado → `fecha_cierre` `NULL`;
    - abierto → `domain.ErrProyectoNoCerrado`;
    - inexistente → `ErrProyectoNoEncontrado`;
    - carrera: dos `Reabrir` seguidos → el segundo `ErrProyectoNoCerrado` (FR-017).
  - Fakes de test: `Reabrir` en `fakeProjectRepository` y `stubProjectRepository` (en este RED, como en T019).
  - En `internal/service/reabrir_proyecto_test.go`, `NewReabrirProyecto(proyectos repository.ProjectRepository)` y `Ejecutar(ctx, proyectoID, usuarioID)`:
    - el Scrum Master reabre → `FechaCierre == nil`, persistido;
    - abierto → `ErrProyectoNoCerrado`;
    - US4: Product Builder → `ErrNoEsScrumMaster`, no integrante → `ErrNoEsIntegrante`, inexistente → `ErrProyectoNoEncontrado`, y el proyecto sigue cerrado;
    - **reabrir no recibe `ErrProyectoCerrado`** (usa `verificarScrumMaster`).
  - Confirmar que FALLA.
- [ ] T050 [US3] COMMIT `RED: agregar tests de Reabrir en repositorios y ReabrirProyecto (falla)`
- [ ] T051 [US3] GREEN: implementar la reapertura.
  - `Reabrir` en `ProjectRepository`, en `internal/repository/sqlite_project.go` (`UPDATE proyectos SET fecha_cierre = NULL WHERE id = ? AND fecha_cierre IS NOT NULL`; con 0 filas, releer: `ErrProyectoNoEncontrado` o `ErrProyectoNoCerrado`) y en el repositorio en memoria.
  - `internal/service/reabrir_proyecto.go`: `ObtenerPorID` → `verificarScrumMaster` → `p.Reabrir()` → `s.proyectos.Reabrir`.
  - `go test ./...` en verde.
- [ ] T052 [US3] COMMIT `GREEN: implementar ReabrirProyecto, test en verde`
- [ ] T053 [US3] [US4] RED: handler y rutas de reabrir.
  - En `internal/http/project_handler_test.go`, `POST /projects/{id}/reabrir`:
    - 200 con `"fecha_cierre": null` y `estado` según HU-13;
    - 409 `{"mensaje":"El proyecto no está cerrado"}`;
    - 403 de Product Builder;
    - 404.
  - `nivelesDeLaSpec`: `"POST /projects/{id}/reabrir": apihttp.SesionScrumMaster`.
  - En `internal/app/proyecto_cerrado_test.go`:
    - `TestProyectoCerrado_AlReabrirVuelveAAceptarEscrituras`: cerrar, reabrir y `POST /backlog` → 201;
    - `TestProyectoCerrado_ReabrirDevuelveElEstadoPrevio` (SC-003): un proyecto «En curso» (con fecha de inicio pasada y fin futura respecto del reloj de prueba, sin Sprints) da el mismo `GET /projects/{id}/status` antes de cerrar y después de reabrir.
  - Ajustar los llamadores de `NewProjectHandler` en test a la firma que suma `*service.ReabrirProyecto`, pasando `service.NewReabrirProyecto(...)`.
  - Confirmar que FALLA.
- [ ] T054 [US3] COMMIT `RED: agregar tests del handler POST /projects/{id}/reabrir (falla)`
- [ ] T055 [US3] GREEN: implementar el handler y la ruta de reabrir.
  - `NewProjectHandler` suma `reabrir *service.ReabrirProyecto`.
  - Handler `Reabrir` en `internal/http/project_handler.go`: responde el proyecto con su estado recalculado, usando el servicio `ObtenerEstadoProyecto` existente para no duplicar la lectura de Sprints.
  - `ErrProyectoNoCerrado` → 409 en `escribirError`.
  - La ruta en `internal/http/rutas.go`.
  - Cableado en `internal/app/app.go`.
  - `go test ./...` en verde.
- [ ] T056 [US3] COMMIT `GREEN: POST /projects/{id}/reabrir, test en verde`
- [ ] T057 [US3] (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: las cuatro historias funcionan por la API.

---

## Phase 6: Escenarios BDD (US1, US2, US3, US4)

- [ ] T058 [US1] [US2] [US3] [US4] Escribir `features/cerrar_reabrir_proyecto.feature`.
  - Encabezado y Característica:
    - `# language: es` y el tag `@300`;
    - Característica con el Como/Quiero/Para de la Descripción de la spec;
    - comentario con las categorías: normal, alternativo, límite y error.
  - Antecedentes, reutilizando los pasos de HU-21:
    - `Dado que inicié sesión como "ana"`;
    - `Y que "ana" creó el proyecto "Tienda"`;
    - `Y que "bruno" tiene cuenta y es Product Builder de "Tienda"`;
    - `Y que "Tienda" tiene la historia "Carrito"`.
  - Un escenario por criterio de aceptación de la spec, sin inventar casos:
    - US1.1 cerrar → «Finalizado» con fecha de cierre `2026-11-11` (el «hoy» de los escenarios);
    - US1.3 Sprint Activo → 409;
    - US1.4 cerrar dos veces → 409;
    - US2 `Esquema del escenario` con escrituras representativas de integrante (crear una historia, votar, registrar horas) y de Scrum Master (crear un Sprint, editar el proyecto, agregar a alguien al equipo) → 409 «El proyecto está cerrado» y no cambia nada;
    - US2.3 lecturas: consultar el backlog del proyecto cerrado funciona;
    - US3.1 y US3.2 reabrir y volver a crear una historia;
    - US3.3 reabrir uno abierto → 409;
    - US4 `Esquema` con «bruno» (403 «Solo el Scrum Master puede realizar esta acción»), un no integrante (403 «No sos integrante de este proyecto») y sin sesión (401), para cerrar y reabrir.
  - Si una categoría no aplica, decirlo en el comentario.
- [ ] T059 [US1] [US2] [US3] [US4] Pasos: nuevo `features/steps_cerrar_reabrir_proyecto.go` con `InitializeScenarioCierreProyecto(ctx, permisos *permisosContext)`.
  - `InitializeScenarioPermisos` (`features/steps_permisos_por_rol.go:559`) pasa a devolver su `*permisosContext`, para reutilizar `proyectos`, `historias`, `pedir`/`esperar` y `noCambiaNada`.
  - Registrar el initializer en `features/features_test.go` con el comentario `// #300`.
  - Reutilizar los pasos existentes donde el texto coincide (`no cambia nada`, `el sistema rechaza la petición con (\d+)`, `el sistema responde 403 con el mensaje "..."`). Pasos nuevos solo para cerrar/reabrir y sus aserciones.
  - `go test ./features/` en verde, incluidos todos los escenarios existentes.
- [ ] T060 COMMIT `Agregar escenarios BDD de cerrar y reabrir un proyecto (Refs #300)`

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T061 `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test ./...` sin errores.
- [ ] T062 Correr la validación manual de [quickstart.md](quickstart.md) (tabla de pasos con curl) y anotar el resultado en el PR.
- [ ] T063 Revisar que `MODELO-DE-DATOS.md` (§2, §3, §4.1, §4.10, §5, §6, §7, §8 y §9) quedó consistente con la migración 008 y el código. Ante una contradicción, frenar y avisar.
- [ ] T064 Marcar las tareas hechas en este archivo y commitear `docs: marcar las tareas de la spec 044 como hechas (Refs #300)`.

---

## Dependencies & Execution Order

- **Foundational (Phase 2)** → **US1 (Phase 3)** → **US2 (Phase 4)** → **US3 (Phase 5)** → **BDD (Phase 6)** → **Polish**.
- US2 depende de US1: necesita un proyecto cerrado vía la API para el test de T042.
- US3 depende de US1 y de T023: usa `verificarScrumMaster`, el renombre que se hace ahí.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT del RED.
- T005/T006 (documentación del esquema) pueden ir en cualquier momento antes del PR, pero se dejan junto a la migración para que el PR los traiga juntos.
- Si el PR #304 (spec 043) se mergea antes, traer `main` a la rama antes de abrir el PR y resolver los conflictos en `internal/domain/project.go`, `MODELO-DE-DATOS.md` §5 y `features/steps_estado_proyecto.go` ([research.md §10](research.md)).

## Parallel Opportunities

- T015 (tests de `Project.Cerrar`) y T011 (tests del estado) tocan archivos distintos del dominio. Se pueden escribir juntos, pero cada uno lleva su propio commit RED.
- T045 (tests de `Project.Reabrir`) no depende de US2: se puede adelantar después de T018.
- Entre fases no hay paralelismo real: los GREEN de los repositorios tocan los mismos fakes y el mismo `sqlite_project.go`.

## Implementation Strategy

1. **Foundational**: la columna, su documentación y la lectura.
2. **US1** (MVP): cerrar desde la API. Se valida con el quickstart, pasos 1, 2, 6, 10 y 11.
3. **US2**: solo lectura. Es lo que le da sentido al cierre; con US1 y US2 juntos ya se podría entregar.
4. **US3**: reabrir.
5. **BDD y cierre.** El PR se abre contra `main` con `Closes #300` y se mergea con «Create a merge commit». El frontend lo hace otra integrante en otro PR.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- Ante cualquier duda sobre el modelo de datos durante la implementación, se frena y se consulta (AGENTS.md).
