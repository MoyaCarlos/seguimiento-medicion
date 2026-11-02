# Feature Specification: Imputación de horas en pantalla (HU-07 Frontend)

**Feature Branch**: `feature/HU-07-frontend`

**Issue**: #63 (parte de #7)

**Created**: 2026-11-02

**Status**: Implementada, pendiente de revisión por PR.

**Input**: User description: "HU-07 Frontend: imputación de horas trabajadas (sub-issue #63 de la HU-07, issue padre #7). El backend ya está en main (spec 019-hu-07-imputacion-horas: POST y GET /backlog/{id}/horas). Alcance: pantalla con formulario de registro de horas (integrante, fecha, actividad y horas) sobre historias del Sprint activo, mostrar los errores de validación que devuelve la API identificando el campo, y mostrar el esfuerzo total acumulado de cada historia. No hay login, por lo que hay que decidir en el clarify cómo se elige quién imputa las horas (propuesta del equipo: elegir al integrante de una lista). Sin cambios de backend ni de esquema."

**Dependencia**: el backend de HU-07 ya está en `main` (spec 019, PR #157). Esta spec solo cubre la pantalla: consume el contrato de la spec 019 sin cambiarlo.

## Clarifications

### Session 2026-11-02

- Q: Como no hay login, ¿cómo se elige quién imputa las horas? → A: De una lista con los integrantes del proyecto, que empieza sin nadie seleccionado; hay que elegir cada vez que se abre el formulario y no se guarda nada en el navegador.
- Q: ¿Dónde va el formulario de registro de horas? → A: En la pantalla de Sprints, dentro de cada historia del Sprint activo: un botón "Registrar horas" lo abre ahí mismo. No se agrega una pantalla ni una entrada de menú.
- Q: ¿En qué historias se muestra el esfuerzo total acumulado? → A: En todas las de la pantalla de Sprints (Pendiente, Activo y Finalizado). La pantalla de Historias (Product Backlog) no cambia.
- Q: Además del total, ¿se muestra el detalle de cada carga? → A: Sí, plegado: junto al total hay un "Ver detalle" que despliega las cargas de esa historia con integrante, fecha, actividad y horas.

## Objetivo

Que un integrante del equipo registre desde la pantalla las horas reales que trabajó en una historia del Sprint activo, entienda qué dato corregir cuando la carga es rechazada, y vea el esfuerzo total acumulado de cada historia, sin recargar la aplicación.

## Entradas y salidas

- **Entradas (de la persona)**: historia sobre la que imputa, integrante, fecha trabajada, actividad y cantidad de horas.
- **Entradas (del sistema)**: integrantes del proyecto, y esfuerzo real de cada historia (total y detalle de cargas).
- **Salida en pantalla**: el esfuerzo total acumulado de cada historia, en horas; una confirmación al guardar; el error junto al campo que corresponde, o como aviso general cuando el error no es de un campo del formulario.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Registrar horas en una historia del Sprint activo (Priority: P1)

Como integrante del equipo, desde una historia del Sprint activo abro el formulario de registro, elijo mi nombre, indico fecha, actividad y horas, y guardo, para que el esfuerzo real de esa historia quede registrado sin salir de la pantalla.

**Why this priority**: Es el núcleo de la sub-issue. Sin el formulario, las horas solo se pueden cargar por fuera de la interfaz y el producto no es usable por el equipo.

**Independent Test**: Con un proyecto que tiene integrantes y un Sprint activo con una historia, se abre el formulario, se completan los cuatro datos y se guarda; se verifica que se envían los datos elegidos y que aparece la confirmación.

**Acceptance Scenarios**:

1. **Given** una historia del Sprint activo, **When** elijo "Registrar horas", **Then** se abre el formulario dentro de esa historia, con la lista de integrantes del proyecto sin nadie seleccionado y los campos de fecha, actividad y horas.
2. **Given** el formulario abierto, **When** elijo un integrante, indico la fecha de hoy, una actividad y 4 horas, y guardo, **Then** la carga se registra y veo una confirmación.
3. **Given** el formulario abierto, **When** escribo las horas con coma decimal (1,5), **Then** la carga se registra como 1,5 horas.
4. **Given** una historia de un Sprint "Pendiente" o "Finalizado", **When** la veo en pantalla, **Then** no se ofrece registrar horas.
5. **Given** una carga en curso, **When** intento guardar otra vez, **Then** no se envía una segunda carga.

---

### User Story 2 - Ver qué dato corregir cuando la carga se rechaza (Priority: P1)

Como integrante del equipo, si la carga es inválida veo el error junto al campo con problemas y conservo lo que escribí, para corregirlo sin volver a cargar todo.

**Why this priority**: La sub-issue lo pide de forma explícita, y el Escenario 2 de HU-07 pide que el error se vea en la interfaz.

**Independent Test**: Se intenta guardar con cada dato inválido por separado y se verifica que el mensaje aparece junto al campo correcto, que lo escrito se conserva y que el total no cambia.

**Acceptance Scenarios**:

1. **Given** el formulario, **When** intento guardar sin elegir un integrante, **Then** veo un error sobre el integrante y no se envía la carga.
2. **Given** el formulario, **When** la API rechaza las horas (negativas, cero, texto no numérico, más de 24 o sin minutos enteros), **Then** veo el mensaje de la API junto al campo de horas.
3. **Given** el formulario, **When** la API rechaza la actividad (vacía, solo espacios o más de 500 caracteres), **Then** veo el mensaje de la API junto al campo de actividad.
4. **Given** el formulario, **When** la API rechaza la fecha (futura o fuera del Sprint), **Then** veo el mensaje de la API junto al campo de fecha.
5. **Given** el formulario, **When** la API rechaza al integrante, **Then** veo el mensaje de la API junto al campo de integrante.
6. **Given** que el Sprint se cerró mientras tenía el formulario abierto, **When** guardo, **Then** veo como aviso general el mensaje de que solo se imputan horas a historias del Sprint activo.
7. **Given** cualquiera de estos rechazos, **When** veo el error, **Then** lo que escribí sigue en el formulario y el total de la historia no cambia.
8. **Given** el formulario con un error mostrado, **When** vuelvo a guardar y la carga se acepta, **Then** el error desaparece.

---

### User Story 3 - Ver el esfuerzo acumulado de cada historia (Priority: P2)

Como integrante del equipo, veo el total de horas trabajadas de cada historia, para saber cuánto trabajo real lleva.

**Why this priority**: La sub-issue lo pide y es lo que hace visible el resultado de cargar horas, pero depende de que existan cargas (P1).

**Independent Test**: Con historias con y sin cargas, se verifica el total mostrado en cada una; al guardar una carga nueva, el total se actualiza sin recargar.

**Acceptance Scenarios**:

1. **Given** una historia con cargas, **When** la veo en pantalla, **Then** veo el total acumulado en horas (por ejemplo, 5,5 h).
2. **Given** una historia sin cargas, **When** la veo en pantalla, **Then** veo un total de 0 h, sin tratarlo como error.
3. **Given** una historia con un total mostrado, **When** guardo una carga válida, **Then** el total se actualiza sin recargar la aplicación.
4. **Given** que falla la consulta del esfuerzo de una historia, **When** veo la historia, **Then** veo un aviso en lugar del total, y el resto de la historia y de las demás historias sigue visible.
5. **Given** una historia con cargas, **When** elijo "Ver detalle", **Then** veo cada carga con integrante, fecha, actividad y horas, en el orden que devuelve la API.
6. **Given** una historia sin cargas, **When** la veo en pantalla, **Then** no se ofrece "Ver detalle".
7. **Given** el detalle desplegado, **When** guardo una carga válida en esa historia, **Then** la carga nueva aparece en el detalle.

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- La pantalla **shall** mostrar los mensajes de la API tal como llegan, sin reescribirlos ni repetir las reglas de validación de la spec 019: esas reglas viven solo en el backend.
- **WHEN** la API informa un error sobre `integrante_id`, `fecha`, `actividad` u `horas`, **THEN** la pantalla **shall** mostrarlo junto a ese campo.
- **IF** la API informa un error sobre otro dato o sin campo (por ejemplo, historia fuera del Sprint activo o no encontrada), **THEN** la pantalla **shall** mostrarlo como aviso general del formulario.
- **IF** no se puede conectar con el servidor al guardar, **THEN** la pantalla **shall** mostrar el aviso de conexión y **shall** conservar lo escrito.
- **IF** se intenta guardar sin un integrante elegido, **THEN** la pantalla **shall** mostrar un error sobre el integrante y **shall not** enviar la carga. Es la única validación propia de la pantalla: sin elección no hay identificador que enviar.
- **WHEN** las horas se escriben con coma decimal, **THEN** la pantalla **shall** enviarlas con punto decimal. Es una traducción de formato, no una validación: la cantidad la valida el backend.
- **WHERE** las horas no son numéricas, la pantalla **shall** enviar el texto tal cual para que el backend responda el error que corresponde.
- **WHILE** una carga está en curso, la pantalla **shall** deshabilitar el guardado e indicar que está enviando.
- **WHEN** se abre el formulario, **THEN** la pantalla **shall** consultar los integrantes del proyecto; **WHILE** esa consulta está en curso, **shall** mostrar un indicador de carga.
- **IF** falla la consulta de integrantes, **THEN** la pantalla **shall** mostrar el error y **shall not** permitir guardar.
- **IF** el proyecto no tiene integrantes, **THEN** la lista **shall** quedar solo con la opción vacía y guardar **shall** mostrar el error sobre el integrante, sin enviar la carga.
- **IF** el total acumulado no da una cantidad exacta con dos decimales (por ejemplo, 61 minutos), **THEN** la pantalla **shall** redondearlo a dos decimales al mostrarlo (1,02 h).
- **WHEN** se abre el formulario, **THEN** la fecha **shall** proponer el día de hoy o, si el Sprint ya pasó su fecha de fin, esa fecha de fin. La persona puede cambiarla.
- **WHEN** una carga se guarda correctamente, **THEN** la pantalla **shall** vaciar actividad y horas y **shall** conservar integrante y fecha, para facilitar otra carga.
- **WHEN** el Sprint se cierra desde esa misma pantalla, **THEN** la pantalla **shall** dejar de ofrecer el registro de horas y **shall** seguir mostrando el total de las historias que quedan en ese Sprint (las "Completada"). Las no completadas vuelven al Product Backlog (HU-05): salen de la tarjeta y su total no se ve hasta que se planifiquen en otro Sprint, aunque sus horas se conservan.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: La pantalla de Sprints MUST ofrecer un botón "Registrar horas" en cada historia de un Sprint "Activo", que abre el formulario dentro de esa historia, y MUST NOT ofrecerlo en historias de Sprints "Pendiente" o "Finalizado".
- **FR-002**: El formulario MUST pedir integrante, fecha, actividad y horas, y MUST enviarlos en una sola acción de guardado.
- **FR-003**: El integrante MUST elegirse de una lista con los integrantes del proyecto, por nombre. La lista MUST empezar sin nadie seleccionado cada vez que se abre el formulario, y la elección MUST NOT guardarse en el navegador.
- **FR-004**: La pantalla MUST mostrar cada error de validación junto al campo que informa la API, y como aviso general si no corresponde a un campo del formulario.
- **FR-005**: Ante un rechazo, la pantalla MUST conservar lo escrito y MUST dejar sin cambios el total mostrado.
- **FR-006**: La pantalla de Sprints MUST mostrar el esfuerzo total acumulado de cada historia de los Sprints "Pendiente", "Activo" y "Finalizado", en horas, incluidas las que no tienen cargas (0 h). La pantalla de Historias (Product Backlog) MUST NOT cambiar.
- **FR-007**: Al guardar una carga válida, la pantalla MUST mostrar una confirmación y actualizar el total de esa historia sin recargar la aplicación.
- **FR-008**: La pantalla MUST aceptar coma o punto como separador decimal de las horas.
- **FR-009**: La pantalla MUST impedir un doble envío mientras una carga está en curso.
- **FR-010**: Un fallo al consultar el esfuerzo de una historia MUST NOT impedir ver esa historia ni las demás.
- **FR-011**: Esta historia MUST NOT cambiar el contrato de la API, el backend ni el esquema de la base. Editar o borrar cargas queda fuera de alcance (spec 019, FR-015).
- **FR-012**: La pantalla MUST permitir desplegar, junto al total, el detalle de las cargas de una historia con integrante, fecha, actividad y horas, sin reordenarlas. El detalle MUST empezar plegado y MUST NOT ofrecerse si la historia no tiene cargas.

### Key Entities *(include if feature involves data)*

- **Carga de horas**: la de la spec 019: integrante, fecha, actividad, duración y Sprint.
- **Esfuerzo real**: total acumulado de una historia y detalle de sus cargas (spec 019).
- **Integrante**: persona asignada al proyecto (HU-04); es quien imputa las horas.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un integrante registra las horas de una actividad en menos de 1 minuto y en una sola acción de guardado.
- **SC-002**: El 100% de los rechazos de horas, actividad, fecha e integrante se muestran junto a su campo, conservando lo escrito.
- **SC-003**: Tras guardar una carga, el total de la historia se ve actualizado sin recargar la aplicación.
- **SC-004**: El total mostrado coincide con la suma de las cargas (0,5 h + 1,5 h + 2 h se ve como 4 h).
- **SC-005**: Cada Criterio de Aceptación de esta spec tiene al menos una prueba de componente con el escenario en su nombre.

## Assumptions

- No hay autenticación: elegir un integrante no verifica que quien carga sea esa persona (igual que en la spec 019).
- "Hoy", para la validación, es la fecha del servidor (spec 019). La fecha propuesta en el formulario usa la del navegador; si difieren, decide el backend.
- La pantalla no agrega límites propios al selector de fecha ni repite el tope de horas: el backend es la única fuente de esas reglas (principio DRY del proyecto).
- El formato de horas mostrado usa coma decimal y hasta dos decimales. El dato exacto sigue siendo el de minutos enteros del backend.
- El total de una historia suma todas sus cargas, también las de Sprints anteriores si volvió al backlog y se replanificó (spec 019, FR-013): es esfuerzo de la historia, no del Sprint en cuya tarjeta se muestra.
- El detalle no muestra el Sprint de cada carga: la API solo devuelve su identificador interno, que no es el número que ve la persona (#129). Mostrarlo traducido no lo pide el issue y puede sumarse después.
- Los comportamientos de la API ya están cubiertos con Godog por la spec 019; esta spec agrega pruebas de componente.
- Métricas de desviación y reportes quedan fuera: son de HU-03 y HU-10.
