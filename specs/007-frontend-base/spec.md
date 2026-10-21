# Feature Specification: Base del frontend (Issue #108)

**Feature Branch**: `feature/frontend-base`

**Created**: 2026-10-19

**Status**: Draft

**Input**: User description: "Base del frontend (Issue #108): dejar `web/` lista para que las cuatro pantallas del Sprint 2 (HU-04 #36, HU-01 #31, HU-05 #26 y HU-13 #41) se desarrollen en paralelo, cada una en su carpeta, sin pisarse. Consume solo la API que existe hoy."

## Clarifications

### Session 2026-10-19

- Q: Cuando el backend rechaza algo con un error sin `campo`, ¿dónde lo ve la persona? → A: Bloque de aviso visible arriba del formulario, que se queda hasta el próximo envío.
- Q: ¿Se aprueban React Router y Vitest con Testing Library como dependencias nuevas? → A: Sí, las dos.
- Q: ¿Qué entrega la base de "prototipado"? → A: Solo el código base; los diseños existentes quedan como referencia en `docs/`.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Arrancar `web/` limpio y conectado a la API (Priority: P1)

Como integrante del equipo, quiero clonar el repositorio, instalar lo necesario y levantar la interfaz con dos comandos, y que desde ahí pueda hablar con el backend que corre en mi máquina, sin pedirle a nadie que toque el código de Go. Hoy `web/` solo tiene el contador de ejemplo de la plantilla y el navegador no puede llamar al backend.

**Why this priority**: Sin una interfaz que arranque y llegue a la API, ninguna de las cuatro pantallas puede empezar. Es el bloqueo compartido del Sprint 2.

**Independent Test**: Con el backend corriendo, levantar la interfaz, abrirla en el navegador y comprobar que una llamada a la API de ejemplo (por ejemplo, consultar un proyecto cargado con los `curl` de `docs/`) devuelve datos y no un error de conexión.

**Acceptance Scenarios**:

1. **Given** un clon limpio del repositorio, **When** se instalan las dependencias y se ejecuta la compilación de `web/`, **Then** termina sin errores.
2. **Given** el mismo clon, **When** se ejecuta el lint de `web/`, **Then** termina sin errores ni advertencias.
3. **Given** la interfaz abierta, **When** se la inspecciona, **Then** no queda ningún rastro del contenido de ejemplo de la plantilla (contador, logos, textos de Vite).
4. **Given** el backend corriendo en el puerto 8080 y la interfaz levantada en modo desarrollo, **When** la interfaz pide `/projects/{id}`, **Then** recibe la respuesta del backend sin error de permisos del navegador (CORS).

---

### User Story 2 - Sumar una pantalla sin pisar a nadie (Priority: P1)

Como Product Builder que desarrolla una pantalla (HU-04, HU-01, HU-05 o HU-13), quiero crear mi carpeta, registrar mi pantalla con una sola línea y que aparezca en la navegación, sin editar archivos de otras pantallas, para que mi PR y los de mis compañeras no tengan conflictos difíciles de resolver.

**Why this priority**: Es el motivo de la base: permitir el trabajo en paralelo. Si dos PRs se pisan, se pierde el beneficio.

**Independent Test**: Crear una pantalla de prueba siguiendo el README y verificar que el cambio fuera de su carpeta es exactamente una línea en el registro de rutas y una en el de navegación.

**Acceptance Scenarios**:

1. **Given** la base mergeada, **When** se agrega una pantalla nueva siguiendo `web/README.md`, **Then** el cambio fuera de `web/src/features/<nombre>/` es una línea en el registro de rutas y una línea en el registro de navegación.
2. **Given** dos pantallas agregadas en ramas distintas, **When** se mergean, **Then** el conflicto (si lo hay) se resuelve dejando las dos líneas, sin reescribir código.
3. **Given** una pantalla con ruta `/projects/:id/...`, **When** se abre la URL con un id y se recarga la página, **Then** la pantalla sigue sabiendo de qué proyecto se trata, porque el id está en la URL y no en memoria compartida.
4. **Given** una URL que no corresponde a ninguna pantalla, **When** se abre, **Then** se muestra un aviso de "página no encontrada" y no una pantalla en blanco.

---

### User Story 3 - Llamar a la API y mostrar errores de forma uniforme (Priority: P1)

Como Product Builder, quiero un único cliente de la API que me devuelva los datos cuando todo sale bien y un error ya interpretado (con `campo` y `mensaje`) cuando el backend responde 400, 404 o 409, para no volver a escribir la lectura de errores en cada pantalla.

