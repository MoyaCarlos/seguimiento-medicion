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
2. **Servicio** después: suma el `ProjectRepository` a `NewIniciarSprint`. El cambio de firma toca código de test (tests existentes y pasos BDD) y de producción (servicio y cableado). Para no mezclar test e implementación en un commit (regla 3 de la plantilla, Principio I), **todo el código de test va en el `RED:`** (que no compila: es su falla) y **el `GREEN:` solo toca `internal/service/iniciar_sprint.go` y `cmd/api/main.go`**, con lo que el árbol vuelve a compilar y queda en verde.
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

- [x] T001 [US1][US2][US3] RED: en `internal/domain/project_test.go`, agregar `TestProject_AdmiteSprint` como prueba de tabla (Arrange-Act-Assert) sobre `func (p Project) AdmiteSprint(inicio, fin time.Time) error`, con fechas `time.Date(..., 0, 0, 0, 0, time.UTC)`. Casos, con proyecto de inicio 2026-10-31 y fin 2026-11-16 salvo que se indique otra cosa:
  - proyecto sin fechas (`FechaInicio` y `FechaFin` nil), Sprint 2026-01-01 a 2027-01-01 → `nil` (US3);
  - solo `FechaInicio` 2026-10-31: Sprint 2026-10-30 a 2026-11-06 → `ValidationError` en `fecha_inicio`; Sprint 2026-10-31 a 2027-03-01 → `nil`; Sprint 2026-12-01 a 2027-03-01 → `nil` (US1, US3);
  - solo `FechaFin` 2026-11-16: Sprint 2026-01-01 a 2026-01-15 → `nil`; Sprint 2026-11-10 a 2026-11-16 → `nil`; Sprint 2026-11-10 a 2026-11-20 → `ValidationError` en `fecha_fin` (US2, US3);
  - ambas fechas: Sprint 2026-11-01 a 2026-11-10 → `nil`; Sprint 2026-10-31 a 2026-11-16 (los dos bordes) → `nil`; Sprint 2026-10-30 a 2026-11-06 → `fecha_inicio`; Sprint 2026-11-10 a 2026-11-17 → `fecha_fin`; Sprint 2026-10-01 a 2026-12-01 (las dos afuera) → `fecha_inicio`.
  - En los casos de error, verificar el mensaje exacto: `el Sprint no puede empezar antes del inicio del proyecto (2026-10-31)` y `el Sprint no puede terminar después del fin del proyecto (2026-11-16)` ([research.md §6](research.md)).
  - Correr `go test ./internal/domain/` y confirmar que FALLA (no compila: `AdmiteSprint` no existe).
