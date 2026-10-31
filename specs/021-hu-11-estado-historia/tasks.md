---

description: "Tareas de HU-11: Actualización de Estado de una Historia en el Sprint (TDD con commits por fase)"
---

# Tasks: Actualización de Estado de una Historia en el Sprint (HU-11)

**Input**: Design documents from `/specs/021-hu-11-estado-historia/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica es
RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT).

**Compilación en RED**: todos los commits compilan. En cada RED se agregan, junto al test, solo
las firmas nuevas que el test necesita con un cuerpo vacío (devolver valor cero / `nil`), de modo
que `go build ./...` pase y el test falle por aserción, no por símbolo indefinido.

**Organization**: por user story de la spec (US1 avanzar el estado, US2 rechazos sin Sprint e
historia inexistente, US3 pantalla). US1 queda completa, con su escenario BDD en verde, antes de
empezar US2. US3 está **bloqueada** hasta que la pantalla de Sprints esté en `main`.

**Sin migraciones**: no se toca `MODELO-DE-DATOS.md` ni `internal/repository/migraciones/`.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3

---

## Phase 1: Setup

**Purpose**: partir de una base en verde.

- [X] T001 Verificar en la rama `feature/HU-11-estado-historia` que `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida) y `go test ./...` están en verde antes de empezar. Sin cambios ni commit

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: el estado `"En progreso"` y el puerto de persistencia que necesitan todas las user stories.

**⚠️ CRITICAL**: ninguna historia empieza hasta terminar esta fase.

- [X] T002 RED: crear internal/domain/estado_test.go con tests de `Estado.EsValida()` y de la constante `EstadoEnProgreso`: `EstadoNueva`, `EstadoEnProgreso` y `EstadoCompletada` devuelven `true`; `domain.Estado("Pendiente")` y `domain.Estado("")` devuelven `false`; `EstadoEnProgreso == "En progreso"` (el valor del enum del modelo). En internal/domain/estado.go agregar `EstadoEnProgreso Estado = "En progreso"` y `func (e Estado) EsValida() bool { return false }` (stub). Confirmar que FALLA
- [X] T003 COMMIT `RED: agregar tests del estado En progreso y su validación (falla)`
- [X] T004 GREEN: implementar `Estado.EsValida()` en internal/domain/estado.go con un `switch` sobre `EstadoNueva`, `EstadoEnProgreso` y `EstadoCompletada`, siguiendo el patrón de `Prioridad.EsValida()` — `go test ./...` en verde
- [X] T005 COMMIT `GREEN: agregar el estado En progreso y su validación, test en verde`
- [X] T006 Declarar en internal/repository/backlog_repository.go el puerto `CambioEstadoHistoriaRepository { ObtenerPorID(ctx context.Context, id int64) (domain.BacklogItem, error); CambiarEstado(ctx context.Context, id int64, desde, hacia domain.Estado, fechaFinalizacion *time.Time) error }`, documentando que `CambiarEstado` solo aplica si la historia sigue en `desde` y devuelve `domain.ErrConflictoEstadoHistoria` si no afectó ninguna fila. Separado de `HistoriaRepository` (ISP, research.md R3): **no** agregar el método a `HistoriaRepository`. Sin lógica
- [X] T007 COMMIT `Agregar puerto CambioEstadoHistoriaRepository`

**Checkpoint**: la fundación está lista; cualquier historia puede empezar.

---

## Phase 3: User Story 1 - Avanzar el estado de una historia dentro del Sprint activo (Priority: P1) 🎯 MVP

**Goal**: `PUT /backlog/{id}/estado` con `{"estado": "En progreso"}` o `{"estado": "Completada"}` mueve una historia de un Sprint `Activo` en la cadena `Nueva → En progreso → Completada`, registra la fecha de finalización al completar, conserva el Sprint y los Story Points, y rechaza las transiciones inválidas con el orden de errores de la spec.

