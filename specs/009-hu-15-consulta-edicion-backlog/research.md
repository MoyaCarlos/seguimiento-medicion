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

- **Decision**: en la consulta SQL del adaptador, como fija `MODELO-DE-DATOS.md` §4.4:
  `WHERE proyecto_id = ? AND sprint_id IS NULL ORDER BY CASE prioridad WHEN 'M' THEN 0
  WHEN 'S' THEN 1 WHEN 'C' THEN 2 WHEN 'W' THEN 3 END, id`. El service no reordena.
- **Rationale**: decisión del equipo ya documentada en el modelo de datos. El orden alfabético
  de `prioridad` no sirve (`C` < `M` < `S` < `W`), por eso el `CASE` explícito; el "orden de
  creación" es `id`. El ciclo TDD de la regla se hace con el test de integración del
  repositorio contra SQLite `:memory:` (RED/GREEN) y el escenario BDD.
- **Alternatives considered**: ordenar en el dominio; descartada para seguir el modelo de datos
  acordado por el equipo.

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
  `domain.ErrHistoriaNoEncontrada`), después valida los campos (`ConDatosEditados`) y por
  último aplica la regla D2b (R10). En el handler, la comprobación de
  claves obligatorias del cuerpo y el parseo de `estimacion_sp` se hacen **después** de
  verificar la existencia, igual que `ProjectHandler.Editar`. Para la consulta, el service
  llama a `ProjectRepository.ObtenerPorID` (→ `ErrProyectoNoEncontrado`) antes de listar.
- **Rationale**: FR-012 y edge case "la existencia domina", mismo criterio que HU-04/HU-05.
  `escribirError` agrega `ErrHistoriaNoEncontrada` al caso 404.

## R7. Puertos de persistencia (ISP)

- **Decision**: dos interfaces nuevas en `internal/repository/backlog_repository.go`, sin
  tocar `BacklogRepository` (crear):
  - `ProductBacklogRepository { ListarSinSprint(ctx, proyectoID) ([]domain.BacklogItem, error) }` (ya ordenado, R2)
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

## R10. Story Points congelados en Sprints Activo o Finalizado (D2b)

- **Decision**: regla de **servicio** en `EditarHistoria`, como fija `MODELO-DE-DATOS.md` §5
  ("Servicio (HU-15 y HU-14)"). Si la historia tiene `SprintID`, el service carga el Sprint con
  `SprintRepository.ObtenerPorID`; si su estado es `Activo` o `Finalizado`, la estimación actual
  no es `nil` y la nueva es distinta (otro valor o `nil`), devuelve
  `domain.ErrEstimacionCongelada`, un `ValidationError{Campo: "estimacion_sp", Mensaje: "los
  Story Points ya cargados no se pueden cambiar en un Sprint Activo o Finalizado"}` (400), y no
  llama a `Actualizar`. Cargarla desde `nil` se permite. Con Sprint `Pendiente` o sin Sprint no
  hay restricción.
- **Rationale**: decisión D2b del equipo (texto: "No se modifican los Story Points ya cargados
  ... Cargarlos por primera vez sí se permite"). Reenviar el mismo valor no modifica los Story
  Points, por eso no se rechaza: el formulario siempre manda todos los campos (Assumptions).
  `EditarHistoria` depende de la interfaz `repository.SprintRepository` ya existente (DIP).
- **Alternatives considered**: ponerla en el dominio (`ConDatosEditados`); el modelo de datos la
  asigna al servicio porque necesita el estado del Sprint, que no es parte de la historia.
