# Tasks: Numerar los Sprints por proyecto (Issue #129)

Cada tarea con lógica lleva su propio commit `RED:` y `GREEN:`. Todos los commits compilan.

## Tareas

- [x] T001 Escribir la especificación con la decisión del Scrum Master.
- [x] T002 COMMIT `docs: agregar la especificación de la numeración de Sprints por proyecto (Refs #129)`
- [x] T003 RED: tests del repositorio SQLite y del handler con Sprints de dos proyectos, y escenario BDD de la numeración.
- [x] T004 COMMIT `RED: agregar tests de la numeración de Sprints por proyecto (fallan, Refs #129)`
- [ ] T005 GREEN: derivar `Numero` en el adaptador SQLite y en el fake en memoria, y exponerlo como `numero` en la respuesta.
- [ ] T006 COMMIT `GREEN: numerar los Sprints por proyecto, tests en verde (Refs #129)`
- [ ] T007 Marcar las tareas y verificar `go build`, `go vet`, `gofmt -l` y `go test ./...`.
- [ ] T008 COMMIT `docs: marcar las tareas completadas de la numeración de Sprints (Refs #129)`
