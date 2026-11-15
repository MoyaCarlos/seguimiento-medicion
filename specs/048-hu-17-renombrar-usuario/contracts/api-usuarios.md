# Contrato API: renombrar el propio usuario (HU-17)

Todas las rutas exigen sesión (`Authorization: Bearer <token>`, HU-20).

## `PUT /auth/yo` — **nuevo**

Cambia el nombre del usuario de la sesión. Nivel: solo sesión. El cuerpo **no** lleva
`usuario_id`: el usuario sale de la sesión.

### Request

```json
{ "nombre": "ana.perez" }
```

`nombre` es obligatorio y se valida con las reglas de HU-20: 3 a 30 caracteres, solo letras
ASCII, números, punto, guion y guion bajo. Los espacios externos no son válidos.

### Responses

| Caso | HTTP | Cuerpo |
|---|---|---|
| Cambio aplicado | 200 | `{"id": 7, "nombre": "ana.perez"}` |
| Sin sesión / token inválido o vencido | 401 | `{"mensaje": "Tu sesión no es válida"}` |
| JSON inválido o `nombre` ausente | 400 | `{"mensaje": "JSON inválido"}` o `{"campo": "nombre", "mensaje": "..."}` |
| Nombre fuera de las reglas | 400 | `{"campo": "nombre", "mensaje": "el nombre de usuario debe tener entre 3 y 30 caracteres"}` o `{"campo": "nombre", "mensaje": "el nombre de usuario solo admite letras sin tildes, números, punto, guion y guion bajo"}` |
| Nombre normalizado ya usado | 409 | `{"mensaje": "Ese nombre de usuario ya existe"}` |
| Usuario de la sesión inexistente (no alcanzable desde la API) | 404 | `{"mensaje": "registro no encontrado"}` |

### Efectos

- `usuarios.nombre` y `usuarios.nombre_normalizado` se actualizan juntos.
- Las sesiones, votos, horas, integrantes y defectos siguen asociados al mismo `usuario_id`.
- La sesión actual sigue vigente; `GET /auth/yo` devuelve el nombre nuevo.

## `GET /auth/yo` — sin cambios

Devuelve `{"id": <id>, "nombre": "<nombre actual>"}` para la sesión. Después del renombrado
refleja el nombre nuevo.
