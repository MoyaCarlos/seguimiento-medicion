# language: es

# Issue #174 (parte de HU-01): criterios de aceptación al crear una historia.
# Cobertura por categoría:
# - Normal: alta con criterios válidos.
# - Alternativo: alta sin criterios (omitidos), solo espacios y criterios con espacios externos.
# - Límite: criterios de 2000 caracteres.
# - Error: criterios de 2001 caracteres, rechazados con campo "criterios_aceptacion".
# La pantalla de creación y la confirmación se cubren con pruebas de componente en web/.

@HU-01
Característica: Criterios de aceptación al crear una historia

  Como integrante del equipo
  Quiero cargar los criterios de aceptación al crear una historia
  Para dejar la historia completa sin pasar por la edición

  Antecedentes:
    Dado que existe un proyecto con identificador 1
    Y que estoy en el panel del Product Backlog del proyecto

  Escenario: Alta con criterios de aceptación
    Cuando ingreso un título, una descripción, una prioridad "M" y los criterios de aceptación "Dado un usuario, cuando ingresa, entonces accede." y presiono "Guardar"
    Entonces la historia se guarda de forma persistente en SQLite
    Y los criterios de aceptación de la historia son "Dado un usuario, cuando ingresa, entonces accede."

  Escenario: Alta sin criterios de aceptación
    Cuando ingreso un título, una descripción y una prioridad "M" sin criterios de aceptación y presiono "Guardar"
    Entonces la historia se guarda de forma persistente en SQLite
    Y la historia no tiene criterios de aceptación

  Escenario: Criterios con solo espacios equivalen a no tener criterios
    Cuando ingreso un título, una descripción, una prioridad "M" y los criterios de aceptación "   " y presiono "Guardar"
    Entonces la historia se guarda de forma persistente en SQLite
    Y la historia no tiene criterios de aceptación

  Escenario: Los criterios se guardan sin espacios externos
    Cuando ingreso un título, una descripción, una prioridad "M" y los criterios de aceptación "  Dado, cuando, entonces  " y presiono "Guardar"
    Entonces la historia se guarda de forma persistente en SQLite
    Y los criterios de aceptación de la historia son "Dado, cuando, entonces"

  Escenario: Criterios en el límite de 2000 caracteres
    Cuando ingreso un título, una descripción, una prioridad "M" y criterios de aceptación de 2000 caracteres y presiono "Guardar"
    Entonces la historia se guarda de forma persistente en SQLite
    Y los criterios de aceptación tienen 2000 caracteres

  Escenario: Criterios de más de 2000 caracteres
    Cuando ingreso un título, una descripción, una prioridad "M" y criterios de aceptación de 2001 caracteres y presiono "Guardar"
    Entonces el sistema muestra un error de validación sobre el campo "criterios_aceptacion"
    Y no registra la historia en el Product Backlog
