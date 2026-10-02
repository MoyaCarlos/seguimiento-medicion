# language: es

@HU-04
Característica: Creación de Proyecto y Asignación de Equipo

  Como Scrum Master
  Quiero crear un proyecto y agregar a los integrantes con sus respectivos roles
  Para poder tener un espacio de trabajo organizado y estructurar las responsabilidades del equipo

  Escenario: Creación del entorno del proyecto
    Dado que me encuentro en la pantalla principal del sistema
    Cuando ingreso el nombre "Software Metrics & Estimation", una descripción y presiono "Crear"
    Entonces el sistema genera el proyecto de forma persistente
    Y confirma su creación devolviendo el proyecto con su identificador
    Y habilita su panel principal, cuya navegación realiza el frontend

  Escenario: Creación con nombre vacío
    Dado que intento crear un proyecto
    Cuando dejo el nombre vacío o compuesto solo por espacios y presiono "Crear"
    Entonces el sistema muestra una advertencia de validación
    Y no genera el proyecto

  Escenario: Asignación de integrante y rol
    Dado que existe un proyecto creado
    Y estoy en la vista de configuración del proyecto
    Cuando ingreso el nombre de un integrante, selecciono su rol "Product Builder" y confirmo
    Entonces el sistema vincula al integrante con el proyecto
    Y le habilita los permisos correspondientes a su rol

  Escenario: Asignación con datos inválidos
    Dado que existe un proyecto creado
    Y estoy en la vista de configuración del proyecto
    Cuando ingreso un integrante sin nombre o con un rol inválido y confirmo
    Entonces el sistema muestra una advertencia de validación
    Y no realiza la vinculación

  Escenario: Edición de un proyecto existente
    Dado que existe un proyecto creado
    Y estoy en la vista de edición del proyecto
    Cuando modifico su nombre, su descripción o sus fechas con datos válidos y guardo
    Entonces el sistema persiste los cambios del proyecto

  Escenario: Edición con nombre vacío
    Dado que existe un proyecto creado
    Y estoy en la vista de edición del proyecto
    Cuando dejo el nombre vacío o compuesto solo por espacios y guardo
    Entonces el sistema muestra una advertencia de validación
    Y no guarda los cambios

  Escenario: Edición con fecha de fin anterior a la fecha de inicio
    Dado que existe un proyecto creado
    Y estoy en la vista de edición del proyecto
    Cuando ingreso una fecha de fin anterior a la fecha de inicio y guardo
    Entonces el sistema muestra una advertencia de validación
    Y no guarda los cambios

  Escenario: Cancelar la edición sin aplicar cambios
    Dado que existe un proyecto creado
    Y estoy en la vista de edición del proyecto
    Cuando decido no aplicar cambios y cancelo
    Entonces el proyecto permanece sin modificaciones
