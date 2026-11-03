# Data Model: Imputación de horas en pantalla (HU-07 Frontend)

No hay cambios de esquema. Las cargas son la tabla `registros_horas`: ver [`MODELO-DE-DATOS.md` §4.8](../../MODELO-DE-DATOS.md) (D8: minutos enteros). Los campos de la API están definidos en el [contrato de la spec 019](../019-hu-07-imputacion-horas/contracts/openapi.yaml).

## Tipos que usa la pantalla

Declarados en `web/src/features/sprints/horas.ts`:

| Tipo | Campos | Origen |
|---|---|---|
| `RegistroHoras` | `id`, `integrante`, `fecha`, `actividad`, `horas` | `RegistroHoras` del contrato (solo los campos que se muestran) |
| `EsfuerzoReal` | `total_horas`, `registros` | `EsfuerzoReal` del contrato |
| `IntegranteDeProyecto` | `id`, `nombre` | `GET /projects/{id}/members` |
| `DatosImputacion` | `integrante_id`, `fecha`, `actividad`, `horas` (string) | `ImputarHorasRequest` del contrato |

El `id` de `GET /projects/{id}/members` es el identificador de usuario (`aIntegranteResponse` lo toma de `Member.UserID`), el mismo que `POST /backlog/{id}/horas` espera en `integrante_id`.

## Estado de la pantalla

**Esfuerzo de una historia**: `cargando` → `listo` (total y cargas) o `error` (aviso). Tras guardar una carga, se reconsulta desde el estado actual sin volver a `cargando`.

**Integrantes del formulario**: `cargando` → `listo` o `error` (no permite guardar).

**Formulario**: valores de los cuatro campos, `enviando`, aviso general, errores por campo y confirmación.

| Evento | Efecto |
|---|---|
| Guardar sin integrante | Error en el campo; no se envía |
| Guardar | Limpia errores y confirmación; `enviando` hasta la respuesta |
| Carga aceptada | Confirmación; vacía actividad y horas; avisa para reconsultar el total |
| Carga rechazada | Error por campo o aviso; conserva los valores |
| Cerrar | Se descarta todo el estado del formulario |
