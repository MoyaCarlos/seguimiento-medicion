# Feature Specification: Pantalla de Sprints (HU-05 Frontend, Issue #26)

**Feature Branch**: `feature/HU-05-frontend`

**Created**: 2026-10-24

**Status**: Implementada (US1, US2, US3 y US4), pendiente de revisión por PR.

**Input**: User description: "HU-05 Frontend (Issue #26, sub-issue de #5): pantalla de Sprints en `web/src/features/sprints/`. Crear un Sprint con un botón, iniciar un Sprint Pendiente con un formulario (Sprint Goal y fechas), cerrar un Sprint Activo con un botón, y mostrar los Sprints agrupados por estado. Consume solo la API que existe hoy; no hay listado de Sprints (#110)."

## Clarifications

### Session 2026-10-24

- Q: ¿Qué pasa con un Sprint que existe en la base pero la pantalla no conoce (por ejemplo, tras recargar)? → A: Obsoleto desde #110: la pantalla carga los Sprints del proyecto al abrirse (ver [`specs/015-hu-05-consultar-sprints`](../015-hu-05-consultar-sprints/spec.md)). Originalmente se aceptaba la limitación y se mostraba un texto que la explicaba; ese texto se quitó.
- Q: ¿Cada Sprint muestra su identificador? → A: Sí, como "Sprint #<id>", para poder cruzarlo con la API y distinguir varios Sprints Finalizados. (Decisión del Scrum Master)
- Q: Cuando falla "Crear Sprint", iniciar o cerrar con un error que no pertenece a un campo, ¿dónde se muestra el aviso? → A: Un aviso por Sprint, arriba de su formulario de inicio o de su botón "Cerrar Sprint"; el error de "Crear Sprint" va en un aviso arriba de la pantalla. Cada aviso permanece hasta el próximo envío de su propia acción.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Crear un Sprint (Priority: P1)

Como Scrum Master, quiero crear un Sprint nuevo para mi proyecto con un solo botón, sin completar ningún campo, para dejarlo preparado antes de iniciarlo.

**Why this priority**: Sin un Sprint creado no se puede iniciar ni cerrar nada; es el punto de entrada del ciclo.

**Independent Test**: Abrir `/projects/<id>/sprints` con el backend corriendo, pulsar el botón de crear y comprobar que aparece un Sprint en la sección "Pendiente", sin Goal ni fechas.

**Acceptance Scenarios**:

1. **Given** la pantalla abierta para un proyecto existente, **When** se pulsa "Crear Sprint", **Then** aparece un Sprint nuevo en la sección "Pendiente", sin Sprint Goal ni fechas.
2. **Given** una creación en curso, **When** todavía no hubo respuesta, **Then** el botón queda deshabilitado y se ve el indicador de carga.
3. **Given** que el proyecto ya tiene un Sprint Pendiente, **When** se pulsa "Crear Sprint", **Then** se muestra en un aviso arriba el mensaje que devuelve la API y no se agrega ningún Sprint.
4. **Given** un id de proyecto inexistente en la URL, **When** se pulsa "Crear Sprint", **Then** se muestra en el aviso el mensaje "proyecto no encontrado" de la API.

---

### User Story 2 - Iniciar un Sprint Pendiente (Priority: P1)

Como Scrum Master, quiero iniciar un Sprint Pendiente definiendo su Sprint Goal y sus fechas de inicio y fin, para que el equipo sepa que la iteración está en marcha.

**Why this priority**: Es la mitad del ciclo iterativo que pide la historia y el único flujo con formulario y errores por campo.

**Independent Test**: Con un Sprint Pendiente en pantalla, completar Goal y fechas válidas, enviar, y comprobar que el Sprint pasa a la sección "Activo" mostrando su Goal y fechas.

**Acceptance Scenarios**:

1. **Given** un Sprint Pendiente en pantalla, **When** se completan Goal, fecha de inicio y fecha de fin y se envía, **Then** el Sprint pasa a la sección "Activo" mostrando Goal y fechas.
2. **Given** que la API rechaza un dato con un campo (por ejemplo, una fecha de fin no posterior a la de inicio), **When** se envía, **Then** el mensaje aparece junto al campo correspondiente y el Sprint sigue Pendiente.
3. **Given** que la API rechaza con un error sin campo visible (por ejemplo, "el proyecto ya tiene un Sprint activo", que viene asociado al campo `estado`), **When** se envía, **Then** el mensaje aparece en un aviso arriba del formulario, que permanece hasta el próximo envío.
4. **Given** un envío en curso, **When** todavía no hubo respuesta, **Then** el botón de envío queda deshabilitado y se ve el indicador de carga.
5. **Given** que hay un Sprint Activo o Finalizado, **When** se mira la pantalla, **Then** ese Sprint no ofrece el formulario de inicio.

---

### User Story 3 - Cerrar un Sprint Activo (Priority: P1)

Como Scrum Master, quiero cerrar el Sprint activo con un botón al terminar la iteración, para marcarlo como Finalizado.

**Why this priority**: Completa el ciclo de vida del Sprint.

**Independent Test**: Con un Sprint Activo en pantalla, pulsar "Cerrar Sprint" y comprobar que pasa a la sección "Finalizado".

**Acceptance Scenarios**:

1. **Given** un Sprint Activo en pantalla, **When** se pulsa "Cerrar Sprint", **Then** el Sprint pasa a la sección "Finalizado".
2. **Given** que la API rechaza el cierre, **When** se pulsa el botón, **Then** el mensaje de la API se muestra en un aviso y el Sprint permanece Activo.
3. **Given** un cierre en curso, **When** todavía no hubo respuesta, **Then** el botón queda deshabilitado y se ve el indicador de carga.
4. **Given** un Sprint Pendiente o Finalizado, **When** se mira la pantalla, **Then** ese Sprint no ofrece el botón "Cerrar Sprint".

---

### User Story 4 - Ver los Sprints agrupados por estado (Priority: P2)

Como Scrum Master, quiero ver los Sprints en tres secciones (Pendiente, Activo, Finalizado), para saber de un vistazo en qué estado está cada iteración.

**Why this priority**: Da el contexto de las acciones anteriores, pero cada acción se usa aun sin una vista agrupada pulida.

**Independent Test**: Crear, iniciar y cerrar Sprints y comprobar que cada uno aparece solo en la sección de su estado actual.

**Acceptance Scenarios**:

1. **Given** la pantalla recién abierta, **When** no se creó ningún Sprint en la sesión, **Then** se ven las tres secciones, vacías, con un texto que lo indica.
2. **Given** Sprints en distintos estados, **When** se mira la pantalla, **Then** cada Sprint aparece únicamente en la sección de su estado.
3. **Given** Sprints creados en la sesión, **When** se recarga la página, **Then** se ven los Sprints del proyecto cargados desde la API (obsoleto: antes se veían vacías; ver spec 015).

---

### Edge Cases

Redactados en formato EARS.

