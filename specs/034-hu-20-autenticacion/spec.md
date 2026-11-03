# Feature Specification: Autenticación de usuarios (HU-20)

**Feature Branch**: `feature/HU-20-autenticacion`

**Issue**: #233 (sub-issues #238 a #242, 5 SP)

**Created**: 2026-11-03

**Status**: Draft

**Dependencia**: usa la tabla `usuarios` ya existente (HU-04). No depende de ninguna
otra historia abierta. La migración 007 está reservada para #236 y no se usa aquí.

**Input**: User description: "HU-20: Autenticación de usuarios. Primera de dos
historias del login: registro, ingreso y sesión SIN imponer permisos (los permisos por
rol son HU-21, Sprint 5). Ningún endpoint actual pasa a exigir sesión en este PR: lo
existente y los escenarios de Godog siguen funcionando. Backend primero; el frontend va
después, en la misma spec."

## Clarifications

### Session 2026-11-03

- Q: ¿Qué pasa con la forma actual de crear usuarios sin contraseña (HU-04) junto con el nuevo registro? → A: Sigue exactamente igual; el registro es un camino nuevo e independiente, y los usuarios creados por la vía actual no pueden ingresar.

## Objetivo

Que una persona pueda crear una cuenta con nombre de usuario y contraseña, ingresar,
mantener su sesión abierta, saber quién es el usuario con sesión y cerrar la sesión.
Esta historia solo establece la identidad: **no restringe** ninguna operación existente.

## Fuera de alcance

- Permisos por rol (HU-21, Sprint 5): ningún endpoint actual exige sesión.
- Correo electrónico, «recuperar contraseña», verificación por correo e ingreso con
  terceros.
- Cambio de contraseña, bloqueo por intentos fallidos y administración de cuentas.
- Habilitar el ingreso de los usuarios que ya existen sin contraseña.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Registrarse (Priority: P1)

Una persona nueva elige un nombre de usuario y una contraseña y crea su cuenta. Al
terminar queda con la sesión iniciada, sin tener que ingresar de nuevo.

**Why this priority**: sin cuentas no hay nada más que autenticar; es la puerta de entrada.

**Independent Test**: registrar un usuario nuevo y verificar que se recibe una sesión
válida con la que se consulta quién es el usuario.

**Acceptance Scenarios**:

1. **Given** que no existe el usuario `ana.perez`, **When** se registra con ese nombre y
   una contraseña de 8 o más caracteres, **Then** la cuenta se crea, la persona queda con
   la sesión iniciada y la respuesta no contiene la contraseña.
2. **Given** que ya existe el usuario `ana.perez`, **When** otra persona se registra con
   `Ana.Perez`, **Then** se rechaza con el mensaje «Ese nombre de usuario ya existe» y no
   se crea nada.
3. **Given** un registro con una contraseña de 7 caracteres, **When** se envía, **Then**
   se rechaza indicando que el campo `contrasena` es inválido (mínimo 8 caracteres) y no
   se crea la cuenta.
4. **Given** un nombre de usuario de 2 caracteres, de 31 caracteres o con caracteres no
   permitidos (espacios, `@`, tildes), **When** se registra, **Then** se rechaza
   indicando el campo `nombre`.

---

### User Story 2 - Ingresar (Priority: P1)

Una persona con cuenta ingresa con su nombre de usuario y su contraseña y obtiene una
sesión.

**Why this priority**: es el uso diario; sin ingreso, la cuenta creada no sirve.

**Independent Test**: registrar un usuario, cerrar su sesión e ingresar de nuevo con las
mismas credenciales.

**Acceptance Scenarios**:

1. **Given** una cuenta existente, **When** se ingresa con el nombre (en cualquier
   combinación de mayúsculas y minúsculas) y la contraseña correctos, **Then** se obtiene
   una sesión válida.
2. **Given** una cuenta existente, **When** se ingresa con contraseña incorrecta, **Then**
   se rechaza con «Usuario o contraseña incorrectos».