**Independent Test**: sembrar un proyecto con un Sprint `Activo` y una historia `Nueva` asignada; pasarla a `En progreso` y luego a `Completada`, verificando estado, fecha de finalización, Sprint y Story Points; cerrar el Sprint y ver que no vuelve al Product Backlog; provocar cada rechazo (salto directo, Sprint `Pendiente`/`Finalizado`, retroceso, estado repetido) y ver el mensaje exacto.

### Escenarios BDD de US1

- [ ] T008 [US1] Escribir features/cambiar_estado_historia.feature en Gherkin español: tag `@HU-11`, Característica con el Como/Quiero/Para de la HU-11 ("Como Product Builder quiero cambiar el estado de una historia asignada al Sprint (En progreso, Completada) para poder reflejar el avance real del trabajo durante la iteración"), Antecedentes `Dado que existe un proyecto con identificador 1` (paso de features/steps_creacion.go, no se redefine), y exactamente 9 escenarios, uno por Acceptance Scenario de US1 en spec.md (US1.10, la carrera, no se escribe en Godog): US1.1 `Nueva → En progreso` conserva Sprint y datos; US1.2 `En progreso → Completada` registra la fecha de finalización; US1.3 la completada aporta sus Story Points como completados (usar una historia con 5 Story Points y verificar el dato que HU-03 va a consumir); US1.4 al cerrar el Sprint la completada no vuelve al Product Backlog; US1.5 el cambio se conserva al volver a consultar; US1.6 `Nueva → Completada` directo se rechaza con "primero debe pasar por En progreso"; US1.7 un solo escenario con Esquema del escenario y dos Ejemplos, uno para Sprint `Pendiente` y otro para `Finalizado`, ambos rechazados con "solo se cambia el estado de historias de un Sprint activo"; US1.8 retroceso (`En progreso → Nueva`) rechaza con "el estado de una historia solo avanza: no se puede volver a un estado anterior"; US1.9 estado repetido (`En progreso → En progreso`) rechaza con "la historia ya está en ese estado". US1.10 (carrera entre dos cambios simultáneos) se cubre en los tests de repositorio y handler porque no es determinista en Godog. Comentario en el archivo con la categoría de cada escenario (normal: US1.1 y US1.2; alternativo: US1.3; límite: US1.4 y US1.5; error: US1.6 a US1.9, con el Esquema de US1.7 contando como un escenario) y con la nota de US1.10. Sin pasos implementados: con el runner estricto, `go test ./features/` FALLA
- [ ] T009 [US1] COMMIT `RED: agregar escenarios BDD de cambio de estado de una historia (fallan)`

### Dominio: regla de la transición

