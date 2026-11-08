# Research: Un único Scrum Master por proyecto (#236)

Evidencia obtenida sobre el código real de la rama (base `origin/main` bca73c6) y con experimentos contra `modernc.org/sqlite` (SQLite 3.53.x, Go 1.27.1), ejecutados con tests temporales que se descartaron.

## 1. Cómo distinguir la violación del índice de la de la clave primaria

**Hoy**: `SQLiteProjectRepository.AgregarIntegrante` traduce `SQLITE_CONSTRAINT_UNIQUE` **y** `SQLITE_CONSTRAINT_PRIMARYKEY` a `ErrMiembroDuplicado` (`internal/repository/sqlite_project.go:151`). El índice nuevo caería en esa rama y respondería "integrante duplicado" (400) en lugar del 409.

**Experimento** (esquema migrado, `integrantes(proyecto_id, usuario_id)` con PK compuesta + índice `UNIQUE (proyecto_id) WHERE rol = 'scrum_master'`, el proyecto 1 ya tiene su Scrum Master):

| Operación | Código extendido | Texto del error |
|---|---|---|
| `INSERT` del mismo usuario y proyecto (PK) | **1555** (`PRIMARYKEY`) | `UNIQUE constraint failed: integrantes.proyecto_id, integrantes.usuario_id` |
| `INSERT` de otro usuario como `scrum_master` | **2067** (`UNIQUE`) | `UNIQUE constraint failed: integrantes.proyecto_id` |
| `UPDATE` de un `product_builder` a `scrum_master` | **2067** | `UNIQUE constraint failed: integrantes.proyecto_id` |
| `INSERT` del mismo usuario que ya es Scrum Master, otra vez como `scrum_master` (viola ambas) | **2067** | `UNIQUE constraint failed: integrantes.proyecto_id` |
| `INSERT` de un `scrum_master` en otro proyecto | sin error | |

**Decisión**: distinguir por **código extendido** (`*sqlite.Error.Code()`), no por el texto. 1555 → duplicado; 2067 → `ErrScrumMasterYaExiste`. El texto cambia según las columnas y no es contrato.

**Hallazgo**: cuando una fila viola las dos a la vez, SQLite informa 2067 (el índice), no 1555. Por eso el orden de FR-006 (el duplicado gana) **no se puede delegar a la base**: lo garantiza el chequeo previo del service. Solo en una carrera extrema (la misma persona agregada dos veces a la vez) el repositorio devolvería el 409 en vez del 400; es aceptable porque el estado final sigue siendo correcto y está declarado.

**Alternativas descartadas**: comparar el texto del error (frágil); consultar el índice después de fallar para reconstruir la causa (más código y más viajes a la base); `INSERT ... ON CONFLICT` (no distingue cuál restricción saltó).

## 2. Dónde se hace cumplir la regla

**Decisión** (recomendación del Scrum Master, 2026-11-08): doble protección. (a) El service comprueba antes de insertar para dar el orden de errores correcto y un mensaje claro sin depender de la base. (b) El índice garantiza la regla ante carreras y escrituras por fuera del dominio, y el repositorio traduce su violación al mismo error de dominio.

**Alternativas**: solo service (no cubre carreras ni escrituras directas); solo índice (no puede garantizar que el 400 gane sobre el 409, ver §1).

Los repositorios en memoria usados por los tests de service imitan el índice: `agregarIntegrante` devuelve `ErrScrumMasterYaExiste` si ya hay uno. Así el service y el repositorio se prueban con el mismo contrato.

## 3. La migración 007 y su criterio

Orden dentro de la misma migración (el runner la ejecuta en una transacción, así que o se aplica todo o nada):

1. `UPDATE` que pasa a `product_builder` a todo Scrum Master que no sea el de `creado_en` más antiguo de su proyecto (empate: menor `usuario_id`).
2. `CREATE UNIQUE INDEX integrantes_un_scrum_master ON integrantes (proyecto_id) WHERE rol = 'scrum_master'`.

**Experimento** (base migrada hasta la 005, sembrada, con el SQL de la 007 aplicado):

| Proyecto sembrado | Resultado |
|---|---|
| P1: tres Scrum Master (usuarios 3, 1, 2; el 3 es el más antiguo) | queda el 3; el 1 y el 2 pasan a `product_builder` |
| P2: dos Scrum Master con el mismo `creado_en` (usuarios 4 y 2) y un builder | queda el 2 (menor `usuario_id`); el 4 pasa a builder |
| P3: sin Scrum Master | sin cambios |
| P4: un solo Scrum Master | sin cambios |

El índice se creó sin error y ninguna fila se borró. El `UPDATE` es seguro aunque su subconsulta lea la misma tabla: quitar el rol a los más nuevos no cambia cuál es el primero.

**Ordenar por `creado_en` como texto** es correcto porque el modelo guarda instantes en RFC 3339 UTC con 'Z' (`MODELO-DE-DATOS.md` §2) y todo el código los escribe así (`time.Now().UTC().Format(time.RFC3339)`). Una fila con otro formato podría ordenar mal; no hay datos así.

