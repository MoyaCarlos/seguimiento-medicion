# Feature Specification: Creación de Proyecto y Asignación de Equipo (HU-04)

**Feature Branch**: `HU-04-crear-proyecto-equipo`

**Created**: 2026-10-01

**Status**: Draft

**Input**: User description: "HU-04: Creación de Proyecto y Asignación de Equipo. Como Scrum Master quiero crear un proyecto y agregar a los integrantes con sus respectivos roles para poder tener un espacio de trabajo organizado y estructurar las responsabilidades del equipo. Prioridad M (Must have). Estado Nueva. Valor de Negocio 21 (Fibonacci). Estimación 8 Story Points."

## Clarifications

### Session 2026-10-01

- Q: ¿Cómo se identifica y vincula a un integrante del equipo? → A: El Scrum Master ingresa el nombre y el sistema crea (o reutiliza) el usuario con ese nombre; no se requiere registro previo.
- Q: ¿Cuántos Scrum Master puede tener un proyecto y es obligatorio que tenga alguno? → A: Al menos uno (el creador); se permiten varios.
- Q: ¿Esta historia aplica los permisos por rol o solo registra la asignación? → A: Solo registra la relación Proyecto–Usuario–Rol; la verificación/aplicación efectiva de permisos se difiere a una historia posterior.
- Q: ¿Cómo determina el sistema que un nombre ingresado corresponde a un integrante ya existente y debe reutilizarlo? → A: Comparación ignorando mayúsculas/minúsculas y espacios externos; "Ana", "ana" y " Ana " son el mismo integrante.
- Q: ¿Un proyecto debe tener estados de ciclo de vida o basta con registrarlo sin estado? → A: No tiene estados; se registra sin estado y queda operativo desde su creación.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Crear el entorno del proyecto (Priority: P1)

Como Scrum Master, desde la pantalla principal del sistema ingreso el nombre del proyecto ("Software Metrics & Estimation"), una descripción y presiono "Crear". El sistema genera el proyecto, lo deja disponible de forma persistente y confirma su creación para que el cliente pueda abrir su panel principal, listo para configurar el equipo y operar sobre él. La navegación/redirección efectiva del cliente es responsabilidad del frontend (ver Assumptions).

**Why this priority**: Sin un proyecto no existe contenedor para el backlog, los sprints, los defectos ni las métricas. Es el punto de entrada de todo el flujo del producto y habilita el resto de las historias (MVP).

**Independent Test**: Puede probarse de forma independiente creando un proyecto con datos válidos y verificando que la respuesta confirma el proyecto persistido (con su identificador) y que sigue disponible al consultarlo luego de recargar.

**Acceptance Scenarios**:

1. **Given** que me encuentro en la pantalla principal del sistema, **When** ingreso el nombre del proyecto y una descripción y presiono "Crear", **Then** el sistema genera el proyecto de forma persistente y confirma su creación devolviendo el proyecto (con su identificador), habilitando su panel principal (la navegación la realiza el frontend).
2. **Given** que intento crear un proyecto, **When** dejo el nombre vacío (o compuesto solo por espacios) y presiono "Crear", **Then** el sistema muestra una advertencia y no genera el proyecto.

---

### User Story 2 - Asignar integrantes y roles al proyecto (Priority: P2)

Como Scrum Master, desde la vista de configuración del proyecto recién creado, ingreso el nombre de un integrante y selecciono su rol (Scrum Master o Product Builder). El sistema vincula al integrante con el proyecto y habilita los permisos correspondientes a su rol, dejando estructurada la responsabilidad del equipo.

**Why this priority**: Sin integrantes ni roles asignados el proyecto existe pero no es operable por el equipo. Aporta el valor de organización y control de responsabilidades que motiva la historia, pero depende de que el proyecto exista (US1).

**Independent Test**: Puede probarse de forma independiente sobre un proyecto existente agregando un integrante con un rol válido y verificando que la vinculación queda persistida con su rol asociado.

**Acceptance Scenarios**:

1. **Given** que estoy en la vista de configuración de un proyecto, **When** ingreso el nombre de un integrante, selecciono su rol (Scrum Master o Product Builder) y confirmo, **Then** el sistema vincula al integrante con el proyecto y le habilita los permisos de su rol.
2. **Given** que estoy en la vista de configuración de un proyecto, **When** intento agregar un integrante sin nombre o sin rol válido, **Then** el sistema muestra una advertencia y no realiza la vinculación.

