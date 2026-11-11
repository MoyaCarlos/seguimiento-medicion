# language: es
@HU-21
Característica: Permisos por rol (segregación de funciones)
  Como Scrum Master
  Quiero que las acciones de gestión del proyecto solo las pueda hacer el Scrum Master
  Para que el resto del equipo participe sin alterar la planificación

  # Categorías: normal (1, 5), alternativo (8, 9), límite (6, 7), error (2, 3, 4).

  Antecedentes:
    Dado que inicié sesión como "ana"
    Y que "ana" creó el proyecto "Tienda"
    Y que "bruno" tiene cuenta y es Product Builder de "Tienda"
    Y que "Tienda" tiene la historia "Carrito"

  Escenario: 1. Acción permitida para el Scrum Master
    Dado que inicié sesión como "ana"
    Cuando creo un Sprint en "Tienda"
    E inicio la votación de "Carrito"
    Y edito el proyecto "Tienda"
    Entonces cada acción se realiza

  Esquema del escenario: 2. Acción no permitida para un Product Builder
    Dado que inicié sesión como "bruno"
    Cuando intento <accion>
    Entonces el sistema responde 403 con un mensaje que empieza con "Solo el Scrum Master puede"
    Y no cambia nada

    Ejemplos:
      | accion                                     |
      | crear un Sprint en "Tienda"                |
      | iniciar la votación de "Carrito"           |
      | editar el proyecto "Tienda"                |
      | agregar a "carla" al equipo de "Tienda"    |

  Escenario: 3. Sin sesión
    Dado que no inicié sesión
    Cuando consulto mis proyectos
    Entonces el sistema rechaza la petición con 401

  Escenario: 4. No integrante
    Dado que "diego" tiene cuenta y no integra "Tienda"
    Y que inicié sesión como "diego"
    Cuando abro el proyecto "Tienda"
    Entonces el sistema responde 403 con el mensaje "No sos integrante de este proyecto"

  Escenario: 5. Mis proyectos con mi rol
    Dado que inicié sesión como "bruno"
    Cuando consulto mis proyectos
    Entonces veo solo "Tienda" con el rol "product_builder"

  Escenario: 6. Un solo Scrum Master
    Dado que inicié sesión como "ana"
    Y que "carla" tiene cuenta
    Cuando agrego a "carla" al equipo de "Tienda" como "scrum_master"
    Entonces el sistema rechaza la petición con 409
    Y "Tienda" sigue teniendo a "ana" como único Scrum Master

  Escenario: 7. Identidad de la sesión
    Dado que "ana" inició la votación de "Carrito"
    Y que inicié sesión como "bruno"
    Y voto 8 en "Carrito" indicando el usuario de "ana"
    Cuando "ana" consulta la votación de "Carrito" indicando el usuario de "bruno"
    Entonces el voto de 8 queda a nombre de "bruno"
    Y "ana" no ve el voto oculto de "bruno"

  Escenario: 8. Configurar el equipo solo con usuarios registrados
    Dado que inicié sesión como "ana"
    Y que "carla" tiene cuenta
    Cuando agrego a "carla" al equipo de "Tienda"
    Entonces "carla" queda como "product_builder" de "Tienda"
    Y si agrego a "nadie", que no tiene cuenta, el sistema lo rechaza

  Escenario: 9. Usuario sin proyectos
    Dado que "eva" tiene cuenta y no integra ningún proyecto
    Y que inicié sesión como "eva"
    Cuando consulto mis proyectos
    Entonces recibo una lista vacía
