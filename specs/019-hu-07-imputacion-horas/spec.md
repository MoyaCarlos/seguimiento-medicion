# Feature Specification: Imputación de Horas Trabajadas (HU-07)

**Feature Branch**: `feature/HU-07-imputacion-horas`

**Created**: 2026-10-28

**Status**: Draft

**Input**: User description: "HU-07: Imputación de Horas Trabajadas (issue #7). Como integrante del equipo quiero registrar las horas reales que invertí trabajando en una Historia de Usuario para poder proveer datos al sistema que permitan comparar el esfuerzo estimado (Story Points convertidos a horas, ver HU-03) contra el esfuerzo real. Escenario 1: registro válido de tiempo (integrante, fecha de hoy, actividad y 4 horas sobre una historia de un Sprint activo; se suman al esfuerzo acumulado de la historia y se registra el Sprint activo en ese momento). Escenario 2: validación de inputs (horas negativas, texto no numérico, actividad vacía o integrante que no pertenece al proyecto bloquean el guardado). Depende de HU-04, HU-06 y HU-05. Las horas se guardan como minutos enteros (D8). El tope de horas por carga lo define esta especificación. Alcance: backend; el frontend es la sub-issue #63."

## Clarifications

### Session 2026-10-28

- Q: ¿Cuál es el tope de horas de una sola carga? → A: 24 horas (un día completo), con el borde inclusivo. Frena errores de tipeo (400 en vez de 4) sin limitar casos reales. Es la regla que D8 de `MODELO-DE-DATOS.md` deja a esta spec; vive en el dominio, no en la base.
- Q: ¿En qué unidad y con qué precisión se ingresan las horas? → A: En horas con decimales, con precisión de un minuto: se convierten a minutos enteros (1,5 h = 90 min) y se rechaza lo que no da minutos enteros (0,01 h = 0,6 min).
- Q: ¿Qué fechas acepta una carga? → A: Desde la fecha de inicio del Sprint activo hasta hoy, ambas inclusive. Ni fechas futuras ni anteriores al Sprint al que se atribuyen las horas. (El límite superior se precisa en la pregunta sobre la fecha de fin, más abajo.)
- Q: ¿Se pueden imputar horas a una historia "Completada" que sigue en el Sprint activo? → A: Sí. Mientras la historia esté en el Sprint "Activo", la carga se acepta cualquiera sea su estado; solo importa el Sprint.
- Q: Si el Sprint sigue "Activo" después de su fecha de fin, ¿se aceptan cargas con fecha posterior a esa fecha de fin? → A: No. El rango aceptado va desde la fecha de inicio del Sprint hasta la menor entre hoy y su fecha de fin, ambas inclusive, para que las horas caigan siempre dentro del período del Sprint al que se atribuyen.
- Q: ¿Se limita el total de horas que un integrante puede imputar en un mismo día, sumando sus cargas? → A: No. Solo rige el tope de 24 horas por carga; no hay tope diario por integrante.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Imputar horas trabajadas en una historia (Priority: P1)

Como integrante del equipo, elijo una historia que está en el Sprint activo del proyecto, selecciono mi nombre, indico la fecha en que trabajé, describo la actividad realizada y la cantidad de horas, y guardo, para que el esfuerzo real de esa historia quede registrado.

**Why this priority**: Es el núcleo de la historia. Sin horas reales no hay forma de calcular la desviación de esfuerzo (HU-03) ni el esfuerzo por Sprint de los reportes.

**Independent Test**: Se prueba sembrando un proyecto con un integrante, un Sprint "Activo" y una historia asignada a ese Sprint; se imputan 4 horas y se verifica que la carga queda registrada con integrante, fecha, actividad, duración y Sprint, y que el esfuerzo acumulado de la historia sube 4 horas.

**Acceptance Scenarios**:

1. **Given** que soy integrante del proyecto y la historia pertenece al Sprint "Activo", **When** imputo 4 horas con mi nombre, la fecha de hoy y la actividad realizada, **Then** el sistema registra la carga con integrante, fecha, actividad, duración y el Sprint activo en ese momento, y el esfuerzo total acumulado de la historia aumenta en 4 horas.
2. **Given** que una historia ya tiene horas imputadas, **When** otro integrante (o el mismo) imputa más horas, **Then** el esfuerzo acumulado es la suma de todas las cargas de la historia.
3. **Given** que imputé horas a una historia durante un Sprint, **When** la historia pasa después a otro Sprint (por ejemplo, vuelve al Product Backlog al cerrar el Sprint y se planifica en el siguiente), **Then** esas horas siguen atribuidas al Sprint en el que se cargaron.

