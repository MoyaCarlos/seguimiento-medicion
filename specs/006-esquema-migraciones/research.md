# Research: Adopción del modelo de datos (Issue #107)

Las decisiones de modelo (tablas, columnas, restricciones) **no** se toman acá: están en
[MODELO-DE-DATOS.md](../../MODELO-DE-DATOS.md) y en las *Clarifications* de la [spec](spec.md).
Esto registra solo las decisiones técnicas de implementación.

## Verificación previa (2026-10-23)

Probado con `modernc.org/sqlite` v1.60.1 en un módulo aparte, fuera del repo:

- `db.Exec` con varias sentencias separadas por `;` las ejecuta todas → un archivo `.sql` se aplica con un solo `ExecContext`.
- Un `tx.Exec("CREATE TABLE c(z); CREATE TABLE a(x);")` que falla en la segunda sentencia, seguido de `Rollback`, no deja `c` → el DDL de SQLite es transaccional y una migración que falla a la mitad no deja rastros.
- Códigos `(*sqlite.Error).Code()`: FK = 787, `CHECK` = 275, `NOT NULL` = 1299 (coinciden con §5).

## R1. Forma del runner

- **Decision**: `AplicarMigraciones(ctx, db) error` (exportada, usa el `embed.FS` de `migraciones/*.sql`) delega en `aplicarMigraciones(ctx, db, fsys fs.FS) error`, que es la que prueban los tests con `fstest.MapFS`. Vive en `migrate.go`, junto a `AbrirSQLite`.
- **Rationale**: inyectar el `fs.FS` permite probar orden, huecos y fallos con migraciones de juguete sin tocar la 001. Una función, sin tipo `Runner` ni interfaz (KISS).
- **Alternatives considered**: `golang-migrate` (dependencia nueva, prohibida por §9 y la constitución VI); un struct `Migrador` configurable (YAGNI).

## R2. Orden, nombres y huecos

- **Decision**: se listan los archivos con `fs.ReadDir` (que devuelve orden por nombre) y cada nombre debe cumplir `NNN_nombre.sql` con `NNN` de tres dígitos; la versión es `NNN` como entero. Un nombre que no cumple da error al arrancar. Se aplica cada versión que no esté en `migraciones_aplicadas`; los huecos no importan porque solo se compara pertenencia.
- **Rationale**: con tres dígitos fijos, el orden por nombre es el orden numérico; no hace falta ordenar a mano. Rechazar nombres mal formados evita que un `1_algo.sql` se cuele fuera de orden.
- **Alternatives considered**: ordenar por entero con `sort.Slice` (código extra que el prefijo fijo ya resuelve); exigir numeración sin huecos (§9 pide tolerarlos: dos ramas reservan números).
- **Versiones desconocidas** (supuesto de la spec): una versión registrada que el binario no trae se ignora.

## R3. Transacción por migración y tabla de control

- **Decision**: antes de nada, `CREATE TABLE IF NOT EXISTS migraciones_aplicadas (version INTEGER PRIMARY KEY, aplicada_en TEXT NOT NULL)` (DDL decidido en *Clarifications*). Por cada migración pendiente: `BeginTx` → `ExecContext(contenido)` → `INSERT INTO migraciones_aplicadas (version, aplicada_en)` con `time.Now().UTC().Format(time.RFC3339)` → `Commit`. Ante error, `Rollback` y `fmt.Errorf("aplicar migración %s: %w", nombre, err)`.
- **Rationale**: el registro y el DDL quedan en la misma transacción: o están los dos o ninguno (FR-002, escenario 5 de US1). Las migraciones anteriores ya confirmadas se conservan.
- **Alternatives considered**: una sola transacción para todas (una falla en la 003 desharía la 002 ya válida); registrar fuera de la transacción (podría quedar registrada una migración revertida).
- **Nota**: §9 punto 4 (recrear tablas con `foreign_keys = OFF` fuera de la transacción) no aplica a la 001, que solo crea. Lo resolverá la primera migración que recree una tabla.

## R4. Detección de base vieja

- **Decision**: antes de crear la tabla de control, `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name IN ('projects','users','project_members','backlog_items')`; si da > 0 → `ErrEsquemaViejo`, con el texto nuevo: pide borrar `seguimiento.db` y reiniciar. `sprints` no está en la lista (*Clarifications*).
- **Rationale**: una sola consulta, sin inspeccionar columnas. Ir antes de crear `migraciones_aplicadas` asegura que una base vieja no quede "a medio migrar".
- **Alternatives considered**: mirar el tipo de `projects.id`, como hoy `verificarEsquemaViejo` (solo detectaba un caso y no ve una base que solo tiene `backlog_items`).

## R5. Traducción de errores en los adaptadores

