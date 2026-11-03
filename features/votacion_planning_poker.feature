# language: es

# Un escenario por Acceptance Scenario de specs/012-hu-02-votacion-planning-poker/spec.md.
# Cobertura por categoría:
# - Normal: votación secreta (US1.1); revelación de votos (US2.1); inicio de la votación (US3.1).
# - Alternativo: cambio de voto antes de revelar (US1.2).
# - Límite: revelar cuando falta votar un integrante (US2.3).
# - Error: carta fuera de la baraja (US1.3); voto de quien no es integrante (US1.4); voto sobre
#   una ronda revelada (US2.2); revelación por un Product Builder (US2.4); segundo inicio (US3.2);
#   historia inexistente (US3.3); inicio por un Product Builder (US3.4); historia Completada (US3.5).
# Los empates del promedio sugerido no son criterios de la spec: los cubren los tests de dominio.

@HU-02
Característica: Votación de esfuerzo en Planning Poker

  Como integrante del equipo (Scrum Master o Product Builder)
  Quiero votar el esfuerzo de una historia en secreto usando cartas virtuales
  Para poder realizar una estimación colaborativa sin sesgos

  Antecedentes:
    Dado que el proyecto de la votación tiene a "Carlos" como Scrum Master y a "Candela" y "Jimena" como Product Builders
    Y que el proyecto de la votación tiene una historia en estado "Nueva"

  Escenario: Votación secreta
    Dado que "Carlos" inició la votación de la historia
    Cuando "Candela" selecciona la carta 8 y confirma su voto
    Entonces el voto de "Candela" queda registrado en la ronda 1
    Y "Carlos" ve que "Candela" votó, con el valor "Oculto"
    Y "Candela" ve su propia carta 8

  Escenario: Cambio de voto antes de revelar
    Dado que "Carlos" inició la votación de la historia
    Y que "Candela" votó la carta 8
    Cuando "Candela" selecciona la carta 13 y confirma su voto
    Entonces "Candela" tiene un solo voto en la ronda, con la carta 13
    Y "Carlos" ve que "Candela" votó, con el valor "Oculto"

  Escenario: Voto con una carta que no pertenece a la baraja
    Dado que "Carlos" inició la votación de la historia
    Cuando "Candela" selecciona la carta 4 y confirma su voto
    Entonces el voto se rechaza con una advertencia sobre el campo "valor"
    Y la ronda de votación no cambia

  Escenario: Voto de una persona que no es integrante del proyecto
    Dado que "Carlos" inició la votación de la historia
    Cuando "Diego", que no es integrante del proyecto, selecciona la carta 8 y confirma su voto
    Entonces el voto se rechaza porque "Diego" no es integrante del proyecto
    Y la ronda de votación no cambia

  Escenario: Revelación de votos
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Cuando "Carlos" presiona "Revelar"
    Entonces se ve la carta elegida por cada integrante:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y el promedio sugerido es 5,33, con la carta sugerida 5

  Escenario: Cambio de voto en una ronda ya revelada
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Cuando "Candela" selecciona la carta 13 y confirma su voto
    Entonces el voto se rechaza con una advertencia sobre el campo "estado"
    Y los votos revelados no cambian

  Escenario: Revelación cuando todavía falta votar un integrante
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 8     |
    Cuando "Carlos" presiona "Revelar"
    Entonces la revelación se rechaza con una advertencia sobre el campo "votos"
    Y los votos siguen ocultos

  Escenario: Revelación por un Product Builder
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Cuando "Candela" presiona "Revelar"
    Entonces la revelación se rechaza porque "Candela" no es Scrum Master del proyecto
    Y los votos siguen ocultos

  Escenario: Inicio de la votación de una historia
    Dado que la historia todavía no tiene votación
    Cuando "Carlos" inicia la votación de la historia
    Entonces la historia tiene abierta la ronda 1, sin votos y sin revelar

  Escenario: Inicio de una votación que ya existe
    Dado que "Carlos" inició la votación de la historia
    Cuando "Carlos" inicia la votación de la historia
    Entonces el inicio de la votación se rechaza porque la historia ya tiene una votación
    Y la historia sigue con una sola ronda de votación

  Escenario: Inicio de la votación de una historia inexistente
    Cuando "Carlos" inicia la votación de la historia con identificador 999, que no existe
    Entonces el inicio de la votación se rechaza con error "historia no encontrada"

  Escenario: Inicio de la votación por un Product Builder
    Dado que la historia todavía no tiene votación
    Cuando "Candela" inicia la votación de la historia
    Entonces el inicio de la votación se rechaza porque "Candela" no es Scrum Master del proyecto
    Y la historia sigue sin rondas de votación

  Escenario: Inicio de la votación de una historia completada
    Dado que la historia está en estado "Completada" y todavía no tiene votación
    Cuando "Carlos" inicia la votación de la historia
    Entonces el inicio de la votación se rechaza con una advertencia sobre el campo "estado"
    Y la historia sigue sin rondas de votación
