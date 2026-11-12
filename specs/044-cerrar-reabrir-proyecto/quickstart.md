# Quickstart: validar «Cerrar y reabrir un proyecto» (spec 044)

## Pruebas automáticas

Desde la raíz del repo:

```bash
go build ./... && go vet ./... && gofmt -l .
go test ./internal/domain/ -run 'Cerrar|Reabrir|EstadoProyecto'
go test ./internal/repository/ -run 'Esquema|Migracion008|Cerrar|Reabrir'
go test ./internal/service/ -run 'CerrarProyecto|ReabrirProyecto|Cerrado'
go test ./internal/app/ -run 'ProyectoCerrado'
go test ./features/ -godog.tags=@300
go test ./...
```

Esperado: todo en verde, `gofmt -l .` sin salida.

## Prueba manual con la API

Levantar el backend (`go run ./cmd/api`), registrarse e ingresar como `ana` (HU-20) y crear
un proyecto: `ana` queda como Scrum Master ([contrato](contracts/cerrar-reabrir-proyecto.yaml)).
`$T` es el token de `ana`; `$P`, el id del proyecto.

| Paso | Pedido | Esperado |
|---|---|---|
| 1 | `POST /projects/$P/cerrar` | 200, `estado: "Finalizado"`, `fecha_cierre` = hoy |
| 2 | `GET /projects/$P/status` | `{"estado":"Finalizado"}` |
| 3 | `POST /backlog` con `proyecto_id: $P` | 409 «El proyecto está cerrado» |
| 4 | `PUT /projects/$P` | 409 «El proyecto está cerrado» |
| 5 | `GET /projects/$P/backlog`, `/metricas`, `/reporte` | 200, como siempre |
| 6 | `POST /projects/$P/cerrar` | 409 «El proyecto ya está cerrado» |
| 7 | `POST /projects/$P/reabrir` | 200, `fecha_cierre: null`, estado según HU-13 |
| 8 | `POST /backlog` con `proyecto_id: $P` | 201 |
| 9 | `POST /projects/$P/reabrir` | 409 «El proyecto no está cerrado» |
| 10 | Crear e iniciar un Sprint; `POST /projects/$P/cerrar` | 409 «Cerrá el Sprint activo antes de cerrar el proyecto» |
| 11 | Con el token de un Product Builder: `POST /projects/$P/cerrar` | 403 «Solo el Scrum Master puede realizar esta acción» |

```bash
curl -s -X POST -H "Authorization: Bearer $T" http://localhost:8080/projects/$P/cerrar
```

## Migración

Con una base `seguimiento.db` previa a esta rama: al arrancar se aplica la 008 y los
proyectos existentes siguen abiertos (`fecha_cierre` en `null` en `GET /projects`).