**Why this priority**: Todas las pantallas llaman a la API y todas deben mostrar los errores de validación igual. Duplicarlo en cuatro lugares rompe la regla DRY.

**Independent Test**: Con las pruebas automáticas del cliente, simular respuestas 200, 400 con campo, 400 sin campo, 404, 409 y caída del servidor, y comprobar el resultado de cada una.

**Acceptance Scenarios**:

1. **Given** una respuesta exitosa del backend, **When** el cliente la procesa, **Then** entrega los datos recibidos.
2. **Given** una respuesta 400 con `{campo, mensaje}`, **When** el cliente la procesa, **Then** entrega un error que incluye el estado, el campo y el mensaje.
3. **Given** una respuesta 400, 404 o 409 sin `campo`, **When** el cliente la procesa, **Then** entrega un error con el estado y el mensaje, y sin campo.
4. **Given** que el backend no responde, **When** el cliente intenta la llamada, **Then** entrega un error de conexión distinguible de los errores de la API.
5. **Given** que el backend responde un error con un cuerpo que no es `{campo, mensaje}` (por ejemplo, un 500 con texto), **When** el cliente lo procesa, **Then** entrega un error con el estado y un mensaje genérico, distinguible del error de conexión.

---

### User Story 4 - Armar formularios con piezas comunes (Priority: P2)

Como Product Builder, quiero un botón, un campo de formulario, un mensaje de error y un indicador de carga ya hechos, para que las pantallas se vean y se comporten igual sin que cada una los reinvente.

**Why this priority**: Evita duplicación y da coherencia, pero cada pantalla podría arrancar sin ellas si hiciera falta.

**Independent Test**: Usar cada pieza en una pantalla de prueba y comprobar su comportamiento.

**Acceptance Scenarios**:

1. **Given** un campo con un error asociado a su nombre, **When** se muestra, **Then** el mensaje de error aparece junto al campo correspondiente.
2. **Given** una operación en curso, **When** se muestra el indicador de carga, **Then** el botón de envío queda deshabilitado hasta que termine.
3. **Given** un campo con etiqueta, **When** se lo usa con un lector de pantalla o se hace clic en la etiqueta, **Then** la etiqueta está asociada al campo.

---

### User Story 5 - Documentación y datos de ejemplo (Priority: P2)

Como integrante del equipo (aunque no conozca Go ni React), quiero un `web/README.md` y un listado de `curl` en `docs/` para cargar un proyecto y un Sprint de ejemplo, para poder probar mi pantalla sin armar los datos a mano.

**Why this priority**: Reduce la barrera de entrada, pero no bloquea el código.

**Independent Test**: Una persona que no escribió la base sigue el README y los `curl` desde cero y llega a ver una pantalla de prueba con datos de ejemplo.

**Acceptance Scenarios**:

1. **Given** una persona sin experiencia en el proyecto, **When** sigue `web/README.md`, **Then** logra correr la interfaz y agregar una pantalla de prueba sin pedir ayuda.
2. **Given** el backend recién levantado, **When** ejecuta los `curl` del listado en orden, **Then** quedan creados un proyecto y un Sprint de ejemplo, y el listado muestra cómo anotar sus ids.

---

### Edge Cases

Redactados en formato EARS.

- **IF** el backend no está corriendo al hacer una llamada, **THEN** el cliente de la API **shall** devolver un error de conexión (FR-009a), y la interfaz **shall** mostrarlo en lugar de quedar cargando indefinidamente.
- **IF** el backend responde con un estado o cuerpo inesperado (por ejemplo, un 500 o un cuerpo que no es JSON), **THEN** el cliente **shall** devolver un error genérico con el estado, sin romper la pantalla.
- **IF** una respuesta de error no trae `campo` (error general, no de un campo del formulario), **THEN** la interfaz **shall** mostrar el mensaje en un bloque de aviso visible arriba del formulario, que permanece hasta el próximo envío.
- **IF** el id del proyecto en la URL no corresponde a ningún proyecto, **THEN** la pantalla **shall** poder mostrar el 404 que devuelve la API usando el mismo mecanismo de errores.
- **IF** se recarga una pantalla que dependía de un id creado en la misma sesión y que no está en la URL (por ejemplo, el Sprint recién creado, porque la API no lista Sprints), **THEN** la pantalla **shall** mostrarse sin ese dato y sin fallar; recuperarlo queda para cuando existan los listados.
- **WHEN** una pantalla futura necesite un listado, **THEN** la base **shall** permitir sumarlo en la misma carpeta de la pantalla, sin rehacer la pantalla ni tocar código compartido.
- **IF** dos pantallas agregan su línea en el mismo lugar de un registro, **THEN** el conflicto **shall** poder resolverse conservando ambas líneas.
- **IF** una pantalla importa código de otra pantalla, **THEN** la revisión del PR **shall** rechazarlo; el código común solo puede venir de `web/src/compartido/`.

