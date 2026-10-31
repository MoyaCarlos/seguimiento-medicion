# language: es

# Cobertura por categoría:
# - Normal: Nueva → En progreso (US1.1); En progreso → Completada con fecha de finalización (US1.2).
# - Alternativo: la historia completada aporta sus Story Points como completados (US1.3).
# - Límite: la historia completada no vuelve al Product Backlog al cerrar el Sprint (US1.4); el cambio se conserva al volver a consultar (US1.5).
# - Error: salto Nueva → Completada (US1.6); Sprint Pendiente o Finalizado, Esquema del escenario con dos ejemplos (US1.7); retroceso de estado (US1.8); estado repetido (US1.9);
#   historia sin Sprint a En progreso (US2.1) y a Completada (US2.2); historia inexistente (US2.3).
# Nota: US1.10 (dos cambios simultáneos sobre la misma historia) no es determinista en Godog;
# se cubre en los tests de repositorio (sqlite_backlog_test.go) y de handler (backlog_handler_test.go).

@HU-11
Característica: Estado de una historia en el Sprint

  Como Product Builder
  Quiero cambiar el estado de una historia asignada al Sprint (En progreso, Completada)
  Para poder reflejar el avance real del trabajo durante la iteración

  Antecedentes:
    Dado que existe un proyecto con identificador 1

  Escenario: Una historia nueva del Sprint activo pasa a En progreso
    Dado que el tablero del proyecto tiene un Sprint activo con la historia "Iniciar sesión" en estado "Nueva"
    Cuando pongo la historia "Iniciar sesión" del tablero en estado "En progreso"
    Entonces la historia "Iniciar sesión" del tablero figura en estado "En progreso"
    Y la historia "Iniciar sesión" del tablero sigue en el Sprint activo con todos sus datos

  Escenario: Una historia en progreso se marca como Completada y registra su fecha de finalización
    Dado que el tablero del proyecto tiene un Sprint activo con la historia "Iniciar sesión" en estado "En progreso"
    Cuando pongo la historia "Iniciar sesión" del tablero en estado "Completada"
    Entonces la historia "Iniciar sesión" del tablero figura en estado "Completada"
    Y la historia "Iniciar sesión" del tablero tiene registrada su fecha de finalización

  Escenario: La historia completada aporta sus Story Points como completados
    Dado que el tablero del proyecto tiene un Sprint activo con la historia "Iniciar sesión" de 5 Story Points en estado "En progreso"
    Cuando pongo la historia "Iniciar sesión" del tablero en estado "Completada"
    Entonces el Sprint activo del tablero cuenta 5 Story Points completados

  Escenario: La historia completada no vuelve al Product Backlog al cerrar el Sprint
    Dado que el tablero del proyecto tiene un Sprint activo con la historia "Iniciar sesión" en estado "En progreso"
    Y que pongo la historia "Iniciar sesión" del tablero en estado "Completada"
    Cuando se cierra el Sprint activo del tablero
    Entonces la historia "Iniciar sesión" del tablero permanece en el Sprint cerrado
    Y la historia "Iniciar sesión" del tablero no aparece en el Product Backlog del proyecto

  Escenario: El cambio de estado se conserva al volver a consultar la historia
    Dado que el tablero del proyecto tiene un Sprint activo con la historia "Iniciar sesión" en estado "Nueva"
    Y que pongo la historia "Iniciar sesión" del tablero en estado "En progreso"
    Cuando consulto de nuevo la historia "Iniciar sesión" del tablero
    Entonces la historia "Iniciar sesión" del tablero figura en estado "En progreso"

  Escenario: No se puede completar directamente una historia nueva
    Dado que el tablero del proyecto tiene un Sprint activo con la historia "Iniciar sesión" en estado "Nueva"
    Cuando pongo la historia "Iniciar sesión" del tablero en estado "Completada"
    Entonces el cambio de estado del tablero se rechaza con el mensaje "primero debe pasar por En progreso"
    Y la historia "Iniciar sesión" del tablero figura en estado "Nueva"

  Esquema del escenario: Solo se cambia el estado de historias de un Sprint activo
    Dado que el tablero del proyecto tiene un Sprint <estado_sprint> con la historia "Iniciar sesión" en estado "Nueva"
    Cuando pongo la historia "Iniciar sesión" del tablero en estado "En progreso"
    Entonces el cambio de estado del tablero se rechaza con el mensaje "solo se cambia el estado de historias de un Sprint activo"
    Y la historia "Iniciar sesión" del tablero figura en estado "Nueva"

    Ejemplos:
      | estado_sprint |
      | Pendiente     |
      | Finalizado    |

  Escenario: No se puede volver a un estado anterior
    Dado que el tablero del proyecto tiene un Sprint activo con la historia "Iniciar sesión" en estado "En progreso"
    Cuando pongo la historia "Iniciar sesión" del tablero en estado "Nueva"
    Entonces el cambio de estado del tablero se rechaza con el mensaje "el estado de una historia solo avanza: no se puede volver a un estado anterior"
    Y la historia "Iniciar sesión" del tablero figura en estado "En progreso"

  Escenario: No se puede pedir el estado que la historia ya tiene
    Dado que el tablero del proyecto tiene un Sprint activo con la historia "Iniciar sesión" en estado "En progreso"
    Cuando pongo la historia "Iniciar sesión" del tablero en estado "En progreso"
    Entonces el cambio de estado del tablero se rechaza con el mensaje "la historia ya está en ese estado"
    Y la historia "Iniciar sesión" del tablero figura en estado "En progreso"

  Escenario: Una historia sin Sprint no puede pasar a En progreso
    Dado que el tablero del proyecto tiene la historia "Iniciar sesión" sin Sprint
    Cuando pongo la historia "Iniciar sesión" del tablero en estado "En progreso"
    Entonces el cambio de estado del tablero se rechaza con el mensaje "primero debe asignarse a un Sprint activo"
    Y la historia "Iniciar sesión" del tablero figura en estado "Nueva"

  Escenario: Una historia sin Sprint no puede pasar a Completada
    Dado que el tablero del proyecto tiene la historia "Iniciar sesión" sin Sprint
    Cuando pongo la historia "Iniciar sesión" del tablero en estado "Completada"
    Entonces el cambio de estado del tablero se rechaza con el mensaje "primero debe asignarse a un Sprint activo"
    Y la historia "Iniciar sesión" del tablero figura en estado "Nueva"

  Escenario: No se puede cambiar el estado de una historia inexistente
    Cuando pongo la historia con identificador 999 del tablero, que no existe, en estado "En progreso"
    Entonces el cambio de estado del tablero se rechaza con el error "historia no encontrada"
