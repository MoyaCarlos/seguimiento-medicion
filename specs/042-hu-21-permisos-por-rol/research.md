# Research: Permisos por rol (HU-21)

No quedaron `NEEDS CLARIFICATION` en el Technical Context: las cuatro dudas de la spec se
cerraron en `/speckit.clarify` (sesión 2026-11-11). Este documento registra las decisiones
de diseño tomadas a partir del código actual de `main`.

## R1. Dónde se exige la sesión (401)

- **Decision**: un envoltorio HTTP `conSesion` en `internal/http` (archivo nuevo
  `sesion_middleware.go`) que lee el token con el `tokenBearer` ya existente, lo resuelve con
  `service.ObtenerSesion` (HU-20) y guarda el usuario en el `context` de la petición. Los
  handlers lo leen con `usuarioDeSesion(r)`. Sin sesión válida responde 401 con
  `domain.ErrNoAutenticado`, el mismo mensaje que `GET /auth/yo`.
- **Rationale**: la autenticación es un detalle de infraestructura (Principio IV): el dominio
  y los servicios reciben un `usuarioID int64`, no un token. Reutiliza todo lo de HU-20.
- **Alternatives considered**: (a) resolver el token dentro de cada servicio — mezcla HTTP
  con casos de uso; (b) un envoltorio global con lista blanca de rutas públicas — hay que
  replicar el enrutamiento de `ServeMux` para decidir si una ruta es pública y es fácil que
  una ruta nueva quede abierta por error en la lista.

## R2. Cómo se garantiza que ninguna ruta quede sin proteger

- **Decision**: el registro de rutas sale de `cmd/api/main.go` a una función
  `apihttp.RegistrarRutas(mux, Handlers, conSesion)` en `internal/http/rutas.go`, que recorre
  una **tabla declarativa** de rutas (patrón, nivel, handler) y envuelve las protegidas.
  `http.ServeMux` no permite listar sus rutas, por eso la tabla es la fuente. Un test
  (`rutas_permisos_test.go`) la recorre contra la tabla «Qué se protege» de la spec: cada
  ruta protegida sin token → 401; cada pública sin token → distinto de 401; y compara la
  lista de patrones con una lista fija del test, de modo que una ruta nueva sin clasificar
  lo hace fallar.
- **Rationale**: Godog y los tests usan **el mismo** enrutamiento que producción (FR-028: sin
  atajos de prueba); una ruta nueva sin envolver hace fallar el test.
- **Alternatives considered**: seguir registrando en `main.go` — no se puede testear sin
  levantar el binario.

## R3. Dónde viven las reglas de pertenencia y de rol (403)

- **Decision**: en los **servicios** (casos de uso), igual que hoy en la votación
  (`internal/service/integrantes.go`: `rolEnProyecto`, `exigirScrumMaster`). Se agrega
  `exigirIntegrante` y `exigirScrumMaster` pasa a devolver `ErrNoEsIntegrante` cuando la
  persona no integra el proyecto (antes devolvía `ErrNoEsScrumMaster`).
- **Rationale**: «quién puede hacer qué» es una regla de negocio (testeable con TDD y fakes
  en memoria, Principio I). Cada servicio ya carga el recurso (Sprint, historia, defecto,
  proyecto), así que sabe a qué proyecto pertenece sin otra consulta.
- **Alternatives considered**: un middleware HTTP por rol — tendría que cargar el recurso
  para saber su proyecto (duplicaría el trabajo del servicio y rompería el orden 404 → 403).

## R4. Orden de errores 401 → 404 → 403 → 400/409

- **Decision**: 401 lo resuelve `conSesion` antes del handler. Dentro de cada servicio, el
  orden es: cargar el recurso (404) → `exigirIntegrante` / `exigirScrumMaster` (403) →
  validación de dominio y conflictos (400/409). Cuando hoy un servicio valida antes de
  cargar (por ejemplo, `CrearSprint` construye el Sprint antes de buscar el proyecto), se
  reordena.
- **Excepción declarada**: un cuerpo que no es JSON válido responde 400 «JSON inválido»
  antes del 403, porque en `POST /sprints` y `POST /backlog` el proyecto viaja en el cuerpo y
  sin leerlo no se sabe sobre qué proyecto evaluar. No revela nada de la acción.
- **Mensajes**: `ErrNoEsIntegrante` = «No sos integrante de este proyecto» (spec, escenario
  4); `ErrNoEsScrumMaster` = «Solo el Scrum Master puede realizar esta acción» (empieza con
  «Solo el Scrum Master puede», FR-016). Ambos ya se traducen a 403 en `escribirError`.

## R5. Identidad que viaja en la petición

- **Decision**: se quitan los campos `usuario_id` de `votacion_dto.go`, `integrante_id` de
  `horas_dto.go` y `Creador` de `crearProyectoRequest`, y el parámetro `usuario_id` de la
  query de `GET /backlog/{id}/votacion`. Como el decodificador JSON ignora campos
  desconocidos, una petición que todavía los mande se procesa a nombre de la sesión
  (clarify Q1). En la votación desaparece `leerPedido` con su 400 «el usuario es
  obligatorio».
- **Rationale**: lo más simple que cumple el escenario 7.

## R6. «Mis proyectos» y el rol

