# Análisis de consistencia — spec 006 (Issue #107)

**Comando**: `/speckit-analyze` (T001 de [tasks.md](tasks.md)) · **Fecha**: 2026-10-23
**Artefactos**: [spec.md](spec.md), [plan.md](plan.md), [tasks.md](tasks.md), contrastados con `.specify/memory/constitution.md` y con el código de `main`.

**Resultado**: 0 CRITICAL, 1 HIGH, 3 MEDIUM, 5 LOW. El HIGH (K1) frenó la implementación según T001 y lo resolvió el Scrum Master; todos los hallazgos quedaron resueltos en los artefactos antes de escribir código.

## Hallazgos

| ID | Categoría | Severidad | Ubicación | Resumen | Recomendación | Resolución |
|----|-----------|-----------|-----------|---------|---------------|------------|
| K1 | Constitución (III, BDD) | HIGH | spec.md SC-002, FR-009, US4; tasks.md | El Principio III exige escenario Godog para comportamiento observable. El 409 por FK es un cambio de la API y la base vieja un error visible, sin escenario; SC-002 congela los `.feature`. | (a) Excepción justificada en la spec y el PR, o (b) escenarios BDD nuevos. | **Opción (a), decisión del Scrum Master.** Agregada a *Clarifications*, a la sección *Proceso* de la spec ("Excepción al Principio III") y a la descripción del PR (T042): el 409 no se alcanza por el flujo normal y se prueba contra SQLite real en la fase 4; la base vieja es un error de arranque, no comportamiento de la API; SC-002 se mantiene. |
| I1 | Inconsistencia | MEDIUM | spec.md *Proceso* vs plan.md fase 3c, research R7, tasks T021–T022 | La spec decía "3a–3d renombre (un par de commits por repositorio)" y la 3c va como un solo commit `test:`. | Nombrar la excepción de la 3c en la spec. | Corregida la línea de *Proceso*: un par `RED:`/`GREEN:` por repositorio salvo la 3c, que va como un único `test:` (research R7). |
| U1 | Subespecificación | MEDIUM | spec.md CA-13; tasks T027–T029 | CA-13 prometía traducir "cada violación anterior" (CA-6 a CA-12); las tareas prueban por tipo y la FK compuesta no se puede provocar desde un adaptador. | Reescribir CA-13 por tipo de violación. | CA-13 reescrito: FK → `ErrIntegridadReferencial`, índice parcial → `ValidationError` con el mensaje del servicio, `CHECK`/`NOT NULL` → genérico; la FK compuesta de una historia en un Sprint ajeno se cubre en el test de esquema (CA-8). |
| C1 | Cobertura | MEDIUM | spec.md caso límite D9, FR-004; tasks T009 | Sin verificación de "sin `ON DELETE`" ni caso de borrado de un padre con hijos. | Chequear `on_delete`/`on_update` y sumar un `DELETE`. | T009: `TestEsquema_Estructura` verifica `on_delete` y `on_update` = `NO ACTION` en cada FK; `TestEsquema_Filas` suma `DELETE` de un proyecto con historias → 787. |
| I2 | Inconsistencia | LOW | data-model.md, quickstart.md paso 6 | Listaban cuatro correcciones al documento; FR-016a tiene cinco (falta la frase de §10 del rojo). | Agregarla. | Agregada en las dos listas. |
| I3 | Inconsistencia | LOW | research.md R5 vs plan.md fase 4 | R5 no decía que el uso de `errorDeEscritura` en `integrantes` es del REFACTOR. | Aclarar la fase. | R5 aclara: historias y Sprints en el `GREEN:`; `sqlite_project.go` en el `REFACTOR:` de la fase 4. |
| I4 | Inconsistencia | LOW | checklists/requirements.md | La nota seguía con "dos supuestos para confirmar". | Actualizar. | Nota actualizada: confirmados en el clarify y pasados a FR-006, FR-006a, FR-011 y FR-016a. |
| U2 | Subespecificación | LOW | tasks T021; sqlite_sprint_test.go | `TestSQLiteSprintRepository_ListarPorProyecto` guarda dos `Pendiente` en el proyecto 1 y choca con el índice parcial. | Nombrar el test y cómo se ajusta. | T021 lo nombra: dos proyectos reales y el primer Sprint del proyecto 1 iniciado antes de crear el segundo; las aserciones de orden no cambian. |
| C2 | Cobertura | LOW | spec.md US3.1, US4.4, FR-013 | Se apoyaban en tests existentes sin nombrarlos. | Nombrarlos en el chequeo final. | T039 lista los tests de handler, servicio y Godog que respaldan "respuestas iguales", "el 404 sigue" y los duplicados. |