3. **Given** un nombre de usuario que no existe, **When** se intenta ingresar, **Then** se
   rechaza con exactamente el mismo mensaje y la misma forma de respuesta que en el caso
   anterior, de modo que no se pueda saber cuál de los dos datos falló.
4. **Given** un usuario previo a esta historia (sin contraseña), **When** se intenta
   ingresar con él, **Then** se rechaza con «Usuario o contraseña incorrectos».

---

### User Story 3 - Sesión vigente y vencida (Priority: P2)

Mientras la sesión es válida, el sistema sabe quién es la persona. Cuando la sesión
venció o no es válida, el sistema lo detecta, pide ingresar de nuevo y, tras ingresar, la
persona vuelve a la pantalla en la que estaba.

**Why this priority**: da sentido a la sesión, pero se puede probar solo con US1 y US2.

**Independent Test**: con una sesión vigente consultar «quién soy»; con una sesión
vencida, inexistente o con token alterado, comprobar que se rechaza.

**Acceptance Scenarios**:

1. **Given** una sesión vigente, **When** se consulta «quién soy», **Then** se devuelve el
   identificador y el nombre del usuario, sin contraseña.
2. **Given** una sesión con más de 7 días desde su creación, **When** se consulta «quién
   soy», **Then** se rechaza como no autenticada.
3. **Given** una petición sin token, con un token inexistente o con un token mal formado,
   **When** se consulta «quién soy», **Then** se rechaza como no autenticada con el mismo
   mensaje en los tres casos.
4. **Given** una persona con sesión vencida que estaba en una pantalla protegida,
   **When** el sistema detecta el vencimiento, **Then** muestra el aviso «Tu sesión
   venció», la lleva a ingresar y, al ingresar, la devuelve a esa pantalla.

---

### User Story 4 - Cerrar sesión (Priority: P2)

La persona cierra su sesión desde el menú de usuario. A partir de ese momento el token
usado deja de ser válido.

**Why this priority**: necesario para equipos compartidos; no bloquea el resto.

**Independent Test**: cerrar la sesión y comprobar que el mismo token es rechazado en
«quién soy».

**Acceptance Scenarios**:

1. **Given** una sesión vigente, **When** se cierra la sesión, **Then** el mismo token
   deja de ser válido y «quién soy» lo rechaza.
2. **Given** una persona con dos sesiones abiertas (dos dispositivos), **When** cierra una,
   **Then** la otra sigue siendo válida.
3. **Given** un token ya inválido, **When** se intenta cerrar la sesión, **Then** se
   rechaza como no autenticada sin efectos.

---

### User Story 5 - Pantallas y navegación (Priority: P2, frontend)

La persona ve pantallas de ingreso y de registro, un menú de usuario con «Cerrar sesión»
y rutas protegidas que, sin sesión, la envían a ingresar.

**Why this priority**: es la cara visible; se entrega después del backend.

**Independent Test**: en la interfaz, registrarse, ver el nombre en el menú, cerrar la
sesión y comprobar que una ruta protegida redirige a ingresar.

**Acceptance Scenarios**:

1. **Given** la pantalla de registro, **When** se completa con datos válidos, **Then** la
   persona queda dentro, con su nombre en el menú de usuario.
2. **Given** la pantalla de registro o de ingreso, **When** el servidor devuelve un error
   (nombre repetido, contraseña corta, credenciales incorrectas), **Then** el mensaje se
   muestra junto al campo correspondiente o en el formulario, sin borrar lo escrito.
3. **Given** una persona sin sesión, **When** abre una ruta protegida, **Then** va a
   ingresar y, tras ingresar, vuelve a la ruta que pidió.
4. **Given** una persona con sesión, **When** elige «Cerrar sesión» en el menú, **Then**
   la sesión se cierra y se la lleva a ingresar.
5. **Given** que la interfaz recibe «no autenticado» en cualquier petición, **When**
   ocurre, **Then** descarta la sesión local, muestra «Tu sesión venció» y lleva a
   ingresar conservando la pantalla de origen.

