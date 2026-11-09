# Contrato: alta de historia con criterios de aceptación

Endpoint existente: **`POST /backlog`** (`cmd/api/main.go`). El `proyecto_id` viaja en el cuerpo; no se agrega una ruta nueva.

## Request

`Content-Type: application/json`

| Campo | Tipo | Obligatorio | Reglas |
|---|---|---|---|
| `proyecto_id` | integer | Sí | Entero positivo; el proyecto debe existir. |
| `titulo` | string | Sí | No vacío después del trim; máximo 200 caracteres. |
| `descripcion` | string | Sí | No vacía después del trim; máximo 2000 caracteres. |
| `prioridad` | string | Sí | `M`, `S`, `C` o `W`. |
| `valor_negocio` | integer \| omitido | No | `1, 2, 3, 5, 8, 13, 21`. |
| `criterios_aceptacion` | string \| null \| omitido | No | Texto libre; máximo 2000 caracteres después del trim. Vacío o solo espacios = `NULL`. |

Ejemplo:

```json
{
  "proyecto_id": 1,
  "titulo": "Iniciar sesión",
  "descripcion": "Autenticación con usuario y contraseña",
  "prioridad": "M",
  "valor_negocio": 8,
  "criterios_aceptacion": "Dado un usuario registrado, cuando ingresa, entonces accede."
}
```

## Response 201

Devuelve la historia creada, con el mismo `historiaResponse` del listado y la edición:

```json
{
  "id": 12,
  "proyecto_id": 1,
  "titulo": "Iniciar sesión",
  "descripcion": "Autenticación con usuario y contraseña",
  "criterios_aceptacion": "Dado un usuario registrado, cuando ingresa, entonces accede.",
  "prioridad": "M",
  "estado": "Nueva",
  "valor_negocio": 8,
  "estimacion_sp": null,
  "sprint_id": null,
  "fecha_finalizacion": null
}
```

Si no se enviaron criterios, `criterios_aceptacion` es `null`.

## Errores

| Caso | HTTP | Cuerpo |
|---|---|---|
| Criterios de más de 2000 caracteres | 400 | `{"campo":"criterios_aceptacion","mensaje":"los criterios de aceptación no pueden superar los 2000 caracteres"}` |
| Título, descripción o prioridad inválidos | 400 | `{"campo":"…","mensaje":"…"}` (sin cambios) |
| Proyecto inexistente | 404 | `{"mensaje":"proyecto no encontrado"}` (sin cambios) |
| JSON inválido o `criterios_aceptacion` de tipo incorrecto | 400 | `{"mensaje":"JSON inválido"}` |
| Error interno al persistir | 500 | `{"mensaje":"error interno"}` (sin cambios) |

## Compatibilidad

El campo es opcional: los clientes que no lo envían siguen funcionando y reciben `criterios_aceptacion: null`, como hoy. No cambia el esquema, ni la ruta, ni los demás campos.