- [ ] T010 [P] [US1] RED: crear internal/domain/estado_historia_test.go con una tabla de casos de `ValidarCambioDeEstado(h BacklogItem, s *Sprint, nuevo Estado) error`, sin base de datos: (a) `Nueva` + Sprint `Activo` → `En progreso` ⇒ `nil`; (b) `En progreso` + Sprint `Activo` → `Completada` ⇒ `nil`; (c) estado destino `domain.Estado("Pendiente")` ⇒ `ErrEstadoInvalido`; (d) sin Sprint (`SprintID` `nil`, `s` `nil`) → `En progreso` ⇒ `ErrHistoriaSinSprint`; (e) sin Sprint → `Completada` ⇒ `ErrHistoriaSinSprint`; (f) sin Sprint → `Nueva` ⇒ `ErrEstadoRepetido`; (g) Sprint `Pendiente` → `En progreso` ⇒ `ErrSprintNoActivoParaEstado`; (h) Sprint `Finalizado` + historia `En progreso` → `Completada` ⇒ `ErrSprintNoActivoParaEstado`; (i) `En progreso` → `En progreso` ⇒ `ErrEstadoRepetido`; (j) `Completada` → `Completada` ⇒ `ErrEstadoRepetido`; (k) `En progreso` → `Nueva` ⇒ `ErrRetrocesoEstado`; (l) `Completada` → `En progreso` ⇒ `ErrRetrocesoEstado`; (m) `Nueva` → `Completada` ⇒ `ErrSaltoDirecto`; precedencia (FR-018): (n) Sprint `Finalizado` + estado repetido ⇒ `ErrSprintNoActivoParaEstado`; (o) sin Sprint + `Nueva → Completada` ⇒ `ErrHistoriaSinSprint`; (p) estado inválido + sin Sprint ⇒ `ErrEstadoInvalido`. Comparar con `errors.Is` y comprobar `Campo == "estado"` en cada rechazo. Crear internal/domain/estado_historia.go con los siete errores de data-model.md (`ErrEstadoInvalido`, `ErrHistoriaSinSprint`, `ErrSprintNoActivoParaEstado`, `ErrEstadoRepetido`, `ErrRetrocesoEstado`, `ErrSaltoDirecto` como `ValidationError` con `Campo: "estado"` y los mensajes exactos; `ErrConflictoEstadoHistoria` como `errors.New` plano) y `ValidarCambioDeEstado` que devuelve `nil`. Confirmar que FALLA
- [ ] T011 [US1] COMMIT `RED: agregar tests de la regla de cambio de estado (falla)`
- [ ] T012 [US1] GREEN: implementar `ValidarCambioDeEstado` en internal/domain/estado_historia.go evaluando en este orden y devolviendo el primer error (research.md R1): 1) `!nuevo.EsValida()`; 2) `h.SprintID == nil`: si `nuevo == EstadoNueva` ⇒ `ErrEstadoRepetido`, si no ⇒ `ErrHistoriaSinSprint`; 3) `s == nil || s.Estado != SprintActivo` ⇒ `ErrSprintNoActivoParaEstado`; 4) `nuevo == h.Estado` ⇒ `ErrEstadoRepetido`; 5) `h.Estado == EstadoCompletada || nuevo == EstadoNueva` ⇒ `ErrRetrocesoEstado`; 6) `h.Estado == EstadoNueva && nuevo == EstadoCompletada` ⇒ `ErrSaltoDirecto` — `go test ./...` en verde
- [ ] T013 [US1] COMMIT `GREEN: implementar la regla de cambio de estado, test en verde`

### Persistencia: escritura condicional y fecha de finalización

- [ ] T014 [US1] RED: en internal/repository/sqlite_backlog_test.go, tests de `SQLiteBacklogRepository.CambiarEstado` y de la lectura de `FechaFinalizacion` contra `nuevaBDMigrada(t)`, con los helpers existentes (`guardarProyecto`, `guardarSprint`, `iniciarSprint`): (a) historia `Nueva` de un Sprint `Activo` → `En progreso` con fecha `nil` ⇒ `nil`, y `ObtenerPorID` devuelve `Estado` `En progreso` y `FechaFinalizacion` `nil`; (b) historia `En progreso` → `Completada` con la fecha `2026-10-29` ⇒ `nil`, y `ObtenerPorID` devuelve `Estado` `Completada` y `FechaFinalizacion` formateada `2026-10-29` (el `time.Time` de entrada puede llevar hora: el repositorio la descarta); (c) carrera simulada: dos llamadas seguidas con `desde = Nueva` sobre la misma historia ⇒ la primera `nil` y la segunda `errors.Is(err, domain.ErrConflictoEstadoHistoria)`, y la historia queda en el estado de la primera; (d) historia inexistente llamada directo ⇒ `errors.Is(err, domain.ErrConflictoEstadoHistoria)` y ninguna fila afectada; (e) una historia sin fecha sigue devolviendo `FechaFinalizacion` `nil` al leerla (las demás columnas no cambian). Agregar el campo `FechaFinalizacion *time.Time` a `domain.BacklogItem` y el stub `CambiarEstado` que devuelve `nil` en internal/repository/sqlite_backlog.go (sin tocar todavía `columnasHistoria`). Confirmar que FALLA
- [ ] T015 [US1] COMMIT `RED: agregar tests de CambiarEstado y de la fecha de finalización (falla)`
- [ ] T016 [US1] GREEN: en internal/repository/sqlite_backlog.go: agregar `fecha_finalizacion` a `columnasHistoria` y leerla en `escanearHistoria` con `sql.NullString` + `fechaDesdeSQL` (helper existente de sqlite_sprint.go), dejando `FechaFinalizacion` en `nil` si no hay valor; implementar `CambiarEstado` con una única sentencia `UPDATE historias SET estado = ?, fecha_finalizacion = ? WHERE id = ? AND estado = ?` usando `fechaASQL` para serializar (o `nil`), `errorDeEscritura` ante error y `RowsAffected() == 0` ⇒ `domain.ErrConflictoEstadoHistoria`. Sin transacción ni columnas nuevas — `go test ./...` en verde
- [ ] T017 [US1] COMMIT `GREEN: implementar CambiarEstado en SQLite, test en verde`

