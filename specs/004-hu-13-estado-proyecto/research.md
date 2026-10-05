# Research: Fechas y Estado del Proyecto (HU-13)

Sin `NEEDS CLARIFICATION` pendientes: la única ambigüedad crítica (desempate fechas vs
Sprints) se resolvió en `spec.md` (sesión 2026-10-04). Este documento registra las decisiones
técnicas y el razonamiento, no investigación externa nueva.

## Decisión: El estado es derivado, nunca persistido

**Decisión**: el estado del proyecto se calcula en cada consulta con una función pura del
dominio; no se agrega columna `estado` a `projects` ni tabla nueva.

**Rationale**:
- El estado depende del **instante de la consulta** (la fecha actual). Persistirlo quedaría
  obsoleto apenas cambia el día, obligando a un recálculo/actualización periódica (lo que la
  spec descarta explícitamente: FR-011, SC-003).
- Cumple YAGNI: no se agrega infraestructura (columna, job, trigger) para un dato que se
  reconstruye en O(cantidad de Sprints) a partir de datos que ya existen.
- Cumple la restricción del usuario: **no se agregan tablas ni campos nuevos**.

**Alternatives considered**:
- Columna `estado` en `projects` actualizada por triggers/jobs — descartada: dato derivado,
  no testeable con TDD puro de dominio, y se desincroniza con el paso del tiempo.
- Campo `Estado` en la struct `domain.Project` — descartada: mezcla un valor calculado con el
  estado persistido de la entidad; rompería el principio de que la entidad refleja la fila de
  la base. El comentario existente en `project.go` menciona "fechas/estado (HU-13)"; se
  interpreta como fechas ya implementadas + estado **calculado**, no como campo nuevo.

## Decisión: Regla de derivación y precedencia

**Decisión** (función pura en `internal/domain`):

1. Si hay al menos un Sprint `Activo` ⇒ **"En curso"**.
2. Si no hay `Activo` y existe al menos un Sprint y **todos** están `Finalizado` ⇒ **"Finalizado"**.
3. Si no hay `Activo` y hay al menos un `Finalizado` junto a al menos un `Pendiente` ⇒ **"En curso"**.
4. En el resto, decide la fecha:
   - sin `FechaInicio`, o `ahora` anterior a `FechaInicio` ⇒ **"Planificado"**;
   - `ahora` dentro de `[FechaInicio, FechaFin]` (bordes inclusivos) ⇒ **"En curso"**;
   - `ahora` posterior a `FechaFin` ⇒ **"Finalizado"**.

**Rationale**: los Sprints reflejan ejecución real y pesan más que la fecha estimada; el
conflicto "Sprint Activo con fecha de fin vencida" se resuelve como "En curso". Las fechas
solo desempatan cuando no hay Sprints iniciados. Coincide con la aceptación y con los casos
borde de `spec.md`.

**Alternatives considered**:
- Fechas dominantes — descartada por el usuario: mostraría "Finalizado" con un Sprint aún activo.
- Regla híbrida con "finalizado por fecha" antes que por Sprints — descartada: agrega una
  rama sin valor observable distinto al caso de uso.

## Decisión: Hueco de especificación `[Finalizado, Pendiente]` (2026-10-05)

**Decisión**: sin Sprint `Activo`, con al menos un `Finalizado` y al menos un `Pendiente`
⇒ **"En curso"** (FR-016). Un Sprint cerrado es evidencia de ejecución; el proyecto no
terminó porque queda un Sprint pendiente. Las fechas solo deciden cuando no hay Sprints
iniciados.

**Rationale**: el estado se había definido como "todos Finalizado ⇒ Finalizado" y "fechas
en el resto", lo que dejaba a `[Finalizado, Pendiente]` cayendo a las fechas (⇒
"Planificado" sin `FechaInicio`), pese a haber ejecución registrada. La regla cierra el
hueco sin alterar D1 (un único Sprint `Finalizado` y sin `Pendiente` sigue dando
"Finalizado" por FR-004, decisión que queda en manos del equipo).

## Decisión: El instante actual se inyecta como parámetro

**Decisión**: `CalcularEstadoProyecto(p domain.Project, sprints []domain.Sprint, ahora time.Time)`
recibe `ahora`; el handler HTTP pasa `time.Now()` y los tests pasan una fecha fija.

**Rationale**: mantiene la función de dominio pura y determinista (sin llamar a `time.Now()`
adentro), lo que permite TDD sin flakiness. No se introduce una interfaz `Clock` — sería
abstracción para una necesidad hipotética (YAGNI); pasar el `time.Time` como argumento es la
solución más simple que cumple el requisito.

**Alternatives considered**:
- Inyectar `func() time.Time` (clock) en el service — descartada por KISS: agrega un campo y
  un constructor extra sin ganar nada sobre pasar el `time.Time` en `Ejecutar`.
- Llamar `time.Now()` dentro del dominio — descartada: acopla el dominio al reloj y hace el
  test dependiente del reloj real.

## Decisión: Caso de uso y endpoint

**Decisión**: un único service `ObtenerEstadoProyecto` (puertos `ProjectRepository` +
`SprintRepository`) y un endpoint de lectura `GET /projects/{id}/status` que devuelve
`{ "estado": "..." }`.

**Rationale**:
- Un caso de uso por servicio (SRP, Principio V). El service es de solo lectura: solo llama
  `ObtenerPorID` y `ListarPorProyecto`, nunca `Actualizar`/`Guardar`.
- Endpoint dedicado en lugar de sumar `estado` a `GET /projects/{id}`: no altera el contrato
  ni el service de HU-04 (`ObtenerProyecto`), respeta SRP y mantiene la lectura de estado
  desacoplada.
- 404 usando `domain.ErrProyectoNoEncontrado`, ya traducido por `escribirError` en
  `internal/http/sprint_handler.go` (mismo criterio que HU-04/HU-05).

**Alternatives considered**:
- Extender `GET /projects/{id}` con `estado` — descartada: obliga a `ObtenerProyecto` (HU-04)
  a depender de `SprintRepository`, ampliando su responsabilidad y su superficie de test.
- Nombre de ruta `/projects/{id}/estado` — descartada por consistencia con los sustantivos en
  inglés ya usados en las rutas (`/members`); el campo de respuesta igualmente es `estado`.
