# Feature Specification: Dashboard de Métricas (HU-09)

**Feature Branch**: `feature/HU-09-dashboard-metricas`

**Issue**: #69 (parte de #9)

**Created**: 2026-11-04

**Status**: Implementada, pendiente de revisión por PR y de la validación de SC-002 con una persona del equipo.

**Input**: User description: "HU-09: Dashboard de Métricas (issue #9, sub-issue #69). Como Scrum Master quiero visualizar en un panel las métricas calculadas del proyecto (HU-03) con representaciones gráficas para poder evaluar de un vistazo el estado del proyecto sin tener que leer datos crudos. Escenario 1: dado que el proyecto tiene al menos un Sprint finalizado y sus métricas calculadas, cuando ingreso a la sección 'Dashboard', el sistema renderiza un gráfico de barras de velocidad (Sprints en eje X, Puntos Comprometidos vs. Completados en eje Y) y al menos un indicador visual por cada métrica restante del punto 7 del enunciado. Escenario 2: dado que un proyecto todavía no tiene Sprints finalizados, cuando ingreso al Dashboard, el sistema muestra un estado vacío informativo, no un gráfico roto ni un error. Fuentes: docs/PRODUCT-BACKLOG.md (HU-03 y HU-09) y specs/032-hu-03-calculo-metricas."