### Caso de uso: CambiarEstadoHistoria

- [ ] T018 [US1] RED: en internal/service/fakes_test.go agregar `fakeCambioEstadoRepository` (mapa de historias con `sync.Mutex`; `ObtenerPorID` devuelve `domain.ErrHistoriaNoEncontrada` si no existe; `CambiarEstado` aplica solo si el estado actual es `desde`, guarda estado y fecha, y si no coincide devuelve `domain.ErrConflictoEstadoHistoria`; cuenta las llamadas y registra los argumentos `desde`, `hacia` y `fecha`). Crear internal/service/cambiar_estado_historia_test.go usando `repository.NewSprintRepositoryEnMemoria()` para los Sprints y un reloj fijo (`func() time.Time { return time.Date(2026, 10, 29, 15, 30, 0, 0, time.UTC) }`): (a) `Nueva → En progreso` ⇒ devuelve la historia con `Estado` `En progreso` y `FechaFinalizacion` `nil`, y el fake recibió `desde: Nueva`, `hacia: En progreso`, `fecha: nil`; (b) `En progreso → Completada` ⇒ `Estado` `Completada` y `FechaFinalizacion` con el día del reloj, y el fake recibió `desde: En progreso`, `hacia: Completada`, `fecha` no nula; (c) historia inexistente ⇒ `ErrHistoriaNoEncontrada` y `CambiarEstado` no se llama; (d) historia sin Sprint ⇒ `ErrHistoriaSinSprint` y `CambiarEstado` no se llama; (e) Sprint `Pendiente` y `Finalizado` ⇒ `ErrSprintNoActivoParaEstado` y `CambiarEstado` no se llama; (f) salto directo, retroceso y estado repetido se propagan tal cual; (g) si el repositorio devuelve `ErrConflictoEstadoHistoria` (carrera), el servicio lo propaga sin convertirlo. Crear internal/service/cambiar_estado_historia.go con `CambiarEstadoHistoria`, `NewCambiarEstadoHistoria(historias repository.CambioEstadoHistoriaRepository, sprints repository.SprintRepository, reloj func() time.Time)` y `Ejecutar(ctx, id int64, nuevo domain.Estado) (domain.BacklogItem, error)` vacío. Confirmar que FALLA
- [ ] T019 [US1] COMMIT `RED: agregar tests del caso de uso CambiarEstadoHistoria (falla)`
- [ ] T020 [US1] GREEN: implementar `Ejecutar` en internal/service/cambiar_estado_historia.go: `historias.ObtenerPorID` (404) → si `SprintID != nil`, `sprints.ObtenerPorID` → `domain.ValidarCambioDeEstado` → si el destino es `Completada`, `fecha := s.reloj()`; si no, `nil` → `historias.CambiarEstado(ctx, id, historia.Estado, nuevo, fecha)` → devolver la historia con `Estado` y `FechaFinalizacion` actualizados. El servicio solo orquesta, sin reglas propias — `go test ./...` en verde
- [ ] T021 [US1] COMMIT `GREEN: implementar CambiarEstadoHistoria, test en verde`

### HTTP: PUT /backlog/{id}/estado

