# Contrato API: Gestión del equipo (HU-16)

Todas las rutas exigen sesión (`Authorization: Bearer <token>`, HU-20).

## `DELETE /projects/{id}/members/{usuario_id}` — **nuevo**

Quita (baja lógica) a un integrante. Nivel: sesión y Scrum Master (HU-21).

| Caso | HTTP | Cuerpo |
|---|---|---|
| Baja aplicada | 204 | — |
| Sin sesión | 401 | `{"mensaje": "<no autenticado, HU-20>"}` |
| `id` o `usuario_id` no son enteros positivos | 400 | `{"campo": "id", "mensaje": "el identificador debe ser un entero positivo"}` |
| Proyecto inexistente | 404 | `{"mensaje": "proyecto no encontrado"}` |
| Quien actúa no integra el proyecto | 403 | `{"mensaje": "No sos integrante de este proyecto"}` |
| Quien actúa es Product Builder | 403 | `{"mensaje": "Solo el Scrum Master puede realizar esta acción"}` |
| Proyecto cerrado (#300) | 409 | `{"mensaje": "El proyecto está cerrado"}` |
| Usuario que no es integrante activo (no existe, nunca integró, ya dado de baja) | 404 | `{"mensaje": "integrante no encontrado"}` |
| Usuario es el Scrum Master | 409 | `{"mensaje": "No se puede quitar al Scrum Master del proyecto"}` |

Efectos: `integrantes.fecha_baja = hoy`; se borran sus votos de rondas no reveladas del
proyecto. Votos revelados, horas y defectos no cambian.

## `POST /projects/{id}/members` — cambia

Mismo contrato que HU-21. Cambio: si el usuario fue dado de baja del proyecto, se lo
**reactiva** como Product Builder (misma fila, conserva `creado_en`) y responde **201** con
el integrante, igual que un alta.

| Caso nuevo / aclarado | HTTP |
|---|---|
| Usuario dado de baja, sin rol o `product_builder` | 201 (reactivado) |
| Usuario dado de baja, `rol: scrum_master` | 409 «el proyecto ya tiene un Scrum Master» (#236) |
| Integrante activo (incluido el Scrum Master, con cualquier rol) | 400 `{"campo": "integrante", "mensaje": "el integrante ya pertenece al proyecto"}` |

## Lecturas — cambian sin cambiar su forma

| Ruta | Cambio |
|---|---|
| `GET /projects/{id}/members` | solo integrantes activos |
| `GET /projects` («Mis proyectos») | no incluye proyectos donde la persona fue dada de baja |
| `GET /backlog/{id}/votacion` | integrantes activos + quien tenga voto en la ronda actual |
| `GET /backlog/{id}/votacion/rondas` (historial de rondas reveladas) | conserva nombre y rol de quien fue dado de baja |
| Reporte PDF del proyecto | el equipo lista solo activos; horas y métricas sin cambios |
