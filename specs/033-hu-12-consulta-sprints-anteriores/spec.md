# Feature Specification: Consulta de Sprints Anteriores (HU-12)

**Feature Branch**: `feature/HU-12-consulta-sprints-anteriores`

**Created**: 2026-11-03

**Status**: Draft

**Input**: User description: "HU-12: Consulta de Sprints Anteriores (issue #12, sub-issue #84). Como Scrum Master quiero consultar el historial de Sprints ya cerrados con su resumen (Sprint Goal, historias completadas, métricas) para poder comparar el desempeño del equipo a lo largo del proyecto. Escenario 1: dado que el proyecto tiene al menos dos Sprints finalizados, cuando ingreso a la sección 'Historial de Sprints', el sistema muestra una lista con cada Sprint cerrado, su Sprint Goal, fechas y métricas principales. Escenario 2: dado que estoy en el historial, cuando selecciono un Sprint, el sistema muestra el detalle completo de ese Sprint (historias, esfuerzo, defectos) sin afectar al Sprint activo actual. Notas: las historias que se muestran como comprometidas son las asignadas al cerrar el Sprint, completadas o no (ver HU-05). Depende de HU-05, HU-03 (métricas, spec 032 todavía en borrador, sin mergear), HU-07, HU-08 y HU-11. Fuentes: docs/PRODUCT-BACKLOG.md (HU-12) y MODELO-DE-DATOS.md (§4.6 historias_arrastradas; HU-12 es solo lectura)."

## Clarifications

### Session 2026-11-03

- Q: ¿Qué métricas principales muestra la lista del historial? → A: Story Points planificados y completados, porcentaje de historias completadas, horas estimadas y reales, desviación y defectos detectados y resueltos.
- Q: ¿Con qué nivel de detalle se muestra el esfuerzo de un Sprint? → A: Total del Sprint y total por historia, sin cada registro de horas.
- Q: ¿Qué debe responder el sistema si se pide, en el historial, el detalle de un Sprint que todavía está "Activo" o "Pendiente"? → A: Rechaza la consulta con un mensaje que indica que el Sprint todavía no está finalizado.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Ver la lista de Sprints cerrados (Priority: P1)

Como Scrum Master, consulto el historial de un proyecto y veo cada Sprint cerrado con su Sprint Goal, sus fechas y sus métricas principales, para comparar el desempeño del equipo a lo largo del proyecto.

**Why this priority**: Es el Escenario 1 del backlog y un requisito explícito del enunciado (punto 3). Permite comparar Sprints sin abrirlos uno por uno.

**Independent Test**: Se prueba sembrando un proyecto con dos Sprints finalizados, uno "Activo" y uno "Pendiente"; la consulta devuelve solo los dos finalizados, con sus datos y con las mismas métricas que informa HU-03.

**Acceptance Scenarios**:

1. **Given** un proyecto con dos Sprints finalizados, **When** consulto el historial, **Then** el sistema lista ambos, ordenados por su número dentro del proyecto, cada uno con su número, Sprint Goal, fecha de inicio, fecha de fin y métricas principales.
2. **Given** un proyecto con un Sprint finalizado, uno "Activo" y uno "Pendiente", **When** consulto el historial, **Then** solo aparece el Sprint finalizado.
3. **Given** un Sprint finalizado, **When** consulto el historial, **Then** cada Sprint informa como métricas principales los Story Points planificados y completados, el porcentaje de historias completadas, las horas estimadas y reales, la desviación y los defectos detectados y resueltos.
4. **Given** un Sprint finalizado, **When** comparo sus métricas en el historial con las que informa HU-03 para ese Sprint, **Then** los valores son idénticos.
5. **Given** un proyecto sin horas por Story Point, **When** consulto el historial, **Then** las horas estimadas y la desviación que se muestren aparecen como "sin datos", nunca como cero.
6. **Given** que el proyecto no existe, **When** consulto su historial, **Then** el sistema responde "proyecto no encontrado".

---

### User Story 2 - Ver el detalle de un Sprint cerrado (Priority: P1)

Como Scrum Master, selecciono un Sprint del historial y veo su detalle completo (historias comprometidas, esfuerzo y defectos), para entender qué pasó en ese Sprint, sin afectar al Sprint activo.

**Why this priority**: Es el Escenario 2 del backlog. Explica los números de la lista: qué se comprometió, qué se completó, cuánto esfuerzo llevó y qué defectos aparecieron.

