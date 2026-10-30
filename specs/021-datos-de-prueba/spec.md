# Feature Specification: Cargar datos de prueba (issue #144, versión 1)

**Feature Branch**: `feature/144-datos-de-prueba`

**Created**: 2026-10-30

**Status**: Borrador

**Input**: User description: "Botón «Cargar datos de prueba» (issue #144), versión 1: carga de un proyecto de demostración completo (usuarios, proyecto y equipo, Sprints, historias, horas imputadas y una votación revelada) con los datos que ya existen en `main`, para poder mostrar y probar la aplicación sin cargar todo a mano. Trabajo en curso: cada historia futura que agregue tablas o datos extiende la carga con su propio bloque."

## Clarifications

### Session 2026-10-30

- Q: ¿«Razonables frente a lo estimado» es verificable? → A: Sí: en las historias Completadas estimadas, el total imputado queda entre el 50 % y el 150 % de Story Points × horas por Story Point (FR-016).
- Las demás decisiones (id fijo 25, atomicidad, 404 con la función apagada, Sprint #3 sin goal ni fechas) ya vienen cerradas por el Scrum Master; no quedan ambigüedades críticas.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Cargar el proyecto de demostración (Priority: P1)

Como integrante del equipo (o profesor en una demo), quiero apretar «Cargar datos de prueba» y obtener un proyecto completo y coherente, para mostrar o probar todas las pantallas sin cargar historias, Sprints y horas a mano.

**Why this priority**: Es el objetivo del issue; sin la carga no hay nada que demostrar.

**Independent Test**: Con la función habilitada y la base vacía, se pide la carga y se verifica que existe el proyecto 25 con su equipo, sus tres Sprints, sus historias, sus horas y su votación revelada, y que la respuesta resume las cantidades creadas.

**Acceptance Scenarios**:

1. **Given** la función habilitada y la base sin el proyecto 25, **When** pido la carga, **Then** el sistema responde 201 con el resumen (proyecto 25 y la cantidad de usuarios, integrantes, Sprints, historias, registros de horas y rondas de votación creados) y el proyecto «Licencia de conducir (datos de prueba)» queda disponible.
2. **Given** la carga exitosa, **When** consulto el proyecto 25, **Then** tiene descripción, horas por Story Point y fechas de inicio (hoy − 15 días) y fin (hoy + 45 días).
3. **Given** la carga exitosa, **When** consulto sus integrantes, **Then** Beta Uno es Scrum Master, Beta Dos y Beta Tres son Product Builder, y Beta Cuatro existe como usuario pero no integra el proyecto.
4. **Given** la carga exitosa, **When** consulto sus Sprints, **Then** el #1 está Finalizado (hoy − 14 a hoy − 8), el #2 Activo (hoy − 7 a hoy + 6) y el #3 Pendiente sin Sprint Goal ni fechas.
5. **Given** la carga exitosa, **When** consulto el Product Backlog y los Sprints, **Then** hay unas 10 historias de licencias de conducir con prioridades M, S, C y W, estados Nueva, En progreso y Completada, algunas estimadas (escala Fibonacci) y otras no, algunas con criterios de aceptación y descripciones con Markdown, repartidas así: 3 Completadas en el Sprint 1, 2 En progreso y 1 Nueva en el Sprint 2, 1 Nueva en el Sprint 3 y 3 en el Product Backlog sin Sprint (al menos una sin estimar), y una historia que quedó sin completar al cerrar el Sprint 1 y volvió al backlog con su registro de arrastrada.
6. **Given** la carga exitosa, **When** consulto el esfuerzo real de las historias de los Sprints 1 y 2, **Then** tienen horas imputadas por distintos integrantes, razonables frente a lo estimado.
7. **Given** la carga exitosa, **When** consulto la votación de una historia sin estimar del Product Backlog, **Then** hay una ronda revelada con los votos de Beta Uno, Beta Dos y Beta Tres.

---

### User Story 2 - Cargar dos veces no duplica ni rompe nada (Priority: P1)

Como integrante del equipo, si la carga ya se hizo o el id 25 lo usa otro proyecto, quiero un mensaje claro y que no cambie nada, para no ensuciar la base ni perder datos reales.

**Why this priority**: La demo vive en bases compartidas (Render); una carga repetida o parcial arruinaría la demostración.

**Independent Test**: Se carga una vez (201) y se vuelve a pedir: debe responder 409 y la base queda idéntica. Aparte, con un proyecto 25 de otro nombre, debe responder un 409 distinto.

**Acceptance Scenarios**:

1. **Given** los datos de prueba ya cargados, **When** pido la carga otra vez, **Then** el sistema responde 409 «los datos de prueba ya están cargados» y no modifica nada.
2. **Given** un proyecto con id 25 y otro nombre, **When** pido la carga, **Then** el sistema responde 409 «el id 25 ya lo usa otro proyecto» y no modifica nada.
3. **Given** que la carga falla a mitad de camino, **When** termina el pedido, **Then** no queda ningún dato de la carga (ni usuarios, ni proyecto, ni nada).
4. **Given** dos cargas simultáneas, **When** ambas terminan, **Then** una responde 201 y la otra 409, sin duplicados ni datos a medias.

---

### User Story 3 - Habilitación controlada (Priority: P1)

Como Scrum Master, quiero que la función esté apagada por defecto y solo se active por configuración, para que nadie cargue datos de prueba en un entorno real por accidente.

**Why this priority**: Evita contaminar entornos reales; es una regla de seguridad operativa.

**Independent Test**: Sin la variable, cualquier pedido a `/demo/*` responde 404; con la variable en «true», `GET /demo/estado` responde 200.

**Acceptance Scenarios**:

1. **Given** la función deshabilitada, **When** pido `GET /demo/estado` o `POST /demo/datos`, **Then** el sistema responde 404 como si la ruta no existiera y no escribe nada.
2. **Given** la función habilitada, **When** pido `GET /demo/estado`, **Then** responde 200 con `{"habilitado":true}`.

---

### User Story 4 - Botón en la lista de proyectos (Priority: P2)

Como integrante del equipo, quiero un botón «Cargar datos de prueba» en la lista de proyectos que aparezca solo cuando la función está habilitada, para cargar la demo con un clic y abrirla enseguida.

**Why this priority**: Es la cara visible, pero la carga es usable por API sin ella.

**Independent Test**: Pruebas de componente con el estado habilitado/deshabilitado y cada respuesta de la carga simulada.

**Acceptance Scenarios**:

1. **Given** que `GET /demo/estado` responde 200, **When** se muestra la lista de proyectos, **Then** aparece el botón «Cargar datos de prueba».
2. **Given** que `GET /demo/estado` responde 404 o falla, **When** se muestra la lista de proyectos, **Then** no aparece el botón.
3. **Given** el botón visible, **When** la carga responde 201, **Then** se muestra un aviso de éxito con el enlace «Abrir el proyecto de prueba» a `/projects/25`.
4. **Given** el botón visible, **When** la carga responde 409 «ya cargado», **Then** se muestra un aviso (no un error) con el mensaje, que incluye el enlace al proyecto de prueba.
5. **Given** el botón visible, **When** la carga responde otro error (incluido el 409 de id ocupado), **Then** se muestra un mensaje de error.

---

### Edge Cases

- Alguno de los usuarios «Beta» ya existe: se reutiliza, no se duplica y no falla la carga (el nombre normalizado es único).
- El id 25 está ocupado por un proyecto con el nombre reservado: «ya cargado», aunque le falten partes (la clave primaria es la única fuente de idempotencia; no se completa ni se repara).
- Los usuarios «Beta» ya existen y el proyecto 25 no: la carga procede y reutiliza los usuarios.
- Cambia el día entre cargas en bases distintas: las fechas son siempre relativas al día de la carga.
- Una carga que falla por una violación de CHECK o clave foránea no deja datos.
- La secuencia de ids de proyectos avanza hasta 25: los proyectos creados después reciben ids mayores (consecuencia aceptada del id fijo).

## Requirements *(mandatory)*

### Functional Requirements

Las reglas se redactan en formato EARS.

**Habilitación (Basado en estado / Ubicuo)**

- **FR-001**: Mientras la función no esté habilitada, el sistema responderá 404 a todas las rutas `/demo/*`, como si no existieran, sin escribir nada.
- **FR-002**: El sistema deberá considerar habilitada la función solo cuando la variable de entorno `DEMO_HABILITADO` valga exactamente «true»; por defecto está apagada.
- **FR-003**: El sistema deberá decidir la habilitación en un único punto, de modo que pueda reemplazarse después por un chequeo de perfil sin tocar el resto.
- **FR-004**: Cuando se pida `GET /demo/estado` con la función habilitada, el sistema responderá 200 con `{"habilitado":true}`.

**Carga (Basado en evento)**

- **FR-005**: Cuando se pida `POST /demo/datos` con la función habilitada y el id 25 libre, el sistema creará el proyecto de demostración con id 25 y nombre reservado «Licencia de conducir (datos de prueba)» y responderá 201 con el resumen (`proyecto_id` y cantidades creadas).
- **FR-006**: El sistema deberá hacer toda la carga en una sola transacción: o se crea todo o no se crea nada.
- **FR-007**: El sistema deberá crear cada dato aplicando las validaciones del dominio (constructores) y respetando los CHECK e índices únicos del esquema, sin cambiarlo.
- **FR-008**: El sistema deberá armar la carga en bloques independientes (usuarios, proyecto y equipo, Sprints, historias, horas, votación), cada uno en su propia función, de modo que una historia futura agregue el suyo sin tocar los demás.
- **FR-009**: El sistema deberá calcular todas las fechas relativas al día de la carga.

**Contenido del bloque de usuarios y equipo**

- **FR-010**: El sistema deberá asegurar los usuarios «Beta Uno», «Beta Dos», «Beta Tres» y «Beta Cuatro», sin rol propio, y no duplicará los que ya existan.
- **FR-011**: El sistema deberá crear el proyecto con descripción, horas por Story Point, inicio en hoy − 15 días y fin en hoy + 45 días.
- **FR-012**: El sistema deberá asignar a Beta Uno como Scrum Master y a Beta Dos y Beta Tres como Product Builder, y dejar a Beta Cuatro sin proyecto.

**Contenido de Sprints, historias, horas y votación**

- **FR-013**: El sistema deberá crear los Sprints #1 Finalizado (hoy − 14 a hoy − 8), #2 Activo (hoy − 7 a hoy + 6) y #3 Pendiente sin Sprint Goal ni fechas (por el CHECK de `MODELO-DE-DATOS.md` §4.4).
- **FR-014**: El sistema deberá crear unas 10 historias del dominio de licencias de conducir con la variedad y el reparto de la User Story 1, escenario 5, donde toda historia Completada tenga fecha de finalización y las estimaciones usen la escala Fibonacci.
- **FR-015**: El sistema deberá registrar la historia sin completar del Sprint 1 como arrastrada y dejarla en el Product Backlog en estado Nueva, con su registro en `historias_arrastradas`.
- **FR-016**: El sistema deberá imputar horas a historias de los Sprints 1 y 2, con fechas dentro de cada Sprint y hasta hoy, repartidas entre los integrantes del proyecto; para las historias Completadas estimadas, el total imputado queda entre el 50 % y el 150 % de «Story Points × horas por Story Point».
- **FR-017**: El sistema deberá crear una ronda de votación revelada en una historia sin estimar del Product Backlog, con un voto de cada integrante con rol (Beta Uno, Beta Dos y Beta Tres).

**Idempotencia y errores (No deseado)**

- **FR-018**: Si el id 25 ya lo usa el proyecto con el nombre reservado, el sistema responderá 409 con el mensaje «los datos de prueba ya están cargados» y no cambiará nada.
- **FR-019**: Si el id 25 ya lo usa un proyecto con otro nombre, el sistema responderá 409 con el mensaje «el id 25 ya lo usa otro proyecto» y no cambiará nada.
- **FR-020**: Si dos cargas llegan a la vez, el sistema deberá completar una (201) y rechazar la otra (409), sin duplicados ni datos a medias.
- **FR-021**: Si cualquier escritura de la carga falla, el sistema deberá deshacer todo y responder con un error, sin dejar datos parciales.
- **FR-022**: El sistema no deberá cargar defectos, métricas ni datos de HU-03, HU-08, HU-09 o HU-10 (todavía no existen).

**Frontend (Basado en evento / Basado en estado)**

- **FR-023**: Mientras `GET /demo/estado` responda 200, la lista de proyectos mostrará el botón «Cargar datos de prueba»; en cualquier otro caso no lo mostrará.
- **FR-024**: Cuando la carga responda 201, la interfaz mostrará un aviso de éxito con el enlace «Abrir el proyecto de prueba» a `/projects/25`.
- **FR-025**: Cuando la carga responda 409 «ya cargado», la interfaz lo mostrará como aviso, no como error.
- **FR-026**: Cuando la carga responda cualquier otro error, la interfaz mostrará un mensaje de error.
- **FR-027**: La interfaz deberá integrar el botón en un componente propio con una única línea de montaje en la lista de proyectos, para no chocar con #148.

### Key Entities

- **Proyecto de demostración**: proyecto con id fijo 25 y nombre reservado; su id es la clave de idempotencia.
- **Usuarios Beta**: cuatro usuarios predeterminados sin rol; el rol existe solo como integrante de un proyecto.
- **Bloque de carga**: unidad independiente de la carga (usuarios, proyecto y equipo, Sprints, historias, horas, votación) que cada historia futura puede ampliar.

Las tablas y restricciones están en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md); esta funcionalidad no agrega tablas ni columnas.

