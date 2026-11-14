# Feature Specification: Cerrar y reabrir un proyecto (frontend)

**Feature Branch**: `feature/312-cerrar-reabrir-frontend`

**Created**: 2026-11-13

**Status**: Draft

**Input**: User description: "#312: frontend de cerrar y reabrir un proyecto (#300). Solo frontend: el backend (spec 044, PR #311) ya está en main. Zona «Cerrar proyecto» dentro de Editar proyecto, solo para el Scrum Master, con confirmación; botón «Reabrir proyecto» si está cerrado; aviso de solo lectura «Proyecto cerrado el AAAA-MM-DD»; acciones de escritura deshabilitadas con una nota en un proyecto cerrado."

## Clarifications

### Session 2026-11-13

- Q: ¿Qué pantallas deshabilitan sus acciones de escritura en un proyecto cerrado? → A: Todas las pantallas del proyecto que escriben (autenticación queda fuera), entregadas en fases independientes por pantalla y por prioridad: 1) zona «Cerrar proyecto», aviso y Editar proyecto; 2) Sprints (con horas y estado de las historias); 3) Historias, Votación, Defectos y Equipo. Se puede cortar por fase y entregar el resto en otro PR. El estado «cerrado» se toma del mismo dato del proyecto que ya consulta cada pantalla, sin pedidos nuevos.

- Q: ¿Cómo se confirma el cierre? → A: Confirmación en línea dentro de la zona: el botón «Cerrar proyecto» pasa a «¿Cerrar el proyecto? Confirmar / Cancelar». El foco va a «Cancelar» (la opción segura); Escape cancela; mientras se envía el pedido, ambos botones quedan deshabilitados (sin doble envío); si el servidor responde 409 (Sprint activo o ya cerrado), se muestra su mensaje junto a la zona y la confirmación se cierra.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - El Scrum Master cierra un proyecto (Priority: P1)

El Scrum Master abre «Editar proyecto» y, en una zona «Cerrar proyecto», cierra el proyecto. Antes de cerrar se le pide confirmación. Si el servidor rechaza el cierre (hay un Sprint Activo, o el proyecto ya estaba cerrado), ve el mensaje del servidor tal cual y el proyecto no cambia.

**Why this priority**: es el núcleo de la funcionalidad: sin poder cerrar, nada de lo demás tiene sentido.

**Independent Test**: con un proyecto abierto y sin Sprint Activo, entrar como Scrum Master a Editar proyecto, cerrar y confirmar; el proyecto pasa a «Finalizado» con su fecha de cierre.

**Acceptance Scenarios**:

1. **Given** un proyecto abierto sin Sprint Activo y el usuario es Scrum Master, **When** pide cerrar y confirma, **Then** el proyecto queda cerrado y la pantalla muestra «Proyecto cerrado el AAAA-MM-DD» con estado «Finalizado».
2. **Given** la confirmación de cierre abierta, **When** el usuario cancela, **Then** no se envía nada y el proyecto sigue abierto.
3. **Given** un proyecto con un Sprint Activo, **When** el Scrum Master confirma el cierre, **Then** ve «Cerrá el Sprint activo antes de cerrar el proyecto» y el proyecto sigue abierto.
4. **Given** un proyecto que otra persona ya cerró (la pantalla estaba desactualizada), **When** el Scrum Master confirma el cierre, **Then** ve «El proyecto ya está cerrado».

---

### User Story 2 - Un proyecto cerrado se ve como solo lectura (Priority: P1)

En las pantallas del proyecto aparece un aviso «Proyecto cerrado el AAAA-MM-DD», y las acciones de escritura de las pantallas que lo respetan quedan deshabilitadas con una nota que explica por qué. Las lecturas siguen funcionando.

**Why this priority**: sin esto, un proyecto cerrado invitaría a acciones que el servidor rechaza con 409.

**Independent Test**: abrir un proyecto cerrado y comprobar que se ve el aviso, que se puede navegar y leer todo, y que las acciones de escritura están deshabilitadas con su nota.

**Acceptance Scenarios**:

