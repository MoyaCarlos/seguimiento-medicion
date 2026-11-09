# Feature Specification: Avance de las Historias del Proyecto (HU-18)

**Feature Branch**: `feature/HU-18-avance-historias`

**Created**: 2026-11-07

**Status**: Draft

**Input**: User description: "HU-18: Avance de las Historias del Proyecto (issue #177). Es la pantalla 'Estado' (/projects/{id}/status, hoy un placeholder): muestra todas las historias del proyecto agrupadas en Nueva, En progreso y Completada, con la cantidad de historias y de Story Points de cada grupo, el Sprint de cada historia (o 'Product Backlog' si no tiene), un mensaje con totales en cero si no hay historias y 'no encontrado' si el proyecto no existe. Es de solo lectura y no recalcula métricas. Criterios completos en la sección HU-18 de docs/PRODUCT-BACKLOG.md."

## Clarifications

### Session 2026-11-07 (decisiones del Scrum Master)

- Q: ¿Qué Sprint se muestra para una historia arrastrada al cerrar un Sprint? → A: Aparece en Nueva como "Product Backlog"; los Sprints anteriores por los que pasó quedan para HU-12.
- Q: ¿Qué se muestra para una historia Nueva asignada a un Sprint que todavía no se inició? → A: Aparece en Nueva con ese Sprint.
- Q: ¿Cómo se identifica el Sprint? → A: "Sprint <número>".
- Q: ¿Cómo se trata una historia sin Story Points? → A: Cuenta en la cantidad, suma 0 y se muestra "sin estimar" (nunca 0).
- Q: ¿Cómo se ordenan las historias dentro de cada grupo y hay filtros? → A: Orden estable por creación, sin filtros ni paginación.
- Q: ¿Los grupos vacíos se muestran? → A: Los tres grupos aparecen siempre, aunque estén vacíos.
- Q: ¿Qué responde el sistema ante un proyecto inexistente o un identificador inválido? → A: Un proyecto inexistente responde 404 "no encontrado". Un identificador con formato inválido NO es "no encontrado": responde 400 con campo `id` (convención de `parsearID` en `internal/http/project_handler.go:238` y de HU-12). Corrige el borrador, que trataba ambos casos como "no encontrado".

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Ver el avance de las historias por estado (Priority: P1)

Como integrante del equipo, abro la pantalla "Estado" del proyecto y veo todas sus historias agrupadas en Nueva, En progreso y Completada, con la cantidad de historias y de Story Points de cada grupo, para saber qué se terminó y qué falta sin revisar los Sprints uno por uno.

**Why this priority**: Es el Escenario 1 del backlog y la razón de ser de la historia: hoy no se puede responder "¿qué está terminado y qué no?" desde la pantalla.

**Independent Test**: Se siembra un proyecto con historias en los tres estados y con distintos Story Points; la consulta devuelve tres grupos con la cantidad y los Story Points esperados de cada uno.

**Acceptance Scenarios**:

1. **Given** un proyecto con historias en distintos estados y Sprints, **When** abro la pantalla "Estado", **Then** el sistema las muestra agrupadas en Nueva, En progreso y Completada.
2. **Given** un grupo con historias, **When** miro el grupo, **Then** veo la cantidad de historias y la suma de sus Story Points.
3. **Given** una historia sin Story Points, **When** miro su grupo, **Then** figura como "sin estimar" y suma cero a los Story Points del grupo.
4. **Given** un grupo sin historias, **When** miro la pantalla, **Then** el grupo se muestra igualmente, con cantidad cero y cero Story Points.

---

### User Story 2 - Ver el Sprint de cada historia (Priority: P1)

Como integrante del equipo, veo junto a cada historia el Sprint al que pertenece, o que está en el Product Backlog si no tiene Sprint, para ubicar cada historia sin abrir otras pantallas.

**Why this priority**: Es el Escenario 2 del backlog. Sin esto el agrupamiento por estado no dice dónde está el trabajo.

**Independent Test**: Se siembra una historia asignada a un Sprint y otra sin Sprint; la consulta informa el Sprint de la primera y "Product Backlog" para la segunda.

**Acceptance Scenarios**:

1. **Given** una historia que pertenece a un Sprint, **When** miro su grupo, **Then** veo el Sprint al que pertenece.
2. **Given** una historia sin Sprint, **When** miro su grupo, **Then** se indica que está en el Product Backlog.

---

### User Story 3 - Proyecto sin historias o inexistente (Priority: P2)

Como integrante del equipo, abro "Estado" de un proyecto sin historias y veo un mensaje de que no hay historias con los totales en cero; si el proyecto no existe, el sistema responde "no encontrado"; si el identificador no es válido, lo rechaza con un error sobre el campo `id`.

**Why this priority**: Son los Escenarios 3 y 4 del backlog; casos límite y de error de las dos historias anteriores.

**Independent Test**: Se consulta un proyecto recién creado (éxito, mensaje y ceros) y un identificador inexistente ("no encontrado", 404) y uno con formato inválido (400, campo `id`).

**Acceptance Scenarios**:

1. **Given** un proyecto sin historias, **When** abro la pantalla "Estado", **Then** el sistema muestra un mensaje de que no hay historias y los totales en cero.
2. **Given** que el proyecto no existe, **When** abro la pantalla "Estado", **Then** el sistema responde "no encontrado" (404).
3. **Given** un identificador de proyecto con formato inválido, **When** abro la pantalla "Estado", **Then** el sistema lo rechaza con un error de validación sobre el campo `id` (400), no con "no encontrado".

