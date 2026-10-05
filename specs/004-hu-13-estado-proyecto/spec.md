# Feature Specification: Fechas y Estado del Proyecto (HU-13)

**Feature Branch**: `feature/HU-13-fechas-estado-proyecto`

**Created**: 2026-10-04

**Status**: Implementada (US1 y US2), pendiente de revisión por PR.

**Input**: User description: "HU-13: Fechas y Estado del Proyecto. Como Scrum Master quiero registrar la fecha de inicio y finalización de un proyecto y consultar su estado general para poder tener visibilidad del ciclo de vida completo del proyecto. El registro de fechas ya lo resuelve HU-04 (editar proyecto), así que esta historia cubre solo la consulta de estado: dado un proyecto con Sprints y fechas cargadas, al consultar su estado el sistema muestra si está "Planificado", "En curso" o "Finalizado" según la fecha actual y el estado de sus Sprints. Se apoya en Project (HU-04) y en el SprintRepository (HU-05), ya en main."

## Clarifications

### Session 2026-10-04

- Q: Cuando el estado de los Sprints y las fechas del proyecto se contradicen (p.ej. un Sprint "Activo" con la fecha de fin ya transcurrida), ¿qué señal manda para el estado? → A: Los Sprints mandan y las fechas son respaldo: cualquier Sprint "Activo" ⇒ "En curso" aunque la fecha de fin haya pasado; si no hay "Activo" y todos los Sprints están "Finalizado" ⇒ "Finalizado"; en el resto deciden las fechas.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Consultar el estado general de un proyecto (Priority: P1)

Como Scrum Master, desde el panel de un proyecto consulto su estado general y el sistema me muestra una única etiqueta —"Planificado", "En curso" o "Finalizado"— que resume en qué punto del ciclo de vida está el proyecto. El estado se calcula a partir de la fecha actual, de las fechas de inicio/fin ya cargadas (HU-04) y de los estados de los Sprints del proyecto (HU-05), sin que yo tenga que interpretar los datos crudos.

**Why this priority**: Es el único valor de la historia: el enunciado pide explícitamente tener visibilidad del ciclo de vida completo del proyecto. Como el registro/edición de fechas ya está resuelto por HU-04, esta consulta es lo que aporta valor nuevo y es testeable de punta a punta por sí sola.

**Independent Test**: Puede probarse de forma independiente sobre un proyecto con fechas y Sprints ya cargados, consultando su estado y verificando que la etiqueta devuelta coincide con la regla de derivación para la fecha actual y el conjunto de Sprints dado.

**Acceptance Scenarios**:

1. **Given** un proyecto con fecha de inicio y fecha de fin cargadas y al menos un Sprint en estado "Activo", **When** consulto su estado, **Then** el sistema muestra "En curso".
2. **Given** un proyecto con al menos un Sprint y todos sus Sprints en estado "Finalizado", **When** consulto su estado, **Then** el sistema muestra "Finalizado".
3. **Given** un proyecto con fecha de inicio futura y sin ningún Sprint iniciado (sin Sprints "Activo" ni "Finalizado"), **When** consulto su estado, **Then** el sistema muestra "Planificado".
4. **Given** un proyecto sin Sprints iniciados y con la fecha actual dentro del rango [fecha_inicio, fecha_fin], **When** consulto su estado, **Then** el sistema muestra "En curso".
5. **Given** un identificador de proyecto que no existe, **When** consulto su estado, **Then** el sistema responde "proyecto no encontrado".

---

### User Story 2 - Ver el estado actualizado sin recálculo manual (Priority: P2)

Como Scrum Master, cada vez que abro o refresco el panel del proyecto el estado que veo refleja el momento de la consulta: si cambió la fecha, si se inició o cerró un Sprint, la etiqueta se actualiza sola, sin que nadie tenga que "recalcular" ni guardar nada.

**Why this priority**: La visibilidad del ciclo de vida pierde sentido si el estado queda congelado. Es un refinamiento del mismo valor de US1 y por eso es secundario, pero define una garantía de comportamiento observable (derivación on-demand).

**Independent Test**: Puede probarse de forma independiente consultando el estado de un proyecto, cambiando la fecha de referencia o el estado de sus Sprints y volviendo a consultar, verificando que la etiqueta cambia sin ninguna acción de guardado.

**Acceptance Scenarios**:

1. **Given** un proyecto consultado cuyo estado era "Planificado", **When** se inicia un Sprint del proyecto y vuelvo a consultar su estado, **Then** el sistema muestra "En curso" sin haber ejecutado ninguna acción de recálculo manual.
2. **Given** un proyecto consultado cuyo estado era "En curso", **When** se cierra su último Sprint y vuelvo a consultar su estado, **Then** el sistema muestra "Finalizado".

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- **WHEN** la fecha actual coincide exactamente con la fecha de inicio del proyecto, **THEN** el sistema **shall** considerar el proyecto como iniciado ("En curso", salvo evidencia de finalización).
- **WHEN** la fecha actual coincide exactamente con la fecha de fin del proyecto, **THEN** el sistema **shall** considerar el proyecto todavía "En curso" (los bordes se tratan de forma inclusiva; el proyecto finaliza recién al día siguiente).
- **IF** el proyecto no tiene ninguna fecha cargada y no tiene Sprints iniciados, **THEN** el sistema **shall** mostrar "Planificado".
- **IF** el proyecto solo tiene Sprints en estado "Pendiente" y su fecha de inicio es futura, **THEN** el sistema **shall** mostrar "Planificado" (un Sprint pendiente no constituye ejecución).
- **IF** el proyecto tiene un Sprint "Activo" pero la fecha de fin del proyecto ya pasó, **THEN** el sistema **shall** mostrar "En curso" (el Sprint "Activo" prevalece sobre la fecha vencida; el calendario solo es respaldo).
- **IF** el proyecto tiene fechas cargadas pero ningún Sprint, **THEN** el sistema **shall** derivar el estado únicamente a partir de la fecha actual y las fechas del proyecto.
- **IF** existe una fecha de fin sin fecha de inicio, **THEN** el sistema **shall** manejarlo como un proyecto sin inicio definido (no puede estar "En curso" por fecha; el estado dependerá de sus Sprints).
- **WHEN** se consulta el estado de un proyecto inexistente, **THEN** el sistema **shall** responder "proyecto no encontrado".
- **WHILE** se consulta el estado, **THEN** el sistema **shall** no modificar las fechas ni los Sprints del proyecto (operación de solo lectura).
- **WHEN** dos consultas del mismo proyecto se hacen en fechas distintas, **THEN** el sistema **shall** devolver el estado correspondiente a cada fecha de consulta (no hay estado persistido).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST calcular y mostrar el estado general de un proyecto al consultarlo, con un valor entre "Planificado", "En curso" y "Finalizado".
- **FR-002**: El estado MUST derivarse, en cada consulta, de la fecha actual, de la fecha de inicio y la fecha de fin del proyecto y de los estados de los Sprints del proyecto.
- **FR-003**: El sistema MUST mostrar "En curso" cuando el proyecto tenga al menos un Sprint en estado "Activo".
- **FR-004**: El sistema MUST mostrar "Finalizado" cuando el proyecto tenga al menos un Sprint y todos sus Sprints estén en estado "Finalizado".
- **FR-005**: Cuando el proyecto no tenga Sprints iniciados (ninguno "Activo" ni "Finalizado"), el sistema MUST derivar el estado a partir de las fechas: "Planificado" si la fecha actual es anterior a la fecha de inicio; "En curso" si está dentro del rango de fechas; "Finalizado" si es posterior a la fecha de fin.
- **FR-006**: El sistema MUST tratar los límites del rango de fechas de forma inclusiva: en la fecha de inicio y en la fecha de fin el proyecto está "En curso".
- **FR-007**: El sistema MUST mostrar "Planificado" cuando el proyecto no tenga ninguna ejecución registrada: sin Sprints iniciados y sin fecha de inicio alcanzada (incluidos los proyectos sin fechas cargadas).
- **FR-008**: El sistema MUST manejar proyectos sin fecha de inicio como no iniciados por fecha: no pueden pasar a "En curso" por fecha, salvo por la existencia de un Sprint "Activo".
- **FR-009**: El sistema MUST permitir consultar el estado de un proyecto por su identificador.
- **FR-010**: Cuando se consulte el estado de un proyecto que no existe, el sistema MUST responder "proyecto no encontrado".
- **FR-011**: El sistema MUST recalcular el estado en cada consulta (on-demand), de modo que refleje la fecha actual y los Sprints vigentes sin persistir el estado ni requerir un recálculo manual.
- **FR-012**: La consulta de estado MUST ser de solo lectura: MUST NOT modificar las fechas del proyecto ni el estado de sus Sprints.
- **FR-013**: El registro y la edición de las fechas de inicio y fin del proyecto quedan cubiertos por HU-04 y MUST NOT formar parte del alcance de esta historia.
- **FR-014**: Cuando la información de fechas y la de Sprints entren en conflicto (por ejemplo, un Sprint "Activo" con la fecha de fin del proyecto ya transcurrida), el sistema MUST dar prioridad a los Sprints sobre las fechas: un Sprint "Activo" implica "En curso" aunque la fecha de fin haya pasado, y las fechas del proyecto solo deciden cuando no hay Sprints iniciados.
- **FR-015**: El sistema MUST devolver el estado usando las etiquetas visibles "Planificado", "En curso" y "Finalizado", con un valor canónico interno estable asociado a cada una.

