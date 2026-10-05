# REQUISITOS DE USUARIO

> **Convención de campos:**
> - **Prioridad** usa escala MoSCoW (Must/Should/Could/Won't have). Al cargar en el
>   tablero de GitHub Projects, mapear: `M` → Alta, `S` → Media, `C` → Baja.
> - **Valor de Negocio** y **Estimación** usan ambas una escala tipo Fibonacci, pero
>   miden cosas distintas e independientes (como en WSJF): Valor de Negocio es el
>   impacto/beneficio relativo para el proyecto; Estimación es el esfuerzo/complejidad
>   relativa (Story Points). No están relacionadas entre sí ni deben coincidir.
> - **Estado** inicial de toda historia recién creada: `Nueva`.
> - **Depende de** declara qué otras historias (y qué piezas de ellas) necesita esta
>   historia para poder completarse. Toda historia lo lleva.

> **Decisiones transversales** (valen para todas las historias; si una historia las
> contradice, mandan estas):
> - Todos los identificadores son números enteros (`int64`).
> - Si un recurso no existe se responde "no encontrado" (404) antes de validar el
>   cuerpo de la petición (400).
> - Estados de una historia de usuario: `Nueva` → `En progreso` → `Completada`. Es el
>   único vocabulario: no usar "Terminado" ni "Done".
> - La fecha de fin debe ser estrictamente posterior a la de inicio, tanto en
>   proyectos como en Sprints.
> - Los estados derivados (como el estado del proyecto) y las métricas se calculan al
>   consultar; no se guardan.
> - Todo requisito mínimo del enunciado tiene una historia dueña. Las historias
>   14 y 15 se agregaron para cubrir requisitos que ninguna historia tenía.

## ***Prioridad 1: Cimientos del Sistema (Dependencia Crítica)***

No se puede gestionar nada si no existe un contenedor. Esta es la base estructural.

* **HU-04: Creación de Proyecto y Asignación de Equipo (Must have).** Es la máxima prioridad operativa. Todo el sistema (historias, sprints, horas) depende de que exista un proyecto y usuarios registrados.  
* **HU-01: Creación de Historias de Usuario (Must have).** Sin esto, no hay Product Backlog.
* **HU-15: Consulta y Edición del Product Backlog (Must have).** Completa el punto 2 del enunciado: listar el backlog y mantener los criterios de aceptación de cada historia. HU-06 necesita este listado.
* **HU-13: Fechas y Estado del Proyecto (Should have).** Extiende HU-04: consultar el estado general del proyecto (Planificado / En curso / Finalizado) a partir de sus fechas y Sprints; el registro de las fechas de inicio/fin lo cubre HU-04 (requisito explícito del punto 1 del enunciado). Para el estado de los Sprints usa HU-05.

## ***Prioridad 2: El Motor de Scrum (Iteraciones)***

Una vez que tienen el proyecto y el Backlog, necesitan el ciclo de trabajo ágil.

* **HU-05: Apertura y Cierre de Sprints (Must have).** Da vida a la metodología.  
* **HU-06: Movimiento de Historias al Sprint Backlog (Must have).** Permite a las Product Builders definir en qué se va a trabajar durante la iteración.
* **HU-11: Actualización de Estado de una Historia en el Sprint (Must have).** Sin poder marcar una historia como completada durante el sprint, no hay datos reales para calcular velocidad ni ninguna métrica.

## ***Prioridad 3: Alimentación de Datos (El Valor de Negocio)***

Para que el diferenciador del sistema (las métricas) funcione, primero necesitan generar y registrar datos reales.

* **HU-02: Votación de Esfuerzo en Planning Poker (Should have).** Aporta un gran valor metodológico: registrar las estimaciones individuales, mantenerlas ocultas y mostrarlas.  
* **HU-14: Consenso y Registro de la Estimación Acordada (Should have).** Completa el Planning Poker del punto 4 del enunciado: detectar diferencias, nuevas rondas y registrar la estimación acordada, que son los "Story Points" estimados.
* **HU-07: Imputación de Horas Trabajadas (Must have).** Requisito indispensable para luego poder cruzar lo estimado vs. lo real.

## ***Prioridad 4: Resultados y Calidad (El Valor Final)***

Estas funcionalidades son el corazón del valor para el cliente final, pero tienen una alta dependencia técnica de las etapas anteriores. Se dividió la antigua HU-03 en tres historias más chicas (cálculo, visualización y exportación son responsabilidades distintas) y se agregó cobertura completa del punto 7 del enunciado.

* **HU-03: Cálculo de Métricas del Proyecto (Must have).** El motor de cálculo (Story Points planificados/completados, velocidad, horas estimadas/reales, desviación, % completadas, defectos detectados/resueltos) — requisito central del enunciado (punto 7).
* **HU-09: Dashboard de Métricas (Must have).** Visualización gráfica de lo que calcula HU-03 (punto 8 del enunciado).
* **HU-10: Generación de Reporte en PDF (Must have).** Exportación a PDF con Maroto de historias, estimaciones, esfuerzo, métricas y defectos (punto 9 del enunciado).
* **HU-08: Registro y Seguimiento de Bugs (Should have).** Completa el ciclo de aseguramiento de calidad del software desarrollado.
* **HU-12: Consulta de Sprints Anteriores (Should have).** Requisito explícito del punto 3 del enunciado. Va al final porque su historial muestra métricas (HU-03), esfuerzo (HU-07) y defectos (HU-08), que tienen que existir antes.

---

# HISTORIAS DE USUARIO

## HU-01: Creación de Historias de Usuario 

### Descripción
Como Product Builder quiero crear historias de usuario con prioridad, estado y estimación para poder alimentar y organizar el Product Backlog del proyecto.

### Criterios de aceptacion (BDD):
* Escenario 1: Creación exitosa.
Dado que me encuentro en el panel del Product Backlog de un proyecto existente,
Cuando ingreso los datos obligatorios (título, descripción, prioridad) y presiono "Guardar",
Entonces la historia debe quedar registrada de forma persistente con estado "Nueva" (el listado y la edición se cubren en HU-15).
* Escenario 2: Faltan campos.
Dado que intento crear una historia,
Cuando dejo el campo "Título" en blanco y presiono "Guardar",
Entonces el sistema debe mostrar una advertencia y no debe registrar la HU.
* Escenario 3: Proyecto inexistente.
Dado que intento crear una historia para un proyecto que no existe,
Cuando presiono "Guardar",
Entonces el sistema responde "proyecto no encontrado" y no registra la HU.

### Depende de
HU-04 (el proyecto debe existir).

### Prioridad
M (Must have) - Esencial para el MVP. 

### Estado
Implementada (el Escenario 3 está pendiente: lo cubre la tarea de integridad referencial).

### Valor de Negocio
21 (Fibonacci) 

### Estimación
5 Story Points 

---

## HU-02: Votación de Esfuerzo en Planning Poker 
### Descripción
Como integrante del equipo (Scrum Master o Product Builder) quiero votar el esfuerzo de una historia en secreto usando cartas virtuales para poder realizar una estimación colaborativa sin sesgos. El análisis de las diferencias, las nuevas rondas y el registro de la estimación acordada los cubre HU-14.
### Criterios de aceptacion (BDD):
Escenario 1: Votación secreta.
Dado que el Scrum Master inicia la votación para la HU seleccionada,
Cuando selecciono una carta de la baraja Fibonacci (ej. 8) y confirmo mi voto,
Entonces mi voto debe registrarse en el backend (Go) pero mostrarse como "Oculto" en el frontend (React) para el resto del equipo.
Escenario 2: Revelación de votos.
Dado que todos los integrantes confirmaron su voto,
Cuando el Scrum Master presiona "Revelar",
Entonces el sistema debe mostrar el valor elegido por cada uno y calcular el promedio sugerido.
### Depende de
HU-04 (integrantes que votan), HU-01 (historia a estimar).
### Prioridad
S (Should have) 
### Estado
Nueva
### Valor de Negocio
13 (Fibonacci) 
### Estimación
8 Story Points (revisar en el refinamiento: al pasar el consenso a HU-14, puede bajar)

---

## HU-03: Cálculo de Métricas del Proyecto 
### Descripción
Como Scrum Master quiero que el sistema calcule automáticamente las métricas del proyecto y de cada Sprint para poder analizar la capacidad real del equipo y la calidad del trabajo entregado.
### Criterios de aceptacion (BDD):
Escenario 1: Cálculo tras cerrar un Sprint.
Dado que un Sprint se cierra (HU-05) con historias marcadas como completadas (HU-11),
Cuando el sistema recalcula las métricas del proyecto,
Entonces debe obtener: Story Points planificados y completados, velocidad del equipo, horas estimadas (Story Points × horas por Story Point del proyecto) y reales, desviación entre esfuerzo estimado y real, porcentaje de historias completadas, y cantidad de defectos detectados y resueltos en ese Sprint.
Escenario 2: Proyecto sin Sprints cerrados.
Dado que un proyecto todavía no tiene ningún Sprint finalizado,
Cuando se solicita el cálculo de métricas,
Entonces el sistema debe devolver todos los valores en cero (o "sin datos") en vez de fallar o mostrar un error.
Escenario 3: Parámetro de horas por Story Point.
Dado que estoy editando un proyecto existente,
Cuando defino cuántas horas equivale un Story Point,
Entonces el sistema lo guarda en el proyecto y lo usa para calcular las horas estimadas. Si todavía no está definido, las horas estimadas y la desviación se informan como "sin datos".
### Depende de
HU-05 (cierre de Sprint), HU-11 (historias completadas), HU-14 (Story Points acordados), HU-07 (horas reales), HU-08 (defectos), HU-04 (edición del proyecto, donde vive el parámetro de horas por Story Point).
### Prioridad
M (Must have) - Requisito central del enunciado (punto 7).
### Estado
Nueva
### Valor de Negocio
21 (Fibonacci)
### Estimación
8 Story Points (cálculos de dominio, sin UI — se testea con TDD directo sobre `/internal/domain`). Pendiente de confirmar con el profesor: de dónde salen las "horas estimadas" (el enunciado no lo define); mientras tanto rige el parámetro del proyecto.

---

## HU-04: Creación de Proyecto y Asignación de Equipo 
### Descripción
Como Scrum Master
quiero crear un proyecto y agregar a los integrantes con sus respectivos roles
para poder tener un espacio de trabajo organizado y estructurar las responsabilidades del equipo.
### Criterios de aceptacion (BDD):
Escenario 1: Creación del entorno del proyecto.
Dado que me encuentro en la pantalla principal del sistema,
Cuando ingreso el nombre del proyecto ("Software Metrics & Estimation"), una descripción y presiono "Crear",
Entonces el sistema genera el proyecto en la base de datos SQLite y me redirige a su panel principal.
Escenario 2: Asignación de integrantes y roles.
Dado que estoy en la vista de configuración del proyecto recién creado,
Cuando ingreso el nombre de un integrante y selecciono su rol (Scrum Master o Product Builder),
Entonces el sistema vincula al usuario al proyecto, habilitando sus permisos correspondientes.
Escenario 3: Fechas del proyecto inválidas.
Dado que estoy editando un proyecto existente,
Cuando ingreso una fecha de fin igual o anterior a la fecha de inicio y guardo,
Entonces el sistema rechaza el cambio con un error de validación y no guarda.
### Depende de
Ninguna (es la base).
### Prioridad
M (Must have) - Requisito fundamental para que el sistema funcione. 
### Estado
Implementada (el Escenario 3 exige fin estrictamente posterior: ajuste pendiente respecto de lo implementado, que acepta fechas iguales).
### Valor de Negocio
21 (Fibonacci) 
### Estimación
8 Story Points (Implica crear las relaciones en base de datos entre Proyecto, Usuario y Roles). 

---

## HU-05: Apertura y Cierre de Sprints 
### Descripción
Como Scrum Master
quiero crear, iniciar y cerrar un Sprint definiendo su objetivo (Sprint Goal) y plazos
para poder organizar y delimitar las iteraciones de desarrollo del equipo.

### Criterios de aceptacion (BDD):
Escenario 1: Iniciar un Sprint.
Dado que existe un Sprint configurado con estado "Pendiente",
Cuando defino el Sprint Goal, selecciono la fecha de inicio/fin y presiono "Iniciar Sprint",
Entonces el estado cambia a "Activo" y el sistema bloquea la posibilidad de iniciar otro Sprint simultáneo en el mismo proyecto.
Escenario 2: Cierre de Sprint con arrastre de trabajo.
Dado que un Sprint se encuentra en estado "Activo",
Cuando el Scrum Master presiona "Cerrar Sprint",
Entonces el sistema cambia su estado a "Finalizado" y mueve automáticamente todas las Historias de Usuario no completadas de vuelta al Product Backlog
Escenario 3: Crear un Sprint.
Dado que existe un proyecto,
Cuando creo un Sprint para ese proyecto,
Entonces el Sprint queda "Pendiente". El sistema rechaza la creación si el proyecto no existe o si el proyecto ya tiene otro Sprint "Pendiente".
Escenario 4: Fechas del Sprint inválidas.
Dado que existe un Sprint "Pendiente",
Cuando intento iniciarlo con una fecha de fin igual o anterior a la de inicio (incluido un Sprint de un día),
Entonces el sistema rechaza la operación por validación de fechas.

### Depende de
HU-04 (el proyecto debe existir), HU-01 (historias a arrastrar al cerrar).

### Prioridad
M (Must have) - Sin esto no hay marco iterativo ágil. 
### Estado
Implementada
### Valor de Negocio
21 (Fibonacci) 

### Estimación
13 Story Points (La lógica de mover las historias no completadas al backlog requiere un manejo cuidadoso de las transacciones en SQLite). 

---

## HU-06: Movimiento de Historias al Sprint Backlog 
### Descripción
Como Product Builder
quiero seleccionar historias de usuario del Product Backlog y asignarlas a un Sprint específico
para poder definir el compromiso de trabajo de la iteración actual.

### Criterios de aceptacion (BDD):
Escenario 1: Asignación exitosa.
Dado que existe un Sprint en estado "Pendiente" o "Activo",
Cuando arrastro o selecciono una HU del Product Backlog (estado "Nueva" o "En progreso") y la asigno al Sprint,
Entonces la HU desaparece de la vista del backlog general y pasa a conformar el Sprint Backlog de esa iteración.
Escenario 2: Restricción de historias completadas.
Dado que estoy planificando un Sprint,
Cuando intento asignar una HU cuyo estado ya es "Completada",
Entonces el sistema me muestra un mensaje de error indicando que solo se pueden planificar historias pendientes.

### Depende de
HU-05 (Sprints), HU-15 (vista del backlog sin Sprint), HU-01 (historias). El estado "Completada" lo define HU-11.

### Prioridad
M (Must have) - Core del flujo de trabajo de Scrum. 
### Estado
Nueva
### Valor de Negocio
13 (Fibonacci) 

### Estimación
5 Story Points

---

## HU-07: Imputación de Horas Trabajadas 

### Descripción
Como integrante del equipo
quiero registrar las horas reales que invertí trabajando en una Historia de Usuario
para poder proveer datos al sistema que permitan comparar el esfuerzo estimado (Story Points) contra el esfuerzo real.

### Criterios de aceptacion (BDD):
Escenario 1: Registro válido de tiempo.
Dado que soy integrante del proyecto y la historia pertenece a un Sprint activo,
Cuando abro el formulario de registro, selecciono mi nombre, ingreso la fecha de hoy, la actividad realizada y "4 horas", y guardo,
Entonces el sistema suma esas 4 horas al esfuerzo total acumulado de esa HU, registrando integrante, fecha y actividad.
Escenario 2: Validación de inputs.
Dado que intento cargar horas a una historia,
Cuando ingreso un valor negativo (-2 horas), texto no numérico, una actividad vacía o un integrante que no pertenece al proyecto,
Entonces el sistema bloquea el guardado y muestra un error de validación en la interfaz de React.

### Depende de
HU-04 (integrantes), HU-06 (historia dentro de un Sprint), HU-05 (Sprint activo).

### Prioridad
M (Must have) - Fundamental para que luego se puedan calcular las métricas de desviación de esfuerzo. 
### Estado
Nueva
### Valor de Negocio
13 (Fibonacci) 

### Estimación
5 Story Points

---

## HU-08: Registro y Seguimiento de Bugs 
### Descripción
Como integrante del equipo
quiero registrar un defecto técnico indicando su severidad, la historia relacionada y el sprint de detección
para poder hacer un seguimiento de la calidad y asegurar que se resuelva antes del despliegue final.

### Criterios de aceptacion (BDD):
Escenario 1: Carga de un nuevo defecto.
Dado que se detecta un error técnico en el software y el proyecto tiene un Sprint "Activo",
Cuando completo el formulario de Defectos con la descripción, nivel de severidad (Alta, Media, Baja) y la historia relacionada, y guardo,
Entonces el Bug se registra en el sistema asociado a esa historia, con el Sprint "Activo" como Sprint de detección y estado "Pendiente".
Escenario 2: Resolución del defecto.
Dado que un desarrollador solucionó un Bug pendiente,
Cuando cambia su estado a "Resuelto",
Entonces el sistema registra internamente el Sprint de resolución para que posteriormente aparezca en el Dashboard de métricas de calidad.
Escenario 3: Sin Sprint activo.
Dado que el proyecto no tiene ningún Sprint "Activo",
Cuando intento registrar un defecto,
Entonces el sistema rechaza la operación e indica que no hay un Sprint activo al que asociarlo.

### Depende de
HU-01 (historia relacionada), HU-05 (Sprint de detección y de resolución).

### Prioridad
S (Should have) - Es muy importante para las métricas de calidad, pero el proyecto podría arrancar los primeros Sprints funcionales sin este módulo. 
### Estado
Nueva
### Valor de Negocio
8 (Fibonacci)
### Estimación
8 Story Points 

---

## HU-09: Dashboard de Métricas 
### Descripción
Como Scrum Master quiero visualizar en un panel las métricas calculadas del proyecto (HU-03) con representaciones gráficas para poder evaluar de un vistazo el estado del proyecto sin tener que leer datos crudos.
### Criterios de aceptacion (BDD):
Escenario 1: Renderizado del Dashboard.
Dado que el proyecto tiene al menos un Sprint finalizado y sus métricas calculadas,
Cuando ingreso a la sección "Dashboard",
Entonces el sistema debe renderizar un gráfico de barras de velocidad (Sprints en eje X, Puntos Comprometidos vs. Completados en eje Y) y al menos un indicador visual por cada métrica restante del punto 7 del enunciado.
Escenario 2: Proyecto sin datos.
Dado que un proyecto todavía no tiene Sprints finalizados,
Cuando ingreso al Dashboard,
Entonces el sistema debe mostrar un estado vacío informativo, no un gráfico roto ni un error.
### Depende de
HU-03 (métricas calculadas).
### Prioridad
M (Must have) - Requisito central del enunciado (punto 8).
### Estado
Nueva
### Valor de Negocio
13 (Fibonacci)
### Estimación
8 Story Points (integración de librería de gráficos en el frontend React)

---

## HU-10: Generación de Reporte en PDF 
### Descripción
Como Scrum Master quiero exportar un reporte en PDF de un proyecto o Sprint con sus historias, estimaciones, esfuerzo, métricas y defectos para poder compartir el estado del proyecto fuera del sistema (ej. con el profesor/Cliente).
### Criterios de aceptacion (BDD):
Escenario 1: Exportación exitosa.
Dado que estoy viendo el resumen de un Sprint cerrado,
Cuando presiono "Exportar a PDF",
Entonces la librería Maroto debe generar un archivo con: historias planificadas y completadas, estimaciones, esfuerzo registrado, métricas (HU-03) y defectos asociados a ese Sprint.
Escenario 2: Error de generación.
Dado que el proceso de generación del PDF falla (ej. dato faltante),
Cuando se solicita el reporte,
Entonces el sistema debe informar el error al usuario en vez de descargar un archivo corrupto o vacío.
### Depende de
HU-05 (Sprint cerrado), HU-14 (estimaciones), HU-07 (esfuerzo registrado), HU-03 (métricas), HU-08 (defectos).
### Prioridad
M (Must have) - Requisito central del enunciado (punto 9).
### Estado
Nueva
### Valor de Negocio
13 (Fibonacci)
### Estimación
8 Story Points (integración con Maroto, depende de HU-03, HU-07 y HU-08 para los datos)

---

## HU-11: Actualización de Estado de una Historia en el Sprint 
### Descripción
Como Product Builder quiero cambiar el estado de una historia asignada al Sprint (En progreso, Completada) para poder reflejar el avance real del trabajo durante la iteración.
### Criterios de aceptacion (BDD):
Escenario 1: Iniciar el trabajo de una historia.
Dado que una historia está asignada a un Sprint activo con estado "Nueva",
Cuando la paso a "En progreso",
Entonces el sistema registra el cambio de estado.
Escenario 2: Marcar historia como completada.
Dado que una historia está asignada a un Sprint activo con estado "En progreso",
Cuando la marco como "Completada",
Entonces el sistema registra la fecha de finalización y la incluye en el cálculo de Story Points completados del Sprint (HU-03).
Escenario 3: Transición inválida.
Dado que una historia todavía no fue asignada a ningún Sprint,
Cuando intento marcarla como "En progreso" o "Completada" directamente desde el Product Backlog,
Entonces el sistema rechaza el cambio y muestra un mensaje indicando que primero debe asignarse a un Sprint activo.
### Depende de
HU-06 (historia asignada a un Sprint), HU-05 (Sprint activo).
### Prioridad
M (Must have) - Sin esto no hay datos reales para ninguna métrica.
### Estado
Nueva
### Valor de Negocio
13 (Fibonacci)
### Estimación
3 Story Points (revisar en el refinamiento: ahora cubre dos transiciones)

---

## HU-12: Consulta de Sprints Anteriores 
### Descripción
Como Scrum Master quiero consultar el historial de Sprints ya cerrados con su resumen (Sprint Goal, historias completadas, métricas) para poder comparar el desempeño del equipo a lo largo del proyecto.
### Criterios de aceptacion (BDD):
Escenario 1: Listado de Sprints cerrados.
Dado que el proyecto tiene al menos dos Sprints finalizados,
Cuando ingreso a la sección "Historial de Sprints",
Entonces el sistema muestra una lista con cada Sprint cerrado, su Sprint Goal, fechas y métricas principales.
Escenario 2: Detalle de un Sprint puntual.
Dado que estoy en el historial de Sprints,
Cuando selecciono uno en particular,
Entonces el sistema muestra el detalle completo de ese Sprint (historias, esfuerzo, defectos) sin afectar al Sprint activo actual.
### Depende de
HU-05 (Sprints cerrados), HU-03 (métricas), HU-07 (esfuerzo), HU-08 (defectos), HU-11 (historias completadas).
### Prioridad
S (Should have) - Requisito explícito del enunciado (punto 3).
### Estado
Nueva
### Valor de Negocio
8 (Fibonacci)
### Estimación
3 Story Points (una vez que existen las historias de las que depende)

---

## HU-13: Fechas y Estado del Proyecto 
### Descripción
Como Scrum Master quiero registrar la fecha de inicio y finalización de un proyecto y consultar su estado general para poder tener visibilidad del ciclo de vida completo del proyecto. El registro/edición de fechas lo cubre HU-04; esta historia implementa la consulta de estado.
### Criterios de aceptacion (BDD):
Escenario 1: Registro de fechas (cubierto por HU-04).
Dado que estoy editando un proyecto existente (HU-04),
Cuando ingreso una fecha de inicio y una fecha de finalización estimada y guardo,
Entonces el sistema valida que la fecha de fin sea posterior a la de inicio y las persiste.
Escenario 2: Consulta de estado.
Dado que un proyecto tiene Sprints y fechas cargadas,
Cuando consulto su estado,
Entonces el sistema muestra si está "Planificado", "En curso" o "Finalizado" según la fecha actual y el estado de sus Sprints.
Escenario 3: Los Sprints mandan sobre las fechas.
Dado que un proyecto tiene un Sprint "Activo" y su fecha de fin ya pasó,
Cuando consulto su estado,
Entonces el sistema muestra "En curso".
Escenario 4: Proyecto sin Sprints ni fechas.
Dado que un proyecto no tiene Sprints ni fechas cargadas,
Cuando consulto su estado,
Entonces el sistema muestra "Planificado".
### Depende de
HU-04 (proyecto y fechas), HU-05 (estado de los Sprints).
### Prioridad
S (Should have) - Extiende HU-04, no bloquea el resto del sistema.
### Estado
Implementada
### Valor de Negocio
8 (Fibonacci)
### Estimación
3 Story Points

---

## HU-14: Consenso y Registro de la Estimación Acordada 
### Descripción
Como Scrum Master quiero que, al revelar los votos de Planning Poker, el sistema detecte las diferencias entre las estimaciones, permita hacer nuevas rondas y registre la estimación acordada en la historia para poder obtener los Story Points sobre los que se calculan las métricas.
### Criterios de aceptacion (BDD):
Escenario 1: Detección de diferencias.
Dado que se revelaron los votos de una historia,
Cuando las estimaciones no coinciden,
Entonces el sistema indica que hay diferencias y cuáles son los valores mínimo y máximo votados.
Escenario 2: Nueva ronda.
Dado que se detectaron diferencias entre las estimaciones,
Cuando el Scrum Master inicia una nueva ronda,
Entonces el sistema descarta los votos de la ronda anterior (que quedan registrados como historial) y permite votar de nuevo.
Escenario 3: Registro de la estimación acordada.
Dado que el equipo llegó a un acuerdo,
Cuando el Scrum Master registra el valor acordado,
Entonces la historia queda con esos Story Points como estimación.
Escenario 4: Estimación inválida.
Dado que el Scrum Master intenta registrar una estimación,
Cuando el valor no pertenece a la baraja Fibonacci o la votación no fue revelada,
Entonces el sistema rechaza la operación y la historia mantiene su estimación anterior.
### Depende de
HU-02 (votos y revelación), HU-01 (campo de Story Points de la historia).
### Prioridad
S (Should have) - Completa el requisito de Planning Poker del enunciado (punto 4).
### Estado
Nueva
### Valor de Negocio
8 (Fibonacci)
### Estimación
5 Story Points

---

## HU-15: Consulta y Edición del Product Backlog 
### Descripción
Como Product Builder quiero ver el Product Backlog de un proyecto ordenado por prioridad y poder editar el título, la descripción y los criterios de aceptación de cada historia para poder mantenerlo actualizado y cumplir con el contenido mínimo que pide el enunciado.
### Criterios de aceptacion (BDD):
Escenario 1: Listado del Product Backlog.
Dado que un proyecto tiene historias, algunas ya asignadas a un Sprint,
Cuando consulto su Product Backlog,
Entonces el sistema muestra solo las historias sin Sprint asignado, ordenadas por prioridad (MoSCoW) y, dentro de la misma prioridad, por orden de creación.
Escenario 2: Edición de una historia.
Dado que existe una historia en el Product Backlog,
Cuando modifico su título, descripción o criterios de aceptación con datos válidos y guardo,
Entonces el sistema persiste los cambios y conserva su identificador, estado y estimación.
Escenario 3: Edición inválida.
Dado que estoy editando una historia,
Cuando dejo el título en blanco,
Entonces el sistema rechaza el cambio, muestra una advertencia y no modifica la historia.
Escenario 4: Historia inexistente.
Dado que intento consultar o editar una historia o un proyecto que no existen,
Cuando envío la petición,
Entonces el sistema responde "no encontrado".
### Depende de
HU-01 (historias), HU-04 (proyecto), HU-05 (para saber qué historias tienen Sprint asignado).
### Prioridad
M (Must have) - El enunciado exige los criterios de aceptación como parte mínima de cada elemento del Product Backlog, y HU-06 necesita el listado.
### Estado
Nueva
### Valor de Negocio
13 (Fibonacci)
### Estimación
5 Story Points
