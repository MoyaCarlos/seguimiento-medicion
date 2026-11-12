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

- [ ] T001 Verificar que el PR #311 (#300, migración `008_cierre_proyecto.sql`) está mergeado en `origin/main` (`git fetch` y `git ls-tree origin/main internal/repository/migraciones/`). **Si no lo está, frenar y avisar**: no se renumera ni se apila sobre `feature/300-cerrar-reabrir-proyecto` por cuenta propia.
- [ ] T002 Mergear `origin/main` en `feature/175-gestion-equipo` (`git merge origin/main`, merge commit) y correr `go test ./...` en verde.

---

## Phase 2: Foundational — columna `fecha_baja` (migración 009), errores y «integrante = activo»

**Goal**: la base guarda la baja, el dominio conoce sus errores y regla, y toda lectura de integrantes devuelve solo activos. Todavía nadie da de baja ([data-model.md](data-model.md)).

**Independent Test**: `go test ./internal/repository/ ./internal/domain/ -v`.

### Esquema

- [ ] T003 RED: esquema y migración 009.
  - En `internal/repository/esquema_test.go`, `TestEsquema_Filas`: casos D12 de `integrantes.fecha_baja` — `'2026-11-1'`, `'basura'` y `'2026-02-30'` → `violaCheck`; `'2026-11-12'` y `NULL` → `aceptada`.
  - En `TestEsquema_Estructura`, la columna `fecha_baja` en `integrantes`.
  - Nuevo `internal/repository/migracion_009_test.go`, con el patrón de `migracion_007_test.go`: helper `baseAntesDeLa009` (todas las migraciones embebidas salvo `009_`); siembra un proyecto con un Scrum Master y un Product Builder, aplica todas y verifica que ambos siguen con `fecha_baja IS NULL`.
  - Correr `go test ./internal/repository/` y confirmar que FALLA.
- [ ] T004 COMMIT `RED: agregar tests de esquema de integrantes.fecha_baja y migración 009 (falla)` (solo `_test.go`)
- [ ] T005 GREEN: crear `internal/repository/migraciones/009_baja_integrantes.sql`: cabecera `-- Baja lógica de integrantes: MODELO-DE-DATOS.md §4.3 (D15, spec 045, #175).` y `ALTER TABLE integrantes ADD COLUMN fecha_baja TEXT CHECK (fecha_baja IS NULL OR date(fecha_baja) IS fecha_baja);`. `go test ./...` en verde.
- [ ] T006 COMMIT `GREEN: migración 009 con integrantes.fecha_baja, test en verde`
- [ ] T007 Documentar el esquema en `MODELO-DE-DATOS.md`, en el mismo PR que la migración (regla de AGENTS.md):
  - §2: `fecha_baja` en las fechas de calendario y en el `CHECK` D12.
  - §3: `text fecha_baja` en `integrantes` del diagrama.
  - §4.3: fila `fecha_baja | TEXT? | Nueva. HU-16. YYYY-MM-DD; NULL ⇔ activo` y párrafo de baja lógica: quitar = fecha de baja; reactivar = misma fila con `fecha_baja = NULL`; el Scrum Master no se da de baja; «integrante» en las reglas = integrante activo.
  - §4.7: al dar de baja se borran los votos de esa persona en rondas no reveladas del proyecto; las reveladas no se tocan.
  - §4.10: fila de `date(col) IS col` con `fecha_baja`.
  - §7: D15 **cerrada el 2026-11-12** (confirmada por el equipo), con la decisión de la spec 045.
  - §8: fila de HU-16 → `integrantes.fecha_baja`, migración 009.
  - §9: fila `009 | integrantes.fecha_baja | HU-16 (spec 045)` y los casos de `fecha_baja` en el punto de tests de esquema.
  - Las filas de §5 van en T049, cuando existan las reglas.
- [ ] T008 COMMIT `docs: integrantes.fecha_baja en MODELO-DE-DATOS.md, cierra D15 (Refs #175)`

### Dominio