- [ ] T022 [US1] RED: en internal/http/backlog_handler_test.go, tests de `BacklogHandler.CambiarEstado` (ruta `PUT /backlog/{id}/estado` en un `ServeMux` de prueba) contra contracts/openapi.yaml: (a) 200 con la historia, `"estado"` nuevo y `"fecha_finalizacion"` `null` al pasar a `En progreso`; (b) 200 al completar con `"fecha_finalizacion"` con el día aportado por el reloj del service; (c) historia inexistente con cuerpo `{}` y con JSON inválido ⇒ 404 `{"mensaje": "historia no encontrada"}` (la existencia domina, FR-007); (d) `estado` ausente o vacío ⇒ 400 `{"campo": "estado", "mensaje": "el campo es obligatorio"}`; (e) `estado` `"Pendiente"` ⇒ 400 `{"campo": "estado", "mensaje": "el estado debe ser Nueva, En progreso o Completada"}`; (f) cada rechazo de la regla ⇒ 400 con el `campo` y el `mensaje` exactos de data-model.md; (g) si el repositorio devuelve `ErrConflictoEstadoHistoria` ⇒ 409 `{"mensaje": "la historia cambió de estado mientras tanto: recargá y volvé a intentar"}` **sin** `campo`; (h) id no numérico ⇒ 400. Agregar `CambiarEstado` a `stubBacklogRepository` (condicional, con un campo `errCambiarEstado` para forzar el conflicto), el parámetro `cambiarEstado *service.CambiarEstadoHistoria` a `NewBacklogHandler`, su llamada en `nuevoHandlerDePrueba` y en cmd/api/main.go (`nil` hasta T026), el método `CambiarEstado` vacío en internal/http/backlog_handler.go, `cambiarEstadoRequest { Estado string \`json:"estado"\` }` y el campo `FechaFinalizacion *string \`json:"fecha_finalizacion"\`` en `historiaResponse` (internal/http/backlog_dto.go). Confirmar que FALLA
- [ ] T023 [US1] COMMIT `RED: agregar tests del handler PUT /backlog/{id}/estado (falla)`
- [ ] T024 [US1] GREEN: implementar `CambiarEstado` en internal/http/backlog_handler.go siguiendo el patrón de `Editar` y `AsignarSprint`: `parsearID` → `obtenerHistoria.Ejecutar` (404 antes que el cuerpo) → decodificar `cambiarEstadoRequest` (error ⇒ 400 "JSON inválido"; `Estado` vacío ⇒ 400 campo `estado` "el campo es obligatorio") → `cambiarEstado.Ejecutar(r.Context(), id, domain.Estado(req.Estado))` → `escribirError` ante error / 200 con `aRespuesta`. En internal/http/backlog_dto.go cargar `FechaFinalizacion` en `aRespuesta` formateada con `time.DateOnly` (`null` si es `nil`). En `escribirError` agregar `case errors.Is(err, domain.ErrConflictoEstadoHistoria)` ⇒ `409` con solo `Mensaje` — `go test ./...` en verde
- [ ] T025 [US1] COMMIT `GREEN: implementar handler PUT /backlog/{id}/estado, test en verde`
- [ ] T026 [US1] Cablear en cmd/api/main.go: construir `service.NewCambiarEstadoHistoria(backlogRepo, sprintRepo, time.Now)` con los repositorios existentes, pasarlo a `NewBacklogHandler` (reemplazar el `nil` de T022) y registrar `mux.HandleFunc("PUT /backlog/{id}/estado", backlogHandler.CambiarEstado)` junto a las rutas de `/backlog`
- [ ] T027 [US1] COMMIT `Cablear PUT /backlog/{id}/estado en cmd/api`

### BDD de US1 en verde

- [ ] T028 [US1] GREEN: crear features/steps_estado_historia.go con `InitializeScenarioEstadoHistoria(ctx *godog.ScenarioContext, base *scenarioContext)`, que use `base.db` y un reloj fijo (`2026-10-29`) para construir `CambiarEstadoHistoria`; los Sprints se preparan con los services de HU-05 (`CrearSprint`, `IniciarSprint`, `CerrarSprint`) y la asignación con `AsignarHistoriaASprint` (HU-06), sin SQL a mano; el estado `En progreso` de partida se alcanza con el propio service de HU-11 (un paso `dado que la historia ... esté "En progreso"` que llama a `CambiarEstado` con la fecha del reloj). Los pasos afirman el `ValidationError` (campo `estado` y mensaje exacto), que la historia no cambió, y que la fecha de finalización es `2026-10-29`. Registrar en features/features_test.go `InitializeScenarioEstadoHistoria(ctx, sc) // HU-11`. No redefinir pasos existentes — `go test ./features/` en verde (9 escenarios @HU-11)
- [ ] T029 [US1] COMMIT `GREEN: implementar los pasos BDD de cambio de estado, escenarios en verde`

