---

description: "Tareas de la adopción del modelo de datos (Issue #107): migraciones numeradas, esquema 001, renombre y traducción de errores"
---

# Tasks: Adopción del modelo de datos — esquema SQLite con migraciones numeradas (Issue #107)

**Input**: `specs/006-esquema-migraciones/` — [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/), [quickstart.md](quickstart.md)

**Fuente de diseño**: [MODELO-DE-DATOS.md](../../MODELO-DE-DATOS.md). Ninguna tarea inventa tablas, columnas ni restricciones: el DDL y la estructura esperada se copian de §4.1–§4.5. Si aparece una contradicción entre el documento y el código, se frena y se consulta al Scrum Master (FR-017).

**Tests**: OBLIGATORIOS (constitución, Principio I). Ciclo RED → COMMIT → GREEN → COMMIT por fase; `REFACTOR:` solo en la fase 4.

**Organización**: por **fases del plan** (0b, 1, 2, 3a, 3b, 3c, 3d, 4, 5, 6), que siguen el orden de MODELO-DE-DATOS.md §10. Cada tarea lleva la User Story de la spec que cubre: US1 runner, US2 esquema 001, US3 renombre, US4 traducción de errores, US5 base vieja.

**Reglas de commit (todas las tareas COMMIT)**:
- `git add` **siempre por nombre de archivo**. Nunca `git add .`, `git add -A`, `git add docs/` ni `git add specs/` (hay archivos sin trackear en el repo que no van en este PR).
- Ninguna tarea hace `git push`, abre el PR ni mergea: eso lo hace Carlos.
- Cada mensaje termina con la línea `Co-Authored-By` que corresponda al agente que commitea.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: se puede hacer en paralelo (archivos distintos, sin dependencias pendientes)
- **[Story]**: User Story de spec.md (US1–US5)

---

## Fase 0b: Spec y documento (sin ciclo RED/GREEN)

**Objetivo**: analizar la consistencia de la spec, dejarla versionada y corregir MODELO-DE-DATOS.md **antes** de escribir tests, para que la estructura esperada se copie de un documento ya corregido (FR-016a, CA-14).

- [ ] T001 Correr `/speckit-analyze` sobre `specs/006-esquema-migraciones/` y guardar el informe en specs/006-esquema-migraciones/analisis.md. **Si hay hallazgos CRITICAL o HIGH, frenar y consultar al Scrum Master antes de seguir**; no se resuelven por cuenta propia.
- [ ] T002 COMMIT de los artefactos de la spec, agregados por nombre: `git add specs/006-esquema-migraciones/spec.md specs/006-esquema-migraciones/plan.md specs/006-esquema-migraciones/research.md specs/006-esquema-migraciones/data-model.md specs/006-esquema-migraciones/quickstart.md specs/006-esquema-migraciones/tasks.md specs/006-esquema-migraciones/contracts/migraciones.md specs/006-esquema-migraciones/contracts/errores-http.md specs/006-esquema-migraciones/checklists/requirements.md specs/006-esquema-migraciones/analisis.md` — mensaje `docs: agregar la especificación SDD de la adopción del modelo de datos (Refs #107)`
- [ ] T003 Corregir MODELO-DE-DATOS.md (FR-016a), solo estas cinco cosas:
  - §4.5: reemplazar `CREATE UNIQUE INDEX ... ON sprints` por los nombres `sprints_un_activo_por_proyecto` (`WHERE estado = 'Activo'`) y `sprints_un_pendiente_por_proyecto` (`WHERE estado = 'Pendiente'`).
  - §5, fila "Índice único parcial de `sprints`" y párrafo siguiente: el mensaje es el del servicio, `"el proyecto ya tiene un Sprint activo"` / `"el proyecto ya tiene un Sprint pendiente"` (en lugar de "ya hay un Sprint activo/pendiente").
  - §9 punto 1: la tabla de control la crea el runner antes de leerla, con `CREATE TABLE IF NOT EXISTS migraciones_aplicadas (version INTEGER PRIMARY KEY, aplicada_en TEXT NOT NULL)`, `aplicada_en` en RFC 3339 UTC.
  - §9 punto 5 y §10 (fila de la fase 5): una base es vieja si tiene al menos una de `projects`, `users`, `project_members` o `backlog_items`; `sprints` no cuenta.
  - §10, párrafo después de la tabla de fases: la frase "quede en rojo desde la fase 2 hasta terminar la 3d" pasa a decir que el rojo arranca en el `GREEN:` de la fase 3a (la fase 2 sigue en verde porque `Migrar` vive hasta la fase 5) y termina al cerrar la 3d.
