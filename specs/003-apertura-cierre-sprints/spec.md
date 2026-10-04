# Feature Specification: Apertura y Cierre de Sprints

**Feature Branch**: `003-apertura-cierre-sprints`

**Created**: 2026-10-01

**Status**: En progreso — NO TERMINADA, NO MERGEAR. Rama publicada solo como
respaldo. Implementadas US1 (iniciar) y US2 (cerrar con arrastre); falta US3
(crear Sprint); HU-04 ya está en `main`. Ver `tasks.md`.

**Input**: User description: "HU-05: Apertura y Cierre de Sprints. Como Scrum
Master quiero crear, iniciar y cerrar un Sprint definiendo su objetivo
(Sprint Goal) y plazos para poder organizar y delimitar las iteraciones de
desarrollo del equipo. Depende de HU-04 (Project), ya existe un contrato
mínimo commiteado (Project/ProjectRepository) para TDD con fake en memoria."

## Clarifications

### Session 2026-10-04

- Q: ¿Qué responde la API al crear un Sprint para un proyecto que no existe? → A: 404 "proyecto no encontrado", igual que HU-04.
- Q: Al iniciar un Sprint inexistente con datos inválidos, ¿qué error se informa primero? → A: 404 primero; el cuerpo se valida (400) solo si el Sprint existe.
- Q: ¿Un Sprint puede tener la misma fecha de inicio y de fin? → A: No; el fin debe ser estrictamente posterior (posible, pero no práctico). HU-04 acepta fechas iguales para proyectos a propósito.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Iniciar un Sprint (Priority: P1)

Como Scrum Master, selecciono un Sprint pendiente de mi proyecto, le defino
un Sprint Goal y un rango de fechas, y lo inicio para que el equipo sepa que
la iteración está en marcha.

**Why this priority**: Sin poder iniciar un Sprint no existe el marco
iterativo que pide la consigna — es el requisito mínimo para que Scrum
funcione en la aplicación.

**Independent Test**: Se puede probar por completo creando un Sprint en
estado "Pendiente" y verificando que, al iniciarlo con un Goal y fechas
válidas, pasa a "Activo".

**Acceptance Scenarios**:

1. **Given** existe un Sprint con estado "Pendiente" en un proyecto sin
   ningún otro Sprint "Activo", **When** se define el Sprint Goal, la fecha
   de inicio y la fecha de fin, y se inicia el Sprint, **Then** su estado
   cambia a "Activo".
2. **Given** un proyecto ya tiene un Sprint en estado "Activo", **When** se
   intenta iniciar otro Sprint del mismo proyecto, **Then** el sistema
   rechaza la operación y el segundo Sprint permanece "Pendiente".

---

### User Story 2 - Cerrar un Sprint con arrastre de trabajo (Priority: P1)

Como Scrum Master, cierro el Sprint activo al terminar la iteración, y las
historias que no se llegaron a completar vuelven solas al Product Backlog
para poder replanificarlas en el próximo Sprint.

**Why this priority**: Es la otra mitad del ciclo iterativo — sin cierre no
hay forma de medir qué se completó ni de recuperar el trabajo pendiente.

**Independent Test**: Se puede probar por completo cerrando un Sprint
"Activo" que tiene historias completadas y no completadas, y verificando
que solo las no completadas quedan sueltas en el Product Backlog.

**Acceptance Scenarios**:

1. **Given** un Sprint está en estado "Activo" con historias completadas y
   no completadas asignadas, **When** el Scrum Master cierra el Sprint,
   **Then** el estado del Sprint cambia a "Finalizado" y las historias no
   completadas quedan sin Sprint asignado (vuelven al Product Backlog).
2. **Given** un Sprint está en estado "Pendiente" o ya "Finalizado",
   **When** se intenta cerrarlo, **Then** el sistema rechaza la operación.

---

### User Story 3 - Crear un Sprint (Priority: P2)

Como Scrum Master, creo un Sprint nuevo asociado a mi proyecto para
prepararlo antes de iniciarlo.

**Why this priority**: Es un prerrequisito técnico de las historias 1 y 2,
pero de menor valor por sí solo — nadie "crea y deja ahí" un Sprint sin
intención de iniciarlo.

**Independent Test**: Se puede probar por completo creando un Sprint para
un proyecto existente y verificando que queda en estado "Pendiente", sin
Goal ni fechas todavía.

**Acceptance Scenarios**:

1. **Given** existe un Project, **When** se crea un Sprint para ese
   proyecto sin más datos, **Then** el Sprint queda persistido en estado
   "Pendiente".
2. **Given** se intenta crear un Sprint para un `ProyectoID` que no existe,
   **When** se ejecuta la creación, **Then** el sistema rechaza la
   operación.
3. **Given** un proyecto ya tiene un Sprint en estado "Pendiente", **When**
   se intenta crear otro Sprint para ese proyecto, **Then** el sistema
   rechaza la operación y el proyecto sigue teniendo un solo Sprint
   "Pendiente".

