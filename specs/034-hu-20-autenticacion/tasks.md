---

description: "Tareas de HU-20 Autenticación de usuarios (issue #233): backend primero, frontend después, TDD con commits por fase"
---

# Tasks: Autenticación de usuarios (HU-20)

**Input**: Design documents from `/specs/034-hu-20-autenticacion/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento con lógica va como
RED → COMMIT → GREEN → COMMIT (→ REFACTOR → COMMIT solo si hay algo que limpiar). Los
COMMIT son tareas propias: no se marcan `[X]` hasta haber corrido `git commit`. Ninguna
tarea hace `git push` ni merge.

**Regla de compilación del RED**: el commit RED de cada fase TIENE que compilar y fallar por aserción (o, en Godog, por pasos sin definir bajo `Strict`), nunca por símbolos indefinidos. Cuando el test usa un símbolo nuevo, el mismo commit RED incluye un **stub sin lógica** (firma, tipos y valores cero; sin reglas de negocio), que luego el GREEN completa. El stub es esqueleto, no implementación: no se commitea ningún comportamiento junto al test.

**Organization**: el backend (Phases 1 a 6) va antes que el frontend (Phase 7), en la misma
rama y el mismo PR. La numeración de la migración es 006; la 007 es de #236 y no se usa.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: se puede hacer en paralelo (archivos distintos, sin depender de tareas incompletas)
- **[Story]**: US1 Registrarse · US2 Ingresar · US3 Sesión vigente y vencida · US4 Cerrar sesión · US5 Pantallas y navegación

## Phase 1: Setup — migración 006 y modelo de datos

- [X] T001 RED: agregar a `internal/repository/esquema_test.go` el test de la migración 006 con el patrón de `columnaEsperada`/`indiceEsperado`: `usuarios.contrasena_hash` TEXT **NULL** (sin default); `sesiones` con `id` INTEGER PK AUTOINCREMENT, `usuario_id` INTEGER NOT NULL, `token_hash` TEXT NOT NULL, `creado_en` TEXT NOT NULL, `vence_en` TEXT NOT NULL, `cerrada_en` TEXT **NULL**; índice UNIQUE sobre `token_hash`; FK `usuario_id → usuarios(id)` sin `ON DELETE`. Agregar el caso «base con usuarios previos»: aplicar las migraciones 001 a 005 (patrón de `migrate_test.go`), insertar un usuario, aplicar la 006 sin error y comprobar que el usuario conserva `id`, `nombre`, `nombre_normalizado`, `creado_en` y que su `contrasena_hash` es NULL. Confirmar que FALLA Compila sin stubs (el test es de SQL) y falla por aserción: faltan la columna y la tabla.
- [X] T002 COMMIT `RED: agregar test de esquema de la migración 006 (falla)` (incluye el stub, si lo hay)
- [X] T003 GREEN: crear `internal/repository/migraciones/006_autenticacion.sql` con `ALTER TABLE usuarios ADD COLUMN contrasena_hash TEXT;` y `CREATE TABLE sesiones` según [data-model.md](data-model.md) (sin `CHECK` de fecha, sin `ON DELETE`); `go test ./internal/repository/` en verde
- [X] T004 COMMIT `GREEN: agregar migración 006 de autenticación, test de esquema en verde`
- [X] T005 Actualizar `MODELO-DE-DATOS.md`: §2 (nota: `vence_en` entra en una regla de negocio y `cerrada_en` evita borrar), §3 (diagrama con `sesiones` y `usuarios.contrasena_hash`), §4.2 (columna `contrasena_hash` TEXT NULL), §4.11 nueva `sesiones` (HU-20, migración 006), §8 (HU-20 → `usuarios`, `sesiones`, 006) y §9 (fila de la 006)
- [X] T006 COMMIT `Documentar en MODELO-DE-DATOS.md la migración 006 de autenticación`

---

## Phase 2: Foundational — dominio, puertos, repositorios y servicios compartidos

**⚠️ CRITICAL**: ninguna historia puede empezar hasta cerrar esta fase.

### Dominio

- [X] T007 [P] RED: crear `internal/domain/usuario_credenciales_test.go` con tabla de casos. `ValidarNombreUsuario`: «ana» válido (3), «a.b-c_9» válido, 30 caracteres válido; 2 caracteres, 31 caracteres, vacío, con espacio interno, con espacio inicial o final (no se recorta), con `@`, con tilde («josé»), con emoji → `ValidationError{Campo: "nombre"}`. `ValidarContrasena`: 8 caracteres válida, 8 runas con tilde/emoji válida, con espacios válida, 128 válida; vacía, 7, 129 → `ValidationError{Campo: "contrasena"}` (contar runas, no bytes). Confirmar que FALLA **Stub (mismo commit RED):** crear `internal/domain/usuario_credenciales.go` con `ValidarNombreUsuario` y `ValidarContrasena` que devuelven `nil`; así compila y los casos inválidos fallan por aserción.
- [X] T008 COMMIT `RED: agregar tests de reglas de nombre de usuario y contraseña (falla)` (incluye el stub, si lo hay)
- [X] T009 GREEN: crear `internal/domain/usuario_credenciales.go` con `ValidarNombreUsuario(nombre string) error` (3 a 30 caracteres, solo letras ASCII, números, `.`, `-`, `_`) y `ValidarContrasena(contrasena string) error` (de 8 a 128 caracteres); `go test ./internal/domain/` en verde
- [X] T010 COMMIT `GREEN: implementar reglas de nombre de usuario y contraseña, test en verde`
- [X] T011 [P] RED: crear `internal/domain/sesion_test.go`: `Sesion.Vigente(ahora)` es verdadera antes de `VenceEn`, falsa en el instante exacto de `VenceEn` y después, falsa si `CerradaEn` no es nil aunque no haya vencido; `DuracionSesion` es 7×24 h. Confirmar que FALLA **Stub (mismo commit RED):** crear `internal/domain/sesion.go` con los tipos `Sesion` y `SesionConUsuario`, `DuracionSesion = 0` y `Vigente` que devuelve `false`; compila y falla por aserción.
- [X] T012 COMMIT `RED: agregar tests de vigencia de la sesión (falla)` (incluye el stub, si lo hay)
- [X] T013 GREEN: crear `internal/domain/sesion.go` (`Sesion{ID, UsuarioID, CreadoEn, VenceEn, CerradaEn *time.Time}`, `SesionConUsuario{Sesion, Usuario User}`, `DuracionSesion`, `Vigente`) y agregar a `internal/domain/errors.go` `ErrCredencialesInvalidas` («Usuario o contraseña incorrectos»), `ErrNoAutenticado` («Tu sesión no es válida») y `ErrNombreUsuarioEnUso` («Ese nombre de usuario ya existe»); `go test ./internal/domain/` en verde
- [X] T014 COMMIT `GREEN: implementar Sesion y errores de autenticación, test en verde`

### Puertos

- [X] T015 Definir los puertos de un método en `internal/repository/credencial_repository.go` (`RegistroCredenciales.Registrar(ctx, domain.User, hash string) (domain.User, error)`, `BuscadorCredenciales.ObtenerPorNombreNormalizado(ctx, string) (domain.User, string, error)` que devuelve el hash, vacío si es NULL) y en `internal/repository/sesion_repository.go` (`CreadorSesiones.Crear(ctx, huella string, domain.Sesion) (domain.Sesion, error)`, `BuscadorSesiones.ObtenerPorHuella(ctx, huella string) (domain.SesionConUsuario, error)`, `CerradorSesiones.Cerrar(ctx, sesionID int64, ahora time.Time) error`)
- [X] T016 COMMIT `Agregar puertos de credenciales y sesiones`

### Repositorios SQLite (contra SQLite real)

- [X] T017 [P] RED: crear `internal/repository/sqlite_credencial_test.go`: `Registrar` inserta usuario con `contrasena_hash`, devuelve `ID` y `CreadoEn`; segundo registro con el mismo nombre normalizado («Ana.Perez» tras «ana.perez») → `ErrUsuarioDuplicado`; nombre igual al de un usuario previo sin contraseña creado con `SQLiteUserRepository.Guardar` → `ErrUsuarioDuplicado`; `ObtenerPorNombreNormalizado` devuelve el hash guardado, devuelve hash vacío para el usuario previo (NULL) y `domain.ErrNoEncontrado` para uno inexistente; **caso límite de concurrencia**: dos goroutines registran a la vez el mismo nombre normalizado y exactamente una crea el usuario mientras la otra recibe `ErrUsuarioDuplicado`. Confirmar que FALLA **Stub (mismo commit RED):** crear `internal/repository/sqlite_credencial.go` con `SQLiteCredencialRepository`, su constructor y los métodos de `RegistroCredenciales` y `BuscadorCredenciales` devolviendo valores cero y `nil`.
- [X] T018 COMMIT `RED: agregar tests de SQLiteCredencialRepository (falla)` (incluye el stub, si lo hay)
- [X] T019 GREEN: crear `internal/repository/sqlite_credencial.go` (traduce `SQLITE_CONSTRAINT_UNIQUE` a `ErrUsuarioDuplicado` como `sqlite_user.go`; instantes en RFC 3339 UTC); `go test ./internal/repository/` en verde
- [X] T020 COMMIT `GREEN: implementar SQLiteCredencialRepository, test en verde`
- [X] T021 [P] RED: crear `internal/repository/sqlite_sesion_test.go`: `Crear` y `ObtenerPorHuella` devuelven la sesión con el usuario (`ID`, `Nombre`); huella inexistente → `domain.ErrNoEncontrado`; huella repetida → error de unicidad; `usuario_id` inexistente → `ErrIntegridadReferencial` (prueba la FK contra SQLite real); `Cerrar` fija `cerrada_en`, la fila **sigue existiendo** y `ObtenerPorHuella` la devuelve con `CerradaEn` no nil; cerrar una sesión no toca otra del mismo usuario. Confirmar que FALLA **Stub (mismo commit RED):** crear `internal/repository/sqlite_sesion.go` con `SQLiteSesionRepository`, su constructor y los métodos de `CreadorSesiones`, `BuscadorSesiones` y `CerradorSesiones` devolviendo valores cero y `nil`.
- [X] T022 COMMIT `RED: agregar tests de SQLiteSesionRepository (falla)` (incluye el stub, si lo hay)
- [X] T023 GREEN: crear `internal/repository/sqlite_sesion.go`; `go test ./internal/repository/` en verde
- [X] T024 COMMIT `GREEN: implementar SQLiteSesionRepository, test en verde`

### Servicios compartidos

- [X] T025 Agregar a `internal/service/fakes_test.go` los fakes en memoria de los puertos nuevos (`fakeCredenciales` implementa `RegistroCredenciales` y `BuscadorCredenciales` con la misma unicidad por nombre normalizado; `fakeSesiones` implementa `CreadorSesiones`, `BuscadorSesiones` y `CerradorSesiones`; `fakeCredenciales` protegido con mutex para el test concurrente de T041) y un `fakeContrasenas` que cuente las llamadas a `Verificar`
- [X] T026 COMMIT `Agregar fakes de credenciales y sesiones para los tests de servicio`
- [X] T027 [P] RED: crear `internal/service/contrasenas_test.go` (iteraciones bajas inyectadas): `Hashear` devuelve `pbkdf2-sha256$<iter>$<sal>$<hash>`; dos hashes de la misma contraseña difieren (sal única); `Verificar` acierta con la contraseña correcta, falla con otra, falla sin error con un formato corrupto o vacío; verifica con las iteraciones guardadas en el hash, no con las del servicio; el constructor de producción usa 600 000 iteraciones, sal de 16 bytes y clave de 32 bytes. Confirmar que FALLA **Stub (mismo commit RED):** crear `internal/service/contrasenas.go` con la interfaz `HasheadorContrasenas` (`Hashear`, `Verificar`), `ServicioContrasenas`, el constructor con iteraciones inyectables, el de producción (constantes 600 000, 16 y 32) y métodos que devuelven `""`, `false` y `nil`.
- [X] T028 COMMIT `RED: agregar tests de ServicioContrasenas (falla)` (incluye el stub, si lo hay)
- [X] T029 GREEN: crear `internal/service/contrasenas.go` con `crypto/pbkdf2` (`pbkdf2.Key(sha256.New, ...)`), `crypto/rand` y `subtle.ConstantTimeCompare`; `go test ./internal/service/` en verde
- [X] T030 COMMIT `GREEN: implementar ServicioContrasenas con PBKDF2, test en verde`
- [X] T031 [P] RED: crear `internal/service/emisor_sesiones_test.go` con reloj fijo: `Emitir(ctx, usuarioID)` devuelve un token de 43 caracteres base64url sin relleno, guarda en `CreadorSesiones` la huella SHA-256 hexadecimal del token (nunca el token), fija `CreadoEn = ahora` y `VenceEn = ahora + 7 días`; dos emisiones dan tokens distintos. Confirmar que FALLA **Stub (mismo commit RED):** crear `internal/service/emisor_sesiones.go` con `EmisorSesiones`, su constructor `(CreadorSesiones, reloj)`, `Emitir` devolviendo `""`, `domain.Sesion{}` y `nil`, y `HuellaDeToken` devolviendo `""`.
- [X] T032 COMMIT `RED: agregar tests de EmisorSesiones (falla)` (incluye el stub, si lo hay)
- [X] T033 GREEN: crear `internal/service/emisor_sesiones.go` (reloj `func() time.Time` inyectado, token de 32 bytes de `crypto/rand`, función exportada para calcular la huella de un token, reutilizable por «quién soy» y «cerrar sesión»); `go test ./internal/service/` en verde
- [X] T034 COMMIT `GREEN: implementar EmisorSesiones, test en verde`

### Infraestructura HTTP

- [ ] T035 [P] RED: crear `internal/http/auth_errores_test.go`: `escribirError` responde 401 `{"mensaje":"Usuario o contraseña incorrectos"}` para `ErrCredencialesInvalidas`, 401 `{"mensaje":"Tu sesión no es válida"}` para `ErrNoAutenticado` y 409 `{"mensaje":"Ese nombre de usuario ya existe"}` para `ErrNombreUsuarioEnUso`; los errores previos mantienen su código. Confirmar que FALLA Compila (usa `escribirError` y los errores de T013) y falla por aserción: hoy responde 500.
- [ ] T036 COMMIT `RED: agregar tests de traducción de errores de autenticación (falla)` (incluye el stub, si lo hay)
- [ ] T037 GREEN: agregar los tres casos a `escribirError` en `internal/http/sprint_handler.go`; `go test ./internal/http/` en verde
- [ ] T038 COMMIT `GREEN: traducir errores de autenticación a 401 y 409, test en verde`
- [ ] T039 Agregar `"/auth"` a `PrefijosAPI` en `internal/http/spa.go` y `'/auth': apiProxy` en `web/vite.config.ts`; correr `go test ./internal/http/` (`TestPrefijosAPI_IncluyenLosQueProxeaVite` en verde)
- [ ] T040 COMMIT `Registrar el prefijo /auth en la SPA y el proxy de Vite`

**Checkpoint**: la base está lista; las historias de usuario pueden avanzar.

---

## Phase 3: User Story 1 — Registrarse (Priority: P1) 🎯 MVP

**Goal**: crear una cuenta y quedar con la sesión iniciada.

**Independent Test**: `POST /auth/registro` devuelve 201 con token y ese token resuelve a un usuario; los escenarios 1 a 3 del issue en verde.

- [X] T041 [US1] RED: crear `internal/service/registrar_usuario_test.go`: éxito → usuario con el nombre tal como se escribió, `NombreNormalizado` en minúsculas, hash guardado distinto de la contraseña (nunca texto plano), sesión emitida con token; nombre inválido → `ValidationError{Campo: "nombre"}` sin llamar al registro ni al hash; contraseña de 7 → `Campo: "contrasena"`; de 128 válida y de 129 inválida; orden de validación nombre y luego contraseña; nombre repetido (otras mayúsculas) → `ErrNombreUsuarioEnUso` y ninguna sesión creada; un error del repositorio que no es duplicado se propaga; **caso límite de concurrencia**: dos registros simultáneos (goroutines) con el mismo nombre sobre `fakeCredenciales` dan exactamente un éxito y un `ErrNombreUsuarioEnUso`, y una sola sesión emitida. Confirmar que FALLA **Stub (mismo commit RED):** crear `internal/service/registrar_usuario.go` con `ResultadoSesion{Usuario domain.User, Token string, VenceEn time.Time}` (se reutiliza en `Ingresar`), las interfaces que consume, `RegistrarUsuario`, `NewRegistrarUsuario` y `Ejecutar` devolviendo `ResultadoSesion{}` y `nil`.
- [X] T042 [US1] COMMIT `RED: agregar tests de RegistrarUsuario (falla)` (incluye el stub, si lo hay)
- [X] T043 [US1] GREEN: crear `internal/service/registrar_usuario.go` (`NewRegistrarUsuario(RegistroCredenciales, hasheador, emisor)`; solo orquesta: validar, hashear, registrar, emitir; convierte `ErrUsuarioDuplicado` en `ErrNombreUsuarioEnUso`); `go test ./internal/service/` en verde
- [X] T044 [US1] COMMIT `GREEN: implementar RegistrarUsuario, test en verde`
- [X] T045 [US1] RED: crear `internal/http/auth_handler_test.go` (servicios reales con iteraciones de PBKDF2 bajas inyectadas) con `POST /auth/registro` según [api-auth.md](contracts/api-auth.md): 201 con `usuario.id`, `usuario.nombre`, `token` y `vence_en`; la respuesta no contiene la contraseña ni el hash; 400 con `campo` `nombre` y con `campo` `contrasena`; 400 `JSON inválido`; 409 con «Ese nombre de usuario ya existe». Confirmar que FALLA **Stub (mismo commit RED):** crear `internal/http/auth_dto.go` y `internal/http/auth_handler.go` con `AuthHandler`, `NewAuthHandler` y un `Registro` que no escribe nada.
- [X] T046 [US1] COMMIT `RED: agregar tests del handler POST /auth/registro (falla)` (incluye el stub, si lo hay)
- [X] T047 [US1] GREEN: crear `internal/http/auth_dto.go` y `internal/http/auth_handler.go` con `AuthHandler.Registro` (el handler no registra cuerpos ni cabeceras); `go test ./internal/http/` en verde
- [X] T048 [US1] COMMIT `GREEN: implementar handler POST /auth/registro, test en verde`
- [X] T049 [US1] Cablear en `cmd/api/main.go` los repositorios SQLite de credenciales y sesiones, `ServicioContrasenas` (600 000 iteraciones), `EmisorSesiones` (con `time.Now`), `RegistrarUsuario`, `AuthHandler` y la ruta `POST /auth/registro`
- [X] T050 [US1] COMMIT `Cablear POST /auth/registro en cmd/api`
- [X] T051 [US1] RED: crear `features/autenticacion.feature` con los escenarios 1 (registro exitoso, queda con sesión), 2 (nombre repetido, «Ese nombre de usuario ya existe», también con otras mayúsculas) y 3 (contraseña de menos de 8 caracteres, se rechaza indicando el campo `contrasena`) crear `features/steps_autenticacion.go` con `InitializeScenarioAutenticacion(ctx *godog.ScenarioContext)` vacía (stub) y registrarla en `features/features_test.go`, de modo que el paquete compile; `go test ./features/` debe FALLAR por pasos sin definir (Strict), no por compilación
- [X] T052 [US1] COMMIT `RED: agregar escenarios BDD de registro (fallan)` (incluye el stub, si lo hay)
- [X] T053 [US1] GREEN: completar `features/steps_autenticacion.go` (stub de T051) con la base de pruebas (SQLite en memoria con migraciones, servicios y `AuthHandler` reales con iteraciones de PBKDF2 bajas inyectadas, reloj controlable) y los pasos de los escenarios 1 a 3; `go test ./features/` en verde
- [X] T054 [US1] COMMIT `GREEN: agregar pasos BDD de registro, escenarios en verde`

**Checkpoint**: US1 funcional y verificable por sí sola.

---

## Phase 4: User Story 2 — Ingresar (Priority: P1)

**Goal**: ingresar con nombre y contraseña, sin revelar cuál falló.

**Independent Test**: `POST /auth/ingreso` da 200 con token con credenciales correctas y el mismo 401 idéntico en los tres casos de falla; escenarios 4 y 5.

- [x] T055 [US2] RED: crear `internal/service/ingresar_test.go`: éxito con el nombre en otras mayúsculas (se normaliza con `NormalizarNombre`) → usuario y sesión nueva; contraseña incorrecta, usuario inexistente y usuario previo sin contraseña (hash vacío) → los tres devuelven el mismo `ErrCredencialesInvalidas`; en los tres casos de falla `Verificar` se llama exactamente una vez (contra el hash real o contra el señuelo), para que el costo sea igual; dos ingresos generan dos sesiones distintas. Confirmar que FALLA **Stub (mismo commit RED):** crear `internal/service/ingresar.go` con `Ingresar`, `NewIngresar` y `Ejecutar` devolviendo `ResultadoSesion{}` y `nil`.
- [x] T056 [US2] COMMIT `RED: agregar tests de Ingresar (falla)` (incluye el stub, si lo hay)
- [x] T057 [US2] GREEN: crear `internal/service/ingresar.go` (`NewIngresar(BuscadorCredenciales, hasheador, emisor)`; el hash señuelo se calcula al construir el servicio con los mismos parámetros); `go test ./internal/service/` en verde
- [x] T058 [US2] COMMIT `GREEN: implementar Ingresar, test en verde`
- [x] T059 [US2] RED: agregar a `internal/http/auth_handler_test.go` `POST /auth/ingreso`: 200 con el mismo cuerpo que el registro; 401 con `{"mensaje":"Usuario o contraseña incorrectos"}` y cuerpo **idéntico byte a byte** para contraseña incorrecta, usuario inexistente y usuario sin contraseña; 400 `JSON inválido`; no valida formato de nombre ni largo de contraseña (un cuerpo vacío da el mismo 401). Confirmar que FALLA **Stub (mismo commit RED):** ampliar `NewAuthHandler` con el caso de uso y agregar un `Ingreso` que no escribe nada.
- [x] T060 [US2] COMMIT `RED: agregar tests del handler POST /auth/ingreso (falla)` (incluye el stub, si lo hay)
- [x] T061 [US2] GREEN: agregar `AuthHandler.Ingreso` en `internal/http/auth_handler.go`; `go test ./internal/http/` en verde
- [x] T062 [US2] COMMIT `GREEN: implementar handler POST /auth/ingreso, test en verde`
- [x] T063 [US2] Cablear `Ingresar` y la ruta `POST /auth/ingreso` en `cmd/api/main.go`
- [x] T064 [US2] COMMIT `Cablear POST /auth/ingreso en cmd/api`
- [x] T065 [US2] RED: agregar a `features/autenticacion.feature` los escenarios 4 (ingreso exitoso, también un usuario previo creado por HU-04 que no puede ingresar) y 5 (credenciales incorrectas: «Usuario o contraseña incorrectos» sin revelar cuál falló); `go test ./features/` debe FALLAR
- [x] T066 [US2] COMMIT `RED: agregar escenarios BDD de ingreso (fallan)` (incluye el stub, si lo hay)
- [x] T067 [US2] GREEN: agregar los pasos en `features/steps_autenticacion.go`; `go test ./features/` en verde
- [x] T068 [US2] COMMIT `GREEN: agregar pasos BDD de ingreso, escenarios en verde`

**Checkpoint**: US1 y US2 funcionan juntas e independientes.

---

## Phase 5: User Story 3 — Sesión vigente y vencida (Priority: P2)

**Goal**: «quién soy» con sesión vigente; rechazo uniforme de toda sesión no válida.

**Independent Test**: `GET /auth/yo` da 200 con sesión vigente y el mismo 401 para token ausente, mal formado, inexistente, vencido o cerrado; escenario 6 (API).

- [X] T069 [US3] RED: crear `internal/service/obtener_sesion_test.go`: token vigente → `domain.SesionConUsuario` (ID de la sesión y usuario con id y nombre); vencida (reloj en `VenceEn` y después) → `ErrNoAutenticado`; cerrada → `ErrNoAutenticado`; inexistente → `ErrNoAutenticado`; vacío, de longitud distinta de 43 o con caracteres fuera de base64url → `ErrNoAutenticado` **sin consultar el repositorio**. Confirmar que FALLA **Stub (mismo commit RED):** crear `internal/service/obtener_sesion.go` con `ObtenerSesion`, `NewObtenerSesion` y `Ejecutar(ctx, token)` devolviendo `domain.SesionConUsuario{}` y `nil`.
- [X] T070 [US3] COMMIT `RED: agregar tests de ObtenerSesion (falla)` (incluye el stub, si lo hay)
- [X] T071 [US3] GREEN: crear `internal/service/obtener_sesion.go` (`NewObtenerSesion(BuscadorSesiones, reloj)`; usa la huella del emisor y `Sesion.Vigente`; `domain.ErrNoEncontrado` del repositorio se convierte en `ErrNoAutenticado`); `go test ./internal/service/` en verde
- [X] T072 [US3] COMMIT `GREEN: implementar ObtenerSesion, test en verde`
- [X] T073 [US3] RED: agregar a `internal/http/auth_handler_test.go` `GET /auth/yo`: 200 `{"id","nombre"}`; 401 con `{"mensaje":"Tu sesión no es válida"}` e idéntico en: sin cabecera, esquema `Basic`, `Bearer` sin valor, token mal formado, token inexistente, sesión vencida, sesión cerrada. Confirmar que FALLA **Stub (mismo commit RED):** ampliar `NewAuthHandler` con el caso de uso y agregar un `Yo` que no escribe nada.
- [X] T074 [US3] COMMIT `RED: agregar tests del handler GET /auth/yo (falla)` (incluye el stub, si lo hay)
- [X] T075 [US3] GREEN: agregar `AuthHandler.Yo` y la lectura de `Authorization: Bearer <token>` en `internal/http/auth_handler.go`; `go test ./internal/http/` en verde
- [X] T076 [US3] COMMIT `GREEN: implementar handler GET /auth/yo, test en verde`
- [X] T077 [US3] Cablear `ObtenerSesion` y la ruta `GET /auth/yo` en `cmd/api/main.go`
- [X] T078 [US3] COMMIT `Cablear GET /auth/yo en cmd/api`
- [X] T079 [US3] RED: agregar a `features/autenticacion.feature` el escenario 6 (sesión vencida o inválida se rechaza, con sesión vigente se sabe quién es) y un escenario de que las rutas existentes siguen respondiendo sin token; `go test ./features/` debe FALLAR
- [X] T080 [US3] COMMIT `RED: agregar escenarios BDD de sesión vencida o inválida (fallan)` (incluye el stub, si lo hay)
- [X] T081 [US3] GREEN: agregar los pasos (avance del reloj de los pasos a 7 días); `go test ./features/` en verde
- [X] T082 [US3] COMMIT `GREEN: agregar pasos BDD de sesión vencida o inválida, escenarios en verde`

---

## Phase 6: User Story 4 — Cerrar sesión (Priority: P2)

**Goal**: el token usado deja de ser válido sin afectar otras sesiones.

**Independent Test**: tras `POST /auth/salir`, `GET /auth/yo` con ese token da 401; escenario 7.

- [x] T083 [US4] RED: crear `internal/service/cerrar_sesion_test.go`: cierra una sesión vigente y `ObtenerSesion` con ese token pasa a `ErrNoAutenticado`; otra sesión del mismo usuario sigue vigente; token inválido, vencido o ya cerrado → `ErrNoAutenticado` sin llamar a `Cerrar`; cerrar dos veces da error la segunda; el token se valida con el `ObtenerSesion` real (sobre los fakes), sin repetir en `CerrarSesion` la validación de formato ni de vigencia. Confirmar que FALLA **Stub (mismo commit RED):** crear `internal/service/cerrar_sesion.go` con `CerrarSesion`, `NewCerrarSesion` y `Ejecutar(ctx, token)` devolviendo `nil`.
- [x] T084 [US4] COMMIT `RED: agregar tests de CerrarSesion (falla)` (incluye el stub, si lo hay)
- [x] T085 [US4] GREEN: crear `internal/service/cerrar_sesion.go` (`NewCerrarSesion(obtener *ObtenerSesion, CerradorSesiones, reloj)`; solo llama a `ObtenerSesion` y luego a `Cerrar` con el ID de la sesión, sin repetir la validación del token ni la vigencia); `go test ./internal/service/` en verde
- [x] T086 [US4] COMMIT `GREEN: implementar CerrarSesion, test en verde`
- [x] T087 [US4] RED: agregar a `internal/http/auth_handler_test.go` `POST /auth/salir`: 204 sin cuerpo y el mismo token da 401 en `GET /auth/yo`; 401 `Tu sesión no es válida` sin efectos con token no vigente; otra sesión del usuario sigue válida. Confirmar que FALLA **Stub (mismo commit RED):** ampliar `NewAuthHandler` con el caso de uso y agregar un `Salir` que no escribe nada.
- [x] T088 [US4] COMMIT `RED: agregar tests del handler POST /auth/salir (falla)` (incluye el stub, si lo hay)
- [x] T089 [US4] GREEN: agregar `AuthHandler.Salir`; `go test ./internal/http/` en verde
- [x] T090 [US4] COMMIT `GREEN: implementar handler POST /auth/salir, test en verde`
- [x] T091 [US4] Cablear `CerrarSesion` y la ruta `POST /auth/salir` en `cmd/api/main.go`
- [x] T092 [US4] COMMIT `Cablear POST /auth/salir en cmd/api`
- [ ] T093 [US4] RED: agregar a `features/autenticacion.feature` el escenario 7 (cierre de sesión: el token deja de ser válido; con dos sesiones, cerrar una no cierra la otra); `go test ./features/` debe FALLAR
- [ ] T094 [US4] COMMIT `RED: agregar escenario BDD de cierre de sesión (falla)` (incluye el stub, si lo hay)
- [ ] T095 [US4] GREEN: agregar los pasos; `go test ./features/` en verde
- [ ] T096 [US4] COMMIT `GREEN: agregar pasos BDD de cierre de sesión, escenario en verde`

**Checkpoint**: backend completo (escenarios 1 a 7 en Godog). Desde acá empieza el frontend.

---

## Phase 7: User Story 5 — Pantallas y navegación (Priority: P2, frontend)

**Goal**: ingreso, registro, menú de usuario, ruta protegida y aviso «Tu sesión venció».

**Independent Test**: `cd web && npm test`; cada prueba lleva el escenario en el nombre (por ejemplo «Escenario 6: …»). Mantener verde `web/src/features/aislamiento.test.ts`.

### Almacenamiento del token

- [ ] T097 [P] [US5] RED: crear `web/src/features/autenticacion/almacenamientoToken.test.ts`: guarda, lee y borra el token; si `localStorage.setItem` lanza, `guardarToken` no lanza; si `getItem` lanza, `leerToken` devuelve `null`; si `removeItem` lanza, `borrarToken` no lanza. Confirmar que FALLA **Stub (mismo commit RED):** crear `web/src/features/autenticacion/almacenamientoToken.ts` con `leerToken` devolviendo `null` y `guardarToken` y `borrarToken` sin efecto.
- [ ] T098 [US5] COMMIT `RED: agregar tests del almacenamiento del token (falla)` (incluye el stub, si lo hay)
- [ ] T099 [US5] GREEN: crear `web/src/features/autenticacion/almacenamientoToken.ts` (clave única, cada acceso en `try/catch`, con copia en memoria si el almacenamiento falla); `npx vitest run src/features/autenticacion/almacenamientoToken.test.ts` en verde
- [ ] T100 [US5] COMMIT `GREEN: implementar almacenamiento del token, test en verde`

### Cliente HTTP compartido

- [ ] T101 [P] [US5] RED: agregar a `web/src/compartido/api/client.test.ts`: con token se envía `Authorization: Bearer <token>`; sin token no se envía la cabecera; un 401 de una petición con token llama a `alNoAutenticado`; un 401 sin token no la llama; un 401 de `/auth/ingreso` o `/auth/registro` aunque haya token no la llama y llega como `ApiError` 401; sin enganche configurado el cliente se comporta como antes. Confirmar que FALLA **Stub (mismo commit RED):** crear `web/src/compartido/api/sesionHttp.ts` con `configurarSesionHttp` sin efecto (`client.ts` no cambia todavía).
- [ ] T102 [US5] COMMIT `RED: agregar tests de cabecera Authorization y 401 en el cliente (fallan)` (incluye el stub, si lo hay)
- [ ] T103 [US5] GREEN: crear `web/src/compartido/api/sesionHttp.ts` (`configurarSesionHttp`) y usarlo en `web/src/compartido/api/client.ts` según [frontend-sesion.md](contracts/frontend-sesion.md); `npm test` en verde (los tests existentes del cliente siguen igual)
- [ ] T104 [US5] COMMIT `GREEN: enviar Authorization y avisar ante 401 en el cliente, test en verde`

### Contexto de sesión

- [ ] T105 [US5] RED: crear `web/src/features/autenticacion/SesionContext.test.tsx` con `fetch` simulado: «Escenario 1: registrarse guarda el token, fija el usuario y navega a la pantalla de origen o a `/`»; «Escenario 4: ingresar guarda el token y deja la sesión autenticada»; «Escenarios 2, 3 y 5: ante un error del servidor la promesa se rechaza con el `ApiError` y no se guarda token»; arranque con token válido (`GET /auth/yo` 200) → `autenticado`; «Escenario 6: arranque con token que da 401 → descarta el token, `anonimo`, navega a `/ingresar` con `state.vencida`»; «Escenario 7: salir pide `POST /auth/salir`, borra el token, `anonimo` y navega a `/ingresar`, incluso si el pedido falla»; «Escenario 6: con la sesión ya iniciada, una petición cualquiera que recibe 401 invoca `alNoAutenticado`: descarta el token, pasa a `anonimo` y navega a `/ingresar` con `state.desde` (la ruta actual) y `state.vencida`»; el contexto registra el enganche del cliente al montar y lo quita al desmontar. Confirmar que FALLA **Stub (mismo commit RED):** crear `web/src/features/autenticacion/api.ts` con `registrar`, `ingresar`, `yo` y `salir` que resuelven `undefined`, y `SesionContext.tsx` con `SesionProvider` que solo renderiza sus hijos y `useSesion` que devuelve estado `anonimo` con funciones sin efecto.
- [ ] T106 [US5] COMMIT `RED: agregar tests del contexto de sesión (fallan)` (incluye el stub, si lo hay)
- [ ] T107 [US5] GREEN: crear `web/src/features/autenticacion/api.ts` (`registrar`, `ingresar`, `yo`, `salir` sobre `compartido/api/client`) y `web/src/features/autenticacion/SesionContext.tsx` (`SesionProvider`, `useSesion`); `npm test` en verde
- [ ] T108 [US5] COMMIT `GREEN: implementar contexto de sesión, test en verde`

### Pantallas de ingreso y registro

- [ ] T109 [US5] RED: crear `web/src/features/autenticacion/Ingresar.test.tsx` y `Registro.test.tsx`: «Escenario 5: error de credenciales se muestra en el formulario sin borrar lo escrito»; «Escenario 6: con `state.vencida` se muestra “Tu sesión venció”»; «Escenario 4: ingresar con datos válidos navega a la pantalla de origen»; «Escenario 1: registrarse con datos válidos deja la sesión iniciada»; «Escenario 2: “Ese nombre de usuario ya existe” se muestra junto al campo nombre»; «Escenario 3: el error del campo `contrasena` se muestra junto a la contraseña»; enlaces entre ingreso y registro. Confirmar que FALLA **Stub (mismo commit RED):** crear `Ingresar.tsx` y `Registro.tsx` que renderizan un fragmento vacío.
- [ ] T110 [US5] COMMIT `RED: agregar tests de las pantallas de ingreso y registro (fallan)` (incluye el stub, si lo hay)
- [ ] T111 [US5] GREEN: crear `web/src/features/autenticacion/Ingresar.tsx` y `Registro.tsx` (etiquetas accesibles, `type="password"`, errores de campo con `ApiError.field`); `npm test` en verde
- [ ] T112 [US5] COMMIT `GREEN: implementar pantallas de ingreso y registro, test en verde`

### Ruta protegida y menú de usuario

- [ ] T113 [US5] RED: crear `web/src/features/autenticacion/RutaProtegida.test.tsx` y `MenuUsuario.test.tsx`: «Escenario 6: sin sesión redirige a `/ingresar` con `state.desde` y, tras ingresar, vuelve a la ruta pedida»; mientras `cargando` muestra un indicador y no el contenido; con sesión muestra el contenido; el menú muestra el nombre del usuario y «Cerrar sesión»; «Escenario 7: elegir “Cerrar sesión” cierra la sesión y lleva a `/ingresar`»; sin sesión el menú ofrece enlaces a ingresar y registrarse. Confirmar que FALLA **Stub (mismo commit RED):** crear `RutaProtegida.tsx` que renderiza `<Outlet />` y `MenuUsuario.tsx` que renderiza un fragmento vacío.
- [ ] T114 [US5] COMMIT `RED: agregar tests de RutaProtegida y MenuUsuario (fallan)` (incluye el stub, si lo hay)
- [ ] T115 [US5] GREEN: crear `web/src/features/autenticacion/RutaProtegida.tsx` y `MenuUsuario.tsx`; `npm test` en verde
- [ ] T116 [US5] COMMIT `GREEN: implementar RutaProtegida y MenuUsuario, test en verde`

### Rutas y encabezado

- [ ] T117 [US5] RED: actualizar `web/src/rutas.test.tsx` con una sesión simulada (token en `localStorage` y `fetch` que responde 200 a `GET /auth/yo`) para las rutas de siempre y sin sesión para los casos de redirección: `it.each` sobre `/`, `/projects/42/proyecto`, `/projects/42/backlog`, `/projects/42/backlog/nueva`, `/projects/42/sprints`, `/projects/42/status` y una ruta inexistente (`/no-existe`) → sin sesión cada una redirige a `/ingresar` conservando la ruta pedida en `state.desde`; `/ingresar` y `/registro` abren sus pantallas sin sesión; con sesión simulada las rutas de siempre abren las mismas pantallas y la ruta inexistente muestra «no encontrada»; tras ingresar se vuelve a la ruta pedida; el encabezado muestra el menú de usuario. Ningún test del backend se modifica. Confirmar que FALLA Compila sin stubs (usa `rutas` y las pantallas de T111 y T115) y falla por aserción: todavía no hay redirección ni rutas `/ingresar` y `/registro`.
- [ ] T118 [US5] COMMIT `RED: agregar tests de rutas protegidas e ingreso (fallan)` (incluye el stub, si lo hay)
- [ ] T119 [US5] GREEN: en `web/src/rutas.tsx` agregar `ingresar` y `registro` y envolver en `RutaProtegida` el resto, incluidos la raíz (`index`), `projects/:id/*` y `*` («no encontrada»), según [research.md](research.md) D12 (decidido); en `web/src/App.tsx` envolver el `Outlet` en `SesionProvider` y mostrar `MenuUsuario` en el encabezado; `npm test` en verde, incluido `aislamiento.test.ts`
- [ ] T120 [US5] COMMIT `GREEN: proteger las rutas de la SPA y mostrar el menú de usuario, test en verde`
- [ ] T121 [US5] Agregar `web/src/features/autenticacion/estilos.css` (formularios, aviso de sesión vencida, menú) e importarlo en las pantallas
- [ ] T122 [US5] COMMIT `Estilizar pantallas de autenticación y menú de usuario`

---

## Phase 8: Polish & Cross-Cutting Concerns

- [ ] T123 [P] Crear `internal/http/auth_logs_test.go` (iteraciones de PBKDF2 bajas inyectadas): capturar la salida de `log` durante registro, ingreso, «quién soy» y cierre de sesión, y comprobar que no aparecen la contraseña, el hash ni el token (SC-006); y que ninguna respuesta de error los contiene
- [ ] T124 COMMIT `Agregar test de que contraseñas y tokens no se filtran a logs ni a errores`
- [ ] T125 Confirmar que no se editó ningún `.feature` ni paso previo (`git diff main -- features/*.feature` solo muestra `autenticacion.feature`) y que no existe ninguna migración 007 en `internal/repository/migraciones/`
- [ ] T126 `go build ./...`, `go vet ./...` y `go test ./...` sin errores (incluye Godog, `internal/demo` y la migración sobre una base con usuarios previos)
- [ ] T127 En `web/`: `npm test`, `npm run lint` y `npm run build` sin errores
- [ ] T128 Recorrer la verificación manual de [quickstart.md](quickstart.md) (API e interfaz) y comprobar a mano que cada ruta protegida (incluida la raíz y una inexistente) lleva a `/ingresar` sin sesión y vuelve a la ruta pedida tras ingresar (D12)
- [ ] T129 Actualizar el `Status` de `spec.md` a «Implementada (pendiente de review)» y marcar estas tareas
- [ ] T130 COMMIT `Cerrar HU-20: spec actualizada y verificación completa`

---

## Dependencies & Execution Order

- **Phase 1 (migración y modelo)** → **Phase 2 (fundación)** → **Phases 3 a 6 (backend, US1 a US4)** → **Phase 7 (frontend, US5)** → **Phase 8**.
- US1 a US4 solo dependen de la Phase 2 y no entre sí; se prefiere el orden de prioridad. Todas comparten `internal/http/auth_handler.go`, `cmd/api/main.go`, `features/autenticacion.feature` y `features/steps_autenticacion.go`, así que se hacen en serie.
- US5 depende del backend terminado (Phases 3 a 6) y de T039 (proxy `/auth`).
- Dentro de cada comportamiento: RED → COMMIT → GREEN → COMMIT → (REFACTOR → COMMIT). Nunca saltear el COMMIT de RED; test e implementación nunca en el mismo commit.

### Oportunidades de paralelismo

- Phase 2: T007, T011, T017, T021, T027, T031 y T035 son tests de archivos distintos sin dependencia entre sí (cada GREEN depende de su propio RED).
- Phase 7: T097 (almacenamiento) y T101 (cliente) son independientes.
- T123 se puede preparar en cualquier momento después de T047.

## Implementation Strategy

1. **MVP**: Phases 1, 2 y 3 (US1): se puede registrar y obtener una sesión, verificable con Godog.
2. Sumar US2 (ingreso), US3 (sesión) y US4 (cierre) de a una, cada una dejando `go test ./...` en verde.
3. Con el backend completo, construir el frontend de adentro hacia afuera: almacenamiento → cliente → contexto → pantallas → ruta protegida → rutas.
4. Cerrar con la fase de Polish y la verificación del quickstart.

## Notes

- Un commit por fase del ciclo TDD: el historial `RED:` / `GREEN:` / `REFACTOR:` es evidencia evaluada por la cátedra.
- Nunca `git push` ni squash desde las tareas; el merge del PR es con «Create a merge commit».
- La contraseña, su hash y el token nunca se registran en logs ni aparecen en respuestas de error.
