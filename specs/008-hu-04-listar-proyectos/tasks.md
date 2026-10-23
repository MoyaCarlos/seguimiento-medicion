# Tasks: Listado de proyectos (Issue #109)

Cada tarea con lógica lleva su propio commit `RED:` y `GREEN:`. Todos los commits deben
compilar.

## Tareas

- [ ] T001 RED: tests del repositorio SQLite y del fake en memoria para listar proyectos ordenados y sin proyectos.
- [ ] T002 COMMIT `RED: agregar tests de listado de proyectos en SQLite y fake en memoria (fallan)`
- [ ] T003 GREEN: agregar `Listar` al puerto `ProjectRepository`, al adaptador SQLite y al fake en memoria; actualizar stubs en el mismo commit.
- [ ] T004 COMMIT `GREEN: implementar Listar en el repositorio de proyectos, tests en verde`
- [ ] T005 RED: tests del servicio `ListarProyectos` con proyectos desordenados, sin proyectos y con error de persistencia.
- [ ] T006 COMMIT `RED: agregar tests de ListarProyectos (fallan)`
- [ ] T007 GREEN: implementar `ListarProyectos`, ordenando por `id` y devolviendo un arreglo no nulo.
- [ ] T008 COMMIT `GREEN: implementar el servicio ListarProyectos, tests en verde`
- [ ] T009 RED: tests del handler `GET /projects` con lista y con lista vacía.
- [ ] T010 COMMIT `RED: agregar tests del handler GET /projects (fallan)`
- [ ] T011 GREEN: implementar el handler y registrar la ruta `GET /projects` en `cmd/api`.
- [ ] T012 COMMIT `GREEN: implementar GET /projects y cablear la ruta, tests en verde`
- [ ] T013 Actualizar el contrato OpenAPI y la trazabilidad de la spec.
- [ ] T014 COMMIT `docs: actualizar el contrato de GET /projects`
