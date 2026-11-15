# Research: Tarjetas plegables y confirmación de estado (#219)

## R1. Plegado con elementos nativos

- **Decisión**: una tarjeta por historia con `<details>` y `<summary>`; el resumen muestra el
  título, una marca con el estado y los Story Points. El detalle conserva prioridad, fecha de
  finalización, control de estado y esfuerzo.
- **Rationale**: `<summary>` es enfocable y se activa con Enter o Espacio; el navegador maneja
  el estado abierto/cerrado sin código ni dependencias.
- **Alternativas**: un `button` con `aria-expanded` y estado controlado (más código y otra
  fuente de verdad); una librería de acordeones (dependencia innecesaria).

## R2. Marca de estado y Story Points

- **Decisión**: una etiqueta corta con el estado (`data-estado` para el color) y el texto
  «N SP» o «Sin estimar» en el resumen.
- **Rationale**: permite leer el avance de un vistazo sin abrir la tarjeta y no repite la
  prioridad, que queda en el detalle.
- **Alternativas**: mostrar todos los datos en el resumen (vuelve al problema de altura);
  mostrar solo el título (pierde la lectura rápida que pide el issue).

## R3. Confirmación en línea

- **Decisión**: al pulsar «Iniciar» o «Completar» se reemplaza la acción por una pregunta con
  «Confirmar» y «Cancelar». El foco inicial queda en «Cancelar» y Escape cierra sin enviar.
  Solo «Confirmar» ejecuta el `PUT` de HU-11.
- **Rationale**: reutiliza el patrón ya probado de `ZonaCierreProyecto`, es accesible y se
  puede probar como componente. Un clic fuera no envía nada porque solo «Confirmar» dispara
  la petición.
- **Alternativas**: `window.confirm` (no se puede probar como componente ni estilizar);
  deshacer después del cambio (la API de HU-11 no admite transiciones inversas).

## R4. Doble envío y errores

- **Decisión**: mientras el pedido está en curso, «Confirmar» y «Cancelar» quedan
  deshabilitados y la función de envío vuelve a comprobar el estado. Ante un error se cierra
  la confirmación, se muestra el mensaje de la API y la acción vuelve a estar disponible.
- **Rationale**: mantiene la protección de HU-11 y la posibilidad de reintentar sin recargar.
- **Alternativas**: dejar la confirmación abierta tras el error (el mensaje y el reintento
  quedan más lejos); recargar la lista (pierde el contexto y el foco).

## R5. Pruebas

- **Decisión**: pruebas de componente con el escenario Given-When-Then en el nombre. Se
  comprueba `details.open` para el plegado, el texto del resumen, la apertura con clic, la
  confirmación, «Cancelar», Escape, el clic fuera y el doble envío.
- **Rationale**: el Principio III se cumple por equivalencia en el frontend. `details.open` es
  observable y estable en jsdom; no se depende de estilos del navegador.
- **Alternativas**: pruebas de extremo a extremo (fuera del stack actual del frontend).
