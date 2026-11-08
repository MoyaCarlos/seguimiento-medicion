---
description: "Tareas de la spec 039: criterios de aceptación al crear una historia"
---

# Tasks: Criterios de aceptación al crear una historia

**Input**: Design documents from `specs/039-criterios-aceptacion-creacion/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/api-backlog.md`, `quickstart.md`

**Tests**: OBLIGATORIOS. La constitución (Principio I) exige RED → COMMIT → GREEN → COMMIT para todo comportamiento con lógica. No hay migración ni cambio de esquema.

**Organization**: por user story. US1 es el alta con criterios (backend); US2 es la visualización en la confirmación y el listado (frontend).

**Compilación en RED**: cada commit RED puede incluir solo la firma o el campo stub necesario para compilar (por ejemplo, la firma de `ConCriteriosAceptacion` y el campo `CriteriosAceptacion` en `CrearHistoriaInput`), como se hizo en HU-20. La lógica va en el GREEN.

**BDD en rojo intermedio**: desde T002 hasta T024, `go test ./features/` falla porque los pasos todavía no existen. Durante esa ventana se exige `go test ./internal/... ./cmd/...` en verde; `go test ./...` se corre desde T024.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: se puede ejecutar en paralelo (archivos distintos, sin dependencia)
- **[Story]**: US1 o US2
- Todas las rutas son relativas a la raíz del repositorio

---

## Phase 1: Setup (Shared Infrastructure)

- [ ] T001 COMMIT `Agregar spec 039 de criterios de aceptación en la creación` — incluye `specs/039-criterios-aceptacion-creacion/spec.md`, `plan.md`, `research.md`, `data-model.md`, `contracts/api-backlog.md`, `quickstart.md`, `checklists/requirements.md` y este `tasks.md`
- [ ] T002 Crear `features/criterios_aceptacion_creacion.feature` (tag `@HU-01`) con escenarios: alta con criterios, alta sin criterios (omitido/vacío), criterios con espacios externos, límite de 2000 y error por 2001. No escribir los pasos todavía
- [ ] T003 COMMIT `Agregar escenarios BDD de criterios de aceptación en la creación`

---

## Phase 2: Foundational (Blocking Prerequisites)

**⚠️ CRITICAL**: sin la validación del dominio no se puede implementar el alta ni el handler.

- [ ] T004 RED: agregar en `internal/domain/backlog_item_test.go` los tests de `ConCriteriosAceptacion(criterios *string) (BacklogItem, error)`: `nil` → `nil`; `""` y `"   "` → `nil`; texto con espacios externos → recortado; 2000 runas → válido; 2001 runas → `ValidationError` con `Campo: "criterios_aceptacion"` y mensaje "los criterios de aceptación no pueden superar los 2000 caracteres"; conserva los demás campos. Incluir solo la firma stub (sin lógica) en `internal/domain/backlog_item.go` para compilar. Confirmar que FALLA
- [ ] T005 COMMIT `RED: agregar tests de ConCriteriosAceptacion (fallan)`
- [ ] T006 GREEN: implementar `ConCriteriosAceptacion` y el helper `normalizarCriterios` compartido con `ConDatosEditados` en `internal/domain/backlog_item.go`. `go test ./internal/domain` en verde
- [ ] T007 COMMIT `GREEN: implementar ConCriteriosAceptacion, test en verde`

**Checkpoint**: el dominio acepta, recorta y rechaza criterios con las mismas reglas que la edición.

---

## Phase 3: User Story 1 - Cargar los criterios al crear (Priority: P1) 🎯 MVP

**Goal**: `POST /backlog` persiste y devuelve `criterios_aceptacion` opcional, con validación de 2000 caracteres.

**Independent Test**: crear por API una historia con criterios y otra sin criterios; la respuesta y la base reflejan el valor normalizado; 2001 caracteres responde 400 sin crear.

- [ ] T008 [US1] RED: agregar en `internal/service/crear_historia_test.go` los tests de `CrearHistoriaBacklog.Ejecutar` con `CrearHistoriaInput.CriteriosAceptacion`: con criterios → se persisten y se devuelven; sin criterios (`nil`) → `nil`; solo espacios → `nil`; 2001 → `ValidationError` con campo `criterios_aceptacion` y sin persistir; proyecto inexistente con criterios → `ErrProyectoNoEncontrado` y sin persistir. Incluir solo el campo stub en `internal/service/crear_historia.go` (sin pasarlo todavía) para compilar. Confirmar que FALLA
- [ ] T009 [US1] COMMIT `RED: agregar tests de criterios de aceptación en CrearHistoriaBacklog (fallan)`
- [ ] T010 [US1] GREEN: agregar `CriteriosAceptacion *string` a `CrearHistoriaInput` y encadenar `NewBacklogItem` → `ConCriteriosAceptacion` antes de validar el proyecto y guardar, en `internal/service/crear_historia.go`. `go test ./internal/service` en verde
- [ ] T011 [US1] COMMIT `GREEN: aceptar criterios de aceptación en CrearHistoriaBacklog, test en verde`
- [ ] T012 [US1] RED: agregar en `internal/http/backlog_handler_test.go` los tests de `POST /backlog`: 201 con criterios en la respuesta y en el repositorio; 201 con `criterios_aceptacion: null` cuando se omite o se envía `null`; 201 con `null` cuando se envía solo espacios; 400 con `campo: "criterios_aceptacion"` y sin persistir cuando tiene 2001 caracteres. Confirmar que FALLA
- [ ] T013 [US1] COMMIT `RED: agregar tests del alta con criterios de aceptación (fallan)`
- [ ] T014 [US1] GREEN: agregar `CriteriosAceptacion *string` a `crearHistoriaRequest` en `internal/http/backlog_dto.go` y pasarlo desde `BacklogHandler.Crear` en `internal/http/backlog_handler.go`. `go test ./internal/http` en verde
- [ ] T015 [US1] COMMIT `GREEN: aceptar criterios de aceptación en POST /backlog, test en verde`

