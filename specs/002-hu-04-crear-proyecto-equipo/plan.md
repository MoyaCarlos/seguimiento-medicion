# Implementation Plan: Creación de Proyecto y Asignación de Equipo (HU-04)

**Branch**: `hu-4-creación-de-proyecto-asignación-de-equipo` | **Date**: 2026-10-01 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-hu-04-crear-proyecto-equipo/spec.md`

## Summary

Implementar la gestión del ciclo básico de un proyecto de punta a punta (backend) con
arquitectura hexagonal liviana (Ports & Adapters): entidad de dominio `Project` con
invariantes reutilizadas por creación y edición, entidades `User` (integrante), `Role`
y `Membership` (relación Proyecto–Usuario–Rol); casos de uso `CrearProyecto`,
`EditarProyecto` y `AsignarIntegrante`; puertos `ProjectRepository` (existente,
extendido) y `UserRepository` (nuevo); adaptadores de persistencia SQLite
(`modernc.org/sqlite`); y adaptadores de entrada HTTP. Se reutilizan la entidad
`domain.Project` y el puerto `ProjectRepository` (`Guardar`/`ObtenerPorID` con `id int64`).

**Alcance de esta iteración (backend)**: creación, edición simple (nombre, descripción y
fechas) y asignación de integrantes con rol, con persistencia y validaciones compartidas.
La navegación/redirección y las vistas del frontend (pantalla principal, configuración y
edición) quedan diferidas a un incremento de frontend (ver Complexity Tracking). La
aplicación efectiva de permisos por rol y la consulta del estado de ciclo de vida del
proyecto (HU-13) quedan fuera de alcance.

## Technical Context

**Language/Version**: Go 1.26.5

**Primary Dependencies**: Librería estándar `net/http` (routing por método y path wildcards, Go ≥ 1.22) y `database/sql`; `modernc.org/sqlite` como driver SQLite puro Go (sin CGO). Testing: paquete `testing` estándar (TDD) y Godog para BDD. Sin dependencias nuevas.

**Storage**: SQLite embebido vía `modernc.org/sqlite` (driver `"sqlite"`), tablas `projects`, `users` y `project_members`.

**Testing**: `go test ./...` para dominio, servicio, repositorio y handler; escenarios BDD con Godog en `/features`.

**Target Platform**: Servidor local (Windows/Linux) que expone una API REST/JSON. El frontend React consume los endpoints; su UI queda fuera del alcance de esta iteración backend.

**Project Type**: Aplicación web (backend Go + frontend React), con estructura hexagonal.

**Performance Goals**: Aplicación local mono-usuario; crear, editar y asignar deben percibirse inmediatos (< 1 s de extremo a extremo). Sin objetivos de throughput.

**Constraints**: Sin CGO; la dependencia apunta siempre hacia el dominio; no incorporar dependencias más allá de las ya presentes (`modernc.org/sqlite`, Godog). Los identificadores de proyecto e integrante son `int64` (AUTOINCREMENT).

**Scale/Scope**: MVP de gestión de proyectos; decenas de proyectos, un equipo pequeño (decenas de integrantes) por proyecto.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

La constitución del proyecto (`.specify/memory/constitution.md`, v1.0.0) está ratificada y
sus principios (Test-First, SDD, BDD, Clean/Hexagonal, Clean Code/SOLID/KISS/YAGNI/DRY,
stack obligatorio) actúan como gates formales de este plan.

| Gate | Fuente | Estado |
|------|--------|--------|
| Dependencia apunta hacia adentro (domain no conoce SQLite/HTTP) | Arquitectura hexagonal / DIP | PASS |
| Un service por caso de uso (Crear, Editar, Asignar); services dependen de interfaces | SRP / DIP | PASS |
| Validaciones de proyecto en un único lugar, reutilizadas por alta y edición | DRY | PASS |
| Sin dependencias nuevas (REST con stdlib, IDs `int64` autoincrementales) | KISS / YAGNI | PASS |
| Ciclo TDD RED→GREEN→REFACTOR en dominio/servicio | TDD obligatorio | PASS |
| Escenario BDD por criterio de aceptación | BDD | PASS |
| Respeto del scaffolding de HU-05 (IDs `int64`) | Contrato entre historias | PASS |
| Sin edición de equipo, borrado de proyecto ni estados de ciclo de vida | YAGNI | PASS |
| Constitución formal ratificada | Gobernanza | PASS |

**Re-evaluación post-diseño (Phase 1)**: los gates siguen en PASS.

## Project Structure

### Documentation (this feature)

```text
specs/002-hu-04-crear-proyecto-equipo/
├── plan.md              # Este archivo
├── research.md          # Fase 0: decisiones técnicas
├── data-model.md        # Fase 1: entidades Project, User, Role, Membership
├── quickstart.md        # Fase 1: guía de validación end-to-end
├── contracts/
│   └── openapi.yaml     # Fase 1: contrato REST de proyectos e integrantes
├── checklists/
│   └── requirements.md  # Calidad de la spec (ya existe)
└── tasks.md             # Fase 2: tareas (generado por /speckit.tasks)
```

### Source Code (repository root)

```text
cmd/
└── api/
    └── main.go                        # wiring: repos→services→handlers; registra POST/GET/PUT /projects y /projects/{id}/members

