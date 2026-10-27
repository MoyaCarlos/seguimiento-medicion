# Feature Specification: Votación de Esfuerzo en Planning Poker

**Feature Branch**: `feature/HU-02-votacion-planning-poker`

**Created**: 2026-10-26

**Status**: Especificada y validada en la reunión de definición (2026-10-27): solo el Scrum Master inicia y revela, no se revela hasta que votaron todos, permisos con 403.

**Input**: User description: "HU-02: Votación de Esfuerzo en Planning Poker
(issue #2, sub-issue de especificación #44). Como integrante del equipo (Scrum
Master o Product Builder) quiero votar el esfuerzo de una historia en secreto
usando cartas virtuales para poder realizar una estimación colaborativa sin
sesgos. El análisis de las diferencias, las nuevas rondas y el registro de la
estimación acordada los cubre HU-14. Depende de HU-04 (integrantes que votan)
y HU-01 (historia a estimar). El modelo de votación contempla rondas desde el
inicio (cada voto pertenece a una ronda). Alcance: backend; el frontend va en
otra sub-issue."

## Clarifications

### Session 2026-10-27

- Q: Si todavía falta votar algún integrante, ¿qué pasa al presionar "Revelar"? → A: Se rechaza; solo se revela cuando votaron todos los integrantes del proyecto.
- Q: Sin autenticación, ¿se valida que quien inicia y revela sea el Scrum Master? → A: Sí; la petición indica qué integrante actúa y, si no tiene rol `scrum_master` en el proyecto, se rechaza.
- Q: ¿Cómo se calcula el "promedio sugerido"? → A: Se informan los dos: el promedio aritmético exacto (3, 5 y 8 → 5,33) y la carta de la baraja más cercana (→ 5).
- Q: Antes de revelar, ¿quien votó puede ver su propio voto al consultar la votación? → A: Sí; la consulta indica quién consulta, le muestra su propio valor y deja el de los demás como "Oculto".
- Q: ¿Se puede iniciar la votación de una historia "Completada"? → A: No; se rechaza. Las historias "Nueva" y "En progreso" se pueden votar.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Votar en secreto (Priority: P1)

Como integrante del proyecto (Scrum Master o Product Builder), una vez que el
Scrum Master abrió la votación de una historia, elijo una carta de la baraja
Fibonacci y confirmo mi voto. Mi voto queda registrado, pero el resto del
equipo solo ve que ya voté, no qué valor elegí.

**Why this priority**: Es el núcleo de Planning Poker: sin votos ocultos no hay
estimación "sin sesgos", que es el valor que pide la historia.

**Independent Test**: Se abre la votación de una historia, un integrante vota
8 y, al consultar la votación, aparece como "votó" con el valor oculto.

**Acceptance Scenarios**:

1. **Given** el Scrum Master inició la votación de una historia del proyecto,
   **When** un integrante selecciona la carta 8 y confirma su voto, **Then** el
   voto queda registrado en la ronda actual y, cuando otro integrante consulta
   la votación, el que votó figura como que ya votó, con el valor "Oculto";
   cuando consulta el propio votante, ve su carta (8).
2. **Given** la votación de una historia está abierta y un integrante ya votó,
   **When** ese mismo integrante confirma otra carta antes de la revelación,
   **Then** su voto se reemplaza por el nuevo (sigue habiendo un solo voto suyo
   en la ronda) y sigue oculto.
3. **Given** la votación de una historia está abierta, **When** se intenta
   votar con un valor que no pertenece a la baraja (por ejemplo 4), **Then** el
   sistema rechaza el voto y la ronda no cambia.
4. **Given** la votación de una historia está abierta, **When** vota una
   persona que no es integrante del proyecto, **Then** el sistema rechaza el
   voto.

---

### User Story 2 - Revelar los votos (Priority: P1)

Como Scrum Master, cuando el equipo terminó de votar, revelo los votos: el
sistema muestra el valor que eligió cada integrante y calcula el promedio
sugerido.