### Key Entities *(include if feature involves data)*

- **Proyecto** (ya existente, HU-04): espacio de trabajo con identificador único, nombre, descripción y fechas de inicio y fin opcionales. Esta historia consume sus fechas como insumo del cálculo y no las modifica.
- **Sprint** (ya existente, HU-05): iteración de trabajo con estado (`Pendiente` / `Activo` / `Finalizado`) y fechas. Esta historia consume el estado de los Sprints de un proyecto como insumo del cálculo y no los modifica.
- **Estado del Proyecto** (derivado): uno de "Planificado", "En curso" o "Finalizado". No es un dato persistido en el proyecto: se calcula on-demand a partir de la fecha actual, las fechas del proyecto y los estados de sus Sprints.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un Scrum Master puede conocer el estado general de un proyecto con una sola consulta, sin pasos intermedios ni acciones de recálculo.
- **SC-002**: El 100% de las consultas devuelven uno de los tres valores válidos y de forma determinista: ante la misma fecha y el mismo conjunto de Sprints, el resultado es siempre el mismo.
- **SC-003**: El 100% de las consultas reflejan el momento de la fecha de consulta: al iniciar o cerrar un Sprint, el estado cambia en la siguiente consulta sin intervención manual.
- **SC-004**: El 100% de las consultas sobre proyectos inexistentes se rechazan con "proyecto no encontrado".
- **SC-005**: La consulta de estado produce cero modificaciones sobre las fechas del proyecto y sobre los estados de sus Sprints (operación de solo lectura).

## Assumptions

- La "fecha actual" es la fecha del servidor al momento de la consulta; el cálculo usa la zona horaria del servidor y no contempla husos horarios del cliente.
- El estado del proyecto es un valor derivado y no persistido: no se almacena en el proyecto ni se congela; se recalcula en cada consulta (coherente con el principio YAGNI y con la decisión de no construir Observer/tiempo real).
- El registro y la edición de las fechas de inicio y fin del proyecto ya están cubiertos por HU-04 (fechas opcionales, con validación de que la fin no sea anterior a la inicio); esta historia no vuelve a especificarlos.
- Las fechas del proyecto son opcionales (HU-04); esta historia define el comportamiento para proyectos sin fechas y sin Sprints iniciados: "Planificado".
- Los bordes del rango de fechas se tratan de forma inclusiva (inicio y fin cuentan como "En curso"), consistente con el criterio laxo de HU-04, que permite fechas de proyecto iguales.
- Un Sprint en estado "Activo" es evidencia de que el proyecto está "En curso"; un Sprint "Pendiente" no es evidencia de ejecución.
- Precedencia acordada (2026-10-04): los Sprints mandan sobre las fechas. Un Sprint "Activo" mantiene el proyecto "En curso" aunque la fecha de fin haya pasado; las fechas del proyecto solo deciden cuando no hay Sprints iniciados (ninguno "Activo" ni "Finalizado").
- Se consideran iniciados los Sprints en estado "Activo" o "Finalizado"; "Pendiente" no cuenta como inicio.
- Se asume la regla de HU-05 de a lo sumo un Sprint "Activo" y a lo sumo un Sprint "Pendiente" por proyecto a la vez.
- El cálculo vive en la capa de dominio/servicio y se apoya en el `Project` de HU-04 y en el `SprintRepository` de HU-05, ambos ya en `main`; el detalle de implementación se define en `/speckit.plan`.
- La presentación del estado en la interfaz (ubicación exacta en el panel del proyecto) queda fuera del alcance de esta historia; esta historia garantiza el valor devuelto por la consulta.
- La consulta de estado no requiere autenticación adicional más allá de la ya asumida como existente para operar el sistema.
