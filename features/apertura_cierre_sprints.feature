# language: es

@HU-05
Característica: Apertura y cierre de Sprints

  Como Scrum Master
  Quiero crear, iniciar y cerrar un Sprint definiendo su objetivo y plazos
  Para organizar y delimitar las iteraciones de desarrollo del equipo

  Antecedentes:
    Dado que existe un proyecto con identificador 1

  Escenario: Creación exitosa de un Sprint
    Cuando creo un Sprint para el proyecto
    Entonces el Sprint queda persistido en estado "Pendiente"

  Escenario: Creación de un Sprint para un proyecto inexistente
    Cuando intento crear un Sprint para el proyecto con identificador 999, que no existe
    Entonces el sistema rechaza la operación

  Escenario: Iniciar un Sprint exitosamente
    Dado que existe un Sprint en estado "Pendiente" para el proyecto, sin otro Sprint "Activo"
    Cuando defino el Sprint Goal "Entregar el MVP", una fecha de inicio y una fecha de fin posterior, y lo inicio
    Entonces el Sprint pasa al estado "Activo"

  Escenario: No se puede iniciar un segundo Sprint Activo en el mismo proyecto
    Dado que el proyecto ya tiene un Sprint en estado "Activo"
    Y existe otro Sprint en estado "Pendiente" para el mismo proyecto
    Cuando intento iniciar ese segundo Sprint
    Entonces el sistema rechaza la operación
    Y el segundo Sprint permanece en estado "Pendiente"

  Escenario: Fecha de fin no posterior a la fecha de inicio
    Dado que existe un Sprint en estado "Pendiente" para el proyecto
    Cuando intento iniciarlo con una fecha de fin igual o anterior a la fecha de inicio
    Entonces el sistema rechaza la operación por validación de fechas

  Escenario: Cerrar un Sprint con arrastre de historias no completadas
    Dado que el proyecto tiene un Sprint en estado "Activo"
    Y ese Sprint tiene una historia completada y otra no completada asignadas
    Cuando el Scrum Master cierra el Sprint
    Entonces el Sprint pasa al estado "Finalizado"
    Y la historia no completada queda sin Sprint asignado en el Product Backlog
    Y la historia completada permanece vinculada a ese Sprint

  Escenario: No se puede cerrar un Sprint que no está Activo
    Dado que existe un Sprint en estado "Pendiente" para el proyecto
    Cuando intento cerrarlo
    Entonces el sistema rechaza la operación
