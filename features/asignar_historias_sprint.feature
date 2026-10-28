# language: es

# Cobertura por categoría:
# - Normal: asignación a un Sprint "Pendiente".
# - Alternativo: asignación a un Sprint "Activo".
# - Límite: la asignación conserva el estado "Nueva" y los Story Points; historia que ya está en el mismo Sprint de destino.
# - Error: historia "Completada"; Sprint "Finalizado"; Sprint de otro proyecto; historia inexistente; Sprint inexistente.

@HU-06
Característica: Movimiento de historias al Sprint Backlog

  Como Product Builder
  Quiero seleccionar historias de usuario del Product Backlog y asignarlas a un Sprint específico
  Para poder definir el compromiso de trabajo de la iteración actual

  Antecedentes:
    Dado que existe un proyecto con identificador 1

  Escenario: Asignación de una historia a un Sprint Pendiente
    Dado que el proyecto tiene un Sprint a planificar en estado "Pendiente"
    Y que existe la historia "Iniciar sesión" en el Product Backlog del proyecto
    Cuando asigno la historia "Iniciar sesión" al Sprint a planificar
    Entonces la historia "Iniciar sesión" ya no aparece en el Product Backlog del proyecto
    Y la historia "Iniciar sesión" forma parte del Sprint Backlog del Sprint a planificar

  Escenario: Asignación de una historia a un Sprint Activo
    Dado que el proyecto tiene un Sprint a planificar en estado "Activo"
    Y que existe la historia "Iniciar sesión" en el Product Backlog del proyecto
    Cuando asigno la historia "Iniciar sesión" al Sprint a planificar
    Entonces la historia "Iniciar sesión" ya no aparece en el Product Backlog del proyecto
    Y la historia "Iniciar sesión" forma parte del Sprint Backlog del Sprint a planificar

  Escenario: La asignación conserva el estado y los Story Points de la historia
    Dado que el proyecto tiene un Sprint a planificar en estado "Pendiente"
    Y que existe la historia "Iniciar sesión" con 5 Story Points en el Product Backlog del proyecto
    Cuando asigno la historia "Iniciar sesión" al Sprint a planificar
    Entonces la historia "Iniciar sesión" conserva el estado "Nueva" y sus 5 Story Points

  Escenario: No se puede asignar una historia completada
    Dado que el proyecto tiene un Sprint a planificar en estado "Pendiente"
    Y que existe la historia "Exportar reporte" en estado "Completada" en otro Sprint del proyecto
    Cuando asigno la historia "Exportar reporte" al Sprint a planificar
    Entonces la asignación se rechaza con el mensaje "solo se pueden planificar historias pendientes"
    Y la historia "Exportar reporte" queda sin cambios

  Escenario: No se puede asignar una historia que ya está en el mismo Sprint
    Dado que el proyecto tiene un Sprint a planificar en estado "Pendiente"
    Y que existe la historia "Iniciar sesión" en el Product Backlog del proyecto
    Y que la historia "Iniciar sesión" ya está asignada al Sprint a planificar
    Cuando asigno la historia "Iniciar sesión" al Sprint a planificar
    Entonces la asignación se rechaza con el mensaje "la historia ya pertenece a un Sprint"
    Y la historia "Iniciar sesión" queda sin cambios

  Escenario: Un Sprint Finalizado no admite historias
    Dado que el proyecto tiene un Sprint a planificar en estado "Finalizado"
    Y que existe la historia "Iniciar sesión" en el Product Backlog del proyecto
    Cuando asigno la historia "Iniciar sesión" al Sprint a planificar
    Entonces la asignación se rechaza con el mensaje "un Sprint finalizado no admite historias"
    Y la historia "Iniciar sesión" queda sin cambios

  Escenario: No se puede asignar una historia a un Sprint de otro proyecto
    Dado que existe otro proyecto con un Sprint en estado "Pendiente"
    Y que existe la historia "Iniciar sesión" en el Product Backlog del proyecto
    Cuando asigno la historia "Iniciar sesión" al Sprint del otro proyecto
    Entonces la asignación se rechaza con el mensaje "el Sprint no pertenece al proyecto de la historia"
    Y la historia "Iniciar sesión" queda sin cambios

  Escenario: Asignación de una historia inexistente
    Cuando asigno la historia con identificador 999, que no existe, al Sprint con identificador 999, que tampoco existe
    Entonces la asignación se rechaza con error "historia no encontrada"

  Escenario: Asignación a un Sprint inexistente
    Dado que existe la historia "Iniciar sesión" en el Product Backlog del proyecto
    Cuando asigno la historia "Iniciar sesión" al Sprint con identificador 999, que no existe
    Entonces la asignación se rechaza con error "sprint no encontrado"
    Y la historia "Iniciar sesión" queda sin cambios
