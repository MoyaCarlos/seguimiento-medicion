# Tasks: Consultar Sprints de un proyecto (Issue #110)

Cada tarea con lógica lleva su propio commit `RED:` y `GREEN:`. Todos los commits compilan.

- [x] T001 Redactar spec, plan y tareas (specify, clarify, plan, tasks, analyze).
- [x] T002 COMMIT `docs: agregar la especificación de consultar Sprints de un proyecto (Refs #110)`
- [x] T003 RED: escenarios BDD en `features/consultar_sprints.feature` sin pasos; fallan con el runner estricto.
- [x] T004 COMMIT `RED: agregar escenarios BDD de consultar Sprints (fallan)`
- [x] T005 RED: tests de `ConsultarSprints` (varios Sprints en desorden, ninguno, proyecto inexistente).
- [x] T006 COMMIT `RED: agregar tests de ConsultarSprints (fallan)`
- [x] T007 GREEN: implementar `ConsultarSprints`.
- [x] T008 COMMIT `GREEN: implementar el servicio ConsultarSprints, tests en verde`
- [x] T009 RED: tests del handler `GET /projects/{id}/sprints` (200, `[]`, 404, 400).
- [x] T010 COMMIT `RED: agregar tests del handler GET /projects/{id}/sprints (fallan)`
- [x] T011 GREEN: implementar `SprintHandler.Listar` y registrar la ruta en `cmd/api`.
- [x] T012 COMMIT `GREEN: implementar GET /projects/{id}/sprints y cablear la ruta, tests en verde`
- [x] T013 GREEN: implementar los pasos BDD, escenarios en verde.
- [x] T014 COMMIT `GREEN: implementar los pasos BDD de consultar Sprints, escenarios en verde`
- [x] T015 Verificar `go build`, `go vet`, `gofmt -l`, `go test ./...` y la ruta con curl; marcar tareas.
- [x] T016 COMMIT `docs: marcar las tareas completadas de consultar Sprints (Refs #110)`
- [x] T017 RED: tests de la pantalla de Sprints (carga y agrupación, lista vacía, error de la API, cambio de proyecto) en `web/src/features/sprints/index.test.tsx`; se quita el test de `cargarSprintsIniciales` y del texto de limitación.
- [x] T018 COMMIT `RED: agregar tests de la carga inicial de Sprints en la pantalla (fallan)`
- [x] T019 GREEN: `listarSprints` en `api.ts`, `useEffect` con indicador de carga y aviso de error en `index.tsx`; se elimina `cargarSprintsIniciales` y el texto de limitación.
- [x] T020 COMMIT `GREEN: cargar los Sprints existentes al abrir la pantalla, tests en verde`
- [x] T021 REFACTOR: nada que limpiar, sin commit.
- [x] T022 Actualizar las specs 011 y 015 con la parte de frontend.
- [x] T023 COMMIT `docs: ampliar la spec 015 con el frontend y actualizar la 011`
