# Contrato: piezas compartidas, registros y carpetas

## Piezas (`web/src/compartido/componentes/`, se conservan)

| Pieza | Qué recibe | Comportamiento |
|---|---|---|
| `Button` | atributos de botón; `variant` (`primary`, `secondary`, `quiet`) | Con `disabled` queda deshabilitado |
| `Field` | `label`, `error?` y atributos de campo | La etiqueta queda asociada al campo (`htmlFor`/`id`); con `error` lo muestra junto al campo y marca el campo como inválido |
| `ErrorMessage` | `message` | Aviso con `role="alert"`; se usa arriba del formulario cuando el error no tiene campo |
| `LoadingIndicator` | — | Texto "Cargando…" con `role="status"` |

Pruebas que se agregan: etiqueta asociada (FR-017), botón deshabilitado e indicador visible durante una operación (FR-018), aviso sin campo arriba del formulario (FR-008).

## Registro de rutas (`web/src/rutas.tsx`)

Una línea por feature, hijas de `projects/:id`:

```text
{ path: 'backlog/*', lazy: () => import('./features/backlog') },    // hoy una sola pantalla
{ path: 'proyecto/*', lazy: () => import('./features/proyecto') },  // con sub-pantallas (settings, edit)
```

## Registro de navegación (`web/src/navegacion.ts`)

Una línea por feature, relativa al proyecto actual:

```text
{ ruta: 'backlog', texto: 'Historias' },
{ ruta: 'proyecto', texto: 'Proyecto' },
```

## Carpetas

- Cada feature en `web/src/features/<nombre>/`, con `index.tsx` como pantalla de entrada.
- `index.tsx` exporta `Component` (lo que pide la carga diferida de React Router).
- Una feature importa solo de `web/src/compartido/` y de su propia carpeta.

## Reglas verificables

- Agregar una pantalla = carpeta nueva en `features/` + 1 línea en cada registro.
- Ninguna feature importa de otra (regla de lint o test de `npm test`).
- `compartido/` y `package.json` no cambian en los PR de pantallas.
