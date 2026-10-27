# Contrato: servido del SPA y la API en el mismo puerto

No se agregan endpoints. Se define a quién responde cada petición cuando el directorio de estáticos existe.

| Petición | Responde |
|---|---|
| `GET`/`HEAD` de un archivo existente en `WEB_DIR` (ej. `/assets/x.js`) | el archivo |
| `GET` + `Accept` con `text/html` a `/backlog…`, `/projects…`, `/sprints…` | `index.html` |
| `GET` + `Accept` con `text/html` a una ruta no-API (ej. `/`) | `index.html` |
| `GET /health` (con o sin `text/html`) | `ok` (API) |
| `GET` sin `text/html` a `/projects…` | API (JSON) |
| `POST`/`PUT` a cualquier ruta (aun con `text/html`) | API |
| `GET` sin `text/html` a un archivo inexistente | 404 |

Si `WEB_DIR` no existe: todo va a la API.
