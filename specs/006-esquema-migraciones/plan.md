# Implementation Plan: Adopción del modelo de datos — esquema SQLite con migraciones numeradas (Issue #107)

**Branch**: `fix/107-esquema-migraciones` | **Date**: 2026-10-23 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/006-esquema-migraciones/spec.md`

## Summary

Reemplazar `Migrar` (DDL sueltos + `asegurarColumna` + `verificarEsquemaViejo`) por un
runner de migraciones numeradas embebidas con `go:embed`, agregar la migración
`001_esquema_base.sql` con el esquema de [MODELO-DE-DATOS.md §4.1–§4.5](../../MODELO-DE-DATOS.md#4-tablas),
renombrar el SQL de los cuatro repositorios SQLite, pasar las fechas del proyecto a
`YYYY-MM-DD` y traducir las violaciones de la base según [§5](../../MODELO-DE-DATOS.md#5-dónde-vive-cada-regla)
(FK → `ErrIntegridadReferencial`/409; índice parcial → `ValidationError` con el mensaje del
servicio/400; `CHECK`/`NOT NULL` → error genérico/500). Una base con las tablas viejas
detiene el arranque pidiendo borrar `seguimiento.db`. Decisiones técnicas en [research.md](research.md).

## Technical Context

**Language/Version**: Go 1.26.5 (`go.mod`)

**Primary Dependencies**: `modernc.org/sqlite` v1.60.1 (ya instalado); biblioteca estándar (`embed`, `io/fs`, `testing/fstest`, `database/sql`). **Sin dependencias nuevas** (constitución VI, §9 punto 1).

**Storage**: SQLite embebido (`seguimiento.db`; `:memory:` en tests y Godog)

**Testing**: `testing` estándar (unitarios y de repositorio contra SQLite real), Godog (`features/`)

**Target Platform**: servidor HTTP local (Linux/Windows del equipo)

**Project Type**: web-service (API REST en Go; el frontend no se toca)

**Performance Goals**: N/A (base de desarrollo; la 001 se aplica una vez por base)

**Constraints**: una sola conexión (`SetMaxOpenConns(1)`), `PRAGMA foreign_keys = ON` por conexión (ya lo hace `AbrirSQLite`); cada migración en su propia transacción; el dominio no conoce SQLite (constitución IV)

**Scale/Scope**: 1 migración (001), 5 tablas + tabla de control, 4 repositorios, 1 traductor HTTP, 4 archivos de steps de Godog

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Cumple | Cómo |
|---|---|---|
| I. Test-First | ✅ | Cada fase con lógica de §10 arranca con `RED:` y sigue con `GREEN:` (tabla de fases abajo). La ventana en rojo, del `GREEN:` de la 3a al cierre de la 3d, es la excepción documentada en §10 y en la spec (*Proceso*); se repite en `tasks.md` y en el PR. |
| II. SDD | ✅ | `specify` → `clarify` → `plan` hechos; faltan `tasks` y `analyze` antes de escribir código. |
| III. BDD | ✅ | No hay comportamiento nuevo para el usuario: los `.feature` no cambian (SC-002). Solo cambia el código de los steps (SQL directo y llamada al runner). |
| IV. Hexagonal | ✅ | SQL, runner y traducción de códigos SQLite quedan en `internal/repository`. Los mensajes del índice parcial salen del dominio (ver research R6), no se copian en el adaptador. |
| V. KISS/YAGNI/DRY | ✅ | Runner de una función, sin "down migrations", sin locking entre procesos y sin una biblioteca de migraciones. El esquema esperado se escribe una vez en el test (FR-015); los mensajes de Sprint, una vez en el dominio. |
| VI. Stack | ✅ | `go:embed` de la biblioteca estándar y `modernc.org/sqlite` ya instalado. |

**Post-design re-check (Phase 1)**: sin violaciones. Complexity Tracking vacío.

## Fases de implementación (MODELO-DE-DATOS.md §10)

| Fase | Commit(s) | `RED:` (falla) | `GREEN:` |
|---|---|---|---|
| 0b. Documento | `docs:` (sin ciclo) | — | Corregir MODELO-DE-DATOS.md (FR-016a): §4.5 nombres de índices, §5 texto real de los mensajes, §9 punto 1 DDL de `migraciones_aplicadas`, §9 punto 5 y §10 las cuatro tablas viejas; §10 la frase de la ventana en rojo, que arranca en el `GREEN:` de la fase 3a (no en la 2). Va primero para que los tests copien del documento ya corregido. |
| 1. Runner | RED + GREEN | `migrate_test.go`: con `fstest.MapFS` aplica en orden, registra versión y `aplicada_en`, es idempotente, tolera huecos, aplica solo las nuevas, revierte y no registra una migración que falla a la mitad y conserva las anteriores; nombre de archivo inválido da error | Solo `aplicarMigraciones(ctx, db, fsys)` en `migrate.go` (research R1–R3), sin `go:embed`: un patrón `go:embed` sin archivos no compila y la 001 llega en la fase 2. `Migrar` sigue existiendo. |
| 2. Esquema 001 | RED + GREEN | `esquema_test.go`: helper de test `nuevaBDMigrada(t)` (usa `AplicarMigraciones`; lo reutilizan las fases 3 y 4), estructura esperada copiada de §4.1–§4.5 + casos CA-5 a CA-12; falla porque no existen la `001` ni `AplicarMigraciones` | `migraciones/001_esquema_base.sql`; en `migrate.go`, `//go:embed migraciones/*.sql` y `AplicarMigraciones(ctx, db)`, que llama a `aplicarMigraciones` con el FS embebido |
| 3a. `usuarios` | RED + GREEN | Los tests de `sqlite_user` pasan al helper `nuevaBDMigrada(t)` (creado en la fase 2) | SQL de `sqlite_user.go` a `usuarios` / `nombre`, `nombre_normalizado`, `creado_en` |
| 3b. `proyectos` + `integrantes` | RED + GREEN | Tests de `sqlite_project` al helper nuevo; uno nuevo lee `fecha_inicio`/`fecha_fin` crudas y espera `YYYY-MM-DD` | SQL de `sqlite_project.go`; fechas con `time.DateOnly` (D3) |
| 3c. `sprints` | un commit `test:` | — | `sprints` conserva nombre y columnas, así que `sqlite_sprint.go` **no cambia**. Solo se adaptan sus tests: helper nuevo y un proyecto real creado antes (por la FK). No hay `RED:` honesto (fallaría el fixture, no el código), así que va como un único commit; se avisa en el PR. Diferencia con §10 ("un par de commits por repositorio"), ver research R7. |
| 3d. `historias` | RED + GREEN | Tests de `sqlite_backlog` y `HistoriasDeSprint` al helper nuevo, con proyecto y Sprints reales y `fecha_finalizacion` en la `Completada` | SQL de `sqlite_backlog.go` a `historias`; steps de Godog: `AplicarMigraciones`, SQL directo con tablas/columnas nuevas, fechas `YYYY-MM-DD`, `fecha_finalizacion` al completar (research R8). **Aquí vuelve el verde de `repository` y `features`.** |
| 4. Traducción de errores | RED + GREEN (+ REFACTOR) | Tests contra SQLite real: historia y Sprint con proyecto inexistente → `ErrIntegridadReferencial`; segundo `Pendiente` por `Guardar` y segundo `Activo` por `Actualizar` → `ValidationError{estado}` con el mensaje del servicio; `Activo` sin goal por `Actualizar` → error que no es ni FK ni `ValidationError`; `escribirError` con `ErrIntegridadReferencial` → 409 | `errorDeEscritura` en `errors.go` (research R5); uso en `sqlite_backlog` y `sqlite_sprint`; caso 409 en `escribirError`. REFACTOR: mensajes de Sprint a `domain` (R6) y `sqlite_project.go` pasa a usar `errorDeEscritura` en lugar de su traducción de FK a mano. |
| 5. Cableado y limpieza | RED + GREEN | `migrate_test.go`: una base con **solo** `projects`, solo `users`, solo `project_members` y solo `backlog_items` (4 casos) → `AplicarMigraciones` devuelve `ErrEsquemaViejo`; base con solo `sprints` → `verificarBaseVieja` no da error (se prueba esa función directamente, no `AplicarMigraciones`, que fallaría después al crear `sprints` en la 001 porque la tabla ya existe) | `verificarBaseVieja` en el runner; `cmd/api/main.go` usa `AplicarMigraciones`; se borran `Migrar`, `asegurarColumna`, `verificarEsquemaViejo`, las constantes `esquema*` y sus tests viejos; `TestAbrirSQLite_EscriturasConcurrentesSinLock` pasa a `proyectos` |
| 6. Cierre | — | — | `go build ./...`, `go vet ./...`, `gofmt -l .`, `go test ./...` (incluye Godog) en verde; `/speckit.analyze` registrado; PR con `Closes #107` y la excepción de proceso en la descripción |

