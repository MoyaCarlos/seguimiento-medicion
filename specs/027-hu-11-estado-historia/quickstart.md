# Quickstart: Actualización de Estado de una Historia en el Sprint (HU-11)

Guía para validar la funcionalidad de punta a punta. Contrato: [contracts/openapi.yaml](./contracts/openapi.yaml).
Reglas y orden de errores: [data-model.md](./data-model.md).

## Prerrequisitos

- Go 1.26+ y el repo con `main` mergeado (migración 001, sin migraciones nuevas).
- Base limpia: si `seguimiento.db` es de un esquema viejo, borrarla antes de arrancar.
- Para la pantalla (paso 3): que la pantalla de Sprints esté en `main`.

## 1. Tests automáticos

```bash
go test ./internal/...                    # dominio, service, repositorio, handlers
go test ./features/... -run TestFeatures  # escenarios BDD @HU-11 (y el resto)
gofmt -l .                                # no debe listar nada
```

Esperado: todo en verde.

## 2. Validación manual con la API

```bash
go run ./cmd/api
```

En otra terminal, un proyecto, dos historias y un Sprint activo:

```bash
curl -s -X POST localhost:8080/projects -d '{"nombre":"Demo","descripcion":"","fecha_inicio":null,"fecha_fin":null,"creador":"Ana"}'
curl -s -X POST localhost:8080/backlog -d '{"proyecto_id":1,"titulo":"Cerrar sesión","descripcion":"d","prioridad":"S"}'   # historia 1
curl -s -X POST localhost:8080/backlog -d '{"proyecto_id":1,"titulo":"Modo oscuro","descripcion":"d","prioridad":"M"}'     # historia 2 (sin Sprint)
curl -s -X POST localhost:8080/sprints -d '{"proyecto_id":1}'                                                              # Sprint 1, Pendiente
curl -s -X POST localhost:8080/sprints/1/iniciar -d '{"sprint_goal":"Meta","fecha_inicio":"2026-10-20","fecha_fin":"2026-11-02"}'
curl -s -X PUT localhost:8080/backlog/1/sprint -d '{"sprint_id":1}'
```

> El cuerpo exacto de `POST /projects` está en `specs/002-hu-04-crear-proyecto-equipo/contracts/`.

| # | Petición | Esperado |
|---|----------|----------|
| 1 | `PUT /backlog/1/estado` con `{"estado":"Completada"}` | 400 `campo: estado`, `primero debe pasar por En progreso` (FR-013) |
| 2 | `PUT /backlog/1/estado` con `{"estado":"En progreso"}` | 200; `estado: En progreso`, `fecha_finalizacion: null` |
| 3 | `PUT /backlog/1/estado` con `{"estado":"En progreso"}` | 400 `campo: estado`, `la historia ya está en ese estado` (FR-016) |
| 4 | `PUT /backlog/1/estado` con `{"estado":"Nueva"}` | 400 `campo: estado`, `el estado de una historia solo avanza...` (FR-015) |
| 5 | `PUT /backlog/1/estado` con `{"estado":"Completada"}` | 200; `estado: Completada`, `fecha_finalizacion` con la fecha de hoy (FR-003) |
| 6 | `GET /projects/1/backlog` | La historia 1 no aparece: sigue en el Sprint (FR-010) |
| 7 | `PUT /backlog/1/estado` con `{"estado":"En progreso"}` | 400 `campo: estado`, `el estado de una historia solo avanza...` (FR-015) |
| 8 | `PUT /backlog/2/estado` con `{"estado":"En progreso"}` | 400 `campo: estado`, `primero debe asignarse a un Sprint activo` (FR-004) |
| 9 | `PUT /backlog/2/estado` con `{"estado":"Completada"}` | Igual que el 8 (FR-004) |
| 10 | `PUT /backlog/1/estado` con `{}` | 400 `campo: estado`, `el campo es obligatorio` (FR-008) |
| 11 | `PUT /backlog/1/estado` con `{"estado":"Pendiente"}` | 400 `campo: estado`, `el estado debe ser Nueva, En progreso o Completada` (FR-008) |
| 12 | `PUT /backlog/99/estado` con `{}` | 404 `historia no encontrada` (no 400, FR-007) |
| 13 | `POST /sprints/1/cerrar` y después `PUT /backlog/1/estado` con `{"estado":"En progreso"}` | 400 `campo: estado`, `solo se cambia el estado de historias de un Sprint activo` (FR-014) |
| 14 | `POST /sprints` (crea el Sprint 2, Pendiente), `PUT /backlog/2/sprint` con `{"sprint_id":2}` y `PUT /backlog/2/estado` con `{"estado":"En progreso"}` | 400 `solo se cambia el estado de historias de un Sprint activo` (FR-014) |
| 15 | Reiniciar `go run ./cmd/api` y repetir el paso 7 sobre la historia 1 | 400 `solo se cambia el estado de historias de un Sprint activo`: el Sprint 1 siguió Finalizado después del reinicio (FR-009). La fecha de finalización se verifica en los tests de repositorio y en la respuesta del paso 5. |

La carrera entre dos cambios simultáneos (FR-017, 409) no se puede forzar desde `curl` de forma
determinista: se cubre en `sqlite_backlog_test.go` y en el test del handler.

## 3. Pantalla (US3, cuando la pantalla de Sprints esté en `main`)

```bash
go run ./cmd/api
cd web && npm run dev
```

1. Entrar al Sprint activo de un proyecto con historias.
2. En una historia "Nueva", "Iniciar": queda "En progreso" sin recargar a mano.
3. En esa historia, "Completar": queda "Completada" y se ve su fecha de finalización.
4. Provocar un rechazo (por ejemplo, intentar completar una historia del Product Backlog): se
   muestra el mensaje de la API y la historia queda como estaba.
5. Ante un 409 (otra persona cambió la historia), la pantalla muestra el aviso de recargar y
   reintentar.

```bash
cd web && npm run lint && npm test && npm run build
```
