# Feature Specification: Adopción del modelo de datos — esquema SQLite con migraciones numeradas (Issue #107, continúa #19)

**Feature Branch**: `fix/107-esquema-migraciones`

**Created**: 2026-10-23

**Status**: Draft

**Input**: User description: "Adopción del modelo de datos: esquema SQLite con migraciones numeradas (Issue #107, continúa #19). Reemplazar los CREATE TABLE sueltos y los ALTER por un runner de migraciones numeradas (go:embed, sin dependencias nuevas) y una migración 001 con el esquema base del documento (proyectos, usuarios, integrantes, historias, sprints), con llaves foráneas dentro del CREATE TABLE, CHECK para enums y fechas (D12) e índices únicos parciales para Sprint Activo y Pendiente. Renombrar tablas y columnas en los 4 repositorios SQLite. Fechas del proyecto a YYYY-MM-DD. Traducción de errores en los adaptadores (§5). Una base con las tablas viejas en inglés da un error claro que pide borrar seguimiento.db."

> **Fuente única de diseño:** [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md). Esta spec
> **no** define tablas, columnas ni restricciones: dice *qué* tiene que cumplir la adopción
> y enlaza a la sección del documento que lo describe. Si esta spec y el documento
> divergen, manda el documento y se corrige la spec en el mismo PR (§11 del documento).
> Ante una duda sobre el modelo se frena y se pregunta al Scrum Master (AGENTS.md).

## Contexto

