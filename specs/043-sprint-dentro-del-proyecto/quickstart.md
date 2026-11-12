# Quickstart: validar Sprints dentro del plazo del proyecto (#168)

Guía para comprobar la historia de punta a punta. Contrato: [contracts/iniciar-sprint-400.yaml](contracts/iniciar-sprint-400.yaml). Regla: [data-model.md](data-model.md).

## 1. Pruebas automáticas

Desde la raíz del repo:

```sh
go build ./... && go vet ./... && gofmt -l .
go test ./internal/domain/ -run AdmiteSprint -v
go test ./internal/service/ -run IniciarSprint -v
go test ./features/ -v            # Godog: incluye features/sprint_dentro_del_proyecto.feature
go test ./...
```

Esperado: todo en verde, `gofmt -l .` sin salida.

## 2. Prueba manual contra la API

Levantar el backend (`go run ./cmd/api`) y, con un proyecto del 2026-10-31 al 2026-11-16 y un Sprint Pendiente creado (`POST /projects/{id}/sprints`), iniciar el Sprint con `POST /sprints/{id}/iniciar`:

| Cuerpo (`sprint_goal` = "Goal") | Esperado |
|---|---|
| `fecha_inicio` 2026-10-30, `fecha_fin` 2026-11-06 | 400, `campo: fecha_inicio`; el caso del beta testing de Render |
| `fecha_inicio` 2026-11-10, `fecha_fin` 2026-11-17 | 400, `campo: fecha_fin` |
| `fecha_inicio` 2026-10-31, `fecha_fin` 2026-11-16 | 200, Sprint Activo (los dos bordes) |

Después de un 400, `GET /projects/{id}/sprints` muestra el Sprint todavía Pendiente.

Con un proyecto sin fechas, cualquier par válido (fin posterior al inicio) da 200.

## 3. Pantalla

En la pantalla de Sprints, iniciar un Sprint con una fecha de inicio anterior al proyecto: el mensaje aparece debajo del campo "Fecha de inicio" (no requiere cambios de frontend, ver [research.md §8](research.md)).
