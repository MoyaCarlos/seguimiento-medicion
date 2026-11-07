# language: es

# Cobertura por categoría:
# - Normal: lista de los Sprints finalizados con sus datos y métricas (US1.1, US1.3) y
#   coincidencia con HU-03 (US1.4); detalle con historias, esfuerzo y defectos (US2.1, US2.2,
#   US2.4).
# - Alternativo: sin horas por Story Point (US1.5); historia arrastrada que se edita y se
#   completa en otro Sprint (US2.3); defecto detectado en un Sprint y resuelto en otro (US2.5)
#   o en el mismo (US2.6).
# - Límite: Sprints "Activo" y "Pendiente" excluidos (US1.2); el Sprint activo no cambia al
#   consultar (US2.7); proyecto sin Sprints finalizados (US3.1); Sprint finalizado sin
#   historias, horas ni defectos (US3.2).
# - Error: proyecto inexistente (US1.6); Sprint inexistente (US2.8); Sprint de otro proyecto
#   (caso límite de la spec); Sprint todavía no finalizado (US2.9).
# US4 es la pantalla: se cubre con pruebas de componente en web/.

@HU-12
Característica: Consulta de Sprints anteriores

  Como Scrum Master
  Quiero consultar el historial de Sprints ya cerrados con su resumen (Sprint Goal, historias completadas, métricas)
  Para poder comparar el desempeño del equipo a lo largo del proyecto

  Antecedentes:
    Dado un proyecto del historial con 4 horas por Story Point

  # US1 - Lista de Sprints cerrados

  Escenario: El historial lista los Sprints finalizados ordenados por número
    Dado el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Y el Sprint del historial "S2" con Sprint Goal "Segundo incremento" está finalizado
    Cuando consulto el historial de Sprints del proyecto
    Entonces el historial lista los Sprints "S1, S2" en ese orden
    Y cada Sprint del historial informa su número, Sprint Goal, fecha de inicio y fecha de fin

  Escenario: Los Sprints "Activo" y "Pendiente" no aparecen en el historial
    Dado el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Y el Sprint del historial "S2" está activo
    Y el Sprint del historial "S3" está pendiente
    Cuando consulto el historial de Sprints del proyecto
    Entonces el historial lista los Sprints "S1" en ese orden

  Escenario: Cada Sprint del historial informa sus métricas principales
    Dado el Sprint del historial "S1" tiene una historia completada de 5 Story Points
    Y el Sprint del historial "S1" tiene una historia arrastrada de 3 Story Points
    Y el Sprint del historial "S1" tiene 30 horas imputadas
    Y el Sprint del historial "S1" tiene 2 defectos detectados y 1 resuelto
    Y el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Cuando consulto el historial de Sprints del proyecto
    Entonces el Sprint "S1" del historial informa estas métricas:
      | Story Points planificados | 8     |
      | Story Points completados  | 5     |
      | historias completadas (%) | 50    |
      | horas estimadas           | 32    |
      | horas reales              | 30    |
      | desviación (h)            | -2    |
      | defectos detectados       | 2     |
      | defectos resueltos        | 1     |

  Escenario: Las métricas del historial son las mismas que informa HU-03
    Dado el Sprint del historial "S1" tiene una historia completada de 5 Story Points
    Y el Sprint del historial "S1" tiene una historia arrastrada de 3 Story Points
    Y el Sprint del historial "S1" tiene 30 horas imputadas
    Y el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Cuando consulto el historial de Sprints del proyecto
    Entonces las métricas de cada Sprint del historial coinciden con las del cálculo de métricas del proyecto

  Escenario: Sin horas por Story Point, las horas estimadas y la desviación quedan sin datos
    Dado el proyecto del historial no tiene horas por Story Point
    Y el Sprint del historial "S1" tiene una historia completada de 5 Story Points
    Y el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Cuando consulto el historial de Sprints del proyecto
    Entonces el Sprint "S1" del historial informa sin datos las horas estimadas y la desviación

  Escenario: Historial de un proyecto inexistente
    Cuando consulto el historial de Sprints de un proyecto inexistente
    Entonces la consulta del historial falla con "proyecto no encontrado"

  # US2 - Detalle de un Sprint cerrado

  Escenario: El detalle muestra las historias comprometidas con su resultado
    Dado el Sprint del historial "S1" tiene la historia completada "Login" de 5 Story Points
    Y el Sprint del historial "S1" tiene la historia arrastrada "Reportes" de 3 Story Points
    Y el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Cuando consulto el detalle del Sprint "S1" en el historial
    Entonces el detalle informa el número, el Sprint Goal "Primer incremento", las fechas y las métricas del Sprint
    Y el detalle informa estas historias comprometidas:
      | historia | Story Points | completada |
      | Login    | 5            | sí         |
      | Reportes | 3            | no         |

  Escenario: Una historia arrastrada conserva sus Story Points y su resultado en cada Sprint
    Dado el Sprint del historial "S1" tiene la historia arrastrada "Reportes" de 3 Story Points
    Y el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Y la historia "Reportes" del historial pasa a tener 8 Story Points
    Y la historia "Reportes" del historial se completa en el Sprint "S2"
    Y el Sprint del historial "S2" con Sprint Goal "Segundo incremento" está finalizado
    Cuando consulto el detalle del Sprint "S1" en el historial
    Entonces el detalle informa estas historias comprometidas:
      | historia | Story Points | completada |
      | Reportes | 3            | no         |
    Cuando consulto el detalle del Sprint "S2" en el historial
    Entonces el detalle informa estas historias comprometidas:
      | historia | Story Points | completada |
      | Reportes | 8            | sí         |

  Escenario: El esfuerzo se informa por historia y su total coincide con las horas reales
    Dado el Sprint del historial "S1" tiene la historia completada "Login" de 5 Story Points
    Y el Sprint del historial "S1" tiene la historia arrastrada "Reportes" de 3 Story Points
    Y la historia "Login" del historial tiene 90 minutos imputados en el Sprint "S1"
    Y la historia "Reportes" del historial tiene 50 minutos imputados en el Sprint "S1"
    Y el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Cuando consulto el detalle del Sprint "S1" en el historial
    Entonces el detalle informa este esfuerzo por historia:
      | historia | minutos | horas |
      | Login    | 90      | 1.5   |
      | Reportes | 50      | 0.83  |
    Y la suma de los minutos del detalle coincide con las horas reales del Sprint

  Escenario: Un defecto detectado en un Sprint y resuelto en otro aparece en ambos detalles
    Dado un defecto del historial "Error al guardar" detectado en el Sprint "S1" y resuelto en el Sprint "S2"
    Y el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Y el Sprint del historial "S2" con Sprint Goal "Segundo incremento" está finalizado
    Cuando consulto el detalle del Sprint "S1" en el historial
    Entonces el detalle informa estos defectos:
      | defecto          | detectado | resuelto |
      | Error al guardar | sí        | no       |
    Cuando consulto el detalle del Sprint "S2" en el historial
    Entonces el detalle informa estos defectos:
      | defecto          | detectado | resuelto |
      | Error al guardar | no        | sí       |

  Escenario: Un defecto detectado y resuelto en el mismo Sprint aparece una sola vez
    Dado un defecto del historial "Error al guardar" detectado en el Sprint "S1" y resuelto en el Sprint "S1"
    Y el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Cuando consulto el detalle del Sprint "S1" en el historial
    Entonces el detalle informa estos defectos:
      | defecto          | detectado | resuelto |
      | Error al guardar | sí        | sí       |

  Escenario: Consultar el detalle no modifica el Sprint activo
    Dado el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Y el Sprint del historial "S2" está activo con historias, horas y defectos
    Cuando consulto el detalle del Sprint "S1" en el historial
    Entonces el Sprint "S2" del historial, sus historias, sus horas y sus defectos no cambiaron

  Escenario: Detalle de un Sprint inexistente
    Cuando consulto el detalle de un Sprint inexistente en el historial
    Entonces la consulta del historial falla con "sprint no encontrado"

  Escenario: Detalle de un Sprint de otro proyecto
    Dado otro proyecto del historial con un Sprint finalizado
    Cuando consulto en el historial del proyecto el detalle del Sprint del otro proyecto
    Entonces la consulta del historial falla con "sprint no encontrado"

  Esquema del escenario: Detalle de un Sprint que todavía no está finalizado
    Dado el Sprint del historial "S1" está <estado>
    Cuando consulto el detalle del Sprint "S1" en el historial
    Entonces la consulta del historial falla en el campo "estado" con "el Sprint todavía no está finalizado"

    Ejemplos:
      | estado    |
      | activo    |
      | pendiente |

  # US3 - Proyecto sin Sprints cerrados

  Escenario: Un proyecto sin Sprints finalizados tiene el historial vacío
    Dado el Sprint del historial "S1" está activo
    Cuando consulto el historial de Sprints del proyecto
    Entonces el historial no lista ningún Sprint

  Escenario: Un Sprint finalizado sin historias, horas ni defectos
    Dado el Sprint del historial "S1" con Sprint Goal "Primer incremento" está finalizado
    Cuando consulto el detalle del Sprint "S1" en el historial
    Entonces el detalle no informa historias ni defectos
    Y el detalle informa las métricas del cálculo de métricas para ese Sprint
