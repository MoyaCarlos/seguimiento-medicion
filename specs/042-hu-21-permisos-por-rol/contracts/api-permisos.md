# Contrato de API: Permisos por rol (HU-21)

Todas las rutas protegidas reciben `Authorization: Bearer <token>` (HU-20). Los errores
mantienen la forma actual `{"campo"?: string, "mensaje": string}`.

## Respuestas comunes

| Situación | Código | Cuerpo |
|---|---|---|
| Sin sesión, token inválido o vencido | 401 | `{"mensaje": "Tu sesión no es válida"}` |
| Recurso de la ruta inexistente | 404 | mensaje actual de cada recurso |
| No integra el proyecto | 403 | `{"mensaje": "No sos integrante de este proyecto"}` |
| No es Scrum Master (en acción de gestión) | 403 | `{"mensaje": "Solo el Scrum Master puede realizar esta acción"}` |

Orden: 401 → 404 → 403 (no integrante) → 403 (no Scrum Master) → 400/409. Excepción: un
cuerpo que no es JSON válido responde 400 «JSON inválido» antes del 403.

## Nivel de cada ruta

La tabla completa está en la spec («Qué se protege») y es la que recorre el test de rutas.
Resumen de cambios de contrato por ruta:

### `POST /projects` — solo sesión

- **Quita** `creador` del cuerpo (si llega, se ignora).
- Quien tiene la sesión queda como Scrum Master.
- Respuesta 201: igual que hoy más `"rol": "scrum_master"`.

### `GET /projects` — «Mis proyectos», solo sesión, nunca 403

```json
[
  { "id": 3, "nombre": "Tienda", "descripcion": "", "fecha_inicio": null, "fecha_fin": null,
    "creado_en": "2026-11-10T12:00:00Z", "horas_por_story_point": null,
    "rol": "product_builder" }
]
```

Vacía (`[]`) si la persona no integra ningún proyecto.

### `GET /projects/{id}` — sesión e integrante

Igual que hoy más `"rol"` de quien consulta. 403 si no integra.

### `GET /users` — nuevo, solo sesión

```json
[ { "id": 1, "nombre": "ana.perez" }, { "id": 2, "nombre": "Bruno" } ]
```

Solo usuarios con cuenta (contraseña). Nunca contraseña, huella ni sesiones.

### `POST /projects/{id}/members` — sesión y Scrum Master

Cuerpo: `{"nombre": "Bruno", "rol"?: "product_builder"}`.

| Caso | Código |
|---|---|
| Usuario registrado, sin `rol` o con `product_builder` | 201, `rol: "product_builder"` |
| `rol: "scrum_master"` | 409 «el proyecto ya tiene un Scrum Master» (#236) |
| Usuario inexistente o sin contraseña | 404 `{"mensaje": "usuario no encontrado"}` |
| Ya integra el proyecto | 400 (regla de duplicado actual) |

Ya **no** crea usuarios.

### Votación (`/backlog/{id}/votacion...`)

- `POST /backlog/{id}/votacion`, `.../revelar`, `.../rondas`, `.../estimacion`: sesión y
  Scrum Master. **Quita** `usuario_id` del cuerpo.
- `POST /backlog/{id}/votacion/votos`: sesión e integrante; cuerpo `{"valor": 5}` (sin
  `usuario_id`). El voto queda a nombre de la sesión.
- `GET /backlog/{id}/votacion`: sesión e integrante. **Quita** la query `?usuario_id=N`: el
  «mi voto» de la respuesta es siempre el de la sesión.
- `GET /backlog/{id}/votacion/rondas`: solo sesión (rondas reveladas).

### `POST /backlog/{id}/horas` — sesión e integrante

**Quita** `integrante_id`: las horas quedan a nombre de la sesión.

### Sprints, backlog y defectos

Sin cambios de cuerpo; solo el nivel de protección de la tabla de la spec.