## Cobertura (después de las resoluciones)

| Requisito | ¿Tarea? | Tareas | Notas |
|---|---|---|---|
| FR-001 | Sí | T007, T011, T037 | |
| FR-002 | Sí | T005–T008 | incluye versión desconocida |
| FR-003, FR-004 | Sí | T009–T012 | FR-004 verificado por C1 |
| FR-005 | Sí | T025, T037 | |
| FR-006, FR-006a | Sí | T035–T038 | |
| FR-007, FR-008 | Sí | T013–T026 | |
| FR-009 | Sí | T039, T040 | tests existentes nombrados (C2) |
| FR-010 a FR-014 | Sí | T027–T034, T039 | |
| FR-015 | Sí | T009 | |
| FR-016 | — | — | cumplido por data-model.md |
| FR-016a | Sí | T003, T004 | |
| FR-017 | Sí | cabecera de tasks.md | |
| CA-1 a CA-12 | Sí | T009–T012, T035, T040 | |
| CA-13 | Sí | T027–T029, T009 | reescrito (U1) |
| CA-14 | Sí | T004, T040 | |
| SC-001 a SC-004 | Sí | T005, T025, T009, T029 | |

## Constitución

- I (Test-First): cumple; ventana en rojo (GREEN de la 3a a la 3d) y fase 3c documentadas.
- II (SDD): cumple; este análisis corre antes de `/speckit-implement`.
- III (BDD): excepción registrada (K1).
- IV, V, VI: sin observaciones.

## Métricas

- Requisitos: 37 (19 FR, 14 CA, 4 SC) · Tareas: 42
- Cobertura: 100 %
- Hallazgos: 9 (0 CRITICAL, 1 HIGH, 3 MEDIUM, 5 LOW), todos resueltos

## Segunda corrida (T001, 2026-10-23)

Corrida de `/speckit-analyze` sobre los artefactos ya corregidos, contrastando además
cada archivo, función, consulta y mensaje que nombran las tareas con el código de `main`
(`3ac8df0`).

**Resultado**: 0 CRITICAL, 0 HIGH, 1 MEDIUM, 0 LOW. Las resoluciones de la primera corrida
siguen vigentes. No frena la implementación (T001 solo frena ante CRITICAL o HIGH).

| ID | Categoría | Severidad | Ubicación | Resumen | Recomendación |
|----|-----------|-----------|-----------|---------|---------------|
| U3 | Subespecificación | MEDIUM | tasks.md T025; features/steps_proyecto.go:143-145 y :202-204 | T025 nombra las líneas 143 y 203 pero solo lista los renombres `projects`, `name` y `users`. Esas dos consultas usan `project_members`, `project_id` y `role`, que T025 no menciona. | Al hacer T025, renombrarlas también según §4.3 (`integrantes`, `proyecto_id`, `rol`). No es una duda de modelo: el destino está definido en §4.3. Se avisa como desvío al cerrar la fase 3d. |

Verificado contra el código, sin hallazgos:

- Los mensajes del servicio (`crear_sprint.go:36`, `iniciar_sprint.go:44`) coinciden con FR-011 y con T003.
- `ErrIntegridadReferencial` ya existe (`errors.go:12`) y `sqlite_project.go:50` y `:136` ya traducen la FK a mano (lo que T033 reemplaza).
- `escribirError` (`sprint_handler.go:113`) es el único traductor HTTP y lo usan los tres handlers, así que el caso 409 de T031 cubre toda la API.
- Solo `steps_creacion.go` y `steps_proyecto.go` llaman a `Migrar`; `steps_estado_proyecto.go` es el único que escribe fechas de proyecto en RFC 3339 (T025).
- Los helpers `abrirBDDePrueba` (solo en `sqlite_project_test.go` y `sqlite_user_test.go`) y `nuevaBD` (solo en `sqlite_sprint_test.go`) se usan donde dicen T013, T017, T021 y T037.

**Métricas**: 37 requisitos, 42 tareas, cobertura 100 %, 0 ambigüedades, 0 duplicaciones, 0 CRITICAL.
