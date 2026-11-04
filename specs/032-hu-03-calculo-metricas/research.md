# Research: Cálculo de Métricas del Proyecto (HU-03)

Decisiones técnicas de la Fase 0. No quedan `NEEDS CLARIFICATION`. Las decisiones de producto están
en las Clarifications de la [spec](./spec.md); acá solo se decide cómo implementarlas.

El código de referencia es `origin/main` del 2026-11-03.

## R1. Dónde vive el cálculo

- **Decision**: una función pura en `internal/domain/metricas.go` recibe los datos ya sumados de cada Sprint finalizado y las horas por Story Point, y devuelve las métricas por Sprint y del proyecto. El service `CalcularMetricas` solo orquesta.
- **Rationale**: el cálculo se prueba con TDD sin base de datos, y HU-09, HU-10 y HU-12 consumen un único resultado (Principio V, DRY).
- **Alternatives considered**: calcular todo en SQL (la desviación y los redondeos quedarían fuera del dominio y sin TDD directo); guardar las métricas al cerrar (contradice `MODELO-DE-DATOS.md` §6 y la decisión de recalcular).

## R2. Lectura de los datos

- **Decision**: un puerto nuevo, `DatosMetricasRepository`, con un solo método que devuelve una fila por Sprint finalizado. El adaptador SQLite usa una sola consulta con subconsultas correlacionadas, con las fuentes de §6.
- **Rationale**: los repositorios actuales no alcanzan: `ListarPorSprint` no incluye las arrastradas, `RegistroHorasRepository` lista por historia y `DefectoRepository` trae todos los defectos con su detalle. Una consulta evita N+1 y lee un estado coherente. El puerto separado respeta ISP.
- **Alternatives considered**: componer los repositorios existentes en el service (varias consultas por Sprint y las horas no se pueden agrupar por Sprint sin un método nuevo); reutilizar la `HistoriasComprometidas` de HU-12 (todavía no existe, y HU-03 solo necesita sumas).

## R3. Número del Sprint

