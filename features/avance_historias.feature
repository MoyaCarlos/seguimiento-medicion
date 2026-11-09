# language: es

# Cobertura por categoría:
# - Normal: historias agrupadas en Nueva, En progreso y Completada con la cantidad y los Story
#   Points de cada grupo (US1.1, US1.2); Sprint de cada historia (US2.1).
# - Alternativo: historia sin Story Points, que cuenta y suma 0 (US1.3); historia sin Sprint
#   en el Product Backlog (US2.2); historia Nueva asignada a un Sprint no iniciado, que
#   conserva ese Sprint (Clarification 2); historia arrastrada al cerrar un Sprint, que vuelve
#   a Nueva como "Product Backlog" (Clarification 1).
# - Límite: los tres grupos aparecen aunque estén vacíos (US1.4); proyecto sin historias con
#   los totales en cero (US3.1).
# - Error: proyecto inexistente, 404 (US3.2); identificador con formato inválido, 400 con
#   campo "id" (US3.3).
# La pantalla "Estado" se cubre con pruebas de componente en web/.

@HU-18
Característica: Avance de las historias del proyecto

  Como integrante del equipo
  Quiero ver todas las historias del proyecto agrupadas por estado y por Sprint
  Para saber qué se terminó y qué falta sin revisar los Sprints uno por uno

  Antecedentes:
    Dado un proyecto para el avance de historias

  # US1 - Historias agrupadas por estado

  Escenario: Las historias se agrupan en Nueva, En progreso y Completada con su cantidad y Story Points
    Dado la historia de avance "H1" de 3 Story Points está Nueva
    Y la historia de avance "H2" de 5 Story Points está En progreso en el Sprint "S1"
    Y la historia de avance "H3" de 8 Story Points está Completada en el Sprint "S1"
    Y la historia de avance "H4" de 2 Story Points está Completada en el Sprint "S1"
    Cuando consulto el avance de las historias del proyecto
    Entonces el avance informa estos grupos en este orden:
      | estado      | historias | Story Points |
      | Nueva       | 1         | 3            |
      | En progreso | 1         | 5            |
      | Completada  | 2         | 10           |
    Y el avance informa un total de 4 historias y 18 Story Points

  Escenario: Una historia sin Story Points cuenta en la cantidad y suma cero
    Dado la historia de avance "H1" de 3 Story Points está Nueva
    Y la historia de avance "H2" sin estimar está Nueva
    Cuando consulto el avance de las historias del proyecto
    Entonces el grupo "Nueva" informa 2 historias y 3 Story Points
    Y la historia "H2" figura como "sin estimar"

  Escenario: Los tres grupos aparecen aunque estén vacíos
    Dado la historia de avance "H1" de 3 Story Points está Nueva
    Cuando consulto el avance de las historias del proyecto
    Entonces el avance informa estos grupos en este orden:
      | estado      | historias | Story Points |
      | Nueva       | 1         | 3            |
      | En progreso | 0         | 0            |
      | Completada  | 0         | 0            |

  # US2 - Sprint de cada historia

  Escenario: Cada historia muestra el Sprint al que pertenece
    Dado la historia de avance "H1" de 5 Story Points está En progreso en el Sprint "S1"
    Cuando consulto el avance de las historias del proyecto
    Entonces la historia "H1" figura en el Sprint "S1"

  Escenario: Una historia sin Sprint figura en el Product Backlog
    Dado la historia de avance "H1" de 3 Story Points está Nueva
    Cuando consulto el avance de las historias del proyecto
    Entonces la historia "H1" figura en el "Product Backlog"

  Escenario: Una historia Nueva asignada a un Sprint no iniciado conserva ese Sprint
    Dado la historia de avance "H1" de 3 Story Points está Nueva y asignada al Sprint pendiente "S2"
    Cuando consulto el avance de las historias del proyecto
    Entonces el grupo "Nueva" informa 1 historias y 3 Story Points
    Y la historia "H1" figura en el Sprint "S2"

  Escenario: Una historia arrastrada al cerrar un Sprint vuelve a Nueva en el Product Backlog
    Dado la historia de avance "H1" de 5 Story Points estaba En progreso en el Sprint "S1"
    Cuando se cierra el Sprint "S1" del avance
    Y consulto el avance de las historias del proyecto
    Entonces el grupo "Nueva" informa 1 historias y 5 Story Points
    Y la historia "H1" figura en el "Product Backlog"

  # US3 - Proyecto sin historias, inexistente o con identificador inválido

  Escenario: Un proyecto sin historias informa los totales en cero
    Cuando consulto el avance de las historias del proyecto
    Entonces el avance informa estos grupos en este orden:
      | estado      | historias | Story Points |
      | Nueva       | 0         | 0            |
      | En progreso | 0         | 0            |
      | Completada  | 0         | 0            |
    Y el avance informa un total de 0 historias y 0 Story Points

  Escenario: Un proyecto inexistente responde "no encontrado"
    Cuando consulto el avance de las historias del proyecto 9999
    Entonces la respuesta del avance es 404 "proyecto no encontrado"

  Escenario: Un identificador con formato inválido se rechaza con un error sobre el campo "id"
    Cuando consulto el avance de las historias del proyecto "abc"
    Entonces la respuesta del avance es 400 con el campo "id"