- [x] T002 COMMIT `RED: agregar tests de Project.AdmiteSprint, plazo del proyecto (falla)` (solo `internal/domain/project_test.go`)
- [x] T003 [US1][US2][US3] GREEN: en `internal/domain/project.go`, implementar `AdmiteSprint`: si `p.FechaInicio != nil && inicio.Before(*p.FechaInicio)` → `ValidationError{Campo: "fecha_inicio", Mensaje: "el Sprint no puede empezar antes del inicio del proyecto (" + p.FechaInicio.Format(time.DateOnly) + ")"}`; después, si `p.FechaFin != nil && fin.After(*p.FechaFin)` → `ValidationError{Campo: "fecha_fin", Mensaje: "el Sprint no puede terminar después del fin del proyecto (" + p.FechaFin.Format(time.DateOnly) + ")"}`; si no, `nil`. Bordes incluidos; comparación directa con `Before`/`After`, sin normalizar ([research.md §3](research.md)). Comentario de doc breve que cite #168. `go test ./...` en verde.
- [x] T004 COMMIT `GREEN: implementar Project.AdmiteSprint, test en verde`
- [x] T005 (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>` — no hizo falta: dos condiciones sin duplicación. Nota: el test reutiliza el helper `dia` de `registro_horas_test.go` (mismo paquete).

**Checkpoint**: la regla existe y está probada; nadie la llama todavía.

---

## Phase 3: `IniciarSprint` aplica el plazo del proyecto (US1, US2, US3) 🎯 MVP

**Goal**: `POST /sprints/{id}/iniciar` rechaza con 400 las fechas fuera del plazo y deja el Sprint Pendiente ([contracts/iniciar-sprint-400.yaml](contracts/iniciar-sprint-400.yaml)).

**Independent Test**: `go test ./internal/service/ -run IniciarSprint -v`.

- [x] T006 [US1][US2][US3] RED, con **todo** el código de test que pide la nueva firma `NewIniciarSprint(sprints, proyectos)`:
  - En `internal/service/iniciar_sprint_test.go`, agregar el helper `sembrarProyectoConPlazo(t *testing.T, repo repository.ProjectRepository, inicio, fin *time.Time) int64`, que guarda `domain.NewProject("Software Metrics & Estimation", "", inicio, fin)` y devuelve su ID. **No** se llama `sembrarProyecto`: ese nombre ya existe en `internal/service/crear_sprint_test.go:12` (sin fechas, devuelve `domain.Project`) y se reutiliza para los casos sin fechas.
  - Casos nuevos en `internal/service/iniciar_sprint_test.go`, con un `repository.NewProjectRepositoryEnMemoria()` y un proyecto del 2026-10-31 al 2026-11-16:
    - US1: Sprint 2026-10-30 a 2026-11-06 → `ValidationError` en `fecha_inicio` (el caso del beta testing de Render) y, releyendo el repositorio, el Sprint sigue `Pendiente` con `SprintGoal` vacío y fechas cero (FR-007);
    - US1: Sprint 2026-10-31 a 2026-11-06 → `Activo`;
    - US2: Sprint 2026-11-10 a 2026-11-17 → `ValidationError` en `fecha_fin` y sigue `Pendiente`;
    - US2: Sprint 2026-11-10 a 2026-11-16 → `Activo`;
    - US3: proyecto sin fechas, Sprint 2026-01-01 a 2026-01-15 → `Activo`;
    - orden de errores (FR-006): Goal vacío con fechas fuera del plazo → `ValidationError` en `sprint_goal`; fin no posterior al inicio con fechas fuera del plazo (2026-10-30 a 2026-10-29) → `ValidationError` en `fecha_fin` por la regla existente; otro Sprint `Activo` en el proyecto con fechas fuera del plazo → `domain.ErrYaHaySprintActivo`.
  - Ajuste de Arrange de los tests existentes (ninguna aserción cambia; los proyectos se dan de alta **sin fechas**, así que el resultado esperado es el mismo):
    - `internal/service/iniciar_sprint_test.go`: los casos viejos crean un `ProjectRepositoryEnMemoria`, siembran el proyecto con `sembrarProyecto` y usan su ID en lugar de `1`/`2` (el caso "otro proyecto" siembra dos).
    - `internal/service/cerrar_sprint_test.go:33`: `sembrarSprintActivo(t, repo, proyectoID)` **conserva su firma**, porque lo llaman `cerrar_sprint_test.go`, `cambiar_estado_historia_test.go` (con `1`) y `asignar_historia_sprint_test.go` (con `proyectoID` y con `2`, el caso de otro proyecto). Pasa a usar `NewIniciarSprint(repo, proyectosSinFechas{})`, con un fake nuevo en el mismo archivo: `type proyectosSinFechas struct{ repository.ProjectRepository }` cuyo único método propio es `ObtenerPorID(_ context.Context, id int64) (domain.Project, error)`, que devuelve `domain.Project{ID: id}` (sin fechas, así que admite cualquier Sprint). Los demás métodos de la interfaz embebida no se usan en estos tests.
    - `internal/service/crear_sprint_test.go:81`: pasar el `proyectos` que el test ya crea: `NewIniciarSprint(sprints, proyectos)`.
    - `internal/http/sprint_handler_test.go`: `nuevoSprintHandlerConProyectos` pasa su `proyectos` a `service.NewIniciarSprint(sprints, proyectos)`, y `nuevoSprintHandlerDePrueba` siembra un proyecto sin fechas con `sembrarProyectoParaSprints` antes de devolver, porque `sembrarPendiente` crea los Sprints con `domain.NewSprint(1)` (el primer proyecto del repositorio en memoria tiene ID 1).
    - Pasos BDD: pasar un `repository.NewSQLiteProjectRepository(p.base.db)` a `NewIniciarSprint` en `features/steps_asignacion.go:47`, `features/steps_defectos.go:51` (ya tiene `proyectos` en la línea 45), `features/steps_estado_historia.go:42`, `features/steps_estado_proyecto.go:38` y `features/steps_sprints.go:45`.
  - Correr `go test ./...` y confirmar que FALLA (no compila: `NewIniciarSprint` todavía recibe un solo argumento).
- [x] T007 COMMIT `RED: agregar tests de IniciarSprint con el plazo del proyecto (falla)` (solo archivos `_test.go` y `features/steps_*.go`)
- [x] T008 [US1][US2][US3] GREEN, solo código de producción:
  - `internal/service/iniciar_sprint.go`: campo `proyectos repository.ProjectRepository`; `NewIniciarSprint(sprints repository.SprintRepository, proyectos repository.ProjectRepository)`. En `Ejecutar`, después de `sprint.Iniciar(...)` y antes de `s.sprints.Actualizar(...)`: `proyecto, err := s.proyectos.ObtenerPorID(ctx, sprint.ProyectoID)` (propagar el error) y `if err := proyecto.AdmiteSprint(sprint.FechaInicio, sprint.FechaFin); err != nil { return domain.Sprint{}, err }`. Actualizar el comentario del tipo ("…y que sus fechas caigan dentro del plazo del proyecto, #168").
  - Cableado `cmd/api/main.go:57`: `service.NewIniciarSprint(sprintRepo, proyectos)` (la variable `proyectos` ya existe en la línea 34).
  - `gofmt -l .` sin salida y `go test ./...` en verde (incluye Godog). Si algún test existente queda en rojo por un proyecto no sembrado, el ajuste es de test: volver a T006 con un commit `test: sembrar el proyecto en <test>` aparte, nunca dentro del `GREEN:`.
- [x] T009 COMMIT `GREEN: IniciarSprint rechaza fechas fuera del plazo del proyecto, test en verde` (solo `internal/service/iniciar_sprint.go` y `cmd/api/main.go`)
- [x] T010 (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>` — no hizo falta.

**Notas de la fase 3 (desvíos respecto de lo planeado):**

- Los tests viejos de `internal/service/iniciar_sprint_test.go` usan el fake `proyectosSinFechas` (definido en ese archivo, no en `cerrar_sprint_test.go`) en lugar de sembrar proyectos: mismo resultado, sin tocar los IDs `1`/`2`. Los casos nuevos sí usan `ProjectRepositoryEnMemoria` con `sembrarProyectoConPlazo`.
- El `GREEN:` dejó en rojo 3 escenarios de `features/estado_proyecto.feature` (HU-13): sus pasos cargan fechas al proyecto y luego iniciaban siempre un Sprint de hoy a hoy+14, fuera de ese plazo. Contradice research.md §7 ("todos los proyectos de `/features` se crean sin fechas"), que pasaba por alto los `UPDATE` de fechas de `steps_estado_proyecto.go`. Se corrigió en un commit `test:` aparte, antes del `GREEN:` (`e408ad1`): `crearActivo` recorta el Sprint al plazo del proyecto. Lo que verifica cada escenario no cambia.

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
- Dentro de T006, los ajustes de los distintos archivos de pasos y de tests son independientes entre sí, pero van en un solo commit `RED:`.

## Implementation Strategy

1. Phase 2 + Phase 3 = MVP: la API ya rechaza los Sprints fuera del plazo (cierra el defecto de Render).
2. Phase 4 agrega la evidencia BDD exigida por la Definition of Done.
3. Phase 5 y 6 dejan la documentación y el cierre.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas; el PR cierra con `Closes #168` y se mergea con "Create a merge commit".
- Sin cambios de frontend: `FormularioInicio.tsx` ya muestra el error por campo ([research.md §8](research.md)).
