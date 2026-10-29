# Feature Specification: Movimiento de Historias al Sprint Backlog (HU-06)

**Feature Branch**: `feature/HU-06-asignar-historias-sprint`

**Created**: 2026-10-28

**Status**: Draft

**Input**: User description: "HU-06 (issue #6, sub-issues #54 a #58): asignar historias del Product Backlog a un Sprint. Como Product Builder quiero seleccionar historias de usuario del Product Backlog y asignarlas a un Sprint específico para poder definir el compromiso de trabajo de la iteración actual. Asignar una historia en estado Nueva a un Sprint Pendiente o Activo: desaparece del Product Backlog y pasa al Sprint Backlog de esa iteración. Rechazar con mensaje claro la asignación de una historia Completada ('solo se pueden planificar historias pendientes') y de una historia que ya tenga Sprint. Sprint inexistente, historia inexistente o Sprint de otro proyecto: error; recurso inexistente = 404 antes que 400. Un Sprint Finalizado no admite asignaciones. Sin migración (historias.sprint_id ya existe en la migración 001). No cambiar los Story Points al asignar (D2b); una historia asignada a un Sprint Activo cuenta como planificada al cerrarlo (D2c). Backend primero; el frontend va después en la misma spec, cuando estén en main GET /projects/{id}/sprints (PR #149) y el listado del backlog en pantalla (#98). Prioridad M (Must have). Valor de Negocio 13. Estimación 5 Story Points. Depende de HU-05, HU-15 y HU-01."

## Clarifications

### Session 2026-10-28

