# Contrato consumido: API de votación (spec 012)

Solo lectura de lo que implementa el backend (`internal/http/votacion_handler.go`, `votacion_dto.go`, `sprint_handler.go#escribirError`). Esta spec no lo modifica.

| Acción | Pedido | Éxito |
|--------|--------|-------|
| Consultar | `GET /backlog/{id}/votacion?usuario_id=N` (el parámetro es opcional) | 200 `Votacion` |
| Iniciar | `POST /backlog/{id}/votacion` `{usuario_id}` | 201 `Votacion` |
| Votar | `POST /backlog/{id}/votacion/votos` `{usuario_id, valor}` | 200 `Votacion` |
| Revelar | `POST /backlog/{id}/votacion/revelar` `{usuario_id}` | 200 `Votacion` |
| Integrantes | `GET /projects/{id}/members` | 200 `[{id, proyecto_id, nombre, rol, vinculado_en}]` |
| Historia (resumen) | `GET /projects/{id}/backlog` (solo historias sin Sprint; se busca la historia por `id`; reemplazable por `GET /backlog/{id}`, issue #222) | 200 `[{id, titulo, descripcion, prioridad, estimacion_sp, ...}]` |

`Votacion` = `{historia_id, ronda, revelada_en|null, integrantes:[{usuario_id, nombre, rol, ya_voto, valor|null}], promedio:{exacto, carta_sugerida}|null}`.

## Errores (todos con `{mensaje}`, y `campo` en 400 de validación)

| Estado | Cuándo | Qué hace la pantalla |
|--------|--------|----------------------|
| 400 | valor fuera de baraja, ronda ya revelada, falta votar, usuario inválido | muestra el mensaje |
| 403 | no es integrante / no es Scrum Master | muestra el mensaje |
| 404 | historia inexistente, o «la historia no tiene votación iniciada» | «sin votación» o mensaje (R1 de [research.md](../research.md)) |
| 409 | la historia ya tiene una votación | muestra el mensaje y vuelve a consultar |
| 502/503/504 o sin red | — | aviso de conexión (`ConnectionError`) |

Nota: el encargo listaba 400 y 404; el backend también responde 403 y 409, que se tratan igual (mostrar `mensaje`).
