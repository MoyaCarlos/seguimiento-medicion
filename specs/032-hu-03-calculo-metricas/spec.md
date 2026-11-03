# Feature Specification: Cálculo de Métricas del Proyecto (HU-03)

**Feature Branch**: `feature/HU-03-calculo-metricas`

**Created**: 2026-11-03

**Status**: Draft

**Input**: User description: "HU-03: Cálculo de Métricas del Proyecto (issue #3, sub-issue #49). Como Scrum Master quiero que el sistema calcule automáticamente las métricas del proyecto y de cada Sprint para poder analizar la capacidad real del equipo y la calidad del trabajo entregado. Escenario 1: tras cerrar un Sprint (HU-05) con historias completadas (HU-11), el sistema obtiene Story Points planificados y completados, velocidad del equipo, horas estimadas (Story Points × horas por Story Point del proyecto) y reales, desviación entre esfuerzo estimado y real, porcentaje de historias completadas, y cantidad de defectos detectados y resueltos en ese Sprint. Escenario 2: un proyecto sin Sprints finalizados devuelve todos los valores en cero (o 'sin datos') en vez de fallar. Escenario 3: al editar un proyecto existente se define cuántas horas equivale un Story Point; el sistema lo guarda y lo usa para las horas estimadas; si no está definido, las horas estimadas y la desviación se informan como 'sin datos'. Fuentes: docs/PRODUCT-BACKLOG.md (HU-03) y MODELO-DE-DATOS.md §6 y §8 (incluye que PUT /projects/{id} todavía no acepta horas_por_story_point y esa tarea es de esta historia; qué SP entran en las horas estimadas lo define esta spec). Decisiones vigentes: D2b, D2c, D7, D8."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Consultar las métricas de cada Sprint cerrado y del proyecto (Priority: P1)

Como Scrum Master, consulto las métricas de un proyecto y veo, para cada Sprint finalizado y para el proyecto en conjunto, los Story Points planificados y completados, la velocidad, las horas estimadas y reales, la desviación, el porcentaje de historias completadas y los defectos detectados y resueltos, para analizar la capacidad real del equipo y la calidad del trabajo entregado.

**Why this priority**: Es el núcleo de la historia y un requisito central del enunciado (punto 7). El dashboard (HU-09), el reporte (HU-10) y el historial de Sprints (HU-12) muestran estos valores y no los recalculan.

**Independent Test**: Se prueba sembrando un proyecto con horas por Story Point definidas y un Sprint finalizado que tiene historias completadas, historias arrastradas, horas imputadas y defectos; la consulta devuelve cada métrica con el valor calculado a mano. Un segundo Sprint finalizado permite verificar los valores del proyecto.

**Acceptance Scenarios**:

1. **Given** un Sprint finalizado con una historia completada de 5 Story Points y una historia arrastrada de 3, **When** consulto las métricas, **Then** el Sprint informa 8 Story Points planificados, 5 completados, velocidad 5 y 50% de historias completadas.
2. **Given** ese Sprint, 4 horas por Story Point en el proyecto y 30 horas imputadas al Sprint, **When** consulto las métricas, **Then** el Sprint informa [NEEDS CLARIFICATION: ¿las horas estimadas salen de los Story Points planificados (32 horas, desviación −2) o de los completados (20 horas, desviación +10)?] y 30 horas reales.
3. **Given** un Sprint finalizado con dos defectos detectados en él y uno resuelto en él, **When** consulto las métricas, **Then** el Sprint informa 2 defectos detectados y 1 resuelto.
4. **Given** un defecto detectado en un Sprint y resuelto en otro, ambos finalizados, **When** consulto las métricas, **Then** cuenta como detectado en el primero y como resuelto en el segundo.
5. **Given** dos Sprints finalizados con 5 y 8 Story Points completados, **When** consulto las métricas, **Then** el proyecto informa 13 Story Points completados en total y una velocidad de 6,5 (promedio por Sprint finalizado).
6. **Given** un proyecto con un Sprint finalizado y uno "Activo", **When** consulto las métricas, **Then** solo el Sprint finalizado aparece y aporta a los valores del proyecto; las horas y los defectos del Sprint "Activo" no se cuentan.
7. **Given** una historia arrastrada que después se edita en el Product Backlog y se completa en otro Sprint, **When** consulto las métricas, **Then** en el primer Sprint sigue contando con los Story Points que tenía al cerrarlo y como no completada, y en el segundo cuenta como completada.
8. **Given** que el proyecto no existe, **When** consulto sus métricas, **Then** el sistema responde "proyecto no encontrado".

