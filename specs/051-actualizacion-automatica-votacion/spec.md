# Feature Specification: Actualización automática de la votación (#225)

**Feature Branch**: `feature/225-actualizacion-automatica-votacion`

**Issue**: #225

**Created**: 2026-11-15

**Status**: Implementada (pendiente de PR)

**Relacionado**: spec 030 (pantalla de Planning Poker), spec 012 (API de votación), #222,
#218. Solo frontend; no cambia el contrato de la API.

**Input**: User description: "Planning Poker: actualización automática de la votación (sin
pulsar «Actualizar»). Hoy, para ver quién ya votó o si el Scrum Master reveló, hay que pulsar
«Actualizar» (decisión del clarify de la spec 030: sin polling). En una reunión de estimación
se siente lento. Consultar la votación cada pocos segundos mientras la ronda esté abierta y
pausar la consulta cuando la pestaña no está visible o la ronda está revelada. Polling simple,
sin SSE ni cambios de contrato."

## Objetivo

Que la pantalla de Planning Poker refleje sola los votos de los demás y la revelación del
Scrum Master mientras la ronda está abierta, sin que nadie tenga que pulsar «Actualizar». El
sondeo se detiene cuando no aporta (pestaña oculta, ronda revelada o acción en curso) para no
gastar pedidos ni mantener despierto al backend de Render innecesariamente.

## Entradas y salidas

- **Entrada**: la votación ya cargada por la pantalla (`GET /backlog/{id}/votacion`), el estado
  de la ronda (abierta o revelada), la visibilidad de la pestaña y el estado de envío de una
  acción.
- **Salida**: la votación se vuelve a consultar sola cada 5 segundos mientras la ronda está
  abierta y la pestaña visible; los votos ajenos y la revelación aparecen sin pulsar nada. Si
  la pestaña está oculta, la ronda está revelada o hay una acción en curso, no se consulta.
- **Sin cambios**: la API, el esquema, las rutas, el contrato de `Votacion` y el botón
  «Actualizar» (sigue disponible como refresco manual).
- **Sin indicador de carga**: las consultas automáticas no muestran el indicador de carga ni
  deshabilitan la pantalla.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Ver los votos de los demás sin pulsar «Actualizar» (Priority: P1)

Como integrante del equipo en una reunión de estimación, mientras la ronda está abierta veo
que otro integrante ya votó sin tener que pulsar «Actualizar».

**Independent Test**: con una ronda abierta y la pestaña visible, dejar pasar el intervalo de
sondeo y comprobar que se vuelve a consultar la votación y que la fila del otro integrante pasa
de «Falta» a «Ya votó» / «Oculto».

**Acceptance Scenarios**:

1. **Given** una ronda abierta y la pestaña visible, **When** pasan 5 segundos, **Then** la
   pantalla vuelve a consultar la votación sin pulsar «Actualizar».
2. **Given** que otro integrante votó después de la última consulta, **When** llega la
   consulta automática, **Then** su fila muestra «Ya votó» y «Oculto» mientras la ronda no
   está revelada.
3. **Given** una consulta automática en curso, **When** todavía no respondió, **Then** la
   pantalla sigue operable y no muestra el indicador de carga.

---

### User Story 2 - Detener el sondeo cuando no hace falta (Priority: P1)

Como equipo, el sondeo se detiene cuando la pestaña no está visible, cuando la ronda ya se
reveló o mientras hay una acción en curso, para no gastar pedidos ni mantener despierto al
servicio.

**Independent Test**: con la pestaña oculta, avanzar el tiempo y comprobar que no hay
consultas; volver a visible y comprobar que se reanudan. Repetir con la ronda revelada y con
una acción pendiente.

**Acceptance Scenarios**:

1. **Given** una ronda abierta y la pestaña oculta, **When** pasan varios intervalos,
   **Then** no se consulta la votación.
2. **Given** la pestaña oculta con una ronda abierta, **When** la pestaña vuelve a estar
   visible y pasa un intervalo, **Then** se reanuda la consulta automática.
3. **Given** una ronda revelada, **When** pasa el intervalo, **Then** no se vuelve a consultar
   sola.
4. **Given** una acción de votar, revelar, iniciar, abrir ronda o registrar estimación en
   curso, **When** pasa el intervalo, **Then** no se consulta sola hasta que la acción termina.

### Edge Cases

Redactados en formato EARS.

- **WHEN** la ronda está abierta y la pestaña visible, **THEN** el sistema **shall** consultar
  la votación cada 5 segundos.
