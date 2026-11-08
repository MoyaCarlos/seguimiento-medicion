# Feature Specification: Un único Scrum Master por proyecto (#236)

**Feature Branch**: `feature/236-scrum-master-unico`

**Created**: 2026-11-08

**Status**: Draft

**Input**: User description: "Issue #236 — Un único Scrum Master por proyecto (índice único parcial, cambio de esquema). Decisión del Scrum Master (2026-11-03): con login y permisos por rol (HU-20 y HU-21), quien tenga el rol scrum_master ejecuta las acciones exclusivas, así que tiene que haber uno solo por proyecto. Quien crea el proyecto queda como Scrum Master y es su dueño, por lo que la decisión D14 de MODELO-DE-DATOS.md deja de hacer falta."

## Clarifications

### Session 2026-11-08

- Q: ¿Qué hace la migración 007 con proyectos que ya tienen más de un Scrum Master? → A: Opción B. Conserva como Scrum Master al de `creado_en` más antiguo de cada proyecto (si empatan, el de menor `usuario_id`) y pasa a `product_builder` a los demás, dentro de la misma migración y antes de crear el índice único parcial. Motivo: la base de Render se vacía en cada deploy y los datos de prueba v2 crean un solo Scrum Master por proyecto, así que la migración tiene que arrancar sola, sin depender de un arreglo manual. La migración 007 explica este criterio en un comentario.
- Q: Si la persona que se agrega como `scrum_master` ya es integrante del proyecto, ¿qué error prevalece? → A: El de integrante duplicado actual (400) prevalece sobre el 409. Confirma el supuesto de FR-006.
- Sin cambios: la migración se llama 007 y se renumera si #236 se mergea antes que HU-20; traspasar el rol y quitar integrantes siguen fuera de alcance (HU-16, #175).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - El proyecto nace con un único Scrum Master (Priority: P1)

Como persona que crea un proyecto, quedo registrada como su Scrum Master y dueña, para que desde el primer momento haya exactamente una persona con las acciones exclusivas del rol.

**Why this priority**: Es el estado base de la regla. Sin él, "uno solo por proyecto" no se cumple desde el origen. Hoy el comportamiento ya existe (el creador queda como Scrum Master); esta historia lo declara como regla y lo protege.

**Independent Test**: Crear un proyecto con un creador y listar sus integrantes: aparece uno solo, con rol `scrum_master`, que es el creador.

**Acceptance Scenarios**:

1. **Given** una persona que crea un proyecto, **When** el proyecto se crea, **Then** el proyecto tiene exactamente un Scrum Master y es esa persona.
2. **Given** dos proyectos distintos, **When** cada uno tiene su Scrum Master, **Then** la misma persona puede ser Scrum Master de ambos (la unicidad es por proyecto, no global).

---

### User Story 2 - Rechazar un segundo Scrum Master (Priority: P1)

Como Scrum Master, al configurar el equipo no puedo agregar a otra persona como Scrum Master mientras el proyecto ya tiene uno, y recibo un mensaje claro, para que las acciones exclusivas del rol nunca queden repartidas entre varias personas.

**Why this priority**: Es el núcleo de la regla de negocio y lo que el login y los permisos por rol (HU-20, HU-21) necesitan como garantía.

**Independent Test**: Con un proyecto que ya tiene Scrum Master, intentar agregar un integrante con rol `scrum_master`: se rechaza con conflicto (409) y mensaje claro, y la lista de integrantes no cambia.

**Acceptance Scenarios**:

1. **Given** un proyecto con Scrum Master, **When** agrego a otra persona con rol `scrum_master`, **Then** el sistema responde 409 con un mensaje que indica que el proyecto ya tiene Scrum Master, y no agrega a nadie.
2. **Given** un proyecto con Scrum Master, **When** agrego a una persona nueva con rol `product_builder`, **Then** queda agregada (201) y el Scrum Master no cambia.
3. **Given** un proyecto con Scrum Master, **When** agrego como `scrum_master` a quien ya es integrante del proyecto, **Then** se aplica la regla de integrante duplicado que ya existe y no se cambia su rol.
4. **Given** un proyecto con Scrum Master, **When** la misma persona intenta agregarse otra vez como `scrum_master`, **Then** el resultado es el de integrante duplicado, no un segundo Scrum Master.

---

### User Story 3 - La base impide un segundo Scrum Master por fuera del dominio (Priority: P1)

Como equipo, queremos que la propia base de datos rechace un segundo Scrum Master en un proyecto, para que la regla se cumpla aunque alguien escriba en la base sin pasar por el dominio, incluso con dos pedidos simultáneos.

**Why this priority**: El cambio de esquema es la parte central del issue. La verificación en el servicio sola no cubre escrituras concurrentes ni escrituras directas.

