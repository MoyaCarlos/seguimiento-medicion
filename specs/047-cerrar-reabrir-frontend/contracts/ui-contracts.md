# Contratos de interfaz: Cerrar y reabrir un proyecto (frontend)

La API es la de la spec 044 (`specs/044-cerrar-reabrir-proyecto/contracts/`); no cambia.

## Cliente HTTP (`features/proyecto/api.ts`)

| Función | Pedido | Éxito | Errores |
|---|---|---|---|
| `cerrarProyecto(id)` | `POST /projects/{id}/cerrar`, sin cuerpo | 200 proyecto con `fecha_cierre` y `estado: Finalizado` | 403 sin permiso, 404, 409 «El proyecto ya está cerrado» o «Cerrá el Sprint activo antes de cerrar el proyecto» |
| `reabrirProyecto(id)` | `POST /projects/{id}/reabrir`, sin cuerpo | 200 proyecto con `fecha_cierre: null` y estado recalculado | 403, 404, 409 «El proyecto no está cerrado» |

## Hook (`compartido/permisos/usePermisosProyecto.ts`)

`usePermisosProyecto(proyectoId: number)` → `{ rol, fechaCierre, cerrado, puedeGestionar, puedeEscribir, aplicarCierre }` (reglas en [data-model.md](../data-model.md)). Un solo `GET /projects/{id}` por montaje.
`useEsScrumMaster(id): boolean` conserva su contrato (devuelve `puedeGestionar`, con la misma lógica de HU-21 para proyectos abiertos).

## Componentes

| Componente | Props | Salida |
|---|---|---|
| `AvisoProyectoCerrado` | `fechaCierre: string \| null` | Sin fecha: nada. Con fecha: `role="status"` con «Proyecto cerrado el AAAA-MM-DD» |
| `NotaCerrado` | `id: string` | `<p id class="nota-gestion">El proyecto está cerrado</p>` para `aria-describedby` |
| `ZonaCierreProyecto` | `proyectoId: string`, `rol: Rol \| null`, `fechaCierre: string \| null`, `onCambio(fechaCierre: string \| null)` | Sección «Cerrar proyecto» (ver abajo) |

### `ZonaCierreProyecto`
- Abierto y Scrum Master: botón «Cerrar proyecto». Al pulsarlo: «¿Cerrar el proyecto?» con «Confirmar» y «Cancelar»; foco en «Cancelar»; Escape cancela.
- Cerrado y Scrum Master: «Proyecto cerrado el AAAA-MM-DD» y botón «Reabrir proyecto» (sin confirmación).
- Product Builder: el botón correspondiente deshabilitado con la nota «Solo el Scrum Master puede cerrar o reabrir el proyecto» (`aria-describedby`).
- En curso: botones deshabilitados y `LoadingIndicator`.
- Errores: 409 y 400 con el mensaje del servidor tal cual; 403 con `AvisoSinPermiso`; 404 con el aviso de proyecto inexistente.
- Éxito: llama `onCambio(respuesta.fecha_cierre)` y muestra el estado devuelto.

## Props de las pantallas (Fases 2 y 3)
Los componentes que escriben reciben `puedeGestionar = true` y `puedeEscribir = true` como props opcionales; la raíz de la pantalla las obtiene de una única llamada a `usePermisosProyecto` y renderiza `AvisoProyectoCerrado`. Si `cerrado`, la nota visible es `NotaCerrado`; si no, la `NotaGestion` de HU-21.