- [ ] T004 COMMIT `git add MODELO-DE-DATOS.md` — mensaje `docs: alinear MODELO-DE-DATOS.md con la spec 006 (índices, mensajes, tabla de control, base vieja y ventana en rojo) (Refs #107)` (CA-14)

**Checkpoint**: el documento ya dice lo que la implementación va a hacer.

---

## Fase 1: Runner de migraciones (US1)

**Objetivo**: aplicar migraciones numeradas en orden, registrarlas, tolerar huecos y revertir una que falle (FR-001, FR-002; research R1–R3).

**Test independiente**: `go test ./internal/repository -run TestAplicarMigraciones -v` con migraciones de juguete en `fstest.MapFS`, sin la 001.

> **Nota técnica**: en esta fase **no** se agrega `//go:embed` ni la función exportada `AplicarMigraciones`: un patrón `go:embed` que no encuentra archivos no compila, y la 001 recién llega en la fase 2. Acá se implementa solo `aplicarMigraciones(ctx, db, fsys fs.FS)`.

- [ ] T005 [US1] RED: agregar en internal/repository/migrate_test.go tests de `aplicarMigraciones(ctx, db, fsys)` usando `fstest.MapFS` y `AbrirSQLite(":memory:")` (FR-002, US1 escenarios 1–5):
  - `TestAplicarMigraciones_EnOrdenYRegistra`: `001_crear.sql` crea `t1(x)` y `002_insertar.sql` inserta en `t1` (solo funciona en orden); `migraciones_aplicadas` queda con versiones 1 y 2 y cada `aplicada_en` se parsea con `time.RFC3339`.
  - `TestAplicarMigraciones_Idempotente`: correr dos veces; sin error y `t1` sigue con una sola fila.
  - `TestAplicarMigraciones_ToleraHuecos`: con `001` y `003` (sin `002`) se aplican las dos.
  - `TestAplicarMigraciones_SoloLasNuevas`: aplicar `{001}`, después `{001, 003}`; `001` (un `CREATE TABLE` sin `IF NOT EXISTS`) no se re-ejecuta y `003` sí.
  - `TestAplicarMigraciones_FallaALaMitad`: `001` válida y `002_rota.sql` = `CREATE TABLE b(x); CREATE TABLE b(x);` → error que contiene `002_rota.sql`; la tabla `b` no existe; la versión 2 no está registrada; la 1 sí.
  - `TestAplicarMigraciones_NombreInvalido`: un archivo `1_x.sql` → error.
  - `TestAplicarMigraciones_IgnoraVersionDesconocida`: la tabla de control ya registra la versión 9, que el FS no trae; con `{001, 002}` no hay error, se aplican solo `001` y `002` y la versión 9 sigue registrada.
  Correr `go test ./internal/repository -run TestAplicarMigraciones` y confirmar que FALLA (no compila: `aplicarMigraciones` no existe).
- [ ] T006 [US1] COMMIT `git add internal/repository/migrate_test.go` — `RED: agregar tests del runner de migraciones (falla)`
- [ ] T007 [US1] GREEN: implementar en internal/repository/migrate.go `aplicarMigraciones(ctx context.Context, db *sql.DB, fsys fs.FS) error` según research R2–R3: `CREATE TABLE IF NOT EXISTS migraciones_aplicadas (version INTEGER PRIMARY KEY, aplicada_en TEXT NOT NULL)`; `fs.ReadDir(fsys, ".")`; cada nombre debe cumplir `NNN_nombre.sql` (tres dígitos) o devuelve error; por cada versión no registrada: `BeginTx` → `ExecContext(contenido)` → `INSERT INTO migraciones_aplicadas (version, aplicada_en)` con `time.Now().UTC().Format(time.RFC3339)` → `Commit`; ante error `Rollback` y `fmt.Errorf("aplicar migración %s: %w", nombre, err)`. `Migrar` no se toca. `go test ./...` en verde.
- [ ] T008 [US1] COMMIT `git add internal/repository/migrate.go` — `GREEN: implementar el runner de migraciones numeradas, test en verde`

**Checkpoint**: runner probado con migraciones de juguete; el resto del sistema sigue con `Migrar`.

---

## Fase 2: Esquema 001 y test de esquema (US2)

**Objetivo**: la migración `001_esquema_base.sql` con §4.1–§4.5 y el test de esquema de §9 punto 6 acotado a la 001 (FR-003, FR-004, FR-015; CA-1, CA-5 a CA-12).

**Test independiente**: `go test ./internal/repository -run TestEsquema -v`.

