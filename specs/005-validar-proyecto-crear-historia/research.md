# Research: Validar proyecto existente al crear historia de backlog (Issue #19)

Sin `NEEDS CLARIFICATION` pendientes: la spec cerró el único punto discutible (orden de
validación, sesión 2026-10-05). Este documento registra las decisiones técnicas y el
razonamiento, incluidos los impactos que aparecieron al inspeccionar el código real.

## D1. Reutilizar el patrón de `CrearSprint`, sin duplicar la validación

**Decisión**: agregar a `CrearHistoriaBacklog` una dependencia `repository.ProjectRepository`
y, dentro de `Ejecutar`, llamar `proyectos.ObtenerPorID(ctx, input.ProyectoID)` después de
`domain.NewBacklogItem` y antes de `repo.Guardar`. Si `ObtenerPorID` devuelve error, se
propaga con `return domain.BacklogItem{}, err`.

**Rationale**:
- `CrearSprint` (`internal/service/crear_sprint.go:25`) ya hace exactamente esto; copiar el
  patrón mantiene KISS/DRY y deja el código consistente (mismo estilo, mismo punto del flujo).
- La verificación de existencia **no se reimplementa**: `ObtenerPorID` ya devuelve
  `domain.ErrProyectoNoEncontrado` en ambas implementaciones
  (`sqlite_project.go:184` por `sql.ErrNoRows`, `project_repository_en_memoria.go:73`).
- Se mantiene el orden invariantes → existencia para no alterar el comportamiento observable
  existente (`ValidationError` por título/prioridad tiene prioridad, FR-005).

**Alternatives considered**:
- Validar en el handler antes de llamar al service — descartada: duplica la regla en la capa
  HTTP y rompe "la lógica vive en el caso de uso"; `CrearSprint` no lo hace así.
- Agregar un método `Existe(id)` a `ProjectRepository` — descartada: YAGNI, no aporta nada
  sobre `ObtenerPorID`, que es el método ya usado por el patrón de referencia.
- Envolver el error con `fmt.Errorf(...)` — descartada: rompería `errors.Is(err,
  domain.ErrProyectoNoEncontrado)` si no se usa `%w`, y complica la traducción en `escribirError`.

## D2. Propagar el error tal cual (FR-006)

**Decisión**: devolver el error de `ObtenerPorID` sin traducirlo.

**Rationale**: cubre con una sola línea tanto `ErrProyectoNoEncontrado` como cualquier error
de persistencia (FR-006). Es el mismo criterio que `CrearSprint`, y es lo que `escribirError`
espera para decidir el status (404 para "no encontrado", 500 para el resto).

**Alternatives considered**:
- Traducir a `ErrProyectoNoEncontrado` solo cuando `errors.Is(...)` — innecesario: el
  repositorio ya devuelve ese error exacto; agregar la traducción sería ruido sin valor.

## D3. El cambio de firma del constructor impacta a **todos** los llamadores

**Decisión**: actualizar los cuatro (4) puntos de construcción de `NewCrearHistoriaBacklog`,
no solo `main.go`:

| Ubicación | Por qué |
|---|---|
| `internal/service/crear_historia_test.go` | Es la suite del propio service (pedido explícito). |
| `cmd/api/main.go` | Wiring de producción (pedido explícito). |
| `internal/http/backlog_handler_test.go:27` | `nuevoHandlerDePrueba()` instancia el service; sin actualizar no compila el paquete `http`. **No estaba en el pedido.** |
| `features/steps_creacion.go:44` | `iniciarBD()` instancia el service; sin actualizar no compila el paquete `features` (rompe toda la suite BDD). **No estaba en el pedido.** |

**Rationale**: al agregar un parámetro a un constructor en Go, cualquier llamador que no se
actualice rompe la compilación. El pedido listó dos, pero hay dos más en el repo. Detectarlos
ahora evita que `/speckit.implement` falle a mitad de camino o, peor, que el equipo no sepa
por qué "no compila" algo fuera de los archivos que tocó.

**Alternatives considered**:
- Mantener una sobrecarga/variante `NewCrearHistoriaBacklogSinProyecto` — descartada: Godot,
  no hay overloading en Go y una segunda función solo para tests contradice KISS/DRY.
- Hacer `ProyectoRepository` opcional (nil-check) — descartada: esconde el requisito y
  permitiría volver a crear huérfanos por configuración; el proyecto debe existir siempre.

## D4. El mapeo a 404 **ya existe**: `escribirError` lo cubre (FR-007/US2)

**Decisión**: no tocar `internal/http/backlog_handler.go` ni `escribirError`; no agregar tests
de handler para el 404 en esta historia.

