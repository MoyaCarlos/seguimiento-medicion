# Research: Cerrar y reabrir un proyecto (spec 044, #300)

Fase 0 del plan. No quedan `NEEDS CLARIFICATION`: las tres decisiones de modelo se
confirmaron el 2026-11-11 (ver Clarifications de la spec). Lo que sigue son las
decisiones de diseño tomadas a partir del código de `main` (con HU-21 ya mergeada).

## 1. Dónde vive el chequeo de «proyecto cerrado» (FR-005)

**Decision**: en `internal/service/integrantes.go`, en una función nueva
`exigirProyectoAbierto(ctx, integrantes, proyectoID)`, llamada **solo** desde los
guardias de escritura `exigirIntegrante` y `exigirScrumMaster`. Los helpers quedan así:

| Helper | Qué comprueba | Lo usan |
|---|---|---|
| `rolEnProyecto` | (sin cambios) rol del usuario | todos los de abajo, `ObtenerProyecto` |
| `verificarIntegrante` *(nuevo nombre del chequeo de rol)* | integra el proyecto (403) | `ConsultarVotacion` (lectura) |
| `verificarScrumMaster` *(ídem)* | integra (403) y es Scrum Master (403) | `CerrarProyecto`, `ReabrirProyecto` |
| `exigirIntegrante` | `verificarIntegrante` + `exigirProyectoAbierto` (409) | toda escritura de integrante |
| `exigirScrumMaster` | `verificarScrumMaster` + `exigirProyectoAbierto` (409) | toda escritura de Scrum Master |

Además, dos escrituras que hoy **no** pasan por `exigir*` se migran a `exigirIntegrante`:
`Votar` (usa `rolEnProyecto` directo) e `ImputarHoras` (usa su propio
`nombreDelIntegrante`). Sin ese ajuste votar y registrar horas quedarían abiertos en un
proyecto cerrado.

**Rationale**: HU-21 ya dejó un único punto por el que pasan todas las escrituras sobre un
proyecto, y el orden de errores de su FR-006 (404 → 403 → 400/409) se apoya en él. Poner el
409 ahí lo ubica automáticamente después de los 403 y antes de las validaciones del caso de
uso, como pide la spec. Separar `verificar*` (rol) de `exigir*` (rol + abierto) evita que la
lectura de la votación y la propia reapertura queden bloqueadas.

**Alternatives considered**:
- *Middleware HTTP por ruta*: la capa HTTP no sabe a qué proyecto pertenece una historia,
  un Sprint o un defecto sin consultar la base; duplicaría la resolución que ya hacen los
  servicios y rompería el orden 404 → 403.
- *Chequear en cada servicio*: 17 servicios de escritura; viola FR-005 y DRY, y una escritura nueva
  podría olvidarlo.
- *Restricción en la base* (`WHERE NOT EXISTS proyecto cerrado` en cada `INSERT`/`UPDATE`):
  cierra también la carrera de §6, pero toca todos los adaptadores SQLite y no da el
  mensaje 409; desproporcionado (KISS).

## 2. De dónde sale el dato «cerrado» para el guardia

**Decision**: agregar un método al puerto existente `repository.IntegrantesRepository`:

```go
// ProyectoCerrado informa si el proyecto tiene fecha de cierre.
ProyectoCerrado(ctx context.Context, proyectoID int64) (bool, error)
```

Lo implementan `SQLiteProjectRepository` (`SELECT fecha_cierre IS NOT NULL FROM proyectos
WHERE id = ?`) y `ProjectRepositoryEnMemoria`; los fakes de test (`fakeProjectRepository`,
`fakeIntegrantesRepository`, `integrantesConActor`, `scrumMasterConActor` y
`stubProjectRepository` de `internal/http`) devuelven `false` salvo que el test diga lo
contrario.

**Rationale**: el guardia ya recibe ese puerto en todos los servicios; agregar un método no
cambia ningún constructor ni el cableado de `internal/app`. Si el proyecto no existe
devuelve `false`: el 404 ya lo resolvió el servicio antes (orden de FR-006 de HU-21).

