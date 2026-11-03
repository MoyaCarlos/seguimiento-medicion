# Research: Registro y Seguimiento de Defectos (HU-08)

**Fecha**: 2026-10-30 | **Spec**: [spec.md](./spec.md)

No quedaron marcas NEEDS CLARIFICATION en el Technical Context: las dudas de negocio se
resolvieron en `/speckit.specify` y `/speckit.clarify` (ver Clarifications de la spec). Este
documento registra las decisiones técnicas tomadas para el diseño.

## R1. Rutas de la API

- **Decision**:
  - `POST /projects/{id}/defectos`: registrar un defecto en el proyecto `{id}`.
  - `GET /projects/{id}/defectos`: listar los defectos del proyecto.
  - `POST /defectos/{id}/resolver`: marcar como resuelto el defecto `{id}`.
  
  Se agrega `/defectos` a `PrefijosAPI` (`internal/http/spa.go`) y al proxy de
  `web/vite.config.ts`.
- **Rationale**: el defecto es un recurso hijo del proyecto (`proyecto_id NOT NULL`), igual que
  `GET /projects/{id}/backlog` y `GET /projects/{id}/sprints`. La historia relacionada va en el
  cuerpo y no en la ruta: es un dato del defecto, no su dueño, y puede ser cualquier historia del
  proyecto (Clarification 2). La resolución es una transición de estado sobre un recurso ya
  identificado, con el mismo estilo que `POST /sprints/{id}/iniciar` y `/cerrar`; el proyecto se
  deduce del defecto. El test `TestPrefijosAPI_IncluyenLosQueProxeaVite` obliga a mantener alineados
  los dos listados.
- **Alternatives considered**: `POST /backlog/{id}/defectos` (ata el defecto a la ruta de la
  historia y obliga a listar por proyecto en otra ruta); `PATCH /defectos/{id}` con
  `{"estado": "Resuelto"}` (sugiere que se pueden editar otros campos o volver a `Pendiente`,
  y FR-016 lo prohíbe); `POST /projects/{id}/defectos/{defectoId}/resolver` (evita el prefijo
  nuevo, pero repite el proyecto y suma el caso "defecto de otro proyecto").

## R2. Orden de las validaciones al registrar

- **Decision**: el orden de la spec (Edge Cases, FR-012) se reparte así:
  1. **Handler**: id de ruta inválido → 400; proyecto inexistente → 404, consultando el proyecto
     (`ObtenerProyecto`) **antes** de leer el cuerpo, igual que `POST /backlog/{id}/horas`
     (HU-07). Cuerpo que no es JSON → 400 "JSON inválido".
  2. **Service** `RegistrarDefecto`: Sprint activo del proyecto (`ListarPorProyecto` +
     `domain.SprintActivoDe`) → si no hay, `ErrSinSprintActivoParaDetectar`. Historia relacionada:
     `historia_id` ausente o 0 → `ErrHistoriaRelacionadaObligatoria`; inexistente
     (`ErrHistoriaNoEncontrada` del repositorio) o de otro proyecto → `ErrHistoriaRelacionadaAjena`.
  3. **Dominio** `NewDefecto`: descripción → severidad.
- **Rationale**: la existencia del proyecto domina sobre la validación (HU-04, 05, 06, 07, 15).
  Las reglas que necesitan leer otras entidades (Sprints, historia) están en el service; las que
  dependen solo de los datos del defecto, en el dominio. Una historia inexistente es un error
  **sobre el campo** `historia_id` (400), no un 404: el recurso pedido es el proyecto, y la spec
  pide "un error sobre la historia" (Escenario 1.5), igual que el integrante ajeno de HU-07.
- **Alternatives considered**: validar la historia antes que el Sprint activo (contradice el orden
  fijado en la spec, que pone primero la condición del proyecto y después los datos del cuerpo).

## R3. Orden de las validaciones al resolver

