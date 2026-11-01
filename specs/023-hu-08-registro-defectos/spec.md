# Feature Specification: Registro y Seguimiento de Defectos (HU-08)

**Feature Branch**: `feature/HU-08-registro-defectos`

**Created**: 2026-10-30

**Status**: Draft

**Input**: User description: "HU-08: Registro y Seguimiento de Bugs (issue #8, sub-issue de spec #64). Como integrante del equipo quiero registrar un defecto técnico indicando su severidad, la historia relacionada y el sprint de detección para poder hacer un seguimiento de la calidad y asegurar que se resuelva antes del despliegue final. Escenario 1 (carga de un nuevo defecto): Dado que se detecta un error técnico y el proyecto tiene un Sprint "Activo", cuando completo el formulario de Defectos con la descripción, nivel de severidad (Alta, Media, Baja) y la historia relacionada, y guardo, entonces el Bug se registra asociado a esa historia, con el Sprint "Activo" como Sprint de detección y estado "Pendiente". Escenario 2 (resolución): Dado que un desarrollador solucionó un Bug pendiente, cuando cambia su estado a "Resuelto", entonces el sistema registra el Sprint de resolución (el Sprint Activo en ese momento) para que luego aparezca en el Dashboard de métricas de calidad. Escenario 3 (sin Sprint activo): Dado que el proyecto no tiene Sprint "Activo", cuando intento registrar un defecto o marcar uno como "Resuelto", entonces el sistema rechaza la operación e indica que no hay un Sprint activo al que asociarlo. Depende de HU-01 (historia relacionada) y HU-05 (Sprints). Prioridad Must, valor 8, estimación 8 SP. Modelo de datos ya definido en MODELO-DE-DATOS.md §4.9 (tabla defectos, migración 005): sin título (D6), historia obligatoria, CHECK (estado='Resuelto') = (sprint_resolucion_id IS NOT NULL), FK compuestas con proyecto_id. Rama: feature/HU-08-registro-defectos desde origin/main."

## Clarifications

### Session 2026-10-30

- Q: Cuando un defecto ya está "Resuelto", ¿se puede volver a ponerlo en "Pendiente" (reabrirlo) si el error reaparece? → A: No. "Resuelto" es definitivo; si el error vuelve, se registra un defecto nuevo. Así los defectos resueltos de un Sprint ya cerrado no cambian después.
- Q: Como historia relacionada de un defecto, ¿se puede elegir cualquier historia del proyecto, o solo una que esté en el Sprint "Activo"? → A: Cualquier historia del proyecto, en cualquier Sprint o en el Product Backlog y en cualquier estado. El Sprint de detección indica cuándo se encontró el defecto, no dónde se construyó la historia.
- Q: En el listado de defectos de un proyecto, ¿en qué orden se muestran? → A: Por orden de registro, el más antiguo primero. Agrupar o filtrar por estado o severidad queda del lado del frontend.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Registrar un defecto (Priority: P1)

Como integrante del equipo, cuando detecto un error técnico en el software, completo el formulario de defectos del proyecto con la descripción del error, su severidad (Alta, Media o Baja) y la historia relacionada, y guardo, para que el defecto quede registrado en el Sprint en curso y se pueda seguir hasta su resolución.

**Why this priority**: Es el núcleo de la historia. Sin defectos registrados no hay métrica de defectos detectados (HU-03), ni sección de defectos en el reporte (HU-10), ni nada que resolver (User Story 2).

**Independent Test**: Se prueba sembrando un proyecto con un Sprint "Activo" y una historia; se registra un defecto y se verifica que queda con la descripción, la severidad y la historia indicadas, el Sprint "Activo" como Sprint de detección, el estado "Pendiente" y sin Sprint de resolución.

**Acceptance Scenarios**:

1. **Given** que el proyecto tiene un Sprint "Activo" y una historia, **When** registro un defecto con una descripción, severidad "Alta" y esa historia como relacionada, **Then** el sistema registra el defecto asociado a la historia, con el Sprint "Activo" como Sprint de detección, estado "Pendiente" y sin Sprint de resolución.
2. **Given** que el proyecto tiene un Sprint "Activo", **When** registro un defecto relacionado con una historia que no está en ese Sprint (por ejemplo, una historia completada en un Sprint anterior o una que sigue en el Product Backlog), **Then** el sistema lo registra igual, con el Sprint "Activo" como Sprint de detección.
3. **Given** que el proyecto no tiene ningún Sprint "Activo" (no tiene Sprints, o todos están "Pendiente" o "Finalizado"), **When** intento registrar un defecto, **Then** el sistema rechaza la operación, indica que no hay un Sprint activo al que asociarlo como Sprint de detección y no registra nada.
4. **Given** que intento registrar un defecto, **When** dejo la descripción vacía, elijo una severidad que no es Alta, Media ni Baja, u omito la historia relacionada, **Then** el sistema rechaza la operación con un error que identifica el campo y no registra nada.
5. **Given** que intento registrar un defecto, **When** elijo como relacionada una historia que no existe o que es de otro proyecto, **Then** el sistema rechaza la operación con un error sobre la historia y no registra nada.
6. **Given** que el proyecto no existe, **When** intento registrarle un defecto, **Then** el sistema responde "proyecto no encontrado".

---

### User Story 2 - Marcar un defecto como resuelto (Priority: P1)

Como desarrollador que solucionó un defecto pendiente, cambio su estado a "Resuelto" para que el sistema registre en qué Sprint se resolvió y el defecto cuente como resuelto en las métricas de calidad.

**Why this priority**: El backlog pide explícitamente el ciclo completo (Escenario 2) y la métrica de defectos resueltos por Sprint de HU-03 depende de ese Sprint de resolución. Depende de que exista un defecto (User Story 1).

**Independent Test**: Se prueba sembrando un defecto "Pendiente" en un proyecto con un Sprint "Activo"; se lo marca como resuelto y se verifica que queda en estado "Resuelto" con ese Sprint como Sprint de resolución, y que su Sprint de detección no cambió.

**Acceptance Scenarios**:

1. **Given** que un defecto está "Pendiente" y el proyecto tiene un Sprint "Activo", **When** lo marco como "Resuelto", **Then** el sistema cambia su estado a "Resuelto" y registra el Sprint "Activo" como Sprint de resolución, sin cambiar el Sprint de detección.
2. **Given** que un defecto se detectó en un Sprint ya finalizado y hoy hay otro Sprint "Activo", **When** lo marco como "Resuelto", **Then** el Sprint de resolución es el Sprint "Activo" actual, distinto del de detección.
3. **Given** que un defecto está "Pendiente" y el proyecto no tiene ningún Sprint "Activo", **When** intento marcarlo como "Resuelto", **Then** el sistema rechaza la operación, indica que no hay un Sprint activo al que asociarlo como Sprint de resolución y el defecto sigue "Pendiente".
4. **Given** que un defecto ya está "Resuelto", **When** intento marcarlo como "Resuelto" otra vez, **Then** el sistema rechaza la operación indicando que el defecto ya está resuelto y no cambia su Sprint de resolución.
5. **Given** que el defecto no existe, **When** intento marcarlo como "Resuelto", **Then** el sistema responde "defecto no encontrado".

---

### User Story 3 - Consultar los defectos de un proyecto (Priority: P2)

Como integrante del equipo, consulto los defectos de un proyecto y veo, para cada uno, su descripción, severidad, estado, la historia relacionada, el Sprint de detección y, si está resuelto, el Sprint de resolución, para saber qué queda pendiente antes del despliegue final.