- [ ] T009 [US2] RED: crear internal/repository/esquema_test.go con:
  - Helper `nuevaBDMigrada(t *testing.T) *sql.DB`: `AbrirSQLite(":memory:")` + `AplicarMigraciones(ctx, db)` + `t.Cleanup(db.Close)`. Lo reutilizan las fases 3 y 4.
  - `TestEsquema_Estructura` (CA-5): la estructura esperada **escrita en el test**, copiada de MODELO-DE-DATOS.md §4.1–§4.5 con un comentario que cite cada sección (FR-015; el test no lee el `.md`): para `proyectos`, `usuarios`, `integrantes`, `historias` y `sprints`, cada columna con tipo, `NOT NULL` y PK (`pragma_table_info`); cada FK, incluidas las compuestas `(sprint_id, proyecto_id) → sprints(id, proyecto_id)` (`pragma_foreign_key_list`); los índices `sprints_un_activo_por_proyecto` y `sprints_un_pendiente_por_proyecto`, únicos y parciales (`pragma_index_list`); que `on_delete` y `on_update` valgan `NO ACTION` en cada FK (FR-004, D9); y `migraciones_aplicadas` con la versión 1 registrada. Falla si sobra o falta algo.
  - `TestEsquema_Filas`: tabla de casos, cada uno con su `INSERT`/`UPDATE` y el resultado esperado — aceptado o el código de `(*sqlite.Error).Code()`: 787 FK, 275 `CHECK`, 1299 `NOT NULL`, 2067 `UNIQUE`. Casos de la spec, sin agregar otros: CA-6 `proyectos`; CA-7 `sprints`; CA-8 `historias`; CA-9 índices parciales (por `INSERT` y por `UPDATE`; varios `Finalizado` válidos); CA-10 fechas D12 (`'2026-10-9'`, `'basura'`, `'2026-02-30'` rechazadas; `'2026-10-09'` y `NULL` válidas) en `proyectos.fecha_inicio`/`fecha_fin`, `sprints.fecha_inicio`/`fecha_fin` e `historias.fecha_finalizacion`; CA-11 enums de `integrantes.rol`, `historias.prioridad`, `historias.estado`, `sprints.estado` y `NULL` en columnas `NOT NULL`; CA-12 cierre de Sprint (el `UPDATE historias SET sprint_id = NULL, estado = 'Nueva'` único se acepta; quitar primero el Sprint de una `En progreso` da 275); y D9: `DELETE` de un proyecto con historias da 787 (caso límite de la spec, FR-004).
  Correr `go test ./internal/repository -run TestEsquema` y confirmar que FALLA (no compila: `AplicarMigraciones` no existe).
- [ ] T010 [US2] COMMIT `git add internal/repository/esquema_test.go` — `RED: agregar el test de esquema de la migración 001 (falla)`
- [ ] T011 [US2] GREEN: crear internal/repository/migraciones/001_esquema_base.sql con el DDL de MODELO-DE-DATOS.md §4.1–§4.5 tal cual (tablas en español, FK dentro del `CREATE TABLE`, sin `ON DELETE`, `CHECK` de enums, de rangos con `typeof`, de fechas `date(col) IS col` y de reglas de forma, `UNIQUE (id, proyecto_id)` en `historias` y `sprints`, los dos índices parciales con sus nombres). Agregar en internal/repository/migrate.go `//go:embed migraciones/*.sql` (variable `migracionesFS embed.FS`) y `AplicarMigraciones(ctx context.Context, db *sql.DB) error`, que llama a `aplicarMigraciones` con el subdirectorio `migraciones` (`fs.Sub`). `go test ./...` en verde (los tests viejos siguen con `Migrar`).
- [ ] T012 [US2] COMMIT `git add internal/repository/migraciones/001_esquema_base.sql internal/repository/migrate.go` — `GREEN: agregar la migración 001 con el esquema base, test en verde`

**Checkpoint**: la 001 existe y el test de esquema está en verde; los repositorios todavía usan el esquema viejo.

---

## Fase 3a: Renombre — `usuarios` (US3)

**Objetivo**: `sqlite_user.go` sobre `usuarios` (`nombre`, `nombre_normalizado`, `creado_en`) (FR-007).

