# Feature Specification: Permisos por rol (HU-21)

**Feature Branch**: `feature/HU-21-permisos-por-rol`

**Issue**: #234 (sub-issue de especificación #243, 8 SP)

**Created**: 2026-11-11

**Status**: Draft

**Dependencias**: HU-20 (autenticación, spec `034-hu-20-autenticacion`) y #236 (un único
Scrum Master por proyecto, spec `040-scrum-master-unico`). Protege las rutas de HU-02,
HU-04, HU-05, HU-06, HU-07, HU-08, HU-11, HU-14 y HU-15.

**Input**: User description: "HU-21: Permisos por rol (segregación de funciones). Fuente de
verdad: la sección «## HU-21» de `docs/PRODUCT-BACKLOG.md` y el issue #234. Los 9 criterios
de aceptación y la lista «Qué se protege» están DEFINIDOS: copiarlos fielmente, no
reabrirlos. Reglas de negocio, casos límite y errores en formato EARS. Godog: paso «Dado que
inicié sesión como X», escenarios ajustados en el mismo PR, sin atajos de prueba en
producción."

## Clarifications

### Session 2026-11-11

- Q: Si una petición trae `usuario_id`, `integrante_id` o `creador`, ¿qué hace el sistema? → A: Se ignora el dato y la acción queda a nombre de la persona con sesión (no se responde 400).
- Q: Con varios errores a la vez, ¿cuál se responde primero? → A: 401 (sin sesión) → 404 (recurso inexistente) → 403 no integrante → 403 no Scrum Master → 400/409 (validación y conflictos), como en #101.
- Q: ¿Qué usuarios aparecen en la lista para configurar el equipo y se pueden agregar? → A: Solo los registrados con contraseña (HU-20). Un usuario sin contraseña no aparece en la lista y agregarlo se rechaza igual que un usuario inexistente.
- Q: ¿Cómo se muestra «No tenés permiso» ante un 403 de rol inesperado? → A: Como aviso en la misma pantalla, junto a la acción, con el mensaje del servidor, sin salir de la pantalla ni perder lo escrito. «No sos integrante» sí es una pantalla completa.
- Q: ¿El equipo acepta los cambios de comportamiento sobre escenarios ya en verde y las decisiones de diseño del plan? → A: Sí (confirmado por el equipo, 2026-11-11). Se aceptan: (1) un no integrante recibe «No sos integrante de este proyecto» (403) antes que el error de rol, lo que cambia el mensaje esperado en `consenso_estimacion.feature` y `votacion_planning_poker.feature`, y en `imputacion_horas.feature` pasa de 400 sobre `integrante_id` a 403; (2) un cuerpo que no es JSON válido responde 400 antes del 403; (3) los integrantes se agregan por `nombre`, elegido de la lista de usuarios registrados; (4) se actualiza solo texto de `MODELO-DE-DATOS.md` (§4.2 y §8), sin cambio de esquema. La regla de fondo de esos escenarios no cambia: un no integrante sigue sin poder actuar.

## Objetivo

Que las acciones de gestión del proyecto (Sprints, votación, edición y equipo) solo las
pueda hacer su Scrum Master; que toda acción quede a nombre de la persona con sesión
iniciada, sin que nadie pueda actuar ni consultar en nombre de otra; y que cada persona vea
solo los proyectos que integra, con su rol. Es el PR 2 de 2 del login: HU-20 estableció la
identidad, esta historia la impone.

**Descripción (backlog)**: Como Scrum Master quiero que las acciones de gestión del proyecto
solo las pueda hacer el Scrum Master para que el resto del equipo participe sin alterar la
planificación.

## Entradas

- La sesión de la persona (token de HU-20 en la cabecera `Authorization: Bearer`), en toda
  petición salvo las rutas públicas.
- Los mismos datos que hoy recibe cada acción, **menos** la identidad: desaparecen
  `usuario_id` (cuerpo de la votación y query de `GET /backlog/{id}/votacion`),
  `integrante_id` (registro de horas) y `creador` (creación de proyecto).
- Al configurar el equipo: la persona elegida de la lista de usuarios registrados.

## Salidas esperadas

