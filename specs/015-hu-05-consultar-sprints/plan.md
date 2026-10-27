# Plan: Consultar Sprints de un proyecto (Issue #110)

Sin cambios de esquema ni de `SprintRepository`. Sin `data-model.md`: no hay tablas nuevas.

## Cambios

- `internal/service/consultar_sprints.go`: `ConsultarSprints.Ejecutar(ctx, proyectoID)`.
  Dependencias mínimas (interfaces locales): `ObtenerPorID` de proyectos y `ListarPorProyecto`
  de Sprints. Valida el proyecto, ordena por `id`, devuelve `[]` si no hay.
- `internal/http/sprint_handler.go`: método `Listar` en `SprintHandler` (recibe el nuevo caso de uso).
- `cmd/api/main.go`: `GET /projects/{id}/sprints` junto a las rutas de Sprints.
- `features/consultar_sprints.feature` y `features/steps_consultar_sprints.go`.

## Decisiones

- Reutiliza `sprintResponse`/`aRespuestaSprint`, `idDeRuta` y `escribirError`.
- No agrega `numero` (#129).
