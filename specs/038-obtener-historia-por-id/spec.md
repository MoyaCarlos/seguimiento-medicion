# Feature Specification: Consultar una historia por su identificador

**Feature Branch**: `feature/222-obtener-historia`

**Issue**: #222 (habilita las mejoras de la pantalla de votación #223, #224 y #225)

**Created**: 2026-11-08

**Status**: Draft

**Input**: User description: "Issue #222: GET /backlog/{id}. Consultar una historia por su id; habilita mejoras en la pantalla de votación (#223 a #225). Solo backend: sin cambio de esquema, sin tocar el frontend y sin exigir sesión (HU-21 impone permisos en el Sprint 5). 200 con la historia con la misma forma que el resto de la API de backlog; 404 si no existe; 400 con campo `id` si el identificador no es un entero positivo. Reutilizar el caso de uso ya existente; el handler solo adapta HTTP."

## Objetivo

Que quien consume la API (hoy, la pantalla de votación) pueda obtener **una sola historia** por su identificador, sin traer el backlog completo del proyecto ni el de un Sprint para buscarla del lado del cliente.

## Entradas y salidas

- **Entrada**: el identificador de la historia, que viaja en la ruta. No hay cuerpo ni parámetros de consulta.
- **Salida (éxito)**: la historia con los mismos campos y formato que ya devuelve el resto de la API de backlog: `id`, `proyecto_id`, `titulo`, `descripcion`, `criterios_aceptacion`, `prioridad`, `estado`, `valor_negocio`, `estimacion_sp`, `sprint_id` y `fecha_finalizacion`.
- **Salida (error)**: un mensaje de error en el formato habitual de la API, para historia inexistente o identificador inválido.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Consultar una historia del Product Backlog (Priority: P1)

Quien consume la API pide una historia por su identificador y recibe sus datos completos. Cuando la historia todavía no está en ningún Sprint, `sprint_id` viene en `null`.

**Why this priority**: es el caso base y lo que desbloquea #223 a #225; sin él no hay feature.

**Independent Test**: crear una historia en un proyecto, consultarla por su id y comprobar que todos los campos coinciden con los cargados y que `sprint_id` es `null`.

**Acceptance Scenarios**:

1. **Given** una historia del Product Backlog sin Sprint asignado, **When** se consulta por su id, **Then** se responde 200 con todos los campos de la historia y `sprint_id` en `null`.
2. **Given** la misma historia, **When** se compara la respuesta con la que entrega el listado del backlog del proyecto para esa historia, **Then** ambas tienen exactamente la misma forma y valores.

---

### User Story 2 - Consultar una historia asignada a un Sprint o ya completada (Priority: P2)

La respuesta refleja la situación actual de la historia: si está asignada a un Sprint, informa cuál (`sprint_id`); si fue completada, informa su fecha de finalización.

**Why this priority**: la pantalla de votación necesita saber en qué Sprint está la historia y si ya se completó; es la variante más frecuente después del caso base.

**Independent Test**: asignar una historia a un Sprint y consultarla; completar otra y consultarla.

**Acceptance Scenarios**:

1. **Given** una historia asignada a un Sprint, **When** se consulta por su id, **Then** se responde 200 y `sprint_id` es el identificador de ese Sprint.
2. **Given** una historia completada, **When** se consulta por su id, **Then** se responde 200, el estado es el de completada y `fecha_finalizacion` trae su fecha en el mismo formato que el resto de la API de backlog.
3. **Given** una historia aún no completada, **When** se consulta por su id, **Then** `fecha_finalizacion` viene en `null`.

---

### User Story 3 - Recibir un error claro ante un id inexistente o inválido (Priority: P2)

Quien consume la API distingue "esa historia no existe" de "el identificador que mandé está mal formado", con el mismo estilo de error que el resto de la API.

**Why this priority**: evita que la pantalla trate un error como si fuera una historia vacía y sigue la convención ya establecida para los demás endpoints de backlog.

**Independent Test**: consultar un id que no existe y luego uno que no es un entero positivo.

**Acceptance Scenarios**:

1. **Given** que no existe una historia con el id pedido, **When** se consulta, **Then** se responde 404 con el mensaje `historia no encontrada`.
2. **Given** un id que no es un entero positivo (por ejemplo `abc`, `0`, `-3` o `1.5`), **When** se consulta, **Then** se responde 400 indicando el campo `id` y el mensaje `identificador inválido`.

---

### Edge Cases