**Rationale**: `escribirError` (`internal/http/sprint_handler.go:118`) ya incluye
`errors.Is(err, domain.ErrProyectoNoEncontrado)` en la rama que responde `StatusNotFound`.
Una vez que el service devuelva ese error, `POST /backlog` responde 404 automáticamente. El
"500 genérico actual" que menciona el issue quedó obsoleto cuando HU-05 introdujo ese mapper
compartido. Confirmado por inspección: `backlog_handler.go:36` ya usa `escribirError`.

**Alternatives considered**:
- Agregar un `case` específico en el handler — descartada: duplicaría lo que ya hace el mapper
  compartido (viola DRY) y estaba fuera de alcance por decisión del usuario.
- Crear una tarea de "mapeo de errores" para este caso — innecesaria: no hay trabajo pendiente
  para este error concreto. Se deja constancia para que el equipo no abra una tarea vacía.

## D5. BDD: escenario nuevo en el feature de creación de historias

**Decisión**: agregar al final de `features/crear_historia_backlog.feature` un escenario
"Rechazo por proyecto inexistente" (Dado un proyecto que no existe / Cuando intento crear la
historia / Entonces responde "proyecto no encontrado" y no registra la historia), con sus
steps en `features/steps_creacion.go`.

**Rationale**: Principio III de la constitución exige escenario Godog para todo comportamiento
observable. El `Antecedentes` actual siembra el proyecto 1; para el caso negativo se necesita
un identificador inexistente (p. ej. 999) sin sembrarlo, apoyándose en los steps existentes
`noRegistraHistoria` y en `existeProyecto` para el camino feliz.

**Alternatives considered**:
- Nuevo archivo `.feature` solo para esta corrección — descartada: AGENTS.md pide un
  `.feature` por historia; esta es una corrección del comportamiento de HU-01, va en el mismo.
- No agregar BDD — descartada: incumpliría la constitución y la Definition of Done.

## D6. Sin `FOREIGN KEY` en el esquema

**Decisión**: no modificar `internal/repository/migrate.go` ni agregar restricciones.

**Rationale**: fuera de alcance explícito (FR-008), acordado con el Scrum Master como paso
posterior. La validación en aplicación resuelve el caso del issue; la FK sería defensa en
profundidad para inserciones por fuera de este caso de uso.

**Alternatives considered**:
- Agregar la FK ahora — descartada por decisión del usuario y porque cambiar el esquema
  exige migración de datos existentes (puede haber huérfanas históricas), lo que es un paso
  propio con su spec.

## /speckit.analyze (remediación, 2026-10-05)

**Resultado**: **0 CRITICAL**. El análisis corrió sobre `spec.md`, `plan.md` y `tasks.md`
(con apoyo en `research.md`, `data-model.md`, `contracts/openapi.yaml`, `quickstart.md` y la
constitución). Sin violaciones de la constitución y sin requisitos en alcance sin cobertura.

Hallazgos (no bloqueantes) y resolución aplicada:

| ID | Categoría | Severidad | Ubicación | Resumen | Resolución |
|---|---|---|---|---|---|
| I1 | Inconsistencia | HIGH | `tasks.md` T003 | El test RED referenciaba `errPersistencia`, declarado en `internal/http/project_handler_test.go` (paquete `http`), inaccesible desde `internal/service`; el test no habría compilado. | Corregido: T003 ahora usa un error local del paquete service (`errFalloPersistencia := errors.New(...)` en `crear_historia_test.go`) y advierte explícitamente no reusar el de `http`. |
| I2 | Inconsistencia | MEDIUM | `spec.md` FR-007 | FR-007 era un MUST a la vez marcado "fuera de alcance de implementación" (tensión interna). | Corregido: FR-007 se reformuló como criterio de aceptación relacionado ya cubierto por `escribirError`; se alinearon US2 y el último supuesto. |
| U1 | Subespecificación | LOW | `tasks.md` T011 | El paso BDD `Dado` se describía como "fija proyectoID y ejecuta ingresoValido", solapándose con el paso `Cuando`. | Corregido: el `Dado` solo fija `s.proyectoID`; la ejecución la hace `ingresoValido` en el `Cuando`. |
| C1 | Cobertura | MEDIUM | `spec.md` FR-007 | FR-007 sin tareas. | Aceptado como cobertura diferida: ya satisfecho por `escribirError`; verificación manual en T018 del quickstart. Sin tareas de handler (fuera de alcance). |
| C2 | Cobertura | LOW | `spec.md` FR-008 | FR-008 sin tareas. | Correcto: fuera de alcance explícito (sin `FOREIGN KEY`). |
| D1 | Duplicación | LOW | `spec.md` SC-001/002/003 vs FR-002/003 | Los criterios de éxito reexpresan requisitos. | Sin acción (derivación normal). |

Sin hallazgos CRITICAL: FR-001..FR-006 quedan cubiertos por T003/T005 (RED/GREEN) y SC-004 por
la verificación manual T018. La remediación de I1/I2/U1 se aplicó en `spec.md` y `tasks.md`
(ver Phase 1 de `tasks.md`).
