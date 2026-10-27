# Tasks: Estado del proyecto en el panel principal (HU-13 Frontend)

**Feature Branch**: `feature/HU-13-frontend-panel`

**Issue**: #41 (parte de #13)

## Tareas

- [X] T001 Escribir la especificación del estado del proyecto en el panel.
- [X] T002 COMMIT `docs: agregar la especificación del estado en el panel (Issue #41)`
- [X] T003 RED: tests del panel para los tres estados, la consulta a `GET /projects/{id}/status` y el error sin ocultar los datos del proyecto.
- [X] T004 COMMIT `RED: agregar tests del estado del proyecto en el panel (fallan)`
- [X] T005 GREEN: mostrar el estado en `Panel.tsx` y agregar `obtenerEstadoProyecto` en `api.ts`.
- [X] T006 COMMIT `GREEN: mostrar el estado del proyecto en el panel, tests en verde`
- [X] T007 Crear la rama `feature/HU-13-frontend-panel` desde `main` y reaplicar los commits propios (spec, RED y GREEN), sin arrastrar los cambios de HU-04 que ya están en `main`.
- [X] T008 COMMIT `docs: mover la spec de HU-13 frontend a specs/014-hu-13-frontend`
- [X] T009 Verificar `npm test`, `npm run lint` y `npm run build` en `web/`.