**Checkpoint**: US1 funciona de punta a punta (dominio, SQLite, servicio, API y BDD). Es el MVP.

---

## Phase 4: User Story 2 - Rechazar el cambio de estado de una historia que no está en un Sprint (Priority: P2)

**Goal**: una historia sin Sprint no se puede pasar a `En progreso` ni a `Completada` (mensaje "primero debe asignarse a un Sprint activo"), y una historia inexistente responde "historia no encontrada" antes de evaluar el cuerpo.

**Independent Test**: sembrar una historia sin Sprint e intentar las dos transiciones: ambas se rechazan con el mismo mensaje y la historia queda intacta; pedir el cambio de una historia inexistente devuelve 404 aunque el cuerpo sea inválido.

> La regla ya quedó implementada en US1 (es una sola función del dominio). Esta fase la consolida
> como historia propia: los escenarios BDD de US2 y la prueba de que la historia no se modifica.

- [ ] T030 [US2] Agregar a features/cambiar_estado_historia.feature exactamente 3 escenarios, uno por Acceptance Scenario de US2 en spec.md: US2.1 historia sin Sprint → `En progreso` ⇒ rechazo "primero debe asignarse a un Sprint activo" y la historia queda con estado `Nueva`; US2.2 historia sin Sprint → `Completada` ⇒ el mismo rechazo; US2.3 historia inexistente (identificador 999) ⇒ "historia no encontrada". Actualizar el comentario de categorías (error: 1 a 3). Con los pasos de US2 sin implementar, `go test ./features/` FALLA
- [ ] T031 [US2] COMMIT `RED: agregar escenarios BDD de rechazo sin Sprint e historia inexistente (fallan)`
- [ ] T032 [US2] GREEN: en features/steps_estado_historia.go, agregar los pasos de los 3 escenarios de US2 (historia sin Sprint, ambas transiciones, y cambio de estado de una historia inexistente) reutilizando el service y `domain.ErrHistoriaNoEncontrada`; comprobar que el estado y la fecha de finalización no cambiaron. No redefinir pasos existentes — `go test ./features/` en verde (12 escenarios @HU-11)
- [ ] T033 [US2] COMMIT `GREEN: implementar los pasos BDD de rechazo sin Sprint, escenarios en verde`

**Checkpoint**: US1 y US2 funcionan de punta a punta; el backend de HU-11 está completo.

---

## Phase 5: User Story 3 - Cambiar el estado desde la pantalla del Sprint (Priority: P3) 🔒 BLOQUEADA

**Bloqueo**: no empezar hasta que la pantalla de Sprints esté mergeada en `main`. Los nombres
exactos de componentes y archivos de esta fase se ajustan al código que traiga esa pantalla; lo que
no cambia es el comportamiento pedido.

**Goal**: desde la pantalla del Sprint activo, cada historia muestra la acción que corresponde a su estado (`Nueva` → "Iniciar", `En progreso` → "Completar"); al confirmar se ve el estado nuevo y la fecha de finalización, y ante un rechazo se muestra el `mensaje` de la API.

**Independent Test**: con la API real, iniciar y completar una historia del Sprint activo sin recargar a mano; provocar un rechazo (una historia del Product Backlog o un 409 de carrera) y ver el mensaje del sistema.

