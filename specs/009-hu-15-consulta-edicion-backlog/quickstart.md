# Quickstart: Consulta y Edición del Product Backlog (HU-15)

Guía para validar la funcionalidad de punta a punta. Contrato: [contracts/openapi.yaml](./contracts/openapi.yaml).
Reglas: [data-model.md](./data-model.md).

## Prerrequisitos

- Go 1.26+ y el repo con `main` mergeado (migración 001 con `criterios_aceptacion`).
- Base limpia: si `seguimiento.db` es de un esquema viejo, borrarla antes de arrancar.

## 1. Tests automáticos

```bash
go test ./internal/...                 # dominio, service, repositorio, handler
go test ./features/... -run TestFeatures  # escenarios BDD @HU-15 (y el resto)
gofmt -l .                             # no debe listar nada
```

Esperado: todo en verde.

## 2. Validación manual con la API

```bash
go run ./cmd/api
```

En otra terminal:

```bash
# Proyecto e historias de distintas prioridades (creadas en orden S, M, W, M)
curl -s -X POST localhost:8080/projects -d '{"nombre":"Demo","descripcion":"","fecha_inicio":null,"fecha_fin":null,"creador":"Ana"}'
curl -s -X POST localhost:8080/backlog -d '{"proyecto_id":1,"titulo":"Cerrar sesión","descripcion":"d","prioridad":"S"}'
curl -s -X POST localhost:8080/backlog -d '{"proyecto_id":1,"titulo":"Iniciar sesión","descripcion":"d","prioridad":"M"}'
curl -s -X POST localhost:8080/backlog -d '{"proyecto_id":1,"titulo":"Modo oscuro","descripcion":"d","prioridad":"W"}'
curl -s -X POST localhost:8080/backlog -d '{"proyecto_id":1,"titulo":"Recuperar clave","descripcion":"d","prioridad":"M"}'
```

> El cuerpo exacto de `POST /projects` está en `specs/002-hu-04-crear-proyecto-equipo/contracts/`.

| # | Petición | Esperado |
|---|----------|----------|
| 1 | `GET /projects/1/backlog` | 200, ids en orden 2, 4, 1, 3 (M, M, S, W); `criterios_aceptacion: null` |
| 2 | `GET /projects/99/backlog` | 404 `proyecto no encontrado` |
| 3 | `PUT /backlog/2` con `{"titulo":"Login","descripcion":"OAuth","criterios_aceptacion":"Dado...","estimacion_sp":5}` | 200; mismo id, estado `Nueva`, prioridad `M` |
| 4 | `GET /projects/1/backlog` | La historia 2 muestra los datos nuevos |
| 5 | `PUT /backlog/2` con `"titulo":"   "` | 400 `campo: titulo`; un GET muestra la historia sin cambios |
| 6 | `PUT /backlog/2` con `"estimacion_sp":4` (y con `2.5`) | 400 `campo: estimacion_sp` |
| 7 | `PUT /backlog/2` con `"estimacion_sp":null` | 200; `estimacion_sp: null` |
| 8 | `PUT /backlog/99` con `"titulo":""` | 404 `historia no encontrada` (no 400) |
| 9 | Reiniciar `go run ./cmd/api` y repetir el paso 4 | Los cambios persisten (FR-013) |

Exclusión por Sprint (FR-002) y Story Points congelados en Sprints `Activo`/`Finalizado`
(FR-015, D2b): hasta HU-06 no hay endpoint para asignar historias; se cubren en
`sqlite_backlog_test.go`, `editar_historia_test.go` y en los escenarios BDD sembrando
`sprint_id` (y el estado del Sprint) directo en la base.
