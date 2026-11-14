# Data Model: Cerrar y reabrir un proyecto (frontend)

No hay almacenamiento nuevo; son tipos del frontend.

## Proyecto (`features/proyecto/api.ts`)
Se agrega a `Proyecto`: `fecha_cierre: string | null` (AAAA-MM-DD si está cerrado, `null` si está abierto) y `estado?: EstadoProyecto`, que solo traen las respuestas de `cerrar` y `reabrir`. `cerrarProyecto` y `reabrirProyecto` devuelven `Proyecto & { estado: EstadoProyecto }`.

## Estado del rol (`useRolEnProyecto`)
`EstadoRol` mantiene sus cuatro variantes; `listo` pasa a `{ estado: 'listo'; rol: Rol; fechaCierre: string | null }`.
Validación: `fecha_cierre` ausente o `null` → `null`; string → se toma tal cual; otro tipo → `error` «No se pudo interpretar el rol en el proyecto».

## Permisos del proyecto (`usePermisosProyecto`)

| Campo | Tipo | Regla |
|---|---|---|
| `rol` | `Rol \| null` | `null` si cargando, error o no integrante |
| `fechaCierre` | `string \| null` | `null` si abierto o desconocido |
| `cerrado` | `boolean` | `fechaCierre !== null` |
| `puedeGestionar` | `boolean` | `!cerrado && rol !== 'product_builder' && estado !== 'no_integrante'` |
| `puedeEscribir` | `boolean` | `!cerrado && estado !== 'no_integrante'` |
| `aplicarCierre` | `(f: string \| null) => void` | reemplaza `fechaCierre` con la respuesta del servidor |

Desconocido (cargando o error) → `puedeGestionar` y `puedeEscribir` son `true`.

## Estado de la zona de cierre (`ZonaCierreProyecto`)

```text
reposo ──(Cerrar proyecto)──▶ confirmando ──(Confirmar)──▶ enviando ──200──▶ reposo (cerrado)
   ▲                              │  ▲                        │
   └────────(Cancelar / Escape)───┘  └── 409 / 403 / 404 / red┘ (vuelve a reposo con el mensaje)
reposo (cerrado) ──(Reabrir proyecto)──▶ enviando ──200──▶ reposo (abierto)
```