---

### User Story 3 - Editar un proyecto existente (Priority: P3)

Como Scrum Master, desde la vista de edición de un proyecto ya creado puedo corregir su nombre, su descripción y sus fechas para arreglar errores de carga, manteniendo el proyecto en el mismo estado operativo y sin conservar versiones ni historial de cambios. La edición reutiliza las mismas validaciones aplicadas al crear el proyecto.

**Why this priority**: Permite corregir datos erróneos de un proyecto ya cargado sin recrearlo. Es un requisito del enunciado ("Crear y modificar proyectos"), pero el MVP ya es funcional con la creación y la asignación de equipo; la edición es una corrección posterior de datos.

**Independent Test**: Puede probarse de forma independiente sobre un proyecto existente modificando su nombre, descripción y fechas con datos válidos, verificando que los cambios quedan persistidos y que un valor inválido es rechazado.

**Acceptance Scenarios**:

1. **Given** que estoy en la vista de edición de un proyecto existente, **When** modifico su nombre, descripción o fechas con datos válidos y guardo, **Then** el sistema persiste los cambios del proyecto.
2. **Given** que estoy editando un proyecto, **When** dejo el nombre vacío (o compuesto solo por espacios) y guardo, **Then** el sistema muestra una advertencia y no guarda los cambios.
3. **Given** que estoy editando un proyecto, **When** intento guardar una fecha de fin anterior a la fecha de inicio, **Then** el sistema muestra una advertencia y no guarda los cambios.
4. **Given** que estoy editando un proyecto, **When** decido no aplicar cambios y cancelo, **Then** el proyecto permanece sin modificaciones.

---

### Edge Cases