- **Decision**: `numero` se calcula con la misma subcuenta que `sqlite_sprint.go`: la posición del Sprint entre **todos** los de su proyecto, por `id`. El filtro por `estado = 'Finalizado'` no afecta esa subcuenta.
- **Rationale**: `sprints` no tiene columna `numero` (#129, spec 024). Con `ROW_NUMBER()` sobre las filas filtradas, o contando solo los finalizados, un Sprint podría informar otro número que en el resto de la aplicación.
- **Alternatives considered**: numerar en Go a partir de la lista filtrada (mismo error).

## R4. Aritmética y redondeo

- **Decision**: las sumas se hacen en enteros (Story Points y minutos). La desviación en horas es `minutos/60 − SP × horas_por_SP`, y el porcentaje es esa desviación sobre las estimadas, ambos sin redondear. Solo el valor final se redondea a dos decimales con `math.Round(x*100)/100`, y un resultado `-0` se normaliza a `0`.
- **Rationale**: FR-011 y D8. El proyecto se calcula desde los totales enteros y no desde los valores redondeados de cada Sprint, así que su desviación puede diferir en 0,01 de la suma de los Sprints; es el valor correcto. La normalización hace falta porque `encoding/json` escribe `-0` para el cero negativo.
- **Alternatives considered**: redondear antes de restar (acumula error); usar un tipo decimal (dependencia nueva, Principio VI).

## R5. "Sin datos" en la respuesta

- **Decision**: `*float64` en el dominio y `null` en el JSON para `horas_estimadas`, `desviacion_horas` y `desviacion_porcentaje`.
- **Rationale**: distingue "sin datos" de cero (SC-003) y coincide con el contrato provisorio de la spec 033.
- **Alternatives considered**: un texto `"sin datos"` (mezcla tipos; el texto es una decisión de cada pantalla).

## R6. Contrato de la consulta

- **Decision**: `GET /projects/{id}/metricas`. Los nombres de los campos son los de `MetricasSprint` en la spec 033, sin cambios. Se agrega `proyecto` con esa misma forma y `sprints` con `sprint_id`, `numero` y `metricas`.
- **Rationale**: HU-12 declaró esos nombres como provisorios a la espera de este plan (su R2), así que adoptarlos evita retrabajo. La ruta en español sigue a `/defectos` y `/historial`. `/projects` ya está en `PrefijosAPI`.
- **Alternatives considered**: `GET /sprints/{id}/metricas` más otra de proyecto (SC-006 pide una sola consulta); `/metrics` (mezcla idiomas en rutas nuevas).

## R7. Los tres estados del parámetro en la edición

- **Decision**: el handler usa el mapa `claves` que ya tiene para distinguir ausente, `null` y valor. El service recibe `HorasPorSPEnviadas bool` y `HorasPorSP *string`, con el valor JSON sin interpretar. `Project.ConHorasPorStoryPoint(*string)` lo convierte con `strconv.ParseFloat` y exige un número finito mayor que cero. La lista de campos obligatorios del handler no cambia.
- **Rationale**:
  - Con `*float64` en el DTO, `"abc"` rompe el `Unmarshal` y responde "JSON inválido" sin campo, lo que incumple FR-017.
  - Un texto JSON llega con comillas, así que `""`, `"abc"` y `"5"` fallan en `ParseFloat` y se rechazan sobre el campo; `1e999` falla por rango.
  - Validar en el dominio, después de `ConDatosEditados`, cumple que el parámetro se valide al final. Parsear en el handler lo adelantaría a las demás reglas.
  - Es el mismo patrón de HU-07, donde el handler pasa `textoDeHoras` y el dominio valida.
- **Alternatives considered**: un tipo opcional genérico (YAGNI, un solo uso); tratar `""` como vacío (descartado en Clarifications).

## R8. Persistencia del parámetro

- **Decision**: agregar `horas_por_story_point` al `SELECT` de `ObtenerPorID` y `Listar`, a `escanearProyecto` y al `UPDATE` de `Actualizar`, todo en el mismo commit GREEN. El `INSERT` no cambia, porque el parámetro no se pide al crear.
- **Rationale**: si `UPDATE` escribe la columna y la lectura no la trae, cualquier edición borra el valor. Un test de repositorio cubre editar otros datos y conservarlo.
- **Alternatives considered**: un `UPDATE` separado para el parámetro (dos escrituras por edición, sin atomicidad).

## R9. Regla de Story Points con Sprint finalizado

- **Decision**: cambiar `validarEstimacionCongelada`, que ya comparten `EditarHistoria` y `RegistrarEstimacion`:
  1. Sin Sprint, o sin cambio de valor (incluido sin estimar → sin estimar): permitido.
  2. Sprint `Finalizado`: rechazado. Si la historia no tenía Story Points, con el error nuevo `ErrCargaEnSprintFinalizado`; si los tenía, con `ErrEstimacionCongelada`.
  3. Sprint `Activo`: rechazado solo si ya tenía Story Points.
- **Rationale**:
  - El orden actual sale antes de mirar el Sprint cuando la historia no tiene Story Points, por eso hay que reordenarlo.
  - El caso sin cambio es necesario para FR-024: HU-15 envía `estimacion_sp: null` al editar otros datos.
  - El mensaje de `ErrEstimacionCongelada` habla de Story Points "ya cargados" y no describe este caso.
  - Ambos son `ValidationError` con campo `estimacion_sp`, así que `escribirError` ya responde 400. No hay cambios de HTTP.
  - Una historia arrastrada tiene `sprint_id` nulo, entra por el paso 1 y su foto no se toca.
- **Límite conocido**: una carga simultánea con el cierre del Sprint puede pasar la validación antes del cierre. Es el mismo límite de la D2b actual; una aplicación de un equipo chico no justifica un `UPDATE` condicional.
- **Alternatives considered**: una validación aparte para cada historia (duplica la regla); frenarlo en la base (la regla es de servicio, §5, y un `CHECK` no puede mirar otra tabla).

## R10. Dónde van los escenarios BDD

- **Decision**: `features/calculo_metricas.feature` (@HU-03) con `steps_metricas.go`, también para la User Story 5. En `consulta_edicion_backlog.feature` se saca el ejemplo `Finalizado` de la carga inicial. `consenso_estimacion.feature` no se toca.
- **Rationale**: el PR abierto #268 modifica `consenso_estimacion.feature` y `steps_consenso.go`. Ese feature no tiene un escenario de primera carga con Sprint `Finalizado`, así que no contradice la regla nueva.
- **Alternatives considered**: agregar el escenario de votación en `consenso_estimacion.feature` (conflicto con #268).

## R11. Pantalla de edición

- **Decision**: agregar `horas_por_story_point: string` a `DatosProyecto`. `aProyectoEditado` envía `null` si está vacío y `Number(valor)` si no; `aNuevoProyecto` lo ignora. `validarEdicion` rechaza lo que no sea un número mayor que cero.
- **Rationale**: `FORMULARIO_VACIO` define qué campos del backend se muestran junto a un input, así que el error se ubica solo al agregar la clave. Es un `Field` con `type="number"`.
- **Alternatives considered**: un tipo aparte para la edición (más código para un campo).

## R12. Orden de trabajo con `main`

- **Decision**: traer `origin/main` a la rama antes de `/speckit-implement`. Los cambios de documentos entran en el mismo PR.
- **Rationale**: la rama está 102 commits atrás y no tiene HU-14, que esta historia modifica. `/speckit-tasks` puede generarse igual, porque este plan se basa en `origin/main`. El PR abierto #250 modifica `docs/PRODUCT-BACKLOG.md`, así que puede haber un conflicto menor ahí.
- **Alternatives considered**: separar la regla de D2b en una rama `fix/` (no es un defecto sino una decisión de la historia, y sus métricas dependen de ella).
