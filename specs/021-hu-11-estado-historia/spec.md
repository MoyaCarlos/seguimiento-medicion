# Feature Specification: Actualización de Estado de una Historia en el Sprint (HU-11)

**Feature Branch**: `feature/HU-11-estado-historia`

**Created**: 2026-10-29

**Status**: Draft

**Input**: User description: "HU-11 (issue #11, sub-issues #79 a #83): actualizar el estado de una historia asignada a un Sprint. Rama feature/HU-11-estado-historia, spec specs/019-hu-11-estado-historia. Criterios (docs/PRODUCT-BACKLOG.md, HU-11): - Una historia asignada a un Sprint Activo con estado 'Nueva' pasa a 'En progreso'. - Una historia 'En progreso' de un Sprint Activo pasa a 'Completada' y el sistema registra su fecha de finalización; queda disponible para el cálculo de Story Points completados (HU-03). - Una historia sin Sprint no puede pasar a 'En progreso' ni a 'Completada': se rechaza con un mensaje que indica que primero debe asignarse a un Sprint activo. - Historia inexistente: 404 antes de validar el cuerpo. Sin migración: historias.estado y historias.fecha_finalizacion existen desde la 001, con el CHECK que ata Completada a su fecha (estado Nueva solo sin Sprint). No cambies el esquema ni MODELO-DE-DATOS.md; ante una duda del modelo, frená y preguntá. Tomá como referencia de estilo la asignación de HU-06 (spec 018): regla pura en el dominio, el servicio solo orquesta, error de validación con campo y mensaje, ruta con la historia en el camino. Dejá para el clarify: Nueva → Completada directo, historias de un Sprint Pendiente o Finalizado, transiciones hacia atrás, estado repetido, y la carrera entre dos cambios simultáneos. Backend primero (servicio, ruta en cmd/api/main.go, escenario BDD con Godog y tests RED/GREEN por capa); el frontend va después, en la misma spec."

## Clarifications

### Session 2026-10-29

- Q: ¿Se permite pasar una historia directamente de «Nueva» a «Completada», sin pasar antes por «En progreso»? → A: No. Solo se acepta la cadena "Nueva → En progreso → Completada"; el salto directo se rechaza con un mensaje que indica que primero debe pasar por "En progreso" y no modifica la historia.
- Q: ¿Qué pasa si la historia pertenece a un Sprint "Pendiente" o "Finalizado" en lugar de uno "Activo"? → A: Se rechaza en ambos casos, con el mensaje "solo se cambia el estado de historias de un Sprint activo", y la historia no se modifica.
- Q: ¿Se permite volver a un estado anterior, es decir pasar de "En progreso" a "Nueva" o reabrir una historia "Completada"? → A: No. El estado solo avanza ("Nueva" → "En progreso" → "Completada"); todo retroceso se rechaza con el mensaje "el estado de una historia solo avanza: no se puede volver a un estado anterior" y no modifica la historia.
- Q: Si se pide el estado que la historia ya tiene, ¿se responde como éxito sin cambios o se rechaza? → A: Se rechaza con el mensaje "la historia ya está en ese estado", sin modificarla (se prefiere el rechazo explícito al éxito idempotente, como en HU-06).
- Q: Si dos personas cambian el estado de la misma historia a la vez, ¿qué cambio prevalece y cómo se evita que quede un estado inconsistente con la fecha de finalización? → A: El cambio se aplica solo si el estado leído sigue siendo el mismo (actualización condicional, sin columnas nuevas). El primero gana; el segundo recibe un conflicto (409) que le pide recargar y reintentar.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Avanzar el estado de una historia dentro del Sprint activo (Priority: P1)

Como Product Builder, durante un Sprint activo cambio el estado de una historia asignada a ese Sprint: la paso a "En progreso" cuando el equipo empieza a trabajarla y a "Completada" cuando termina, para que el avance real del trabajo quede reflejado y alimente las métricas del proyecto (HU-03).

**Why this priority**: Es el núcleo de la historia y la única fuente de datos reales de avance del producto. Sin estas transiciones no hay Story Points completados, ni velocidad, ni porcentaje de historias completadas; además, el cierre de un Sprint (HU-05) necesita el estado de cada historia para decidir cuáles vuelven al Product Backlog.

