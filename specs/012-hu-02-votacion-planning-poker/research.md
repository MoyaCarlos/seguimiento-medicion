# Research: Votación de Esfuerzo en Planning Poker (HU-02)

**Fecha**: 2026-10-27 | **Spec**: [spec.md](./spec.md)

No quedaron marcas NEEDS CLARIFICATION en el Technical Context: las cinco dudas de
negocio se resolvieron en `/speckit.specify` y `/speckit.clarify` (ver Clarifications de la
spec). Este documento registra las decisiones técnicas del diseño.

## R1. Rutas de la API

- **Decision**: la votación cuelga de la historia, como recurso único por historia:
  - `POST /backlog/{id}/votacion` — iniciar (abre la ronda 1).
  - `POST /backlog/{id}/votacion/votos` — votar o cambiar el voto.
  - `POST /backlog/{id}/votacion/revelar` — revelar.
  - `GET /backlog/{id}/votacion?usuario_id=N` — consultar la ronda actual.
- **Rationale**: las historias ya viven en `/backlog/{id}` (HU-01, HU-15). Las acciones de
  estado usan un verbo como subrecurso, igual que `POST /sprints/{id}/iniciar` y `/cerrar`.
  "La votación" en singular expresa la ronda actual; HU-14 podrá agregar
  `POST /backlog/{id}/votacion/rondas` sin romper estas rutas.
- **Alternatives considered**: `/rondas/{id}/votos` (obliga al cliente a conocer el id de la
  ronda, que en HU-02 no aporta nada); `PUT` para el voto (se descartó para no mezclar
  estilos: el resto de las acciones del proyecto son `POST`).

## R2. Identificación de quien actúa (sin autenticación)

- **Decision**: el cuerpo de cada `POST` lleva `usuario_id`; la consulta lo recibe como
  parámetro opcional `?usuario_id=`. El rol se lee de los integrantes del proyecto de la
  historia con `ProjectRepository.ListarIntegrantes` (puerto existente, HU-04).
- **Rationale**: la API no tiene sesión; `usuario_id` es el `id` que ya devuelve
  `GET /projects/{id}/members`. Sin `usuario_id`, la consulta muestra todos los valores
  ocultos: no hay voto propio que mostrar y no exponerse nada es el comportamiento seguro.
- **Alternatives considered**: identificar por nombre (como el alta de integrantes): obliga a
  normalizar y buscar en cada voto; un header `X-Usuario-Id`: un mecanismo nuevo que nadie
  más usa en la API.

## R3. Dónde vive cada regla

- **Decision**:
  - **Dominio** (`internal/domain/votacion.go`): `RondaVotacion` con `Votar` (ronda no
    revelada, valor en la baraja, reemplazo del voto propio) y `Revelar` (no revelada, votaron
    todos los integrantes recibidos); `CalcularPromedio` (exacto a 2 decimales y carta más
    cercana, empate hacia la mayor). Reusa `EsEstimacionValida` (misma escala, DRY).
  - **Servicio**: existencia de la historia, historia no `Completada`, permisos (integrante /
    Scrum Master), una sola votación por historia y el ocultamiento en la consulta
    (`MODELO-DE-DATOS.md` §4.7: "esa regla es del servicio/handler").
- **Rationale**: las reglas puras (estado de la ronda, cálculo) se testean sin fakes; las que
  necesitan leer otros agregados quedan en el caso de uso, como D2b en `EditarHistoria`.
- **Alternatives considered**: Strategy para el cálculo del promedio (constitución IV la
  menciona para el consenso). Hoy hay una sola forma de calcularlo: YAGNI. Si HU-14 necesita
  otra regla de consenso, se extrae ahí.

## R4. Casos de uso (SRP)

- **Decision**: cuatro services: `IniciarVotacion`, `Votar`, `RevelarVotacion`,
  `ConsultarVotacion`. Un helper privado `rolEnProyecto` comparte la búsqueda del integrante.
- **Rationale**: un service por caso de uso (constitución V), como `CrearSprint`,
  `IniciarSprint` y `CerrarSprint`.

## R5. Puerto de persistencia

