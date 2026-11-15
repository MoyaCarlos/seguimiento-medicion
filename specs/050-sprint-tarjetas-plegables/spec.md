# Feature Specification: Tarjetas plegables y confirmación de estado en el Sprint (#219)

**Feature Branch**: `feature/219-sprint-tarjetas-confirmacion`

**Issue**: #219

**Created**: 2026-11-15

**Status**: Implementada (PR #327, pendiente de review)

**Relacionado**: PR #217 (HU-11 frontend), spec 027, #215, #218. Solo frontend; no cambia la
API, el esquema ni las reglas de negocio de HU-11.

**Input**: User description: "Sprint Backlog: historias plegables como tarjetas y confirmación
antes de iniciar o completar. En su estado normal, cada historia asignada a un Sprint muestra
solo su título; al hacer clic se despliegan los datos adicionales, incluido el botón de
Iniciar o Completar. La fila plegada muestra además una marca chica con el estado y los Story
Points. Un clic accidental en Iniciar o Completar no debe cambiar el estado: se pide una
confirmación en la propia pantalla, preferentemente con `<details>`/`<summary>` y sin
dependencias nuevas."

## Objetivo

Reducir la altura de la tarjeta de cada Sprint y evitar cambios de estado accidentales, sin
tocar la API ni las reglas de HU-11. Las historias se leen como tarjetas plegables y el cambio
de estado exige una confirmación explícita en la pantalla.

## Entradas y salidas

- **Entrada**: la lista de historias que devuelve `GET /sprints/{id}/backlog`, con título,
  prioridad, estado, Story Points y fecha de finalización; y los clics o el teclado sobre una
  tarjeta y sobre los controles de estado.
- **Salida**: una tarjeta plegada por historia con título, marca de estado y Story Points; al
  desplegarla, prioridad, fecha de finalización cuando corresponde, el control de estado y el
  esfuerzo. La confirmación muestra «Confirmar» y «Cancelar» antes de enviar el `PUT` de HU-11.
- **Sin cambios**: la API, el esquema, las transiciones de estado y sus mensajes de error.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Leer el Sprint como tarjetas plegables (Priority: P1)

Quien consulta un Sprint ve cada historia como una tarjeta plegada con su título, una marca
chica de estado y sus Story Points. Al abrirla, aparecen la prioridad, la fecha de finalización
si está completada, el control para cambiar el estado y el esfuerzo real.

**Independent Test**: renderizar un Sprint con historias en distintos estados; comprobar que
cada tarjeta está cerrada, que el resumen muestra título, estado y Story Points, y que al
activar el resumen se abren los datos adicionales y las acciones.

**Acceptance Scenarios**:

1. **Given** una historia asignada a un Sprint, **When** se muestra la tarjeta del Sprint,
   **Then** la historia aparece plegada con su título, la marca de estado y sus Story Points.
2. **Given** una historia plegada, **When** se activa su resumen con el ratón o el teclado,
   **Then** se despliegan la prioridad, la fecha de finalización si está completada, el
   control de estado y el esfuerzo real.
3. **Given** una historia completada, **When** se muestra plegada, **Then** la marca dice
   «Completada» y la fecha de finalización queda dentro del detalle desplegado.
4. **Given** una historia sin estimación, **When** se muestra plegada, **Then** en lugar de
   Story Points se lee «Sin estimar».

---

### User Story 2 - Confirmar antes de iniciar o completar (Priority: P1)

Quien cambia el estado de una historia pulsa primero «Iniciar» o «Completar» y luego confirma
en la propia pantalla. Un clic en «Cancelar», en Escape o en cualquier otro lugar no envía
nada; mientras el pedido está en curso no se puede enviar dos veces.

**Independent Test**: pulsar la acción de una historia y comprobar que todavía no se envió
ningún pedido; confirmar y comprobar el `PUT`; repetir cancelando con «Cancelar», con Escape y
con un clic fuera; y comprobar el doble envío con la respuesta pendiente.

**Acceptance Scenarios**:

1. **Given** una historia «Nueva» del Sprint activo, **When** pulso «Iniciar», **Then** se
   muestra una confirmación con «Confirmar» y «Cancelar» y no se envía ningún pedido.
2. **Given** una confirmación abierta, **When** pulso «Confirmar», **Then** se envía el cambio
   de estado de HU-11 y la tarjeta refleja el resultado.
3. **Given** una confirmación abierta, **When** pulso «Cancelar» o Escape, **Then** la
   confirmación se cierra y no se envía ningún pedido.
4. **Given** una confirmación abierta, **When** pulso fuera de la confirmación, **Then** no se
   envía ningún pedido.
5. **Given** un pedido de cambio de estado en curso, **When** se vuelve a pulsar «Confirmar»,
   **Then** no se envía un segundo pedido y los botones quedan deshabilitados.
6. **Given** una historia «En progreso», **When** pulso «Completar» y confirmo, **Then** se
   envía el cambio a «Completada» con la misma protección y confirmación.

### Edge Cases

Redactados en formato EARS.

- **WHEN** una historia no está completada, **THEN** su tarjeta plegada **shall** mostrar el
  título, la marca de estado y los Story Points, sin la prioridad ni la fecha de finalización.
- **WHEN** una historia está completada, **THEN** la tarjeta plegada **shall** mostrar la marca
  «Completada» y **shall** mantener la fecha de finalización dentro del detalle.
- **IF** la historia no tiene estimación, **THEN** la tarjeta **shall** mostrar «Sin estimar»
  en lugar de Story Points.
- **WHEN** se activa el resumen de una historia, **THEN** su detalle **shall** desplegarse sin
  enviar pedidos ni cambiar el estado.
- **WHEN** se pulsa «Iniciar» o «Completar», **THEN** el sistema **shall** pedir confirmación y
  **shall not** enviar ningún pedido todavía.
- **WHEN** se pulsa «Confirmar», **THEN** el sistema **shall** enviar el cambio una sola vez y
  **shall** mantener el botón deshabilitado hasta que termine el pedido.
- **IF** se pulsa «Cancelar», Escape o cualquier otro lugar durante la confirmación, **THEN**
  el sistema **shall not** enviar ningún pedido.
- **IF** el proyecto está cerrado o el usuario no puede escribir, **THEN** la acción de estado
  **shall** quedar deshabilitada con la nota existente y **shall not** abrir la confirmación.
- **IF** el cambio de estado se rechaza, **THEN** la confirmación **shall** cerrarse, **shall**
  mostrarse el mensaje de la API y la historia **shall** quedar como estaba.

## Requirements *(mandatory)*

### Reglas de comportamiento (EARS)

**Ubicuas**

- **FR-001**: Cada historia del Sprint MUST mostrarse como una tarjeta plegable implementada con
  los elementos nativos `<details>` y `<summary>`.
- **FR-002**: El resumen plegado MUST mostrar el título, una marca con el estado («Nueva», «En
  progreso» o «Completada») y los Story Points o «Sin estimar».
- **FR-003**: El detalle desplegado MUST mostrar la prioridad con su etiqueta, la fecha de
  finalización cuando la historia está completada, el control de estado cuando corresponde y
  el esfuerzo real.
- **FR-004**: El cambio de estado MUST pedir confirmación en la propia pantalla con «Confirmar»
  y «Cancelar», sin usar `window.confirm`.
- **FR-005**: Mientras el pedido está en curso, «Confirmar» y «Cancelar» MUST quedar
  deshabilitados y MUST NOT enviarse más de un pedido.
- **FR-006**: La funcionalidad MUST ser solo de frontend, sin dependencias nuevas ni cambios de
  API, esquema o reglas de negocio.

**Basadas en evento**

- **WHEN** se activa el resumen, el sistema SHALL alternar el detalle de esa historia sin enviar
  pedidos.
- **WHEN** se pulsa «Iniciar» o «Completar», el sistema SHALL mostrar la confirmación y SHALL
  NOT enviar el cambio.
- **WHEN** se pulsa «Confirmar», el sistema SHALL enviar el cambio de estado de HU-11 una sola
  vez y SHALL reflejar la historia devuelta.
- **WHEN** se pulsa «Cancelar» o Escape, el sistema SHALL cerrar la confirmación sin enviar
  nada.
- **WHEN** se pulsa fuera de la confirmación, el sistema SHALL NOT enviar nada.

**No deseadas**

- **IF** el proyecto está cerrado o el usuario no puede escribir, THEN la acción SHALL quedar
  deshabilitada con la nota existente y SHALL NOT abrir la confirmación.
- **IF** la API rechaza el cambio, THEN la confirmación SHALL cerrarse, el mensaje de la API
  SHALL mostrarse y la historia SHALL quedar como estaba.

### Requisitos de la interfaz

- **FR-007**: La confirmación MUST ser accesible: texto de pregunta, botones con nombre,
  foco inicial en «Cancelar» y cierre con Escape.
- **FR-008**: Las tarjetas MUST conservar los datos y el orden que devuelve la API, y MUST
  mantener las pruebas de HU-11 y de esfuerzo en verde.

## Fuera de alcance

- Cambios en la API, el esquema, `MODELO-DE-DATOS.md` o las transiciones de HU-11.
- Transiciones inversas o reapertura de historias completadas.
- Un diálogo modal nuevo del navegador o una dependencia de terceros.
- Cambios en la asignación de historias a un Sprint (HU-06) o en el esfuerzo (HU-07).

## Success Criteria *(mandatory)*

- **SC-001**: El 100 % de las historias de un Sprint se muestran plegadas con título, estado y
  Story Points, y despliegan sus datos adicionales con ratón o teclado.
- **SC-002**: El 100 % de los clics en «Iniciar» o «Completar» pide confirmación y no envía el
  cambio hasta pulsar «Confirmar».
- **SC-003**: «Cancelar», Escape y un clic fuera no envían ningún pedido.
- **SC-004**: Con la respuesta pendiente se envía un único pedido y los botones quedan
  deshabilitados.
- **SC-005**: Las pruebas de componente de `EstadoHistoria` y `SprintBacklog` pasan en Vitest,
  junto con el resto de la suite, el lint y el build.

## Assumptions

- #217 ya está en `main` y dejó el control de estado de HU-11 en `EstadoHistoria`; esta mejora
  cambia solo su interacción y la presentación de `SprintBacklog`.
- `<details>` y `<summary>` son accesibles con teclado y no requieren dependencias nuevas.
- La confirmación reutiliza el patrón en línea de `ZonaCierreProyecto` (pregunta, «Confirmar»,
  «Cancelar», foco en «Cancelar» y Escape), que ya usa el proyecto.
- No hay cambios de backend ni de contrato: el `PUT /backlog/{id}/estado` y sus mensajes
  siguen igual.
- Solo se tocan archivos del frontend (`web/src/features/sprints/`) y la documentación.
