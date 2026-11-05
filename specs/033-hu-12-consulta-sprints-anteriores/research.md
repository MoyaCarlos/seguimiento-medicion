# Research: Consulta de Sprints Anteriores (HU-12)

Decisiones técnicas de la Fase 0. No quedan `NEEDS CLARIFICATION`.

## R1. Origen de las métricas

- **Decision**: los services de HU-12 reciben el cálculo de HU-03 a través de una interfaz angosta declarada en `internal/service`, y no leen fuentes de métricas por su cuenta.
- **Rationale**: FR-012 y SC-003 exigen los mismos valores que HU-03, y la constitución pide que el cálculo viva en un único lugar. La interfaz permite probar con un fake.
- **Alternatives considered**: calcular las métricas también en HU-12 (duplica la regla y puede divergir); guardarlas al cerrar el Sprint (contradice `MODELO-DE-DATOS.md` §6).

## R2. Dependencia de HU-03

- **Decision**: `/speckit-implement` esperaba el merge de HU-03. **Resuelto (2026-11-05)**: HU-03 ya está en `main`. Los services de HU-12 dependen de una interfaz angosta con la firma de `service.CalcularMetricas.Ejecutar(ctx, proyectoID) (domain.MetricasProyecto, error)` y toman las métricas de cada Sprint de `MetricasProyecto.Sprints` por `SprintID`. El objeto `metricas` de [openapi.yaml](./contracts/openapi.yaml) coincide campo por campo con `metricasResponse` de HU-03 (`internal/http/metricas_dto.go`) y se arma con su mismo mapeo (`aMetricasResponse`).
- **Rationale**: implementar antes del merge obligaba a duplicar el cálculo o a inventar un contrato que HU-03 no asumió.
- **Alternatives considered**: implementar primero con métricas simuladas (el producto no funcionaría hasta HU-03); limitar HU-12 a los datos sin métricas (incumple el Escenario 1).

## R3. Rutas

- **Decision**: `GET /projects/{id}/historial` y `GET /projects/{id}/historial/{sprintId}`. La SPA usa las mismas URL.
- **Rationale**: incluir el proyecto en la ruta del detalle permite rechazar un Sprint de otro proyecto (edge case de la spec). `/projects` ya está en `PrefijosAPI` y en el proxy de Vite, y `ServirSPA` distingue página y API por el encabezado `Accept`, igual que en `/projects/{id}/sprints`.
- **Alternatives considered**: `GET /sprints/{id}/historial` (no expresa el proyecto consultado); ampliar `GET /projects/{id}/sprints` con un filtro (mezcla una consulta simple con métricas y detalle).

## R4. Sprint no finalizado y errores

- **Decision**: `ErrSprintNoFinalizado = ValidationError{Campo: "estado", Mensaje: "el Sprint todavía no está finalizado"}`, que responde 400. Orden de evaluación: identificador inválido (400), proyecto inexistente (404), Sprint inexistente o de otro proyecto (404), Sprint no finalizado (400).
- **Rationale**: Clarification 3, con el mismo patrón que `ErrSprintNoActivo`. Se reutiliza `ErrSprintNoEncontrado`, cuyo mensaje vigente es "sprint no encontrado": es el error que la spec llama "Sprint no encontrado".
- **Alternatives considered**: 409 (el proyecto usa 400 para reglas de estado); 404 (descartado en la clarificación).

## R5. Historias comprometidas

- **Decision**: la consulta une las historias `Completada` con `sprint_id = S` y las de `historias_arrastradas` con `sprint_id = S`, y toma los Story Points de la foto en estas últimas. Orden: completadas primero, luego arrastradas, cada grupo por `id`.
- **Rationale**: es la definición de `MODELO-DE-DATOS.md` §4.6 (D2c). Una historia no aparece dos veces porque un Sprint finalizado no admite nuevas asignaciones.
- **Riesgo**: HU-03 también aplica esta definición para sus sumas. Si su plan expone una consulta de comprometidas, se reutiliza en lugar de `HistoriasComprometidas`. En cualquier caso, un escenario BDD compara el detalle con las métricas (FR-010).
- **Alternatives considered**: reconstruir desde `GET /sprints/{id}/backlog` (no incluye las arrastradas).

## R6. Lectura sin transacción

- **Decision**: el detalle hace lecturas separadas, sin transacción.
- **Rationale**: en un Sprint finalizado no cambian las historias ni las horas. Solo podría cambiar el estado de un defecto entre lecturas, y eso no altera los conteos del Sprint consultado de forma relevante en una aplicación de uso local. Un error en cualquier lectura devuelve 500 sin datos parciales.
- **Alternatives considered**: una transacción de lectura (agrega complejidad sin una necesidad demostrada).

## R7. Horas

- **Decision**: la respuesta incluye `minutos` y `horas` con dos decimales por historia, calculadas por separado a partir de los minutos. El total del Sprint es el `horas_reales` de HU-03, no la suma de valores redondeados. Se reutiliza el redondeo de HU-03 (su FR-011).
- **Rationale**: D8 evita errores acumulados. Con dos decimales, las horas por historia pueden diferir en centésimas de su suma; los minutos sí suman exacto y permiten verificar FR-010.
- **Alternatives considered**: sumar horas redondeadas (podría no coincidir con HU-03).

## R8. Métricas en lista y detalle

- **Decision**: la API devuelve el mismo objeto `metricas` en ambas respuestas. La pantalla de la lista muestra las métricas de FR-003 y la de detalle, todas.
- **Rationale**: un solo DTO evita mantener dos formas de un mismo dato. La velocidad de un Sprint es igual a sus Story Points completados, así que omitirla en la lista no pierde información.
- **Alternatives considered**: un DTO reducido para la lista (más código sin beneficio, ya que HU-03 calcula todo igual).
