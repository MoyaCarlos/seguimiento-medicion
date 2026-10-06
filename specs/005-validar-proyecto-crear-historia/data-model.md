# Data Model: Validar proyecto existente al crear historia de backlog (Issue #19)

## Sin cambios en la persistencia

Esta corrección **no agrega tablas, columnas, campos ni restricciones**. Se reutilizan las
entidades y tablas existentes:

- `projects` (HU-04): `id, name, description, start_date, end_date, created_at`.
- `backlog_items` (HU-01): `id, proyecto_id, titulo, descripcion, prioridad, estado, valor_negocio, ...`.

`internal/repository/migrate.go`, `internal/domain/backlog_item.go` y
`internal/repository/backlog_repository.go` **no se modifican**. No se agrega `FOREIGN KEY`
(FR-008, fuera de alcance).

## Entidades involucradas (sin cambios de forma)

- **Proyecto** (ya existente): entidad referenciada por la historia a través de
  `BacklogItem.ProyectoID`. Esta corrección solo **consulta su existencia** por identificador
  (`ProjectRepository.ObtenerPorID`); no la crea, edita ni elimina.
- **Historia del Product Backlog** (ya existente, `domain.BacklogItem`): se intenta crear. Sus
  invariantes no cambian (`NewBacklogItem` valida `proyecto_id > 0`, título, descripción,
  prioridad y valor de negocio). Se agrega una **precondición de existencia del proyecto
  referenciado**, no un atributo nuevo.

## Relación (sin cambios de esquema)

```text
Project (1) ----< (N) BacklogItem
```

La relación ya existía por `BacklogItem.ProyectoID`; lo que faltaba era **hacerla cumplir en
la aplicación**. Esta corrección la refuerza en el caso de uso, no en el motor de base de datos.

## Regla de validación del caso de uso (orden)

`CrearHistoriaBacklog.Ejecutar(ctx, input)`:

1. `domain.NewBacklogItem(input.ProyectoID, input.Titulo, input.Descripcion, input.Prioridad,
   input.ValorNegocio)`.
   - Si falla alguna invariante → devolver `domain.ValidationError` (no se consulta el proyecto).
2. `s.proyectos.ObtenerPorID(ctx, input.ProyectoID)`.
   - Si el proyecto no existe → `ObtenerPorID` devuelve `domain.ErrProyectoNoEncontrado` y el
     caso de uso lo **propaga tal cual** (no persiste nada).
   - Si devuelve otro error de persistencia → se **propaga tal cual** (no persiste nada, FR-006).
3. `s.repo.Guardar(ctx, item)` → persiste la historia y devuelve el `BacklogItem` con `ID`.

Comportamiento observable por caso:

| Entrada | Resultado del caso de uso | ¿Persiste? |
|---|---|---|
| Historia inválida (p. ej. título vacío), proyecto exista o no | `ValidationError` del dominio | No |
| Historia válida, `ProyectoID` inexistente | `domain.ErrProyectoNoEncontrado` | No |
| Historia válida, error genérico de `ObtenerPorID` | El error original propagado | No |
| Historia válida, proyecto existente | `BacklogItem` creado (estado `Nueva`) | Sí |

Nota: `ObtenerPorID` ya devuelve `domain.ErrProyectoNoEncontrado` en ambas implementaciones
(`sqlite_project.go` y `project_repository_en_memoria.go`); el caso de uso no traduce ni
envuelve el error.
