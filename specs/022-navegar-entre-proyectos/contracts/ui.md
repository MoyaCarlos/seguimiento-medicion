# UI Contract: Navegar entre proyectos desde la interfaz

**Feature**: #148 (cierra #141) | **Date**: 2026-10-27 | **Spec**: [spec.md](../spec.md)

## Purpose

Esta historia no expone una API nueva ni cambia un contrato existente. Define el contrato de la
interfaz que consumen las pantallas y el uso que hace del endpoint existente `GET /projects`.

## Consumed API Contract

- **Endpoint**: `GET /projects`, ya existente y sin cambios (ver
  [spec 010](../../010-hu-04-listar-proyectos/spec.md)).
- **Cliente**: `web/src/compartido/api/client.ts` (`get`) con `Accept: application/json`.
- **Respuesta exitosa**: `200` con un arreglo de proyectos. La interfaz solo usa `id` y `nombre`.
  El orden del arreglo se conserva en la lista y en el selector.
- **Lista vacía**: `200` con `[]`; se muestra la invitación a crear el primer proyecto en el
  inicio y la opción sintética del proyecto actual en el encabezado.
- **Errores**: el cliente entrega `ApiError` (mensaje de la API o genérico) o `ConnectionError`.
  La interfaz muestra el mensaje y no inventa datos.
- **Respuesta no interpretable**: JSON inválido lo detecta el cliente; un `200` que no sea un
  arreglo lo detecta `useProyectos` y pasa a estado de error.

No se modifica `GET /projects`, ni su contrato OpenAPI, ni el esquema de la base.

## Pure Functions

### `rutaDeSeccion(pathname, idNuevo)`

- **Entrada**: pathname actual (`string`) e identificador del proyecto elegido (`string` o
  `number`).
- **Salida**: `/projects/<idNuevo>/<sección>`.
- **Regla**: `<sección>` es el primer tramo después del `:id`; si no hay tramo, es `proyecto`.
  Los tramos posteriores se descartan.

| Entrada | Salida |
|---|---|
| `('/projects/1/backlog', 2)` | `/projects/2/backlog` |
| `('/projects/1/proyecto/edit', 2)` | `/projects/2/proyecto` |
| `('/projects/1', 2)` | `/projects/2/proyecto` |
| `('/projects/1/otra', 2)` | `/projects/2/otra` |
| `('/projects/1/backlog/', '2')` | `/projects/2/backlog` |

### `etiquetaProyecto({ id, nombre })`

- **Entrada**: `{ id: number, nombre: string }`.
- **Salida**: `<nombre> (#<id>)`.
- **Regla**: no recorta ni transforma el nombre; dos nombres iguales producen etiquetas distintas
  por el `#id`.

| Entrada | Salida |
|---|---|
| `{ id: 7, nombre: 'Alfa' }` | `Alfa (#7)` |
| `{ id: 8, nombre: 'Alfa' }` | `Alfa (#8)` |

## `useProyectos()` Contract

Devuelve una unión discriminada:

```text
{ estado: 'cargando' }
{ estado: 'listo'; proyectos: ProyectoResumen[] }
{ estado: 'error'; mensaje: string }
```

- Hace un `GET /projects` al montarse.
- No actualiza estado después del desmontaje.
- Conserva el orden recibido.
- Un error del cliente se muestra con su mensaje; una respuesta que no es arreglo se trata como
  error.
- No reintenta, no hace polling y no comparte estado con otros consumidores.

## `SelectorProyecto` Contract

- Se monta en el encabezado de toda pantalla `/projects/:id/...`.
- **Control**: `<select>` con etiqueta accesible **"Proyecto"**.
- **Opciones**: una por proyecto, con etiqueta `<nombre> (#<id>)`.
- **Actual**: `value` igual al `:id` de la URL; se compara como string.
- **Actual ausente del listado**: se agrega y selecciona una opción con etiqueta `Proyecto <id>`.
- **Cambio**: si el id elegido es distinto, navega con `rutaDeSeccion`; si es el mismo, no navega.
- **Carga**: muestra `LoadingIndicator` junto a un `<select>` deshabilitado con la única opción
  `Proyecto <id>` del proyecto actual (seleccionada).
- **Error**: muestra `ErrorMessage` con el mensaje del hook junto a un `<select>` deshabilitado
  con la única opción `Proyecto <id>` del proyecto actual (seleccionada).
- **Listo**: el `<select>` se habilita con la lista completa de proyectos en el orden de la API.
- **Enlace permanente**: "Nuevo proyecto" apunta a `/` y está visible y operativo en carga, error
  y listo; los enlaces de sección del encabezado tampoco dependen del listado.

## `ListaProyectos` Contract

- Se monta en `/`, junto al formulario de creación (composición en `rutas.tsx`).
- **Lista**: un enlace por proyecto, con etiqueta `<nombre> (#<id>)`, a
  `/projects/<id>/proyecto`, en el orden de la API.
- **Carga**: muestra `LoadingIndicator`.
- **Error**: muestra `ErrorMessage` con el mensaje del hook; el formulario de creación sigue
  disponible.
- **Vacía**: muestra una invitación a crear el primer proyecto, sin tratarla como error.

## Home Composition Contract

- `/` sigue mostrando la pantalla existente `features/proyecto/crear` sin modificarla.
- La ruta índice es hija de un layout pathless en `rutas.tsx` que renderiza el formulario
  (`Outlet`) y `ListaProyectos`.
- La composición no altera el comportamiento de creación ni sus pruebas.

## Accessibility and Test Hooks

- El selector debe ser encontrable como `combobox` con nombre accesible "Proyecto", incluso
  cuando está deshabilitado durante la carga o el error.
- El enlace debe ser encontrable como `link` con nombre "Nuevo proyecto".
- Los estados de carga usan `role="status"` (provisto por `LoadingIndicator`).
- Los errores usan `role="alert"` (provisto por `ErrorMessage`).
- Las pruebas usan `getByRole`/`findByRole` y `fetch` simulado; no dependen de clases CSS.

## Out of Scope

- Nuevos endpoints o parámetros de `GET /projects`.
- Autenticación, permisos, paginación, búsqueda, filtros y ordenamiento distinto.
- Persistencia del proyecto o de la sección elegida.
- Cambios en los contratos de creación, edición o consulta de proyectos existentes.
