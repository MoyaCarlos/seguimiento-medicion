# Data Model: Renombrar el propio usuario (HU-17)

## Sin cambios de esquema

La tabla `usuarios` ya existe desde la migración 001:

| Columna | Tipo | Notas |
|---|---|---|
| `id` | PK | |
| `nombre` | TEXT NOT NULL | nombre visible, tal como lo escribió la persona |
| `nombre_normalizado` | TEXT NOT NULL **UK** | minúsculas y sin espacios externos; es la identidad y el dato de ingreso |
| `contrasena_hash` | TEXT? | HU-20; no se toca |
| `creado_en` | TEXT NOT NULL | no se toca |

No hay migración nueva. El renombrado actualiza `nombre` y `nombre_normalizado` en la misma
sentencia:

```sql
UPDATE usuarios
   SET nombre = ?, nombre_normalizado = ?
 WHERE id = ?
```

## Reglas de integridad

- `usuarios.nombre_normalizado UNIQUE` garantiza que dos usuarios no compartan identidad,
  también en renombrados simultáneos. La violación `SQLITE_CONSTRAINT_UNIQUE` se traduce a
  `ErrUsuarioDuplicado` y el servicio responde `409`.
- Si el `id` no existe, la sentencia no afecta filas y el repositorio devuelve
  `domain.ErrNoEncontrado` (404). No es alcanzable desde la API porque el usuario sale de una
  sesión vigente, pero se define por completitud.
- `nombre` y `nombre_normalizado` se escriben juntos; no puede quedar una identidad
  desincronizada.

## Relaciones que no cambian

```text
usuarios (1) ----< (N) integrantes
usuarios (1) ----< (N) votos
usuarios (1) ----< (N) registros_horas
usuarios (1) ----< (N) sesiones
usuarios (1) ----< (N) defectos (referencias por usuario_id en el modelo vigente)
```

Todas las relaciones usan `usuario_id`, por lo que el historial y las sesiones siguen
asociados a la misma persona. Las lecturas que muestran el nombre hacen JOIN con `usuarios`
y ven el nombre nuevo automáticamente.

## Entidades de dominio

- `domain.User{ID, Nombre, NombreNormalizado, CreadoEn}` ya existe. El renombrado reutiliza
  `domain.ValidarNombreUsuario` y `domain.NormalizarNombre`; no se agrega un campo ni un
  método de dominio nuevo.

## Documentación asociada

- `MODELO-DE-DATOS.md` §4.2: nota de que el renombrado actualiza `nombre` y
  `nombre_normalizado` juntos.
- `MODELO-DE-DATOS.md` §5: fila «El renombrado del propio usuario es atómico y rechaza
  duplicados sin distinguir mayúsculas (HU-17)» → repositorio (una sola sentencia, índice
  único).
- `MODELO-DE-DATOS.md` §8: fila `HU-17 Renombrar usuario | usuarios.nombre,
  usuarios.nombre_normalizado | usuarios | sin migración`.