### Edge Cases

- Nombre de usuario con espacios al inicio o al final: se rechaza por caracteres no
  permitidos, no se recorta en silencio (el ingreso, en cambio, normaliza el nombre antes de buscarlo).
- Dos registros simultáneos con el mismo nombre: solo uno se crea; el otro recibe «Ese
  nombre de usuario ya existe».
- Un nombre que coincide, sin distinguir mayúsculas, con un usuario previo sin
  contraseña: se rechaza como repetido (no se «reclama» la cuenta previa).
- Contraseña de exactamente 8 caracteres: válida. Contraseña con espacios, tildes o
  emojis: válida mientras tenga al menos 8 caracteres.
- Contraseña vacía o ausente: se rechaza indicando el campo `contrasena`.
- Cabecera de autorización con esquema distinto de `Bearer` o sin valor: no autenticada.
- Base con usuarios previos al aplicar la migración: la migración termina sin error y esos
  usuarios conservan sus datos y su historial.
- La contraseña nunca aparece en respuestas, en el registro de la aplicación ni en los
  mensajes de error.
- Endpoints existentes (proyectos, sprints, backlog, votación, horas, defectos): siguen
  funcionando sin token, igual que antes.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir registrar una cuenta con nombre de usuario y
  contraseña, y devolver una sesión iniciada en la misma respuesta.
- **FR-002**: El nombre de usuario MUST tener entre 3 y 30 caracteres y contener solo
  letras, números, punto, guion y guion bajo; la unicidad MUST ser sin distinguir
  mayúsculas ni minúsculas, apoyada en la identidad ya existente del usuario.
- **FR-003**: El sistema MUST conservar el nombre tal como lo escribió la persona para
  mostrarlo, y usar la forma normalizada para identificar y comparar.
- **FR-004**: La contraseña MUST tener al menos 8 caracteres; si no, el registro se
  rechaza indicando el campo `contrasena`.
- **FR-005**: Un nombre ya en uso MUST rechazarse con el mensaje «Ese nombre de usuario
  ya existe» sin crear nada.
- **FR-006**: El sistema MUST permitir ingresar con nombre (sin distinguir mayúsculas) y
  contraseña, y devolver una sesión nueva.
- **FR-007**: Ante usuario inexistente, contraseña incorrecta o usuario sin contraseña,
  el ingreso MUST responder idénticamente «Usuario o contraseña incorrectos», sin revelar
  cuál falló, incluso en la forma y el contenido de la respuesta.
- **FR-008**: La contraseña MUST almacenarse solo como resultado de un hash con sal única
  por usuario; la contraseña original nunca se guarda.
- **FR-009**: La contraseña MUST NOT aparecer en ninguna respuesta, mensaje de error ni
  registro de la aplicación.
- **FR-010**: Cada sesión MUST identificarse por un token aleatorio imposible de
  adivinar; el sistema MUST guardar solo una huella irreversible del token, nunca el
  token en claro.
- **FR-011**: La sesión MUST vencer a los 7 días de su creación; una sesión vencida se
  trata como no válida.
- **FR-012**: El token MUST viajar en la cabecera `Authorization: Bearer <token>`.
- **FR-013**: El sistema MUST ofrecer una consulta «quién soy» que, con sesión válida,
  devuelva el identificador y el nombre del usuario, y sin sesión válida responda «no
  autenticado» con un mensaje único para token ausente, mal formado, inexistente o vencido.
- **FR-014**: El sistema MUST permitir cerrar la sesión; el token usado MUST dejar de ser
  válido de inmediato, sin afectar otras sesiones del mismo usuario.
- **FR-015**: Un usuario MUST poder tener varias sesiones a la vez.
- **FR-016**: Los usuarios existentes antes de esta historia MUST conservarse sin
  contraseña y MUST NOT poder ingresar; la migración MUST completarse sin error aunque
  ya existan usuarios.