- ¿Qué ocurre si el nombre del proyecto contiene solo espacios en blanco? Se considera inválido: se muestra la advertencia y no se crea el proyecto.
- ¿Qué ocurre si el nombre del proyecto supera los 100 caracteres? Se muestra una advertencia y se solicita acortarlo; no se crea el proyecto.
- ¿Qué ocurre si la descripción del proyecto está vacía? El proyecto se crea igualmente: la descripción es opcional.
- ¿Qué ocurre si dos proyectos tienen el mismo nombre? Se permiten nombres repetidos (no hay unicidad de nombre de proyecto).
- ¿Qué ocurre si intento asignar un rol distinto de Scrum Master o Product Builder? El sistema rechaza la selección con un mensaje de validación.
- ¿Qué ocurre si intento agregar al mismo integrante dos veces al mismo proyecto? El sistema lo rechaza indicando que el integrante ya pertenece al proyecto.
- ¿Qué ocurre si intento agregar integrantes a un proyecto que no existe? La acción no está disponible: la vista de configuración solo opera sobre un proyecto creado.
- ¿Qué ocurre si el nombre del integrante supera los 200 caracteres? Se muestra una advertencia y no se realiza la vinculación.
- ¿Qué ocurre si dos Scrum Masters crean proyectos al mismo tiempo? Ambos proyectos quedan registrados de forma independiente, sin sobrescribirse.
- ¿Qué ocurre al editar si dejo la descripción vacía? Se permite: la descripción sigue siendo opcional al modificar.
- ¿Qué ocurre al editar si dejo las fechas sin completar? Se permite: las fechas son opcionales.
- ¿Qué ocurre al editar si la fecha de fin es anterior a la de inicio? Se muestra una advertencia y no se guardan los cambios.
- ¿Qué ocurre si edito un proyecto y le pongo un nombre ya usado por otro proyecto? Se permite: no hay unicidad de nombre de proyecto.
- ¿Qué ocurre si intento editar un proyecto que no existe? La operación no está disponible; solo se puede editar un proyecto existente.
- ¿Qué ocurre si dos personas editan el mismo proyecto al mismo tiempo? No hay versionado ni historial: prevalece el último guardado que resulte válido.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir a un Scrum Master crear un proyecto indicando un nombre obligatorio (≤ 100 caracteres) y una descripción opcional (≤ 2000 caracteres).
- **FR-002**: El sistema MUST asignar un identificador único e interno a cada proyecto creado.
- **FR-003**: El sistema MUST persistir todo proyecto creado de forma que permanezca disponible luego de cerrar y volver a abrir la aplicación.
- **FR-004**: Después de crear un proyecto con datos válidos, el sistema MUST confirmar la creación devolviendo el proyecto (con su identificador y datos), de modo que el cliente pueda abrir su panel principal. La navegación/redirección efectiva es responsabilidad del frontend y queda fuera del alcance de esta historia (ver Assumptions).
- **FR-005**: El sistema MUST exigir el nombre del proyecto como campo obligatorio, no vacío (ni compuesto solo por espacios) y de hasta 100 caracteres como máximo.
- **FR-006**: Cuando falten datos obligatorios o se excedan las longitudes máximas, el sistema MUST mostrar una advertencia que identifique el problema y MUST NOT crear el proyecto.
- **FR-007**: El sistema MUST admitir proyectos con nombres repetidos.
- **FR-008**: El sistema MUST vincular automáticamente al Scrum Master creador como integrante con rol Scrum Master del proyecto creado.
- **FR-009**: El sistema MUST permitir, desde la vista de configuración de un proyecto existente, vincular integrantes indicando el nombre del integrante y su rol.
- **FR-010**: El sistema MUST aceptar, para esta historia, únicamente los roles Scrum Master y Product Builder al vincular integrantes.
- **FR-011**: El sistema MUST asociar a cada integrante vinculado un rol dentro del proyecto, de modo que la vinculación quede persistida y disponible posteriormente.
- **FR-012**: El sistema MUST dejar habilitados los roles correspondientes al vincular a cada integrante (relación Proyecto–Usuario–Rol), registrando la asociación que habilita los permisos de su rol. La verificación y aplicación efectiva de permisos en las operaciones queda fuera del alcance de esta historia (ver Assumptions).
- **FR-013**: El sistema MUST exigir el nombre del integrante como campo obligatorio y no vacío (ni compuesto solo por espacios), de hasta 200 caracteres como máximo.
- **FR-014**: Cuando falten datos obligatorios del integrante o el rol no sea válido, el sistema MUST mostrar una advertencia y MUST NOT realizar la vinculación.
- **FR-015**: El sistema MUST rechazar la vinculación de un integrante que ya pertenece al mismo proyecto.
- **FR-016**: El sistema MUST identificar al integrante comparando el nombre sin distinguir mayúsculas/minúsculas ni espacios externos, creándolo si no existe o reutilizándolo si ya existe, y MUST permitir que un mismo integrante pertenezca a distintos proyectos.
- **FR-017**: El sistema MUST garantizar que todo proyecto tenga al menos un integrante con rol Scrum Master (el creador) y MUST permitir asignar más de un Scrum Master al mismo proyecto.
- **FR-018**: El sistema MUST permitir a un Scrum Master editar un proyecto existente modificando su nombre, su descripción y sus fechas.
- **FR-019**: La edición MUST reutilizar las mismas validaciones aplicadas al crear el proyecto: nombre obligatorio, no vacío (ni compuesto solo por espacios) y de hasta 100 caracteres; descripción de hasta 2000 caracteres; fechas opcionales.
- **FR-020**: El sistema MUST validar que la fecha de fin no sea anterior a la fecha de inicio y MUST rechazar un rango de fechas inválido.
- **FR-021**: Cuando la edición tenga datos inválidos, el sistema MUST mostrar una advertencia que identifique el problema y MUST NOT persistir los cambios, conservando el proyecto con sus valores previos.
- **FR-022**: El sistema MUST persistir los cambios válidos de un proyecto editado de forma que permanezcan disponibles luego de recargar la aplicación.
- **FR-023**: El sistema MUST aplicar la edición sin conservar versiones ni historial de cambios del proyecto (edición simple).

### Key Entities *(include if feature involves data)*

