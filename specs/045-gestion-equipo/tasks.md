---

description: "Tareas de la spec 045: Gestión del equipo del proyecto (HU-16, #175)"
---

# Tasks: Gestión del equipo del proyecto (HU-16)

**Input**: Design documents from `/specs/045-gestion-equipo/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/api-equipo.md](contracts/api-equipo.md)

**Tests**: OBLIGATORIOS. La constitución del proyecto (Principio I, NON-NEGOTIABLE) exige TDD estricto con un commit por fase.

**Organization**: por historia de usuario, en el orden en que se pueden probar:

- La fundación (Phase 2) agrega la columna `fecha_baja` y el filtro «integrante = activo», que usan todas las historias.
- US1 (quitar) y US2 (Scrum Master) comparten el servicio, el repositorio y el endpoint de baja: US2 no tiene fase propia de código, sus casos (409 al quitar al Scrum Master, 400 al reagregarlo) se prueban dentro de las tareas de US1 y US3. Una fase propia tendría tests que pasan apenas se escriben, sin `RED:` real.
- US4 (no encontrado) tampoco tiene fase propia: sus casos (404 proyecto, 404 integrante, baja repetida) se prueban en el servicio y el handler de US1.
- US5 (alta como Product Builder) ya está en `permisos_por_rol.feature` (HU-21): solo se verifica que siga en verde.
- Godog va al final, con el `.feature` y sus pasos juntos: `TestFeatures` corre con `Strict: true` y un `.feature` sin pasos rompe el CI (mismo criterio que las specs 040, 043 y 044).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 quitar, US2 Scrum Master, US3 reactivar, US4 no encontrado, US6 frontend

---

## Phase 1: Setup — base con la migración 008

**Goal**: la rama tiene la migración 008 de #300 antes de agregar la 009 ([research.md §8](research.md)).

- [x] T001 Confirmar con `git fetch` que `origin/main` tiene `internal/repository/migraciones/008_cierre_proyecto.sql` (PR #311, mergeado el 2026-11-12) y que no apareció otra `009_` en `main` mientras tanto. Si apareció, frenar y avisar: no se renumera por cuenta propia.
- [x] T002 Mergear `origin/main` en `feature/175-gestion-equipo` (`git merge origin/main`, merge commit) y correr `go test ./...` en verde.

---

## Phase 2: Foundational — columna `fecha_baja` (migración 009), errores y «integrante = activo»

**Goal**: la base guarda la baja, el dominio conoce sus errores y regla, y toda lectura de integrantes devuelve solo activos. Todavía nadie da de baja ([data-model.md](data-model.md)).

**Independent Test**: `go test ./internal/repository/ ./internal/domain/ -v`.

### Esquema

- [x] T003 RED: esquema y migración 009.
  - En `internal/repository/esquema_test.go`, `TestEsquema_Filas`: casos D12 de `integrantes.fecha_baja` — `'2026-11-1'`, `'basura'` y `'2026-02-30'` → `violaCheck`; `'2026-11-12'` y `NULL` → `aceptada`.
  - En `TestEsquema_Estructura`, la columna `fecha_baja` en `integrantes`.
  - Nuevo `internal/repository/migracion_009_test.go`, con el patrón de `migracion_007_test.go`: helper `baseAntesDeLa009` (todas las migraciones embebidas salvo `009_`); siembra un proyecto con un Scrum Master y un Product Builder, aplica todas y verifica que ambos siguen con `fecha_baja IS NULL`.
  - Correr `go test ./internal/repository/` y confirmar que FALLA.
- [x] T004 COMMIT `RED: agregar tests de esquema de integrantes.fecha_baja y migración 009 (falla)` (solo `_test.go`)
- [x] T005 GREEN: crear `internal/repository/migraciones/009_baja_integrantes.sql`: cabecera `-- Baja lógica de integrantes: MODELO-DE-DATOS.md §4.3 (D15, spec 045, #175).` y `ALTER TABLE integrantes ADD COLUMN fecha_baja TEXT CHECK (fecha_baja IS NULL OR date(fecha_baja) IS fecha_baja);`. `go test ./...` en verde.
- [x] T006 COMMIT `GREEN: migración 009 con integrantes.fecha_baja, test en verde`
- [x] T007 Documentar el esquema en `MODELO-DE-DATOS.md`, en el mismo PR que la migración (regla de AGENTS.md):
  - §2: `fecha_baja` en las fechas de calendario y en el `CHECK` D12.
  - §3: `text fecha_baja` en `integrantes` del diagrama.
  - §4.3: fila `fecha_baja | TEXT? | Nueva. HU-16. YYYY-MM-DD; NULL ⇔ activo` y párrafo de baja lógica: quitar = fecha de baja; reactivar = misma fila con `fecha_baja = NULL`; el Scrum Master no se da de baja; «integrante» en las reglas = integrante activo.
  - §4.7: al dar de baja se borran los votos de esa persona en rondas no reveladas del proyecto; las reveladas no se tocan.
  - §4.10: fila de `date(col) IS col` con `fecha_baja`.
  - §7: D15 **cerrada el 2026-11-12** (confirmada por el equipo), con la decisión de la spec 045.
  - §8: fila de HU-16 → `integrantes.fecha_baja`, migración 009.
  - §9: fila `009 | integrantes.fecha_baja | HU-16 (spec 045)` y los casos de `fecha_baja` en el punto de tests de esquema.
  - Las filas de §5 van en T049, cuando existan las reglas.
- [x] T008 COMMIT `docs: integrantes.fecha_baja en MODELO-DE-DATOS.md, cierra D15 (Refs #175)`

### Dominio

- [x] T009 [P] RED: en `internal/domain/member_test.go`, `TestMember_ValidarBaja`: un `Member` con `Role: RolScrumMaster` → `errors.Is(err, ErrScrumMasterNoSeQuita)`; con `RolProductBuilder` → `nil`. Confirmar que FALLA (no compila).
- [x] T010 COMMIT `RED: agregar test de Member.ValidarBaja (falla)`
- [x] T011 GREEN: en `internal/domain/errors.go`, `ErrIntegranteNoEncontrado = errors.New("integrante no encontrado")` y `ErrScrumMasterNoSeQuita = errors.New("No se puede quitar al Scrum Master del proyecto")`; en `internal/domain/member.go`, `func (m Member) ValidarBaja() error`. Tests en verde.
- [x] T012 COMMIT `GREEN: implementar Member.ValidarBaja, test en verde`

### Puerto y lecturas de activos

- [x] T013 Ampliar el puerto `ProjectRepository` en `internal/repository/project_repository.go`:
  - `ListarIntegrantesHistoricos(ctx, proyectoID int64) ([]domain.Member, error)` — activos y dados de baja, orden de alta.
  - `DarDeBajaIntegrante(ctx, proyectoID, usuarioID int64, fecha time.Time) error`.
  - Comentario de `ListarIntegrantes` y `ListarPorIntegrante`: «solo integrantes activos (`fecha_baja IS NULL`)».
  - Agregar los dos métodos a todos los que implementan el puerto, devolviendo lo mínimo para compilar (`nil`/`errors.New("no implementado")` en los repos reales hasta su GREEN). La lista sale de buscar `func (.*) ListarIntegrantes(` **después del merge de T002** (el #311 sumó `ProyectoCerrado`, `Cerrar` y `Reabrir` al puerto y pudo agregar stubs); antes del #311 eran `internal/repository/project_repository_en_memoria.go`, `internal/repository/sqlite_project.go`, `internal/service/fakes_test.go`, `internal/service/iniciar_sprint_test.go` e `internal/http/project_handler_test.go`. En `origin/main` (después del #311) la búsqueda da además los fakes de `IntegrantesRepository` de `internal/service/fakes_test.go` (`integrantesConActor`, `scrumMasterConActor`, `fakeIntegrantesRepository`): ese puerto no cambia, así que no llevan los métodos nuevos. `fakeProjectRepository` lleva stubs acá y su comportamiento real en T015. `go build ./... && go vet ./...` sin errores.
- [x] T014 COMMIT `Agregar ListarIntegrantesHistoricos y DarDeBajaIntegrante al puerto ProjectRepository`
- [x] T015 RED: lecturas de activos en el fake y en SQLite.
  - `internal/repository/project_repository_en_memoria_test.go` y `internal/repository/sqlite_project_test.go`: con un Product Builder marcado como dado de baja (en SQLite, `UPDATE integrantes SET fecha_baja = '2026-11-12'` directo; en memoria, con un helper de test sobre el estado interno, porque `DarDeBajaIntegrante` recién se implementa en T022), verificar que `ListarIntegrantes` no lo devuelve y `ListarIntegrantesHistoricos` sí (con nombre y rol).
  - `internal/repository/sqlite_project_listar_test.go` y `project_repository_en_memoria_listar_test.go`: `ListarPorIntegrante` no devuelve el proyecto donde la persona está dada de baja.
  - Soporte de test en `internal/service/fakes_test.go`: `fakeProjectRepository` es el fake que usan los tests de `ConsultarVotacion`, `RevelarVotacion`, `ConsultarHistorialRondas`, `GenerarReporte` y `AsignarIntegrante`. Sin esto, T032, T032b, T035b y T045 no pueden expresar una baja. Agregarle:
    - Un registro de bajas por vínculo (p. ej. `bajas map[[2]int64]time.Time`, clave `{proyectoID, usuarioID}`).
    - `DarDeBajaIntegrante` que lo marca, con las mismas reglas que el repositorio en memoria (`ErrIntegranteNoEncontrado`, `ErrScrumMasterNoSeQuita`).
    - `ListarIntegrantes` filtra a los dados de baja, tanto en la lista derivada como en `listaMiembros`.
    - `ListarIntegrantesHistoricos` devuelve la lista sin filtrar.
    - `AgregarIntegrante` reactiva a un dado de baja (borra la marca, toma el rol nuevo y conserva `CreadoEn`) y sigue devolviendo `ErrMiembroDuplicado` para un activo.

    Va en este commit RED porque es solo código de test.
  - Confirmar que FALLA.
- [x] T016 COMMIT `RED: agregar tests de integrantes activos e históricos (falla)`
- [x] T017 GREEN: en `internal/repository/sqlite_project.go`, `AND i.fecha_baja IS NULL` en `ListarIntegrantes` y `ListarPorIntegrante`; `ListarIntegrantesHistoricos` = la consulta sin el filtro. En el fake en memoria, un campo de baja por vínculo (p. ej. `fechaBaja *time.Time` en un tipo interno, sin tocar `domain.Membership`) y el mismo filtrado. `go test ./...` en verde.
- [x] T018 COMMIT `GREEN: listar solo integrantes activos y agregar ListarIntegrantesHistoricos, test en verde`
- [x] T019 — no hizo falta: no hubo nada que limpiar. (solo si hay algo que limpiar, p. ej. consulta duplicada entre activos e históricos) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: la base tiene `fecha_baja`, las reglas existentes (pertenencia, rol, votación, horas, reporte, «Mis proyectos») ya ignoran a quien tenga fecha de baja, y `go test ./...` sigue en verde.

---

## Phase 3: User Story 1 - Quitar un integrante conservando su historial (P1) 🎯 MVP — incluye US2 y US4

**Goal**: el Scrum Master quita a un Product Builder con baja lógica; sus aportes se conservan; el Scrum Master no se puede quitar; inexistentes dan 404.

**Independent Test**: `go test ./internal/repository/ ./internal/service/ ./internal/http/ -run 'Baja|Quitar' -v` y la prueba manual de [quickstart.md](quickstart.md) pasos 1–8 y 10.

### Repositorio: `DarDeBajaIntegrante`

- [x] T020 [US1] RED: en `internal/repository/project_repository_en_memoria_test.go`, `TestProjectRepositoryEnMemoria_DarDeBajaIntegrante`:
  - Product Builder activo → `nil`; deja de salir en `ListarIntegrantes` y sigue en `ListarIntegrantesHistoricos`.
  - Segunda baja del mismo → `ErrIntegranteNoEncontrado`.
  - Usuario que nunca integró → `ErrIntegranteNoEncontrado`.
  - Scrum Master → `ErrScrumMasterNoSeQuita` y sigue activo.
  - Confirmar que FALLA.
- [x] T021 [US1] COMMIT `RED: agregar tests de DarDeBajaIntegrante en memoria (falla)`
- [x] T022 [US1] GREEN: `DarDeBajaIntegrante` en `internal/repository/project_repository_en_memoria.go` (bajo el mutex). Tests en verde.
- [x] T023 [US1] COMMIT `GREEN: implementar DarDeBajaIntegrante en memoria, test en verde`
- [x] T024 [US1] RED: en `internal/repository/sqlite_project_test.go` (o nuevo `sqlite_project_baja_test.go`), `TestSQLiteProjectRepository_DarDeBajaIntegrante`:
  - Mismos cuatro casos que T020; la fila queda con `fecha_baja = '2026-11-12'` (formato `YYYY-MM-DD`) y no se borra. La fecha del test se arma en UTC (`time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)`), porque T026 guarda el día en UTC.
  - Votos: con una historia del proyecto, una ronda **revelada** y otra **no revelada** donde votó el Product Builder (y votó también otra persona), después de la baja el voto de la ronda no revelada no existe y el de la revelada sigue igual; el voto de la otra persona sigue en las dos. Un voto del mismo usuario en una ronda no revelada de **otro** proyecto no se toca.
  - Atomicidad: la baja del Scrum Master no borra sus votos de rondas abiertas.
  - Horas: un `registros_horas` del Product Builder sigue existiendo (la FK a `integrantes` no se rompe), y `ListarDatosSprintsFinalizados` del proyecto devuelve las mismas horas antes y después de la baja (SC-002).
  - Confirmar que FALLA.
- [x] T025 [US1] COMMIT `RED: agregar tests de DarDeBajaIntegrante en SQLite (falla)`
- [x] T026 [US1] GREEN: `DarDeBajaIntegrante` en `internal/repository/sqlite_project.go`, en una transacción ([research.md §4](research.md)):
  1. `UPDATE integrantes SET fecha_baja = ? WHERE proyecto_id = ? AND usuario_id = ? AND fecha_baja IS NULL AND rol <> 'scrum_master'`. La fecha va con el helper existente `tiempoOpcionalAValor(&fecha)` (`t.UTC().Format(time.DateOnly)`, la misma convención que `fecha_cierre`), no con `Format` directo, para no guardar el día local.
  2. Si `RowsAffected` es 0: leer la fila; activa y `scrum_master` → `domain.ErrScrumMasterNoSeQuita`; si no → `domain.ErrIntegranteNoEncontrado`; rollback.
  3. `DELETE FROM votos WHERE usuario_id = ? AND ronda_id IN (SELECT r.id FROM rondas_votacion r JOIN historias h ON h.id = r.historia_id WHERE h.proyecto_id = ? AND r.revelada_en IS NULL)`.
  4. Commit. `go test ./...` en verde.
- [x] T027 [US1] COMMIT `GREEN: implementar DarDeBajaIntegrante en SQLite con UPDATE condicional y descarte de votos no revelados, test en verde`

### Servicio `QuitarIntegrante`

- [x] T028 [US1] RED: nuevo `internal/service/quitar_integrante_test.go` (con `fakeProjectRepository`, que desde T015 guarda la fecha de cada baja; con `ProjectRepositoryEnMemoria` la fecha no se puede ver desde el paquete `service`), tabla de casos:
  - Scrum Master quita a un Product Builder → `nil`; ya no está en `ListarIntegrantes`; la fecha guardada en el registro de bajas del fake es la del reloj inyectado.
  - Proyecto inexistente → `ErrProyectoNoEncontrado` (US4).
  - Actor no integrante → `ErrNoEsIntegrante`; actor Product Builder → `ErrNoEsScrumMaster`.
  - Usuario que no integra / ya dado de baja → `ErrIntegranteNoEncontrado` (US4).
  - Scrum Master se quita a sí mismo → `ErrScrumMasterNoSeQuita` y el equipo no cambia (US2).
  - Orden de errores (FR-013): con proyecto inexistente **y** actor no integrante, gana 404; con actor Product Builder **y** objetivo inexistente, gana 403.
  - Proyecto cerrado (guardia `exigirProyectoAbierto` de #300, que llama `exigirScrumMaster`) → `ErrProyectoCerrado`, después de los 403 y antes del 404 de integrante.
  - Confirmar que FALLA (no compila).
- [x] T029 [US1] COMMIT `RED: agregar tests de QuitarIntegrante (falla)`
- [x] T030 [US1] GREEN: nuevo `internal/service/quitar_integrante.go`:
  - `QuitarIntegranteInput{ProyectoID, UsuarioID, ActorID int64}`; `NewQuitarIntegrante(proyectos repository.ProjectRepository, reloj func() time.Time)` (usar el patrón de reloj que ya usen otros servicios).
  - `Ejecutar`: `ObtenerPorID` → `exigirScrumMaster` → buscar al objetivo en `ListarIntegrantes` (no está → `ErrIntegranteNoEncontrado`) → `Member.ValidarBaja()` → `DarDeBajaIntegrante(…, reloj())`.
  - `go test ./...` en verde.
- [x] T031 [US1] COMMIT `GREEN: implementar QuitarIntegrante, test en verde`

### Votación e historial con los dados de baja

- [x] T032b [US1] Tests de regresión. Hacerla **antes** de T032, para no dejar tests en rojo sin commitear en el árbol. La baja se marca con `DarDeBajaIntegrante` del `fakeProjectRepository` (T015):
  - En `internal/service/revelar_votacion_test.go`: con un integrante dado de baja que no votó, la ronda se revela.
  - En `internal/service/consultar_votacion_test.go`: en una ronda abierta donde no votó, la vista no lo lista.

  No es un RED: los dos pasan porque esos servicios ya usan `ListarIntegrantes`, que filtra; los tests protegen esa dependencia. COMMIT `test: votación sin esperar ni listar a un integrante dado de baja (Refs #175)`.
- [x] T032 [P] [US1] RED: en `internal/service/consultar_historial_rondas_test.go`, una ronda revelada con el voto de un Product Builder que después se dio de baja: el historial lo sigue mostrando con su nombre y rol. En `internal/service/consultar_votacion_test.go`: ronda revelada con su voto y baja posterior → la vista lo incluye (con valor) y el promedio coincide. En los dos archivos, la baja se marca con `DarDeBajaIntegrante` del `fakeProjectRepository` (T015), y el nombre sale de `listaMiembros`, que `ListarIntegrantesHistoricos` devuelve sin filtrar. Correr y confirmar que los dos casos FALLAN.
- [x] T033 [US1] COMMIT `RED: agregar tests de votación e historial con integrantes dados de baja (falla)`
- [x] T034 [US1] GREEN: `internal/service/consultar_historial_rondas.go` usa `ListarIntegrantesHistoricos` para resolver nombres; `internal/service/consultar_votacion.go` lista a los activos más quien tenga voto en la ronda actual (nombre y rol desde `ListarIntegrantesHistoricos`), ordenados por id como hoy. `go test ./...` en verde.
- [x] T035 [US1] COMMIT `GREEN: conservar a los dados de baja en el historial y la vista de votación, test en verde`
- [x] T035b [US1] Test de regresión en `internal/service/generar_reporte_test.go`: con un Product Builder dado de baja que imputó horas en un Sprint finalizado, `EjecutarProyecto` no lo incluye en `Integrantes` y las horas y métricas del reporte son las mismas que antes de la baja (spec US1 esc. 4, FR-001, SC-002). La baja se marca con `DarDeBajaIntegrante` del `fakeProjectRepository` (T015), no fijando `listaMiembros` sin esa persona: eso daría un test que no prueba nada. Pasa porque `EjecutarProyecto` ya usa `ListarIntegrantes`, que filtra, y las horas salen de `ListarDatosSprintsFinalizados`, que no depende del equipo. COMMIT `test: el reporte excluye a los integrantes dados de baja sin cambiar sus horas (Refs #175)`.

### HTTP `DELETE /projects/{id}/members/{usuario_id}`

- [x] T036 [US1] RED: en `internal/http/project_handler_test.go` (o nuevo `project_handler_equipo_test.go`), tabla según [contracts/api-equipo.md](contracts/api-equipo.md): 204 sin cuerpo; 400 `id` o `usuario_id` no numérico o ≤ 0; 404 «proyecto no encontrado»; 403 «No sos integrante de este proyecto»; 403 «Solo el Scrum Master puede realizar esta acción»; 409 «El proyecto está cerrado»; 404 «integrante no encontrado»; 409 «No se puede quitar al Scrum Master del proyecto». En `internal/http/rutas_test.go` (o donde se prueben los niveles), la ruta es `SesionScrumMaster` y responde 401 sin sesión. Confirmar que FALLA.
- [x] T037 [US1] COMMIT `RED: agregar tests del handler DELETE /projects/{id}/members/{usuario_id} (falla)`
- [x] T038 [US1] GREEN:
  - `internal/http/project_handler.go`: campo `quitar *service.QuitarIntegrante` y handler `QuitarIntegrante` (`actorDeSesion`, `parsearID` de `id` y de `usuario_id`, `Ejecutar`, `w.WriteHeader(204)`).
  - `internal/http/rutas.go`: `{"DELETE /projects/{id}/members/{usuario_id}", SesionScrumMaster, h.Proyectos.QuitarIntegrante}`.
  - `internal/http/sprint_handler.go` (`escribirError`): `ErrIntegranteNoEncontrado` → 404; `ErrScrumMasterNoSeQuita` → 409.
  - Cablear `service.NewQuitarIntegrante` en `internal/app/app.go` y en el constructor `NewProjectHandler`.
  - `go test ./...` en verde.
- [x] T039 [US1] COMMIT `GREEN: implementar DELETE /projects/{id}/members/{usuario_id}, test en verde`
- [x] T040 [US1] — no hizo falta: no hubo nada que limpiar. (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: US1, US2 (quitar al Scrum Master) y US4 funcionan por API.

---

## Phase 4: User Story 3 - Volver a sumar a un integrante dado de baja (P2) — incluye US2 (cambio de rol)

**Goal**: el alta existente reactiva la misma fila como Product Builder; un integrante activo (incluido el Scrum Master) sigue dando duplicado.

**Independent Test**: `go test ./internal/repository/ ./internal/service/ -run 'Reactiv|AgregarIntegrante|AsignarIntegrante' -v` y [quickstart.md](quickstart.md) paso 9.

- [x] T041 [US3] RED: reactivación en el repositorio.
  - `internal/repository/project_repository_en_memoria_test.go` y `internal/repository/sqlite_project_test.go`: alta de un Product Builder dado de baja → vuelve a `ListarIntegrantes` como `product_builder`, una sola vez, con el `creado_en` original; en SQLite, `SELECT COUNT(*)` de la fila = 1 y `fecha_baja IS NULL`.
  - Alta de un integrante activo → `ErrMiembroDuplicado` (como hoy).
  - Los votos revelados y las horas anteriores a la baja siguen asociados a esa persona.
  - Confirmar que FALLA.
- [x] T042 [US3] COMMIT `RED: agregar tests de reactivación de integrantes (falla)`
- [x] T043 [US3] GREEN: `AgregarIntegrante` en `internal/repository/sqlite_project.go` pasa a `INSERT INTO integrantes (...) VALUES (...) ON CONFLICT (proyecto_id, usuario_id) DO UPDATE SET fecha_baja = NULL, rol = excluded.rol WHERE integrantes.fecha_baja IS NOT NULL`; con `RETURNING creado_en` (`QueryRowContext`). Si no devuelve fila (`sql.ErrNoRows`), la persona ya era integrante activo → `ErrMiembroDuplicado`. Si devuelve fila, ese `creado_en` es el de la membresía: en un alta, el nuevo; en una reactivación, el original, sin una segunda lectura. Mantener la traducción de `SQLITE_CONSTRAINT_UNIQUE` → `ErrScrumMasterYaExiste`. Mismo comportamiento en `project_repository_en_memoria.go`. `go test ./...` en verde.
- [x] T044 [US3] COMMIT `GREEN: reactivar al integrante dado de baja en AgregarIntegrante, test en verde`
- [x] T045 [US3] RED: en `internal/service/asignar_integrante_test.go` (con `fakeProjectRepository`, que desde T015 sabe dar de baja y reactivar; la baja se marca con su `DarDeBajaIntegrante`): (a) alta de un dado de baja sin rol → reactivado como Product Builder; (b) dado de baja con `rol: scrum_master` → `ErrScrumMasterYaExiste` (409) y no se reactiva; (c) el Scrum Master agregándose de nuevo como `product_builder` y como `scrum_master` → error de integrante duplicado (400) y el equipo no cambia (US2, Clarification Q2). Escribir los tres casos, correrlos y separarlos: los que fallan van en el commit RED (T046); los que ya pasan por el GREEN de T043 van en un commit aparte `test: casos de AsignarIntegrante con integrantes dados de baja (Refs #175)`. Si ninguno falla, no hay RED/GREEN: solo ese commit `test:`, y se saltean T046–T048.
- [x] T046 [US3] — no hizo falta: los tres casos de T045 pasaron con el GREEN de T043 y fueron en el commit `test:`. COMMIT `RED: agregar tests de AsignarIntegrante con integrantes dados de baja (falla)`
- [x] T047 [US3] — no hizo falta: `asignar_integrante.go` no necesitó cambios. GREEN: ajustar `internal/service/asignar_integrante.go` solo si T045 lo pide (p. ej. `verificarScrumMasterLibre` con un dado de baja). `go test ./...` en verde. Si no hace falta código, saltear este GREEN y su commit y anotarlo en el PR.
- [x] T048 [US3] — no hizo falta (no hubo T047). COMMIT `GREEN: <ajuste> en AsignarIntegrante, test en verde` (solo si hubo T047)

**Checkpoint**: baja y reactivación completas por API.

---

## Phase 5: BDD y documentación de reglas

- [x] T049 Completar `MODELO-DE-DATOS.md` §5 (dónde vive cada regla): «integrante = activo» (consultas de `ListarIntegrantes`/`ListarPorIntegrante`); «al Scrum Master no se lo quita» (`Member.ValidarBaja` + `UPDATE` condicional, 409); «una baja no se aplica dos veces» (`UPDATE` condicional, 404); «la reactivación reutiliza la fila» (`ON CONFLICT ... WHERE fecha_baja IS NOT NULL`); «al dar de baja se descartan sus votos no revelados» (misma transacción). Agregar el límite conocido: `Votar` comprueba la pertenencia y después llama a `GuardarVoto` en dos pasos separados. Si la baja borra los votos de la persona justo entre esos dos pasos, su voto queda guardado y cuenta en el promedio de esa ronda. FR-014 cubre solo dos bajas al mismo tiempo; se acepta por lo improbable (equipos chicos, un solo Scrum Master). COMMIT `docs: reglas de la baja de integrantes en MODELO-DE-DATOS.md §5 (Refs #175)`.
- [x] T050 Escribir `features/gestion_equipo.feature` (Gherkin en español, `@HU-16`, Característica con el Como/Quiero/Para de la spec, Antecedentes con «Dado que inicié sesión como …» y un proyecto con Scrum Master y un Product Builder) con un escenario por criterio de aceptación de la spec, sin inventar casos:
  - Escenario 2 (normal): quitar a un integrante que votó una ronda revelada e imputó horas → no figura en el equipo, una ronda nueva se revela sin su voto, el historial y las horas lo siguen mostrando, y las métricas del proyecto son las mismas antes y después de la baja (reutilizar «Cuando consulto las métricas del proyecto» / «Entonces el Sprint "S1" informa estas métricas:» de `calculo_metricas.feature`; SC-002).
  - Escenario 2 (alternativo, US1.2): el integrante quitado intenta abrir el proyecto, votar e imputar horas → 403 «No sos integrante de este proyecto» (`Esquema del escenario` con las tres acciones).
  - Escenario 2 (alternativo, US1.3): el integrante quitado consulta «Mis proyectos» → el proyecto no aparece (reutilizar los pasos de `steps_listado_proyectos.go`).
  - Edge case (límite): baja con una ronda abierta donde ya votó → al revelar, su voto no está.
  - Edge case (límite): baja con una ronda abierta y reactivación con la ronda todavía abierta → puede volver a votar.
  - Escenario 3 (alternativo, US3.1 y US3.2): reactivarlo → figura una sola vez como Product Builder, y en «Mis proyectos» el proyecto vuelve a aparecer con el rol Product Builder.
  - Escenario 3 (error, US3.3): agregar de nuevo a un Product Builder activo → 400 «el integrante ya pertenece al proyecto».
  - Escenario 4 (error): quitar al Scrum Master → 409; agregarlo de nuevo como Product Builder → 400; el equipo no cambia.
  - Escenario 5 (error, US4.1 a US4.3): proyecto inexistente → 404; usuario que nunca integró el proyecto → 404; integrante ya dado de baja que se quiere quitar otra vez → 404.
- [x] T051 Step definitions en `features/steps_gestion_equipo.go`, reutilizando los pasos de sesión, votación y horas que ya existen (`steps_autenticacion.go`, `steps_votacion.go`, `steps_horas.go`, `steps_scrum_master_unico.go`). `go test ./features/` en verde, incluidos los `.feature` existentes (`permisos_por_rol` cubre el Escenario 1).
- [x] T052 COMMIT `Agregar escenarios BDD de gestión del equipo (Refs #175)`

---

## Phase 6: User Story 6 - Quitar desde la pantalla del equipo (P3, frontend)

**Goal**: el Scrum Master quita a un Product Builder desde «Configurar equipo».

**Independent Test**: `cd web && npm test -- Configuracion`.

**Fuera de esta rama (decisión del 2026-11-12)**: el frontend no entra en `feature/175-gestion-equipo`. T053–T056 pasan a la sub-issue #182, asignada a Ana Iriarte Lopez, que los hace en su propia rama. Según la Definition of Done, la sub-issue de Frontend no frena el cierre de #175.

- [ ] T053 [US6] RED: en `web/src/features/proyecto/Configuracion.test.tsx` (escenario en el nombre de cada prueba): cada Product Builder tiene «Quitar»; el Scrum Master no; al confirmar, se llama a `quitarIntegrante` y la fila desaparece; al cancelar, no se llama; si el servidor responde error, se muestra su mensaje y la fila sigue. Confirmar que FALLA.
- [ ] T054 [US6] COMMIT `RED: agregar pruebas de quitar integrante en Configurar equipo (falla)`
- [ ] T055 [US6] GREEN: `quitarIntegrante(id, usuarioId)` (`DELETE /projects/${id}/members/${usuarioId}`) en `web/src/features/proyecto/api.ts` (con su prueba en `api.test.ts` si existe el patrón) y el botón con confirmación en `web/src/features/proyecto/Configuracion.tsx`. No se deshabilita el botón para los Product Builders: ese escenario quedó fuera de alcance (spec, «Fuera de alcance») y va con el frontend de HU-21 (#247). `npm test`, `npm run lint` y `npm run build` en verde.
- [ ] T056 [US6] COMMIT `GREEN: quitar integrante desde Configurar equipo, test en verde`

---

## Phase 7: Polish & Cross-Cutting Concerns

- [x] T057 `go build ./...`, `go vet ./...`, `gofmt -l .` (vacío) y `go test ./...` sin errores.
- [ ] T058 Correr la validación manual de [quickstart.md](quickstart.md) con el servidor local.
- [x] T059 Actualizar `docs/PRODUCT-BACKLOG.md` (HU-16: estado, nota de que el Escenario 1 lo cubre HU-21, D15 cerrada, sin cambio de rol, votos no revelados descartados) y la guía de curl si lista las rutas del equipo. COMMIT `docs: HU-16 implementada en el backlog (Refs #175)`.
- [x] T060 Marcar las tareas hechas en este archivo. COMMIT `docs: marcar las tareas de la spec 045 como hechas (Refs #175)`.

---

## Dependencies & Execution Order

- **Setup (T001–T002)** → **Foundational (T003–T019)** → **US1+US2+US4 (T020–T040)** → **US3 (T041–T048)** → **BDD y docs (T049–T052)** → **US6 frontend (T053–T056)** → **Polish (T057–T060)**.
- T001 bloquea todo lo que toca el esquema: sin la 008 en `main`, no se crea la 009.
- T013 (puerto) antes de T015; T017 antes de T020; T015 (fake con bajas) antes de T028, T032b, T032, T035b y T045; T032b antes de T032; T030 antes de T036.
- US3 depende de la baja (T026) para poder probar la reactivación.
- Godog (T050–T051) después de US1 y US3, con `.feature` y pasos en el mismo commit (`Strict: true`).
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.

### Paralelizables

- T009–T012 (dominio) en paralelo con T003–T008 (esquema): archivos distintos.
- T032b/T032 (tests de votación e historial) en paralelo con T028 (tests del servicio de baja).
- US6 (frontend) puede empezar apenas exista el contrato (T038), en paralelo con US3.

## Implementation Strategy

1. Setup + Foundational → la columna existe y todo lee solo activos.
2. US1 (con US2 y US4) → MVP: quitar por API conservando el historial.
3. US3 → reactivar.
4. BDD → escenarios del backlog automatizados.
5. US6 → pantalla.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- El fake en memoria no reproduce el descarte de votos (los votos viven en otro fake): esa garantía se prueba en SQLite (T024) y en Godog (T050).
