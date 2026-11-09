# Contrato: `GET /backlog/{id}`

Consulta una historia por su identificador. Sin cuerpo, sin parámetros de consulta, sin sesión (FR-009).

## Petición

| Parte | Valor |
|-------|-------|
| Método y ruta | `GET /backlog/{id}` |
| `{id}` | entero positivo (`int64`); `007` equivale a `7` |

## Respuestas

### 200 — la historia existe

`Content-Type: application/json`

```json
{
  "id": 7,
  "proyecto_id": 1,
  "titulo": "Iniciar sesión",
  "descripcion": "Como usuario quiero iniciar sesión",
  "criterios_aceptacion": null,
  "prioridad": "Must",
  "estado": "Pendiente",
  "valor_negocio": 8,
  "estimacion_sp": 5,
  "sprint_id": null,
  "fecha_finalizacion": null
}
```

- Misma forma que `historiaResponse`, la de `GET /projects/{id}/backlog`, `PUT /backlog/{id}` y las demás.
- `sprint_id`: `null` o el id del Sprint.
- `fecha_finalizacion`: `null` o la fecha, con el formato del resto de la API.

### 404 — no existe

```json
{ "mensaje": "historia no encontrada" }
```

### 400 — identificador inválido

Para `abc`, `0`, `-3`, `1.5` o un número fuera del rango de enteros:

```json
{ "campo": "id", "mensaje": "identificador inválido" }
```

No consulta si la historia existe (FR-007).

## Compatibilidad

- No cambia ningún endpoint existente ni el texto de sus errores (`parsearID` queda intacto).
- No colisiona con `GET /backlog/{id}/horas`.