- **Decision**:
  1. **Handler**: id de ruta inválido → 400. No hay cuerpo.
  2. **Service** `ResolverDefecto`: `ObtenerPorID` → `ErrDefectoNoEncontrado` (404). Busca el
     Sprint activo del proyecto del defecto y llama a `defecto.Resolver(activo)`.
  3. **Dominio** `(*Defecto).Resolver(activo *Sprint)`: si ya está `Resuelto` →
     `ErrDefectoYaResuelto`; si `activo` es `nil` → `ErrSinSprintActivoParaResolver`; si no, pasa
     a `Resuelto` con `SprintResolucionID = activo.ID`.
- **Rationale**: "ya resuelto" va antes que "sin Sprint activo" porque es la información más
  útil y no depende del estado del proyecto (orden de la spec). La transición vive en el dominio,
  así que el invariante "Resuelto ⇔ con Sprint de resolución" se prueba con TDD directo y la base
  lo respalda con su `CHECK` (D13: si se viola, 500).

## R4. Carreras: cierre del Sprint y resolución doble

- **Decision**: el adaptador repite la condición en la misma sentencia.
  - Registrar:
    ```sql
    INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, creado_en)
    SELECT ?, ?, ?, ?, 'Pendiente', ?, ?
    WHERE EXISTS (SELECT 1 FROM sprints WHERE id = ? AND estado = 'Activo')
    ```
    Si no inserta ninguna fila → `domain.ErrSinSprintActivoParaDetectar`.
  - Resolver:
    ```sql
    UPDATE defectos SET estado = 'Resuelto', sprint_resolucion_id = ?
    WHERE id = ? AND estado = 'Pendiente'
      AND EXISTS (SELECT 1 FROM sprints WHERE id = ? AND estado = 'Activo')
    ```
    Si no actualiza ninguna fila, relee el defecto: si ya está `Resuelto` →
    `domain.ErrDefectoYaResuelto`; si no → `domain.ErrSinSprintActivoParaResolver`.
- **Rationale**: entre que el service lee el Sprint y que escribe, otro pedido puede cerrar el
  Sprint o resolver el mismo defecto. Una sola sentencia es atómica en SQLite, así que el defecto
  queda asociado a un Sprint que estaba `Activo` o se rechaza con el mismo mensaje que la
  validación previa (Edge Cases). En la resolución doble, solo una sentencia encuentra
  `estado = 'Pendiente'`, y el Sprint de resolución de la primera no se pisa. Es el mismo criterio
  que HU-07 (R5) y que el `UPDATE ... WHERE estado = 'Activo'` del fix #145.
- **Alternatives considered**: transacción con `SELECT` previo (más código para el mismo
  resultado con una sola conexión); confiar solo en el service (deja las carreras abiertas).

## R5. Búsqueda del Sprint activo

- **Decision**: función del dominio `SprintActivoDe(sprints []Sprint) (Sprint, bool)`, que usan
  los dos services con los Sprints de `SprintRepository.ListarPorProyecto`.
- **Rationale**: el proyecto tiene a lo sumo un Sprint `Activo` (índice parcial de HU-05), y
  `ListarPorProyecto` ya existe y está probado. Una sola función evita repetir el bucle en los dos
  services (DRY); `ObtenerEstadoProyecto` (HU-13) e `IniciarSprint` lo hacen hoy a mano y pueden
  pasar a usarla en un REFACTOR si queda limpio, sin obligación en esta historia.
- **Alternatives considered**: método nuevo `ObtenerActivo(proyectoID)` en el repositorio (otra
  consulta SQL para un volumen chico; YAGNI).

## R6. Severidad y estado

- **Decision**: tipos del dominio `Severidad` (`Alta`, `Media`, `Baja`) y `EstadoDefecto`
  (`Pendiente`, `Resuelto`), como `Prioridad` y `EstadoSprint`. La severidad se compara exacta
  (sensible a mayúsculas, sin recortar espacios): `"alta"` o `" Alta"` se rechazan con
  `ValidationError{Campo: "severidad"}`.
