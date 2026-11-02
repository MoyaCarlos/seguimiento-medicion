# Data Model: Actualización de Estado de una Historia en el Sprint (HU-11)

**Sin cambios de esquema.** La fuente única es [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md);
este documento enlaza a sus secciones y no redefine tablas.

## Tablas que usa

| Tabla | Sección | Uso en HU-11 |
|---|---|---|
| `historias` | [§4.4](../../MODELO-DE-DATOS.md#44-historias--hu-01-hu-15-hu-06-hu-11-hu-14-existe-hoy-como-backlog_items) | Se escriben **solo** `estado` y `fecha_finalizacion`. Se lee para conocer el estado actual, el Sprint asignado y la fecha. |
| `sprints` | [§4.5](../../MODELO-DE-DATOS.md#45-sprints--hu-05-existe-hoy-mismo-nombre) | Solo lectura: `estado` del Sprint de la historia (`Activo` es la única que admite el cambio). |

Mapeo de HU-11 en el documento: [§8](../../MODELO-DE-DATOS.md#8-mapa-historia--tablas), fila
"HU-11 Estado de historia" → `historias.estado`, `fecha_finalizacion`, migración 001.

## Restricciones de base que intervienen

- `CHECK (estado IN ('Nueva','En progreso','Completada'))`: el estado destino siempre es uno de
  los tres; el dominio lo valida antes (400), la base queda como última barrera.
- `CHECK (estado = 'Nueva' OR sprint_id IS NOT NULL)`: impide `En progreso` o `Completada` sin
  Sprint. El servicio lo valida antes (`ErrHistoriaSinSprint`, 400).
- `CHECK ((estado = 'Completada') = (fecha_finalizacion IS NOT NULL))`: obliga a escribir la fecha
  al completar y a dejarla en `NULL` en los demás estados. HU-11 escribe las dos columnas en el
  mismo `UPDATE`, así que la restricción se cumple siempre.
- `CHECK (date(fecha_finalizacion) IS fecha_finalizacion)` (D12): la fecha es de calendario
  (`AAAA-MM-DD`); el repositorio la serializa con el helper `fechaASQL`, que descarta la hora.

Una violación de `CHECK` que llegue igual a la base se informa como error interno (§5): es un
defecto a corregir, no un error del cliente.

## Entidades de dominio

- `domain.BacklogItem`: se usan `ID`, `Estado`, `SprintID` y el **campo nuevo de lectura**
  `FechaFinalizacion *time.Time` (la columna ya existía desde la 001; faltaba exponerla). El resto
  de los campos no cambia al cambiar de estado (FR-006).
- `domain.Sprint`: se usa `Estado` (y `ID`, para verificar que el Sprint de la historia exista).
- `domain.Estado`: suma `EstadoEnProgreso = "En progreso"` y `Estado.EsValida()`, siguiendo el
  patrón de `Prioridad.EsValida()`.

## Regla de transición (`domain.ValidarCambioDeEstado`)

Se evalúa después de comprobar que la historia existe (404 primero, FR-007). Devuelve el primer
error en este orden (FR-018):

| # | Condición | Error | Campo | Mensaje |
|---|---|---|---|---|
| 1 | el estado destino no es uno de los tres | `ErrEstadoInvalido` | `estado` | el estado debe ser Nueva, En progreso o Completada |
| 2 | `historia.SprintID == nil` y el destino no es `Nueva` | `ErrHistoriaSinSprint` | `estado` | primero debe asignarse a un Sprint activo |
| 3 | el Sprint no está `Activo` | `ErrSprintNoActivoParaEstado` | `estado` | solo se cambia el estado de historias de un Sprint activo |
| 4 | el destino es igual al estado actual | `ErrEstadoRepetido` | `estado` | la historia ya está en ese estado |
| 5 | retroceso (`Completada → otra` o `→ Nueva`) | `ErrRetrocesoEstado` | `estado` | el estado de una historia solo avanza: no se puede volver a un estado anterior |
| 6 | `Nueva → Completada` directo | `ErrSaltoDirecto` | `estado` | primero debe pasar por En progreso |

`ErrConflictoEstadoHistoria` (error plano, no `ValidationError`) lo devuelve el repositorio cuando
el `UPDATE` condicional no afecta filas (carrera, [research.md](./research.md) R2) y el handler lo
traduce a 409 sin campo.

## Transiciones válidas

```text
Nueva ──▶ En progreso ──▶ Completada
  │                          ▲
  └──────────────────────────┘
        (rechazado: Q1)

Retrocesos (→ Nueva, Completada → En progreso): rechazados (Q3)
Mismo estado (Nueva → Nueva, etc.): rechazado (Q4)
Sprint Pendiente o Finalizado: rechazado (Q2)
```

Al completar se escribe `fecha_finalizacion` con el día del reloj. La historia **sigue asignada al
Sprint** (no se toca `sprint_id`), que es lo que necesitan HU-03 (Story Points completados) y
HU-05 (no devolverla al Product Backlog al cerrar).

## Vistas (no se guardan)

- **Product Backlog** (HU-15): `historias WHERE proyecto_id = ? AND sprint_id IS NULL`. Toda
  historia sin Sprint está en `Nueva` y no admite esta transición (`ErrHistoriaSinSprint`).
- **Sprint Backlog** (HU-06): `historias WHERE sprint_id = ?`, con el estado y la fecha de
  finalización ya expuestos en la respuesta.

## Métricas (sin cálculo nuevo en esta historia)

HU-11 solo deja los datos que HU-03 necesita: `historia.SprintID` y `estado = Completada` con su
fecha. El cálculo de Story Points completados sigue viviendo en HU-03 (D2b/D2c, §6); esta historia
no agrega fórmulas.
