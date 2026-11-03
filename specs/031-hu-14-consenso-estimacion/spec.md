# Feature Specification: Consenso y Registro de la Estimación Acordada

**Feature Branch**: `feature/HU-14-consenso-estimacion`

**Created**: 2026-11-03

**Status**: Especificada y validada por el equipo en la reunión de definición.

**Input**: User description: "HU-14: Consenso y Registro de la Estimación Acordada
(issue #42, sub-issue de especificación #89). Como Scrum Master quiero que, al
revelar los votos de Planning Poker, el sistema detecte las diferencias entre las
estimaciones, permita hacer nuevas rondas y registre la estimación acordada en la
historia para poder obtener los Story Points sobre los que se calculan las métricas.
Depende de HU-02 (votos y revelación) y HU-01 (campo de Story Points de la historia).
No bloquea al cálculo de métricas ni al reporte: los Story Points también se cargan de
forma directa al editar el Product Backlog (HU-15). Alcance: backend; el frontend va en
otra sub-issue (#93)."

## Clarifications

### Session 2026-11-03

- Q: ¿Qué se considera una "diferencia" entre estimaciones? → A: Cualquier voto que no coincida, sin importar la distancia entre las cartas.
- Q: ¿Se puede iniciar una nueva ronda si la ronda revelada no tiene diferencias? → A: Sí; lo decide el Scrum Master.
- Q: ¿El historial de rondas se puede consultar o solo queda guardado? → A: Se puede consultar en la API; la pantalla no entra en esta historia.
- Q: ¿El valor acordado puede ser cualquier carta de la baraja, o tiene que ser una de las votadas en la ronda? → A: Cualquier carta de la baraja, aunque nadie la haya votado.
- Q: ¿Se puede registrar la estimación acordada de una historia "Completada"? → A: Sí, sujeto solo a D2b: se permite si no tenía estimación, o si no está en un Sprint `Activo` o `Finalizado`.
- Q: Después de registrar la estimación, ¿se puede abrir otra ronda y registrar otro valor? → A: Sí; registrar no cierra la votación y siguen valiendo las mismas reglas, incluida D2b.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Detectar diferencias entre los votos (Priority: P1)

Como Scrum Master, al consultar una ronda revelada veo si las estimaciones coinciden y,
si no coinciden, cuáles son los valores mínimo y máximo votados, para abrir la discusión
del equipo sin comparar los votos a mano.

**Why this priority**: Es lo que le dice al equipo si tiene que seguir conversando o si
ya puede registrar la estimación.

**Independent Test**: Se revela una ronda con votos 3, 5 y 8, y la consulta indica que
hay diferencias, con mínimo 3 y máximo 8.

**Acceptance Scenarios**:

1. **Given** se revelaron los votos 3, 5 y 8 de una historia, **When** se consulta la
   votación, **Then** el sistema indica que hay diferencias, con mínimo 3 y máximo 8.
2. **Given** se revelaron los votos 5, 5 y 5 de una historia, **When** se consulta la
   votación, **Then** el sistema indica que no hay diferencias.
3. **Given** una ronda todavía no fue revelada, **When** se consulta la votación,
   **Then** el sistema no informa diferencias, mínimo ni máximo, para no dar pistas de
   los votos ocultos.

---

### User Story 2 - Registrar la estimación acordada (Priority: P1)

Como Scrum Master, cuando el equipo llegó a un acuerdo registro el valor acordado, y la
historia queda con esos Story Points como estimación.

**Why this priority**: Es el resultado de la historia: los Story Points acordados son los
que usan las métricas.

**Independent Test**: Con una ronda revelada, se registra el valor 5 y se verifica que la
historia queda con 5 Story Points.

**Acceptance Scenarios**:

1. **Given** una historia sin estimación y con su ronda actual revelada, **When** el Scrum
   Master registra el valor acordado 5, **Then** la historia queda con 5 Story Points.
2. **Given** una historia con 3 Story Points cargados de forma directa (HU-15), sin Sprint
   o con un Sprint `Pendiente`, y con su ronda actual revelada, **When** el Scrum Master
   registra el valor acordado 8, **Then** la historia queda con 8 Story Points.
3. **Given** una historia sin estimación, asignada a un Sprint `Activo` y con su ronda
   actual revelada, **When** el Scrum Master registra el valor acordado 5, **Then** la
   historia queda con 5 Story Points, porque se carga por primera vez.
4. **Given** una ronda revelada, **When** se intenta registrar un valor fuera de la baraja
   Fibonacci (por ejemplo 4), **Then** el sistema rechaza la operación y la historia
   mantiene su estimación anterior.
5. **Given** la ronda actual de una historia no fue revelada, **When** el Scrum Master
   intenta registrar una estimación, **Then** el sistema rechaza la operación y la
   historia mantiene su estimación anterior.
6. **Given** una historia con estimación cargada, asignada a un Sprint `Activo` o
   `Finalizado` y con su ronda actual revelada, **When** el Scrum Master intenta registrar
   un valor distinto, **Then** el sistema rechaza la operación y la historia mantiene su
   estimación anterior.
7. **Given** una ronda revelada, **When** un Product Builder intenta registrar la
   estimación, **Then** el sistema rechaza la operación y la historia mantiene su
   estimación anterior.
8. **Given** una ronda revelada con votos 3 y 8, **When** el Scrum Master registra el valor
   acordado 5, **Then** la historia queda con 5 Story Points, aunque nadie haya votado esa carta.
9. **Given** una historia "Completada" sin estimación, con su ronda actual revelada, **When** el
   Scrum Master registra el valor acordado 5, **Then** la historia queda con 5 Story Points.

---

### User Story 3 - Abrir una nueva ronda (Priority: P2)

Como Scrum Master, cuando las estimaciones no coinciden inicio una nueva ronda: los votos
anteriores dejan de contar, quedan como historial, y el equipo vota de nuevo.

**Why this priority**: Es la vía para llegar al consenso, pero el equipo puede acordar un
valor conversando sobre la primera ronda y registrarlo sin votar otra vez.

**Independent Test**: Con la ronda 1 revelada y con diferencias, se inicia una nueva ronda
y se verifica que la ronda actual es la 2, sin votos y oculta, y que los votos de la ronda
1 se conservan.

**Acceptance Scenarios**:

1. **Given** la ronda 1 de una historia fue revelada y tiene diferencias, **When** el Scrum
   Master inicia una nueva ronda, **Then** la ronda actual pasa a ser la 2, abierta y sin
   votos, y los votos de la ronda 1 quedan registrados como historial.
2. **Given** la ronda 2 está abierta, **When** los integrantes votan y el Scrum Master
   revela, **Then** las diferencias y el promedio se calculan solo con los votos de la
   ronda 2.
3. **Given** la ronda actual todavía no fue revelada, **When** el Scrum Master intenta
   iniciar una nueva ronda, **Then** el sistema rechaza la operación y la ronda no cambia.
4. **Given** una ronda revelada con diferencias, **When** un Product Builder intenta iniciar
   una nueva ronda, **Then** el sistema rechaza la operación y no se crea ninguna ronda.
5. **Given** una historia sin votación iniciada, **When** el Scrum Master intenta iniciar
   una nueva ronda, **Then** el sistema rechaza la operación y le indica que la votación no
   está iniciada.
6. **Given** una historia en estado "Completada" con una ronda revelada, **When** el Scrum
   Master intenta iniciar una nueva ronda, **Then** el sistema rechaza la operación y no se
   crea ninguna ronda.
7. **Given** una ronda revelada sin diferencias, **When** el Scrum Master inicia una nueva
   ronda, **Then** la nueva ronda queda abierta y sin votos, igual que cuando hay diferencias.
8. **Given** la ronda 1 fue revelada y la ronda 2 está abierta con votos ocultos, **When** se
   consulta el historial de rondas de la historia, **Then** el sistema devuelve la ronda 1 con
   su número, su momento de revelación y el voto de cada integrante, y no devuelve la ronda 2.

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** considerar que hay diferencias cuando los votos de la ronda revelada
  no son todos iguales, sin importar la distancia entre las cartas (3 y 5 son una diferencia).
- **WHEN** se consulta una ronda revelada, **THEN** el sistema **shall** informar si hay
  diferencias; **IF** las hay, **THEN** **shall** informar también el mínimo y el máximo votados.
- **WHILE** la ronda actual no fue revelada, el sistema **shall** omitir la indicación de
  diferencias, el mínimo y el máximo.
- **IF** la ronda revelada tiene un solo voto, **THEN** el sistema **shall** indicar que no hay
  diferencias.
- **WHEN** el Scrum Master inicia una nueva ronda, **THEN** el sistema **shall** crear la ronda
  con el número siguiente al de la ronda actual, sin votos y sin revelar, y **shall** conservar
  sin cambios las rondas anteriores y sus votos.
- **WHEN** se consulta el historial de rondas de una historia, **THEN** el sistema **shall**
  devolver solo las rondas reveladas, ordenadas por número; una ronda sin revelar no aparece,
  para no exponer sus votos.
- **IF** la historia existe y no tiene rondas reveladas, **THEN** el historial **shall** devolverse
  vacío, sin error.
- **IF** dos pedidos simultáneos intentan abrir la misma ronda, **THEN** el sistema **shall**
  crear una sola y rechazar el otro.
- **WHEN** el Scrum Master registra un valor acordado válido, **THEN** el sistema **shall**
  guardarlo como la estimación en Story Points de la historia, reemplazando la anterior.
- **IF** el valor acordado no pertenece a la baraja Fibonacci (1, 2, 3, 5, 8, 13, 21), **THEN**
  el sistema **shall** rechazarlo con un error de validación sobre la estimación.
- **IF** la historia no tiene votación iniciada o su ronda actual no fue revelada, **THEN** el
  sistema **shall** rechazar el registro de la estimación. Una ronda anterior revelada no
  alcanza si la actual sigue abierta.
- **IF** la historia ya tiene estimación y está asignada a un Sprint `Activo` o `Finalizado`,
  **THEN** el sistema **shall** rechazar un valor distinto del cargado (decisión D2b).
- **IF** el valor acordado es igual a la estimación cargada, **THEN** el sistema **shall**
  aceptar la operación sin cambiar la historia, también con un Sprint `Activo` o `Finalizado`
  (mismo criterio que la edición de HU-15).
- **IF** quien registra la estimación o inicia una nueva ronda no es integrante del proyecto
  con rol Scrum Master, **THEN** el sistema **shall** rechazar la operación sin modificar nada.
- **IF** la historia indicada no existe, **THEN** el sistema **shall** responder "historia no
  encontrada", antes de validar el resto de los datos (decisión transversal del backlog).
- **IF** el registro de la estimación incumple más de una regla, **THEN** el sistema **shall**
  informar la primera en este orden: permiso, `estimacion_sp`, votación iniciada y ronda actual
  revelada, D2b.
- **WHERE** el valor acordado es una carta válida, el sistema **shall** aceptarlo aunque nadie
  lo haya votado en la ronda.
- **IF** ya se registró una estimación, **THEN** el sistema **shall** seguir permitiendo otra
  ronda o un nuevo registro, sujeto a las mismas reglas: registrar no cierra la votación.
- **WHILE** se consulta la votación o el historial, el sistema **shall** no modificar rondas,
  votos ni estimaciones (operación de solo lectura).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: La consulta de una ronda revelada MUST indicar si hay diferencias entre los votos
  y, si las hay, el valor mínimo y el máximo.
- **FR-002**: La consulta de una ronda no revelada MUST NOT informar diferencias, mínimo ni máximo.
- **FR-003**: El sistema MUST permitir que el Scrum Master del proyecto inicie una nueva ronda
  cuando la ronda actual de la historia está revelada, haya o no diferencias.
- **FR-004**: La nueva ronda MUST quedar abierta, sin votos, con el número siguiente, y pasar a
  ser la ronda actual: votar, revelar y consultar operan sobre ella con las reglas de HU-02.
- **FR-005**: El sistema MUST conservar las rondas anteriores y sus votos, sin borrarlos ni
  modificarlos.
- **FR-006**: El sistema MUST ofrecer una consulta del historial de rondas de una historia, con el
  número, el momento de revelación y el voto de cada integrante de cada ronda revelada. Las rondas
  sin revelar MUST NOT aparecer en esa consulta.
- **FR-007**: El sistema MUST rechazar iniciar una nueva ronda si la ronda actual no está revelada,
  si la historia no tiene votación iniciada o si está en estado "Completada".
- **FR-008**: El sistema MUST permitir que el Scrum Master del proyecto registre el valor acordado
  como estimación en Story Points de la historia cuando su ronda actual está revelada.
- **FR-009**: El sistema MUST rechazar el registro, sin cambiar la estimación, si el valor no
  pertenece a la baraja Fibonacci, si la ronda actual no está revelada o si la historia no tiene
  votación iniciada.
- **FR-010**: El sistema MUST rechazar, sin cambiar la estimación, un valor distinto del cargado
  cuando la historia ya tiene estimación y está asignada a un Sprint `Activo` o `Finalizado`
  (D2b). Cargarla por primera vez sí se permite.
- **FR-011**: El sistema MUST rechazar iniciar una nueva ronda o registrar la estimación cuando
  quien lo pide no sea integrante del proyecto con rol Scrum Master.
- **FR-012**: El sistema MUST responder "historia no encontrada" para cualquier operación sobre
  una historia inexistente, antes de validar el resto de los datos.
- **FR-013**: La estimación registrada MUST ser la misma que usan el Product Backlog (HU-15) y las
  métricas: no se crea un segundo dato de estimación.

### Key Entities

- **Ronda de votación** (HU-02, ya existente): una vuelta de votación sobre una historia, con un
  número correlativo y un momento de revelación. La ronda actual es la de mayor número y las
  anteriores son el historial. Detalle en
  [MODELO-DE-DATOS.md §4.7](../../MODELO-DE-DATOS.md#47-rondas_votacion-y-votos--hu-02-hu-14-tablas-nuevas-migración-003).
- **Voto** (HU-02, ya existente): la carta que eligió un integrante en una ronda. Mismo apartado §4.7.
- **Estimación acordada**: no es una entidad propia. Se guarda como los Story Points de la historia
  ([MODELO-DE-DATOS.md §4.4](../../MODELO-DE-DATOS.md#44-historias--hu-01-hu-15-hu-06-hu-11-hu-14-antes-backlog_items)).
- **Diferencias**: dato calculado al consultar una ronda revelada (si hay, mínimo y máximo). No se guarda.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Tras revelar, el equipo ve si hay diferencias y sus extremos en la misma consulta, sin
  comparar los votos a mano.
- **SC-002**: En el 100% de las consultas de una ronda sin revelar, las diferencias no dan ninguna
  pista sobre los votos ocultos.
- **SC-003**: Después de varias rondas, el 100% de los votos de las rondas reveladas se puede
  consultar en el historial, sin cambios.
- **SC-004**: El 100% de los registros rechazados deja la estimación de la historia como estaba.
- **SC-005**: Una estimación registrada se ve de inmediato como los Story Points de la historia en
  el Product Backlog, sin cargarla de nuevo.

## Assumptions

- El esquema ya está definido en MODELO-DE-DATOS.md §4.4 y §4.7, y §5 ya asigna a HU-14 las reglas
  de revelar antes de abrir otra ronda y de no cambiar una estimación congelada. Esta historia no
  agrega tablas, columnas ni migraciones.
- El mínimo y el máximo se calculan con los votos de la ronda, igual que el promedio de HU-02.
- La condición de D2b del Escenario 4 sale de `docs/PRODUCT-BACKLOG.md`. El texto del issue #42 no
  la incluye, pero mandan el backlog y sus decisiones transversales.
- La aplicación no tiene autenticación: quien actúa se identifica con su usuario en la petición,
  como en HU-02. Si el proyecto tiene más de un Scrum Master, cualquiera puede hacerlo.
- Abrir una nueva ronda se rechaza en una historia "Completada", igual que el inicio de la votación
  en HU-02. Registrar la estimación no mira ese estado: alcanza con la regla D2b.
- No se guarda qué ronda originó la estimación: §4.7 no prevé ese dato y la historia no lo pide.
- La carga directa de Story Points desde HU-15 sigue funcionando igual, también durante una votación.
- El historial lo puede consultar cualquiera, sin indicar quién consulta: solo muestra rondas
  reveladas, cuyos votos ya son visibles para todo el equipo.
- La pantalla es otra sub-issue (#93) y queda fuera de alcance, incluida la vista del historial.
