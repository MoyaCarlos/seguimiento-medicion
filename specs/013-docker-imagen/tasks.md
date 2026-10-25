---

description: "Tareas de la imagen Docker (Issue #132): TDD con commits por fase para la lógica; empaquetado sin RED/GREEN"
---

# Tasks: Imagen Docker multi-stage sin root

**Input**: `/specs/013-docker-imagen/` (spec.md, plan.md, research.md, data-model.md, contracts/servido-spa.md, quickstart.md)

**Tests**: obligatorios para la lógica (config y servido del SPA). Dockerfile, `.dockerignore` y README no llevan RED/GREEN.

**Reglas**: cada fase del ciclo es su propio commit (los COMMIT son tareas); archivos agregados por nombre (nunca `git add .`); sin `git push`. Go en `/usr/local/go/bin/go`.

## Phase 1: Foundational — configuración por entorno (US4, P2)

**Goal**: `PORT`, `DB_PATH` y `WEB_DIR` leídos por una función que recibe un `getenv`.

**Independent Test**: `go test ./cmd/api/`.

- [ ] T001 [US4] RED: test de `cargarConfig` en cmd/api/config_test.go — sin variables devuelve `8080`, `seguimiento.db`, `web/dist`; con variables definidas usa esas; variable vacía cae al default (FR-002, FR-003). Confirmar que FALLA
- [ ] T002 [US4] COMMIT `RED: agregar tests de la configuración por entorno (falla)` (solo cmd/api/config_test.go)
- [ ] T003 [US4] GREEN: implementar `cargarConfig(getenv func(string) string)` en cmd/api/config.go — `go test ./cmd/api/` en verde
- [ ] T004 [US4] COMMIT `GREEN: implementar la configuración por entorno, test en verde`
- [ ] T005 [US4] Cablear la config en cmd/api/main.go (`AbrirSQLite(cfg.RutaBase)`, escuchar en `":" + cfg.Puerto`, log con el puerto) — `go build ./...` y `go vet ./...` sin errores
- [ ] T006 [US4] COMMIT `Leer puerto y ruta de la base desde variables de entorno en cmd/api`

---

## Phase 2: User Story 5 — el servidor Go sirve el SPA sin romper la API (P1)

**Goal**: `ServirSPA(api, dirWeb)` replica la regla de `web/vite.config.ts` (contrato en contracts/servido-spa.md).

**Independent Test**: `go test ./internal/http/` con un `dirWeb` temporal (`index.html` y `assets/app.js`).

- [ ] T007 [US5] RED: tests de `ServirSPA` en internal/http/spa_test.go con una API falsa que responde `api` — `GET /projects/1/status` con `Accept: text/html` → index.html; sin ese `Accept` → API; `GET /assets/app.js` → el archivo; `GET /` con `text/html` → index.html; `GET /health` con `text/html` → API; `GET /assets/no-existe.js` sin `text/html` → 404; `POST` y `PUT` a `/projects` aun con `text/html` → API; `dirWeb` inexistente → todo a la API (FR-004..FR-007). Confirmar que FALLA
- [ ] T008 [US5] COMMIT `RED: agregar tests del servido del SPA (falla)` (solo internal/http/spa_test.go)
- [ ] T009 [US5] GREEN: implementar `ServirSPA` en internal/http/spa.go, con los prefijos exportados desde el inicio como `var PrefijosAPI = []string{"/backlog", "/projects", "/sprints", "/health"}` (así no hace falta un REFACTOR aparte) — `go test ./internal/http/` en verde
- [ ] T010 [US5] COMMIT `GREEN: implementar el servido del SPA, test en verde`
- [ ] T011 [US5] (solo si hay algo que limpiar) REFACTOR de spa.go + COMMIT `REFACTOR: <qué se limpió>`
- [ ] T012 [US5] Envolver el mux con `apihttp.ServirSPA(mux, cfg.DirWeb)` en cmd/api/main.go — `go build ./...`, `go vet ./...`, `gofmt -l .` y `go test ./...` sin errores
- [ ] T013 [US5] COMMIT `Servir el frontend desde el binario Go en cmd/api`
- [ ] T014 [US5] Test de coherencia en internal/http/spa_prefijos_test.go (FR-014): lee `../../web/vite.config.ts`, extrae con una regex simple las claves que apuntan a `apiProxy` (ej. `'/projects': apiProxy`) y verifica que cada una esté en `PrefijosAPI`; `/health` es solo del servidor Go y no se exige en Vite. Si la regex encuentra 0 prefijos, falla con "no se encontraron prefijos en vite.config.ts: ¿cambió el formato?". Pasa de entrada (no hay RED natural). Comprobar que detecta el problema con una prueba descartable: sacar temporalmente `"/sprints"` de `PrefijosAPI`, confirmar que el test FALLA y revertir a mano sin commitear
- [ ] T015 [US5] COMMIT `test: verificar que los prefijos del servidor Go incluyen los del proxy de Vite` (solo internal/http/spa_prefijos_test.go); antes, `go test ./...`, `go vet ./...` y `gofmt -l .` sin errores