- **FR-017**: Ningún endpoint existente MUST pasar a exigir sesión en esta historia; los
  escenarios de Godog actuales MUST seguir en verde sin cambios. En particular, la
  creación de usuarios de HU-04 (solo nombre, sin contraseña) MUST seguir igual; los
  usuarios que crea no pueden ingresar y comparten la unicidad de nombre con el registro.
- **FR-018**: El esquema nuevo (columna de contraseña y tabla de sesiones) MUST entregarse
  en una única migración numerada 006, junto con la actualización de
  `MODELO-DE-DATOS.md` en el mismo PR.
- **FR-019** *(frontend)*: La interfaz MUST ofrecer pantallas de ingreso y de registro que
  muestren los errores del servidor sin perder lo escrito.
- **FR-020** *(frontend)*: La interfaz MUST guardar el token en el almacenamiento local
  del navegador, enviarlo en cada petición y descartarlo al cerrar sesión o al recibir
  «no autenticado».
- **FR-021** *(frontend)*: La interfaz MUST ofrecer un menú de usuario con el nombre de
  la persona y la opción «Cerrar sesión».
- **FR-022** *(frontend)*: Las rutas protegidas MUST redirigir a ingresar cuando no hay
  sesión, y devolver a la persona a la pantalla de origen tras ingresar.
- **FR-023** *(frontend)*: Ante sesión vencida o inválida la interfaz MUST mostrar «Tu
  sesión venció» y llevar a ingresar.

### Key Entities

- **Usuario**: persona identificada por su nombre normalizado (único). Gana un dato
  opcional de contraseña protegida; vacío en los usuarios previos, que por eso no pueden
  ingresar.
- **Sesión**: periodo en que un usuario está autenticado. Atributos: usuario al que
  pertenece, huella del token (única), fecha de creación y fecha de vencimiento (creación
  + 7 días). Un usuario puede tener varias; se invalida al cerrar sesión, sin borrarse.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Una persona nueva completa el registro y queda dentro del sistema en menos
  de 1 minuto y un solo envío de formulario.
- **SC-002**: Una persona con cuenta ingresa en menos de 30 segundos.
- **SC-003**: Los 7 escenarios del issue #233 pasan automatizados (Godog para la API,
  pruebas de componente para la interfaz).
- **SC-004**: En el 100 % de los intentos de ingreso fallidos, la respuesta es idéntica
  sin importar si falló el usuario o la contraseña.
- **SC-005**: Tras cerrar sesión, el 100 % de los usos posteriores de ese token son
  rechazados.
- **SC-006**: Una revisión de respuestas, mensajes de error y registros de la aplicación
  no encuentra la contraseña ni el token en claro en ningún caso.
- **SC-007**: El 100 % de los escenarios de Godog previos a esta historia siguen pasando
  sin modificarlos, y la migración se aplica sin error sobre una base con usuarios
  previos.

## Assumptions

- El token viaja en `Authorization: Bearer` y el frontend lo guarda en el almacenamiento
  local del navegador (decisión confirmada por el Scrum Master).
- Las sesiones se guardan en la base de datos; solo se almacena la huella del token, y el
  vencimiento es fijo a 7 días desde la creación, sin renovación automática.
- El hash de contraseña se calcula con la biblioteca estándar de Go, sin dependencias
  nuevas.
- Migración `006_autenticacion.sql`: columna de contraseña (nula) en `usuarios` y tabla
  `sesiones`; se actualiza `MODELO-DE-DATOS.md` en el mismo PR. La 007 no se usa.
- Un nombre que coincide con un usuario previo sin contraseña se trata como repetido; no
  existe flujo para «reclamar» esas cuentas en esta historia.
- La contraseña admite hasta 128 caracteres (límite fijado en el plan para evitar abusos);
  una más larga se rechaza indicando el campo `contrasena`.
- Backend primero, frontend después, ambos en esta spec y en un solo PR.
- Cualquier duda de modelo que surja en el plan se frena y se consulta al Scrum Master.
