# Feature Specification: Creación de Historias de Usuario (HU-01)

**Feature Branch**: `HU-01-crear-historias-usuario`

**Created**: 2026-09-30

**Status**: Draft

**Input**: User description: "HU-01: Creación de Historias de Usuario. Como Product Builder quiero crear historias de usuario con prioridad, estado y estimación para poder alimentar y organizar el Product Backlog del proyecto. Prioridad M (Must have - MVP). Estado Nuevo. Valor de Negocio 21 (Fibonacci). Estimación 5 Story Points."

## Clarifications

### Session 2026-09-30

- Q: Al crear una historia de usuario, ¿la estimación en Story Points se ingresa directamente en el formulario de creación o solo se define después mediante Planning Poker (HU-02)? → A: Solo vía Planning Poker; la estimación queda fuera del alcance de esta historia.
- Q: ¿El formulario de creación debe capturar el Valor de Negocio de la historia, o ese dato se completa en otro momento? → A: Capturar el Valor de Negocio de forma opcional en la creación.
- Q: ¿Qué roles pueden crear historias de usuario en el Product Backlog? → A: Product Builder y Scrum Master.
- Q: ¿Qué límites máximos de longitud deben aplicarse al título y a la descripción? → A: Título 200 caracteres; descripción 2000 caracteres.
- Q: ¿Qué valores exactos de la escala tipo Fibonacci son válidos para el Valor de Negocio? → A: 1, 2, 3, 5, 8, 13, 21.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Crear una historia de usuario en el Product Backlog (Priority: P1)

Como integrante del equipo (Product Builder o Scrum Master), desde el panel del Product Backlog del proyecto, completo los datos obligatorios de una historia de usuario (título, descripción y prioridad), presiono "Guardar" y la historia queda registrada en el backlog con estado inicial "Nueva", disponible para su posterior planificación y estimación.

**Why this priority**: Sin la creación de historias de usuario no existe Product Backlog y, por lo tanto, no es posible planificar Sprints, estimar esfuerzo ni calcular ninguna métrica. Es el punto de entrada de todo el flujo de trabajo del producto (MVP).

**Independent Test**: Puede probarse de forma independiente creando una historia con datos válidos y verificando que queda persistida con estado "Nueva" y que sigue disponible luego de recargar la aplicación.

**Acceptance Scenarios**:

1. **Given** que me encuentro en el panel del Product Backlog de un proyecto existente, **When** ingreso los datos obligatorios (título, descripción, prioridad) y presiono "Guardar", **Then** la historia queda registrada de forma persistente con estado "Nueva".
2. **Given** que intento crear una historia, **When** dejo el campo "Título" en blanco y presiono "Guardar", **Then** el sistema muestra una advertencia y no registra la historia.

---

### User Story 2 - Registrar el Valor de Negocio de la historia (Priority: P2)

Como integrante del equipo, al crear una historia de usuario puedo indicar su valor de negocio (escala tipo Fibonacci) para priorizar el backlog en función del impacto de cada historia.

**Why this priority**: Aporta el valor de negocio usado para priorizar el backlog y calcular métricas de planificación, pero el MVP del backlog puede existir sin este dato (puede completarse luego).

**Independent Test**: Puede probarse creando una historia indicando su valor de negocio y verificando que queda asociado a la historia persistida.

**Acceptance Scenarios**:

1. **Given** que estoy completando el formulario de una nueva historia, **When** indico un valor de negocio de la escala permitida y guardo, **Then** el valor queda asociado a la historia.
2. **Given** que estoy completando el formulario de una nueva historia, **When** dejo el valor de negocio sin completar y guardo, **Then** la historia se registra igualmente con ese campo vacío (no es obligatorio).

---

### Edge Cases

