# Tasks: Listado de proyectos (Issue #109)

Cada tarea con lógica lleva su propio commit `RED:` y `GREEN:`. Todos los commits compilan.

## Tareas

- [x] T001 RED: tests del repositorio SQLite y del fake en memoria para listar proyectos ordenados y sin proyectos.
- [x] T002 COMMIT `RED: agregar tests de listado de proyectos en SQLite y fake en memoria (fallan)`
- [x] T003 GREEN: agregar `Listar` al puerto `ProjectRepository`, al adaptador SQLite y al fake en memoria; actualizar stubs en el mismo commit.
- [x] T004 COMMIT `GREEN: implementar Listar en el repositorio de proyectos, tests en verde`
- [x] T005 RED: tests del servicio `ListarProyectos` con proyectos desordenados, sin proyectos y con error de persistencia.
- [x] T006 COMMIT `RED: agregar tests de ListarProyectos (fallan)`
- [x] T007 GREEN: implementar `ListarProyectos`, ordenando por `id` y devolviendo un arreglo no nulo.
- [x] T008 COMMIT `GREEN: implementar el servicio ListarProyectos, tests en verde`
- [x] T009 RED: tests del handler `GET /projects` con lista y con lista vacía.
- [x] T010 COMMIT `RED: agregar tests del handler GET /projects (fallan)`
- [x] T011 GREEN: implementar el handler y registrar la ruta `GET /projects` en `cmd/api`.
- [x] T012 COMMIT `GREEN: implementar GET /projects y cablear la ruta, tests en verde`
- [x] T013 Actualizar el contrato OpenAPI de HU-04 con `GET /projects`.
- [x] T014 COMMIT `docs: agregar la especificación y el contrato de listado de proyectos (Issue #109)`
- [x] T015 Marcar las tareas y verificar `go build`, `go vet`, `gofmt -l` y `go test ./...`.
- [x] T016 COMMIT `docs: marcar las tareas completadas de listar proyectos (Refs #109)`
- [x] T017 Habilitar `Strict` en el runner de Godog para que los pasos indefinidos fallen.
- [x] T018 COMMIT `test: hacer estricto el runner BDD para que los pasos indefinidos fallen`
- [x] T019 RED: escenarios BDD del listado (orden por id, lista vacía y error 500) sin pasos; fallan con el runner estricto.
- [x] T020 COMMIT `RED: agregar escenarios BDD del listado de proyectos (fallan)`
- [x] T021 GREEN: implementar los pasos BDD sobre SQLite y el handler, escenarios en verde.
- [x] T022 COMMIT `GREEN: implementar los pasos BDD del listado, escenarios en verde`