**Why this priority**: Es el "seguimiento" de la historia y lo necesita la pantalla de la sub-issue de Frontend (#68, "Listado de defectos con su estado") para ofrecer la acción de resolver. Depende de que existan defectos (P1).

**Independent Test**: Se prueba sembrando varios defectos en un proyecto (pendientes y resueltos, en distintos Sprints) y uno en otro proyecto; la consulta devuelve solo los del proyecto, con todos sus datos y en orden de registro. Un proyecto sin defectos devuelve una lista vacía.

**Acceptance Scenarios**:

1. **Given** que un proyecto tiene defectos pendientes y resueltos, **When** consulto sus defectos, **Then** el sistema muestra todos los del proyecto, en orden de registro, cada uno con descripción, severidad, estado, historia relacionada, Sprint de detección, Sprint de resolución (vacío si está pendiente) y fecha de registro.
2. **Given** que un proyecto existe y no tiene defectos, **When** consulto sus defectos, **Then** el sistema muestra una lista vacía (no es un error).
3. **Given** que el proyecto no existe, **When** consulto sus defectos, **Then** el sistema responde "proyecto no encontrado".

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** guardar cada defecto con una única historia relacionada y un único Sprint de detección, ambos del mismo proyecto que el defecto.
- El sistema **shall** mantener siempre la coherencia entre estado y Sprint de resolución: un defecto "Pendiente" no tiene Sprint de resolución y uno "Resuelto" siempre lo tiene.
- **WHEN** se registra un defecto válido en un proyecto con Sprint "Activo", **THEN** el sistema **shall** guardarlo en estado "Pendiente", con ese Sprint como Sprint de detección y sin Sprint de resolución.
- **WHEN** se marca como "Resuelto" un defecto "Pendiente" de un proyecto con Sprint "Activo", **THEN** el sistema **shall** registrar ese Sprint como Sprint de resolución, aunque sea el mismo que el de detección.
- **WHILE** el proyecto no tiene ningún Sprint "Activo", **IF** se intenta registrar un defecto o marcar uno como "Resuelto", **THEN** el sistema **shall** rechazar la operación con un mensaje que indique que no hay un Sprint activo al que asociarlo (como Sprint de detección o de resolución, según la operación).
- **WHILE** el proyecto tiene un Sprint "Activo", el sistema **shall** aceptar como historia relacionada cualquier historia del proyecto, esté en ese Sprint, en otro o en el Product Backlog, y cualquiera sea su estado.
- **WHEN** la historia relacionada cambia de Sprint o vuelve al Product Backlog después de registrar el defecto, **THEN** el sistema **shall** mantener el Sprint de detección y el de resolución tal como se registraron; nunca se recalculan.
- **IF** la descripción falta, está vacía, tiene solo espacios o supera los 2000 caracteres, **THEN** el sistema **shall** rechazar el registro con un error sobre la descripción.
- **WHERE** la descripción tiene espacios al principio o al final, el sistema **shall** guardarla sin ellos.
- **IF** la severidad falta o no es exactamente "Alta", "Media" o "Baja", **THEN** el sistema **shall** rechazar el registro con un error sobre la severidad.
- **IF** la historia relacionada falta, no existe o pertenece a otro proyecto, **THEN** el sistema **shall** rechazar el registro con un error sobre la historia.
- **IF** el pedido de registro incluye un estado, un Sprint de detección o un Sprint de resolución, **THEN** el sistema **shall** ignorarlos: el estado inicial es siempre "Pendiente" y los Sprints los asigna el sistema.
- **IF** se intenta marcar como "Resuelto" un defecto que ya está "Resuelto", **THEN** el sistema **shall** rechazar la operación indicando que ya está resuelto y **shall** conservar su Sprint de resolución original.
- **WHEN** se opera sobre un proyecto o un defecto que no existe, **THEN** el sistema **shall** responder "proyecto no encontrado" o "defecto no encontrado" antes de validar cualquier otro dato (la existencia domina sobre la validación, igual que en HU-04, HU-05, HU-06, HU-07 y HU-15).
- **IF** un registro incumple varias reglas a la vez, **THEN** el sistema **shall** responder con la primera según este orden: identificador de proyecto inválido, proyecto inexistente, cuerpo de la petición mal formado, proyecto sin Sprint activo, historia relacionada inválida, descripción inválida, severidad inválida.
- **IF** una resolución incumple varias reglas a la vez, **THEN** el sistema **shall** responder con la primera según este orden: identificador de defecto inválido, defecto inexistente, defecto ya resuelto, proyecto sin Sprint activo.
- **IF** una operación es inválida, **THEN** el sistema **shall** rechazarla entera y los defectos del proyecto **shall** quedar igual que antes.
- **WHEN** el Sprint se cierra mientras se registra o se resuelve un defecto, **THEN** la operación **shall** quedar asociada a ese Sprint si llegó antes del cierre, o rechazada por falta de Sprint activo si llegó después; nunca registrada sin Sprint ni asociada a un Sprint que ya no estaba "Activo".
- **WHEN** dos personas marcan como "Resuelto" el mismo defecto al mismo tiempo, **THEN** el sistema **shall** aceptar solo una de las dos; la otra **shall** recibir el rechazo de "ya está resuelto".

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir registrar un defecto en un proyecto indicando descripción, severidad (Alta, Media o Baja) y la historia relacionada.
- **FR-002**: Todo defecto nuevo MUST quedar en estado "Pendiente", con el Sprint "Activo" del proyecto en ese momento como Sprint de detección y sin Sprint de resolución, junto con la fecha y hora de registro.
- **FR-003**: El sistema MUST rechazar el registro cuando el proyecto no tiene un Sprint "Activo", con un mensaje que indique que no hay un Sprint activo al que asociarlo como Sprint de detección.
- **FR-004**: La historia relacionada MUST ser obligatoria y pertenecer al proyecto del defecto; puede estar en cualquier Sprint o en el Product Backlog y en cualquier estado.
- **FR-005**: El sistema MUST rechazar el registro cuando la descripción falta, está vacía, tiene solo espacios o supera los 2000 caracteres, y MUST guardarla sin los espacios del principio y del final.
- **FR-006**: El sistema MUST rechazar el registro cuando la severidad falta o no es "Alta", "Media" o "Baja".
- **FR-007**: El sistema MUST permitir marcar como "Resuelto" un defecto "Pendiente", registrando como Sprint de resolución el Sprint "Activo" del proyecto en ese momento, sin modificar el Sprint de detección ni el resto de los datos.
- **FR-008**: El sistema MUST rechazar la resolución cuando el proyecto no tiene un Sprint "Activo", con un mensaje que indique que no hay un Sprint activo al que asociarlo como Sprint de resolución, y el defecto MUST seguir "Pendiente".
- **FR-009**: El sistema MUST rechazar la resolución de un defecto que ya está "Resuelto", conservando su Sprint de resolución original.
- **FR-010**: El sistema MUST responder "proyecto no encontrado" o "defecto no encontrado" cuando el proyecto o el defecto no existen, antes de validar cualquier otro dato.
- **FR-011**: Cada rechazo MUST identificar el problema (Sprint activo, historia, descripción, severidad o estado del defecto) y MUST dejar sin cambios los defectos del proyecto.
- **FR-012**: Ante varias reglas incumplidas, el sistema MUST informar la primera según el orden de los Edge Cases, para que la respuesta sea siempre la misma ante la misma petición.
- **FR-013**: El sistema MUST permitir consultar los defectos de un proyecto existente, en orden de registro (el más antiguo primero), mostrando de cada uno la descripción, la severidad, el estado, la historia relacionada, el Sprint de detección, el Sprint de resolución (vacío si está pendiente) y la fecha de registro.
- **FR-014**: Los Sprints de detección y de resolución de un defecto MUST mantenerse aunque la historia relacionada cambie de Sprint después.
- **FR-015**: Los defectos MUST persistir y seguir disponibles después de cerrar y volver a abrir la aplicación.
- **FR-016**: En esta historia los defectos no se editan ni se borran. "Resuelto" es un estado definitivo: un defecto resuelto nunca vuelve a "Pendiente" y, si el error reaparece, se registra un defecto nuevo.

### Key Entities *(include if feature involves data)*

- **Defecto**: un error técnico detectado en el software. Atributos: descripción, severidad (Alta, Media, Baja), estado (Pendiente, Resuelto), historia relacionada, Sprint de detección, Sprint de resolución (solo si está resuelto) y fecha de registro. Pertenece a un proyecto; la historia y los Sprints son siempre del mismo proyecto. No tiene título (D6).
- **Historia de Usuario**: la de HU-01/HU-15. Puede tener varios defectos relacionados; el defecto no cambia la historia.
- **Sprint**: el de HU-05. Solo un Sprint "Activo" puede ser Sprint de detección o de resolución. Los defectos detectados y resueltos en cada Sprint alimentan las métricas de calidad (HU-03), el dashboard (HU-09) y el reporte (HU-10).
- **Proyecto**: el de HU-04. Agrupa sus defectos; tiene como máximo un Sprint "Activo".

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un integrante puede registrar un defecto en menos de 1 minuto y en una sola acción de guardado.
- **SC-002**: El 100% de los defectos registrados quedan con estado "Pendiente" y con el Sprint que estaba "Activo" al registrarlos como Sprint de detección.
- **SC-003**: El 100% de los defectos marcados como resueltos quedan con el Sprint que estaba "Activo" al resolverlos como Sprint de resolución, y su Sprint de detección no cambia.
- **SC-004**: El 100% de los intentos de registrar o resolver sin Sprint "Activo", y de registrar con descripción, severidad o historia inválidas, se rechazan con un mensaje que identifica el problema, y los defectos del proyecto quedan igual.
- **SC-005**: Para cualquier Sprint, la cantidad de defectos detectados y resueltos en él se puede contar sin ambigüedad a partir de los defectos registrados.
- **SC-006**: Los defectos siguen disponibles, con sus Sprints de detección y resolución, después de recargar la aplicación.

## Assumptions

- No hay autenticación: igual que en el resto del sistema, no se registra quién cargó o resolvió el defecto. La tabla diseñada no tiene esa columna y esta historia no la agrega.
- El modelo de datos ya está diseñado en `MODELO-DE-DATOS.md` (§4.9 `defectos`, migración 005, D6): historia obligatoria, severidad y estado como enums, coherencia entre estado y Sprint de resolución garantizada por la base, y claves compuestas que aseguran que historia y Sprints son del mismo proyecto. Esta historia agrega esa migración y no cambia el esquema.
- Que los Sprints de detección y de resolución sean el Sprint "Activo" en ese momento lo valida el servicio, no la base (§5 del modelo de datos, "Dónde vive cada regla": "Servicio (HU-08)").
- Un proyecto tiene como máximo un Sprint "Activo" (HU-05), así que "el Sprint activo" del proyecto está siempre bien definido.
- El largo máximo de la descripción (2000 caracteres) es una regla del dominio, igual que el de la descripción de una historia (HU-01).
- La historia relacionada puede estar fuera del Sprint "Activo": un defecto se puede detectar en una historia completada en un Sprint anterior. El Sprint de detección indica cuándo se detectó, no dónde se construyó la historia.
- La resolución es de un solo sentido y definitiva: si el error reaparece, se registra un defecto nuevo (decidido en Clarifications). Editar o borrar un defecto queda fuera de esta historia (ninguna historia del backlog borra datos, D9); si hace falta, será una historia nueva.
- La consulta de defectos no tiene filtros ni paginación (YAGNI): el volumen de un proyecto académico es chico y el frontend puede filtrar por estado si lo necesita.
- Las métricas de defectos detectados y resueltos por Sprint son de HU-03, y su visualización de HU-09; esta historia solo provee los datos.
- El frontend (formulario, listado con estado, acción de resolver y mensaje sin Sprint activo) es la sub-issue #68 y se entrega en otro Sprint; esta spec cubre el comportamiento del sistema que esa pantalla va a usar.