## Alcance de la versión 1 y evolución

Es un trabajo en curso. La versión 1 carga solo lo que ya existe en `main`:

| Bloque | Historia de origen | Incluye |
|---|---|---|
| Usuarios | HU-04 | Beta Uno a Beta Cuatro |
| Proyecto y equipo | HU-04 | Proyecto 25 y tres integrantes |
| Sprints | HU-05 | #1 Finalizado, #2 Activo, #3 Pendiente |
| Historias | HU-01, HU-15, HU-06, HU-05 | ~10 historias y una arrastrada |
| Horas | HU-07 | Registros en los Sprints 1 y 2 |
| Votación | HU-02 | Una ronda revelada |

Queda fuera hasta que existan: defectos y métricas (HU-03, HU-08, HU-09, HU-10) y la transición de estado de historias de HU-11 (los estados se cargan ya resueltos).

**Evolución**: cada historia futura que agregue tablas o datos extiende la carga en su propio PR agregando su bloque (una función nueva y su prueba) y actualiza esta tabla, sin modificar los bloques existentes. Esta sección se mantiene al día con cada historia.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Con la función habilitada, un usuario obtiene el proyecto de demostración completo con un solo clic y lo abre en un paso más.
- **SC-002**: El 100 % de las cargas repetidas, o con el id ocupado, termina sin modificar la base.
- **SC-003**: En el 100 % de las cargas fallidas no queda ningún dato.
- **SC-004**: Con dos cargas simultáneas, exactamente una se completa y no quedan duplicados.
- **SC-005**: Con la función apagada, el 100 % de las rutas `/demo/*` responde como inexistente.
- **SC-006**: Tras la carga, el estado del proyecto (Sprint activo vigente, horas dentro de las fechas) es coherente con el día de la carga.

## Assumptions

- La habilitación por `DEMO_HABILITADO` es temporal; se reemplazará por un chequeo de perfil (administrador o beta tester).
- Sin cambio de esquema ni de `MODELO-DE-DATOS.md`.
- Los estados de las historias (En progreso, Completada) se cargan directamente, porque la transición es de HU-11.
- Las horas de los Sprints 1 y 2 se cargan con las reglas de fecha del dominio; la restricción de «solo Sprint activo» de HU-07 rige para la API, no para este dato histórico.
- Dos cargas simultáneas se resuelven por la clave primaria y la transacción; la espera máxima la fija el `busy_timeout` existente.
- El id fijo 25 hace avanzar la secuencia de ids de proyectos; es aceptable.
