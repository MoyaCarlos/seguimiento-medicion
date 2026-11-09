# language: es

# Cobertura por categoría:
# - Normal: reporte de un Sprint finalizado con todas sus secciones (US1.1); reporte
#   completo del proyecto (US2.1); generación exitosa de un PDF válido (US4.2).
# - Alternativo: Sprint finalizado sin horas ni defectos (US1.2); proyecto sin Sprint
#   Activo (US2.2); proyecto sin historias en el Product Backlog (US2.3); proyecto sin
#   Sprints finalizados con y sin Sprint Activo (US3.1 y US3.2); historia arrastrada y
#   completada en otro Sprint (US1.5); defecto detectado y resuelto en Sprints distintos.
# - Límite: proyecto sin Sprints finalizados y sin historias (US3.3); Sprint finalizado
#   sin historias comprometidas (US1.2).
# - Error: Sprint Pendiente o Activo (US1.3); Sprint inexistente (US1.4); proyecto
#   inexistente (US2.4 y US3.4); fallo real de generación sin archivo (US4.1 y US4.3).

@HU-10
Característica: Generación de reporte en PDF

  Como Scrum Master
  Quiero exportar un reporte en PDF de un proyecto o Sprint con sus historias, estimaciones, esfuerzo, métricas y defectos
  Para poder compartir el estado del proyecto fuera del sistema

  Antecedentes:
    Dado un proyecto de reporte con identificador 1
    Y el proyecto de reporte tiene 4 horas por Story Point
    Y el proyecto de reporte tiene al integrante "Ana" con rol "product_builder"

  Escenario: Reporte de un Sprint finalizado con todas sus secciones
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Y la historia "Carrito de compras" está en el Sprint "Sprint 1" con 5 Story Points y estado "Completada"
    Y la historia "Pago" se arrastró al cerrar el Sprint "Sprint 1" con 3 Story Points
    Y el Sprint "Sprint 1" tiene 6 horas registradas
    Y el Sprint "Sprint 1" tiene el defecto "El total no se actualiza" con severidad "Alta" detectado
    Y el Sprint "Sprint 1" tiene el defecto "El total no se actualiza" con severidad "Alta" resuelto
    Cuando solicito el reporte del Sprint "Sprint 1"
    Entonces obtengo un PDF válido con el reporte del Sprint "Sprint 1"
    Y el reporte del Sprint incluye la historia completada "Carrito de compras"
    Y el reporte del Sprint incluye la historia planificada "Pago" como no completada
    Y el reporte del Sprint informa 6 horas reales
    Y el reporte del Sprint informa los defectos detectados y resueltos

  Escenario: Reporte de un Sprint finalizado sin horas, defectos ni historias
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Sprint vacío"
    Cuando solicito el reporte del Sprint "Sprint 1"
    Entonces obtengo un PDF válido con el reporte del Sprint "Sprint 1"
    Y el reporte del Sprint informa 0 horas reales
    Y el reporte del Sprint no informa defectos

  Esquema del escenario: Reporte de un Sprint no finalizado
    Dado el Sprint "Sprint 1" en estado "<estado>"
    Cuando solicito el reporte del Sprint "Sprint 1"
    Entonces la solicitud se rechaza con el mensaje "solo se puede generar el reporte de un Sprint finalizado"
    Y no se descarga ningún archivo

    Ejemplos:
      | estado   |
      | Pendiente |
      | Activo    |

  Escenario: Reporte de un Sprint inexistente
    Cuando solicito el reporte del Sprint con identificador 999, que no existe
    Entonces la solicitud se rechaza con el mensaje "sprint no encontrado"
    Y no se descarga ningún archivo

  Escenario: Historia arrastrada planificada en el Sprint anterior
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Y la historia "Pago" se arrastró al cerrar el Sprint "Sprint 1" con 3 Story Points
    Cuando solicito el reporte del Sprint "Sprint 1"
    Entonces el reporte del Sprint incluye la historia planificada "Pago" como no completada

  Escenario: Historia arrastrada completada en otro Sprint
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Y la historia "Pago" se arrastró al cerrar el Sprint "Sprint 1" con 3 Story Points
    Y el Sprint "Sprint 2" finalizado con el objetivo "Terminar el pago"
    Y la historia "Pago" está en el Sprint "Sprint 2" con 3 Story Points y estado "Completada"
    Cuando solicito el reporte del Sprint "Sprint 2"
    Entonces el reporte del Sprint incluye la historia completada "Pago"

  Escenario: Defecto detectado en un Sprint finalizado
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Y el Sprint "Sprint 2" finalizado con el objetivo "Terminar el pago"
    Y el defecto "El total no se actualiza" fue detectado en el Sprint "Sprint 1" y resuelto en el Sprint "Sprint 2"
    Cuando solicito el reporte del Sprint "Sprint 1"
    Entonces el reporte del Sprint lista el defecto "El total no se actualiza" como detectado

  Escenario: Defecto resuelto en otro Sprint finalizado
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Y el Sprint "Sprint 2" finalizado con el objetivo "Terminar el pago"
    Y el defecto "El total no se actualiza" fue detectado en el Sprint "Sprint 1" y resuelto en el Sprint "Sprint 2"
    Cuando solicito el reporte del Sprint "Sprint 2"
    Entonces el reporte del Sprint lista el defecto "El total no se actualiza" como resuelto

  Escenario: Reporte completo del proyecto
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Y la historia "Carrito de compras" está en el Sprint "Sprint 1" con 5 Story Points y estado "Completada"
    Y el Sprint "Sprint 1" tiene 6 horas registradas
    Y el Sprint "Sprint 2" finalizado con el objetivo "Terminar el pago"
    Y la historia "Pago" está en el Sprint "Sprint 2" con 3 Story Points y estado "Completada"
    Y el Sprint "Sprint 3" en estado "Activo" con el objetivo "Pulir la entrega"
    Y la historia "Ajustes finales" está en el Sprint "Sprint 3" con 2 Story Points y estado "En progreso"
    Y el Sprint "Sprint 3" tiene 9 horas registradas
    Y la historia "Mejoras futuras" está en el Product Backlog con prioridad "S" y 8 Story Points
    Cuando solicito el reporte del proyecto
    Entonces obtengo un PDF válido con el reporte del proyecto
    Y el reporte incluye los datos del proyecto, sus integrantes y el Product Backlog
    Y el reporte incluye 2 Sprints finalizados en orden
    Y el reporte incluye el Sprint Activo "Sprint 3" con sus historias
    Y el reporte no suma el Sprint Activo a las métricas

  Escenario: Proyecto sin Sprint Activo
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Y el Sprint "Sprint 2" finalizado con el objetivo "Terminar el pago"
    Y la historia "Mejoras futuras" está en el Product Backlog con prioridad "S" y 8 Story Points
    Cuando solicito el reporte del proyecto
    Entonces el reporte no incluye la sección del Sprint Activo

  Escenario: Proyecto sin historias en el Product Backlog
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Y la historia "Carrito de compras" está en el Sprint "Sprint 1" con 5 Story Points y estado "Completada"
    Cuando solicito el reporte del proyecto
    Entonces el reporte no incluye historias en el Product Backlog

  Escenario: Proyecto inexistente
    Cuando solicito el reporte del proyecto con identificador 999, que no existe
    Entonces la solicitud se rechaza con el mensaje "proyecto no encontrado"
    Y no se descarga ningún archivo

  Escenario: Sprint finalizado sin horas ni defectos en el reporte del proyecto
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Sprint vacío"
    Cuando solicito el reporte del proyecto
    Entonces obtengo un PDF válido con el reporte del proyecto
    Y el reporte del Sprint "Sprint 1" informa 0 horas reales y sin defectos

  Escenario: Defecto detectado y resuelto en Sprints distintos dentro del proyecto
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Y el Sprint "Sprint 2" finalizado con el objetivo "Terminar el pago"
    Y el defecto "El total no se actualiza" fue detectado en el Sprint "Sprint 1" y resuelto en el Sprint "Sprint 2"
    Cuando solicito el reporte del proyecto
    Entonces el reporte del proyecto lista el defecto "El total no se actualiza" como detectado en el Sprint "Sprint 1" y resuelto en el Sprint "Sprint 2"

  Escenario: Product Backlog y Sprint Activo en secciones distintas
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Y el Sprint "Sprint 2" en estado "Activo" con el objetivo "Pulir la entrega"
    Y la historia "Ajustes finales" está en el Sprint "Sprint 2" con 2 Story Points y estado "En progreso"
    Y la historia "Mejoras futuras" está en el Product Backlog con prioridad "S" y 8 Story Points
    Cuando solicito el reporte del proyecto
    Entonces el reporte incluye la historia "Mejoras futuras" en el Product Backlog
    Y el reporte incluye la historia "Ajustes finales" solo en el Sprint Activo

  Escenario: Proyecto sin Sprints finalizados con Sprint Activo
    Dado el Sprint "Sprint 1" en estado "Activo" con el objetivo "Primer Sprint"
    Y la historia "Ajustes finales" está en el Sprint "Sprint 1" con 2 Story Points y estado "En progreso"
    Y la historia "Mejoras futuras" está en el Product Backlog con prioridad "S" y 8 Story Points
    Cuando solicito el reporte del proyecto
    Entonces obtengo un PDF válido con el reporte del proyecto
    Y el reporte incluye el Sprint Activo "Sprint 1" con sus historias
    Y el reporte informa las métricas, el esfuerzo y los defectos como "sin datos"

  Escenario: Proyecto sin Sprints finalizados ni Sprint Activo
    Dado la historia "Mejoras futuras" está en el Product Backlog con prioridad "S" y 8 Story Points
    Cuando solicito el reporte del proyecto
    Entonces el reporte no incluye la sección del Sprint Activo
    Y el reporte informa las métricas, el esfuerzo y los defectos como "sin datos"

  Escenario: Proyecto sin Sprints finalizados ni historias en el Product Backlog
    Cuando solicito el reporte del proyecto
    Entonces obtengo un PDF válido con el reporte del proyecto
    Y el reporte no incluye historias en el Product Backlog
    Y el reporte informa las métricas, el esfuerzo y los defectos como "sin datos"

  Escenario: Proyecto inexistente sin Sprints finalizados
    Cuando solicito el reporte del proyecto con identificador 999, que no existe
    Entonces la solicitud se rechaza con el mensaje "proyecto no encontrado"
    Y no se descarga ningún archivo

  Escenario: Error real en la generación del PDF
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Cuando solicito el reporte del proyecto con un generador que falla
    Entonces la solicitud se rechaza con el mensaje "no se pudo generar el reporte"
    Y no se descarga ningún archivo

  Escenario: La generación exitosa devuelve un PDF válido
    Dado el Sprint "Sprint 1" finalizado con el objetivo "Cerrar el carrito"
    Cuando solicito el reporte del proyecto
    Entonces obtengo un PDF válido con el reporte del proyecto
