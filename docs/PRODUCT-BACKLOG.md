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
>   único vocabulario: no usar "Terminado" ni "Done". `En progreso` solo existe dentro
>   de un Sprint: si una historia no completada vuelve al Product Backlog al cerrar su
>   Sprint, vuelve a `Nueva`.
> - La fecha de fin debe ser estrictamente posterior a la de inicio, tanto en
>   proyectos como en Sprints.
> - Los Story Points planificados de un Sprint son los de las historias **asignadas al
>   cerrarlo** (incluidas las agregadas con el Sprint ya `Activo`): cuáles quedaron sin
>   completar se registra al cerrar, porque esas historias vuelven al Product Backlog.
> - No se modifican los Story Points ya cargados de una historia asignada a un Sprint
>   `Activo` o `Finalizado`, para que las métricas históricas se puedan reproducir. Cargarlos
>   por primera vez (la historia estaba sin estimar) sí se permite.
> - Los estados derivados (como el estado del proyecto) y las métricas se calculan al
>   consultar; no se guardan.
> - El enunciado habla de registrar esfuerzo "para cada historia o tarea": en este
>   proyecto no existe la entidad Tarea y el esfuerzo (HU-07) se registra por historia.
> - Todo requisito mínimo del enunciado tiene una historia dueña. Las historias
>   14 y 15 se agregaron para cubrir requisitos que ninguna historia tenía.

## ***Prioridad 1: Cimientos del Sistema (Dependencia Crítica)***

No se puede gestionar nada si no existe un contenedor. Esta es la base estructural.

* **HU-04: Creación de Proyecto y Asignación de Equipo (Must have).** Es la máxima prioridad operativa. Todo el sistema (historias, sprints, horas) depende de que exista un proyecto y usuarios registrados.  
* **HU-01: Creación de Historias de Usuario (Must have).** Sin esto, no hay Product Backlog.

## ***Prioridad 2: El Motor de Scrum (Iteraciones)***

Una vez que tienen el proyecto y el Backlog, necesitan el ciclo de trabajo ágil.

* **HU-05: Apertura y Cierre de Sprints (Must have).** Da vida a la metodología.  
* **HU-15: Consulta y Edición del Product Backlog (Must have).** Completa el punto 2 del enunciado: listar el backlog y mantener los criterios de aceptación y los Story Points de cada historia. HU-06 necesita este listado.
* **HU-06: Movimiento de Historias al Sprint Backlog (Must have).** Permite a las Product Builders definir en qué se va a trabajar durante la iteración.
* **HU-11: Actualización de Estado de una Historia en el Sprint (Must have).** Sin poder marcar una historia como completada durante el sprint, no hay datos reales para calcular velocidad ni ninguna métrica.
* **HU-13: Fechas y Estado del Proyecto (Should have).** Extiende HU-04: consultar el estado general del proyecto (Planificado / En curso / Finalizado) a partir de sus fechas y Sprints; el registro de las fechas de inicio/fin lo cubre HU-04 (requisito explícito del punto 1 del enunciado). Va en esta prioridad porque usa el estado de los Sprints (HU-05).

## ***Prioridad 3: Alimentación de Datos (El Valor de Negocio)***

Para que el diferenciador del sistema (las métricas) funcione, primero necesitan generar y registrar datos reales.

* **HU-02: Votación de Esfuerzo en Planning Poker (Should have).** Aporta un gran valor metodológico: registrar las estimaciones individuales, mantenerlas ocultas y mostrarlas.  
* **HU-14: Consenso y Registro de la Estimación Acordada (Should have).** Completa el Planning Poker del punto 4 del enunciado: detectar diferencias, nuevas rondas y registrar la estimación acordada, que son los "Story Points" estimados.
* **HU-07: Imputación de Horas Trabajadas (Must have).** Requisito indispensable para luego poder cruzar lo estimado vs. lo real.
* **HU-08: Registro y Seguimiento de Bugs (Must have).** Los defectos son requisito mínimo del enunciado (punto 6) y alimentan las métricas (HU-03) y el reporte (HU-10), por eso se carga junto con el resto de los datos.

## ***Prioridad 4: Resultados y Calidad (El Valor Final)***

Estas funcionalidades son el corazón del valor para el cliente final, pero tienen una alta dependencia técnica de las etapas anteriores. Se dividió la antigua HU-03 en tres historias más chicas (cálculo, visualización y exportación son responsabilidades distintas) y se agregó cobertura completa del punto 7 del enunciado.

* **HU-03: Cálculo de Métricas del Proyecto (Must have).** El motor de cálculo (Story Points planificados/completados, velocidad, horas estimadas/reales, desviación, % completadas, defectos detectados/resueltos) — requisito central del enunciado (punto 7).
* **HU-09: Dashboard de Métricas (Must have).** Visualización gráfica de lo que calcula HU-03 (punto 8 del enunciado).
* **HU-10: Generación de Reporte en PDF (Must have).** Exportación a PDF con Maroto, de un Sprint o del proyecto completo, con historias, estimaciones, esfuerzo, métricas y defectos (punto 9 del enunciado; el Sprint 4 pide el informe del proyecto).
* **HU-12: Consulta de Sprints Anteriores (Should have).** Requisito explícito del punto 3 del enunciado. Va al final porque su historial muestra métricas (HU-03), esfuerzo (HU-07) y defectos (HU-08), que tienen que existir antes.

## ***Prioridad 5: Mejoras de calidad del beta testing (Could have / Should have)***

Surgieron al probar la aplicación desplegada en Render (2026-10-30). No son requisitos del enunciado. HU-16, HU-17 y HU-18 se estiman con el equipo y entran solo si sobra tiempo, siempre después de las historias de las prioridades anteriores. HU-20 (Sprint 4) y HU-21 (Sprint 5), junto con las mejoras #235 y #236 (Sprint 4), figuran en el plan de cada Sprint por decisión del Scrum Master (2026-11-03): no dependen de que sobre tiempo. HU-19 quedó absorbida por HU-20 y HU-21. Prioridad, valor y estimación de estas seis historias son una propuesta del Scrum Master, a confirmar con el equipo.

