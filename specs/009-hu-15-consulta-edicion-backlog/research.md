# Research: Consulta y Edición del Product Backlog (HU-15)

**Fecha**: 2026-10-25 | **Spec**: [spec.md](./spec.md)

No quedaron marcas NEEDS CLARIFICATION en el Technical Context: las tres dudas de
negocio se resolvieron en `/speckit.clarify` (ver Clarifications de la spec). Este
documento registra las decisiones técnicas tomadas para el diseño.

## R1. Rutas de la API

- **Decision**: `GET /projects/{id}/backlog` para consultar y `PUT /backlog/{id}` para editar.
- **Rationale**: el backlog es un recurso del proyecto (igual que `GET /projects/{id}/members`
  y `/status`). La edición sigue el recurso ya existente `POST /backlog`, donde `{id}` es el
  identificador de la historia. `PUT` porque la edición envía todos los campos editables
  juntos (Assumptions: "no hay edición parcial"), igual que `PUT /projects/{id}`.
- **Alternatives considered**: `GET /backlog?proyecto_id=` (rompe el estilo de recursos
  anidados ya usado); `PATCH /backlog/{id}` (implica edición parcial, descartada en la spec).

## R2. Dónde vive la regla de orden MoSCoW

- **Decision**: en el dominio. `Prioridad.Orden()` devuelve 0..3 (M, S, C, W) y una función
  `domain.OrdenarBacklog([]BacklogItem)` ordena estable por (orden de prioridad, id). El
  repositorio devuelve las historias sin Sprint ordenadas por `id`; el service aplica el orden.
- **Rationale**: es una regla de negocio (FR-003) → Principio I (TDD sobre el dominio) y DRY:
  el fake en memoria y SQLite no tienen que replicar el `ORDER BY CASE`. El volumen es chico
  (sin paginación), así que ordenar en memoria no tiene costo.
- **Alternatives considered**: `ORDER BY CASE prioridad WHEN 'M' THEN 0 ... END, id` en SQL
  (la regla quedaría en infraestructura y sin test unitario de dominio).

## R3. Edición en el dominio

- **Decision**: método `func (b BacklogItem) ConDatosEditados(titulo, descripcion string,
  criterios *string, estimacionSP *int) (BacklogItem, error)`, análogo a
  `Project.ConDatosEditados`. Devuelve una copia; conserva ID, proyecto, estado, prioridad,
  valor de negocio y Sprint (FR-011). Las validaciones de título y descripción se extraen a
  helpers compartidos con `NewBacklogItem` (REFACTOR, DRY con HU-01).
- **Rationale**: valor inmutable, una sola validación, mismo patrón que proyectos. Como
  devuelve error antes de tocar la copia, una edición inválida no guarda nada (FR-010).
- **Alternatives considered**: setters por campo (permitirían ediciones parciales y estados
  intermedios inválidos).

## R4. Escala de Story Points

- **Decision**: `domain.EsEstimacionValida(v int) bool` reutiliza la misma escala Fibonacci
  del valor de negocio (1, 2, 3, 5, 8, 13, 21). Se renombra el mapa interno a
  `escalaFibonacci` y ambas funciones lo consultan.
- **Rationale**: Clarification 1: misma baraja para Valor de Negocio, HU-02 y HU-14. Una sola
  fuente (DRY) y una función con nombre propio para que el mensaje de error hable de Story Points.
- **Alternatives considered**: llamar directo a `EsValorNegocioValido` (nombre engañoso).

## R5. Story Points no enteros en el JSON

- **Decision**: el DTO decodifica `estimacion_sp` como `*json.Number`. Si no es un entero,
  el handler responde 400 con `campo: "estimacion_sp"` y el mismo mensaje que el dominio
  (`domain.ErrEstimacionInvalida`, un `ValidationError` exportado). `null` = sin estimación.
- **Rationale**: con `*int`, un `2.5` haría fallar el decode y saldría "JSON inválido" sin
  campo, cuando la spec pide una advertencia sobre los Story Points (edge case "números no
  enteros"). El mensaje sigue viviendo en el dominio.
- **Alternatives considered**: `*float64` en el dominio (filtra un detalle de transporte al
  modelo).

## R6. Existencia antes que validación

- **Decision**: `EditarHistoria.Ejecutar` primero llama a `ObtenerPorID` (que devuelve
  `domain.ErrHistoriaNoEncontrada`) y recién después valida. En el handler, la comprobación de
  claves obligatorias del cuerpo y el parseo de `estimacion_sp` se hacen **después** de
  verificar la existencia, igual que `ProjectHandler.Editar`. Para la consulta, el service
  llama a `ProjectRepository.ObtenerPorID` (→ `ErrProyectoNoEncontrado`) antes de listar.
- **Rationale**: FR-012 y edge case "la existencia domina", mismo criterio que HU-04/HU-05.
  `escribirError` agrega `ErrHistoriaNoEncontrada` al caso 404.

## R7. Puertos de persistencia (ISP)

- **Decision**: dos interfaces nuevas en `internal/repository/backlog_repository.go`, sin
  tocar `BacklogRepository` (crear):
  - `ProductBacklogRepository { ListarSinSprint(ctx, proyectoID) ([]domain.BacklogItem, error) }`
  - `HistoriaRepository { ObtenerPorID(ctx, id) (domain.BacklogItem, error); Actualizar(ctx, item) error }`
  `SQLiteBacklogRepository` implementa las tres.
- **Rationale**: mismo criterio que `HistoriasDeSprintRepository` (ISP): los stubs existentes
  de crear historia no tienen que implementar métodos que no usan.

## R8. Criterios de aceptación en la persistencia

- **Decision**: la columna `historias.criterios_aceptacion TEXT` ya existe en
  `001_esquema_base.sql` (#107). No hace falta migración. `Guardar`, los `SELECT` y
  `Actualizar` pasan a incluirla; NULL ↔ `nil` (`*string`). Las historias previas quedan con
  NULL (FR-014).
- **Rationale**: Assumption "la implementación se hace después de #107", ya cumplida.

## R9. Concurrencia

- **Decision**: `UPDATE ... WHERE id = ?` sin versión: gana el último guardado válido.
- **Rationale**: edge case explícito de la spec (mismo criterio que HU-04).