- **IF** la pestaña está oculta, **THEN** el sistema **shall not** consultar la votación.
- **IF** la ronda está revelada, **THEN** el sistema **shall not** consultar la votación
  automáticamente.
- **IF** hay una acción en curso, **THEN** el sistema **shall not** iniciar consultas
  automáticas.
- **WHEN** una consulta automática responde, **THEN** el sistema **shall** actualizar la
  votación sin mostrar el indicador de carga ni deshabilitar la pantalla.
- **IF** una consulta automática falla, **THEN** el sistema **shall** mostrar el mensaje de
  error de consulta existente y **shall** conservar la última votación conocida.
- **IF** una acción empieza mientras hay una consulta automática en curso, **THEN** la
  respuesta de la consulta **shall** descartarse para no pisar el resultado de la acción.
- **IF** la consulta automática detecta que la ronda se reveló, **THEN** el sistema **shall**
  mostrar el resultado y **shall** detener el sondeo.

## Requirements *(mandatory)*

### Reglas de comportamiento (EARS)

**Ubicuas**

- **FR-001**: El sistema MUST consultar la votación automáticamente cada 5 segundos mientras
  la ronda esté abierta y la pestaña visible.
- **FR-002**: El sistema MUST NOT consultar automáticamente con la pestaña oculta, con la
  ronda revelada o mientras una acción está en curso.
- **FR-003**: La consulta automática MUST actualizar el estado de la votación usando el mismo
  contrato de la API, sin cambiar el backend.
- **FR-004**: La consulta automática MUST NOT mostrar el indicador de carga ni deshabilitar la
  pantalla.
- **FR-005**: El botón «Actualizar» MUST seguir funcionando como refresco manual.
- **FR-006**: La funcionalidad MUST ser solo de frontend, sin dependencias nuevas, sin SSE y
  sin cambios de esquema.

**Basadas en evento**

- **WHEN** pasan 5 segundos con la ronda abierta y la pestaña visible, el sistema SHALL
  consultar la votación.
- **WHEN** la consulta automática responde, el sistema SHALL reemplazar la votación mostrada.
- **WHEN** la pestaña vuelve a estar visible con la ronda abierta, el sistema SHALL reanudar
  las consultas automáticas.
- **WHEN** la ronda se revela, la acción en curso termina o la pestaña vuelve a estar visible,
  el sistema SHALL recalcular si corresponde seguir sondeando.

**No deseadas**

- **IF** la consulta automática falla, THEN el sistema SHALL mostrar el mensaje de error de
  consulta y SHALL conservar la última votación conocida.
- **IF** una acción empieza durante una consulta automática, THEN la respuesta de esa consulta
  SHALL descartarse.

## Fuera de alcance

- Server-Sent Events, WebSockets o cualquier cambio de contrato o de backend.
- Sondeo de la lista de integrantes o de la historia: solo se reconsulta la votación.
- Polling cuando no hay votación iniciada o cuando la pantalla no está montada.
- Cambiar la decisión de producto de la spec 030 sobre el botón «Actualizar»: se conserva.

## Success Criteria *(mandatory)*

- **SC-001**: En una ronda abierta y visible, los votos ajenos y la revelación aparecen sin
  pulsar «Actualizar» en un máximo de 5 segundos más el tiempo de respuesta.
- **SC-002**: Con la pestaña oculta, la ronda revelada o una acción en curso, no se hace
  ninguna consulta automática.
- **SC-003**: La consulta automática no muestra el indicador de carga ni bloquea la pantalla.
- **SC-004**: Las pruebas de componente del hook y de la pantalla pasan en Vitest, junto con
  el resto de la suite, el lint y el build.
- **SC-005**: No se agregan dependencias, endpoints ni columnas.

## Assumptions

- La spec 030 ya dejó la consulta en `useVotacion` y el botón «Actualizar»; esta mejora agrega
  un efecto de sondeo en ese hook.
- El polling simple alcanza para el TP: no hay SSE ni WebSockets (KISS/YAGNI).
- El intervalo de 5 segundos es un equilibrio entre la sensación de tiempo real y el costo de
  Render Free; al pausar con la pestaña oculta, el servicio puede dormir si nadie lo usa.
- Los errores de la consulta automática se muestran con el mensaje de consulta existente y no
  pisan los errores de acciones.
- Solo se tocan archivos del frontend (`web/src/features/votacion/`) y la documentación.
