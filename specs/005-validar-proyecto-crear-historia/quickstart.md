# Quickstart: Validar proyecto existente al crear historia de backlog (Issue #19)

## Prerrequisitos

- Go 1.26 instalado (`go build ./...` corre limpio en el repo).
- HU-01, HU-04 y HU-05 ya están en `main`: se reutilizan `fakeProjectRepository`
  (`internal/service/fakes_test.go`) en los tests de service y `stubProjectRepository`
  (`internal/http/project_handler_test.go`) en los tests del handler. No se crean fakes nuevos.
- No hay cambios de esquema: no hace falta recrear `seguimiento.db`.

## Validar la regla en el caso de uso

```bash
go test ./internal/service/... -run CrearHistoria -v
```

Escenarios a verificar (los tres nuevos, en `internal/service/crear_historia_test.go`):

| Caso | Entrada | Resultado esperado |
|---|---|---|
| Proyecto inexistente | historia válida + `ProyectoID` sin proyecto en el fake | `errors.Is(err, domain.ErrProyectoNoEncontrado)` y `len(repo.guardados) == 0` |
| Error genérico de `ObtenerPorID` | fake con `err` seteado (fallo de persistencia) | el mismo error se propaga y `len(repo.guardados) == 0` |
| Historia inválida con proyecto inexistente | título vacío + `ProyectoID` inexistente | `domain.ValidationError` con `campo == "titulo"` (no se consulta el proyecto) |
| Camino feliz (existente) | historia válida + proyecto precargado | sin error, `len(repo.guardados) == 1` |

Los tests existentes (`Exitosa`, `NoPersisteSiEsInvalida`, `ValorNegocio`) deben seguir
pasando: ahora reciben además un `fakeProjectRepository` con el proyecto precargado.

## Validar el handler (sin cambio de comportamiento pedido)

```bash
go test ./internal/http/... -run Backlog -v
```

- `POST /backlog` con proyecto existente → `201`.
- `POST /backlog` con `proyecto_id` ausente o inválido → `400` (invariante de `BacklogItem`,
  se evalúa antes de consultar el proyecto).
- El caso 404 por proyecto inexistente **no se agrega acá** (FR-007 fuera de alcance); el
  mapper compartido `escribirError` ya lo traduce a `404` cuando el service devuelve
  `domain.ErrProyectoNoEncontrado`.

## Correr el escenario BDD (Godog)

```bash
go test ./features/... -run TestFeatures -v
```

- Escenario nuevo en `features/crear_historia_backlog.feature`: crear una historia para un
  proyecto inexistente → el sistema responde "proyecto no encontrado" y no registra la
  historia en el Product Backlog.
- Los pasos reutilizan los helpers existentes (`existeProyecto`, `noRegistraHistoria`) en
  `features/steps_creacion.go`.

## Verificación manual end-to-end (opcional)

```bash
go run ./cmd/api
# en otra terminal (el proyecto 999 no existe):
curl -i -X POST http://localhost:8080/backlog -H "Content-Type: application/json" \
  -d '{"proyecto_id":999,"titulo":"Historia huérfana","descripcion":"No debe persistir","prioridad":"M"}'
```

Resultado esperado: `HTTP/1.1 404 Not Found` con `{"campo":"","mensaje":"proyecto no encontrado"}`,
y `backlog_items` sin filas nuevas. (El mapper compartido `escribirError` usa `errorResponse`,
por eso el cuerpo incluye `campo` vacío.)

## Suite completa

```bash
go build ./... ; go test ./...
```
