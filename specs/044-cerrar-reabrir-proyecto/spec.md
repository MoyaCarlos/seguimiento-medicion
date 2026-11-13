# Feature Specification: Cerrar y reabrir un proyecto (solo el Scrum Master)

**Feature Branch**: `feature/300-cerrar-reabrir-proyecto`

**Issue**: #300

**Created**: 2026-11-11

**Status**: Draft

**Dependencias**: HU-21 (permisos por rol, spec `042-hu-21-permisos-por-rol`, PR #301 ya
mergeado) y HU-13 (estado del proyecto, spec `004-hu-13-estado-proyecto`). Cierra la
decisión **D16** de `MODELO-DE-DATOS.md`.

**Alcance de este trabajo**: solo backend. El frontend (zona «Cerrar proyecto» en Editar
proyecto, aviso de solo lectura y acciones deshabilitadas) lo hace otra integrante, en otro
PR.

**Input**: User description: "Cerrar y reabrir un proyecto (solo el Scrum Master), issue
#300. El Scrum Master del proyecto puede cerrar un proyecto abierto: queda «Finalizado» con
la fecha de cierre de hoy y de solo lectura. Puede reabrirlo: la fecha de cierre vuelve a
vacía y el estado se recalcula con las reglas de HU-13, que no cambian. Cerrar un proyecto
con un Sprint Activo se rechaza. Toda escritura sobre un proyecto cerrado responde 409 «el
proyecto está cerrado» y no cambia datos; el chequeo vive en un solo lugar, el control de
acceso de HU-21. Fuera de alcance: eliminar un proyecto (#232)."

## Clarifications

### Session 2026-11-11

Decisiones abiertas del issue #300, confirmadas por la usuaria (backend) el 2026-11-11:

- Q: ¿Cerrar un proyecto que tiene un Sprint Activo se rechaza o se permite? → A: Se
  **rechaza** con 409 «Cerrá el Sprint activo antes de cerrar el proyecto».
- Q: ¿Cómo se guarda el cierre? → A: Columna anulable `proyectos.fecha_cierre`, fecha de
  calendario `YYYY-MM-DD` con el `CHECK (date(fecha_cierre) IS fecha_cierre)` de D12, en la
  migración 008. Un solo PR con `MODELO-DE-DATOS.md` y la migración. Cierra D16.
- Q: D16 propone «sin cierre manual, el proyecto no pasa a Finalizado»; el issue propone que
  el cierre tenga prioridad y las reglas de HU-13 sigan iguales. ¿Cuál vale? → A: **La del
  issue.** Cerrado ⇒ «Finalizado» siempre; abierto ⇒ las reglas de HU-13 sin cambios.

## Objetivo

Que el Scrum Master pueda dar un proyecto por terminado: al cerrarlo queda «Finalizado» y
nadie puede modificarlo, pero se sigue consultando (backlog, Sprints, métricas, historial,
reportes). Si hizo falta cerrarlo antes de tiempo o hay que corregir algo, el Scrum Master
lo reabre y el proyecto vuelve a funcionar como antes.

**Descripción**: Como Scrum Master quiero cerrar un proyecto terminado, y poder reabrirlo,
para que sus datos queden fijos y el equipo no los modifique por error.

## Entradas

- La sesión de la persona (HU-20/HU-21).
- El proyecto a cerrar o reabrir (su identificador en la ruta). No hay cuerpo.

## Salidas esperadas

- **Cerrar**: el proyecto con su estado «Finalizado» y su fecha de cierre (la de hoy).
- **Reabrir**: el proyecto con la fecha de cierre vacía y el estado recalculado con las
  reglas de HU-13.
- **Consultar un proyecto** (detalle y «Mis proyectos», y la respuesta de crearlo y
  editarlo): incluye la fecha de cierre, vacía si el proyecto está abierto, para que la
  interfaz sepa si está en solo lectura. La consulta del estado mantiene su forma y da
  «Finalizado».
- **Escritura sobre un proyecto cerrado**: 409 «El proyecto está cerrado», sin cambios.

## Fuera de alcance

- Eliminar un proyecto (#232).
- El frontend (otro PR, a cargo de otra integrante).
- Cerrar automáticamente un proyecto (por fecha o por Sprints): el cierre es siempre manual.
- Guardar quién cerró el proyecto o un historial de cierres y reaperturas.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Cerrar un proyecto (Priority: P1)

El Scrum Master cierra un proyecto abierto que no tiene un Sprint Activo. El proyecto pasa
a «Finalizado» con la fecha de cierre de hoy.

**Why this priority**: es el corazón de la historia; sin cierre no hay nada que reabrir ni
que proteger.

**Independent Test**: con un proyecto abierto sin Sprint Activo, el Scrum Master lo cierra
y al consultarlo el estado es «Finalizado» y la fecha de cierre es la de hoy.

**Acceptance Scenarios**:

1. **Given** que soy el Scrum Master de un proyecto abierto sin Sprint Activo, **When** lo
   cierro, **Then** el proyecto queda «Finalizado» con la fecha de cierre de hoy.
2. **Given** un proyecto abierto «Planificado» o «En curso» por sus fechas, **When** el
   Scrum Master lo cierra, **Then** su estado pasa a ser «Finalizado», sin importar las
   fechas ni los Sprints.
3. **Given** que soy el Scrum Master de un proyecto con un Sprint Activo, **When** intento
   cerrarlo, **Then** el sistema responde 409 «Cerrá el Sprint activo antes de cerrar el
   proyecto» y el proyecto sigue abierto.
4. **Given** un proyecto ya cerrado, **When** el Scrum Master intenta cerrarlo de nuevo,
   **Then** el sistema responde 409 «El proyecto ya está cerrado» y la fecha de cierre no
   cambia.

---

### User Story 2 - Proyecto cerrado en solo lectura (Priority: P1)

Mientras el proyecto está cerrado, ninguna escritura sobre él se realiza: historias,
Sprints, horas, votación, defectos, equipo y edición del proyecto responden 409. Las
lecturas siguen funcionando.

**Why this priority**: es lo que le da sentido al cierre; sin esto, «cerrado» sería solo
una etiqueta.

**Independent Test**: con un proyecto cerrado, cada ruta de escritura sobre él responde 409
«El proyecto está cerrado» y los datos no cambian; las lecturas responden como antes.

**Acceptance Scenarios**:

1. **Given** un proyecto cerrado, **When** un integrante intenta crear, editar, asignar a un
   Sprint o cambiar de estado una historia, registrar horas, votar, o registrar o resolver
   un defecto, **Then** el sistema responde 409 «El proyecto está cerrado» y no cambia nada.
2. **Given** un proyecto cerrado, **When** el Scrum Master intenta crear, iniciar o cerrar un
   Sprint, iniciar o revelar una votación, abrir una nueva ronda, registrar la estimación,
   editar el proyecto o agregar un integrante, **Then** el sistema responde 409 «El proyecto
   está cerrado» y no cambia nada.
3. **Given** un proyecto cerrado, **When** cualquier persona con sesión consulta su
   backlog, Sprints, horas, votaciones, defectos, métricas, avance, historial o reportes,
   **Then** obtiene la misma respuesta que con el proyecto abierto (con el estado
   «Finalizado»).

---

### User Story 3 - Reabrir un proyecto (Priority: P1)

El Scrum Master reabre un proyecto cerrado: la fecha de cierre queda vacía, el estado se
vuelve a calcular con las reglas de HU-13 y se puede volver a escribir.

**Why this priority**: sin reapertura, un cierre por error deja el proyecto bloqueado para
siempre.

**Independent Test**: con un proyecto cerrado, el Scrum Master lo reabre; al consultarlo la
fecha de cierre está vacía, el estado es el que dan las reglas de HU-13 y una escritura
vuelve a funcionar.

**Acceptance Scenarios**:

1. **Given** que soy el Scrum Master de un proyecto cerrado, **When** lo reabro, **Then** la
   fecha de cierre queda vacía y el estado vuelve a calcularse según los Sprints y las
   fechas (HU-13).
2. **Given** un proyecto reabierto, **When** un integrante crea una historia, **Then** la
   historia se crea como en cualquier proyecto abierto.
3. **Given** un proyecto abierto, **When** el Scrum Master intenta reabrirlo, **Then** el
   sistema responde 409 «El proyecto no está cerrado» y no cambia nada.

---

### User Story 4 - Solo el Scrum Master cierra y reabre (Priority: P1)

**Why this priority**: es una acción de gestión (HU-21); un Product Builder no puede
congelar ni descongelar el proyecto.

**Independent Test**: con la sesión de un Product Builder y la de una persona que no
integra el proyecto, cerrar y reabrir responden el 403 de HU-21 y el proyecto no cambia.

**Acceptance Scenarios**:

1. **Given** que soy Product Builder del proyecto, **When** intento cerrarlo o reabrirlo,
   **Then** el sistema responde 403 «Solo el Scrum Master puede realizar esta acción» y el
   proyecto no cambia.
2. **Given** que no integro el proyecto, **When** intento cerrarlo o reabrirlo, **Then** el
   sistema responde 403 «No sos integrante de este proyecto» y el proyecto no cambia.
3. **Given** que no inicié sesión, **When** intento cerrarlo o reabrirlo, **Then** el
   sistema responde 401.
4. **Given** un proyecto que no existe, **When** el Scrum Master de otro proyecto intenta
   cerrarlo o reabrirlo, **Then** el sistema responde 404.

### Edge Cases

- **Reabrir un proyecto con la fecha de fin vencida o con todos sus Sprints
  finalizados**: al reabrirlo, HU-13 lo sigue mostrando «Finalizado», pero está abierto
  (fecha de cierre vacía) y acepta escrituras. «Finalizado» es el estado derivado; lo que
  bloquea las escrituras es solo el cierre.
- **Cerrar un proyecto sin Sprints ni fechas**: se permite; queda «Finalizado».
- **Cerrar un proyecto con Sprints Pendientes**: se permite (solo el Sprint Activo lo
  impide); los Sprints Pendientes quedan como están, sin poder iniciarse mientras el
  proyecto siga cerrado.
- **Votación abierta al cerrar**: se permite cerrar; la ronda queda como está, se puede
  consultar y nadie puede votar ni revelar hasta reabrir.
- **Cerrar, reabrir y cerrar el mismo día o otro día**: cada cierre guarda la fecha de ese
  cierre; la anterior no se conserva.
- **«Hoy»**: el día civil del servidor, el mismo criterio que usa HU-13 para comparar
  fechas.
- **Orden de errores** (extiende FR-006 de HU-21): 401 → 404 → 403 no integrante → 403 no
  Scrum Master → 409 proyecto cerrado → 400/409 propios del caso de uso. Un cuerpo mal
  formado sigue respondiendo 400 antes del 403, como en HU-21.
- **Un Product Builder que escribe sobre un proyecto cerrado**: en las acciones de gestión
  recibe 403 (el rol va antes); en las acciones de integrante recibe 409.
- **Crear un proyecto nuevo** (`POST /projects`) y **cargar los datos de demostración** no
  se ven afectados: no escriben sobre un proyecto existente.

## Requirements *(mandatory)*

### Reglas de negocio (EARS)

**Ubicuas**

- **FR-001**: Un proyecto SHALL estar *cerrado* si y solo si tiene fecha de cierre.
- **FR-002**: El estado de un proyecto cerrado SHALL ser «Finalizado», con prioridad sobre
  todas las reglas de HU-13.
- **FR-003**: El estado de un proyecto abierto SHALL calcularse con las reglas de HU-13, sin
  cambios.
- **FR-004**: La fecha de cierre SHALL ser una fecha de calendario (`YYYY-MM-DD`), con la
  misma validación de formato que las demás fechas del modelo (D12).
- **FR-005**: El chequeo de «proyecto cerrado» SHALL vivir en un único lugar, el control de
  acceso de HU-21 por el que pasan todas las escrituras sobre un proyecto.

**Basadas en evento**

- **FR-006**: WHEN el Scrum Master cierra un proyecto abierto sin Sprint Activo, el sistema
  SHALL guardar como fecha de cierre la fecha de hoy y responder con el proyecto
  «Finalizado».
- **FR-007**: WHEN el Scrum Master reabre un proyecto cerrado, el sistema SHALL dejar vacía
  la fecha de cierre y responder con el proyecto y su estado recalculado según HU-13.
- **FR-008**: WHEN se consulta un proyecto (detalle o «Mis proyectos») o se lo crea o
  edita, el sistema SHALL incluir su fecha de cierre en la respuesta, vacía si está
  abierto.

**Basadas en estado**

- **FR-009**: WHILE un proyecto está cerrado, el sistema SHALL rechazar toda escritura sobre
  él (historias, horas, votación, defectos, Sprints, equipo y edición del proyecto) con 409
  «El proyecto está cerrado», sin cambiar datos.
- **FR-010**: WHILE un proyecto está cerrado, el sistema SHALL seguir respondiendo todas sus
  lecturas como con el proyecto abierto.
- **FR-011**: WHILE un proyecto está cerrado, reabrirlo SHALL ser la única acción de
  escritura permitida sobre él.

**No deseadas**

- **FR-012**: IF el Scrum Master intenta cerrar un proyecto con un Sprint Activo, THEN el
  sistema SHALL responder 409 «Cerrá el Sprint activo antes de cerrar el proyecto» sin
  cerrarlo.
- **FR-013**: IF se intenta cerrar un proyecto ya cerrado, THEN el sistema SHALL responder
  409 «El proyecto ya está cerrado» sin cambiar la fecha de cierre.
- **FR-014**: IF se intenta reabrir un proyecto abierto, THEN el sistema SHALL responder 409
  «El proyecto no está cerrado» sin cambiar nada.
- **FR-015**: IF quien intenta cerrar o reabrir no es el Scrum Master del proyecto, THEN el
  sistema SHALL responder el 403 de HU-21 («No sos integrante de este proyecto» o «Solo el
  Scrum Master puede realizar esta acción») sin cambiar nada.
- **FR-016**: IF el proyecto no existe, THEN el sistema SHALL responder 404.
- **FR-017**: IF dos pedidos de cierre (o de reapertura) del mismo proyecto llegan a la vez,
  THEN solo uno SHALL realizarse y el otro SHALL responder el 409 de FR-013 (o FR-014).

### Requisitos de datos

- **FR-018**: El cambio de esquema (fecha de cierre en `proyectos`) MUST entrar en un solo
  PR con la migración 008 y la actualización de `MODELO-DE-DATOS.md` (§4.1, §5, §6, D16 en
  §7, §8 y §9); los proyectos existentes quedan abiertos.

### Requisitos de prueba

- **FR-019**: Los escenarios de Godog MUST cubrir cerrar, reabrir, el rechazo por Sprint
  Activo, los 409 de cerrar dos veces y reabrir uno abierto, los permisos (403/401/404) y
  al menos una escritura de integrante y una de Scrum Master rechazadas con 409 sobre un
  proyecto cerrado.
- **FR-020**: Las reglas de dominio (estado con cierre, transiciones cerrar/reabrir) MUST
  tener tests unitarios con ciclo RED/GREEN/REFACTOR, y el bloqueo de escrituras MUST
  probarse para cada ruta de escritura de un proyecto.

### Key Entities

- **Proyecto**: suma la **fecha de cierre** (vacía mientras está abierto). Su estado sigue
  sin guardarse: se calcula, y la fecha de cierre pasa a ser la primera regla.
- **Integrante del proyecto** (HU-21): sin cambios; define quién es el Scrum Master.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: El 100 % de las rutas de escritura sobre un proyecto responde 409 «El proyecto
  está cerrado» cuando el proyecto está cerrado, y ningún dato cambia.
- **SC-002**: El 100 % de las lecturas de un proyecto cerrado responde igual que con el
  proyecto abierto, salvo el estado («Finalizado») y la fecha de cierre.
- **SC-003**: Tras cerrar y reabrir un proyecto, su estado es exactamente el que tenía antes
  del cierre (mismas reglas de HU-13, mismos datos).
- **SC-004**: Ninguna persona que no sea el Scrum Master del proyecto logra cerrarlo ni
  reabrirlo.
- **SC-005**: Todos los escenarios de Godog (los nuevos y los existentes) y los tests de
  `go test` pasan en el CI.

## Assumptions

- **Decisiones confirmadas** (ver Clarifications): rechazo por Sprint Activo, columna
  `fecha_cierre` en la migración 008 y prioridad del cierre sobre HU-13. Cualquier otra duda
  sobre el modelo de datos se consulta al Scrum Master antes de seguir.
- Las rutas son `POST /projects/{id}/cerrar` y `POST /projects/{id}/reabrir`, de nivel
  «Sesión y Scrum Master» en la tabla de rutas de HU-21; no llevan cuerpo.
- El 409 por proyecto cerrado se responde después de los 403 de HU-21 y antes de las
  validaciones del caso de uso: no se le revela a quien no tiene permiso que el proyecto
  está cerrado.
- No se guarda quién cerró el proyecto (no lo pide el issue; YAGNI).
- La interfaz, que hará otra integrante, se apoya en la fecha de cierre que devuelven las
  consultas del proyecto (FR-008).
