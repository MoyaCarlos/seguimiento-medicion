# Contrato: API de autenticación

Ninguna ruta existente cambia ni pasa a exigir sesión. Errores en el formato de siempre:
`{"campo": "...", "mensaje": "..."}` (el campo solo en validaciones).

El token viaja en `Authorization: Bearer <token>`. Las respuestas de sesión incluyen el
token **una sola vez**; la contraseña no aparece en ninguna respuesta.

## POST /auth/registro

Cuerpo: `{"nombre": "ana.perez", "contrasena": "una-clave-larga"}`

| Estado | Cuándo | Cuerpo |
|---|---|---|
| 201 | cuenta creada y sesión iniciada | `{"usuario": {"id": 7, "nombre": "ana.perez"}, "token": "<43 caracteres>", "vence_en": "2026-11-10T12:00:00Z"}` |
| 400 | JSON inválido | `{"mensaje": "JSON inválido"}` |
| 400 | nombre fuera de reglas | `{"campo": "nombre", "mensaje": "..."}` |
| 400 | contraseña de menos de 8 o más de 128 caracteres, o ausente | `{"campo": "contrasena", "mensaje": "..."}` |
| 409 | nombre en uso (sin distinguir mayúsculas; incluye usuarios previos) | `{"mensaje": "Ese nombre de usuario ya existe"}` |

El nombre se devuelve tal como se escribió. Orden de validación: `nombre`, luego `contrasena`.

## POST /auth/ingreso

Cuerpo: `{"nombre": "Ana.Perez", "contrasena": "una-clave-larga"}`

| Estado | Cuándo | Cuerpo |
|---|---|---|
| 200 | credenciales correctas | igual que el registro |
| 400 | JSON inválido | `{"mensaje": "JSON inválido"}` |
| 401 | usuario inexistente, usuario sin contraseña o contraseña incorrecta | `{"mensaje": "Usuario o contraseña incorrectos"}` (idéntico en los tres casos) |

El ingreso no valida el formato del nombre ni el largo de la contraseña: cualquier cuerpo
que no coincida da el mismo 401.

## GET /auth/yo

Cabecera: `Authorization: Bearer <token>`

| Estado | Cuándo | Cuerpo |
|---|---|---|
| 200 | sesión vigente | `{"id": 7, "nombre": "ana.perez"}` |
| 401 | sin cabecera, esquema distinto de `Bearer`, token mal formado, inexistente, vencido o cerrado | `{"mensaje": "Tu sesión no es válida"}` (idéntico en todos) |

## POST /auth/salir

Cabecera: `Authorization: Bearer <token>`. Sin cuerpo.

| Estado | Cuándo | Cuerpo |
|---|---|---|
| 204 | sesión cerrada; ese token deja de ser válido, las demás sesiones del usuario siguen | — |
| 401 | token no vigente | `{"mensaje": "Tu sesión no es válida"}`; sin efectos |

## Notas

- `GET /auth/yo` y `POST /auth/salir` con `Accept: text/html` caen en la SPA (comportamiento
  general de `ServirSPA`); el cliente siempre pide `application/json`.
- `/auth` se agrega a `PrefijosAPI` y al proxy de `web/vite.config.ts`. Las pantallas de la
  SPA son `/ingresar` y `/registro`, para no chocar con ese prefijo.
