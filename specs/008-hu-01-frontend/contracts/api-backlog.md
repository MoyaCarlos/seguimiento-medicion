# Contrato: creación de historias desde el frontend (HU-01)

La pantalla consume un único endpoint existente del backend. No expone endpoints
nuevos. Fuente de verdad del backend:
[`internal/http/backlog_dto.go`](../../../internal/http/backlog_dto.go) y
[`internal/http/backlog_handler.go`](../../../internal/http/backlog_handler.go).

## `POST /backlog`

### Cuerpo enviado por la pantalla

```json
{
  "proyecto_id": 1,
  "titulo": "Como usuario quiero iniciar sesión",
  "descripcion": "Autenticación con email y contraseña",
  "prioridad": "M",
  "valor_negocio": 13
}
```

| Campo | Tipo | Obligatorio | Regla de la pantalla |
|---|---|---|---|
| `proyecto_id` | número entero | Sí | Entero positivo tomado de la URL; si no es válido no se envía |
| `titulo` | string | Sí | Recortado; 1–200 puntos de código |
| `descripcion` | string | Sí | Recortada; 1–2000 puntos de código |
| `prioridad` | `"M" \| "S" \| "C" \| "W"` | Sí | Código de la opción MoSCoW elegida; por defecto `"M"` |
| `valor_negocio` | número entero | No | `1, 2, 3, 5, 8, 13` o `21`; **se omite del cuerpo** si la opción es "Sin valor" |

La pantalla hace **una sola** llamada por guardado y la construye con `post()` de
`web/src/compartido/api/client.ts` (agrega `Accept` y `Content-Type:
application/json`).

### Respuesta exitosa (`201 Created`)

```json
{
  "id": 7,
  "proyecto_id": 1,
  "titulo": "Como usuario quiero iniciar sesión",
  "descripcion": "Autenticación con email y contraseña",
  "prioridad": "M",
  "estado": "Nueva",
  "valor_negocio": 13,
  "estimacion_sp": null
}
```

La confirmación muestra: título, prioridad legible (`M` → "Must have", etc.), estado
"Nueva", ausencia de estimación (`estimacion_sp: null`) y el valor de negocio solo si no
es `null`.

### Errores

Formato general: `{ "campo"?: string, "mensaje": string }`.

| Estado | Cuerpo típico | Dónde se muestra |
|---|---|---|
| `400` con `campo` conocido (`titulo`, `descripcion`, `prioridad`, `valor_negocio`) | `{ "campo": "titulo", "mensaje": "el título es obligatorio" }` | Junto a ese campo |
| `400` sin `campo` o con `campo` desconocido (incluido `proyecto_id`) | `{ "mensaje": "..." }` | Aviso general arriba del formulario |
| `404` proyecto inexistente | `{ "mensaje": "proyecto no encontrado" }` | Aviso general |
| `500` | `{ "mensaje": "error interno" }` | Aviso general |
| Sin conexión | `ConnectionError` del cliente | Aviso general |
| Cuerpo de error no interpretable | `ApiError` genérica del cliente | Aviso general |

Los mensajes de error no se borran al escribir y se reemplazan recién en el próximo
intento de guardado.

## Contrato de la pantalla (frontend)

- **Ruta**: `/projects/:id/backlog` (ya registrada en `web/src/rutas.tsx`; no cambia).
- **Módulo**: `web/src/features/backlog/index.tsx` exporta `Component`, sin props, como
  espera el registro de rutas.
- **Entradas**: solo `:id` de la URL.
- **Salidas observables**: creación en la API, confirmación de éxito, advertencias locales
  y errores del servidor según la tabla anterior.
- **No expone** endpoints, no lee estado de otra pantalla y no modifica `compartido/`.