## Requirements *(mandatory)*

### Functional Requirements

**Arranque y compilación**

- **FR-001**: El sistema **shall** compilar `web/` sin errores tras instalar las dependencias.
- **FR-002**: El sistema **shall** pasar el lint de `web/` sin errores.
- **FR-003**: El sistema **shall** eliminar de `web/` el contenido de ejemplo de la plantilla (contador, logos, estilos y textos de demostración).

**Conexión con la API**

- **FR-004**: **WHILE** corre en modo desarrollo, la interfaz **shall** reenviar las llamadas a la API al backend en el puerto 8080, sin requerir cambios en el código de Go.
- **FR-005**: El sistema **shall** ofrecer un cliente de la API genérico, con operaciones GET, POST y PUT que envían y leen JSON y aplican el manejo de errores de FR-006 a FR-009b. Las llamadas a cada endpoint (por ejemplo, crear un proyecto o iniciar un Sprint) viven en la carpeta de la pantalla que las usa y se construyen sobre ese cliente; el código compartido no conoce los endpoints ni sus datos.
- **FR-006**: **WHEN** la API responde con éxito, el cliente **shall** entregar los datos recibidos.
- **FR-007**: **WHEN** la API responde 400, 404 o 409 con `{campo, mensaje}`, el cliente **shall** entregar un error con estado, campo y mensaje.
- **FR-008**: **WHEN** la API responde un error sin `campo`, el cliente **shall** entregar un error con estado y mensaje, y campo ausente; la pieza de mensaje de error **shall** mostrarlo en un bloque de aviso arriba del formulario, que permanece hasta el próximo envío.
- **FR-009a**: **IF** no hay conexión con el backend, **THEN** el cliente **shall** entregar un error de conexión, distinguible de los errores de la API.
- **FR-009b**: **IF** el backend responde con un estado de error y un cuerpo que no es `{campo, mensaje}` interpretable (por ejemplo, un 500 o un cuerpo que no es JSON), **THEN** el cliente **shall** entregar un error con ese estado y un mensaje genérico, sin romper la pantalla.
- **FR-010**: El cliente de la API **shall** tener pruebas automáticas que cubran éxito, 400 con campo, 400 sin campo, 404, 409, respuesta de error no interpretable y falta de conexión. Depende de la decisión de FR-026 (librería de pruebas).

**Ruteo y navegación**

- **FR-011**: El sistema **shall** tener un registro de rutas donde cada pantalla ocupa una sola línea.
- **FR-012**: El sistema **shall** tener un registro de navegación donde cada pantalla ocupa una sola línea.
- **FR-013**: El proyecto actual **shall** viajar en la URL (`/projects/:id/...`) y no en un estado compartido entre pantallas.
- **FR-014**: **WHEN** se abre una URL sin pantalla asociada, el sistema **shall** mostrar un aviso de "página no encontrada".
- **FR-015**: Agregar una pantalla nueva **shall** requerir, fuera de su carpeta, solo una línea en el registro de rutas y una en el de navegación.

**Piezas compartidas**

- **FR-016**: El sistema **shall** ofrecer un botón, un campo de formulario, un mensaje de error y un indicador de carga en `web/src/compartido/`.
- **FR-017**: El campo de formulario **shall** asociar su etiqueta al campo y poder mostrar el mensaje de error de su campo.
- **FR-018**: **WHILE** una operación está en curso, el botón de envío **shall** poder mostrarse deshabilitado y el indicador de carga visible.

**Independencia entre pantallas**

- **FR-019**: Cada historia **shall** vivir en su propia carpeta `web/src/features/<nombre>/`.
- **FR-020**: Ninguna carpeta de `features/` **shall** importar código de otra carpeta de `features/`.
- **FR-021**: El código compartido **shall** vivir en `web/src/compartido/` y solo esta base lo modifica; las pantallas lo usan sin editarlo.
- **FR-022**: Solo esta base **shall** agregar dependencias a `package.json`.
- **FR-023**: El panel del proyecto **shall** ser el único punto de contacto entre HU-04 y HU-13, integrado al final del Sprint.

**Documentación**

- **FR-024**: `web/README.md` **shall** explicar cómo instalar, correr, compilar y testear la interfaz, y los pasos para agregar una pantalla nueva.
- **FR-025**: `docs/` **shall** incluir un listado de comandos `curl` para cargar un proyecto y un Sprint de ejemplo, con los datos que la API pide hoy (por ejemplo, `creador` y fechas `AAAA-MM-DD`) y dónde anotar los ids devueltos.