1. **Given** un proyecto cerrado, **When** cualquier integrante abre una pantalla del proyecto, **Then** ve el aviso «Proyecto cerrado el AAAA-MM-DD» y puede seguir leyendo los datos.
2. **Given** un proyecto cerrado, **When** se muestra una acción de escritura, **Then** está deshabilitada y una nota visible explica que el proyecto está cerrado.
3. **Given** un proyecto abierto, **When** se muestran las pantallas, **Then** no hay aviso y las acciones se comportan como antes (los usos actuales de permisos de HU-21 no cambian).
4. **Given** una acción que igual llega al servidor sobre un proyecto cerrado, **When** el servidor responde 409 «El proyecto está cerrado», **Then** se muestra ese mensaje tal cual.

---

### User Story 3 - El Scrum Master reabre un proyecto cerrado (Priority: P2)

En la misma zona, con el proyecto cerrado, el Scrum Master ve un botón «Reabrir proyecto». Al reabrir, el aviso desaparece, el estado se recalcula y vuelve a poder escribir.

**Why this priority**: corrige cierres por error; es menos frecuente que cerrar.

**Independent Test**: con un proyecto cerrado, reabrir como Scrum Master; el aviso desaparece y las acciones se habilitan.

**Acceptance Scenarios**:

1. **Given** un proyecto cerrado y el usuario es Scrum Master, **When** pulsa «Reabrir proyecto», **Then** el proyecto queda abierto, sin fecha de cierre, con el estado recalculado, y las acciones de escritura se habilitan.
2. **Given** un proyecto que otra persona ya reabrió, **When** el Scrum Master pulsa «Reabrir proyecto», **Then** ve «El proyecto no está cerrado».

---

### User Story 4 - El Product Builder no puede cerrar ni reabrir (Priority: P2)

El Product Builder no ve la zona «Cerrar proyecto», o la ve con sus botones deshabilitados y una nota que lo explica.

**Why this priority**: respeta el modelo de permisos de HU-21; la protección real es del servidor.

**Independent Test**: entrar como Product Builder a Editar proyecto y comprobar que no puede cerrar ni reabrir.

**Acceptance Scenarios**:

1. **Given** el usuario es Product Builder, **When** abre Editar proyecto, **Then** no puede cerrar ni reabrir (zona oculta o deshabilitada con nota; ver FR-009).
2. **Given** un 403 al cerrar o reabrir, **When** llega la respuesta, **Then** se muestra el aviso de falta de permiso.
3. **Given** un 404 al cerrar o reabrir, **When** llega la respuesta, **Then** se muestra el aviso de que el proyecto no existe.

---

### Edge Cases

