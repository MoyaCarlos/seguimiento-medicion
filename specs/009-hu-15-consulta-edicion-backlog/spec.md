# Feature Specification: Consulta y Edición del Product Backlog (HU-15)

**Feature Branch**: `feature/HU-15-consulta-edicion-backlog`

**Created**: 2026-10-13

**Status**: Draft

**Input**: User description: "HU-15: Consulta y Edición del Product Backlog. Como Product Builder quiero ver el Product Backlog de un proyecto ordenado por prioridad y poder editar el título, la descripción, los criterios de aceptación y los Story Points de cada historia para poder mantenerlo actualizado y cumplir con el contenido mínimo que pide el enunciado (puntos 2 y 4: estimar historias mediante Story Points). Prioridad M (Must have). Estado Nueva. Valor de Negocio 13 (Fibonacci). Estimación 5 Story Points. Depende de HU-01, HU-04 y HU-05."

## Clarifications

### Session 2026-10-13

- Q: ¿Qué valores exactos componen la escala Fibonacci de los Story Points? → A: 1, 2, 3, 5, 8, 13, 21 (la misma escala que el Valor de Negocio de HU-01). HU-02 y HU-14 deben usar la misma baraja.
- Q: ¿Los criterios de aceptación se guardan como un texto libre o como una lista de criterios? → A: Un solo texto libre opcional, máximo 2000 caracteres, con el mismo límite que la descripción. Vacío o ausente se guarda como NULL. Si más adelante hace falta tratar los criterios por separado, se migra a una tabla hija.
- Q: ¿La edición incluye también la prioridad y el valor de negocio? → A: No. Solo se editan título, descripción, criterios de aceptación y Story Points, como dice la historia. Si hace falta repriorizar, se agrega como un cambio aparte.

### Session 2026-10-25

- Q: ¿Se pueden cambiar los Story Points de una historia asignada a un Sprint `Activo` o `Finalizado`? → A: No, si ya estaban cargados. Es la decisión D2b del equipo (aprobada el 2026-10-13, aclarada el 2026-10-14), registrada en `MODELO-DE-DATOS.md` (§4.4, §5 y D2b) y en las decisiones transversales y el Escenario 3 de HU-15 de `docs/PRODUCT-BACKLOG.md`: "No se modifican los Story Points ya cargados de una historia asignada a un Sprint `Activo` o `Finalizado`, para que las métricas históricas se puedan reproducir. Cargarlos por primera vez (la historia estaba sin estimar) sí se permite." Es una regla de servicio, compartida con HU-14.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Consultar el Product Backlog de un proyecto (Priority: P1)

Como Product Builder, entro al Product Backlog de un proyecto y veo las historias que todavía no están asignadas a ningún Sprint, ordenadas por prioridad, para saber qué trabajo queda pendiente y en qué orden conviene abordarlo.

**Why this priority**: Sin el listado no hay forma de ver el backlog. HU-06 (mover historias al Sprint Backlog) lo necesita para elegir qué historias entran en un Sprint, y es la base sobre la que se edita cada historia.

**Independent Test**: Se prueba sembrando un proyecto con historias de distintas prioridades, algunas asignadas a un Sprint, y verificando que el listado muestra solo las que no tienen Sprint, en el orden esperado.

**Acceptance Scenarios**:

