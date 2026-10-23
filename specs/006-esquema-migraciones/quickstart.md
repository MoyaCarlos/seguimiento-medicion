# Quickstart: validar la adopción (Issue #107)

Guía para quien revise el PR. No hace falta saber Go: cada paso dice qué correr y qué
esperar. Si algo no da lo esperado, pedile al agente que lo investigue.

## 1. Chequeos automáticos (CA-1, CA-2)

Desde la raíz del repo:

```bash
go build ./...
go vet ./...
gofmt -l .          # no debe listar ningún archivo
go test ./...       # incluye el test de esquema y los escenarios de Godog
```

Esperado: todo sin errores y `gofmt -l .` sin salida. Para ver solo el test de esquema y el
runner: `go test ./internal/repository -run 'Esquema|Migraciones' -v`.

## 2. Los `.feature` no cambiaron (SC-002)

```bash
git diff main --stat -- features/*.feature
```

Esperado: sin salida. (Los `steps_*.go` sí cambian; ver [research.md R8](research.md#r8-ajustes-en-los-tests-existentes-no-en-los-feature).)

## 3. Base vieja (CA-3, US5)

Con una `seguimiento.db` creada antes de este PR (o cualquier base que tenga `projects`,
`users`, `project_members` o `backlog_items`):

```bash
go run ./cmd/api
```

Esperado: el programa termina enseguida con un mensaje que pide borrar `seguimiento.db` y
reiniciar. Después de borrarla, `go run ./cmd/api` arranca normalmente.

## 4. Migraciones sin borrar la base (SC-001)

Con la API ya arrancada una vez sobre una base nueva, detenela y volvé a arrancarla.
Esperado: arranca sin errores y conserva los proyectos creados.

## 5. API sin cambios, salvo el 409 (US3, US4)

Con la API corriendo sobre una base nueva:

```bash
curl -s -X POST localhost:8080/projects -H 'Content-Type: application/json' \
  -d '{"nombre":"Demo","creador":"Ana","fecha_inicio":"2026-10-01","fecha_fin":"2026-12-01"}'
```

Esperado: 201 con las fechas en `AAAA-MM-DD`, igual que antes. Un pedido a un proyecto
inexistente sigue respondiendo 404. El 409 no se puede provocar desde la API en uso normal
(el servicio valida antes); lo prueban los tests de la fase 4
([contracts/errores-http.md](contracts/errores-http.md)).

## 6. Documento alineado (CA-4, CA-14)

```bash
git diff main -- MODELO-DE-DATOS.md
```

Esperado: cambios en §4.5 (nombres de índices), §5 (texto de los mensajes), §9 puntos 1 y 5
y §10 (tabla de control, las cuatro tablas viejas y la frase de la ventana en rojo, que arranca en el `GREEN:` de la fase 3a), y nada más que contradiga la migración 001.