**Independent Test**: Se prueba sembrando un proyecto con un Sprint "Activo" y una historia "Nueva" asignada a él, pasándola a "En progreso", verificando que conserva su Sprint y sus datos, y luego a "Completada", verificando que queda registrada su fecha de finalización, que sigue asignada al Sprint y que sus Story Points cuentan como completados.

**Acceptance Scenarios**:

1. **Given** una historia asignada a un Sprint "Activo" con estado "Nueva", **When** la paso a "En progreso", **Then** el sistema registra el cambio y la historia conserva su Sprint, su título, su descripción, sus criterios, su prioridad, su valor de negocio y sus Story Points.
2. **Given** una historia asignada a un Sprint "Activo" con estado "En progreso", **When** la marco como "Completada", **Then** el sistema registra el cambio y la fecha de finalización de la historia.
3. **Given** una historia completada de un Sprint "Activo", **When** se calculan las métricas del proyecto (HU-03), **Then** la historia aporta sus Story Points como completados de ese Sprint.
4. **Given** una historia completada de un Sprint "Activo", **When** el Sprint se cierra (HU-05), **Then** la historia permanece en el Sprint y no vuelve al Product Backlog.
5. **Given** un cambio de estado aceptado, **When** se vuelve a consultar la historia, **Then** el nuevo estado (y su fecha de finalización, si corresponde) se conserva después de recargar la aplicación.
6. **Given** una historia asignada a un Sprint "Activo" con estado "Nueva", **When** intento marcarla directamente como "Completada", **Then** el sistema rechaza el cambio con un mensaje que indica que primero debe pasar por "En progreso" y no modifica la historia.
7. **Given** una historia asignada a un Sprint "Pendiente" o "Finalizado", **When** intento cambiar su estado, **Then** el sistema rechaza el cambio con el mensaje "solo se cambia el estado de historias de un Sprint activo" y no modifica la historia.
8. **Given** una historia "En progreso" o "Completada", **When** intento volverla a un estado anterior, **Then** el sistema rechaza el cambio con el mensaje "el estado de una historia solo avanza: no se puede volver a un estado anterior" y la historia queda como estaba.
9. **Given** una historia en un estado, **When** pido ese mismo estado otra vez, **Then** el sistema rechaza el cambio con el mensaje "la historia ya está en ese estado" y la historia queda como estaba.
10. **Given** dos personas que cambian el estado de la misma historia a la vez, **When** ambos pedidos se procesan, **Then** solo uno se aplica y el otro recibe un aviso de conflicto que le pide recargar y reintentar.

---

### User Story 2 - Rechazar el cambio de estado de una historia que no está en un Sprint (Priority: P2)

Como Product Builder, cuando intento pasar a "En progreso" o a "Completada" una historia que todavía está en el Product Backlog (sin Sprint), el sistema rechaza el cambio y me indica que primero debo asignarla a un Sprint activo.

**Why this priority**: Es la otra mitad del criterio y la que evita estados imposibles: por la regla transversal del modelo, "En progreso" y "Completada" solo existen dentro de un Sprint. Además le dice al usuario qué hacer para continuar.

**Independent Test**: Se prueba sembrando una historia sin Sprint en estado "Nueva" e intentando las dos transiciones ("En progreso" y "Completada"), verificando que ambas se rechazan con el mismo mensaje y que la historia queda exactamente como estaba. Se prueba también pedir el cambio de estado de una historia inexistente.

**Acceptance Scenarios**:

1. **Given** una historia sin Sprint (Product Backlog) en estado "Nueva", **When** intento pasarla a "En progreso", **Then** el sistema rechaza el cambio con un mensaje que indica que primero debe asignarse a un Sprint activo y no modifica la historia.
2. **Given** una historia sin Sprint, **When** intento pasarla a "Completada", **Then** el sistema rechaza el cambio con el mismo mensaje y no modifica la historia.
3. **Given** que la historia no existe, **When** intento cambiar su estado, **Then** el sistema responde "historia no encontrada" antes de evaluar el cuerpo o cualquier otra regla.

---

### User Story 3 - Cambiar el estado desde la pantalla del Sprint (Priority: P3)

Como Product Builder, desde la pantalla del Sprint activo veo las historias que lo componen con su estado actual y el control para pasarlas a "En progreso" o a "Completada"; al confirmar, la pantalla refleja el nuevo estado y, si el cambio se rechaza, veo el mensaje que devuelve el sistema.