- Acción permitida: la misma respuesta que hoy, a nombre de la persona con sesión.
- Sin sesión: 401 (mensaje de «no autenticado» de HU-20).
- Product Builder en acción de gestión: 403 con un mensaje que empieza con «Solo el Scrum
  Master puede…».
- No integrante en una ruta que exige pertenencia: 403 «No sos integrante de este
  proyecto».
- «Mis proyectos»: la lista de los proyectos que integra la persona, cada uno con su rol;
  vacía si no integra ninguno.
- Lista de usuarios registrados: identificador y nombre de cada uno, sin datos sensibles.

## Qué se protege *(definido en el backlog, no se reabre)*

| Nivel | Rutas |
|---|---|
| **Sin sesión** | `POST /auth/registro`, `POST /auth/ingreso`, `GET /health`, el `index.html` con los archivos estáticos del frontend, `GET /demo/estado` y `POST /demo/datos` (protegidas por `DEMO_HABILITADO`). |
| **Solo sesión** | `POST /projects` (quien lo crea sale de la sesión y queda como Scrum Master) y `GET /users`. |
| **Sesión e integrante** | `POST /backlog`, `PUT /backlog/{id}`, `PUT /backlog/{id}/estado`, `PUT /backlog/{id}/sprint`, `POST /backlog/{id}/horas`, `POST /projects/{id}/defectos`, `POST /defectos/{id}/resolver`, `POST /backlog/{id}/votacion/votos`. |
| **Sesión y Scrum Master** | `POST /sprints`, `POST /sprints/{id}/iniciar`, `POST /sprints/{id}/cerrar`, `POST /backlog/{id}/votacion` (iniciar), `POST /backlog/{id}/votacion/revelar`, `POST /backlog/{id}/votacion/rondas` (nueva ronda), `POST /backlog/{id}/votacion/estimacion` (registrar la estimación, HU-14), `PUT /projects/{id}` y `POST /projects/{id}/members`. |
| **Lecturas (`GET`)** | Todas exigen sesión (401). `GET /projects` es «Mis proyectos»: lista filtrada, nunca 403. La pertenencia (403) se comprueba **solo** en `GET /projects/{id}` (panel) y `GET /backlog/{id}/votacion` (voto secreto). El resto de las lecturas exige solo sesión. |

Las rutas de sesión de HU-20 (`GET /auth/yo`, `POST /auth/salir`) siguen como las dejó esa
historia.

**Límite conocido (se declara en la documentación)**: con sesión iniciada, quien no integra
un proyecto puede leer sus demás datos (historias, Sprints, horas, defectos, métricas,
historial, reportes).

## Fuera de alcance

- Traspasar el rol de Scrum Master, quitarlo o cambiarle el rol: no se puede hasta después
  de la entrega (límite declarado; HU-16).
