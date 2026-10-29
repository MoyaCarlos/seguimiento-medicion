# Tasks: Product Backlog en pantalla (HU-15 Frontend)

**Feature Branch**: `feature/HU-15-frontend`

**Issue**: #98 (parte de #43)

## Fase 1: Especificación

- [X] T001 Escribir `specs/020-hu-15-frontend/spec.md` con objetivo, entradas/salidas, reglas EARS, restricciones, casos límite, errores y criterios de aceptación.
- [X] T002 Escribir `specs/020-hu-15-frontend/plan.md` con el diseño de la pantalla y las decisiones.
- [X] T003 COMMIT `docs: agregar la especificación del Product Backlog en pantalla (Issue #98)`

## Fase 2: RED (pruebas que fallan)

- [X] T004 Mover la pantalla de creación de HU-01 a `crear.tsx` y su test a `crear.test.tsx`; dejar `index.tsx` como shell de rutas.
- [X] T005 Escribir `edicionHistoria.test.ts` con la validación de título, descripción, criterios y Story Points, y el armado del cuerpo (`null` para vacíos).
- [X] T006 Escribir `Listado.test.tsx` con los escenarios de listado, backlog vacío, proyecto inexistente, error de red y edición válida.
- [X] T007 Agregar a `Listado.test.tsx` los escenarios de error: título/descripción vacíos, criterios largos, error de campo de la API, 404, cancelación y cambio de proyecto.
- [X] T008 Ajustar `rutas.test.tsx` para simular el listado y cubrir `backlog/nueva`.
- [X] T009 COMMIT `RED: agregar pruebas del listado y la edición del backlog (fallan)`

## Fase 3: GREEN (implementación)

- [X] T010 Agregar los tipos y las llamadas `listarBacklog` y `editarHistoria` en `api.ts`/`historia.ts`.
- [X] T011 Implementar `Edicion.tsx` con el formulario controlado, la validación y el mapeo de errores de la API.
- [X] T012 Implementar `Listado.tsx` con carga, listado, edición inline, confirmación y descarte al cambiar de proyecto.
- [X] T013 Convertir `index.tsx` en el shell de rutas (`index` → listado, `nueva` → creación) y agregar los estilos.
- [X] T014 COMMIT `GREEN: implementar el listado y la edición del Product Backlog, tests en verde`

## Fase 4: REFACTOR y verificación

- [X] T015 REFACTOR: validación pura en `edicionHistoria.ts` y reinicio del listado por `key={id}` en lugar de `setState` dentro del efecto.
- [X] T016 REFACTOR integrado en GREEN: no hubo un commit vacío; el reinicio por `key={id}` quedó en el commit de implementación.
- [X] T017 Verificar `npm test`, `npm run lint` y `npm run build` en `web/`.
- [X] T018 Revisar que las pruebas de HU-01 y de rutas sigan en verde.
- [ ] T019 Actualizar la rama con `origin/main` antes de abrir el PR.
- [ ] T020 Abrir el PR con `Closes #98`, pedir review y mergear con "Create a merge commit".
