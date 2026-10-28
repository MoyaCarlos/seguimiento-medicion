# Quickstart: Movimiento de Historias al Sprint Backlog (HU-06)

Guía para validar la funcionalidad de punta a punta. Contrato: [contracts/openapi.yaml](./contracts/openapi.yaml).
Reglas y orden de errores: [data-model.md](./data-model.md).

## Prerrequisitos

- Go 1.26+ y el repo con `main` mergeado (migración 001, sin migraciones nuevas).
- Base limpia: si `seguimiento.db` es de un esquema viejo, borrarla antes de arrancar.
- Para la pantalla (paso 3): en `main` el PR #149 (`GET /projects/{id}/sprints`) y el #98.

## 1. Tests automáticos

```bash
go test ./internal/...                    # dominio, service, repositorio, handlers
go test ./features/... -run TestFeatures  # escenarios BDD @HU-06 (y el resto)
gofmt -l .                                # no debe listar nada
```

Esperado: todo en verde.

## 2. Validación manual con la API

```bash
go run ./cmd/api
```

En otra terminal, dos proyectos, historias y Sprints:

```bash
curl -s -X POST localhost:8080/projects -d '{"nombre":"Demo","descripcion":"","fecha_inicio":null,"fecha_fin":null,"creador":"Ana"}'
curl -s -X POST localhost:8080/projects -d '{"nombre":"Otro","descripcion":"","fecha_inicio":null,"fecha_fin":null,"creador":"Ana"}'
curl -s -X POST localhost:8080/backlog -d '{"proyecto_id":1,"titulo":"Cerrar sesión","descripcion":"d","prioridad":"S"}'   # historia 1
curl -s -X POST localhost:8080/backlog -d '{"proyecto_id":1,"titulo":"Iniciar sesión","descripcion":"d","prioridad":"M"}'  # historia 2
curl -s -X POST localhost:8080/backlog -d '{"proyecto_id":1,"titulo":"Modo oscuro","descripcion":"d","prioridad":"W"}'     # historia 3
curl -s -X POST localhost:8080/sprints -d '{"proyecto_id":1}'   # Sprint 1, Pendiente
curl -s -X POST localhost:8080/sprints -d '{"proyecto_id":2}'   # Sprint 2, de otro proyecto
```

> El cuerpo exacto de `POST /projects` está en `specs/002-hu-04-crear-proyecto-equipo/contracts/`.

| # | Petición | Esperado |
|---|----------|----------|
| 1 | `GET /sprints/1/backlog` | 200 `[]` (FR-012) |
| 2 | `PUT /backlog/1/sprint` con `{"sprint_id":1}` | 200; `sprint_id: 1`, `estado: Nueva`, Story Points sin cambios |
| 3 | `GET /projects/1/backlog` | La historia 1 ya no aparece (FR-002) |
| 4 | `PUT /backlog/2/sprint` con `{"sprint_id":1}` y `GET /sprints/1/backlog` | 200; orden 2, 1 (M antes que S, FR-011) |
| 5 | `PUT /backlog/1/sprint` con `{"sprint_id":1}` | 400 `la historia ya pertenece a un Sprint` (mismo Sprint, FR-005) |
| 6 | `PUT /backlog/3/sprint` con `{"sprint_id":2}` | 400 `campo: sprint_id`, `el Sprint no pertenece al proyecto de la historia` |
| 7 | `PUT /backlog/3/sprint` con `{"sprint_id":99}` | 404 `sprint no encontrado` |
| 8 | `PUT /backlog/99/sprint` con `{}` | 404 `historia no encontrada` (no 400) |
| 9 | `PUT /backlog/3/sprint` con `{}` | 400 `campo: sprint_id`, `el campo es obligatorio` |
| 10 | `GET /sprints/99/backlog` | 404 `sprint no encontrado` |
| 11 | `POST /sprints/1/iniciar` (Goal y fechas) y `PUT /backlog/3/sprint` con `{"sprint_id":1}` | 200: un Sprint Activo admite historias |
| 12 | `POST /sprints/1/cerrar`, crear historia 4 y `PUT /backlog/4/sprint` con `{"sprint_id":1}` | 400 `un Sprint finalizado no admite historias` |
| 13 | Reiniciar `go run ./cmd/api` y repetir `GET /sprints/1/backlog` | Las asignaciones persisten (FR-013) |

Historia `Completada` (FR-004): HU-11 todavía no permite completar historias; se cubre en
`asignacion_sprint_test.go`, en el service y en el escenario BDD sembrando el estado directo en
la base. La carrera entre dos asignaciones se cubre en `sqlite_backlog_test.go`.

## 3. Pantalla (US3, cuando estén #149 y #98)

```bash
go run ./cmd/api
cd web && npm run dev
```

1. Entrar a un proyecto con historias en el backlog y un Sprint `Pendiente` o `Activo`.
2. En una historia, "Asignar a Sprint": el selector ofrece solo Sprints `Pendiente`/`Activo`.
3. Confirmar: el botón queda deshabilitado mientras la petición está en curso; al terminar, la
   historia desaparece del backlog y aparece en el Sprint Backlog del Sprint elegido.
4. Provocar un rechazo (por ejemplo, cerrar el Sprint en otra pestaña y asignar): se muestra el
   mensaje de la API y los listados quedan como estaban.

```bash
cd web && npm run lint && npm test && npm run build
```
