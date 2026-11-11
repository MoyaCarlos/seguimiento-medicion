# Feature Specification: Sprints dentro del plazo del proyecto (#168)

**Feature Branch**: `feature/168-sprint-dentro-del-proyecto`

**Created**: 2026-11-11

**Status**: Draft

**Input**: User description: "Issue #168 — Rechazar Sprints con fechas fuera del plazo del proyecto. Al iniciar un Sprint (donde se cargan el Sprint Goal y las fechas), el backend no compara sus fechas con las del proyecto (HU-04): se puede iniciar un Sprint que empiece antes o termine después del proyecto. Lo encontró el Agile Enabler en el beta testing en Render (2026-10-30): con un proyecto desde el 31/10/2026 se inició un Sprint desde el 30/10/2026."

## Clarifications

### Session 2026-11-11

Decisiones abiertas del issue, resueltas por el equipo antes de especificar:

- Q: El issue habla de `CrearSprint`, pero crear un Sprint no recibe fechas. ¿Dónde se valida? → A: Al **iniciar** el Sprint, que es la única operación que carga sus fechas (pasa de Pendiente a Activo con Goal, inicio y fin). Crear un Sprint no cambia.
- Q: Las fechas del proyecto son opcionales y cada una puede faltar por separado. ¿Contra qué se valida? → A: Contra **cada fecha que exista**. Si el proyecto tiene inicio, el Sprint no puede empezar antes; si tiene fin, el Sprint no puede terminar después; si no tiene ninguna, no se valida.
- Q: ¿Bordes y Sprint parcialmente afuera? → A: Los **bordes se incluyen**: un Sprint puede empezar el mismo día que el proyecto y terminar el mismo día que él. Se rechaza si **cualquiera** de las dos fechas del Sprint queda afuera.
- Q: Si después se acortan las fechas del proyecto (HU-04) y un Sprint ya iniciado queda afuera, ¿qué pasa? → A: **Se tolera**. Queda fuera de alcance y se registra en el issue #302.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Rechazar un Sprint que empieza antes que el proyecto (Priority: P1)

Como Scrum Master, al iniciar un Sprint no puedo ponerle una fecha de inicio anterior al comienzo del proyecto, y recibo un mensaje que me dice cuál fecha corregir, para que el seguimiento y las métricas no cuenten trabajo fuera de la vida del proyecto.

**Why this priority**: Es exactamente el defecto encontrado en el beta testing en Render.

**Independent Test**: Con un proyecto que empieza el 31/10/2026 y un Sprint Pendiente, iniciar el Sprint desde el 30/10/2026: se rechaza con 400 en el campo `fecha_inicio` y el Sprint sigue Pendiente.

**Acceptance Scenarios**:

1. **Given** un proyecto que empieza el 2026-10-31 y un Sprint Pendiente, **When** el Scrum Master lo inicia del 2026-10-30 al 2026-11-06, **Then** el sistema responde 400 con el campo `fecha_inicio` y un mensaje que indica que el Sprint no puede empezar antes que el proyecto, y el Sprint sigue Pendiente, sin Goal ni fechas.
2. **Given** un proyecto que empieza el 2026-10-31 y un Sprint Pendiente, **When** el Scrum Master lo inicia del 2026-10-31 al 2026-11-06, **Then** el Sprint queda Activo (borde incluido).

---

### User Story 2 - Rechazar un Sprint que termina después que el proyecto (Priority: P1)

Como Scrum Master, al iniciar un Sprint no puedo ponerle una fecha de fin posterior al fin del proyecto, y recibo un mensaje que me dice cuál fecha corregir.

**Why this priority**: Es la otra mitad de la misma regla; sin ella el plazo del proyecto se puede exceder igual.

**Independent Test**: Con un proyecto que termina el 2026-11-16 y un Sprint Pendiente, iniciar el Sprint hasta el 2026-11-17: se rechaza con 400 en el campo `fecha_fin`.

**Acceptance Scenarios**:

