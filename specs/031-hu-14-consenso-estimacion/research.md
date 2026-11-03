# Research: Consenso y Registro de la Estimación Acordada (HU-14)

**Fecha**: 2026-11-03 | **Spec**: [spec.md](./spec.md)

No quedaron marcas NEEDS CLARIFICATION en el Technical Context: las dudas de negocio se
resolvieron en `/speckit.specify` y `/speckit.clarify` (ver Clarifications de la spec). Este
documento registra las decisiones técnicas del diseño.

## R1. Rutas de la API

- **Decision**: extender el recurso `votacion` de la historia:
  - `POST /backlog/{id}/votacion/rondas`: abrir una nueva ronda.
  - `GET /backlog/{id}/votacion/rondas`: consultar el historial de rondas reveladas.
  - `POST /backlog/{id}/votacion/estimacion`: registrar la estimación acordada.
  - `GET /backlog/{id}/votacion` suma el campo `diferencias`.
- **Rationale**: `rondas` es la ruta que HU-02 dejó prevista (R1 de la spec 012), y no rompe
  las rutas existentes. Las diferencias van en la consulta de la ronda porque son un dato de
  esa ronda: el equipo las ve junto con los votos y el promedio (SC-001).
- **Alternatives considered**: registrar con `PUT /backlog/{id}` (HU-15): esa edición no valida
  rol ni revelación, y mezclaría dos casos de uso; un `GET .../diferencias` aparte: obliga a una
  segunda petición para un dato que siempre se muestra junto a los votos.

## R2. Regla de consenso (Strategy)

- **Decision**: interfaz `ReglaConsenso` en `internal/domain/votacion.go`, con un método
  `HayDiferencias(valores []int) bool` y una implementación `CoincidenciaExacta` (clarify: cualquier
  voto que no coincida es una diferencia). Se inyecta en `ConsultarVotacion` y se elige en
  `cmd/api/main.go`. El mínimo y el máximo no dependen de la regla: los calcula
  `RondaVotacion.Diferencias`.
- **Rationale**: la constitución (IV) y `AGENTS.md` nombran este patrón para el consenso de Planning
  Poker, la sub-issue #91 lo pide, y R3 de la spec 012 difirió la extracción a esta historia. Se
  limita a una interfaz de un método para no sumar ceremonia.
- **Alternatives considered**: una función suelta sin interfaz: es lo más simple (YAGNI), porque hoy
  hay una sola regla, pero no cumple lo que piden la sub-issue y la constitución. Ver Complexity
  Tracking del plan.

## R3. Forma de `diferencias`

- **Decision**:
  - Ronda sin revelar: `"diferencias": null`.
  - Revelada con diferencias: `{"hay_diferencias": true, "minimo": 3, "maximo": 8}`.
  - Revelada sin diferencias: `{"hay_diferencias": false, "minimo": null, "maximo": null}`.