- El resto de la gestión de usuarios (crear, editar): HU-17, a redefinir.
- Eliminar un proyecto (#232): depende de estos permisos, pero es otra historia.
- Restringir las lecturas a integrantes más allá de las dos rutas indicadas.

## User Scenarios & Testing *(mandatory)*

Los escenarios del backlog (1 a 9) se mapean a historias según su prioridad. El texto de
cada escenario del backlog se conserva.

### User Story 1 - Identidad de la sesión (Priority: P1) — Escenario 7

Toda acción que hoy recibe la identidad en la petición pasa a tomarla de la sesión. Nadie
puede votar, registrar horas, consultar una votación ni crear un proyecto a nombre de otra
persona, y nadie puede ver el voto oculto de otra.

**Why this priority**: hoy cualquiera puede ver el voto oculto de otra persona pasando su
`usuario_id`; es el agujero más grave y la base para los permisos por rol.

**Independent Test**: con dos personas integrantes, una vota; la otra consulta la votación
y no ve ese voto oculto, y aunque mande `usuario_id` de la primera, el resultado es el
mismo.

**Acceptance Scenarios** *(Escenario 7 del backlog)*:

1. **Given** que inicié sesión, **When** voto, **Then** el voto queda a mi nombre.
2. **Given** que inicié sesión, **When** registro horas en una historia, **Then** las horas
   quedan a mi nombre.
3. **Given** que inicié sesión e integro el proyecto, **When** consulto la votación de una
   historia, **Then** veo mi propio voto y no el voto oculto de otra persona, aunque mande
   `usuario_id` de otra persona en la query.
4. **Given** que inicié sesión, **When** creo un proyecto, **Then** quedo como su Scrum
   Master, aunque mande `creador` con otro nombre en el cuerpo.
5. **Given** que inicié sesión, **When** mando `usuario_id` o `integrante_id` de otra
   persona en el cuerpo de un voto o de un registro de horas, **Then** no se acepta otro
   usuario: la acción queda a mi nombre (ver FR-009).

---

### User Story 2 - Acciones de gestión solo para el Scrum Master (Priority: P1) — Escenarios 1 y 2

Crear, iniciar y cerrar Sprints; iniciar y revelar la votación; abrir una nueva ronda;
registrar la estimación acordada; editar el proyecto y configurar el equipo: solo el Scrum
Master del proyecto.

**Why this priority**: es el objetivo de la historia (segregación de funciones).

**Independent Test**: con un proyecto con Scrum Master y un Product Builder, cada acción de
gestión se realiza con la sesión del Scrum Master y responde 403 con la del Product Builder.

**Acceptance Scenarios**:

1. *(Escenario 1)* **Given** que soy el Scrum Master del proyecto, **When** creo, inicio o
   cierro un Sprint, inicio o revelo una votación, inicio una nueva ronda, registro la
   estimación acordada, edito el proyecto o configuro el equipo, **Then** la acción se
   realiza.
2. *(Escenario 2)* **Given** que soy Product Builder, **When** intento una de esas acciones,
   **Then** el sistema responde 403 y la pantalla la muestra deshabilitada con la nota
   «Solo el Scrum Master puede…».
3. **Given** que soy Product Builder, **When** el sistema rechaza una acción de gestión,
   **Then** no cambia nada (ni Sprint, ni votación, ni proyecto, ni equipo).

---

### User Story 3 - Sin sesión (Priority: P1) — Escenario 3

**Why this priority**: sin esto, cualquiera con `curl` saltea la pantalla.

**Independent Test**: llamar a cada ruta protegida sin token: 401; llamar a cada ruta
pública sin token: funciona.

**Acceptance Scenarios**:

1. *(Escenario 3)* **Given** que no inicié sesión, **When** uso una de las acciones que
   exigen sesión, **Then** el sistema responde 401.
2. **Given** que no inicié sesión, **When** uso una ruta pública (registro, ingreso,
   `GET /health`, la página y sus estáticos, `GET /demo/estado`, `POST /demo/datos`),
   **Then** funciona igual que antes.
3. **Given** una sesión vencida o un token inválido, **When** uso una acción que exige
   sesión, **Then** el sistema responde 401, igual que sin token.

---

### User Story 4 - Pertenencia al proyecto (Priority: P1) — Escenario 4

**Why this priority**: sin pertenencia, cualquiera con cuenta podría escribir en proyectos
ajenos y ver votos ocultos.

**Independent Test**: una persona con sesión que no integra el proyecto recibe 403 «No sos
integrante de este proyecto» en el panel, en la votación y en cada escritura que exige
pertenencia.

**Acceptance Scenarios**:

1. *(Escenario 4)* **Given** que no soy integrante del proyecto, **When** intento acceder a
   él, **Then** el sistema responde 403 con el mensaje «No sos integrante de este
   proyecto».
2. **Given** que no soy integrante del proyecto, **When** intento una escritura de los
   niveles «sesión e integrante» o «sesión y Scrum Master» sobre ese proyecto, **Then**
   recibo 403 «No sos integrante de este proyecto» y no cambia nada.
3. **Given** que no soy integrante del proyecto, **When** leo sus historias, Sprints, horas
   o defectos, **Then** puedo leerlos (límite conocido, declarado).

---

### User Story 5 - Mis proyectos (Priority: P2) — Escenarios 5 y 9

**Why this priority**: es la puerta de entrada de cada persona; sin filtro, la lista mezcla
proyectos ajenos donde todo da 403.

**Independent Test**: dos personas con proyectos distintos consultan la lista y cada una ve
solo los suyos, con su rol; una persona sin proyectos ve el mensaje y puede crear uno.

**Acceptance Scenarios**:

1. *(Escenario 5)* **Given** que inicié sesión, **When** consulto mis proyectos, **Then**
   veo solo aquellos donde soy integrante, con mi rol.
2. *(Escenario 9)* **Given** que inicié sesión y no soy integrante de ningún proyecto,
   **When** consulto mis proyectos, **Then** el sistema muestra un mensaje de que no hay
   proyectos y me permite crear uno.
3. **Given** que acabo de crear un proyecto, **When** consulto mis proyectos, **Then**
   aparece con el rol Scrum Master.

---

### User Story 6 - Configurar el equipo con usuarios registrados (Priority: P2) — Escenarios 6 y 8

**Why this priority**: completa la segregación: el equipo solo crece con Product Builders,
y el único Scrum Master no se duplica.

**Independent Test**: el Scrum Master elige de la lista a un usuario registrado y queda
como Product Builder; un usuario inexistente o un segundo Scrum Master se rechazan.

**Acceptance Scenarios**:

1. *(Escenario 8)* **Given** que soy el Scrum Master del proyecto, **When** configuro el
   equipo, **Then** elijo de una lista a un usuario ya registrado, y entra como Product
   Builder; si el usuario no existe, el sistema lo rechaza.
2. *(Escenario 6)* **Given** que quien crea un proyecto queda como Scrum Master, **When**
   se intenta agregar un segundo Scrum Master, **Then** el sistema lo rechaza.
3. **Given** que inicié sesión, **When** pido la lista de usuarios registrados, **Then**
   recibo el identificador y el nombre de cada uno, sin contraseña ni datos de sesión.
4. **Given** que soy el Scrum Master, **When** agrego a quien ya integra el proyecto,
   **Then** se aplica la regla de integrante duplicado que ya existe.

---

### User Story 7 - Pantallas según el rol (Priority: P2, frontend)

**Why this priority**: es la cara visible de los permisos; el backend es el que los hace
cumplir.

**Independent Test**: con sesión de Product Builder, las acciones de gestión aparecen
deshabilitadas con la nota; con sesión de Scrum Master, habilitadas; un no integrante ve la
pantalla «No sos integrante».

**Acceptance Scenarios**:

1. **Given** que soy Product Builder en un proyecto, **When** veo una pantalla con acciones
   de gestión, **Then** cada una aparece deshabilitada con la nota «Solo el Scrum Master
   puede…».
2. **Given** que soy Scrum Master, **When** veo esas pantallas, **Then** las acciones están
   habilitadas.
3. **Given** que voto o registro horas, **When** veo el formulario, **Then** ya no hay
   selector de «Quién actúa» ni de integrante: actúo yo.
4. **Given** que abro un proyecto que no integro, **When** la pantalla recibe el 403 de
   pertenencia, **Then** muestra la pantalla «No sos integrante de este proyecto».
5. **Given** que la pantalla recibe un 403 de rol en una acción (por ejemplo, porque el
   rol cambió en otra pestaña), **When** ocurre, **Then** muestra el aviso «No tenés
   permiso» junto a la acción, con el mensaje del servidor, sin salir de la pantalla ni
   perder lo escrito.
6. **Given** «Mis proyectos», **When** se muestra, **Then** cada proyecto lleva el rol de la
   persona; sin proyectos, el mensaje y la opción de crear uno.

### Edge Cases

- Un Scrum Master de un proyecto que es Product Builder en otro: el rol se evalúa **por
  proyecto**; en el segundo recibe 403 en las acciones de gestión.
- Una acción de gestión sobre un recurso inexistente (Sprint, historia, proyecto): 404, como
  hoy (ver orden de errores en FR-006).
- Una acción sobre una historia: el proyecto que se evalúa es el de la historia, no uno
  enviado en la petición.
- `POST /sprints` y `POST /backlog` reciben el proyecto en el cuerpo: se evalúa la
  pertenencia (y el rol, en Sprints) sobre ese proyecto.
- Un Product Builder con datos inválidos en una acción de gestión: recibe 403, no el error
  de validación (no se le revela nada de la acción). Única excepción: un cuerpo que no es
  JSON válido responde 400 «JSON inválido» (FR-006).
- La persona que creó el proyecto es la única Scrum Master; no puede dejar de serlo en esta
  historia.
- `GET /projects` nunca responde 403, aunque la persona no integre ningún proyecto.
- `GET /users` no incluye contraseñas, huellas de token ni sesiones.

## Requirements *(mandatory)*

### Reglas de negocio (EARS)

**Ubicuas**

- **FR-001**: El sistema SHALL tomar la identidad de quien actúa exclusivamente de la
  sesión.
- **FR-002**: El sistema SHALL evaluar el rol de una persona por proyecto, según su
  vínculo de integrante.
- **FR-003**: Cada proyecto SHALL tener exactamente un Scrum Master, que es quien lo creó
  (#236).
- **FR-004**: El sistema SHALL agregar a todo integrante nuevo con el rol Product Builder.
- **FR-005**: La lista de usuarios registrados SHALL incluir solo usuarios con contraseña
  (cuenta de HU-20) y contener solo identificador y nombre de cada uno.
- **FR-006**: Ante varias condiciones de error a la vez, el sistema SHALL responder la
  primera en este orden: 401 (sin sesión) → 404 (recurso inexistente) → 403 no integrante
  → 403 no Scrum Master → 400/409 (validación y conflictos del caso de uso). Excepción: un
  cuerpo que no es JSON válido responde 400 «JSON inválido» antes del 403, porque en
  algunas rutas (`POST /sprints`, `POST /backlog`) el proyecto viaja en el cuerpo y sin
  leerlo no se sabe sobre qué proyecto evaluar; esa respuesta no revela nada de la acción.

**Basadas en evento**

- **FR-007**: WHEN una persona con sesión crea un proyecto, el sistema SHALL registrarla como
  su Scrum Master.
- **FR-008**: WHEN una persona con sesión vota o registra horas, el sistema SHALL guardar el
  voto o las horas a su nombre.
- **FR-009**: WHEN una petición trae `usuario_id`, `integrante_id` o `creador`, el sistema
  SHALL ignorar ese dato (sin responder 400) y usar la persona de la sesión.
- **FR-010**: WHEN una persona con sesión consulta la votación de una historia, el sistema
  SHALL devolverle su propio voto y nunca el voto oculto de otra persona.
- **FR-011**: WHEN una persona con sesión consulta «Mis proyectos», el sistema SHALL
  devolver solo los proyectos que integra, cada uno con su rol en él.
- **FR-012**: WHEN el Scrum Master agrega un integrante, el sistema SHALL aceptar solo un
  usuario registrado y vincularlo como Product Builder.

**Basadas en estado**

- **FR-013**: WHILE la persona no integra ningún proyecto, «Mis proyectos» SHALL devolver una
  lista vacía y la pantalla SHALL mostrar un mensaje de que no hay proyectos y la opción de
  crear uno.
- **FR-014** *(frontend)*: WHILE la persona es Product Builder en el proyecto abierto, la
  pantalla SHALL mostrar las acciones de gestión deshabilitadas con la nota «Solo el Scrum
  Master puede…».

**No deseadas**

- **FR-015**: IF una petición a una ruta protegida llega sin sesión válida (ausente,
  inválida o vencida), THEN el sistema SHALL responder 401 sin ejecutar la acción.
- **FR-016**: IF una persona que no es Scrum Master del proyecto intenta una acción de
  gestión, THEN el sistema SHALL responder 403 con un mensaje que empiece con «Solo el
  Scrum Master puede» y no cambiar nada.
- **FR-017**: IF una persona que no integra el proyecto intenta una escritura sobre él, o
  consulta `GET /projects/{id}` o `GET /backlog/{id}/votacion`, THEN el sistema SHALL
  responder 403 «No sos integrante de este proyecto» y no cambiar nada.
- **FR-018**: IF el Scrum Master intenta agregar un usuario que no existe o que no tiene
  contraseña, THEN el sistema SHALL rechazarlo igual en ambos casos, sin agregar a nadie.
- **FR-019**: IF se intenta agregar un integrante con rol Scrum Master, THEN el sistema SHALL
  rechazarlo con el conflicto de #236 (409) sin agregar a nadie.
- **FR-020**: *(Retirado en `/speckit.analyze`: quitar o cambiar el rol del único Scrum
  Master no es una regla a construir sino un límite; ver «Fuera de alcance».)*

**Opcionales**

- **FR-021**: WHERE `DEMO_HABILITADO` está activo, `GET /demo/estado` y `POST /demo/datos`
  SHALL funcionar sin sesión.

### Requisitos de la interfaz *(frontend)*

- **FR-022**: La interfaz MUST quitar el selector «Quién actúa» de la votación (#229) y el de
  integrante del registro de horas (#227).
- **FR-023**: La interfaz MUST mostrar «Mis proyectos» con el rol de la persona en cada uno.
- **FR-024**: La interfaz MUST mostrar, según el 403 recibido, la pantalla completa «No sos
  integrante de este proyecto» o el aviso en la misma pantalla «No tenés permiso» (con el
  mensaje del servidor, sin perder lo escrito).
- **FR-025**: La interfaz MUST, al configurar el equipo, ofrecer la lista de usuarios
  registrados para elegir, sin selector de rol (siempre Product Builder).
- **FR-026**: La interfaz MUST dejar de enviar `creador` al crear un proyecto.

### Requisitos de prueba

- **FR-027**: Los escenarios de Godog MUST usar un paso de contexto «Dado que inicié sesión
  como X» que guarda el token y lo envía en cada llamada; los escenarios existentes que
  mandan `usuario_id`, `integrante_id` o `creador` MUST ajustarse **en este mismo PR**, con
  el check `backend` en verde.
- **FR-028**: El código de producción MUST NOT contener atajos de prueba (usuarios fijos,
  cabeceras de prueba, desactivar la autenticación por configuración).

### Key Entities

- **Usuario**: persona registrada (HU-20). Se expone en la lista solo con identificador y
  nombre.
- **Sesión**: identifica a quien actúa (HU-20). No cambia.
- **Integrante del proyecto**: vínculo usuario–proyecto con rol (Scrum Master o Product
  Builder). Ya existe (HU-04, #236); es la fuente de verdad de la pertenencia y del rol. No
  hay cambios de esquema previstos en esta historia (`MODELO-DE-DATOS.md`, decisión D4 se
  cierra y D14 deja de hacer falta).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Los 9 escenarios del backlog para HU-21 pasan automatizados (Godog para la API,
  pruebas de componente para la interfaz).
- **SC-002**: El 100 % de las rutas de la tabla «Qué se protege» responde según su nivel:
  las protegidas devuelven 401 sin sesión y las públicas funcionan sin ella.
- **SC-003**: En el 100 % de las acciones de gestión, un Product Builder recibe 403 y el
  estado del proyecto no cambia.
- **SC-004**: En ningún caso una persona obtiene el voto oculto de otra, cualquiera sea el
  dato que mande en la petición.
- **SC-005**: El 100 % de los escenarios de Godog (los existentes, ajustados, y los nuevos)
  pasan en el CI, sin atajos de prueba en producción.
- **SC-006**: Una persona ve en «Mis proyectos» exclusivamente los proyectos que integra.

## Assumptions

- **Fuente de verdad**: la sección HU-21 de `docs/PRODUCT-BACKLOG.md` y el issue #234. La
  tabla «Qué se protege», los códigos y «Mis proyectos» no se reabren en `/speckit.clarify`.
- Las rutas de HU-20 (`GET /auth/yo`, `POST /auth/salir`) conservan su comportamiento.
- El campo de rol al agregar un integrante es opcional: ausente o `product_builder` entra
  como Product Builder; `scrum_master` se rechaza con el 409 de #236 (escenario 6).
- Ninguna ruta crea ya usuarios sin contraseña y la base de Render se vacía en cada deploy,
  así que los usuarios sin contraseña solo pueden existir en bases locales viejas.
- Sin cambios de esquema: la pertenencia y el rol salen de la tabla de integrantes ya
  existente. Si el plan detecta que hace falta un cambio, se frena y se consulta al Scrum
  Master.
- **Plan B** (backlog): si el 13-nov los permisos del backend no están completos, se entrega
  con las acciones ocultas en pantalla y se declara el límite; eso no cumple el Definition
  of Done (quedan incumplidos los escenarios 2, 3 y 7).
- Conviene que #229, #227 y #231 estén mergeados antes de implementar.
