# Data Model: Autenticación de usuarios (HU-20)

Migración `internal/repository/migraciones/006_autenticacion.sql`. La 007 es de #236 y no
se usa. Se documenta en `MODELO-DE-DATOS.md` en el mismo PR.

## `usuarios` (cambia: HU-04 → HU-20)

| Columna | Tipo | Notas |
|---|---|---|
| `contrasena_hash` | TEXT **NULL** | **nueva.** `pbkdf2-sha256$<iter>$<sal>$<hash>`. `NULL` = usuario previo o creado por HU-04: no puede ingresar. Nunca se devuelve. |

Sin `CHECK` de largo (§2: los largos viven en el dominio). El resto de la tabla no cambia.

## `sesiones` (nueva)

| Columna | Tipo | Notas |
|---|---|---|
| `id` | PK AUTOINCREMENT | |
| `usuario_id` | INTEGER NOT NULL | FK → `usuarios(id)`; sin `ON DELETE` (§2) |
| `token_hash` | TEXT NOT NULL **UK** | SHA-256 hexadecimal del token; el token en claro no se guarda |
| `creado_en` | TEXT NOT NULL | RFC 3339 UTC |
| `vence_en` | TEXT NOT NULL | RFC 3339 UTC; `creado_en` + 7 días. **Entra en una regla de negocio** (desvío de §2, ver [research.md](research.md) D5) |
| `cerrada_en` | TEXT **NULL** | instante del cierre de sesión; `NULL` = no cerrada. Cerrar no borra la fila (D6) |

Sin `CHECK` de fecha de calendario: son instantes, no fechas (§2 limita el `date()` a fechas
de calendario). `UNIQUE (token_hash)` ya crea el índice que usa «quién soy».

## Entidades de dominio

- **`Sesion`**: `ID`, `UsuarioID`, `CreadoEn`, `VenceEn`, `CerradaEn *time.Time`.
  `Vigente(ahora) bool` = no cerrada y `ahora` anterior a `VenceEn`. `DuracionSesion = 7 * 24h`.
- **`SesionConUsuario`**: `Sesion` + `Usuario` (id y nombre), lo que devuelve la búsqueda
  por huella.
- **Reglas puras**: `ValidarNombreUsuario(string) error` y `ValidarContrasena(string) error`.

## Transiciones de una sesión

```text
(emitida) ──ahora ≥ vence_en──▶ vencida
    │
    └──cerrar──▶ cerrada      (ambas: no vigente, no se reactiva)
```

## Reglas de integridad

- Unicidad de identidad: `usuarios.nombre_normalizado UNIQUE` (ya existe) cubre el registro
  repetido, también contra usuarios previos sin contraseña y contra los creados por HU-04.
- Un usuario puede tener varias sesiones; cerrar una no toca las otras.
- Violación `SQLITE_CONSTRAINT_UNIQUE` en `usuarios` → `ErrUsuarioDuplicado` → 409.
- Violación de FK en `sesiones.usuario_id` → `ErrIntegridadReferencial` (ya mapeada), no
  alcanzable desde la API.

## Cambios en `MODELO-DE-DATOS.md`

§2 (nota sobre `vence_en` y `cerrada_en`), §3 (diagrama: `sesiones`, `usuarios.contrasena_hash`),
§4.2 (columna nueva), §4.11 (tabla `sesiones`), §8 (HU-20 → `usuarios`, `sesiones`, 006) y §9
(fila de la migración 006).
