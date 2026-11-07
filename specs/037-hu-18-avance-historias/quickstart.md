# Quickstart: Avance de las Historias del Proyecto (HU-18)

Guía para validar la historia de punta a punta. Contrato: [contracts/openapi.yaml](./contracts/openapi.yaml). Modelo: [data-model.md](./data-model.md).

## Tests automáticos

```bash
go test ./...                 # dominio, service (fakes), repositorio (SQLite :memory:), handler
go test ./features/...        # escenarios Godog, incluido features/avance_historias.feature (@HU-18)
gofmt -l . && go vet ./...    # lo mismo que el check `backend` del CI
cd web && npm run lint && npm test && npm run build   # el check `frontend`
```

Esperado: todo en verde y `gofmt -l .` sin salida.

## Prueba manual contra la API

Prerrequisito: el servidor corriendo (`go run ./cmd/api`) sobre una base nueva.

1. Crear un proyecto y consultar `GET /projects/<proyecto>/avance`.
   Esperado: `200`, `total` en cero y los tres grupos con `cantidad: 0` y `historias: []`.
2. Crear tres historias (una sin estimar, una de 5 y una de 3 Story Points). Consultar de nuevo.
   Esperado: las tres en `Nueva`, `sprint: null`, cantidad 3 y 8 Story Points.
3. Crear e iniciar un Sprint, asignarle la de 5, pasarla a En progreso y luego a Completada.
   Esperado: una historia en cada grupo con su cantidad y Story Points; la completada con `sprint` `{id, numero: 1}`.
4. Cerrar el Sprint con una historia no completada asignada.
   Esperado: la no completada vuelve a `Nueva` con `sprint: null`.
5. `GET /projects/abc/avance` → `400` con `campo: id`; `GET /projects/9999/avance` → `404`.

## Prueba manual en la interfaz

1. Abrir `/projects/<proyecto>/status`: se ven los tres grupos con cantidad y Story Points, y el Sprint de cada historia ("Sprint 1" o "Product Backlog").
2. Proyecto sin historias: mensaje de que no hay historias y totales en cero.
3. `/projects/abc/status`: mensaje de error sin pedir datos al servidor.
4. Con el servidor detenido: mensaje de error de conexión y ningún dato parcial.
