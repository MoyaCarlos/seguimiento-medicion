# Research: Gestión del equipo del proyecto (HU-16)

Decisiones de diseño tomadas en `/speckit.plan` (2026-11-12) a partir del código de `main`
(merge de #304) y de la spec.

## 1. Baja lógica: columna `integrantes.fecha_baja` (D15)

- **Decision**: columna anulable `fecha_baja TEXT` con `CHECK (fecha_baja IS NULL OR
  date(fecha_baja) IS fecha_baja)` (D12, día `YYYY-MM-DD`), migración `009`. `NULL` = activo.
- **Rationale**: es la propuesta de D15, confirmada por el equipo el 2026-11-12. La PK
  `(proyecto_id, usuario_id)` se mantiene: un solo vínculo por persona y proyecto, también
  después de la baja, y las FK de `registros_horas` hacia `integrantes` siguen válidas.
- **Alternatives considered**: borrar la fila (rompe la FK de `registros_horas` y pierde
  el historial); booleano `activo` (no dice desde cuándo; D15 pide fecha).

## 2. «Integrante» pasa a significar «integrante activo»

- **Decision**: `ProjectRepository.ListarIntegrantes` y `ListarPorIntegrante` filtran
  `fecha_baja IS NULL`. Se agrega `ListarIntegrantesHistoricos(proyectoID)` (activos y dados
  de baja), que usan solo las vistas que muestran aportes pasados.
- **Rationale**: todas las reglas que hoy usan `ListarIntegrantes` deben ignorar a quien fue
  dado de baja: pertenencia y rol (`exigirIntegrante`, `exigirScrumMaster`,
  `ObtenerProyecto.ComoIntegrante`), «Mis proyectos», equipo listado, revelar la votación
  (FR-007), imputar horas y el equipo del reporte. Filtrar en un solo lugar las corrige todas
  sin tocar cada servicio (DRY).
- **Usuarios de `ListarIntegrantesHistoricos`**:
  - `ConsultarHistorialRondas`: resuelve nombre y rol de cada voto revelado; sin esto, el
    voto de quien fue dado de baja saldría sin nombre (FR-002).
  - `ConsultarVotacion`: lista a los integrantes activos **más** quien tenga un voto en la
    ronda actual. Si la ronda se reveló y después se quitó a alguien, el promedio (que se
    calcula de los votos guardados) y la lista mostrada siguen coincidiendo.
- Las horas (`sqlite_registro_horas` hace `JOIN usuarios`) y los defectos ya resuelven el
  nombre sin pasar por `integrantes`: no cambian.

## 3. Voto en una ronda abierta (Clarification Q1)

- **Decision**: al dar de baja, en la misma transacción, `DELETE FROM votos WHERE usuario_id = ?
  AND ronda_id IN (rondas no reveladas de historias del proyecto)`.
- **Rationale**: el resultado de una ronda (promedio, diferencias) se recalcula siempre desde
  `votos`; descartar el voto secreto, que nadie vio, es la forma de que no cuente sin otra
  columna. Las rondas reveladas no se tocan (FR-002, FR-018). Elegido por la usuaria.
- **Alternatives considered**: marcar el voto como excluido (segundo cambio de esquema);
  filtrar al calcular (el resultado de una ronda ya revelada cambiaría si alguien se quita
  después).

## 4. Atomicidad de la baja

- **Decision**: `ProjectRepository.DarDeBajaIntegrante(ctx, proyectoID, usuarioID, fecha)`
  en SQLite abre una transacción con:
  1. `UPDATE integrantes SET fecha_baja = ? WHERE proyecto_id = ? AND usuario_id = ? AND
     fecha_baja IS NULL AND rol <> 'scrum_master'`.
  2. Si no afectó filas: lee la fila para distinguir `ErrScrumMasterNoSeQuita` (activa y
     Scrum Master) de `ErrIntegranteNoEncontrado` (no existe o ya estaba dada de baja).
  3. El `DELETE` de votos de §3.
- **Rationale**: la regla del Scrum Master y «no dos bajas» quedan en el `WHERE` del `UPDATE`
  condicional, como #145/#158: dos pedidos simultáneos aplican uno solo (FR-014). El
  servicio igual valida antes para dar el mensaje (patrón «servicio da el mensaje, la base
  cubre la carrera»).
- **Fake en memoria**: aplica la baja sobre su lista; no tiene votos (viven en otro fake), así
  que el descarte de votos se prueba en el test del repositorio SQLite y en Godog. Mismo
  límite que el cierre de Sprint (§4.6 de `MODELO-DE-DATOS.md`).

## 5. Reactivación por el alta existente

- **Decision**: `AgregarIntegrante` en SQLite pasa a `INSERT ... ON CONFLICT (proyecto_id,
  usuario_id) DO UPDATE SET fecha_baja = NULL, rol = excluded.rol WHERE integrantes.fecha_baja
  IS NOT NULL`. Si no afectó filas, la persona ya era integrante activo →
  `ErrMiembroDuplicado` (400 «el integrante ya pertenece al proyecto», como hoy). Se conserva
  `creado_en` original.
- **Rationale**: Escenario 3 «sin duplicarlo» y D15 «la reactivación reutiliza la fila»;
  atómico y sin un endpoint nuevo. El rol siempre es `product_builder`: el servicio ya
  rechaza `scrum_master` con `ErrScrumMasterYaExiste` (409) porque hay un Scrum Master activo.
- **Respuesta**: 201 con el integrante, igual que un alta (la spec no pide distinguirlos).

## 6. Endpoint de baja

- **Decision**: `DELETE /projects/{id}/members/{usuario_id}`, nivel `SesionScrumMaster`,
  responde 204 sin cuerpo.
- **Orden de errores** (FR-013): 401 → 404 proyecto → 403 no integrante → 403 no Scrum
  Master → 409 proyecto cerrado (cuando entre #300) → 404 «integrante no encontrado» → 409
  «No se puede quitar al Scrum Master del proyecto».
- `usuario_id` no numérico o ≤ 0: 400 (como `parsearID`).

## 7. Errores nuevos de dominio

| Error | Mensaje | HTTP |
|---|---|---|
| `domain.ErrIntegranteNoEncontrado` | «integrante no encontrado» | 404 |
| `domain.ErrScrumMasterNoSeQuita` | «No se puede quitar al Scrum Master del proyecto» | 409 |

La regla «al Scrum Master no se lo quita» vive en el dominio (`Member.ValidarBaja()`), con
TDD; el `WHERE` del repositorio la respalda.

## 8. Dependencia con #300 (PR #311)

- La migración es la `009`: va después de la `008` de #300. La rama sale de `main`; antes de
  agregar la migración se mergea `main` con #311 adentro. Si #311 no entró todavía, se frena
  y se avisa (no se renumera por cuenta propia).
- La guardia de proyecto cerrado llega sola: `exigirScrumMaster` la llama desde #300.

## 9. Frontend

- `Configuracion.tsx` (equipo): botón «Quitar» en cada Product Builder, con confirmación, que
  llama a `quitarIntegrante(id, usuarioId)` y saca la fila de la lista; el Scrum Master no
  tiene botón; un error del servidor se muestra con su mensaje.
- La deshabilitación para Product Builders («Solo el Scrum Master puede…») depende del
  frontend de HU-21 (#247); si no está mergeado, se aplica el patrón que deje esa historia.
