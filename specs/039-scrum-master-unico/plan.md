# Implementation Plan: Un único Scrum Master por proyecto (#236)

**Branch**: `feature/236-scrum-master-unico` | **Date**: 2026-11-08 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/039-scrum-master-unico/spec.md`

## Summary

Garantizar a lo sumo un Scrum Master por proyecto con **doble protección**: el service lo comprueba antes de insertar y responde un error de dominio nuevo (`ErrScrumMasterYaExiste`, 409); un índice único parcial en la base (migración 007) lo garantiza ante carreras y escrituras directas, y el repositorio traduce su violación al mismo error. La 007 corrige primero los proyectos que ya tienen varios Scrum Master (criterio de la clarificación del 2026-11-08) y recién después crea el índice. Todo en un solo PR junto con `MODELO-DE-DATOS.md`.

## Technical Context

**Language/Version**: Go 1.27 (backend). Sin cambios en el frontend.

**Primary Dependencies**: la biblioteca estándar y `modernc.org/sqlite`, ya en el proyecto. No se agrega ninguna dependencia.

**Storage**: SQLite. Una migración nueva (`007_scrum_master_unico.sql`, `go:embed`) con un `UPDATE` y un `CREATE UNIQUE INDEX` parcial. Sin tablas ni columnas nuevas. Ver [data-model.md](./data-model.md).

**Testing**: `go test` (dominio, service, repositorio, handler) y Godog en `/features`. Ciclo RED → GREEN → REFACTOR con un commit por fase.

**Target Platform**: Linux (servidor único, Docker en Render).

**Project Type**: aplicación web (backend Go + SPA React); este cambio toca solo el backend, el esquema y la documentación.

**Performance Goals**: no aplica (una consulta de integrantes extra por alta de integrante).

**Constraints**: la migración debe arrancar sola sobre cualquier base (sin intervención manual); una migración mergeada no se edita; cambio de esquema y `MODELO-DE-DATOS.md` en el mismo PR; no se rediseña el runner de migraciones.

**Scale/Scope**: 1 error de dominio, 1 chequeo en 1 service, 1 traducción de error en 1 repositorio (+ el repositorio en memoria), 1 línea en el mapeo HTTP, 1 migración, 1 documento.

## Constitution Check

| Principio | Cumplimiento |
|---|---|
| I. Test-First | Cumple. RED/GREEN por capa (dominio/service, repositorio, migración, handler) en commits propios; los tests existentes que contradicen la regla se corrigen en su RED (ver research.md §5). |
| II. SDD | Cumple. specify → clarify → plan → tasks → analyze antes de implementar. |
| III. BDD | Cumple. Escenarios Godog para crear, rechazar el segundo (409), agregar Product Builder y el índice por fuera del dominio. La corrección de la 007 se cubre con un test de migración (research.md §4). |
| IV. Hexagonal | Cumple. La regla vive en el service (con el error en `domain`); el repositorio solo traduce la violación; el dominio no conoce SQLite. |
| V. Clean/SOLID/KISS/YAGNI/DRY | Cumple. Sin abstracciones nuevas; el error se mapea con una línea en el `switch` existente. |
| VI. Stack | Cumple. Sin dependencias nuevas. |

Sin violaciones; no hay tabla de complejidad. Re-chequeo posterior al diseño: igual.

## Project Structure

### Documentation (this feature)

```text
specs/039-scrum-master-unico/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── members-409.yaml
├── checklists/requirements.md
└── spec.md
```

### Source Code (repository root)

```text
internal/domain/errors.go                                   # + ErrScrumMasterYaExiste
internal/service/asignar_integrante.go                      # chequeo previo (duplicado → 400, luego Scrum Master → 409)
internal/repository/sqlite_project.go                       # separar código 1555 (PK) de 2067 (índice)
internal/repository/project_repository_en_memoria.go       # imita el índice para los tests
internal/repository/migraciones/007_scrum_master_unico.sql  # UPDATE de corrección + índice, con el criterio comentado
internal/http/sprint_handler.go                             # escribirError: ErrScrumMasterYaExiste → 409
MODELO-DE-DATOS.md                                          # D4, D14, §4.3, §9 (mismo PR)
specs/002-hu-04-crear-proyecto-equipo/contracts/openapi.yaml # + respuesta 409 en POST /projects/{id}/members
features/creacion_proyecto_equipo.feature (+ steps)         # escenarios BDD nuevos
internal/repository/esquema_test.go, migrate_test.go        # esquema + corrección de la 007
internal/service/asignar_integrante_test.go y fakes_test.go # service y fake con el mismo orden de errores
```

**Structure Decision**: sin carpetas nuevas; todo entra en los paquetes existentes de la arquitectura hexagonal. Sin cambios en `web/`.

## Diseño de la protección doble

1. **Service** (`AsignarIntegrante.Ejecutar`), después de resolver el usuario y antes de insertar: lista los integrantes del proyecto (`ListarIntegrantes`, que ya usa `rolEnProyecto`).
   - Si el usuario ya es integrante → `ValidationError{integrante}` actual (400). **Gana siempre** (FR-006).
   - Si no, y el rol pedido es `scrum_master` y ya hay uno → `domain.ErrScrumMasterYaExiste` (409).
   - Si no, inserta.
2. **Repositorio SQLite** (`AgregarIntegrante`): ante una violación `UNIQUE`, distingue por código extendido (evidencia en research.md §1): `SQLITE_CONSTRAINT_PRIMARYKEY` (1555) → `ErrMiembroDuplicado`; `SQLITE_CONSTRAINT_UNIQUE` (2067) → `domain.ErrScrumMasterYaExiste`. Hoy ambos caen en "duplicado".
3. **HTTP**: `escribirError` suma `ErrScrumMasterYaExiste` al grupo de los 409 (mensaje en `mensaje`, sin `campo`).
4. **Creación del proyecto** no cambia: `GuardarConScrumMaster` inserta el proyecto y su Scrum Master en una transacción; el índice no se puede violar ahí porque el proyecto es nuevo.

## Fuera de alcance / límites declarados

Traspaso de rol y baja de integrantes (HU-16, #175). Frontend: el formulario de equipo ya muestra en el error general el mensaje de cualquier respuesta sin `campo` (`Configuracion.tsx`, `errores.ts`), así que el 409 se ve sin tocar nada (research.md §6).
