# language: es
@HU-16
Característica: Gestión del equipo del proyecto
  Como Scrum Master
  Quiero quitar integrantes del equipo de un proyecto, y volver a sumarlos, conservando todo lo que aportaron
  Para mantener el equipo al día sin perder el historial del proyecto

  # Categorías: normal (US1.1), alternativo (US1.2, US1.3, US3.1 y US3.2),
  # límite (baja con una ronda abierta, reactivación con la ronda abierta),
  # error (US2.1, US3.3, US4.1 a US4.3). El Escenario 1 del backlog (alta como
  # Product Builder) lo cubre permisos_por_rol.feature (HU-21).
  # «Hoy» en los escenarios es el 2026-11-11.

  Antecedentes:
    Dado que inicié sesión como "ana"
    Y que "ana" creó el proyecto "Tienda"
    Y que "bruno" tiene cuenta y es Product Builder de "Tienda"
    Y que "Tienda" tiene la historia "Carrito"

  Escenario: US1.1 Quitar a un integrante conserva sus votos, sus horas y las métricas
    Dado que "bruno" votó 5 en la ronda revelada de "Carrito"
    Y que "bruno" imputó 2 horas en "Carrito" en un Sprint ya finalizado
    Y que "Tienda" tiene la historia "Pago"
    Y que anoté las métricas de "Tienda"
    Cuando quito a "bruno" del equipo de "Tienda"
    Entonces la baja se confirma sin contenido
    Y "bruno" no figura en el equipo de "Tienda"
    Y una votación nueva de "Pago" se revela solo con el voto de "ana"
    Y el historial de rondas de "Carrito" muestra el voto 5 de "bruno"
    Y las horas de "Carrito" muestran 2 horas de "bruno"
    Y las métricas de "Tienda" son las mismas que antes de la baja

  Esquema del escenario: US1.2 Quien fue quitado ya no actúa sobre el proyecto
    Dado que "ana" inició la votación de "Carrito"
    Y que "ana" quitó a "bruno" del equipo de "Tienda"
    Y que inicié sesión como "bruno"
    Cuando <accion>
    Entonces el sistema responde 403 con el mensaje "No sos integrante de este proyecto"

    Ejemplos:
      | accion                                 |
      | abro el proyecto "Tienda"              |
      | intento votar 5 en "Carrito"           |
      | intento registrar 2 horas en "Carrito" |

  Escenario: US1.3 Quien fue quitado ya no ve el proyecto en «Mis proyectos»
    Dado que "ana" quitó a "bruno" del equipo de "Tienda"
    Y que inicié sesión como "bruno"
    Cuando consulto mis proyectos
    Entonces recibo una lista vacía

  Escenario: Límite: quitar a quien ya votó en una ronda abierta descarta su voto
    Dado que "ana" inició la votación de "Carrito"
    Y que "bruno" ya votó 8 en "Carrito"
    Y que "ana" ya votó 3 en "Carrito"
    Cuando quito a "bruno" del equipo de "Tienda"
    Y "ana" revela la votación de "Carrito"
    Entonces la votación revelada de "Carrito" muestra solo el voto 3 de "ana"

  Escenario: Límite: quien vuelve con la ronda todavía abierta puede votar otra vez
    Dado que "ana" inició la votación de "Carrito"
    Y que "bruno" ya votó 8 en "Carrito"
    Y que "ana" quitó a "bruno" del equipo de "Tienda"
    Cuando agrego a "bruno" al equipo de "Tienda"
    Entonces "bruno" puede votar 5 en "Carrito"

  Escenario: US3.1 y US3.2 Volver a sumar a quien fue quitado lo reactiva sin duplicarlo
    Dado que "ana" quitó a "bruno" del equipo de "Tienda"
    Cuando agrego a "bruno" al equipo de "Tienda"
    Entonces "bruno" figura una sola vez en el equipo de "Tienda" como "product_builder"
    Y "bruno" ve solo "Tienda" en sus proyectos con el rol "product_builder"

  Escenario: US3.3 Agregar de nuevo a un integrante activo
    Cuando intento agregar a "bruno" al equipo de "Tienda"
    Entonces el sistema rechaza la petición con 400 y el mensaje "el integrante ya pertenece al proyecto"
    Y no cambia nada

  Escenario: US2.1 Al Scrum Master no se lo quita ni se le cambia el rol
    Cuando intento quitar a "ana" del equipo de "Tienda"
    Entonces el sistema rechaza la petición con 409 y el mensaje "No se puede quitar al Scrum Master del proyecto"
    Y no cambia nada
    Cuando intento agregar a "ana" al equipo de "Tienda"
    Entonces el sistema rechaza la petición con 400 y el mensaje "el integrante ya pertenece al proyecto"
    Y no cambia nada

  Escenario: US4.1 Proyecto inexistente
    Cuando intento quitar a "bruno" de un proyecto que no existe
    Entonces el sistema rechaza la petición con 404 y el mensaje "proyecto no encontrado"

  Escenario: US4.2 Usuario que nunca integró el proyecto
    Dado que "diego" tiene cuenta y no integra "Tienda"
    Cuando intento quitar a "diego" del equipo de "Tienda"
    Entonces el sistema rechaza la petición con 404 y el mensaje "integrante no encontrado"
    Y no cambia nada

  Escenario: US4.3 Quitar dos veces a la misma persona
    Dado que "ana" quitó a "bruno" del equipo de "Tienda"
    Cuando intento quitar a "bruno" del equipo de "Tienda"
    Entonces el sistema rechaza la petición con 404 y el mensaje "integrante no encontrado"
    Y no cambia nada
