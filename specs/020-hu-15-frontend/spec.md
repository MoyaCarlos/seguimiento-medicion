# Feature Specification: Product Backlog en pantalla (HU-15 Frontend)

**Feature Branch**: `feature/HU-15-frontend`

**Issue**: #98 (parte de #43)

**Created**: 2026-10-29

**Status**: Implementada (pendiente de PR y review)

**Dependencia**: el backend de HU-15 ya está en `main` (PR #135). Esta spec solo cubre
la pantalla: consume `GET /projects/{id}/backlog` y `PUT /backlog/{id}` sin cambiar el
contrato. La pantalla de creación de historias de HU-01 se conserva, movida a la
sub-ruta `backlog/nueva`.

## Objetivo

Que un Product Builder vea el Product Backlog de un proyecto ordenado por prioridad y
pueda editar el título, la descripción, los criterios de aceptación y los Story Points
de una historia desde la interfaz, sin salir de la pantalla ni recargar la aplicación.

## Entrada y salida

- **Entrada (listado)**: `GET /projects/{id}/backlog`, que devuelve un arreglo de
  historias sin Sprint asignado, ya ordenado por prioridad MoSCoW y, a igual prioridad,
  por fecha de creación (`id` ascendente). Cada historia incluye `id`, `proyecto_id`,
  `titulo`, `descripcion`, `criterios_aceptacion`, `prioridad`, `estado`,
  `valor_negocio`, `estimacion_sp` y `sprint_id`.
- **Entrada (edición)**: `PUT /backlog/{id}` con
  `{ titulo, descripcion, criterios_aceptacion, estimacion_sp }`. Los cuatro campos son
  obligatorios como claves; `criterios_aceptacion` y `estimacion_sp` aceptan `null`.
- **Salida en pantalla**: el listado del backlog con todos sus datos y, al elegir una
  historia, un formulario de edición con los cuatro campos editables. Al guardar, la
  historia actualizada se refleja en el listado y se muestra una confirmación.
- **Salida ante errores**: el mensaje de la API tal cual, junto al campo cuando la API
  informa `campo`, o como aviso general cuando no hay campo.

## Reglas de comportamiento (EARS)

- **WHEN** se abre el Product Backlog de un proyecto, **THEN** el sistema **shall**
  consultar `GET /projects/{id}/backlog`.
- **WHILE** el listado está en curso, **THEN** la pantalla **shall** mostrar un indicador
  de carga.
- **WHEN** la API devuelve historias, **THEN** la pantalla **shall** mostrarlas en el
  mismo orden recibido, sin reordenarlas ni filtrarlas en el cliente.
- **WHEN** la API devuelve un arreglo vacío, **THEN** la pantalla **shall** mostrar un
  mensaje de backlog vacío, sin tratarlo como error.
- **WHEN** se elige una historia, **THEN** la pantalla **shall** abrir su formulario de
  edición con los valores actuales de título, descripción, criterios de aceptación y
  Story Points.
- **WHEN** se guarda una edición con datos válidos, **THEN** el sistema **shall** enviar
  los cuatro campos juntos con `PUT /backlog/{id}` y **shall** reemplazar la historia
  del listado por la que devuelve la API.
- **WHEN** una edición se guarda correctamente, **THEN** la pantalla **shall** conservar
  visible el identificador y el estado de la historia y mostrar una confirmación.
- **WHEN** se guardan los criterios de aceptación vacíos o compuestos solo por espacios,
  **THEN** el sistema **shall** enviar `null`.
- **WHERE** una historia no tiene Story Points, la pantalla **shall** ofrecer la opción
  "Sin estimar" y **shall** enviar `null`.
- **IF** el título queda vacío, compuesto solo por espacios o supera los 200 caracteres,
  **THEN** la pantalla **shall** rechazar el envío con una advertencia junto al campo y
  **shall** not llamar a la API.
- **IF** la descripción queda vacía, compuesta solo por espacios o supera los 2000
  caracteres, **THEN** la pantalla **shall** rechazar el envío con una advertencia junto
  al campo y **shall** not llamar a la API.
- **IF** los criterios de aceptación superan los 2000 caracteres, **THEN** la pantalla
  **shall** rechazar el envío con una advertencia junto al campo y **shall** not llamar a
  la API.
- **IF** la API rechaza la edición con un `campo` conocido, **THEN** la pantalla
  **shall** mostrar el mensaje de la API junto a ese campo y **shall** conservar los
  datos escritos por la persona.
- **IF** la API rechaza la edición sin `campo` (por ejemplo, 404 de historia no
  encontrada, error 500 o red no disponible), **THEN** la pantalla **shall** mostrar el
  mensaje como aviso general y **shall** dejar el formulario usable.
- **WHEN** se cancela la edición, **THEN** la pantalla **shall** cerrar el formulario sin
  llamar a la API y sin modificar la historia.
- **WHEN** cambia el proyecto de la URL, **THEN** la pantalla **shall** descartar el
  listado, la edición y los avisos del proyecto anterior antes de consultar el nuevo.
- El cálculo del orden, la validación definitiva y la regla D2b (no cambiar Story Points
  ya cargados de una historia de un Sprint `Activo` o `Finalizado`) **shall** residir en
  la API; la pantalla solo muestra el resultado o el error recibido.

## Restricciones

- La pantalla **shall** reutilizar el cliente compartido
  (`web/src/compartido/api/client.ts`) para todas las llamadas HTTP.
- La pantalla **shall** seguir el patrón de `web/src/features/` (feature aislada, sin
  importar otras features).
- La pantalla **shall not** agregar dependencias nuevas ni cambiar el contrato de la API.
- La creación de historias de HU-01 **shall** seguir disponible en
  `/projects/:id/backlog/nueva`; el listado de HU-15 vive en `/projects/:id/backlog`.
- La pantalla **shall not** editar prioridad, valor de negocio, estado ni Sprint: quedan
  fuera de HU-15.

## Casos límite

- **IF** el proyecto no existe y la API responde 404, **THEN** la pantalla **shall**
  mostrar "proyecto no encontrado" y **shall** dejar usable el enlace a crear historia.
- **IF** el identificador de la URL no es un número válido, **THEN** la pantalla **shall**
  mostrar el error de la API sin inventar un listado.
- **WHEN** una historia tiene criterios de aceptación `null`, **THEN** el formulario
  **shall** mostrarlos vacíos y **shall** enviar `null` si no se completan.
- **WHEN** una historia tiene un Story Point de la escala Fibonacci, **THEN** el selector
  **shall** mostrarlo seleccionado.
- **WHEN** dos personas editan la misma historia, **THEN** la pantalla **shall** mostrar
  el último guardado válido que devuelva la API, sin versionado propio.
- **WHEN** se edita una historia y luego otra, **THEN** la pantalla **shall** mantener el
  listado actualizado y no mezclar los errores de una con la otra.

## Condiciones de error

- **IF** falla la consulta del listado (red, 404 o 500), **THEN** la pantalla **shall**
  mostrar el mensaje y **shall not** mostrar historias de una consulta anterior.
- **IF** la API responde 400 por Story Points fuera de escala o por la regla D2b,
  **THEN** la pantalla **shall** mostrar el mensaje junto al selector de Story Points y
  **shall** conservar los demás valores del formulario.
- **IF** la API responde 404 al guardar, **THEN** la pantalla **shall** mostrar "historia
  no encontrada" como aviso general y **shall not** cerrar el formulario.
- **IF** la API responde un error de campo desconocido, **THEN** la pantalla **shall**
  mostrarlo como aviso general, sin romper el formulario.

## Criterios de aceptación

1. **Given** un proyecto con historias en su Product Backlog, **When** abro
   `/projects/:id/backlog`, **Then** veo las historias en el orden que devuelve la API,
   con título, prioridad, estado, valor de negocio y Story Points.
2. **Given** un proyecto sin historias sin Sprint, **When** abro su Product Backlog,
   **Then** veo un mensaje de backlog vacío y ningún error.
3. **Given** un proyecto inexistente, **When** abro su Product Backlog, **Then** veo el
   mensaje de error de la API.
4. **Given** una historia del listado, **When** abro su edición y guardo cambios válidos,
   **Then** se envía `PUT /backlog/{id}` con los cuatro campos, la historia se actualiza
   en el listado y veo una confirmación.
5. **Given** el formulario de edición, **When** dejo el título vacío, **Then** veo una
   advertencia junto al campo y no se llama a la API.
6. **Given** el formulario de edición, **When** ingreso criterios de más de 2000
   caracteres, **Then** veo una advertencia y no se llama a la API.
7. **Given** una historia con Story Points, **When** la API rechaza cambiarlos por la
   regla D2b, **Then** veo el mensaje de la API junto al selector y la historia queda sin
   modificar.
8. **Given** una historia existente, **When** la API responde 404 al guardar, **Then** veo
   el aviso general y el formulario sigue abierto con los datos.
9. **Given** el Product Backlog de un proyecto, **When** cambio a otro proyecto,
   **Then** no veo las historias del anterior.
10. **Given** la pantalla del listado, **When** elijo "Nueva historia", **Then** llego a
    la pantalla de creación de HU-01.

## Fuera de alcance

- Crear historias desde esta pantalla (sigue en `backlog/nueva`, HU-01).
- Editar prioridad, valor de negocio, estado o Sprint.
- Mostrar descripción y criterios con Markdown (issue #147).
- Paginación, búsqueda o filtros.
- Reordenar el listado en el cliente.

## Trazabilidad

| Criterio | Implementación | Prueba |
|---|---|---|
| Listado y orden de la API | `Listado` + `listarBacklog` | `US1.1` de `Listado.test.tsx` |
| Backlog vacío | `Listado` | `US1.2` de `Listado.test.tsx` |
| Proyecto inexistente y conexión | `Listado` + cliente compartido | `US1.3` y caso de red de `Listado.test.tsx` |
| Edición válida | `Edicion` + `editarHistoria` | `US2.1` de `Listado.test.tsx` |
| Título y descripción obligatorios | `validarEdicion` | `US2.2`, `US2.3` y `edicion.test.ts` |
| Criterios de hasta 2000 | `validarEdicion` | `US2.4` y `edicion.test.ts` |
| Error de campo de la API | `Edicion` | `US2.5` de `Listado.test.tsx` |
| Error general / 404 al guardar | `Edicion` | `US2.6` de `Listado.test.tsx` |
| Cambio de proyecto | `Listado` (`useEffect` por `id`) | caso de cambio de proyecto |
| Creación de HU-01 conservada | ruta `backlog/nueva` | `rutas.test.tsx` |
