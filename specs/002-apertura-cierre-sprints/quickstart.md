# Quickstart: Apertura y Cierre de Sprints

## Prerrequisitos

- Go 1.26 instalado (`go build ./...` corre limpio en el repo).
- No hace falta que HU-04 esté implementada: se usa
  `repository.NewProjectRepositoryFake()` (en memoria) cargado con un
  `Project` de prueba, siguiendo `internal/domain/project.go` /
  `internal/repository/project_repository.go` ya commiteados.

## Validar el flujo completo (una vez implementado)

```bash
go test ./internal/... -run Sprint -v
```

### Escenario 1 — Iniciar un Sprint

1. Crear un `Project` de prueba en el fake de `ProjectRepository`.
2. `CrearSprint(ctx, proyectoID)` → Sprint en estado `Pendiente`.
3. `IniciarSprint(ctx, sprintID, "Entregar MVP", fechaInicio, fechaFin)` →
   Sprint pasa a `Activo`.
4. Verificar: un segundo `IniciarSprint` sobre otro Sprint del mismo
   proyecto devuelve error (regla de un solo Activo).

### Escenario 2 — Cerrar un Sprint con arrastre

1. Con el Sprint `Activo` del paso anterior, asignarle (vía
   `BacklogRepository`, campo `SprintID`) dos historias: una marcada
   "completada" y otra no.
2. `CerrarSprint(ctx, sprintID)` → Sprint pasa a `Finalizado`.
3. Verificar: la historia no completada tiene `SprintID == nil` (volvió al
   backlog); la completada conserva su `SprintID` original.

## Correr el escenario BDD (Godog)

```bash
go test ./features/... -v
```
(una vez generado `features/apertura_cierre_sprints.feature` y sus step
definitions, después de `/speckit.tasks` + `/speckit.analyze`).