- **Rationale**: `null` sin revelar no da pistas sobre los votos ocultos (FR-002, SC-002). Los extremos
  se informan solo si hay diferencias (FR-001). El campo se agrega sin cambiar los existentes, así no
  rompe a quien ya consume la consulta, como la pantalla de HU-02 (PR #229).

## R4. Dónde vive cada regla

- **Decision**:
  - **Dominio**: `RondaVotacion.Siguiente` (la ronda está revelada; número + 1, sin votos),
    `RondaVotacion.Diferencias`, `ReglaConsenso`, `BacklogItem.ConEstimacion` (valor en la baraja).
  - **Servicio**: existencia de la historia, permisos, historia no `Completada` al abrir ronda,
    ronda actual revelada al registrar, D2b, y el armado del historial.
- **Rationale**: mismo criterio que R3 de la spec 012: las reglas puras se testean sin fakes y las que
  leen otros agregados quedan en el caso de uso.

## R5. Casos de uso (SRP)

- **Decision**: tres services nuevos: `AbrirNuevaRonda`, `RegistrarEstimacion` y
  `ConsultarHistorialRondas`. `ConsultarVotacion` suma las diferencias a su vista.
- **Rationale**: un service por caso de uso (constitución V).

## R6. Regla D2b compartida

- **Decision**: mover `validarEstimacionCongelada` de método de `EditarHistoria` a una función del
  paquete `service` (`estimacion_congelada.go`), que usan `EditarHistoria` y `RegistrarEstimacion`.
- **Rationale**: `MODELO-DE-DATOS.md` §5 dice que la regla es del "Servicio (HU-15 y HU-14)", y DRY pide
  un único lugar. Su comportamiento ya incluye aceptar el mismo valor (edge case de la spec). Es un
  commit `REFACTOR:` con los tests de HU-15 en verde.

## R7. Escritura de la estimación

- **Decision**: puerto nuevo `EstimacionHistoriaRepository` con `ObtenerPorID` y
  `GuardarEstimacion(ctx, id, sp)`, implementado en `SQLiteBacklogRepository` con
  `UPDATE historias SET estimacion_sp = ? WHERE id = ?`.
- **Rationale**: `HistoriaRepository.Actualizar` también reescribe título, descripción y criterios, y
  pisaría una edición simultánea de HU-15. La separación sigue el patrón de
  `AsignacionHistoriaRepository` y `CambioEstadoHistoriaRepository` (ISP).
- **Alternatives considered**: reutilizar `Actualizar` (riesgo anterior).

## R8. Concurrencia

- **Decision**:
  - **Abrir ronda**: el `UNIQUE (historia_id, numero)` frena dos aperturas simultáneas; el perdedor
    recibe `ErrVotacionYaIniciada` (409), que `CrearRonda` ya traduce.
  - **Registrar estimación**: chequeo y escritura separados, igual que `EditarHistoria`. Se acepta la
    ventana de milisegundos en la que un Sprint se inicia entre el chequeo de D2b y la escritura, o se
    abre una ronda entre el chequeo de revelación y la escritura.
- **Rationale**: es el mismo criterio que la edición de HU-15 y que R6 de la spec 012 (equipo chico). Un
  `UPDATE` condicional repetiría la regla D2b en SQL, y se escribe un solo campo.
- **Alternatives considered**: transacción en el servicio (el servicio no conoce transacciones).

## R9. Historial de rondas

- **Decision**: `VotacionRepository.ListarRondasReveladas(ctx, historiaID)` devuelve las rondas con
  `revelada_en` cargada, ordenadas por número, con sus votos, en **una sola consulta con `LEFT JOIN`**.
  Los nombres y roles salen de `ProjectRepository.ListarIntegrantes`.
- **Rationale**: la conexión es única (`SetMaxOpenConns(1)`): reutilizar `votosDe` con las filas de las
  rondas abiertas bloquearía la consulta. Hoy no se pueden quitar integrantes, así que todo votante está
  en la lista. El historial no pide `usuario_id` porque solo muestra votos ya revelados (Assumptions).
- **Alternatives considered**: `ListarRondas` con todas las rondas y filtrar en el servicio: traería votos
  ocultos a una consulta que no debe exponerlos.

## R10. Errores y códigos HTTP

| Error | Tipo | HTTP | Cuándo |
|---|---|---|---|
| `ErrHistoriaNoEncontrada` (existente) | sentinel | 404 | historia inexistente (primero de todo) |
| `ErrVotacionNoIniciada` (existente) | sentinel | 404 | abrir ronda o registrar sin votación |
| `ErrNoEsScrumMaster` (existente) | sentinel | 403 | abre o registra quien no es Scrum Master |
| `ErrNoEsIntegrante` (existente) | sentinel | 403 | abre o registra quien no integra el proyecto |
| `ErrVotacionYaIniciada` (existente) | sentinel | 409 | dos aperturas simultáneas de la misma ronda |
| `ErrRondaNoRevelada` (**nuevo**) | `ValidationError{estado}` | 400 | abrir ronda o registrar con la ronda actual oculta |
| `ErrHistoriaCompletada` (existente) | `ValidationError{estado}` | 400 | abrir ronda en una historia `Completada` |
| `ErrEstimacionInvalida` (existente) | `ValidationError{estimacion_sp}` | 400 | valor fuera de la baraja |
| `ErrEstimacionCongelada` (existente) | `ValidationError{estimacion_sp}` | 400 | D2b |

- **Rationale**: se reutilizan los errores existentes y `escribirError` no cambia. El mensaje de
  `ErrHistoriaCompletada` ("no se puede votar una historia completada") también sirve para abrir una ronda.
- **Orden de validación**: id de ruta → historia existe (404) → JSON → `usuario_id` → permiso (403) →
  estado de la historia (solo al abrir ronda) → valor → votación iniciada (404) → ronda revelada (400) → D2b.

## R11. Esquema

- **Decision**: sin migración y sin cambios en `MODELO-DE-DATOS.md`.
- **Rationale**: §4.4 y §4.7 ya cubren todo: las rondas anteriores son el historial, la acordada se escribe
  en `historias.estimacion_sp`, y §5 asigna las reglas a HU-14.
