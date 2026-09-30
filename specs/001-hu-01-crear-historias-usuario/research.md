# Research: HU-01 — Creación de Historias de Usuario

**Fecha**: 2026-09-30
**Feature**: [spec.md](./spec.md) | [plan.md](./plan.md)

Este documento resuelve los puntos técnicos abiertos del plan. No quedan
`NEEDS CLARIFICATION` en el Technical Context.

## R1. Enrutamiento HTTP

- **Decisión**: usar la librería estándar `net/http` con `http.ServeMux` y patrones por método (`"POST /backlog"`), disponibles desde Go 1.22.
- **Rationale**: KISS/YAGNI; el proyecto ya usa `net/http` en `cmd/api/main.go`. Evita sumar un framework para un único endpoint.
- **Alternativas consideradas**: `chi` y `gin`. Descartadas: agregan dependencia y ceremonia sin aportar valor en el MVP.

## R2. Driver de persistencia SQLite

- **Decisión**: `modernc.org/sqlite` (driver registrado como `"sqlite"`), con `database/sql`.
- **Rationale**: es el driver puro Go recomendado en `AGENTS.md`; no requiere CGO ni compilador C en las máquinas del equipo.
- **Alternativas consideradas**: `mattn/go-sqlite3` (requiere CGO, descartado) y `glebarez/sqlite` (GORM, innecesario: traería un ORM no requerido).

## R3. Representación de la estimación "sin asignar"

- **Decisión**: modelar `EstimacionSP *int` en el dominio (`nil` = sin estimar), mapeado a `NULL` en SQLite. `ValorNegocio` también es `*int` por ser opcional.
- **Rationale**: distingue "sin estimar" de cualquier valor numérico; mantiene el dominio libre de `database/sql`.
- **Alternativas consideradas**: centinela `0` (ambiguo, 0 no integra la escala Fibonacci), `sql.NullInt64` en el dominio (filtraría detalles de persistencia al núcleo).

## R4. Modelo de errores de validación

- **Decisión**: tipo `domain.ValidationError` con `Campo` y `Mensaje`; el handler HTTP lo detecta con `errors.As` y responde `400` con el campo y el mensaje.
- **Rationale**: FR-006 exige que la advertencia identifique el campo faltante. Un tipo tipado permite construir el contrato de error sin que el handler conozca reglas de negocio.
- **Alternativas consideradas**: errores centinela (`var ErrTituloVacio = errors.New(...)`): dificultan asociar el campo de forma limpia; `panic`/recover: inadecuado para validación esperada.

## R5. Asociación a proyecto

- **Decisión**: la entidad guarda `ProyectoID int64` y la tabla `backlog_items` incluye la columna `proyecto_id` `NOT NULL`. La presencia se valida en el dominio (`ProyectoID > 0`, invariante del constructor). No se declara `FOREIGN KEY` a `projects` porque esa tabla aún no existe (HU-04). La validación de que el proyecto exista en la base queda diferida a HU-04.
- **Rationale**: preserva la relación exigida por FR-012 sin inventar el modelo de proyectos (YAGNI).
- **Alternativas consideradas**: omitir `proyecto_id` (rompería FR-012); declarar FK a tabla inexistente (fallaría la migración).

## R6. Gestión del esquema

- **Decisión**: función `repository.Migrate(ctx, db)` idempotente con `CREATE TABLE IF NOT EXISTS`, invocada al arrancar `cmd/api`.
- **Rationale**: KISS; alcanza para el MVP y para los tests (SQLite `:memory:`).
- **Alternativas consideradas**: herramienta de migraciones (goose, golang-migrate): sobre-ingeniería para una sola tabla.

## R7. Inversión de dependencias (puerto/adaptador)

- **Decisión**: la interfaz `repository.BacklogRepository` se declara en el paquete `repository`; `service.CrearHistoriaBacklog` depende de esa interfaz. El adaptador concreto `SQLiteBacklogRepository` se inyecta solo en `cmd/api`. Los tests de `service` usan un fake en memoria.
- **Rationale**: SOLID (DIP) y testabilidad; el dominio y los casos de uso no conocen SQLite.

## R8. Autorización por rol (FR-014)

- **Decisión**: diferida. El handler asume que el llamador ya fue autenticado/autorizado; la verificación de rol (Product Builder / Scrum Master) se implementará cuando HU-04 introduzca usuarios y sesión.
- **Rationale**: no existe aún el mecanismo de identidad; implementarlo ahora excede el alcance solicitado (YAGNI).

## R9. Valor canónico del estado inicial

- **Decisión**: valor canónico `"Nueva"` (enum `EstadoNueva`), referido a "Historia de Usuario" (femenino). **RESUELTO** en `/speckit.analyze`: el equipo confirmó `"Nueva"`.
- **Acción aplicada**: se normalizó `spec.md`, `plan.md`, `tasks.md`, `data-model.md`, `contracts/openapi.yaml` y `docs/PRODUCT-BACKLOG.md` (convención de campos y valores de estado) a `"Nueva"`. No quedan referencias a `"Nuevo"` como estado.

## R11. Listado/vista del Product Backlog

- **Decisión**: fuera del alcance de HU-01. La historia solo crea, persiste y responde; el orden de creación se preserva con `id` incremental (`AUTOINCREMENT`).
- **Rationale**: la iteración es backend-only; la vista/orden visual y la actualización en vivo sin recargar (ex SC-006) son responsabilidad de un incremento de frontend.
- **Efecto en artefactos**: FR-003 y SC-002 se reformularon a lo verificable por backend; SC-006 se retiró de la spec y se documentó en Assumptions. No se agrega endpoint `GET /backlog` en esta historia.

## R12. Autorización por rol (FR-014)

- **Decisión**: diferida a HU-04 (identidad y roles). HU-01 no aplica control de acceso por rol.
- **Rationale**: no existe aún el mecanismo de identidad; implementarlo ahora excede el alcance (YAGNI). Se marcó explícitamente en `spec.md` (FR-014) y en `plan.md` (Complexity Tracking).

## R10. Dependencias a agregar

- `modernc.org/sqlite` (runtime).
- `github.com/cucumber/godog` (testing BDD).
- Se incorporan con `go get` y se fijan en `go.mod`/`go.sum` durante la implementación (fase `/speckit.tasks` → `/speckit.implement`).

## Resumen de decisiones

| ID | Tema | Decisión |
|----|------|----------|
| R1 | Routing | `net/http` stdlib (Go 1.22+) |
| R2 | SQLite | `modernc.org/sqlite` (sin CGO) |
| R3 | Estimación opcional | `*int` → `NULL` |
| R4 | Errores de validación | `domain.ValidationError` → 400 |
| R5 | Proyecto | guardar `proyecto_id`; validación diferida a HU-04 |
| R6 | Esquema | `CREATE TABLE IF NOT EXISTS` al iniciar |
| R7 | Puertos | `BacklogRepository` + fake en tests |
| R8 | Rol | diferido a HU-04 |
| R9 | Estado inicial | `"Nueva"` (normalizado en todos los artefactos) |
| R10 | Dependencias | `modernc.org/sqlite`, Godog |
| R11 | Listado/UI | fuera de alcance; order por `id` incremental |
| R12 | FR-014 autorización | diferida a HU-04 |