**Why this priority**: Es la entrega de la interfaz (sub-issue de Frontend), posterior al backend en la misma spec. Por la Definition of Done no frena el cierre del issue padre.

**Independent Test**: Se prueba con la aplicación levantada: se toma una historia "Nueva" del Sprint activo, se la pasa a "En progreso" y luego a "Completada", y se verifica que la pantalla muestra el estado nuevo y la fecha de finalización sin recargar a mano; se provoca un rechazo (por ejemplo, una historia del Product Backlog) y se verifica que se muestra el mensaje del sistema.

**Acceptance Scenarios**:

1. **Given** que estoy en la pantalla de un Sprint con historias, **When** paso una historia "Nueva" a "En progreso", **Then** la pantalla muestra el nuevo estado de esa historia.
2. **Given** que estoy en la pantalla de un Sprint con una historia "En progreso", **When** la marco como "Completada", **Then** la pantalla muestra el estado "Completada" y su fecha de finalización.
3. **Given** que el sistema rechaza un cambio de estado, **When** recibo la respuesta, **Then** la pantalla muestra el mensaje de error del sistema y la historia queda como estaba.

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** reconocer exactamente tres estados de historia: "Nueva", "En progreso" y "Completada".
- **IF** se pide pasar una historia de "Nueva" a "Completada" sin pasar por "En progreso", **THEN** el sistema **shall** rechazar el cambio con un mensaje que indique que primero debe pasar por "En progreso" y **shall** no modificar la historia.
- **IF** la historia pertenece a un Sprint "Pendiente" o "Finalizado", **THEN** el sistema **shall** rechazar el cambio de estado con el mensaje "solo se cambia el estado de historias de un Sprint activo" y **shall** no modificar la historia.
- **IF** se pide un retroceso de estado ("En progreso" → "Nueva" o "Completada" → "En progreso"), **THEN** el sistema **shall** rechazarlo con el mensaje "el estado de una historia solo avanza: no se puede volver a un estado anterior" y **shall** no modificar la historia ni su fecha de finalización.
- **IF** se pide el estado que la historia ya tiene, **THEN** el sistema **shall** rechazarlo con el mensaje "la historia ya está en ese estado" y **shall** no modificarla.
- **IF** un pedido incumple varias reglas a la vez, **THEN** el sistema **shall** responder con la primera según este orden: historia inexistente, cuerpo sin estado o con estado inválido, historia sin Sprint a la que se pide pasar a "En progreso" o "Completada", Sprint no activo, estado repetido, retroceso de estado, salto directo de "Nueva" a "Completada".
- **WHEN** se pide cambiar el estado de una historia que no existe, **THEN** el sistema **shall** responder "historia no encontrada" antes de evaluar el cuerpo o cualquier otra regla (la existencia domina sobre la validación, igual que en HU-04, HU-05, HU-06 y HU-15).
- **IF** el cuerpo del pedido no indica el estado destino o indica un valor que no es uno de los tres estados, **THEN** el sistema **shall** rechazarlo como error de validación que identifica el campo del estado y explica el motivo, después de verificar que la historia existe.
- **WHEN** se cambia el estado de una historia asignada a un Sprint "Activo", **THEN** el sistema **shall** persistir el cambio y no modificar ningún otro dato de la historia (proyecto, Sprint, título, descripción, criterios de aceptación, prioridad, valor de negocio ni Story Points).
- **WHEN** una historia pasa a "Completada", **THEN** el sistema **shall** registrar su fecha de finalización como fecha de calendario (día, sin hora) y **shall** mantener siempre la correspondencia entre el estado "Completada" y la existencia de esa fecha.
- **WHEN** una historia se completa, **THEN** la historia **shall** seguir asignada a su Sprint, de modo que sus Story Points cuenten como completados de ese Sprint (HU-03) y no vuelvan al Product Backlog al cerrarlo (HU-05).
- **IF** un cambio de estado se rechaza, **THEN** el sistema **shall** dejar la historia exactamente como estaba (mismo estado, misma fecha de finalización, mismo Sprint y mismos datos).
- **IF** una historia no tiene Sprint, **THEN** el sistema **shall** rechazar el paso a "En progreso" y a "Completada" con un mensaje que indique que primero debe asignarse a un Sprint activo.
- **WHEN** dos cambios de estado sobre la misma historia se piden a la vez, **THEN** el sistema **shall** aplicar solo uno y **shall** responder al otro con un conflicto que indique que la historia cambió de estado mientras tanto y que debe recargarse y reintentarse; el estado resultante **shall** ser válido y consistente con la fecha de finalización.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir pasar de "Nueva" a "En progreso" una historia asignada a un Sprint "Activo".
- **FR-002**: El sistema MUST permitir pasar de "En progreso" a "Completada" una historia asignada a un Sprint "Activo".
- **FR-003**: Al pasar una historia a "Completada", el sistema MUST registrar su fecha de finalización como fecha de calendario.
- **FR-004**: El sistema MUST rechazar el paso a "En progreso" y a "Completada" de una historia sin Sprint, con un mensaje que indique que primero debe asignarse a un Sprint activo.
- **FR-005**: Todo cambio de estado rechazado MUST dejar la historia exactamente como estaba.
- **FR-006**: Todo cambio de estado aceptado MUST conservar el Sprint asignado y todos los demás datos de la historia (proyecto, título, descripción, criterios de aceptación, prioridad, valor de negocio y Story Points); lo único que cambia es el estado y, al completar, la fecha de finalización.
- **FR-007**: El sistema MUST responder "historia no encontrada" cuando la historia no existe, antes de evaluar el cuerpo o cualquier otra regla. Un recurso inexistente se informa como "no encontrado"; las demás reglas, como un rechazo de la petición.
- **FR-008**: El sistema MUST rechazar con un error de validación que identifica el campo del estado y el motivo cuando el pedido no indica el estado destino o indica un valor que no es "Nueva", "En progreso" ni "Completada".
- **FR-009**: El sistema MUST persistir los cambios de estado aceptados, de modo que se sigan viendo después de cerrar y volver a abrir la aplicación.
- **FR-010**: Una historia "Completada" de un Sprint MUST aportar sus Story Points al cálculo de Story Points completados de ese Sprint (HU-03) y MUST permanecer en el Sprint al cerrarlo (HU-05).
- **FR-011**: El sistema MUST mantener la correspondencia entre el estado "Completada" y la existencia de la fecha de finalización: ninguna historia "Completada" puede quedar sin fecha ni ninguna historia sin completar puede tener fecha.
- **FR-012**: La pantalla del Sprint MUST permitir cambiar el estado de sus historias y mostrar el resultado; ante un rechazo, MUST mostrar el mensaje de error que devuelve el sistema, sin reemplazarlo por un texto genérico.
- **FR-013**: El sistema MUST rechazar el paso directo de "Nueva" a "Completada", con un mensaje que indique que primero debe pasar por "En progreso".
- **FR-014**: El sistema MUST rechazar el cambio de estado de una historia que pertenece a un Sprint "Pendiente" o "Finalizado", con el mensaje "solo se cambia el estado de historias de un Sprint activo", y MUST dejar la historia sin cambios.
- **FR-015**: El sistema MUST rechazar todo retroceso de estado ("En progreso" → "Nueva" o "Completada" → "En progreso"), con el mensaje "el estado de una historia solo avanza: no se puede volver a un estado anterior", y MUST dejar la historia y su fecha de finalización sin cambios.
- **FR-016**: El sistema MUST rechazar el pedido del estado que la historia ya tiene, con el mensaje "la historia ya está en ese estado", y MUST dejar la historia sin cambios.
- **FR-017**: Ante dos cambios de estado simultáneos sobre la misma historia, el sistema MUST aplicar uno solo y MUST informar el otro como conflicto, indicando que la historia cambió de estado mientras tanto y que debe recargarse y reintentarse. El resultado MUST ser siempre un estado válido y consistente con la fecha de finalización.
- **FR-018**: Cuando un pedido incumple varias reglas a la vez, el sistema MUST informar la primera según el orden de los Edge Cases (historia inexistente, cuerpo sin estado o con estado inválido, historia sin Sprint a la que se pide pasar a "En progreso" o "Completada", Sprint no activo, estado repetido, retroceso de estado, salto directo de "Nueva" a "Completada"), para que la respuesta sea siempre la misma ante la misma petición.

