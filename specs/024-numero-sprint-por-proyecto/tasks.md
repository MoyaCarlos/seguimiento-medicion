# Tasks: Numerar los Sprints por proyecto (Issue #129)

Cada tarea con lógica lleva su propio commit `RED:` y `GREEN:`. Todos los commits compilan.

## Tareas

- [x] T001 Escribir la especificación con la decisión del Scrum Master.
- [x] T002 COMMIT `docs: agregar la especificación de la numeración de Sprints por proyecto (Refs #129)`
- [x] T003 RED: tests del repositorio SQLite y del handler con Sprints de dos proyectos, y escenario BDD de la numeración.
- [x] T004 COMMIT `RED: agregar tests de la numeración de Sprints por proyecto (fallan, Refs #129)`
- [x] T005 GREEN: derivar `Numero` en el adaptador SQLite y en el fake en memoria, y exponerlo como `numero` en la respuesta.
- [x] T006 COMMIT `GREEN: numerar los Sprints por proyecto, tests en verde (Refs #129)`
- [x] T007 Marcar las tareas y verificar `go build`, `go vet`, `gofmt -l` y `go test ./...`.
- [x] T008 COMMIT `docs: marcar las tareas completadas de la numeración de Sprints (Refs #129)`

## Frontend (fix/129-numero-sprint-frontend)

- [x] T009 RED: tests de componente con `numero` distinto del `id` en las tarjetas de Sprint, el selector de asignación y el aviso de asignación.
- [x] T010 COMMIT `RED: agregar tests del número de Sprint en la pantalla (fallan, Refs #129)`
- [x] T011 GREEN: mostrar `numero` en `SeccionSprints`, `AsignarASprint` y `Listado`, y agregarlo al contrato de `specs/003-apertura-cierre-sprints`.
- [x] T012 COMMIT `GREEN: mostrar el número del Sprint por proyecto, tests en verde (Refs #129)`
- [x] T013 Verificar `npm run lint`, `npm test` y `npm run build`.
