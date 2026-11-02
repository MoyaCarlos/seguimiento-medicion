# language: es

@HU-05
Característica: Consulta de los Sprints de un proyecto

  Como Scrum Master
  Quiero consultar los Sprints de un proyecto
  Para ver el estado de las iteraciones del equipo

  Escenario: Sprints ordenados por identificador
    Dado que el proyecto 1 tiene los Sprints 9 y 2
    Cuando consulto los Sprints del proyecto 1
    Entonces el sistema responde 200
    Y recibo los Sprints 2 y 9 en ese orden

  Escenario: Cada proyecto numera sus Sprints desde 1
    Dado que el proyecto 1 tiene los Sprints 1 y 2
    Y que el proyecto 2 tiene los Sprints 3 y 4
    Cuando consulto los Sprints del proyecto 2
    Entonces el sistema responde 200
    Y los Sprints recibidos tienen los números 1 y 2

  Escenario: Proyecto sin Sprints
    Dado que el proyecto 1 no tiene Sprints
    Cuando consulto los Sprints del proyecto 1
    Entonces el sistema responde 200
    Y recibo un arreglo de Sprints vacío

  Escenario: Proyecto inexistente
    Cuando consulto los Sprints del proyecto 999
    Entonces el sistema responde 404

  Escenario: Identificador no numérico
    Cuando consulto los Sprints del proyecto "abc"
    Entonces el sistema responde 400
