---

description: "Task list template for feature implementation (override del proyecto: TDD con commits por fase)"
---

# Tasks: [FEATURE NAME]

**Input**: Design documents from `/specs/[###-feature-name]/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: OBLIGATORIOS. La constitución del proyecto (Principio I, NON-NEGOTIABLE)
exige TDD estricto con un commit por fase. No son opcionales en este proyecto.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

<!--
  ============================================================================
  REGLAS DE GENERACIÓN (para /speckit.tasks) — OBLIGATORIAS EN ESTE PROYECTO

  1. Todo comportamiento con lógica (reglas de negocio, cálculos, validaciones,
     casos de uso, repositorios, handlers) se genera como una secuencia de
     tareas separadas, en este orden exacto:
       a. RED     — escribir el test; correrlo y confirmar que FALLA.
       b. COMMIT  — `RED: <qué se testea> (falla)` (solo los archivos de test).
       c. GREEN   — implementación mínima; correr `go test ./...` en verde.
       d. COMMIT  — `GREEN: implementar <qué>, test en verde`.
       e. REFACTOR (opcional, solo si hay algo que limpiar) + COMMIT
          `REFACTOR: <qué se limpió>`, con los tests todavía en verde.
     Los COMMIT son tareas propias del tasks.md, con su propio ID: no se
     marcan [X] hasta haber corrido `git commit`.
  2. Tareas sin lógica (setup, migraciones, cableado de rutas, puertos/
     interfaces, documentación, escenarios .feature) llevan UNA tarea de
     commit al final con mensaje descriptivo, sin prefijo RED/GREEN.
  3. Nunca agrupar test e implementación en el mismo commit.
  4. Ninguna tarea hace `git push` ni merge: eso lo decide el equipo vía PR
     (merge commit, nunca squash).
  5. Rutas reales del proyecto (ver plan.md): internal/domain, internal/service,
     internal/repository, internal/http, cmd/api, features/.

  Las tareas de abajo son EJEMPLOS del formato; /speckit.tasks DEBE
  reemplazarlas por las reales de la historia.
  ============================================================================
-->

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

- [ ] T001 Agregar tabla/columna [x] a la migración en internal/repository/migrate.go
- [ ] T002 COMMIT `Agregar [x] a la migración`

---

## Phase 2: Foundational (Blocking Prerequisites)

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T003 RED: test de [Entidad] (invariantes del Factory) en internal/domain/[entidad]_test.go — confirmar que FALLA
- [ ] T004 COMMIT `RED: agregar tests de [Entidad] (falla)`
- [ ] T005 GREEN: implementar [Entidad] + New[Entidad] en internal/domain/[entidad].go — tests en verde
- [ ] T006 COMMIT `GREEN: implementar [Entidad], test en verde`
- [ ] T007 Definir puerto [Entidad]Repository en internal/repository/[entidad]_repository.go
- [ ] T008 COMMIT `Agregar puerto [Entidad]Repository`

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - [Title] (Priority: P1) 🎯 MVP

**Goal**: [Brief description of what this story delivers]

**Independent Test**: [How to verify this story works on its own]

- [ ] T009 [US1] RED: test del caso de uso [CasoDeUso] — [comportamiento] en internal/service/[caso]_test.go — confirmar que FALLA
- [ ] T010 [US1] COMMIT `RED: agregar tests de [CasoDeUso] (falla)`
- [ ] T011 [US1] GREEN: implementar [CasoDeUso] en internal/service/[caso].go — `go test ./...` en verde
- [ ] T012 [US1] COMMIT `GREEN: implementar [CasoDeUso], test en verde`
- [ ] T013 [US1] RED: test del handler [MÉTODO /ruta] en internal/http/[x]_handler_test.go — confirmar que FALLA
- [ ] T014 [US1] COMMIT `RED: agregar tests del handler [MÉTODO /ruta] (falla)`
- [ ] T015 [US1] GREEN: implementar el handler en internal/http/[x]_handler.go — tests en verde
- [ ] T016 [US1] COMMIT `GREEN: implementar handler [MÉTODO /ruta], test en verde`
- [ ] T017 [US1] (solo si hay algo que limpiar) REFACTOR de [qué] + COMMIT `REFACTOR: [qué se limpió]`
- [ ] T018 [US1] Cablear [MÉTODO /ruta] en cmd/api/main.go
- [ ] T019 [US1] COMMIT `Cablear [MÉTODO /ruta] en cmd/api`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

[Una fase por user story, mismo patrón RED → COMMIT → GREEN → COMMIT]

---

## Phase N: Polish & Cross-Cutting Concerns

- [ ] TXXX Step definitions de Godog para features/[historia].feature — `go test ./features/` en verde
- [ ] TXXX COMMIT `Agregar pasos BDD de [historia]`
- [ ] TXXX `go build ./...`, `go vet ./...` y `go test ./...` sin errores
- [ ] TXXX Run quickstart.md validation

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **User Stories (Phase 3+)** → **Polish**
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- User stories en orden de prioridad (P1 → P2 → P3); cada una independientemente testeable.

## Implementation Strategy

1. Setup + Foundational → fundación lista.
2. User Story 1 → validar de forma independiente (MVP).
3. Agregar las siguientes de a una, cada una dejando el sistema funcionando.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas.
- [P] tasks = different files, no dependencies.
