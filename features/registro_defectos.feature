# language: es

# Cobertura por categoría:
# - Normal: registro de un defecto en el Sprint activo (US1.1); resolución de un defecto
#   pendiente (US2.1); listado de los defectos de un proyecto (US3.1).
# - Alternativo: historia relacionada fuera del Sprint activo (US1.2); defecto resuelto en
#   un Sprint distinto del de detección (US2.2); proyecto sin defectos (US3.2).
# - Límite: ningún criterio de aceptación es de borde. El largo de 2000 caracteres de la
#   descripción y el mismo Sprint de detección y de resolución se prueban con tests
#   unitarios del dominio.
# - Error: sin Sprint activo, datos inválidos, historia inexistente o ajena y proyecto
#   inexistente al registrar (US1.3–US1.6); sin Sprint activo, defecto ya resuelto y
#   defecto inexistente al resolver (US2.3–US2.5); proyecto inexistente al listar (US3.3).

@HU-08
Característica: Registro y seguimiento de defectos

  Como integrante del equipo
  Quiero registrar un defecto técnico indicando su severidad, la historia relacionada y el sprint de detección
  Para poder hacer un seguimiento de la calidad y asegurar que se resuelva antes del despliegue final

  Antecedentes:
    Dado que existe un proyecto con identificador 1
    Y que el proyecto tiene la historia "Carrito de compras" en el Product Backlog

  Escenario: Registro de un defecto en el Sprint activo
    Dado que el proyecto tiene el Sprint "Sprint 2" en estado "Activo"
    Cuando registro el defecto "El total no se actualiza" con severidad "Alta" en la historia "Carrito de compras"
    Entonces el defecto queda registrado con la descripción "El total no se actualiza", severidad "Alta" y la historia "Carrito de compras"
    Y el defecto queda en estado "Pendiente", detectado en el Sprint "Sprint 2" y sin Sprint de resolución

  Esquema del escenario: Defecto de una historia que no está en el Sprint activo
    Dado que la historia "Carrito de compras" <ubicacion>
    Y que el proyecto tiene el Sprint "Sprint 2" en estado "Activo"
    Cuando registro el defecto "El total no se actualiza" con severidad "Media" en la historia "Carrito de compras"
    Entonces el defecto queda en estado "Pendiente", detectado en el Sprint "Sprint 2" y sin Sprint de resolución

    Ejemplos:
      | ubicacion                                         |
      | se completó en el Sprint "Sprint 1", ya finalizado |
      | sigue en el Product Backlog                       |

  Esquema del escenario: Registro sin un Sprint activo
    Dado <sprints>
    Cuando intento registrar el defecto "El total no se actualiza" con severidad "Media" en la historia "Carrito de compras"
    Entonces el registro se rechaza con el mensaje "el proyecto no tiene un Sprint activo al que asociar el defecto como Sprint de detección"
    Y el proyecto sigue sin defectos

    Ejemplos:
      | sprints                                                    |
      | que el proyecto no tiene Sprints                           |
      | que el proyecto tiene el Sprint "Sprint 3" en estado "Pendiente"  |
      | que el proyecto tiene el Sprint "Sprint 1" en estado "Finalizado" |

  Esquema del escenario: Registro con datos inválidos
    Dado que el proyecto tiene el Sprint "Sprint 2" en estado "Activo"
    Cuando intento registrar un defecto <variante>
    Entonces el registro se rechaza con un error sobre el campo "<campo>"
    Y el proyecto sigue sin defectos

    Ejemplos:
      | variante                   | campo       |
      | con la descripción vacía   | descripcion |
      | con la severidad "Crítica" | severidad   |
      | sin historia relacionada   | historia_id |

  Esquema del escenario: Registro con una historia relacionada que no es del proyecto
    Dado que el proyecto tiene el Sprint "Sprint 2" en estado "Activo"
    Cuando intento registrar un defecto relacionado con <historia>
    Entonces el registro se rechaza con un error sobre el campo "historia_id"
    Y el proyecto sigue sin defectos

    Ejemplos:
      | historia                            |
      | la historia con identificador 999   |
      | la historia "Login" de otro proyecto |

  Escenario: Registro en un proyecto inexistente
    Cuando intento registrar un defecto en el proyecto con identificador 999, que no existe
    Entonces el registro se rechaza con error "proyecto no encontrado"

  Escenario: Resolución de un defecto pendiente
    Dado que el proyecto tiene el Sprint "Sprint 2" en estado "Activo"
    Y que se registró el defecto "El total no se actualiza" con severidad "Alta"
    Cuando marco el defecto como "Resuelto"
    Entonces el defecto queda en estado "Resuelto", resuelto en el Sprint "Sprint 2"
    Y el defecto sigue detectado en el Sprint "Sprint 2"

  Escenario: Resolución en un Sprint distinto del de detección
    Dado que el proyecto tiene el Sprint "Sprint 1" en estado "Activo"
    Y que se registró el defecto "El total no se actualiza" con severidad "Alta"
    Y que se cerró el Sprint "Sprint 1"
    Y que el proyecto tiene el Sprint "Sprint 2" en estado "Activo"
    Cuando marco el defecto como "Resuelto"
    Entonces el defecto queda en estado "Resuelto", resuelto en el Sprint "Sprint 2"
    Y el defecto sigue detectado en el Sprint "Sprint 1"

  Escenario: Resolución sin un Sprint activo
    Dado que el proyecto tiene el Sprint "Sprint 1" en estado "Activo"
    Y que se registró el defecto "El total no se actualiza" con severidad "Alta"
    Y que se cerró el Sprint "Sprint 1"
    Cuando intento marcar el defecto como "Resuelto"
    Entonces la resolución se rechaza con el mensaje "el proyecto no tiene un Sprint activo al que asociar el defecto como Sprint de resolución"
    Y el defecto queda en estado "Pendiente", detectado en el Sprint "Sprint 1" y sin Sprint de resolución

  Escenario: Resolución de un defecto ya resuelto
    Dado que el proyecto tiene el Sprint "Sprint 2" en estado "Activo"
    Y que se registró el defecto "El total no se actualiza" con severidad "Alta"
    Y que se marcó como "Resuelto" el defecto "El total no se actualiza"
    Y que se cerró el Sprint "Sprint 2"
    Y que el proyecto tiene el Sprint "Sprint 3" en estado "Activo"
    Cuando intento marcar el defecto como "Resuelto"
    Entonces la resolución se rechaza con el mensaje "el defecto ya está resuelto"
    Y el defecto queda en estado "Resuelto", resuelto en el Sprint "Sprint 2"

  Escenario: Resolución de un defecto inexistente
    Cuando intento marcar como "Resuelto" el defecto con identificador 999
    Entonces la resolución se rechaza con error "defecto no encontrado"

  Escenario: Consulta de los defectos de un proyecto
    Dado que el proyecto tiene el Sprint "Sprint 1" en estado "Activo"
    Y que se registró el defecto "El total no se actualiza" con severidad "Alta"
    Y que se registró el defecto "El botón no responde" con severidad "Baja"
    Y que se cerró el Sprint "Sprint 1"
    Y que el proyecto tiene el Sprint "Sprint 2" en estado "Activo"
    Y que se marcó como "Resuelto" el defecto "El total no se actualiza"
    Y que se registró el defecto "Falta un mensaje de error" con severidad "Media"
    Y que otro proyecto tiene un defecto registrado
    Cuando consulto los defectos del proyecto
    Entonces veo estos defectos, en este orden, cada uno con su fecha de registro:
      | descripcion               | severidad | estado    | historia           | detectado | resuelto |
      | El total no se actualiza  | Alta      | Resuelto  | Carrito de compras | Sprint 1  | Sprint 2 |
      | El botón no responde      | Baja      | Pendiente | Carrito de compras | Sprint 1  |          |
      | Falta un mensaje de error | Media     | Pendiente | Carrito de compras | Sprint 2  |          |

  Escenario: Consulta de un proyecto sin defectos
    Cuando consulto los defectos del proyecto
    Entonces veo una lista de defectos vacía, sin error

  Escenario: Consulta de los defectos de un proyecto inexistente
    Cuando consulto los defectos del proyecto con identificador 999, que no existe
    Entonces la consulta de defectos se rechaza con error "proyecto no encontrado"
