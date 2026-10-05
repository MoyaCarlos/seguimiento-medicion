# language: es

@HU-13
Característica: Consulta del estado general del proyecto

  Como Scrum Master
  Quiero registrar la fecha de inicio y finalización de un proyecto y consultar su estado general
  Para poder tener visibilidad del ciclo de vida completo del proyecto

  Antecedentes:
    Dado que existe un proyecto con identificador 1

  Escenario: Proyecto con un Sprint Activo muestra "En curso"
    Dado que el proyecto tiene fecha de inicio y fecha de fin cargadas
    Y que el proyecto cuenta con un Sprint en estado "Activo"
    Cuando consulto el estado del proyecto
    Entonces el estado del proyecto es "En curso"

  Escenario: Proyecto con todos sus Sprints Finalizados muestra "Finalizado"
    Dado que el proyecto tiene todos sus Sprints en estado "Finalizado"
    Cuando consulto el estado del proyecto
    Entonces el estado del proyecto es "Finalizado"

  Escenario: Proyecto con fecha de inicio futura y sin Sprints iniciados muestra "Planificado"
    Dado que el proyecto tiene una fecha de inicio futura
    Y que el proyecto no tiene Sprints iniciados
    Cuando consulto el estado del proyecto
    Entonces el estado del proyecto es "Planificado"

  Escenario: Proyecto sin Sprints iniciados y con la fecha actual dentro del rango muestra "En curso"
    Dado que el proyecto tiene la fecha actual dentro de su rango de fechas
    Y que el proyecto no tiene Sprints iniciados
    Cuando consulto el estado del proyecto
    Entonces el estado del proyecto es "En curso"

  Escenario: Consulta del estado de un proyecto inexistente
    Cuando consulto el estado de un proyecto que no existe
    Entonces el sistema responde "proyecto no encontrado"

  Escenario: La fecha actual coincide exactamente con la fecha de inicio
    Dado que el proyecto tiene la fecha actual como fecha de inicio
    Y que el proyecto no tiene Sprints iniciados
    Cuando consulto el estado del proyecto
    Entonces el estado del proyecto es "En curso"

  Escenario: La fecha actual coincide exactamente con la fecha de fin
    Dado que el proyecto tiene la fecha actual como fecha de fin
    Y que el proyecto no tiene Sprints iniciados
    Cuando consulto el estado del proyecto
    Entonces el estado del proyecto es "En curso"

  Escenario: El estado se actualiza al iniciar un Sprint sin recálculo manual
    Dado que el proyecto tiene una fecha de inicio futura
    Y que consulté el estado del proyecto y era "Planificado"
    Cuando inicio un Sprint del proyecto y consulto su estado nuevamente
    Entonces el estado del proyecto es "En curso"

  Escenario: El estado se actualiza al cerrar el último Sprint sin recálculo manual
    Dado que el proyecto cuenta con un Sprint en estado "Activo"
    Y que consulté el estado del proyecto y era "En curso"
    Cuando cierro su último Sprint y consulto su estado nuevamente
    Entonces el estado del proyecto es "Finalizado"

  Escenario: Proyecto con un Sprint Finalizado y uno Pendiente muestra "En curso"
    Dado que el proyecto tiene un Sprint en estado "Finalizado" y uno "Pendiente"
    Cuando consulto el estado del proyecto
    Entonces el estado del proyecto es "En curso"
