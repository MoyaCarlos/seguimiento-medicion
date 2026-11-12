# language: es

@HU-04
Característica: Listado de proyectos

  Como integrante de proyectos
  Quiero obtener los proyectos que integro
  Para poder elegir cuál consultar o editar sin conocer su identificador

  Antecedentes:
    Dado que inicié sesión como "Ana Valentina"

  Escenario: Listado ordenado por identificador
    Dado que integro proyectos guardados con identificadores 9 y 2
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
