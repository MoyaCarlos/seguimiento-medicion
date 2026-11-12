# language: es

@HU-05 @168
Característica: Sprints dentro del plazo del proyecto

  Como Scrum Master
  Quiero que al iniciar un Sprint no pueda ponerle fechas fuera del plazo del proyecto
  Para que el seguimiento y las métricas no cuenten trabajo fuera de la vida del proyecto

  Antecedentes:
    Dado que existe un proyecto con identificador 1
    Y el proyecto tiene un Sprint pendiente de iniciar

  Escenario: Rechazar un Sprint que empieza antes que el proyecto
    Dado que el plazo del proyecto empieza el "2026-10-31"
    Cuando el Scrum Master inicia el Sprint del "2026-10-30" al "2026-11-06"
    Entonces el inicio se rechaza en el campo "fecha_inicio"
    Y el Sprint sigue pendiente, sin Sprint Goal ni fechas

  Escenario: Aceptar un Sprint que empieza el mismo día que el proyecto
    Dado que el plazo del proyecto empieza el "2026-10-31"
    Cuando el Scrum Master inicia el Sprint del "2026-10-31" al "2026-11-06"
    Entonces el Sprint queda activo

  Escenario: Rechazar un Sprint que termina después que el proyecto
    Dado que el plazo del proyecto termina el "2026-11-16"
    Cuando el Scrum Master inicia el Sprint del "2026-11-10" al "2026-11-17"
    Entonces el inicio se rechaza en el campo "fecha_fin"
    Y el Sprint sigue pendiente, sin Sprint Goal ni fechas

  Escenario: Aceptar un Sprint que termina el mismo día que el proyecto
    Dado que el plazo del proyecto termina el "2026-11-16"
    Cuando el Scrum Master inicia el Sprint del "2026-11-10" al "2026-11-16"
    Entonces el Sprint queda activo

  Escenario: Un proyecto sin fechas no restringe el Sprint
    Dado que el proyecto no tiene fechas de plazo
    Cuando el Scrum Master inicia el Sprint del "2026-01-01" al "2026-01-15"
    Entonces el Sprint queda activo

  Escenario: Un proyecto solo con fecha de inicio no limita el fin del Sprint
    Dado que el plazo del proyecto empieza el "2026-10-31"
    Cuando el Scrum Master inicia el Sprint del "2026-12-01" al "2027-03-01"
    Entonces el Sprint queda activo

  Escenario: Un proyecto solo con fecha de fin no limita el inicio del Sprint
    Dado que el plazo del proyecto termina el "2026-11-16"
    Cuando el Scrum Master inicia el Sprint del "2026-01-01" al "2026-01-15"
    Entonces el Sprint queda activo

  Escenario: Un proyecto solo con fecha de fin rechaza un Sprint que termina después
    Dado que el plazo del proyecto termina el "2026-11-16"
    Cuando el Scrum Master inicia el Sprint del "2026-11-10" al "2026-11-20"
    Entonces el inicio se rechaza en el campo "fecha_fin"