- **Decision**: interfaz nueva `VotacionRepository` en `internal/repository/votacion_repository.go`
  con `CrearRonda`, `ObtenerRondaActual`, `GuardarVoto` y `MarcarRevelada`; adaptador
  `SQLiteVotacionRepository`. Los tests de service usan un `fakeVotacionRepository` en
  `internal/service/fakes_test.go` que simula las escrituras condicionales de R6 (un fake sin
  esa regla ocultaría defectos, punto 2 de #101); el BDD usa SQLite `:memory:`, como HU-15.
  No se agrega un fake de producción en `internal/repository` porque ninguna otra historia lo
  necesita todavía (YAGNI). Ver [data-model.md](./data-model.md).
- **Rationale**: ISP, igual que `HistoriaRepository` / `ProductBacklogRepository`. La ronda
  actual es la de mayor `numero` (§4.7), resuelto con `ORDER BY numero DESC LIMIT 1`.

## R6. Concurrencia (diferido desde clarify)

- **Decision**: las escrituras sobre una ronda son **condicionales** en el adaptador:
  - Voto: `INSERT INTO votos (...) SELECT ?, ?, ? WHERE EXISTS (SELECT 1 FROM rondas_votacion
    WHERE id = ? AND revelada_en IS NULL) ON CONFLICT (ronda_id, usuario_id) DO UPDATE SET
    valor = excluded.valor`; 0 filas afectadas ⇒ `ErrRondaRevelada`.
  - Revelación: `UPDATE rondas_votacion SET revelada_en = ? WHERE id = ? AND revelada_en IS
    NULL`; 0 filas ⇒ `ErrRondaRevelada`.
  - Inicio: el `UNIQUE (historia_id, numero)` frena dos inicios simultáneos; el adaptador
    traduce el 2067 a `ErrVotacionYaIniciada` (tabla de traducción de §5: error propio, 409).
- **Rationale**: un voto que llega justo después de revelar cambiaría un resultado ya visto;
  la condición en la misma sentencia lo evita sin transacciones explícitas. El caso "se suma
  un integrante entre el chequeo de 'votaron todos' y la revelación" se acepta (ventana de
  milisegundos, equipo chico), con el mismo criterio que el chequeo de Sprint activo de
  `IniciarSprint`.
- **Alternatives considered**: transacción con relectura en el servicio (el servicio no
  conoce transacciones, arquitectura hexagonal); no proteger nada (un voto tardío alteraría
  una ronda revelada, que es justo lo que FR-010 prohíbe).

## R7. Instante de revelación

- **Decision**: el handler pasa `time.Now()` al service (`ahora time.Time`), como
  `ObtenerEstadoProyecto`; se guarda en RFC 3339 UTC (§2, instantes).
- **Rationale**: tests deterministas sin un reloj inyectado nuevo.

## R8. Errores y códigos HTTP

| Error | Tipo | HTTP | Cuándo |
|---|---|---|---|
| `ErrHistoriaNoEncontrada` (existente) | sentinel | 404 | historia inexistente (primero de todo) |
| `ErrVotacionNoIniciada` | sentinel | 404 | votar, revelar o consultar sin ronda |
| `ErrNoEsIntegrante` | sentinel | 403 | vota quien no es integrante del proyecto |
| `ErrNoEsScrumMaster` | sentinel | 403 | inicia o revela quien no es `scrum_master` |
| `ErrVotacionYaIniciada` | sentinel | 409 | la historia ya tiene ronda (servicio o `UNIQUE`) |
| `ErrVotoInvalido` | `ValidationError{valor}` | 400 | carta fuera de la baraja |
| `ErrHistoriaCompletada` | `ValidationError{estado}` | 400 | iniciar sobre historia `Completada` |
| `ErrRondaRevelada` | `ValidationError{estado}` | 400 | votar o revelar una ronda revelada |
| `ErrFaltanVotos` | `ValidationError{votos}` | 400 | revelar sin que voten todos |
| `usuario_id` ausente o ≤ 0 | `ValidationError{usuario_id}` | 400 | cuerpo sin quien actúa |

- **Rationale**: las reglas de estado son `ValidationError` → 400, como "solo se puede
  iniciar un Sprint pendiente". 403 es nuevo en la API: es el código que corresponde a "sos
  quien decís, pero no tenés permiso"; se agrega en `escribirError`. 409 para el `UNIQUE`,
  como fija la tabla de traducción de `MODELO-DE-DATOS.md` §5.
- **Orden de validación**: id de ruta → historia existe (404) → JSON → `usuario_id` →
  permiso (403) → reglas de estado y de valor. Así la historia inexistente domina (FR-011),
  con el mismo patrón que `SprintHandler.Iniciar`.

## R9. Migración 003 y su orden respecto de la 002

- **Decision**: `internal/repository/migraciones/003_votacion.sql` crea `rondas_votacion` y
  `votos` tal como §4.7 (número reservado para HU-02 en §9). No modifica `MODELO-DE-DATOS.md`
  porque no cambia el esquema documentado.
- **Rationale**: la 002 (#117) puede no estar mergeada cuando entre HU-02; el runner aplica
  las que falten aunque haya huecos (§9, punto 1), así que el orden de merge no importa. La 003
  no depende de `historias_arrastradas`.
- **Test de esquema**: se amplía `esquema_test.go` (§9, punto 6) con estructura, FK, PK y
  filas válidas e inválidas de las dos tablas.

## R10. Promedio: representación

- **Decision**: `Promedio{Exacto float64; CartaSugerida int}`; `Exacto =
  math.Round(suma/n*100) / 100`. JSON: `"promedio": {"exacto": 5.33, "carta_sugerida": 5}`,
  `null` mientras la ronda no esté revelada.
- **Rationale**: la revelación exige que voten todos (FR-015) y cada integrante vota una vez,
  así que `n ≥ 1` siempre que hay promedio (un proyecto siempre tiene al menos su Scrum
  Master, HU-04). La carta sugerida se busca sobre la misma escala Fibonacci del dominio.
