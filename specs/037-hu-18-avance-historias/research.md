# Research: Avance de las Historias del Proyecto (HU-18)

Decisiones de la Fase 0. No quedan `NEEDS CLARIFICATION`.

## R1. Un endpoint de lectura nuevo

- **Decision**: `GET /projects/{id}/avance`, que devuelve la agrupación ya armada.
- **Rationale**: ningún endpoint actual lista todas las historias de un proyecto: `GET /projects/{id}/backlog` (HU-15) solo trae las sin Sprint y `GET /sprints/{id}/backlog` trae las de un Sprint. Armar la vista en el frontend obligaría a duplicar agrupación y sumas (viola DRY, `AGENTS.md`) y a hacer una consulta por Sprint.
- **Alternatives considered**: ampliar `/projects/{id}/backlog` con un parámetro (cambia el contrato de HU-15, que debe seguir mostrando solo las sin Sprint); componer en el cliente `backlog` + `sprints` + `sprints/{id}/backlog` (N+1 consultas y lógica duplicada); `/projects/{id}/status` (colisiona con la ruta de la SPA; `/projects` ya está en `PrefijosAPI` y la distinción por `Accept` es un riesgo innecesario).

## R2. Agrupación y sumas en el dominio

- **Decision**: función pura `domain.AgruparAvance(historias, sprints)` que devuelve `AvanceHistorias` con los tres grupos siempre presentes, en el orden Nueva, En progreso, Completada. Cada grupo trae `Cantidad` y `StoryPoints` (suma, con `nil` = 0). El DTO solo copia.
- **Rationale**: única fuente de verdad (DRY) y testeable con TDD sin base de datos. Las historias conservan el orden recibido (creación).
- **Alternatives considered**: `GROUP BY` en SQL (la regla "sin estimar suma 0" quedaría en el adaptador y no en el dominio); sumar en el frontend (duplica la lógica).

## R3. Unión con los Sprints en el service

- **Decision**: el service lee las historias con un puerto nuevo y los Sprints con `listadorSprints.ListarPorProyecto` (ya existe, lo implementa `sprintRepo` y entrega el `Numero` derivado). `AgruparAvance` resuelve `SprintID → Numero`. Si una historia apunta a un Sprint que no está en la lista, las fuentes son inconsistentes: error interno (500), igual que `metricasDelSprint` en HU-12.
- **Rationale**: reutiliza el cálculo de `Numero` (#129) en vez de repetirlo en SQL.
- **Alternatives considered**: `JOIN` con `ROW_NUMBER()` en la consulta (duplica la regla del número); devolver solo `sprint_id` y que el frontend busque el número (más consultas).

## R4. Puerto nuevo, método en el adaptador existente

- **Decision**: puerto `AvanceHistoriasRepository` con `ListarPorProyecto(ctx, proyectoID) ([]domain.BacklogItem, error)`, ordenado por `id`, implementado en `SQLiteBacklogRepository` reutilizando `listarHistorias` y `columnasHistoria`. Los slices vacíos no son `nil`.
- **Rationale**: ISP, como `ProductBacklogRepository`; ningún puerto actual lista todas las historias del proyecto. No hace falta un adaptador nuevo: la tabla y el escaneo ya están ahí.
- **Alternatives considered**: ampliar `ProductBacklogRepository` (cambia su significado "sin Sprint"); un adaptador SQLite aparte (duplica `escanearHistoria`).

## R5. Orden de errores

- **Decision**: id inválido (400, campo `id`, con `parsearID`) → proyecto inexistente (404, `ErrProyectoNoEncontrado` del repositorio de proyectos) → lectura. Un proyecto sin historias responde 200 con los tres grupos en cero.
- **Rationale**: misma convención de HU-12 y de `ConsultarProductBacklog`. Clarification 7.

## R6. Forma de la respuesta

- **Decision**: objeto `{ total, grupos: [ { estado, cantidad, story_points, historias: [...] } ] }`. `total` trae `historias` y `story_points` del proyecto. Cada historia: `id`, `titulo`, `estimacion_sp` (`null` si no está estimada), `estado`, `sprint` (`{id, numero}` o `null`). Arreglos nunca `null`.
- **Rationale**: el mensaje "no hay historias" y los totales en cero (FR-005) salen de `total.historias == 0`; los grupos siempre presentes (Clarification 6) sin que el frontend tenga que completarlos.
- **Alternatives considered**: objeto con claves por estado (`nueva`, `en_progreso`, `completada`; el orden y los nombres con espacio son frágiles); omitir `total` (el frontend sumaría).

## R7. Frontend

- **Decision**: la feature `estado` ya está registrada como `status/*` con carga lazy y la entrada "Estado" de navegación: solo se reemplaza `index.tsx`. Valida el id (`/^[1-9]\d*$/`) antes de pedir, como la pantalla de votación, y distingue carga (`LoadingIndicator`), vacío (mensaje y totales en cero), error (`ErrorMessage`, sin datos parciales) y datos. Estilos en `estado.css`.
- **Rationale**: patrón de `votacion` y `dashboard`; sin librerías nuevas ni cambios en `rutas.tsx`.
- **Alternatives considered**: ruta nueva (innecesaria: la pestaña existe).
