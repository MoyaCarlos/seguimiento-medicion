# language: es

# Cobertura por categoría:
# - Normal: imputación válida de 4 horas; consulta del esfuerzo con varias cargas.
# - Alternativo: cargas de dos integrantes que se suman; horas que siguen en su Sprint
#   después de que la historia pasa a otro; historia sin cargas.
# - Límite: cero horas. Los bordes del tope de 24 horas y del rango de fechas no son
#   criterios de aceptación: se prueban con tests unitarios del dominio.
# - Error: horas negativas o no numéricas, actividad vacía, integrante ajeno al proyecto,
#   historia fuera del Sprint activo e historia inexistente (al imputar y al consultar).

@HU-07
Característica: Imputación de horas trabajadas

  Como integrante del equipo
  Quiero registrar las horas reales que invertí trabajando en una Historia de Usuario
  Para poder proveer datos al sistema que permitan comparar el esfuerzo estimado contra el esfuerzo real

  Antecedentes:
    Dado que existe un proyecto con identificador 1
    Y que "Ana" es integrante del proyecto
    Y que el proyecto tiene un Sprint activo
    Y que existe una historia asignada al Sprint activo

  Escenario: Registro válido de tiempo
    Cuando "Ana" imputa 4 horas a la historia con la fecha de hoy y la actividad "Implementar el login"
    Entonces la carga queda registrada con integrante "Ana", la fecha de hoy, la actividad "Implementar el login", 240 minutos y el Sprint activo
    Y el esfuerzo acumulado de la historia es de 4 horas

  Escenario: Las cargas de distintos integrantes se suman
    Dado que "Jimena" es integrante del proyecto
    Y que "Ana" imputó 1.5 horas a la historia
    Cuando "Jimena" imputa 2.5 horas a la historia con la fecha de hoy y la actividad "Revisar el PR"
    Entonces el esfuerzo acumulado de la historia es de 4 horas

  Escenario: Las horas siguen atribuidas al Sprint en que se cargaron
    Dado que "Ana" imputó 3 horas a la historia
    Cuando se cierra el Sprint activo y la historia se planifica en un Sprint nuevo
    Entonces las 3 horas siguen atribuidas al Sprint cerrado

  Esquema del escenario: Horas negativas o cero
    Cuando "Ana" intenta imputar "<horas>" horas a la historia
    Entonces la carga se rechaza con un error sobre el campo "horas"
    Y el esfuerzo acumulado de la historia no cambia

    Ejemplos:
      | horas |
      | -2    |
      | 0     |

  Escenario: Horas no numéricas
    Cuando "Ana" intenta imputar "cuatro" horas a la historia
    Entonces la carga se rechaza con un error sobre el campo "horas"
    Y el esfuerzo acumulado de la historia no cambia

  Esquema del escenario: Actividad vacía o compuesta solo por espacios
    Cuando "Ana" intenta imputar horas a la historia con una actividad <actividad>
    Entonces la carga se rechaza con un error sobre el campo "actividad"
    Y el esfuerzo acumulado de la historia no cambia

    Ejemplos:
      | actividad         |
      | vacía             |
      | solo con espacios |

  Escenario: Integrante que no pertenece al proyecto
    Cuando "Pedro", que no es integrante del proyecto, intenta imputar horas a la historia
    Entonces la carga se rechaza con un error sobre el campo "integrante_id"
    Y el esfuerzo acumulado de la historia no cambia

  Esquema del escenario: Historia fuera del Sprint activo
    Dado que existe otra historia <ubicacion>
    Cuando "Ana" intenta imputar horas a esa otra historia
    Entonces la carga se rechaza con el mensaje "solo se imputan horas a historias del Sprint activo"
    Y esa otra historia no tiene horas imputadas

    Ejemplos:
      | ubicacion                            |
      | en el Product Backlog                |
      | en un Sprint en estado "Pendiente"   |
      | en un Sprint en estado "Finalizado"  |

  Escenario: Imputación a una historia inexistente
    Cuando "Ana" intenta imputar horas a la historia con identificador 999
    Entonces la carga se rechaza con error "historia no encontrada"

  Escenario: Consulta del esfuerzo real con varias cargas
    Dado que "Jimena" es integrante del proyecto
    Y que la historia tiene estas cargas:
      | integrante | horas | actividad      |
      | Ana        | 0.5   | Análisis       |
      | Jimena     | 1.5   | Implementación |
      | Ana        | 2     | Pruebas        |
    Cuando consulto el esfuerzo real de la historia
    Entonces el total es de 4 horas
    Y veo 3 cargas con integrante, fecha, actividad, duración y Sprint

  Escenario: Consulta del esfuerzo real de una historia sin cargas
    Cuando consulto el esfuerzo real de la historia
    Entonces el total es de 0 horas
    Y no veo ninguna carga

  Escenario: Consulta del esfuerzo real de una historia inexistente
    Cuando consulto el esfuerzo real de la historia con identificador 999
    Entonces la consulta del esfuerzo se rechaza con error "historia no encontrada"
