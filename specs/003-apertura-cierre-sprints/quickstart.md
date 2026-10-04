# Quickstart: Apertura y Cierre de Sprints

## Prerrequisitos

- Go 1.26 instalado (`go build ./...` corre limpio en el repo).
- HU-04 ya está en `main`: en los tests se usa
  `repository.NewProjectRepositoryEnMemoria()` con un `Project` de prueba
  guardado vía `Guardar`.
- Si tenés un `seguimiento.db` local de antes de HU-04, borralo: la API
  detecta el esquema viejo y termina con error.

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
(escenarios en `features/apertura_cierre_sprints.feature`, pasos en
`features/steps_sprints.go`).