---

### User Story 2 - Consultar un proyecto sin Sprints cerrados (Priority: P2)

Como Scrum Master, consulto las métricas de un proyecto que todavía no cerró ningún Sprint y recibo valores en cero, en vez de un error, para que las pantallas que las muestran puedan presentar un estado vacío.

**Why this priority**: Es el estado inicial de todo proyecto; el dashboard (HU-09, Escenario 2) depende de que la consulta no falle.

**Independent Test**: Se prueba consultando un proyecto recién creado y otro con un único Sprint "Activo" con horas y defectos; ambos devuelven una respuesta exitosa, sin Sprints y con los valores del proyecto en cero.

**Acceptance Scenarios**:

1. **Given** un proyecto sin Sprints finalizados y con horas por Story Point definidas, **When** consulto las métricas, **Then** el sistema responde con éxito, sin Sprints, y con todos los valores del proyecto en cero.
2. **Given** un proyecto sin Sprints finalizados y sin horas por Story Point, **When** consulto las métricas, **Then** el sistema responde con éxito, con los demás valores en cero y las horas estimadas y la desviación como "sin datos".

---

### User Story 3 - Definir las horas por Story Point del proyecto (Priority: P2)

Como Scrum Master, al editar un proyecto defino cuántas horas equivale un Story Point, para que el sistema pueda calcular las horas estimadas y compararlas con las reales.

**Why this priority**: Sin este parámetro no hay horas estimadas ni desviación (decisión D7). El resto de las métricas se calcula sin él, por eso es independiente de P1.

**Independent Test**: Se prueba editando un proyecto para definir el parámetro, consultándolo de nuevo para verificar que quedó guardado y consultando las métricas de un Sprint finalizado antes y después de definirlo.

**Acceptance Scenarios**:

1. **Given** un proyecto existente, **When** lo edito y defino 4 horas por Story Point, **Then** el sistema guarda el valor, lo devuelve al consultar el proyecto y lo usa para calcular las horas estimadas.
2. **Given** un proyecto con un Sprint finalizado y sin horas por Story Point, **When** consulto las métricas, **Then** las horas estimadas y la desviación se informan como "sin datos", y las horas reales y las demás métricas se informan con su valor.
3. **Given** un proyecto con horas por Story Point, **When** lo edito enviando solo los datos que ya existían antes de esta historia, sin mencionar el parámetro, **Then** el parámetro conserva su valor.
4. **Given** un proyecto con horas por Story Point, **When** lo edito y dejo el parámetro explícitamente vacío, **Then** el parámetro queda sin definir y las horas estimadas y la desviación vuelven a informarse como "sin datos".
5. **Given** que estoy editando un proyecto, **When** ingreso como horas por Story Point un valor cero, negativo o no numérico, **Then** el sistema rechaza la edición con un error sobre ese campo y el proyecto queda como estaba.
6. **Given** un proyecto con Sprints finalizados y 4 horas por Story Point, **When** cambio el parámetro a 5, **Then** [NEEDS CLARIFICATION: ¿las horas estimadas de los Sprints ya cerrados se recalculan con 5, o conservan las calculadas con 4?].

---

### User Story 4 - Definir el parámetro desde la pantalla de edición del proyecto (Priority: P3)

Como Scrum Master, en la pantalla de edición del proyecto veo y modifico el campo "horas por Story Point"; si el sistema rechaza el valor, veo el mensaje de error que devuelve.

