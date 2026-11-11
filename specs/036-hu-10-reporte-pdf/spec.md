# Feature Specification: Generación de Reporte en PDF (HU-10)

**Feature Branch**: `feature/HU-10-reporte-pdf`

**Issue**: #10 (HU-10), sub-issues #74 a #78

**Created**: 2026-11-03

**Status**: Implementada (frontend #78 en Sprint 5), pendiente de revisión por PR.

**Input**: User description: "HU-10: Generación de Reporte en PDF (issue #10, sub-issues #74 a #78). Como Scrum Master quiero exportar un reporte en PDF de un proyecto o Sprint con sus historias, estimaciones, esfuerzo, métricas y defectos para poder compartir el estado del proyecto fuera del sistema (por ejemplo, con el profesor o el cliente). Escenario 1: al ver el resumen de un Sprint finalizado, presiono 'Exportar a PDF' y obtengo un archivo con historias planificadas y completadas, estimaciones, esfuerzo registrado, métricas (HU-03) y defectos asociados a ese Sprint. Escenario 2: al ver el resumen de un proyecto, presiono 'Exportar reporte del proyecto' y obtengo los datos del proyecto (nombre, fechas, estado e integrantes), su Product Backlog actual (historias sin Sprint, prioridad y estimación), la lista de Sprints finalizados con Sprint Goal y fechas, el Sprint Activo con sus historias si lo hay y, para los Sprints finalizados, historias planificadas y completadas, estimaciones, esfuerzo, métricas y defectos detectados y resueltos. Escenario 3: un proyecto sin Sprints finalizados (con o sin Sprint Activo) genera el reporte con los datos del proyecto, su Product Backlog y el Sprint Activo si existe, e informa métricas, esfuerzo y defectos como 'sin datos' en vez de fallar. Escenario 4: si la generación falla por un error real, el sistema informa el error en vez de descargar un archivo corrupto o vacío. Planificación del Sprint 4 (2026-11-03): la especificación, los escenarios BDD, el backend y los tests (#74 a #77) se adelantan sin depender de HU-03 en su totalidad; hasta que HU-03 esté lista, las métricas se informan como 'sin datos'. El frontend (#78) y el cierre de la historia quedan en el Sprint 5. El reporte es de solo lectura y no modifica datos."

## Clarifications

### Session 2026-11-03

- Q: ¿Qué pasa si se solicita el reporte de un Sprint que no está finalizado? → A: Se rechaza con un mensaje claro y no se genera ningún archivo. El reporte de Sprint es solo para Sprints finalizados; la interfaz solo ofrece el botón en el resumen de uno finalizado.
- Q: ¿El reporte del proyecto incluye el esfuerzo y los defectos del Sprint Activo? → A: No. El Sprint Activo se lista con sus historias (título, estado y estimación); su esfuerzo, sus métricas y sus defectos no se incluyen ni suman a los totales. Ya está definido en el issue #10 y en las notas del backlog.
- Q: ¿Qué datos de cada defecto incluye el reporte? → A: Descripción, severidad, estado e historia relacionada, más los Sprints de detección y resolución. Se toma la definición de la spec 023 de HU-08.

### Session 2026-11-07

- Q: ¿El esfuerzo se informa por historia o agregado? → A: Agregado por Sprint, con horas reales, horas estimadas y desviación. El reporte no muestra horas por historia.
- Q: ¿El reporte incluye las historias de un Sprint "Pendiente"? → A: No. El reporte del proyecto lista solo los Sprints finalizados y, si existe, el Sprint Activo; un Sprint "Pendiente" no aporta historias, métricas, esfuerzo ni defectos.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Exportar el reporte de un Sprint finalizado (Priority: P1)

Como Scrum Master, desde el resumen de un Sprint finalizado presiono "Exportar a PDF" y descargo un reporte con las historias planificadas y completadas, sus estimaciones, el esfuerzo registrado, las métricas del Sprint y los defectos asociados, para compartir el resultado de esa iteración fuera del sistema.

**Why this priority**: Es el primer criterio de aceptación y la porción más acotada del valor: prueba el armado del PDF con un conjunto de datos limitado. Es la base del reporte del proyecto, que reutiliza las mismas secciones por Sprint.

**Independent Test**: Se siembra un proyecto con un Sprint finalizado, historias completadas y arrastradas, estimaciones, horas imputadas y defectos; se solicita el reporte del Sprint y se verifica el PDF contra los valores calculados a mano. También se verifica que un Sprint no finalizado se rechace y que un Sprint inexistente responda "Sprint no encontrado".

**Acceptance Scenarios**:

1. **Given** un Sprint finalizado con historias planificadas y completadas, estimaciones, horas imputadas y defectos, **When** solicito su reporte, **Then** obtengo un PDF válido con la identificación del proyecto y del Sprint (goal y fechas), las historias planificadas y completadas, sus estimaciones, el esfuerzo registrado, las métricas de HU-03 y los defectos asociados.
2. **Given** un Sprint finalizado sin horas imputadas, sin defectos o sin historias en alguna sección, **When** solicito su reporte, **Then** el PDF se genera con esa sección vacía, en cero o como "sin datos" según corresponda, sin fallar.
3. **Given** un Sprint "Pendiente" o "Activo", **When** solicito su reporte, **Then** el sistema rechaza la operación con un mensaje claro y no genera ningún archivo.
4. **Given** un Sprint que no existe, **When** solicito su reporte, **Then** el sistema responde "Sprint no encontrado" antes de generar nada.
5. **Given** un Sprint finalizado con una historia arrastrada desde otro Sprint, **When** solicito su reporte, **Then** esa historia aparece como planificada con los Story Points que tenía al cerrar y como no completada en ese Sprint; si se completó en otro Sprint, aparece como completada en el reporte de ese otro Sprint.
6. **Given** un defecto detectado en el Sprint y otro resuelto en el Sprint, **When** solicito su reporte, **Then** ambos aparecen asociados al Sprint con su descripción, severidad, estado y la historia relacionada; un defecto detectado en otro Sprint y resuelto en este aparece como resuelto en este reporte y como detectado en el del otro Sprint.

---

### User Story 2 - Exportar el reporte completo de un proyecto (Priority: P1)

Como Scrum Master, desde el resumen de un proyecto presiono "Exportar reporte del proyecto" y descargo un PDF con los datos del proyecto, su Product Backlog, sus Sprints finalizados y el Sprint Activo si existe, para compartir el estado completo del proyecto.

**Why this priority**: Es el segundo criterio de aceptación y el valor principal a nivel de proyecto: reúne el estado actual y el historial en un solo documento. Depende de que el reporte de Sprint esté resuelto, porque repite sus secciones por cada Sprint finalizado.

**Independent Test**: Se siembra un proyecto con integrantes, Product Backlog, al menos dos Sprints finalizados y un Sprint Activo, con historias, estimaciones, horas y defectos; se solicita el reporte del proyecto y se verifica que incluya todas las secciones, que el Sprint Activo se liste con sus historias y que no sume a las métricas, el esfuerzo ni los defectos.

**Acceptance Scenarios**:

1. **Given** un proyecto con Sprints finalizados, Product Backlog y un Sprint Activo, **When** solicito el reporte del proyecto, **Then** el PDF incluye los datos del proyecto, el Product Backlog actual, cada Sprint finalizado con su detalle y el Sprint Activo con sus historias; las métricas, el esfuerzo y los defectos se informan solo con los Sprints finalizados.
2. **Given** un proyecto sin Sprint Activo, **When** solicito el reporte, **Then** el PDF se genera sin la sección del Sprint Activo (o la marca como "sin Sprint activo") y el resto se incluye normalmente.
3. **Given** un proyecto sin historias en el Product Backlog, **When** solicito el reporte, **Then** esa sección aparece vacía o como "sin historias" y el resto se genera normalmente.
4. **Given** un proyecto que no existe, **When** solicito su reporte, **Then** el sistema responde "proyecto no encontrado" antes de generar nada.
5. **Given** un proyecto con un Sprint finalizado que no tiene horas ni defectos, **When** solicito el reporte, **Then** esas secciones aparecen en cero o "sin datos" y el PDF se genera.
6. **Given** un defecto detectado en un Sprint finalizado y resuelto en otro Sprint finalizado, **When** solicito el reporte del proyecto, **Then** aparece como detectado en el primer Sprint y como resuelto en el segundo, igual que en las métricas de HU-03.
7. **Given** un proyecto con historias en el Product Backlog y un Sprint Activo, **When** solicito el reporte, **Then** las historias sin Sprint aparecen en el Product Backlog con su prioridad y estimación, y las del Sprint Activo aparecen en la sección de ese Sprint con su estado y estimación, sin sumar a los totales del proyecto.
8. **Given** un proyecto con varios Sprints finalizados, **When** solicito el reporte, **Then** el PDF incluye todos, en orden, cada uno con su Sprint Goal, fechas y su detalle.

---

### User Story 3 - Reporte de un proyecto sin Sprints finalizados (Priority: P2)

Como Scrum Master, genero el reporte de un proyecto que todavía no cerró ningún Sprint (con o sin Sprint Activo) y obtengo los datos del proyecto, su Product Backlog y, si existe, el Sprint Activo, con las métricas, el esfuerzo y los defectos como "sin datos", sin errores.

**Why this priority**: Es el estado inicial de todo proyecto y el Escenario 3 del issue. Sin esto, el reporte de un proyecto nuevo parecería roto o generaría un archivo incompleto.

**Independent Test**: Se siembra un proyecto sin Sprints finalizados, una vez con Sprint Activo y otra sin él; se solicita el reporte y se verifica que incluya los datos del proyecto, el Product Backlog y el Sprint Activo cuando existe, y que las métricas, el esfuerzo y los defectos se informen como "sin datos" sin error.

**Acceptance Scenarios**:

1. **Given** un proyecto sin Sprints finalizados y con un Sprint Activo, **When** solicito el reporte, **Then** el PDF incluye los datos del proyecto, el Product Backlog y el Sprint Activo con sus historias; las métricas, el esfuerzo y los defectos se informan como "sin datos"; no hay error.
2. **Given** un proyecto sin Sprints finalizados y sin Sprint Activo, **When** solicito el reporte, **Then** el PDF incluye los datos del proyecto y el Product Backlog; no hay sección de Sprint Activo; las métricas, el esfuerzo y los defectos se informan como "sin datos"; no hay error.
3. **Given** un proyecto sin Sprints finalizados y sin historias en el Product Backlog, **When** solicito el reporte, **Then** el PDF se genera con las secciones vacías o "sin datos" y no falla.
4. **Given** un proyecto que no existe, **When** solicito su reporte, **Then** el sistema responde "proyecto no encontrado".

---

### User Story 4 - Informar un error real de generación (Priority: P2)

Como Scrum Master, si la generación del PDF falla por un error real (lectura de datos, armado del documento o escritura del archivo), veo un mensaje de error y no descargo un archivo corrupto, vacío o parcial.

**Why this priority**: Es el Escenario 4 del issue. Un archivo corrupto o vacío es peor que un error visible, porque el Scrum Master podría compartirlo creyendo que es válido.

**Independent Test**: Se fuerza un fallo real de generación (por ejemplo, un fallo simulado de lectura de datos o de escritura del archivo); se verifica que el sistema responda con un error y que no entregue ningún archivo. También se verifica que una generación exitosa produzca un PDF válido y no vacío.

**Acceptance Scenarios**:

1. **Given** un error real durante la generación, **When** solicito el reporte, **Then** el sistema informa el error y no entrega un archivo corrupto, vacío o parcial.
2. **Given** una generación exitosa, **When** solicito el reporte, **Then** el archivo descargado es un PDF válido con el contenido solicitado.
3. **Given** que el error ocurre después de comenzar a armar el documento, **When** la generación falla, **Then** no se descarga ningún archivo parcial.

---

### User Story 5 - Exportar desde la interfaz (Priority: P3)

Como Scrum Master, veo el botón "Exportar a PDF" en el resumen de un Sprint finalizado y "Exportar reporte del proyecto" en el proyecto; al presionarlos, descargo el archivo o veo el mensaje de error.

**Why this priority**: La historia se completa cuando el Scrum Master puede disparar el reporte desde la interfaz, pero el frontend es la sub-issue #78 y está planificado para el Sprint 5. No frena el avance del backend, los escenarios BDD ni los tests del Sprint 4.

**Independent Test**: Con la aplicación y el backend funcionando, se abre el resumen de un Sprint finalizado y el resumen de un proyecto; se presionan los dos botones y se verifica la descarga del PDF. Se simula un error del backend y se verifica el mensaje en pantalla y que no se descargue un archivo corrupto.

**Acceptance Scenarios**:

1. **Given** el resumen de un Sprint finalizado, **When** presiono "Exportar a PDF", **Then** el PDF del Sprint se descarga.
2. **Given** el resumen de un proyecto, **When** presiono "Exportar reporte del proyecto", **Then** el PDF del proyecto se descarga.
3. **Given** que el backend informa un error de generación, **When** presiono el botón, **Then** la interfaz muestra el error y no descarga un archivo corrupto o vacío.
4. **Given** un Sprint que no está finalizado, **When** abro su resumen, **Then** el botón de exportar no se ofrece, porque el reporte de Sprint es solo para Sprints finalizados.
5. **Given** una descarga exitosa, **When** se guarda el archivo, **Then** su nombre identifica el proyecto y, en el caso del Sprint, también el Sprint.

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** generar el reporte sin modificar ningún dato del proyecto, Sprint, historia, hora o defecto.
- El sistema **shall** tomar las métricas de HU-03 tal como esa historia las define; el reporte no **shall** recalcularlas con reglas propias.
- **WHILE** un Sprint está "Pendiente" o "Activo", el sistema **shall** excluirlo de las métricas, el esfuerzo y los defectos del reporte del proyecto; el Sprint Activo solo **shall** listarse con sus historias.
- **WHEN** un Sprint está "Pendiente", **THEN** el sistema **shall** no incluir sus historias en el reporte del proyecto.
- **WHEN** se solicita el reporte de un Sprint que no está "Finalizado", **THEN** el sistema **shall** rechazar la operación con un mensaje claro y **shall** no generar archivo.
- **WHEN** se solicita el reporte de un proyecto sin Sprints finalizados, **THEN** el sistema **shall** generarlo con los datos del proyecto, el Product Backlog y el Sprint Activo si existe, y **shall** informar métricas, esfuerzo y defectos como "sin datos".
- **WHEN** se solicita el reporte de un proyecto o Sprint inexistente, **THEN** el sistema **shall** responder "proyecto no encontrado" o "Sprint no encontrado" antes de generar el archivo.
- **IF** un Sprint finalizado no tiene historias comprometidas, **THEN** el reporte **shall** mostrar la sección vacía o con cero Story Points, sin fallar.
- **IF** un Sprint finalizado no tiene horas imputadas, **THEN** el reporte **shall** informar cero horas reales; el resto de las métricas **shall** seguir las reglas de HU-03.
- **WHEN** el reporte informa el esfuerzo, **THEN** el sistema **shall** mostrarlo agregado por Sprint (horas reales, estimadas y desviación) y **shall** no listar horas por historia.
- **IF** un Sprint finalizado no tiene defectos detectados ni resueltos, **THEN** el reporte **shall** mostrar la sección vacía o "sin defectos", sin fallar.
- **IF** una historia no tiene estimación, **THEN** el reporte **shall** mostrarla como "sin estimar" y **shall** sumar cero a los totales de Story Points, igual que HU-03.
- **WHEN** una historia fue arrastrada desde un Sprint anterior, **THEN** el reporte del Sprint anterior **shall** mostrarla como planificada y no completada, con los Story Points que tenía al cerrar; si se completó después, **shall** aparecer como completada en el Sprint donde se completó.
- **WHEN** un defecto fue detectado en un Sprint y resuelto en otro, **THEN** **shall** aparecer como detectado en el reporte del primero y como resuelto en el del segundo.
- **IF** el proyecto no tiene una fecha de inicio o de fin, **THEN** el reporte **shall** mostrarla como "sin definir", igual que el panel del proyecto.
- **IF** el proyecto no tiene integrantes, **THEN** el reporte **shall** mostrar "sin integrantes" en la sección del equipo.
- **IF** el Product Backlog no tiene historias, **THEN** el reporte **shall** mostrar la sección vacía o "sin historias", sin fallar.
- **IF** el proyecto no tiene Sprint Activo, **THEN** el reporte **shall** omitir esa sección o marcarla como "sin Sprint activo", sin fallar.
- **IF** el proyecto no tiene horas por Story Point o HU-03 no está disponible, **THEN** las métricas que dependan de ese dato **shall** informarse como "sin datos"; el resto del reporte **shall** generarse igual.
- **WHERE** el texto incluya tildes, eñes u otros caracteres UTF-8, el reporte **shall** conservarlos correctamente en el PDF.
- **WHEN** el reporte tenga muchas historias, Sprints o defectos, **THEN** el sistema **shall** incluirlos todos y **shall** paginar el documento; no **shall** truncar datos.
- **IF** la generación falla por un error real, **THEN** el sistema **shall** informar el error y **shall** no entregar un archivo corrupto, vacío o parcial.
- **WHEN** se solicita el mismo reporte dos veces con los mismos datos, **THEN** el contenido **shall** ser equivalente y **shall** no haber cambios en la base.

## Requirements *(mandatory)*

### Functional Requirements

**Reporte de un Sprint**

- **FR-001**: El sistema MUST permitir solicitar el reporte en PDF de un Sprint finalizado existente.
- **FR-002**: El reporte de Sprint MUST incluir el nombre del proyecto, el Sprint Goal, las fechas del Sprint, las historias planificadas y completadas, sus estimaciones, el esfuerzo registrado, las métricas de HU-03 y los defectos asociados.
- **FR-003**: Las historias planificadas MUST incluir las completadas y las arrastradas al cerrar el Sprint, con los Story Points que tenían en ese momento (D2c).
- **FR-004**: El reporte MUST usar las mismas definiciones y valores que HU-03 para Story Points planificados y completados, velocidad, horas estimadas y reales, desviación, porcentaje de historias completadas y defectos detectados y resueltos; MUST NOT recalcularlos con otra regla.
- **FR-005**: Si el Sprint no está "Finalizado", el sistema MUST rechazar la solicitud con un mensaje claro y MUST NOT generar ningún archivo.
- **FR-006**: Si el Sprint no existe, el sistema MUST responder "Sprint no encontrado" antes de generar.
- **FR-007**: Si el Sprint finalizado no tiene historias, horas o defectos en una sección, el reporte MUST mostrarla vacía, en cero o como "sin datos" según corresponda, sin fallar.
- **FR-008**: El reporte MUST incluir los defectos detectados en el Sprint y los resueltos en el Sprint, con descripción, severidad, estado y la historia relacionada; un defecto detectado en otro Sprint y resuelto en este MUST aparecer como resuelto en este reporte.
- **FR-009**: El reporte MUST mostrar cada historia con su estimación o "sin estimar"; una historia sin estimación MUST sumar cero a los totales de Story Points.
- **FR-010**: El reporte de Sprint MUST ser un PDF válido, no vacío ni corrupto, con todos los datos y paginado si hace falta.

**Reporte de un proyecto**

- **FR-011**: El sistema MUST permitir solicitar el reporte en PDF de un proyecto existente.
- **FR-012**: El reporte de proyecto MUST incluir el nombre del proyecto, sus fechas (o "sin definir"), su estado derivado (HU-13) y sus integrantes (o "sin integrantes").
- **FR-013**: El reporte MUST incluir el Product Backlog actual: las historias sin Sprint, con su prioridad y su estimación (o "sin estimar"), en el orden de HU-15.
- **FR-014**: El reporte MUST incluir todos los Sprints finalizados, en orden, con su Sprint Goal y fechas; para cada uno, MUST incluir las historias planificadas y completadas, sus estimaciones, el esfuerzo, las métricas de HU-03 y los defectos detectados y resueltos.
- **FR-015**: El reporte MUST incluir el Sprint Activo, si existe, con sus historias (título, estado y estimación); MUST NOT incluir su esfuerzo, sus defectos ni sus métricas, ni sumarlos a los totales del proyecto.
- **FR-016**: Los totales del proyecto (métricas, esfuerzo y defectos) MUST calcularse solo con los Sprints finalizados.
- **FR-017**: Si el proyecto no tiene Sprints finalizados, el reporte MUST generarse con los datos del proyecto, el Product Backlog y el Sprint Activo si existe, e informar las métricas, el esfuerzo y los defectos como "sin datos", sin fallar.
- **FR-018**: Si no hay historias en el Product Backlog, Sprint Activo o defectos, el reporte MUST mostrar la sección vacía o con el mensaje correspondiente ("sin historias", "sin Sprint activo", "sin defectos"), sin fallar.
- **FR-019**: Si el proyecto no existe, el sistema MUST responder "proyecto no encontrado" antes de generar.
- **FR-020**: El reporte MUST incluir un defecto detectado en un Sprint finalizado y uno resuelto en un Sprint finalizado; si son Sprints distintos, MUST aparecer en cada uno según su detección o resolución, y en los totales del proyecto como en HU-03.
- **FR-021**: El reporte de proyecto MUST ser un PDF válido, no vacío ni corrupto, con todos los datos y paginado si hace falta.

**Generación, errores y descarga**

- **FR-022**: Si la generación falla por un error real (lectura de datos, armado del PDF o escritura del archivo), el sistema MUST informar el error y MUST NOT entregar un archivo corrupto, vacío o parcial.
- **FR-023**: La generación de un reporte MUST NOT modificar ningún dato.
- **FR-024**: El archivo descargado MUST ser un PDF válido y su nombre MUST identificar el proyecto y, cuando corresponda, el Sprint.
- **FR-025**: El reporte MUST reflejar los datos al momento de la solicitud; dos solicitudes con los mismos datos MUST producir contenido equivalente.
- **FR-026**: La interfaz MUST ofrecer "Exportar a PDF" en el resumen de un Sprint finalizado y "Exportar reporte del proyecto" en el resumen del proyecto; el botón de Sprint MUST NOT ofrecerse para Sprints no finalizados.
- **FR-027**: Si el backend informa un error de generación, la interfaz MUST mostrar el mensaje y MUST NOT descargar un archivo corrupto o vacío.
- **FR-028**: Mientras HU-03 no esté integrada, las métricas que dependan de ella MUST informarse como "sin datos" sin romper el reporte; cuando HU-03 esté disponible, el reporte MUST usar sus valores.
- **FR-029**: El reporte MUST conservar los caracteres UTF-8 (tildes, eñes y otros) en todas las secciones.

### Key Entities *(include if feature involves data)*

- **Reporte**: documento PDF generado a demanda, de un Sprint finalizado o de un proyecto; no se guarda en la base.
- **Proyecto**: el de HU-04. Aporta nombre, fechas, estado derivado (HU-13) e integrantes.
- **Sprint**: el de HU-05. Aporta goal, fechas, estado y, si está finalizado, sus métricas y su detalle.
- **Historia**: la de HU-01/HU-15. Aporta título, estado, prioridad, estimación y pertenencia al Sprint o al Product Backlog.
- **Historia comprometida**: completada en el Sprint o arrastrada al cerrarlo (D2c); la arrastrada conserva los Story Points que tenía al cerrar.
- **Registro de horas**: el de HU-07. Aporta el esfuerzo atribuido a cada Sprint.
- **Defecto**: el de HU-08. Aporta descripción, severidad, estado, historia relacionada y Sprints de detección y resolución.
- **Métricas**: las de HU-03; el reporte las consume y no las recalcula.
- **Product Backlog**: las historias sin Sprint, con prioridad y estimación, en el orden de HU-15.
- **Integrante**: los del proyecto (HU-04).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Para un Sprint finalizado de referencia, el 100% de las secciones y valores del PDF coincide con los datos y las métricas de HU-03 calculados a mano.
- **SC-002**: Para un proyecto de referencia, el 100% de las secciones requeridas aparece: datos del proyecto, Product Backlog, todos los Sprints finalizados y el Sprint Activo si existe; el Sprint Activo no altera los totales.
- **SC-003**: El 100% de los proyectos sin Sprints finalizados genera un PDF con métricas, esfuerzo y defectos como "sin datos", sin error.
- **SC-004**: El 100% de los fallos reales de generación muestra un error y no descarga un archivo corrupto, vacío o parcial.
- **SC-005**: El 100% de las solicitudes de un proyecto o Sprint inexistente responde "no encontrado" antes de generar.
- **SC-006**: El 100% de los reportes conserva los caracteres UTF-8 y incluye todos los elementos, sin truncarlos, con paginación cuando hace falta.
- **SC-007**: Un Scrum Master obtiene el reporte con una sola acción y el PDF descargado abre correctamente.
- **SC-008**: El 100% de las solicitudes de reporte de un Sprint no finalizado se rechaza con un mensaje claro y sin archivo.
- **SC-009**: Dos solicitudes del mismo reporte con los mismos datos devuelven contenido equivalente y no modifican la base.

## Constraints

- **PDF**: el reporte se entrega como archivo PDF, generado con la librería que fija la constitución del proyecto (Maroto). No se introducen dependencias nuevas fuera del stack acordado.
- **Solo lectura**: el reporte no crea, edita ni borra proyectos, Sprints, historias, horas, defectos ni miembros.
- **Sin cambios de esquema**: esta historia no agrega tablas ni columnas; usa los datos que ya existen (MODELO-DE-DATOS.md).
- **Datos del mismo proyecto**: el reporte de un proyecto solo incluye sus propios Sprints, historias, horas, defectos e integrantes.
- **Una historia a la vez**: el Sprint 4 adelanta spec, BDD, backend y tests (#74 a #77); el frontend (#78) y el cierre quedan en el Sprint 5.
- **Idioma**: el reporte se genera en español y conserva tildes y eñes.
- **Sin autenticación**: igual que el resto del sistema en este Sprint, no hay login ni permisos por rol; el reporte es accesible desde la aplicación.

## Assumptions

- **Origen de las métricas**: las métricas del reporte son las de HU-03. Durante el Sprint 4, mientras HU-03 no esté integrada, se muestran como "sin datos"; cuando HU-03 esté mergeada, el reporte las consume sin duplicar el cálculo.
- **Solo Sprints finalizados**: las métricas, el esfuerzo y los defectos del reporte del proyecto se calculan solo con Sprints finalizados; el Sprint Activo se lista con sus historias, pero no suma.
- **Reporte de Sprint solo para finalizados**: la pantalla de Sprint solo ofrece el botón en un Sprint finalizado; si se pide por otro medio, el sistema lo rechaza.
- **Defectos**: se listan con descripción, severidad, estado y la historia relacionada. Los conteos de detectados y resueltos coinciden con HU-03.
- **Product Backlog**: se muestran solo las historias sin Sprint, como en HU-15, con prioridad y estimación.
- **Proyecto sin integrantes**: no debería ocurrir porque al crear el proyecto se registra su creador (HU-04); si ocurriera, el reporte muestra "sin integrantes".
- **Proyecto sin fechas**: se muestran como "sin definir", igual que el panel del proyecto.
- **Sin Sprints finalizados**: no es un error; es un estado válido y el reporte se genera con "sin datos" en las métricas, el esfuerzo y los defectos.
- **Paginación**: el reporte incluye todos los elementos; si no entran en una página, se pagina. No se trunca información.
- **Nombre del archivo**: identifica el proyecto y, si es un Sprint, el Sprint; el formato exacto se define en el plan/contrato.
- **Diseño del PDF**: la disposición visual (secciones, orden, tablas) se define en el plan y el contrato; esta spec fija el contenido y el comportamiento.
- **Dependencias**: HU-04, HU-13, HU-05, HU-06, HU-15, HU-07, HU-03 y HU-08. HU-03 puede no estar completa durante el Sprint 4; el reporte se adelanta con "sin datos" en lo que dependa de ella.
- **Frontend**: la sub-issue #78 (botones, descarga y mensajes) se entrega en el Sprint 5; el backend, los escenarios BDD y los tests se entregan en el Sprint 4.

## Out of Scope

- Enviar el reporte por correo, guardarlo en el servidor o mantener un historial de reportes generados.
- Editar o borrar datos desde el reporte.
- Gráficos, dashboard (HU-09) e historial de Sprints (HU-12).
- Autenticación, permisos por rol o control de acceso (HU-20/HU-21).
- Reportes de Sprints no finalizados, de Sprints Activos con métricas parciales o de varios proyectos a la vez.
- Plantillas personalizadas, filtros, selección de secciones, programación o formatos distintos de PDF.
- Cambios de esquema o migraciones.
- Firma digital, código QR o marca de agua.
