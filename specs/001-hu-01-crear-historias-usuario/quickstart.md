# Quickstart: Validación de HU-01 — Creación de Historias de Usuario

**Feature**: [spec.md](./spec.md) | [plan.md](./plan.md) | [data-model.md](./data-model.md) | [contracts/openapi.yaml](./contracts/openapi.yaml)

Guía para validar la creación de historias de usuario de punta a punta una vez
implementada la feature. No reemplaza a `tasks.md` ni contiene el código de
implementación.

## Prerrequisitos

- Go 1.26.5 instalado.
- Dependencias del proyecto descargadas (`modernc.org/sqlite`, Godog).
- Puerto `8080` libre.

## Preparación

```powershell
go get modernc.org/sqlite
go get github.com/cucumber/godog
go mod tidy
```

## Ejecutar la API

```powershell
go run ./cmd/api
```

Salida esperada: log de arranque escuchando en `:8080` y base SQLite creada con
la tabla `backlog_items`.

## Validación manual del endpoint

### Escenario 1 — Creación exitosa

```powershell
curl -i -X POST http://localhost:8080/backlog `
  -H "Content-Type: application/json" `
  -d '{"proyecto_id":1,"titulo":"Como usuario quiero iniciar sesión","descripcion":"Autenticación con email y contraseña","prioridad":"M","valor_negocio":13}'
```

Resultado esperado:

- HTTP `201 Created`.
- Cuerpo JSON con `id` asignado, `estado: "Nueva"`, `valor_negocio: 13` y
  `estimacion_sp: null`.
- Al reiniciar la API, la historia sigue persistida (verificar directamente en
  la base SQLite; el listado vía API queda fuera de alcance, ver R11) (FR-007).

### Escenario 2 — Faltan campos

```powershell
curl -i -X POST http://localhost:8080/backlog `
  -H "Content-Type: application/json" `
  -d '{"proyecto_id":1,"titulo":"","descripcion":"Sin título","prioridad":"M"}'
```

Resultado esperado:

- HTTP `400 Bad Request`.
- Cuerpo `{"campo":"titulo","mensaje":"el título es obligatorio"}`.
- No se crea ningún registro nuevo.

### Casos borde adicionales

- Título solo con espacios: `400` (`titulo`).
- Título > 200 caracteres o descripción > 2000: `400` (campo correspondiente).
- `prioridad` distinta de `M/S/C/W`: `400` (`prioridad`).
- `valor_negocio` fuera de `{1,2,3,5,8,13,21}`: `400` (`valor_negocio`).
- JSON malformado: `400` con mensaje de JSON inválido.

## Pruebas automatizadas

```powershell
go test ./...
```

Esperado: en verde los tests de TDD de dominio y servicio, el test de
integración del repositorio SQLite (`:memory:`) y los tests del handler
(`httptest`).

```powershell
go test ./features/...
```

Esperado: los escenarios BDD de `features/crear_historia_backlog.feature`
(creación exitosa y falta de título) pasan a través de Godog.

## Trazabilidad

| Criterio de aceptación | Escenario BDD | Test |
|------------------------|---------------|------|
| Creación exitosa (estado "Nueva", persistencia) | `features/crear_historia_backlog.feature` | handler + service + repo |
| Falta el título → advertencia, sin registro | `features/crear_historia_backlog.feature` | dominio (`NewBacklogItem`) + handler |
