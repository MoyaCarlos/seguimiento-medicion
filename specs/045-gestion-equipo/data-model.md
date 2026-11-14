# Data Model: Gestión del equipo del proyecto (HU-16)

La fuente de verdad es [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md). Esta historia
agrega **una columna** y cierra la decisión **D15**; el cambio va en el mismo PR que la
migración (regla de `AGENTS.md`).

## Cambio de esquema — migración `009_baja_integrantes.sql`

| Tabla | Columna | Tipo | Notas |
|---|---|---|---|
| `integrantes` (§4.3) | `fecha_baja` | TEXT, anulable | `NULL` = activo. Día `YYYY-MM-DD`, `CHECK (fecha_baja IS NULL OR date(fecha_baja) IS fecha_baja)` (D12) |

- `ALTER TABLE integrantes ADD COLUMN fecha_baja TEXT CHECK (...)`: las filas existentes
  quedan en `NULL` (todas activas); no hace falta corregir datos.
- La PK `(proyecto_id, usuario_id)` y el índice `integrantes_un_scrum_master` no cambian: el
  Scrum Master nunca se da de baja, así que el índice parcial no necesita filtrar por
  `fecha_baja`.

## Cambios de texto en `MODELO-DE-DATOS.md`

- §4.3: fila `fecha_baja`; párrafo «baja lógica»: quitar = fecha de baja; reactivar = misma
  fila con `fecha_baja = NULL`; el Scrum Master no se da de baja.
- §4.7 (votación): al dar de baja se borran los votos de esa persona en rondas no reveladas
  del proyecto; las reveladas no se tocan.
- §5 (dónde vive cada regla): «Integrante = activo» (consultas), «al Scrum Master no se lo
  quita» (dominio + `UPDATE` condicional), «una baja no se aplica dos veces» (`UPDATE`
  condicional).
- §7: D15 pasa a **Cerrada el 2026-11-12** (confirmada por el equipo).
- §8 y lista de migraciones: `009` → HU-16 (#175, spec 045).

## Entidades (dominio)

| Concepto | Tipo Go | Cambio |
|---|---|---|
| Integrante (vínculo) | `domain.Membership` | sin cambio (la baja es una operación del repositorio) |
| Integrante con nombre | `domain.Member` | método `ValidarBaja() error`: `ErrScrumMasterNoSeQuita` si `Role == RolScrumMaster` |

## Puertos (`internal/repository`)

| Método | Cambio |
|---|---|
| `ProjectRepository.ListarIntegrantes` | solo activos (`fecha_baja IS NULL`) |
| `ProjectRepository.ListarPorIntegrante` | solo proyectos donde la persona está activa |
| `ProjectRepository.ListarIntegrantesHistoricos` | **nuevo**: activos y dados de baja (historial de rondas, votación) |
| `ProjectRepository.AgregarIntegrante` | reactiva la fila si estaba dada de baja; duplicado si estaba activa |
| `ProjectRepository.DarDeBajaIntegrante(ctx, proyectoID, usuarioID, fecha)` | **nuevo**: `UPDATE` condicional + borrado de votos no revelados, en una transacción |

## Estados del vínculo

```text
(no existe) --alta--> Activo --baja--> Dado de baja --alta (reactivar)--> Activo
Scrum Master: Activo siempre (baja → 409)
```