---

### User Story 2 - Rechazar cargas inválidas (Priority: P1)

Como integrante del equipo, si cargo un dato inválido el sistema bloquea el guardado y me dice qué campo está mal, para que el esfuerzo real no se contamine con datos erróneos.

**Why this priority**: Una carga inválida distorsiona todas las métricas que se calculan a partir de las horas reales; el backlog lo pide como criterio de aceptación explícito.

**Independent Test**: Se prueba intentando imputar con cada dato inválido por separado y verificando que la respuesta identifica el campo y que el esfuerzo acumulado de la historia no cambia.

**Acceptance Scenarios**:

1. **Given** que intento cargar horas a una historia, **When** ingreso una cantidad negativa (-2 horas) o cero, **Then** el sistema rechaza la carga con un error sobre las horas y no registra nada.
2. **Given** que intento cargar horas a una historia, **When** ingreso un texto no numérico como cantidad de horas, **Then** el sistema rechaza la carga con un error sobre las horas y no registra nada.
3. **Given** que intento cargar horas a una historia, **When** dejo la actividad vacía o compuesta solo por espacios, **Then** el sistema rechaza la carga con un error sobre la actividad y no registra nada.
4. **Given** que intento cargar horas a una historia, **When** elijo una persona que no es integrante del proyecto de la historia, **Then** el sistema rechaza la carga con un error sobre el integrante y no registra nada.
5. **Given** que una historia no está en un Sprint "Activo" (está en el Product Backlog, en un Sprint "Pendiente" o en uno "Finalizado"), **When** intento imputarle horas, **Then** el sistema rechaza la carga con un mensaje que indica que solo se imputan horas a historias del Sprint activo.
6. **Given** que la historia no existe, **When** intento imputarle horas, **Then** el sistema responde "historia no encontrada".

---

### User Story 3 - Consultar el esfuerzo real de una historia (Priority: P2)

Como integrante del equipo, consulto una historia y veo su esfuerzo total acumulado y el detalle de las cargas (quién, cuándo, qué actividad, cuánto y en qué Sprint), para saber cuánto trabajo real lleva.

**Why this priority**: El Escenario 1 pide que las horas "se sumen al esfuerzo total acumulado" y la sub-issue de Frontend pide mostrarlo; sin una consulta, el acumulado no se puede ver ni verificar. Depende de que exista la carga (P1).

**Independent Test**: Se prueba sembrando varias cargas en una historia (de distintos integrantes y Sprints) y verificando que la consulta devuelve el total exacto y el detalle de cada carga; una historia sin cargas devuelve total cero y detalle vacío.

**Acceptance Scenarios**:

1. **Given** que una historia tiene varias cargas, **When** consulto su esfuerzo real, **Then** el sistema muestra el total acumulado en horas y el detalle de cada carga con integrante, fecha, actividad, duración y Sprint.
2. **Given** que una historia existe y no tiene cargas, **When** consulto su esfuerzo real, **Then** el sistema muestra un total de 0 horas y ningún detalle (no es un error).
3. **Given** que la historia no existe, **When** consulto su esfuerzo real, **Then** el sistema responde "historia no encontrada".

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** guardar la duración de cada carga como minutos enteros y mostrarla convertida a horas (D8).
- **WHEN** se imputan horas válidas a una historia del Sprint "Activo" de su proyecto, **THEN** el sistema **shall** registrar la carga atribuida a ese Sprint y sumarla al esfuerzo acumulado de la historia.
- **WHEN** una historia con horas imputadas cambia de Sprint o vuelve al Product Backlog, **THEN** el sistema **shall** mantener cada carga atribuida al Sprint en el que se registró; el Sprint de una carga no se recalcula nunca.
- **WHILE** la historia no tiene Sprint asignado o su Sprint está "Pendiente" o "Finalizado", **IF** se intenta imputarle horas, **THEN** el sistema **shall** rechazar la carga con un mensaje que indique que solo se imputan horas a historias del Sprint activo.
- **WHILE** la historia está en el Sprint "Activo", el sistema **shall** aceptar la carga cualquiera sea el estado de la historia ("Nueva", "En progreso" o "Completada").
- **IF** la cantidad de horas es cero, negativa, no numérica o no se informa, **THEN** el sistema **shall** rechazar la carga con un error sobre las horas.
- **IF** la cantidad de horas supera 24 horas (el tope por carga), **THEN** el sistema **shall** rechazar la carga con un error sobre las horas que indique el tope. Exactamente 24 horas se acepta.
- **WHERE** la cantidad de horas tiene decimales, el sistema **shall** convertirla a minutos (1,5 h = 90 min; 0,25 h = 15 min).
- **IF** la cantidad de horas no se puede expresar en minutos enteros (por ejemplo, 0,01 h = 0,6 min), **THEN** el sistema **shall** rechazar la carga con un error sobre las horas. La carga mínima válida es 1 minuto.
- **IF** la actividad está vacía, compuesta solo por espacios o supera los 500 caracteres, **THEN** el sistema **shall** rechazar la carga con un error sobre la actividad.
- **IF** la fecha falta, no tiene el formato `AAAA-MM-DD` o no es una fecha real (por ejemplo, 2026-02-30), **THEN** el sistema **shall** rechazar la carga con un error sobre la fecha.
- **IF** la fecha es posterior a hoy, posterior a la fecha de fin del Sprint activo o anterior a su fecha de inicio, **THEN** el sistema **shall** rechazar la carga con un error sobre la fecha. Los bordes se aceptan: la fecha de inicio del Sprint y la menor entre hoy y su fecha de fin.
- **WHILE** el Sprint sigue "Activo" después de su fecha de fin, **WHEN** se imputan horas con una fecha entre su inicio y su fin, **THEN** el sistema **shall** aceptar la carga; solo se rechazan las fechas posteriores al fin.
- **WHEN** un integrante imputa varias cargas con la misma fecha (en la misma historia o en distintas), **THEN** el sistema **shall** aceptarlas todas aunque sumen más de 24 horas: el tope es por carga, no por día.
- **IF** la persona elegida no es integrante del proyecto de la historia (existe en otro proyecto o no existe), **THEN** el sistema **shall** rechazar la carga con un error sobre el integrante.
- **WHEN** se imputan horas a una historia que no existe, **THEN** el sistema **shall** responder "historia no encontrada" antes de validar los datos de la carga (la existencia domina sobre la validación, igual que en HU-04, HU-05, HU-06 y HU-15).
- **IF** una carga incumple varias reglas a la vez, **THEN** el sistema **shall** responder con la primera según este orden: historia inexistente, historia fuera del Sprint activo, integrante ajeno al proyecto, horas inválidas, actividad inválida, fecha inválida.
- **IF** una carga es inválida, **THEN** el sistema **shall** rechazarla entera y el esfuerzo acumulado de la historia **shall** quedar igual que antes.
- **WHEN** dos integrantes imputan horas a la misma historia al mismo tiempo, **THEN** el sistema **shall** registrar ambas cargas y el acumulado **shall** ser la suma de las dos.
- **WHEN** el Sprint se cierra mientras se envía una carga, **THEN** la carga **shall** quedar registrada en ese Sprint si llegó antes del cierre, o rechazada por estar fuera del Sprint activo si llegó después; nunca registrada sin Sprint ni en un Sprint distinto al de la historia en ese momento.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST permitir imputar horas a una historia que pertenece al Sprint "Activo" de su proyecto, indicando integrante, fecha, actividad y cantidad de horas.
- **FR-002**: Cada carga MUST registrar el integrante, la historia, la fecha trabajada, la actividad, la duración y el Sprint activo en el momento de la carga.
- **FR-003**: El Sprint de una carga MUST ser el Sprint "Activo" al que pertenecía la historia al cargarla, y MUST mantenerse aunque la historia cambie de Sprint después.
- **FR-004**: El sistema MUST rechazar la carga cuando la historia no pertenece a un Sprint "Activo".
- **FR-005**: El sistema MUST rechazar la carga cuando la persona elegida no es integrante del proyecto de la historia.
- **FR-006**: La cantidad de horas MUST ingresarse en horas, admitiendo decimales, y guardarse convertida a minutos enteros. El sistema MUST rechazar la carga cuando la cantidad no es un número, es cero o negativa, supera las 24 horas (tope por carga) o no da un número entero de minutos.
- **FR-007**: El sistema MUST rechazar la carga cuando la actividad está vacía, tiene solo espacios o supera los 500 caracteres.
- **FR-008**: El sistema MUST rechazar la carga cuando la fecha falta, no es una fecha válida, es posterior a hoy, es posterior a la fecha de fin del Sprint activo o es anterior a su fecha de inicio.
- **FR-009**: Cada rechazo MUST identificar el campo con problemas (historia, integrante, horas, actividad o fecha) y MUST dejar sin cambios las cargas y el esfuerzo acumulado.
- **FR-010**: El sistema MUST responder "historia no encontrada" cuando la historia no existe, antes de validar cualquier dato de la carga.
- **FR-011**: Ante varias reglas incumplidas, el sistema MUST informar la primera según el orden de los Edge Cases, para que la respuesta sea siempre la misma ante la misma petición.
- **FR-012**: El sistema MUST permitir consultar el esfuerzo real de una historia existente: el total acumulado en horas y el detalle de cada carga (integrante, fecha, actividad, duración y Sprint), ordenado por fecha trabajada y, ante la misma fecha, por orden de carga.
- **FR-013**: El esfuerzo acumulado de una historia MUST ser la suma exacta de las duraciones de todas sus cargas, sin errores de redondeo.
- **FR-014**: Las cargas válidas MUST persistir y seguir disponibles después de cerrar y volver a abrir la aplicación.
- **FR-015**: Las cargas no se editan ni se borran en esta historia.