- ¿Qué ocurre si el título contiene solo espacios en blanco? Se considera inválido: se muestra la advertencia y no se registra la historia.
- ¿Qué ocurre si falta la descripción o la prioridad? Se muestra una advertencia indicando los campos faltantes y no se registra la historia.
- ¿Qué ocurre si el título supera los 200 caracteres o la descripción supera los 2000? Se muestra una advertencia y se solicita acortarlos; no se registra la historia.
- ¿Qué ocurre si dos historias tienen el mismo título? Se permiten títulos repetidos (no hay unicidad de título).
- ¿Qué ocurre si intento guardar con un valor de negocio fuera de la escala Fibonacci permitida? El sistema rechaza el valor con un mensaje de validación.
- ¿Qué ocurre si falta el identificador de proyecto (`proyecto_id`), es cero o negativo? El sistema rechaza la creación con una advertencia de validación y no registra la historia (la existencia del proyecto en la base se valida en HU-04).
- ¿Qué ocurre si dos integrantes guardan una historia al mismo tiempo? Ambas quedan registradas sin sobrescribirse, preservando su orden de creación.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir a un Product Builder o a un Scrum Master crear una historia de usuario indicando como mínimo su título, descripción y prioridad.
- **FR-002**: El sistema MUST asignar automáticamente el estado inicial "Nueva" a toda historia recién creada, sin permitir elegir otro estado en la creación.
- **FR-003**: El sistema MUST preservar el orden de creación de las historias mediante un identificador incremental, de modo que cualquier listado posterior del Product Backlog pueda presentarlas al final (la interfaz de listado queda fuera del alcance de esta historia; ver Assumptions).
- **FR-004**: El sistema MUST exigir el título como campo obligatorio y no vacío (ni compuesto solo por espacios), con un máximo de 200 caracteres.
- **FR-005**: El sistema MUST exigir descripción y prioridad como campos obligatorios, con un máximo de 2000 caracteres para la descripción.
- **FR-006**: Cuando falten campos obligatorios o se excedan las longitudes máximas, el sistema MUST mostrar una advertencia que identifique el problema y MUST NOT registrar la historia.
- **FR-007**: El sistema MUST persistir toda historia creada de forma que permanezca disponible luego de cerrar y volver a abrir la aplicación.
- **FR-008**: El sistema MUST permitir registrar opcionalmente el valor de negocio de la historia en una escala tipo Fibonacci.
- **FR-009**: El sistema MUST validar que el valor de negocio, cuando se informe, sea uno de los valores permitidos de la escala Fibonacci: 1, 2, 3, 5, 8, 13 o 21.
- **FR-010**: La prioridad MUST registrarse según la escala MoSCoW (Must / Should / Could / Won't have).
- **FR-011**: El sistema MUST admitir títulos repetidos entre historias del mismo proyecto.
- **FR-012**: El sistema MUST asociar cada historia a un proyecto mediante su identificador (`proyecto_id`, obligatorio y entero positivo `> 0`), rechazando la creación si falta o no es válido. La validación de que el proyecto exista en la base se difiere a HU-04, cuando se introduzca la entidad Proyecto (ver Assumptions).
- **FR-013**: El sistema MUST registrar toda historia creada con su estimación en Story Points vacía; la estimación queda fuera del alcance de esta historia y se define en la votación de Planning Poker (HU-02).
- **FR-014** (diferido a HU-04): Solo los roles Product Builder y Scrum Master podrán crear historias; otros roles deberán recibir un rechazo al intentarlo. No se aplica en HU-01 por no existir aún el mecanismo de identidad y roles (ver Assumptions).

### Key Entities *(include if feature involves data)*

- **Historia de Usuario**: unidad de trabajo del Product Backlog. Atributos relevantes: identificador (incremental, preserva el orden de creación), título (≤ 200 caracteres), descripción (≤ 2000 caracteres), prioridad (MoSCoW), estado (inicial "Nueva"), valor de negocio (Fibonacci ∈ {1,2,3,5,8,13,21}, opcional), estimación en Story Points (vacía al crear; se completa en HU-02) y proyecto al que pertenece (`proyecto_id`, entero positivo obligatorio). Los criterios de aceptación por historia quedan fuera del alcance de esta historia.
- **Product Backlog**: conjunto de las historias de usuario de un proyecto, ordenado por creación; la vista de listado se difiere a un incremento de frontend.
- **Proyecto**: contenedor al que pertenecen las historias y el backlog; debe existir previamente (depende de HU-04).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un Product Builder o Scrum Master puede crear una historia completando únicamente los tres datos obligatorios en menos de 1 minuto.
- **SC-002**: El 100% de las creaciones con datos válidos quedan persistidas con estado "Nueva" y con un identificador incremental posterior al de las historias existentes (orden de creación preservado).
- **SC-003**: El 100% de las historias creadas correctamente siguen disponibles después de recargar la aplicación (persistencia verificable).
- **SC-004**: El 100% de los intentos de creación con el título en blanco son rechazados con una advertencia visible y no generan un registro.
- **SC-005**: Cero historias registradas sin los campos obligatorios (título, descripción y prioridad) o excediendo las longitudes máximas.

## Assumptions

- Existe autenticación previa y el usuario que crea la historia ya está identificado con rol Product Builder o Scrum Master; el inicio de sesión no es parte de esta historia.
- Ya existe al menos un proyecto creado (HU-04); esta historia opera dentro del panel del Product Backlog de un proyecto existente.
- La prioridad se registra con la escala MoSCoW y luego se mapea a Alta/Media/Baja en el tablero de gestión (convención del Product Backlog del proyecto).
- La estimación en Story Points no se ingresa al crear la historia; se define mediante Planning Poker (HU-02) y puede permanecer vacía.
- Valor de Negocio y Estimación son independientes entre sí (estilo WSJF) aunque ambos usen una escala tipo Fibonacci; no deben coincidir necesariamente.
- La escala tipo Fibonacci adoptada para el Valor de Negocio es 1, 2, 3, 5, 8, 13, 21.
- No se exige unicidad de título; se permiten historias con títulos repetidos.
- La edición y el borrado de historias existentes están fuera del alcance de esta historia.
- El listado/vista del Product Backlog, su orden visual y la actualización en vivo sin recargar corresponden a un incremento de frontend posterior; esta historia garantiza la persistencia y el orden de creación mediante identificador incremental.
- La autorización por rol (FR-014) se difiere a HU-04 (identidad y roles); HU-01 no aplica control de acceso por rol.
- La validación de que el proyecto referenciado exista (FR-012) se difiere a HU-04; esta historia solo persiste `proyecto_id`.
- La captura de criterios de aceptación por historia no forma parte de esta historia.