**Dependencia**: HU-03 (issue #3, PR #270), ya en `main`. Esta historia muestra los valores que calcula HU-03 y no los recalcula.

## Clarifications

### Session 2026-11-04

Decisiones propuestas por Candela Vulcano y confirmadas en la reunión de definición de la historia:

- Q: ¿Hace falta backend nuevo? La sub-issue #71 pide un endpoint "listo para graficar", pero la consulta de métricas de HU-03 ya entrega lo necesario. → A: No. La pantalla consume la consulta de HU-03, sin cambios de backend ni de esquema (DRY, YAGNI). #70, #71 y #72 se ajustan a esta decisión.
- Q: ¿Con qué se dibujan los gráficos? La estimación habla de una librería, y #230 pide no agregar librerías salvo que el equipo lo decida. → A: Con la librería Recharts, cargada solo al entrar al Dashboard. Es la decisión de equipo que piden #230 y el Principio VI de la constitución.
- Q: ¿Qué alcance tiene el panel además del gráfico de velocidad? → A: Indicadores con los valores del proyecto en conjunto, con elementos gráficos (proporción para el porcentaje, barras comparadas para horas y defectos). Sin más gráficos por Sprint ni selector de Sprint; el detalle por Sprint es de HU-12.

**Punto abierto (no se trató en la reunión; no bloquea esta historia)**: la spec de HU-12 (033, en su propia rama) deja los "gráficos y comparaciones visuales" a HU-09, pero con este alcance solo se grafican los Story Points por Sprint. Las horas y los defectos por Sprint quedan en el resumen de HU-12, sin gráfico en ninguna historia. Falta confirmar si eso alcanza o si se suman esos gráficos; se implementa el alcance de arriba, y sumarlos después no obliga a rehacerlo.

## Objetivo

Que el Scrum Master evalúe de un vistazo el estado del proyecto: cuánto se comprometió y cuánto se completó en cada Sprint, y cómo van el esfuerzo y la calidad, sin leer datos crudos.

## Entradas y salidas

- **Entradas (de la persona)**: el proyecto y el ingreso a la sección "Dashboard". No hay filtros ni datos para cargar.
- **Entradas (del sistema)**: las métricas calculadas por HU-03: el detalle de cada Sprint finalizado, los valores del proyecto en conjunto y las horas por Story Point del proyecto.
- **Salida en pantalla**: el gráfico de velocidad y los indicadores; o el estado vacío informativo si no hay Sprints finalizados; o un mensaje de error si las métricas no se pudieron obtener.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Comparar lo comprometido y lo completado en cada Sprint (Priority: P1)

Como Scrum Master, ingreso a la sección "Dashboard" de un proyecto y veo un gráfico de barras con los Story Points comprometidos y completados de cada Sprint finalizado, para evaluar de un vistazo la capacidad real del equipo.

**Why this priority**: Es el gráfico que pide de forma explícita el Escenario 1 de la historia y el requisito central del punto 8 del enunciado.

**Independent Test**: Con un proyecto de tres Sprints finalizados con valores conocidos, se ingresa al Dashboard y se verifica que cada Sprint aparece en orden con sus dos valores.

**Acceptance Scenarios**:

1. **Given** un proyecto con Sprints finalizados, **When** elijo "Dashboard" en la navegación del proyecto, **Then** veo un gráfico de barras con un grupo por Sprint en el eje X y los Story Points en el eje Y.
2. **Given** un Sprint finalizado con 8 Story Points planificados y 5 completados, **When** veo el gráfico, **Then** el grupo de ese Sprint tiene una barra "Comprometidos" de 8 y una barra "Completados" de 5.
3. **Given** tres Sprints finalizados, **When** veo el gráfico, **Then** aparecen ordenados por su número dentro del proyecto.
4. **Given** un Sprint finalizado con cero Story Points comprometidos, **When** veo el gráfico, **Then** ese Sprint sigue en el eje X, con sus valores en cero.
5. **Given** un proyecto con Sprints finalizados y un Sprint "Activo" o "Pendiente", **When** veo el gráfico, **Then** solo aparecen los finalizados.

---

### User Story 2 - Ver el resto de las métricas en indicadores (Priority: P1)

Como Scrum Master, veo en el mismo panel un indicador visual por cada métrica restante, para saber cómo van la velocidad, el esfuerzo y la calidad sin abrir cada Sprint.

**Why this priority**: Es la otra mitad del Escenario 1: sin los indicadores, el panel no cubre el punto 7 del enunciado.

**Independent Test**: Con un proyecto de valores calculados a mano, se verifica que cada indicador muestra el valor informado por HU-03, con el mismo redondeo.

**Acceptance Scenarios**:

1. **Given** un proyecto con Sprints finalizados, **When** ingreso al Dashboard, **Then** veo un indicador de velocidad, de porcentaje de historias completadas, de horas estimadas, de horas reales, de desviación, de defectos detectados y de defectos resueltos, con los valores del proyecto en conjunto.
2. **Given** un proyecto con 50% de historias completadas, **When** veo su indicador, **Then** el valor se lee también de forma gráfica, como una proporción de un total.
3. **Given** 32 horas estimadas y 30 horas reales, **When** veo los indicadores de horas, **Then** se comparan visualmente en una misma escala.
4. **Given** una desviación de -2 horas (-6,25%), **When** veo su indicador, **Then** veo ambos valores y la aclaración de que se trabajó menos de lo estimado.
5. **Given** 2 defectos detectados y 1 resuelto, **When** veo sus indicadores, **Then** se comparan visualmente en una misma escala.
6. **Given** un proyecto sin horas por Story Point, **When** ingreso al Dashboard, **Then** las horas estimadas y la desviación se informan como "sin datos", nunca como cero; las horas reales y los demás indicadores se muestran igual, y se indica que el valor se define al editar el proyecto.

---

### User Story 3 - Entender por qué todavía no hay métricas (Priority: P1)

Como Scrum Master de un proyecto que todavía no cerró ningún Sprint, veo un estado vacío que me explica la situación, para no confundir la falta de datos con un error.

**Why this priority**: Es el Escenario 2 de la historia y el estado inicial de todo proyecto.

**Independent Test**: Con un proyecto sin Sprints finalizados, se ingresa al Dashboard y se verifica que aparece el estado vacío, sin gráfico, sin indicadores en cero y sin error.

**Acceptance Scenarios**:

1. **Given** un proyecto sin Sprints, **When** ingreso al Dashboard, **Then** veo un mensaje que explica que las métricas aparecen cuando se cierra el primer Sprint, y no veo el gráfico ni los indicadores.
2. **Given** un proyecto cuyo único Sprint está "Activo", **When** ingreso al Dashboard, **Then** veo el mismo estado vacío.
3. **Given** el estado vacío, **When** lo leo, **Then** tengo un acceso a la sección de Sprints del proyecto.

---

### User Story 4 - Saber qué pasó si las métricas no se pueden obtener (Priority: P2)

Como Scrum Master, si la consulta falla veo un mensaje claro, para no interpretar un panel incompleto como el estado real del proyecto.

**Why this priority**: Es un camino de excepción, pero mostrar valores parciales o en cero llevaría a evaluar mal el proyecto.

**Independent Test**: Se simula cada falla de la consulta y se verifica el mensaje, sin gráfico ni indicadores.

**Acceptance Scenarios**:

1. **Given** una consulta en curso, **When** ingreso al Dashboard, **Then** veo un indicador de carga.
2. **Given** un proyecto que no existe, **When** ingreso a su Dashboard, **Then** veo el mensaje "proyecto no encontrado".
3. **Given** un servidor que no responde, **When** ingreso al Dashboard, **Then** veo un mensaje de error de conexión, sin gráfico ni indicadores.

---

### Edge Cases

Redactados en formato EARS (Ubicuo / Evento / Estado / No deseado / Opcional).

- El sistema **shall** mostrar los valores tal como los informa el cálculo de HU-03, y **shall** no recalcular, redondear de otra forma ni guardar ninguna métrica.
- El sistema **shall** decidir el estado vacío por la ausencia de Sprints finalizados, y no por los valores en cero del proyecto.
- **WHILE** no hay Sprints finalizados, el sistema **shall** mostrar el estado vacío y **shall** no mostrar gráficos ni indicadores en cero.
- **IF** el proyecto no tiene horas por Story Point, **THEN** el sistema **shall** informar las horas estimadas y la desviación como "sin datos", y **shall** no dibujarlas como cero.
- **IF** la desviación porcentual es "sin datos" y la desviación en horas tiene valor, **THEN** el sistema **shall** mostrar el valor en horas y "sin datos" para el porcentaje.
- **WHERE** la desviación es positiva, el sistema **shall** indicar que se trabajó más de lo estimado; **WHERE** es negativa, que se trabajó menos; **WHERE** es cero, que no hubo desviación.
- **IF** todos los Sprints finalizados tienen cero Story Points, **THEN** el sistema **shall** mostrar el gráfico con sus Sprints y sus valores en cero, sin fallar.
- **IF** un proyecto con Sprints finalizados no tiene defectos u horas imputadas, **THEN** el sistema **shall** mostrar esos indicadores en cero, porque cero es el valor real, y **shall** no tratarlos como "sin datos".
- **IF** las horas o los defectos comparados suman cero, **THEN** el sistema **shall** mostrar el indicador sin barras rotas ni proporciones indefinidas.
- **WHEN** se ingresa al Dashboard, **THEN** el sistema **shall** consultar las métricas en ese momento; **shall** no actualizarlas solo mientras la pantalla está abierta.
- **WHEN** se cierra un Sprint o se cambian las horas por Story Point y se vuelve al Dashboard, **THEN** el sistema **shall** mostrar los valores nuevos.
- **WHILE** el proyecto tiene muchos Sprints finalizados, el sistema **shall** mantener legibles los rótulos y **shall** no desbordar el ancho de la página, incluso en un ancho de 320 px.
- **IF** la consulta responde "proyecto no encontrado", **THEN** el sistema **shall** mostrar ese mensaje, y **shall** no mostrar el estado vacío.
- **IF** la consulta falla por conexión o por error interno, **THEN** el sistema **shall** mostrar un mensaje de error y **shall** no mostrar el estado vacío ni valores anteriores.
- **IF** la respuesta llega sin métricas o el error no proviene de la API, **THEN** el sistema **shall** mostrar "Ocurrió un error inesperado", y **shall** no quedar cargando ni mostrar textos técnicos.

## Requirements *(mandatory)*

### Functional Requirements

**Sección Dashboard**

- **FR-001**: El sistema MUST ofrecer la sección "Dashboard" en la navegación de cada proyecto.
- **FR-002**: El Dashboard MUST ser de solo lectura: ingresar a él no modifica ningún dato.
- **FR-003**: El sistema MUST obtener las métricas al ingresar a la sección, con una sola consulta, y MUST mostrar un indicador de carga mientras tanto.

**Gráfico de velocidad**

- **FR-004**: El sistema MUST mostrar un gráfico de barras con un grupo por cada Sprint finalizado en el eje X, y dos valores en el eje Y: "Comprometidos" (Story Points planificados de HU-03) y "Completados".
- **FR-005**: Los Sprints MUST aparecer ordenados por su número dentro del proyecto e identificados con ese número.
- **FR-006**: El gráfico MUST distinguir las dos series sin depender solo del color, con una leyenda siempre a la vista, aunque el gráfico se desplace, y sus valores MUST poder leerse como texto, sin depender del puntero.

**Indicadores**

- **FR-007**: El sistema MUST mostrar, con los valores del proyecto en conjunto, al menos un indicador visual por cada métrica: velocidad, porcentaje de historias completadas, horas estimadas, horas reales, desviación, defectos detectados y defectos resueltos.
- **FR-008**: El porcentaje de historias completadas MUST representarse gráficamente como una proporción de un total, además de su valor.
- **FR-009**: Las horas estimadas y reales MUST compararse en una misma escala, igual que los defectos detectados y resueltos, además de mostrar sus valores. Si falta uno de los dos términos, no se dibuja ninguna barra, porque no hay contra qué comparar.
- **FR-010**: La desviación MUST mostrarse en horas y en porcentaje, con una aclaración en palabras de si se trabajó más, menos o lo mismo que lo estimado, que no dependa solo del color.
- **FR-011**: Los valores decimales MUST mostrarse con coma decimal, igual que en el resto de la aplicación.

**Datos ausentes**

- **FR-012**: Cuando una métrica se informa como "sin datos", su indicador MUST mostrar "sin datos" y MUST no dibujarse como cero.
- **FR-013**: Si el proyecto no tiene horas por Story Point, el Dashboard MUST indicar que ese valor se define al editar el proyecto.

**Estado vacío**

- **FR-014**: Si el proyecto no tiene Sprints finalizados, el sistema MUST mostrar un estado vacío que explique que las métricas aparecen al cerrar el primer Sprint, con un acceso a la sección de Sprints, y MUST no mostrar el gráfico ni los indicadores.

**Errores**

- **FR-015**: Si el proyecto no existe, el sistema MUST mostrar "proyecto no encontrado".
- **FR-016**: Si la consulta falla, el sistema MUST mostrar un mensaje de error, sin gráfico, sin indicadores y sin estado vacío.

**Restricciones**

- **FR-017**: El sistema MUST mostrar los valores que informa HU-03 y MUST no recalcularlos, para que coincidan con los del reporte (HU-10) y el historial (HU-12).
- **FR-018**: El Dashboard MUST usarse sin desbordar la página en un ancho de 320 px en adelante.

### Key Entities *(include if feature involves data)*

- **Métricas del proyecto**: las define HU-03; se calculan al consultar y no se guardan. Esta historia no agrega entidades, tablas ni columnas.
- **Métricas de un Sprint finalizado**: las mismas métricas para cada Sprint, identificado por su número dentro del proyecto.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Para un proyecto de referencia con valores calculados a mano, el 100% de los valores del gráfico y de los indicadores coincide con los que informa HU-03.
- **SC-002**: Un Scrum Master identifica en menos de 30 segundos, sin ayuda, en qué Sprint se completó menos de lo comprometido y si el proyecto trabajó más o menos horas de lo estimado.
- **SC-003**: El 100% de los proyectos sin Sprints finalizados muestra el estado vacío, sin gráfico roto ni error.
- **SC-004**: El 100% de los proyectos sin horas por Story Point informa las horas estimadas y la desviación como "sin datos", nunca como cero.
- **SC-005**: Cada una de las métricas restantes del punto 7 del enunciado tiene al menos un indicador visual en el panel.
- **SC-006**: El panel se usa completo en un ancho de 320 px, sin desplazamiento horizontal de la página.

## Assumptions

- **Valor que se muestra**: los indicadores muestran los valores del proyecto en conjunto y el gráfico aporta el detalle por Sprint. El resumen completo de cada Sprint es de HU-12.
- **Métricas del punto 7**: el enunciado no está en el repo; se toma la lista del backlog (HU-03): Story Points planificados y completados, velocidad, horas estimadas y reales, desviación, porcentaje de historias completadas y defectos detectados y resueltos.
- **"Comprometidos"**: son los Story Points planificados de HU-03, las historias asignadas al cerrar el Sprint, completadas o no (nota de HU-09 en el backlog).
- **Actualización**: las métricas se recalculan a pedido y no en tiempo real (`AGENTS.md`, YAGNI).
- **Acceso**: el Dashboard no restringe por rol; cuando HU-20 y HU-21 protejan la consulta, los errores de sesión o de permisos se muestran como cualquier otro error de la consulta.
- **Sección "Estado"**: el placeholder de la navegación no se reemplaza en esta historia; el estado del proyecto (HU-13) ya está en el panel principal.
- **Rediseño visual (#230)**: el Dashboard usa la identidad visual vigente; si el rediseño cambia la navegación, se ajusta ahí.
- **Fuera de alcance**: cambios en el cálculo de las métricas (HU-03), reporte en PDF (HU-10), historial de Sprints (HU-12), filtros por fecha o por Sprint, exportar los gráficos, y métricas de un Sprint "Activo".
- **Dependencias**: HU-03, con el backend y el campo de horas por Story Point ya en `main`.
