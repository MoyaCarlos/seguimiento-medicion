# language: es

@Issue-144 @Issue-235
Característica: Cargar datos de prueba

  Como integrante del equipo
  Quiero cargar un proyecto de demostración con datos de ejemplo y usuarios con contraseña conocida
  Para mostrar y probar la aplicación sin cargar todo a mano

  Escenario: Carga exitosa
    Dado que la carga de datos de prueba está habilitada
    Y que la base no tiene datos de prueba
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 201 con el resumen del proyecto 25
    Y el proyecto 25 tiene 12 historias, 4 Sprints y 4 integrantes
    Y el resumen informa 4 usuarios, 4 integrantes, 4 Sprints, 12 historias, 20 registros de horas, 2 rondas de votación, 6 votos y 5 defectos
    Y el proyecto 25 tiene 5 defectos, 3 Resueltos y 2 Pendientes
    Y "carlos.moya" es el único Scrum Master del proyecto 25

  Esquema del escenario: Los usuarios de la demo ingresan con la contraseña de demostración
    Dado que la carga de datos de prueba está habilitada
    Y que los datos de prueba ya fueron cargados
    Cuando ingreso con el usuario de la demo "<usuario>" y la contraseña "demo1234"
    Entonces el sistema abre la sesión de "<usuario>"

    Ejemplos:
      | usuario         |
      | carlos.moya     |
      | ana.iriarte     |
      | candela.vulcano |
      | jimena.martinez |

  Escenario: La contraseña de la demo no se guarda en claro
    Dado que la carga de datos de prueba está habilitada
    Y que los datos de prueba ya fueron cargados
    Entonces ningún usuario de la demo guarda su contraseña en claro
    Y no existen usuarios «Beta»

  Escenario: Datos de prueba ya cargados
    Dado que la carga de datos de prueba está habilitada
    Y que los datos de prueba ya fueron cargados
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 409 con el mensaje "los datos de prueba ya están cargados"
    Y el mensaje no nombra a ningún usuario
    Y la carga no agregó ni cambió ningún dato

  Escenario: Un usuario de la demo ya existe
    Dado que la carga de datos de prueba está habilitada
    Y que el usuario "ana.iriarte" ya existe
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 409 con el mensaje "ya existe el usuario ana.iriarte"
    Y la carga no agregó ni cambió ningún dato

  Escenario: Varios usuarios de la demo existen y se nombra al primero de la lista
    Dado que la carga de datos de prueba está habilitada
    Y que el usuario "jimena.martinez" ya existe
    Y que el usuario "candela.vulcano" ya existe
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 409 con el mensaje "ya existe el usuario candela.vulcano"
    Y la carga no agregó ni cambió ningún dato

  Escenario: El conflicto de usuario ignora las mayúsculas
    Dado que la carga de datos de prueba está habilitada
    Y que el usuario "Ana.Iriarte" ya existe
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 409 con el mensaje "ya existe el usuario ana.iriarte"
    Y la carga no agregó ni cambió ningún dato

  Escenario: Un usuario existente gana al identificador 25 ocupado
    Dado que la carga de datos de prueba está habilitada
    Y que el usuario "carlos.moya" ya existe
    Y que otro proyecto ya usa el identificador 25
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 409 con el mensaje "ya existe el usuario carlos.moya"
    Y la carga no agregó ni cambió ningún dato

  Escenario: Identificador 25 ocupado por otro proyecto
    Dado que la carga de datos de prueba está habilitada
    Y que otro proyecto ya usa el identificador 25
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 409 con el mensaje "el id 25 ya lo usa otro proyecto"
    Y la carga no agregó ni cambió ningún dato

  Escenario: Una falla a mitad de la carga no deja ningún dato
    Dado que la carga de datos de prueba está habilitada
    Y que la base rechaza los defectos de severidad Baja
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 500 sin filtrar el detalle del error
    Y la carga no agregó ni cambió ningún dato

  Escenario: Carga de datos de prueba deshabilitada
    Dado que la carga de datos de prueba está deshabilitada
    Cuando pido cargar los datos de prueba
    Entonces el sistema responde 404 como si la ruta no existiera
    Y la carga no agregó ni cambió ningún dato