- [ ] T013 [US3] RED: en internal/repository/sqlite_user_test.go, los tres tests (`TestSQLiteUserRepository_CrearYBuscar`, `_NoEncuentra`, `_NoDuplicaNombreNormalizado`) pasan de `abrirBDDePrueba(t)` a `nuevaBDMigrada(t)`. Correr `go test ./internal/repository -run TestSQLiteUserRepository` y confirmar que FALLA (`no such table: users`).
- [ ] T014 [US3] COMMIT `git add internal/repository/sqlite_user_test.go` — `RED: tests de sqlite_user sobre el esquema 001 (falla)`
- [ ] T015 [US3] GREEN: en internal/repository/sqlite_user.go, el SQL de `Guardar` y `ObtenerPorNombreNormalizado` pasa a `usuarios (nombre, nombre_normalizado, creado_en)`. Los tests de `sqlite_user` quedan en verde. **Excepción de proceso (spec, *Proceso*; MODELO-DE-DATOS.md §10): desde este GREEN, el paquete `repository` y `features` quedan en rojo, porque los tests y los steps que todavía usan `Migrar` (esquema viejo) dependen del repositorio de usuarios ya renombrado. Es esperado; vuelven al verde al terminar la fase 3d. No se intenta un verde intermedio.** Verificar solo `go test ./internal/repository -run TestSQLiteUserRepository`.
- [ ] T016 [US3] COMMIT `git add internal/repository/sqlite_user.go` — `GREEN: sqlite_user usa la tabla usuarios, test en verde`

---

## Fase 3b: Renombre — `proyectos` + `integrantes` (US3)

**Objetivo**: `sqlite_project.go` sobre `proyectos` e `integrantes`, con fechas `YYYY-MM-DD` (FR-007, FR-008, D3; research R9).

- [ ] T017 [US3] RED: en internal/repository/sqlite_project_test.go: todos los tests que usan `abrirBDDePrueba(t)` pasan a `nuevaBDMigrada(t)` y se borra `abrirBDDePrueba`; el SQL directo de los tests pasa a las tablas y columnas nuevas (`SELECT COUNT(*) FROM proyectos`, `... WHERE nombre = ?`); los usuarios que necesiten los tests se crean con `SQLiteUserRepository`. Nuevo `TestSQLiteProjectRepository_FechasComoDia`: guarda un proyecto con inicio y fin, lee `fecha_inicio` y `fecha_fin` crudas con `SELECT` y espera `"2026-10-01"` y `"2026-12-01"` (FR-008). Confirmar que FALLA con `go test ./internal/repository -run TestSQLiteProjectRepository`.
- [ ] T018 [US3] COMMIT `git add internal/repository/sqlite_project_test.go` — `RED: tests de sqlite_project sobre el esquema 001 y fechas AAAA-MM-DD (falla)`
- [ ] T019 [US3] GREEN: en internal/repository/sqlite_project.go: `projects` → `proyectos (nombre, descripcion, fecha_inicio, fecha_fin, creado_en)` en `insertarProyecto`, `ObtenerPorID` y `Actualizar`; `project_members` → `integrantes (proyecto_id, usuario_id, rol, creado_en)` en `GuardarConScrumMaster`, `AgregarIntegrante` y `ListarIntegrantes` (el `JOIN` a `usuarios u` con `u.nombre`); `tiempoOpcionalAValor`/`valorATiempoOpcional` usan `time.DateOnly` en vez de `time.RFC3339` (`creado_en` sigue en RFC 3339). Tests de `sqlite_user` y `sqlite_project` en verde (el resto del paquete y `features` siguen en rojo por la excepción).
- [ ] T020 [US3] COMMIT `git add internal/repository/sqlite_project.go` — `GREEN: sqlite_project usa proyectos e integrantes con fechas AAAA-MM-DD, test en verde`

---

## Fase 3c: Renombre — `sprints` (US3) — un solo commit `test:`

**Objetivo**: los tests de `sqlite_sprint` corren sobre el esquema 001. `sprints` conserva nombre y columnas, así que **`sqlite_sprint.go` no cambia** y no hay un `RED:` honesto: fallaría el fixture, no el código (research R7). Va como un único commit `test:`.

- [ ] T021 [US3] En internal/repository/sqlite_sprint_test.go: `TestSQLiteSprintRepository_GuardarYObtenerPendiente`, `_ObtenerInexistente`, `_ActualizarPersisteLaTransicion` y `_ListarPorProyecto` pasan a `nuevaBDMigrada(t)` y crean antes el proyecto con `NewSQLiteProjectRepository(db).Guardar` (lo exige la FK); `guardarSprint` usa el id real de ese proyecto. `TestSQLiteSprintRepository_ListarPorProyecto` hoy guarda dos Sprints `Pendiente` en el proyecto 1 y uno en el 2: crear dos proyectos reales e iniciar el primer Sprint del proyecto 1 (`Iniciar` del dominio + `Actualizar`) antes de crear el segundo, para no chocar con `sprints_un_pendiente_por_proyecto` (el test que la viola a propósito es de la fase 4); las aserciones de orden no cambian. `TestMigrar_EsIdempotente`, `TestSQLiteBacklogRepository_HistoriasDeSprint` y el helper `nuevaBD` quedan como están (se ocupan las fases 3d y 5). `go test ./internal/repository -run TestSQLiteSprintRepository` en verde.
- [ ] T022 [US3] COMMIT `git add internal/repository/sqlite_sprint_test.go` — `test: tests de sqlite_sprint sobre el esquema 001 (Refs #107)`

