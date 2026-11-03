# Data Model: Pantalla de Sprints

La pantalla no persiste datos ni define tablas; el esquema vive en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md) y no se toca. Estas son las formas que maneja el navegador.

## Sprint
Tal como lo devuelve la API (ver el esquema `Sprint` en [`specs/003-apertura-cierre-sprints/contracts/openapi.yaml`](../003-apertura-cierre-sprints/contracts/openapi.yaml)).
- `id`: número; se muestra como "Sprint #<id>".
- `proyecto_id`: número.
- `sprint_goal`, `fecha_inicio`, `fecha_fin`: texto o ausentes mientras el Sprint está Pendiente.
- `estado`: `Pendiente` | `Activo` | `Finalizado`.

## Sprints de la pantalla
Lista en memoria de la pantalla. Reglas:
- Una respuesta de la API reemplaza al Sprint con el mismo `id`, o lo agrega si es nuevo.
- Se agrupa por `estado`; cada grupo se ordena por `id` ascendente.
- Empieza vacía y se carga con `GET /projects/{id}/sprints` al abrir la pantalla o cambiar de proyecto (spec 015); los Sprints creados, iniciados o cerrados se siguen actualizando con `guardarSprint`.

## Transiciones que la pantalla ofrece
| Estado | Acción ofrecida | Resultado confirmado por la API |
|---|---|---|
| (nuevo) | "Crear Sprint" | Pendiente |
| Pendiente | formulario de inicio | Activo |
| Activo | "Cerrar Sprint" | Finalizado |
| Finalizado | ninguna | — |

Las reglas que permiten o rechazan cada transición las aplica la API.

## Clasificación de un error
Entrada: error lanzado por el cliente + nombres de los campos del formulario.
Salida: `{ campo, mensaje }` si es un `ApiError` cuyo `field` es un campo del formulario; si no, `{ aviso }`.
