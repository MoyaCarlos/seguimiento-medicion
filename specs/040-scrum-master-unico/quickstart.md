# Quickstart: validar el Scrum Master único (#236)

Prerrequisitos: Go 1.26.5 o posterior (el `go.mod` declara 1.26.5), el repositorio en la rama `feature/236-scrum-master-unico`.

## 1. Pruebas automáticas

```sh
go test ./...        # dominio, service, repositorio (esquema + migración), handler y Godog
gofmt -l . && go vet ./...
```

Esperado: todo en verde. Deben fallar antes de implementar (RED): el test de service que rechaza el segundo Scrum Master, el de repositorio que traduce el código 2067, el de migración y el de esquema del índice.

## 2. Prueba manual con la API

Con una base temporal y un puerto libre (el puerto sale de `PORT`, por defecto 8080, y la base de `DB_PATH`; ver `cmd/api/config.go`):

```sh
DB_PATH=/tmp/qs.db PORT=9090 go run ./cmd/api &
B=localhost:9090
curl -s -i -XPOST $B/projects -d '{"nombre":"P","creador":"Ana"}'                                # 201, Ana es Scrum Master
curl -s -i -XPOST $B/projects/1/members -d '{"nombre":"Beto","rol":"product_builder"}'           # 201
curl -s -i -XPOST $B/projects/1/members -d '{"nombre":"Cami","rol":"scrum_master"}'              # 409
curl -s -i -XPOST $B/projects/1/members -d '{"nombre":"Ana","rol":"scrum_master"}'               # 400 integrante
curl -s $B/projects/1/members                                                                    # Ana (scrum_master) y Beto (product_builder)
kill %1; rm -f /tmp/qs.db*
```

Cuerpos esperados:

- 409: `{"mensaje":"el proyecto ya tiene un Scrum Master"}` (sin `campo`).
- 400: `{"campo":"integrante","mensaje":"el integrante ya pertenece al proyecto"}`; es el mismo para quien ya es Scrum Master y para un Product Builder repetido.

Contrato: [contracts/members-409.yaml](./contracts/members-409.yaml).

## 3. Migración sobre una base con duplicados

Cubierta por el test de migración (ver [research.md](./research.md) §4): migra hasta la 005, siembra proyectos con 0, 1, 2 y 3 Scrum Master y un empate de `creado_en`, aplica la 007 y comprueba el resultado.