- **Decision**: `ProjectRepository` suma `ListarPorIntegrante(ctx, usuarioID)` que devuelve
  cada proyecto con el rol de esa persona (consulta `JOIN integrantes`, sin cambio de
  esquema). `ListarProyectos` pasa a recibir el `usuarioID`. La respuesta de
  `GET /projects` y la de `GET /projects/{id}` suman el campo `rol` de quien consulta.
- **Rationale**: la pantalla necesita el rol para habilitar o deshabilitar acciones; el panel
  ya carga `GET /projects/{id}` (que exige pertenencia), así que tener el rol ahí evita
  depender de la lista.
- **Alternatives considered**: un endpoint aparte `GET /projects/{id}/mi-rol` — una llamada
  más sin ganancia.

## R7. Lista de usuarios y configurar el equipo

- **Decision**: `UserRepository` suma `ListarRegistrados(ctx)` (usuarios con
  `contrasena_hash` no nulo, orden por nombre normalizado). Servicio nuevo `ListarUsuarios`
  y handler `GET /users` → `[{id, nombre}]`. `POST /projects/{id}/members` **conserva** el
  cuerpo `{"nombre", "rol"}`: la pantalla elige un nombre de la lista. `AsignarIntegrante`
  deja de crear usuarios (ya no llama a `resolverUsuario`): busca por nombre normalizado y,
  si no existe o no tiene contraseña, responde 404 «usuario no encontrado» (clarify Q3: el
  mismo error en ambos casos). `rol` es opcional: vacío o `product_builder` → Product
  Builder; `scrum_master` → 409 `ErrScrumMasterYaExiste` (siempre hay uno, el creador).
- **Rationale**: conservar `nombre` evita confundir el identificador del usuario elegido con
  el `usuario_id` de identidad que se elimina en R5, y no cambia el contrato de HU-04.
- **Alternatives considered**: `{"usuario_id": N}` — choca en nombre con R5.

## R8. Crear proyecto

- **Decision**: `CrearProyectoInput.Creador string` pasa a `CreadorID int64` (de la sesión).
  `resolverUsuario` queda sin uso y se elimina (ninguna ruta crea ya usuarios sin
  contraseña).
- **Rationale**: escenario 7 y FR-007.

## R9. Godog

- **Decision**:
  - **Feature nueva** `features/permisos_por_rol.feature` (`@HU-21`), un escenario por
    criterio de aceptación, que ejercita la API real vía `RegistrarRutas` +
    `httptest`, con el paso de contexto «Dado que inicié sesión como "X"» (registra o
    ingresa con `/auth/...` y guarda el token) y que envía `Authorization: Bearer` en cada
    llamada. Los usuarios se crean con el registro real de HU-20 (sin atajos).
  - **Features existentes**: la mayoría de los steps llaman a los servicios directamente.
    Donde un servicio pasa a recibir el `usuarioID` de quien actúa, el step usa el mismo
    paso de contexto «Dado que inicié sesión como "X"» (que en esos steps fija quién
    actúa) y deja de leer `usuario_id`/`integrante_id`/`creador` de la tabla del escenario.
    Los que usan `httptest` (consenso, consultar Sprints, datos de prueba, avance, obtener
    historia, listado) pasan a usar `RegistrarRutas` y mandar el token.
- **Rationale**: FR-027 y FR-028; el check `backend` en verde en el mismo PR.

## R10. Frontend

- **Decision**:
  - `client.ts` ya manda el token y maneja el 401 (HU-20). Se suma la traducción del 403:
    el mensaje «No sos integrante de este proyecto» lleva a una pantalla completa
    `NoIntegrante`; cualquier otro 403 se muestra como aviso «No tenés permiso» junto a la
    acción, con el mensaje del servidor (clarify Q4), reutilizando los `errores.ts` de cada
    feature.
  - Un hook `useRolProyecto` (o el `rol` que ya trae el proyecto cargado) decide si las
    acciones de gestión se habilitan; las deshabilitadas muestran «Solo el Scrum Master
    puede…» (sprints: crear/iniciar/cerrar; votación: `AccionesScrumMaster.tsx`; proyecto:
    `Editar.tsx`, `Configuracion.tsx`).
  - Se elimina `votacion/SelectorActor.tsx` y el selector de integrante de
    `sprints/FormularioHoras.tsx`; `proyecto/crear.tsx` deja de mandar `creador`;
    `Configuracion.tsx` elige de `GET /users` y no ofrece rol.
  - `compartido/proyectos/ListaProyectos.tsx` muestra el rol; sin proyectos, el mensaje y el
    botón de crear (escenario 9).
- **Rationale**: Principio III (pruebas de componente con el escenario en el nombre).

## R11. Esquema

- **Decision**: **sin cambios de esquema ni migración**. Pertenencia y rol salen de
  `integrantes` (§4.3 de `MODELO-DE-DATOS.md`), «registrado» de
  `usuarios.contrasena_hash` (§4.2) y la sesión de `sesiones` (§4.11). Las consultas nuevas
  (R6, R7) son lecturas sobre columnas existentes.
- Solo se actualiza texto del documento: la fila de HU-21 en el mapa historia → tablas (§8)
  y la nota de §4.2 de que HU-04 ya no crea usuarios sin contraseña. Si en la
  implementación aparece la necesidad de un índice o columna nuevos, se **frena y se
  consulta** al Scrum Master.
