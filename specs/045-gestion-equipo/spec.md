# Feature Specification: Gestión del equipo del proyecto (HU-16)

**Feature Branch**: `feature/175-gestion-equipo`

**Issue**: #175 (sub-issues #178 a #182)

**Created**: 2026-11-12

**Status**: Draft

**Dependencias**: HU-04 (proyecto e integrantes), HU-20 (autenticación), HU-21 (permisos
por rol, spec `042-hu-21-permisos-por-rol`) y #236 (un único Scrum Master, spec
`040-scrum-master-unico`). Afecta a HU-02 (votación), HU-07 (horas) y HU-03 (métricas).
La migración nueva va después de la 008 de #300 (PR #311).

**Input**: User description: "HU-16 (#175): Gestión del Equipo del Proyecto. Como Scrum
Master quiero quitar integrantes del equipo de un proyecto, y volver a sumarlos,
conservando todo lo que aportaron, para mantener el equipo al día sin perder el historial.
El alta como Product Builder de un usuario registrado ya la cubre HU-21 (escenario 8);
esta historia agrega quitar con baja lógica, reactivar, proteger al único Scrum Master y
responder «no encontrado». Decisión D15 confirmada por el equipo: columna anulable
`integrantes.fecha_baja`, vacía = activo; la reactivación reutiliza la fila."

## Clarifications

### Session 2026-11-12

- Q: Si se quita a un integrante que ya votó en una ronda todavía abierta, ¿su voto cuenta en esa ronda? → A: No. Al dar de baja se descartan sus votos de las rondas no reveladas (eran secretos y nadie los vio), así no cuentan en el promedio ni en las diferencias; las rondas ya reveladas no cambian. Se eligió descartar en vez de marcar el voto porque el resultado de una ronda se recalcula siempre desde los votos guardados y marcarlo exigiría otro cambio de esquema (decidido en `/speckit.plan`, 2026-11-12).
- Q: ¿Cómo se cubre «cambiarle el rol» al Scrum Master (Escenario 4) si no hay una acción de cambio de rol? → A: No se agrega una acción nueva (YAGNI). Se prueba con los caminos existentes (agregar de nuevo al Scrum Master con cualquier rol responde 400 «el integrante ya pertenece al proyecto» sin cambiar el equipo; corregido en `/speckit.plan`: la regla de duplicado existente responde 400, no 409) y se declara que no hay cambio de rol.

## Objetivo

Que el Scrum Master pueda mantener el equipo de un proyecto al día: quitar a quien deja de
participar y volver a sumarlo si regresa, sin que se pierda nada de lo que esa persona
aportó (votos, horas, defectos) ni cambien las métricas ya calculadas. El proyecto conserva
siempre a su único Scrum Master.

**Descripción (backlog)**: Como Scrum Master quiero agregar a un usuario ya registrado como
Product Builder y quitar a un integrante del equipo de un proyecto, conservando todo lo que
aportó, para poder mantener el equipo al día sin perder el historial del proyecto.

## Entradas

- La sesión del Scrum Master del proyecto (HU-20, HU-21).
- Para quitar: el proyecto y el integrante (por su identificador de usuario).
- Para volver a sumar: el proyecto y el nombre del usuario registrado, igual que el alta de
  HU-21.

## Salidas esperadas

- Baja exitosa: el integrante deja de figurar en el equipo; la respuesta confirma la baja
  sin contenido.
- Reactivación exitosa: el integrante vuelve al equipo como Product Builder; la respuesta es
  la misma que la de un alta.
- Rechazos: 404 «no encontrado» (proyecto o integrante inexistentes), 403 (no Scrum
  Master, no integrante; HU-21), 409 con un mensaje claro (quitar al Scrum Master), 400
  «el integrante ya pertenece al proyecto» (agregar de nuevo a un integrante activo,
  incluido el Scrum Master con otro rol).

## Fuera de alcance

- Agregar un usuario registrado como Product Builder: ya lo cubre HU-21 (escenario 8). Esta
  historia solo extiende ese alta para reactivar a quien fue dado de baja.
