# Quickstart: validar la pantalla de creación de historias (HU-01)

Guía para comprobar la historia de punta a punta. El contrato del endpoint y las reglas
exactas están en [contracts/api-backlog.md](contracts/api-backlog.md) y
[data-model.md](data-model.md).

## Requisitos

- Go (versión del repo) y Node.js 24.15.0 o superior.
- Backend en `http://localhost:8080` y frontend en `http://localhost:5173`.
- Un proyecto existente con id conocido (los ids no se listan todavía).

## Puesta en marcha

Desde la raíz del repositorio:

```bash
go run ./cmd/api
```

En otra terminal, desde `web/`:

```bash
npm install
npm run dev
```

Cargar datos de ejemplo si hace falta (crear proyecto, integrantes, Sprint e historia de
muestra) con los `curl` de
[`docs/curl-frontend-sprint-2.md`](../../docs/curl-frontend-sprint-2.md). Anotar el
`id` del proyecto.

## Validación manual

Abrir `http://localhost:5173/projects/<id>/backlog`.

1. **Creación feliz**: completar título, descripción y prioridad; presionar "Guardar".
   Esperado: una sola llamada `POST /backlog`, confirmación con la historia, estado
   "Nueva", sin estimación, y formulario vacío con prioridad "Must have" y valor
   "Sin valor".
2. **Recarga**: recargar la URL. Esperado: se sigue usando el proyecto de la dirección.
3. **Obligatorios**: guardar sin título o con solo espacios; repetir con descripción.
   Esperado: advertencia junto al campo, sin llamada a la API.
4. **Varios errores**: dejar título vacío y descripción de más de 2000 caracteres.
   Esperado: las dos advertencias a la vez.
5. **Límites**: título de 200 y descripción de 2000 caracteres (por ejemplo con emojis en
   los bordes). Esperado: se aceptan; 201 y 2001 se rechazan localmente.
6. **Prioridad y valor**: probar las cuatro prioridades y cada valor admitido; guardar con
   "Sin valor". Esperado: el cuerpo lleva `M/S/C/W` y el `valor_negocio` correcto, u
   omitido en "Sin valor".
7. **Error de campo**: forzar un rechazo del servidor (por ejemplo, enviar un título que
   el backend considere inválido o simular la respuesta en la prueba automatizada).
   Esperado: mensaje junto al campo, datos conservados.
8. **Proyecto inexistente**: abrir `/projects/999999/backlog` y guardar. Esperado: 404 y
   aviso general "proyecto no encontrado", datos conservados.
9. **Id inválido**: abrir `/projects/abc/backlog` o `/projects/0/backlog`. Esperado:
   aviso general y formulario deshabilitado.
10. **Backend caído**: detener el backend e intentar guardar. Esperado: aviso general y
    pantalla usable (sin quedar cargando).

## Validación automatizada

Desde `web/`:

```bash
npm test
npm run lint
npm run build
```

Esperado: todos los tests en verde (funciones puras del formulario y pruebas de componente
por escenario), lint sin errores y compilación exitosa. Las pruebas de componente simulan
`fetch`; no requieren el backend levantado.

## Criterio de aceptación

La historia se considera verificada cuando los 10 puntos manuales se comportan como se
describe, `npm test` pasa completo y no hay cambios en `web/src/compartido/`,
`web/package.json` ni en las demás pantallas.
