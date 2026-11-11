---

description: "Tareas de la spec 043: Sprints dentro del plazo del proyecto (#168)"
---

# Tasks: Sprints dentro del plazo del proyecto (#168)

**Input**: Design documents from `/specs/043-sprint-dentro-del-proyecto/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/)

**Tests**: OBLIGATORIOS. La constitución del proyecto (Principio I, NON-NEGOTIABLE) exige TDD estricto con un commit por fase.

**Organization**: por capa y no por historia, porque las tres historias de la spec (US1 inicio antes del proyecto, US2 fin después del proyecto, US3 fechas incompletas) son **una misma regla** (`Project.AdmiteSprint`) con una misma llamada en el servicio. Separarlas en fases propias obligaría a escribir tests que pasan apenas se escriben (no habría `RED:` real). Cada tarea indica qué historias cubre.

## Orden de las fases (y por qué)

1. **Dominio** primero: la regla no depende de nada y el servicio la necesita.
2. **Servicio** después: suma el `ProjectRepository` a `NewIniciarSprint`. El cambio de firma rompe la compilación de todos los llamadores, por eso el `GREEN:` ajusta en el mismo commit el cableado, los pasos BDD y los tests existentes (si no, el árbol queda sin compilar entre commits).
3. **Godog al final**: sin pasos definidos, un `.feature` nuevo rompe `TestFeatures` (`Strict: true`), así que el `.feature` y sus pasos van juntos cuando la regla ya está implementada (mismo criterio que la spec 040).
4. **`MODELO-DE-DATOS.md`** y cierre.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)

---

## Phase 1: Setup

Sin tareas: la rama `feature/168-sprint-dentro-del-proyecto` ya existe, sale de `main` y no hay esquema ni dependencias nuevas.

---

## Phase 2: Regla de dominio `Project.AdmiteSprint` (US1, US2, US3)

**Goal**: el proyecto decide si un rango de fechas de Sprint cabe en su plazo ([data-model.md](data-model.md)).

**Independent Test**: `go test ./internal/domain/ -run AdmiteSprint -v`.

- [ ] T001 [US1][US2][US3] RED: en `internal/domain/project_test.go`, agregar `TestProject_AdmiteSprint` como prueba de tabla (Arrange-Act-Assert) sobre `func (p Project) AdmiteSprint(inicio, fin time.Time) error`, con fechas `time.Date(..., 0, 0, 0, 0, time.UTC)`. Casos, con proyecto de inicio 2026-10-31 y fin 2026-11-16 salvo que se indique otra cosa:
  - proyecto sin fechas (`FechaInicio` y `FechaFin` nil), Sprint 2026-01-01 a 2027-01-01 → `nil` (US3);
  - solo `FechaInicio` 2026-10-31: Sprint 2026-10-30 a 2026-11-06 → `ValidationError` en `fecha_inicio`; Sprint 2026-10-31 a 2027-03-01 → `nil`; Sprint 2026-12-01 a 2027-03-01 → `nil` (US1, US3);
  - solo `FechaFin` 2026-11-16: Sprint 2026-01-01 a 2026-01-15 → `nil`; Sprint 2026-11-10 a 2026-11-16 → `nil`; Sprint 2026-11-10 a 2026-11-20 → `ValidationError` en `fecha_fin` (US2, US3);
  - ambas fechas: Sprint 2026-11-01 a 2026-11-10 → `nil`; Sprint 2026-10-31 a 2026-11-16 (los dos bordes) → `nil`; Sprint 2026-10-30 a 2026-11-06 → `fecha_inicio`; Sprint 2026-11-10 a 2026-11-17 → `fecha_fin`; Sprint 2026-10-01 a 2026-12-01 (las dos afuera) → `fecha_inicio`.
  - En los casos de error, verificar el mensaje exacto: `el Sprint no puede empezar antes del inicio del proyecto (2026-10-31)` y `el Sprint no puede terminar después del fin del proyecto (2026-11-16)` ([research.md §6](research.md)).
  - Correr `go test ./internal/domain/` y confirmar que FALLA (no compila: `AdmiteSprint` no existe).
- [ ] T002 COMMIT `RED: agregar tests de Project.AdmiteSprint, plazo del proyecto (falla)` (solo `internal/domain/project_test.go`)
- [ ] T003 [US1][US2][US3] GREEN: en `internal/domain/project.go`, implementar `AdmiteSprint`: si `p.FechaInicio != nil && inicio.Before(*p.FechaInicio)` → `ValidationError{Campo: "fecha_inicio", Mensaje: "el Sprint no puede empezar antes del inicio del proyecto (" + p.FechaInicio.Format(time.DateOnly) + ")"}`; después, si `p.FechaFin != nil && fin.After(*p.FechaFin)` → `ValidationError{Campo: "fecha_fin", Mensaje: "el Sprint no puede terminar después del fin del proyecto (" + p.FechaFin.Format(time.DateOnly) + ")"}`; si no, `nil`. Bordes incluidos; comparación directa con `Before`/`After`, sin normalizar ([research.md §3](research.md)). Comentario de doc breve que cite #168. `go test ./...` en verde.
- [ ] T004 COMMIT `GREEN: implementar Project.AdmiteSprint, test en verde`
- [ ] T005 (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: la regla existe y está probada; nadie la llama todavía.

---

## Phase 3: `IniciarSprint` aplica el plazo del proyecto (US1, US2, US3) 🎯 MVP

**Goal**: `POST /sprints/{id}/iniciar` rechaza con 400 las fechas fuera del plazo y deja el Sprint Pendiente ([contracts/iniciar-sprint-400.yaml](contracts/iniciar-sprint-400.yaml)).

**Independent Test**: `go test ./internal/service/ -run IniciarSprint -v`.

- [ ] T006 [US1][US2][US3] RED: en `internal/service/iniciar_sprint_test.go`:
  - Agregar un helper `sembrarProyecto(t, proyectos *repository.ProjectRepositoryEnMemoria, inicio, fin *time.Time) int64` que guarda un `domain.NewProject("Proyecto", "", inicio, fin)` y devuelve su ID.
  - Pasar a `NewIniciarSprint(repo, proyectos)` en **los tests nuevos** (los existentes se ajustan en T008).
  - Casos nuevos, con proyecto del 2026-10-31 al 2026-11-16:
    - US1: Sprint 2026-10-30 a 2026-11-06 → `ValidationError` en `fecha_inicio` (el caso del beta testing de Render) y, releyendo el repositorio, el Sprint sigue `Pendiente` con `SprintGoal` vacío y fechas cero (FR-007);
    - US1: Sprint 2026-10-31 a 2026-11-06 → `Activo`;
    - US2: Sprint 2026-11-10 a 2026-11-17 → `ValidationError` en `fecha_fin` y sigue `Pendiente`;
    - US2: Sprint 2026-11-10 a 2026-11-16 → `Activo`;
    - US3: proyecto sin fechas, Sprint 2026-01-01 a 2026-01-15 → `Activo`;
    - orden de errores (FR-006): Goal vacío con fechas fuera del plazo → `ValidationError` en `sprint_goal`; otro Sprint `Activo` en el proyecto con fechas fuera del plazo → `domain.ErrYaHaySprintActivo`.
  - Correr `go test ./internal/service/` y confirmar que FALLA (no compila: `NewIniciarSprint` todavía recibe un solo argumento).
- [ ] T007 COMMIT `RED: agregar tests de IniciarSprint con el plazo del proyecto (falla)` (solo `internal/service/iniciar_sprint_test.go`)
- [ ] T008 [US1][US2][US3] GREEN, en un solo commit para que el árbol compile:
  - `internal/service/iniciar_sprint.go`: campo `proyectos repository.ProjectRepository`; `NewIniciarSprint(sprints repository.SprintRepository, proyectos repository.ProjectRepository)`. En `Ejecutar`, después de `sprint.Iniciar(...)` y antes de `s.sprints.Actualizar(...)`: `proyecto, err := s.proyectos.ObtenerPorID(ctx, sprint.ProyectoID)` (propagar el error) y `if err := proyecto.AdmiteSprint(sprint.FechaInicio, sprint.FechaFin); err != nil { return domain.Sprint{}, err }`. Actualizar el comentario del tipo ("…y que sus fechas caigan dentro del plazo del proyecto, #168").
  - Cableado `cmd/api/main.go`: `service.NewIniciarSprint(sprintRepo, projectRepo)` (usar el nombre de la variable del `ProjectRepository` que ya existe en ese archivo).
  - Pasos BDD existentes: pasar un `repository.NewSQLiteProjectRepository(<db del escenario>)` a `NewIniciarSprint` en `features/steps_asignacion.go`, `features/steps_defectos.go`, `features/steps_estado_historia.go`, `features/steps_estado_proyecto.go` y `features/steps_sprints.go`.
  - Tests existentes: `internal/service/iniciar_sprint_test.go` (los casos viejos siembran el proyecto sin fechas con `sembrarProyecto` y usan su ID en lugar de `1`/`2`), `internal/service/cerrar_sprint_test.go`, `internal/service/crear_sprint_test.go` y `internal/http/sprint_handler_test.go` (pasarle el `proyectos` en memoria que ya crea `nuevoSprintHandlerConProyectos`; si algún test inicia un Sprint de un proyecto que no dio de alta, darlo de alta sin fechas). Es ajuste de Arrange: ninguna aserción existente cambia.
  - `gofmt -l .` sin salida y `go test ./...` en verde (incluye Godog).
- [ ] T009 COMMIT `GREEN: IniciarSprint rechaza fechas fuera del plazo del proyecto, test en verde`
- [ ] T010 (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: la regla rige en la API. El caso de Render (proyecto desde 31/10, Sprint desde 30/10) da 400 en `fecha_inicio`.

---

## Phase 4: Escenarios BDD (US1, US2, US3)

**Goal**: los escenarios Given-When-Then de la spec automatizados con Godog (FR-009).

**Independent Test**: `go test ./features/ -v`.

- [ ] T011 [US1][US2][US3] Crear `features/sprint_dentro_del_proyecto.feature` en Gherkin en español (`# language: es`), tags `@HU-05 @168`, Característica con el Como/Quiero/Para de la US1 de la spec. Escenarios, uno por criterio de aceptación de la spec (sin inventar casos):
  - US1.1 error: proyecto que empieza el 2026-10-31; iniciar del 2026-10-30 al 2026-11-06 → rechazo en `fecha_inicio` y el Sprint sigue Pendiente;
  - US1.2 límite: iniciar del 2026-10-31 al 2026-11-06 → Activo;
  - US2.1 error: proyecto que termina el 2026-11-16; iniciar del 2026-11-10 al 2026-11-17 → rechazo en `fecha_fin` y sigue Pendiente;
  - US2.2 límite: iniciar del 2026-11-10 al 2026-11-16 → Activo;
  - US3.1 alternativo: proyecto sin fechas → Activo;
  - US3.2 alternativo: proyecto solo con inicio 2026-10-31, iniciar del 2026-12-01 al 2027-03-01 → Activo;
  - US3.3 alternativo: proyecto solo con fin 2026-11-16, iniciar del 2026-01-01 al 2026-01-15 → Activo;
  - US3.4 error: proyecto solo con fin 2026-11-16, iniciar del 2026-11-10 al 2026-11-20 → rechazo en `fecha_fin`.
  - Las 4 categorías (normal, alternativo, límite, error) quedan cubiertas; el "normal" es un Sprint dentro del plazo, cubierto por los límites, que son el caso más exigente.
  - Usar `Esquema del escenario` con `Ejemplos` si agrupa casos del mismo paso sin perder la trazabilidad a los criterios.
