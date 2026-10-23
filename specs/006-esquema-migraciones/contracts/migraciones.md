# Contrato: runner de migraciones

Contrato para quien agregue una migración (002 en adelante: #117, HU-02, HU-07, HU-08).
El detalle de las tablas está en [MODELO-DE-DATOS.md §9](../../../MODELO-DE-DATOS.md#9-migraciones).

## Archivos

- Ruta: `internal/repository/migraciones/NNN_nombre.sql`, con `NNN` de **tres dígitos** y el número reservado en la tabla de §9 punto 3.
- Un archivo puede tener varias sentencias separadas por `;`. Se ejecuta entero dentro de una transacción.
- Una migración mergeada a `main` **no se edita**; un cambio es una migración nueva (§9 punto 2).
- Un nombre que no cumple `NNN_nombre.sql` hace fallar el arranque.
- Dos archivos con el mismo número hacen fallar el arranque.

## Función

```go
// AplicarMigraciones aplica, en orden, las migraciones embebidas que la base todavía no tiene.
func AplicarMigraciones(ctx context.Context, db *sql.DB) error
```

| Situación | Resultado |
|---|---|
| Base vacía | Crea `migraciones_aplicadas`, aplica todas en orden y registra cada una. `nil`. |
| Todas aplicadas | No ejecuta nada. `nil`. |
| Faltan algunas (con o sin huecos de número) | Aplica solo las que faltan, en orden. `nil`. |
| Una migración falla | Esa migración se revierte entera y no se registra; las anteriores quedan. Error que nombra el archivo. |
| La base tiene `projects`, `users`, `project_members` o `backlog_items` | No aplica nada. `ErrEsquemaViejo` (pide borrar `seguimiento.db` y reiniciar). |
| La tabla de control tiene una versión que el binario no trae | Se ignora. |

Lo llaman `cmd/api/main.go`, los steps de Godog y el helper de los tests de repositorio,
siempre después de `AbrirSQLite`.
