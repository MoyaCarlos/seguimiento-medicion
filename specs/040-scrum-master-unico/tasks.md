# Tasks: Un único Scrum Master por proyecto (#236)

**Input**: `specs/040-scrum-master-unico/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: plan.md, spec.md

**Tests**: OBLIGATORIOS (Principio I, NON-NEGOTIABLE). Cada comportamiento con lógica va como RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT si hay algo que limpiar). Los COMMIT no se marcan [X] hasta haber corrido `git commit`. Ninguna tarea hace `git push` ni merge.

**Organization**: por capa y no por historia, porque las historias comparten archivos (ver "Orden de las fases"). Cada fase deja `go test ./...` en verde al cerrar su GREEN.

## Orden de las fases (y por qué)

1. **Migración 007** primero, porque el repositorio (fase 3) y los Godog (fase 5) necesitan el índice en la base.
2. **Service** antes que el repositorio: el chequeo previo es el que garantiza el orden de errores de FR-006 (el 400 gana al 409), y no depende de la base.
3. **Repositorio** (separar 1555 de 2067) y el repositorio en memoria.
4. **Mapeo HTTP y contrato**, que necesitan el error del service.
5. **Godog al final**: sin pasos definidos, un `.feature` nuevo rompe `TestFeatures` (`Strict: true`), así que el `.feature` y sus pasos van juntos cuando todo ya está implementado.
6. **MODELO-DE-DATOS.md** en el mismo PR.
7. **Cierre**.

**Ajuste respecto del orden propuesto**: el test de esquema (`esquemaEsperado` con el índice y las filas del FR-014) va en la **fase 1**, no en la 6. `TestEsquema_Estructura` compara la lista exacta de índices de cada tabla (`esquema_test.go:233`); si el índice de la 007 llega sin actualizarlo, ese test queda en rojo al cerrar el GREEN de la fase 1. La fase 6 conserva solo el documento.

Formato: `[ID] [US] Descripción`. Rutas reales según plan.md.

## Phase 1: Migración 007 y esquema (US3, US4)

**Goal**: la 007 corrige los proyectos con varios Scrum Master y crea el índice único parcial.

**Independent Test**: base migrada hasta la 005 y sembrada; al aplicar la 007 queda a lo sumo un Scrum Master por proyecto y el índice rechaza un segundo.

- [x] T001 [US4] RED: test de la corrección en `internal/repository/migrate_test.go` (o `migracion_007_test.go`): armar un `fstest.MapFS` con las migraciones embebidas menos la 007, aplicarlo, sembrar proyectos con 0, 1, 2 y 3 Scrum Master y un empate de `creado_en`, aplicar `AplicarMigraciones` y verificar quién queda, que nadie se borra y que los demás pasan a `product_builder` — confirmar que FALLA
- [x] T002 [US3] RED: test de esquema en `internal/repository/esquema_test.go`: agregar a `esquemaEsperado["integrantes"]` el índice `integrantes_un_scrum_master` (columna `proyecto_id`, `WHERE rol = 'scrum_master'`) y las filas del FR-014 (segundo `scrum_master` por `INSERT` y por `UPDATE` rechazados con `UNIQUE`; varios `product_builder` y un `scrum_master` por proyecto distinto aceptados) — confirmar que FALLA
- [x] T003 COMMIT `RED: agregar tests de la migración 007 y del índice de Scrum Master único (falla)`
- [x] T004 [US3][US4] GREEN: `internal/repository/migraciones/007_scrum_master_unico.sql`: comentario inicial con el criterio (de `creado_en` más antiguo; empate, menor `usuario_id`; por qué: la base de Render se vacía en cada deploy y los datos v2 crean un SM por proyecto), `UPDATE` de corrección y después `CREATE UNIQUE INDEX integrantes_un_scrum_master ON integrantes (proyecto_id) WHERE rol = 'scrum_master'`. `go test ./internal/repository/` en verde
- [x] T005 COMMIT `GREEN: implementar la migración 007 (Scrum Master único), test en verde`
- [x] T006 (sin cambios: no había nada que limpiar) (solo si hay algo que limpiar) REFACTOR de los tests de la fase 1 + COMMIT `REFACTOR: <qué se limpió>`

**Checkpoint**: `go test ./...` en verde con la 007 aplicada.

---

## Phase 2: Error de dominio y chequeo previo del service (US2)

**Goal**: agregar un segundo Scrum Master se rechaza con `ErrScrumMasterYaExiste`; si ya es integrante, gana el duplicado.

**Independent Test**: con un proyecto que ya tiene Scrum Master, `AsignarIntegrante` con rol `scrum_master` devuelve `ErrScrumMasterYaExiste` y no agrega a nadie.

- [x] T007 [US2] RED: tests en `internal/service/asignar_integrante_test.go`: invertir `TestAsignarIntegrante_VariosScrumMaster` (el segundo se rechaza con `domain.ErrScrumMasterYaExiste` y la lista no cambia); agregar Product Builder con Scrum Master presente (éxito); persona que ya es integrante pedida como `scrum_master`, incluida la que ya es el Scrum Master, devuelve el duplicado y no el nuevo error (FR-006); cambiar a `product_builder` el rol de Jimena en `TestAsignarIntegrante_ReutilizaUsuario`; carrera (FR-008, US3 esc. 5): con un fake cuyo `AgregarIntegrante` devuelve `ErrScrumMasterYaExiste` aunque el chequeo previo no vio Scrum Master (`fakes_test.go`), el service devuelve ese mismo error y no un 500 — confirmar que FALLA
- [x] T008 COMMIT `RED: agregar tests del rechazo del segundo Scrum Master en AsignarIntegrante (falla)`
- [x] T009 [US2] GREEN: `ErrScrumMasterYaExiste` en `internal/domain/errors.go` (mensaje "el proyecto ya tiene un Scrum Master") y chequeo previo en `internal/service/asignar_integrante.go`: con `ListarIntegrantes`, primero duplicado (400), después Scrum Master repetido. `go test ./...` en verde
- [x] T010 COMMIT `GREEN: rechazar un segundo Scrum Master en AsignarIntegrante, test en verde`
- [x] T011 (sin cambios: no había nada que limpiar) (solo si hay algo que limpiar, p. ej. reutilizar `rolEnProyecto`) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

---

## Phase 3: Repositorios (US2, US3)

**Goal**: el repositorio traduce la violación del índice (2067) a `ErrScrumMasterYaExiste` y la de la clave primaria (1555) sigue siendo duplicado; el repositorio en memoria imita el índice.

**Independent Test**: con SQLite migrado, `AgregarIntegrante` de un segundo `scrum_master` devuelve `domain.ErrScrumMasterYaExiste`; repetir a la misma persona devuelve `ErrMiembroDuplicado`.

- [x] T012 [US3] RED: tests en `internal/repository/sqlite_project_test.go` (segundo SM → `ErrScrumMasterYaExiste` y no `ErrMiembroDuplicado`; mismo usuario y proyecto → `ErrMiembroDuplicado`; builders y SM en otro proyecto aceptados; dos altas simultáneas de un `scrum_master` al mismo proyecto desde goroutines: exactamente una se acepta y la otra recibe `ErrScrumMasterYaExiste`, SC-003) y en `project_repository_en_memoria_test.go` (segundo SM → `ErrScrumMasterYaExiste`) — confirmar que FALLAN
- [x] T013 COMMIT `RED: agregar tests de la traducción del índice de Scrum Master único (falla)`
- [x] T014 [US3] GREEN: en `internal/repository/sqlite_project.go` distinguir por código extendido (`SQLITE_CONSTRAINT_PRIMARYKEY` → `ErrMiembroDuplicado`; `SQLITE_CONSTRAINT_UNIQUE` → `domain.ErrScrumMasterYaExiste`) y en `project_repository_en_memoria.go` rechazar el segundo SM. `go test ./...` en verde
- [x] T015 COMMIT `GREEN: separar el duplicado del segundo Scrum Master en el repositorio, test en verde`
- [x] T016 (sin cambios: no había nada que limpiar) (solo si hay algo que limpiar) REFACTOR + COMMIT `REFACTOR: <qué se limpió>`

---

## Phase 4: Mapeo HTTP y contrato (US2)

**Goal**: `POST /projects/{id}/members` responde 409 con un mensaje claro.

**Independent Test**: segundo Scrum Master → 409 con `mensaje` y sin `campo`; duplicado → 400 sin cambios.

- [x] T017 [US2] RED: tests en `internal/http/project_handler_test.go`: segundo `scrum_master` → 409 con el mensaje; persona ya integrante → 400 `integrante`; Product Builder → 201 — confirmar que FALLA
- [x] T018 COMMIT `RED: agregar tests del 409 de POST /projects/{id}/members (falla)`
- [x] T019 [US2] GREEN: `internal/http/sprint_handler.go` (`escribirError`): sumar `domain.ErrScrumMasterYaExiste` al grupo de los 409. `go test ./...` en verde
- [x] T020 COMMIT `GREEN: responder 409 ante un segundo Scrum Master, test en verde`
- [x] T021 [US2] Agregar la respuesta 409 de `contracts/members-409.yaml` a `specs/002-hu-04-crear-proyecto-equipo/contracts/openapi.yaml`
- [x] T022 COMMIT `Agregar el 409 al contrato de POST /projects/{id}/members`

---

## Phase 5: Escenarios BDD (US1, US2, US3)

**Goal**: los 4 escenarios Godog, escritos al final para no romper `TestFeatures`.

- [x] T023 [US1][US2][US3] `features/scrum_master_unico.feature` (Gherkin en español, tag `@HU-04`, Característica con el Como/Quiero/Para de la spec) con 4 escenarios: crear un proyecto deja un único Scrum Master; agregar un segundo Scrum Master (409 y lista sin cambios); agregar un Product Builder (éxito); el índice rechaza un segundo Scrum Master insertado por fuera del dominio. Pasos en `features/steps_scrum_master_unico.go` (siguiendo el patrón de `steps_creacion.go`) y registro en `features/features_test.go`. `go test ./features/` en verde
- [x] T024 COMMIT `Agregar escenarios BDD del Scrum Master único`

---

## Phase 6: MODELO-DE-DATOS.md (US5)

**Goal**: documento y esquema no divergen; va en el MISMO PR que la migración.

- [x] T025 [US5] `MODELO-DE-DATOS.md`: §4.3 (reemplazar la nota "no se impone un único Scrum Master" por la regla y el nombre del índice); §7 (D4 cerrada el 2026-11-03, D14 "no necesaria" porque el dueño es el Scrum Master, y el párrafo de encabezado sin D4 ni D14); §9 (la 007 en la tabla de números con el criterio de desempate, nota de que la 006 es de HU-20, y los casos del FR-014 en el punto 6); §5 (fila de la regla)
- [x] T026 [US5] Contrastar el documento con `esquema_test.go` y la 007: mismo nombre de índice, mismo `WHERE`, mismos casos. `go test ./internal/repository/` en verde
- [x] T027 COMMIT `Documentar el Scrum Master único en MODELO-DE-DATOS.md (D4, D14, 007)`

---

## Phase 7: Cierre

- [x] T028 `go build ./...`, `go vet ./...` y `gofmt -l .` sin salida
- [x] T029 `go test ./...` y `go test -race ./...` sin errores
- [x] T030 Recorrer `quickstart.md` (API manual y migración) y confirmar los resultados esperados
- [x] T031 Confirmar que `git log` muestra los commits RED/GREEN de cada fase y que el árbol queda limpio (sin `git push`)

---

## Dependencies & Execution Order

- Fase 1 → 2 → 3 → 4 → 5 → 6 → 7, en orden. La 3 necesita el índice de la 1; la 4, el error de la 2; la 5, todo lo anterior.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.
- Un RED en Go puede fallar por compilación (`undefined: domain.ErrScrumMasterYaExiste`); cuenta como "falla".

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra. Merge con "Create a merge commit".
- HU-20 (la 006) ya está en `main`: la 007 sigue a la 006 y no hay que renumerar (research.md §3).