- [ ] T009 [P] RED: en `internal/domain/member_test.go`, `TestMember_ValidarBaja`: un `Member` con `Role: RolScrumMaster` → `errors.Is(err, ErrScrumMasterNoSeQuita)`; con `RolProductBuilder` → `nil`. Confirmar que FALLA (no compila).
- [ ] T010 COMMIT `RED: agregar test de Member.ValidarBaja (falla)`
- [ ] T011 GREEN: en `internal/domain/errors.go`, `ErrIntegranteNoEncontrado = errors.New("integrante no encontrado")` y `ErrScrumMasterNoSeQuita = errors.New("No se puede quitar al Scrum Master del proyecto")`; en `internal/domain/member.go`, `func (m Member) ValidarBaja() error`. Tests en verde.
- [ ] T012 COMMIT `GREEN: implementar Member.ValidarBaja, test en verde`

### Puerto y lecturas de activos

- [ ] T013 Ampliar el puerto `ProjectRepository` en `internal/repository/project_repository.go`:
  - `ListarIntegrantesHistoricos(ctx, proyectoID int64) ([]domain.Member, error)` — activos y dados de baja, orden de alta.
  - `DarDeBajaIntegrante(ctx, proyectoID, usuarioID int64, fecha time.Time) error`.
  - Comentario de `ListarIntegrantes` y `ListarPorIntegrante`: «solo integrantes activos (`fecha_baja IS NULL`)».
  - Agregar los dos métodos a todos los que implementan el puerto, devolviendo lo mínimo para compilar (`nil`/`errors.New("no implementado")` en los repos reales hasta su GREEN): `internal/repository/project_repository_en_memoria.go`, `internal/repository/sqlite_project.go`, `internal/service/fakes_test.go`, `internal/service/iniciar_sprint_test.go`, `internal/http/project_handler_test.go`. `go build ./... && go vet ./...` sin errores.
- [ ] T014 COMMIT `Agregar ListarIntegrantesHistoricos y DarDeBajaIntegrante al puerto ProjectRepository`
- [ ] T015 RED: lecturas de activos en el fake y en SQLite.
  - `internal/repository/project_repository_en_memoria_test.go` y `internal/repository/sqlite_project_test.go`: con un Product Builder marcado como dado de baja (en SQLite, `UPDATE integrantes SET fecha_baja = '2026-11-12'` directo; en memoria, con `DarDeBajaIntegrante` una vez implementado o un helper de test sobre el estado interno), verificar que `ListarIntegrantes` no lo devuelve y `ListarIntegrantesHistoricos` sí (con nombre y rol).
  - `internal/repository/sqlite_project_listar_test.go` y `project_repository_en_memoria_listar_test.go`: `ListarPorIntegrante` no devuelve el proyecto donde la persona está dada de baja.
  - Confirmar que FALLA.
- [ ] T016 COMMIT `RED: agregar tests de integrantes activos e históricos (falla)`
- [ ] T017 GREEN: en `internal/repository/sqlite_project.go`, `AND i.fecha_baja IS NULL` en `ListarIntegrantes` y `ListarPorIntegrante`; `ListarIntegrantesHistoricos` = la consulta sin el filtro. En el fake en memoria, un campo de baja por vínculo (p. ej. `fechaBaja *time.Time` en un tipo interno, sin tocar `domain.Membership`) y el mismo filtrado. `go test ./...` en verde.
- [ ] T018 COMMIT `GREEN: listar solo integrantes activos y agregar ListarIntegrantesHistoricos, test en verde`
- [ ] T019 (solo si hay algo que limpiar, p. ej. consulta duplicada entre activos e históricos) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: la base tiene `fecha_baja`, las reglas existentes (pertenencia, rol, votación, horas, reporte, «Mis proyectos») ya ignoran a quien tenga fecha de baja, y `go test ./...` sigue en verde.

---

## Phase 3: User Story 1 - Quitar un integrante conservando su historial (P1) 🎯 MVP — incluye US2 y US4

**Goal**: el Scrum Master quita a un Product Builder con baja lógica; sus aportes se conservan; el Scrum Master no se puede quitar; inexistentes dan 404.

**Independent Test**: `go test ./internal/repository/ ./internal/service/ ./internal/http/ -run 'Baja|Quitar' -v` y la prueba manual de [quickstart.md](quickstart.md) pasos 1–8 y 10.

### Repositorio: `DarDeBajaIntegrante`

