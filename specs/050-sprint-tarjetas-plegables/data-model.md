# Data Model: Tarjetas plegables y confirmación de estado (#219)

## Sin cambios de datos

Esta mejora es solo de presentación e interacción en el frontend. No agrega tablas, columnas,
campos, endpoints ni cambios en el contrato de las historias.

- `historias.estado` y `historias.fecha_finalizacion` siguen con las reglas de HU-11
  (`MODELO-DE-DATOS.md` §4.4).
- El `PUT /backlog/{id}/estado` y sus mensajes no cambian.
- La tarjeta reutiliza los campos que ya devuelve `GET /sprints/{id}/backlog`: `titulo`,
  `prioridad`, `estado`, `estimacion_sp` y `fecha_finalizacion`.
- El estado abierto/cerrado de `<details>` y la confirmación son estado de UI, no se
  persisten ni se envían a la API.

No se modifica `MODELO-DE-DATOS.md`.
