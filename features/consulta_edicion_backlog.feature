# language: es

# Cobertura por categoría:
# - Normal: listado del Product Backlog; edición válida.
# - Alternativo: backlog vacío; carga inicial de Story Points en un Sprint "Activo" o "Finalizado".
# - Límite: Story Points fuera de la escala Fibonacci; Story Points ya cargados en un Sprint "Activo" o "Finalizado".
# - Error: proyecto inexistente; título en blanco; historia inexistente (el 404 domina sobre la validación).

@HU-15
Característica: Consulta y edición del Product Backlog

  Como Product Builder
  Quiero ver el Product Backlog de un proyecto ordenado por prioridad y poder editar el título, la descripción, los criterios de aceptación y los Story Points de cada historia
  Para poder mantenerlo actualizado y cumplir con el contenido mínimo que pide el enunciado

  Antecedentes:
    Dado que existe un proyecto con identificador 1

  Escenario: Listado del Product Backlog ordenado por prioridad y creación
    Dado que el proyecto tiene estas historias, en este orden de creación:
      | titulo           | prioridad | en_sprint |
      | Cerrar sesión    | S         | no        |
      | Iniciar sesión   | M         | no        |
      | Exportar reporte | M         | si        |
      | Modo oscuro      | W         | no        |
      | Ayuda en línea   | C         | no        |
      | Recuperar clave  | M         | no        |
    Cuando consulto el Product Backlog del proyecto
    Entonces veo las historias en este orden:
      | titulo          |
      | Iniciar sesión  |
      | Recuperar clave |
      | Cerrar sesión   |
      | Ayuda en línea  |
      | Modo oscuro     |

  Escenario: Product Backlog vacío
    Dado que el proyecto no tiene historias sin Sprint asignado
    Cuando consulto el Product Backlog del proyecto
    Entonces veo un Product Backlog vacío, sin error

  Escenario: Consulta del Product Backlog de un proyecto inexistente
    Cuando consulto el Product Backlog del proyecto con identificador 999, que no existe
    Entonces la consulta se rechaza con error "proyecto no encontrado"

  Escenario: Edición válida de una historia
    Dado que existe una historia en el Product Backlog del proyecto
    Cuando cambio su título a "Iniciar sesión con Google", su descripción, sus criterios de aceptación y sus Story Points a 5, y guardo
    Entonces la historia queda guardada con los datos nuevos
    Y conserva su identificador y su estado "Nueva"

  Escenario: Edición con el título en blanco
    Dado que existe una historia en el Product Backlog del proyecto
    Cuando cambio su título a "   " y guardo
    Entonces la edición se rechaza con una advertencia sobre el campo "titulo"
    Y la historia queda sin modificaciones

  Escenario: Edición con Story Points fuera de la escala Fibonacci
    Dado que existe una historia en el Product Backlog del proyecto
    Cuando cambio sus Story Points a 4 y guardo
    Entonces la edición se rechaza con una advertencia sobre el campo "estimacion_sp"
    Y la historia queda sin modificaciones

  Escenario: Edición de una historia inexistente
    Cuando intento editar la historia con identificador 999, que no existe, dejando su título en blanco
    Entonces la edición se rechaza con error "historia no encontrada"

  Esquema del escenario: No se cambian los Story Points ya cargados en un Sprint "<estado>"
    Dado que existe una historia con 3 Story Points asignada a un Sprint en estado "<estado>"
    Cuando cambio sus Story Points a 5 y guardo
    Entonces la edición se rechaza con una advertencia sobre el campo "estimacion_sp"
    Y la historia queda sin modificaciones

    Ejemplos:
      | estado     |
      | Activo     |
      | Finalizado |

  Esquema del escenario: Carga inicial de Story Points en un Sprint "<estado>"
    Dado que existe una historia sin estimar asignada a un Sprint en estado "<estado>"
    Cuando cambio sus Story Points a 5 y guardo
    Entonces la historia queda guardada con 5 Story Points

    Ejemplos:
      | estado     |
      | Activo     |
      | Finalizado |
