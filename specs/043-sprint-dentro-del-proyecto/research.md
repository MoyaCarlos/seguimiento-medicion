# Research: Sprints dentro del plazo del proyecto (#168)

No quedaron `NEEDS CLARIFICATION` en el contexto técnico: las decisiones de producto se tomaron con el equipo (spec, Clarifications). Esta sección registra las decisiones de diseño.

## 1. Dónde se valida: al iniciar, no al crear

- **Decision**: la regla se evalúa en el caso de uso `IniciarSprint` (`POST /sprints/{id}/iniciar`).
- **Rationale**: `CrearSprint` solo crea un Sprint Pendiente sin Goal ni fechas (`domain.NewSprint`); las fechas entran únicamente en `Sprint.Iniciar`. El issue mencionaba `CrearSprint` porque lo redactó antes de leer el flujo.
- **Alternatives considered**: validar en `CrearSprint` — no hay fechas que comparar.

## 2. Dónde vive la regla: un método de `Project` en el dominio

- **Decision**: `func (p Project) AdmiteSprint(inicio, fin time.Time) error` en `internal/domain/project.go`. Devuelve `ValidationError` en `fecha_inicio` o `fecha_fin`, o `nil`. El servicio la llama después de `sprint.Iniciar(...)` y antes de `Actualizar`.
- **Rationale**: el plazo es un dato del proyecto, así que la regla "¿este rango cabe en mi plazo?" es del proyecto (SRP, DRY). Queda testeable sin repositorios. Llamarla después de `sprint.Iniciar` mantiene el orden de errores existente (FR-006) y, como `sprint` es una copia en memoria, un rechazo no persiste nada (FR-007).
- **Alternatives considered**:
  - Pasar el proyecto a `Sprint.Iniciar` — acopla `Sprint` a `Project` y obliga a tocar todos los llamadores de `Iniciar` (incluido `internal/demo`), sin beneficio.
  - Validar en el handler HTTP — viola FR-005 (la regla vale para cualquier cliente y para los pasos BDD, que llaman al servicio).

## 3. Comparación de fechas

- **Decision**: comparar con `Before`/`After` sobre los `time.Time` tal como llegan: inicio del Sprint `Before(*p.FechaInicio)` → error; fin del Sprint `After(*p.FechaFin)` → error. Igualdad = aceptado (bordes incluidos).
- **Rationale**: desde la migración 001 las fechas del proyecto también se guardan como `AAAA-MM-DD` (`MODELO-DE-DATOS.md` D3 resuelta; `sqlite_project.go` usa `time.DateOnly`), y la API parsea ambas con `2006-01-02`. Todas quedan a las 00:00 UTC, así que la comparación de instantes equivale a la de días. La preocupación del issue por RFC 3339 ya no aplica.
- **Alternatives considered**: normalizar a día con `time.Date(...)` antes de comparar — defensivo pero innecesario (YAGNI); si algún día entrara una hora distinta de 00:00, el test de bordes lo detectaría.

## 4. Orden de las validaciones del inicio

- **Decision**: igual al de hoy, con el plazo al final: Sprint inexistente (404) → otro Sprint Activo (409) → validaciones propias de `Sprint.Iniciar` (estado Pendiente, Goal, inicio obligatorio, fin posterior al inicio; 400) → **plazo del proyecto (400)**. Dentro del plazo, primero `fecha_inicio` y después `fecha_fin`.
- **Rationale**: FR-006 y el caso límite "las dos fechas afuera → `fecha_inicio`". Si se mergea antes HU-21 (#301), su 403 queda antes del 409 como define esa spec; esta regla sigue al final.

## 5. Dependencia nueva del servicio y convivencia con HU-21 (#301)

- **Decision**: `NewIniciarSprint(sprints repository.SprintRepository, proyectos repository.ProjectRepository)`. El servicio obtiene el proyecto con `proyectos.ObtenerPorID(sprint.ProyectoID)`.
- **Rationale**: es el puerto existente que devuelve `domain.Project` con sus fechas; no hace falta un puerto nuevo (YAGNI). El PR #301 (abierto) cambia el mismo constructor a `NewIniciarSprint(sprints, integrantes repository.IntegrantesRepository)` y, en su cableado (`internal/app/app.go`), ya le pasa el `ProjectRepository`. **Resolución del conflicto, para la rama que se mergee segunda**: el campo queda como `proyectos repository.ProjectRepository` (que también satisface `IntegrantesRepository`), y se usa tanto para `exigirScrumMaster` como para `ObtenerPorID`. La firma pública queda `NewIniciarSprint(sprints, proyectos)` en ambos casos, así que los llamadores no cambian dos veces.
- **Alternatives considered**: un puerto angosto `ProyectoLector` con solo `ObtenerPorID` — más limpio en teoría, pero agrega una interfaz que nadie más usa y complica la fusión con #301.

## 6. Mensajes de error

- **Decision**:
  - `fecha_inicio`: `el Sprint no puede empezar antes del inicio del proyecto (AAAA-MM-DD)`
  - `fecha_fin`: `el Sprint no puede terminar después del fin del proyecto (AAAA-MM-DD)`
- **Rationale**: "mensaje claro" del issue: incluir la fecha límite le dice a quien usa la pantalla qué valor elegir, sin otra consulta.

## 7. Impacto sobre pruebas y datos existentes

- **Escenarios BDD**: los proyectos de `/features` se crean sin fechas (`domain.NewProject(..., nil, nil)`), y hay que pasar el `ProjectRepository` a los `NewIniciarSprint` de los pasos (5 archivos). **Corrección (implementación, 2026-11-11):** los pasos de HU-13 (`features/steps_estado_proyecto.go`) sí cargan fechas al proyecto con un `UPDATE` y luego iniciaban un Sprint de hoy a hoy+14 fuera de ese plazo; 3 escenarios de `estado_proyecto.feature` quedaron en rojo y se ajustaron para iniciar el Sprint dentro del plazo (ver tasks.md, notas de la fase 3).
- **Tests unitarios del servicio y del handler**: hoy siembran Sprints de proyectos que no existen en el repositorio de proyectos (`proyectoID` 1 sin alta). Con la nueva lectura del proyecto darían 404; hay que sembrar el proyecto (sin fechas) en esos tests. Es un ajuste de arreglo (Arrange), no de comportamiento.
- **Datos de prueba v2 (`internal/demo`)**: proyecto de −30 a +45 días y Sprints entre −28 y +6: quedan dentro del plazo. Además la carga usa `Sprint.Iniciar` directo, no el servicio.

## 8. Frontend

- **Decision**: sin cambios de frontend en esta historia.
- **Rationale**: `web/src/features/sprints/FormularioInicio.tsx` ya clasifica el 400 por `campo` (`sprint_goal`, `fecha_inicio`, `fecha_fin`) y muestra el `mensaje` debajo del campo correspondiente, así que el nuevo error se ve en el lugar correcto sin tocar la pantalla. El "issue hermano" que menciona el #168 no existe en el repositorio (búsqueda del 2026-11-11); lo único que podría aportar es limitar los selectores de fecha al plazo del proyecto (`min`/`max`), una mejora de usabilidad que no hace falta para cumplir esta regla.
