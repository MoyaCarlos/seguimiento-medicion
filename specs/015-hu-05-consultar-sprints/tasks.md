# Tasks: Consultar Sprints de un proyecto (Issue #110)

Cada tarea con lógica lleva su propio commit `RED:` y `GREEN:`. Todos los commits compilan.

- [x] T001 Redactar spec, plan y tareas (specify, clarify, plan, tasks, analyze).
- [x] T002 COMMIT `docs: agregar la especificación de consultar Sprints de un proyecto (Refs #110)`
- [ ] T003 RED: escenarios BDD en `features/consultar_sprints.feature` sin pasos; fallan con el runner estricto.
- [ ] T004 COMMIT `RED: agregar escenarios BDD de consultar Sprints (fallan)`
- [ ] T005 RED: tests de `ConsultarSprints` (varios Sprints en desorden, ninguno, proyecto inexistente).
- [ ] T006 COMMIT `RED: agregar tests de ConsultarSprints (fallan)`
- [ ] T007 GREEN: implementar `ConsultarSprints`.
- [ ] T008 COMMIT `GREEN: implementar el servicio ConsultarSprints, tests en verde`
- [ ] T009 RED: tests del handler `GET /projects/{id}/sprints` (200, `[]`, 404, 400).
- [ ] T010 COMMIT `RED: agregar tests del handler GET /projects/{id}/sprints (fallan)`
- [ ] T011 GREEN: implementar `SprintHandler.Listar` y registrar la ruta en `cmd/api`.
- [ ] T012 COMMIT `GREEN: implementar GET /projects/{id}/sprints y cablear la ruta, tests en verde`
- [ ] T013 GREEN: implementar los pasos BDD, escenarios en verde.
- [ ] T014 COMMIT `GREEN: implementar los pasos BDD de consultar Sprints, escenarios en verde`
- [ ] T015 Verificar `go build`, `go vet`, `gofmt -l`, `go test ./...` y la ruta con curl; marcar tareas.
- [ ] T016 COMMIT `docs: marcar las tareas completadas de consultar Sprints (Refs #110)`