---

## Fase 3d: Renombre — `historias` (US3) — vuelve el verde

**Objetivo**: `sqlite_backlog.go` sobre `historias`; los steps de Godog sobre el esquema 001 (FR-005, FR-007, SC-002; research R8).

- [ ] T023 [US3] RED: en internal/repository/sqlite_backlog_test.go: `nuevaBDDePrueba` usa `nuevaBDMigrada` y crea un proyecto real antes de guardar historias; el SQL directo pasa a `historias`. En internal/repository/sqlite_sprint_test.go, `TestSQLiteBacklogRepository_HistoriasDeSprint` pasa a `nuevaBDMigrada`, crea el proyecto y dos Sprints reales (uno `Activo`, con `Iniciar` del dominio + `Actualizar`, y otro `Pendiente`, para no violar `sprints_un_pendiente_por_proyecto`) en lugar de los ids fijos 7 y 8; el SQL directo pasa a `historias`, y la historia `Completada` se marca junto con `fecha_finalizacion = '2026-10-20'` (lo exige el `CHECK` de §4.4). Confirmar que FALLA con `go test ./internal/repository -run 'TestSQLiteBacklogRepository'` (`no such table: backlog_items`).
- [ ] T024 [US3] COMMIT `git add internal/repository/sqlite_backlog_test.go internal/repository/sqlite_sprint_test.go` — `RED: tests de sqlite_backlog sobre el esquema 001 (falla)`
- [ ] T025 [US3] GREEN: en internal/repository/sqlite_backlog.go, `backlog_items` → `historias` en `Guardar`, `ListarPorSprint` y `QuitarDeSprint`. En los steps de Godog (los `.feature` **no** se tocan, SC-002):
  - features/steps_creacion.go: `repository.Migrar` → `repository.AplicarMigraciones`; `backlog_items` → `historias` en las consultas.
  - features/steps_proyecto.go: `Migrar` → `AplicarMigraciones`; `projects` → `proyectos`, `name` → `nombre`, `users` → `usuarios` y `project_members` → `integrantes` (columnas `proyecto_id`, `usuario_id`, `rol`, §4.3) en las consultas (líneas 67, 105, 143, 192, 203, 287).
  - features/steps_sprints.go: `backlog_items` → `historias`; la historia `Completada` se actualiza junto con `fecha_finalizacion` (`CHECK` de §4.4).
  - features/steps_estado_proyecto.go: `fijarFechas` escribe `UPDATE proyectos SET fecha_inicio = ?, fecha_fin = ?` con `time.DateOnly` (el `CHECK` de D12 rechaza RFC 3339); renombrar `fechaOpcionalRFC3339` a `fechaOpcionalComoDia`.
  **Con este GREEN termina la excepción de proceso: el paquete `repository` y `features`, en rojo desde el GREEN de la fase 3a, vuelven al verde.** Correr `go test ./...` (incluye Godog) y confirmar todo en verde.
- [ ] T026 [US3] COMMIT `git add internal/repository/sqlite_backlog.go features/steps_creacion.go features/steps_proyecto.go features/steps_sprints.go features/steps_estado_proyecto.go` — `GREEN: sqlite_backlog usa historias y los steps de Godog el esquema 001, test en verde`

**Checkpoint**: los cuatro repositorios y Godog sobre el esquema 001, todo en verde. `cmd/api` todavía llama a `Migrar` (fase 5).

---

## Fase 4: Traducción de errores (US4)

**Objetivo**: FK → `ErrIntegridadReferencial`/409; índice parcial → `ValidationError` sobre `estado` con el mensaje del servicio/400; `CHECK`/`NOT NULL` → error genérico/500 (FR-010 a FR-014; CA-13; research R5–R6; [contracts/errores-http.md](contracts/errores-http.md)).

**Test independiente**: `go test ./internal/repository ./internal/http -run 'Integridad|IndiceParcial|Check|Conflicto' -v`.

