# Quickstart: validar el Scrum Master único (#236)

Prerrequisitos: Go 1.27, el repositorio en la rama `feature/236-scrum-master-unico`.

## 1. Pruebas automáticas

```sh
go test ./...        # dominio, service, repositorio (esquema + migración), handler y Godog
gofmt -l . && go vet ./...
```

Esperado: todo en verde. Deben fallar antes de implementar (RED): el test de service que rechaza el segundo Scrum Master, el de repositorio que traduce el código 2067, el de migración y el de esquema del índice.

## 2. Prueba manual con la API

```sh
go run ./cmd/api &
curl -s -XPOST localhost:8080/projects -d '{"nombre":"P","creador":"Ana"}'      # 201, Ana es Scrum Master
curl -s -XPOST localhost:8080/projects/1/members -d '{"nombre":"Beto","rol":"product_builder"}'  # 201
curl -s -i -XPOST localhost:8080/projects/1/members -d '{"nombre":"Cami","rol":"scrum_master"}'  # 409 + mensaje
curl -s -i -XPOST localhost:8080/projects/1/members -d '{"nombre":"Ana","rol":"scrum_master"}'   # 400 integrante duplicado
```

Contrato: [contracts/members-409.yaml](./contracts/members-409.yaml). Puerto y rutas según `cmd/api/main.go`.

## 3. Migración sobre una base con duplicados

Cubierta por el test de migración (ver [research.md](./research.md) §4): migra hasta la 005, siembra proyectos con 0, 1, 2 y 3 Scrum Master y un empate de `creado_en`, aplica la 007 y comprueba el resultado.