### Key Entities *(include if feature involves data)*

- **Historia de Usuario**: la de HU-01/HU-15/HU-06. En esta historia cambian su **estado** (entre "Nueva", "En progreso" y "Completada") y, al completarse, su **fecha de finalización**. Conserva su Sprint asignado y el resto de sus datos.
- **Sprint**: el de HU-05, con estado "Pendiente", "Activo" o "Finalizado". Solo un Sprint "Activo" admite cambios de estado de sus historias; "Pendiente" y "Finalizado" los rechazan.
- **Product Backlog**: vista de las historias sin Sprint (HU-15). Toda historia sin Sprint está en estado "Nueva" y no admite los estados "En progreso" ni "Completada".

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un Product Builder puede pasar una historia de "Nueva" a "Completada" durante un Sprint en menos de 30 segundos y en dos acciones de confirmación, una por transición.
- **SC-002**: El 100% de las historias completadas queda con su fecha de finalización registrada y aporta sus Story Points como completados de su Sprint.
- **SC-003**: El 100% de los intentos de pasar a "En progreso" o a "Completada" una historia sin Sprint se rechaza con un mensaje que indica que primero debe asignarse a un Sprint activo, y la historia queda intacta.
- **SC-004**: El 100% de los pedidos de cambio de estado sobre historias inexistentes responde "historia no encontrada" sin evaluar el cuerpo del pedido.
- **SC-005**: El 100% de los cambios aceptados conserva los Story Points y el resto de los datos de la historia, y se sigue viendo después de recargar la aplicación.
- **SC-006**: Ante la misma petición inválida, el sistema responde siempre con el mismo motivo de rechazo.
- **SC-007**: El 100% de los cambios simultáneos deja un único estado aplicado y consistente con la fecha de finalización, y el pedido que no se aplicó recibe un aviso de conflicto.

