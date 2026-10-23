# Feature Specification: Pantalla de creación de historias de usuario (HU-01, frontend)

**Feature Branch**: `feature/HU-01-frontend`

**Created**: 2026-10-22

**Status**: Implementada, pendiente de revisión por PR (Issue #31)

**Input**: User description: "[HU-01] Frontend (issue #31, parte de #1). Pantalla para crear historias de usuario en el Product Backlog de un proyecto, sobre la base del frontend (#108).

Alcance:
- Ruta /projects/:id/backlog; el proyecto sale de la URL.
- Formulario con título, descripción, prioridad y valor de negocio (opcional).
- Prioridad: se muestran las cuatro opciones MoSCoW con sus etiquetas canónicas (Must have, Should have, Could have, Won't have) y se envía la letra M, S, C o W.
- Valor de negocio opcional: 1, 2, 3, 5, 8, 13 o 21.
- Al guardar: POST /backlog con { proyecto_id, titulo, descripcion, prioridad, valor_negocio }.
- Éxito (201): confirmación con la historia creada (estado "Nueva", sin estimación) y formulario listo para cargar otra.
- Validación: título y descripción obligatorios, con advertencia sin enviar nada; máximo 200 y 2000 caracteres. Los errores del servidor (400 con campo y mensaje, 404 proyecto no encontrado) se muestran junto al campo o como mensaje general.
- Fuera de alcance: listado del Product Backlog (HU-15), edición, estimación y Story Points.

Restricciones: sin dependencias nuevas; código en web/src/features/backlog/; usar el cliente y los componentes de la base compartida; la lógica en funciones puras probadas con vitest.

Usá la carpeta specs/008-hu-01-frontend: la 006 está reservada para #107 y la 007 para la base del frontend."

## Clarifications

### Session 2026-10-22

- Q: Cuando el formulario tiene más de un error de validación a la vez (por ejemplo, título vacío y descripción de más de 2000 caracteres), ¿la pantalla debe mostrar todas las advertencias juntas o solo la primera? → A: Mostrar todas las advertencias de los campos inválidos en el mismo intento.
- Q: Cuando la persona escribe el título o la descripción con espacios en los extremos, ¿la aplicación debe enviarlos al servidor recortados o tal como los escribió? → A: Enviar título y descripción recortados (sin espacios en los extremos).
- Q: Después de una creación exitosa, si el siguiente intento de guardar se frena por validación local (sin llamar al servidor), ¿la confirmación de la historia creada desaparece o se mantiene? → A: Desaparece al presionar "Guardar"; ese intento reemplaza la confirmación con su propio resultado.
- Q: Si la dirección no trae un identificador de proyecto válido, ¿la pantalla debe mostrar el formulario deshabilitado o editable pero impedir el envío? → A: Mostrar el formulario deshabilitado (campos y "Guardar") junto al aviso general.
- Q: Cuando la persona deja el valor de negocio en "Sin valor", ¿el cuerpo del `POST` debe omitir el campo `valor_negocio` o enviarlo como `null`? → A: Omitir `valor_negocio` del cuerpo.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Crear una historia de usuario desde el Product Backlog (Priority: P1)

Como integrante del equipo, abro el Product Backlog de un proyecto existente desde su dirección, completo la historia de usuario y la guardo. La aplicación confirma que la historia fue creada con estado "Nueva" y sin estimación, deja los campos vacíos y me permite cargar la siguiente historia sin salir de la pantalla.

**Why this priority**: Es el objetivo central de la historia: sin un formulario que cree historias, el Product Backlog no se puede alimentar desde la interfaz. Es el incremento mínimo que aporta valor.

**Independent Test**: Con el backend corriendo y un proyecto existente, abrir la dirección del backlog de ese proyecto, completar los campos obligatorios, guardar y comprobar que se hizo una única creación, que la confirmación muestra la historia con estado "Nueva" y sin estimación, y que el formulario quedó listo para otra carga.

**Acceptance Scenarios**:

1. **Given** que abro la dirección del Product Backlog de un proyecto existente, **When** completo título, descripción y prioridad y presiono "Guardar", **Then** la aplicación envía una única creación y muestra una confirmación con los datos de la historia creada, su estado "Nueva" y sin estimación.
2. **Given** que la historia se creó correctamente, **When** miro el formulario, **Then** título, descripción, valor de negocio y prioridad volvieron a su estado inicial, la confirmación sigue visible y el proyecto de la dirección se conserva para cargar otra historia.
3. **Given** que abro o recargo la dirección del backlog de un proyecto, **When** la pantalla termina de cargar, **Then** el proyecto utilizado para crear la historia es el que figura en la dirección, sin depender de datos guardados en otra pantalla.

---

### User Story 2 - Evitar el envío de historias incompletas o demasiado largas (Priority: P1)

Como integrante del equipo, si intento guardar una historia sin título o sin descripción, o con un texto que supera el máximo permitido, la aplicación me lo advierte antes de llamar al servidor y no registra nada.

**Why this priority**: Evita crear registros inválidos y da una respuesta inmediata, sin depender de un viaje al servidor. Además, deja claro qué corregir antes de reintentar.

**Independent Test**: Abrir el formulario, intentar guardar con cada caso inválido y comprobar que aparece la advertencia correspondiente junto al campo y que no se realiza ninguna creación.

**Acceptance Scenarios**:

1. **Given** que el formulario no tiene título, **When** presiono "Guardar", **Then** aparece una advertencia junto al campo título y no se envía ninguna creación.
2. **Given** que el título o la descripción contienen solo espacios en blanco, **When** presiono "Guardar", **Then** se consideran vacíos, se advierte junto al campo correspondiente y no se envía ninguna creación.
3. **Given** que el título supera los 200 caracteres o la descripción supera los 2000, **When** presiono "Guardar", **Then** se advierte junto al campo correspondiente y no se envía ninguna creación.
4. **Given** un título de exactamente 200 caracteres y una descripción de exactamente 2000, **When** presiono "Guardar", **Then** la validación de longitud los acepta y la creación se envía.

---

### User Story 3 - Entender por qué el servidor rechazó la historia (Priority: P2)

Como integrante del equipo, cuando el servidor rechaza la creación, quiero ver el motivo en el lugar correcto —junto al campo que lo originó o como aviso general— para poder corregir y reintentar sin adivinar.

**Why this priority**: Sin esta retroalimentación, un rechazo del servidor deja a la persona sin saber qué corregir. Mejora la experiencia, pero la creación exitosa ya aporta valor por sí sola.

**Independent Test**: Simular cada respuesta de error del servidor y comprobar que el mensaje aparece junto al campo indicado o en el aviso general, que permanece visible hasta el siguiente intento y que la pantalla no se rompe.

**Acceptance Scenarios**:

1. **Given** que el servidor rechaza la creación con un error de un campo del formulario, **When** la pantalla recibe la respuesta, **Then** el mensaje aparece junto a ese campo, sin enviar ni borrar el resto de los datos ingresados.
2. **Given** que el servidor rechaza la creación con un error general (sin campo), **When** la pantalla recibe la respuesta, **Then** el mensaje aparece en un aviso general visible por encima del formulario.
3. **Given** que el proyecto de la dirección no existe, **When** el servidor responde "no encontrado", **Then** la pantalla muestra un mensaje general que indica que el proyecto no fue encontrado y no presenta el error como si fuera de un campo.
4. **Given** que el backend no responde o devuelve una falla inesperada, **When** la creación falla, **Then** la pantalla muestra un aviso general y sigue usable, sin quedar cargando ni en blanco.
5. **Given** un aviso de error visible, **When** presiono "Guardar" de nuevo, **Then** el aviso anterior se reemplaza por el resultado del nuevo intento.

---

### User Story 4 - Elegir la prioridad y el valor de negocio (Priority: P2)

Como integrante del equipo, quiero elegir la prioridad entre las cuatro opciones de MoSCoW y, si corresponde, indicar un valor de negocio de la escala permitida, para que la historia quede clasificada desde su creación.

**Why this priority**: Ambos datos enriquecen el backlog y habilitan su priorización, pero la historia puede crearse sin valor de negocio. La prioridad siempre tiene un valor por defecto válido.

**Independent Test**: Elegir cada una de las cuatro prioridades y cada valor de negocio permitido, guardar y comprobar que la historia creada refleja la opción elegida; guardar sin valor de negocio y comprobar que se crea igualmente.

**Acceptance Scenarios**:

1. **Given** el formulario abierto, **When** despliego la prioridad, **Then** veo las cuatro opciones "Must have", "Should have", "Could have" y "Won't have", con "Must have" seleccionada por defecto.
2. **Given** que elijo una prioridad y guardo una historia válida, **When** la creación se completa, **Then** la historia queda registrada con el código que corresponde a esa opción (M, S, C o W) y la confirmación muestra la opción elegida.
3. **Given** el formulario abierto, **When** reviso el valor de negocio, **Then** puedo dejar "Sin valor" o elegir 1, 2, 3, 5, 8, 13 o 21, y no hay otras opciones.
4. **Given** que dejo el valor de negocio en "Sin valor", **When** guardo una historia válida, **Then** la historia se crea igualmente, la confirmación no muestra valor de negocio y la historia sigue sin estimación.

---

### Edge Cases

Los casos de campos obligatorios vacíos o con solo espacios, de longitud excedida y de varios campos inválidos en el mismo intento están definidos en FR-011, FR-012 y FR-012a y no se repiten acá. Los casos restantes están redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- **WHEN** el título tiene exactamente 200 caracteres y la descripción exactamente 2000, **THEN** la validación **shall** aceptarlos como válidos.
- **IF** la dirección no trae un identificador de proyecto, no son dígitos que representen un entero positivo o el valor supera 9007199254740991 (no es representable sin pérdida de precisión en JavaScript), **THEN** la pantalla **shall** mostrar el formulario deshabilitado (campos y "Guardar") junto a un aviso general, sin permitir cargar datos ni enviar la creación.
- **WHILE** una creación está en curso, **THEN** la pantalla **shall** deshabilitar el botón de guardar y mostrar el indicador de carga, y **shall** ignorar nuevos intentos de envío.
- **WHEN** el servidor responde con un error asociado a un campo conocido del formulario (título, descripción, prioridad o valor de negocio), **THEN** la pantalla **shall** mostrar el mensaje junto a ese campo.
- **IF** el error del servidor se refiere al proyecto o a un campo desconocido, **THEN** la pantalla **shall** mostrarlo en el aviso general.
- **WHEN** el servidor responde "proyecto no encontrado", **THEN** la pantalla **shall** mostrar un aviso general que lo indique, sin borrar los datos ya ingresados.
- **IF** el backend no responde o devuelve una respuesta inesperada, **THEN** la pantalla **shall** mostrar un aviso general y **shall** permanecer usable.
- **WHEN** la creación termina con éxito, **THEN** la pantalla **shall** limpiar el formulario y **shall** mantener la confirmación visible hasta el siguiente intento de guardado; todo nuevo intento, incluso uno detenido por validación local, **shall** reemplazarla por el resultado de ese intento.
- **IF** se crean dos historias con el mismo título, **THEN** la pantalla **shall** permitirlas sin tratarlas como error, porque el título no es único.

## Requirements *(mandatory)*

### Functional Requirements

**Pantalla y contexto del proyecto**

- **FR-001**: La aplicación MUST ofrecer la pantalla de creación de historias en la dirección del Product Backlog del proyecto (`/projects/:id/backlog`) y MUST tomar el identificador del proyecto desde esa dirección.
- **FR-002**: La pantalla MUST conservar el proyecto de la dirección al recargar o abrir la URL directamente, sin depender de datos guardados por otra pantalla.
- **FR-002a**: Un identificador de proyecto válido son dígitos que representan un entero positivo representable sin pérdida de precisión en JavaScript (`Number.isSafeInteger`, hasta 9007199254740991). Cuando la dirección no traiga un identificador válido, la pantalla MUST mostrar un aviso general con el formulario deshabilitado (campos y "Guardar") y MUST NOT enviar la creación.

**Formulario**

- **FR-003**: La pantalla MUST ofrecer los campos título, descripción, prioridad y valor de negocio, con sus etiquetas visibles y asociadas a cada campo.
- **FR-004**: El título MUST ser obligatorio y MUST admitir como máximo 200 caracteres; la descripción MUST ser obligatoria y MUST admitir como máximo 2000 caracteres. Las longitudes MUST contarse sobre el texto recortado en los extremos y por caracteres Unicode, no por unidades UTF-16.
- **FR-005**: La prioridad MUST ofrecerse entre las cuatro opciones MoSCoW con las etiquetas "Must have", "Should have", "Could have" y "Won't have", y MUST tener "Must have" como valor inicial.
- **FR-006**: La prioridad MUST registrarse con los códigos M, S, C y W correspondientes a esas cuatro opciones, respectivamente.
- **FR-007**: El valor de negocio MUST ser opcional y, cuando se informe, MUST limitarse a los valores 1, 2, 3, 5, 8, 13 o 21; la opción inicial MUST ser "Sin valor".

**Envío y confirmación**

- **FR-008**: Al guardar, la pantalla MUST enviar una única creación con el identificador del proyecto de la dirección, el título y la descripción recortados en los extremos, la prioridad elegida y el valor de negocio (campo omitido cuando se eligió "Sin valor").
- **FR-009**: Cuando la creación responde con éxito, la pantalla MUST mostrar una confirmación con la historia creada que incluya el título, la prioridad elegida, el estado "Nueva" y la ausencia de estimación; cuando se haya informado, MUST mostrar también el valor de negocio.
- **FR-010**: Después de una creación exitosa, la pantalla MUST dejar el formulario listo para cargar otra historia, conservando el proyecto de la dirección y volviendo los campos a su estado inicial.
- **FR-010a**: Todo nuevo intento de guardado, incluso uno detenido por validación local sin llamar al servidor, MUST reemplazar la confirmación anterior por el resultado de ese intento.

**Validación antes de enviar**

- **FR-011**: Cuando falte el título o la descripción, o contengan solo espacios en blanco, la pantalla MUST mostrar una advertencia junto al campo y MUST NOT enviar la creación.
- **FR-012**: Cuando el título supere los 200 caracteres o la descripción supere los 2000, la pantalla MUST mostrar una advertencia junto al campo y MUST NOT enviar la creación.
- **FR-012a**: Cuando haya más de un campo inválido en el mismo intento, la pantalla MUST mostrar todas las advertencias de esos campos a la vez, no solo la primera.

**Errores del servidor**

- **FR-013**: Cuando el servidor rechace la creación con un error asociado a un campo conocido, la pantalla MUST mostrar ese mensaje junto al campo correspondiente.
- **FR-014**: Cuando el servidor rechace la creación sin campo, con un campo desconocido o porque el proyecto no existe, la pantalla MUST mostrar el mensaje en un aviso general por encima del formulario.
- **FR-015**: Los mensajes de error MUST permanecer visibles hasta el siguiente intento de guardado y MUST NOT borrar los datos ya ingresados.
- **FR-016**: Mientras la creación está en curso, la pantalla MUST deshabilitar el guardado y mostrar un indicador de carga, y MUST evitar envíos duplicados.
- **FR-017**: Cuando no haya conexión con el backend o la respuesta sea inesperada, la pantalla MUST mostrar un aviso general y MUST permanecer usable.

**Alcance y restricciones de entrega**

- **FR-018**: Los textos de la interfaz MUST estar en español, conservando las etiquetas MoSCoW en su forma canónica (Must have, Should have, Could have, Won't have), y la pantalla MUST NOT incluir el listado del Product Backlog, la edición de historias, la estimación ni los Story Points.
- **FR-019**: La pantalla MUST reutilizar el cliente y las piezas de la base compartida del frontend, MUST NOT agregar dependencias y MUST NOT modificar el comportamiento de las demás pantallas.

### Key Entities *(include if feature involves data)*

- **Historia de Usuario**: unidad del Product Backlog que se crea. Datos visibles en la pantalla: título, descripción, prioridad MoSCoW, valor de negocio opcional; al crearse recibe estado "Nueva" y no tiene estimación.
- **Proyecto**: contenedor del backlog al que pertenece la historia. La pantalla lo identifica por el valor de la dirección y no lo crea ni lo modifica.
- **Error de la API**: rechazo del servidor con un estado, un mensaje y, a veces, el nombre del campo al que se refiere. Determina si el mensaje se muestra junto a un campo o en el aviso general.
- **Opciones del formulario**: las cuatro prioridades MoSCoW con su etiqueta y código, y los valores permitidos para el valor de negocio (incluida la opción de no informarlo).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Una persona puede crear una historia completando solo los tres datos obligatorios en menos de 1 minuto.
- **SC-002**: El 100% de las creaciones con datos válidos muestra una confirmación con estado "Nueva", sin estimación, y deja el formulario listo para otra carga.
- **SC-003**: El 100% de los intentos con título o descripción vacíos, solo espacios o fuera de longitud produce cero llamadas de creación.
- **SC-004**: El 100% de los rechazos del servidor muestra su mensaje en el lugar correcto (junto al campo o en el aviso general), sin perder los datos ingresados.
- **SC-005**: Las cuatro prioridades se muestran con su etiqueta canónica MoSCoW y el 100% de las veces se registra el código correcto (M, S, C o W).
- **SC-006**: El 100% de las aperturas y recargas de la dirección del backlog usa el proyecto indicado en la URL.
- **SC-007**: La pantalla se entrega sin dependencias nuevas y sin cambios de comportamiento en las demás pantallas del frontend.
- **SC-008**: El 100% de las historias creadas sin valor de negocio queda registrado sin ese dato, sin impedir la creación.

## Assumptions

- La base del frontend (#108) está mergeada: existen el cliente de la API, las piezas compartidas (campo, botón, mensaje de error, indicador de carga), el registro de rutas y la dirección `/projects/:id/backlog` con una pantalla provisional. Esta historia reemplaza esa pantalla en la ruta existente y ubica su código en `web/src/features/backlog/`, sin tocar el código compartido.
- Se usa el endpoint existente de creación de historias del backend (`POST /backlog`) con el cuerpo `{ proyecto_id, titulo, descripcion, prioridad, valor_negocio }`; la respuesta exitosa (201) es la historia creada con estado "Nueva" y `estimacion_sp` vacío. El formato de errores `{campo, mensaje}` es el de la base (#108, Issue #21).
- El identificador de proyecto se envía como número entero. Si la dirección no trae un identificador válido, se avisa en general, el formulario queda deshabilitado y no se llama al servidor.
- "Must have" es la prioridad inicial porque el campo siempre debe tener un valor válido y la validación pedida solo alcanza a título y descripción.
- "Sin valor" es la opción inicial del valor de negocio; en ese caso el campo se omite del cuerpo del envío, que el backend ya admite como opcional.
- Los máximos de 200 y 2000 caracteres son inclusivos; se cuentan sobre el texto recortado en los extremos y por caracteres Unicode (puntos de código), no por unidades UTF-16, igual que el backend. Título y descripción se envían recortados. La obligatoriedad se evalúa después de descartar espacios en blanco en los extremos.
- La confirmación de éxito permanece visible hasta el siguiente intento de guardado —aunque ese intento se frene por validación local, que la reemplaza con su resultado— o hasta salir de la pantalla; no se exige que desaparezca al empezar a escribir otra historia.
- Los errores de campo conocidos son los que corresponden a título, descripción, prioridad y valor de negocio; cualquier otro campo (incluido el proyecto) se muestra como aviso general.
- No se consulta ni se muestra el nombre del proyecto: la pantalla solo usa el identificador de la dirección. El aviso de "proyecto no encontrado" surge del rechazo del servidor al crear, porque no se pide el proyecto antes de guardar.
- La pantalla es de una sola página, en español, y no incluye listado, edición, estimación ni Story Points; el listado del backlog corresponde a HU-15.
- Restricciones técnicas de la pantalla: el código vive en `web/src/features/backlog/`; se reutilizan el cliente y los componentes de `web/src/compartido/`; no se agregan dependencias; la validación del formulario y el mapeo de opciones se implementan como funciones puras y se prueban con Vitest. La base ya incluye Testing Library y jsdom: la pantalla se prueba con pruebas de componente que cubren los escenarios de aceptación de esta spec, con el escenario en el nombre de cada prueba, además de las funciones puras para la validación.
- El componente `Field` de la base solo maneja `<input>`; la descripción (`<textarea>`) y la prioridad y el valor de negocio (`<select>`) se implementan dentro de `features/backlog` usando `ErrorMessage` y las clases de estilo de la base, sin modificar el código de `compartido/`.
- La gestión de estado y los estilos se resuelven con lo que ya provee la base, sin incorporar librerías de formularios ni de estilos.
