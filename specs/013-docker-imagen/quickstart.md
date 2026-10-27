# Quickstart: verificar la imagen

Prerequisitos: Docker. Desde la raíz del repo:

```sh
docker build -t seguimiento-medicion .
docker volume create seguimiento-data
docker run -d --name seg -p 8080:8080 -v seguimiento-data:/data seguimiento-medicion

curl -i localhost:8080/health                                   # 200 ok
curl -s -H 'Accept: text/html' localhost:8080/projects/1/status | head -c 100   # <!doctype html>
docker inspect --format '{{.Config.User}}' seguimiento-medicion  # nonroot:nonroot (no root)
curl -s -H 'Accept: text/html' localhost:8080/ | head -c 100                        # <!doctype html>
curl -i -X POST localhost:8080/projects -H 'Content-Type: application/json' \
  -d '{"nombre":"Demo","creador":"Ana","fecha_inicio":"2026-10-01","fecha_fin":"2026-12-01"}'   # 201
docker rm -f seg
docker run -d --name seg -p 8080:8080 -v seguimiento-data:/data seguimiento-medicion
curl -s localhost:8080/projects/1                                # el proyecto sigue (200)
docker rm -f seg
```

Tests: `go test ./...` y `cd web && npm test`.
