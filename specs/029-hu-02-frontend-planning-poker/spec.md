# Feature Specification: Pantalla de Planning Poker (HU-02, Frontend)

**Feature Branch**: `feature/HU-02-frontend-planning-poker`

**Created**: 2026-11-02

**Status**: Draft

**Input**: User description: "HU-02 Frontend (sub-issue #48): pantalla de Planning Poker. El Scrum Master inicia la votación de una historia del Product Backlog; cada integrante elige una carta Fibonacci (1, 2, 3, 5, 8, 13, 21) en secreto; la pantalla muestra quién votó y quién falta (votos de otros «Oculto»); el Scrum Master revela cuando votaron todos y se muestran todos los valores, el promedio exacto y la carta sugerida. Consume la API de la spec 012. Sin inicio de sesión: se elige quién actúa entre los integrantes del proyecto. Votar no modifica Story Points (HU-14 fuera de alcance). Solo frontend (`web/`)."

## Clarifications

### Session 2026-11-02

- Q: ¿Cómo se elige quién actúa si la aplicación no tiene inicio de sesión? → A: Con un selector de integrantes del proyecto (nombre y rol), sin recordar la elección entre visitas; mientras no se elija a nadie no hay cartas ni acciones.
- Q: ¿Desde dónde se llega a la pantalla de votación? → A: Con un botón «Votar» en cada historia del Product Backlog que navega a una ruta nueva del proyecto (`/projects/:id/votacion/:historiaId`; la pantalla vive en su propia feature). El backlog solo enlaza por URL; ninguna feature importa a la otra.
- Q: ¿Cómo se actualiza la vista mientras otros votan? → A: Con un botón «Actualizar» que vuelve a consultar; no hay consulta periódica automática.
- Q: ¿Qué se muestra cuando el servidor rechaza con 403 (sin permiso) o 409 (votación ya iniciada), que el encargo no listaba? → A: El `mensaje` del servidor tal cual, igual que con 400 y 404; y tras un rechazo la pantalla vuelve a consultar la votación para reflejar el estado real.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Votar en secreto con cartas (Priority: P1)

Como integrante del proyecto, abro la pantalla de votación de una historia, elijo quién soy entre los integrantes y selecciono una carta de la baraja Fibonacci (1, 2, 3, 5, 8, 13, 21). Veo mi carta elegida, veo quién del equipo ya votó y quién falta, y los valores de los demás aparecen como «Oculto».

**Why this priority**: Es el núcleo de Planning Poker y el valor de HU-02: estimar en secreto, sin sesgos.

**Independent Test**: Con una votación abierta, elegir un integrante y la carta 8, y verificar que la pantalla muestra «ya votó» para esa persona, su carta 8 visible solo para ella y «Oculto» en los demás.

**Acceptance Scenarios**:

1. **Given** una votación abierta y un integrante elegido, **When** selecciona la carta 8, **Then** se registra su voto, la pantalla marca que ya votó y muestra su carta; los demás integrantes figuran «ya votó» u «falta» sin mostrar su valor («Oculto»).
2. **Given** un integrante que ya votó y la ronda no fue revelada, **When** elige otra carta, **Then** la pantalla muestra la nueva carta como su voto.
3. **Given** que cambia el integrante elegido, **When** se actualiza la vista, **Then** se muestra el voto propio del nuevo integrante y se oculta el del anterior.
4. **Given** un envío de voto en curso, **When** se vuelve a pulsar una carta, **Then** no se envía un segundo pedido.
5. **Given** que el servidor rechaza el voto (por ejemplo, ronda ya revelada), **When** ocurre el error, **Then** la pantalla muestra el mensaje y sigue utilizable.

---

### User Story 2 - Iniciar la votación (Priority: P2)

Como Scrum Master, desde una historia del Product Backlog abro la pantalla de votación y la inicio, lo que abre la ronda 1 para que el equipo vote.

**Why this priority**: Prerrequisito de las demás, pero por sí sola no estima nada.

**Independent Test**: Para una historia sin votación, elegir al Scrum Master, pulsar «Iniciar votación» y ver las cartas habilitadas y a todos los integrantes como «falta».

**Acceptance Scenarios**:

1. **Given** una historia sin votación y el Scrum Master elegido, **When** pulsa «Iniciar votación», **Then** aparece la ronda 1 abierta con todos los integrantes como «falta».
2. **Given** una historia sin votación y un Product Builder elegido, **When** se muestra la pantalla, **Then** no se ofrece el botón «Iniciar votación» y se indica que solo el Scrum Master puede iniciarla.
3. **Given** una historia sin votación, **When** se muestra la pantalla, **Then** no se ofrecen cartas hasta que se inicie la votación.
4. **Given** un envío de inicio en curso, **When** se vuelve a pulsar «Iniciar votación», **Then** no se envía un segundo pedido.
5. **Given** que el servidor rechaza el inicio (historia Completada, ya iniciada, no es Scrum Master), **When** ocurre el error, **Then** se muestra el mensaje del servidor.

---

### User Story 3 - Revelar y ver el resultado (Priority: P1)

Como Scrum Master, cuando todos votaron, pulso «Revelar»: la pantalla muestra el valor de cada integrante, el promedio exacto y la carta sugerida.

**Why this priority**: Sin revelación los votos no sirven; habilita HU-14.

**Independent Test**: Con todos los votos cargados, revelar y verificar valores, promedio (5,33) y carta sugerida (5).

**Acceptance Scenarios**:

1. **Given** que votaron todos y el Scrum Master está elegido, **When** pulsa «Revelar», **Then** se muestran el valor de cada integrante, el promedio exacto y la carta sugerida.
2. **Given** que falta votar al menos un integrante, **When** se muestra la pantalla al Scrum Master, **Then** el botón «Revelar» está deshabilitado y se indica quién falta.
3. **Given** un Product Builder elegido, **When** se muestra la pantalla, **Then** no se ofrece «Revelar».
4. **Given** una ronda revelada, **When** se muestra la pantalla, **Then** las cartas ya no se pueden elegir y se aclara que registrar la estimación acordada se hace en otra historia (HU-14).
5. **Given** un envío de revelación en curso, **When** se vuelve a pulsar «Revelar», **Then** no se envía un segundo pedido.
6. **Given** que el servidor rechaza la revelación, **When** ocurre el error, **Then** se muestra el mensaje del servidor.

---

### User Story 4 - Llegar a la pantalla y mantenerla al día (Priority: P2)

Como integrante, desde cada historia del Product Backlog llego a su pantalla de votación, y con un botón «Actualizar» veo lo que votaron los demás mientras la ronda sigue abierta.

**Why this priority**: Sin acceso desde el backlog ni actualización la pantalla no es usable por un equipo.

**Independent Test**: En el Product Backlog pulsar «Votar» en una historia y llegar a su pantalla; pulsar «Actualizar» y ver reflejado un voto nuevo.

**Acceptance Scenarios**:

1. **Given** el Product Backlog con historias, **When** pulso «Votar» en una, **Then** navego a la pantalla de votación de esa historia dentro del proyecto.
2. **Given** la pantalla de votación abierta, **When** pulso «Actualizar», **Then** se vuelve a consultar y se reflejan los votos nuevos.
3. **Given** una historia inexistente o un proyecto sin integrantes, **When** se abre la pantalla, **Then** se muestra un mensaje claro y la pantalla no se rompe.
4. **Given** que el servidor no responde, **When** se abre o actualiza la pantalla, **Then** se muestra el aviso de conexión y se puede reintentar.

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** ofrecer exactamente las cartas 1, 2, 3, 5, 8, 13 y 21.
- **WHILE** la ronda no fue revelada, el sistema **shall** mostrar como «Oculto» el valor de los demás integrantes que ya votaron, y como «Falta» a los que no.
- **WHILE** no hay integrante elegido, el sistema **shall** no ofrecer cartas ni acciones de Scrum Master, y pedir que se elija quién actúa.
- **WHEN** se elige un integrante, el sistema **shall** consultar la votación como ese integrante.
- **IF** el rol del integrante elegido no es `scrum_master`, **THEN** el sistema **shall** no mostrar «Iniciar votación» ni «Revelar».
- **IF** falta el voto de algún integrante, **THEN** el sistema **shall** deshabilitar «Revelar».
- **IF** hay un pedido de iniciar, votar o revelar en curso, **THEN** el sistema **shall** ignorar un nuevo accionamiento de esa misma acción.
- **IF** el servidor responde un error (validación, no encontrado, permiso, conflicto), **THEN** el sistema **shall** mostrar el mensaje del servidor sin romper la pantalla y conservar el último estado conocido.
- **IF** no se puede conectar con el servidor, **THEN** el sistema **shall** mostrar un aviso de conexión.
- **IF** el proyecto no tiene integrantes, **THEN** el sistema **shall** informarlo y no ofrecer acciones.
- **WHEN** la ronda está revelada, el sistema **shall** mostrar el valor de cada integrante, el promedio exacto y la carta sugerida.
- **IF** la respuesta llega para un integrante o historia que ya no es el actual, **THEN** el sistema **shall** descartarla.
- El sistema **shall** no modificar los Story Points de la historia.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: La pantalla MUST permitir elegir quién actúa entre los integrantes del proyecto, mostrando su nombre y rol.
- **FR-002**: La pantalla MUST consultar la votación de la historia en nombre del integrante elegido y mostrar la ronda, y por integrante si ya votó o falta.
- **FR-003**: La pantalla MUST mostrar «Oculto» como valor de los demás integrantes que ya votaron mientras la ronda no esté revelada, y el valor propio del integrante elegido.
- **FR-004**: La pantalla MUST ofrecer las cartas 1, 2, 3, 5, 8, 13 y 21 cuando hay una ronda abierta y hay un integrante elegido, y enviar el voto al elegir una carta.
- **FR-005**: La pantalla MUST permitir cambiar el voto propio mientras la ronda no esté revelada.
- **FR-006**: La pantalla MUST ofrecer «Iniciar votación» solo al integrante con rol `scrum_master` y solo si la historia no tiene votación.
- **FR-007**: La pantalla MUST ofrecer «Revelar» solo al integrante con rol `scrum_master`, habilitado únicamente cuando votaron todos los integrantes y la ronda no está revelada.
- **FR-008**: Tras revelar, la pantalla MUST mostrar el valor de cada integrante, el promedio exacto y la carta sugerida, y dejar de ofrecer cartas.
- **FR-009**: La pantalla MUST evitar el doble envío de iniciar, votar y revelar mientras haya un pedido en curso.
- **FR-010**: La pantalla MUST mostrar un mensaje claro por cada error del servidor (validación, no encontrado, sin permiso, conflicto, conexión) sin romperse.
- **FR-011**: La pantalla MUST ofrecer «Actualizar» para volver a consultar la votación.
- **FR-012**: El Product Backlog MUST ofrecer en cada historia un botón «Votar» que navega a la pantalla de votación de esa historia dentro del proyecto.
- **FR-013**: La pantalla MUST NOT modificar los Story Points de la historia ni ofrecer registrar la estimación acordada (HU-14).
- **FR-014**: Cada feature del frontend MUST permanecer aislada: la nueva pantalla no importa a otras features y el Product Backlog no importa a la nueva.

### Key Entities

- **Votación de una historia**: ronda actual, momento de revelación, integrantes con su estado de voto y valor, y promedio sugerido. Definida por el backend en la spec 012; el esquema está en [MODELO-DE-DATOS.md §4.7](../../MODELO-DE-DATOS.md#47-rondas_votacion-y-votos--hu-02-hu-14-tablas-nuevas-migración-003) y esta spec no lo redefine.
- **Integrante**: persona del proyecto con rol `scrum_master` o `product_builder`; es quien actúa en la pantalla.
- **Historia**: lo que se estima, identificada por la ruta.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un integrante vota en dos pasos: elegir quién es y pulsar una carta.
- **SC-002**: En el 100% de las vistas con la ronda sin revelar, ningún valor de otro integrante es visible.
- **SC-003**: Tras revelar, el equipo ve valores, promedio y carta sugerida sin calcular nada a mano.
- **SC-004**: El 100% de los errores del servidor se muestran con un mensaje comprensible y la pantalla sigue operable.
- **SC-005**: Una acción repetida por doble clic genera un único pedido.

## Assumptions

- Las decisiones de selector de integrante, botón «Votar» y botón «Actualizar» están registradas en Clarifications.
- El backend de la spec 012 está implementado y no se modifica; sus rechazos incluyen, además de 400 y 404, 403 (sin permiso) y 409 (votación ya iniciada), todos con un campo `mensaje`.
- No hay inicio de sesión: se elige quién actúa entre los integrantes del proyecto (selector).
- Se llega a la pantalla con un botón «Votar» en cada historia del Product Backlog, que navega a una ruta nueva del proyecto.
- La vista se actualiza manualmente con el botón «Actualizar», sin consulta periódica automática.
- Los textos de la interfaz están en español, como el resto de la aplicación.
- Abrir nuevas rondas, comparar diferencias y registrar la estimación acordada son HU-14 y quedan fuera de alcance.