**Why this priority**: Sin revelación los votos no sirven para estimar, y es lo
que habilita a HU-14 (diferencias, nuevas rondas y estimación acordada).

**Independent Test**: Con una ronda en la que votaron todos los integrantes, se
revela y se verifica que la consulta devuelve cada valor y el promedio.

**Acceptance Scenarios**:

1. **Given** todos los integrantes del proyecto confirmaron su voto (3, 5 y 8),
   **When** el Scrum Master presiona "Revelar", **Then** el sistema muestra el
   valor elegido por cada integrante y el promedio sugerido.
2. **Given** una ronda ya fue revelada, **When** un integrante intenta votar o
   cambiar su voto, **Then** el sistema rechaza la operación y los votos
   revelados no cambian.
3. **Given** todavía falta votar al menos un integrante, **When** el Scrum
   Master intenta revelar, **Then** el sistema rechaza la revelación y los
   votos siguen ocultos.
4. **Given** todos los integrantes votaron, **When** un Product Builder
   intenta revelar, **Then** el sistema rechaza la operación porque solo el
   Scrum Master del proyecto puede revelar.

---

### User Story 3 - Iniciar la votación de una historia (Priority: P2)

Como Scrum Master, selecciono una historia del Product Backlog e inicio su
votación, lo que abre la primera ronda para que el equipo vote.

**Why this priority**: Es el prerrequisito de las otras dos, pero de poco valor
por sí solo: abrir una votación sin votar no estima nada.

**Independent Test**: Se inicia la votación de una historia sin rondas y se
verifica que queda una ronda 1 abierta, sin votos.

**Acceptance Scenarios**:

1. **Given** una historia del proyecto que todavía no tiene votación, **When**
   el Scrum Master inicia la votación, **Then** queda abierta la ronda 1 de
   esa historia, sin votos y sin revelar.
2. **Given** una historia que ya tiene una votación (abierta o revelada),
   **When** se intenta iniciar la votación otra vez, **Then** el sistema rechaza
   la operación (abrir rondas siguientes es parte de HU-14).
3. **Given** el identificador de una historia que no existe, **When** se
   intenta iniciar su votación, **Then** el sistema responde "historia no
   encontrada".
4. **Given** una historia sin votación, **When** un integrante con rol
   Product Builder intenta iniciarla, **Then** el sistema rechaza la operación
   y la historia sigue sin rondas.
5. **Given** una historia en estado "Completada" y sin votación, **When** el
   Scrum Master intenta iniciar su votación, **Then** el sistema rechaza la
   operación y la historia sigue sin rondas.

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** asociar cada voto a una ronda de votación de una
  historia, y **shall** admitir como máximo un voto por integrante y ronda.
- El sistema **shall** aceptar como valor de voto únicamente las cartas de la
  baraja Fibonacci del proyecto: 1, 2, 3, 5, 8, 13 y 21 (la misma escala que
  la estimación en Story Points).
- **WHEN** el Scrum Master inicia la votación de una historia sin rondas,
  **THEN** el sistema **shall** crear la ronda número 1, sin votos y sin
  revelar.
- **IF** la historia ya tiene al menos una ronda, **THEN** el sistema **shall**
  rechazar el inicio de la votación (la ronda 2 en adelante la abre HU-14).
- **WHILE** la ronda actual no fue revelada, **WHEN** se consulta la votación,
  **THEN** el sistema **shall** informar qué integrantes ya votaron y cuáles
  faltan, e incluir solo el valor del voto de quien consulta (el de los demás
  queda "Oculto").
- **WHILE** la ronda actual no fue revelada, **WHEN** un integrante que ya votó
  confirma otra carta, **THEN** el sistema **shall** reemplazar su voto
  anterior.
- **WHEN** el Scrum Master revela la ronda actual, **THEN** el sistema
  **shall** registrar el momento de la revelación y, desde entonces, informar
  el valor de cada voto y el promedio sugerido.
