# Quickstart: validar `GET /backlog/{id}`

Contrato en [contracts/obtener-historia.md](contracts/obtener-historia.md).

## Prerrequisitos

- Go instalado y la API levantada en `http://localhost:8080` (`go run ./cmd/api`).
- Una historia creada (ver `docs/curl-frontend-sprint-2.md`, "Crear una historia"). Los ejemplos usan el id `1`.

## Pruebas automáticas

```bash
go test ./internal/http/... -run Obtener
go test ./features/...
go test ./...
```

Resultado esperado: todo en verde; las suites existentes no cambian.

## Validación manual (bash)

```bash
curl -i http://localhost:8080/backlog/1      # 200, historia completa
curl -i http://localhost:8080/backlog/999    # 404 {"mensaje":"historia no encontrada"}
curl -i http://localhost:8080/backlog/abc    # 400 {"campo":"id","mensaje":"identificador inválido"}
curl -i http://localhost:8080/backlog/0      # 400
curl -i http://localhost:8080/backlog/007    # igual que /backlog/7
```

## Validación manual (cmd)

```bat
curl -i http://localhost:8080/backlog/1
curl -i http://localhost:8080/backlog/999
curl -i http://localhost:8080/backlog/abc
```

## Resultados esperados

| Caso | Estado | Cuerpo |
|------|--------|--------|
| Historia sin Sprint | 200 | `sprint_id: null`, `fecha_finalizacion: null` |
| Historia en un Sprint | 200 | `sprint_id` con el id del Sprint |
| Historia completada | 200 | `fecha_finalizacion` con la fecha |
| Inexistente | 404 | `historia no encontrada` |
| Id inválido | 400 | campo `id`, `identificador inválido` |
