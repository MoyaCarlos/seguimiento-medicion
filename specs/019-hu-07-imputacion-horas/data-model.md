# Data Model: Imputación de Horas Trabajadas (HU-07)

**Fecha**: 2026-10-28 | **Spec**: [spec.md](./spec.md) | **Research**: [research.md](./research.md)

## Esquema (persistencia)

La tabla ya está diseñada; este documento **no la redefine**:

- [`MODELO-DE-DATOS.md` §4.8 `registros_horas`](../../MODELO-DE-DATOS.md#48-registros_horas--hu-07-tabla-nueva-migración-004):
  columnas, FK compuestas y `CHECK`.
- §2 (convenciones), §5 (dónde vive cada regla y traducción de errores), D8 (minutos
  enteros), D12 (formato de fecha), §9 (migraciones; la 004 está reservada para HU-07).

Esta historia agrega `internal/repository/migraciones/004_registros_horas.sql` con ese DDL y
no cambia el documento.

## Entidades del dominio

### `RegistroHoras` (nueva, `internal/domain/registro_horas.go`)

| Campo | Tipo Go | Origen | Regla |
|---|---|---|---|
| `ID` | `int64` | base | asignado al guardar |
| `ProyectoID` | `int64` | historia | el de la historia |
| `HistoriaID` | `int64` | ruta | historia existente |
| `UsuarioID` | `int64` | cuerpo (`integrante_id`) | integrante del proyecto (service) |
| `SprintID` | `int64` | historia | Sprint `Activo` de la historia al cargar; no se recalcula (FR-003) |
| `Fecha` | `time.Time` (00:00 UTC) | cuerpo | `AAAA-MM-DD` real; `inicio del Sprint ≤ fecha ≤ min(hoy, fin del Sprint)` |
| `Actividad` | `string` | cuerpo | sin espacios externos; no vacía; ≤ 500 runas |
| `Minutos` | `int` | cuerpo (`horas`) | `1 ≤ minutos ≤ 1440`; horas decimales que den minutos enteros (R2) |
| `CreadoEn` | `time.Time` | base | auditoría (RFC 3339 UTC); sin reglas |

**Factory**:

```text
NewRegistroHoras(historia BacklogItem, sprint Sprint, usuarioID int64,
                 horas, actividad, fecha string, ahora time.Time) (RegistroHoras, error)
```

Valida, en este orden, horas → actividad → fecha (formato, después rango). Ante el primer
error devuelve un `ValidationError` y ninguna carga. Supone que el service ya verificó el
Sprint activo y el integrante (R4).

**Funciones del dominio**:

- `MinutosDesdeHoras(texto string) (int, error)`: conversión exacta con `big.Rat` (R2).
- `ValidarSprintParaImputar(historia BacklogItem, sprint Sprint) error`: la historia está
  asignada a ese Sprint y el Sprint está `Activo`; si no, `ErrHistoriaFueraDelSprintActivo`.
- `TotalMinutos(registros []RegistroHorasDetalle) int`: suma exacta (FR-013).

**Errores nuevos** (`internal/domain/errors.go` o junto a la entidad):

| Error | Campo | Mensaje |
|---|---|---|
| `ErrHistoriaFueraDelSprintActivo` | `historia_id` | "solo se imputan horas a historias del Sprint activo" |
| `ErrIntegranteAjeno` | `integrante_id` | "el integrante no pertenece al proyecto de la historia" |
| (horas) | `horas` | "las horas son obligatorias" / "las horas deben ser un número" / "las horas deben ser mayores que cero" / "las horas no pueden superar 24 por carga" / "las horas deben equivaler a minutos enteros" |
| (actividad) | `actividad` | "la actividad es obligatoria" / "la actividad no puede superar los 500 caracteres" |
| (fecha) | `fecha` | "formato esperado AAAA-MM-DD" / "la fecha no puede ser posterior a hoy" / "la fecha no puede ser anterior al inicio del Sprint" / "la fecha no puede ser posterior al fin del Sprint" |

Los textos son orientativos; los fija el test RED de cada caso.

### `RegistroHorasDetalle` (vista de lectura)

`RegistroHoras` + `NombreIntegrante string`, armado con `JOIN usuarios`. Igual que
`domain.Member`, es una vista para mostrar; no se construye por factory.

### `EsfuerzoReal` (vista, no se guarda)

Total de una historia = `TotalMinutos(detalle)`; se muestra en horas (`minutos / 60`).
Una historia sin cargas tiene total 0 y detalle vacío.

## Puertos (`internal/repository/registro_horas_repository.go`)

```text
RegistroHorasRepository
  Registrar(ctx, r domain.RegistroHoras) (domain.RegistroHoras, error)
      INSERT condicional (R5): si la historia ya no está en ese Sprint Activo
      → domain.ErrHistoriaFueraDelSprintActivo. FK → ErrIntegridadReferencial.
  ListarPorHistoria(ctx, historiaID int64) ([]domain.RegistroHorasDetalle, error)
      ORDER BY fecha, id. Lista vacía (no nil) si no hay cargas.

IntegrantesRepository
  ListarIntegrantes(ctx, proyectoID int64) ([]domain.Member, error)
      Ya implementado por SQLiteProjectRepository (ISP, R6).
```

Lecturas que se reutilizan sin cambios: `HistoriaRepository.ObtenerPorID`
(`SQLiteBacklogRepository`) y `SprintRepository.ObtenerPorID`.

## Casos de uso (`internal/service`)

| Service | Depende de | Hace |
|---|---|---|
| `ImputarHoras` | `HistoriaRepository`, `SprintRepository`, `IntegrantesRepository`, `RegistroHorasRepository` | historia → Sprint activo → integrante → `NewRegistroHoras` → `Registrar` |
| `ConsultarEsfuerzoReal` | `HistoriaRepository`, `RegistroHorasRepository` | historia existe → `ListarPorHistoria` → total |

## Estados

Una carga no tiene estados ni transiciones: se crea y no se edita ni se borra (FR-015, D9).
Los estados que importan son de otras entidades: el Sprint tiene que estar `Activo`; el
estado de la historia no influye (Clarifications).