1. **Given** un proyecto que termina el 2026-11-16 y un Sprint Pendiente, **When** el Scrum Master lo inicia del 2026-11-10 al 2026-11-17, **Then** el sistema responde 400 con el campo `fecha_fin` y un mensaje que indica que el Sprint no puede terminar después que el proyecto, y el Sprint sigue Pendiente.
2. **Given** un proyecto que termina el 2026-11-16 y un Sprint Pendiente, **When** el Scrum Master lo inicia del 2026-11-10 al 2026-11-16, **Then** el Sprint queda Activo (borde incluido).

---

### User Story 3 - Proyectos con fechas incompletas (Priority: P2)

Como Scrum Master de un proyecto que todavía no tiene definidas sus fechas, o solo una de ellas, puedo iniciar Sprints con la única restricción de la fecha que sí está definida.

**Why this priority**: Las fechas del proyecto son opcionales (HU-04); la regla no puede bloquear a los proyectos que no las cargaron.

**Independent Test**: Iniciar Sprints en proyectos sin fechas, solo con inicio y solo con fin, y comprobar que solo se valida la fecha existente.

**Acceptance Scenarios**:

1. **Given** un proyecto sin fecha de inicio ni de fin, **When** el Scrum Master inicia un Sprint con cualquier par de fechas válido, **Then** el Sprint queda Activo.
2. **Given** un proyecto con fecha de inicio 2026-10-31 y sin fecha de fin, **When** el Scrum Master inicia un Sprint del 2026-12-01 al 2027-03-01, **Then** el Sprint queda Activo.
3. **Given** un proyecto con fecha de fin 2026-11-16 y sin fecha de inicio, **When** el Scrum Master inicia un Sprint del 2026-01-01 al 2026-01-15, **Then** el Sprint queda Activo.
4. **Given** un proyecto con fecha de fin 2026-11-16 y sin fecha de inicio, **When** el Scrum Master inicia un Sprint del 2026-11-10 al 2026-11-20, **Then** el sistema responde 400 con el campo `fecha_fin`.

---

### Edge Cases

Redactados en formato EARS.

