# Feature Specification: Validar proyecto existente al crear historia de backlog (Issue #19)

**Feature Branch**: `fix/validar-proyecto-crear-historia`

**Created**: 2026-10-05

**Status**: Draft

**Input**: User description: "Issue #19: CrearHistoriaBacklog no valida que el proyecto exista antes de guardar. Como Product Builder, cuando intento crear una historia de usuario con un proyecto_id que no corresponde a ningún proyecto existente, el sistema no debería crearla — hoy la crea igual, generando una historia huérfana. Objetivo: que CrearHistoriaBacklog rechace la creación de una historia si el proyecto_id no corresponde a un proyecto existente. Salida esperada: el caso de uso devuelve domain.ErrProyectoNoEncontrado y no persiste nada en backlog_items. Referencia: CrearSprint ya hace esta validación llamando a proyectos.ObtenerPorID antes de persistir. Fuera de alcance: agregar FOREIGN KEY en el esquema SQLite."

## Clarifications

### Session 2026-10-05

- Q: ¿Qué debe ocurrir cuando el proyecto no existe pero la historia también tiene datos inválidos (p. ej. título vacío)? → A: Se mantiene el orden actual de validación: primero se valida la construcción de la historia (invariantes del `BacklogItem`) y recién después la existencia del proyecto. Si el título es inválido se devuelve el `ValidationError` de dominio; si la historia es válida pero el proyecto no existe se devuelve `ErrProyectoNoEncontrado`.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Rechazar la creación de una historia si el proyecto no existe (Priority: P1)

Como Product Builder, cuando intento crear una historia de usuario indicando un proyecto que no existe, el sistema rechaza la operación y me informa que el proyecto no fue encontrado, en lugar de guardar la historia. Hoy la historia se guarda igual y queda huérfana, es decir, asociada a un proyecto inexistente.

**Why this priority**: Es el objetivo central de la historia. Una historia huérfana corrompe la integridad de los datos del backlog y no puede consultarse ni asociarse a un proyecto real. Sin esta validación, el resto del flujo de backlog opera sobre datos inconsistentes.

**Independent Test**: Puede probarse de forma independiente ejecutando el caso de uso de creación de historia con un identificador de proyecto inexistente y verificando que devuelve el error de proyecto no encontrado y que el repositorio de backlog no registró ninguna historia.

**Acceptance Scenarios**:

1. **Given** un repositorio de proyectos donde no existe el proyecto indicado, **When** se intenta crear una historia de usuario con ese identificador de proyecto y datos válidos, **Then** el sistema rechaza la operación devolviendo el error "proyecto no encontrado" y no persiste ninguna historia.
2. **Given** un repositorio de proyectos donde existe el proyecto indicado, **When** se crea una historia de usuario con ese identificador y datos válidos, **Then** el sistema crea y persiste la historia asociada al proyecto, como hasta ahora.
3. **Given** un identificador de proyecto inexistente y una historia con datos inválidos (por ejemplo, título vacío), **When** se intenta crearla, **Then** el sistema rechaza la operación con el error de validación de la historia y no llega a consultar la existencia del proyecto.

---

### User Story 2 - Responder 404 en POST /backlog ante un proyecto inexistente (Priority: P2)

Como consumidor de la API, cuando envío la creación de una historia con un proyecto inexistente, recibo una respuesta "no encontrado" (404) en lugar de un error interno (500), para distinguir un dato de entrada inválido de una falla del servidor.

**Why this priority**: Mejora la experiencia de quien consume la API, pero depende de que la validación de US1 exista. El mapeo del error a 404 ya lo resuelve el traductor compartido `escribirError`, así que no genera implementación propia: queda registrado como criterio de aceptación relacionado y se verifica de punta a punta en el quickstart.

**Independent Test**: Puede probarse de forma independiente enviando una solicitud de creación de historia con un proyecto inexistente y verificando que la respuesta HTTP es 404 (y no 500), una vez implementado el mapeo de errores.

**Acceptance Scenarios**:

1. **Given** el mapeo de errores de la API implementado, **When** la creación de una historia falla porque el proyecto no existe, **Then** el endpoint responde con estado HTTP 404.
2. **Given** el mapeo de errores de la API implementado, **When** la creación de una historia falla por una falla interna no prevista, **Then** el endpoint sigue respondiendo 500 (ese mapeo no cambia).

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- **WHEN** se intenta crear una historia de usuario con un identificador de proyecto que no corresponde a ningún proyecto existente, **THEN** el sistema **shall** rechazar la operación devolviendo el error "proyecto no encontrado", sin persistir la historia.
- **WHEN** se intenta crear una historia con un proyecto existente y datos válidos, **THEN** el sistema **shall** persistir la historia exactamente como lo hacía antes de esta corrección.
- **IF** la historia no cumple las invariantes de creación (por ejemplo, título vacío), **THEN** el sistema **shall** devolver el error de validación de la historia antes de verificar la existencia del proyecto.
- **IF** el repositorio de proyectos falla al consultar el proyecto por una causa distinta a "no encontrado", **THEN** el sistema **shall** propagar ese error sin persistir la historia.
- **WHEN** la validación de proyecto falla, **THEN** el sistema **shall** no invocar la persistencia de la historia (cero escrituras en el repositorio de backlog).
- **WHILE** se crea una historia, **THEN** el sistema **shall** validar la existencia del proyecto dentro de la misma operación de creación, sin pasos manuales ni recálculos posteriores.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Cuando se intente crear una historia de usuario con un identificador de proyecto que no corresponde a ningún proyecto existente, el sistema MUST rechazar la operación devolviendo el error "proyecto no encontrado" (`domain.ErrProyectoNoEncontrado`, ya existente y usado por la creación de Sprints).
- **FR-002**: Cuando la creación de la historia sea rechazada por proyecto inexistente, el sistema MUST NOT persistir ninguna historia en el backlog (cero escrituras).
- **FR-003**: Cuando el proyecto indicado exista y la historia sea válida, el sistema MUST crear y persistir la historia, preservando el comportamiento previo.
- **FR-004**: El sistema MUST verificar la existencia del proyecto antes de persistir la historia, siguiendo el mismo patrón ya usado por la creación de Sprints (consultar el repositorio de proyectos por identificador antes de guardar).
- **FR-005**: El sistema MUST mantener el orden de validación actual: primero las invariantes de la historia (título, prioridad, valor de negocio) y luego la existencia del proyecto.
- **FR-006**: Cuando la consulta de existencia del proyecto falle por una causa distinta a "no encontrado", el sistema MUST propagar el error original sin persistir la historia.
- **FR-007***: *(criterio de aceptación relacionado — no genera trabajo de implementación en esta historia)* Cuando la creación de la historia falle por proyecto inexistente, la respuesta esperada del endpoint `POST /backlog` es HTTP 404 en lugar del 500 genérico reportado originalmente en el issue. Este comportamiento ya lo provee el traductor compartido `escribirError` (`internal/http/sprint_handler.go`), que mapea `domain.ErrProyectoNoEncontrado` a 404; por eso no se agregan tareas ni tests de handler en esta historia y se verifica de punta a punta en el quickstart.
- **FR-008**: La adición de restricciones de clave foránea (`FOREIGN KEY`) en el esquema de persistencia MUST NOT formar parte del alcance de esta historia; se aborda en un paso posterior acordado con el Scrum Master.

### Key Entities *(include if feature involves data)*

- **Proyecto** (ya existente): entidad referenciada por la historia de usuario a través de su identificador. Esta historia consulta su existencia; no la crea, modifica ni elimina.
- **Historia del Product Backlog** (ya existente): entidad que se intenta crear. Esta historia agrega una precondición (el proyecto referenciado debe existir) sin cambiar sus atributos ni su ciclo de vida.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: El 100% de los intentos de crear una historia con un proyecto inexistente se rechazan, sin que se registre ninguna historia nueva.
- **SC-002**: El 100% de las creaciones con un proyecto existente y datos válidos siguen funcionando como antes de la corrección (cero regresiones).
- **SC-003**: El sistema deja de generar historias huérfanas: cero historias asociadas a proyectos inexistentes creadas a partir de esta funcionalidad.
- **SC-004**: Ante un proyecto inexistente, el consumidor de la API recibe una respuesta distinguible de "no encontrado" frente a una falla del servidor (404 en lugar de 500), una vez aplicado el mapeo de errores.
- **SC-005**: Todo intento rechazado por proyecto inexistente produce cero escrituras en el almacenamiento del backlog, verificable por inspección del repositorio.

## Assumptions

- El error `domain.ErrProyectoNoEncontrado` ya existe y es el que corresponde devolver; el repositorio de proyectos ya devuelve ese error (o uno equivalente) al consultar un identificador inexistente, como lo aprovecha la creación de Sprints.
- La historia es válida según las invariantes actuales de `BacklogItem`; esta corrección no modifica esas reglas ni sus mensajes.
- El orden de validación se mantiene igual al actual (invariantes de la historia primero, existencia del proyecto después); no se reordena para no alterar el comportamiento observable existente.
- Esta corrección aplica al caso de uso de creación de historia del backlog. Otros casos de uso que creen entidades asociadas a un proyecto pueden tener el mismo problema, pero quedan fuera del alcance salvo que un issue posterior lo indique.
- El mapeo del error a HTTP 404 corresponde a la capa de presentación y ya lo provee el traductor compartido `escribirError` (`internal/http/sprint_handler.go`); esta spec solo lo deja registrado como criterio de aceptación relacionado (FR-007) y no requiere trabajo adicional en la capa HTTP.
- La verificación de existencia del proyecto se realiza contra el repositorio de proyectos existente, sin introducir cachés ni mecanismos nuevos (KISS/YAGNI).
- La integridad referencial en el motor de persistencia (claves foráneas) queda explícitamente fuera de alcance; esta historia resuelve la validación a nivel de aplicación.
