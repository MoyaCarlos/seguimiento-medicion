# Data Model: Consenso y Registro de la Estimación Acordada (HU-14)

**Fecha**: 2026-11-03 | **Spec**: [spec.md](./spec.md) | **Research**: [research.md](./research.md)

## Esquema

Sin cambios. Esta historia usa lo definido en
[MODELO-DE-DATOS.md §4.7](../../MODELO-DE-DATOS.md#47-rondas_votacion-y-votos--hu-02-hu-14-tablas-nuevas-migración-003)
(`rondas_votacion`, `votos`) y
[§4.4](../../MODELO-DE-DATOS.md#44-historias--hu-01-hu-15-hu-06-hu-11-hu-14-antes-backlog_items)
(`historias.estimacion_sp`), sin migración nueva.

| Restricción o regla | Origen | Uso en HU-14 |
|---|---|---|
| `UNIQUE (historia_id, numero)` | §4.7 | frena dos aperturas simultáneas de la misma ronda |
| Ronda actual = la de mayor `numero` | §4.7 | votar, revelar y consultar pasan solos a la ronda nueva |
| La ronda anterior se reveló antes de abrir otra | §5, servicio | `RondaVotacion.Siguiente` |
| D2b sobre `estimacion_sp` | §5, servicio | función compartida con HU-15 (R6) |

## Dominio (`internal/domain/votacion.go`)

### `RondaVotacion` (métodos nuevos)

- `(r RondaVotacion) Siguiente() (RondaVotacion, error)`
  - sin revelar → `ErrRondaNoRevelada`;
  - si no, devuelve una ronda de la misma historia con `Numero + 1`, sin votos ni revelar.
- `(r RondaVotacion) Diferencias(regla ReglaConsenso) (Diferencias, bool)`: `false` si no está revelada.

### `ReglaConsenso` (Strategy)

```go
type ReglaConsenso interface {
    HayDiferencias(valores []int) bool
}
```

| Implementación | Regla |
|---|---|
| `CoincidenciaExacta` | `true` si algún valor es distinto de otro, sin importar la distancia |

Casos de test: `[3,5,8]` y `[3,5]` → `true`; `[5,5,5]`, `[8]` y `[]` → `false`.

### `Diferencias`

| Campo | Tipo | Regla |
|---|---|---|
| `Hay` | `bool` | resultado de la regla |
| `Minimo`, `Maximo` | `int` | extremos de los votos; solo tienen sentido si `Hay` es `true` |

### `BacklogItem` (método nuevo, `backlog_item.go`)

- `(b BacklogItem) ConEstimacion(sp int) (BacklogItem, error)`: copia con los nuevos Story Points;
  `ErrEstimacionInvalida` si el valor está fuera de la baraja.

### Errores (`errors.go`)

`ErrRondaNoRevelada = ValidationError{Campo: "estado", Mensaje: "la ronda actual todavía no fue revelada"}` y
`ErrRondaEnHistoriaCompletada = ValidationError{Campo: "estado", Mensaje: "no se puede abrir una ronda en una historia completada"}` y
`ErrRondaYaAbierta = errors.New("otra persona abrió esta ronda mientras tanto: recargá la votación")`.
El resto se reutiliza: ver [research.md R10](./research.md#r10-errores-y-códigos-http).

### Transiciones de una ronda

```text
                                       Revelar (votaron todos)
 (ronda N) Abierta ───────────────────────────────▶ Revelada
                                                        │
                                 Siguiente              ▼
 (ronda N+1) Abierta ◀───────────────────────── queda como historial
```

Registrar la estimación no cambia ninguna ronda.

## Servicios (`internal/service/`)

| Service | Entrada | Pasos |
|---|---|---|
| `AbrirNuevaRonda` | `historiaID`, `usuarioID` | historia existe → Scrum Master → no `Completada` → ronda actual → `Siguiente` → `CrearRonda` (`ErrVotacionYaIniciada` → `ErrRondaYaAbierta`) |
| `RegistrarEstimacion` | `historia` (la lee el handler, que comprueba que existe), `usuarioID`, `valor` | Scrum Master → `ConEstimacion` → ronda actual revelada → D2b → `GuardarEstimacion` |
| `ConsultarHistorialRondas` | `historiaID` | historia existe → `ListarRondasReveladas` → agrega nombre y rol a cada voto |
| `ConsultarVotacion` (cambia) | sin cambios | `VistaVotacion` suma `Diferencias *domain.Diferencias`, `nil` sin revelar |

### Vista del historial

| Campo | Tipo |
|---|---|
| `HistoriaID` | `int64` |
| `Rondas[]` | `Numero`, `ReveladaEn`, `Votos[]` (`UsuarioID`, `Nombre`, `Rol`, `Valor`) |

## Puertos (`internal/repository/`)

```go
// votacion_repository.go: se agregan dos métodos a VotacionRepository.
ListarRondasReveladas(ctx context.Context, historiaID int64) ([]domain.RondaVotacion, error)
RondaActualRevelada(ctx context.Context, historiaID int64) (bool, error)

// backlog_repository.go: puerto nuevo, separado de HistoriaRepository (ISP).
type EstimacionHistoriaRepository interface {
    GuardarEstimacion(ctx context.Context, id int64, sp int) error
}
```

| Método | SQLite | Errores |
|---|---|---|
| `ListarRondasReveladas` | una consulta `rondas_votacion LEFT JOIN votos`, `WHERE revelada_en IS NOT NULL`, por `numero` y `usuario_id` (R9) | sin filas → lista vacía |
| `RondaActualRevelada` | la ronda de mayor `numero`, sin leer `votos` | sin filas → `ErrVotacionNoIniciada` |
| `GuardarEstimacion` | `UPDATE historias SET estimacion_sp = ? WHERE id = ?` (R7) | 0 filas → `ErrHistoriaNoEncontrada` |

**Fake de los tests**: `fakeVotacionRepository` guarda hoy una ronda por historia. Debe conservar las
rondas anteriores, mantener el `ErrVotacionYaIniciada` cuando el número no es mayor que el de la ronda
actual, y asignar ids únicos. Los tests existentes siguen sembrando la ronda actual en `rondas`.
