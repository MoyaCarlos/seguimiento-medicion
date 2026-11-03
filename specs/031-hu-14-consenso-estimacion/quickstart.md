# Quickstart: Consenso y Registro de la Estimación Acordada (HU-14)

Guía para validar la funcionalidad de punta a punta. Contrato: [contracts/openapi.yaml](./contracts/openapi.yaml).
Reglas y errores: [data-model.md](./data-model.md) y [research.md R10](./research.md#r10-errores-y-códigos-http).

## Prerrequisitos

- Go 1.26+ y la rama al día con `main` (HU-02 ya mergeada).
- No hay migración nueva: no hace falta borrar `seguimiento.db`.

## 1. Tests automáticos

```bash
go test ./internal/...                     # dominio, service, repositorio, handler
go test ./features/... -run TestFeatures   # escenarios BDD @HU-14 (y el resto)
gofmt -l .                                 # no debe listar nada
```

Esperado: todo en verde, también los escenarios `@HU-02` y los tests de edición de HU-15.

## 2. Validación manual con la API

```bash
go run ./cmd/api
```

En otra terminal, un proyecto con Carlos (Scrum Master) y dos Product Builders, una historia y
su votación con la ronda 1 abierta:

```bash
curl -s -X POST localhost:8080/projects -d '{"nombre":"Demo","descripcion":"","creador":"Carlos","fecha_inicio":null,"fecha_fin":null}'
curl -s -X POST localhost:8080/projects/1/members -d '{"nombre":"Candela","rol":"product_builder"}'
curl -s -X POST localhost:8080/projects/1/members -d '{"nombre":"Jimena","rol":"product_builder"}'
curl -s localhost:8080/projects/1/members        # anotar los id: Carlos=1, Candela=2, Jimena=3
curl -s -X POST localhost:8080/backlog -d '{"proyecto_id":1,"titulo":"Login","descripcion":"d","prioridad":"M"}'
curl -s -X POST localhost:8080/backlog/1/votacion -d '{"usuario_id":1}'
```

| # | Petición | Esperado |
|---|----------|----------|
| 1 | `GET /backlog/1/votacion` | `diferencias: null` (ronda sin revelar) |
| 2 | `POST /backlog/1/votacion/rondas` con `{"usuario_id":1}` | 400 `campo: estado` (ronda sin revelar) |
| 3 | `POST /backlog/1/votacion/estimacion` con `{"usuario_id":1,"estimacion_sp":5}` | 400 `campo: estado`; `estimacion_sp` sigue sin cargar |
| 4 | Votan 3 (Carlos), 8 (Candela) y 5 (Jimena); Carlos revela | `diferencias: {hay_diferencias: true, minimo: 3, maximo: 8}` |
| 5 | `POST /backlog/1/votacion/estimacion` con `{"usuario_id":2,"estimacion_sp":5}` | 403 (Candela no es Scrum Master) |
| 6 | `POST /backlog/1/votacion/estimacion` con `{"usuario_id":1,"estimacion_sp":4}` | 400 `campo: estimacion_sp` |
| 7 | `POST /backlog/99/votacion/estimacion` con cuerpo `{}` | 404 `historia no encontrada` (no 400) |
| 8 | `POST /backlog/1/votacion/rondas` con `{"usuario_id":2}` | 403 |
| 9 | `POST /backlog/1/votacion/rondas` con `{"usuario_id":1}` | 201, `ronda: 2`, nadie votó, `diferencias: null` |
| 10 | `GET /backlog/1/votacion/rondas` | solo la ronda 1, con los votos 3, 8 y 5 |
| 11 | Los tres votan 5 y Carlos revela | `diferencias: {hay_diferencias: false, minimo: null, maximo: null}` |
| 12 | `GET /backlog/1/votacion/rondas` | rondas 1 y 2, cada una con sus votos |
| 13 | `POST /backlog/1/votacion/estimacion` con `{"usuario_id":1,"estimacion_sp":13}` | 200, `estimacion_sp: 13`, aunque nadie votó 13 |
| 14 | `POST /backlog/1/votacion/rondas` con `{"usuario_id":1}` | 201, `ronda: 3`: se puede abrir sin diferencias y después de registrar |
| 15 | `GET /projects/1/backlog` | la historia figura con `estimacion_sp: 13` |
| 16 | Reiniciar `go run ./cmd/api` y repetir el paso 12 | el historial persiste |

> Los casos con Sprint `Activo` o `Finalizado` (D2b), historia `Completada` y dos aperturas simultáneas
> se validan con los tests de service y repositorio y con los escenarios BDD.

## Resultado de la validación (2026-11-03)

- **Tests automáticos**: `go build`, `go vet`, `gofmt -l .` (sin salida) y `go test ./...` en verde,
  incluidos los 20 escenarios `@HU-14` (152 escenarios BDD en total, contando cada fila de un esquema).
- **API manual**: los pasos 1 a 16 dieron lo esperado, contra un servidor con base nueva. Paso 4: `diferencias`
  con mínimo 3 y máximo 8. Paso 11: `hay_diferencias: false`, con mínimo y máximo en null. Paso 13: la historia
  queda con 13 Story Points, aunque nadie votó esa carta. Paso 14: se abre la ronda 3 después de registrar la
  estimación. Paso 16: tras reiniciar el servidor, el historial conserva las rondas 1 y 2.
