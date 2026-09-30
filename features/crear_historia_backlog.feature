# language: es

@HU-01
Característica: Creación de historias de usuario en el Product Backlog

  Como integrante del equipo (Product Builder o Scrum Master)
  Quiero crear historias de usuario con prioridad y estado
  Para alimentar y organizar el Product Backlog del proyecto

  Antecedentes:
    Dado que existe un proyecto con identificador 1

  Escenario: Creación exitosa de una historia de usuario
    Dado que estoy en el panel del Product Backlog del proyecto
    Cuando ingreso un título, una descripción y una prioridad "M" válidos y presiono "Guardar"
    Entonces la historia se guarda de forma persistente en SQLite
    Y queda al final del Product Backlog según el orden de creación (identificador incremental posterior)
    Y su estado es "Nueva"
    Y no tiene estimación en Story Points asignada

  Escenario: Falta el título
    Dado que intento crear una historia
    Cuando dejo el campo "Título" en blanco y presiono "Guardar"
    Entonces el sistema muestra una advertencia de validación
    Y no registra la historia en el Product Backlog

  Escenario: Prioridad inválida
    Dado que intento crear una historia
    Cuando ingreso la prioridad "Urgente", que no pertenece al enum MoSCoW (Must have / Should have / Could have / Won't have), y presiono "Guardar"
    Entonces el sistema muestra un error de validación
    Y no registra la historia en el Product Backlog
