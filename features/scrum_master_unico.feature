# language: es

@HU-04
Característica: Un único Scrum Master por proyecto

  Como Scrum Master
  Quiero que cada proyecto tenga un único Scrum Master
  Para que las acciones exclusivas del rol nunca queden repartidas entre varias personas


  # Desde HU-21 el proyecto lo crea la persona con sesión iniciada.
  Antecedentes:
    Dado que inicié sesión como "Ana Valentina"

  Escenario: El proyecto nace con un único Scrum Master
    Dado que me encuentro en la pantalla principal del sistema
    Cuando ingreso el nombre "Software Metrics & Estimation", una descripción y presiono "Crear"
    Entonces el proyecto tiene un único Scrum Master
    Y ese Scrum Master es quien lo creó

  Escenario: Rechazo de un segundo Scrum Master
    Dado que existe un proyecto creado
    Y estoy en la vista de configuración del proyecto
    Cuando intento agregar a otra persona con el rol "Scrum Master"
    Entonces el sistema rechaza la operación con un conflicto que indica que el proyecto ya tiene Scrum Master
    Y la lista de integrantes del proyecto no cambia

  Escenario: Alta de un Product Builder con Scrum Master presente
    Dado que existe un proyecto creado
    Y estoy en la vista de configuración del proyecto
    Cuando ingreso el nombre de un integrante, selecciono su rol "Product Builder" y confirmo
    Entonces el sistema vincula al integrante con el proyecto
    Y el Scrum Master del proyecto no cambia

  Escenario: La base rechaza un segundo Scrum Master insertado por fuera del dominio
    Dado que existe un proyecto creado
    Cuando se inserta directamente en la base un segundo Scrum Master para ese proyecto
    Entonces la base rechaza la inserción
    Y el proyecto sigue con su único Scrum Master