- Id `0` o negativo: se rechaza como identificador inválido (400), no como inexistente (404).
- Id numérico válido pero muy grande (que no entra en el rango de enteros admitido): se trata como identificador inválido (400), nunca como un error interno.
- Id con ceros a la izquierda (por ejemplo `007`): se interpreta como el entero 7.
- Una historia de otro proyecto se devuelve igual si el id existe: el identificador es único en todo el sistema y esta consulta no filtra por proyecto.
- La consulta no modifica nada: repetirla devuelve el mismo resultado mientras la historia no cambie.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir consultar una historia por su identificador mediante `GET /backlog/{id}`.
- **FR-002**: Si la historia existe, el sistema MUST responder 200 con sus campos `id`, `proyecto_id`, `titulo`, `descripcion`, `criterios_aceptacion`, `prioridad`, `estado`, `valor_negocio`, `estimacion_sp`, `sprint_id` y `fecha_finalizacion`.
- **FR-003**: La respuesta MUST tener la misma forma, nombres de campo y formatos (incluida la fecha) que la historia devuelta por el resto de la API de backlog, sin campos nuevos ni omitidos.
- **FR-004**: `sprint_id` MUST ser `null` cuando la historia no está asignada a un Sprint, y el identificador del Sprint cuando sí lo está.
- **FR-005**: `fecha_finalizacion` MUST ser `null` cuando la historia no está completada, y su fecha cuando sí lo está.
- **FR-006**: Si no existe una historia con ese id, el sistema MUST responder 404 con el mensaje `historia no encontrada`.
- **FR-007**: Si el id de la ruta no es un entero positivo, el sistema MUST responder 400 con el campo `id` y el mensaje `identificador inválido`, sin consultar si la historia existe.
- **FR-008**: La consulta MUST ser de solo lectura: no modifica historias, Sprints ni ningún otro dato.
- **FR-009**: La consulta MUST funcionar sin sesión ni credenciales, igual que el resto de los endpoints de backlog hoy; los permisos los incorpora HU-21 en el Sprint 5.
- **FR-010**: La lógica de "obtener una historia o informar que no existe" MUST reutilizar la que ya existe y usan hoy la edición de historias y la imputación de horas; el endpoint solo traduce entre HTTP y esa lógica, sin duplicar reglas.

### Key Entities

- **Historia de usuario (item del backlog)**: unidad de trabajo del proyecto. Atributos relevantes: identificador, proyecto al que pertenece, título, descripción, criterios de aceptación, prioridad, estado, valor de negocio, estimación en Story Points, Sprint asignado (opcional) y fecha de finalización (opcional). Esta feature solo la lee; no cambia su definición.

## Fuera de alcance

- Cambios en la pantalla de votación (se resuelven en #223 a #225).
- Listados de historias (ya existen: backlog del proyecto y backlog del Sprint).
- Cualquier cambio de esquema o de modelo de datos.
- Control de permisos o de sesión (HU-21, Sprint 5).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Para el 100 % de las historias existentes, la consulta por id devuelve datos idénticos a los de esa misma historia en los listados de backlog (verificado sobre historia sin Sprint, asignada a un Sprint y completada).
- **SC-002**: El 100 % de las consultas con id inexistente o mal formado reciben un error con el estado y el mensaje definidos (404 y 400), y ninguna termina en un error interno del servidor.
- **SC-003**: Los cinco escenarios BDD de esta feature (historia del Product Backlog, historia en Sprint, historia completada, historia inexistente, identificador inválido) pasan, y las suites existentes siguen en verde sin modificaciones de comportamiento.
- **SC-004**: La pantalla de votación puede obtener una historia con una sola consulta, sin descargar un listado completo.
- **SC-005**: La ruta queda documentada con un ejemplo ejecutable en la guía de pruebas con curl del Sprint 2.

## Assumptions

- No hay duda de modelo: la historia ya tiene todos los campos pedidos y no hace falta tocar el esquema.
- El mensaje exacto `identificador inválido` es un requisito de este issue. Hoy el helper compartido que valida ids devuelve otro texto (`el identificador debe ser un entero positivo`) en los demás endpoints; esta feature no cambia ese comportamiento existente, y cómo cumplir el texto pedido sin alterar los otros endpoints se decide en el plan.
- El endpoint no valida ni filtra por proyecto: el id de historia identifica una única historia en todo el sistema.
- No se agregan límites de tamaño ni de frecuencia: es una lectura simple por clave.
- Se agrega la ruta con ejemplos en `docs/curl-frontend-sprint-2.md`, en el mismo estilo que las demás.
- Los escenarios BDD se escriben en el idioma y estilo de los `.feature` existentes (Godog), siguiendo el Principio III de la constitución, y el desarrollo es Test-First (Principio I).