### Edge Cases

- ¿Qué pasa si se intenta cerrar un Sprint que nunca se inició (sigue en
  "Pendiente")? → Se rechaza, solo se cierran Sprints "Activo".
- ¿Qué pasa si la fecha de fin ingresada es anterior o igual a la fecha de
  inicio al iniciar el Sprint? → Se rechaza con error de validación,
  incluido el Sprint de un día. Difiere a propósito de HU-04, que acepta
  fechas iguales para un proyecto: un Sprint de un día no es práctico.
- ¿Qué pasa si el Sprint que se cierra no tiene ninguna historia asignada?
  → Se cierra igual, no es un error; simplemente no hay nada que mover.
- ¿Qué pasa si se intenta iniciar o cerrar un Sprint que no existe
  (`SprintID` inexistente)? → Se rechaza con error de "no encontrado",
  aunque el cuerpo de la petición también sea inválido: la existencia se
  verifica antes que los datos (mismo criterio que HU-04).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir crear un Sprint asociado a un
  `Project` existente, quedando en estado inicial "Pendiente".
- **FR-002**: El sistema MUST rechazar la creación de un Sprint si el
  `Project` indicado no existe, con error de "proyecto no encontrado"
  (no de "datos inválidos"), mismo criterio que HU-04.
- **FR-003**: El sistema MUST permitir iniciar un Sprint en estado
  "Pendiente", definiendo en ese momento el Sprint Goal, la fecha de inicio
  y la fecha de fin, pasando su estado a "Activo".
- **FR-004**: El sistema MUST rechazar iniciar un Sprint si ya existe otro
  Sprint en estado "Activo" para el mismo `Project`.
- **FR-005**: El sistema MUST validar que la fecha de fin sea estrictamente posterior a
  la fecha de inicio al iniciar un Sprint.
- **FR-006**: El sistema MUST permitir cerrar un Sprint que esté en estado
  "Activo", cambiando su estado a "Finalizado".
- **FR-007**: El sistema MUST rechazar el cierre de un Sprint que no esté
  en estado "Activo".
- **FR-008**: Al cerrar un Sprint, el sistema MUST desvincular del Sprint a
  todas las Historias de Usuario asignadas que no estén completadas,
  dejándolas disponibles en el Product Backlog sin Sprint asignado.
- **FR-009**: Al cerrar un Sprint, las Historias de Usuario que sí estén
  completadas MUST permanecer vinculadas a ese Sprint.
- **FR-010**: El sistema MUST poder consultar los Sprints de un `Project`
  filtrando por estado, para poder validar la regla de "un solo Sprint
  Activo por proyecto" (FR-004).
- **FR-011**: El sistema MUST rechazar la creación de un nuevo Sprint
  "Pendiente" para un `Project` si ya existe otro Sprint "Pendiente" para
  ese mismo proyecto (máximo un Sprint "Pendiente" por proyecto a la vez,
  igual criterio que la regla de "un solo Activo").

### Key Entities

- **Sprint**: iteración de trabajo del equipo. Atributos: identificador,
  `ProyectoID` al que pertenece, Sprint Goal (texto), fecha de inicio, fecha
  de fin, estado (`Pendiente` / `Activo` / `Finalizado`).
- **Historia de Usuario** (`BacklogItem`, ya existente): se asocia a un
  Sprint mediante una referencia opcional (una historia puede no tener
  Sprint asignado — está en el Product Backlog — o pertenecer a uno).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un Scrum Master puede iniciar un Sprint (Goal + fechas) en un
  solo paso, sin pantallas intermedias.
- **SC-002**: El sistema previene el 100% de los intentos de tener dos
  Sprints activos simultáneos en el mismo proyecto.
- **SC-003**: Al cerrar un Sprint, el 100% de las historias no completadas
  quedan visibles de nuevo en el Product Backlog sin intervención manual
  adicional.
- **SC-004**: Las historias completadas de Sprints cerrados siguen siendo
  consultables asociadas a su Sprint de origen (sin pérdida de datos para
  HU-03/HU-12).

## Assumptions

- Las Historias de Usuario completadas en un Sprint cerrado permanecen
  vinculadas a ese Sprint (no se desvinculan) — es lo que necesitan HU-03
  (cálculo de métricas por Sprint) y HU-12 (consulta de Sprints anteriores)
  para funcionar.
- No hay cierre automático de Sprints por fecha vencida; cerrar un Sprint
  es siempre una acción manual del Scrum Master, consistente con el
  Escenario 2 original ("el Scrum Master presiona Cerrar Sprint").
- Esta historia usa el `ProjectRepository` real de HU-04 (ya en `main`,
  IDs `int64`) para verificar que el proyecto existe; en los tests se usa
  su fake en memoria (`ProjectRepositoryEnMemoria`).
- Un proyecto tiene como máximo un Sprint "Pendiente" a la vez (FR-011) —
  coincide con la práctica real de Scrum de planificar un Sprint por vez,
  no una cola de varios por adelantado.
