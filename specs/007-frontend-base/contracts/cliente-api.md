# Contrato: cliente de la API

Código en `web/src/compartido/api/client.ts`. El cliente **no conoce endpoints** (FR-005): cada pantalla arma sus llamadas en su carpeta.

## Operaciones

| Operación | Qué hace |
|---|---|
| `get(ruta)` | GET; devuelve el JSON recibido |
| `post(ruta, cuerpo)` | POST con cuerpo JSON; devuelve el JSON recibido |
| `put(ruta, cuerpo)` | PUT con cuerpo JSON; devuelve el JSON recibido |

Todas envían `Accept: application/json` (necesario para que el proxy de desarrollo reenvíe la llamada al backend) y, con cuerpo, `Content-Type: application/json`.

## Resultado

- **Éxito (2xx)**: el cuerpo JSON ya leído; `undefined` si el cuerpo viene vacío.
- **Error de la API**: lanza `ApiError { status, message, field? }`.
- **Sin respuesta**: lanza `ConnectionError`.
- **Respuesta no interpretable** (cuerpo no JSON, éxito o error): lanza `ApiError` con el estado recibido y un mensaje genérico.

## Casos que cubren las pruebas (FR-010)

| Respuesta simulada | Resultado esperado |
|---|---|
| 200 con JSON | datos |
| 400 `{campo:"nombre", mensaje:"..."}` | `ApiError` estado 400, campo "nombre" |
| 400 `{mensaje:"..."}` | `ApiError` estado 400, sin campo |
| 404 `{mensaje:"..."}` | `ApiError` estado 404 |
| 409 `{mensaje:"..."}` | `ApiError` estado 409 |
| 500 con cuerpo no JSON | `ApiError` estado 500, mensaje genérico |
| `fetch` rechazado | `ConnectionError` |

## Endpoints existentes (referencia para las features)

`POST /backlog` · `POST /projects` · `GET/PUT /projects/{id}` · `GET /projects/{id}/status` · `POST/GET /projects/{id}/members` · `POST /sprints` · `POST /sprints/{id}/iniciar` · `POST /sprints/{id}/cerrar`. No existen `GET /projects` (#109), listado de Sprints (#110) ni listado del backlog (HU-15).