- Cerrar o reabrir con la conexión caída: se muestra el error de conexión y el estado no cambia.
- Doble clic en confirmar: se envía una sola solicitud (los botones se deshabilitan mientras el pedido está en curso).
- Escape con la confirmación abierta: se cancela sin enviar nada.
- Mientras no se conoce el rol o el estado de cierre (cargando), no se bloquea de más: decide el servidor, igual que en HU-21.
- Una acción deshabilitada por cierre y por falta de permiso a la vez: se muestra una sola nota, la más relevante.
- Cambiar de proyecto en la misma sesión: el estado de cierre no se arrastra de uno a otro.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: La pantalla Editar proyecto MUST mostrar una zona «Cerrar proyecto» al Scrum Master.
- **FR-002**: El cierre MUST pedir confirmación antes de enviarse, y MUST poder cancelarse sin efecto. La confirmación es en línea dentro de la zona: «¿Cerrar el proyecto?» con «Confirmar» y «Cancelar». Al mostrarse, el foco MUST ir a «Cancelar»; Escape MUST cancelar; mientras el pedido está en curso, ambos botones MUST estar deshabilitados; ante un 409 del servidor, la confirmación MUST cerrarse y el mensaje mostrarse junto a la zona.
- **FR-003**: Los errores 409 de cerrar y reabrir MUST mostrarse con el mensaje del servidor, sin reescribirlo.
- **FR-004**: Con el proyecto cerrado, la misma zona MUST ofrecer «Reabrir proyecto» al Scrum Master, sin confirmación previa (acción reversible y poco riesgosa).
- **FR-005**: Tras cerrar o reabrir con éxito, la interfaz MUST reflejar de inmediato la fecha de cierre y el estado devueltos por el servidor, sin recargar.
- **FR-006**: Las pantallas del proyecto MUST mostrar el aviso «Proyecto cerrado el AAAA-MM-DD» cuando el proyecto tiene fecha de cierre, y no mostrarlo si no la tiene.
- **FR-007**: En un proyecto cerrado, las acciones de escritura de todas las pantallas del proyecto que escriben (excepto autenticación) MUST estar deshabilitadas con una nota visible asociada. Se entrega en tres fases independientes, cada una completa y probable por separado: **Fase 1** zona «Cerrar proyecto», aviso y Editar proyecto; **Fase 2** Sprints (incluidas horas y estado de las historias); **Fase 3** Historias, Votación, Defectos y Equipo.
- **FR-007a**: Cada pantalla MUST obtener el estado «cerrado» del mismo dato del proyecto que ya consulta, sin pedidos nuevos al servidor.
- **FR-008**: Las lecturas MUST seguir funcionando en un proyecto cerrado.
- **FR-009**: El Product Builder MUST NOT poder cerrar ni reabrir; se le oculta la zona o se le muestra deshabilitada con nota (se decide en planificación, siguiendo el patrón de HU-21 para acciones de gestión).
- **FR-010**: El estado «cerrado» MUST sumarse a los permisos compartidos de forma retrocompatible: los usos actuales de HU-21 no cambian su comportamiento.
- **FR-011**: Un 403 o un 404 al cerrar o reabrir MUST mostrar el aviso correspondiente (sin permiso / proyecto inexistente).
- **FR-012**: Esta funcionalidad MUST ser solo de frontend: sin cambios de API, esquema ni backend, y sin tocar App.tsx, App.css, rutas.tsx ni LayoutProyecto.tsx; los componentes nuevos van en archivos nuevos.

### Key Entities

- **Proyecto**: ahora expone su fecha de cierre (vacía si está abierto, fecha AAAA-MM-DD si está cerrado) y su estado («Finalizado» cuando está cerrado).
- **Rol en el proyecto**: Scrum Master (puede cerrar y reabrir) o Product Builder (no puede).
- **Estado de cierre**: abierto o cerrado (con fecha); alimenta el aviso y la deshabilitación de acciones.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: El Scrum Master cierra un proyecto abierto en no más de 3 interacciones (abrir la zona, cerrar, confirmar).
- **SC-002**: El 100% de los intentos de cierre cancelados en la confirmación terminan sin ningún cambio en el proyecto.
- **SC-003**: El 100% de los mensajes 409 y de los avisos 403/404 de cerrar y reabrir se muestran con el texto del servidor o el aviso estándar, sin pantalla en blanco ni mensaje genérico.
- **SC-004**: En un proyecto cerrado, el 100% de las acciones de escritura de las pantallas alcanzadas aparecen deshabilitadas con una nota, y el 0% de ellas provoca un 409 por acción del usuario.
- **SC-005**: Tras reabrir, el 100% de esas acciones vuelve a estar habilitada sin recargar la página.
- **SC-006**: Las pruebas de permisos existentes de HU-21 pasan sin modificarse.
- **SC-007**: Cada escenario normal, alternativo, límite y de error de esta spec tiene una prueba automatizada que lo cubre.

## Assumptions

- El backend de la spec 044 está en main y su contrato no cambia.
- La rama sale de la de HU-21 (PR #321) porque amplía los permisos compartidos.
- «Reabrir» no pide confirmación, porque se puede volver a cerrar y no destruye datos.
- La fecha de cierre se muestra tal como la entrega el servidor (AAAA-MM-DD).
- Si falta tiempo, se corta por fase: lo no entregado pasa a otro PR, y esas pantallas siguen mostrando el 409 del servidor mientras tanto.
- El armazón del rediseño (spec 046) queda intacto; el aviso y la zona se montan desde archivos nuevos.
