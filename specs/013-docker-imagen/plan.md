# Implementation Plan: Imagen Docker multi-stage sin root

**Branch**: `feature/132-docker-imagen` | **Date**: 2026-10-24 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/013-docker-imagen/spec.md`

## Summary

Una sola imagen: el binario Go (`cmd/api`) sirve la API y `web/dist`. Se agregan dos piezas
pequeñas con TDD (lectura de configuración por entorno y handler que sirve el SPA replicando
la regla del proxy de Vite) y el empaquetado (Dockerfile multi-stage, `.dockerignore`, README).

## Technical Context

**Language/Version**: Go 1.26.5 (`go.mod`); Node 24 (`web/.nvmrc`, `engines >=24.15.0`) solo en build

**Primary Dependencies**: ninguna nueva (stdlib `net/http`, `os`, `path`)

**Storage**: SQLite embebido (`modernc.org/sqlite`, sin CGO); esquema sin cambios

**Testing**: `go test` (paquete `testing`, `httptest`); `npm test` (vitest) sin cambios

**Target Platform**: Linux contenedor (`linux/amd64`), imagen final distroless `static`, usuario no root

**Project Type**: web-service + SPA en un solo binario/imagen

**Constraints**: sin shell ni curl en la imagen final (sin HEALTHCHECK); `/data` escribible por el usuario no root

**Scale/Scope**: 2 archivos Go nuevos con sus tests, `main.go` (wiring), `Dockerfile`, `.dockerignore`, README

## Constitution Check

| Principio | Estado |
|---|---|
| I. Test-First | Config y servido del SPA con RED/GREEN en commits propios; Dockerfile/.dockerignore/README sin ciclo (no son lógica) |
| II. SDD | Spec → clarify → plan → tasks → analyze antes de implementar |
| III. BDD | Sin comportamiento de negocio nuevo: no aplica `.feature` (las reglas se prueban con tests HTTP) |
| IV. Clean/Hexagonal | Config en `cmd/api` (wiring); handler HTTP en `internal/http`; el dominio no se toca |
| V. KISS/YAGNI | Sin librerías ni abstracciones nuevas; sin segunda imagen |
| VI. Stack | Sin dependencias fuera del stack |

Sin violaciones. Re-chequeo post-diseño: igual.

## Project Structure

### Documentation (this feature)

```text
specs/013-docker-imagen/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/servido-spa.md
├── checklists/requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/api/config.go        # cargarConfig(getenv) -> {Puerto, RutaBase, DirWeb}
cmd/api/config_test.go
cmd/api/main.go          # usa la config y envuelve el mux con el handler del SPA
internal/http/spa.go     # ServirSPA(api http.Handler, dirWeb string) http.Handler
internal/http/spa_test.go
internal/http/spa_prefijos_test.go   # coherencia con web/vite.config.ts (FR-014)
Dockerfile
.dockerignore
README.md                # sección "Correr con Docker"
```

**Structure Decision**: la regla del SPA vive en `internal/http` (es HTTP), la lectura de entorno en `cmd/api` (es wiring). `ServirSPA` envuelve al mux en vez de registrar rutas, así no se tocan los handlers ni sus tests.

## Diseño

- `cargarConfig(getenv func(string) string)`: cada valor vacío cae al default (`8080`, `seguimiento.db`, `web/dist`). `main` escucha en `":" + Puerto`.
- `ServirSPA(api, dirWeb)`: si `dirWeb` no existe o no es directorio, devuelve `api` tal cual (FR-007). Si no, por request:
  1. `GET`/`HEAD` de un archivo regular existente en `dirWeb` → se sirve (`http.FileServer`).
  2. `GET` con `Accept` que contiene `text/html` y ruta no-API (`/health` es API) → `index.html`; las rutas bajo `/backlog`, `/projects`, `/sprints` también → `index.html` (FR-004, FR-005b).
  3. Todo lo demás → `api`.
- `PrefijosAPI` es una variable exportada de `internal/http`; `spa_prefijos_test.go` lee `web/vite.config.ts` y verifica que sus prefijos estén incluidos (FR-014).
- Prefijos de la API: los mismos tres de `web/vite.config.ts` más `/health`; coincide la ruta exacta o una subruta (`/projects`, `/projects/1`), no `/projectsx`.
- Dockerfile (3 stages): `node:24-alpine` (`npm ci && npm run build`) → `golang:1.26.5-alpine` (`CGO_ENABLED=0 go build -trimpath -ldflags="-s -w"`) → `gcr.io/distroless/static-debian13:nonroot`. En la final: `COPY --chown=nonroot:nonroot` del binario y `web/dist`, y `/data` creado en un stage previo y copiado con dueño `nonroot` (sin shell en la final), `ENV PORT=8080 DB_PATH=/data/seguimiento.db WEB_DIR=/app/web/dist`, `VOLUME /data`, `EXPOSE 8080`, `USER nonroot:nonroot`, `WORKDIR /app`.
