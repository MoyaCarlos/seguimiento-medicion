# Contrato de UI: polling de la votación (#225)

## Disparo

- Intervalo: **5 segundos** (`INTERVALO_ACTUALIZACION_MS`).
- Se activa solo cuando:
  - hay una votación cargada y `revelada_en === null`;
  - `document.visibilityState === 'visible'`;
  - no hay una acción en curso (`enviando === false`).
- Se detiene al revelarse la ronda, al ocultarse la pestaña o al empezar una acción.

## Consulta

- Usa `consultarVotacion(historiaId)`, el mismo cliente de la spec 030.
- No modifica `recarga` ni `cargadaPara`: no muestra `LoadingIndicator` ni deshabilita la
  pantalla.
- Actualiza `votacion`, limpia `sinVotacion` y el error de votación al responder bien.

## Carreras

- Cada consulta automática captura `accionesIniciadas.current`.
- Si una acción incrementó el contador mientras la consulta estaba en curso, la respuesta se
  descarta.
- Se evita solapar dos consultas automáticas con una bandera local de la consulta en curso.

## Errores

- Se muestra el mensaje de consulta existente (`mensajeDeError`).
- La última votación conocida se conserva.
- El error de acción no se borra ni se mezcla con el de la consulta.

## Lo que no cambia

- El botón «Actualizar» y su comportamiento.
- La API, el esquema y el tipo `Votacion`.
- Las consultas de integrantes e historia: el sondeo solo consulta la votación.
