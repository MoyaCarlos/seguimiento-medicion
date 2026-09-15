# REQUISITOS DE USUARIO

> **Convención de campos:**
> - **Prioridad** usa escala MoSCoW (Must/Should/Could/Won't have). Al cargar en el
>   tablero de GitHub Projects, mapear: `M` → Alta, `S` → Media, `C` → Baja.
> - **Valor de Negocio** y **Estimación** usan ambas una escala tipo Fibonacci, pero
>   miden cosas distintas e independientes (como en WSJF): Valor de Negocio es el
>   impacto/beneficio relativo para el proyecto; Estimación es el esfuerzo/complejidad
>   relativa (Story Points). No están relacionadas entre sí ni deben coincidir.
> - **Estado** inicial de toda historia recién creada: `Nuevo`.

## ***Prioridad 1: Cimientos del Sistema (Dependencia Crítica)***

No se puede gestionar nada si no existe un contenedor. Esta es la base estructural.

* **HU-04: Creación de Proyecto y Asignación de Equipo (Must have).** Es la máxima prioridad operativa. Todo el sistema (historias, sprints, horas) depende de que exista un proyecto y usuarios registrados.  
* **HU-01: Creación de Historias de Usuario (Must have).** Sin esto, no hay Product Backlog.
* **HU-13: Fechas y Estado del Proyecto (Should have).** Extiende HU-04: registrar fecha de inicio/fin y consultar el estado general del proyecto (requisito explícito del punto 1 del enunciado).

## ***Prioridad 2: El Motor de Scrum (Iteraciones)***

Una vez que tienen el proyecto y el Backlog, necesitan el ciclo de trabajo ágil.

* **HU-05: Apertura y Cierre de Sprints (Must have).** Da vida a la metodología.  
* **HU-06: Movimiento de Historias al Sprint Backlog (Must have).** Permite a las Product Builders definir en qué se va a trabajar durante la iteración.
* **HU-11: Actualización de Estado de una Historia en el Sprint (Must have).** Sin poder marcar una historia como completada durante el sprint, no hay datos reales para calcular velocidad ni ninguna métrica.
* **HU-12: Consulta de Sprints Anteriores (Should have).** Requisito explícito del punto 3 del enunciado; de baja complejidad técnica una vez que existen sprints cerrados.

## ***Prioridad 3: Alimentación de Datos (El Valor de Negocio)***

Para que el diferenciador del sistema (las métricas) funcione, primero necesitan generar y registrar datos reales.

* **HU-02: Votación de Esfuerzo en Planning Poker (Should have).** Aporta un gran valor metodológico y permite tener los "Story Points" estimados.  
* **HU-07: Imputación de Horas Trabajadas (Must have).** Requisito indispensable para luego poder cruzar lo estimado vs. lo real.

## ***Prioridad 4: Resultados y Calidad (El Valor Final)***

Estas funcionalidades son el corazón del valor para el cliente final, pero tienen una alta dependencia técnica de las etapas anteriores. Se dividió la antigua HU-03 en tres historias más chicas (cálculo, visualización y exportación son responsabilidades distintas) y se agregó cobertura completa del punto 7 del enunciado.

* **HU-03: Cálculo de Métricas del Proyecto (Must have).** El motor de cálculo (Story Points planificados/completados, velocidad, horas estimadas/reales, desviación, % completadas, defectos detectados/resueltos) — requisito central del enunciado (punto 7).
* **HU-09: Dashboard de Métricas (Must have).** Visualización gráfica de lo que calcula HU-03 (punto 8 del enunciado).
* **HU-10: Generación de Reporte en PDF (Must have).** Exportación a PDF con Maroto de historias, estimaciones, esfuerzo, métricas y defectos (punto 9 del enunciado).
* **HU-08: Registro y Seguimiento de Bugs (Should have).** Completa el ciclo de aseguramiento de calidad del software desarrollado.

---

# HISTORIAS DE USUARIO

## HU-01: Creación de Historias de Usuario 

### Descripción
Como Product Builder quiero crear historias de usuario con prioridad, estado y estimación para poder alimentar y organizar el Product Backlog del proyecto.

### Criterios de aceptacion (BDD):
* Escenario 1: Creación exitosa.
Dado que me encuentro logueado en el panel del Product Backlog,
Cuando ingreso los datos obligatorios (título, descripción, prioridad) y presiono "Guardar",
Entonces la historia debe aparecer al final de la lista con estado "Nuevo" y guardarse en la base de datos SQLite.
* Escenario 2: Faltan campos.
Dado que intento crear una historia,
Cuando dejo el campo "Título" en blanco y presiono "Guardar",
Entonces el sistema debe mostrar una advertencia y no debe registrar la HU.

### Prioridad
M (Must have) - Esencial para el MVP. 

### Estado
Nuevo

### Valor de Negocio
21 (Fibonacci) 

### Estimación
5 Story Points 

---

## HU-02: Votación de Esfuerzo en Planning Poker 
### Descripción
Como integrante del equipo (Scrum Master o Product Builder) quiero votar el esfuerzo de una tarea en secreto usando cartas virtuales para poder realizar una estimación colaborativa sin sesgos.
### Criterios de aceptacion (BDD):
Escenario 1: Votación secreta.
Dado que el Scrum Master inicia la votación para la HU seleccionada,
Cuando selecciono una carta de la baraja Fibonacci (ej. 8) y confirmo mi voto,
Entonces mi voto debe registrarse en el backend (Go) pero mostrarse como "Oculto" en el frontend (React) para el resto del equipo.
Escenario 2: Revelación de votos.
Dado que todos los integrantes confirmaron su voto,
Cuando el Scrum Master presiona "Revelar",
Entonces el sistema debe mostrar el valor elegido por cada uno y calcular el promedio sugerido.
### Prioridad
S (Should have) 
### Estado
Nuevo
### Valor de Negocio
13 (Fibonacci) 
### Estimación
8 Story Points 

---

## HU-03: Cálculo de Métricas del Proyecto 
### Descripción
Como Scrum Master quiero que el sistema calcule automáticamente las métricas del proyecto y de cada Sprint para poder analizar la capacidad real del equipo y la calidad del trabajo entregado.
### Criterios de aceptacion (BDD):
Escenario 1: Cálculo tras cerrar un Sprint.
Dado que un Sprint se cierra (HU-05) con historias marcadas como completadas (HU-11),
Cuando el sistema recalcula las métricas del proyecto,
Entonces debe obtener: Story Points planificados y completados, velocidad del equipo, horas estimadas y reales, desviación entre esfuerzo estimado y real, porcentaje de historias completadas, y cantidad de defectos detectados y resueltos en ese Sprint.
Escenario 2: Proyecto sin Sprints cerrados.
Dado que un proyecto todavía no tiene ningún Sprint finalizado,
Cuando se solicita el cálculo de métricas,
Entonces el sistema debe devolver todos los valores en cero (o "sin datos") en vez de fallar o mostrar un error.
### Prioridad
M (Must have) - Requisito central del enunciado (punto 7).
### Estado
Nuevo
### Valor de Negocio
21 (Fibonacci)
### Estimación
8 Story Points (cálculos de dominio, sin UI — se testea con TDD directo sobre `/internal/domain`) 


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
### Prioridad
M (Must have) - Requisito fundamental para que el sistema funcione. 
### Estado
Nuevo
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

### Prioridad
M (Must have) - Sin esto no hay marco iterativo ágil. 
### Estado
Nuevo
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
Cuando arrastro o selecciono una HU del Product Backlog y la asigno al Sprint,
Entonces la HU desaparece de la vista del backlog general y pasa a conformar el Sprint Backlog de esa iteración.
Escenario 2: Restricción de historias terminadas.
Dado que estoy planificando un Sprint,
Cuando intento asignar una HU cuyo estado ya es "Terminado" (Done),
Entonces el sistema me muestra un mensaje de error indicando que solo se pueden planificar tareas pendientes.

### Prioridad
M (Must have) - Core del flujo de trabajo de Scrum. 
### Estado
Nuevo
### Valor de Negocio
13 (Fibonacci) 

### Estimación
5 Story Points

---

## HU-07: Imputación de Horas Trabajadas 

### Descripción
Como integrante del equipo
quiero registrar las horas reales que invertí trabajando en una Historia de Usuario o Tarea
para poder proveer datos al sistema que permitan comparar el esfuerzo estimado (Story Points) contra el esfuerzo real.

### Criterios de aceptacion (BDD):
Escenario 1: Registro válido de tiempo.
Dado que estoy asignado a una HU en un Sprint activo,
Cuando abro el formulario de registro, ingreso "4 horas" correspondientes al día de hoy y guardo,
Entonces el sistema suma esas 4 horas al esfuerzo total acumulado de esa HU.
Escenario 2: Validación de inputs.
Dado que intento cargar horas a una tarea,
Cuando ingreso un valor negativo (-2 horas) o texto no numérico,
Entonces el sistema bloquea el guardado y muestra un error de validación en la interfaz de React.

### Prioridad
M (Must have) - Fundamental para que luego se puedan calcular las métricas de desviación de esfuerzo. 
### Estado
Nuevo
### Valor de Negocio
13 (Fibonacci) 

### Estimación
5 Story Points

---

## HU-08: Registro y Seguimiento de Bugs 
### Descripción
Como integrante del equipo
quiero registrar un defecto técnico indicando su severidad y sprint de detección
para poder hacer un seguimiento de la calidad y asegurar que se resuelva antes del despliegue final.

### Criterios de aceptacion (BDD):
Escenario 1: Carga de un nuevo defecto.
Dado que se detecta un error técnico en el software,
Cuando completo el formulario de Defectos con el título, descripción, nivel de severidad (Alta, Media, Baja) y guardo,
Entonces el Bug se registra en el sistema asociado al Sprint actual como "Pendiente".
Escenario 2: Resolución del defecto.
Dado que un desarrollador solucionó un Bug pendiente,
Cuando cambia su estado a "Resuelto",
Entonces el sistema registra internamente el Sprint de resolución para que posteriormente aparezca en el Dashboard de métricas de calidad.

### Prioridad
S (Should have) - Es muy importante para las métricas de calidad, pero el proyecto podría arrancar los primeros Sprints funcionales sin este módulo. 
### Estado
Nuevo
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
### Prioridad
M (Must have) - Requisito central del enunciado (punto 8).
### Estado
Nuevo
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
### Prioridad
M (Must have) - Requisito central del enunciado (punto 9).
### Estado
Nuevo
### Valor de Negocio
13 (Fibonacci)
### Estimación
8 Story Points (integración con Maroto, depende de HU-03 y HU-08 para los datos)

---

## HU-11: Actualización de Estado de una Historia en el Sprint 
### Descripción
Como Product Builder quiero cambiar el estado de una historia asignada al Sprint (En progreso, Completada) para poder reflejar el avance real del trabajo durante la iteración.
### Criterios de aceptacion (BDD):
Escenario 1: Marcar historia como completada.
Dado que una historia está asignada a un Sprint activo con estado "En progreso",
Cuando la marco como "Completada",
Entonces el sistema registra la fecha de finalización y la incluye en el cálculo de Story Points completados del Sprint (HU-03).
Escenario 2: Transición inválida.
Dado que una historia todavía no fue asignada a ningún Sprint,
Cuando intento marcarla como "Completada" directamente desde el Product Backlog,
Entonces el sistema rechaza el cambio y muestra un mensaje indicando que primero debe asignarse a un Sprint activo.
### Prioridad
M (Must have) - Sin esto no hay datos reales para ninguna métrica.
### Estado
Nuevo
### Valor de Negocio
13 (Fibonacci)
### Estimación
3 Story Points

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
### Prioridad
S (Should have) - Requisito explícito del enunciado (punto 3), baja complejidad técnica.
### Estado
Nuevo
### Valor de Negocio
8 (Fibonacci)
### Estimación
3 Story Points

---

## HU-13: Fechas y Estado del Proyecto 
### Descripción
Como Scrum Master quiero registrar la fecha de inicio y finalización de un proyecto y consultar su estado general para poder tener visibilidad del ciclo de vida completo del proyecto.
### Criterios de aceptacion (BDD):
Escenario 1: Registro de fechas.
Dado que estoy editando un proyecto existente (HU-04),
Cuando ingreso una fecha de inicio y una fecha de finalización estimada y guardo,
Entonces el sistema valida que la fecha de fin sea posterior a la de inicio y las persiste.
Escenario 2: Consulta de estado.
Dado que un proyecto tiene Sprints y fechas cargadas,
Cuando consulto su estado,
Entonces el sistema muestra si está "Planificado", "En curso" o "Finalizado" según la fecha actual y el estado de sus Sprints.
### Prioridad
S (Should have) - Extiende HU-04, no bloquea el resto del sistema.
### Estado
Nuevo
### Valor de Negocio
8 (Fibonacci)
### Estimación
3 Story Points