* **HU-16: Gestión del Equipo del Proyecto (Could have).** Agregar a un usuario ya registrado como Product Builder y quitar integrantes sin perder lo que aportaron; el único Scrum Master no se quita ni cambia de rol. Extiende HU-04.
* **HU-17: Gestión de Usuarios (Could have).** Crear y editar usuarios y elegir a los integrantes de una lista, en vez de escribir nombres a mano; incluye los usuarios predeterminados de la demo. Extiende HU-04.
* **HU-18: Avance de las Historias del Proyecto (Should have).** Ver todas las historias agrupadas por estado y por Sprint, para saber qué se terminó y qué falta. Extiende HU-15 y HU-11.
* **HU-19: Administración de Mis Proyectos (Could have).** Absorbida por HU-20 y HU-21 (decisión del Scrum Master, 2026-11-03). Su nivel 1, la lista y el selector de proyectos, es #148; los niveles 2 y 3 (mis proyectos, permisos por rol y autenticación) pasan a HU-20 y HU-21. Se descarta la transferencia de la propiedad: el dueño del proyecto es su Scrum Master. Extiende HU-04.
* **HU-20: Autenticación de usuarios (Should have).** Registrarse e iniciar sesión con un nombre de usuario y una contraseña, para que el sistema sepa quién usa la aplicación. Extiende HU-04.
* **HU-21: Permisos por rol (Should have).** Que las acciones de gestión del proyecto (Sprints, votación, edición y equipo) solo las pueda hacer el Scrum Master, con la identidad tomada de la sesión, y que cada persona vea solo sus proyectos. Depende de HU-20.

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
Implementada (el Escenario 3, de integridad referencial, se resolvió con el fix #19, PR #103; el frontend se entregó con el PR #125).

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
### Notas
El modelo de votación debe contemplar rondas desde el inicio (cada voto pertenece a una ronda), porque la historia de consenso (HU-14) agrega nuevas rondas y conserva el historial.
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
HU-05 (cierre de Sprint), HU-11 (historias completadas), HU-15 (Story Points de cada historia), HU-07 (horas reales), HU-08 (defectos), HU-04 (edición del proyecto, donde vive el parámetro de horas por Story Point).
### Notas
Los Story Points planificados de un Sprint son los de las historias asignadas al cerrarlo (incluidas las que quedaron sin completar y volvieron al Product Backlog; ver HU-05). Como los Story Points de una historia de un Sprint `Activo` o `Finalizado` no se editan (HU-15), las métricas de Sprints pasados no cambian.
### Prioridad
M (Must have) - Requisito central del enunciado (punto 7).
### Estado
Nueva
### Valor de Negocio
21 (Fibonacci)
### Estimación
8 Story Points (cálculos de dominio, sin UI — se testea con TDD directo sobre `/internal/domain`). **Decidido por el equipo (2026-10-28):** el enunciado no define de dónde salen las "horas estimadas" y el equipo votó la opción A: horas estimadas = Story Points × "horas por Story Point" fijo del proyecto. Si el profesor aclara otra fuente, se ajusta con un fix.

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
Entonces el sistema vincula al usuario al proyecto y registra su rol (la aplicación de permisos por rol queda fuera del alcance por ahora).
Escenario 3: Fechas del proyecto inválidas.
Dado que estoy editando un proyecto existente,
Cuando ingreso una fecha de fin igual o anterior a la fecha de inicio y guardo,
Entonces el sistema rechaza el cambio con un error de validación y no guarda.
### Depende de
Ninguna (es la base).
### Prioridad
M (Must have) - Requisito fundamental para que el sistema funcione. 
### Estado
Implementada (el ajuste de fecha de fin estrictamente posterior se resolvió con el fix #100, PR #111; GET /projects con #126 y el frontend con el PR #121).
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
Entonces el sistema cambia su estado a "Finalizado" y mueve automáticamente todas las Historias de Usuario no completadas de vuelta al Product Backlog, sin Sprint asignado y con estado "Nueva", y deja registrado cuáles eran y sus Story Points de ese momento.
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

### Notas
Un Sprint puede iniciarse con fechas de inicio y de fin anteriores a hoy (el equipo pudo empezarlo y olvidar declararlo); no se restringe (#20). Solo se exige que el fin sea posterior al inicio.
Al cerrar, el sistema registra qué historias quedaron sin completar y sus Story Points de ese momento; sin ese registro no se podrían calcular los Story Points planificados (HU-03), el gráfico de HU-09 ni el detalle de HU-12 de un Sprint ya cerrado. El cierre es una sola operación: o se aplica todo (estado del Sprint, registro y regreso al Product Backlog) o nada.
### Prioridad
M (Must have) - Sin esto no hay marco iterativo ágil. 
### Estado
Implementada (las historias arrastradas y el cierre atómico se resolvieron con el issue #117, PR #136; el frontend con el PR #127; falta la consulta de Sprints, #110, PR #149 abierto)
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
Cuando arrastro o selecciono una HU del Product Backlog (estado "Nueva") y la asigno al Sprint,
Entonces la HU desaparece de la vista del backlog general y pasa a conformar el Sprint Backlog de esa iteración.
Escenario 2: Restricción de historias completadas.
Dado que estoy planificando un Sprint,
Cuando intento asignar una HU cuyo estado ya es "Completada",
Entonces el sistema me muestra un mensaje de error indicando que solo se pueden planificar historias pendientes.

### Depende de
HU-05 (Sprints y estado "Completada", que ya existe), HU-15 (vista del backlog sin Sprint), HU-01 (historias).

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
para poder proveer datos al sistema que permitan comparar el esfuerzo estimado (Story Points convertidos a horas, ver HU-03) contra el esfuerzo real.

### Criterios de aceptacion (BDD):
Escenario 1: Registro válido de tiempo.
Dado que soy integrante del proyecto y la historia pertenece a un Sprint activo,
Cuando abro el formulario de registro, selecciono mi nombre, ingreso la fecha de hoy, la actividad realizada y "4 horas", y guardo,
Entonces el sistema suma esas 4 horas al esfuerzo total acumulado de esa HU, registrando integrante, fecha, actividad y el Sprint activo en ese momento (así las horas se atribuyen al Sprint correcto aunque la historia pase luego a otro).
Escenario 2: Validación de inputs.
Dado que intento cargar horas a una historia,
Cuando ingreso un valor negativo (-2 horas), texto no numérico, una actividad vacía o un integrante que no pertenece al proyecto,
Entonces el sistema bloquea el guardado y muestra un error de validación en la interfaz de React.

### Depende de
HU-04 (integrantes), HU-06 (historia dentro de un Sprint), HU-05 (Sprint activo).

### Notas
Las horas se guardan como minutos enteros y se convierten a horas solo al mostrarlas. El tope de horas por carga lo define la especificación de esta historia.
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
Cuando intento registrar un defecto o marcar uno como "Resuelto",
Entonces el sistema rechaza la operación e indica que no hay un Sprint activo al que asociarlo (como Sprint de detección o de resolución).

### Depende de
HU-01 (historia relacionada), HU-05 (Sprint de detección y de resolución).

### Prioridad
M (Must have) - Los defectos son requisito mínimo del enunciado (punto 6) y HU-03 y HU-10 dependen de ellos.
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
### Notas
"Puntos Comprometidos" son los Story Points de las historias asignadas al cerrar el Sprint, completadas o no (ver HU-05).
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
Escenario 1: Exportación exitosa de un Sprint.
Dado que estoy viendo el resumen de un Sprint finalizado,
Cuando presiono "Exportar a PDF",
Entonces la librería Maroto debe generar un archivo con: historias planificadas y completadas, estimaciones, esfuerzo registrado, métricas (HU-03) y defectos asociados a ese Sprint.
Escenario 2: Exportación exitosa del proyecto.
Dado que estoy viendo el resumen de un proyecto,
Cuando presiono "Exportar reporte del proyecto",
Entonces la librería Maroto debe generar un archivo con: los datos del proyecto (nombre, fechas, estado e integrantes), su Product Backlog actual (las historias sin Sprint, con prioridad y estimación), el listado de sus Sprints finalizados con Sprint Goal y fechas, el Sprint Activo con sus historias si lo hay y, para los Sprints finalizados, historias planificadas y completadas, estimaciones, esfuerzo registrado, métricas (HU-03) y defectos detectados y resueltos.
Escenario 3: Proyecto sin Sprints finalizados.
Dado que un proyecto todavía no tiene ningún Sprint finalizado (con o sin un Sprint Activo),
Cuando solicito el reporte del proyecto,
Entonces el sistema genera el reporte con los datos del proyecto, su Product Backlog y, si lo hay, el Sprint Activo con sus historias, e informa las métricas, el esfuerzo y los defectos como "sin datos" en vez de fallar (igual que HU-03 y HU-09).
Escenario 4: Error de generación.
Dado que el proceso de generación del PDF falla por un error real (por ejemplo, Maroto no puede armar el documento, falla la escritura del archivo o la lectura de la base),
Cuando se solicita el reporte,
Entonces el sistema debe informar el error al usuario en vez de descargar un archivo corrupto o vacío.
### Depende de
HU-04 (datos del proyecto e integrantes), HU-13 (estado del proyecto), HU-05 (Sprint finalizado), HU-06 (historias del Sprint Activo), HU-15 (estimaciones), HU-07 (esfuerzo registrado), HU-03 (métricas), HU-08 (defectos).
### Notas
Las historias "planificadas" son las asignadas al cerrar el Sprint, completadas o no (ver HU-05).
El reporte del proyecto incluye siempre el Product Backlog actual (decisión del Scrum Master, 2026-11-01). Las métricas, el esfuerzo y los defectos se agregan solo con los Sprints finalizados, igual que las métricas de HU-03; un Sprint Activo se lista con sus historias pero no suma a las métricas. El enunciado pide "un reporte de un proyecto o Sprint" y, en el Sprint 4, "exportar un informe del proyecto a un archivo PDF": por eso la historia cubre las dos variantes (decisión del equipo, 2026-10-28).
### Prioridad
M (Must have) - Requisito central del enunciado (punto 9).
### Estado
Nueva
### Valor de Negocio
13 (Fibonacci)
### Estimación
8 Story Points (integración con Maroto, depende de HU-03, HU-07 y HU-08 para los datos). Revisar en el refinamiento del Sprint 4: ahora cubre el reporte de un Sprint y el del proyecto, que comparten el armado del documento.

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
Entonces el sistema registra la fecha de finalización; la historia queda disponible para el cálculo de Story Points completados del Sprint (HU-03).
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
### Notas
Las historias que se muestran como comprometidas son las asignadas al cerrar el Sprint, completadas o no (ver HU-05).
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
Implementada (el frontend se entregó con el PR #150, issue #41; queda el fix #151, PR #152 abierto)
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
Dado que se revelaron los votos de una historia, haya o no diferencias entre las estimaciones,
Cuando el Scrum Master inicia una nueva ronda,
Entonces el sistema descarta los votos de la ronda anterior (que quedan registrados como historial, que se puede consultar) y permite votar de nuevo.
Escenario 3: Registro de la estimación acordada.
Dado que el equipo llegó a un acuerdo,
Cuando el Scrum Master registra el valor acordado,
Entonces la historia queda con esos Story Points como estimación, reemplazando la que tuviera cargada (por ejemplo, de forma directa desde HU-15).
Escenario 4: Estimación inválida.
Dado que el Scrum Master intenta registrar una estimación,
Cuando el valor no pertenece a la baraja Fibonacci, la votación no fue revelada o la historia ya tiene estimación y está asignada a un Sprint "Activo" o "Finalizado",
Entonces el sistema rechaza la operación y la historia mantiene su estimación anterior.
### Depende de
HU-02 (votos y revelación), HU-01 (campo de Story Points de la historia).
### Notas
No bloquea al cálculo de métricas ni al reporte: los Story Points también se cargan de forma directa al editar el Product Backlog (HU-15).
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
Como Product Builder quiero ver el Product Backlog de un proyecto ordenado por prioridad y poder editar el título, la descripción, los criterios de aceptación y los Story Points de cada historia para poder mantenerlo actualizado y cumplir con el contenido mínimo que pide el enunciado (puntos 2 y 4: estimar historias mediante Story Points).
### Criterios de aceptacion (BDD):
Escenario 1: Listado del Product Backlog.
Dado que un proyecto tiene historias, algunas ya asignadas a un Sprint,
Cuando consulto su Product Backlog,
Entonces el sistema muestra solo las historias sin Sprint asignado, ordenadas por prioridad (MoSCoW) y, dentro de la misma prioridad, por orden de creación.
Escenario 2: Edición de una historia.
Dado que existe una historia en el Product Backlog,
Cuando modifico su título, descripción, criterios de aceptación o Story Points (valor de la escala Fibonacci) con datos válidos y guardo,
Entonces el sistema persiste los cambios y conserva su identificador y estado.
Escenario 3: Edición inválida.
Dado que estoy editando una historia,
Cuando dejo el título en blanco, ingreso unos Story Points que no pertenecen a la escala Fibonacci o intento cambiar los Story Points ya cargados de una historia asignada a un Sprint "Activo" o "Finalizado",
Entonces el sistema rechaza el cambio, muestra una advertencia y no modifica la historia.
Escenario 4: Historia inexistente.
Dado que intento consultar o editar una historia o un proyecto que no existen,
Cuando envío la petición,
Entonces el sistema responde "no encontrado".
### Depende de
HU-01 (historias), HU-04 (proyecto), HU-05 (campo de Sprint de la historia, que ya existe).
### Notas
La asignación de historias a un Sprint llega después, con la historia de movimiento al Sprint Backlog; hasta entonces el listado se prueba sembrando historias con Sprint.
Los criterios de aceptación son un texto libre opcional de hasta 2000 caracteres (el mismo límite que la descripción); vacío o solo espacios equivale a no tener criterios.
### Prioridad
M (Must have) - El enunciado exige los criterios de aceptación y los Story Points como parte mínima de cada elemento del Product Backlog, HU-06 necesita el listado y HU-03 y HU-10 necesitan los Story Points.
### Estado
Implementada (backend, PR #135; falta el frontend: issue #98)
### Valor de Negocio
13 (Fibonacci)
### Estimación
5 Story Points

---

## HU-16: Gestión del Equipo del Proyecto
### Descripción
Como Scrum Master quiero agregar a un usuario ya registrado como Product Builder y quitar a un integrante del equipo de un proyecto, conservando todo lo que aportó, para poder mantener el equipo al día sin perder el historial del proyecto.
### Criterios de aceptacion (BDD):
Escenario 1: Agregar a un usuario registrado como Product Builder.
Dado que soy el Scrum Master del proyecto y existe un usuario registrado que no integra el equipo,
Cuando lo agrego al proyecto,
Entonces queda como Product Builder y el sistema lo muestra en el equipo del proyecto, sin ascenso a Scrum Master.
Escenario 2: Quitar un integrante conservando su historial.
Dado que un integrante ya votó una estimación e imputó horas,
Cuando lo quito del equipo,
Entonces deja de figurar en el equipo y no cuenta en las votaciones nuevas, y sus votos y sus horas siguen en el historial y en las métricas.
Escenario 3: Volver a sumar a un integrante dado de baja.
Dado que un integrante fue quitado del equipo,
Cuando lo agrego de nuevo al proyecto como Product Builder,
Entonces el sistema lo reactiva con ese rol, sin duplicarlo ni perder su historial.
Escenario 4: Siempre hay un Scrum Master (único por proyecto).
Dado que un proyecto tiene un único Scrum Master activo,
Cuando intento quitarlo o cambiarle el rol,
Entonces el sistema rechaza el cambio con un mensaje claro y no modifica el equipo.
Escenario 5: Integrante o proyecto inexistente.
Dado que intento modificar un integrante o un proyecto que no existen,
Cuando envío la petición,
Entonces el sistema responde "no encontrado".
### Depende de
HU-04 (proyecto e integrantes). Afecta a HU-02 (la votación cuenta solo a los integrantes activos), a HU-07 y a HU-03 (el historial de quien fue dado de baja se conserva).
### Notas
La baja es lógica: requiere una columna nueva en `integrantes` (por ejemplo, la fecha de baja, vacía mientras el integrante está activo). Un cambio de esquema se hace en un solo PR que modifica `MODELO-DE-DATOS.md` y agrega la migración. Propuesta del Scrum Master del 2026-10-30, a confirmar con el equipo.
Al crear un proyecto siempre se registra a su creador como Scrum Master (HU-04; con HU-21 es quien tiene la sesión iniciada), así que un proyecto nunca nace sin uno y la regla se cumple desde el principio. Con el Scrum Master único (#236, HU-21) la decisión D4 de `MODELO-DE-DATOS.md` queda cerrada en esos issues: esta historia ya no la reabre.
Con un único Scrum Master por proyecto (#236), no hay ascenso a Scrum Master y no se puede quitar al Scrum Master ni cambiar su rol: el traspaso del rol queda fuera de alcance hasta definirlo.
La regla del último Scrum Master se verifica dentro de la misma transacción que el cambio, para que dos pedidos simultáneos no la salteen (la misma carrera que los fixes #145 y #158).
El nombre de la persona no se edita desde el equipo del proyecto: pertenece al usuario y se edita en HU-17.
Los permisos por rol quedaron fuera del alcance de HU-04 y los cubre HU-21.
### Prioridad
C (Could have) - Mejora de calidad surgida del beta testing; no es requisito del enunciado.
### Estado
Nueva
### Valor de Negocio
5 (Fibonacci) - propuesta, a confirmar con el equipo
### Estimación
5 Story Points - propuesta, a estimar con el equipo

---

## HU-17: Gestión de Usuarios
### Descripción
Como Scrum Master quiero crear y editar usuarios y elegir a los integrantes de un proyecto de entre los usuarios ya creados, para no asignar roles a personas ficticias ni por un error de tipeo, y para que el equipo de pruebas tenga usuarios listos después de cada despliegue.
### Criterios de aceptacion (BDD):
Escenario 1: Alta de un usuario.
Dado que ingreso un nombre válido que no existe,
Cuando creo el usuario,
Entonces el sistema lo agrega a la lista de usuarios, sin rol (el rol se asigna por proyecto).
Escenario 2: Usuario repetido.
Dado que ya existe un usuario con ese nombre, sin distinguir mayúsculas ni espacios externos,
Cuando intento crearlo otra vez,
Entonces el sistema rechaza el alta con un mensaje claro y no crea un duplicado.
Escenario 3: Elegir integrantes de la lista.
Dado que existen usuarios creados,
Cuando asigno un integrante a un proyecto,
Entonces elijo a la persona de la lista y le asigno un rol, y no puedo escribir un nombre libre.
Escenario 4: Elegir al Scrum Master al crear el proyecto.
Dado que existen usuarios creados,
Cuando creo un proyecto,
Entonces elijo de la lista a quien será su Scrum Master.
Escenario 5: Edición del nombre.
Dado que existe un usuario que participa en varios proyectos,
Cuando cambio su nombre a uno que no está en uso,
Entonces el nombre cambia en todos los proyectos y el historial (votos, horas, equipos) sigue asociado a la misma persona.
Escenario 6: Edición con un nombre en uso.
Dado que existe otro usuario con el nombre que ingreso,
Cuando guardo el cambio,
Entonces el sistema lo rechaza y no modifica al usuario.
Escenario 7: Usuarios predeterminados de la demo.
Dado que uso el botón de datos de prueba,
Cuando se cargan los datos,
Entonces se crean los usuarios predeterminados sin rol, y si ya existen no se duplican.
### Depende de
HU-04 (usuarios e integrantes) y el issue #144 (botón de datos de prueba).
### Notas
No requiere cambio de esquema: `usuarios` ya existe y es independiente de los proyectos (ver `MODELO-DE-DATOS.md` §4.2). La identidad es `nombre_normalizado`, y las demás tablas apuntan por `usuario_id`, por eso renombrar es seguro.
El contrato de asignación cambia: hoy `POST /projects/{id}/members` recibe un nombre. Se hace en dos pasos para no romper `main`: primero el backend acepta `usuario_id` además del nombre y agrega `GET /users`; después el frontend pasa a elegir de la lista y, en ese mismo PR, se retira el camino del nombre. Lo mismo con el Scrum Master de `POST /projects`.
Quién puede crear usuarios mientras no existan perfiles: cualquiera. La restricción por perfil ("admin" o "betaTester") queda para cuando se agreguen perfiles, y se discute con el equipo.
**Alineación pendiente con HU-20 y HU-21** (decisión de producto, sin resolver; HU-17 queda como está hasta que el Scrum Master decida):
1. Escenario 4: elige de una lista a quien será el Scrum Master del proyecto, pero en HU-21 el Scrum Master es quien crea el proyecto (con la sesión iniciada).
2. Escenario 3: asigna a un integrante cualquier rol, pero en HU-21 "configurar equipo" solo agrega usuarios ya registrados como Product Builder.
3. Escenario 1: crea usuarios sin contraseña, pero con HU-20 cada persona se registra sola con la suya.
4. Si el nombre de usuario de HU-20 es `usuarios.nombre`, renombrar un usuario (escenario 5) cambia su dato de ingreso.
<!-- decidir: qué queda de HU-17 con HU-20 y HU-21 (se descarta, se reduce o se redefine) -->
### Prioridad
C (Could have) - Mejora de calidad surgida del beta testing; no es requisito del enunciado.
### Estado
Nueva
### Valor de Negocio
8 (Fibonacci) - propuesta, a confirmar con el equipo
### Estimación
8 Story Points - propuesta, a estimar con el equipo

---

## HU-18: Avance de las Historias del Proyecto
### Descripción
Como integrante del equipo quiero ver todas las historias del proyecto agrupadas por estado y por Sprint, para saber qué se terminó y qué falta sin revisar los Sprints uno por uno.
### Criterios de aceptacion (BDD):
Escenario 1: Vista de avance.
Dado que un proyecto tiene historias en distintos estados y Sprints,
Cuando abro la pantalla "Estado" del proyecto,
Entonces el sistema las muestra agrupadas en Nueva, En progreso y Completada, con la cantidad de historias y de Story Points de cada grupo.
Escenario 2: Sprint de cada historia.
Dado que una historia pertenece a un Sprint,
Cuando miro su grupo,
Entonces veo el Sprint al que pertenece; si no tiene Sprint, se indica que está en el Product Backlog.
Escenario 3: Proyecto sin historias.
Dado que un proyecto no tiene historias,
Cuando abro la pantalla "Estado",
Entonces el sistema muestra un mensaje de que no hay historias y los totales en cero.
Escenario 4: Proyecto inexistente.
Dado que el proyecto no existe,
Cuando abro la pantalla "Estado",
Entonces el sistema responde "no encontrado".
### Depende de
HU-01 (historias), HU-05 y HU-06 (Sprints y su backlog), HU-11 (estado de la historia), HU-15 (Story Points).
### Notas
Es una vista de solo lectura: no recalcula métricas (eso es HU-03) ni cambia el Product Backlog de HU-15, que sigue mostrando solo las historias sin Sprint.
Reemplaza el texto "esta pantalla será implementada" de la pestaña Estado (`/projects/{id}/status`). El panel del proyecto conserva el estado Planificado / En curso / Finalizado de HU-13.
### Prioridad
S (Should have) - Responde a la pregunta básica de seguimiento ("¿qué está terminado y qué no?") que hoy no se puede contestar desde la pantalla.
### Estado
Nueva
### Valor de Negocio
8 (Fibonacci) - propuesta, a confirmar con el equipo
### Estimación
5 Story Points - propuesta, a estimar con el equipo

---

## HU-19: Administración de Mis Proyectos
### Descripción
Como usuario quiero ver los proyectos que creé y los que tengo asignados, y administrarlos según mi rol (dueño, Scrum Master o Product Builder), para trabajar solo con lo que me corresponde y poder transferir la propiedad de un proyecto a otra persona.
### Criterios de aceptacion (BDD):
Escenario 1: Mis proyectos.
Dado que soy un usuario con proyectos creados por mí y proyectos en los que soy integrante,
Cuando abro "Administrar proyectos",
Entonces el sistema lista cada proyecto con mi rol en él (dueño, Scrum Master o Product Builder).
Escenario 2: El creador es el dueño.
Dado que creo un proyecto,
Cuando se crea,
Entonces quedo registrado como su dueño.
Escenario 3: Acciones según el rol.
Dado que abro un proyecto de mi lista,
Cuando miro las acciones disponibles,
Entonces el dueño puede editar el proyecto, transferirlo y cerrarlo; el Scrum Master administra el equipo y los Sprints; y el Product Builder trabaja con las historias, sin ver las acciones que no le corresponden.
Escenario 4: Transferir la propiedad.
Dado que soy el dueño de un proyecto,
Cuando lo transfiero a otro integrante del proyecto,
Entonces esa persona pasa a ser la dueña y yo conservo mi rol de integrante.
Escenario 5: Transferencia inválida.
Dado que soy el dueño de un proyecto,
Cuando intento transferirlo a un usuario que no existe, que no es integrante o que soy yo mismo,
Entonces el sistema rechaza la transferencia y el dueño no cambia.
Escenario 6: Usuario sin proyectos.
Dado que no creé ni tengo asignado ningún proyecto,
Cuando abro "Administrar proyectos",
Entonces el sistema muestra un mensaje de que no hay proyectos y me permite crear uno.
### Depende de
HU-04 (proyectos e integrantes), HU-17 (usuarios y elegir quién soy), HU-16 (equipo del proyecto), #148 (navegar entre proyectos desde la interfaz) y la decisión del cierre manual del proyecto (#173).
### Notas
Decisión del Scrum Master (2026-11-03): esta historia queda absorbida por HU-20 y HU-21, porque la entrega final es el 16-nov y el login real entra antes de esa fecha.
1. **Nivel 1:** lista de proyectos y selector para navegar entre ellos (#148, PR #231). Es lo que queda de esta historia.
2. **Nivel 2 (usuario, "mis proyectos" y acciones según el rol):** pasa a HU-21.
3. **Nivel 3 (autenticación y permisos en el backend):** pasa a HU-20 (autenticación) y HU-21 (permisos).
Se descarta la transferencia de la propiedad (escenarios 2, 4 y 5): el dueño del proyecto es quien lo crea y queda como su único Scrum Master, así que la decisión D14 de `MODELO-DE-DATOS.md` (dónde guardar al dueño) deja de hacer falta. Un solo Scrum Master por proyecto es el issue #236.
Cerrar el proyecto de forma manual (#173) y eliminarlo (#232) son decisiones aparte.
### Prioridad
C (Could have) - Mejora de calidad surgida del beta testing; los niveles 2 y 3 pasaron a HU-21 y HU-20.
### Estado
Absorbida por HU-20 y HU-21
### Valor de Negocio
8 (Fibonacci) - propuesta, a confirmar con el equipo
### Estimación
Sin estimación propia: el trabajo restante se estima en HU-20 y HU-21

---

## HU-20: Autenticación de usuarios
### Descripción
Como integrante del equipo quiero registrarme e iniciar sesión con un nombre de usuario y una contraseña para que el sistema sepa quién soy y mis acciones queden a mi nombre.
### Criterios de aceptacion (BDD):
Escenario 1: Registro exitoso.
Dado un nombre de usuario que no existe,
Cuando me registro con una contraseña válida,
Entonces se crea mi cuenta y quedo con la sesión iniciada.
Escenario 2: Nombre de usuario repetido.
Dado que el nombre de usuario ya existe,
Cuando intento registrarme,
Entonces el sistema lo rechaza con el mensaje "Ese nombre de usuario ya existe".
Escenario 3: Contraseña inválida.
Dado que la contraseña tiene menos de 8 caracteres,
Cuando me registro,
Entonces el sistema lo rechaza indicando el campo.
Escenario 4: Ingreso exitoso.
Dado que tengo una cuenta,
Cuando ingreso con mi usuario y contraseña,
Entonces accedo a mis proyectos.
Escenario 5: Credenciales incorrectas.
Dado un usuario o una contraseña incorrectos,
Cuando intento ingresar,
Entonces el sistema responde "Usuario o contraseña incorrectos" sin revelar cuál de los dos falló.
Escenario 6: Sesión vencida.
Dado que mi sesión venció o es inválida,
Cuando uso la aplicación,
Entonces se me pide ingresar de nuevo y vuelvo a la pantalla en la que estaba.
Escenario 7: Cierre de sesión.
Cuando cierro sesión,
Entonces la sesión deja de ser válida.
### Depende de
HU-04 (usuarios y equipo).
### Notas
Es el PR 1 de 2 del login: registro, ingreso y sesión **sin imponer permisos**, para que lo existente y los escenarios de Godog sigan funcionando. Los permisos por rol son HU-21.
Usa nombre de usuario y contraseña, sin correo electrónico (no existe ese dato en el modelo). La contraseña se guarda con hash y sal usando `crypto/pbkdf2` de la biblioteca estándar de Go, sin dependencias nuevas, y nunca se devuelve ni se registra. Fuera de alcance: recuperar la contraseña, ingreso con terceros y verificación por correo.
Cambia el esquema (credenciales y sesiones): un solo PR con `MODELO-DE-DATOS.md` y la migración. Se deciden en la spec con el Scrum Master, no quien implementa: cómo viaja el token y su vencimiento, dónde viven las sesiones (en la base se vacían en cada deploy de Render), cómo se dan de alta las contraseñas de los usuarios que hoy no tienen y las reglas del nombre de usuario.
### Prioridad
S (Should have) - Hace que cada persona actúe a su nombre. El voto de Planning Poker realmente secreto se logra recién con HU-21, que quita el `usuario_id` del cuerpo y de la query; con HU-20 sola el servidor sigue aceptando el id que mande el cliente. Se planifica en el Sprint 4, antes de la entrega final, por decisión del Scrum Master (2026-11-03).
### Estado
Nueva
### Valor de Negocio
8 (Fibonacci) - propuesta, a confirmar con el equipo
### Estimación
5 Story Points - propuesta, a estimar con el equipo

---

## HU-21: Permisos por rol (segregación de funciones)
### Descripción
Como Scrum Master quiero que las acciones de gestión del proyecto solo las pueda hacer el Scrum Master para que el resto del equipo participe sin alterar la planificación.
### Criterios de aceptacion (BDD):
Escenario 1: Acción permitida.
Dado que soy el Scrum Master del proyecto,
Cuando creo, inicio o cierro un Sprint, inicio o revelo una votación, inicio una nueva ronda, registro la estimación acordada, edito el proyecto o configuro el equipo,
Entonces la acción se realiza.
Escenario 2: Acción no permitida.
Dado que soy Product Builder,
Cuando intento una de esas acciones,
Entonces el sistema responde 403 y la pantalla la muestra deshabilitada con la nota "Solo el Scrum Master puede…".
Escenario 3: Sin sesión.
Dado que no inicié sesión,
Cuando uso una de las acciones que exigen sesión (ver "Qué se protege" en las Notas),
Entonces el sistema responde 401.
Escenario 4: No integrante.
Dado que no soy integrante del proyecto,
Cuando intento acceder a él,
Entonces el sistema responde 403 con el mensaje "No sos integrante de este proyecto".
Escenario 5: Mis proyectos.
Dado que inicié sesión,
Cuando consulto mis proyectos,
Entonces veo solo aquellos donde soy integrante, con mi rol.
Escenario 6: Un solo Scrum Master.
Dado que quien crea un proyecto queda como Scrum Master,
Cuando se intenta agregar un segundo Scrum Master,
Entonces el sistema lo rechaza.
Escenario 7: Identidad de la sesión.
Dado que inicié sesión,
Cuando voto, registro horas o consulto la votación de una historia,
Entonces queda a nombre de la persona con sesión iniciada, y no se acepta otro usuario ni en el cuerpo ni en la query de la petición (`usuario_id`), así que nadie puede ver el voto oculto de otra persona.
Escenario 8: Configurar equipo solo con usuarios registrados.
Dado que soy el Scrum Master del proyecto,
Cuando configuro el equipo,
Entonces solo puedo agregar usuarios ya registrados, y entran como Product Builder.
Escenario 9: Usuario sin proyectos.
Dado que inicié sesión y no soy integrante de ningún proyecto,
Cuando consulto mis proyectos,
Entonces el sistema muestra un mensaje de que no hay proyectos y me permite crear uno.
### Depende de
HU-20 (autenticación), HU-02 (Planning Poker), HU-04, HU-05 (Sprints) y HU-07 (horas). Se apoya en el issue #236 (índice único de Scrum Master).
### Notas
Es el PR 2 de 2 del login. Los endpoints toman el usuario de la sesión: se elimina `usuario_id` del cuerpo de la votación y `integrante_id` del registro de horas, y también la identidad que viaja en la query de `GET /backlog/{id}/votacion?usuario_id=N` (`internal/http/votacion_handler.go`), que hoy devuelve el voto del que consulta: con el id de otra persona se vería su voto oculto. Quien crea un proyecto queda como su Scrum Master. "Configurar equipo" solo agrega usuarios ya registrados, como Product Builder. Hay un único Scrum Master por proyecto, lo que cierra la decisión D4 de `MODELO-DE-DATOS.md`; la decisión D14 (dueño) deja de hacer falta.
**Qué se protege.** Todas las rutas que escriben datos exigen sesión (y ser integrante del proyecto, escenario 4); las de gestión exigen además ser Scrum Master (escenario 2):
- Solo sesión: `POST /projects` (quien lo crea queda como Scrum Master), `POST /backlog`, `PUT /backlog/{id}`, `PUT /backlog/{id}/estado`, `PUT /backlog/{id}/sprint`, `POST /backlog/{id}/horas`, `POST /projects/{id}/defectos`, `POST /defectos/{id}/resolver`, `POST /backlog/{id}/votacion/votos` y `POST /demo/datos`.
- Sesión y Scrum Master: `POST /sprints`, `POST /sprints/{id}/iniciar`, `POST /sprints/{id}/cerrar`, `POST /backlog/{id}/votacion` (iniciar), `POST /backlog/{id}/votacion/revelar`, `POST /backlog/{id}/votacion/rondas` (nueva ronda), `POST /backlog/{id}/votacion/estimacion` (registrar la estimación, HU-14), `PUT /projects/{id}` y `POST /projects/{id}/members` (equipo).
- Lecturas (`GET`): a decidir si exigen sesión. <!-- decidir: ¿los GET exigen sesión? (sí/no, y si se limitan al integrante del proyecto) -->

Los escenarios de Godog que hoy mandan `usuario_id` se ajustan **en el mismo PR**, con un paso de contexto "Dado que inicié sesión como X" (el check `backend` tiene que estar en verde para mergear). Sin atajos de prueba en el código de producción.
Si el viernes 13 los permisos en el backend no están completos, se entrega con las acciones ocultas en la pantalla y se declara el límite en la documentación.
Absorbe los niveles 2 y 3 de HU-19. Eliminar un proyecto (#232) depende de estos permisos.
### Prioridad
S (Should have) - Complementa a HU-20: sin permisos en el backend, cualquiera con `curl` puede saltear la pantalla.
### Estado
Nueva
### Valor de Negocio
8 (Fibonacci) - propuesta, a confirmar con el equipo
### Estimación
5 Story Points - propuesta, a estimar con el equipo (toca muchos endpoints)

---

# FIXES DEL SPRINT 1 INCLUIDOS EN EL SPRINT 2

Defectos encontrados durante el Sprint 1 (que terminó el 2026-10-05). Entran al
backlog del Sprint 2 y se reparan cada uno en su propia rama `fix/<n°-issue>-nombre`
(ver la convención de Git en `AGENTS.md`). Cada issue tiene la etiqueta `bug`.

| Issue | Fix | Historia | Rama sugerida |
|---|---|---|---|
| #19 | Integridad referencial de `proyecto_id` en historias y Sprints: llave foránea en la base, el caso de uso `CrearHistoriaBacklog` no valida que el proyecto exista, y `ErrIntegridadReferencial` se responde como 500 | HU-01, HU-05 | `fix/19-integridad-proyecto-id` |
| #21 | Los errores 404 de la API incluyen `"campo": ""` vacío | transversal | `fix/21-404-sin-campo-vacio` |
| #99 | Al cerrar un Sprint, las historias no completadas vuelven a "Nueva" | HU-05 | `fix/99-estado-nueva-al-arrastrar` |
| #100 | La fecha de fin de un proyecto debe ser estrictamente posterior a la de inicio | HU-04 | `fix/100-fecha-fin-posterior-proyecto` |
| #101 | Pendientes del code review de HU-04: fake en memoria de proyectos y 404 antes que 400 al asignar integrantes | HU-04 | `fix/101-pendientes-review-hu04` |
| #102 | `.gitattributes` para que `gofmt` no marque todos los archivos en Windows (CRLF) | transversal | `fix/102-gitattributes-lf` |

**Continuación de #99:** el PR #115 ya resolvió que las historias no completadas vuelvan a "Nueva". El registro de las historias arrastradas, con sus Story Points, y el cierre atómico son el issue #117, que depende de #107.

**Decidido por el equipo (2026-10-14), no es un defecto:** #20. Un Sprint se puede iniciar
con fechas de inicio y de fin anteriores a hoy, porque el equipo pudo haberlo empezado de
hecho y olvidar declararlo en la aplicación. No se restringe; solo el fin debe ser posterior al inicio.

---

# FIXES DEL SPRINT 2 INCLUIDOS EN EL SPRINT 3

Defectos encontrados durante el Sprint 2 (que terminó el 2026-10-26) y deudas que dejaron
sus reviews. Cada uno tiene su issue y se repara en su propia rama `fix/<n°-issue>-nombre`
(ver la convención de Git en `AGENTS.md`); los que se mergean después de `v0.2.0` reciben
un tag `v0.2.P`.

| Issue | Fix | Historia | Estado |
|---|---|---|---|
| #117 | El cierre de un Sprint registra las historias arrastradas con sus Story Points, de forma atómica (migración 002) | HU-05 | Mergeado (PR #136, `v0.2.1`) |
| #151 | Completar la spec y los tests del estado del proyecto en el panel | HU-13 | Mergeado (PR #152, `v0.2.2`) |
| #145 | Dos cierres simultáneos de un Sprint pasaban la validación | HU-05 | Mergeado (PR #155, `v0.2.3`) |
| #158 | Dos inicios simultáneos del mismo Sprint pasaban la validación | HU-05 | Mergeado (PR #160, `v0.2.4`) |
| #128 | Con el backend apagado, el cliente compartido muestra «502» en vez del error de conexión | frontend base | Abierto, tarea libre |
| #129 | Los Sprints se numeran de forma global y no por proyecto | HU-05 | Abierto, tarea libre |
| #130 | `docs/curl-frontend-sprint-2.md` no tiene el comando para cerrar un Sprint | transversal | Abierto, tarea libre |
| #137 | Ajustes al workflow de CI sugeridos en la review del PR #133 | transversal | Abierto (Carlos) |
| #146 | El paso BDD del cierre con arrastre usa SQL directo en vez de `EditarHistoria` | HU-05 | Mergeado (PR #163) |
| #159 | Documentar el cierre condicional en `MODELO-DE-DATOS.md` §4.6 y en el plan de la spec 003 | HU-05 | Abierto, Sprint 4 |
| #172 | Al recargar la página, las listas de Sprints y del backlog pueden no cargar (falta `Vary: Accept` en `ServirSPA`) | transversal | Abierto, Sprint 3 |

---

# MEJORAS TÉCNICAS Y DE USABILIDAD (NO SON HISTORIAS DE USUARIO)

Trabajo que no agrega requisitos del enunciado y por eso no es una historia, pero tiene
issue y entra en los Sprints.

| Issue | Mejora | Estado |
|---|---|---|
| #139 | Despliegue continuo en Render: un servicio con la imagen Docker, que se actualiza en cada merge a `main` | Abierto (Carlos) |
| #144 | Botón «Cargar datos de prueba» (endpoint de demo y pantalla), para el servicio de Render, que no conserva la base | Abierto (Carlos) |
| #148 | Navegar entre proyectos desde la interfaz (lista y selector) | Abierto (Jimena) |
| #147 | Mostrar con Markdown la descripción y los criterios de aceptación de las historias | Abierto, tarea libre (depende de #98) |
| #140 | Test de las rutas HTTP y del cableado de `cmd/api/main.go` | Abierto, tarea libre |
| #141 | Eliminar `Inicio.tsx`, que ningún archivo importa | Abierto, tarea libre |
| #168 | Rechazar Sprints con fechas fuera del plazo del proyecto (backend) | Abierto, a priorizar |
| #169 | Mostrar las fechas del proyecto en la pantalla de Sprints y limitar los selectores | Abierto, a priorizar |
| #173 | Regla del estado del proyecto: cierre manual en vez de «Finalizado» derivado de los Sprints (HU-13) | Abierto, a decidir |
| #174 | Cargar criterios de aceptación, notas y «depende de» al crear una historia | Abierto, Sprint 4 (solo los criterios de aceptación; notas y «depende de» quedan para después) |
| #230 | Rediseño visual del frontend y mejor aprovechamiento del ancho de pantalla (a partir de un prototipo) | Abierto, Sprint 5 por pantallas |
| #232 | Eliminar un proyecto: solo el Scrum Master, con triple confirmación | Abierto, Sprint 5 |
| #235 | Datos de prueba v2: usuarios con contraseña, equipo, historias con criterios, Sprints, horas, votaciones y defectos | Abierto, Sprint 4 |
| #236 | Un único Scrum Master por proyecto (índice único parcial, cambio de esquema) | Abierto, Sprint 4 |