**Checkpoint**: `WEB_DIR=web/dist` local sirve la UI y la API en `:8080`.

---

## Phase 3: User Stories 1, 2 y 3 — empaquetado (P1)

**Goal**: imagen única, persistente y no root.

**Independent Test**: pasos de quickstart.md.

- [ ] T016 [P] [US1] Crear .dockerignore (`**/node_modules`, `.git`, `*.db`, `docs/`, `specs/`, `web/dist`) — FR-011
- [ ] T017 [US1] [US2] [US3] Crear Dockerfile multi-stage según plan.md: `node:24-alpine` (`npm ci && npm run build`), `golang:1.26.5-alpine` (`CGO_ENABLED=0 go build`), `gcr.io/distroless/static-debian13:nonroot` con binario + `web/dist` + `/data` con dueño nonroot, `ENV PORT=8080 DB_PATH=/data/seguimiento.db WEB_DIR=/app/web/dist`, `VOLUME /data`, `EXPOSE 8080`, `USER nonroot:nonroot`, sin HEALTHCHECK
- [ ] T018 [US1] Agregar la sección "Correr con Docker" a README.md (`docker build`, `docker run -p 8080:8080 -v seguimiento-data:/data`, nota sobre bind mounts y permisos)
- [ ] T019 [US1] COMMIT `Agregar Dockerfile, .dockerignore y guía para correr con Docker` (por nombre: Dockerfile, .dockerignore, README.md)
- [ ] T020 [US1] [US2] [US3] Verificar con quickstart.md: `docker build`, `GET /health` 200, `Accept: text/html` devuelve el index, `docker inspect` muestra usuario no root, `GET /` y `/projects/1/status` con `Accept: text/html` devuelven el index, `POST /projects` (curl del quickstart) responde 201 y `GET /projects/1` devuelve el proyecto tras reiniciar con el mismo volumen. Si Docker no corre, dejar los pasos para Carlos

---

## Phase 4: Polish

- [ ] T021 `go build ./...`, `go vet ./...`, `gofmt -l .` y `go test ./...` sin errores; `cd web && npm test` en verde
- [ ] T022 Actualizar el estado de specs/013-docker-imagen/spec.md y marcar las tareas hechas + COMMIT `Cerrar la spec de la imagen Docker: estado y tareas (Issue #132)`

---

## Dependencies & Execution Order

- Phase 1 → Phase 2 (ambas tocan cmd/api/main.go) → Phase 3 (el Dockerfile necesita las variables y el servido) → Polish.
- T016 es paralelizable con T017/T018 (archivos distintos).
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT).

## Implementation Strategy

MVP = Phase 1 + 2 + 3 (todas las historias P1 dependen del servido y del Dockerfile); cada fase deja `go test ./...` en verde.