- **IF** el backend no está corriendo, **THEN** la pantalla **shall** mostrar el error de conexión en el aviso de la acción que falló, sin quedar cargando indefinidamente.
- **IF** la API responde con un error no interpretable (por ejemplo, un 500), **THEN** la pantalla **shall** mostrar el mensaje genérico en el aviso de la acción que falló, sin romperse.
- **IF** la API rechaza un envío con un `campo` que el formulario no tiene (por ejemplo `estado`), **THEN** la pantalla **shall** tratarlo como error sin campo visible y mostrarlo en el aviso arriba del formulario.
- **IF** hay un Sprint Pendiente o Activo en la base que la pantalla no conoce, **THEN** crear otro Sprint o iniciar uno **shall** mostrar el error que devuelve la API, y la pantalla ya muestra los Sprints existentes gracias al listado (#110).
- **WHEN** empieza un nuevo envío, la pantalla **shall** borrar los errores del envío anterior.
- **WHEN** hay una operación en curso sobre un Sprint, la pantalla **shall** impedir una segunda operación sobre ese mismo Sprint hasta que termine.
- **IF** el id del proyecto no está en la URL o no es válido, **THEN** la pantalla **shall** mostrar el error que devuelva la API al operar, sin validar ella misma el id.
- **WHEN** exista un listado de Sprints (#110), la pantalla **shall** poder cargarlo en el lugar previsto sin rehacer las secciones ni las acciones.
- **WHEN** el contrato de cierre cambie (#117, cierre con arrastre de historias), el cambio **shall** quedar limitado al componente del botón "Cerrar Sprint" y a su llamada a la API.

## Requirements *(mandatory)*

### Functional Requirements

**Crear**

- **FR-001**: La pantalla **shall** ofrecer un botón "Crear Sprint" sin campos, que envía solo el id del proyecto tomado de la URL.
- **FR-002**: **WHEN** la API confirma la creación, la pantalla **shall** agregar el Sprint devuelto a la sección "Pendiente".

**Iniciar**

- **FR-003**: Cada Sprint Pendiente **shall** ofrecer un formulario con Sprint Goal, fecha de inicio y fecha de fin (formato `AAAA-MM-DD`) y un botón para iniciarlo.
- **FR-004**: **WHEN** la API confirma el inicio, la pantalla **shall** mostrar el Sprint devuelto en la sección "Activo", con su Sprint Goal y fechas.

**Cerrar**

- **FR-005**: Cada Sprint Activo **shall** ofrecer un botón "Cerrar Sprint".
- **FR-006**: **WHEN** la API confirma el cierre, la pantalla **shall** mostrar el Sprint devuelto en la sección "Finalizado".
- **FR-007**: El botón "Cerrar Sprint" y su llamada a la API **shall** vivir aislados en su propio componente y archivo, para que un cambio de contrato (#117) no toque el resto de la pantalla.

**Mostrar**

- **FR-008**: La pantalla **shall** mostrar tres secciones fijas (Pendiente, Activo, Finalizado), siempre visibles, y cada Sprint **shall** aparecer solo en la de su estado.
- **FR-009** *(reemplazado por la spec 015, FR-F01)*: ~~La pantalla muestra solo los Sprints de la sesión.~~ Ahora recupera los Sprints del proyecto al abrirse y al recargar.
- **FR-010** *(cumplido por la spec 015)*: el punto de carga inicial es un `useEffect` que pide `GET /projects/{id}/sprints`; las secciones y acciones no cambiaron.

- **FR-018** *(retirado por la spec 015, FR-F05)*: el texto de limitación "Los Sprints existentes no se listan todavía…" ya no se muestra, porque el listado existe.
- **FR-019**: Cada Sprint mostrado **shall** exhibir su identificador como "Sprint #<id>".

**Errores y carga**

- **FR-011**: La pantalla **shall** mostrar el mensaje que devuelve la API tal cual, sin replicar validaciones del backend (fechas, estado, existencia del proyecto).
- **FR-012**: **WHEN** la API devuelve un error con campo que corresponde a un campo del formulario de inicio, la pantalla **shall** mostrarlo junto a ese campo.
- **FR-013**: **WHEN** el error no tiene campo, tiene un campo que el formulario no muestra, es de conexión o no es interpretable, la pantalla **shall** mostrarlo en un aviso propio: el de iniciar o cerrar un Sprint, arriba del formulario o del botón de ese Sprint; el de crear, arriba de la pantalla. Cada aviso permanece hasta el próximo envío de su propia acción y no se muestra en otros Sprints.
- **FR-014**: **WHILE** una operación está en curso, la pantalla **shall** mostrar el indicador de carga y deshabilitar el botón de esa operación.

**Alcance técnico**

- **FR-015**: La implementación **shall** modificar solo `web/src/features/sprints/` y `specs/011-*/`; todas las llamadas **shall** pasar por el cliente compartido y los componentes comunes se **shall** reutilizar de `web/src/compartido/` sin editarlos.
- **FR-016**: La pantalla **shall** seguir exportando `Component` desde `web/src/features/sprints/index.tsx`, para que el registro de rutas no cambie.
- **FR-017**: La interfaz **shall** estar en español.
- **FR-020**: La pantalla **shall** conservar el contrato con la base del frontend: un encabezado de nivel 1 con el texto exacto "Sprints" y el id del proyecto de la URL en un único elemento `<strong>` dentro de "Proyecto actual: <id>"; ningún otro elemento muestra solo ese número. Lo verifican los tests existentes de `web/src/rutas.test.tsx`, que no se modifican.

### Key Entities

- **Sprint**: iteración del proyecto, según lo devuelve la API: identificador, proyecto, Sprint Goal y fechas (ausentes mientras está Pendiente) y estado (Pendiente, Activo, Finalizado).
- **Sprints de la sesión**: los Sprints que la pantalla conoce porque los creó o modificó en esta visita; se pierden al recargar.
- **Error de la API**: estado, mensaje y, a veces, campo (definido por la base del frontend, spec 007).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Una persona crea un Sprint con un solo clic y lo ve en "Pendiente" en menos de 2 segundos con el backend local.
- **SC-002**: Una persona crea, inicia y cierra un Sprint de ejemplo, en ese orden, en menos de 2 minutos sin ayuda.
- **SC-003**: El 100 % de los errores que devuelve la API en esos tres flujos se ve en pantalla, junto al campo o en el aviso; ninguno se pierde ni deja la pantalla en blanco.
- **SC-004**: El 100 % de los Sprints mostrados está en la sección de su estado.
- **SC-005**: Un cambio en el contrato de cierre toca un solo componente de la pantalla.
- **SC-006**: El PR modifica únicamente `web/src/features/sprints/` y `specs/011-*/`; `npm test`, `npm run lint` y `npm run build` terminan sin errores.

## Assumptions

- Depende de la base del frontend (spec 007, Issue #108): cliente de la API, piezas comunes y registros de rutas ya existen; esta historia no los modifica.
- Endpoints disponibles (spec 003, contrato `openapi.yaml`): `POST /sprints` con `{proyecto_id}`, `POST /sprints/{id}/iniciar` con `{sprint_goal, fecha_inicio, fecha_fin}` y `POST /sprints/{id}/cerrar` sin cuerpo; todos devuelven el Sprint actualizado.
- No existe `GET` de Sprints (#110): por eso solo se ven los Sprints de la sesión. Tampoco hay listado de historias, así que la pantalla no muestra historias ni el efecto del arrastre al cerrar.
- El proyecto actual sale de la URL (`/projects/:id/sprints`); la pantalla no valida que exista, lo hace la API.
- Los controles de fecha entregan `AAAA-MM-DD`; si la API rechaza un valor, se muestra su mensaje.
- Un Sprint Pendiente aún no tiene Goal ni fechas, y la pantalla lo muestra sin ellos.
- Las reglas (un solo Pendiente, un solo Activo, fecha de fin posterior, cierre solo de Activos) las aplica la API (spec 003); la pantalla no las duplica.
- El cierre con arrastre de historias (#117) puede cambiar el contrato de cierre; esta spec cubre el contrato actual.
- Ciclo RED/GREEN/REFACTOR: aplica a la lógica de la pantalla (mapeo de errores a campo o aviso, estados de los botones, agrupación por estado); las piezas sin lógica no necesitan RED. BDD con Godog no aplica al frontend.
- Fuera de alcance: listado de Sprints (#110), asignar historias a un Sprint (HU-06), métricas del Sprint, autenticación, cambios en el backend y estilos elaborados.