internal/
├── domain/
│   ├── project.go                     # Project (existente) + fechas y constructor/edición validados
│   ├── project_test.go                # TDD de invariantes de Project (alta y edición)
│   ├── member.go                       # User, Role (Scrum Master/Product Builder), Membership
│   ├── member_test.go                  # TDD de validaciones de integrante y rol
│   ├── errors.go                       # ValidationError (existente, se reutiliza)
│   └── doc.go
├── service/
│   ├── crear_proyecto.go               # CrearProyecto (crea proyecto + asigna creador como Scrum Master)
│   ├── crear_proyecto_test.go
│   ├── obtener_proyecto.go             # ObtenerProyecto (lectura de soporte para verificar persistencia)
│   ├── obtener_proyecto_test.go
│   ├── editar_proyecto.go              # EditarProyecto (reutiliza validaciones del alta)
│   ├── editar_proyecto_test.go
│   ├── asignar_integrante.go           # AsignarIntegrante (find-or-create usuario + membership)
│   ├── asignar_integrante_test.go
│   └── doc.go
├── repository/
│   ├── project_repository.go           # interfaz (existente, extendida con Actualizar/AgregarIntegrante/ListarIntegrantes)
│   ├── user_repository.go              # interfaz UserRepository (find-or-create por nombre normalizado)
│   ├── sqlite_project.go               # adaptador SQLite de proyectos/memberships
│   ├── sqlite_user.go                  # adaptador SQLite de usuarios
│   ├── sqlite_project_test.go          # integración en :memory:
│   ├── sqlite_user_test.go             # integración en :memory:
│   ├── migrate.go                      # (existente) agrega tablas projects/users/project_members
│   └── doc.go
└── http/
    ├── project_handler.go              # POST /projects, GET /projects/{id}, PUT /projects/{id}, POST/GET /projects/{id}/members
    ├── project_dto.go                  # request/response DTOs
    ├── project_handler_test.go         # tests con httptest
    └── doc.go

features/
├── creacion_proyecto_equipo.feature            # escenarios BDD (crear, editar, asignar, errores)
└── steps_proyecto.go
```

**Structure Decision**: Se respeta la estructura ya definida en `AGENTS.md`. El paquete `internal/http` conserva su nombre actual (documentado en `internal/http/doc.go`). El fake en memoria de los repositorios vive como test double no exportado dentro de los tests de `internal/service` (no en producción), según YAGNI; `ProjectRepositoryEnMemoria` (no-test, en `internal/repository`) queda disponible para HU-05/HU-13. `domain.Project` y `ProjectRepository` ya existían (scaffolding para HU-05): se extendieron, con IDs `int64`.

## Complexity Tracking

> Sin violaciones a los principios del proyecto.

Notas de alcance (decisiones documentadas, no violaciones):

| Tema | Decisión | Justificación |
|------|----------|---------------|
| Redirección al panel / vistas (US1, US2, US3) | Esta iteración expone la API que **devuelve el proyecto creado/actualizado** (FR-004 reformulado); la navegación del SPA (pantalla principal, configuración, edición) se difiere al incremento de frontend | Coherente con HU-01 (backend-only); una redirección es responsabilidad de la UI |
| Identificador de proyecto `int64` | Se usa `id int64` asignado por SQLite con `AUTOINCREMENT` (`LastInsertId`) | Simple (KISS), suficiente para una app local y sin dependencias nuevas |
| Fechas del proyecto | HU-04 agrega `FechaInicio`/`FechaFin` opcionales y valida fin ≥ inicio al crear/editar; el cálculo del estado (Planificado/En curso/Finalizado) queda para HU-13 | La consigna pide editar fechas en HU-04; HU-13 (Should have) solo consulta el estado derivado |
| Auto-asociación del creador (FR-008) | Sin autenticación, el request de creación incluye el nombre del creador (`creador`) y el service lo resuelve/crea y lo vincula como Scrum Master | Garantiza SC-006 (≥ 1 Scrum Master) sin inventar login |
| Aplicación de permisos por rol (FR-012) | Se persiste la relación Proyecto–Usuario–Rol; la verificación en operaciones se difiere | Confirmado en `/speckit.clarify`; no hay identidad/sesión todavía |
| Borrado de proyecto y edición de equipo | Fuera de alcance | YAGNI; la spec solo pide crear, asignar y editar datos del proyecto |
| `GET /projects/{id}` y `GET /projects/{id}/members` | Se incluyen como lectura de soporte para verificar persistencia (SC-002, SC-009) | Sin lectura no hay forma de comprobar la recarga; alcance mínimo |
| No versionado ni historial (FR-023) | `UPDATE` simple, prevalece el último guardado válido | Requisito explícito de "edición simple" |