**Independent Test**: Se prueba sembrando un Sprint finalizado con una historia completada, una arrastrada, horas imputadas y defectos, además de un Sprint "Activo"; la consulta devuelve el detalle esperado y, al repetir la consulta del Sprint "Activo", sus datos no cambiaron.

**Acceptance Scenarios**:

1. **Given** un Sprint finalizado, **When** consulto su detalle, **Then** el sistema muestra su número, Sprint Goal, fechas, todas sus métricas de HU-03, sus historias comprometidas, su esfuerzo y sus defectos.
2. **Given** un Sprint con una historia completada de 5 Story Points y una arrastrada de 3, **When** consulto su detalle, **Then** ambas aparecen como comprometidas, la primera como completada con 5 Story Points y la segunda como no completada con 3.
3. **Given** una historia arrastrada que después se editó en el Product Backlog y se completó en otro Sprint, **When** consulto el detalle del primer Sprint, **Then** sigue figurando como no completada y con los Story Points que tenía al cerrarlo; en el detalle del segundo Sprint figura como completada.
4. **Given** un Sprint con horas imputadas, **When** consulto su detalle, **Then** el sistema muestra el esfuerzo como total del Sprint y total por historia, y el total del Sprint coincide con las horas reales de HU-03.
5. **Given** un defecto detectado en un Sprint y resuelto en otro, **When** consulto el detalle de cada uno, **Then** aparece en ambos, identificado como detectado en el primero y como resuelto en el segundo.
6. **Given** un defecto detectado y resuelto en el mismo Sprint, **When** consulto su detalle, **Then** el defecto aparece una sola vez, identificado como detectado y resuelto en ese Sprint.
7. **Given** un proyecto con un Sprint "Activo", **When** consulto el detalle de un Sprint finalizado, **Then** el Sprint "Activo", sus historias, sus horas y sus defectos no cambian.
8. **Given** que el Sprint no existe, **When** consulto su detalle, **Then** el sistema responde "Sprint no encontrado".
9. **Given** un Sprint "Activo" o "Pendiente", **When** consulto su detalle en el historial, **Then** el sistema rechaza la consulta indicando que el Sprint todavía no está finalizado.

---

### User Story 3 - Consultar un proyecto sin Sprints cerrados (Priority: P2)

Como Scrum Master, consulto el historial de un proyecto que todavía no cerró ningún Sprint y recibo una lista vacía, en vez de un error, para entender que todavía no hay nada que comparar.

**Why this priority**: Es el estado inicial de todo proyecto, pero no aporta valor sin las dos historias anteriores.

**Independent Test**: Se prueba consultando un proyecto recién creado y otro con un único Sprint "Activo"; ambos responden con éxito y sin Sprints.

**Acceptance Scenarios**:

1. **Given** un proyecto sin Sprints finalizados, **When** consulto el historial, **Then** el sistema responde con éxito y una lista vacía.
2. **Given** un Sprint finalizado sin historias, horas ni defectos, **When** consulto su detalle, **Then** el sistema responde con éxito, con listas vacías y con las métricas que informa HU-03 para ese caso.

---

### User Story 4 - Navegar el historial desde la interfaz (Priority: P3)

Como Scrum Master, entro a la sección "Historial de Sprints" del proyecto, veo la lista y selecciono un Sprint para ver su detalle.

