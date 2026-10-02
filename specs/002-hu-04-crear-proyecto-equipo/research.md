# Research: HU-04 — Creación de Proyecto y Asignación de Equipo

**Fecha**: 2026-10-01
**Feature**: [spec.md](./spec.md) | [plan.md](./plan.md)

Este documento resuelve los puntos técnicos abiertos del plan. No quedan
`NEEDS CLARIFICATION` en el Technical Context.

## R1. Identificador del proyecto

- **Decisión**: mantener `Project.ID` como `string` y generar un UUID v4 con `crypto/rand` desde el adaptador de persistencia al crear (si `ID == ""`). No se cambian las firmas existentes `Create(p *domain.Project) error` ni `GetByID(id string) (*domain.Project, error)`.
- **Rationale**: el scaffolding de HU-04/HU-05 ya define `id string` y `GetByID(id string)`; cambiarlo a `int64` rompería el contrato con HU-05. Un UUID evita colisiones y no requiere autoincremento.
- **Alternativas consideradas**: `int64` con `AUTOINCREMENT` (rompería HU-05); `github.com/google/uuid` (dependencia ya presente como indirecta, pero `crypto/rand` evita promover una dependencia nueva, alineado con KISS/YAGNI).

## R2. Representación de fechas opcionales

- **Decisión**: `StartDate *time.Time` y `EndDate *time.Time` en el dominio; en SQLite se guardan como `TEXT` en formato RFC3339 y `NULL` cuando son `nil`. Validación: si ambas existen, `EndDate` no puede ser anterior a `StartDate`.
- **Rationale**: distingue "sin fecha" de una fecha cero; mantiene el dominio libre de tipos de `database/sql` (mismo criterio que `*int` en HU-01).
- **Alternativas consideradas**: `time.Time` con centinela (ambiguo); `sql.NullTime` en el dominio (filtraría persistencia); guardar fechas como `string` sin parseo (pierde validación de orden).
- **Límite con HU-13**: HU-04 solo registra y edita fechas con la validación de orden. HU-13 (Should have) es quien consulta el estado derivado (Planificado/En curso/Finalizado) a partir de fechas y sprints.

## R3. Identidad de integrante por nombre (reutilización)

- **Decisión**: cada `User` guarda `Name` (nombre tal como se ingresó, recortado) y `NormalizedName` = `strings.ToLower(strings.TrimSpace(Name))` con restricción `UNIQUE`. La vinculación busca por `NormalizedName`; si existe, reutiliza el usuario; si no, lo crea. "Ana", "ana" y " Ana " resuelven al mismo usuario.
- **Rationale**: implementa la clarificación de `/speckit.clarify` (comparación sin distinguir mayúsculas/minúsculas ni espacios externos) de forma determinista y testeable a nivel de base.
- **Alternativas consideradas**: comparar con `COLLATE NOCASE` sobre `Name` (no maneja el *trim* de forma explícita); buscar con `LOWER(TRIM(name))` en cada consulta (no indexable y propenso a duplicados). La columna normalizada + `UNIQUE` resuelve ambos.
- **Concurrencia**: ante colisión del `UNIQUE` al crear, el repositorio re-consulta por nombre normalizado y devuelve el existente, evitando duplicados.

## R4. Modelo de integrantes y roles

- **Decisión**: entidades `User` (integrante), `Role` (enum) y `Membership` (relación Proyecto–Usuario–Rol). Un `User` puede pertenecer a varios proyectos; dentro de un proyecto tiene exactamente un rol. Restricción `PRIMARY KEY (project_id, user_id)` en `project_members`, que rechaza la vinculación duplicada (FR-015).
- **Rationale**: FR-016 (un integrante en varios proyectos) y SC-007 (un rol por integrante y proyecto) se expresan directamente en el modelo.
- **Alternativas consideradas**: guardar el rol dentro de `users` (impediría pertenecer a varios proyectos con roles distintos); tabla de roles separada (innecesaria para dos valores fijos, YAGNI).

## R5. Valores canónicos de rol

- **Decisión**: códigos canónicos `"scrum_master"` y `"product_builder"` en dominio y persistencia; la API los expone/acepta con esos valores. Se documenta el mapeo a las etiquetas visibles "Scrum Master" y "Product Builder".
- **Rationale**: códigos estables independientes del idioma de la UI; validación simple contra el enum.
- **Alternativas consideradas**: usar las etiquetas con espacio (`"Scrum Master"`) como valor (más frágil para comparaciones y para la DB); números (menos legibles en la base).

## R6. Reutilización de validaciones entre alta y edición

- **Decisión**: una única función de dominio (`NewProject` y su helper de validación) es usada por `CrearProyecto` y por `EditarProyecto`. La edición valida nombre (obligatorio, ≤ 100), descripción (≤ 2000) y orden de fechas antes de persistir.
- **Rationale**: FR-019 exige reutilizar las mismas validaciones; el dominio es el único lugar de reglas (DRY, SRP).
- **Alternativas consideradas**: validar en cada service (duplicaría reglas); confiar solo en validación HTTP (filtra reglas a la capa de transporte).