**Independent Test**: Con el esquema aplicado, insertar directamente dos filas `scrum_master` para el mismo proyecto: la segunda falla. Insertar dos `product_builder` o un `scrum_master` en otro proyecto: se aceptan.

**Acceptance Scenarios**:

1. **Given** un proyecto con un Scrum Master en la base, **When** se inserta directamente un segundo `scrum_master` para ese proyecto, **Then** la base lo rechaza.
2. **Given** un proyecto con un Scrum Master, **When** se actualiza a `scrum_master` el rol de un `product_builder` del mismo proyecto, **Then** la base lo rechaza.
3. **Given** un proyecto, **When** se insertan varios `product_builder`, **Then** todos se aceptan.
4. **Given** dos proyectos, **When** cada uno tiene su `scrum_master`, **Then** ambos se aceptan.
5. **Given** dos pedidos simultáneos que agregan un `scrum_master` al mismo proyecto, **When** ambos se ejecutan, **Then** exactamente uno se acepta y el otro recibe el conflicto de la historia 2.

---

### User Story 4 - La migración corrige proyectos con varios Scrum Master (Priority: P1)

Como equipo, queremos que la migración que crea el índice arranque sola aunque la base tenga proyectos con más de un Scrum Master, para que un deploy nunca quede bloqueado por datos viejos.

**Why this priority**: Sin este tratamiento, el índice no se puede crear sobre una base con duplicados y la aplicación no arranca.

**Independent Test**: Con una base en el esquema anterior a la 007, sembrar proyectos con 0, 1, 2 y 3 Scrum Master (y un empate de `creado_en`), aplicar la 007 y comprobar quién queda como Scrum Master en cada uno.

**Acceptance Scenarios**:

1. **Given** un proyecto con dos Scrum Master de distinto `creado_en`, **When** corre la 007, **Then** queda como Scrum Master el de `creado_en` más antiguo y el otro pasa a `product_builder`, sin dejar de ser integrante.
2. **Given** un proyecto con dos Scrum Master con el mismo `creado_en`, **When** corre la 007, **Then** queda el de menor `usuario_id`.
3. **Given** un proyecto con tres Scrum Master, **When** corre la 007, **Then** queda uno solo y los otros dos pasan a `product_builder`.
4. **Given** un proyecto con un solo Scrum Master, **When** corre la 007, **Then** no se toca.
5. **Given** un proyecto sin Scrum Master, **When** corre la 007, **Then** no se toca y la migración no le asigna uno.
6. **Given** una base ya corregida por la 007, **When** se intenta insertar un segundo Scrum Master, **Then** el índice lo rechaza.

---

### User Story 5 - Documentar la decisión en el modelo de datos (Priority: P2)

Como equipo, queremos que `MODELO-DE-DATOS.md` refleje la regla, para que el documento y el esquema no diverjan.

**Why this priority**: La regla del proyecto exige que el cambio de esquema y el documento vayan en el mismo PR. No agrega comportamiento, pero sin él el PR no se puede mergear.

**Independent Test**: Revisar el documento: D4 figura cerrada, D14 figura como no necesaria, la tabla de migraciones incluye la 007 y el test de esquema (§9) cubre la regla.

**Acceptance Scenarios**:

1. **Given** el documento, **When** se revisa la §7, **Then** D4 figura cerrada con la regla "un solo Scrum Master por proyecto" y su fecha, y D14 figura como "no necesaria" con el motivo.
2. **Given** el documento, **When** se revisa la §4.3, **Then** la nota "no se impone un único Scrum Master" se reemplaza por la regla y por el índice que la impone.
3. **Given** el documento, **When** se revisa la §9, **Then** la tabla de números reservados incluye la migración 007 y los casos mínimos del test de esquema incluyen la regla.
4. **Given** el esquema migrado, **When** corre el test de esquema, **Then** falla si el documento y la migración divergen en esta regla.

---

### Edge Cases

Redactados en formato EARS.