- [ ] T020 [US1] RED: en `internal/repository/project_repository_en_memoria_test.go`, `TestProjectRepositoryEnMemoria_DarDeBajaIntegrante`:
  - Product Builder activo → `nil`; deja de salir en `ListarIntegrantes` y sigue en `ListarIntegrantesHistoricos`.
  - Segunda baja del mismo → `ErrIntegranteNoEncontrado`.
  - Usuario que nunca integró → `ErrIntegranteNoEncontrado`.
  - Scrum Master → `ErrScrumMasterNoSeQuita` y sigue activo.
  - Confirmar que FALLA.
- [ ] T021 [US1] COMMIT `RED: agregar tests de DarDeBajaIntegrante en memoria (falla)`
- [ ] T022 [US1] GREEN: `DarDeBajaIntegrante` en `internal/repository/project_repository_en_memoria.go` (bajo el mutex). Tests en verde.
- [ ] T023 [US1] COMMIT `GREEN: implementar DarDeBajaIntegrante en memoria, test en verde`
- [ ] T024 [US1] RED: en `internal/repository/sqlite_project_test.go` (o nuevo `sqlite_project_baja_test.go`), `TestSQLiteProjectRepository_DarDeBajaIntegrante`:
  - Mismos cuatro casos que T020; la fila queda con `fecha_baja = '2026-11-12'` (formato `YYYY-MM-DD`) y no se borra.
  - Votos: con una historia del proyecto, una ronda **revelada** y otra **no revelada** donde votó el Product Builder (y votó también otra persona), después de la baja el voto de la ronda no revelada no existe y el de la revelada sigue igual; el voto de la otra persona sigue en las dos. Un voto del mismo usuario en una ronda no revelada de **otro** proyecto no se toca.
  - Atomicidad: la baja del Scrum Master no borra sus votos de rondas abiertas.
  - Horas: un `registros_horas` del Product Builder sigue existiendo (la FK a `integrantes` no se rompe).
  - Confirmar que FALLA.
- [ ] T025 [US1] COMMIT `RED: agregar tests de DarDeBajaIntegrante en SQLite (falla)`
- [ ] T026 [US1] GREEN: `DarDeBajaIntegrante` en `internal/repository/sqlite_project.go`, en una transacción ([research.md §4](research.md)):
  1. `UPDATE integrantes SET fecha_baja = ? WHERE proyecto_id = ? AND usuario_id = ? AND fecha_baja IS NULL AND rol <> 'scrum_master'` (fecha con `Format(time.DateOnly)`).
  2. Si `RowsAffected` es 0: leer la fila; activa y `scrum_master` → `domain.ErrScrumMasterNoSeQuita`; si no → `domain.ErrIntegranteNoEncontrado`; rollback.
  3. `DELETE FROM votos WHERE usuario_id = ? AND ronda_id IN (SELECT r.id FROM rondas_votacion r JOIN historias h ON h.id = r.historia_id WHERE h.proyecto_id = ? AND r.revelada_en IS NULL)`.
  4. Commit. `go test ./...` en verde.
- [ ] T027 [US1] COMMIT `GREEN: implementar DarDeBajaIntegrante en SQLite con UPDATE condicional y descarte de votos no revelados, test en verde`

### Servicio `QuitarIntegrante`