- Traspasar el rol de Scrum Master: límite declarado hasta después de la entrega.
- Editar el nombre de la persona: pertenece al usuario (HU-17).
- Que un integrante se dé de baja a sí mismo: solo el Scrum Master gestiona el equipo.
- Borrar físicamente a un integrante o sus aportes.
- Mostrar al Product Builder la acción de quitar deshabilitada, con la nota «Solo el Scrum
  Master puede…» (antes US6, escenario 3). La pantalla del equipo todavía no conoce el rol de
  quien la ve. Va con el frontend de HU-21 (#247). Mientras tanto, si un Product Builder pide
  la baja, el servidor responde 403 (FR-012) y la interfaz muestra ese mensaje (FR-016).

## User Scenarios & Testing *(mandatory)*

Los escenarios del backlog (sección HU-16 de `docs/PRODUCT-BACKLOG.md`) se mapean a
historias según su prioridad. El escenario 1 ya está implementado por HU-21 y se mantiene
como regresión.

### User Story 1 - Quitar un integrante conservando su historial (Priority: P1) — Escenario 2

El Scrum Master quita del equipo a un Product Builder. Esa persona deja de figurar en el
equipo, ya no actúa sobre el proyecto y no se la espera en las votaciones nuevas, pero sus
votos, horas y defectos siguen en el historial y en las métricas.

**Why this priority**: es el núcleo de la historia; hoy no hay forma de sacar a nadie del
equipo.

**Independent Test**: con un Product Builder que votó una ronda revelada e imputó horas, el
Scrum Master lo quita; el equipo ya no lo lista, una ronda nueva se revela sin su voto, y el
historial de rondas, las horas y las métricas siguen mostrando sus aportes con su nombre.

**Acceptance Scenarios**:

1. *(Escenario 2)* **Given** que un integrante ya votó una estimación e imputó horas,
   **When** lo quito del equipo, **Then** deja de figurar en el equipo y no cuenta en las
   votaciones nuevas, y sus votos y sus horas siguen en el historial y en las métricas.
2. **Given** que quité a un integrante, **When** esa persona intenta votar, imputar horas o
   abrir el panel del proyecto, **Then** el sistema responde 403 «No sos integrante de este
   proyecto».
3. **Given** que quité a un integrante, **When** esa persona consulta «Mis proyectos»,
   **Then** el proyecto ya no aparece.
4. **Given** que quité a un integrante, **When** se genera el reporte del proyecto, **Then**
   el equipo del reporte no lo incluye y las horas y métricas no cambian.

---

### User Story 2 - El Scrum Master no se quita ni cambia de rol (Priority: P1) — Escenario 4

**Why this priority**: sin esta regla el proyecto puede quedar sin Scrum Master y nadie
podría gestionarlo (HU-21).

**Independent Test**: el Scrum Master intenta quitarse a sí mismo (409) y agregarse de nuevo
como Product Builder (400, integrante duplicado); los dos pedidos se rechazan con un mensaje
claro y el equipo queda igual.

**Acceptance Scenarios**:

1. *(Escenario 4)* **Given** que un proyecto tiene un único Scrum Master activo, **When**
   intento quitarlo o cambiarle el rol, **Then** el sistema rechaza el cambio con un mensaje
   claro y no modifica el equipo.

---

### User Story 3 - Volver a sumar a un integrante dado de baja (Priority: P2) — Escenario 3

**Why this priority**: completa el ciclo; sin ella, quitar a alguien por error no tiene
vuelta atrás.

**Independent Test**: el Scrum Master quita a un Product Builder y lo vuelve a agregar por su
nombre; aparece una sola vez en el equipo, como Product Builder, y sus aportes previos
siguen a su nombre.

**Acceptance Scenarios**:

1. *(Escenario 3)* **Given** que un integrante fue quitado del equipo, **When** lo agrego de
   nuevo al proyecto como Product Builder, **Then** el sistema lo reactiva con ese rol, sin
   duplicarlo ni perder su historial.
2. **Given** que reactivé a un integrante, **When** consulta «Mis proyectos», **Then** el
   proyecto vuelve a aparecer con el rol Product Builder.
3. **Given** un integrante activo, **When** intento agregarlo otra vez, **Then** se aplica la
   regla de integrante duplicado que ya existe (400 «el integrante ya pertenece al
   proyecto»).

---

### User Story 4 - Integrante o proyecto inexistente (Priority: P2) — Escenario 5

**Acceptance Scenarios**:

1. *(Escenario 5)* **Given** que intento modificar un integrante o un proyecto que no
   existen, **When** envío la petición, **Then** el sistema responde «no encontrado».
2. **Given** un usuario que existe pero nunca integró el proyecto, **When** intento quitarlo,
   **Then** el sistema responde «no encontrado».
3. **Given** un integrante ya dado de baja, **When** intento quitarlo de nuevo, **Then** el
   sistema responde «no encontrado» (para el equipo, ya no es integrante).

---

### User Story 5 - Agregar un usuario registrado como Product Builder (regresión) — Escenario 1

Cubierto por HU-21 (escenario 8, `features/` de esa historia). Esta historia no lo cambia:
se verifica que siga en verde.

1. *(Escenario 1)* **Given** que soy el Scrum Master del proyecto y existe un usuario
   registrado que no integra el equipo, **When** lo agrego al proyecto, **Then** queda como
   Product Builder y el sistema lo muestra en el equipo del proyecto, sin ascenso a Scrum
   Master.

---

### User Story 6 - Quitar desde la pantalla del equipo (Priority: P3, frontend)

**Acceptance Scenarios**:

1. **Given** que soy el Scrum Master y veo el equipo, **When** elijo quitar a un Product
   Builder y confirmo, **Then** desaparece de la lista.
2. **Given** que veo el equipo, **When** se muestra el Scrum Master, **Then** no tiene la
   acción de quitar.

El escenario «el Product Builder ve la acción deshabilitada» queda fuera de alcance (ver
«Fuera de alcance»): va con el frontend de HU-21 (#247).

### Edge Cases

- Ronda de votación abierta al momento de la baja: si el integrante no votó, la ronda se
  revela sin esperarlo; si ya votó, su voto se descarta y no cuenta en el resultado de esa
  ronda (promedio y diferencias) ni se muestra en ella (FR-018). Las rondas ya reveladas
  antes de la baja no cambian. Si se lo reactiva con la ronda todavía abierta, puede volver
  a votar.
- Baja y reactivación el mismo día: se permite; la fila vuelve a quedar activa.
- Proyecto cerrado (#300): quitar o reactivar es una escritura y responde 409 «El proyecto
  está cerrado», por la guardia de proyecto cerrado que ya llaman los chequeos de Scrum
  Master.
- Dos pedidos simultáneos sobre el mismo integrante: la baja se aplica una sola vez; el
  segundo responde «no encontrado».
- Un integrante dado de baja que el Scrum Master intenta agregar como Scrum Master: se
  rechaza con el conflicto de #236 (409), igual que a cualquier otro.
- Defectos que registró o resolvió el integrante: se conservan con su nombre.

## Requirements *(mandatory)*

### Reglas de negocio (EARS)

**Ubicuas**

- **FR-001**: El sistema SHALL tratar como integrante del proyecto solo a quien no fue dado
  de baja (integrante activo), para la pertenencia, el rol, el equipo listado, «Mis
  proyectos», la votación y el reporte.
- **FR-002**: El sistema SHALL conservar sin cambios los votos, las horas y los defectos de un
  integrante dado de baja, y seguir mostrándolos con su nombre en el historial de rondas,
  las horas y las métricas.
- **FR-003**: La baja SHALL ser lógica: registra la fecha de baja en el vínculo existente
  y no borra el vínculo ni sus aportes.
- **FR-004**: Cada proyecto SHALL conservar siempre exactamente un Scrum Master activo.

**Basadas en evento**

- **FR-005**: WHEN el Scrum Master quita a un Product Builder activo, el sistema SHALL
  registrar la fecha de baja (el día de hoy) y responder sin contenido.
- **FR-006**: WHEN el Scrum Master agrega a un usuario registrado que fue dado de baja del
  proyecto, el sistema SHALL reactivar el mismo vínculo como Product Builder (borrar la
  fecha de baja) sin crear otro, y responder igual que un alta.
- **FR-007**: WHEN se revela una ronda de votación, el sistema SHALL exigir solo los votos de
  los integrantes activos.

**Basadas en estado**

- **FR-008**: WHILE un integrante está dado de baja, el sistema SHALL responderle 403 «No sos
  integrante de este proyecto» en toda ruta que exige pertenencia (HU-21).
- **FR-018**: WHEN el Scrum Master quita a un integrante, el sistema SHALL descartar, en la
  misma operación, los votos de esa persona en las rondas del proyecto que todavía no se
  revelaron (votos secretos que nadie vio); los votos de rondas ya reveladas SHALL
  conservarse sin cambios.

**No deseadas**

- **FR-009**: IF se intenta quitar al Scrum Master del proyecto, THEN el sistema SHALL
  responder 409 con un mensaje claro (por ejemplo, «No se puede quitar al Scrum Master del
  proyecto») y no modificar el equipo.
- **FR-010**: IF se intenta cambiar el rol del Scrum Master, THEN el sistema SHALL rechazarlo
  con un mensaje claro y no modificar el equipo. No se agrega una acción para cambiar roles:
  el escenario se cubre con los caminos que existen (agregar de nuevo al Scrum Master, con
  cualquier rol, responde 400 «el integrante ya pertenece al proyecto» y el equipo no
  cambia), y se declara que el sistema no ofrece cambio de rol.
- **FR-011**: IF el proyecto no existe, o el usuario no es integrante activo del proyecto,
  THEN el sistema SHALL responder 404 «no encontrado» sin modificar nada.
- **FR-012**: IF quien pide la baja no es el Scrum Master del proyecto, THEN el sistema SHALL
  responder 403 según HU-21 (no integrante, o «Solo el Scrum Master puede…») sin modificar
  nada.
- **FR-013**: Ante varias condiciones de error a la vez, el sistema SHALL seguir el orden de
  HU-21: 401 → 404 proyecto → 403 no integrante → 403 no Scrum Master → 409 proyecto
  cerrado → 404 integrante → 409 Scrum Master.
- **FR-014**: IF dos pedidos de baja del mismo integrante llegan a la vez, THEN el sistema
  SHALL aplicar solo uno; la regla del Scrum Master y la baja SHALL verificarse en la misma
  operación atómica que el cambio.

### Requisitos de la interfaz *(frontend)*

- **FR-015**: La pantalla del equipo MUST ofrecer al Scrum Master la acción de quitar en cada
  Product Builder, con confirmación, y no ofrecerla en el Scrum Master.
- **FR-016**: La interfaz MUST mostrar el mensaje del servidor si la baja se rechaza.

### Requisitos de datos

- **FR-017**: El cambio de esquema (fecha de baja anulable en el vínculo de integrante,
  formato día `YYYY-MM-DD`, decisión D15) MUST ir en un solo PR que actualiza
  `MODELO-DE-DATOS.md` y agrega la migración numerada.

### Key Entities

- **Integrante del proyecto**: vínculo usuario–proyecto con rol (HU-04, #236). Gana la fecha
  de baja: vacía = activo; con fecha = dado de baja. Un usuario tiene un solo vínculo por
  proyecto, también después de la baja y la reactivación.
- **Voto, registro de horas, defecto**: aportes de un usuario en el proyecto. No cambian con
  la baja.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Los 5 escenarios del backlog para HU-16 pasan automatizados (Godog para la API,
  pruebas de componente para la interfaz).
- **SC-002**: Después de quitar a un integrante, el 100 % de sus votos, horas y defectos sigue
  visible con su nombre, y las métricas del proyecto dan exactamente los mismos valores que
  antes de la baja.
- **SC-003**: En ningún caso un proyecto queda sin Scrum Master activo.
- **SC-004**: Después de quitar y reactivar a un integrante, figura exactamente una vez en el
  equipo.
- **SC-005**: Todos los escenarios de Godog existentes siguen en verde.

## Assumptions

- **Fuente de verdad**: la sección HU-16 de `docs/PRODUCT-BACKLOG.md` y el issue #175.
- **D15 confirmada por el equipo (2026-11-12)**: fecha de baja anulable en el vínculo de
  integrante; la reactivación reutiliza el mismo vínculo. Se cierra la D15 en
  `MODELO-DE-DATOS.md`.
- La reactivación usa el mismo alta de HU-21 (por nombre de un usuario registrado); no hay
  una acción aparte.
- Quien queda dado de baja pierde el acceso de integrante a ese proyecto, pero conserva su
  cuenta y sus otros proyectos.
- El listado de usuarios registrados (HU-21) no cambia.
- La guardia de proyecto cerrado (#300) se aplica sola si ese PR está mergeado al
  implementar; si no, se agrega cuando entre.