### Key Entities *(include if feature involves data)*

- **Registro de horas (carga)**: una imputación de trabajo. Atributos: integrante, historia, Sprint al que se atribuye, fecha trabajada, actividad y duración (en minutos enteros; se muestra en horas). Pertenece a un proyecto; el integrante, la historia y el Sprint son siempre del mismo proyecto.
- **Historia de Usuario**: la de HU-01/HU-15. Acumula cargas; su esfuerzo real es la suma de ellas. Para recibir una carga, tiene que estar en el Sprint "Activo".
- **Sprint**: el de HU-05. Solo un Sprint "Activo" admite cargas. Las cargas atribuidas a un Sprint forman su esfuerzo real, que usan las métricas (HU-03).
- **Integrante**: persona asignada al proyecto (HU-04). Solo los integrantes del proyecto de la historia pueden imputarle horas.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un integrante puede registrar las horas de una actividad en menos de 1 minuto y en una sola acción de guardado.
- **SC-002**: El 100% de las cargas con horas negativas, cero o no numéricas, actividad vacía, fecha inválida, integrante ajeno al proyecto o historia fuera del Sprint activo se rechazan con un error que identifica el campo, y el esfuerzo acumulado queda igual.
- **SC-003**: El esfuerzo acumulado de una historia coincide exactamente con la suma de sus cargas (por ejemplo, tres cargas de 0,5 h, 1,5 h y 2 h dan exactamente 4 h).
- **SC-004**: El 100% de las cargas siguen atribuidas al Sprint en que se registraron después de que la historia cambie de Sprint o se cierre ese Sprint.
- **SC-005**: Las cargas válidas siguen disponibles después de recargar la aplicación.

## Assumptions

- No hay autenticación: igual que en el resto del sistema, "selecciono mi nombre" significa elegir un integrante del proyecto en la carga. El sistema no verifica que quien carga sea esa persona.
- Asignar historias a un Sprint llega con HU-06 (rama `feature/HU-06-asignar-historias-sprint`, todavía sin mergear). Hasta entonces, los tests siembran historias con Sprint directamente en la base o en los fakes, igual que hizo HU-15.
- "Hoy" es la fecha local del servidor al recibir la carga, igual que en el estado del proyecto (HU-13). Los tests la fijan en lugar de leer el reloj.
- Un Sprint "Activo" siempre tiene fecha de inicio y de fin (lo garantiza HU-05), así que el rango de fechas aceptado está siempre definido.
- Un proyecto tiene como máximo un Sprint "Activo" (HU-05), así que "el Sprint activo" de una historia es el Sprint al que está asignada, cuando ese Sprint está "Activo".
- La tabla de cargas ya está diseñada en `MODELO-DE-DATOS.md` (§4.8 `registros_horas`, migración 004, D8 y D12): minutos enteros mayores que cero, fecha con formato validado y claves compuestas que garantizan que integrante, historia y Sprint son del mismo proyecto. Esta historia agrega esa migración y no cambia el esquema.
- El largo máximo de la actividad (500 caracteres) es una regla del dominio, no de la base, igual que los largos de título y descripción (§2 del modelo de datos).
- Corregir o borrar una carga queda fuera de esta historia: ninguna historia del backlog borra datos (D9). Si hace falta, será una historia nueva.
- Las métricas que comparan esfuerzo estimado y real (desviación, horas por Sprint) son de HU-03; esta historia solo provee los datos y el acumulado por historia.
- El frontend (formulario, errores en pantalla y acumulado visible) es la sub-issue #63 y se entrega en otro Sprint; esta spec cubre el comportamiento del sistema que esa pantalla va a usar. El mensaje "muestra un error de validación en la interfaz de React" del Escenario 2 se cumple en el backend devolviendo un error que identifica el campo.