- El sistema **shall** garantizar a lo sumo un Scrum Master por proyecto, desde el dominio y desde la base.
- **WHEN** se agrega un integrante con rol `scrum_master` a un proyecto que ya tiene uno, **THEN** el sistema **shall** rechazar la operación con 409 y un mensaje que indique que el proyecto ya tiene Scrum Master.
- **WHEN** se agrega un integrante con rol `product_builder`, **THEN** el sistema **shall** aceptarlo sin importar si hay o no Scrum Master.
- **WHILE** un proyecto tiene Scrum Master, "Configurar equipo" (`POST /projects/{id}/members`) **shall** asignar únicamente Product Builders.
- **IF** la persona ya es integrante del proyecto, **THEN** el sistema **shall** responder con el error de integrante duplicado existente, y no con el de Scrum Master repetido.
- **IF** dos pedidos agregan un Scrum Master al mismo proyecto a la vez, **THEN** el sistema **shall** aceptar uno solo y responder 409 al otro.
- **WHEN** se crea un proyecto, **THEN** el sistema **shall** registrar al creador como único Scrum Master en la misma operación que crea el proyecto, de modo que nunca quede un proyecto sin Scrum Master ni con dos.
- **IF** la escritura en la base viola el índice único porque dos pedidos simultáneos pasaron el chequeo previo, **THEN** el sistema **shall** traducir la violación al mismo conflicto 409, y no a un error interno ni a "integrante duplicado".
- **WHEN** la migración corre sobre un proyecto con dos o más Scrum Master, **THEN** el sistema **shall** conservar como Scrum Master al de `creado_en` más antiguo, **shall** pasar a `product_builder` a los demás y **shall** hacerlo antes de crear el índice, en la misma migración.
- **IF** dos Scrum Master de un proyecto tienen el mismo `creado_en`, **THEN** la migración **shall** conservar al de menor `usuario_id`.
- **IF** un proyecto no tiene Scrum Master, **THEN** la migración **shall** dejarlo como está y **shall** no asignarle uno.
- **IF** un proyecto tiene un solo Scrum Master, **THEN** la migración **shall** dejarlo sin tocar.
- **WHEN** la migración pasa a un Scrum Master a `product_builder`, **THEN** esa persona **shall** seguir siendo integrante del proyecto, con sus votos y horas intactos.
- **IF** se intenta cambiar el rol de un integrante existente o reactivarlo, **THEN** el sistema **shall** comportarse como hoy: no existe esa operación (es de HU-16, #175), así que la regla no tiene un camino adicional que cubrir.

## Requirements *(mandatory)*

### Functional Requirements

**Regla de negocio**

- **FR-001**: El sistema MUST garantizar que un proyecto tenga a lo sumo un integrante con rol `scrum_master`.
- **FR-002**: Al crear un proyecto, el sistema MUST dejar al creador como su único Scrum Master y dueño, en la misma operación que crea el proyecto.
- **FR-003**: El sistema MUST rechazar con 409 agregar un integrante con rol `scrum_master` a un proyecto que ya tiene uno. El mensaje MUST indicar con claridad que el proyecto ya tiene Scrum Master.
- **FR-004**: Agregar un integrante con rol `product_builder` MUST seguir funcionando igual que hoy, sin importar quién sea el Scrum Master.
- **FR-005**: "Configurar equipo" (`POST /projects/{id}/members`) MUST asignar solo Product Builders cuando el proyecto ya tiene Scrum Master.
- **FR-006**: Si la persona ya es integrante del proyecto, el sistema MUST responder con el error de integrante duplicado actual (prevalece sobre el de Scrum Master repetido), sin cambiar su rol.
- **FR-007**: Los demás errores de la asignación (proyecto inexistente 404, rol o nombre inválidos 400) MUST mantener su comportamiento actual.

**Esquema**

- **FR-008**: El esquema MUST impedir, con un índice único parcial sobre `integrantes(proyecto_id)` restringido a `rol = 'scrum_master'`, dos Scrum Master en el mismo proyecto, tanto por `INSERT` como por `UPDATE`, y también ante escrituras simultáneas.
- **FR-009**: El índice MUST crearse en una migración numerada, embebida en el binario, con el **número 007** (la 006 es `006_autenticacion.sql` de HU-20, en otra rama).
- **FR-010**: El índice MUST aceptar varios `product_builder` por proyecto y un `scrum_master` en cada proyecto distinto.
- **FR-011**: Una violación del índice MUST traducirse al conflicto de FR-003 y MUST distinguirse del error de integrante duplicado (clave primaria).
- **FR-012**: Antes de crear el índice y dentro de la misma migración, la 007 MUST dejar un solo Scrum Master en cada proyecto que tenga más de uno: el de `creado_en` más antiguo (si empatan, el de menor `usuario_id`), y MUST pasar a `product_builder` a los demás. MUST dejar sin tocar a los proyectos con un solo Scrum Master o con ninguno. MUST no borrar ninguna fila de `integrantes`. MUST explicar este criterio en un comentario de la propia migración.

**Documentación y pruebas**

- **FR-013**: `MODELO-DE-DATOS.md` MUST actualizarse en el mismo PR que la migración: D4 cerrada ("un solo Scrum Master por proyecto"), D14 marcada "no necesaria" porque el dueño es el Scrum Master, nota de §4.3 reemplazada por la regla, el encabezado de §7 sin D4 ni D14 como pendientes, una fila de §5 con la regla, y la 007 agregada a la tabla de §9 con una línea que describe el criterio de desempate de FR-012.
- **FR-014**: El test de esquema (§9, punto 6) MUST incluir los casos de esta regla: segundo `scrum_master` por `INSERT` y por `UPDATE` (rechazados), varios `product_builder` y un `scrum_master` por proyecto distinto (aceptados). Además MUST cubrir la corrección de la 007 sobre una base sembrada antes de migrar: proyecto con 2 o más Scrum Master, empate de `creado_en`, proyecto sin Scrum Master y proyecto con uno solo.
- **FR-015**: Los escenarios BDD (Godog) MUST cubrir: crear un proyecto deja un único Scrum Master; agregar un segundo Scrum Master da 409; agregar un Product Builder tiene éxito; el índice rechaza un segundo Scrum Master insertado por fuera del dominio. La corrección de la 007 sobre proyectos con varios Scrum Master se cubre con el test de migración (en `internal/repository`, junto al de esquema), porque Godog arranca sobre el esquema ya migrado.

### Key Entities *(include if feature involves data)*

- **Integrante (membresía)**: relación persona–proyecto–rol. Cada persona tiene un solo rol por proyecto (ya garantizado) y, desde esta historia, cada proyecto tiene a lo sumo una membresía con rol `scrum_master`. Detalle de tablas en `MODELO-DE-DATOS.md` §4.3.
- **Scrum Master**: el integrante del proyecto con rol `scrum_master`. Es quien creó el proyecto y su dueño. Ejecuta las acciones exclusivas del rol (HU-20 y HU-21).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: El 100% de los proyectos creados tiene exactamente un Scrum Master, el creador.
- **SC-002**: El 100% de los intentos de agregar un segundo Scrum Master es rechazado con un mensaje que dice que el proyecto ya tiene Scrum Master, y deja la lista de integrantes sin cambios.
- **SC-003**: Con dos pedidos simultáneos para agregar un Scrum Master, exactamente uno se acepta.
- **SC-004**: Una escritura directa en la base que intente un segundo Scrum Master es rechazada en el 100% de los casos.
- **SC-005**: Agregar Product Builders sigue funcionando en el 100% de los casos válidos, con y sin Scrum Master.
- **SC-006**: Aplicar la 007 sobre una base con proyectos de 0, 1, 2 y 3 Scrum Master termina sin error en el 100% de los casos y deja a lo sumo un Scrum Master por proyecto, sin borrar ningún integrante.
- **SC-007**: `MODELO-DE-DATOS.md` y la migración no divergen: el test de esquema pasa y D4 y D14 figuran resueltas.

## Assumptions

- **Decisión de producto**: tomada por el Scrum Master el 2026-11-03; es del equipo y no se consulta a la cátedra.
- **Dependencia con HU-20/HU-21**: esta historia no exige login; la regla vale con o sin él. Los permisos por rol (HU-20, HU-21) se apoyan en ella.
- **Numeración de la migración**: 007. Si #236 se mergea antes que HU-20, hay que renumerar (la 006 queda libre para HU-20 o se reasigna): anotarlo en la tabla de §9 en el PR que se mergee segundo. La spec 039 no se renumera.
- **Quién es el dueño**: el Scrum Master del proyecto. No se guarda una referencia aparte al dueño (D14 deja de hacer falta).
- **Cómo se crea hoy el proyecto**: `POST /projects` con `creador`, que inserta el proyecto y la membresía de Scrum Master en una sola transacción. Ese camino ya cumple la regla y no cambia.
- **Cómo se cambia el rol hoy**: no existe. No hay endpoint ni operación para cambiar el rol ni reactivar a un integrante; la clave primaria impide repetir a una persona. Por eso el único camino a cubrir además de la creación es `POST /projects/{id}/members`.
- **Criterio de la migración**: el de `creado_en` más antiguo es una heurística aprobada por el Scrum Master (2026-11-08). Si el creador no fuera el más antiguo, la migración podría dejar a otra persona; se acepta porque la base de Render se vacía en cada deploy y los datos de prueba v2 crean un solo Scrum Master por proyecto.
- **Datos de demostración**: el bloque de datos de prueba carga un solo Scrum Master por proyecto, así que no se ve afectado.
- **Traducción de errores**: hoy el repositorio traduce a "integrante duplicado" cualquier violación `UNIQUE` o de clave primaria; el nuevo índice no puede caer en ese caso (FR-011). El plan los separa por código extendido (research.md §1).
- **Fuera de alcance**: traspasar el rol de Scrum Master y quitar integrantes (HU-16, #175). Mientras no exista el traspaso, el único Scrum Master de un proyecto es siempre su creador. Es un límite declarado, no un defecto de esta historia.
- **Frontend**: sin pantallas nuevas. Si la pantalla de configurar equipo ofrece el rol `scrum_master`, mostrará el mensaje del 409; ajustar la pantalla queda fuera de esta historia.
