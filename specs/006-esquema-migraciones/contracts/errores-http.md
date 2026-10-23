# Contrato: errores de persistencia en la API

Fuente: [MODELO-DE-DATOS.md §5, "Traducción de violaciones a errores"](../../../MODELO-DE-DATOS.md#5-dónde-vive-cada-regla).
Los endpoints, cuerpos de pedido y respuestas exitosas **no cambian**. El cuerpo de error
sigue siendo el `errorResponse` de hoy (`{"campo": ..., "mensaje": ...}`).

| Causa | Error que devuelve el adaptador | HTTP | ¿Cambia? |
|---|---|---|---|
| Recurso inexistente detectado por el servicio antes de guardar | `ErrProyectoNoEncontrado`, `ErrSprintNoEncontrado` | 404 | No |
| Violación de FK en la base (simple o compuesta) | `ErrIntegridadReferencial` | **409**, `{"mensaje": "violación de integridad referencial"}` | **Sí** (antes caía en 500) |
| Índice parcial de `sprints` (carrera) | `ValidationError{Campo: "estado"}` con `"el proyecto ya tiene un Sprint activo"` o `"... pendiente"` | 400 | No (mismo cuerpo que la validación previa del servicio) |
| `CHECK` o `NOT NULL` | error envuelto, sin tipo propio | 500, `{"mensaje": "error interno"}` | No (D13) |
| Usuario o integrante duplicado | `ErrUsuarioDuplicado`, `ErrMiembroDuplicado` | como hoy | No |

En el uso normal de la API el 409 no debería aparecer: el servicio valida antes de guardar.
Es la red de seguridad para pedidos concurrentes y para defectos del servicio.