- **Ubicuo**: El sistema **shall** garantizar que todo Sprint iniciado tenga sus fechas dentro del rango definido del proyecto, con bordes incluidos.
- **Basado en evento**: **WHEN** se inicia un Sprint de un proyecto con fecha de inicio, **THEN** el sistema **shall** rechazar con 400, en el campo `fecha_inicio`, una fecha de inicio del Sprint anterior a la del proyecto.
- **Basado en evento**: **WHEN** se inicia un Sprint de un proyecto con fecha de fin, **THEN** el sistema **shall** rechazar con 400, en el campo `fecha_fin`, una fecha de fin del Sprint posterior a la del proyecto.
- **Basado en estado**: **WHILE** el proyecto no tenga fecha de inicio, el sistema **shall** no imponer un límite inferior a la fecha de inicio del Sprint; **WHILE** no tenga fecha de fin, **shall** no imponer un límite superior a la fecha de fin del Sprint.
- **No deseado**: **IF** las dos fechas del Sprint quedan fuera del proyecto, **THEN** el sistema **shall** informar el error de `fecha_inicio` (se valida primero el inicio, como en las validaciones existentes).
- **No deseado**: **IF** el Sprint ya incumple una regla existente (Goal vacío, fecha de inicio faltante, fin no posterior al inicio), **THEN** el sistema **shall** responder ese error, sin cambios: la regla del plazo se evalúa después.
- **No deseado**: **IF** el Sprint no existe o el proyecto ya tiene un Sprint Activo, **THEN** el sistema **shall** responder como hoy (404, 409), antes de evaluar el plazo. Lo mismo vale para el 403 de quien no es Scrum Master cuando HU-21 (#301) esté mergeada.
- **No deseado**: **IF** se rechaza el inicio por el plazo, **THEN** el Sprint **shall** quedar Pendiente, sin Goal ni fechas guardadas.
- **Basado en evento**: **WHEN** un Sprint empieza exactamente el día de inicio del proyecto o termina exactamente el día de fin, **THEN** el sistema **shall** aceptarlo.
- **Opcional**: **WHERE** se editan después las fechas del proyecto, el sistema **shall** comportarse como hoy (no revisa los Sprints existentes); queda fuera de alcance.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Al iniciar un Sprint, el sistema MUST rechazar con 400 una fecha de inicio del Sprint anterior a la fecha de inicio del proyecto, informando el campo `fecha_inicio` y un mensaje que diga que el Sprint no puede empezar antes que el proyecto.
- **FR-002**: Al iniciar un Sprint, el sistema MUST rechazar con 400 una fecha de fin del Sprint posterior a la fecha de fin del proyecto, informando el campo `fecha_fin` y un mensaje que diga que el Sprint no puede terminar después que el proyecto.
- **FR-003**: Las comparaciones MUST ser por día calendario y con bordes incluidos.
- **FR-004**: El sistema MUST validar solo contra las fechas del proyecto que estén definidas; un proyecto sin fechas no impone restricción.
- **FR-005**: La regla MUST vivir en la lógica de negocio (dominio o caso de uso), no solo en la pantalla, para que valga para cualquier cliente de la API.
- **FR-006**: El orden de los errores existentes (404, 409 por Sprint Activo, 400 de Goal y fechas propias del Sprint; y el 403 de HU-21 cuando #301 esté mergeada) MUST mantenerse; la regla del plazo se evalúa después de ellos.
- **FR-007**: Un inicio rechazado MUST dejar el Sprint sin cambios (Pendiente, sin Goal ni fechas).
- **FR-008**: Crear un Sprint, cerrarlo, editar el proyecto y los Sprints ya iniciados MUST no cambiar de comportamiento.
- **FR-009**: Los escenarios BDD (Godog) MUST cubrir: inicio antes del proyecto (error), fin después del proyecto (error), los dos bordes (aceptados) y un proyecto sin fechas (aceptado). La combinatoria de fechas parciales se cubre con tests unitarios.

### Key Entities *(include if feature involves data)*

- **Proyecto**: tiene fecha de inicio y de fin opcionales (HU-04), que definen su plazo.
- **Sprint**: iteración de un proyecto; recibe Goal, fecha de inicio y fecha de fin al iniciarse. Detalle de tablas en `MODELO-DE-DATOS.md` (sin cambios de esquema).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: El 100% de los intentos de iniciar un Sprint con alguna fecha fuera del plazo definido del proyecto es rechazado, con el campo a corregir y un mensaje claro.
- **SC-002**: El 100% de los Sprints con fechas dentro del plazo, incluidos los que coinciden con los bordes, se inicia sin error.
- **SC-003**: En proyectos sin fechas, iniciar Sprints funciona igual que antes del cambio en el 100% de los casos.
- **SC-004**: Reproducir el caso del beta testing (proyecto desde 31/10/2026, Sprint desde 30/10/2026) da error en lugar de aceptar el Sprint.
- **SC-005**: Ningún escenario BDD ni test existente deja de pasar.

## Assumptions

- **Decisiones de producto**: tomadas por el equipo el 2026-11-11 (ver Clarifications). Las propuestas del issue se confirmaron con un ajuste: con una sola fecha del proyecto se valida contra esa fecha.
- **Formato de las fechas**: ambas llegan por la API como `AAAA-MM-DD` y se interpretan a las 00:00 UTC, así que se comparan por día sin problemas de zona horaria (`MODELO-DE-DATOS.md` D3).
- **Sin cambio de esquema**: no hay migración ni cambios en `MODELO-DE-DATOS.md` más allá de, a lo sumo, mencionar dónde vive la regla.
- **Datos de prueba**: los datos de prueba v2 y los escenarios existentes deben seguir siendo válidos; si alguno inicia Sprints fuera del plazo de su proyecto, se ajustan sus fechas.
- **Frontend**: sin cambios. El formulario de inicio ya muestra el error debajo del campo que informa la API (`fecha_inicio` o `fecha_fin`); ver research.md §8.
- **Fuera de alcance**: la edición posterior de las fechas del proyecto, registrada en el issue #302.