- [ ] T027 [P] [US4] RED: en internal/repository/sqlite_backlog_test.go, `TestSQLiteBacklogRepository_Guardar_ProyectoInexistenteEsIntegridad`: `Guardar` de una historia con `proyecto_id` 999 → `errors.Is(err, ErrIntegridadReferencial)` (FR-010, contra SQLite real, FR-014).
- [ ] T028 [P] [US4] RED: en internal/repository/sqlite_sprint_test.go (todos contra SQLite real con `nuevaBDMigrada`):
  - `TestSQLiteSprintRepository_Guardar_ProyectoInexistenteEsIntegridad`: proyecto 999 → `ErrIntegridadReferencial` (FR-010).
  - `TestSQLiteSprintRepository_IndiceParcial_SegundoPendiente`: segundo `Guardar` de un `Pendiente` del mismo proyecto → `domain.ValidationError{Campo: "estado", Mensaje: "el proyecto ya tiene un Sprint pendiente"}` (FR-011).
  - `TestSQLiteSprintRepository_IndiceParcial_SegundoActivo`: con un Sprint ya `Activo`, `Actualizar` otro a `Activo` (goal y fechas válidas) → `ValidationError{Campo: "estado", Mensaje: "el proyecto ya tiene un Sprint activo"}` (FR-011).
  - `TestSQLiteSprintRepository_CheckEsErrorGenerico`: `Actualizar` un Sprint a `Activo` con `SprintGoal` vacío → hay error, y no es `ErrIntegridadReferencial` ni `domain.ValidationError` (FR-012, D13).
- [ ] T029 [P] [US4] RED: en internal/http/sprint_handler_test.go, `TestEscribirError_IntegridadReferencialEsConflicto`: con el fake de servicio del archivo devolviendo `repository.ErrIntegridadReferencial`, la respuesta es 409 con `{"mensaje": "violación de integridad referencial"}` (FR-010, contracts/errores-http.md).
- [ ] T030 [US4] Confirmar que T027–T029 FALLAN con `go test ./internal/repository ./internal/http`. COMMIT `git add internal/repository/sqlite_backlog_test.go internal/repository/sqlite_sprint_test.go internal/http/sprint_handler_test.go` — `RED: agregar tests de traducción de violaciones de la base (falla)`
- [ ] T031 [US4] GREEN:
  - internal/repository/errors.go: `errorDeEscritura(err error) error` que devuelve `ErrIntegridadReferencial` si `errors.As(err, &*sqlite.Error)` y `Code() == sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY`, y si no `err` sin cambios (research R5).
  - internal/repository/sqlite_backlog.go: `Guardar` pasa el error por `errorDeEscritura` antes de envolverlo.
  - internal/repository/sqlite_sprint.go: `Guardar` y `Actualizar` pasan el error por `errorDeEscritura`; si el código es `SQLITE_CONSTRAINT_UNIQUE` (2067), devuelven `domain.ValidationError{Campo: "estado", ...}` con `"el proyecto ya tiene un Sprint activo"` si `s.Estado == domain.SprintActivo` o `"el proyecto ya tiene un Sprint pendiente"` si es `Pendiente`.
  - internal/http/sprint_handler.go: `escribirError` gana `case errors.Is(err, repository.ErrIntegridadReferencial):` → `escribirJSON(w, nethttp.StatusConflict, errorResponse{Mensaje: err.Error()})`.
  `CHECK`/`NOT NULL` no necesitan código: caen en el `default` (500). `go test ./...` en verde.
- [ ] T032 [US4] COMMIT `git add internal/repository/errors.go internal/repository/sqlite_backlog.go internal/repository/sqlite_sprint.go internal/http/sprint_handler.go` — `GREEN: traducir violaciones de FK e índices parciales, test en verde`
- [ ] T033 [US4] REFACTOR (research R6, DRY):
  - internal/domain/sprint.go: `var ErrYaHaySprintActivo = ValidationError{Campo: "estado", Mensaje: "el proyecto ya tiene un Sprint activo"}` y `ErrYaHaySprintPendiente` (mismo campo, `"... pendiente"`).
  - internal/service/iniciar_sprint.go y internal/service/crear_sprint.go: usan esos valores en lugar del literal.
  - internal/repository/sqlite_sprint.go: usa esos valores en lugar del literal.
  - internal/repository/sqlite_project.go: `GuardarConScrumMaster` y `AgregarIntegrante` usan `errorDeEscritura` para la FK en lugar de la comparación a mano (se mantiene la traducción a `ErrMiembroDuplicado`, FR-013).
  `go test ./...` sigue en verde sin cambiar ningún test.
- [ ] T034 [US4] COMMIT `git add internal/domain/sprint.go internal/service/iniciar_sprint.go internal/service/crear_sprint.go internal/repository/sqlite_sprint.go internal/repository/sqlite_project.go` — `REFACTOR: mensajes de Sprint en el dominio y traducción de FK en un solo helper`

---

## Fase 5: Base vieja, cableado y limpieza (US5, US1)

**Objetivo**: detectar una base vieja (FR-006, FR-006a; CA-3; research R4), arrancar `cmd/api` con el runner (FR-005) y borrar el código viejo (FR-001).

