# Data Model: Consultar una historia por su identificador

**Sin cambios de modelo.** No hay migración y `MODELO-DE-DATOS.md` no se modifica. La feature solo lee la entidad existente.

## Entidad: Historia de usuario (`domain.BacklogItem`)

| Campo de la API | Tipo | Nulo | Notas |
|-----------------|------|------|-------|
| `id` | entero | no | único en todo el sistema |
| `proyecto_id` | entero | no | |
| `titulo` | texto | no | |
| `descripcion` | texto | no | |
| `criterios_aceptacion` | texto | sí | |
| `prioridad` | texto | no | MoSCoW |
| `estado` | texto | no | |
| `valor_negocio` | entero | sí | |
| `estimacion_sp` | entero | sí | |
| `sprint_id` | entero | sí | `null` si no está en un Sprint (FR-004) |
| `fecha_finalizacion` | texto | sí | `null` si no está completada (FR-005); mismo formato que el resto de la API |

## Validación de entrada

- `id` de la ruta: entero positivo en `int64`. Cualquier otro valor (`abc`, `0`, `-3`, `1.5`, fuera de rango) → 400, campo `id`, `identificador inválido`. Ceros a la izquierda se aceptan (`007` = 7).

## Transiciones de estado

Ninguna: la consulta es de solo lectura (FR-008).
