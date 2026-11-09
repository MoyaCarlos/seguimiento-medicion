# Data Model: Un único Scrum Master por proyecto (#236)

La definición de tablas y columnas vive en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md), §4.3 (`integrantes`). Este archivo **no** las redefine: solo describe el cambio.

## Cambio de esquema

Sin tablas ni columnas nuevas. Una migración `internal/repository/migraciones/007_scrum_master_unico.sql`, en dos pasos y en este orden:

1. **Corrección de datos**: en cada proyecto con más de un `scrum_master`, queda el de `creado_en` más antiguo (empate: menor `usuario_id`) y los demás pasan a `product_builder`. Proyectos con uno o con ningún Scrum Master no se tocan. No se borra ninguna fila.
2. **Índice único parcial**: `integrantes_un_scrum_master`, único sobre `proyecto_id`, solo para las filas con `rol = 'scrum_master'`.

El archivo comienza con un comentario que explica el criterio (FR-012).

## Reglas que impone

| Regla | Dónde se impone |
|---|---|
| Un usuario tiene un solo rol por proyecto | PK `(proyecto_id, usuario_id)` (ya existente) |
| A lo sumo un `scrum_master` por proyecto | Índice nuevo (base) + chequeo previo del service (dominio) |
| Varios `product_builder` por proyecto | Sin restricción |
| El Scrum Master es el dueño del proyecto | Convención: no se guarda aparte (D14 deja de hacer falta) |

## Cambios en `MODELO-DE-DATOS.md` (en el mismo PR que la migración)

- **§4.3**: reemplazar "Hoy no se impone un único Scrum Master por proyecto (ver D4)" por la regla y el nombre del índice.
- **§7, D4**: pasar a **cerrada el 2026-11-03** (decisión del Scrum Master): un solo Scrum Master por proyecto; la regla es del dominio y de la base; la verifica el service antes de insertar. Ajustar el encabezado de la sección ("hoy están pendientes…") para sacar D4 y D14.
- **§7, D14**: marcar **"no necesaria"**: el dueño es el Scrum Master del proyecto, que es quien lo creó.
- **§9, tabla de números reservados**: agregar `007 | Índice único de Scrum Master por proyecto, con corrección previa de duplicados (el de creado_en más antiguo; empate, el de menor usuario_id) | #236`. Anotar que la 006 es de HU-20.
- **§9, punto 6 (test de esquema)**: sumar los casos del FR-014.
- **§5 (dónde vive cada regla)**: agregar la fila de esta regla.

## Entidades

- **Integrante (membresía)**: `proyecto_id`, `usuario_id`, `rol` (`scrum_master` | `product_builder`), `creado_en`. Sin campos nuevos.

## Transiciones

Ninguna nueva: no existe cambio de rol ni baja (HU-16, #175).