- [ ] T034 [US3] Verificar que la pantalla de Sprints está en `main` y traer `main` a la rama con `git merge origin/main` (merge local, sin push); `cd web && npm test` y `go test ./...` en verde. Si falta, detener la fase y avisar
- [ ] T035 [P] [US3] RED: crear web/src/features/sprints/api.test.ts con tests de `cambiarEstadoHistoria(historiaId, estado)` (hace `PUT /backlog/{historiaId}/estado` con cuerpo `{"estado": estado}` y devuelve la historia) y de que un 400/404/409 rechaza con el `mensaje` de la respuesta, usando el cliente de web/src/compartido/api/client.ts. Exportación vacía en web/src/features/sprints/api.ts. Confirmar que FALLA
- [ ] T036 [US3] COMMIT `RED: agregar tests del cliente de cambio de estado (falla)`
- [ ] T037 [US3] GREEN: implementar `cambiarEstadoHistoria` en web/src/features/sprints/api.ts — `npm test` en verde
- [ ] T038 [US3] COMMIT `GREEN: implementar el cliente de cambio de estado, test en verde`
- [ ] T039 [P] [US3] RED: crear web/src/features/sprints/EstadoHistoria.test.tsx para el control `EstadoHistoria` (props: la historia y un callback al cambiar): (a) con estado `Nueva` ofrece "Iniciar" y llama a `cambiarEstadoHistoria` con `"En progreso"` (FR-012); (b) con estado `En progreso` ofrece "Completar" y llama con `"Completada"`; (c) con estado `Completada` no ofrece acciones; (d) el botón queda deshabilitado mientras la petición está en curso y no envía dos veces; (e) tras un 200 llama al callback; (f) ante un error muestra el `mensaje` de la API tal cual (incluido el 409) y no llama al callback. Componente vacío en web/src/features/sprints/EstadoHistoria.tsx. Confirmar que FALLA
- [ ] T040 [US3] COMMIT `RED: agregar tests del componente EstadoHistoria (falla)`
- [ ] T041 [US3] GREEN: implementar `EstadoHistoria` en web/src/features/sprints/EstadoHistoria.tsx — `npm test` en verde
- [ ] T042 [US3] COMMIT `GREEN: implementar el componente EstadoHistoria, test en verde`
- [ ] T043 [US3] RED: en web/src/features/sprints/index.test.tsx (o el archivo de integración de la pantalla de Sprints, según cómo haya quedado), test de integración: cada historia del Sprint activo muestra su control de estado; al iniciar/completar se reflejan el estado nuevo y la fecha de finalización; ante un rechazo se muestra el mensaje y la lista queda como estaba; con un 409 se ve el aviso de recargar y reintentar. Confirmar que FALLA
- [ ] T044 [US3] COMMIT `RED: agregar test de integración del cambio de estado en pantalla (falla)`
- [ ] T045 [US3] GREEN: integrar `EstadoHistoria` en la lista de historias de la pantalla de Sprints, refrescando la lista tras un 200; estilos en los `.css` existentes de la feature — `npm test` en verde
- [ ] T046 [US3] COMMIT `GREEN: integrar el cambio de estado en la pantalla del Sprint, test en verde`
- [ ] T047 [US3] `cd web && npm run lint && npm test && npm run build` sin errores y validación manual del paso 3 de specs/021-hu-11-estado-historia/quickstart.md

**Checkpoint**: la sub-issue de Frontend de HU-11 queda lista para cerrarse.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T048 `go build ./...`, `go vet ./...`, `gofmt -l .` (sin salida), `go test ./...` y `go test -race ./internal/...` en verde
- [ ] T049 Validación manual de los pasos 1 y 2 de specs/021-hu-11-estado-historia/quickstart.md con la API real (`go run ./cmd/api` sobre una base nueva), anotando el resultado en esta tarea. Para cuerpos con tildes desde Git Bash en Windows, usar `curl --data-binary @archivo.json`
- [ ] T050 Revisar que contracts/openapi.yaml coincide con las respuestas reales (códigos, `campo`, `mensaje` y `fecha_finalizacion`); si hubo diferencias, corregir el contrato + COMMIT `docs: alinear el contrato de HU-11 con la implementación`

---

