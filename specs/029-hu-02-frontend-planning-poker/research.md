# Research: Pantalla de Planning Poker (spec 029)

## R1 — Distinguir «sin votación» de «historia inexistente»

**Decision**: un 404 de `GET /backlog/{id}/votacion` con mensaje exacto «la historia no tiene votación iniciada» se interpreta como «sin votación» (se ofrece Iniciar al Scrum Master); cualquier otro error se muestra como mensaje.
**Rationale**: la API devuelve 404 `{mensaje}` para ambos casos y no hay un `GET /backlog/{id}`. Comparar el mensaje es lo único disponible sin tocar el backend.
**Alternatives**: tratar todo 404 como «sin votación» (oculta «historia no encontrada»); agregar un código al backend (fuera de alcance).
**Riesgo**: si el texto cambia en el backend, la pantalla mostraría el error en vez de «Iniciar». Lo cubre una prueba con el texto real.

## R2 — Quién consulta

**Decision**: la consulta envía `usuario_id` del integrante elegido; sin elección se consulta sin `usuario_id` (todo oculto). Al cambiar de integrante se vuelve a consultar.
**Rationale**: FR-007 de la spec 012; la API acepta la consulta sin `usuario_id`.

## R3 — Estado tras cada acción

**Decision**: iniciar, votar y revelar responden con la vista de quien actúa; la pantalla la usa tal cual como nuevo estado. Tras un error se vuelve a consultar para reflejar el estado real.
**Rationale**: evita un GET extra en el caso feliz.

## R4 — Respuestas atrasadas

**Decision**: cada consulta lleva una marca de vigencia (patrón `vigente` ya usado en `backlog`); se descartan respuestas de otro integrante o historia.

## R5 — Doble envío

**Decision**: un único estado `enviando` bloquea las tres acciones y las cartas mientras hay un pedido en curso (mismo patrón que `EstadoHistoria`).

## R6 — Sin polling

**Decision**: botón «Actualizar» (clarify). Descartado el sondeo periódico por YAGNI.

## R7 — Aislamiento

**Decision**: el backlog enlaza con `Link to={`/projects/${id}/votacion/${historia.id}`}`; la ruta se registra en `rutas.tsx` con `lazy`. Ninguna feature importa a otra.
