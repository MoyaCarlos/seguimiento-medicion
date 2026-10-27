# Contrato: llamadas a la API y piezas de la pantalla

## API
No se redefine. Contrato en [`specs/003-apertura-cierre-sprints/contracts/openapi.yaml`](../../003-apertura-cierre-sprints/contracts/openapi.yaml):

| Acción | Llamada | Cuerpo | Respuesta |
|---|---|---|---|
| Crear | `POST /sprints` | `{ proyecto_id }` | Sprint (Pendiente) |
| Iniciar | `POST /sprints/{id}/iniciar` | `{ sprint_goal, fecha_inicio, fecha_fin }` | Sprint (Activo) |
| Cerrar | `POST /sprints/{id}/cerrar` | sin cuerpo | Sprint (Finalizado) |

Errores: `{ campo?, mensaje }` con 400/404, ya interpretados por el cliente compartido ([`specs/007-frontend-base/contracts/cliente-api.md`](../../007-frontend-base/contracts/cliente-api.md)). El cierre no lleva cuerpo: se envía con `post(ruta, undefined)`, que el cliente compartido ya admite (sin cuerpo ni `Content-Type`).

Campos que devuelve la API en errores de estos endpoints: `proyecto_id`, `estado`, `sprint_goal`, `fecha_inicio`, `fecha_fin`, `id`. El formulario de inicio muestra junto al campo solo `sprint_goal`, `fecha_inicio` y `fecha_fin`; el resto va al aviso.

## Módulos de la pantalla (`web/src/features/sprints/`)

| Módulo | Expone | Recibe / devuelve |
|---|---|---|
| `api.ts` | `crearSprint(proyectoId)`, `iniciarSprint(id, datos)` | `Promise<Sprint>` |
| `BotonCerrar.tsx` | `BotonCerrar` (único export). `cerrarSprint(id)` es una función interna del mismo archivo, no exportada: exportarla junto al componente dispara la regla `react/only-export-components` de oxlint. FR-007 se mantiene: la llamada y el botón siguen aislados en un solo archivo, y se prueban a través del componente | props `{ sprint, alActualizar(sprint) }` |
| `FormularioInicio.tsx` | `FormularioInicio` | props `{ sprint, alActualizar(sprint) }` |
| `SeccionSprints.tsx` | `SeccionSprints` | props `{ estado, sprints, alActualizar }` |
| `sprint.ts` | `Sprint`, `ESTADOS`, `agruparPorEstado`, `guardarSprint` | funciones puras |
| `errores.ts` | `clasificarError(error, campos)` | `{ campo, mensaje } \| { aviso }` |
| `index.tsx` | `Component` | lee `:id` de la URL |