- [ ] T035 [US5] RED: en internal/repository/migrate_test.go:
  - `TestAplicarMigraciones_BaseVieja`: tabla con cuatro casos, uno por tabla: una base con **solo** `projects`, solo `users`, solo `project_members` y solo `backlog_items` (cada una con un `CREATE TABLE` mínimo) → `AplicarMigraciones` devuelve `errors.Is(err, ErrEsquemaViejo)`, el mensaje contiene `seguimiento.db` y no se creó `migraciones_aplicadas` (FR-006a).
  - `TestVerificarBaseVieja_SoloSprints`: una base con solo una tabla `sprints` → `verificarBaseVieja(ctx, db)` devuelve `nil`. Se prueba **esa función directamente, no `AplicarMigraciones`**: con una `sprints` previa, `AplicarMigraciones` fallaría después, al crear `sprints` en la 001, porque la tabla ya existe.
  Confirmar que FALLA (`verificarBaseVieja` no existe).
- [ ] T036 [US5] COMMIT `git add internal/repository/migrate_test.go` — `RED: agregar tests de detección de base vieja (falla)`
- [ ] T037 [US5] GREEN:
  - internal/repository/migrate.go: `verificarBaseVieja(ctx, db) error` con `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name IN ('projects','users','project_members','backlog_items')` → `ErrEsquemaViejo` si es mayor que cero; `aplicarMigraciones` la llama **antes** de crear `migraciones_aplicadas`. Borrar `Migrar`, `asegurarColumna`, `verificarEsquemaViejo` y las constantes `esquemaBacklogItems`, `esquemaSprints`, `esquemaProjects`, `esquemaUsers`, `esquemaProjectMembers` (FR-001).
  - internal/repository/errors.go: el texto de `ErrEsquemaViejo` pasa a pedir borrar `seguimiento.db` y reiniciar, sin mencionar IDs TEXT.
  - internal/repository/migrate_test.go: borrar `TestMigrar_EsquemaViejoIDsTexto`; `TestAbrirSQLite_EscriturasConcurrentesSinLock` usa `AplicarMigraciones` e inserta en `proyectos (nombre, descripcion, creado_en)`.
  - internal/repository/sqlite_sprint_test.go: borrar `TestMigrar_EsIdempotente` (lo cubre `TestAplicarMigraciones_Idempotente`) y el helper `nuevaBD` si quedó sin uso.
  - cmd/api/main.go: `repository.Migrar` → `repository.AplicarMigraciones` (FR-005).
  `go build ./...` y `go test ./...` en verde.
- [ ] T038 [US5] COMMIT `git add internal/repository/migrate.go internal/repository/errors.go internal/repository/migrate_test.go internal/repository/sqlite_sprint_test.go cmd/api/main.go` — `GREEN: detectar base vieja, arrancar con el runner y quitar Migrar, test en verde`

---

## Fase 6: Cierre

- [ ] T039 Correr desde la raíz: `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test ./...` (incluye Godog). Todo en verde (CA-2). Si algo falla, se corrige en un commit propio antes de seguir. Estos tests existentes respaldan lo que no tiene tarea propia (deben seguir en verde sin cambios de aserción):
  - "Respuestas iguales" (US3.1, FR-009): los tests de internal/http/project_handler_test.go, internal/http/backlog_handler_test.go e internal/http/sprint_handler_test.go, y los cuatro `.feature` de Godog.
  - "El 404 sigue" (US4.4): `TestSprintHandler_Crear_ProyectoInexistente`, `TestSprintHandler_Iniciar_NoEncontrado`, `TestProjectHandler_Obtener_Inexistente`, `TestProjectHandler_Editar_Inexistente`, `TestProjectHandler_AsignarIntegrante_ProyectoInexistente`, `TestProjectHandler_ObtenerEstado_Inexistente`; en el servicio, `TestCrearHistoriaBacklog_Ejecutar_ProyectoInexistente`, `TestCrearSprint_ProyectoInexistente`, `TestIniciarSprint_NoEncontrado`, `TestCerrarSprint_NoEncontrado`, `TestAsignarIntegrante_ProyectoInexistente`; en Godog, los escenarios "Creación de un Sprint para un proyecto inexistente", "Rechazo de creación por proyecto inexistente" y "Consulta del estado de un proyecto inexistente".
  - Duplicados (FR-013): `TestSQLiteUserRepository_NoDuplicaNombreNormalizado`, `TestSQLiteProjectRepository_MiembroDuplicado`, `TestAsignarIntegrante_Duplicado` y `TestProjectHandler_AsignarIntegrante_Duplicado`.
