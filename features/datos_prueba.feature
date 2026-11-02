# language: es

@Issue-144
Característica: Cargar datos de prueba

  Como integrante del equipo
  Quiero cargar un proyecto de demostración con datos de ejemplo
  Para mostrar y probar la aplicación sin cargar todo a mano

  Escenario: Carga exitosa
    Dado que la carga de datos de prueba está habilitada
    Y que la base no tiene datos de prueba
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 201 con el resumen del proyecto 25
    Y el proyecto 25 tiene 10 historias, 3 Sprints y 3 integrantes

  Escenario: Datos de prueba ya cargados
    Dado que la carga de datos de prueba está habilitada
    Y que los datos de prueba ya fueron cargados
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 409 con el mensaje "los datos de prueba ya están cargados"
    Y la carga no agregó ni cambió ningún dato

  Escenario: Identificador 25 ocupado por otro proyecto
    Dado que la carga de datos de prueba está habilitada
    Y que otro proyecto ya usa el identificador 25
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 409 con el mensaje "el id 25 ya lo usa otro proyecto"
    Y la carga no agregó ni cambió ningún dato

  Escenario: Carga de datos de prueba deshabilitada
    Dado que la carga de datos de prueba está deshabilitada
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 404 como si la ruta no existiera
    Y la carga no agregó ni cambió ningún dato