**Excepción de proceso (spec, *Proceso*; §10).** `repository` y `features` quedan en rojo
desde el `GREEN:` de la fase 3a hasta cerrar la 3d. La fase 2 sigue en verde porque `Migrar`
vive hasta la fase 5 y los tests existentes lo siguen usando. Desde la 3a, cada
repositorio que pasa al esquema nuevo rompe a los tests todavía en el viejo que lo usan
(p. ej. los de proyectos crean usuarios). No se intenta un verde por repositorio. El
binario (`cmd/api`) tampoco anda entre la 3a y la 5; nadie lo corre a mitad del PR.

## Project Structure

### Documentation (this feature)

```text
specs/006-esquema-migraciones/
├── spec.md
├── plan.md              # este archivo
├── research.md          # decisiones técnicas R1–R8
├── data-model.md        # enlaces a MODELO-DE-DATOS.md (no redefine)
├── quickstart.md        # cómo validar de punta a punta
├── contracts/
│   ├── migraciones.md   # contrato del runner y de los archivos NNN_nombre.sql
│   └── errores-http.md  # traducción de violaciones a HTTP (único cambio: 409)
├── checklists/requirements.md
└── tasks.md             # /speckit.tasks (todavía no)
```

### Source Code (repository root)

```text
MODELO-DE-DATOS.md                         # fase 0b (FR-016a)
cmd/api/main.go                            # fase 5: AplicarMigraciones
internal/domain/sprint.go                  # fase 4 REFACTOR: mensajes "ya tiene un Sprint ..."
internal/service/crear_sprint.go           # fase 4 REFACTOR: usa el mensaje del dominio
internal/service/iniciar_sprint.go         # ídem
internal/http/sprint_handler.go            # fase 4: escribirError → 409
internal/http/sprint_handler_test.go       # fase 4
internal/repository/
├── migraciones/001_esquema_base.sql       # fase 2 (nuevo)
├── migrate.go                             # fase 1 runner; fase 5 base vieja y borrado de Migrar
├── migrate_test.go                        # fases 1 y 5
├── esquema_test.go                        # fase 2 (nuevo)
├── errors.go                              # fase 4 errorDeEscritura; fase 5 texto de ErrEsquemaViejo
├── sqlite_user.go / _test.go              # 3a
├── sqlite_project.go / _test.go           # 3b, 4
├── sqlite_sprint.go / _test.go            # 3c, 4
└── sqlite_backlog.go / _test.go           # 3d, 4
features/steps_creacion.go                 # 3d
features/steps_proyecto.go                 # 3d
features/steps_sprints.go                  # 3d
features/steps_estado_proyecto.go          # 3d
```

**Structure Decision**: se mantiene la estructura hexagonal existente. Todo lo nuevo vive en
`internal/repository` (runner, SQL, traducción de códigos); `domain` solo gana los dos
mensajes de Sprint como valores reutilizables; `http` agrega un caso al traductor. Los
fakes en memoria (`*_en_memoria.go`, `sprint_repository_memoria.go`) no cambian.

## Complexity Tracking

Sin violaciones de la constitución que justificar.