- [ ] T040 Validar con [quickstart.md](quickstart.md), pasos 2 a 6: los `.feature` sin cambios contra `main` (SC-002), base vieja con `go run ./cmd/api` (CA-3), reinicio sin borrar la base (SC-001), `POST /projects` con fechas `AAAA-MM-DD` (US3) y el diff de MODELO-DE-DATOS.md (CA-14). Anotar el resultado de cada paso en la descripción del PR.
- [ ] T041 Marcar con `[X]` las tareas hechas en specs/006-esquema-migraciones/tasks.md. COMMIT `git add specs/006-esquema-migraciones/tasks.md` — `docs: marcar las tareas completadas de la spec 006 (Refs #107)`
- [ ] T042 Redactar la descripción del PR y entregarla a Carlos en el chat (él abre el PR desde `fix/107-esquema-migraciones` hacia `main`). Debe incluir:
  - `Closes #107` y `Closes #19`.
  - Resumen: runner de migraciones con `go:embed`, migración 001, renombre de los cuatro repositorios, fechas del proyecto a `AAAA-MM-DD`, traducción de errores (único cambio de la API: 409 ante una violación de FK) y detección de base vieja. Avisar que cada integrante tiene que borrar su `seguimiento.db` **una última vez**.
  - **Excepción de proceso**: `repository` y `features` quedan en rojo desde el `GREEN:` de la fase 3a hasta cerrar la 3d (los tests y steps que seguían con `Migrar` dependían de repositorios ya renombrados); en el resto de los commits el paquete está en verde.
  - **Excepción al Principio III (BDD)**, decisión del Scrum Master: el 409 por FK y el error de base vieja no llevan escenario de Godog. El 409 no se alcanza por el flujo normal (el servicio valida antes) y se prueba contra SQLite real en la fase 4; la base vieja es un error de arranque del proceso, no comportamiento de la API. Los `.feature` no se tocan (SC-002).
  - La fase 3c va como un solo commit `test:`, sin `RED:`: `sprints` conserva nombre y columnas, `sqlite_sprint.go` no cambia y un `RED:` habría fallado por el fixture, no por el código (research R7).
  - MODELO-DE-DATOS.md se corrigió en el mismo PR (§4.5, §5, §9 puntos 1 y 5, §10), según §11.
  - Resultado del quickstart (T040) y de `/speckit-analyze` (T001).
  - Pedido de revisión: aprobación de otra integrante **y** una revisión por agente en una **sesión nueva** (el cambio es grande).
  - Merge **solo** con "Create a merge commit", nunca squash ni rebase: se perdería la evidencia RED/GREEN/REFACTOR.
  - Terminar con la línea de atribución de Claude Code.

---

## Dependencies & Execution Order

- **Orden estricto**: 0b (análisis, spec y documento) → 1 → 2 → 3a → 3b → 3c → 3d → 4 → 5 → 6. Cada fase depende de la anterior: la 2 necesita el runner; la 3 el helper `nuevaBDMigrada` y la 001; cada renombre sigue las dependencias entre FK (`usuarios` → `proyectos`/`integrantes` → `sprints` → `historias`); la 4 necesita los cuatro repositorios sobre la 001; la 5 borra `Migrar`, que hasta ahí sostiene a los tests viejos y a `cmd/api`.
- Dentro de cada fase: RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT en la 4). Nunca saltear el COMMIT de RED.
- **Ventana en rojo** (excepción documentada): del GREEN de la 3a (T015) al GREEN de la 3d (T025). Fuera de esa ventana, cada commit deja `go test ./...` en verde.

## Parallel Opportunities

- T027, T028 y T029 (tests RED de la fase 4) tocan archivos distintos y se pueden escribir en paralelo; se commitean juntos en T030.
- Fuera de eso no hay paralelismo útil: un solo responsable (§10) y fases encadenadas.

## Implementation Strategy

1. **Fases 0b–2 (MVP técnico)**: documento corregido, runner y 001 con su test de esquema; el sistema sigue andando con `Migrar`.
2. **Fase 3**: renombre en orden de dependencias; se acepta el rojo de la 3a a la 3d y se recupera el verde al cerrar la 3d.
3. **Fases 4–5**: traducción de errores y limpieza, cada una dejando el sistema en verde.
4. **Fase 6**: chequeos finales, quickstart y descripción del PR (el análisis de consistencia ya corrió en la fase 0b).

## Notes

- El historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra: un commit por fase del ciclo.
- Fuera de alcance (no hay tareas): frontend, migraciones 002 a 005 (#117, HU-02, HU-07, HU-08), GET nuevos (#109, #110), HU-15 y el requisito Docker/Azure.
- `git add` siempre por nombre de archivo; nunca `git push` ni merge desde las tareas.