- Q: Si la historia y el Sprint existen pero son de proyectos distintos, ¿la asignación responde "no encontrado" (404) o se rechaza como petición inválida (400)? → A: Se rechaza como petición inválida (400) con un mensaje que indica que el Sprint no pertenece al proyecto de la historia, evaluado después de las existencias de la historia y del Sprint. Los dos recursos existen, así que no es un "no encontrado" (`MODELO-DE-DATOS.md` §5).
- Q: ¿HU-06 incluye en el backend una consulta del Sprint Backlog de un Sprint, o alcanza con lo existente? → A: Se incluye en HU-06 (User Story 2): consulta de solo lectura de las historias de un Sprint, "sprint no encontrado" si no existe y lista vacía si no tiene historias. La consulta de Sprints de un proyecto (PR #149) no devuelve sus historias, y la pantalla la necesita. Reutiliza el listado de historias por Sprint que ya existe, sin cambios de esquema.
- Q: Si se asigna una historia al mismo Sprint en el que ya está, ¿se rechaza o se toma como éxito sin cambios? → A: Se rechaza con "la historia ya pertenece a un Sprint", igual que si el destino fuera otro Sprint. No hay caso idempotente.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Asignar una historia del Product Backlog a un Sprint (Priority: P1)

Como Product Builder, elijo una historia del Product Backlog (sin Sprint, en estado "Nueva") y la asigno a un Sprint "Pendiente" o "Activo" del mismo proyecto, para definir el compromiso de trabajo de esa iteración.

**Why this priority**: Es el núcleo de la historia y del flujo de Scrum del producto. Sin asignación no hay Sprint Backlog, y HU-11 (marcar historias como completadas), HU-07 (imputar horas) y las métricas de HU-03 necesitan historias dentro de un Sprint.

**Independent Test**: Se prueba sembrando un proyecto con un Sprint "Pendiente" (o "Activo") y una historia sin Sprint, asignándola y verificando que deja de aparecer en el Product Backlog, que aparece en el Sprint Backlog de ese Sprint y que conserva su estado "Nueva" y sus Story Points.

**Acceptance Scenarios**:

1. **Given** que existe un Sprint en estado "Pendiente" y una historia del mismo proyecto en el Product Backlog (estado "Nueva"), **When** asigno la historia al Sprint, **Then** la historia desaparece del Product Backlog y pasa a conformar el Sprint Backlog de esa iteración.
2. **Given** que existe un Sprint en estado "Activo" y una historia del mismo proyecto en el Product Backlog, **When** asigno la historia al Sprint, **Then** la historia desaparece del Product Backlog y pasa a conformar el Sprint Backlog de esa iteración.
3. **Given** que asigno una historia a un Sprint, **When** la asignación se completa, **Then** la historia conserva su estado "Nueva", sus Story Points (estén cargados o no) y el resto de sus datos.
4. **Given** que una historia tiene estado "Completada", **When** intento asignarla a un Sprint, **Then** el sistema rechaza la asignación con el mensaje "solo se pueden planificar historias pendientes" y no modifica la historia.
5. **Given** que una historia ya está asignada a un Sprint (cualquiera), **When** intento asignarla a un Sprint, incluso al mismo en el que ya está, **Then** el sistema rechaza la asignación con un mensaje que indica que la historia ya pertenece a un Sprint y no la modifica.
6. **Given** que un Sprint está "Finalizado", **When** intento asignarle una historia del Product Backlog, **Then** el sistema rechaza la asignación con un mensaje que indica que un Sprint finalizado no admite historias y no modifica la historia.
7. **Given** que el Sprint pertenece a otro proyecto que el de la historia, **When** intento asignarla, **Then** el sistema rechaza la asignación con un mensaje que indica que el Sprint no pertenece al proyecto de la historia y no la modifica.
8. **Given** que la historia no existe, **When** intento asignarla a un Sprint (exista o no), **Then** el sistema responde "historia no encontrada".
9. **Given** que la historia existe pero el Sprint no, **When** intento asignarla, **Then** el sistema responde "sprint no encontrado".

---

### User Story 2 - Consultar el Sprint Backlog de una iteración (Priority: P2)

Como Product Builder, consulto el Sprint Backlog de un Sprint y veo las historias que le pertenecen, para saber cuál es el compromiso de trabajo de esa iteración.

**Why this priority**: Es la otra mitad visible de la asignación: el criterio de HU-06 dice que la historia "pasa a conformar el Sprint Backlog de esa iteración", y la pantalla necesita mostrarlo. Sin la asignación (P1) estaría siempre vacío.

**Independent Test**: Se prueba sembrando un Sprint con historias asignadas y otro sin historias, y verificando que la consulta de cada uno devuelve exactamente sus historias, en el orden definido, y un listado vacío para el segundo.

**Acceptance Scenarios**:

1. **Given** que un Sprint tiene historias asignadas, **When** consulto su Sprint Backlog, **Then** el sistema muestra solo las historias de ese Sprint, ordenadas por prioridad (Must, Should, Could, Won't) y, dentro de la misma prioridad, por orden de creación.
2. **Given** que un Sprint existe pero no tiene historias, **When** consulto su Sprint Backlog, **Then** el sistema muestra un Sprint Backlog vacío (no es un error).
3. **Given** que el Sprint no existe, **When** consulto su Sprint Backlog, **Then** el sistema responde "sprint no encontrado".

---

### User Story 3 - Planificar el Sprint desde la pantalla (Priority: P3)

Como Product Builder, desde la pantalla del proyecto elijo una historia del Product Backlog sin Sprint, elijo un Sprint "Pendiente" o "Activo" y la asigno; veo el Sprint Backlog de cada iteración, y si la asignación se rechaza veo el mensaje que devuelve el sistema.

**Why this priority**: Es la entrega de la interfaz (sub-issue de Frontend). Depende de que estén en `main` la consulta de los Sprints de un proyecto (PR #149) y el listado del backlog en pantalla (#98). Por la Definition of Done, no frena el cierre del issue padre.

**Independent Test**: Se prueba con la aplicación levantada: se asigna una historia del backlog a un Sprint y se verifica que desaparece del listado del backlog y aparece en el Sprint Backlog de ese Sprint sin recargar a mano; se provoca un rechazo (por ejemplo, un Sprint que se finalizó mientras tanto) y se verifica que se muestra el mensaje del sistema.

**Acceptance Scenarios**:

1. **Given** que estoy en la pantalla de un proyecto con historias en el Product Backlog y un Sprint "Pendiente" o "Activo", **When** elijo una historia y un Sprint y confirmo la asignación, **Then** la historia deja de verse en el Product Backlog y se ve en el Sprint Backlog de ese Sprint.
2. **Given** que el proyecto tiene varios Sprints, **When** consulto la pantalla, **Then** puedo ver el Sprint Backlog de cada uno.
3. **Given** que elijo el destino de una asignación, **When** se listan los Sprints disponibles, **Then** solo se ofrecen los Sprints "Pendiente" y "Activo" del proyecto.
4. **Given** que el sistema rechaza una asignación, **When** recibo la respuesta, **Then** la pantalla muestra el mensaje de error del sistema y el backlog y los Sprint Backlogs quedan como estaban.

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** permitir asignar a un Sprint únicamente historias sin Sprint asignado; por la regla transversal del modelo de datos, toda historia sin Sprint está en estado "Nueva".
- **WHEN** se asigna una historia a un Sprint "Pendiente" o "Activo" de su mismo proyecto, **THEN** el sistema **shall** asociarla a ese Sprint, sacarla del Product Backlog y conservar su estado "Nueva".
- **WHEN** se asigna una historia, **THEN** el sistema **shall** conservar su identificador, proyecto, título, descripción, criterios de aceptación, prioridad, valor de negocio, estado y Story Points (D2b: asignar no modifica los Story Points).
- **WHERE** la historia no tiene Story Points cargados, el sistema **shall** permitir asignarla igual; la estimación se puede cargar después (D2b permite cargarla por primera vez con el Sprint "Activo").
- **WHEN** se asigna una historia que no existe, **THEN** el sistema **shall** responder "historia no encontrada" antes de evaluar el Sprint o cualquier otra regla (la existencia domina sobre la validación, igual que en HU-04, HU-05 y HU-15).
- **IF** falta el Sprint destino (cuerpo mal formado o sin sprint_id), **THEN** el sistema **shall** rechazar la asignación con un 400 que indique el campo, después de verificar que la historia existe.
- **WHEN** se asigna una historia existente a un Sprint que no existe, **THEN** el sistema **shall** responder "sprint no encontrado" antes de evaluar cualquier otra regla.
- **IF** la historia y el Sprint existen pero pertenecen a proyectos distintos, **THEN** el sistema **shall** rechazar la asignación con un mensaje que indique que el Sprint no pertenece al proyecto de la historia y no modificar la historia.
- **IF** el Sprint está "Finalizado", **THEN** el sistema **shall** rechazar la asignación con un mensaje que indique que un Sprint finalizado no admite historias y no modificar la historia.
- **IF** la historia está "Completada", **THEN** el sistema **shall** rechazar la asignación con el mensaje "solo se pueden planificar historias pendientes" y no modificar la historia. Este mensaje tiene prioridad sobre el de "ya pertenece a un Sprint" (una historia completada siempre tiene Sprint).
- **IF** la historia ya tiene un Sprint asignado y no está "Completada" (está "Nueva" o "En progreso" dentro de un Sprint), **THEN** el sistema **shall** rechazar la asignación con un mensaje que indique que la historia ya pertenece a un Sprint, aunque el destino sea el mismo Sprint, y no modificar la historia.
- **IF** una asignación incumple varias reglas a la vez, **THEN** el sistema **shall** responder con la primera según este orden: historia inexistente, Sprint destino faltante (cuerpo mal formado o sin sprint_id), Sprint inexistente, Sprint de otro proyecto, Sprint finalizado, historia completada, historia que ya tiene Sprint.
- **WHILE** un Sprint está "Activo", **WHEN** se le asigna una historia, **THEN** el sistema **shall** contarla como planificada al cerrar ese Sprint, igual que las asignadas con el Sprint "Pendiente" (D2c: "comprometidas" = lo asignado al cerrar). Esta historia no agrega cálculo: lo resuelven el cierre (HU-05) y las métricas (HU-03).
- **WHEN** se consulta el Sprint Backlog de un Sprint, **THEN** el sistema **shall** mostrar las historias que tienen ese Sprint asignado, en orden de prioridad MoSCoW y luego de creación, con los mismos datos que el Product Backlog.
- **WHEN** se consulta el Sprint Backlog de un Sprint "Finalizado", **THEN** el sistema **shall** mostrar las historias que siguen asignadas a él (las completadas); las no completadas ya volvieron al Product Backlog al cerrarlo (HU-05).
- **WHEN** se consulta el Sprint Backlog de un Sprint que no existe, **THEN** el sistema **shall** responder "sprint no encontrado".
- **WHEN** dos personas asignan la misma historia al mismo tiempo, **THEN** la historia **shall** quedar asignada a un único Sprint, nunca a dos (una historia tiene un solo Sprint asignado). Cómo se resuelve la carrera se define en el plan.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir asignar una historia sin Sprint a un Sprint "Pendiente" o "Activo" de su mismo proyecto.
- **FR-002**: Al asignarse, la historia MUST dejar de aparecer en el Product Backlog del proyecto (HU-15) y MUST aparecer en el Sprint Backlog de ese Sprint.
- **FR-003**: La asignación MUST conservar el estado "Nueva" de la historia y todos sus demás datos, incluidos los Story Points (D2b). Lo único que cambia es el Sprint asignado.
- **FR-004**: El sistema MUST rechazar la asignación de una historia "Completada" con el mensaje "solo se pueden planificar historias pendientes".
- **FR-005**: El sistema MUST rechazar la asignación de una historia que ya tiene un Sprint asignado, sea el mismo Sprint de destino u otro, con el mensaje "la historia ya pertenece a un Sprint". Mover una historia de un Sprint a otro o devolverla al Product Backlog queda fuera de esta historia (D5).
- **FR-006**: El sistema MUST rechazar toda asignación a un Sprint "Finalizado", con un mensaje que lo indique.
- **FR-007**: El sistema MUST rechazar la asignación como petición inválida (no como "no encontrado") cuando el Sprint pertenece a un proyecto distinto del de la historia, con un mensaje que indique que el Sprint no pertenece al proyecto de la historia.
- **FR-008**: El sistema MUST responder "historia no encontrada" o "sprint no encontrado" cuando el recurso no existe, antes de evaluar cualquier otra regla. Un recurso inexistente se informa como "no encontrado"; las demás reglas, como un rechazo de la petición.
- **FR-009**: Cuando una asignación incumple varias reglas, el sistema MUST informar la primera según el orden de los Edge Cases (historia inexistente, Sprint destino faltante, Sprint inexistente, Sprint de otro proyecto, Sprint finalizado, historia completada, historia que ya tiene Sprint), para que la respuesta sea siempre la misma ante la misma petición.
- **FR-010**: Una asignación rechazada MUST dejar la historia exactamente como estaba.
- **FR-011**: El sistema MUST permitir consultar el Sprint Backlog de un Sprint existente: sus historias, ordenadas por prioridad MoSCoW (Must, Should, Could, Won't) y, dentro de la misma prioridad, por orden de creación, con los mismos datos que muestra el Product Backlog.
- **FR-012**: El sistema MUST mostrar un Sprint Backlog vacío, sin error, cuando el Sprint existe y no tiene historias.
- **FR-013**: Las asignaciones válidas MUST persistir y seguir disponibles después de cerrar y volver a abrir la aplicación.
- **FR-014**: La pantalla del proyecto MUST permitir elegir una historia del Product Backlog y un Sprint "Pendiente" o "Activo" del proyecto, confirmar la asignación y ver el Sprint Backlog de cada Sprint del proyecto.
- **FR-015**: Ante un rechazo, la pantalla MUST mostrar el mensaje de error que devuelve el sistema, sin reemplazarlo por un texto genérico.

### Key Entities *(include if feature involves data)*

- **Historia de Usuario**: la de HU-01/HU-15. En esta historia solo cambia su **Sprint asignado** (de vacío a un Sprint). Sin Sprint, está en el Product Backlog y siempre en estado "Nueva".
- **Sprint**: el de HU-05, con estado "Pendiente", "Activo" o "Finalizado" y perteneciente a un proyecto. Solo los dos primeros admiten asignaciones.
- **Product Backlog**: vista de las historias de un proyecto sin Sprint (HU-15). No se guarda: se arma al consultar.
- **Sprint Backlog**: vista de las historias que tienen asignado un Sprint dado. No se guarda: se arma al consultar.
- **Proyecto**: contenedor de historias y Sprints; una historia solo puede entrar a Sprints de su propio proyecto.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un Product Builder puede pasar una historia del Product Backlog al Sprint Backlog de una iteración en menos de 30 segundos y en una sola acción de confirmación.
- **SC-002**: El 100% de las historias asignadas dejan de verse en el Product Backlog y se ven en el Sprint Backlog de su Sprint, también después de recargar la aplicación.
- **SC-003**: El 100% de las asignaciones de historias completadas, de historias que ya tienen Sprint, a Sprints finalizados o a Sprints de otro proyecto se rechazan con un mensaje que explica el motivo y dejan la historia exactamente como estaba.
- **SC-004**: El 100% de las asignaciones válidas conservan el estado "Nueva" y los Story Points que tenía la historia.
- **SC-005**: Ante la misma petición inválida, el sistema responde siempre con el mismo motivo de rechazo.

## Assumptions

- **Sin cambios de esquema**: la columna del Sprint asignado de una historia ya existe en la migración 001, con la restricción que impide asociarla a un Sprint de otro proyecto (`MODELO-DE-DATOS.md` §4.4). Esta historia no agrega migraciones ni modifica `MODELO-DE-DATOS.md`; el `data-model.md` de la spec enlaza a §4.4 y §4.5.
- **Reutilización**: se usan las operaciones ya existentes para obtener y actualizar una historia, obtener un Sprint y listar las historias de un Sprint (esta última, para la consulta del Sprint Backlog de la User Story 2). Si al planificar hace falta una operación nueva de persistencia, se agrega sobre las tablas existentes, sin columnas nuevas.
- Las reglas de flujo (estado del Sprint, estado de la historia, mismo proyecto) las valida el servicio antes de guardar, como pide `MODELO-DE-DATOS.md` §5; la restricción de la base queda como red de seguridad (una violación que llegue igual a la base se informa como conflicto de integridad, §5).
- Mover una historia entre Sprints o devolverla al Product Backlog queda fuera de alcance (D5 de `MODELO-DE-DATOS.md`): HU-06 no lo necesita y ninguna historia lo pide.
- La asignación es de a una historia por vez. Seleccionar varias y asignarlas juntas no se pide y queda fuera de alcance; la pantalla puede repetir la acción.
- El criterio de HU-06 que menciona "arrastrar" se cumple seleccionando historia y Sprint; el arrastrar y soltar es opcional en la interfaz.
- No hay control de acceso por rol: igual que en HU-01 y HU-15, el sistema todavía no tiene identidad de usuario, así que cualquier integrante puede asignar.
- No hay tope de historias ni de Story Points por Sprint: el backlog no lo pide.
- El orden del Sprint Backlog es el mismo que el del Product Backlog (prioridad M, S, C, W y luego identificador), para que ambas vistas se lean igual.
- El frontend (User Story 3) se implementa cuando estén en `main` la consulta de Sprints de un proyecto (PR #149) y el listado del backlog en pantalla (#98). Hasta entonces se entregan el backend, el escenario BDD y los tests.
- Las historias que vuelven al Product Backlog al cerrar un Sprint (arrastre, HU-05 / #117) se pueden volver a asignar a otro Sprint como cualquier historia "Nueva" sin Sprint.
