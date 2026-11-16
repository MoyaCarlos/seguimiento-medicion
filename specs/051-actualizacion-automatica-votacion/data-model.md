# Data Model: Actualización automática de la votación (#225)

## Sin cambios de datos

Esta mejora es solo de frontend. No agrega tablas, columnas, campos, endpoints ni cambios en
el contrato de la votación.

- Se sigue usando `GET /backlog/{id}/votacion`, con el tipo `Votacion` ya existente
  (`historia_id`, `ronda`, `revelada_en`, `integrantes`, `promedio`, `diferencias`).
- No se toca `MODELO-DE-DATOS.md` ni las tablas de rondas y votos (HU-02/HU-14).
- El intervalo de sondeo, la visibilidad de la pestaña y el estado de la acción son estado de
  UI; no se persisten ni se envían al servidor.
- El botón «Actualizar» sigue haciendo las mismas consultas que antes.
