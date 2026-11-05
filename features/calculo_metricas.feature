# language: es

# Cobertura por categoría:
# - Normal: métricas de un Sprint finalizado y del proyecto (US1.1–US1.5); definir el
#   parámetro (US3.1).
# - Alternativo: sin parámetro (US3.2); conservar y borrar el parámetro (US3.3, US3.4);
#   recálculo al cambiarlo (US3.6); editar una historia sin cargar Story Points (US5.2);
#   primera carga con Sprint activo (US5.4) y en una historia arrastrada (US5.5).
# - Límite: Sprint activo excluido (US1.6); historia arrastrada que se edita (US1.7) y que
#   cuenta en dos Sprints (US1.9); proyecto sin Sprints finalizados (US2.1, US2.2).
# - Error: proyecto inexistente (US1.8); parámetro inválido (US3.5); carga rechazada con
#   Sprint finalizado, por edición y por votación (US5.1, US5.3).
# US4 es la pantalla: se cubre con pruebas de componente en web/.

@HU-03
Característica: Cálculo de métricas del proyecto

  Como Scrum Master
  Quiero que el sistema calcule automáticamente las métricas del proyecto y de cada Sprint
  Para poder analizar la capacidad real del equipo y la calidad del trabajo entregado

  Escenario: Story Points, velocidad y porcentaje de un Sprint finalizado
    Dado un proyecto de métricas con 4 horas por Story Point
    Y el Sprint de métricas "S1" tiene una historia completada de 5 Story Points
    Y el Sprint de métricas "S1" tiene una historia arrastrada de 3 Story Points
    Y el Sprint de métricas "S1" está finalizado
    Cuando consulto las métricas del proyecto
    Entonces el Sprint "S1" informa estas métricas:
      | Story Points planificados | 8  |
      | Story Points completados  | 5  |
      | velocidad                 | 5  |
      | historias completadas (%) | 50 |

  Escenario: Horas y desviación de un Sprint finalizado
    Dado un proyecto de métricas con 4 horas por Story Point
    Y el Sprint de métricas "S1" tiene una historia completada de 5 Story Points
    Y el Sprint de métricas "S1" tiene una historia arrastrada de 3 Story Points
    Y el Sprint de métricas "S1" tiene 30 horas imputadas
    Y el Sprint de métricas "S1" está finalizado
    Cuando consulto las métricas del proyecto
    Entonces el Sprint "S1" informa estas métricas:
      | horas estimadas | 32    |
      | horas reales    | 30    |
      | desviación (h)  | -2    |
      | desviación (%)  | -6.25 |

  Escenario: Defectos detectados y resueltos en un Sprint
    Dado un proyecto de métricas con 4 horas por Story Point
    Y el Sprint de métricas "S1" tiene 2 defectos detectados
    Y el Sprint de métricas "S1" tiene 1 defecto resuelto
    Y el Sprint de métricas "S1" está finalizado
    Cuando consulto las métricas del proyecto
    Entonces el Sprint "S1" informa estas métricas:
      | defectos detectados | 2 |
      | defectos resueltos  | 1 |

  Escenario: Un defecto detectado en un Sprint y resuelto en otro
    Dado un proyecto de métricas con 4 horas por Story Point
    Y un defecto de métricas detectado en el Sprint "S1" y resuelto en el Sprint "S2"
    Y el Sprint de métricas "S1" está finalizado
    Y el Sprint de métricas "S2" está finalizado
    Cuando consulto las métricas del proyecto
    Entonces el Sprint "S1" informa estas métricas:
      | defectos detectados | 1 |
      | defectos resueltos  | 0 |
    Y el Sprint "S2" informa estas métricas:
      | defectos detectados | 0 |
      | defectos resueltos  | 1 |

  Escenario: Valores del proyecto con dos Sprints finalizados
    Dado un proyecto de métricas con 4 horas por Story Point
    Y el Sprint de métricas "S1" tiene una historia completada de 5 Story Points
    Y el Sprint de métricas "S2" tiene una historia completada de 8 Story Points
    Y el Sprint de métricas "S1" está finalizado
    Y el Sprint de métricas "S2" está finalizado
    Cuando consulto las métricas del proyecto
    Entonces el proyecto informa estas métricas:
      | Story Points completados | 13  |
      | velocidad                | 6.5 |

  Escenario: Un Sprint activo no aporta a las métricas
    Dado un proyecto de métricas con 4 horas por Story Point
    Y el Sprint de métricas "S1" tiene una historia completada de 5 Story Points
    Y el Sprint de métricas "S1" está finalizado
    Y el Sprint de métricas "S2" tiene 10 horas imputadas
    Y el Sprint de métricas "S2" tiene 1 defecto detectado
    Y el Sprint de métricas "S2" está activo
    Cuando consulto las métricas del proyecto
    Entonces las métricas incluyen solo el Sprint "S1"
    Y el proyecto informa estas métricas:
      | Story Points completados | 5 |
      | horas reales             | 0 |
      | defectos detectados      | 0 |

  Escenario: Editar una historia arrastrada no cambia el Sprint que la arrastró
    Dado un proyecto de métricas con 4 horas por Story Point
    Y la historia de métricas "H" de 3 Story Points fue arrastrada en el Sprint "S1"
    Y la historia de métricas "H" se edita a 5 Story Points en el Product Backlog
    Y la historia de métricas "H" se completa en el Sprint "S2"
    Y el Sprint de métricas "S2" está finalizado
    Cuando consulto las métricas del proyecto
    Entonces el Sprint "S1" informa estas métricas:
      | Story Points planificados | 3 |
      | Story Points completados  | 0 |
      | historias completadas (%) | 0 |
    Y el Sprint "S2" informa estas métricas:
      | Story Points completados  | 5   |
      | historias completadas (%) | 100 |

  Escenario: Consultar las métricas de un proyecto inexistente
    Cuando consulto las métricas del proyecto con identificador 9999, que no existe
    Entonces la consulta de métricas se rechaza con error "proyecto no encontrado"

  Escenario: Una historia arrastrada cuenta en cada Sprint en los totales del proyecto
    Dado un proyecto de métricas con 4 horas por Story Point
    Y la historia de métricas "H" de 3 Story Points fue arrastrada en el Sprint "S1"
    Y la historia de métricas "H" se completa en el Sprint "S2"
    Y el Sprint de métricas "S2" está finalizado
    Cuando consulto las métricas del proyecto
    Entonces el proyecto informa estas métricas:
      | Story Points planificados | 6  |
      | Story Points completados  | 3  |
      | horas estimadas           | 24 |
      | historias completadas (%) | 50 |

  Escenario: Proyecto sin Sprints finalizados y con horas por Story Point
    Dado un proyecto de métricas con 4 horas por Story Point
    Y el Sprint de métricas "S1" tiene 10 horas imputadas
    Y el Sprint de métricas "S1" está activo
    Cuando consulto las métricas del proyecto
    Entonces las métricas no incluyen Sprints
    Y el proyecto informa estas métricas:
      | Story Points planificados | 0         |
      | Story Points completados  | 0         |
      | velocidad                 | 0         |
      | historias completadas (%) | 0         |
      | horas estimadas           | 0         |
      | horas reales              | 0         |
      | desviación (h)            | 0         |
      | desviación (%)            | sin datos |
      | defectos detectados       | 0         |
      | defectos resueltos        | 0         |

  Escenario: Proyecto sin Sprints finalizados y sin horas por Story Point
    Dado un proyecto de métricas sin horas por Story Point
    Cuando consulto las métricas del proyecto
    Entonces las métricas no incluyen Sprints
    Y el proyecto informa estas métricas:
      | Story Points planificados | 0         |
      | horas reales              | 0         |
      | horas estimadas           | sin datos |
      | desviación (h)            | sin datos |
      | desviación (%)            | sin datos |

  Escenario: Definir las horas por Story Point
    Dado un proyecto de métricas sin horas por Story Point
    Y el Sprint de métricas "S1" tiene una historia completada de 5 Story Points
    Y el Sprint de métricas "S1" está finalizado
    Cuando edito el proyecto de métricas y defino 4 horas por Story Point
    Y consulto las métricas del proyecto
    Entonces el proyecto de métricas tiene 4 horas por Story Point
    Y el Sprint "S1" informa estas métricas:
      | horas estimadas | 20 |

  Escenario: Métricas de un Sprint finalizado sin horas por Story Point
    Dado un proyecto de métricas sin horas por Story Point
    Y el Sprint de métricas "S1" tiene una historia completada de 5 Story Points
    Y el Sprint de métricas "S1" tiene 30 horas imputadas
    Y el Sprint de métricas "S1" está finalizado
    Cuando consulto las métricas del proyecto
    Entonces el Sprint "S1" informa estas métricas:
      | horas estimadas          | sin datos |
      | desviación (h)           | sin datos |
      | desviación (%)           | sin datos |
      | horas reales             | 30        |
      | Story Points completados | 5         |

  Escenario: Editar el proyecto sin mencionar el parámetro lo conserva
    Dado un proyecto de métricas con 4 horas por Story Point
    Cuando edito el proyecto de métricas sin mencionar las horas por Story Point
    Entonces el proyecto de métricas tiene 4 horas por Story Point

  Escenario: Editar el proyecto dejando el parámetro vacío lo borra
    Dado un proyecto de métricas con 4 horas por Story Point
    Y el Sprint de métricas "S1" tiene una historia completada de 5 Story Points
    Y el Sprint de métricas "S1" está finalizado
    Cuando edito el proyecto de métricas y dejo vacías las horas por Story Point
    Y consulto las métricas del proyecto
    Entonces el proyecto de métricas no tiene horas por Story Point
    Y el Sprint "S1" informa estas métricas:
      | horas estimadas | sin datos |
      | desviación (h)  | sin datos |

  Esquema del escenario: Las horas por Story Point inválidas se rechazan
    Dado un proyecto de métricas con 4 horas por Story Point
    Cuando intento editar el proyecto de métricas con las horas por Story Point <valor>
    Entonces la edición del proyecto de métricas se rechaza con un error sobre el campo "horas_por_story_point"
    Y el proyecto de métricas tiene 4 horas por Story Point

    Ejemplos:
      | valor   |
      | 0       |
      | -1      |
      | "abc"   |
      | ""      |

  Escenario: Cambiar el parámetro recalcula los Sprints finalizados
    Dado un proyecto de métricas con 4 horas por Story Point
    Y el Sprint de métricas "S1" tiene una historia completada de 5 Story Points
    Y el Sprint de métricas "S1" tiene una historia arrastrada de 3 Story Points
    Y el Sprint de métricas "S1" tiene 30 horas imputadas
    Y el Sprint de métricas "S1" está finalizado
    Cuando edito el proyecto de métricas y defino 5 horas por Story Point
    Y consulto las métricas del proyecto
    Entonces el Sprint "S1" informa estas métricas:
      | horas estimadas | 40  |
      | desviación (h)  | -10 |
      | desviación (%)  | -25 |

  Escenario: No se cargan Story Points por edición con el Sprint finalizado
    Dado un proyecto de métricas con 4 horas por Story Point
    Y la historia de métricas "H" sin estimar se completó en el Sprint finalizado "S1"
    Cuando intento cargar 5 Story Points a la historia de métricas "H" por edición
    Entonces la carga se rechaza con una advertencia sobre el campo "estimacion_sp"
    Y la historia de métricas "H" sigue sin estimar

  Escenario: Editar otros datos sin cargar Story Points con el Sprint finalizado
    Dado un proyecto de métricas con 4 horas por Story Point
    Y la historia de métricas "H" sin estimar se completó en el Sprint finalizado "S1"
    Cuando cambio el título de la historia de métricas "H" a "Título nuevo" sin cargar Story Points
    Entonces la historia de métricas "H" se llama "Título nuevo"
    Y la historia de métricas "H" sigue sin estimar

  Escenario: No se registra la estimación acordada con el Sprint finalizado
    Dado un proyecto de métricas con 4 horas por Story Point
    Y la historia de métricas "H" sin estimar se completó en el Sprint finalizado "S1"
    Y la historia de métricas "H" tiene su ronda de votación revelada
    Cuando el Scrum Master intenta registrar 5 como estimación acordada de la historia de métricas "H"
    Entonces la carga se rechaza con una advertencia sobre el campo "estimacion_sp"
    Y la historia de métricas "H" sigue sin estimar

  Esquema del escenario: La primera carga sigue permitida con el Sprint activo
    Dado un proyecto de métricas con 4 horas por Story Point
    Y la historia de métricas "H" sin estimar está en el Sprint activo "S1"
    Y la historia de métricas "H" tiene su ronda de votación revelada
    Cuando cargo 5 Story Points a la historia de métricas "H" por <vía>
    Entonces la historia de métricas "H" tiene 5 Story Points

    Ejemplos:
      | vía      |
      | edición  |
      | votación |

  Escenario: Estimar una historia arrastrada no cambia el Sprint que la arrastró
    Dado un proyecto de métricas con 4 horas por Story Point
    Y la historia de métricas "H" sin estimar fue arrastrada en el Sprint "S1"
    Cuando cargo 5 Story Points a la historia de métricas "H" por edición
    Y consulto las métricas del proyecto
    Entonces la historia de métricas "H" tiene 5 Story Points
    Y el Sprint "S1" informa estas métricas:
      | Story Points planificados | 0 |