**Alternatives considered**:
- *Puerto nuevo `CierreProyectoRepository`*: obliga a cambiar el constructor y el cableado
  de 17 servicios, los pasos de Godog y los tests; mucho ruido para un booleano.
- *Devolver el cierre en `ListarIntegrantes`*: mezcla dos conceptos en `domain.Member`.
- *Renombrar el puerto a `ControlAccesoRepository`*: más preciso, pero toca todos los
  servicios; se actualiza solo su comentario.

## 3. Regla de estado (FR-002/FR-003)

**Decision**: `Project` suma `FechaCierre *time.Time` y `CalcularEstadoProyecto` agrega
una primera regla: `if p.FechaCierre != nil { return ProyectoFinalizado }`. El resto de la
función no cambia.

**Rationale**: todos los consumidores del estado (status, reporte PDF, métricas, etc.)
llaman a esa función con el `Project` leído del repositorio, así que basta con que el
repositorio lea la columna. Cumple «cerrado ⇒ Finalizado; abierto ⇒ HU-13 sin cambios».

**Alternatives considered**: la propuesta original de D16 («sin cierre manual no hay
Finalizado») se descartó en Clarifications.

## 4. Reglas de transición en el dominio

**Decision**: dos métodos de `Project`, puros y con TDD:

- `Cerrar(ahora time.Time, sprints []Sprint) (Project, error)`:
  `ErrProyectoYaCerrado` si ya tiene fecha de cierre; `ErrSprintActivoAlCerrarProyecto`
  si algún Sprint está `Activo`; si no, copia con `FechaCierre` = día civil de `ahora`
  (mismo criterio que `inicioDelDiaCivil` de HU-13).
- `Reabrir() (Project, error)`: `ErrProyectoNoCerrado` si no tiene fecha de cierre; si
  no, copia con `FechaCierre = nil`.

Orden de los dos 409 al cerrar: primero «ya está cerrado» (un proyecto cerrado no puede
tener un Sprint Activo de todos modos, porque iniciarlo queda bloqueado).

**Rationale**: Factory/métodos de entidad que validan invariantes (constitución IV), igual
que `ConDatosEditados` y `ConHorasPorStoryPoint`. Testeable sin base.

## 5. Atomicidad y carreras del cierre (FR-017)

