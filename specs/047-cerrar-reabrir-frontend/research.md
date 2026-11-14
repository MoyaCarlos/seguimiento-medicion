# Research: Cerrar y reabrir un proyecto (frontend)

No quedaban `NEEDS CLARIFICATION`. Se verificó el código actual y el contrato de la spec 044.

## Hallazgos del código

- `useRolEnProyecto` ya hace `GET /projects/{id}` y lee solo `rol`; la respuesta de la spec 044 incluye `fecha_cierre` (`internal/http/project_dto.go`). No hace falta ningún pedido nuevo.
- `useEsScrumMaster` lo usan Editar, Configuracion, Sprints y Votación; `permisos.test.tsx` cubre ambos hooks y no se modifica (SC-006).
- Editar ya consulta el proyecto con `obtenerProyecto` y el hook hace su propio `GET` (dos `GET` iguales). Es el comportamiento de HU-21; no se cambia.
- Backlog, Defectos y la pantalla de Equipo (Configuracion) todavía no consultan el rol en todas sus acciones: entran en la Fase 3 llamando al hook.

## Decisiones

### D1. Un solo hook, `usePermisosProyecto(id)`
- **Decisión**: `useRolEnProyecto` suma `fechaCierre: string | null` al estado `listo`; `usePermisosProyecto` lo envuelve y calcula `cerrado = fechaCierre != null`, `puedeGestionar = !cerrado && rol !== 'product_builder'`, `puedeEscribir = !cerrado`. Con rol o cierre desconocidos (cargando o error) ambos son `true`. `no_integrante` mantiene `puedeGestionar = false`, como hoy.
- **Rationale**: un pedido por pantalla (FR-007a) y retrocompatible: `useEsScrumMaster` devuelve lo mismo que hoy con proyectos abiertos.
- **Alternativas**: un Context de proyecto (cambia `LayoutProyecto`, prohibido); un segundo hook con otro `GET` (viola FR-007a).

### D2. `aplicarCierre` en el resultado del hook (único ajuste)
- **Problema**: tras cerrar o reabrir hay que reflejar la respuesta del servidor sin recargar (FR-005). El hook guarda el dato en su estado y no hay forma de actualizarlo desde la zona.
- **Decisión**: `usePermisosProyecto` expone también `aplicarCierre(fechaCierre | null)`, que reemplaza la fecha en su estado. Editar lo llama con `fecha_cierre` de la respuesta de `cerrar` o `reabrir`.
- **Alternativas**: estado local en Editar que pise al hook (duplica la regla `cerrado`); volver a pedir el proyecto (pedido extra y ventana de datos viejos).
- **Alcance**: otras pantallas ya montadas no se enteran, pero cada pantalla se monta de nuevo al navegar y vuelve a pedir el dato.

### D3. La zona no usa `puedeGestionar`
Con el proyecto cerrado `puedeGestionar` es `false` para todos, pero el Scrum Master tiene que poder reabrir. La zona usa `rol`: Scrum Master (o rol desconocido, decide el servidor) → habilitada; Product Builder → deshabilitada con la nota «Solo el Scrum Master puede cerrar o reabrir el proyecto» (patrón de HU-21, FR-009).

### D4. Una sola nota cuando coinciden cierre y falta de permiso
Las pantallas muestran `NotaCerrado` si `cerrado`; si no, la `NotaGestion` de HU-21 si falta permiso. Nunca las dos. `NotaCerrado` reutiliza el elemento `p.nota-gestion` con el texto «El proyecto está cerrado», para heredar el estilo sin tocar `App.css`; va en archivo nuevo para dejar `NotaGestion` intacta.

### D5. Confirmación en línea
`ZonaCierreProyecto` tiene tres estados: `reposo`, `confirmando`, `enviando`. Al pasar a `confirmando` el foco va a «Cancelar» (ref + `useEffect`); Escape (keydown en la zona) cancela salvo en `enviando`; en `enviando` ambos botones quedan `disabled` (sin doble envío). Un 409 cierra la confirmación y muestra el mensaje del servidor junto a la zona con `ErrorMessage`; 403 muestra `AvisoSinPermiso`; 404, el aviso estándar de proyecto inexistente; la conexión caída, el mensaje de `ConnectionError`. Los errores se mapean con `mensajeDeError` y `ApiError.tipo` ya existentes en `features/proyecto/errores.ts`.

### D6. Bloqueo de escritura por pantalla
Cada pantalla de las Fases 2 y 3 llama una vez a `usePermisosProyecto` en su raíz. Las acciones que ya dependían de `puedeGestionar` siguen igual (ahora incluye el cierre); las que cualquiera puede hacer (crear Historia, Defecto) reciben `puedeEscribir`. Los componentes hijos reciben ambas como props con default `true`.

### D7. Sin dependencias, sin tocar el armazón
El aviso se monta desde la raíz de cada pantalla (no desde `LayoutProyecto`). Los estilos del aviso y de la zona van en `features/proyecto/estilos.css`; los de cada pantalla, en su propio CSS.

## Riesgos

- **Rol cargando**: durante la carga todo queda habilitado y puede aparecer un parpadeo del aviso. Es el comportamiento acordado (decide el servidor).
- **Dato viejo en otra pestaña**: si otra persona cierra el proyecto, se ve el 409 «El proyecto está cerrado» tal cual (US2.4).
