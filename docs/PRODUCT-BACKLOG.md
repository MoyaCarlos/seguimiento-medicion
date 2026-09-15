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

## ***Prioridad 2: El Motor de Scrum (Iteraciones)***

Una vez que tienen el proyecto y el Backlog, necesitan el ciclo de trabajo ágil.

* **HU-05: Apertura y Cierre de Sprints (Must have).** Da vida a la metodología.  
* **HU-06: Movimiento de Historias al Sprint Backlog (Must have).** Permite a las Product Builders definir en qué se va a trabajar durante la iteración.

## ***Prioridad 3: Alimentación de Datos (El Valor de Negocio)***

Para que el diferenciador del sistema (las métricas) funcione, primero necesitan generar y registrar datos reales.

* **HU-02: Votación de Esfuerzo en Planning Poker (Should have).** Aporta un gran valor metodológico y permite tener los "Story Points" estimados.  
* **HU-07: Imputación de Horas Trabajadas (Must have).** Requisito indispensable para luego poder cruzar lo estimado vs. lo real.

## ***Prioridad 4: Resultados y Calidad (El Valor Final)***

Estas funcionalidades son el corazón del valor para el cliente final, pero tienen una alta dependencia técnica de las etapas anteriores.

* **HU-03: Visualización de Velocidad del Equipo y Reportes (Must have).** El Dashboard y los reportes PDF (Maroto) consumirán la información generada en la Prioridad 3\.  
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

## HU-03: Visualización de Velocidad del Equipo (Velocity) 
### Descripción
Como Scrum Master quiero visualizar un gráfico que compare los Story Points estimados vs. completados por cada Sprint para poder analizar la capacidad real del equipo y mejorar la planificación futura.
### Criterios de aceptacion (BDD):
Escenario 1: Generación del gráfico de Velocity.
Dado que el proyecto tiene al menos un Sprint finalizado,
Cuando ingreso a la sección "Dashboard de Métricas",
Entonces el sistema debe renderizar un gráfico de barras donde el eje X sean los Sprints y el eje Y muestre dos barras por Sprint: "Puntos Comprometidos" (gris) y "Puntos Completados" (verde).
Escenario 2: Generación de reporte PDF.
Dado que estoy visualizando las métricas,
Cuando presiono el botón "Exportar a PDF",
Entonces la librería Maroto debe procesar los datos y descargar un archivo local con el resumen del gráfico.
### Prioridad
M (Must have) - Requisito central del enunciado.
### Estado
Nuevo
### Valor de Negocio
21 (Fibonacci)
### Estimación
13 Story Points (Requiere integración de gráficos y librería Maroto) 


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