- **Decision**: un helper en `errors.go`, `errorDeEscritura(err error) error`, que con `errors.As(err, &*sqlite.Error)` devuelve `ErrIntegridadReferencial` si `Code() == sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY` (787) y si no devuelve `err` tal cual. Lo usan `Guardar` de historias y Sprints (en el `GREEN:` de la fase 4) y, en el `REFACTOR:` de la fase 4, las dos escrituras de `integrantes` en `sqlite_project.go`, que ya hacían esa traducción a mano. El índice parcial (2067, `SQLITE_CONSTRAINT_UNIQUE`) lo traduce **solo** `sqlite_sprint.go` (`Guardar` y `Actualizar`), según `s.Estado`. Las unicidades con error propio (`ErrUsuarioDuplicado`, `ErrMiembroDuplicado`) siguen como están (FR-013).
- **`CHECK`/`NOT NULL`**: no necesitan código: el error queda envuelto con `fmt.Errorf` y `escribirError` ya responde 500 en el caso `default` (D13). El test de la fase 4 verifica que **no** se traduce a `ErrIntegridadReferencial` ni a `ValidationError`.
- **HTTP**: `escribirError` (`internal/http/sprint_handler.go`) gana `case errors.Is(err, repository.ErrIntegridadReferencial): 409` con `errorResponse{Mensaje: err.Error()}`. `http` ya importa `repository`.
- **Rationale**: un solo lugar para la regla "FK → integridad" (DRY); por código, nunca por texto (§5).
- **Alternatives considered**: un tipo de error genérico nuevo para `CHECK` (el `default` ya hace lo que pide D13; YAGNI); traducir 2067 en un helper común (solo `sprints` tiene índices parciales y el mensaje depende del estado del Sprint).

## R6. Mensajes del índice parcial sin duplicar texto

- **Decision**: en la fase 4, el `GREEN:` arma el `ValidationError` en el adaptador; el `REFACTOR:` mueve los dos mensajes a `internal/domain/sprint.go` como valores (`ErrYaHaySprintActivo`, `ErrYaHaySprintPendiente`, ambos `ValidationError{Campo: "estado", ...}` con el texto actual del servicio) y `crear_sprint.go`, `iniciar_sprint.go` y `sqlite_sprint.go` los usan.
- **Rationale**: el texto queda en un solo lugar (DRY) y así la validación previa y la carrera responden igual por construcción (*Clarifications*). El dominio ya es dueño de `ValidationError` y el adaptador ya depende del dominio.
- **Alternatives considered**: copiar el literal en el adaptador (dos fuentes del mismo texto).

## R7. Fase 3c sin cambios de código

- **Decision**: `sprints` conserva el nombre y las columnas, así que `sqlite_sprint.go` no cambia en la fase 3. La 3c se reduce a adaptar sus tests (helper nuevo y un proyecto real por la FK) en un único commit `test:`.
- **Rationale**: un `RED:` ahí fallaría por el fixture, no por el código: no sería evidencia de TDD.
- **Consulta**: §10 pide "un par de commits por repositorio". Es una diferencia de **proceso**, no de modelo; se avisa al Scrum Master y en el PR.

## R8. Ajustes en los tests existentes (no en los `.feature`)

Los tests y los steps que escriben SQL a mano o dependen de que no hubiera FK cambian:

- `sqlite_sprint_test.go` (`HistoriasDeSprint`): hoy usa los Sprints 7 y 8 y el proyecto 1 sin crearlos; con la FK compuesta hay que crearlos y poner `fecha_finalizacion` en la historia `Completada` (`CHECK` de §4.4).
- `sqlite_backlog_test.go`: crear el proyecto antes de guardar historias.
- `features/steps_sprints.go`: el SQL directo pasa a `historias` y la historia `Completada` lleva `fecha_finalizacion`.
- `features/steps_estado_proyecto.go`: `fijarFechas` escribe `proyectos.fecha_inicio`/`fecha_fin` en `YYYY-MM-DD` (con RFC 3339 lo rechaza el `CHECK` de D12).
- `features/steps_creacion.go`, `steps_proyecto.go`: `Migrar` → `AplicarMigraciones`; SQL directo a `historias`.
- Ninguno de los escenarios fija fechas iguales en proyecto o Sprint, así que el `CHECK fin > inicio` no altera ningún resultado esperado.

## R9. Fechas del proyecto (D3)

- **Decision**: `tiempoOpcionalAValor`/`valorATiempoOpcional` de `sqlite_project.go` pasan a `time.DateOnly`. `creado_en` sigue en RFC 3339.
- **Rationale**: el DTO ya parsea y formatea `AAAA-MM-DD` (`formatoFecha`), así que el valor `time.Time` es el mismo (medianoche UTC) y la API no cambia (FR-008).
