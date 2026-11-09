# Research: Consultar una historia por su identificador

No quedan `NEEDS CLARIFICATION`. Un único punto estaba abierto en la spec (Assumptions): cómo cumplir el mensaje `identificador inválido` sin alterar los demás endpoints.

## D1. Validación del id y texto del 400

- **Decision**: el handler nuevo usa `idDeRuta` (`internal/http/sprint_handler.go`). No se modifica `parsearID` ni el texto de ningún endpoint existente.
- **Rationale**: `idDeRuta` ya responde 400 con `{"campo":"id","mensaje":"identificador inválido"}`, exactamente lo que pide el issue. Rechaza `abc`, `0`, negativos, decimales y desbordes de `int64`, y `strconv.ParseInt` interpreta `007` como 7. Reutilizarlo es DRY y no cambia comportamiento previo.
- **Alternatives considered**:
  - Usar `parsearID` (texto `el identificador debe ser un entero positivo`): incumple el texto pedido.
  - Cambiar el texto de `parsearID`: altera todos los endpoints y sus tests; fuera de alcance.
  - Un helper nuevo: duplicaría `idDeRuta`.

## D2. Caso de uso y traducción de errores

- **Decision**: usar `obtenerHistoria` (`service.ObtenerHistoria`) que ya tiene `BacklogHandler`, y `escribirError` sin `case` nuevos.
- **Rationale**: `escribirError` ya mapea `domain.ErrHistoriaNoEncontrada` a 404 con su mensaje `historia no encontrada`. Cumple FR-006 y FR-010.
- **Alternatives considered**: consultar el repositorio directo desde el handler (rompe la separación hexagonal).

## D3. Forma de la respuesta

- **Decision**: `aRespuesta(item)` → `historiaResponse`.
- **Rationale**: es la misma función que usan el listado del backlog, la edición y las asignaciones; garantiza FR-003 (mismos campos y formato de fecha) y SC-001.
- **Alternatives considered**: DTO propio (riesgo de divergencia).

## D4. Ruta

- **Decision**: `GET /backlog/{id}` registrado junto a `PUT /backlog/{id}` en `cmd/api/main.go`.
- **Rationale**: Go 1.22+ distingue por método y por cantidad de segmentos; no choca con `GET /backlog/{id}/horas`.
- **Alternatives considered**: `GET /stories/{id}` (rompe la convención de la API).

## D5. Estrategia de pruebas

- **Decision**: tests de handler con el `stubBacklogRepository` existente y escenarios Godog sobre SQLite en memoria, con el mismo patrón que `steps_consultar_sprints.go`.
- **Rationale**: Principios I y III. El RED con un `Obtener` vacío compila y falla por aserción, no por compilación.
- **Alternatives considered**: escenarios con fakes en memoria: no ejercitarían el mapeo real de fechas y `sprint_id` nulos.

## D6. Datos y documentación

- **Decision**: sin migración ni cambios en `MODELO-DE-DATOS.md`; ejemplos en `docs/curl-frontend-sprint-2.md` (bash y cmd).
- **Rationale**: la historia ya tiene todos los campos; solo se lee.
