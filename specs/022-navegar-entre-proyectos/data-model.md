# Data Model: Navegar entre proyectos desde la interfaz

**Feature**: #148 (cierra #141) | **Date**: 2026-10-27 | **Spec**: [spec.md](./spec.md)

## Scope

Esta historia es solo de frontend y no modifica el esquema de la base ni el contrato de la API.
El único dato leído es el listado existente `GET /projects` (ver
[spec 010](../010-hu-04-listar-proyectos/spec.md)). Los "modelos" de este documento describen
el dato que consume la interfaz y el estado de presentación; no son tablas nuevas.

## Entities

### ProyectoResumen

Proyección de solo lectura de un proyecto para la navegación.

| Campo | Tipo | Origen | Reglas |
|---|---|---|---|
| `id` | `number` | `GET /projects` | Identificador único; se usa como `value` de la opción y para armar la ruta. |
| `nombre` | `string` | `GET /projects` | Texto visible en la etiqueta; no se recorta ni se transforma. |

La etiqueta visible se deriva con `etiquetaProyecto`:

```text
<nombre> (#<id>)
```

Ejemplo: `{ id: 7, nombre: 'Alfa' }` → `Alfa (#7)`.

Si el proyecto actual no está en el listado, el selector agrega una opción sintética con la
etiqueta `Proyecto <id>`. No es un `ProyectoResumen` del listado y no se usa en la lista del
inicio.

### SeccionDeProyecto

Tramo de navegación que se conserva al cambiar de proyecto. No se persiste.

| Valor | Origen | Regla |
|---|---|---|
| `proyecto` | Primer tramo o valor por defecto | Se usa cuando la URL es `/projects/<id>` sin tramo. |
| `backlog` | Primer tramo | Se conserva al cambiar de proyecto. |
| `sprints` | Primer tramo | Se conserva al cambiar de proyecto. |
| `status` | Primer tramo | Se conserva al cambiar de proyecto. |
| Otro texto | Primer tramo | Se conserva tal cual; la pantalla de destino puede no existir. |

`rutaDeSeccion(pathname, idNuevo)` devuelve siempre `/projects/<idNuevo>/<sección>`. Los tramos
posteriores al primero (por ejemplo `edit` en `/proyecto/edit`) se descartan.

### EstadoProyectos

Estado de presentación que produce `useProyectos()`. Es una unión discriminada para que no existan
combinaciones imposibles.

| Estado | Campos | Significado |
|---|---|---|
| `cargando` | — | La lista todavía no respondió. |
| `listo` | `proyectos: ProyectoResumen[]` | La API respondió un arreglo válido (puede estar vacío). |
| `error` | `mensaje: string` | La API falló o respondió algo que no es un arreglo. |

## State Transitions

| Desde | Evento | Hacia | Efecto visible |
|---|---|---|---|
| (montaje) | Se inicia `GET /projects` | `cargando` | `LoadingIndicator` en el selector o en la lista. |
| `cargando` | Respuesta `200` con arreglo | `listo` | Opciones del selector o enlaces de la lista. |
| `cargando` | Error del cliente o respuesta no interpretable | `error` | `ErrorMessage` con el mensaje; el enlace "Nuevo proyecto" y la navegación de secciones siguen disponibles. |
| `listo` | Desmontaje | — | No hay transición; el componente ya no existe. |
| `error` | Desmontaje | — | No hay transición; el componente ya no existe. |

No hay transición de `listo` a `error` ni reintentos automáticos: el hook se ejecuta una vez por
montaje. Volver a montar la pantalla vuelve a pedir el listado.

## Derived Values

| Función | Entrada | Salida | Regla |
|---|---|---|---|
| `etiquetaProyecto` | `{ id: number, nombre: string }` | `string` | `${nombre} (#${id})`; dos nombres iguales producen etiquetas distintas por el id. |
| `rutaDeSeccion` | `pathname: string`, `idNuevo: string \| number` | `string` | Toma el primer tramo después del id; si no hay, usa `proyecto`; devuelve `/projects/${idNuevo}/${seccion}`. |

## Validation Rules

- El `:id` de la URL se compara con `String(proyecto.id)` para marcar el actual; la URL entrega
  strings y la API números.
- Si el proyecto actual no está en el listado y hay un `:id` no vacío, se agrega la opción
  `Proyecto <id>` y se selecciona.
- Una respuesta exitosa que no sea un arreglo se trata como error; no se renderiza una lista
  vacía como si no hubiera datos.
- El listado se muestra en el orden devuelto por la API; no se ordena ni se pagina en el cliente.
- No se valida el nombre del proyecto en el frontend: se muestra tal como llega.

## Relationships

- Un `ProyectoResumen` puede aparecer en el selector del encabezado y en la lista del inicio.
- El selector usa el `:id` de la ruta actual para marcar el proyecto actual.
- La lista del inicio no tiene proyecto actual; todos sus enlaces abren
  `/projects/<id>/proyecto`.
- Cada pantalla monta su propio consumidor de `useProyectos`; no hay estado compartido entre el
  selector y la lista.

## Out of Scope

- Persistencia de proyectos, del último proyecto o de la última sección.
- Paginación, filtros, búsqueda u ordenamiento distinto al de la API.
- Creación, edición o borrado de proyectos desde el selector o la lista.
- Campos adicionales del proyecto (descripción, fechas, integrantes): no se muestran en esta
  historia.

## Traceability

| Requisito | Modelo / función |
|---|---|
| FR-001, FR-002 | `ProyectoResumen`, `etiquetaProyecto`, `:id` actual |
| FR-003, FR-004, FR-005 | `SeccionDeProyecto`, `rutaDeSeccion` |
| FR-006 | Enlace "Nuevo proyecto" (sin entidad nueva) |
| FR-007, FR-008, FR-009 | `ProyectoResumen`, `ListaProyectos`, `EstadoProyectos` |
| FR-010, FR-011, FR-012 | `EstadoProyectos` |
| FR-013 | Reinicio por `:id` en `LayoutProyecto` |
| FR-014 | Orden del arreglo de `GET /projects` |
| FR-015 | Eliminación de `Inicio.tsx` |
| FR-016, FR-017 | Sin cambios de API/esquema; reutilización de piezas |
| FR-018, FR-019, FR-020 | Funciones puras y pruebas de componente |
