---

description: "Tareas de HU-17 Renombrar el propio usuario (#176): backend, BDD, frontend y documentación con TDD"
---

# Tasks: Renombrar el propio usuario (HU-17)

**Input**: Design documents from `/specs/048-hu-17-renombrar-usuario/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento va como
RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT solo si hay algo que limpiar). No se hace
`git push` ni merge desde las tareas.

**Regla de compilación del RED**: el commit RED debe compilar y fallar por aserción (o, en
Godog, por pasos sin definir bajo `Strict`), nunca por símbolos indefinidos. Si el test usa un
símbolo nuevo, el mismo commit RED incluye un stub sin lógica.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: se puede hacer en paralelo (archivos distintos).
- **[Story]**: US1 cambia el nombre · US2 nombre en uso · US3 sesión obligatoria.

---

## Phase 1: Documentación — spec, modelo y backlog

- [x] T001 Crear `specs/048-hu-17-renombrar-usuario/` con spec, plan, research, data-model,
  quickstart, contrato y estas tareas.
- [x] T002 COMMIT `docs: agregar la spec 048 de HU-17 (renombrar usuario) (Refs #176)`.
- [x] T003 Actualizar `MODELO-DE-DATOS.md`: §4.2 (regla de renombrado), §5 (escritura atómica
  y unicidad) y §8 (fila HU-17, sin migración).
- [x] T004 COMMIT `docs: documentar el renombrado de usuario en MODELO-DE-DATOS.md (Refs #176)`.
- [x] T005 Redefinir la sección HU-17 de `docs/PRODUCT-BACKLOG.md`: descripción de opción A,
  escenarios 5 y 6, notas de alineación, prioridad C, estado, valor 8 y estimación 2 SP.
- [x] T006 COMMIT `docs: redefinir HU-17 en el Product Backlog (Refs #176)`.

---

## Phase 2: Backend — puerto y repositorio

**Goal**: persistir el cambio de `nombre` y `nombre_normalizado` de forma atómica, con
unicidad case-insensitive.

- [x] T007 [US1] RED: en `internal/repository/sqlite_credencial_test.go`, `CambiarNombre`
  actualiza ambos campos y conserva `id`, `creado_en` y `contrasena_hash`; un nombre con
  otras mayúsculas de otro usuario devuelve `ErrUsuarioDuplicado` y no cambia nada; un `id`
  inexistente devuelve `domain.ErrNoEncontrado`; dos goroutines que intentan el mismo nombre
  dejan exactamente un éxito y un `ErrUsuarioDuplicado`. Confirmar que FALLA (compila con un
  stub sin lógica).
- [x] T008 [US1] COMMIT `RED: agregar tests de CambiarNombre en SQLite (falla) (Refs #176)`.
- [x] T009 [US1] GREEN: agregar `ActualizadorNombreUsuario` a
  `internal/repository/credencial_repository.go` e implementar `CambiarNombre` en
  `internal/repository/sqlite_credencial.go` con un solo `UPDATE`; traducir
  `SQLITE_CONSTRAINT_UNIQUE` a `ErrUsuarioDuplicado` y 0 filas a `domain.ErrNoEncontrado`.
  `go test ./internal/repository/` en verde.
- [x] T010 [US1] COMMIT `GREEN: implementar CambiarNombre atómico, test en verde (Refs #176)`.

---

## Phase 3: Backend — caso de uso

- [x] T011 [US1] RED: `internal/service/cambiar_nombre_usuario_test.go`: nombre válido →
  devuelve `domain.User` con `Nombre` y `NombreNormalizado` y llama al repositorio con el id
  de la sesión; nombre inválido (2, 31, espacio, `@`, tilde) → `ValidationError{Campo:
  "nombre"}` sin llamar al repositorio; `ErrUsuarioDuplicado` → `ErrNombreUsuarioEnUso`;
  `domain.ErrNoEncontrado` se propaga; cualquier otro error se envuelve.
- [x] T012 [US1] COMMIT `RED: agregar tests de CambiarNombreUsuario (falla) (Refs #176)`.
- [x] T013 [US1] GREEN: crear `internal/service/cambiar_nombre_usuario.go` con
  `CambiarNombreUsuario`, su constructor y `Ejecutar`; agregar el fake del puerto a
  `internal/service/fakes_test.go` si hace falta. `go test ./internal/service/` en verde.
- [x] T014 [US1] COMMIT `GREEN: implementar CambiarNombreUsuario, test en verde (Refs #176)`.

---

## Phase 4: Backend — handler, ruta y cableado

- [x] T015 [US1] RED: en `internal/http/auth_handler_test.go`, `PUT /auth/yo`: 200 con el
  nombre nuevo y sin contraseña; 401 sin sesión/token inválido; 400 JSON inválido, nombre
  ausente o inválido con `campo: "nombre"`; 409 con «Ese nombre de usuario ya existe». En
  `internal/app/rutas_permisos_test.go`, `PUT /auth/yo` es `SoloSesion`. Confirmar que FALLA.
- [x] T016 [US1] COMMIT `RED: agregar tests del handler PUT /auth/yo (falla) (Refs #176)`.
- [x] T017 [US1] GREEN: agregar `cambiarNombreRequest` en `auth_dto.go`, `AuthHandler.CambiarNombre`
  y el caso de uso al constructor; registrar `{"PUT /auth/yo", SoloSesion, h.Auth.CambiarNombre}`
  en `rutas.go`; cablear en `internal/app/app.go` y `cmd/api/main.go` si corresponde.
  `go test ./...` en verde.
- [x] T018 [US1] COMMIT `GREEN: implementar PUT /auth/yo, test en verde (Refs #176)`.

---

## Phase 5: BDD — escenarios de renombrado

- [x] T019 [US1] RED: crear `features/gestion_usuarios.feature` con `@HU-17` y los escenarios:
  (1) cambio exitoso con sesión vigente y nuevo ingreso; (2) nombre en uso con otras
  mayúsculas → 409 y el nombre anterior intacto; (3) sin sesión → 401; (4) nombre inválido →
  400 con `campo: "nombre"`. Crear `features/steps_gestion_usuarios.go` con el registro en
  `features/features_test.go`; `go test ./features/` debe FALLAR por pasos sin definir.
- [x] T020 [US1] COMMIT `RED: agregar escenarios BDD de renombrar usuario (fallan) (Refs #176)`.
- [x] T021 [US1] GREEN: implementar los pasos reutilizando la base de autenticación y los
  helpers existentes; `go test ./features/` en verde.
- [x] T022 [US1] COMMIT `GREEN: agregar pasos BDD de renombrar usuario, escenarios en verde (Refs #176)`.

---

## Phase 6: Frontend — API y contexto de sesión

- [x] T023 [US1] RED: en `web/src/features/autenticacion/SesionContext.test.tsx`, con sesión
  autenticada, `cambiarNombre('ana.perez')` llama a `PUT /auth/yo` con `{nombre}` y actualiza
  el usuario del contexto; si el servidor falla, la promesa se rechaza y el usuario no cambia.
  Confirmar que FALLA.
- [x] T024 [US1] COMMIT `RED: agregar tests de cambiarNombre en el contexto de sesión (fallan) (Refs #176)`.
- [x] T025 [US1] GREEN: agregar `cambiarNombre` a `web/src/features/autenticacion/api.ts`
  (`put<Usuario>('/auth/yo', { nombre })`) y al `SesionContext` (actualiza `usuario` con la
  respuesta); `npm test -- SesionContext` en verde.
- [x] T026 [US1] COMMIT `GREEN: implementar Cambiar nombre en el menú, tests en verde (Refs #176)` (incluye la API y el contexto).

---

## Phase 7: Frontend — formulario en el menú

- [x] T027 [US1] RED: en `web/src/features/autenticacion/MenuUsuario.test.tsx`: «Given una
  sesión, When elijo Cambiar nombre, Then veo el formulario con el nombre actual»; «Given un
  nombre válido, When guardo, Then se envía PUT /auth/yo y el menú muestra el nombre nuevo»;
  «Given un 409, When guardo, Then veo «Ese nombre de usuario ya existe» y el nombre anterior
  sigue»; «Given un 400 de campo nombre, When guardo, Then veo el mensaje y conservo lo
  escrito»; «Given el formulario, When cancelo, Then no se llama a la API». Confirmar que FALLA.
- [x] T028 [US1] COMMIT `RED: agregar tests del formulario Cambiar nombre (fallan) (Refs #176)`.
- [x] T029 [US1] GREEN: implementar el formulario en `MenuUsuario.tsx` (botón, input, guardar,
  cancelar, errores), los estilos en `estilos.css` y la actualización del nombre en el menú.
  `npm test`, `npm run lint` y `npm run build` en verde.
- [x] T030 [US1] (incluido en el commit de T026).

---

## Phase 8: Cierre

- [x] T031 Marcar las tareas hechas en este archivo y actualizar el `Status` de la spec.
- [x] T032 COMMIT `docs: cerrar HU-17 (spec y tareas) (Refs #176)`.
- [x] T033 `gofmt -l .`, `go build ./...`, `go vet ./...`, `go test ./...`; `cd web && npm test`,
  `npm run lint` y `npm run build` en verde.
- [ ] T034 Crear el PR con «Closes #176» y actualizar el cuerpo del issue con la redefinición
  acordada.