**Why this priority**: Es la sub-issue de Frontend (#88), que por la Definition of Done no frena el cierre de la historia.

**Independent Test**: Se prueba con pruebas de componente y con la aplicación levantada: se abre la sección, se ve la lista, se selecciona un Sprint y se ve su detalle; un proyecto sin Sprints cerrados muestra un estado vacío informativo.

**Acceptance Scenarios**:

1. **Given** un proyecto con Sprints cerrados, **When** ingreso a "Historial de Sprints", **Then** veo la lista con cada Sprint identificado como "Sprint <número>".
2. **Given** que estoy en el historial, **When** selecciono un Sprint, **Then** veo su detalle y puedo volver a la lista.
3. **Given** un proyecto sin Sprints cerrados, **When** ingreso al historial, **Then** veo un mensaje que indica que todavía no hay Sprints cerrados.
4. **Given** una historia, una estimación o una métrica sin datos, **When** se muestra en el historial, **Then** aparece como "sin estimar" o "sin datos", nunca como cero.
5. **Given** que la consulta falla, **When** recibo la respuesta, **Then** veo un mensaje de error y ningún dato parcial.

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** resolver el historial únicamente con lectura y **shall** no modificar ningún dato.
- El sistema **shall** incluir en el historial únicamente los Sprints en estado "Finalizado".
- El sistema **shall** informar los mismos valores de métricas que HU-03, sin recalcularlos con otra regla.
- **WHILE** un Sprint está "Pendiente" o "Activo", el sistema **shall** excluirlo de la lista y **shall** rechazar la consulta de su detalle en el historial.
- **WHEN** una historia se arrastra a lo largo de varios Sprints, **THEN** el sistema **shall** mostrarla como no completada en cada Sprint que la arrastró, con los Story Points que tenía al cerrar cada uno, y como completada solo en el Sprint donde se completó.
- **IF** una historia comprometida no tiene Story Points, **THEN** el sistema **shall** mostrarla como "sin estimar" y **shall** sumar cero a los Story Points del Sprint.
- **WHEN** se imputan horas a una historia en un Sprint, **THEN** esas horas **shall** mostrarse en el detalle de ese Sprint, aunque la historia haya pasado después a otro Sprint o haya vuelto al Product Backlog.
- **IF** una historia comprometida no tiene horas imputadas en el Sprint, **THEN** el sistema **shall** mostrar cero horas para esa historia.
- **WHEN** una historia acumula horas en varios Sprints, **THEN** el detalle de cada Sprint **shall** mostrar solo las horas imputadas en ese Sprint.
- **WHILE** un defecto detectado en un Sprint finalizado sigue pendiente, el sistema **shall** mostrarlo como pendiente en el detalle de ese Sprint.
- **IF** el proyecto no existe, **THEN** el sistema **shall** responder "proyecto no encontrado".
- **IF** el Sprint no existe o pertenece a otro proyecto que el consultado, **THEN** el sistema **shall** responder "Sprint no encontrado".
- **IF** el identificador del proyecto o del Sprint no es válido, **THEN** el sistema **shall** rechazar la consulta con un error sobre ese campo.
- **IF** falla la lectura de los datos, **THEN** el sistema **shall** responder con un error y **shall** no devolver datos parciales.

## Requirements *(mandatory)*

### Functional Requirements

**Lista del historial**

- **FR-001**: El sistema MUST permitir consultar los Sprints finalizados de un proyecto existente, ordenados de forma ascendente por su número dentro del proyecto.
- **FR-002**: Cada Sprint de la lista MUST incluir su número, Sprint Goal, fecha de inicio, fecha de fin y métricas principales.
- **FR-003**: Las métricas principales MUST ser los Story Points planificados y completados, el porcentaje de historias completadas, las horas estimadas y reales, la desviación y los defectos detectados y resueltos.
- **FR-004**: Para un proyecto sin Sprints finalizados, el sistema MUST responder con éxito y una lista vacía.

**Detalle de un Sprint**

- **FR-005**: El sistema MUST permitir consultar el detalle de un Sprint finalizado: número, Sprint Goal, fechas, todas sus métricas de HU-03, historias comprometidas, esfuerzo y defectos.
- **FR-006**: Las historias comprometidas MUST ser las asignadas al cerrar el Sprint: las completadas en él y las arrastradas al cerrarlo (D2c).
- **FR-007**: Cada historia comprometida MUST mostrar su título, sus Story Points y si se completó o no en ese Sprint. Una historia arrastrada MUST mostrar los Story Points que tenía al cerrarlo, aunque se hayan editado después.
- **FR-008**: El esfuerzo MUST incluir las horas imputadas a ese Sprint, como total del Sprint y total por historia. Las horas MUST calcularse a partir de los minutos imputados (D8) e informarse con dos decimales, como en HU-03.
- **FR-009**: Los defectos MUST incluir los detectados y los resueltos en ese Sprint, cada uno una sola vez, con su descripción, severidad, historia relacionada, estado actual e indicación de si se detectó, se resolvió o ambas en ese Sprint.
- **FR-010**: La cantidad de defectos detectados y resueltos y el total de horas del detalle MUST coincidir con las métricas de HU-03 para ese Sprint.
- **FR-011**: El sistema MUST rechazar la consulta del detalle de un Sprint "Pendiente" o "Activo" indicando que todavía no está finalizado, y MUST responder "Sprint no encontrado" si no existe o pertenece a otro proyecto.

**Reglas comunes**

- **FR-012**: Las métricas MUST provenir del cálculo único de HU-03; esta historia no define ni recalcula métricas.
- **FR-013**: Los valores que HU-03 informa como "sin datos" MUST mantenerse como "sin datos" en el historial.
- **FR-014**: Las consultas del historial MUST no modificar ningún dato: ni del Sprint consultado ni del Sprint activo.
- **FR-015**: El sistema MUST responder "proyecto no encontrado" cuando el proyecto no existe.

**Interfaz**

- **FR-016**: El proyecto MUST ofrecer una sección "Historial de Sprints" con la lista y el acceso al detalle de cada Sprint. Cada Sprint MUST identificarse como "Sprint <número>".
- **FR-017**: La sección MUST mostrar un estado vacío informativo cuando no hay Sprints cerrados, y un mensaje de error, sin datos parciales, cuando la consulta falla.

### Key Entities *(include if feature involves data)*

- **Sprint cerrado**: un Sprint de HU-05 en estado "Finalizado". Siempre tiene número, Sprint Goal y fechas.
- **Historia comprometida**: una historia completada en el Sprint o arrastrada al cerrarlo (D2c). La arrastrada conserva los Story Points que tenía en ese momento.
- **Métricas del Sprint**: las que calcula HU-03; no se guardan.
- **Registro de horas**: el de HU-07. Se atribuye al Sprint que estaba activo al cargarlo.
- **Defecto**: el de HU-08. Aparece en su Sprint de detección y, si está resuelto, en su Sprint de resolución.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un Scrum Master puede comparar los Sprints cerrados de un proyecto desde una sola pantalla, sin abrirlos uno por uno.
- **SC-002**: Un Scrum Master llega del historial al detalle de cualquier Sprint cerrado con una sola selección.
- **SC-003**: En un proyecto de referencia, el 100% de las métricas del historial coincide con las que informa HU-03.
- **SC-004**: En un proyecto de referencia, el 100% de las historias comprometidas de cada Sprint cerrado aparece en su detalle, con el resultado y los Story Points que tenía al cerrarlo.
- **SC-005**: Consultar el historial deja sin cambios el 100% de los datos del proyecto, incluido el Sprint activo.
- **SC-006**: El 100% de las consultas sobre proyectos sin Sprints finalizados responde con éxito, sin errores.

## Assumptions

- **Dependencia de HU-03**: HU-12 muestra las métricas de HU-03 (spec 032), ya mergeada en `main`.
- **Sin cambios de esquema**: la historia es de solo lectura y no agrega tablas ni migraciones (`MODELO-DE-DATOS.md` §8). Las fuentes son `sprints`, `historias`, `historias_arrastradas`, `registros_horas` y `defectos` (§4.5, §4.6, §4.8 y §4.9).
- **Fechas del Sprint**: se muestran la fecha de inicio y la de fin definidas al iniciarlo. El modelo no guarda la fecha real de cierre.
- **Orden de la lista**: ascendente por número, igual que el detalle por Sprint de HU-03 y más natural para ver la evolución del equipo.
- **Título y estado actuales**: el modelo solo guarda una foto de los Story Points de las arrastradas. El título de la historia y el estado del defecto se muestran con su valor actual, así que un defecto de un Sprint cerrado puede figurar como resuelto si se resolvió después.
- **Reproducibilidad (D2b)**: los Story Points ya cargados no se editan con Sprint "Activo" o "Finalizado", pero cargarlos por primera vez sí se permite. Una historia completada "sin estimar" puede pasar a tener Story Points, igual que en HU-03.
- **Sin escritura**: el estado de las historias de un Sprint finalizado no cambia (HU-11), así que el resultado de cada historia comprometida es estable.
- **Horas siempre de historias comprometidas**: solo se imputan horas a historias del Sprint "Activo" (HU-07), y una historia asignada solo sale de su Sprint al cerrarlo, cuando queda como arrastrada (HU-05, HU-06). Por eso el total por historia suma exactamente las horas reales del Sprint. Si otra historia permite sacar historias de un Sprint activo (D5), hay que revisar este punto.
- **Sin filtros ni paginación**: un proyecto académico tiene pocos Sprints (YAGNI).
- **Sin control de acceso**: igual que en el resto del sistema, no hay autenticación ni permisos por rol.
- **Fuera de alcance**: gráficos y comparaciones visuales (HU-09), reporte en PDF (HU-10), edición de datos de un Sprint cerrado y cada registro de horas por separado, que sigue disponible por historia en HU-07.
- **Dependencias**: HU-03, HU-05, HU-07, HU-08 y HU-11, con el backend ya mergeado en `main`.