- [ ] T012 [US1][US2][US3] Pasos en `features/steps_sprint_dentro_del_proyecto.go` (patrón de `features/steps_sprints.go`): un proyecto con fechas opcionales creado con `repository.NewSQLiteProjectRepository` sobre la base del escenario, un Sprint Pendiente, el inicio vía `service.NewIniciarSprint(sprints, proyectos)` y las verificaciones de estado y del `ValidationError.Campo`. Registrar `InitializeScenarioSprintDentroDelProyecto` en `features/features_test.go` con el comentario `// #168`. `go test ./features/` en verde.
- [ ] T013 COMMIT `Agregar escenarios BDD de Sprints dentro del plazo del proyecto (Refs #168)`

---

## Phase 5: Documentación

- [ ] T014 [P] Agregar en `MODELO-DE-DATOS.md` §5 ("Dónde vive cada regla") la fila definida en [data-model.md](data-model.md) (regla en el dominio, sin restricción de base, menciona #168 y #302), después de la fila "Fin posterior a inicio (proyecto y Sprint)".
- [ ] T015 COMMIT `docs: regla del plazo del proyecto en MODELO-DE-DATOS.md §5 (Refs #168)`

---

## Phase 6: Cierre

- [ ] T016 `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test -count=1 ./...` en verde.
- [ ] T017 Validar con [quickstart.md](quickstart.md) §2 contra la API levantada localmente (los tres casos de la tabla) y anotar el resultado en el PR.
- [ ] T018 Marcar las tareas hechas en este archivo + COMMIT `docs: marcar las tareas de la spec 043 como hechas (Refs #168)`

---

## Dependencies & Execution Order

- Phase 2 → Phase 3 → Phase 4 → Phase 5 → Phase 6. Phase 5 (T014) es independiente del código y puede hacerse en cualquier momento después del plan [P].
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- **Convivencia con #301 (HU-21)**: si se mergea antes que esta rama, traer `main` antes del PR y resolver según [research.md §5](research.md): firma única `NewIniciarSprint(sprints, proyectos repository.ProjectRepository)`, el cableado pasa a `internal/app/app.go` y el campo `proyectos` se usa para `exigirScrumMaster` y para `ObtenerPorID`.

## Parallel Opportunities

- T014 (documentación) en paralelo con las fases 2 a 4.
- Dentro de T008, los ajustes de los distintos archivos de pasos y de tests son independientes entre sí, pero van en un solo commit.

## Implementation Strategy

1. Phase 2 + Phase 3 = MVP: la API ya rechaza los Sprints fuera del plazo (cierra el defecto de Render).
2. Phase 4 agrega la evidencia BDD exigida por la Definition of Done.
3. Phase 5 y 6 dejan la documentación y el cierre.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas; el PR cierra con `Closes #168` y se mergea con "Create a merge commit".
- Sin cambios de frontend: `FormularioInicio.tsx` ya muestra el error por campo ([research.md §8](research.md)).