**Why this priority**: Es la sub-issue de Frontend (#53), que por la Definition of Done no frena el cierre de la historia. La visualización de las métricas es de HU-09, no de esta historia.

**Independent Test**: Se prueba con la aplicación levantada: se abre la edición de un proyecto, se carga el valor, se guarda y al volver a abrir la pantalla el valor sigue ahí; se ingresa un valor inválido y se ve el mensaje del sistema.

**Acceptance Scenarios**:

1. **Given** un proyecto con horas por Story Point, **When** abro la pantalla de edición, **Then** el campo muestra el valor actual.
2. **Given** que estoy en la pantalla de edición, **When** cargo un valor válido y guardo, **Then** el valor queda guardado.
3. **Given** que el sistema rechaza el valor, **When** recibo la respuesta, **Then** la pantalla muestra el mensaje de error junto al campo.

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** calcular las métricas al consultarlas y **shall** no guardarlas (`MODELO-DE-DATOS.md` §6).
- El sistema **shall** considerar únicamente los Sprints en estado "Finalizado", tanto en el detalle por Sprint como en los valores del proyecto.
- **WHILE** un Sprint está "Pendiente" o "Activo", el sistema **shall** excluirlo del cálculo, junto con sus horas y sus defectos.
- **IF** una historia completada o arrastrada no tiene Story Points, **THEN** el sistema **shall** contarla como historia comprometida en el porcentaje y **shall** sumar cero a los Story Points.
- **IF** un Sprint finalizado no tiene historias comprometidas, **THEN** el sistema **shall** informar cero Story Points planificados y completados, velocidad cero y 0% de historias completadas, sin fallar.
- **IF** el proyecto no tiene horas por Story Point, **THEN** el sistema **shall** informar las horas estimadas y la desviación como "sin datos", y nunca como cero.
- **WHERE** las horas por Story Point están definidas y los Story Points de la base de cálculo suman cero, el sistema **shall** informar cero horas estimadas, la desviación en horas igual a las horas reales y la desviación porcentual como "sin datos".
- **IF** un Sprint finalizado no tiene horas imputadas, **THEN** el sistema **shall** informar cero horas reales.
- **WHEN** una historia se arrastra a lo largo de varios Sprints, **THEN** el sistema **shall** contarla como comprometida en cada Sprint que la arrastró, con los Story Points que tenía al cerrar cada uno, y como completada solo en el Sprint donde se completó.
- **WHEN** se imputan horas a una historia que después vuelve al Product Backlog o pasa a otro Sprint, **THEN** esas horas **shall** seguir contando en el Sprint al que se imputaron.
- **WHEN** se consultan las métricas de un proyecto que no existe, **THEN** el sistema **shall** responder "proyecto no encontrado".
- **IF** las horas por Story Point enviadas no son un número mayor que cero, **THEN** el sistema **shall** rechazar la edición con un error sobre ese campo y **shall** no modificar el proyecto.
- **IF** una edición incumple varias reglas a la vez, **THEN** el sistema **shall** responder con la primera según el orden que ya rige en la edición del proyecto (HU-04), y las horas por Story Point **shall** validarse al final.
- **WHEN** se edita un proyecto sin mencionar las horas por Story Point, **THEN** el sistema **shall** conservar el valor guardado.

## Requirements *(mandatory)*

### Functional Requirements

**Consulta de métricas**

- **FR-001**: El sistema MUST permitir consultar las métricas de un proyecto existente indicando solo el proyecto. La respuesta MUST incluir el detalle de cada Sprint finalizado, ordenado por su número dentro del proyecto, y los valores del proyecto en conjunto.
- **FR-002**: Los Story Points planificados de un Sprint MUST ser la suma de los Story Points de sus historias comprometidas: las completadas en él y las arrastradas al cerrarlo, estas últimas con los Story Points que tenían en ese momento (D2c).
- **FR-003**: Los Story Points completados de un Sprint MUST ser la suma de los Story Points de sus historias en estado "Completada".
- **FR-004**: La velocidad de un Sprint MUST ser igual a sus Story Points completados. La velocidad del proyecto MUST ser el promedio de los Story Points completados por Sprint finalizado.
- **FR-005**: El porcentaje de historias completadas MUST calcularse como historias completadas sobre historias comprometidas, por cantidad de historias y no por Story Points.
- **FR-006**: Las horas estimadas de un Sprint MUST calcularse como [NEEDS CLARIFICATION: Story Points planificados o completados, ver User Story 1] × horas por Story Point del proyecto (D7).
- **FR-007**: Las horas reales de un Sprint MUST ser la suma de las horas imputadas a ese Sprint, sin importar dónde esté hoy la historia.
- **FR-008**: La desviación MUST informarse en horas (reales − estimadas) y en porcentaje sobre las estimadas. Un valor positivo significa que se trabajó más de lo estimado.
- **FR-009**: Los defectos detectados y resueltos de un Sprint MUST contarse por su Sprint de detección y su Sprint de resolución, respectivamente.
- **FR-010**: Los valores del proyecto MUST ser la suma de los valores de sus Sprints finalizados, salvo la velocidad (FR-004) y los porcentajes, que MUST recalcularse sobre los totales.
- **FR-011**: El sistema MUST calcular las horas y la desviación a partir de los minutos imputados, sin acumular errores de redondeo (D8), y MUST informarlas redondeadas a dos decimales. Los porcentajes y la velocidad del proyecto MUST informarse también con dos decimales.
- **FR-012**: Para un proyecto sin Sprints finalizados, el sistema MUST responder con éxito, sin Sprints y con los valores del proyecto en cero, salvo lo que indica FR-013.
- **FR-013**: Si el proyecto no tiene horas por Story Point, las horas estimadas y la desviación MUST informarse como "sin datos", en el detalle por Sprint y en el proyecto. La desviación porcentual MUST informarse como "sin datos" también cuando las horas estimadas son cero.
- **FR-014**: El sistema MUST responder "proyecto no encontrado" cuando el proyecto no existe.
- **FR-015**: La consulta de métricas MUST no modificar ningún dato.

**Parámetro de horas por Story Point**

- **FR-016**: El sistema MUST permitir definir, cambiar y dejar sin definir las horas por Story Point al editar un proyecto existente. El valor admite decimales y MUST ser mayor que cero.
- **FR-017**: El sistema MUST rechazar un valor cero, negativo o no numérico con un error que identifica el campo, y MUST dejar el proyecto sin cambios.
- **FR-018**: Una edición que no menciona el parámetro MUST conservar su valor; una que lo envía explícitamente vacío MUST dejarlo sin definir.
- **FR-019**: La consulta de un proyecto MUST devolver sus horas por Story Point, o indicar que no están definidas.
- **FR-020**: El parámetro MUST persistir después de cerrar y volver a abrir la aplicación.
- **FR-021**: Al cambiar el parámetro, las horas estimadas de los Sprints finalizados MUST [NEEDS CLARIFICATION: recalcularse con el valor nuevo o conservar las calculadas con el valor anterior, ver User Story 3].
- **FR-022**: La pantalla de edición del proyecto MUST mostrar y permitir modificar el parámetro y, ante un rechazo, MUST mostrar el mensaje de error que devuelve el sistema.

### Key Entities *(include if feature involves data)*

- **Métricas de un Sprint**: valores calculados de un Sprint finalizado; no se guardan. Story Points planificados y completados, velocidad, horas estimadas y reales, desviación (en horas y en porcentaje), porcentaje de historias completadas, y defectos detectados y resueltos.
- **Métricas del proyecto**: los mismos valores agregados sobre los Sprints finalizados del proyecto, junto con el detalle por Sprint; no se guardan.
- **Proyecto**: el de HU-04. Suma las horas por Story Point, opcionales y mayores que cero.
- **Sprint**: el de HU-05. Solo los "Finalizado" aportan a las métricas.
- **Historia comprometida**: una historia completada en el Sprint o arrastrada al cerrarlo (D2c). La arrastrada conserva los Story Points que tenía en ese momento.
- **Registro de horas**: el de HU-07. Se atribuye al Sprint que estaba activo al cargarlo.
- **Defecto**: el de HU-08. Aporta a su Sprint de detección y, si está resuelto, a su Sprint de resolución.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Para un proyecto de referencia con los valores calculados a mano, el 100% de las métricas consultadas coincide con el resultado esperado.
- **SC-002**: El 100% de las consultas sobre proyectos sin Sprints finalizados responde con éxito, sin errores.
- **SC-003**: El 100% de los proyectos sin horas por Story Point informa las horas estimadas y la desviación como "sin datos", nunca como cero.
- **SC-004**: Con los mismos datos, dos consultas sucesivas devuelven exactamente los mismos valores, y editar en el Product Backlog una historia arrastrada no cambia ninguna métrica del Sprint que la arrastró.
- **SC-005**: El 100% de los valores inválidos de horas por Story Point se rechaza con un error sobre ese campo, y el proyecto queda como estaba.
- **SC-006**: Un Scrum Master obtiene las métricas de todos los Sprints finalizados y del proyecto con una sola consulta.

## Assumptions

- **Origen de las horas estimadas (D7)**: el equipo votó el 2026-10-28 que las horas estimadas son Story Points × horas por Story Point del proyecto. El issue #49 pide además confirmarlo con el profesor; si el profesor aclara otra fuente, se ajusta con un fix.
- **Sin cambios de esquema**: `proyectos.horas_por_story_point` existe desde la migración 001, opcional y mayor que cero, y las demás fuentes también existen (`MODELO-DE-DATOS.md` §6 y §8). Esta historia no agrega migraciones siempre que el cambio del parámetro no requiera guardar valores anteriores (FR-021).
- **Trabajo pendiente en la edición del proyecto**: hoy ni la entidad, ni el repositorio, ni la edición, ni la consulta del proyecto manejan el parámetro (`MODELO-DE-DATOS.md` §8). Es tarea de esta historia, y el repositorio no debe perder el valor al editar otros datos.
- **Solo Sprints finalizados**: el backlog pide calcular tras el cierre y §6 indica que un Sprint sin cerrar no tiene datos. Las métricas parciales de un Sprint "Activo" no se piden y quedan fuera de alcance.
- **Valores en cero o "sin datos"**: el backlog admite ambos para un proyecto sin Sprints cerrados. Se usa cero salvo donde el cero sería engañoso: horas estimadas y desviación sin parámetro, y desviación porcentual sin horas estimadas.
- **Velocidad del proyecto**: el backlog no la define; se toma el promedio por Sprint finalizado, que es la definición habitual en Scrum.
- **Historias sin Story Points**: cuentan como comprometidas en el porcentaje y suman cero a los Story Points; no se agrega un aviso específico.
- **Reproducibilidad (D2b)**: los Story Points ya cargados de una historia de un Sprint "Activo" o "Finalizado" no se editan. Cargarlos por primera vez sí se permite, así que completar la estimación de una historia completada que no la tenía cambia las métricas de su Sprint. Es una excepción acotada a la nota del backlog "las métricas de Sprints pasados no cambian".
- **Edición compatible**: el parámetro es opcional en la edición para que los clientes actuales sigan funcionando hasta que se entregue la sub-issue #53.
- **Sin tope superior**: las horas por Story Point no tienen máximo, porque ni el backlog ni el modelo lo piden.
- **Consulta única**: una sola consulta por proyecto devuelve el detalle por Sprint y los valores del proyecto. No hay filtros ni paginación (YAGNI).
- **Sin control de acceso**: igual que en el resto del sistema, no hay autenticación ni permisos por rol.
- **Fuera de alcance**: gráficos e indicadores (HU-09), reporte en PDF (HU-10) e historial de Sprints (HU-12). El parámetro no se pide al crear el proyecto, solo al editarlo (Escenario 3).
- **Dependencias**: HU-04, HU-05, HU-07, HU-08, HU-11 y HU-15, con el backend ya mergeado en `main`.
