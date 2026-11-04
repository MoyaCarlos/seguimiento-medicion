# language: es

# Un escenario por Acceptance Scenario de specs/031-hu-14-consenso-estimacion/spec.md.
# Cobertura por categoría:
# - Normal: diferencias (US1.1); registro de la estimación (US2.1); nueva ronda (US3.1).
# - Alternativo: sin diferencias (US1.2); reemplazo de una estimación cargada (US2.2); carta que
#   nadie votó (US2.8); votar y revelar la ronda nueva (US3.2); nueva ronda sin diferencias (US3.7);
#   historial (US3.8).
# - Límite: consulta de una ronda oculta (US1.3); primera carga con Sprint Activo (US2.3); primera
#   carga en una historia Completada (US2.9).
# - Error: carta fuera de la baraja (US2.4); ronda oculta (US2.5); D2b (US2.6); permisos (US2.7, US3.4);
#   ronda oculta (US3.3); sin votación (US3.5); historia Completada (US3.6).
# La historia inexistente no es un criterio de la spec: la cubren los tests de service y handler
# (la de registrar la estimación, los del handler, que es quien lee la historia).

@HU-14
Característica: Consenso y registro de la estimación acordada

  Como Scrum Master
  Quiero que, al revelar los votos de Planning Poker, el sistema detecte las diferencias entre las estimaciones, permita hacer nuevas rondas y registre la estimación acordada en la historia
  Para poder obtener los Story Points sobre los que se calculan las métricas

  Antecedentes:
    Dado que el proyecto de la votación tiene a "Carlos" como Scrum Master y a "Candela" y "Jimena" como Product Builders
    Y que el proyecto de la votación tiene una historia en estado "Nueva"

  Escenario: Votos revelados con diferencias
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Cuando se consulta la votación de la historia
    Entonces el sistema indica que hay diferencias, con mínimo 3 y máximo 8

  Escenario: Votos revelados sin diferencias
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 5     |
      | Candela    | 5     |
      | Jimena     | 5     |
    Y que "Carlos" reveló los votos
    Cuando se consulta la votación de la historia
    Entonces el sistema indica que no hay diferencias

  Escenario: Consulta de una ronda sin revelar
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 8     |
    Cuando se consulta la votación de la historia
    Entonces el sistema no informa diferencias, mínimo ni máximo

  Escenario: Registro de la estimación acordada
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Cuando "Carlos" registra el valor acordado 5
    Entonces la historia queda con 5 Story Points

  Escenario: Registro que reemplaza una estimación cargada de forma directa
    Dado que la historia tiene 3 Story Points cargados
    Y que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Cuando "Carlos" registra el valor acordado 8
    Entonces la historia queda con 8 Story Points

  Escenario: Primera carga de la estimación con un Sprint Activo
    Dado que la historia está asignada a un Sprint "Activo"
    Y que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Cuando "Carlos" registra el valor acordado 5
    Entonces la historia queda con 5 Story Points

  Escenario: Estimación con un valor que no pertenece a la baraja
    Dado que la historia tiene 3 Story Points cargados
    Y que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Cuando "Carlos" registra el valor acordado 4
    Entonces la estimación se rechaza con una advertencia sobre el campo "estimacion_sp"
    Y la historia mantiene sus 3 Story Points

  Escenario: Estimación con la ronda sin revelar
    Dado que la historia tiene 3 Story Points cargados
    Y que "Carlos" inició la votación de la historia
    Cuando "Carlos" registra el valor acordado 5
    Entonces la estimación se rechaza con una advertencia sobre el campo "estado"
    Y la historia mantiene sus 3 Story Points

  Esquema del escenario: Cambio de una estimación cargada con un Sprint <estado>
    Dado que la historia tiene 3 Story Points cargados
    Y que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Y que la historia está asignada a un Sprint "<estado>"
    Cuando "Carlos" registra el valor acordado 5
    Entonces la estimación se rechaza con una advertencia sobre el campo "estimacion_sp"
    Y la historia mantiene sus 3 Story Points

    Ejemplos:
      | estado     |
      | Activo     |
      | Finalizado |

  Escenario: Estimación registrada por un Product Builder
    Dado que la historia tiene 3 Story Points cargados
    Y que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Cuando "Candela" registra el valor acordado 5
    Entonces la estimación se rechaza porque "Candela" no es Scrum Master del proyecto
    Y la historia mantiene sus 3 Story Points

  Escenario: Estimación con una carta que nadie votó
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 8     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Cuando "Carlos" registra el valor acordado 5
    Entonces la historia queda con 5 Story Points

  Escenario: Primera carga de la estimación en una historia completada
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Y que la historia pasó a estado "Completada"
    Cuando "Carlos" registra el valor acordado 5
    Entonces la historia queda con 5 Story Points

  Escenario: Nueva ronda con diferencias
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Cuando "Carlos" inicia una nueva ronda
    Entonces la ronda actual es la 2, abierta y sin votos
    Y la ronda 1 conserva sus votos en el historial:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |

  Escenario: Votar y revelar la ronda nueva
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 13    |
    Y que "Carlos" reveló los votos
    Y que "Carlos" inició una nueva ronda
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 5     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Cuando "Carlos" presiona "Revelar"
    Entonces el sistema indica que hay diferencias, con mínimo 5 y máximo 8
    Y el promedio sugerido es 6,00, con la carta sugerida 5

  Escenario: Nueva ronda con la ronda actual sin revelar
    Dado que "Carlos" inició la votación de la historia
    Cuando "Carlos" inicia una nueva ronda
    Entonces la nueva ronda se rechaza con una advertencia sobre el campo "estado"
    Y la ronda de votación no cambia

  Escenario: Nueva ronda iniciada por un Product Builder
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Cuando "Candela" inicia una nueva ronda
    Entonces la nueva ronda se rechaza porque "Candela" no es Scrum Master del proyecto
    Y la historia sigue con una sola ronda de votación

  Escenario: Nueva ronda de una historia sin votación
    Dado que la historia todavía no tiene votación
    Cuando "Carlos" inicia una nueva ronda
    Entonces la nueva ronda se rechaza porque la votación no está iniciada
    Y la historia sigue sin rondas de votación

  Escenario: Nueva ronda de una historia completada
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Y que la historia pasó a estado "Completada"
    Cuando "Carlos" inicia una nueva ronda
    Entonces la nueva ronda se rechaza con una advertencia sobre el campo "estado"
    Y la historia sigue con una sola ronda de votación

  Escenario: Nueva ronda sin diferencias
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 5     |
      | Candela    | 5     |
      | Jimena     | 5     |
    Y que "Carlos" reveló los votos
    Cuando "Carlos" inicia una nueva ronda
    Entonces la ronda actual es la 2, abierta y sin votos

  Escenario: Consulta del historial de rondas
    Dado que "Carlos" inició la votación de la historia
    Y que votaron estas cartas:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
    Y que "Carlos" reveló los votos
    Y que "Carlos" inició una nueva ronda
    Y que votaron estas cartas:
      | integrante | carta |
      | Candela    | 13    |
    Cuando se consulta el historial de rondas de la historia
    Entonces el historial tiene solo la ronda 1
    Y la ronda 1 conserva sus votos en el historial:
      | integrante | carta |
      | Carlos     | 3     |
      | Candela    | 5     |
      | Jimena     | 8     |