- **Proyecto**: espacio de trabajo que agrupa backlog, sprints, defectos y métricas. Atributos relevantes: identificador único, nombre (obligatorio, ≤ 100 caracteres) y descripción (opcional, ≤ 2000 caracteres), más fechas de inicio y de fin (opcionales). No tiene estados de ciclo de vida: se crea operativo y se registra con al menos un integrante con rol Scrum Master (el creador). El nombre, la descripción y las fechas son editables; la edición no conserva versiones ni historial.
- **Usuario / Integrante**: persona que participa en un proyecto. Atributos relevantes: nombre (obligatorio, ≤ 200 caracteres; se compara ignorando mayúsculas/minúsculas y espacios externos para decidir su reutilización). Un mismo integrante puede pertenecer a varios proyectos.
- **Rol**: categoría de responsabilidad dentro de un proyecto. Para esta historia los valores válidos son Scrum Master (código `scrum_master`) y Product Builder (código `product_builder`); el código es el valor canónico que usa el sistema, y la etiqueta es el texto visible en la interfaz.
- **Asignación (Proyecto–Usuario–Rol)**: relación que vincula un integrante con un proyecto y le asigna exactamente un rol, habilitando sus permisos correspondientes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un Scrum Master puede crear un proyecto completando el nombre (y opcionalmente la descripción) en menos de 1 minuto, y la respuesta confirma el proyecto creado con su identificador.
- **SC-002**: El 100% de los proyectos creados con datos válidos quedan persistidos y siguen disponibles después de recargar la aplicación.
- **SC-003**: El 100% de los intentos de creación con el nombre vacío o inválido son rechazados con una advertencia visible y no generan un proyecto.
- **SC-004**: El 100% de las vinculaciones de integrantes con datos válidos quedan persistidas con su rol asociado.
- **SC-005**: El 100% de los intentos de vinculación sin nombre, con rol inválido o con un integrante ya perteneciente al proyecto son rechazados sin modificar el equipo.
- **SC-006**: Cero proyectos operables sin un integrante con rol Scrum Master asignado.
- **SC-007**: El 100% de los integrantes vinculados a un proyecto tienen exactamente un rol asignado dentro de ese proyecto.
- **SC-008**: Un Scrum Master puede editar un proyecto existente (nombre, descripción y fechas) y confirmar los cambios en menos de 1 minuto.
- **SC-009**: El 100% de las ediciones con datos válidos quedan persistidas y siguen disponibles después de recargar la aplicación.
- **SC-010**: El 100% de los intentos de edición con el nombre vacío o con un rango de fechas inválido son rechazados con una advertencia visible, conservando el proyecto con sus valores previos.

## Assumptions

- Existe autenticación previa y el usuario que opera la historia ya está identificado como Scrum Master; el inicio de sesión y el registro de usuarios no forman parte de esta historia.
- La creación del proyecto y la asignación de integrantes se realizan en dos momentos: primero el proyecto (pantalla principal) y luego el equipo (vista de configuración del proyecto).
- La descripción del proyecto es opcional; solo el nombre es obligatorio.
- Se adoptan longitudes máximas: 100 caracteres para el nombre del proyecto, 2000 para su descripción y 200 caracteres para el nombre del integrante (valores razonables de la industria).
- Los roles se almacenan con códigos canónicos (`scrum_master`, `product_builder`) y se muestran con las etiquetas "Scrum Master" y "Product Builder" en la interfaz; el mapeo código↔etiqueta es responsabilidad de la capa de presentación.
- No se exige unicidad del nombre de proyecto; se permiten nombres repetidos.
- El proyecto no tiene estados de ciclo de vida: se registra sin estado y queda operativo desde su creación.
- Los únicos roles gestionados en esta historia son Scrum Master y Product Builder, según el criterio de aceptación. Otros roles del equipo (Product Architect, Agile Enabler) quedan fuera de alcance.
- El integrante creador del proyecto queda asociado automáticamente como Scrum Master del mismo.
- Un integrante se identifica por su nombre: al vincularlo, el sistema crea el usuario si no existe o reutiliza el existente; no se requiere registro previo. La comparación ignora mayúsculas/minúsculas y espacios externos.
- Todo proyecto debe tener al menos un Scrum Master (el creador); se permiten varios.
- Esta historia registra la relación Proyecto–Usuario–Rol y la habilitación de roles; la verificación y aplicación efectiva de permisos por rol en las operaciones se difiere a una historia posterior.
- El sistema adopta dos fechas del proyecto: fecha de inicio y fecha de fin, ambas opcionales; cuando ambas se informan, la fecha de fin no puede ser anterior a la fecha de inicio.
- La edición del proyecto es simple: se modifican nombre, descripción y fechas, sin versionado ni historial de cambios; prevalece el último guardado que resulte válido.
- La edición de la composición del equipo (alta, baja o cambio de rol de integrantes) y el borrado de proyectos quedan fuera del alcance de esta historia.
- El listado/navegación entre proyectos, la redirección al panel principal (FR-004), la vista de configuración y la vista detallada del equipo quedan fuera del alcance de esta historia; esta historia garantiza la persistencia del proyecto, de las asignaciones de equipo y devuelve el proyecto para que el frontend navegue.
- La existencia del proyecto referenciado por historias de usuario (HU-01, FR-012) queda habilitada por esta historia.