1. **Given** que un proyecto tiene historias, algunas ya asignadas a un Sprint, **When** consulto su Product Backlog, **Then** el sistema muestra solo las historias sin Sprint asignado, ordenadas por prioridad (MoSCoW: Must, Should, Could, Won't) y, dentro de la misma prioridad, por orden de creación.
2. **Given** que un proyecto existe pero no tiene historias sin Sprint, **When** consulto su Product Backlog, **Then** el sistema muestra un backlog vacío (no es un error).
3. **Given** que el proyecto no existe, **When** consulto su Product Backlog, **Then** el sistema responde "no encontrado".

---

### User Story 2 - Editar una historia del Product Backlog (Priority: P1)

Como Product Builder, modifico el título, la descripción, los criterios de aceptación o los Story Points de una historia y guardo, para mantener el backlog actualizado a medida que el equipo refina las historias.

**Why this priority**: El enunciado exige criterios de aceptación y Story Points como contenido mínimo de cada elemento del Product Backlog. HU-03 (métricas) y HU-10 (reporte) necesitan los Story Points de cada historia.

**Independent Test**: Se prueba editando una historia existente con datos válidos y verificando que los cambios quedan persistidos y que la historia conserva su identificador, estado, prioridad y Sprint.

**Acceptance Scenarios**:

1. **Given** que existe una historia en el Product Backlog, **When** modifico su título, descripción, criterios de aceptación o Story Points (valor de la escala Fibonacci) con datos válidos y guardo, **Then** el sistema persiste los cambios y conserva su identificador y estado.
2. **Given** que estoy editando una historia, **When** dejo el título en blanco (o compuesto solo por espacios), **Then** el sistema rechaza el cambio, muestra una advertencia y no modifica la historia.
3. **Given** que estoy editando una historia, **When** ingreso unos Story Points que no pertenecen a la escala Fibonacci, **Then** el sistema rechaza el cambio, muestra una advertencia y no modifica la historia.
4. **Given** que intento editar una historia que no existe, **When** envío la edición (aunque sus datos también sean inválidos), **Then** el sistema responde "no encontrado".
5. **Given** que una historia con Story Points ya cargados está asignada a un Sprint "Activo" o "Finalizado", **When** intento cambiar sus Story Points, **Then** el sistema rechaza el cambio, muestra una advertencia y no modifica la historia.
6. **Given** que una historia sin estimar está asignada a un Sprint "Activo" o "Finalizado", **When** cargo sus Story Points con un valor de la escala Fibonacci, **Then** el sistema persiste la estimación.

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** listar el Product Backlog con prioridad Must antes que Should, Should antes que Could y Could antes que Won't; ante la misma prioridad, la historia creada antes va primero.
- **WHEN** se consulta el Product Backlog de un proyecto, **THEN** el sistema **shall** excluir toda historia que tenga un Sprint asignado, cualquiera sea el estado de ese Sprint.
- **WHEN** se consulta el Product Backlog de un proyecto que no existe, **THEN** el sistema **shall** responder "proyecto no encontrado".
- **WHEN** se edita una historia que no existe, **THEN** el sistema **shall** responder "historia no encontrada" antes de validar los datos enviados (la existencia domina sobre la validación, igual que en HU-04 y HU-05).
- **IF** el título editado queda vacío, compuesto solo por espacios o supera los 200 caracteres, **THEN** el sistema **shall** rechazar la edición con una advertencia sobre el título y no modificar la historia.
- **IF** la descripción editada queda vacía o supera los 2000 caracteres, **THEN** el sistema **shall** rechazar la edición con una advertencia sobre la descripción y no modificar la historia (mismas reglas que al crearla en HU-01).
- **IF** los Story Points informados no pertenecen a la escala Fibonacci (1, 2, 3, 5, 8, 13, 21), incluidos 0, valores negativos y números no enteros, **THEN** el sistema **shall** rechazar la edición con una advertencia sobre los Story Points y no modificar la historia.
- **WHERE** los Story Points se envían vacíos, el sistema **shall** dejar la historia sin estimación (la estimación es opcional hasta que se acuerde en Planning Poker, HU-14).
- **IF** los criterios de aceptación superan los 2000 caracteres, **THEN** el sistema **shall** rechazar la edición con una advertencia sobre los criterios y no modificar la historia.
- **WHERE** los criterios de aceptación se envían vacíos, compuestos solo por espacios o sin informar, el sistema **shall** dejar la historia sin criterios de aceptación.
- **IF** una edición tiene varios campos inválidos, **THEN** el sistema **shall** rechazarla entera, sin guardar parcialmente los campos válidos.
- **WHILE** una historia está asignada a un Sprint `Activo` o `Finalizado` y tiene Story Points cargados, **IF** la edición informa Story Points distintos de los cargados (incluido dejarlos vacíos), **THEN** el sistema **shall** rechazar la edición con una advertencia sobre los Story Points y no modificar la historia (D2b).
- **WHILE** una historia está asignada a un Sprint `Activo` o `Finalizado` y no tiene Story Points cargados, **WHEN** la edición informa Story Points válidos, **THEN** el sistema **shall** guardarlos (D2b: cargarlos por primera vez se permite).
- **WHILE** una historia está asignada a un Sprint `Pendiente`, el sistema **shall** permitir cambiar sus Story Points como si estuviera en el Product Backlog.
- **WHEN** se edita una historia, **THEN** el sistema **shall** conservar su identificador, proyecto, estado y Sprint asignado.
- **WHEN** dos personas editan la misma historia al mismo tiempo, **THEN** el sistema **shall** conservar el último guardado válido (sin versionado ni historial, mismo criterio que HU-04).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir consultar el Product Backlog de un proyecto existente.
- **FR-002**: El Product Backlog MUST incluir solo las historias del proyecto que no tienen Sprint asignado.
- **FR-003**: El Product Backlog MUST ordenarse por prioridad MoSCoW (Must, Should, Could, Won't) y, dentro de la misma prioridad, por orden de creación (de la más antigua a la más nueva).
- **FR-004**: Cada historia del listado MUST mostrar su identificador, título, descripción, criterios de aceptación, prioridad, estado, valor de negocio y Story Points (los dos últimos pueden estar vacíos).
- **FR-005**: El sistema MUST mostrar un backlog vacío, sin error, cuando el proyecto existe y no tiene historias sin Sprint.
- **FR-006**: El sistema MUST permitir editar el título, la descripción, los criterios de aceptación y los Story Points de una historia existente.
- **FR-007**: La edición MUST validar el título y la descripción con las mismas reglas que la creación (HU-01): obligatorios, sin ser solo espacios, con un máximo de 200 y 2000 caracteres respectivamente.
- **FR-008**: La edición MUST aceptar los Story Points vacíos o con un valor de la escala Fibonacci, y rechazar cualquier otro valor. La escala es 1, 2, 3, 5, 8, 13, 21 (la misma del Valor de Negocio, HU-01).
- **FR-009**: El sistema MUST registrar los criterios de aceptación de cada historia como un único texto libre opcional, de hasta 2000 caracteres (el mismo límite que la descripción). Si se envían vacíos o no se informan, la historia queda sin criterios de aceptación.
- **FR-010**: Una edición inválida MUST rechazarse entera con una advertencia que identifique el campo con problemas, y la historia MUST quedar sin modificaciones.
- **FR-011**: La edición MUST conservar el identificador, el proyecto, el estado y el Sprint asignado de la historia. La prioridad y el valor de negocio no se editan en esta historia.
- **FR-012**: El sistema MUST responder "no encontrado" cuando el proyecto consultado o la historia editada no existen, antes de validar cualquier dato de la petición.
- **FR-013**: Los cambios de una edición válida MUST persistir y seguir disponibles después de cerrar y volver a abrir la aplicación.
- **FR-014**: Las historias creadas antes de esta funcionalidad MUST aparecer en el listado con los criterios de aceptación vacíos.
- **FR-015**: El sistema MUST rechazar, con una advertencia sobre los Story Points y sin modificar la historia, toda edición que cambie los Story Points ya cargados de una historia asignada a un Sprint `Activo` o `Finalizado`. Cargarlos cuando la historia no tenía estimación MUST permitirse (decisión D2b del equipo).

### Key Entities *(include if feature involves data)*

- **Historia de Usuario**: la misma de HU-01, ampliada con **criterios de aceptación**. Atributos que se editan en esta historia: título, descripción, criterios de aceptación y Story Points. Atributos que no cambian al editar: identificador, proyecto, estado, prioridad, valor de negocio y Sprint asignado.
- **Product Backlog**: vista de las historias de un proyecto sin Sprint asignado, ordenada por prioridad y orden de creación. No se guarda: se arma al consultar.
- **Proyecto**: contenedor de las historias (HU-04); tiene que existir para consultar su backlog.
- **Sprint**: solo interesa si la historia tiene uno asignado (HU-05); en ese caso queda fuera del Product Backlog.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un Product Builder puede ver el Product Backlog completo de un proyecto en un solo paso, sin filtrar ni ordenar a mano.
- **SC-002**: El 100% de las historias con Sprint asignado quedan fuera del listado y el 100% de las historias sin Sprint aparecen, en el orden de prioridad y creación definido.
- **SC-003**: Un Product Builder puede corregir el título, la descripción, los criterios o los Story Points de una historia en menos de 1 minuto.
- **SC-004**: El 100% de las ediciones con título vacío, Story Points fuera de la escala o que cambian los Story Points ya cargados de una historia de un Sprint `Activo` o `Finalizado` se rechazan con una advertencia y dejan la historia exactamente como estaba.
- **SC-005**: El 100% de las ediciones válidas conservan el identificador y el estado de la historia, y siguen disponibles después de recargar la aplicación.

## Assumptions

- No hay control de acceso por rol: igual que en HU-01, el sistema todavía no tiene identidad de usuario, así que cualquier integrante puede consultar y editar.
- La asignación de historias a un Sprint llega con HU-06. Hasta entonces, el listado se prueba sembrando historias con Sprint directamente en la base de datos o en los fakes de los tests.
- La edición se hace por identificador de historia y aplica a cualquier historia existente, esté o no en el Product Backlog, con la restricción de D2b sobre los Story Points (FR-015). Cambiar el estado lo cubre HU-11 y cambiar el Sprint lo cubre HU-06.
- Cargar los Story Points desde HU-15 es una carga directa. HU-14 puede reemplazarlos más adelante con la estimación acordada en Planning Poker.
- Los Story Points se pueden borrar (dejar vacíos), porque la estimación es opcional hasta que se acuerde.
- La edición envía todos los campos editables juntos: lo que se deja vacío queda vacío. No hay edición parcial campo por campo.
- No hay paginación ni búsqueda: el volumen de historias de un proyecto académico es chico.
- Por la regla transversal del backlog, una historia sin Sprint está en "Nueva": al cerrar un Sprint, las no completadas vuelven a "Nueva" (fix #99, PR #115) y la base lo garantiza con `CHECK (estado = 'Nueva' OR sprint_id IS NOT NULL)`. El listado muestra el estado que tenga la historia, sin filtrar por él.
- La columna de criterios de aceptación ya existe en la migración 001 del modelo de datos (#107, D10), así que esta historia no agrega migraciones.
- Los criterios de aceptación se guardan como un solo texto. Si más adelante hace falta tratar cada criterio por separado, se migra a una entidad hija de la historia.
- Repriorizar el backlog (editar prioridad o valor de negocio) queda fuera de esta historia.
- El orden de creación se toma del identificador incremental de la historia (HU-01, FR-003).
