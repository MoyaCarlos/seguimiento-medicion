# Feature Specification: Datos de prueba v2 (issue #235)

**Feature Branch**: `feature/235-datos-de-prueba-v2`

**Created**: 2026-11-09

**Status**: Borrador, pendiente de revisión.

**Input**: User description: "Datos de prueba v2 (issue #235): ampliar `POST /demo/datos` con usuarios reales de demostración con contraseña conocida (hasheada con el servicio de HU-20), un solo Scrum Master, cuatro Sprints (dos finalizados, uno activo, uno pendiente), un bloque de defectos (HU-08) y la pantalla de ingreso con el botón «Cargar datos de prueba» y la tabla de usuarios de la demo. Sin cambio de esquema."

**Spec base**: [`specs/021-datos-de-prueba`](../021-datos-de-prueba/spec.md) (v1, issue #144). Esta spec **no la redefine**: las reglas de habilitación por `DEMO_HABILITADO` (FR-001 a FR-004), el proyecto de id fijo 25, la transacción única, los bloques independientes, las fechas relativas, la idempotencia y la concurrencia (FR-005 a FR-009, FR-018 a FR-021) siguen vigentes. Aquí solo se listan los cambios. Donde esta spec contradice a la 021, vale esta.

## Clarifications

### Session 2026-11-09

- Decisiones ya cerradas por el Scrum Master (issue #235 y encargo de hoy): usuarios `carlos.moya`, `ana.iriarte`, `candela.vulcano`, `jimena.martinez` con contraseña `demo1234`; un solo Scrum Master; cuatro Sprints; defectos de 4 a 6; 409 con el nombre del usuario si ya existe.
- El esquema no cambia: todo lo que sigue entra en las tablas y restricciones de [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md) (§4.2 `usuarios`, §4.3 `integrantes`, §4.5 `sprints`, §4.6 `historias_arrastradas`, §4.7 votación, §4.9 `defectos`).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Entrar a la demo con usuarios reales (Priority: P1)

Como integrante del equipo o profesor en una demo, quiero que la carga cree usuarios con contraseña conocida y que la pantalla de ingreso me las muestre, para iniciar sesión con cualquiera de ellos sin registrarme a mano.

**Why this priority**: Con el login real (HU-20), los usuarios «Beta» sin contraseña ya no pueden ingresar y, además, no cumplen la regla del nombre de usuario; sin esto la demo no es navegable.

**Independent Test**: Con la función habilitada y la base vacía, se pide la carga y se inicia sesión con `carlos.moya` / `demo1234` y con los otros tres usuarios.

**Acceptance Scenarios**:

1. **Given** la función habilitada y la base vacía, **When** pido la carga, **Then** existen los usuarios `carlos.moya`, `ana.iriarte`, `candela.vulcano` y `jimena.martinez`, todos integrantes del proyecto 25: Carlos como Scrum Master y los otros tres como Product Builder.
2. **Given** la carga exitosa, **When** inicio sesión con cualquiera de los cuatro usuarios y la contraseña `demo1234`, **Then** el sistema me deja ingresar.
3. **Given** la carga exitosa, **When** se inspeccionan los usuarios en la base, **Then** ninguno guarda la contraseña en claro: guardan el hash PBKDF2 de HU-20.
4. **Given** la carga exitosa, **When** consulto el equipo del proyecto 25, **Then** hay exactamente un Scrum Master.
5. **Given** la carga exitosa, **When** se buscan usuarios «Beta Uno» a «Beta Cuatro», **Then** no existen.

---

### User Story 2 - Un usuario de la demo ya existe (Priority: P1)

Como integrante del equipo, si alguno de los usuarios de la demo ya existe (por ejemplo, alguien se registró como `ana.iriarte`), quiero un mensaje que nombre al usuario y que no quede nada a medias, para saber qué resolver sin ensuciar la base.

**Why this priority**: La demo vive en bases compartidas; una carga que pise o duplique usuarios reales sería peor que no cargar.

**Independent Test**: Se crea `ana.iriarte` a mano y se pide la carga: responde 409 con «ya existe el usuario ana.iriarte» y las tablas de la carga quedan idénticas.

**Acceptance Scenarios**:

1. **Given** que existe el usuario `ana.iriarte` y no el proyecto 25, **When** pido la carga, **Then** el sistema responde 409 con el mensaje «ya existe el usuario ana.iriarte» y no cambia nada (ni siquiera los otros tres usuarios).
2. **Given** que existen varios usuarios de la demo, **When** pido la carga, **Then** el mensaje nombra al primero en conflicto, en el orden de la lista de la demo.
3. **Given** que el proyecto 25 de la demo ya existe, **When** pido la carga, **Then** el sistema responde 409 `ya_cargado` como en la v1, aunque sus usuarios existan.
4. **Given** que la carga falla a mitad de camino, **When** termina el pedido, **Then** no queda ningún dato de la carga, incluidos los usuarios.
5. **Given** dos cargas simultáneas, **When** ambas terminan, **Then** una responde 201 y la otra 409, sin duplicados ni datos a medias.

---

### User Story 3 - Historial y métricas con comparación (Priority: P1)

Como integrante del equipo, quiero que el proyecto de demostración tenga dos Sprints finalizados, uno activo y uno pendiente, para que el historial de Sprints y las métricas muestren una comparación real.

**Why this priority**: Con un solo Sprint finalizado no hay comparación de velocidad ni historial que mostrar.

**Independent Test**: Tras la carga, se consultan los Sprints, el historial y el dashboard del proyecto 25.

**Acceptance Scenarios**:

1. **Given** la carga exitosa, **When** consulto los Sprints, **Then** el #1 está Finalizado (hoy − 28 a hoy − 15), el #2 Finalizado (hoy − 14 a hoy − 8), el #3 Activo (hoy − 7 a hoy + 6) y el #4 Pendiente sin Sprint Goal ni fechas.
2. **Given** la carga exitosa, **When** consulto el proyecto 25, **Then** su fecha de inicio es hoy − 30 y la de fin hoy + 45.
3. **Given** la carga exitosa, **When** consulto las historias, **Then** hay 12 historias repartidas así: Sprint 1, 3 Completadas; Sprint 2, 2 Completadas; Sprint 3, 2 En progreso y 1 Nueva; Sprint 4, 1 Nueva; Product Backlog sin Sprint, 3 Nuevas (al menos una sin estimar).
4. **Given** la carga exitosa, **When** comparo los Sprints finalizados, **Then** sus Story Points completados difieren (el Sprint 1 completa 13 y el 2 completa 8), de modo que la velocidad no es plana.
5. **Given** la carga exitosa, **When** consulto la historia que no se completó en el Sprint 2, **Then** volvió al Product Backlog en estado Nueva y tiene su registro de arrastrada asociado al Sprint 2.
6. **Given** la carga exitosa, **When** consulto las horas, **Then** hay registros de los cuatro integrantes en los Sprints 1, 2 y 3, con fechas dentro de cada Sprint y hasta hoy, y para las historias Completadas estimadas el total imputado queda entre el 50 % y el 150 % de «Story Points × horas por Story Point».
7. **Given** la carga exitosa, **When** consulto la votación, **Then** hay una ronda revelada con los cuatro votos en una historia sin estimar del backlog y una ronda abierta con votos parciales (dos de cuatro) en otra.

---

### User Story 4 - Defectos de ejemplo (Priority: P2)

Como integrante del equipo, quiero defectos pendientes y resueltos con distintas severidades, ligados a historias de los Sprints finalizados y del activo, para mostrar el registro y la métrica de defectos (HU-08).

**Why this priority**: Es la tabla que la v1 dejó afuera (FR-022 de la 021) y que hoy ya existe.

**Independent Test**: Tras la carga, se consultan los defectos del proyecto 25 y el resumen.

**Acceptance Scenarios**:

1. **Given** la carga exitosa, **When** consulto los defectos, **Then** hay 5, con severidades Alta, Media y Baja, de los cuales 3 están Resueltos y 2 Pendientes.
2. **Given** la carga exitosa, **When** miro cada defecto, **Then** está ligado a una historia del proyecto, y su Sprint de detección es uno de los Sprints 1, 2 o 3.
3. **Given** un defecto Resuelto, **When** lo consulto, **Then** tiene Sprint de resolución, nunca anterior a su Sprint de detección; un Pendiente no lo tiene.
4. **Given** la carga exitosa, **When** consulto el resumen de la respuesta, **Then** incluye el conteo de defectos creados.

---

### User Story 5 - Botón y tabla de usuarios en la pantalla de ingreso (Priority: P2)

Como persona que llega a `/ingresar` sin cuenta, quiero cargar los datos de prueba desde ahí y ver con qué usuarios y contraseña entrar, para empezar la demo sin leer el README.

**Why this priority**: Es la cara visible; la carga es usable por API sin ella.

**Independent Test**: Pruebas de componente con `GET /demo/estado` habilitado o no y cada respuesta de la carga simulada.

**Acceptance Scenarios**:

1. **Given** que `GET /demo/estado` responde 200, **When** se muestra `/ingresar`, **Then** aparece el botón «Cargar datos de prueba».
2. **Given** que `GET /demo/estado` responde 404 o falla, **When** se muestra `/ingresar`, **Then** no aparece el botón ni la tabla.
3. **Given** el botón visible, **When** la carga responde 201, **Then** se muestra el aviso de éxito y una tabla con usuario, contraseña y rol de los cuatro usuarios de la demo.
4. **Given** el botón visible, **When** la carga responde 409 `ya_cargado`, **Then** se muestra un aviso (no un error) y la misma tabla.
5. **Given** el botón visible, **When** la carga responde 409 por usuario existente, **Then** se muestra un aviso claro que nombra al usuario, sin tabla de éxito.
6. **Given** el botón visible, **When** la carga responde 409 por id ocupado u otro error, **Then** se muestra un mensaje de error.
7. **Given** una carga en curso, **When** miro el botón, **Then** está deshabilitado.

---

### Edge Cases

- El usuario en conflicto difiere solo en mayúsculas (`Ana.Iriarte`): el nombre de usuario se compara normalizado, así que también es conflicto.
- Existe el proyecto 25 de la demo y además uno de los usuarios: gana `ya_cargado`.
- Existe un proyecto 25 con otro nombre y no hay conflicto de usuarios: se mantiene el 409 de id ocupado de la v1.
- Una base con los usuarios «Beta» de la v1: no es conflicto (nombres distintos); los usuarios Beta quedan como están.
- Cambia el día entre cargas en bases distintas: todas las fechas, incluidas las de los defectos, son relativas al día de la carga.
- Hashear cuatro contraseñas es costoso: no debe agotar el tiempo de espera de la carga ni mantener tomada la escritura de la base más de lo necesario.

## Requirements *(mandatory)*

### Functional Requirements

Las reglas se redactan en formato EARS.

**Usuarios y equipo (Ubicuo / No deseado)**

- **FR-001**: El sistema deberá crear los usuarios `carlos.moya`, `ana.iriarte`, `candela.vulcano` y `jimena.martinez`, cada uno con la contraseña de demostración `demo1234`, y los reemplaza a los usuarios «Beta Uno» a «Beta Cuatro» de la v1.
- **FR-002**: El sistema deberá guardar cada contraseña hasheada con el servicio PBKDF2 de HU-20, recibido por el cargador como dependencia (no instanciado dentro del bloque), nunca en claro.
- **FR-003**: El sistema deberá asignar a `carlos.moya` como Scrum Master y a los otros tres como Product Builder del proyecto 25, con un único Scrum Master.
- **FR-004**: El sistema deberá exigir que los nombres de usuario de la demo cumplan la regla de nombre de HU-20 (3 a 30 caracteres: letras sin tildes, números, punto, guion y guion bajo).
- **FR-005**: Si alguno de los usuarios de la demo ya existe (comparando el nombre normalizado), el sistema responderá 409 con el mensaje «ya existe el usuario X», donde X es el primero en conflicto en el orden de la lista de la demo, y deshará toda la carga.
- **FR-006**: Si el proyecto de demostración ya existe, el sistema responderá 409 `ya_cargado` como en la v1, incluso si también hay usuarios en conflicto.
- **FR-007**: El sistema no deberá reutilizar ni modificar usuarios preexistentes: la carga v2 crea sus cuatro usuarios o falla.

**Sprints, historias, horas y votación (Ubicuo)**

- **FR-008**: El sistema deberá crear el proyecto 25 con inicio en hoy − 30 y fin en hoy + 45, y cuatro Sprints: #1 Finalizado (hoy − 28 a hoy − 15), #2 Finalizado (hoy − 14 a hoy − 8), #3 Activo (hoy − 7 a hoy + 6) y #4 Pendiente sin Sprint Goal ni fechas (CHECK de §4.5).
- **FR-009**: El sistema deberá crear 12 historias repartidas como en la User Story 3, escenario 3: las 10 de la v1 reubicadas y dos nuevas. Toda historia Completada tiene fecha de finalización dentro de su Sprint y las estimaciones usan la escala Fibonacci.
- **FR-010**: El sistema deberá dejar una historia sin completar al cerrar el Sprint 2, de vuelta en el Product Backlog en estado Nueva y con su registro en `historias_arrastradas` asociado al Sprint 2.
- **FR-011**: El sistema deberá imputar horas de los cuatro integrantes en historias de los Sprints 1, 2 y 3, con fechas dentro de cada Sprint y hasta hoy; en las Completadas estimadas, el total queda entre el 50 % y el 150 % de «Story Points × horas por Story Point».
- **FR-012**: El sistema deberá crear una ronda de votación revelada con un voto de cada uno de los cuatro integrantes en una historia sin estimar del backlog, y una ronda abierta con dos votos en otra historia sin estimar.

**Defectos (Ubicuo)**

- **FR-013**: El sistema deberá crear 5 defectos del dominio de licencias de conducir, 3 Resueltos y 2 Pendientes, con severidades Alta, Media y Baja, ligados a historias de los Sprints 1, 2 y 3.
- **FR-014**: El sistema deberá registrar para cada defecto su Sprint de detección (el Sprint 1, 2 o 3 en que se detectó) y, si está Resuelto, un Sprint de resolución igual o posterior al de detección; un Pendiente no lleva Sprint de resolución (CHECK de §4.9).
- **FR-015**: El sistema deberá aplicar las validaciones del dominio de HU-08 al crear cada defecto y respetar las FK compuestas, sin cambiar el esquema.
- **FR-016**: El sistema deberá sumar el conteo de defectos creados al resumen de la respuesta 201 (`defectos`).

**Idempotencia, atomicidad y concurrencia (se mantienen de la 021)**

- **FR-017**: El sistema deberá seguir haciendo toda la carga en una sola transacción, con la idempotencia dada por la clave primaria del proyecto 25 y los usuarios, y la concurrencia resuelta por la transacción: de dos cargas simultáneas, una responde 201 y la otra 409.
- **FR-018**: Si cualquier escritura de la carga falla, el sistema deberá deshacer todo, incluidos los usuarios, y responder con un error.

**Frontend (Basado en evento / Basado en estado)**

- **FR-019**: Mientras `GET /demo/estado` responda 200, la pantalla `/ingresar` mostrará el botón «Cargar datos de prueba», reutilizando el componente existente; en cualquier otro caso no lo mostrará.
- **FR-020**: Cuando la carga responda 201 o 409 `ya_cargado`, la interfaz mostrará una tabla con usuario, contraseña y rol de los cuatro usuarios de la demo.
- **FR-021**: Cuando la carga responda el 409 de usuario existente, la interfaz mostrará un aviso claro con el nombre del usuario, no un error genérico.
- **FR-022**: Mientras la carga esté en curso, el botón estará deshabilitado; con la función apagada (`DEMO_HABILITADO` sin «true») no se muestra.
- **FR-023**: El README deberá documentar los usuarios y la contraseña de demostración.

### Key Entities

- **Usuarios de la demo**: cuatro usuarios con contraseña conocida y rol en el proyecto 25; reemplazan a los «Beta».
- **Proyecto de demostración**: se mantiene (id fijo 25), ahora con cuatro Sprints.
- **Defecto de demostración**: nuevo bloque de carga, ligado a una historia y a los Sprints de detección y resolución.

Las tablas y restricciones están en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md); esta funcionalidad no agrega tablas ni columnas.

## Impacto en la carga actual

Qué asume hoy tres Sprints o los usuarios «Beta» (rutas relativas a la raíz; líneas de `origin/main` al 2026-11-09).

**Código de la carga**

| Lugar | Qué asume | Cambio |
|---|---|---|
| `internal/demo/bloque_usuarios.go:11` | `nombresUsuarios` con los cuatro «Beta» | Lista nueva con nombre y contraseña, hash inyectado |
| `internal/demo/bloque_usuarios.go:16-39` | `INSERT OR IGNORE` sin contraseña; reutiliza el usuario si existe | Insertar con `contrasena_hash`; si existe, error de conflicto con el nombre |
| `internal/demo/bloque_proyecto.go:23-27` | `equipo` con tres integrantes Beta (un solo SM) | Cuatro integrantes con los usuarios nuevos |
| `internal/demo/bloque_proyecto.go:32` | Inicio del proyecto `c.dia(-15)` | `c.dia(-30)` |
| `internal/demo/bloque_proyecto.go:13-14` | Id 25 y nombre reservado | Sin cambio |
| `internal/demo/bloque_sprints.go:12-31` | Tres Sprints fijos (Finalizado, Activo, Pendiente), fechas −14/−8 y −7/+6 | Cuatro Sprints; el Activo conserva −7/+6 |
| `internal/demo/bloque_historias.go:24` | Campo `sprint` con rango «1 a 3» | Rango 1 a 4 |
| `internal/demo/bloque_historias.go:30-61` | Tabla `historiasDemo` con 10 historias y su reparto por Sprint (`sprint: 1/2/3`, `finalizada: -9..-11`) | 12 historias; reubicadas en los Sprints 1 a 4 con fechas de finalización nuevas |
| `internal/demo/bloque_historias.go:78-83` | La arrastrada se registra en `c.sprints[0]` | Registrarla en `c.sprints[1]` (Sprint 2) |
| `internal/demo/bloque_historias.go:106` | `c.sprints[h.sprint-1]` | Sigue valiendo con 4 Sprints |
| `internal/demo/bloque_horas.go:14,19-37` | `sprint` 1 o 2; 15 registros de «Beta» con fechas −13 a −1 | Registros de los cuatro integrantes en los Sprints 1 a 3 |
| `internal/demo/bloque_horas.go:42-45` | Solo los Sprints `c.sprints[0]` (−14/−8) y `[1]` (−7/+6) | Tres Sprints con sus fechas nuevas |
| `internal/demo/bloque_votacion.go:16-22` | Tres votos «Beta» y `historiaVotada = "aviso"` | Cuatro votos nuevos, más la ronda abierta con dos |
| `internal/demo/carga.go:25-34` | `Resumen` sin defectos | Campo `Defectos` (`json:"defectos"`) |
| `internal/demo/carga.go:87-97` | `bloquesV1` sin bloque de defectos | Agregar `bloqueDefectos` al final (hereda de usuarios, sprints e historias) |
| `internal/demo/carga.go:105-107` | `NewCargador(db, ahora)` sin hasheador | Recibe el hasheador (HU-20) |
| `internal/demo/carga.go:109-143` | `Cargar`/`clasificarConflicto` solo distinguen id del proyecto (`ya_cargado`/id ocupado) | Distinguir además el usuario en conflicto; el proyecto gana |
| `internal/http/demo_handler.go:56-64` | Mapea `ErrYaCargado` y `ErrIdOcupado` a 409 | Mapear también el error de usuario existente |
| `cmd/api/main.go:119` | `demo.NewCargador(db, time.Now)` | Pasar el servicio de contraseñas de producción |

**Pruebas**

| Lugar | Qué asume |
|---|---|
| `internal/demo/bloque_usuarios_test.go:5-38` | Cuatro usuarios «Beta», reutilización de `Beta Dos` con id 1 |
| `internal/demo/bloque_proyecto_test.go:31-46` | Roles de Beta Uno/Dos/Tres, Beta Cuatro sin proyecto, 3 integrantes |
| `internal/demo/bloque_sprints_test.go:5-36` | 3 Sprints y su estado por índice |
| `internal/demo/bloque_historias_test.go:12-89` | 10 historias y reparto por Sprint (líneas 15-27), 3 fechas de finalización en el Sprint 1 (línea 69), arrastrada en `c.sprints[0]` (línea 79) |
| `internal/demo/bloque_horas_test.go:8-37` | 15 registros, solo Sprints 1 y 2, Beta Uno/Dos/Tres y Beta Cuatro sin horas |
| `internal/demo/bloque_votacion_test.go:14-22` | Votos de tres «Beta»; 1 ronda y 3 votos |
| `internal/demo/carga_test.go:18,25,29,41` | `Resumen` esperado y su JSON (4 usuarios, 3 integrantes, 3 Sprints, 10 historias, 15 horas, 1 ronda, 3 votos), lista de tablas propias |
| `internal/demo/concurrencia_test.go:55` | Conteos esperados por tabla |
| `internal/demo/idempotencia_test.go:11,66` | Lista de tablas de la carga, mensaje «usuarios Beta» (hay que sumar `defectos`) |
| `internal/demo/atomicidad_test.go:18` | Usa `bloquesV1` (sin cambio de lógica) |
| `internal/demo/ayudas_test.go` | Ayudas compartidas: revisar si arman los usuarios |
| `features/datos_prueba.feature:14-15` y `features/steps_datos_prueba.go:19,51,120-121,183` | «10 historias, 3 Sprints y 3 integrantes»; `NewCargador(p.db, time.Now)`; lista de tablas |
| `web/src/compartido/demo/BotonDatosDePrueba.test.tsx:33-111` | Ocho pruebas con el 409 de id ocupado y el aviso sin tabla |
| `web/src/compartido/demo/api.ts:3-14` | `PROYECTO_DE_PRUEBA = 25` y el tipo de respuesta sin defectos |
| `web/src/features/proyecto/crear.tsx:5,124` | Monta `<BotonDatosDePrueba />` (único uso hoy) |
| `web/src/features/autenticacion/Ingresar.tsx:6-27` | Sin botón ni tabla de la demo |
| `web/src/features/votacion/Integrantes.test.tsx:19-22` | Usa «Beta Uno»/«Beta Dos» solo como datos de prueba del componente; **no** depende de la carga y no se toca |
| `README.md:64-66` | Describe el reinicio de la base y `DEMO_HABILITADO`, sin contraseñas |

**Dimensión del cambio**: 8 archivos de carga y 1 de handler, 1 de wiring, unas 9 pruebas de `internal/demo` más el escenario Godog, 2 archivos del frontend (botón y pantalla de ingreso) y el README. El bloque nuevo (defectos) y su prueba son los únicos archivos agregados en el backend.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Con la función habilitada, una persona sin cuenta carga la demo desde `/ingresar` y entra con un usuario de la tabla en menos de un minuto.
- **SC-002**: El 100 % de las cargas con un usuario de la demo existente responde 409 nombrando al usuario y deja la base sin cambios.
- **SC-003**: El 100 % de las cargas repetidas o con el id ocupado termina sin modificar la base.
- **SC-004**: En el 100 % de las cargas fallidas no queda ningún dato, tampoco usuarios.
- **SC-005**: Con dos cargas simultáneas, exactamente una se completa.
- **SC-006**: Tras la carga, el historial muestra dos Sprints finalizados con velocidades distintas y el registro de defectos tiene casos Pendientes y Resueltos de las tres severidades.
- **SC-007**: La carga completa termina en tiempo razonable aun con el hashing de cuatro contraseñas (a validar en el plan con el costo real de PBKDF2).

## Assumptions

- Sin cambio de esquema ni de `MODELO-DE-DATOS.md`. Si surge una duda de modelo, se frena y se pregunta al Scrum Master.
- La ronda abierta con votos parciales figura en el issue #235 y se incluye; entra en el esquema actual (`revelada_en NULL`, §4.7).
- Las contraseñas de la demo son públicas y se documentan en el README y en la pantalla de ingreso; la URL de Render es pública, y el endpoint sigue apagado salvo `DEMO_HABILITADO=true`.
- La tabla de usuarios del frontend es una constante del frontend (datos públicos de demostración), no se pide a la API.
- El botón sigue montado en la pantalla de creación de proyectos (`crear.tsx`) además de `/ingresar`; la tabla vive dentro del mismo componente, así aparece en ambos lugares.
- El hasheo (costoso) se hace antes de abrir la transacción de escritura; el cómo se decide en el plan.
- Se mantiene el id fijo 25 y el nombre reservado de la v1.
- No hay cambio de rutas, ni de `GET /demo/estado`.
