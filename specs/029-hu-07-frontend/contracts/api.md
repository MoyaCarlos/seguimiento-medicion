# Contrato consumido: Imputación de horas en pantalla (HU-07 Frontend)

Esta historia no define ni cambia endpoints. La fuente de verdad es el [contrato de la spec 019](../../019-hu-07-imputacion-horas/contracts/openapi.yaml).

| Llamada | Cuándo | Uso |
|---|---|---|
| `GET /backlog/{id}/horas` | Al mostrar una historia y después de guardar una carga | `total_horas` y `registros` |
| `POST /backlog/{id}/horas` | Al guardar | Cuerpo `{ integrante_id, fecha, actividad, horas }`, con `horas` como string. Con `201` no se usa el cuerpo: se reconsulta el total |
| `GET /projects/{id}/members` | Al abrir el formulario | `id` y `nombre` para la lista |

## Errores del `POST`

| `campo` de la respuesta | Dónde se muestra |
|---|---|
| `integrante_id` | Junto a la lista de integrantes |
| `fecha` | Junto a la fecha |
| `actividad` | Junto a la actividad |
| `horas` | Junto a las horas |
| `historia_id`, otro valor o sin `campo` (incluye 404 y 409) | Aviso general del formulario |
| Sin conexión o 502/503/504 | Aviso general ("No se pudo conectar con el servidor") |