- **Rationale**: los valores son los mismos que el `CHECK` de §4.9; aceptar variantes obligaría a
  normalizar y la spec pide "exactamente". El formulario (#68) envía uno de tres valores fijos.
- **Alternatives considered**: normalizar mayúsculas (regla extra sin pedido).

## R7. Descripción

- **Decision**: se recortan los espacios externos (`strings.TrimSpace`); vacía →
  "la descripción es obligatoria"; más de 2000 runas → "la descripción no puede superar los 2000
  caracteres". Se cuentan runas, como `BacklogItem`.
- **Rationale**: mismo límite y mismo criterio que la descripción de una historia (HU-01).

## R8. Listado

- **Decision**: `ListarPorProyecto` devuelve `DefectoDetalle` (el defecto más el título de la
  historia relacionada, con `JOIN historias`), ordenado por `id` (orden de registro,
  Clarification 3). Slice vacío, no `nil`, si no hay defectos. La respuesta expone
  `sprint_resolucion_id: null` para los pendientes.
- **Rationale**: el título evita que el frontend haga un pedido por historia para mostrar el
  listado; `id` es monotónico (`AUTOINCREMENT`), así que coincide con el orden de registro sin
  depender de `creado_en`.
- **Alternatives considered**: ordenar por `creado_en` (dos defectos en el mismo segundo empatan).

## R9. Errores y códigos HTTP

| Caso | Error | HTTP |
|---|---|---|
| Proyecto inexistente (registrar, listar) | `domain.ErrProyectoNoEncontrado` | 404 |
| Defecto inexistente (resolver) | `domain.ErrDefectoNoEncontrado` (nuevo) | 404 |
| Sin Sprint activo al registrar (incluida la carrera de R4) | `ValidationError{Campo: "sprint_deteccion_id"}` | 400 |
| Sin Sprint activo al resolver (incluida la carrera de R4) | `ValidationError{Campo: "sprint_resolucion_id"}` | 400 |
| Historia relacionada ausente, inexistente o de otro proyecto | `ValidationError{Campo: "historia_id"}` | 400 |
| Descripción o severidad inválidas | `ValidationError{Campo: "descripcion" / "severidad"}` | 400 |
| Defecto ya resuelto (incluida la resolución doble de R4) | `ValidationError{Campo: "estado"}` | 400 |
| FK violada (no debería pasar: el service valida antes) | `repository.ErrIntegridadReferencial` | 409 |
| `CHECK` violado (defecto: el dominio dejó pasar algo) | genérico | 500 (D13) |

- **Rationale**: es la tabla de `MODELO-DE-DATOS.md` §5. "Sin Sprint activo" y "ya resuelto"
  son reglas de estado, como "solo se puede cerrar un Sprint activo" (`ErrSprintNoActivo`), que el
  proyecto ya responde con `ValidationError` y 400. Usar como `Campo` el nombre del Sprint que
  falta deja al frontend distinguir los dos mensajes del Escenario 3 sin comparar textos.
  `escribirError` solo suma `ErrDefectoNoEncontrado` al caso 404.

## R10. Migración 005, tests de esquema y BDD

- **Decision**: `internal/repository/migraciones/005_defectos.sql` con la tabla de
  `MODELO-DE-DATOS.md` §4.9 tal cual (FK compuestas, `CHECK` de severidad y estado,
  `CHECK ((estado = 'Resuelto') = (sprint_resolucion_id IS NOT NULL))`). Se agregan a
  `esquema_test.go` los casos mínimos de §9.6 para `defectos`, más historia o Sprint de otro
  proyecto. Los pasos Godog siembran el proyecto, la historia y el Sprint con los services de
  HU-01, HU-04 y HU-05, como `steps_horas.go`, y llaman a los services de esta historia.
- **Rationale**: el número 005 está reservado para HU-08 en §9.3 y el esquema ya está definido,
  así que el documento no cambia. HU-01, HU-04 y HU-05 están en `main`.
- **Alternatives considered**: ninguna; el esquema es la fuente de verdad.