**Decisiones tomadas en `/speckit-clarify`**

- **FR-026**: El sistema **shall** usar React Router para el ruteo y Vitest con Testing Library para las pruebas de la interfaz; el equipo aprobó estas dependencias nuevas fuera del stack (constitución, Principio VI), y solo esta base las agrega.
- **FR-027**: La base **shall** entregar solo el código base de `web/`; los diseños de pantalla existentes (`docs/prototipo-sprint2/`) quedan como referencia fuera de su alcance.

### Key Entities

- **Error de la API**: lo que devuelve el backend cuando algo falla. Tiene un estado (400, 404, 409), un mensaje y, a veces, el nombre del campo del formulario al que se refiere.
- **Registro de rutas**: lista con una línea por pantalla que dice qué dirección abre qué pantalla.
- **Registro de navegación**: lista con una línea por pantalla que dice qué enlace aparece en el menú.
- **Pantalla (feature)**: conjunto de archivos de una historia de usuario, dentro de su propia carpeta.
- **Pieza compartida**: botón, campo, mensaje de error o indicador de carga, reutilizable por todas las pantallas.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Tras instalar las dependencias y compilar `web/`, la compilación termina con 0 errores y el lint con 0 errores.
- **SC-002**: Una persona del equipo que no escribió la base llega, siguiendo solo `web/README.md`, a tener la interfaz corriendo y conectada al backend en menos de 15 minutos.
- **SC-003**: Agregar una pantalla de prueba toca exactamente 2 líneas fuera de su carpeta (1 de rutas, 1 de navegación).
- **SC-004**: Las cuatro pantallas del Sprint 2 pueden desarrollarse en ramas separadas sin que ninguna modifique `web/src/compartido/` ni `package.json`.
- **SC-005**: Los 7 casos del cliente de la API (éxito, 400 con campo, 400 sin campo, 404, 409, respuesta de error no interpretable y falta de conexión) pasan en las pruebas automáticas.
- **SC-006**: Recargar una pantalla con `/projects/:id/...` en la URL conserva el proyecto en el 100% de los casos.
- **SC-007**: Ejecutando los `curl` de `docs/` en orden sobre un backend recién levantado, se crean un proyecto y un Sprint de ejemplo sin errores.

## Assumptions

- El equipo no conoce Go ni necesariamente React: README y `curl` están escritos sin jerga.
- El backend corre en `localhost:8080` durante el desarrollo; el reenvío de llamadas se resuelve del lado de la interfaz (decisión recomendada en el Issue).
- La interfaz consume solo la API existente al 2026-10-14. Quedan fuera `GET /projects` (#109), listado de Sprints (#110), listado del backlog (HU-15) y asignar historias a un Sprint (HU-06); la base no los implementa, pero la estructura permite sumarlos sin rehacer pantallas.
- Límites de la API actual: `POST /projects` pide `creador` (nombre) y fechas `AAAA-MM-DD`; los ids de lo recién creado se pierden al recargar; sin HU-06, cerrar un Sprint no tiene efecto visible.
- El formato de errores (`campo` omitido si no aplica) es el entregado por el Issue #21 (PR #112), ya mergeado.
- Estilos: solo los mínimos necesarios para que las piezas sean usables.
- Fuera de alcance: las pantallas en sí, autenticación, sistema de diseño, cambios en el backend y endpoints nuevos.
- Depende de que #112 esté mergeado; no depende de #107 (modelo de datos). Bloquea #26, #31, #36 y #41.
- Los endpoints que usará cada pantalla son los existentes al 2026-10-14: `POST /backlog` (HU-01), `POST/GET/PUT /projects` y `POST/GET /projects/{id}/members` (HU-04), `POST /sprints`, `POST /sprints/{id}/iniciar` y `POST /sprints/{id}/cerrar` (HU-05) y `GET /projects/{id}/status` (HU-13). Cada pantalla define sus llamadas en su carpeta.
- El ciclo RED/GREEN/REFACTOR aplica al cliente de la API, que es la única lógica de la base. El botón, el campo, el mensaje de error, el indicador de carga y el router no requieren un commit RED por componente.
- La interfaz está en español.
- Decisión: el navegador llega al backend por el proxy de desarrollo de la interfaz, sin habilitar CORS en el backend (recomendado en el Issue #108).
- Si la configuración del lint permite prohibir que una pantalla importe a otra, se evalúa en el plan para reforzar FR-020 sin depender solo de la revisión del PR.