## Dependencies & Execution Order

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **US1 (Phase 3)** → **US2 (Phase 4)** → **Polish (Phase 6)**.
- **US3 (Phase 5)** depende de US1 y **de que la pantalla de Sprints esté en `main`**. Si el backend se entrega antes, Polish se corre sobre US1 + US2 y la Phase 5 queda abierta (la sub-issue de Frontend no frena el cierre del issue padre, AGENTS.md).
- Dentro de US1: T008–T009 (BDD) → dominio T010–T013 → SQLite T014–T017 → servicio T018–T021 → handler T022–T025 → main T026–T027 → pasos BDD T028–T029.
- Dentro de US2: T030–T033 (solo BDD: la regla ya vive en el dominio de US1).
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED.

## Parallel Opportunities

- **US1**: T010 (dominio) y T014 (SQLite) tocan archivos distintos y no dependen entre sí; el servicio (T018) necesita el dominio y el puerto; el handler (T022), el servicio; el `.feature` (T008) no depende de nada.
- **US2**: T030 no depende de código nuevo; los pasos de T032 necesitan el service.
- **US3**: T035 (cliente) y T039 (componente) son archivos distintos y pueden escribirse en paralelo contra un mock del cliente; T043 necesita ambos.
- Cada commit sigue siendo propio de su fase aunque las tareas se escriban en paralelo.

## Implementation Strategy

1. Setup + Foundational (T001–T007): la fundación (`En progreso` + puerto) lista.
2. **MVP = US1** (T008–T029): las dos transiciones con todas sus reglas y su BDD en verde. Se puede demostrar con curl (quickstart, pasos 2–12).
3. US2 (T030–T033): los rechazos sin Sprint y la historia inexistente.
4. Polish (T048–T050) sobre el backend → PR con `Closes #11` (merge commit, nunca squash), con las sub-issues de Especificación SDD, Escenarios BDD, Backend y Tests cerradas.
5. US3 (T034–T047) cuando la pantalla de Sprints esté en `main`: en la misma rama si el PR sigue abierto, o en un PR siguiente de la misma spec; cierra la sub-issue de Frontend.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Los commits del `.feature` (T009, T031) y de sus pasos (T029, T033) usan `RED:` / `GREEN:` siguiendo el precedente de las specs 009, 010 y 018, aunque la regla 2 de la plantilla pide commit sin prefijo para escenarios.
- Nunca `git push` ni squash desde las tareas.
- [P] tasks = different files, no dependencies.

## Excepción de proceso

Entre T009 y T033 el paquete `features` falla a propósito: los escenarios de HU-11 no tienen pasos (RED del BDD) y los GREEN de dominio, repositorio, servicio y handler se verifican por paquete con `go test ./internal/... ./cmd/...`. T029 pone en verde los 10 corridas de US1 (9 escenarios, con el esquema de US1.7 contado dos veces) y T033 las 3 de US2. Los pasos están en `features/steps_estado_historia.go`.

Orden de los commits: el RED de US2 (T031) se commiteó antes del GREEN de los pasos de US1 (T029). No cambia el resultado: ambos GREEN dejan `features` en verde.

Ajuste del escenario "Solo se cambia el estado de historias de un Sprint activo" después de su RED: un Sprint "Finalizado" no puede tener una historia "Nueva" (al cerrar un Sprint, `CerrarConArrastre` devuelve al Product Backlog las historias que no están "Completada"), así que el ejemplo "Finalizado" usa una historia "Completada", el único estado que el sistema puede producir en un Sprint cerrado. El rechazo y su mensaje no cambian, porque la regla evalúa "Sprint no activo" antes que estado repetido y retroceso (FR-018). Se corrigió en un commit aparte del `.feature`, sin SQL a mano en los pasos.

El escenario de carrera (US1.10 de la spec) no se puede reproducir de forma determinista en Godog: se cubre en `sqlite_backlog_test.go` y en `backlog_handler_test.go`, como en HU-06. `go test -race` no se pudo correr localmente (el compilador de C instalado es de 32 bits); la carrera se prueba con el UPDATE condicional y con dos llamadas seguidas.