Es un **fix**, no una historia: entra por la convención de fixes de `AGENTS.md` (issue
`bug` #107, rama `fix/107-esquema-migraciones`, PR con `Closes #107`). Es el arreglo de
raíz de #19: sus partes 2 y 3 (FK en el esquema, base vieja, traducción del error) se
absorben acá ([§10 Encuadre de proceso](../../MODELO-DE-DATOS.md#10-cambios-sobre-lo-que-ya-está-en-main)).
El problema que resuelve está en [§1](../../MODELO-DE-DATOS.md#1-por-qué-existe).

Las decisiones D1–D13 del documento ([§7](../../MODELO-DE-DATOS.md#7-decisiones-pendientes-e-inconsistencias-del-backlog))
están cerradas. Las que bloqueaban la 001 son D12 (`CHECK` de formato de fechas) y D13
(500 ante `CHECK`/`NOT NULL`); D3 fija las fechas del proyecto en `YYYY-MM-DD`.

## Clarifications

### Session 2026-10-23

- Q: ¿Qué tablas delatan una base vieja? → A: `projects`, `users`, `project_members` y `backlog_items` (las que crea el `Migrar` actual y que la 001 no usa); `sprints` no sirve porque conserva el nombre. El test de detección cubre cada una de las cuatro por separado (una base con solo `backlog_items` también da el error), y el PR corrige §9 punto 5 y §10 de MODELO-DE-DATOS.md para nombrar las cuatro (decisión del Scrum Master).
- Q: ¿Qué texto devuelve el índice parcial de `sprints`? → A: el del servicio (`"el proyecto ya tiene un Sprint activo"` / `"el proyecto ya tiene un Sprint pendiente"`), para que la validación previa y la carrera respondan igual y no cambie la API. El PR corrige §5 de MODELO-DE-DATOS.md para que diga ese texto real, en el mismo PR (§11) (decisión del Scrum Master).
- Q: ¿Cómo detecta el test de esquema que el documento y la migración divergen (CA-1)? → A: el test lleva adentro la estructura esperada copiada de §4 (citando la sección) y la compara con la base migrada; no parsea MODELO-DE-DATOS.md. La coherencia entre el test y el documento se revisa en el PR (CA-4, §11).
- Q: ¿Quién crea la tabla de control `migraciones_aplicadas` y con qué DDL (§9 punto 1 no lo dice)? → A: la crea el runner, antes de leerla, con `CREATE TABLE IF NOT EXISTS migraciones_aplicadas (version INTEGER PRIMARY KEY, aplicada_en TEXT NOT NULL)`; `aplicada_en` en RFC 3339 UTC (§2, instantes). El PR agrega ese DDL a §9 punto 1 (decisión del Scrum Master, durante `/speckit-plan`).
- Q: ¿El 409 por FK y el error de base vieja necesitan escenarios de Godog (Principio III), aunque SC-002 congela los `.feature`? → A: no. El 409 no se alcanza por el flujo normal (el servicio valida antes de guardar) y se prueba contra SQLite real en la fase 4; la base vieja es un error de arranque del proceso, no comportamiento de la API. SC-002 se mantiene: los `.feature` no se tocan. Se registra como excepción en *Proceso* y en el PR (decisión del Scrum Master, durante `/speckit-analyze`).
- Q: ¿Qué nombre llevan los índices únicos parciales de `sprints` (§4.5 los escribe con `...`)? → A: `sprints_un_activo_por_proyecto` y `sprints_un_pendiente_por_proyecto`. El PR los escribe en §4.5 y el test de esquema los verifica (decisión del Scrum Master, durante `/speckit-plan`).

## User Scenarios & Testing *(mandatory)*

Los "usuarios" de este fix son el equipo (quien corre la aplicación o los tests en su
máquina) y el consumidor de la API. El comportamiento observable de la API **no cambia**,
salvo el 409 de la User Story 4.

### User Story 1 - El esquema se crea con migraciones numeradas (Priority: P1)

Como integrante del equipo, quiero que la base se cree y se actualice sola con migraciones
numeradas al arrancar la aplicación, para no tener que borrar `seguimiento.db` cada vez que
otra historia cambia el esquema.

**Why this priority**: Sin el mecanismo de migraciones no hay dónde poner la 001 ni las
002–005 que vienen (#117, HU-02, HU-07, HU-08). Es la base de todo lo demás.

**Independent Test**: Con migraciones de prueba (no la 001) sobre una base vacía: se
aplican en orden, quedan registradas, una segunda ejecución no hace nada, tolera huecos de
numeración y una migración que falla a la mitad no deja rastros.

**Acceptance Scenarios**:

1. **Given** una base vacía y un conjunto de migraciones numeradas, **When** arranca la aplicación, **Then** se aplican todas en orden ascendente de número y cada una queda registrada en la tabla de control ([§9 punto 1](../../MODELO-DE-DATOS.md#9-migraciones)).
2. **Given** una base en la que ya se aplicaron todas las migraciones, **When** arranca de nuevo, **Then** no se aplica ninguna y el esquema no cambia (idempotencia).
3. **Given** migraciones numeradas con huecos (p. ej. 001 y 003, sin 002), **When** arranca, **Then** se aplican las existentes en orden, sin error por el hueco.
4. **Given** una base con la 001 aplicada y una 003 nueva, **When** arranca, **Then** se aplica solo la 003.
5. **Given** una migración cuya segunda sentencia falla, **When** arranca, **Then** el arranque falla con un error que identifica la migración, ninguna sentencia de esa migración queda aplicada y la migración no queda registrada.

---

### User Story 2 - La base garantiza la forma de los datos (Priority: P1)

Como equipo, quiero que la migración 001 cree el esquema base del documento (tablas en
español, FK dentro del `CREATE TABLE`, `CHECK` de enums, de fechas y de reglas de forma, e
índices únicos parciales de Sprint), para que un dato inválido que se escape del dominio no
llegue a guardarse (cierra #19).

**Why this priority**: Es el objetivo del issue. El test de esquema es además el que
detecta cualquier divergencia entre el documento y la migración.

**Independent Test**: El **test de esquema** ([§9 punto 6](../../MODELO-DE-DATOS.md#9-migraciones)),
acotado a las tablas de la 001: aplica las migraciones sobre una base vacía, compara
estructura contra el documento e inserta filas válidas e inválidas, incluidas variantes con
`NULL`.

**Acceptance Scenarios**:

1. **Given** una base vacía, **When** se aplican las migraciones, **Then** existen las tablas `proyectos`, `usuarios`, `integrantes`, `historias` y `sprints` con exactamente las columnas, obligatoriedad y FK de [§4.1–§4.5](../../MODELO-DE-DATOS.md#4-tablas), verificado con `pragma_table_info` y `pragma_foreign_key_list`.
2. **Given** el esquema 001, **When** se insertan las filas de los casos mínimos de §9 punto 6 que corresponden a la 001 (ver *Casos del test de esquema* abajo), **Then** cada fila se acepta o se rechaza como indica el documento, y cada rechazo trae el código de violación esperado (FK, `CHECK`, `NOT NULL` o `UNIQUE`).
3. **Given** que alguien cambia una columna, un `CHECK` o una FK en la migración sin actualizar la estructura esperada del test (copiada de §4), **When** corre el test de esquema, **Then** el test falla. Un cambio del documento sin cambiar el test se detecta en la revisión del PR (CA-4).

---

### User Story 3 - Los repositorios usan el esquema nuevo sin cambiar el comportamiento (Priority: P1)

Como consumidor de la API, quiero que todo lo que hoy funciona siga funcionando igual
después del renombre, para que la adopción no rompa HU-01, HU-04, HU-05 ni HU-13.

**Why this priority**: Sin el renombre, la aplicación no anda sobre el esquema 001.

**Independent Test**: Los tests de los cuatro repositorios SQLite corren contra las tablas y
columnas nuevas, y los escenarios de Godog existentes pasan sin modificar los `.feature`.

**Acceptance Scenarios**:

1. **Given** el esquema 001, **When** se crean, leen y actualizan usuarios, proyectos, integrantes, Sprints e historias por la API, **Then** las respuestas (estado HTTP y cuerpo) son las mismas que antes de la adopción.
2. **Given** un proyecto creado con fechas de inicio y fin, **When** se lee de la base, **Then** sus fechas están guardadas como `YYYY-MM-DD` (D3) y la API las devuelve igual que antes.
3. **Given** los archivos `.feature` de `features/` sin cambios, **When** corre Godog, **Then** todos los escenarios pasan.

---

### User Story 4 - Las violaciones de la base se traducen a errores claros (Priority: P2)

Como consumidor de la API, quiero que, si la base rechaza un dato que el servicio dejó
pasar, la respuesta distinga un conflicto de integridad (409), una regla de Sprint (400) y
un defecto del servidor (500), en lugar de un 500 genérico para todo.

**Why this priority**: Los 404 precisos ya los da el servicio al validar antes de guardar;
esto es la red de seguridad para la carrera y para defectos del dominio.

**Independent Test**: Tests de los adaptadores contra SQLite real (no fakes) que provocan
cada tipo de violación y verifican el error de dominio; y un test del traductor HTTP.

**Acceptance Scenarios**:

1. **Given** una violación de FK (simple o compuesta) al guardar, **When** el adaptador recibe el error, **Then** devuelve `ErrIntegridadReferencial` y la API responde 409.
2. **Given** un segundo Sprint `Activo` (o `Pendiente`) para el mismo proyecto que pasa la verificación del servicio (carrera), **When** se guarda, **Then** el adaptador devuelve un `ValidationError` sobre `estado` con el mensaje del estado que se estaba guardando, y la API responde 400.
3. **Given** una violación de `CHECK` o `NOT NULL`, **When** se guarda, **Then** el adaptador devuelve un error genérico y la API responde 500 (D13).
4. **Given** un recurso inexistente referenciado en un pedido, **When** el servicio valida antes de guardar, **Then** la API sigue respondiendo el 404 preciso de hoy (proyecto o Sprint no encontrado), no el 409.

---

### User Story 5 - Una base vieja se detecta con un mensaje claro (Priority: P2)

Como integrante del equipo con una `seguimiento.db` creada antes de la adopción, quiero un
error que me diga qué hacer, en lugar de fallas raras en el primer pedido.

**Why this priority**: Es la última vez que se pide borrar la base ([§9 punto 5](../../MODELO-DE-DATOS.md#9-migraciones));
afecta a cada máquina del equipo una sola vez.

**Independent Test**: Abrir una base que tiene las tablas viejas y verificar que el
arranque falla con el error que pide borrar `seguimiento.db`.

**Acceptance Scenarios**:

1. **Given** una base que tiene una sola de las tablas viejas (`projects`, `users`, `project_members` o `backlog_items`; un caso por tabla), **When** arranca la aplicación, **Then** el arranque falla con un error que pide borrar `seguimiento.db` y reiniciar, y no se aplica ninguna migración.
2. **Given** una base vacía o una base ya migrada, **When** arranca, **Then** no se reporta base vieja.

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional). Las
restricciones concretas de cada tabla **no** se repiten acá: viven en
[§4](../../MODELO-DE-DATOS.md#4-tablas) y se prueban con los casos de §9 punto 6.

**Runner de migraciones**

- El sistema **shall** aplicar, al arrancar, cada migración embebida que no figure en la tabla de control, en orden ascendente de número.
- **WHEN** todas las migraciones embebidas ya figuran en la tabla de control, **THEN** el sistema **shall** arrancar sin ejecutar ninguna sentencia de esquema.
- **WHEN** la numeración de las migraciones tiene huecos, **THEN** el sistema **shall** aplicar las existentes en orden sin reportar error.
- **IF** una sentencia de una migración falla, **THEN** el sistema **shall** revertir esa migración completa, no registrarla y detener el arranque con un error que nombre la migración.
- **IF** una migración falla, **THEN** el sistema **shall** conservar aplicadas y registradas las migraciones anteriores que ya habían terminado.
- **WHERE** la base es en memoria (tests), el sistema **shall** aplicar las migraciones igual que sobre un archivo.
- El sistema **shall** obtener las migraciones del propio binario, sin leer archivos del disco en tiempo de ejecución y sin dependencias nuevas.

**Base vieja**

- **IF** la base contiene al menos una de las tablas `projects`, `users`, `project_members` o `backlog_items`, **THEN** el sistema **shall** detener el arranque, antes de aplicar migraciones, con un error que pide borrar `seguimiento.db` y reiniciar.
- **WHILE** la base está vacía o solo tiene tablas del esquema nuevo (incluida `sprints`), el sistema **shall** no reportar base vieja.

**Esquema 001 (forma)**

- El sistema **shall** crear las tablas, columnas, FK, `CHECK` e índices únicos parciales de [§4.1–§4.5](../../MODELO-DE-DATOS.md#4-tablas), con nombres en español y valores de enum tal cual el dominio ([§2](../../MODELO-DE-DATOS.md#2-convenciones-valen-para-todas-las-tablas)).
- **IF** una fila tiene una columna que admite `NULL` en `NULL`, **THEN** cada `CHECK` que la usa **shall** evaluarse como lo indica la auditoría de [§4.10](../../MODELO-DE-DATOS.md#410-auditoría-de-los-check-frente-a-null) (un `CHECK` que da `NULL` no debe dejar pasar una fila inválida).
- **IF** una fecha de calendario no tiene el formato `YYYY-MM-DD` o no es una fecha real (`'2026-10-9'`, `'basura'`, `'2026-02-30'`), **THEN** la base **shall** rechazarla (D12).
- **WHEN** se guarda un Sprint con estado `Finalizado`, **THEN** la base **shall** permitir varios por proyecto.
- **IF** se intenta borrar o cambiar el id de un padre con hijos, **THEN** la base **shall** rechazarlo (sin cascadas, D9).

**Repositorios y fechas**

- El sistema **shall** guardar `fecha_inicio` y `fecha_fin` del proyecto como `YYYY-MM-DD` y los instantes de auditoría (`creado_en`) como RFC 3339 en UTC ([§2](../../MODELO-DE-DATOS.md#2-convenciones-valen-para-todas-las-tablas)).
- **WHEN** se lee un proyecto con fechas, **THEN** la API **shall** devolverlas con el mismo formato que antes de la adopción.
- **WHEN** se guarda una historia con `sprint_id` en `NULL`, **THEN** la base **shall** aceptarla (la FK compuesta no se evalúa).

**Traducción de errores** ([§5](../../MODELO-DE-DATOS.md#5-dónde-vive-cada-regla))

- **IF** la base rechaza una escritura por cualquier FK (simple o compuesta), **THEN** el adaptador **shall** devolver `ErrIntegridadReferencial`, sin intentar adivinar qué recurso faltaba, y la API **shall** responder 409.
- **IF** la base rechaza un Sprint por un índice único parcial, **THEN** el adaptador **shall** devolver un `ValidationError` sobre `estado` con el mismo mensaje que da el servicio para el estado del Sprint que se estaba guardando (`Activo` → "el proyecto ya tiene un Sprint activo", `Pendiente` → "el proyecto ya tiene un Sprint pendiente"), y la API **shall** responder 400.
- **IF** la base rechaza una escritura por `CHECK` o `NOT NULL`, **THEN** el adaptador **shall** devolver un error genérico y la API **shall** responder 500 (D13).
- **IF** la base rechaza una escritura por la unicidad de `usuarios.nombre_normalizado` o la PK de `integrantes`, **THEN** el adaptador **shall** seguir devolviendo `ErrUsuarioDuplicado` / `ErrMiembroDuplicado`, como hoy.
- El adaptador **shall** distinguir la violación por el código del error del motor, nunca por el texto del mensaje.
- **WHEN** el pedido referencia un proyecto o Sprint inexistente, **THEN** el servicio **shall** seguir respondiendo el 404 preciso de hoy, validando antes de guardar.

## Requirements *(mandatory)*

### Functional Requirements

**Migraciones** ([§9](../../MODELO-DE-DATOS.md#9-migraciones))

- **FR-001**: El sistema MUST reemplazar la creación actual del esquema (`Migrar`, `CREATE TABLE IF NOT EXISTS`, `asegurarColumna`, `verificarEsquemaViejo` y las constantes `esquema*`) por un runner de migraciones numeradas en archivos `NNN_nombre.sql`, embebidos en el binario con la biblioteca estándar, sin agregar dependencias (constitución VI).
- **FR-002**: El runner MUST crear la tabla de control `migraciones_aplicadas (version INTEGER PRIMARY KEY, aplicada_en TEXT NOT NULL)` si no existe, registrar en ella cada migración aplicada (`aplicada_en` en RFC 3339 UTC), aplicar en orden las que falten (tolerando huecos) y ejecutar cada una en su propia transacción.
- **FR-003**: El sistema MUST incluir la migración `001_esquema_base.sql` con las tablas `proyectos`, `usuarios`, `integrantes`, `historias` y `sprints` tal como las definen [§4.1–§4.5](../../MODELO-DE-DATOS.md#4-tablas), incluidas las columnas nuevas que ya prevé el documento para historias siguientes (`criterios_aceptacion`, `fecha_finalizacion`, `horas_por_story_point`, [§8](../../MODELO-DE-DATOS.md#8-mapa-historia--tablas)).
- **FR-004**: Las FK MUST declararse dentro del `CREATE TABLE`, sin `ON DELETE CASCADE` ni `RESTRICT` (§2, D9).
- **FR-005**: El arranque de la aplicación (`cmd/api`) y el arranque de los escenarios de Godog MUST usar el runner en lugar de `Migrar`.

**Base vieja**

- **FR-006**: Antes de aplicar migraciones, el sistema MUST detectar una base vieja y detener el arranque con un error claro que pide borrar `seguimiento.db` y reiniciar (§9 punto 5). Una base es vieja si tiene **al menos una** de estas tablas: `projects`, `users`, `project_members`, `backlog_items`. `sprints` no cuenta (conserva el nombre).
- **FR-006a**: El test de detección MUST cubrir cada una de las cuatro tablas por separado: una base con solo esa tabla da el error (incluida una base con solo `backlog_items`).

**Renombre** ([§10 Fases](../../MODELO-DE-DATOS.md#10-cambios-sobre-lo-que-ya-está-en-main))

- **FR-007**: Los cuatro repositorios SQLite MUST usar las tablas y columnas de la 001, en este orden de trabajo: usuarios → proyectos e integrantes → Sprints → historias.
- **FR-008**: El repositorio de proyectos MUST guardar las fechas de inicio y fin como `YYYY-MM-DD` (D3), sin cambiar lo que expone la API.
- **FR-009**: El comportamiento observable de la API MUST quedar igual, con la única excepción del 409 de FR-010.

**Traducción de errores** ([§5](../../MODELO-DE-DATOS.md#5-dónde-vive-cada-regla))

- **FR-010**: Los adaptadores SQLite MUST traducir cualquier violación de FK a `ErrIntegridadReferencial`, y la API MUST responderla con 409. Reemplaza lo acordado antes para #19 (FK → `ErrProyectoNoEncontrado`/404), según §10.
- **FR-011**: El adaptador de Sprints MUST traducir la violación de un índice único parcial a un `ValidationError` sobre `estado` con el mismo mensaje que ya devuelve el servicio para ese estado: `Activo` → `"el proyecto ya tiene un Sprint activo"`, `Pendiente` → `"el proyecto ya tiene un Sprint pendiente"`. La API lo responde con 400, igual que la validación previa.
- **FR-012**: Los adaptadores MUST traducir una violación de `CHECK` o `NOT NULL` a un error genérico, que la API responde con 500 (D13).
- **FR-013**: Las violaciones de unicidad que ya tienen error propio (`ErrUsuarioDuplicado`, `ErrMiembroDuplicado`) MUST conservar su traducción actual.
- **FR-014**: La traducción MUST basarse en el código del error del motor, nunca en el texto, y MUST probarse contra SQLite real.

**Pruebas y documentación**

- **FR-015**: El sistema MUST incluir un test de esquema según §9 punto 6, acotado a las tablas de la 001. La estructura esperada (tablas, columnas, obligatoriedad, FK) MUST estar escrita en el propio test, copiada de §4.1–§4.5 y citando la sección; el test MUST NOT leer ni parsear MODELO-DE-DATOS.md. Si la migración se aparta de esa estructura, el test falla; que el test siga coincidiendo con el documento se verifica en la revisión del PR (CA-4).
- **FR-016**: El `data-model.md` de esta spec (que genera `/speckit.plan`) MUST enlazar a §4, §5, §9 y §10 del documento en lugar de redefinir tablas, columnas o restricciones.
- **FR-016a**: El mismo PR MUST corregir MODELO-DE-DATOS.md (§11) para que coincida con esta spec: §9 punto 5 y §10 nombran las cuatro tablas que delatan una base vieja; §5 cita el texto real de los mensajes del índice parcial; §9 punto 1 incluye el DDL de `migraciones_aplicadas` y que la crea el runner; §4.5 nombra los dos índices parciales; la frase de §10 sobre la ventana en rojo dice que arranca en el `GREEN:` de la fase 3a (no en la fase 2). Verificable en el diff del PR.
- **FR-017**: Si durante la implementación aparece una contradicción entre el documento y el código (que no sea la diferencia esperada de nombres viejos en `main`), MUST frenarse y consultarse al Scrum Master; no se elige uno en silencio.

### Proceso (excepción documentada)

- Cada fase con lógica de [§10](../../MODELO-DE-DATOS.md#10-cambios-sobre-lo-que-ya-está-en-main) lleva su commit `RED:` y su `GREEN:` por separado; `REFACTOR:` solo si hay algo que limpiar. Fases: 1 runner, 2 esquema 001, 3a–3d renombre (un par de commits `RED:`/`GREEN:` por repositorio, salvo la 3c: `sprints` conserva nombre y columnas, `sqlite_sprint.go` no cambia y la fase va como un único commit `test:` que adapta sus tests, research R7), 4 traducción de errores, 5 cableado y limpieza, 6 cierre.
- **Excepción de proceso:** desde el `GREEN:` de la fase 3a hasta terminar la 3d, el paquete `repository` —y lo que depende de él, como `features`— **queda en rojo**. La fase 2 sigue en verde porque `Migrar` (esquema viejo) vive hasta la fase 5 y los tests existentes lo siguen usando. El rojo empieza cuando el primer repositorio pasa al esquema 001: los tests que todavía usan el esquema viejo y dependen de ese repositorio (p. ej. los de proyectos, que crean usuarios) dejan de coincidir, y no se puede garantizar un verde intermedio por repositorio. El verde se recupera al cerrar la fase 3d. Esta excepción MUST figurar también en `tasks.md` y en la descripción del PR (como en HU-13 y #103).
- **Excepción al Principio III (BDD):** el 409 por violación de FK y el error de base vieja no llevan escenario de Godog. El 409 no se alcanza por el flujo normal (el servicio valida antes de guardar) y se prueba contra SQLite real en la fase 4; la base vieja es un error de arranque del proceso, no comportamiento de la API. SC-002 se mantiene: los `.feature` no se tocan. Decisión del Scrum Master; MUST figurar también en la descripción del PR.
- Es un solo PR, hecho de principio a fin por un solo responsable (Carlos, asignado en #107), mergeado con "Create a merge commit" y aprobado por otra integrante.

### Key Entities *(include if feature involves data)*

Ninguna entidad de dominio nueva. Las tablas están definidas en el documento y no se
repiten acá:

- **Esquema base (001)**: `proyectos`, `usuarios`, `integrantes`, `historias`, `sprints` → [§4.1–§4.5](../../MODELO-DE-DATOS.md#4-tablas).
- **Tabla de control de migraciones**: registro de las migraciones aplicadas → [§9 punto 1](../../MODELO-DE-DATOS.md#9-migraciones).
- **Dónde vive cada regla y traducción de errores** → [§5](../../MODELO-DE-DATOS.md#5-dónde-vive-cada-regla).

## Success Criteria *(mandatory)*

### Criterios de aceptación (issue #107)

- **CA-1**: El test de esquema pasa y falla si la migración se aparta de la estructura de §4 copiada en el test (FR-015); la divergencia test↔documento se controla en el PR (CA-4).
- **CA-2**: `go build`, `go vet`, `gofmt -l` y `go test ./...` en verde, y los escenarios de Godog sin cambios de comportamiento.
- **CA-3**: Una base con las tablas viejas da un error claro que pide borrar `seguimiento.db` (probado con cada una de las cuatro tablas por separado, FR-006a).
- **CA-4**: El documento de modelo de datos y las migraciones coinciden: un cambio de esquema modifica los dos en el mismo PR.

### Casos del test de esquema (§9 punto 6, acotado a la 001)

Los casos son los del documento; se listan solo para que el alcance sea verificable.

- **CA-5 Estructura**: cada tabla, columna y FK de §4.1–§4.5 existe (`pragma_table_info`, `pragma_foreign_key_list`).
- **CA-6 `proyectos`**: fin = inicio y fin < inicio rechazados; solo inicio, solo fin y ambas `NULL` válidas; `horas_por_story_point` 0, −1 y `'abc'` rechazados, `NULL` y `2.5` válidos.
- **CA-7 `sprints`**: `Activo` y `Finalizado` con `fecha_fin NULL`, con `fecha_inicio NULL` y con `sprint_goal NULL` (los seis rechazados); `Pendiente` con goal o con alguna fecha rechazado; fin = inicio y fin < inicio rechazados.
- **CA-8 `historias`**: `En progreso` y `Completada` sin Sprint rechazadas; `Completada` sin `fecha_finalizacion` rechazada; `Nueva` y `En progreso` con `fecha_finalizacion` rechazadas; Sprint de otro proyecto rechazado; Sprint `NULL` válido.
- **CA-9 Índices parciales**: existen `sprints_un_activo_por_proyecto` y `sprints_un_pendiente_por_proyecto`; segundo `Activo` y segundo `Pendiente` del mismo proyecto rechazados, por `INSERT` y por `UPDATE`; varios `Finalizado` válidos.
- **CA-10 Fechas (D12)** en `proyectos`, `sprints` e `historias.fecha_finalizacion`: `'2026-10-9'`, `'basura'` y `'2026-02-30'` rechazadas; `'2026-10-09'` válida; `NULL` válido donde la columna lo admite.
- **CA-11 Enums y obligatoriedad**: valores fuera de los enums de `integrantes.rol`, `historias.prioridad`, `historias.estado` y `sprints.estado` rechazados; `NULL` en columnas `NOT NULL` rechazado.
- **CA-12 Cierre de Sprint**: el `UPDATE` único (quitar Sprint y volver a `Nueva`) se acepta; quitar primero el Sprint de una historia `En progreso` se rechaza.
- **CA-13 Traducción**: cada tipo de violación se traduce al error de dominio de §5, probado en los adaptadores contra `modernc.org/sqlite`: FK → `ErrIntegridadReferencial`; índice parcial → `ValidationError` sobre `estado` con el mensaje del servicio; `CHECK`/`NOT NULL` → error genérico. La FK compuesta de una historia en un Sprint de otro proyecto no se puede provocar desde ningún adaptador actual (ninguno asigna `sprint_id`; es de HU-06): se cubre en el test de esquema (CA-8).
- **CA-14 Documento**: el diff del PR incluye la corrección de §4.5, §5, §9 puntos 1 y 5 y §10 de MODELO-DE-DATOS.md, incluida en §10 la frase de la ventana en rojo, que pasa a arrancar en el `GREEN:` de la fase 3a (FR-016a).

### Measurable Outcomes

- **SC-001**: Después de la adopción, ningún cambio de esquema posterior (002 en adelante) obliga a borrar `seguimiento.db`: una base con la 001 aplicada recibe la siguiente migración al arrancar.
- **SC-002**: El 100 % de los escenarios de Godog existentes pasa sin editar los archivos `.feature`.
- **SC-003**: Una fila que viola una regla de forma de §4.1–§4.5 se rechaza en el 100 % de los casos de CA-6 a CA-12, incluidas las variantes con `NULL`.
- **SC-004**: Ninguna violación de FK responde 500 ni 404 desde la base: responde 409.

## Out of scope

- Migraciones 002 a 005: `historias_arrastradas` y `CerrarConArrastre` (#117), rondas y votos (HU-02), registros de horas (HU-07), defectos (HU-08). Por eso los casos de §9 punto 6 de `registros_horas`, `defectos` e `historias_arrastradas` no entran.
- Los GET nuevos que pide el frontend (#109, #110).
- HU-15 (aunque la 001 ya trae `criterios_aceptacion`, su uso es de HU-15) y aceptar `horas_por_story_point` en `PUT /projects/{id}` (es de HU-03, §8).
- Un único Scrum Master por proyecto (D4), sacar o mover historias entre Sprints (D5).
- Migrar datos de bases viejas: son datos de desarrollo, se borran (§9 punto 5).

## Assumptions

- **Godog**: "sin cambios de comportamiento" significa que no se tocan los `.feature`; el código de los steps sí cambia lo mínimo para llamar al runner en vez de `Migrar` (FR-005).
- **Versiones desconocidas**: si la tabla de control registra una versión que el binario no trae (base migrada por una rama más nueva), el runner la ignora y aplica solo las que faltan. Es una base de desarrollo; no hay producción.
- `AbrirSQLite` ya activa `PRAGMA foreign_keys = ON` por conexión y usa una sola conexión; no cambia.
- #100 ya está mergeado (PR #111), así que el dominio rechaza fecha de fin igual al inicio antes de llegar al `CHECK` de la 001 (§10 Orden y dependencias).