**Checkpoint**: la API crea historias con o sin criterios y responde 400 ante el exceso de longitud.

---

## Phase 4: User Story 2 - Ver los criterios en la confirmación y el listado (Priority: P2)

**Goal**: el formulario de creación permite cargar criterios y la confirmación los muestra con Markdown; el listado conserva su visualización actual.

**Independent Test**: crear desde la pantalla con y sin criterios y ver el resultado; el listado sigue mostrando los criterios de cada historia.

- [ ] T016 [US2] RED: agregar a `web/src/features/backlog/historia.test.ts` los casos de `FormularioHistoria.criterios`: `ESTADO_INICIAL` arranca con `''`; `armarCuerpo` recorta e incluye `criterios_aceptacion` cuando hay texto y lo omite cuando queda vacío; `validarHistoria` acepta 2000 caracteres, rechaza 2001 con el mensaje y acepta vacío/solo espacios. Confirmar que FALLA
- [ ] T017 [US2] COMMIT `RED: agregar tests de criterios en el formulario de creación (fallan)`
- [ ] T018 [US2] GREEN: actualizar `FormularioHistoria`, `CuerpoCrearHistoria`, `HistoriaCreada`, `ESTADO_INICIAL`, `armarCuerpo`, `validarHistoria`, `ErrorDeCampo` y `mapearError` en `web/src/features/backlog/historia.ts` (el error de API `criterios_aceptacion` se muestra en el campo `criterios`). `npm test -- historia` en verde
- [ ] T019 [US2] COMMIT `GREEN: aceptar criterios en el modelo del formulario, test en verde`
- [ ] T020 [US2] RED: agregar a `web/src/features/backlog/crear.test.tsx` los casos de componente: el formulario muestra el campo opcional "Criterios de aceptación"; al enviar con criterios el cuerpo los incluye; la confirmación los muestra con Markdown; sin criterios muestra "Sin criterios de aceptación"; 2001 caracteres muestra el error bajo el campo y no llama a la API; un 400 de API con `criterios_aceptacion` se muestra bajo el campo. Confirmar que FALLA
- [ ] T021 [US2] COMMIT `RED: agregar pruebas del campo y la confirmación de criterios (fallan)`
- [ ] T022 [US2] GREEN: agregar el `<textarea>` opcional, el manejo del error y la sección de criterios en la confirmación con `Markdown` en `web/src/features/backlog/crear.tsx`. Ajustar `web/src/features/backlog/api.test.ts` solo si el fixture de `HistoriaCreada` lo requiere. `npm test -- crear` en verde
- [ ] T023 [US2] COMMIT `GREEN: implementar criterios de aceptación en la pantalla de creación, test en verde`

**Checkpoint**: la pantalla crea historias con criterios y los muestra en la confirmación; el listado no cambia su comportamiento.

---

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T024 Agregar los pasos de `features/criterios_aceptacion_creacion.feature` en `features/steps_creacion.go`, reutilizando `scenarioContext` y el service real. `go test ./features` en verde
- [ ] T025 COMMIT `Agregar pasos BDD de criterios de aceptación en la creación`
- [ ] T026 Documentar en `docs/PRODUCT-BACKLOG.md` (sección HU-01 y, si corresponde, nota en HU-15) que los criterios se pueden cargar al crear (issue #174). Sin cambios en `MODELO-DE-DATOS.md`
- [ ] T027 COMMIT `Documentar criterios de aceptación en la creación (issue #174)`
- [ ] T028 Correr la verificación completa: `go build ./...`, `go vet ./...`, `gofmt -l .`, `go test ./...`; en `web/`, `npm run lint`, `npm test` y `npm run build`. Todo en verde
- [ ] T029 Recorrer `specs/039-criterios-aceptacion-creacion/quickstart.md` contra la API real y `npm run dev` (o el build servido por el backend) y registrar el resultado

---

## Dependencies & Execution Order

- **Setup (T001-T003)** → **Foundational (T004-T007)** → **US1 (T008-T015)** → **US2 (T016-T023)** → **Polish (T024-T029)**
- US1 depende de `ConCriteriosAceptacion` (T006); US2 depende de que la respuesta de la API incluya el campo (T014), aunque el frontend se puede probar con mocks.
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT. No se saltea el commit de RED.
- No se hace `git push` ni merge: el PR lo decide el equipo (merge commit, nunca squash).

## Implementation Strategy

1. Spec + escenarios BDD definidos.
2. Dominio validado (T004-T007).
3. US1: service y handler; API funcional (MVP).
4. US2: modelo del formulario y pantalla; confirmación y listado.
5. Pasos BDD, documentación y validación completa.

## Notes

- El repositorio SQLite no cambia: `Guardar` ya persiste `criterios_aceptacion` y ya tiene tests (`PersisteCriterios`, `CriteriosNull`).
- El listado ya muestra los criterios con Markdown desde HU-15/#147; no se modifica.
- No se agregan notas ni `depende_de` (D17): fuera de alcance.
- Los commits de esta lista no se ejecutan hasta que el equipo lo autorice.
