# Quickstart: Imputación de Horas Trabajadas (HU-07)

Guía para validar la historia de punta a punta. Contrato: [contracts/openapi.yaml](./contracts/openapi.yaml).
Modelo: [data-model.md](./data-model.md).

## Tests automáticos

```bash
go test ./...                 # dominio, services (fakes), repositorio (SQLite :memory:), handlers
go test ./features/...        # escenarios Godog, incluido features/imputacion_horas.feature (@HU-07)
gofmt -l . && go vet ./...    # lo mismo que el check `backend` del CI
```

Esperado: todo en verde y `gofmt -l .` sin salida.

## Prueba manual contra la API

Prerrequisito: el servidor corriendo (`go run ./cmd/api`) sobre una base nueva (la migración
004 se aplica sola al arrancar).

1. Crear un proyecto con su creadora (`POST /projects`) y anotar el `id` del proyecto y el
   del integrante (`GET /projects/{id}/members`).
2. Crear una historia (`POST /backlog`) y un Sprint (`POST /sprints`), e iniciarlo
   (`POST /sprints/{id}/iniciar`) con una fecha de inicio anterior o igual a hoy.
3. Asignar la historia al Sprint (HU-06):
   ```bash
   curl -s -X PUT localhost:8080/backlog/<historia>/sprint \n     -H 'Content-Type: application/json' \n     -d '{"sprint_id": <sprint>}'
   ```
4. Imputar 4 horas:
   ```bash
   curl -s -X POST localhost:8080/backlog/<historia>/horas \
     -H 'Content-Type: application/json' \
     -d '{"integrante_id": <usuario>, "fecha": "<hoy>", "actividad": "Implementar", "horas": 4}'
   ```
   Esperado: `201`, con `minutos: 240`, `horas: 4` y el `sprint_id` del Sprint activo.
5. Consultar el esfuerzo: `curl -s localhost:8080/backlog/<historia>/horas`.
   Esperado: `total_minutos: 240` y una carga en `registros`.
6. Casos de error (cada uno `400` con el `campo` indicado, y el total no cambia):

   | Cuerpo | `campo` |
   |---|---|
   | `"horas": -2` | `horas` |
   | `"horas": "abc"` | `horas` |
   | `"horas": 0.01` | `horas` |
   | `"horas": 25` | `horas` |
   | `"actividad": "   "` | `actividad` |
   | `"fecha": "<mañana>"` | `fecha` |
   | `"integrante_id": 999` | `integrante_id` |

7. Cerrar el Sprint (`POST /sprints/{id}/cerrar`) e imputar otra vez: `400` con
   `campo: historia_id`. Consultar el esfuerzo: las cargas anteriores siguen con el
   `sprint_id` del Sprint cerrado.
8. `POST /backlog/9999/horas` y `GET /backlog/9999/horas`: `404` "historia no encontrada".

## Resultado de la validación (2026-10-28)

- **Tests automáticos**: `go build`, `go vet`, `gofmt -l .` (sin salida) y `go test ./...` en verde,
  incluidos los 12 escenarios `@HU-07` (16 casos con las filas de los esquemas; 64 casos BDD en
  total).
- **API manual**: los pasos 1–8 dieron lo esperado contra un servidor con base nueva. El paso 3 se
  hizo entonces con un `UPDATE` directo en la base, porque HU-06 todavía no estaba en `main`. Paso 6: los siete casos dieron `400`
  con el `campo` indicado y `total_minutos` siguió en 240. Paso 7: después del cierre, `400` con
  `campo: historia_id`, y la carga anterior sigue con `sprint_id: 1`.
- Para probar con `curl` en Windows, mandar los cuerpos con acentos desde un archivo
  (`--data-binary @cuerpo.json`): pasados como argumento, la consola los recodifica y llegan al
  servidor como texto inválido.
