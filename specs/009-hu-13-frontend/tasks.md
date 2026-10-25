# Tasks: Estado del proyecto en el panel principal (HU-13 Frontend)

**Feature Branch**: `feature/HU-13-frontend`

**Issue**: #41 (parte de #13)

## Tareas

- [X] T001 Escribir la especificación del estado del proyecto en el panel.
- [X] T002 COMMIT `docs: agregar la especificación del estado en el panel (Issue #41)`
- [X] T003 RED: tests del panel para los tres estados, la consulta a `GET /projects/{id}/status` y el error sin ocultar los datos del proyecto.
- [X] T004 COMMIT `RED: agregar tests del estado del proyecto en el panel (fallan)`
- [X] T005 GREEN: mostrar el estado en `Panel.tsx` y agregar `obtenerEstadoProyecto` en `api.ts`.
- [X] T006 COMMIT `GREEN: mostrar el estado del proyecto en el panel, tests en verde`
- [X] T007 Actualizar la rama sobre la versión corregida de HU-04 (`feature/HU-04-frontend`), descartando las copias viejas de la base, el README y las rutas.
- [X] T008 COMMIT `Merge branch 'feature/HU-04-frontend' into feature/HU-13-frontend`
- [X] T009 Verificar `npm test`, `npm run lint` y `npm run build` en `web/`.
