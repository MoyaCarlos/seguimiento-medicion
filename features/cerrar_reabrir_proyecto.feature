# language: es
@300
Característica: Cerrar y reabrir un proyecto (solo el Scrum Master)
  Como Scrum Master
  Quiero cerrar un proyecto terminado, y poder reabrirlo
  Para que sus datos queden fijos y el equipo no los modifique por error

  # Categorías: normal (US1.1, US3.1, US3.2), alternativo (US2.3: las lecturas
  # siguen), límite (US1.4 y US3.3: cerrar dos veces y reabrir uno abierto),
  # error (US1.3, US2.1, US2.2, US4.1 a US4.4).
  # «Hoy» en los escenarios es el 2026-11-11.

  Antecedentes:
    Dado que inicié sesión como "ana"
    Y que "ana" creó el proyecto "Tienda"
    Y que "bruno" tiene cuenta y es Product Builder de "Tienda"
    Y que "Tienda" tiene la historia "Carrito"

  Escenario: US1.1 El Scrum Master cierra un proyecto abierto
    Cuando cierro el proyecto "Tienda"
    Entonces "Tienda" queda "Finalizado" con la fecha de cierre "2026-11-11"

  Escenario: US1.3 No se cierra un proyecto con un Sprint activo
    Dado que "Tienda" tiene un Sprint activo
    Cuando intento cerrar el proyecto "Tienda"
    Entonces el sistema rechaza la petición con 409 y el mensaje "Cerrá el Sprint activo antes de cerrar el proyecto"
    Y "Tienda" sigue abierto

  Escenario: US1.4 Cerrar dos veces
    Dado que "ana" cerró el proyecto "Tienda"
    Cuando intento cerrar el proyecto "Tienda"
    Entonces el sistema rechaza la petición con 409 y el mensaje "El proyecto ya está cerrado"
    Y no cambia nada

  Esquema del escenario: US2.1 y US2.2 Un proyecto cerrado no acepta escrituras
    Dado que "ana" inició la votación de "Carrito"
    Y que "ana" cerró el proyecto "Tienda"
    Y que inicié sesión como "<persona>"
    Cuando intento <accion>
    Entonces el sistema rechaza la petición con 409 y el mensaje "El proyecto está cerrado"
    Y no cambia nada

    Ejemplos:
      | persona | accion                                   |
      | bruno   | crear la historia "Pago" en "Tienda"     |
      | bruno   | votar 5 en "Carrito"                     |
      | bruno   | registrar 2 horas en "Carrito"           |
      | ana     | crear un Sprint en "Tienda"              |
      | ana     | editar el proyecto "Tienda"              |
      | ana     | agregar a "carla" al equipo de "Tienda"  |

  Escenario: US2.3 Las lecturas de un proyecto cerrado siguen funcionando
    Dado que "ana" cerró el proyecto "Tienda"
    Y que inicié sesión como "bruno"
    Cuando consulto el backlog de "Tienda"
    Entonces veo la historia "Carrito" en el backlog

  Escenario: US3.1 y US3.2 El Scrum Master reabre y se vuelve a escribir
    Dado que "ana" cerró el proyecto "Tienda"
    Cuando reabro el proyecto "Tienda"
    Entonces "Tienda" queda abierto con el estado "Planificado"
    Y "bruno" puede crear la historia "Pago" en "Tienda"

  Escenario: US3.3 Reabrir un proyecto abierto
    Cuando intento reabrir el proyecto "Tienda"
    Entonces el sistema rechaza la petición con 409 y el mensaje "El proyecto no está cerrado"
    Y no cambia nada

  Esquema del escenario: US4.1 y US4.2 Solo el Scrum Master cierra y reabre
    Dado que "diego" tiene cuenta y no integra "Tienda"
    Y que "Tienda" está <situacion>
    Y que inicié sesión como "<persona>"
    Cuando intento <accion> el proyecto "Tienda"
    Entonces el sistema responde 403 con el mensaje "<mensaje>"
    Y no cambia nada

    Ejemplos:
      | persona | situacion | accion  | mensaje                                         |
      | bruno   | abierto   | cerrar  | Solo el Scrum Master puede realizar esta acción |
      | bruno   | cerrado   | reabrir | Solo el Scrum Master puede realizar esta acción |
      | diego   | abierto   | cerrar  | No sos integrante de este proyecto              |
      | diego   | cerrado   | reabrir | No sos integrante de este proyecto              |

  Esquema del escenario: US4.3 Sin sesión
    Dado que "Tienda" está <situacion>
    Y que no inicié sesión
    Cuando intento <accion> el proyecto "Tienda"
    Entonces el sistema rechaza la petición con 401
    Y no cambia nada

    Ejemplos:
      | situacion | accion  |
      | abierto   | cerrar  |
      | cerrado   | reabrir |

  Esquema del escenario: US4.4 Proyecto inexistente
    Cuando intento <accion> un proyecto que no existe
    Entonces el sistema rechaza la petición con 404

    Ejemplos:
      | accion  |
      | cerrar  |
      | reabrir |
