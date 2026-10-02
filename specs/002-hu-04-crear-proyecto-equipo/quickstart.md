# Quickstart: Validación de HU-04 — Creación de Proyecto y Asignación de Equipo

**Feature**: [spec.md](./spec.md) | [plan.md](./plan.md) | [data-model.md](./data-model.md) | [contracts/openapi.yaml](./contracts/openapi.yaml)

Guía para validar la creación, edición y asignación de equipo de un proyecto de punta a
punta, una vez implementada la feature. No reemplaza a `tasks.md` ni contiene el código de
implementación.

## Prerrequisitos

- Go 1.26.5 instalado.
- Dependencias del proyecto descargadas (`modernc.org/sqlite`, Godog) — no se agregan dependencias nuevas.
- Puerto `8080` libre.
- Formato de fechas en la API: `YYYY-MM-DD` (RFC3339 en la base).

## Preparación

```powershell
go mod tidy
```

## Ejecutar la API

```powershell
go run ./cmd/api
```

Salida esperada: log de arranque escuchando en `:8080` y base SQLite creada con las tablas
`projects`, `users` y `project_members`.

## Escenario 1 — Creación del entorno del proyecto

```powershell
curl -i -X POST http://localhost:8080/projects `
  -H "Content-Type: application/json" `
  -d '{"nombre":"Software Metrics & Estimation","descripcion":"Sistema de estimación","creador":"Ana Valentina","fecha_inicio":"2026-03-01","fecha_fin":"2026-11-30"}'
```

Resultado esperado:

- HTTP `201 Created`.
- Cuerpo con `id` (UUID), `nombre`, `descripcion`, `fecha_inicio`, `fecha_fin`, `creado_en`.
- El creador queda vinculado: `GET /projects/{id}/members` devuelve a "Ana Valentina" con rol `scrum_master`.
- Al reiniciar la API, el proyecto y su integrante siguen disponibles (SC-002).

### Casos borde del alta

- Nombre vacío o solo espacios: `400` (`campo: "nombre"`).
- Nombre > 100 caracteres: `400` (`campo: "nombre"`).
- `fecha_fin` anterior a `fecha_inicio`: `400` (`campo: "fecha_fin"`).
- `descripcion` omitida: el proyecto se crea igual (opcional).
- JSON malformado: `400` con `{"mensaje":"JSON inválido"}`.

## Escenario 2 — Asignación de integrantes y roles

```powershell
curl -i -X POST http://localhost:8080/projects/{id}/members `
  -H "Content-Type: application/json" `
  -d '{"nombre":"Jimena Martinez","rol":"product_builder"}'
```

Resultado esperado:

- HTTP `201 Created`.
- Cuerpo con `id`, `proyecto_id`, `nombre` y `rol: "product_builder"`.
- `GET /projects/{id}/members` lista a Ana Valentina (`scrum_master`) y Jimena Martinez (`product_builder`).

### Casos borde de la asignación

- Sin nombre: `400` (`campo: "nombre"`).
- Rol distinto de `scrum_master`/`product_builder`: `400` (`campo: "rol"`).
- El mismo integrante dos veces en el mismo proyecto: `400` (`campo: "integrante"`).
- Reutilización de identidad: `"jimena martinez"` (minúsculas) resuelve al mismo integrante ya vinculado → `400` por duplicado, sin crear un usuario nuevo.
- Proyecto inexistente: `404`.

## Escenario 3 — Edición simple del proyecto

```powershell
curl -i -X PUT http://localhost:8080/projects/{id} `
  -H "Content-Type: application/json" `
  -d '{"nombre":"Software Metrics & Estimation (v2)","descripcion":"Descripción corregida","fecha_inicio":"2026-03-02","fecha_fin":"2026-12-15"}'
```

Resultado esperado:

- HTTP `200 OK` con el proyecto actualizado.
- `GET /projects/{id}` devuelve los nuevos valores tras recargar (SC-009).

### Casos borde de la edición

- Nombre vacío o solo espacios: `400` (`campo: "nombre"`); el proyecto conserva sus valores previos.
- `fecha_fin` anterior a `fecha_inicio`: `400` (`campo: "fecha_fin"`); sin cambios persistidos.
- `descripcion` vacía: se permite (opcional).
- Nombre repetido de otro proyecto: se permite (no hay unicidad).
- Proyecto inexistente: `404`.
- Sin versionado ni historial: no hay endpoints de versiones; prevalece el último guardado válido.

## Pruebas automatizadas

```powershell
go test ./...
```

Esperado: en verde los tests TDD de dominio (`Project`, `User`, `Role`, `Membership`), los
tests de casos de uso con fakes en memoria, los tests de integración de repositorios SQLite
(`:memory:`) y los tests del handler (`httptest`).

```powershell
go test ./features/...
```

Esperado: los escenarios BDD de `features/creacion_proyecto_equipo.feature` (crear, editar, asignar
y sus errores) pasan a través de Godog.

## Trazabilidad

| Criterio de aceptación | Escenario BDD | Test |
|------------------------|---------------|------|
| Crear proyecto persistido + creador Scrum Master (US1) | `features/creacion_proyecto_equipo.feature` | dominio + service + repo + handler |
| Alta con nombre vacío → advertencia, sin registro (US1) | `features/creacion_proyecto_equipo.feature` | dominio (`NewProject`) + handler |
| Asignar integrante con rol (US2) | `features/creacion_proyecto_equipo.feature` | dominio + service + repo + handler |
| Integrante sin nombre / rol inválido / duplicado (US2) | `features/creacion_proyecto_equipo.feature` | dominio + service + handler |
| Editar proyecto con datos válidos (US3) | `features/creacion_proyecto_equipo.feature` | dominio (`ConDatosEditados`) + service + repo |
| Editar con nombre vacío o fin < inicio (US3) | `features/creacion_proyecto_equipo.feature` | dominio + handler |