**El comentario del criterio** va al principio del archivo de la 007: qué corrige, el desempate y por qué (la base de Render se vacía en cada deploy; los datos de prueba v2 crean un Scrum Master por proyecto).

**Numeración**: 007. `origin/feature/HU-20-autenticacion` ya tiene `006_autenticacion.sql`; `main` llega hasta la 005. El runner tolera huecos, así que la 007 funciona aunque la 006 todavía no esté en `main`. Si #236 se mergea antes, la tabla de §9 queda con un hueco en la 006 (reservada a HU-20); si se mergea después, no hay que renumerar nada. No hace falta renumerar mientras no choque el número.

## 4. Cómo se testea la migración

El runner interno `aplicarMigraciones(ctx, db, fsys fs.FS)` recibe el sistema de archivos como parámetro (`migrate.go`), y `migrate_test.go` ya lo usa con `fstest.MapFS`. Por lo tanto **no hace falta rediseñarlo**:

1. Armar un `fstest.MapFS` con las migraciones embebidas **menos la 007** (leídas de `migracionesFS`).
2. Aplicarlo sobre una base en memoria, sembrar proyectos con 0, 1, 2 y 3 Scrum Master y un empate de `creado_en`.
3. Aplicar el conjunto completo (`AplicarMigraciones`): el runner salta las ya aplicadas y corre solo la 007.
4. Comprobar quién quedó como Scrum Master en cada proyecto, que no se borró ningún integrante y que el índice rechaza un segundo Scrum Master.

Es el experimento de §3 convertido en test. El test de esquema (§9, punto 6) suma el índice a `esquemaEsperado` (campo `indices`, con su `WHERE`) y las filas válidas e inválidas del FR-014.

Godog arranca sobre el esquema ya migrado, por eso la corrección no va ahí; los 4 escenarios Godog usan la base migrada.

## 5. Quién crea Scrum Master hoy y qué rompería el índice

| Lugar | ¿Crea más de un Scrum Master por proyecto? | Acción |
|---|---|---|
| `POST /projects` → `GuardarConScrumMaster` (`sqlite_project.go`) | No: uno por proyecto nuevo | Sin cambios |
| `internal/demo/bloque_proyecto.go` (`/demo/datos`): Beta Uno SM, Beta Dos y Tres builders | No | Sin cambios; los datos v2 (#235) deben mantener un solo SM por proyecto (el test existente `bloque_proyecto_test.go` ya lo verifica) |
| `features/steps_votacion.go`, `steps_metricas.go`, `steps_historial_sprints.go` | No: un SM + builders | Sin cambios |
| `internal/repository/esquema_test.go` (`datosBase`), `internal/http/votacion_handler_test.go` (`equipoDeVotacion`) | No | Sin cambios |
| `internal/service/asignar_integrante_test.go` `TestAsignarIntegrante_VariosScrumMaster` | **Sí: afirma que "se permiten varios Scrum Master"** | **Se rompe, y debe romperse**: es el RED. Se invierte: el segundo Scrum Master se rechaza con `ErrScrumMasterYaExiste` |
| `asignar_integrante_test.go` `TestAsignarIntegrante_ReutilizaUsuario`: agrega a Jimena como `scrum_master` a un proyecto cuyo creador (Ana) ya es Scrum Master | **Sí** | Pasar a Jimena a `product_builder` en el proyecto "Otro" (el test verifica la reutilización del usuario, no el rol) |
| `internal/service/revelar_votacion_test.go:106` (cambia el rol en una lista en memoria) | No toca el repositorio | Sin cambios |
| Tests de `project_repository_en_memoria` y `sqlite_project` que agregan un `scrum_master` | Proyectos sin Scrum Master previo (usan `Guardar`) | Sin cambios; `..._GuardarConScrumMaster_RollbackSiFallaMembresia` usa un integrante pre-sembrado con otro usuario y falla por duplicado: sin cambios |

Ninguna fuente de datos de producción ni de demo se rompe con el índice.

## 6. Frontend

El formulario de equipo (`web/src/features/proyecto/Configuracion.tsx`) captura el error de `asignarIntegrante`; si el `campo` es `nombre` o `rol` lo muestra en el campo y, en cualquier otro caso (incluido sin `campo`), lo muestra como error general con `mensajeDeError`. El 409 trae solo `mensaje`, así que **se muestra solo, sin tocar el frontend**. El selector de rol sigue ofreciendo `scrum_master`; ocultarlo cuando ya hay uno sería una mejora de interfaz fuera de esta historia.

## 7. Decisiones que quedan abiertas

Ninguna del modelo de datos: D4 y la política de la migración están decididas. Pendiente solo de confirmar en el PR: el texto exacto del mensaje del 409 (propuesto: "el proyecto ya tiene un Scrum Master").