## Assumptions

- **Sin cambios de esquema**: `historias.estado` y `historias.fecha_finalizacion` existen desde la migración 001, con las restricciones que limitan el estado a "Nueva", "En progreso" y "Completada", atan la fecha de finalización a "Completada" e impiden "En progreso" o "Completada" sin Sprint (`MODELO-DE-DATOS.md` §4.4). Esta historia no agrega migraciones ni modifica `MODELO-DE-DATOS.md`; el `data-model.md` de la spec enlazará a §4.4, §5 y §6.
- **Reutilización**: se usan las operaciones ya existentes para obtener y actualizar una historia y para obtener el Sprint (HU-06). Si al planificar hace falta una operación nueva de persistencia, se agrega sobre las tablas existentes, sin columnas nuevas.
- Las reglas de flujo (estado del Sprint, estado de la historia) las valida el servicio antes de guardar, como pide `MODELO-DE-DATOS.md` §5; la restricción de la base queda como red de seguridad (una violación que llegue igual a la base se informa como error interno, §5).
- La fecha de finalización es una fecha de calendario (día, sin hora), igual que las demás fechas de calendario del modelo (`MODELO-DE-DATOS.md` §2); se registra el día en que la historia se marca como "Completada".
- Los Story Points de una historia con Sprint "Activo" o "Finalizado" no se editan (D2b); por eso los Story Points de una historia completada no cambian después y las métricas históricas se mantienen reproducibles.
- El cambio de estado es de a una historia por vez. Cambiar el estado de varias historias en una sola operación no se pide y queda fuera de alcance; la pantalla puede repetir la acción.
- **Carrera entre cambios simultáneos**: se resuelve con una actualización condicional que compara el estado leído con el guardado, sin agregar columnas ni migraciones. Si otro cambio ya modificó la historia, el pedido se informa como conflicto (409) en lugar de pisar el estado.
- No hay control de acceso por rol ni identidad de usuario, igual que en HU-01, HU-05, HU-06 y HU-15: cualquier integrante puede cambiar el estado.
- Se respeta el estilo de HU-06 (spec 018) para la forma del pedido y de los errores: la regla vive en el dominio, el servicio solo orquesta, el error de validación indica campo y mensaje, y la historia va en el camino de la ruta. Es una referencia de consistencia, no una decisión de esta spec.
- El frontend (User Story 3) se implementa después del backend, dentro de esta misma spec, cuando la pantalla del Sprint esté disponible en `main`. Hasta entonces se entregan el backend, el escenario BDD y los tests.
- Las cinco decisiones que quedaban abiertas en el primer borrador se resolvieron en la sesión de clarificación del 2026-10-29 y están integradas en `Clarifications`, los escenarios, los casos límite y los requisitos.