---

### Edge Cases

Redactados en formato EARS.

- El sistema **shall** resolver la pantalla únicamente con lectura y **shall** no modificar ningún dato.
- El sistema **shall** no recalcular métricas (HU-03): solo cuenta historias y suma los Story Points ya cargados.
- **WHEN** una historia de un Sprint finalizado quedó completada, **THEN** el sistema **shall** mostrarla en Completada con el Sprint en el que se completó.
- **WHEN** una historia no completada volvió al Product Backlog al cerrar un Sprint (HU-05), **THEN** el sistema **shall** mostrarla en Nueva y como "Product Backlog", sin historial de Sprints anteriores.
- **WHILE** una historia está en Nueva y asignada a un Sprint todavía no cerrado (HU-06), el sistema **shall** mostrarla en Nueva con ese Sprint.
- **IF** una historia no tiene Story Points, **THEN** el sistema **shall** mostrarla como "sin estimar" y **shall** sumar cero.
- **IF** el proyecto no existe, **THEN** el sistema **shall** responder "no encontrado" (404).
- **IF** el identificador del proyecto no es un entero positivo, **THEN** el sistema **shall** rechazar la consulta con un error de validación sobre el campo `id` (400), sin responder "no encontrado".
- **IF** falla la lectura de los datos, **THEN** el sistema **shall** responder con un error y **shall** no mostrar datos parciales.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir consultar todas las historias de un proyecto existente, agrupadas en Nueva, En progreso y Completada.
- **FR-002**: Cada grupo MUST informar la cantidad de historias y la suma de sus Story Points; los tres grupos MUST estar siempre presentes, aunque estén vacíos.
- **FR-003**: Cada historia MUST mostrar su título, sus Story Points (o "sin estimar") y el Sprint al que pertenece; si no tiene Sprint, MUST indicar "Product Backlog".
- **FR-004**: Una historia sin Story Points MUST aportar cero a la suma del grupo y MUST contarse en la cantidad.
- **FR-005**: Para un proyecto sin historias, el sistema MUST responder con éxito, con un mensaje de que no hay historias y los totales en cero.
- **FR-006**: El sistema MUST responder "no encontrado" (404) cuando el proyecto no existe, y MUST rechazar con un error de validación sobre el campo `id` (400) un identificador que no sea un entero positivo, antes de buscar el proyecto.
- **FR-007**: La consulta MUST ser de solo lectura y MUST no recalcular ni modificar métricas ni el Product Backlog de HU-15, que sigue mostrando solo las historias sin Sprint.
- **FR-008**: La pestaña "Estado" del proyecto (`/projects/{id}/status`) MUST reemplazar el texto "esta pantalla será implementada" por esta vista. El estado Planificado / En curso / Finalizado del panel del proyecto (HU-13) MUST quedar sin cambios.
- **FR-009**: La interfaz MUST mostrar un mensaje de error, sin datos parciales, cuando la consulta falla.

### Key Entities *(include if feature involves data)*

- **Historia**: la de HU-01, con estado (Nueva, En progreso, Completada), Story Points (opcionales) y Sprint (opcional). Ver `MODELO-DE-DATOS.md` §4.4.
- **Grupo de estado**: agrupamiento calculado al consultar (estado, cantidad, suma de Story Points); no se guarda.
- **Sprint**: el de HU-05, del que se muestra su identificación.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un integrante responde "¿qué se terminó y qué falta?" desde una sola pantalla, sin abrir Sprints individuales.
- **SC-002**: En un proyecto de referencia, el 100% de las historias aparece en exactamente un grupo, y la suma de las cantidades de los tres grupos es igual al total de historias del proyecto.
- **SC-003**: En un proyecto de referencia, la suma de Story Points de cada grupo coincide con la de sus historias.
- **SC-004**: El 100% de las consultas sobre proyectos sin historias responde con éxito, mensaje y totales en cero.
- **SC-005**: Consultar la pantalla deja sin cambios el 100% de los datos del proyecto.

## Assumptions

- **Sin cambios de esquema**: historia de solo lectura, sin tablas ni migraciones; la fuente es `historias` y `sprints` (`MODELO-DE-DATOS.md` §4.4). El `data-model.md` de la spec enlazará a esa sección.
- **Dependencias** (todas mergeadas en `main`): HU-01, HU-05, HU-06, HU-11 y HU-15.
- **Sprint actual de la historia**: se muestra el Sprint en el que está hoy (`sprint_id`). Una historia arrastrada no conserva el rastro de sus Sprints anteriores en esta vista (eso es HU-12).
- **Invariante del modelo**: En progreso y Completada solo existen dentro de un Sprint (§4.4), así que esos grupos nunca muestran "Product Backlog"; solo Nueva puede mostrarlo.
- **Identificación del Sprint**: "Sprint <número>", igual que en HU-12.
- **Sin filtros, orden configurable ni paginación**: proyecto académico con pocas historias (YAGNI). Dentro de cada grupo, orden estable por creación (decidido en Clarifications).
- **Sin control de acceso**: igual que el resto del sistema.
- **Fuera de alcance**: gráficos (HU-09), reporte PDF (HU-10), métricas (HU-03), cambiar el estado o el Sprint de una historia desde esta pantalla (HU-11, HU-06).