- **IF** la ronda actual ya fue revelada, **THEN** el sistema **shall**
  rechazar nuevos votos, cambios de voto y una segunda revelación.
- **IF** el valor del voto no pertenece a la baraja, **THEN** el sistema
  **shall** rechazarlo con un error de validación sobre el valor.
- **IF** quien vota no es integrante del proyecto de la historia, **THEN** el
  sistema **shall** rechazar el voto.
- **IF** la historia está en estado "Completada", **THEN** el sistema
  **shall** rechazar el inicio de su votación.
- **IF** se vota sobre una historia que no tiene votación iniciada, **THEN**
  el sistema **shall** rechazar el voto.
- **IF** la historia indicada no existe, **THEN** el sistema **shall**
  responder "historia no encontrada", antes de validar el resto de los datos
  (mismo criterio que HU-04 y HU-05).
- **IF** quien inicia o revela la votación no es integrante del proyecto con
  rol `scrum_master`, **THEN** el sistema **shall** rechazar la operación sin
  modificar la ronda.
- **IF** al revelar falta el voto de al menos un integrante del proyecto,
  **THEN** el sistema **shall** rechazar la revelación y mantener los votos
  ocultos.
- **WHEN** se revela una ronda, **THEN** el sistema **shall** informar el
  promedio aritmético exacto de los votos, redondeado a 2 decimales (3, 5 y
  8 → 5,33), y la carta de la baraja más cercana a ese promedio (→ 5).
- **IF** el promedio queda a la misma distancia de dos cartas (3 y 5 →
  promedio 4), **THEN** el sistema **shall** sugerir la carta mayor (→ 5).
- **IF** todos votaron la misma carta, **THEN** el promedio exacto y la carta
  sugerida **shall** ser esa carta.
- **WHILE** se consulta la votación, **THEN** el sistema **shall** no
  modificar rondas ni votos (operación de solo lectura).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir iniciar la votación de una historia
  existente que no tenga rondas, creando la ronda número 1 abierta (sin
  revelar).
- **FR-002**: El sistema MUST rechazar iniciar la votación de una historia que
  ya tenga alguna ronda.
- **FR-003**: El sistema MUST permitir que un integrante del proyecto de la
  historia registre un voto en la ronda actual mientras no esté revelada.
- **FR-004**: El sistema MUST aceptar solo valores de la baraja Fibonacci (1,
  2, 3, 5, 8, 13, 21) y rechazar cualquier otro con un error de validación.
- **FR-005**: El sistema MUST guardar como máximo un voto por integrante y
  ronda; un segundo voto del mismo integrante antes de revelar reemplaza al
  anterior.
- **FR-006**: El sistema MUST rechazar el voto de quien no sea integrante del
  proyecto de la historia.
- **FR-007**: Mientras la ronda no esté revelada, la consulta de la votación
  MUST informar quién votó y quién falta, y exponer solo el valor del voto de
  quien consulta; el resto de los valores se informa como "Oculto".
- **FR-008**: El sistema MUST permitir revelar la ronda actual de una historia,
  registrando el momento de la revelación.
- **FR-009**: Una vez revelada, la consulta de la votación MUST informar el
  valor elegido por cada integrante y el promedio sugerido.
- **FR-010**: El sistema MUST rechazar votos, cambios de voto y nuevas
  revelaciones sobre una ronda ya revelada.
- **FR-011**: El sistema MUST responder "historia no encontrada" para
  cualquier operación sobre una historia inexistente, antes de validar el
  resto de los datos.
- **FR-012**: La votación MUST conservar las rondas y sus votos sin borrarlos,
  identificando como ronda actual la de mayor número, para que HU-14 agregue
  rondas y las anteriores queden como historial.
- **FR-013**: Esta historia MUST NOT modificar la estimación en Story Points
  de la historia; registrar la estimación acordada es parte de HU-14.
