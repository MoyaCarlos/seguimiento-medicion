# Feature Specification: Navegar entre proyectos desde la interfaz

**Feature Branch**: `feature/148-navegacion-proyectos`

**Issue**: #148 (cierra además #141)

**Created**: 2026-10-27

**Status**: Draft

**Input**: User description: "Issue #148: Navegar entre proyectos desde la interfaz (lista y selector). Solo frontend; no cambia la API ni el esquema. Esta historia modifica archivos de la base (compartido/router/LayoutProyecto.tsx y web/src/rutas.test.tsx) porque el issue lo asigna expresamente.

Alcance:
- Selector de proyecto en el encabezado de LayoutProyecto: lista los proyectos con GET /projects, marca el actual (el del :id de la URL) y, al elegir otro, navega a la misma sección del nuevo proyecto (/projects/:id/<seccion>, con la sección tomada del primer tramo después del id).
- Enlace a "Nuevo proyecto" (la pantalla de /) en ese mismo encabezado.
- En la pantalla de inicio (/), junto al formulario de crear proyecto, la lista de proyectos existentes con un enlace a cada uno.
- Estados de la lista: cargando (LoadingIndicator), vacía (invita a crear un proyecto) y error de la API (ErrorMessage) sin romper la pantalla: el encabezado y la navegación de secciones siguen funcionando.
- Al cambiar de proyecto, el estado de las pantallas se reinicia por el :id (ya existe el reinicio por key; no romperlo).
- Se elimina web/src/compartido/router/Inicio.tsx (sin uso, issue #141) y este trabajo cierra también #141.

Restricciones: sin dependencias nuevas; reutilizar LoadingIndicator, ErrorMessage y el cliente de compartido/api; Button no aplica porque esta historia solo navega con enlaces y un <select>; la lógica pura (por ejemplo, armar la ruta del nuevo proyecto a partir de la ruta actual) en funciones con pruebas de Vitest, y las pantallas con pruebas de componente con Testing Library y fetch simulado. Las pruebas existentes de rutas.test.tsx deben seguir pasando (con un fetch simulado de GET /projects)."

## Contexto

Hoy todas las pantallas de un proyecto viven bajo `/projects/:id/<sección>` y comparten un
encabezado con los enlaces a Proyecto, Historias, Sprints y Estado. Ya existe la creación de
proyectos en `/` y la API devuelve el listado completo con `GET /projects` (#109), pero la
interfaz no lo usa: para pasar de un proyecto a otro hay que editar la dirección a mano, y la
pantalla de inicio solo muestra el formulario. Esta historia conecta ese listado con la
navegación, sin tocar el contrato de la API ni el esquema de la base.

## Objetivo

Permitir que una persona cambie de proyecto sin editar la URL y sin perder la sección en la que
estaba, que pueda ir a crear otro proyecto desde cualquier pantalla de proyecto, y que desde el
inicio pueda ver y abrir los proyectos existentes, con estados de carga, lista vacía y error que
no bloqueen el resto de la pantalla.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Cambiar de proyecto conservando la sección (Priority: P1)

Como integrante del equipo, estando en una sección de un proyecto, quiero elegir otro proyecto
desde el encabezado y seguir en la misma sección del proyecto elegido, para comparar o trabajar
en varios proyectos sin editar la dirección ni volver a navegar desde el principio.

**Why this priority**: Es el objetivo central del issue. Sin el selector, la interfaz sigue
obligando a escribir la URL para cambiar de proyecto.

**Independent Test**: Abrir `/projects/1/backlog` con dos proyectos cargados, elegir el proyecto
2 en el selector y comprobar que la dirección pasa a `/projects/2/backlog`, que se ve el backlog
del proyecto 2 y que el selector marca el 2 como actual.

**Acceptance Scenarios**:

1. **Given** el encabezado con la lista de proyectos y el proyecto 1 como actual, **When** se
   elige el proyecto 2, **Then** la aplicación navega a `/projects/2/<sección actual>` y el
   selector marca el proyecto 2 como actual.
2. **Given** la pantalla en `/projects/1/proyecto/edit`, **When** se elige el proyecto 2,
   **Then** la aplicación navega a `/projects/2/proyecto` (solo se conserva el primer tramo
   después del id, no la sub-ruta de edición).
3. **Given** la pantalla abierta en cualquier sección de un proyecto, **When** se despliega el
   selector, **Then** el proyecto de la URL aparece marcado como actual y los demás se listan
   como `<nombre> (#<id>)`.
4. **Given** una pantalla con estado propio del proyecto 1, **When** se elige el proyecto 2,
   **Then** las pantallas del proyecto 2 arrancan sin datos ni estado del proyecto 1.

---

### User Story 2 - Ir a crear un proyecto desde cualquier proyecto (Priority: P1)

Como integrante del equipo, desde cualquier sección de un proyecto quiero un enlace visible a
"Nuevo proyecto" que me lleve a la pantalla de creación, para dar de alta otro proyecto sin
editar la dirección ni volver al inicio a mano.

**Why this priority**: Es un acceso directo pedido por el issue y se apoya en la pantalla de
creación que ya existe; sin él, el selector resuelve solo la mitad de la navegación.

**Independent Test**: Abrir cualquier sección de un proyecto y pulsar "Nuevo proyecto";
comprobar que se llega a `/` y que el formulario de creación está disponible.

**Acceptance Scenarios**:

1. **Given** una pantalla de proyecto, **When** se pulsa el enlace "Nuevo proyecto" del
   encabezado, **Then** se abre la pantalla de creación en `/`.
2. **Given** que la lista de proyectos está cargando, vacía o falló, **When** se mira el
   encabezado, **Then** el enlace "Nuevo proyecto" sigue visible y funciona.

---

### User Story 3 - Ver y abrir los proyectos desde el inicio (Priority: P1)

Como integrante del equipo, en la pantalla de inicio quiero ver, junto al formulario de creación,
la lista de proyectos existentes y abrir cualquiera de ellos, para elegir un proyecto sin
recordar ni escribir su identificador.

**Why this priority**: Es la otra mitad de la navegación pedida: el inicio deja de ser solo un
formulario y pasa a ser el punto de entrada a los proyectos existentes.

**Independent Test**: Abrir `/` con proyectos cargados, comprobar que se listan junto al
formulario y que el enlace de uno de ellos abre su panel.

**Acceptance Scenarios**:

1. **Given** que existen proyectos, **When** se abre `/`, **Then** se ve la lista de proyectos
   junto al formulario de creación, cada uno como `<nombre> (#<id>)` y con un enlace para
   abrirlo.
2. **Given** la lista de proyectos en `/`, **When** se activa el enlace de un proyecto,
   **Then** se abre el panel de ese proyecto en `/projects/<id>/proyecto`.
3. **Given** que no existe ningún proyecto, **When** se abre `/`, **Then** se muestra un mensaje
   que invita a crear el primero, junto al formulario.
4. **Given** que la lista de proyectos falla, **When** se abre `/`, **Then** el formulario de
   creación sigue disponible y usable.

**Nota de trazabilidad**: el escenario 4 se verifica en las tareas de la lista de US4
(T035-T038), junto con los demás estados de error de la lista; no se duplica en las tareas de
US3.

---

### User Story 4 - Estados de la lista que no rompen la pantalla (Priority: P2)

Como integrante del equipo, quiero que la lista de proyectos muestre claramente si está
cargando, vacía o con error, sin que ese estado impida navegar entre secciones ni crear un
proyecto.

**Why this priority**: Da confianza sobre lo que se ve y evita que una falla del listado deje la
interfaz bloqueada; se apoya en las piezas compartidas que ya existen.

**Independent Test**: Simular las respuestas de `GET /projects` (demora, arreglo vacío y error) y
comprobar que cada estado se ve y que la navegación de secciones y el formulario siguen
funcionando.

**Acceptance Scenarios**:

1. **Given** que el listado todavía no respondió, **When** se abre una pantalla de proyecto,
   **Then** se ve el indicador de carga, el selector deshabilitado con la única opción
   "Proyecto <id>" del proyecto actual y los enlaces de sección del encabezado siguen
   funcionando.
2. **Given** que `GET /projects` responde con un error, **When** se abre una pantalla de
   proyecto, **Then** se ve el mensaje de error de la API, el selector deshabilitado con la
   única opción "Proyecto <id>" del proyecto actual y los enlaces de sección y "Nuevo
   proyecto" siguen funcionando.
3. **Given** que `GET /projects` devuelve una lista vacía, **When** se abre el inicio,
   **Then** se ve la invitación a crear un proyecto y no un mensaje de error.
4. **Given** que el listado falla en el inicio, **When** se completa y envía el formulario,
   **Then** la creación se comporta igual que antes de esta historia.

---

### Edge Cases

- **IF** `GET /projects` falla o devuelve una respuesta no interpretable, **THEN** la interfaz
  **shall** mostrar un mensaje de error de la API sin romper la pantalla, **shall** mantener el
  selector deshabilitado con la única opción "Proyecto `<id>`" del proyecto actual y **shall**
  mantener operativos los enlaces de sección y el enlace "Nuevo proyecto" (o el formulario, en
  el inicio).
- **IF** la lista de proyectos está vacía, **THEN** la interfaz **shall** mostrar una invitación a
  crear un proyecto, no un error.
- **WHILE** la lista de proyectos se está cargando, **THEN** la interfaz **shall** mostrar el
  indicador de carga, **shall** mostrar el selector deshabilitado con la única opción
  "Proyecto `<id>`" del proyecto actual y **shall** permitir seguir usando la navegación
  existente.
- **WHEN** la carga del listado se resuelve con éxito, **THEN** el selector **shall** habilitarse
  con la lista completa de proyectos.
- **IF** el `:id` de la URL no aparece en el listado (por ejemplo, un proyecto borrado o una
  lista incompleta), **THEN** el encabezado **shall** seguir mostrando el proyecto actual como
  "Proyecto `<id>`", **shall not** marcar otro proyecto como actual y **shall** permitir cambiar
  a cualquiera de los proyectos listados.
- **IF** la dirección es `/projects/<id>` sin sección, **THEN** la sección usada al cambiar de
  proyecto **shall** ser `proyecto`.
- **IF** la sección actual no es una de las conocidas (por ejemplo, `/projects/1/otra`),
  **THEN** al cambiar de proyecto **shall** conservarse ese mismo primer tramo
  (`/projects/2/otra`), aunque la pantalla de destino muestre "Página no encontrada".
- **WHEN** se elige el mismo proyecto que ya está abierto, **THEN** la aplicación **shall not**
  navegar ni reiniciar el estado de la pantalla.
- **IF** dos proyectos tienen el mismo nombre, **THEN** el selector y la lista **shall**
  mostrarlos como `<nombre> (#<id>)`, con identificadores distintos, y **shall** abrir el
  proyecto cuyo enlace o ítem se eligió.
- **IF** se cambia de proyecto mientras una pantalla tiene datos o estado local, **THEN** la
  pantalla del proyecto nuevo **shall** arrancar limpia, sin mezclar datos del anterior.
- **IF** el listado tarda o falla, **THEN** el enlace "Nuevo proyecto" del encabezado **shall**
  estar disponible igual.

## Requirements *(mandatory)*

### Functional Requirements

**Selector del encabezado**

- **FR-001**: Toda pantalla de proyecto (`/projects/:id/...`) MUST mostrar en su encabezado un
  control de selección de un solo proyecto que liste todos los proyectos obtenidos con
  `GET /projects`, mostrando cada uno como `<nombre> (#<id>)`.
- **FR-002**: El control MUST marcar como actual el proyecto cuyo identificador coincide con el
  `:id` de la URL, de forma visible y accesible.
- **FR-003**: **WHEN** se elige un proyecto distinto del actual, la aplicación MUST navegar a
  `/projects/<id elegido>/<sección>`, donde `<sección>` es el primer tramo de la ruta actual
  después del `:id`; si no hay ningún tramo, MUST usar `proyecto`.
- **FR-004**: Al cambiar de proyecto MUST conservarse únicamente el primer tramo de la sección
  actual; los tramos posteriores (sub-rutas) MUST NOT trasladarse al proyecto nuevo.
- **FR-005**: **WHEN** se elige el proyecto que ya está abierto, la aplicación MUST NOT navegar
  ni reiniciar el estado de la pantalla.
- **FR-006**: El encabezado de toda pantalla de proyecto MUST incluir un enlace visible
  "Nuevo proyecto" que lleve a la pantalla de creación en `/`.

**Listado en el inicio**

- **FR-007**: La pantalla de inicio (`/`) MUST mostrar, junto al formulario de creación, la
  lista de proyectos existentes obtenida con `GET /projects`.
- **FR-008**: Cada proyecto de la lista del inicio MUST mostrarse como `<nombre> (#<id>)` y
  MUST ofrecer un enlace que abra su panel en `/projects/<id>/proyecto`.
- **FR-009**: **WHEN** no hay proyectos, el inicio MUST mostrar una invitación a crear el
  primero, sin tratarlo como error.

**Estados de la lista**

- **FR-010**: **WHILE** la lista de proyectos se está obteniendo, la interfaz MUST mostrar el
  indicador de carga compartido y el selector MUST mostrarse deshabilitado con una única opción
  "Proyecto `<id>`", correspondiente al proyecto actual. **WHEN** la carga se resuelve con
  éxito, el selector MUST habilitarse con la lista completa de proyectos.
- **FR-011**: **IF** la obtención de la lista falla, la interfaz MUST mostrar el mensaje de error
  de la API mediante el componente compartido de errores, sin romper la pantalla, y el selector
  MUST permanecer visible y deshabilitado con la única opción "Proyecto `<id>`" del proyecto
  actual.
- **FR-012**: Los estados de carga, vacío y error de la lista MUST NOT impedir la navegación de
  secciones del encabezado, el enlace "Nuevo proyecto" ni el uso del formulario de creación.

**Cambio de proyecto y estado**

- **FR-013**: **WHEN** cambia el `:id` de la URL, las pantallas del proyecto MUST reiniciarse y
  MUST NOT mostrar datos ni estado del proyecto anterior; el mecanismo de reinicio por
  identificador que ya existe MUST conservarse.
- **FR-014**: El listado de proyectos MUST mostrarse en el orden devuelto por la API, sin
  reordenarlo ni paginarlo en esta historia.

**Limpieza y alcance**

- **FR-015**: MUST eliminarse `web/src/compartido/router/Inicio.tsx`, que no tiene uso, y el
  trabajo MUST cerrar el issue #141.
- **FR-016**: La historia MUST NOT modificar la API, el esquema de la base ni el contrato de
  `GET /projects`; solo consume el listado existente.
- **FR-017**: La interfaz MUST reutilizar `LoadingIndicator`, `ErrorMessage` y el cliente de
  `compartido/api`, y MUST NOT agregar dependencias nuevas. `Button` no aplica porque esta
  historia solo navega con enlaces y un `<select>`; MUST NOT agregar acciones de botón.
- **FR-018**: La lógica pura de navegación (por ejemplo, construir la ruta del proyecto elegido a
  partir de la ruta actual) MUST implementarse en funciones separadas y probarse con Vitest.
- **FR-019**: Las pantallas MUST probarse con pruebas de componente con Testing Library y
  `fetch` simulado, cubriendo los escenarios de aceptación de esta spec.
- **FR-020**: Las pruebas existentes de `web/src/rutas.test.tsx` MUST seguir pasando, con un
  `fetch` simulado que responda a `GET /projects`.

### Key Entities

- **Proyecto**: proyecto existente que se lista y se elige. Para esta historia solo se usan su
  nombre e identificador, que se muestran juntos como `<nombre> (#<id>)`; no se crea, edita ni
  borra desde la navegación.
- **Sección del proyecto**: primer tramo de la ruta después del identificador (`proyecto`,
  `backlog`, `sprints`, `status`, o el que exista). Determina a dónde se navega al cambiar de
  proyecto.
- **Pantalla de inicio**: punto de entrada en `/` que contiene el formulario de creación y, a
  partir de esta historia, la lista de proyectos existentes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Una persona puede cambiar de un proyecto a otro desde cualquier sección en dos
  interacciones o menos (abrir el selector y elegir), y el 100% de las veces termina en la misma
  sección del proyecto elegido.
- **SC-002**: El 100% de las pantallas de proyecto muestran el proyecto actual y permiten
  cambiarlo, sin editar la dirección ni recargar la página.
- **SC-003**: El 100% de las pantallas de proyecto ofrecen el acceso "Nuevo proyecto" a un clic,
  incluso cuando la lista está cargando, vacía o con error.
- **SC-004**: El inicio muestra los N proyectos existentes y cada uno abre su panel; con cero
  proyectos muestra la invitación a crear el primero.
- **SC-005**: Ante una falla del listado, el 100% de las pantallas conserva operativas la
  navegación de secciones y la creación de proyectos, y muestra un mensaje de error visible.
- **SC-006**: Al cambiar de proyecto, el 100% de las pantallas arranca sin datos del proyecto
  anterior (cero mezcla de estado entre proyectos).
- **SC-007**: La suite existente de rutas y las nuevas pruebas de lógica pura y de componente
  pasan en verde, sin dependencias nuevas.

## Assumptions

- `GET /projects` ya existe y devuelve todos los proyectos ordenados por identificador (#109);
  esta historia no agrega endpoints ni campos. El nombre y el identificador son los datos que se
  muestran en el selector y en la lista del inicio, con el formato `<nombre> (#<id>)`.
- La pantalla de creación sigue siendo la de `/`; el enlace "Nuevo proyecto" navega a ella y no
  abre un diálogo ni una pantalla nueva.
- El enlace de un proyecto desde el inicio abre su panel en `/projects/<id>/proyecto`, porque es
  la sección por defecto del proyecto.
- Las secciones que existen hoy son `proyecto`, `backlog`, `sprints` y `status`. Si la ruta trae
  otra sección, se conserva tal cual como primer tramo; que la pantalla de destino exista o no es
  responsabilidad del registro de rutas actual.
- Si el proyecto actual no está en el listado, se muestra como "Proyecto `<id>`" para no perder
  el contexto ni marcar otro proyecto por error. No se inventa un nombre.
- La lista se obtiene al abrir cada pantalla. No se exige actualización en vivo después de crear
  un proyecto: al crearlo se navega a su panel y la lista se vuelve a leer cuando corresponda.
- Los nombres repetidos se distinguen por el identificador visible en la etiqueta
  `<nombre> (#<id>)`; la marca de actual se basa en el `:id`, no en el nombre.
- BDD (Principio III de la constitución) cumple por equivalencia: los escenarios de aceptación
  Given-When-Then de esta spec se automatizan como pruebas de componente (Vitest + Testing
  Library), con `fetch` simulado.
- El ciclo RED/GREEN/REFACTOR aplica a la lógica pura y también a cada comportamiento de las
  pantallas (navegación, carga, vacío y error).
- El control de selección es de un solo proyecto, con etiqueta accesible y operable por teclado;
  no se especifica si es una lista desplegable nativa u otro control equivalente.

## Restricciones de entrega

- El trabajo es exclusivamente de frontend. No se modifican handlers, servicios, repositorios,
  migraciones ni el esquema de la base.
- Archivos permitidos: `web/src/compartido/router/LayoutProyecto.tsx`, `web/src/rutas.tsx` (para
  componer el inicio con el formulario de creación más la lista de proyectos),
  `web/src/rutas.test.tsx`, archivos nuevos en `web/src/compartido/proyectos/` (API, hook,
  selector, lista y la función pura de ruta) y la eliminación de
  `web/src/compartido/router/Inicio.tsx`.
- `web/src/features/proyecto/crear.tsx` y `web/src/features/proyecto/crear.test.tsx` NO se
  modifican: sus pruebas afirman cuántas veces se llama a `fetch` y agregar allí el
  `GET /projects` las rompería.
- La eliminación de `web/src/compartido/router/Inicio.tsx` no debe dejar importaciones ni
  referencias a ese archivo.
- Sin dependencias nuevas. Se reutilizan `LoadingIndicator`, `ErrorMessage` y el cliente de la
  API; no hay acciones de botón (solo enlaces y un `<select>`).
- Las pruebas de `web/src/rutas.test.tsx` deben simular `GET /projects` para que el encabezado
  pueda renderizarse en todas las rutas existentes.

## Fuera de alcance

- Cambios en `GET /projects` o en cualquier otro endpoint; paginación, búsqueda, filtros u
  ordenamiento distinto al de la API.
- Crear, editar o borrar proyectos desde el selector o desde la lista del inicio.
- Mostrar métricas, estado o datos del proyecto en el selector o en la lista.
- Autenticación, autorización por rol o selección de integrante.
- Persistencia de la última sección o del último proyecto visitado.
- Rediseño general del encabezado, de la navegación de secciones o del formulario de creación.
- Actualización automática de la lista después de crear un proyecto.

## Trazabilidad

| Requisito | Verificación prevista |
|---|---|
| FR-001, FR-002 | Prueba de componente del encabezado con listado simulado, etiqueta `<nombre> (#<id>)` y `:id` actual |
| FR-003, FR-004, FR-005 | Pruebas de Vitest de la función pura que arma la ruta + prueba de navegación del selector |
| FR-006, FR-012 | Prueba de componente: enlace "Nuevo proyecto" visible y operativo en cada estado |
| FR-007, FR-008, FR-009 | Pruebas de componente del inicio con lista con datos (etiqueta `<nombre> (#<id>)`) y lista vacía |
| FR-010, FR-011, FR-012 | Pruebas de componente con `fetch` demorado, fallido y exitoso; selector deshabilitado con la única opción "Proyecto `<id>`" |
| US3 escenario 4 (falla del listado en `/`) | Se verifica en las tareas de la lista de US4 (T035-T038), no en las tareas de US3 |
| FR-013 | Prueba de navegación entre dos `:id` distintos, sin datos del anterior (incluye la prueba existente de reinicio) |
| FR-014 | Prueba de componente que conserva el orden devuelto por la API |
| FR-015 | Búsqueda de referencias a `Inicio.tsx` y compilación del frontend |
| FR-016, FR-017, FR-020 | Revisión del diff, suite existente y ausencia de dependencias nuevas |
| FR-018, FR-019 | Pruebas de Vitest de funciones puras y pruebas de componente con Testing Library |