- [ ] T028 [US1] RED: nuevo `internal/service/quitar_integrante_test.go` (con `ProjectRepositoryEnMemoria`), tabla de casos:
  - Scrum Master quita a un Product Builder → `nil`; ya no está en `ListarIntegrantes`; la fecha es la del reloj inyectado.
  - Proyecto inexistente → `ErrProyectoNoEncontrado` (US4).
  - Actor no integrante → `ErrNoEsIntegrante`; actor Product Builder → `ErrNoEsScrumMaster`.
  - Usuario que no integra / ya dado de baja → `ErrIntegranteNoEncontrado` (US4).
  - Scrum Master se quita a sí mismo → `ErrScrumMasterNoSeQuita` y el equipo no cambia (US2).
  - Orden de errores (FR-013): con proyecto inexistente **y** actor no integrante, gana 404; con actor Product Builder **y** objetivo inexistente, gana 403.
  - Si el repositorio tiene la guardia de proyecto cerrado (#300 mergeado): proyecto cerrado → `ErrProyectoCerrado`, después de los 403 y antes del 404 de integrante.
  - Confirmar que FALLA (no compila).
- [ ] T029 [US1] COMMIT `RED: agregar tests de QuitarIntegrante (falla)`
- [ ] T030 [US1] GREEN: nuevo `internal/service/quitar_integrante.go`:
  - `QuitarIntegranteInput{ProyectoID, UsuarioID, ActorID int64}`; `NewQuitarIntegrante(proyectos repository.ProjectRepository, reloj func() time.Time)` (usar el patrón de reloj que ya usen otros servicios).
  - `Ejecutar`: `ObtenerPorID` → `exigirScrumMaster` → buscar al objetivo en `ListarIntegrantes` (no está → `ErrIntegranteNoEncontrado`) → `Member.ValidarBaja()` → `DarDeBajaIntegrante(…, reloj())`.
  - `go test ./...` en verde.
- [ ] T031 [US1] COMMIT `GREEN: implementar QuitarIntegrante, test en verde`

### Votación e historial con los dados de baja

- [ ] T032 [P] [US1] RED: en `internal/service/consultar_historial_rondas_test.go`, una ronda revelada con el voto de un Product Builder que después se dio de baja: el historial lo sigue mostrando con su nombre y rol. En `internal/service/consultar_votacion_test.go`: ronda revelada con su voto y baja posterior → la vista lo incluye (con valor) y el promedio coincide; ronda abierta donde no votó → la vista no lo lista. En `internal/service/revelar_votacion_test.go`: con un integrante dado de baja que no votó, la ronda se revela. Confirmar que FALLA lo que corresponda (historial y vista; revelar ya pasa por el filtro de T017 — si pasa al escribirlo, dejarlo como test de regresión en el mismo commit RED y anotarlo).
- [ ] T033 [US1] COMMIT `RED: agregar tests de votación e historial con integrantes dados de baja (falla)`
- [ ] T034 [US1] GREEN: `internal/service/consultar_historial_rondas.go` usa `ListarIntegrantesHistoricos` para resolver nombres; `internal/service/consultar_votacion.go` lista a los activos más quien tenga voto en la ronda actual (nombre y rol desde `ListarIntegrantesHistoricos`), ordenados por id como hoy. `go test ./...` en verde.
- [ ] T035 [US1] COMMIT `GREEN: conservar a los dados de baja en el historial y la vista de votación, test en verde`

### HTTP `DELETE /projects/{id}/members/{usuario_id}`

- [ ] T036 [US1] RED: en `internal/http/project_handler_test.go` (o nuevo `project_handler_equipo_test.go`), tabla según [contracts/api-equipo.md](contracts/api-equipo.md): 204 sin cuerpo; 400 `usuario_id` no numérico o ≤ 0; 404 «proyecto no encontrado»; 403 «No sos integrante de este proyecto»; 403 «Solo el Scrum Master puede realizar esta acción»; 404 «integrante no encontrado»; 409 «No se puede quitar al Scrum Master del proyecto». En `internal/http/rutas_test.go` (o donde se prueben los niveles), la ruta es `SesionScrumMaster` y responde 401 sin sesión. Confirmar que FALLA.
- [ ] T037 [US1] COMMIT `RED: agregar tests del handler DELETE /projects/{id}/members/{usuario_id} (falla)`
- [ ] T038 [US1] GREEN:
  - `internal/http/project_handler.go`: campo `quitar *service.QuitarIntegrante` y handler `QuitarIntegrante` (`actorDeSesion`, `parsearID` de `id` y de `usuario_id`, `Ejecutar`, `w.WriteHeader(204)`).
  - `internal/http/rutas.go`: `{"DELETE /projects/{id}/members/{usuario_id}", SesionScrumMaster, h.Proyectos.QuitarIntegrante}`.
  - `internal/http/sprint_handler.go` (`escribirError`): `ErrIntegranteNoEncontrado` → 404; `ErrScrumMasterNoSeQuita` → 409.
  - Cablear `service.NewQuitarIntegrante` en `internal/app/app.go` y en el constructor `NewProjectHandler`.
  - `go test ./...` en verde.
- [ ] T039 [US1] COMMIT `GREEN: implementar DELETE /projects/{id}/members/{usuario_id}, test en verde`
- [ ] T040 [US1] (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: US1, US2 (quitar al Scrum Master) y US4 funcionan por API.

---

## Phase 4: User Story 3 - Volver a sumar a un integrante dado de baja (P2) — incluye US2 (cambio de rol)

**Goal**: el alta existente reactiva la misma fila como Product Builder; un integrante activo (incluido el Scrum Master) sigue dando duplicado.

**Independent Test**: `go test ./internal/repository/ ./internal/service/ -run 'Reactiv|AgregarIntegrante|AsignarIntegrante' -v` y [quickstart.md](quickstart.md) paso 9.

- [ ] T041 [US3] RED: reactivación en el repositorio.
  - `internal/repository/project_repository_en_memoria_test.go` y `internal/repository/sqlite_project_test.go`: alta de un Product Builder dado de baja → vuelve a `ListarIntegrantes` como `product_builder`, una sola vez, con el `creado_en` original; en SQLite, `SELECT COUNT(*)` de la fila = 1 y `fecha_baja IS NULL`.
  - Alta de un integrante activo → `ErrMiembroDuplicado` (como hoy).
  - Los votos revelados y las horas anteriores a la baja siguen asociados a esa persona.
  - Confirmar que FALLA.
- [ ] T042 [US3] COMMIT `RED: agregar tests de reactivación de integrantes (falla)`
- [ ] T043 [US3] GREEN: `AgregarIntegrante` en `internal/repository/sqlite_project.go` pasa a `INSERT INTO integrantes (...) VALUES (...) ON CONFLICT (proyecto_id, usuario_id) DO UPDATE SET fecha_baja = NULL, rol = excluded.rol WHERE integrantes.fecha_baja IS NOT NULL`; `RowsAffected` 0 → `ErrMiembroDuplicado`; devolver la membresía con el `creado_en` guardado (leerlo si se reactivó). Mantener la traducción de `SQLITE_CONSTRAINT_UNIQUE` → `ErrScrumMasterYaExiste`. Mismo comportamiento en `project_repository_en_memoria.go`. `go test ./...` en verde.
- [ ] T044 [US3] COMMIT `GREEN: reactivar al integrante dado de baja en AgregarIntegrante, test en verde`
- [ ] T045 [US3] RED: en `internal/service/asignar_integrante_test.go`: (a) alta de un dado de baja sin rol → reactivado como Product Builder; (b) dado de baja con `rol: scrum_master` → `ErrScrumMasterYaExiste` (409) y no se reactiva; (c) el Scrum Master agregándose de nuevo como `product_builder` y como `scrum_master` → error de integrante duplicado (400) y el equipo no cambia (US2, Clarification Q2). Confirmar que FALLA lo que corresponda; los casos que ya pasan por el GREEN de T043 quedan como regresión en el mismo commit y se anota.
- [ ] T046 [US3] COMMIT `RED: agregar tests de AsignarIntegrante con integrantes dados de baja (falla)`
- [ ] T047 [US3] GREEN: ajustar `internal/service/asignar_integrante.go` solo si T045 lo pide (p. ej. `verificarScrumMasterLibre` con un dado de baja). `go test ./...` en verde. Si no hace falta código, saltear este GREEN y su commit y anotarlo en el PR.
- [ ] T048 [US3] COMMIT `GREEN: <ajuste> en AsignarIntegrante, test en verde` (solo si hubo T047)

**Checkpoint**: baja y reactivación completas por API.

---

## Phase 5: BDD y documentación de reglas

- [ ] T049 Completar `MODELO-DE-DATOS.md` §5 (dónde vive cada regla): «integrante = activo» (consultas de `ListarIntegrantes`/`ListarPorIntegrante`); «al Scrum Master no se lo quita» (`Member.ValidarBaja` + `UPDATE` condicional, 409); «una baja no se aplica dos veces» (`UPDATE` condicional, 404); «la reactivación reutiliza la fila» (`ON CONFLICT ... WHERE fecha_baja IS NOT NULL`); «al dar de baja se descartan sus votos no revelados» (misma transacción). COMMIT `docs: reglas de la baja de integrantes en MODELO-DE-DATOS.md §5 (Refs #175)`.
- [ ] T050 Escribir `features/gestion_equipo.feature` (Gherkin en español, `@HU-16`, Característica con el Como/Quiero/Para de la spec, Antecedentes con «Dado que inicié sesión como …» y un proyecto con Scrum Master y un Product Builder) con un escenario por criterio de aceptación de la spec, sin inventar casos:
  - Escenario 2 (normal): quitar a un integrante que votó una ronda revelada e imputó horas → no figura en el equipo, una ronda nueva se revela sin su voto, el historial y las horas lo siguen mostrando.
  - Escenario 2 (alternativo): el integrante quitado intenta abrir el proyecto → 403.
  - Edge case (límite): baja con una ronda abierta donde ya votó → al revelar, su voto no está.
  - Escenario 3 (alternativo): reactivarlo → figura una sola vez como Product Builder.
  - Escenario 4 (error): quitar al Scrum Master → 409; agregarlo de nuevo como Product Builder → 400; el equipo no cambia.
  - Escenario 5 (error): proyecto inexistente → 404; integrante inexistente → 404.
- [ ] T051 Step definitions en `features/steps_gestion_equipo.go`, reutilizando los pasos de sesión, votación y horas que ya existen (`steps_autenticacion.go`, `steps_votacion.go`, `steps_horas.go`, `steps_scrum_master_unico.go`). `go test ./features/` en verde, incluidos los `.feature` existentes (`permisos_por_rol` cubre el Escenario 1).
- [ ] T052 COMMIT `Agregar escenarios BDD de gestión del equipo (Refs #175)`

---

## Phase 6: User Story 6 - Quitar desde la pantalla del equipo (P3, frontend)

**Goal**: el Scrum Master quita a un Product Builder desde «Configurar equipo».

**Independent Test**: `cd web && npm test -- Configuracion`.

- [ ] T053 [US6] RED: en `web/src/features/proyecto/Configuracion.test.tsx` (escenario en el nombre de cada prueba): cada Product Builder tiene «Quitar»; el Scrum Master no; al confirmar, se llama a `quitarIntegrante` y la fila desaparece; al cancelar, no se llama; si el servidor responde error, se muestra su mensaje y la fila sigue. Confirmar que FALLA.
- [ ] T054 [US6] COMMIT `RED: agregar pruebas de quitar integrante en Configurar equipo (falla)`
- [ ] T055 [US6] GREEN: `quitarIntegrante(id, usuarioId)` (`DELETE /projects/${id}/members/${usuarioId}`) en `web/src/features/proyecto/api.ts` (con su prueba en `api.test.ts` si existe el patrón) y el botón con confirmación en `web/src/features/proyecto/Configuracion.tsx`. Si el frontend de HU-21 (#247) ya deshabilita acciones de gestión para Product Builders, aplicar el mismo patrón al botón. `npm test`, `npm run lint` y `npm run build` en verde.
- [ ] T056 [US6] COMMIT `GREEN: quitar integrante desde Configurar equipo, test en verde`

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T057 `go build ./...`, `go vet ./...`, `gofmt -l .` (vacío) y `go test ./...` sin errores.
- [ ] T058 Correr la validación manual de [quickstart.md](quickstart.md) con el servidor local.
- [ ] T059 Actualizar `docs/PRODUCT-BACKLOG.md` (HU-16: estado, nota de que el Escenario 1 lo cubre HU-21, D15 cerrada, sin cambio de rol, votos no revelados descartados) y la guía de curl si lista las rutas del equipo. COMMIT `docs: HU-16 implementada en el backlog (Refs #175)`.
- [ ] T060 Marcar las tareas hechas en este archivo. COMMIT `docs: marcar las tareas de la spec 045 como hechas (Refs #175)`.

---

## Dependencies & Execution Order

- **Setup (T001–T002)** → **Foundational (T003–T019)** → **US1+US2+US4 (T020–T040)** → **US3 (T041–T048)** → **BDD y docs (T049–T052)** → **US6 frontend (T053–T056)** → **Polish (T057–T060)**.
- T001 bloquea todo lo que toca el esquema: sin la 008 en `main`, no se crea la 009.
- T013 (puerto) antes de T015; T017 antes de T020; T022/T026 antes de T028; T030 antes de T036.
- US3 depende de la baja (T026) para poder probar la reactivación.
- Godog (T050–T051) después de US1 y US3, con `.feature` y pasos en el mismo commit (`Strict: true`).
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.

### Paralelizables

- T009–T012 (dominio) en paralelo con T003–T008 (esquema): archivos distintos.
- T032 (tests de votación e historial) en paralelo con T028 (tests del servicio de baja).
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
