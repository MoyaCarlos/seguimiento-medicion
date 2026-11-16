# Research: Actualización automática de la votación (#225)

## R1. Polling simple en lugar de SSE

- **Decisión**: consultar `GET /backlog/{id}/votacion` cada 5 segundos desde el hook.
- **Rationale**: el contrato ya existe y la pantalla ya sabe interpretarlo. SSE exigiría un
  endpoint nuevo, manejo de conexiones y cambios en el backend; para una pantalla de pocos
  integrantes el polling alcanza y cumple KISS/YAGNI.
- **Alternativas**: SSE o WebSockets (más costo y complejidad en Render Free); no hacer nada
  (el problema de usabilidad sigue).

## R2. Consulta silenciosa

- **Decisión**: la consulta automática no incrementa `recarga` ni `cargadaPara`; actualiza
  `votacion` directamente y no toca `cargando`.
- **Rationale**: si usara el mismo camino que «Actualizar», el indicador de carga aparecería
  cada 5 segundos y la pantalla se deshabilitaría sin necesidad.
- **Alternativas**: reutilizar `recargar()` (parpadeo y bloqueo); ocultar el indicador solo
  para el sondeo (más estado y más casos).

## R3. Cuándo sondear

- **Decisión**: solo con `votacion` cargada y sin revelar, `document.visibilityState ===
  'visible'` y sin acción en curso (`enviando`).
- **Rationale**: son los tres casos en que una consulta puede aportar algo. Con la pestaña
  oculta no se gasta red ni se mantiene despierto al servicio; con la ronda revelada ya no hay
  cambios que traer; con una acción en curso se evita pisar su resultado.
- **Alternativas**: sondear siempre (costo innecesario); sondear solo al revelar (no resuelve
  el problema principal).

## R4. Carreras con acciones

- **Decisión**: cada consulta automática captura `accionesIniciadas.current` y descarta su
  respuesta si una acción incrementó el contador mientras estaba en curso.
- **Rationale**: reutiliza el mecanismo que ya protege las consultas manuales y las acciones;
  una respuesta vieja no puede pisar un voto, una revelación o una nueva ronda.
- **Alternativas**: abortar el `fetch` (el cliente actual no expone `AbortController`); un
  contador nuevo (duplicaría el existente).

## R5. Errores del sondeo

- **Decisión**: mostrar el mensaje de error de consulta existente y conservar la última
  votación conocida.
- **Rationale**: la pantalla no se rompe ni pierde el contexto de la reunión; el usuario puede
  seguir viendo los votos ya cargados y usar «Actualizar».
- **Alternativas**: silenciar el error (el usuario no sabe por qué no se actualiza); recargar
  toda la pantalla (pierde el contexto y el foco).

## R6. Pruebas

- **Decisión**: `renderHook` con temporizadores falsos, `document.visibilityState` simulado y
  mocks de `./api`; escenarios Given-When-Then en el nombre.
- **Rationale**: el sondeo es comportamiento del hook y se puede verificar de forma
  determinista sin depender del reloj real ni del navegador. Las pruebas de la pantalla siguen
  cubriendo el renderizado y «Actualizar».
- **Alternativas**: esperas reales de 5 segundos (lentas y frágiles); pruebas E2E (fuera del
  stack actual).
