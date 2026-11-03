# Research: Imputación de horas en pantalla (HU-07 Frontend)

## R1. Dónde vive el código

- **Decision**: dentro de `features/sprints`, con `listarIntegrantes` declarada en `sprints/api.ts`.
- **Rationale**: el formulario se muestra en la tarjeta del Sprint y `aislamiento.test.ts` prohíbe que una feature importe a otra. Es una sola línea sobre el cliente compartido.
- **Alternatives considered**: crear `features/horas` (Sprints tendría que importarla, y el test lo prohíbe); mover la función a `compartido/` (mueve código de `proyecto` sin beneficio real por una línea).

## R2. Cómo se obtiene el total de cada historia

- **Decision**: cada historia consulta `GET /backlog/{id}/horas` al montarse. Después de guardar, se vuelve a consultar.
- **Rationale**: el total y el orden del detalle los decide el backend, sin sumar ni ordenar en el cliente (DRY, FR-013 de la spec 019). La respuesta trae total y detalle juntos, así "Ver detalle" no necesita otra consulta.
- **Alternatives considered**: sumar localmente la respuesta del `POST` (duplica la suma y el orden); agregar un endpoint por Sprint (cambia el backend, fuera de alcance).
- **Consecuencia**: mientras se reconsulta se conserva el total anterior, sin indicador de carga, para que no parpadee. Si la reconsulta falla, se aplica el mismo aviso que en la primera consulta (US3.4).

## R3. Formato en que se envían las horas

- **Decision**: siempre como string, después de recortar espacios y reemplazar la coma por punto.
- **Rationale**: el backend acepta número o string (spec 019, R3), pero su formato solo admite punto: `1,5` responde "las horas deben ser un número". Un string conserva el texto no numérico, así el backend responde su error propio (US2.2).
- **Alternatives considered**: `<input type="number">` impide escribir texto no numérico y el navegador vacía el valor, de modo que el error sería "obligatorias" en lugar de "deben ser un número", y el escenario US2.2 no se podría ejercitar. Por eso se usa `type="text"` con `inputMode="decimal"`.

## R4. Integrante sin elegir

- **Decision**: se frena en la pantalla, con el mensaje "Elegí un integrante", sin enviar la carga.
- **Rationale**: sin elección, `integrante_id` sería 0 y el backend respondería "el integrante no pertenece al proyecto de la historia", un mensaje engañoso.
- **Alternatives considered**: enviar y mostrar el error del backend (mensaje confuso); `required` nativo del `<select>` (mensaje del navegador, distinto del resto de los errores).

## R5. Fecha propuesta

- **Decision**: la menor entre hoy y la fecha de fin del Sprint, calculando "hoy" con la fecha **local** del navegador.
- **Rationale**: `toISOString()` da la fecha UTC: en Argentina (UTC-3), a partir de las 21:00 propondría mañana y el backend rechazaría la carga por fecha futura. Si el Sprint sigue Activo pasada su fecha de fin, proponer hoy también daría error siempre.
- **Alternatives considered**: atributos `min`/`max` del selector, descartados porque repiten la regla del backend y el navegador la frena con su propio mensaje.

## R6. Formato de horas mostrado

- **Decision**: usar `horas` y `total_horas` de la API con `toLocaleString('es-AR', { maximumFractionDigits: 2 })`, seguido de " h".
- **Rationale**: la API ya hace la conversión "solo para mostrar" (contrato de la spec 019). Da `5,5 h`, `4 h` y `1,02 h` para 61 minutos.
- **Alternatives considered**: mostrar "1 h 1 min" (la spec 019 y el issue hablan de horas).

## R7. Estado del formulario

- **Decision**: el formulario es un componente que se monta al abrirlo y se desmonta al cerrarlo. Consulta integrantes al montarse. Al guardar con éxito queda abierto, vacía actividad y horas, y conserva integrante y fecha.
- **Rationale**: desmontar garantiza "sin nadie seleccionado cada vez que se abre" sin lógica de reinicio, y evita consultar integrantes para todas las historias del Sprint.
- **Alternatives considered**: consultar integrantes una vez en la pantalla (estado compartido a través de varios componentes, sin necesidad).

## R8. Pruebas existentes que simulan `fetch`

- **Decision**: ajustar `SprintBacklog.test.tsx` para simular `fetch` según la URL. `index.test.tsx` ya lo hace con `simularApi`: solo se le suma la respuesta del esfuerzo. `rutas.test.tsx` se revisa al correr toda la suite.
- **Rationale**: `SprintBacklog.test.tsx` usa `mockResolvedValue` con **una sola** `Response`, cuyo cuerpo solo se lee una vez. Al sumar una consulta por historia, esa segunda lectura fallaría y las pruebas dejarían de reflejar el comportamiento real.
- **Alternatives considered**: simular `EsfuerzoHistoria` en esas pruebas (oculta la integración con la tarjeta).
