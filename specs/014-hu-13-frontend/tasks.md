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

## Fix #151 (post-merge de #150)

**Fix Branch**: `fix/151-spec-hu-13-frontend` (desde `main`)

- [X] T010 Completar `specs/014-hu-13-frontend/spec.md` con Restricciones, Casos límite y Condiciones de error; actualizar el estado y la dependencia con #121.
- [X] T011 COMMIT `docs: completar la spec de HU-13 frontend (Refs #151)`
- [X] T012 RED: tests del cambio de `id`, de `Consultando…` y de `Sin datos` en `web/src/features/proyecto/Panel.test.tsx`; confirmar que fallan.
- [X] T013 COMMIT `RED: agregar tests del cambio de proyecto y del estado no inventado (fallan)`
- [X] T014 GREEN: reiniciar el panel al cambiar de proyecto para no mostrar el estado anterior.
- [X] T015 COMMIT `GREEN: reiniciar el estado al cambiar de proyecto, tests en verde`
- [X] T016 REFACTOR: reiniciar por `key={id}` en vez de `setState` en los efectos.
- [X] T017 COMMIT `REFACTOR: reiniciar el panel por key en vez de setState en los efectos`
- [X] T018 Verificar `npm run lint`, `npm test` y `npm run build` en `web/`.