## R7. Semántica de la edición (endpoint)

- **Decisión**: `PUT /projects/{id}` con reemplazo de los campos editables (`nombre` obligatorio; `descripcion`, `fecha_inicio`, `fecha_fin` opcionales). No se versiona ni se guarda historial; `UPDATE` simple (prevalece el último guardado válido).
- **Rationale**: "edición simple" de FR-018/FR-023; `PUT` refleja reemplazo total de los campos editables sin requerir parches parciales.
- **Alternativas consideradas**: `PATCH` parcial (agrega complejidad de campos presentes/ausentes sin valor en el MVP); endpoints por campo (sobre-ingeniería).

## R8. Manejo de errores

- **Decisión**: se reutiliza `domain.ValidationError{Campo, Mensaje}`; el handler responde `400` con `{campo, mensaje}`. Proyecto inexistente → `404 {mensaje}`. JSON malformado → `400 {mensaje: "JSON inválido"}`. Fallo de persistencia → `500`.
- **Rationale**: consistente con HU-01; mantiene al handler sin reglas de negocio.
- **Alternativas consideradas**: errores centinela por caso (más código, menos contexto de campo).

## R9. Puertos y separación de interfaces

- **Decisión**: `ProjectRepository` (existente) se extiende con `Update`, `AddMember`, `ListMembers`; se agrega `UserRepository` con `FindByNormalizedName` y `Create` (o `FindOrCreate`). El `SQLiteProjectRepository` y `SQLiteUserRepository` se inyectan solo en `cmd/api`; los tests de `service` usan fakes en memoria.
- **Rationale**: DIP/SRP; la búsqueda/creación de usuarios es una responsabilidad distinta de la de proyectos.
- **Alternativas consideradas**: un único repositorio con todos los métodos (violaría SRP); declarar la interfaz en el paquete `service` (el proyecto ya la ubica en `repository`, se mantiene por consistencia con HU-01).

## R10. Migración del esquema

- **Decisión**: ampliar `repository.Migrar` con `CREATE TABLE IF NOT EXISTS` para `projects`, `users` y `project_members` (idempotente), y activar `PRAGMA foreign_keys = ON` al abrir la conexión.
- **Rationale**: KISS; mismo mecanismo ya usado por HU-01; alcanza para los tests con SQLite `:memory:`.
- **Alternativas consideradas**: herramienta de migraciones (goose/golang-migrate): sobre-ingeniería para tres tablas.

## R11. Generación de IDs

- **Decisión**: helper `nuevoID()` en `internal/repository` que arma un UUID v4 a partir de 16 bytes de `crypto/rand` (versión/variante fijadas, formateado canónico). El repositorio asigna el ID en `Create` cuando llega vacío.
- **Rationale**: sin dependencias nuevas; IDs no predecibles; el dominio no se contamina con generación de IDs.
- **Alternativas consideradas**: `math/rand` (no criptográfico, colisionable y sin semilla garantizada); autoincremento entero (rompería `id string`).

## R12. Alcance de frontend

- **Decisión**: fuera de esta iteración. La API devuelve el proyecto creado/actualizado y el integrante vinculado; la pantalla principal, la vista de configuración y la redirección del SPA se implementan en un incremento de frontend.
- **Rationale**: coherente con HU-01 (backend-only) y con la Definition of Done por capas; la redirección es comportamiento de UI.

## R13. Dependencias a agregar

- Ninguna. Se reutilizan `modernc.org/sqlite`, Godog y la librería estándar (`crypto/rand` para IDs).

## Resumen de decisiones

| ID | Tema | Decisión |
|----|------|----------|
| R1 | ID de proyecto | `string` UUID, firmas de HU-05 intactas |
| R2 | Fechas | `*time.Time`, RFC3339/NULL, fin ≥ inicio |
| R3 | Identidad de integrante | `NormalizedName` (lower+trim) único |
| R4 | Integrantes/roles | `User` + `Role` + `Membership`, PK (project,user) |
| R5 | Valores de rol | `scrum_master`, `product_builder` |
| R6 | Validaciones | Dominio único, reutilizado por alta y edición |
| R7 | Edición | `PUT` reemplazo, sin versionado ni historial |
| R8 | Errores | `ValidationError` → 400; 404/500 |
| R9 | Puertos | `ProjectRepository` extendido + `UserRepository` |
| R10 | Esquema | `CREATE TABLE IF NOT EXISTS`, FK ON |
| R11 | IDs | UUID v4 con `crypto/rand` |
| R12 | Frontend | diferido |
| R13 | Dependencias | ninguna nueva |
