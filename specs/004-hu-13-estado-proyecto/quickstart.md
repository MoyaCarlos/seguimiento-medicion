# Quickstart: Fechas y Estado del Proyecto (HU-13)

## Prerrequisitos

- Go 1.26 instalado (`go build ./...` corre limpio en el repo).
- HU-04 y HU-05 ya están en `main`: se usan los fakes versionados
  `repository.NewProjectRepositoryEnMemoria()` y `repository.NewSprintRepositoryEnMemoria()`.
- No hay cambios de esquema en esta historia: no hace falta recrear `seguimiento.db`.

## Validar el flujo completo (una vez implementado)

```bash
go test ./internal/domain/... ./internal/service/... -run Estado -v
```

### Escenario 1 — Consulta de estado derivado

1. Guardar un `Project` de prueba en el fake de `ProjectRepository` (con `FechaInicio`/`FechaFin`).
2. Guardar Sprints del proyecto en el fake de `SprintRepository` en el estado deseado.
3. `ObtenerEstadoProyecto.Ejecutar(ctx, proyectoID, ahoraFijo)` → `domain.EstadoProyecto`.
4. Verificar contra la tabla de decisión (precedencia Sprints > fechas):

| Sprints del proyecto | Fecha actual vs proyecto | Estado esperado |
|---|---|---|
| Al menos uno `Activo` | cualquiera (aunque `FechaFin` vencida) | `En curso` |
| Ninguno `Activo`, todos `Finalizado` | cualquiera | `Finalizado` |
| Sin iniciar (ninguno Activo/Finalizado) | `ahora < FechaInicio` | `Planificado` |
| Sin iniciar | `FechaInicio ≤ ahora ≤ FechaFin` | `En curso` |
| Sin iniciar | `ahora > FechaFin` | `Finalizado` |
| Sin Sprints y sin `FechaInicio` | — | `Planificado` |

5. Verificar 404: `Ejecutar(ctx, idInexistente, ahora)` devuelve `domain.ErrProyectoNoEncontrado`.

### Escenario 2 — Endpoint HTTP (solo lectura)

```bash
go test ./internal/http/... -run Estado -v
```

- `GET /projects/{id}/status` → `200` con `{ "estado": "En curso" }`.
- `GET /projects/{id}/status` con `id` inexistente → `404` "proyecto no encontrado".
- `GET /projects/{id}/status` con `id` no numérico → `400`.
- Comprobar que la consulta no modifica filas: el `Project` y los `Sprint` quedan iguales.

## Correr el escenario BDD (Godog)

```bash
go test ./features/... -v
```

(escenario en `features/estado_proyecto.feature`; pasos reutilizando/agregando sobre los
helpers de `features/steps_sprints.go`).