**Decision**: dos métodos nuevos en `ProjectRepository`, con `UPDATE` condicional (el mismo
patrón de #158 para iniciar un Sprint):

- `Cerrar(ctx, id, fecha)`:
  `UPDATE proyectos SET fecha_cierre = ? WHERE id = ? AND fecha_cierre IS NULL AND NOT EXISTS (SELECT 1 FROM sprints WHERE proyecto_id = proyectos.id AND estado = 'Activo')`.
  Con 0 filas afectadas, relee el proyecto y sus Sprints y devuelve el error que
  corresponda: `ErrProyectoNoEncontrado`, `ErrProyectoYaCerrado` o
  `ErrSprintActivoAlCerrarProyecto`.
- `Reabrir(ctx, id)`:
  `UPDATE proyectos SET fecha_cierre = NULL WHERE id = ? AND fecha_cierre IS NOT NULL`;
  con 0 filas, `ErrProyectoNoEncontrado` o `ErrProyectoNoCerrado`.

`Actualizar` (editar el proyecto) **no** escribe `fecha_cierre`, así que una edición nunca
la pisa.

**Rationale**: el dominio da el mensaje; la condición del `UPDATE` cubre dos cierres
simultáneos y un Sprint que se inicia entre la lectura y la escritura. Sin transacción
explícita: es una sola sentencia.

## 6. Límite conocido: escritura concurrente con el cierre

Una escritura que pasa `exigirProyectoAbierto` justo antes de que otro pedido cierre el
proyecto puede completarse después del cierre (chequeo y escritura son sentencias
distintas). La ventana es de milisegundos, la base usa una sola conexión y el efecto es el
mismo que si la escritura hubiera llegado un instante antes. Se **declara** como límite en
`MODELO-DE-DATOS.md` §5 y no se cierra (ver alternativa descartada en §1). Si el equipo lo
considera necesario, es un fix aparte.

## 7. Esquema: migración 008

**Decision**: `internal/repository/migraciones/008_cierre_proyecto.sql`:

```sql
ALTER TABLE proyectos ADD COLUMN fecha_cierre TEXT CHECK (date(fecha_cierre) IS fecha_cierre);
```

`ADD COLUMN` con `CHECK` está permitido en SQLite (la columna es anulable y sin valor por
defecto, así que las filas existentes quedan en `NULL`, que pasa el `CHECK`). No hace falta
recrear la tabla (§9 punto 4 no aplica). Los proyectos existentes quedan abiertos.

`MODELO-DE-DATOS.md`, en el mismo PR: §2 (lista de fechas de calendario), §3 (diagrama),
§4.1 (columna y `CHECK`), §4.10 (fila de `NULL`), §5 (filas: «solo lectura» y «cerrar sin
Sprint Activo», con el límite de §6), §6 (fuente del estado: `+ fecha_cierre`), §7 (D16
**cerrada el 2026-11-11** con la decisión del issue), §8 (HU-13/#300) y §9 (008 en la
tabla de números reservados y casos del test de esquema).

Tests: `esquema_test.go` suma los casos de D12 para `fecha_cierre` (`'2026-10-9'`,
`'basura'`, `'2026-02-30'` rechazadas; `'2026-10-09'` y `NULL` válidas), y un test de
migración (como `migracion_007_test.go`) verifica que una base con proyectos previos migra
y los deja abiertos.

## 8. API y traducción de errores

**Decision**:
- Rutas `POST /projects/{id}/cerrar` y `POST /projects/{id}/reabrir`, nivel
  `SesionScrumMaster`, sin cuerpo. Responden 200 con el proyecto (`proyectoResponse`) más
  `estado` (ver [contracts/](contracts/)).
- `proyectoResponse` suma `fecha_cierre` (`YYYY-MM-DD` o `null`): lo devuelven el detalle,
  «Mis proyectos», crear y editar. `GET /projects/{id}/status` no cambia de forma (sigue
  `{"estado": ...}`), solo su valor.
- `escribirError` (único traductor, `internal/http/sprint_handler.go`) mapea a 409 los
  cuatro errores nuevos: `ErrProyectoCerrado`, `ErrProyectoYaCerrado`,
  `ErrProyectoNoCerrado`, `ErrSprintActivoAlCerrarProyecto`.
- `nivelesDeLaSpec` (`internal/app/rutas_permisos_test.go`) suma las dos rutas.

**Rationale**: verbos de acción como `POST /sprints/{id}/iniciar|cerrar`, que ya existen.

## 9. Cómo se prueba que *todas* las escrituras quedan bloqueadas (SC-001, FR-020)

**Decision**: un test de integración en `internal/app` (sobre `clienteDePrueba`, como
`pertenencia_test.go`) con una tabla de **cada** ruta de nivel `SesionIntegrante` y
`SesionScrumMaster` (salvo cerrar y reabrir), cuerpo válido, contra un proyecto cerrado:
espera 409 «El proyecto está cerrado» y verifica que no cambió nada. Un segundo test
compara esa tabla con `apihttp.Rutas`: si mañana se agrega una ruta de escritura y no se
suma a la tabla, falla (mismo patrón que `TestRutas_CadaRutaTieneElNivelDeLaSpec`).

Godog cubre los escenarios de la spec con un `Esquema del escenario` de escrituras
representativas (al menos una de integrante y una de Scrum Master, FR-019), no las 17
rutas: la cobertura exhaustiva queda en el test de `internal/app`.

## 10. Convivencia con el PR #304 (spec 043, abierto)

#304 toca `internal/domain/project.go`, `internal/service/iniciar_sprint.go`,
`features/steps_estado_proyecto.go` y `MODELO-DE-DATOS.md` §5. Los cambios no se pisan en
lógica (043 agrega `AdmiteSprint`; 044 agrega `FechaCierre`, `Cerrar` y `Reabrir`), pero
habrá conflictos de texto. Si #304 se mergea primero, se trae `main` a esta rama antes del
PR. La migración 008 no la reserva nadie más (#304 no cambia el esquema).