- **FR-014**: El sistema MUST rechazar iniciar o revelar la votación cuando
  quien lo pide no sea integrante del proyecto con rol `scrum_master`.
- **FR-015**: El sistema MUST rechazar la revelación mientras falte el voto
  de algún integrante del proyecto.
- **FR-016**: El promedio sugerido MUST informarse de dos formas: el promedio
  aritmético exacto, redondeado a 2 decimales, y la carta de la baraja más
  cercana a él; ante un empate entre dos cartas se sugiere la mayor.
- **FR-017**: El sistema MUST rechazar iniciar la votación de una historia en
  estado "Completada"; las historias "Nueva" y "En progreso" se pueden votar.

### Key Entities

- **Ronda de votación**: una vuelta de votación sobre una historia. Tiene un
  número correlativo por historia (1, 2, 3…) y un momento de revelación
  (vacío mientras los votos están ocultos). Detalle del esquema en
  [MODELO-DE-DATOS.md §4.7](../../MODELO-DE-DATOS.md#47-rondas_votacion-y-votos--hu-02-hu-14-tablas-nuevas-migración-003).
- **Voto**: la carta que eligió un integrante en una ronda. Uno por integrante
  y ronda. Mismo apartado §4.7.
- **Integrante** (HU-04, ya existente): usuario de un proyecto con rol
  `scrum_master` o `product_builder`. Son quienes votan.
- **Historia de Usuario** (HU-01, ya existente): lo que se estima.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Antes de la revelación, en el 100% de las consultas de una
  votación ningún integrante puede ver el valor que eligió otro (solo el
  suyo).
- **SC-002**: Un integrante puede votar en un solo paso (elegir carta y
  confirmar), sin datos adicionales.
- **SC-003**: Tras revelar, el equipo ve los valores de todos y el promedio en
  la misma consulta, sin calcular nada a mano.
- **SC-004**: El 100% de los votos con valores fuera de la baraja o de
  personas ajenas al proyecto se rechazan sin alterar la ronda.
- **SC-005**: Ninguna ronda ni voto se pierde: HU-14 puede consultar el
  historial completo de rondas de una historia.

## Assumptions

- El esquema ya está definido en MODELO-DE-DATOS.md §4.7 (`rondas_votacion` y
  `votos`, migración 003); esta historia lo usa sin cambiarlo.
- La aplicación no tiene autenticación: quien vota, inicia o revela se
  identifica con su usuario en la petición, igual que en el resto de la API;
  el rol se toma del integrante en ese proyecto (HU-04).
- Si un proyecto tiene más de un Scrum Master (no está prohibido, D4),
  cualquiera de ellos puede iniciar y revelar.
- Ante un empate entre dos cartas se sugiere la mayor: en estimación conviene
  errar hacia arriba. Es una sugerencia; la estimación la registra HU-14.
- Votan todos los integrantes del proyecto, incluido el Scrum Master (la
  historia dice "integrante del equipo (Scrum Master o Product Builder)").
  "Todos" se evalúa al revelar: si se suma un integrante con la ronda abierta,
  también tiene que votar antes de que se pueda revelar.
- Se puede votar cualquier historia del proyecto que no esté "Completada"
  (FR-017), tenga o no Sprint o estimación cargada: votar no modifica la
  historia (FR-013). La restricción
  de no reemplazar la estimación con Sprint `Activo` o `Finalizado` (D2b) es
  de HU-14, al registrar la estimación acordada.
- Votar no mira el estado de la historia: FR-017 solo frena el inicio. Hoy una
  historia no puede completarse con la votación abierta (completar historias es
  de HU-11); si HU-11 lo permite, decide ella si cierra la votación.
- Cambiar el voto antes de revelar está permitido: en Planning Poker real se
  puede cambiar la carta mientras no se dio vuelta.
- Detectar diferencias, abrir nuevas rondas y registrar la estimación acordada
  quedan fuera de alcance (HU-14). La pantalla con las cartas es otra
  sub-issue (Frontend).
