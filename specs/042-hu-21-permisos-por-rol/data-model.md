# Data Model: Permisos por rol (HU-21)

**Sin cambios de esquema ni migración nueva.** La fuente de verdad es
[`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md); esta historia solo **lee** tablas que ya
existen.

| Concepto de la spec | Dónde vive | Sección |
|---|---|---|
| Usuario registrado | `usuarios` con `contrasena_hash IS NOT NULL` | §4.2 |
| Integrante y su rol en un proyecto | `integrantes (proyecto_id, usuario_id, rol)`, PK compuesta | §4.3 |
| Un único Scrum Master (D4, #236) | índice único parcial `integrantes_un_scrum_master` | §4.3 y §5 |
| Sesión (quién actúa) | `sesiones` (HU-20) | §4.11 |
| Dueño del proyecto (D14) | el Scrum Master, que es quien lo creó | §7, D14 «No necesaria» |

## Lecturas nuevas (sin columnas nuevas)

- **Mis proyectos** — `ProjectRepository.ListarPorIntegrante(usuarioID)`: proyectos
  unidos con `integrantes` filtrando por `usuario_id`, con el `rol` de esa fila, orden por
  `proyectos.id`. Implementación SQLite y fake en memoria.
- **Usuarios registrados** — `UserRepository.ListarRegistrados()`: `id` y `nombre` de
  `usuarios` con `contrasena_hash IS NOT NULL`, orden por `nombre_normalizado`.
- **Usuario registrado por nombre** — al configurar el equipo se busca por
  `nombre_normalizado` (lectura ya existente) y se exige que tenga contraseña.

## Reglas (dónde se aplican)

| Regla | Capa | Error → HTTP |
|---|---|---|
| Hay sesión vigente | `internal/http` (`conSesion`) | `ErrNoAutenticado` → 401 |
| Quien actúa integra el proyecto | servicio (`exigirIntegrante`) | `ErrNoEsIntegrante` → 403 «No sos integrante de este proyecto» |
| Quien actúa es Scrum Master del proyecto | servicio (`exigirScrumMaster`) | `ErrNoEsScrumMaster` → 403 «Solo el Scrum Master puede realizar esta acción» |
| El integrante nuevo es un usuario registrado | servicio `AsignarIntegrante` | `ErrUsuarioNoEncontrado` → 404 |
| Todo integrante nuevo es Product Builder | servicio `AsignarIntegrante` | `rol=scrum_master` → `ErrScrumMasterYaExiste` → 409 |
| El creador es el Scrum Master | servicio `CrearProyecto` (ya existe `GuardarConScrumMaster`) | — |

## Cambios de texto en `MODELO-DE-DATOS.md` (sin esquema)

- §4.2: `contrasena_hash` `NULL` = usuario previo; desde HU-21 ninguna ruta crea usuarios sin
  contraseña.
- §8: fila «HU-21 Permisos por rol | lee `integrantes`, `usuarios`, `sesiones` | — | sin
  migración».
