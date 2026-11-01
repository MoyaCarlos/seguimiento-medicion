# Quickstart: validar la carga de datos de prueba

## Prerrequisitos
`go`, `node` y `npm` instalados; base local descartable.

## Backend
```sh
DEMO_HABILITADO=true DB_PATH=/tmp/demo.db go run ./cmd/api &
curl -i localhost:8080/demo/estado          # 200 {"habilitado":true}
curl -i -X POST localhost:8080/demo/datos   # 201 con el resumen
curl -i -X POST localhost:8080/demo/datos   # 409 ya cargado
curl -s localhost:8080/projects/25/backlog  # 10 historias
```
Sin la variable, ambos endpoints devuelven 404.

## Pruebas
```sh
go test ./...
cd web && npm run lint && npm test && npm run build
```

## Frontend
`DEMO_HABILITADO=true` + `npm run dev` en `web/`: el inicio muestra el botón; tras cargar, el enlace abre `/projects/25`.
