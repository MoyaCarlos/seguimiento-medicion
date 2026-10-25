# language: es

@HU-04
Característica: Listado de proyectos

  Como Scrum Master
  Quiero obtener todos los proyectos existentes
  Para poder elegir cuál consultar o editar sin conocer su identificador

  Escenario: Listado ordenado por identificador
    Dado que existen proyectos guardados con identificadores 9 y 2
    Cuando consulto el listado de proyectos
    Entonces recibo los proyectos con identificadores 2 y 9 en ese orden
    Y cada proyecto del listado tiene los mismos campos que su consulta individual

  Escenario: Listado sin proyectos
    Dado que no hay proyectos guardados
    Cuando consulto el listado de proyectos
    Entonces recibo un arreglo vacío

  Escenario: Error de persistencia al listar
    Dado que la persistencia de proyectos falla
    Cuando consulto el listado de proyectos
    Entonces el sistema responde con un error interno 500
