# Quickstart: Votación de Esfuerzo en Planning Poker (HU-02)

Guía para validar la funcionalidad de punta a punta. Contrato: [contracts/openapi.yaml](./contracts/openapi.yaml).
Reglas y errores: [data-model.md](./data-model.md) y [research.md R8](./research.md#r8-errores-y-códigos-http).

## Prerrequisitos

- Go 1.26+ y la rama al día con `main` (migración 001).
- La migración 003 se aplica sola al arrancar; no hace falta borrar `seguimiento.db`.

## 1. Tests automáticos

```bash
go test ./internal/...                     # dominio, service, repositorio (incluye test de esquema), handler
go test ./features/... -run TestFeatures   # escenarios BDD @HU-02 (y el resto)
gofmt -l .                                 # no debe listar nada
```

Esperado: todo en verde.

## 2. Validación manual con la API

```bash
go run ./cmd/api
```

En otra terminal, un proyecto con Ana (Scrum Master, creadora) y dos Product Builders, y
una historia:

```bash
curl -s -X POST localhost:8080/projects -d '{"nombre":"Demo","descripcion":"","creador":"Ana","fecha_inicio":null,"fecha_fin":null}'
curl -s -X POST localhost:8080/projects/1/members -d '{"nombre":"Bruno","rol":"product_builder"}'
curl -s -X POST localhost:8080/projects/1/members -d '{"nombre":"Carla","rol":"product_builder"}'
curl -s localhost:8080/projects/1/members        # anotar los id: Ana=1, Bruno=2, Carla=3
curl -s -X POST localhost:8080/backlog -d '{"proyecto_id":1,"titulo":"Login","descripcion":"d","prioridad":"M"}'
```

| # | Petición | Esperado |
|---|----------|----------|
| 1 | `POST /backlog/1/votacion` con `{"usuario_id":2}` | 403: Bruno no es Scrum Master |
| 2 | `POST /backlog/1/votacion` con `{"usuario_id":1}` | 201, `ronda: 1`, `revelada_en: null`, los 3 con `ya_voto: false` |
| 3 | Repetir el paso 2 | 409: la historia ya tiene votación |
| 4 | `POST /backlog/1/votacion/votos` con `{"usuario_id":2,"valor":4}` | 400 `campo: valor` |
| 5 | `POST /backlog/1/votacion/votos` con `{"usuario_id":99,"valor":8}` | 403: no es integrante |
| 6 | `POST /backlog/1/votacion/votos` con `{"usuario_id":2,"valor":8}` | 200; Bruno ve `valor: 8`, Ana y Carla `valor: null` |
| 7 | `GET /backlog/1/votacion?usuario_id=1` | Bruno `ya_voto: true, valor: null` (Oculto); Ana sin voto |
| 8 | Ana vota 3 y Carla vota 13, luego Carla vota 5 | 200; Carla ve solo su 5 (reemplazó el 13) |
| 9 | `POST /backlog/1/votacion/revelar` con `{"usuario_id":1}` antes del paso 8 | 400 `campo: votos` (falta votar) |
| 10 | `POST /backlog/1/votacion/revelar` con `{"usuario_id":1}` | 200; valores 3, 8, 5 y `promedio: {exacto: 5.33, carta_sugerida: 5}` |
| 11 | `POST /backlog/1/votacion/votos` con `{"usuario_id":2,"valor":13}` | 400 `campo: estado` (ronda revelada) |
| 12 | `POST /backlog/99/votacion` con cuerpo `{}` | 404 `historia no encontrada` (no 400) |
| 13 | `GET /projects/1/backlog` | `estimacion_sp` sin cambios (FR-013) |
| 14 | Reiniciar `go run ./cmd/api` y repetir el paso 7 | La ronda revelada persiste |

> El caso "historia Completada" (400 `campo: estado` al iniciar) se valida con el test de
> service y el escenario BDD: hoy no hay endpoint para completar una historia (HU-11).
