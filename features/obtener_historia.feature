# language: es

@Issue-222
Característica: Obtener una historia por su identificador

  Como miembro del equipo
  Quiero consultar una historia por su identificador
  Para ver su detalle completo

  Escenario: Historia del Product Backlog
    Dado que existe la historia 1 en el Product Backlog
    Cuando consulto la historia "1"
    Entonces la consulta de la historia responde 200
    Y la historia recibida es la 1 y tiene sprint_id nulo

  Escenario: Historia asignada a un Sprint
    Dado que existe la historia 1 asignada al Sprint 7
    Cuando consulto la historia "1"
    Entonces la consulta de la historia responde 200
    Y la historia recibida es la 1 y tiene sprint_id 7

  Escenario: Historia completada
    Dado que existe la historia 1 completada el "2026-10-19" en el Sprint 7
    Cuando consulto la historia "1"
    Entonces la consulta de la historia responde 200
    Y la historia recibida tiene fecha_finalizacion "2026-10-19"

  Escenario: Historia inexistente
    Cuando consulto la historia "999"
    Entonces la consulta de la historia responde 404
    Y el mensaje de la consulta de la historia es "historia no encontrada"

  Escenario: Identificador inválido
    Cuando consulto la historia "abc"
    Entonces la consulta de la historia responde 400
    Y el error de la consulta de la historia es del campo "id" con mensaje "identificador inválido"
